// The-Nutrimancers-Codex/amplify/backend/services/nutritionixService.go:
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

/*=================================================================================================*/
// Mapping of WIKI essential nutrients to Nutritionix attr_ids
var nutrientMapping = map[string]int{
	"Potassium":            306, // mg
	"Sodium":               307,
	"Calcium":              301,
	"Phosphorus":           305,
	"Magnesium":            304,
	"Iron":                 303,
	"Zinc":                 309,
	"Manganese":            315,
	"Copper":               312,
	"Selenium":             317, // µg
	"Histidine":            512, // g
	"Isoleucine":           503,
	"Leucine":              504,
	"Lysine":               505,
	"Methionine":           506,
	"Phenylalanine":        508,
	"Threonine":            502, // g
	"Tryptophan":           501, // g
	"Valine":               510,
	"Alpha-Linolenic Acid": 851, // mg Omega-3
	"Linoleic Acid":        675, // mg Omega-6
	"EPA":                  629, // g Omega-3
	"DHA":                  621, // g Omega-3
	"Vitamin A":            320, // µg RAE
	"Vitamin B1":           404, // mg
	"Vitamin B2":           405, // mg
	"Vitamin B3":           406, // mg
	"Vitamin B5":           410, // mg
	"Vitamin B6":           415, // mg
	"Vitamin B9":           417, // µg
	"Vitamin B12":          418, // µg
	"Vitamin C":            401, // mg
	"Vitamin D":            324, // µg
	"Vitamin E":            323, // mg
	"Vitamin K":            430, // µg
	"Choline":              421, // mg
}

/*=================================================================================================*/
// Fetch nutrient data for each individual ingredient from Nutritionix API
func FetchNutrientDataForEachIngredient(ingredients []string) (map[string]map[string]float64, error) {
	nutrientsPerIngredient := make(map[string]map[string]float64)

	for _, ingredient := range ingredients {
		nutrientData, err := FetchNutrientData([]string{ingredient}) // one ingredient at a time
		if err != nil {
			return nil, fmt.Errorf("error fetching nutrient data for %s: %v", ingredient, err)
		}
		nutrientsPerIngredient[ingredient] = nutrientData[ingredient]
	}
	return nutrientsPerIngredient, nil
}

func FetchNutrientData(ingredients []string) (map[string]map[string]float64, error) {
	appID := os.Getenv("NUTRITIONIX_APP_ID") // Load from environment
	appKey := os.Getenv("NUTRITIONIX_APP_KEY")
	if appID == "" || appKey == "" {
		return nil, errors.New("missing Nutritionix API credentials in environment variables")
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
			return nil, fmt.Errorf("nutritionix API error: %s", string(bodyBytes))
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
					nutrients[nutrient] = fn.Value
					break
				}
			}
		}
		nutrientsPerIngredient[ingredient] = nutrients
	}
	return nutrientsPerIngredient, nil
}
