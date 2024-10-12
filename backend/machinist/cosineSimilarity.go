// backend/machinist/cosineSimilarity.go
package machinist

import "math"

func CosineSimilarity(vecA, vecB []float64) float64 {
	var dotProduct, normA, normB float64
	for i := 0; i < len(vecA); i++ {
		dotProduct += vecA[i] * vecB[i]
		normA += vecA[i] * vecA[i]
		normB += vecB[i] * vecB[i]
	}
	if normA == 0 || normB == 0 {
		return 0.0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
