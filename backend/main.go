// backend/main.go

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RidwanSharkar/Bioessence/backend/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Initialize router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/api/process-food", handlers.ProcessFoodHandler).Methods("POST")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
