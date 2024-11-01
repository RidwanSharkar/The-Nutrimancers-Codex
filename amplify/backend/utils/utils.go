// The-Nutrimancers-Codex/amplify/backend/utils/utils.go:
package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// ErrorResponse struct for error messages
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError creates a standardized error response
func RespondWithError(resp events.APIGatewayProxyResponse, statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	errorResponse := ErrorResponse{
		Error: message,
	}
	body, _ := json.Marshal(errorResponse)
	resp.StatusCode = statusCode
	resp.Body = string(body)
	resp.Headers = map[string]string{"Content-Type": "application/json"}
	return resp, nil
}

// LogError logs errors (can be expanded to use structured logging)
func LogError(err error, context string) {
	if err != nil {

	}
}
