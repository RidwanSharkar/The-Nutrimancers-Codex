// src/components/NutrientsPanel.tsx

import React from 'react';

interface NutrientsPanelProps {
  nutrients: string[];
}

const NutrientsPanel: React.FC<NutrientsPanelProps> = ({ nutrients }) => {
  return (
    <div className="bg-[#F48668] rounded-lg p-4">
      <h2 className="text-xl font-semibold mb-4 text-white">Nutrients</h2>
      {nutrients.length > 0 ? (
        <ul className="list-disc list-inside space-y-1">
          {nutrients.map((nutrient) => (
            <li key={nutrient} className="text-white">
              {nutrient}
            </li>
          ))}
        </ul>
      ) : (
        <p className="text-white">No nutrients to display.</p>
      )}
    </div>
  );
};

export default NutrientsPanel;
