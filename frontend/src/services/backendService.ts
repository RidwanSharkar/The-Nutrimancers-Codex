// Bioessence/frontend/src/services/backendService.ts
import axios from 'axios';

interface ProcessFoodResponse {
  ingredients: string[];
  nutrients: { [key: string]: number };
  missingNutrients: string[];
  suggestions: string[];
}

export const processFood = async (foodDescription: string): Promise<ProcessFoodResponse> => {
  try {
    const response = await axios.post<ProcessFoodResponse>('http://localhost:5000/api/process-food', {
      foodDescription,
    });

    return response.data;
  } catch (error: unknown) {
    if (axios.isAxiosError(error)) {
      // If it's an Axios error, check for response data
      throw new Error(error.response?.data?.error || 'An error occurred while processing the food.');
    } else {
      // Fallback for non-Axios errors
      throw new Error('An unexpected error occurred.');
    }
  }
};
