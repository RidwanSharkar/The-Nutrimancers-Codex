// src/components/IngredientsPanel.tsx

import React from 'react';

interface IngredientsPanelProps {
  ingredients: string[];
  onSelectIngredient: (ingredient: string) => void;
}

const IngredientsPanel: React.FC<IngredientsPanelProps> = ({ ingredients, onSelectIngredient }) => {
  return (
    <div className="bg-[#F48668] rounded-lg p-4">
      <h2 className="text-xl font-semibold mb-4 text-white">Ingredients</h2>
      {ingredients.length > 0 ? (
        <ul className="space-y-2">
          {ingredients.map((ingredient) => (
            <li key={ingredient}>
              <button
                className="w-full px-4 py-2 bg-amber-500 text-white rounded hover:bg-amber-600 transition-colors"
                onClick={() => onSelectIngredient(ingredient)}
              >
                {ingredient}
              </button>
            </li>
          ))}
        </ul>
      ) : (
        <p className="text-white">No ingredients to display.</p>
      )}
    </div>
  );
};

export default IngredientsPanel;
