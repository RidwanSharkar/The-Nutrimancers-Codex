// src/App.tsx

import React, { useState } from 'react';
import IngredientsPanel from './components/IngredientsPanel';
import NutrientsPanel from './components/NutrientsPanel';
import SuggestionPanel from './components/SuggestionPanel';
import { processFood } from './services/backendService';
import './App.css';

type ProcessFoodResponse = {
  ingredients: string[];
  nutrients: { [ingredient: string]: { [nutrient: string]: number } };
  missingNutrients: string[];
  suggestions: string[];
};

const App: React.FC = () => {
  const [food, setFood] = useState<string>('');
  const [ingredients, setIngredients] = useState<string[]>([]);
  const [nutrients, setNutrients] = useState<{ [ingredient: string]: { [key: string]: number } }>({});
  const [selectedIngredient, setSelectedIngredient] = useState<string>('');
  const [selectedNutrientData, setSelectedNutrientData] = useState<{ [key: string]: number }>({});
  const [missingNutrients, setMissingNutrients] = useState<string[]>([]);
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [highlightedNutrients, setHighlightedNutrients] = useState<string[]>([]);
  const [normalMealNutrients, setNormalMealNutrients] = useState<{ [key: string]: number }>({});
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

/*=================================================================================================*/

  const handleFoodSubmit = async () => {
    if (food.trim()) {
      setLoading(true);
      setError(null);
      try {
        const response: ProcessFoodResponse = await processFood(food.trim());
        console.log('Received ingredients:', response.ingredients);
        console.log('Received nutrients:', response.nutrients);
        console.log('Received missingNutrients:', response.missingNutrients);
        console.log('Received suggestions:', response.suggestions);
        setIngredients(response.ingredients || []);
        setNutrients(response.nutrients);
        setMissingNutrients(response.missingNutrients);
        setSuggestions(response.suggestions);
        setSelectedIngredient('Full Meal');
        setSelectedNutrientData({});
        setHighlightedNutrients([]);

        // Total Nutrients Full Meal
        const totalNutrients: { [key: string]: number } = {};
        response.ingredients.forEach((ing) => {
          const nutrientData = response.nutrients[ing] || {};
          for (const nutrient in nutrientData) {
            if (Object.prototype.hasOwnProperty.call(nutrientData, nutrient)) {
              totalNutrients[nutrient] = (totalNutrients[nutrient] || 0) + nutrientData[nutrient];
            }
          }
        });
        setNormalMealNutrients(totalNutrients);
      } catch (err: unknown) {
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError('An unexpected error occurred.');
        }
      } finally {
        setLoading(false);
      }
    }
  };

  const handleIngredientClick = (ingredient: string) => {
    console.log('Clicked ingredient:', ingredient);

    // Reset highlighted
    setHighlightedNutrients([]);

    if (ingredient === 'Full Meal') {
      setSelectedIngredient('Full Meal');
      setSelectedNutrientData(normalMealNutrients);
    } else {
      console.log('Nutrients for this ingredient:', nutrients[ingredient]);
      setSelectedIngredient(ingredient);
      setSelectedNutrientData(nutrients[ingredient] || {});
    }
  };

  const handleRecommendationClick = async (suggestion: string) => {
    try {
      const response = await fetch('http://localhost:5000/fetch-nutrient-data', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          foodDescription: suggestion,
          currentNutrients: normalMealNutrients,
        }),
      });
      const data = await response.json();

      // Update state with new + highlighted nutrients
      setSelectedIngredient('Full Meal');
      setSelectedNutrientData(data.nutrients);
      setHighlightedNutrients(data.changedNutrients);
    } catch (error) {
      console.error('Error fetching nutrient data for recommendation:', error);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#F48668]">
      <div className="w-[1024px] max-w-full h-auto bg-[#FFC09F] p-8 rounded-lg shadow-lg overflow-hidden">
        <h1 className="text-4xl font-bold text-center mb-8">The Nutrimancer's Codex Vol. I</h1>

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
            {loading ? 'Extracting...' : 'Extract Essence'}
          </button>
        </div>

        {error && <p className="text-center text-red-500 mb-4">{error}</p>}

        {!loading && !error && ingredients && ingredients.length > 0 && (
          <div className="flex flex-col lg:flex-row justify-center gap-8 mx-auto w-full max-w-5xl">
            <IngredientsPanel ingredients={ingredients} onIngredientClick={handleIngredientClick} />
            <NutrientsPanel
              ingredient={selectedIngredient}
              nutrients={selectedNutrientData}
              highlightedNutrients={highlightedNutrients}
            />
            <SuggestionPanel
              missingNutrients={missingNutrients}
              suggestions={suggestions}
              onRecommendationClick={handleRecommendationClick}
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default App;

