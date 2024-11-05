// The-Nutrimancers-Codex/backend/handlers/processFoodHandler.go
package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RidwanSharkar/The-Nutrimancers-Codex/backend/machinist"
	"github.com/RidwanSharkar/The-Nutrimancers-Codex/backend/services"
	"github.com/RidwanSharkar/The-Nutrimancers-Codex/backend/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ProcessFoodRequest struct {
	FoodDescription string `json:"foodDescription"`
}

type ProcessFoodResponse struct {
	Ingredients      []string                      `json:"ingredients"`
	Nutrients        map[string]map[string]float64 `json:"nutrients"`
	MissingNutrients []string                      `json:"missingNutrients"`
	Suggestions      []string                      `json:"suggestions"`
}

func HandleProcessFood(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req ProcessFoodRequest
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

	// Load food data (consider loading once during initialization if possible)
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
	response := ProcessFoodResponse{
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
