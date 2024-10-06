package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/RidwanSharkar/Bioessence/backend/services"
	"github.com/joho/godotenv"
)

// Request payload structure for Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content CandidateContent `json:"content"`
}

type CandidateContent struct {
	Parts []Part `json:"parts"`
}

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("API_KEY environment variable is not set.")
		return
	}

	// Taking user input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a food description: ")
	scanner.Scan()
	foodDescription := scanner.Text()

	// Prompt
	promptText := fmt.Sprintf("Extract the essential ingredients from the following food description: '%s'. For complex foods like pizza, include the base components (e.g., dough, cheese, tomato sauce). Exclude spices and minor ingredients.", foodDescription)

	// Extract ingredients using Gemini LLM
	ingredients, err := extractIngredientsFromGemini(apiKey, promptText)
	if err != nil {
		fmt.Println("Error extracting ingredients from Gemini:", err)
		return
	}

	// Clean
	cleanedIngredients := cleanIngredientList(ingredients)
	fmt.Println("Cleaned Ingredients:", cleanedIngredients)

	// Fetch nutrient data for each cleaned ingredient - Nutritionix
	nutrientData, err := services.FetchNutrientDataForEachIngredient(cleanedIngredients)
	if err != nil {
		fmt.Println("Error fetching nutrient data from Nutritionix:", err)
		return
	}

	// Aggregate Data
	aggregatedNutrients := aggregateNutrients(nutrientData)
	missingNutrients := determineMissingNutrients(aggregatedNutrients)

	// Generate Suggestions
	suggestions := generateSuggestions(missingNutrients)

	fmt.Println("Aggregated Nutrients:", aggregatedNutrients)
	fmt.Println("Missing Nutrients:", missingNutrients)
	fmt.Println("Suggestions:", suggestions)
}

// Additional functions (extractIngredientsFromGemini, cleanIngredientList, aggregateNutrients, determineMissingNutrients, generateSuggestions) remain the same as previously defined.

// Function to extract ingredients from Gemini
func extractIngredientsFromGemini(apiKey, prompt string) (string, error) {
	// Create the request payload
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
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
	var geminiResp GeminiResponse
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

/*=================================================================================================*/

// Aggregate Nutrition Data
func aggregateNutrients(nutrientData map[string]map[string]float64) map[string]float64 {
	return make(map[string]float64)
}

// Determine Missing
func determineMissingNutrients(aggregated map[string]float64) []string {
	return []string{}
}

// Generate Suggestions
func generateSuggestions(missing []string) []string {
	return []string{}
}
