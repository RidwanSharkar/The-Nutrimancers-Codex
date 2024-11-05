// The-Nutrimancers-Codex/amplify/backend/machinist/dataLoader.go
package machinist

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/RidwanSharkar/The-Nutrimancers-Codex/amplify/backend/models"
)

func LoadFoodData(filePath string) ([]models.FoodItem, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(records) < 2 {
		return nil, nil, fmt.Errorf("CSV file must have at least a header and one data row")
	}

	// Header processing
	header := records[0]
	nutrientNames := header[2:] // Nutrient names start from column 3

	var foodItems []models.FoodItem

	for _, record := range records[1:] {
		if len(record) < 2 {
			continue // Skip if not enough fields
		}

		foodItem := models.FoodItem{
			FdcID:       record[0],
			Description: record[1],
			Nutrients:   make(map[string]float64),
		}

		for i, nutrientName := range nutrientNames {
			if i+2 >= len(record) {
				break
			}
			valueStr := record[i+2]
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				value = 0.0 // missing/invalid
			}
			foodItem.Nutrients[nutrientName] = value
		}
		foodItems = append(foodItems, foodItem)
	}

	return foodItems, nutrientNames, nil
}
