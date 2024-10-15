import React, { useEffect, useRef, useState } from 'react';
import { gsap } from 'gsap';

type NutrientCategory = 'Minerals' | 'Vitamins' | 'Amino Acids' | 'Fatty Acids & Choline' | 'Total';

interface OrbsPanelProps {
  nutrientData: {
    [category in NutrientCategory]: {
      total: number;
      satisfied: number;
      color: string;
    };
  };
  selectedIngredient: string;
  selectedNutrientData: { [key: string]: number };
  highlightedNutrients: string[];
  missingNutrients: string[];
}

const nutrientCategoryList = {
  Minerals: [
    'Potassium', 'Sodium', 'Calcium', 'Phosphorus', 'Magnesium',
    'Iron', 'Zinc', 'Manganese', 'Copper', 'Selenium',
  ],
  Vitamins: [
    'Vitamin A', 'Vitamin B1', 'Vitamin B2', 'Vitamin B3', 'Vitamin B5',
    'Vitamin B6', 'Vitamin B9', 'Vitamin B12', 'Vitamin C', 'Vitamin D',
    'Vitamin E', 'Vitamin K',
  ],
  'Amino Acids': [
    'Histidine', 'Isoleucine', 'Leucine', 'Lysine', 'Methionine',
    'Phenylalanine', 'Threonine', 'Tryptophan', 'Valine',
  ],
  'Fatty Acids & Choline': [
    'Alpha-Linolenic Acid', 'Linoleic Acid', 'EPA', 'DHA', 'Choline',
  ],
  Total: [],
};

const OrbsPanel: React.FC<OrbsPanelProps> = ({
  nutrientData,
  selectedIngredient,
  selectedNutrientData,
  highlightedNutrients,
  missingNutrients,
}) => {
  const orbRefs = useRef<{ [key in NutrientCategory]: HTMLDivElement | null }>({
    Minerals: null,
    Vitamins: null,
    'Amino Acids': null,
    'Fatty Acids & Choline': null,
    Total: null,
  });

  const [selectedCategory, setSelectedCategory] = useState<NutrientCategory | null>(null);

  useEffect(() => {
    (Object.keys(nutrientData) as NutrientCategory[]).forEach((category) => {
      const orb = orbRefs.current[category];
      if (orb) {
        gsap.to(orb, {
          '--fill-percentage': `${(nutrientData[category].satisfied / nutrientData[category].total) * 100}%`,
          duration: 1.5,
          ease: 'power2.out',
        });
      }
    });
  }, [nutrientData]);

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
    <div className="flex flex-col items-center">
      <div className="flex flex-wrap justify-center gap-8 mb-8">
        {(Object.keys(nutrientData) as NutrientCategory[]).map((category) => (
          <div key={category} className="flex flex-col items-center">
            <div className="relative w-32 h-32 cursor-pointer" onClick={() => setSelectedCategory(category)}>
              <div
                className="absolute inset-0 rounded-full orb-container"
                style={{
                  background: `conic-gradient(${nutrientData[category].color} var(--fill-percentage, 0%), #e0e0e0 0%)`,
                }}
                ref={(el) => (orbRefs.current[category] = el)}
              ></div>
              <div className="absolute inset-0 flex items-center justify-center rounded-full bg-white bg-opacity-70">
                <span className="text-lg font-bold text-gray-800 text-center px-2">{category}</span>
              </div>
            </div>
            <span className="mt-2 text-sm text-gray-700">
              {nutrientData[category].satisfied}/{nutrientData[category].total} Nutrients
            </span>
          </div>
        ))}
      </div>
      
      {selectedCategory && (
        <div className="bg-[#F48668] rounded-lg p-4 w-full">
          <h2 className="text-xl font-semibold mb-4 text-white">
            Extracted Bioessence from: {selectedIngredient || 'Selected Ingredient'}
          </h2>
          <h3 className="text-lg font-medium mb-2 text-white">{selectedCategory}</h3>
          <ul className="space-y-1">
            {nutrientCategoryList[selectedCategory].map((nutrient, index) => {
              const percentage = selectedNutrientData ? selectedNutrientData[nutrient] : undefined;
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
      )}
    </div>
  );
};

export default OrbsPanel;