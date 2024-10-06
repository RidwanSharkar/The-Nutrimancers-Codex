// backend/handlers/foodHandler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RidwanSharkar/Bioessence/backend/models"
	"github.com/RidwanSharkar/Bioessence/backend/services"
	"github.com/RidwanSharkar/Bioessence/backend/utils"
)

var essentialNutrients = []string{
	"Potassium",
	"Chloride",
	"Sodium",
	"Calcium",
	"Phosphorus",
	"Magnesium",
	"Iron",
	"Zinc",
	"Manganese",
	"Copper",
	"Iodine",
	"Chromium",

	"Molybdenum",
	"Selenium",

	"Histidine",
	"Isoleucine",
	"Leucine",
	"Lysine",
	"Methionine",
	"Phenylalanine",
	"Threonine",
	"Tryptophan",
	"Valine",

	"Alpha-Linolenic Acid", // Omega-3
	"Linoleic Acid",        // Omega-6

	"Vitamin A",
	"Vitamin B1",
	"Vitamin B2",
	"Vitamin B3",
	"Vitamin B5",
	"Vitamin B6",
	"Vitamin B7",
	"Vitamin B9",
	"Vitamin B12",
	"Vitamin C",
	"Vitamin D",
	"Vitamin E",
	"Vitamin K",

	"choline",
}

var suggestionData = map[string]string{
	"Fiber":     "Consider eating whole grains or fruits.",
	"Protein":   "How about some lean meat or legumes?",
	"Calcium":   "Dairy products or leafy greens can help.",
	"Vitamin D": "Sunlight exposure or fortified foods can boost Vitamin D.",
	// Add more nutrient-suggestion mappings as needed
}

func ProcessFoodHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ProcessFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.FoodDescription) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Food description cannot be empty")
		return
	}

	// Extract ingredients using Gemini
	ingredients, err := services.ExtractIngredients(req.FoodDescription)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to extract ingredients")
		return
	}

	// Fetch nutrient data using Nutritionix
	nutrientsData, err := services.FetchNutrientData(ingredients)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch nutrient data")
		return
	}

	// Aggregate nutrients
	aggregatedNutrients := aggregateNutrients(nutrientsData)

	// Determine missing nutrients
	missing := determineMissingNutrients(aggregatedNutrients)

	// Generate suggestions
	suggestions := generateSuggestions(missing)

	// Prepare response
	response := models.ProcessFoodResponse{
		Ingredients:      ingredients,
		Nutrients:        aggregatedNutrients,
		MissingNutrients: missing,
		Suggestions:      suggestions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func aggregateNutrients(nutrientsData map[string]map[string]float64) map[string]float64 {
	aggregated := make(map[string]float64)
	for _, nutrientMap := range nutrientsData {
		for nutrient, value := range nutrientMap {
			aggregated[nutrient] += value
		}
	}
	return aggregated
}

func determineMissingNutrients(aggregated map[string]float64) []string {
	var missing []string
	for _, nutrient := range essentialNutrients {
		if _, exists := aggregated[nutrient]; !exists || aggregated[nutrient] == 0 {
			missing = append(missing, nutrient)
		}
	}
	return missing
}

func generateSuggestions(missing []string) []string {
	var suggestions []string
	for _, nutrient := range missing {
		if suggestion, exists := suggestionData[nutrient]; exists {
			suggestions = append(suggestions, suggestion)
		}
	}
	return suggestions
}
