// src/services/backendService.ts

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
  } catch (error: any) {
    throw new Error(error.response?.data?.error || 'An error occurred while processing the food.');
  }
};
