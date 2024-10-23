export type AmplifyDependentResourcesAttributes = {
  "api": {
    "FoodHandler": {
      "ApiId": "string",
      "ApiName": "string",
      "RootUrl": "string"
    },
    "NutrientHandler": {
      "ApiId": "string",
      "ApiName": "string",
      "RootUrl": "string"
    }
  },
  "function": {
    "FoodProcess": {
      "Arn": "string",
      "LambdaExecutionRole": "string",
      "LambdaExecutionRoleArn": "string",
      "Name": "string",
      "Region": "string"
    },
    "NutrientFetch": {
      "Arn": "string",
      "LambdaExecutionRole": "string",
      "LambdaExecutionRoleArn": "string",
      "Name": "string",
      "Region": "string"
    }
  }
}