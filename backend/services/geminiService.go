// services/geminiService.go
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/RidwanSharkar/Bioessence/backend/utils"
)

/*=================================================================================================*/

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

type GeminiChoice struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Choices []GeminiChoice `json:"choices"`
}

/*=================================================================================================*/

// Primary Prompt: Accepts user food description dynamically and sends to Gemini API
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
	reqBody := GeminiRequest{
		Contents: []Content{{Parts: []Part{{
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
	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		utils.LogError(err, "ExtractIngredients: Decode")
		return nil, err
	}

	// Optional choices later?
	if len(geminiResp.Choices) == 0 {
		errMsg := "no choices returned from gemini"
		utils.LogError(errors.New(errMsg), "no gemini choices")
		return nil, errors.New(errMsg)
	}
	text := geminiResp.Choices[0].Text
	ingredients := parseIngredients(text)
	return ingredients, nil
}

/*=================================================================================================*/

// Parse and clean ingredients from Gemini response
func parseIngredients(text string) []string {
	var ingredients []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		cleaned := strings.TrimSpace(line)
		cleaned = strings.Trim(cleaned, "-•,.") // Remove bullet points, punctuation
		cleaned = strings.ToLower(cleaned)      // Convert to lowercase
		if len(cleaned) > 0 {
			ingredients = append(ingredients, cleaned) // Append valid ingredients
		}
	}
	return ingredients
}
