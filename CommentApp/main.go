package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	"encoding/json"
	"errors"
	"sync"

	"github.com/IBM/sarama"
)

/*
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Notification struct {
	From    User   `json:"from"`
	To      User   `json:"to"`
	Message string `json:"message"`
}
*/

type MyId struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const (
	ConsumerGroup      = "notifications-group"
	ConsumerTopic      = "delete-recipe-topic"
	ConsumerPort       = ":8085"
	KafkaServerAddress = "localhost:9092"
)

type Config struct {
	Username     string
	Password     string
	Address      string
	DatabaseName string
}

var cfg = Config{
	Username:     "admin",
	Password:     "passwordHere",
	Address:      "localhost:5434",
	DatabaseName: "CommentDB",
}

var connStr = fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName)

type Comment struct {
	ID       string `json:"id"`
	RecipeId string `json:"recipeId"`
	Body     string `json:"body"`
}

func main() {
	store := &NotificationStore{
		data: make(UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go setupConsumerGroup(ctx, store)
	defer cancel()

	//gin.SetMode(gin.ReleaseMode)
	fmt.Print("Made it to default")
	router := gin.Default()
	router.GET("/comments", getAllComments)
	router.GET("/comment/:id", getCommentById)
	router.POST("/comment", postComment)
	router.PUT("/comment", updateComment)
	router.DELETE("/comment/:id", deleteCommentById)

	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		handleNotifications(ctx, store)
	})

	router.Run("localhost:8085")

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+
		"started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	if err := router.Run(ConsumerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

func getAllComments(c *gin.Context) {
	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Query example
	query := `
	SELECT
		comment_id AS Id,
		recipe_id AS RecipeId,
		body AS Body
	FROM
		comment
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed getAllComments: %v", err)
	}
	defer rows.Close()

	comments, err := pgx.CollectRows(rows, pgx.RowToStructByName[Comment])
	if err != nil {
		log.Fatalf("Conversion failed getAllComments: %v", err)
	}
	c.IndentedJSON(http.StatusOK, comments)
}

func getCommentById(c *gin.Context) {
	id := c.Param("id")

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	query := `
	SELECT
		comment_id AS Id,
		recipe_id AS RecipeId,
		body AS Body
	FROM
		comment
	WHERE
		comment_id = @id
	`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := conn.Query(context.Background(), query, args)
	if err != nil {
		log.Fatalf("Query failed getCommentById): %v", err)
		return
	}
	defer rows.Close()

	comment, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Comment])
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "comment not found"})
	}
	c.IndentedJSON(http.StatusOK, comment)
}

func postComment(c *gin.Context) {
	var newComment Comment

	if err := c.BindJSON(&newComment); err != nil {
		log.Fatal("Cannot bind newComment")
		return
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return
	}
	defer conn.Close(context.Background())

	query := `
	INSERT INTO comment(recipe_id, body)
	VALUES (@recipeId, @body)
	RETURNING recipe_id
	`

	args := pgx.NamedArgs{
		"recipeId": newComment.RecipeId,
		"body":     newComment.Body,
	}

	var id int
	qErr := conn.QueryRow(context.Background(), query, args).Scan(&id)
	if qErr != nil {
		log.Fatalf("Query failed insertAlbum: %v", qErr)
		return
	}

	c.IndentedJSON(http.StatusCreated, newComment)
}

func deleteCommentById(c *gin.Context) {
	id := c.Param("id")

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	query := `
	DELETE
	FROM comment
	WHERE
		comment_id = @id
	RETURNING comment_id
	`
	args := pgx.NamedArgs{
		"id": id,
	}

	commandTag, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		log.Fatal("Query failed delete comment.")
		return
	}

	if commandTag.RowsAffected() != 1 {
		log.Fatalf("No comment found with ID %d", id)
		return
	}

	c.IndentedJSON(http.StatusOK, id)
}

func updateComment(c *gin.Context) {
	var newComment Comment

	if err := c.BindJSON(&newComment); err != nil {
		log.Fatal("Cannot bind newComment")
		return
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return
	}
	defer conn.Close(context.Background())

	query := `
	UPDATE comment
	SET body= @body
	WHERE comment_id = @id
	RETURNING comment_id
	`

	args := pgx.NamedArgs{
		"id":   newComment.ID,
		"body": newComment.Body,
	}

	var id int
	qErr := conn.QueryRow(context.Background(), query, args).Scan(&id)
	if qErr != nil {
		log.Fatalf("Query failed updateComment: %v", qErr)
		return
	}

	c.IndentedJSON(http.StatusCreated, newComment)
}

// ============== HELPER FUNCTIONS ==============
var ErrNoMessagesFound = errors.New("no messages found")

func getUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessagesFound
	}
	return userID, nil
}

// ====== NOTIFICATION STORAGE ======
type UserNotifications map[string][]MyId

type NotificationStore struct {
	data UserNotifications
	mu   sync.RWMutex
}

func (ns *NotificationStore) Add(userID string,
	notification MyId) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.data[userID] = append(ns.data[userID], notification)
}

func (ns *NotificationStore) Get(userID string) []MyId {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return ns.data[userID]
}

// ============== KAFKA RELATED FUNCTIONS ==============
type Consumer struct {
	store *NotificationStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Print("consume claim here")
		userID := string(msg.Key)
		var notification MyId
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}
		consumer.store.Add(userID, notification)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {
	fmt.Print("made it to initialize consumer group")
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{KafkaServerAddress}, ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	return consumerGroup, nil
}

func setupConsumerGroup(ctx context.Context, store *NotificationStore) {
	consumerGroup, err := initializeConsumerGroup()
	if err != nil {
		log.Printf("initialization error: %v", err)
	}
	defer consumerGroup.Close()

	consumer := &Consumer{
		store: store,
	}
	fmt.Print("made it to right before consume \n")
	for {
		fmt.Print("Weee idk \n")
		err = consumerGroup.Consume(ctx, []string{ConsumerTopic}, consumer)
		if err != nil {
			fmt.Print("error for err != nill consume")
			log.Printf("error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			fmt.Print("mawp")
			return
		}
	}
}

func handleNotifications(ctx *gin.Context, store *NotificationStore) {
	userID, err := getUserIDFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	notes := store.Get(userID)
	if len(notes) == 0 {
		ctx.JSON(http.StatusOK,
			gin.H{
				"message":       "No notifications found for user",
				"notifications": []MyId{},
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notes})
}
