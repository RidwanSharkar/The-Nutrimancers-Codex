// backend/models/model.go
package models

// ProcessFoodRequest represents the incoming request payload
type ProcessFoodRequest struct {
	FoodDescription string `json:"foodDescription"`
}

// ProcessFoodResponse represents the response payload
type ProcessFoodResponse struct {
	Ingredients      []string           `json:"ingredients"`
	Nutrients        map[string]float64 `json:"nutrients"`
	MissingNutrients []string           `json:"missingNutrients"`
	Suggestions      []string           `json:"suggestions"`
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
}
