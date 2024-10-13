// src/components/NutrientsPanel.tsx

import React from 'react';

interface NutrientsPanelProps {
  ingredient: string;
  nutrients: { [key: string]: number };
  highlightedNutrients: string[]; 
  missingNutrients: string[];
}

const NutrientsPanel: React.FC<NutrientsPanelProps> = ({
  ingredient,
  nutrients,
  highlightedNutrients,
  missingNutrients,
}) => {

  const allNutrients = [
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
    'Histidine',
    'Isoleucine',
    'Leucine',
    'Lysine',
    'Methionine',
    'Phenylalanine',
    'Threonine',
    'Tryptophan',
    'Valine',
    'Alpha-Linolenic Acid',
    'Linoleic Acid',
    'EPA',
    'DHA',
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
    'Choline',
  ];

  // Severity Classification
  const classifyNutrient = (percentage: number | undefined) => {
    if (percentage === undefined || percentage === 0) {
      return 'none';
    } else if (percentage > 0 && percentage <= 5) {
      return 'low';
    } else if (percentage > 5 && percentage <= 20) {
      return 'average';
    } else {
      return 'high';
    }
  };

  // Severity Coloring
  const getColor = (classification: string, nutrient: string) => {
    if (highlightedNutrients.includes(nutrient) && missingNutrients.includes(nutrient)) {
      return 'black';
    }
    switch (classification) {
      case 'none':
        return 'gray';
      case 'low':
        return 'red';
      case 'average':
        return 'yellow';
      case 'high':
        return 'green';
      default:
        return 'gray';
    }
  };

  return (
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">
        Extracted Bioessence from: {ingredient || 'Selected Ingredient'}
      </h2>
      <ul className="space-y-1">
        {allNutrients.map((nutrient, index) => {
          const percentage = nutrients ? nutrients[nutrient] : undefined;
          const classification = classifyNutrient(percentage);
          const color = getColor(classification, nutrient);

          return (
            <li key={index} className="text-white">
              <span style={{ color: color }}>
                {nutrient}:{' '}
                {classification === 'none' ? ' ' : `${percentage?.toFixed(1)}% of RDA`}
              </span>
            </li>
          );
        })}
      </ul>
    </div>
  );
};

export default NutrientsPanel;