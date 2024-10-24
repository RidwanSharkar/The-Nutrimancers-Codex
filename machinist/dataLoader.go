// backend/machinist/dataLoader.go
package machinist

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/RidwanSharkar/Bioessence/backend/models"
)

var (
	foodItems     []models.FoodItem
	nutrientNames []string
	once          sync.Once
	loadErr       error
)

func LoadFoodData(filePath string) ([]models.FoodItem, []string, error) {
	once.Do(func() {
		file, err := os.Open(filePath)
		if err != nil {
			loadErr = err
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1

		records, err := reader.ReadAll()
		if err != nil {
			loadErr = err
			return
		}

		if len(records) < 2 {
			loadErr = fmt.Errorf("CSV file must have at least a header and one data row")
			return
		}

		// Header processing
		header := records[0]
		nutrientNames = header[2:] // Nutrient names start from column 3

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
	})

	return foodItems, nutrientNames, loadErr
}
