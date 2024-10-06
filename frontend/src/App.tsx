// src/App.tsx

import React, { useState } from 'react';
import IngredientsPanel from './components/IngredientsPanel';
import NutrientsPanel from './components/NutrientsPanel';
import SuggestionPanel from './components/SuggestionPanel';
import { processFood } from './services/backendService';
import './App.css';

type ProcessFoodResponse = {
  ingredients: string[];
  nutrients: { [key: string]: number };
  missingNutrients: string[];
  suggestions: string[];
};

const App: React.FC = () => {
  const [food, setFood] = useState<string>(''); // User input
  const [ingredients, setIngredients] = useState<string[]>([]); // Processed ingredients
  const [nutrients, setNutrients] = useState<{ [key: string]: number }>({}); // Nutrient data
  const [missingNutrients, setMissingNutrients] = useState<string[]>([]); // Missing nutrients
  const [suggestions, setSuggestions] = useState<string[]>([]); // Suggestions
  const [loading, setLoading] = useState<boolean>(false); // Loading state
  const [error, setError] = useState<string | null>(null); // Error state

  const handleFoodSubmit = async () => {
    if (food.trim()) {
      setLoading(true);
      setError(null);
      try {
        const response: ProcessFoodResponse = await processFood(food.trim());
        console.log('Response:', response); // Log the response for debugging
        setIngredients(response.ingredients); // Update state with ingredients
        setNutrients(response.nutrients); // Update state with nutrients
        setMissingNutrients(response.missingNutrients); // Update state with missing nutrients
        setSuggestions(response.suggestions); // Update state with suggestions
      } catch (err: unknown) {
        if (err instanceof Error) {
          console.error('Error in handleFoodSubmit:', err.message); // Log error
          setError(err.message); // Set error state
        } else {
          setError('An unexpected error occurred.');
        }
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#F48668]">
      <div className="w-[1024px] max-w-full h-auto bg-[#FFC09F] p-8 rounded-lg shadow-lg overflow-hidden">
        <h1 className="text-4xl font-bold text-center mb-8">Bioessence</h1>

        {/* Input Box */}
        <div className="flex justify-center mb-8">
          <input
            type="text"
            className="w-1/3 p-3 rounded-md focus:outline-none bg-[#EF8354] text-white mr-5"
            placeholder="What have you eaten today?"
            value={food}
            onChange={(e) => setFood(e.target.value)}
            onKeyPress={(e) => {
              if (e.key === 'Enter') handleFoodSubmit();
            }}
          />
          <button
            onClick={handleFoodSubmit}
            disabled={loading || !food.trim()}
            className={`p-3 bg-green-500 text-white rounded-md ${
              loading || !food.trim() ? 'bg-green-300 cursor-not-allowed' : 'hover:bg-green-600'
            }`}
          >
            {loading ? 'Processing...' : 'Submit'}
          </button>
        </div>

        {/* Error State */}
        {error && <p className="text-center text-red-500 mb-4">{error}</p>}

        {/* Panels */}
        {!loading && !error && ingredients.length > 0 && (
          <div className="flex flex-col lg:flex-row justify-center gap-8 mx-auto w-full max-w-5xl">
            <IngredientsPanel ingredients={ingredients} />
            <NutrientsPanel nutrients={nutrients} />
            <SuggestionPanel missingNutrients={missingNutrients} suggestions={suggestions} />
          </div>
        )}
      </div>
    </div>
  );
};

export default App;
