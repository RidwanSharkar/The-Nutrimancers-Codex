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

	"github.com/RidwanSharkar/Bioessence/backend/models"
	"github.com/RidwanSharkar/Bioessence/backend/services"
	"github.com/RidwanSharkar/Bioessence/backend/utils"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

/*==================================================================================*/

// Redundant combine - MAP ALL 4 together{nutrient, unit, RDA, nutrtionixAPI}*
var essentialNutrients = []string{
	"Potassium",
	//"Chloride",
	"Sodium",
	"Calcium",
	"Phosphorus",
	"Magnesium",
	"Iron",
	"Zinc",
	"Manganese",
	"Copper",
	//"Iodine",
	//"Chromium",
	//"Molybdenum",
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
	"Alpha-Linolenic Acid",
	"Linoleic Acid",
	"Vitamin A",
	"Vitamin B1",
	"Vitamin B2",
	"Vitamin B3",
	"Vitamin B5",
	"Vitamin B6",
	//"Vitamin B7",
	"Vitamin B9",
	"Vitamin B12",
	"Vitamin C",
	"Vitamin D",
	"Vitamin E",
	"Vitamin K",
	"Choline",
}

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Allow requests from React app
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	// Define the HTTP endpoint
	http.HandleFunc("/process-food", processFoodHandler)

	// Wrap the multiplexer with CORS middleware
	handler := c.Handler(http.DefaultServeMux)

	// Start the server
	fmt.Println("Server is running on :5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}

/*=================================================================================*/

var nutrientRDA = map[string]float64{
	// Ions
	"Potassium": 4700, // mg
	// "Chloride":   2300,
	"Sodium":     2300, // mg
	"Calcium":    1000, // mg
	"Phosphorus": 700,  // mg
	"Magnesium":  400,  // mg
	"Iron":       10,   // mg
	"Zinc":       10,   // mg
	"Manganese":  2.3,  // mg
	"Copper":     0.9,  // mg
	//"Iodine":     150,  // µg
	// "Chromium":   35,   // µg
	// "Molybdenum": 45,   // µg
	"Selenium": 55, // µg

	// Essential Amino-Acids
	"Histidine":     10000, // mg
	"Isoleucine":    19000, // mg
	"Leucine":       42000, // mg
	"Lysine":        38000, // mg
	"Methionine":    15000, // mg
	"Phenylalanine": 25000, // mg
	"Threonine":     20000, // mg
	"Tryptophan":    5000,  // mg
	"Valine":        24000, // mg

	// Essential Omega Fatty Acids
	"Alpha-Linolenic Acid": 1300, // mg
	"Linoleic Acid":        1400, // mg

	// Vitamins
	"Vitamin A":  0.9, // mg
	"Vitamin B1": 1.2, // mg
	"Vitamin B2": 1.3, // mg
	"Vitamin B3": 16,  // mg
	"Vitamin B5": 5,   // mg
	"Vitamin B6": 1.5, // mg
	// "Vitamin B7":  30,  // µg
	"Vitamin B9":  400, // µg
	"Vitamin B12": 2.4, // µg
	"Vitamin C":   90,  // mg
	"Vitamin D":   15,  // µg
	"Vitamin E":   15,  // mg
	"Vitamin K":   120, // µg

	// Semi-Essential
	"Choline": 550, // mg
	//"Boron" :
}

// Conserve - UNIT CONVERSIONS ====================================================================================
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
	//"Iodine":     "µg",
	//"Chromium":   "µg",
	//"Molybdenum": "µg",
	"Selenium": "µg",

	"Histidine":     "g",
	"Isoleucine":    "g",
	"Leucine":       "g",
	"Lysine":        "g",
	"Methionine":    "g",
	"Phenylalanine": "g",
	"Threonine":     "g",
	"Tryptophan":    "g",
	"Valine":        "g",

	"Alpha-Linolenic Acid": "mg",
	"Linoleic Acid":        "mg",

	"Vitamin A":  "µg",
	"Vitamin B1": "mg",
	"Vitamin B2": "mg",
	"Vitamin B3": "mg",
	"Vitamin B5": "mg",
	"Vitamin B6": "mg",
	//"Vitamin B7":  "µg",
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
				// match RDA units
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

	// Fetch nutrient data for each ingredient
	nutrientData, err := services.FetchNutrientDataForEachIngredient(cleanedIngredients)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching nutrient data: "+err.Error())
		return
	}

	// Aggregate Data
	aggregatedNutrients := aggregateNutrients(nutrientData)
	missingNutrients := determineMissingNutrients(aggregatedNutrients)

	// Generate Suggestions
	suggestions := generateSuggestions(missingNutrients)

	// RDA Colors
	nutrientPercentages := calculateNutrientPercentages(nutrientData)

	// Prepare the response using models.ProcessFoodResponse
	response := models.ProcessFoodResponse{
		Ingredients:      cleanedIngredients,
		Nutrients:        nutrientPercentages,
		MissingNutrients: missingNutrients,
		Suggestions:      suggestions,
	}

	// Send Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/*=================================================================================*/

// Function to extract Gemini output
func extractIngredientsFromGemini(apiKey, prompt string) (string, error) {
	// Create the request payload
	reqBody := models.GeminiRequest{
		Contents: []models.Content{
			{
				Parts: []models.Part{
					{Text: prompt},
				},
			},
		},
	}

	// Convert the request payload to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Send the request to the Gemini API
	endpoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + apiKey
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Decode the JSON response
	var geminiResp models.GeminiResponse
	err = json.Unmarshal(bodyBytes, &geminiResp)
	if err != nil {
		return "", err
	}

	// Extract the ingredients from the response
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

// Aggregate Nutrition Data
func aggregateNutrients(nutrientData map[string]map[string]float64) map[string]map[string]float64 {
	//Check
	return nutrientData
}

// Determine Missing Nutrients
func determineMissingNutrients(nutrientData map[string]map[string]float64) []string {
	var missing []string
	nutrientSet := make(map[string]bool)
	for _, nutrients := range nutrientData {
		for nutrient := range nutrients {
			nutrientSet[nutrient] = true
		}
	}
	for _, nutrient := range essentialNutrients {
		if !nutrientSet[nutrient] {
			missing = append(missing, nutrient)
		}
	}
	return missing
}

// Generate Suggestions Based on Missing Nutrients
func generateSuggestions(missing []string) []string {
	suggestionsMap := map[string]string{
		"Vitamin D": "Include more fatty fish or fortified dairy products.",
		"Calcium":   "Consider adding more leafy greens or dairy products.",
		// Add more mappings as needed
	}

	var suggestions []string
	for _, nutrient := range missing {
		if suggestion, exists := suggestionsMap[nutrient]; exists {
			suggestions = append(suggestions, suggestion)
		} else {
			suggestions = append(suggestions, fmt.Sprintf("Consider adding sources rich in %s.", nutrient))
		}
	}
	return suggestions
}
