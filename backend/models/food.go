// backend/models/food.go
package models

type FoodItem struct {
	FdcID       string
	Description string
	Nutrients   map[string]float64
}
