package main

import (
	"log"
	"net/http"
	"os"
	"x-tract/data"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	conn := os.Getenv("DB_CONNECTION")
	config := data.PostgresConfig{
		URL:          conn,
		MaxOpenConns: 25,
		MaxIdleConns: 10,
	}

	db, err := data.NewPostgres(config)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := data.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
