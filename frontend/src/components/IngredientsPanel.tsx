// src/components/IngredientsPanel.tsx

import React from 'react';

interface IngredientsPanelProps {
  ingredients: string[];
  onIngredientClick: (ingredient: string) => void;
}

const IngredientsPanel: React.FC<IngredientsPanelProps> = ({ ingredients, onIngredientClick }) => {
  // Combine "Full Meal" with the list of ingredients
  const allItems = ['Full Meal', ...ingredients];

  return (
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">Ingredients</h2>
      {allItems.length > 0 ? (
        <div className="flex flex-col gap-2">
          {allItems.map((item, index) => (
            <button
              key={index}
              onClick={() => onIngredientClick(item)}
              className="bg-[#FFC09F] hover:bg-[#EF8354] text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
            >
              {item.replace(/^- /, '')}
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
