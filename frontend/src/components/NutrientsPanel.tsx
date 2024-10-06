// src/components/NutrientsPanel.tsx

import React from 'react';

interface NutrientsPanelProps {
  nutrients: { [key: string]: number };
}

const NutrientsPanel: React.FC<NutrientsPanelProps> = ({ nutrients }) => {
  return (
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">Nutrients</h2>
      {Object.keys(nutrients).length > 0 ? (
        <ul className="list-disc list-inside space-y-1">
          {Object.entries(nutrients).map(([nutrient, value]) => (
            <li key={nutrient} className="text-white">
              {nutrient}: {value}
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
