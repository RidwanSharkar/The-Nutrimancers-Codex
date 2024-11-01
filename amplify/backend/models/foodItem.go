// The-Nutrimancers-Codex/amplify/backend/models/foodItem.go
package models

type FoodItem struct {
	FdcID       string
	Description string
	Nutrients   map[string]float64
}
