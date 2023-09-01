package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Thread struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
}

func main() {
	router := gin.Default()
	router.GET("/threads", getThreads)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

func getThreads(c *gin.Context) {
	connectDB()

	topics, err := getThreadsDb()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.IndentedJSON(http.StatusOK, topics)
}

func postAlbums(c *gin.Context) {
	/*var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)*/
}

func connectDB() {
	if db != nil {
		return
	}

	// Capture connection properties.
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "definitive_forum",
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to DB...")
}

func getThreadsDb() ([]Thread, error) {
	var threads []Thread

	rows, err := db.Query("select id, title, created_at from thread where deleted_at is null")
	if err != nil {
		return nil, fmt.Errorf("Failed to query threads: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.CreatedAt); err != nil {
			return nil, fmt.Errorf("Failed to map threads: %v", err)
		}
		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Failed to map threads: %v", err)
	}

	return threads, nil
}
