export type AmplifyDependentResourcesAttributes = {
  "api": {
    "fetchNutrient": {
      "ApiId": "string",
      "ApiName": "string",
      "RootUrl": "string"
    },
    "processFood": {
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