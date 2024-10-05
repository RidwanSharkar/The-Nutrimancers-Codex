// src/App.tsx

import React, { useState } from 'react';
import IngredientsPanel from './components/IngredientsPanel';
import NutrientsPanel from './components/NutrientsPanel';
import SuggestionPanel from './components/SuggestionPanel';
import './App.css'; // Tailwind CSS

type NutrientData = {
  [key: string]: string[];
};

type SuggestionData = {
  [key: string]: string;
};

const App: React.FC = () => {
  const [food, setFood] = useState<string>('');
  const [ingredients, setIngredients] = useState<string[]>([]);
  const [nutrients, setNutrients] = useState<string[]>([]);
  const [missingNutrients, setMissingNutrients] = useState<string[]>([]);
  const [suggestions, setSuggestions] = useState<string[]>([]);

  const nutrientData: NutrientData = {
    Dough: ['Carbohydrates', 'Fiber'],
    Sauce: ['Vitamins', 'Antioxidants'],
    Cheese: ['Minerals', 'Fatty acids', 'Choline'],
    // DISCONTINTUED TEST --> Gemini
  };

  const allEssentialNutrients: string[] = [
    'Potassium', 
    'Chloride', 
    'Sodium', 
    'Calcium', 
    'Phosphorus', 
    'Magnesium', 
    'Iron', 
    'Zinc', 
    'Manganese', 
    'Copper', 
    'Iodine', 
    'Chromium', 
    'Molybdenum', 
    'Selenium',

    'Histidine',
    'Isoleucine',
    'Leucine',
    'Lysine',
    'Methionine',
    'Phenylalanine',
    'Threonine',
    'Tryptophan',
    'Valine',

    'Alpha-Linolenic Acid', // Omega-3
    'Linoleic Acid',       // Omega-6


    'Vitamin A', //(all-trans-retinols, all-trans-retinyl-esters, as well as all-trans-β-carotene and other provitamin A carotenoids)
    'Vitamin B1', //(thiamine)
    'Vitamin B2', //(riboflavin)
    'Vitamin B3', //(niacin)
    'Vitamin B5', //(pantothenic acid)
    'Vitamin B6', //(pyridoxine)
    'Vitamin B7', //(biotin)
    'Vitamin B9', //(folic acid and folates)
    'Vitamin B12', //(cobalamins)
    'Vitamin C', //(ascorbic acid and ascorbates)
    'Vitamin D', //(calciferols)
    'Vitamin E', //(tocopherols and tocotrienols)
    'Vitamin K', // (phylloquinones, menaquinones, and menadiones)

    'choline'
  ];

  const suggestionData: SuggestionData = {
    Fiber: 'Consider eating whole grains or fruits.',
    Protein: 'How about some lean meat or legumes?',
    Calcium: 'Dairy products or leafy greens can help.',

  };

  const handleFoodSubmit = () => {
    if (food.trim()) {
      setFood(food.trim());
      // Define ingredients based on the food input.
      const mockIngredients: { [key: string]: string[] } = {
        pizza: ['Dough', 'Sauce', 'Cheese'],
        salad: ['Lettuce', 'Tomatoes', 'Cucumbers'],
        burger: ['Bun', 'Patty', 'Cheese'],
  
      };
      const fetchedIngredients = mockIngredients[food.toLowerCase()] || [];
      setIngredients(fetchedIngredients);
      setNutrients([]);
      setMissingNutrients([]);
      setSuggestions([]);
    }
  };

  const handleIngredientSelect = (ingredient: string) => {
    const fetchedNutrients = nutrientData[ingredient] || [];
    setNutrients(fetchedNutrients);
    // Determine missing nutrients
    const missing = allEssentialNutrients.filter(
      (nutrient) => !fetchedNutrients.includes(nutrient)
    );
    setMissingNutrients(missing);
    // Generate suggestions
    const generatedSuggestions = missing
      .filter((nutrient) => suggestionData[nutrient])
      .map((nutrient) => suggestionData[nutrient]);
    setSuggestions(generatedSuggestions);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#F48668]"> 
      <div className="w-[1024px] h-[768px] bg-[#FFC09F] p-8 rounded-lg shadow-lg overflow-hidden">
        <h1 className="text-4xl font-bold text-center mb-8">Bioessence</h1>
        
      {/* Input Box */}
      <div className="flex justify-center mb-8">
        <input
          type="text"
          className="w-1/3 p-3 rounded-md focus:outline-none bg-[#EF8354] text-white mr-5"  // Added "mr-2" for margin-right
          placeholder="What have you eaten today?"
          value={food}
          onChange={(e) => setFood(e.target.value)}
          onKeyPress={(e) => {
            if (e.key === 'Enter') handleFoodSubmit();
          }}
        />
        <button
          onClick={handleFoodSubmit}
          disabled={!food.trim()}
          className={`p-3 bg-green-500 text-white rounded-md ${
            !food.trim() ? 'bg-green-300 cursor-not-allowed' : 'hover:bg-green-600'
          }`}  
        >
          Submit
        </button>
      </div>

        {/* Panels */}
        <div className="flex justify-center lg:justify-center gap-8 mx-auto w-full max-w-5xl">
        {/* Ingredients Panel */}
        <IngredientsPanel
          ingredients={ingredients}
          onSelectIngredient={handleIngredientSelect}
        />
        {/* Nutrients Panel */}
        <NutrientsPanel nutrients={nutrients} />
        {/* Suggestion Panel */}
        <SuggestionPanel missingNutrients={missingNutrients} suggestions={suggestions} />
      </div>
      </div>
    </div>



  );
};

export default App;
