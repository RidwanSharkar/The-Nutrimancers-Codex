// backend/machinist/percentageCalculations.go
package machinist

var nutrientRDA = map[string]float64{
	// Ions
	"Potassium":  4700, // mg
	"Sodium":     2300, // mg
	"Calcium":    1000, // mg
	"Phosphorus": 700,  // mg
	"Magnesium":  400,  // mg
	"Iron":       10,   // mg
	"Zinc":       10,   // mg
	"Manganese":  2.3,  // mg
	"Copper":     0.9,  // mg
	"Selenium":   0.4,  // µg

	// Essential Amino-Acids
	"Histidine":     10000, // mg
	"Isoleucine":    19000, // mg
	"Leucine":       39000, // mg
	"Lysine":        30000, // mg
	"Methionine":    14000, // mg
	"Phenylalanine": 25000, // mg
	"Threonine":     15000, // mg
	"Tryptophan":    5000,  // mg
	"Valine":        24000, // mg

	// Essential Omega Fatty Acids
	"Alpha-Linolenic Acid": 1.2,  // g (Plant Omega-3)
	"Linoleic Acid":        1.0,  // g (Omega-6)
	"EPA":                  5000, // mg (Omega-3 fish oil)
	"DHA":                  3750, // mg (Omega-3 fish oil)

	// Vitamins
	"Vitamin A":   0.9,  // mg RAE
	"Vitamin B1":  1.2,  // mg
	"Vitamin B2":  1.3,  // mg
	"Vitamin B3":  16,   // mg
	"Vitamin B5":  5,    // mg
	"Vitamin B6":  1.5,  // mg
	"Vitamin B9":  0.4,  // µg
	"Vitamin B12": 0.06, // µg
	"Vitamin C":   90,   // mg
	"Vitamin D":   2.0,  // µg
	"Vitamin E":   15,   // mg
	"Vitamin K":   0.18, // µg

	"Choline": 550, // mg
}

var nutrientUnits = map[string]string{
	"Potassium":  "mg",
	"Sodium":     "mg",
	"Calcium":    "mg",
	"Phosphorus": "mg",
	"Magnesium":  "mg",
	"Iron":       "mg",
	"Zinc":       "mg",
	"Manganese":  "mg",
	"Copper":     "mg",
	"Selenium":   "µg",

	"Histidine":     "g",
	"Isoleucine":    "g",
	"Leucine":       "g",
	"Lysine":        "g",
	"Methionine":    "g",
	"Phenylalanine": "g",
	"Threonine":     "g",
	"Tryptophan":    "g",
	"Valine":        "g",

	"Alpha-Linolenic Acid": "mg",
	"Linoleic Acid":        "mg",
	"EPA":                  "g",
	"DHA":                  "g",

	"Vitamin A":   "µg",
	"Vitamin B1":  "mg",
	"Vitamin B2":  "mg",
	"Vitamin B3":  "mg",
	"Vitamin B5":  "mg",
	"Vitamin B6":  "mg",
	"Vitamin B9":  "µg",
	"Vitamin B12": "µg",
	"Vitamin C":   "mg",
	"Vitamin D":   "µg",
	"Vitamin E":   "mg",
	"Vitamin K":   "µg",

	"Choline": "mg",
}

func adjustUnits(amount float64, unit string) float64 {
	switch unit {
	case "mg":
		return amount
	case "µg":
		return amount / 1000.0
	case "g":
		return amount * 1000.0
	case "IU":
		return convertIUtoMg(amount)
	default:
		return amount
	}
}

func convertIUtoMg(amount float64) float64 {
	micrograms := amount * 0.025
	milligrams := micrograms / 1000.0
	return milligrams
}

// Calculate percentage of RDA
func CalculateNutrientPercentages(nutrientData map[string]map[string]float64) map[string]map[string]float64 {
	percentagesPerIngredient := make(map[string]map[string]float64)
	for ingredient, nutrients := range nutrientData {
		percentages := make(map[string]float64)
		for nutrient, amount := range nutrients {
			rda, rdaExists := nutrientRDA[nutrient]
			unit, unitExists := nutrientUnits[nutrient]
			if rdaExists && unitExists {
				// match units
				adjustedAmount := adjustUnits(amount, unit)
				percentage := (adjustedAmount / rda) * 100
				percentages[nutrient] = percentage
			} else {
				percentages[nutrient] = 0
			}
		}
		percentagesPerIngredient[ingredient] = percentages
	}
	return percentagesPerIngredient
}

func CalculateTotalNutrients(nutrientPercentages map[string]map[string]float64) map[string]float64 {
	totalNutrients := make(map[string]float64)

	for _, nutrients := range nutrientPercentages {
		for nutrient, percentage := range nutrients {
			totalNutrients[nutrient] += percentage
		}
	}

	// Cap @ 100%
	for nutrient, percentage := range totalNutrients {
		if percentage > 100 {
			totalNutrients[nutrient] = 100
		}
	}

	return totalNutrients
}

/*=================================================================================*/

// Determine Deficiencies
func DetermineLowAndMissingNutrients(totalNutrients map[string]float64) []string {
	var lowAndMissingNutrients []string

	// Threshold
	const lowThreshold = 3.5

	// Iterate over all
	for nutrient := range nutrientRDA {
		percentage, exists := totalNutrients[nutrient]
		if !exists || percentage <= lowThreshold {
			lowAndMissingNutrients = append(lowAndMissingNutrients, nutrient)
		}
	}

	return lowAndMissingNutrients
}
