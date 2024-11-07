// The-Nutrimancers-Codex/frontend/src/services/backendService.ts

import axios from 'axios';

interface ProcessFoodResponse {
  ingredients: string[];
  nutrients: { [ingredient: string]: { [nutrient: string]: number } };
  missingNutrients: string[];
  suggestions: string[];
}

export const processFood = async (foodDescription: string): Promise<ProcessFoodResponse> => {
  try {
    const response = await axios.post<ProcessFoodResponse>(`https://bioessence-env.elasticbeanstalk.com/process-food`, {
      foodDescription,
    });
    return response.data;
  } catch (error: unknown) {
    if (axios.isAxiosError(error)) {
      console.error('Full error response:', error.response);
      let detailedError = 'An error occurred while processing the food.';
      if (error.response?.data?.error) {
        detailedError = error.response.data.error;
      } else if (typeof error.response?.data === 'string') {
        detailedError = error.response.data;
      }
      throw new Error(detailedError);
    } else {
      throw new Error('An unexpected error occurred.');
    }
  }
};


