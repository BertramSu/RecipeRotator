package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Username     string
	Password     string
	Address      string
	DatabaseName string
}

var cfg = Config{
	Username:     "admin",
	Password:     "root",
	Address:      "localhost:5432",
	DatabaseName: "commentdb",
}

var connStr = fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName)

type Comment struct {
	ID       string `json:"id"`
	RecipeId string `json:"recipeId"`
	Body     string `json:"body"`
}

func main() {
	router := gin.Default()
	router.GET("/comments", getAllComments)
	router.GET("/comment/:id", getCommentById)
	router.POST("/comment", postComment)

	router.Run("localhost:8080")
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
