// backend/machinist/recommendTron.go
package machinist

import (
	"sort"

	"github.com/RidwanSharkar/Bioessence/backend/models"
)

type Recommendation struct {
	Description     string
	SimilarityScore float64
	Nutrients       map[string]float64
}

func RecommendFoods(foodItems []models.FoodItem, nutrientNames []string, deficiencies []string, topN int) []string {
	deficiencyVector := createDeficiencyVector(nutrientNames, deficiencies)

	var recommendations []Recommendation

	for _, food := range foodItems {
		foodVector := make([]float64, len(nutrientNames))
		for i, nutrientName := range nutrientNames {
			foodVector[i] = food.Nutrients[nutrientName]
		}
		similarity := CosineSimilarity(foodVector, deficiencyVector)
		if similarity > 0 {
			recommendation := Recommendation{
				Description:     food.Description,
				SimilarityScore: similarity,
				Nutrients:       make(map[string]float64),
			}
			for _, nutrient := range deficiencies {
				recommendation.Nutrients[nutrient] = food.Nutrients[nutrient]
			}
			recommendations = append(recommendations, recommendation)
		}
	}

	// Sort recommendations by similarity score in descending order
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].SimilarityScore > recommendations[j].SimilarityScore
	})

	// Get top N recommendations
	topRecommendations := recommendations
	if len(recommendations) > topN {
		topRecommendations = recommendations[:topN]
	}

	// Extract descriptions
	var suggestedFoods []string
	for _, rec := range topRecommendations {
		suggestedFoods = append(suggestedFoods, rec.Description)
	}

	return suggestedFoods
}

func createDeficiencyVector(nutrientNames []string, deficiencies []string) []float64 {
	deficiencyVector := make([]float64, len(nutrientNames))
	deficiencySet := make(map[string]bool)
	for _, nutrient := range deficiencies {
		deficiencySet[nutrient] = true
	}

	for i, nutrientName := range nutrientNames {
		if deficiencySet[nutrientName] {
			deficiencyVector[i] = 1.0
		} else {
			deficiencyVector[i] = 0.0
		}
	}
	return deficiencyVector
}
