export type AmplifyDependentResourcesAttributes = {
  "api": {
    "FetchNutrient": {
      "ApiId": "string",
      "ApiName": "string",
      "RootUrl": "string"
    },
    "ProcessFood": {
      "ApiId": "string",
      "ApiName": "string",
      "RootUrl": "string"
    }
  },
  "function": {
    "FetchNutrient": {
      "Arn": "string",
      "LambdaExecutionRole": "string",
      "LambdaExecutionRoleArn": "string",
      "Name": "string",
      "Region": "string"
    },
    "ProcessFood": {
      "Arn": "string",
      "LambdaExecutionRole": "string",
      "LambdaExecutionRoleArn": "string",
      "Name": "string",
      "Region": "string"
    }
  }
}