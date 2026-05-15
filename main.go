package main

import (
	"log"
	"os"
	"time"
	"x-tract/auth"
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
	router.Use(ginCORS())
	jwtManager := newJWTManager()

	endpoints.Register(router, endpoints.Handlers{
		Auth:           controllers.NewAuthController(data.NewAuthStore(db), jwtManager),
		Admins:         controllers.NewAdminController(data.NewAdminStore(db)),
		Customers:      controllers.NewCustomerController(data.NewCustomerStore(db)),
		Groups:         controllers.NewGroupController(data.NewGroupStore(db)),
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

func newJWTManager() *auth.Manager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-insecure-secret-change-me"
	}

	return auth.NewManager(auth.Config{
		Secret:     secret,
		TTL:        24 * time.Hour,
		IssuerName: "x-tract",
	})
}

// In your main.go or where you setup routes
func ginCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
