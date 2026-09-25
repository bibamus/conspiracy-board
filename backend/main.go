package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// Initialize database
	db, err := NewDatabase("conspiracy-board.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize API
	api := NewAPI(db)

	// Setup router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", api.Health)

	// Graph endpoints
	router.GET("/api/graph", api.GetGraph)

	// Connection Type endpoints
	router.POST("/api/connection-types", api.CreateConnectionType)
	router.GET("/api/connection-types", api.GetConnectionTypes)
	router.PUT("/api/connection-types/:id", api.UpdateConnectionType)
	router.DELETE("/api/connection-types/:id", api.DeleteConnectionType)

	// Person endpoints
	router.POST("/api/people", api.CreatePerson)
	router.GET("/api/people", api.GetPeople)
	router.GET("/api/people/:id", api.GetPerson)
	router.DELETE("/api/people/:id", api.DeletePerson)

	// Connection endpoints
	router.POST("/api/connections", api.CreateConnection)
	router.GET("/api/connections", api.GetConnections)
	router.GET("/api/people/:id/connections", api.GetPersonConnections)
	router.DELETE("/api/connections/:id", api.DeleteConnection)

	// Start server
	port := 8000
	fmt.Printf("Starting backend server on http://localhost:%d\n", port)
	if err := router.Run(fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
