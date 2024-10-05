// services/geminiService.go

package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type GeminiRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type GeminiChoice struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Choices []GeminiChoice `json:"choices"`
}

func ExtractIngredients(foodDescription string) ([]string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY not set")
	}

	prompt := `Extract the list of ingredients from the following food description:
    
"` + foodDescription + `"

List the ingredients as bullet points without additional text.`

	reqBody := GeminiRequest{
		Model:       "gemini-1.5-flash",
		Prompt:      prompt,
		MaxTokens:   150,
		Temperature: 0.5,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.*gemini-endpoint.com/v1/completions", bytes.NewBuffer(jsonData)) // Replace with your actual endpoint
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("Gemini API error: " + string(bodyBytes))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Choices) == 0 {
		return nil, errors.New("No choices returned from Gemini")
	}

	text := geminiResp.Choices[0].Text
	ingredients := parseIngredients(text)
	return ingredients, nil
}

func parseIngredients(text string) []string {
	var ingredients []string
	lines := bytes.Split([]byte(text), []byte("\n"))
	for _, line := range lines {
		cleaned := bytes.TrimPrefix(line, []byte("- "))
		cleaned = bytes.TrimSpace(cleaned)
		if len(cleaned) > 0 {
			ingredients = append(ingredients, string(cleaned))
		}
	}
	return ingredients
}
