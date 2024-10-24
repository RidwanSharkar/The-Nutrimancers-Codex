// backend/handlers/processFoodHandler/processFoodHandler.go
package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RidwanSharkar/Bioessence-Codex/machinist"
	"github.com/RidwanSharkar/Bioessence-Codex/models"
	"github.com/RidwanSharkar/Bioessence-Codex/services"
	"github.com/RidwanSharkar/Bioessence-Codex/utils"
	"github.com/aws/aws-lambda-go/events"
)

func HandleProcessFood(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.FoodRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusBadRequest, "Invalid request payload")
	}

	if req.FoodDescription == "" {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusBadRequest, "Food description is required")
	}

	// Extract ingredients using Gemini LLM
	ingredients, err := services.ExtractIngredients(req.FoodDescription)
	if err != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error extracting ingredients: "+err.Error())
	}

	cleanedIngredients := services.CleanIngredientList(ingredients)

	// Load food data
	foodItems, nutrientNames, loadErr := machinist.LoadFoodData("machinist/dataset.csv")
	if loadErr != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error loading food data: "+loadErr.Error())
	}

	// Fetch nutrient data for each ingredient using Nutritionix API
	nutrientData, fetchErr := services.FetchNutrientDataForEachIngredient(cleanedIngredients)
	if fetchErr != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error fetching nutrient data: "+fetchErr.Error())
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
