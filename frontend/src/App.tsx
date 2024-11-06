// The-Nutrimancers-Codex/frontend/src/App.tsx

import React, { useState, useMemo } from 'react';
import IngredientsPanel from './grimoire/IngredientsPanel';
import SuggestionPanel from './grimoire/SuggestionPanel';
import OrbsPanel from './grimoire/OrbsPanel';
import { processFood } from './services/backendService';
import './App.css';

const nutrientCategoryList = {
  Minerals: [
    'Potassium',
    'Sodium',
    'Calcium',
    'Phosphorus',
    'Magnesium',
    'Iron',
    'Zinc',
    'Manganese',
    'Copper',
    'Selenium',
  ],
  Vitamins: [
    'Vitamin A',
    'Vitamin B1',
    'Vitamin B2',
    'Vitamin B3',
    'Vitamin B5',
    'Vitamin B6',
    'Vitamin B9',
    'Vitamin B12',
    'Vitamin C',
    'Vitamin D',
    'Vitamin E',
    'Vitamin K',
  ],
  'Amino Acids': [
    'Histidine',
    'Isoleucine',
    'Leucine',
    'Lysine',
    'Methionine',
    'Phenylalanine',
    'Threonine',
    'Tryptophan',
    'Valine',
  ],
  'Fatty Acids & Choline': [
    'Linoleic Acid',
    'Alpha-Linolenic Acid',
    'EPA',
    'DHA',
    'Choline',
  ],
  Total: [],
} as const;

// Orb Colors
const nutrientCategoryColors: { [key in keyof typeof nutrientCategoryList]: string } = {
  Minerals: '#2ecc71', //green
  Vitamins: '#3498db', // blue
  'Amino Acids': '#e67e22', // Orange
  'Fatty Acids & Choline': '#B47EB3', //purp
  Total: '#e74c3c', // Red
};

type NutrientCategory = keyof typeof nutrientCategoryList;

