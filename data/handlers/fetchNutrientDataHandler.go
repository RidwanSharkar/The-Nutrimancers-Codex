// The-Nutrimancers-Codex/backend/handlers/fetchNutrientDataHandler.go
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

type FetchNutrientDataRequest struct {
	FoodDescription  string             `json:"foodDescription"`
	CurrentNutrients map[string]float64 `json:"currentNutrients"`
}

type FetchNutrientDataResponse struct {
	Nutrients        map[string]float64 `json:"nutrients"`
	ChangedNutrients []string           `json:"changedNutrients"`
}

func HandleFetchNutrientData(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req FetchNutrientDataRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusBadRequest, "Invalid request payload")
	}

	if req.FoodDescription == "" || req.CurrentNutrients == nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusBadRequest, "Food description and current nutrients are required")
	}

	// Fetch nutrient data for the suggested food using Nutritionix API
	nutrientData, err := services.FetchNutrientData([]string{req.FoodDescription})
	if err != nil {
		return utils.RespondWithError(events.APIGatewayProxyResponse{}, http.StatusInternalServerError, "Error fetching nutrient data: "+err.Error())
	}

	// Calculate RDA percentages
	nutrientPercentages := machinist.CalculateNutrientPercentages(nutrientData)

	// Combine current nutrients with new nutrients
	newTotalNutrients := make(map[string]float64)
	for nutrient, amount := range req.CurrentNutrients {
		newTotalNutrients[nutrient] = amount
	}
	for nutrient, amount := range nutrientPercentages[req.FoodDescription] {
		newTotalNutrients[nutrient] += amount
		if newTotalNutrients[nutrient] > 100 {
			newTotalNutrients[nutrient] = 100
		}
	}

	// Determine which nutrients have changed
	changedNutrients := []string{}
	for nutrient := range nutrientPercentages[req.FoodDescription] {
		changedNutrients = append(changedNutrients, nutrient)
	}

	// Prepare the response
	response := FetchNutrientDataResponse{
		Nutrients:        newTotalNutrients,
		ChangedNutrients: changedNutrients,
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
	lambda.Start(HandleFetchNutrientData)
}
