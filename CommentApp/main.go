package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	"encoding/json"

	"github.com/IBM/sarama"
)

type DeletedRecipe struct {
	ID int `json:"id"`
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
	ctx, cancel := context.WithCancel(context.Background())
	go setupConsumerGroup(ctx)
	defer cancel()

	router := gin.Default()
	router.GET("/comments", getAllComments)
	router.GET("/comment/:id", getCommentById)
	router.POST("/comment", postComment)
	router.PUT("/comment", updateComment)
	router.DELETE("/comment/:id", deleteCommentById)
	router.GET("/recipe/:id/comments", getCommentsByRecipeId)

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

func getCommentsByRecipeId(c *gin.Context) {
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
		recipe_id = @id
	`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := conn.Query(context.Background(), query, args)
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

// Kafka functions.
type Consumer struct {
}

// These three are require by the interface.
func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var deletedRecipe DeletedRecipe
		err := json.Unmarshal(msg.Value, &deletedRecipe)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}
		print("Consumed deleted recipe id: ")
		fmt.Println(deletedRecipe.ID)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup([]string{KafkaServerAddress}, ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	return consumerGroup, nil
}

func setupConsumerGroup(ctx context.Context) {
	consumerGroup, err := initializeConsumerGroup()
	if err != nil {
		log.Printf("initialization error: %v", err)
	}
	defer consumerGroup.Close()

	consumer := &Consumer{}
	for {
		err = consumerGroup.Consume(ctx, []string{ConsumerTopic}, consumer)
		if err != nil {
			log.Printf("error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}
