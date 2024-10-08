// src/components/NutrientsPanel.tsx

import React from 'react';

// Define the props interface
interface NutrientsPanelProps {
  ingredient: string;
  nutrients: { [key: string]: number };
}

// Define the NutrientsPanel component
const NutrientsPanel: React.FC<NutrientsPanelProps> = ({ ingredient, nutrients }) => {
  console.log("NutrientsPanel props:", { ingredient, nutrients }); 

  return (
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">
        Nutrients for {ingredient || 'Selected Ingredient'}
      </h2>
      {Object.keys(nutrients).length > 0 ? (
        <ul className="list-disc list-inside space-y-1">
          {Object.entries(nutrients).map(([nutrient, value], index) => (
            <li key={index} className="text-white">
              {nutrient}: {value.toFixed(2)}
            </li>
          ))}
        </ul>
      ) : (
        <p className="text-white">Select an ingredient to view its nutrients.</p>
      )}
    </div>
  );
};

export default NutrientsPanel;
