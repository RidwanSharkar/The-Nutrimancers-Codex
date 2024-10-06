// src/components/IngredientsPanel.tsx

import React from 'react';

interface IngredientsPanelProps {
  ingredients: string[];
}

const IngredientsPanel: React.FC<IngredientsPanelProps> = ({ ingredients }) => {
  return (
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">Ingredients</h2>
      {ingredients.length > 0 ? (
        <ul className="space-y-2">
          {ingredients.map((ingredient) => (
            <li key={ingredient} className="text-white">
              • {ingredient}
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

