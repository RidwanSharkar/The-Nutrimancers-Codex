// The-Nutrimancers-Codex/amplify/backend/function/ProcessFood/src/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/machinist"
	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/models"
	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/services"
	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Global variable to store the Gemini API Key
var geminiAPIKey string

// InitializeSecrets retrieves secrets during the initialization phase
func InitializeSecrets() {
	secretName := "nutrimancer/google-service-account-key"
	region := "us-east-1"

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // Defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatalf("Unable to retrieve secret %s: %v", secretName, err)
	}

	// Decrypts secret using the associated KMS key.
	geminiAPIKey = *result.SecretString
}

type contextKey string

func HandleProcessFood(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("HandleProcessFood invoked")
	log.Printf("Request Body: %s\n", request.Body)

	// Ensure the API key has been initialized
	if geminiAPIKey == "" {
		log.Println("Gemini API key is not initialized")
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "API configuration error")
	}

	// Set the API key in the context if needed by downstream services
	ctx = context.WithValue(ctx, "geminiAPIKeyContextKey", geminiAPIKey)

	var req models.FoodRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Printf("Error unmarshalling request body: %v\n", err)
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusBadRequest, "Invalid request payload")
	}

	if req.FoodDescription == "" {
		log.Println("Food description is empty")
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusBadRequest, "Food description is required")
	}

	// Extract ingredients using Gemini LLM
	ingredients, err := services.ExtractIngredients(req.FoodDescription)
	if err != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error extracting ingredients: "+err.Error())
	}

	cleanedIngredients := services.CleanIngredientList(ingredients)

	// Load food data
	foodItems, nutrientNames, err := machinist.LoadFoodData()
	if err != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error loading food data: "+err.Error())
	}

	// Fetch nutrient data for each ingredient using Nutritionix API
	nutrientData, err := services.FetchNutrientDataForEachIngredient(cleanedIngredients)
	if err != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error fetching nutrient data: "+err.Error())
	}

	// Calculate RDA percentages
	nutrientPercentages := machinist.CalculateNutrientPercentages(nutrientData)

	// Calculate total nutrients
	totalNutrients := machinist.CalculateTotalNutrients(nutrientPercentages)

	// Determine Deficiencies
	lowAndMissingNutrients := machinist.DetermineLowAndMissingNutrients(totalNutrients)

	// Generate Recommendations
	topRecommendations := machinist.RecommendFoods(foodItems, nutrientNames, lowAndMissingNutrients, 5)

	// Prepare the response
	response := models.ProcessFoodResponse{
		Ingredients:      cleanedIngredients,
		Nutrients:        nutrientPercentages,
		MissingNutrients: lowAndMissingNutrients,
		Suggestions:      topRecommendations,
	}

	respBody, err := json.Marshal(response)
	if err != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error forming response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	// Initialize secrets before starting the Lambda handler
	InitializeSecrets()

	// Start the Lambda handler
	lambda.Start(HandleProcessFood)
}
