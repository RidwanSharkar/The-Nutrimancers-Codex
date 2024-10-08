// src/components/IngredientsPanel.tsx

import React from 'react';

interface IngredientsPanelProps {
  ingredients: string[];
  onIngredientClick: (ingredient: string) => void;
}

const IngredientsPanel: React.FC<IngredientsPanelProps> = ({ ingredients, onIngredientClick }) => {
  return (
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">Ingredients</h2>
      {ingredients.length > 0 ? (
        <div className="flex flex-wrap gap-2">
          {ingredients.map((ingredient, index) => (
            <button
              key={index}
              onClick={() => onIngredientClick(ingredient)}
              className="bg-[#FFC09F] hover:bg-[#EF8354] text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
            >
              {ingredient.replace(/^- /, '')}
            </button>
          ))}
        </div>
      ) : (
        <p className="text-white">No ingredients to display.</p>
      )}
    </div>
  );
};

export default IngredientsPanel;
