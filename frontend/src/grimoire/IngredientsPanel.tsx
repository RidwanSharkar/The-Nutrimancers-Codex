// src/grimoire/IngredientsPanel.tsx

import React from 'react';

interface IngredientsPanelProps {
  ingredients: string[];
  onIngredientClick: (ingredient: string) => void;
}

const IngredientsPanel: React.FC<IngredientsPanelProps> = ({ ingredients, onIngredientClick }) => {
  const allItems = ['Full Meal', ...ingredients];

   // style={{ minWidth: 'fit-content' }}
  return (
    <div className="parchment rounded-lg p-4 fade-in-up flex-1" style={{ minWidth: 'fit-content' }} >
      <h2 className="text-xl font-semibold mb-4 text-[#5d473a]" style={{ whiteSpace: 'nowrap' }} >
        Detected Bio-Sources:
      </h2>
      {allItems.length > 0 ? (
        <div className="flex flex-col gap-2 scroll-container">
          {allItems.map((item, index) => (
            <button
              key={index}
              onClick={() => onIngredientClick(item)}
              className="button-magical bg-[#fff8e1] hover:bg-[#c9a66b] text-[#5d473a] font-semibold py-2 px-4 rounded-lg transition duration-300"
            >
              {item.replace(/^- /, '')}
            </button>
          ))}
        </div>
      ) : (
        <p className="text-[#5d473a]">No Bioessence Detected</p>
      )}
    </div>
  );
};

export default IngredientsPanel;
