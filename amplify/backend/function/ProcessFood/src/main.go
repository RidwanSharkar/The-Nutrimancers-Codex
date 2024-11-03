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

// getGeminiAPIKey retrieves the API key from AWS Secrets Manager
func getGeminiAPIKey(ctx context.Context) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return "", err
	}

	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("nutrimancer/google-service-account-key"),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}

func HandleProcessFood(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("HandleProcessFood invoked")
	log.Printf("Request Body: %s\n", request.Body)

	// Get API key at the start of the function
	apiKey, err := getGeminiAPIKey(ctx)
	if err != nil {
		log.Printf("Error retrieving Gemini API key: %v\n", err)
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error with API configuration")
	}

	// Set the API key in the context
	ctx = context.WithValue(ctx, "GEMINI_API_KEY", apiKey)

	var req models.FoodRequest
	if err = json.Unmarshal([]byte(request.Body), &req); err != nil {
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
	lambda.Start(HandleProcessFood)
}