// Categorization
const categorizeNutrients = (
  nutrients: { [key: string]: number }
): { [category in NutrientCategory]: { total: number; satisfied: number; color: string } } => {
  const categorized: {
    [category in NutrientCategory]: { total: number; satisfied: number; color: string };
  } = {
    Minerals: {
      total: nutrientCategoryList.Minerals.length,
      satisfied: 0,
      color: nutrientCategoryColors.Minerals,
    },
    Vitamins: {
      total: nutrientCategoryList.Vitamins.length,
      satisfied: 0,
      color: nutrientCategoryColors.Vitamins,
    },
    'Amino Acids': {
      total: nutrientCategoryList['Amino Acids'].length,
      satisfied: 0,
      color: nutrientCategoryColors['Amino Acids'],
    },
    'Fatty Acids & Choline': {
      total: nutrientCategoryList['Fatty Acids & Choline'].length,
      satisfied: 0,
      color: nutrientCategoryColors['Fatty Acids & Choline'],
    },
    Total: { total: 0, satisfied: 0, color: nutrientCategoryColors.Total },
  };

  (Object.keys(nutrientCategoryList) as NutrientCategory[]).forEach((category) => {
    if (category !== 'Total') {
      const nutrientsInCategory = nutrientCategoryList[category];
      const satisfied = nutrientsInCategory.reduce((count, nutrient) => {
        return nutrients[nutrient] >= 5 ? count + 1 : count;
      }, 0);
      categorized[category].satisfied = satisfied;
    }
  });

  const allNutrients = (Object.keys(nutrientCategoryList) as NutrientCategory[])
    .flatMap((category) => {
      if (category !== 'Total') {
        return nutrientCategoryList[category];
      }
      return [];
    }).length;

  const satisfiedTotal = (Object.keys(nutrientCategoryList) as NutrientCategory[])
    .flatMap((category) => {
      if (category !== 'Total') {
        return nutrientCategoryList[category];
      }
      return [];
    })
    .reduce((count, nutrient) => {
      return nutrients[nutrient] >= 5 ? count + 1 : count;
    }, 0);

  categorized['Total'].total = allNutrients;
  categorized['Total'].satisfied = satisfiedTotal;

  return categorized;
};

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
  const [selectedIngredient, setSelectedIngredient] = useState<string>('Full Meal');
  const [selectedNutrientData, setSelectedNutrientData] = useState<{ [key: string]: number }>({});
  const [missingNutrients, setMissingNutrients] = useState<string[]>([]);
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [highlightedNutrients, setHighlightedNutrients] = useState<string[]>([]);
  const [normalMealNutrients, setNormalMealNutrients] = useState<{ [key: string]: number }>({});
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  /*=================================================================================================*/

  const [baseNutrients, setBaseNutrients] = useState<{ [key: string]: number }>({});
  const [originalMissingNutrients, setOriginalMissingNutrients] = useState<string[]>([]);

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
        setNutrients(response.nutrients || {});
        setMissingNutrients(response.missingNutrients || []);
        setOriginalMissingNutrients(response.missingNutrients || []);
        setSuggestions(response.suggestions || []);
        setSelectedIngredient('Full Meal');
        setSelectedNutrientData({});
        setHighlightedNutrients([]);

        const totalNutrients: { [key: string]: number } = {};
        (response.ingredients || []).forEach((ing) => {
          const nutrientData = response.nutrients[ing] || {};
          for (const nutrient in nutrientData) {
            if (Object.prototype.hasOwnProperty.call(nutrientData, nutrient)) {
              totalNutrients[nutrient] = (totalNutrients[nutrient] || 0) + nutrientData[nutrient];
            }
          }
        });
        setBaseNutrients(totalNutrients);
        setNormalMealNutrients(totalNutrients);
        setSelectedNutrientData(totalNutrients);
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

  /*=================================================================================================*/

  const handleIngredientClick = (ingredient: string) => {
    console.log('Clicked ingredient:', ingredient);
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

  /*=================================================================================================*/

  // Determine Deficiencies
  const determineLowAndMissingNutrients = (nutrients: { [key: string]: number }) => {
    const missing = [];
    for (const nutrient in nutrients) {
      if (nutrients[nutrient] < 20) {
        missing.push(nutrient);
      }
    }
    return missing;
  };

  const handleRecommendationClick = async (suggestion: string) => {
    try {
      const response = await fetch('http://3.84.94.100:5000/fetch-nutrient-data', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          foodDescription: suggestion,
          currentNutrients: baseNutrients,
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Server Error: ${errorText}`);
      }
      const data = await response.json();
      const updatedNutrients = data.nutrients || {};
      const changedNutrients = data.changedNutrients || [];

      // Update State
      setSelectedIngredient(suggestion);
      setSelectedNutrientData(updatedNutrients); 
      setHighlightedNutrients(changedNutrients);
      const updatedMissingNutrients = determineLowAndMissingNutrients(updatedNutrients);
      setMissingNutrients(updatedMissingNutrients);
    } catch (error) {
      console.error('Error fetching nutrient data for recommendation:', error);
      setError('Failed to fetch nutrient data for the recommendation.');
    }
  };

  /*=================================================================================================*/

  // Memo-izing
  const categorizedSelectedNutrients = useMemo(() => {
    return categorizeNutrients(selectedNutrientData);
  }, [selectedNutrientData]);

  return (
    <div className="min-h-screen flex items-center justify-center">

        <div
          className="relative w-full max-w-7.5xl h-auto bg-cover bg-center rounded-lg shadow-lg overflow-hidden"
        >
          <div className="absolute inset-0 bg-black opacity-20"></div>
          <div className="relative z-10 p-8">
          <h1 className="text-3xl font-bold text-center mb-8 text-white">
            The Nutrimancer's Codex Vol. II
          </h1>

          {/* Food Input Section */}
          <div className="flex justify-center mb-8">
            <input
              type="text"
              className="w-1/3 p-3 rounded-md focus:outline-none bg-white text-black mr-5 text-xl"
              placeholder="What have you consumed today?"
              value={food}
              onChange={(e) => setFood(e.target.value)}
              onKeyPress={(e) => {
                if (e.key === 'Enter') handleFoodSubmit();
              }}
            />
            <button
              onClick={handleFoodSubmit}
              disabled={loading || !food.trim()}
              className={`p-3 bg-green-500 text-white rounded-md text-xl ${
                loading || !food.trim() ? 'bg-green-300 cursor-not-allowed' : 'hover:bg-green-600'
              }`}
            >
              {loading ? 'Extracting...' : 'Extract Essence'}
            </button>
          </div>
          {error && <p className="text-center text-red-500 mb-4">{error}</p>}



          {/* Panels Container */}
          {!loading && !error && ingredients && ingredients.length > 0 && (
            <div className="flex justify-center gap-8 w-full px-4">


              {/* Ingredients Panel - Left */}
              <div className="w-[200px] flex-shrink-0">
                <IngredientsPanel
                  ingredients={ingredients}
                  onIngredientClick={handleIngredientClick}
                />
              </div>
              {/* Orbs Panel - Center */}
              <div className="w-[650px] flex-shrink-0">
                <OrbsPanel
                  nutrientData={categorizedSelectedNutrients}
                  selectedIngredient={selectedIngredient}
                  selectedNutrientData={selectedNutrientData}
                  highlightedNutrients={highlightedNutrients}
                  missingNutrients={missingNutrients}
                />
              </div>
              {/* Suggestion Panel - Right */}
              <div className="w-[200px] flex-shrink-0">
                <SuggestionPanel
                  missingNutrients={originalMissingNutrients}
                  suggestions={suggestions}
                  onRecommendationClick={handleRecommendationClick}
                />

              </div>
            </div>)}
        </div>
      </div>
    </div>
  );
};

export default App;