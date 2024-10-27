// backend/services/geminiService.go
package services

import (
	"The-Nutrimancers-Codex/amplify/backend/models"
	"The-Nutrimancers-Codex/amplify/backend/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/*=================================================================================================*/

// Primary Prompt: Accepts food description in any format and sends to Gemini API
func ExtractIngredients(foodDescription string) ([]string, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		err := errors.New("API_KEY not set")
		utils.LogError(err, "ExtractIngredients")
		return nil, err
	}

	// PROMPT
	promptText := fmt.Sprintf("Extract the list of ingredients from the following food description: '%s'. List the main ingredients as bullet points with no descriptions.", foodDescription)

	// Prep Request Body
	reqBody := models.GeminiRequest{
		Contents: []models.Content{{Parts: []models.Part{{
			Text: promptText,
		}}}},
	}

	// Convert request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		utils.LogError(err, "ExtractIngredients: Marshal")
		return nil, err
	}

	// Gemini API Request
	endpoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + apiKey
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.LogError(err, "ExtractIngredients: NewRequest")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogError(err, "ExtractIngredients: DoRequest")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("Gemini API error: %s", string(bodyBytes))
		utils.LogError(errors.New(errMsg), "ExtractIngredients: API Error")
		return nil, errors.New(errMsg)
	}

	// Parsing Response
	var geminiResp models.GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		utils.LogError(err, "ExtractIngredients: Decode")
		return nil, err
	}

	// Optional choices later?
	if len(geminiResp.Candidates) == 0 {
		errMsg := "no choices returned from gemini"
		utils.LogError(errors.New(errMsg), "no gemini choices")
		return nil, errors.New(errMsg)
	}
	text := geminiResp.Candidates[0].Content.Parts[0].Text
	ingredients := parseIngredients(text)
	return ingredients, nil
}

/*=================================================================================================*/

// Parse & Clean Ingredients
func parseIngredients(text string) []string {
	var ingredients []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		cleaned := strings.TrimSpace(line)
		cleaned = strings.Trim(cleaned, "-•,.")
		cleaned = strings.ToLower(cleaned)
		if len(cleaned) > 0 {
			ingredients = append(ingredients, cleaned)
		}
	}
	return ingredients
}

// Clean Ingredient List
func CleanIngredientList(ingredients []string) []string {
	var cleaned []string
	for _, ingredient := range ingredients {
		ingredient = strings.TrimSpace(ingredient)
		if ingredient != "" {
			cleaned = append(cleaned, ingredient)
		}
	}
	return cleaned
}
