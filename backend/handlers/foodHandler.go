// handlers/foodHandler.go

package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"your-project/backend/services"
)

type ProcessFoodRequest struct {
	FoodDescription string `json:"foodDescription"`
}

type ProcessFoodResponse struct {
	Ingredients      []string           `json:"ingredients"`
	Nutrients        map[string]float64 `json:"nutrients"`
	MissingNutrients []string           `json:"missingNutrients"`
	Suggestions      []string           `json:"suggestions"`
}

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
	"Linoleic Acid",        // Omega-6 (CHECK NUTRIONIX USAGE)

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

	"choline"
}

var suggestionData = map[string]string{
	"Fiber":   "Consider eating whole grains or fruits.",
	"Protein": "How about some lean meat or legumes?",
	"Calcium": "Dairy products or leafy greens can help.",
	// TENTATIVE
}

func ProcessFoodHandler(w http.ResponseWriter, r *http.Request) {
	var req ProcessFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.FoodDescription) == "" {
		http.Error(w, "Food description cannot be empty", http.StatusBadRequest)
		return
	}

	// Extract ingredients using Gemini LLM API
	ingredients, err := services.ExtractIngredients(req.FoodDescription)
	if err != nil {
		http.Error(w, "Failed to extract ingredients: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch nutrient data using Nutritionix API
	nutrientsData, err := services.FetchNutrientData(ingredients)
	if err != nil {
		http.Error(w, "Failed to fetch nutrient data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Aggregate Nutrients
	aggregatedNutrients := aggregateNutrients(nutrientsData)

	// Determine Missing
	missing := determineMissingNutrients(aggregatedNutrients)

	// Generate Suggestions
	suggestions := generateSuggestions(missing)

	// Response Format
	response := ProcessFoodResponse{
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
