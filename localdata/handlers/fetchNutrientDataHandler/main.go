// backend/handlers/fetchNutrientDataHandler/main.go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleFetchNutrientData)
}
