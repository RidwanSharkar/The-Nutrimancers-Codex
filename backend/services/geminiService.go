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

	"github.com/RidwanSharkar/Bioessence/backend/utils"
)

// GeminiRequest defines the structure of the request payload for Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// GeminiChoice represents each choice in the Gemini API response
type GeminiChoice struct {
	Text string `json:"text"`
}

// GeminiResponse represents the overall response structure from Gemini API
type GeminiResponse struct {
	Choices []GeminiChoice `json:"choices"`
}

// ExtractIngredients sends a request to the Gemini API to extract ingredients from a food description
func ExtractIngredients(foodDescription string) ([]string, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		err := errors.New("API_KEY not set")
		utils.LogError(err, "ExtractIngredients")
		return nil, err
	}

	// Construct the prompt as per Gemini API's expected format
	promptText := fmt.Sprintf(`Extract the list of ingredients from the following food description:

"%s"

List the ingredients as bullet points without additional text.`, foodDescription)

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: promptText,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		utils.LogError(err, "ExtractIngredients: Marshal")
		return nil, err
	}

	// Gemini API Endpoint
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

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		utils.LogError(err, "ExtractIngredients: Decode")
		return nil, err
	}

	if len(geminiResp.Choices) == 0 {
		errMsg := "No choices returned from Gemini"
		utils.LogError(errors.New(errMsg), "ExtractIngredients: NoChoices")
		return nil, errors.New(errMsg)
	}

	text := geminiResp.Choices[0].Text
	ingredients := parseIngredients(text)
	return ingredients, nil
}

// parseIngredients processes the Gemini API response text to extract ingredients
func parseIngredients(text string) []string {
	var ingredients []string
	lines := bytes.Split([]byte(text), []byte("\n"))
	for _, line := range lines {
		cleaned := bytes.TrimSpace(line)
		// Remove common bullet point prefixes
		cleaned = bytes.TrimPrefix(cleaned, []byte("- "))
		cleaned = bytes.TrimPrefix(cleaned, []byte("* "))
		cleaned = bytes.TrimPrefix(cleaned, []byte("• "))
		cleaned = bytes.TrimSpace(cleaned)
		if len(cleaned) > 0 {
			ingredients = append(ingredients, string(cleaned))
		}
	}
	return ingredients
}
