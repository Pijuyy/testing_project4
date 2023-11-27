package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Pijuyy/testing_project4/config"
	"github.com/Pijuyy/testing_project4/controllers"
	"github.com/Pijuyy/testing_project4/routes"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database connection
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Create a new router
	router := mux.NewRouter()

	// Register routes
	routes.RegisterRoutes(router, db)

	// Initialize admin user
	controllers.InitializeAdminUser(db)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// try a lot of
	// try 2.0 a lot of

	// Start the server
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
