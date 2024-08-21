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
	//router.GET("/albums/:id", getAlbumByID)
	//router.POST("/albums", postAlbums)

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
	query := "SELECT comment_id AS Id, recipe_id AS RecipeId, body AS Body FROM comment"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed getAllComments: %v", err)
	}
	defer rows.Close()

	comments, err := pgx.CollectRows(rows, pgx.RowToStructByName[Comment])
	c.IndentedJSON(http.StatusOK, comments)
}
