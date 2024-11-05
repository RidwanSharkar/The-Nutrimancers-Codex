// The-Nutrimancers-Codex/backend/main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/machinist"
	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/models"
	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/services"
	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/utils"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

/*==================================================================================*/

var (
	foodItems     []models.FoodItem
	nutrientNames []string
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// Machinist
	dataFilePath := "machinist/dataset.csv" // Adjust the path as needed
	foodItems, nutrientNames, err = machinist.LoadFoodData(dataFilePath)
	if err != nil {
		log.Fatal("Error loading food data:", err)
	}
	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Allow requests from React app
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	// HTTP endpoint
	http.HandleFunc("/process-food", processFoodHandler)
	http.HandleFunc("/fetch-nutrient-data", fetchNutrientDataHandler)

	handler := c.Handler(http.DefaultServeMux)

	// Start server
	fmt.Println("Server is running on :5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}

/*=================================================================================*/

func determineLowAndMissingNutrients(totalNutrients map[string]float64) []string {
	var lowAndMissingNutrients []string

	// Threshold
	const lowThreshold = 3.5

	// Iterate over all
	for nutrient := range nutrientRDA {
		percentage, exists := totalNutrients[nutrient]
		if !exists || percentage <= lowThreshold {
			lowAndMissingNutrients = append(lowAndMissingNutrients, nutrient)
		}
	}

	return lowAndMissingNutrients
}

// Redundant combine - MAP ALL 4 together{nutrient, unit, RDA, nutrtionixAPI}*
var nutrientRDA = map[string]float64{
	// Ions
	"Potassium":  4700, // mg
	"Sodium":     2300, // mg
	"Calcium":    1000, // mg
	"Phosphorus": 700,  // mg
	"Magnesium":  400,  // mg
	"Iron":       10,   // mg
	"Zinc":       10,   // mg
	"Manganese":  2.3,  // mg
	"Copper":     0.9,  // mg
	"Selenium":   0.4,  // Âµg

	// Essential Amino-Acids
	"Histidine":     10000, // mg
	"Isoleucine":    19000, // mg
	"Leucine":       39000, // mg
	"Lysine":        30000, // mg
	"Methionine":    14000, // mg
	"Phenylalanine": 25000, // mg
	"Threonine":     15000, // mg
	"Tryptophan":    5000,  // mg
	"Valine":        24000, // mg

	// Essential Omega Fatty Acids
	"Alpha-Linolenic Acid": 1.2,  // g (Plant Omega-3)
	"Linoleic Acid":        1.0,  // g (Omega- 6)
	"EPA":                  5000, // (Omega-3 fish oil)
	"DHA":                  3750, // (Omega-3 fish oil)

	// Vitamins
	"Vitamin A":   0.9,  // mg       	Âµg
	"Vitamin B1":  1.2,  // mg
	"Vitamin B2":  1.3,  // mg
	"Vitamin B3":  16,   // mg
	"Vitamin B5":  5,    // mg
	"Vitamin B6":  1.5,  // mg
	"Vitamin B9":  0.4,  // µg	Âµg check api documentation mgiht be outdated
	"Vitamin B12": 0.06, // µg	Âµg
	"Vitamin C":   90,   // mg
	"Vitamin D":   2.0,  // µg	IU
	"Vitamin E":   15,   // mg
	"Vitamin K":   0.18, // µg	Âµg mg callibrating to standard serving

	"Choline": 550, // mg
}

// Conserve - UNIT CONVERSIONS =================================================================
var nutrientUnits = map[string]string{
	"Potassium": "mg",
	//"Chloride":   "mg",
	"Sodium":     "mg",
	"Calcium":    "mg",
	"Phosphorus": "mg",
	"Magnesium":  "mg",
	"Iron":       "mg",
	"Zinc":       "mg",
	"Manganese":  "mg",
	"Copper":     "mg",
	"Selenium":   "µg",

	"Histidine":     "g",
	"Isoleucine":    "g",
	"Leucine":       "g",
	"Lysine":        "g",
	"Methionine":    "g",
	"Phenylalanine": "g",
	"Threonine":     "g",
	"Tryptophan":    "g",
	"Valine":        "g",

	"Alpha-Linolenic Acid": "mg", // Omega-3
	"Linoleic Acid":        "mg", // Omega-6
	"EPA":                  "g",  // Omega-3
	"DHA":                  "g",  // Omega-3

	"Vitamin A":   "µg",
	"Vitamin B1":  "mg",
	"Vitamin B2":  "mg",
	"Vitamin B3":  "mg",
	"Vitamin B5":  "mg",
	"Vitamin B6":  "mg",
	"Vitamin B9":  "µg",
	"Vitamin B12": "µg",
	"Vitamin C":   "mg",
	"Vitamin D":   "µg",
	"Vitamin E":   "mg",
	"Vitamin K":   "µg",

	"Choline": "mg",
}

func adjustUnits(amount float64, unit string) float64 {
	switch unit {
	case "mg":
		return amount
	case "µg":
		return amount / 1000.0
	case "g":
		return amount * 1000.0
	case "IU":
		return convertIUtoMg(amount)
	default:
		return amount
	}
}

func convertIUtoMg(amount float64) float64 {
	micrograms := amount * 0.025
	milligrams := micrograms / 1000.0
	return milligrams
}

// ================================================================================================================

// Calculate percentage of RDA
func calculateNutrientPercentages(nutrientData map[string]map[string]float64) map[string]map[string]float64 {
	percentagesPerIngredient := make(map[string]map[string]float64)
	for ingredient, nutrients := range nutrientData {
		percentages := make(map[string]float64)
		for nutrient, amount := range nutrients {
			rda, rdaExists := nutrientRDA[nutrient]
			unit, unitExists := nutrientUnits[nutrient]
			if rdaExists && unitExists {
				// match units
				adjustedAmount := adjustUnits(amount, unit)
				percentage := (adjustedAmount / rda) * 100
				percentages[nutrient] = percentage
			} else {
				percentages[nutrient] = 0
			}
		}
		percentagesPerIngredient[ingredient] = percentages
	}
	return percentagesPerIngredient
}

func calculateTotalNutrients(nutrientPercentages map[string]map[string]float64) map[string]float64 {
	totalNutrients := make(map[string]float64)

	for _, nutrients := range nutrientPercentages {
		for nutrient, percentage := range nutrients {
			totalNutrients[nutrient] += percentage
		}
	}

	// Cap @ 100%
	for nutrient, percentage := range totalNutrients {
		if percentage > 100 {
			totalNutrients[nutrient] = 100
		}
	}

	return totalNutrients
}

/*=================================================================================*/

func fetchNutrientDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		FoodDescription  string             `json:"foodDescription"`
		CurrentNutrients map[string]float64 `json:"currentNutrients"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding request: "+err.Error())
		return
	}

	// Fetch nutrient data for the suggested food using Nutritionix API
	nutrientData, err := services.FetchNutrientData([]string{req.FoodDescription})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching nutrient data: "+err.Error())
		return
	}

	// Calculate RDA percentages
	nutrientPercentages := calculateNutrientPercentages(nutrientData)

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

	// Response
	response := struct {
		Nutrients        map[string]float64 `json:"nutrients"`
		ChangedNutrients []string           `json:"changedNutrients"`
	}{
		Nutrients:        newTotalNutrients,
		ChangedNutrients: changedNutrients,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/*=================================================================================*/

func processFoodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.FoodRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error extracting ingredients: "+err.Error())
		return
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		utils.RespondWithError(w, http.StatusInternalServerError, "API_KEY not set")
		return
	}

	promptText := fmt.Sprintf("Extract the essential ingredients from the following food description: '%s'. For complex foods like pizza, include the base components (e.g., dough, cheese, tomato sauce). Exclude spices and minor ingredients.", req.FoodDescription)

	// Extract ingredients using Gemini LLM
	ingredients, err := extractIngredientsFromGemini(apiKey, promptText)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error extracting ingredients: "+err.Error())
		return
	}

	cleanedIngredients := cleanIngredientList(ingredients)

	// Fetch nutrient data for each ingredient using Nutritionix API
	nutrientData, err := services.FetchNutrientDataForEachIngredient(cleanedIngredients)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching nutrient data: "+err.Error())
		return
	}

	// Calculate RDA percentages
	nutrientPercentages := calculateNutrientPercentages(nutrientData)

	// Calculate total nutrients
	totalNutrients := calculateTotalNutrients(nutrientPercentages)

	// Determine Deficiencies
	lowAndMissingNutrients := determineLowAndMissingNutrients(totalNutrients)

	// Generate Recommendations
	topRecommendations := machinist.RecommendFoods(foodItems, nutrientNames, lowAndMissingNutrients, 5)

	// Prepare the response
	response := models.ProcessFoodResponse{
		Ingredients:      cleanedIngredients,
		Nutrients:        nutrientPercentages,
		MissingNutrients: lowAndMissingNutrients,
		Suggestions:      topRecommendations,
	}

	// Send Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/*=================================================================================*/

// Function to extract Gemini output
func extractIngredientsFromGemini(apiKey, prompt string) (string, error) {
	// Request payload
	reqBody := models.GeminiRequest{
		Contents: []models.Content{
			{
				Parts: []models.Part{
					{Text: prompt},
				},
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Send request to Gemini API
	endpoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + apiKey
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var geminiResp models.GeminiResponse
	err = json.Unmarshal(bodyBytes, &geminiResp)
	if err != nil {
		return "", err
	}

	// Extract ingredients
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no ingredients returned from Gemini")
}

// Extract True Ingredients
func cleanIngredientList(ingredients string) []string {
	lines := strings.Split(ingredients, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "*", "")
		line = strings.TrimPrefix(line, "• ")
		if len(line) > 0 && !strings.Contains(line, ":") && len(line) < 50 {
			cleaned = append(cleaned, line)
		}
	}
	return cleaned
}

/*=================================================================================*/
