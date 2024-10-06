// services/nutritionixService.go

package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type NutritionixRequest struct {
	Query string `json:"query"`
}

type FullNutrient struct {
	AttrID int     `json:"attr_id"`
	Value  float64 `json:"value"`
}

type NutritionixFood struct {
	FoodName      string         `json:"food_name"`
	ServingQty    float64        `json:"serving_qty"`
	ServingUnit   string         `json:"serving_unit"`
	FullNutrients []FullNutrient `json:"full_nutrients"`
}

type NutritionixResponse struct {
	Foods []NutritionixFood `json:"foods"`
}

// Mapping of essential nutrients to Nutritionix [attr_ids]
var nutrientMapping = map[string]int{
	"Potassium":            306,
	"Chloride":             307, // verify
	"Sodium":               307,
	"Calcium":              301, // verify
	"Phosphorus":           305,
	"Magnesium":            304,
	"Iron":                 303, // verify
	"Zinc":                 309,
	"Manganese":            315,
	"Copper":               312,
	"Iodine":               318,
	"Chromium":             317,
	"Molybdenum":           319,
	"Selenium":             317, // Verify
	"Histidine":            512,
	"Isoleucine":           503,
	"Leucine":              504,
	"Lysine":               505,
	"Methionine":           506,
	"Phenylalanine":        508,
	"Threonine":            507,
	"Tryptophan":           509,
	"Valine":               510,
	"Alpha-Linolenic Acid": 645,
	"Linoleic Acid":        646,
	"Vitamin A":            318,
	"Vitamin B1":           415,
	"Vitamin B2":           418,
	"Vitamin B3":           417,
	"Vitamin B5":           415,
	"Vitamin B6":           415,
	"Vitamin B7":           415,
	"Vitamin B9":           417,
	"Vitamin B12":          418,
	"Vitamin C":            401,
	"Vitamin D":            324,
	"Vitamin E":            323,
	"Vitamin K":            428,
	"Choline":              430,
}

func FetchNutrientData(ingredients []string) (map[string]map[string]float64, error) {
	appID := os.Getenv("NUTRITIONIX_APP_ID")
	appKey := os.Getenv("NUTRITIONIX_APP_KEY")
	if appID == "" || appKey == "" {
		return nil, errors.New("Nutritionix API credentials not set")
	}

	nutrientsPerIngredient := make(map[string]map[string]float64)

	for _, ingredient := range ingredients {
		reqBody := NutritionixRequest{
			Query: ingredient,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", "https://trackapi.nutritionix.com/v2/natural/nutrients", bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-app-id", appID)
		req.Header.Set("x-app-key", appKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			return nil, errors.New("Nutritionix API error: " + string(bodyBytes))
		}

		var nutritionixResp NutritionixResponse
		if err := json.NewDecoder(resp.Body).Decode(&nutritionixResp); err != nil {
			return nil, err
		}

		if len(nutritionixResp.Foods) == 0 {
			nutrientsPerIngredient[ingredient] = map[string]float64{}
			continue
		}

		food := nutritionixResp.Foods[0]
		nutrients := make(map[string]float64)

		for nutrient, attrID := range nutrientMapping {
			for _, fn := range food.FullNutrients {
				if fn.AttrID == attrID {
					nutrients[nutrient] += fn.Value
				}
			}
		}

		nutrientsPerIngredient[ingredient] = nutrients
	}

	return nutrientsPerIngredient, nil
}
