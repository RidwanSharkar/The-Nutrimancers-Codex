// backend/models/model.go
package models

// Represents the response payload
type ProcessFoodResponse struct {
	Ingredients      []string                      `json:"ingredients"`
	Nutrients        map[string]map[string]float64 `json:"nutrients"`
	MissingNutrients []string                      `json:"missingNutrients"`
	Suggestions      []string                      `json:"suggestions"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

/*==================================================================================*/

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

/*==================================================================================*/

// outgoing responses
type FoodResponse struct {
	Ingredients      []string                      `json:"ingredients"`
	Nutrients        map[string]map[string]float64 `json:"nutrients"`
	MissingNutrients []string                      `json:"missingNutrients"`
	Suggestions      []string                      `json:"suggestions"`
}
type FoodRequest struct {
	FoodDescription string `json:"foodDescription"`
}
