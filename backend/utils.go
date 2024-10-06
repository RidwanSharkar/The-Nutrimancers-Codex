// utils.go
package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithError sends an error response in JSON format
func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"error": message}
	jsonResp, _ := json.Marshal(response)
	w.Write(jsonResp)
}

// LogError logs the error with a standard message
func LogError(err error, context string) {
	if err != nil {
		log.Printf("Error in %s: %v\n", context, err)
	}
}
