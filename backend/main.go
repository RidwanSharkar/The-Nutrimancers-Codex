// main.go

package main

import (
    "log"
    "net/http"
    "os"

    "your-project/backend/handlers"

    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
)

func main() {
    // Load environment variables
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
