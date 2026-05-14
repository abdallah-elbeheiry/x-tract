package main

import (
	"log"
	"os"
	"x-tract/controllers"
	"x-tract/data"
	"x-tract/endpoints"

	"github.com/gin-gonic/gin"
)

// main wires the database, HTTP handlers, and Gin router together.
func main() {
	db := mustDatabase()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("database close error: %v", err)
		}
	}()

	router := gin.Default()

	endpoints.Register(router, endpoints.Handlers{
		Admins:         controllers.NewAdminController(data.NewAdminStore(db)),
		Customers:      controllers.NewCustomerController(data.NewCustomerStore(db)),
		Salesmen:       controllers.NewSalesmanController(data.NewSalesmanStore(db)),
		GuestEmployees: controllers.NewGuestEmployeeController(data.NewGuestEmployeeStore(db)),
	})

	if err := router.Run(serverAddress()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// mustDatabase opens the database and applies migrations before the server starts.
func mustDatabase() *data.Database {
	config := data.PostgresConfig{
		URL:          os.Getenv("DB_CONNECTION"),
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

	return db
}

// serverAddress returns the listening address for the HTTP server.
func serverAddress() string {
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		return addr
	}
	return ":8080"
}
