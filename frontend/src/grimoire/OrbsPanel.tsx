// src/grimoire/OrbsPanel.tsx

import React, { useEffect, useRef } from 'react';
import { gsap } from 'gsap';
import Particles from "react-tsparticles";

type NutrientCategory = 'Vitamins' | 'Minerals' | 'Amino Acids' | 'Fatty Acids & Choline' | 'Total';

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

// Abbreviate/Change Display 
const displayNameMap: { [key: string]: string } = {
  'Alpha-Linolenic Acid': 'ALA',
  'Linoleic Acid': 'LA',
};

const nutrientCategoryList: { [key in Exclude<NutrientCategory, 'Total'>]: string[] } = {
  Vitamins: [
    'Vitamin A', 'Vitamin B1', 'Vitamin B2', 'Vitamin B3', 'Vitamin B5',
    'Vitamin B6', 'Vitamin B9', 'Vitamin B12', 'Vitamin C', 'Vitamin D',
    'Vitamin E', 'Vitamin K',
  ],
  Minerals: [
    'Potassium', 'Sodium', 'Calcium', 'Phosphorus', 'Magnesium',
    'Iron', 'Zinc', 'Manganese', 'Copper', 'Selenium',
  ],
  'Amino Acids': [
    'Histidine', 'Isoleucine', 'Leucine', 'Lysine', 'Methionine',
    'Phenylalanine', 'Threonine', 'Tryptophan', 'Valine',
  ],
  'Fatty Acids & Choline': [
    'Linoleic Acid', 'Alpha-Linolenic Acid', 'EPA', 'DHA', 'CHOLINE_BREAK', 'Choline',
  ],
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

  const classifyNutrient = (percentage: number | undefined): 'none' | 'low' | 'average' | 'high' => {
    if (percentage === undefined || percentage === 0) return 'none';
    if (percentage > 0 && percentage <= 4) return 'low';
    if (percentage > 4 && percentage <= 15) return 'average';
    return 'high';
  };

  const getColor = (classification: 'none' | 'low' | 'average' | 'high', nutrient: string): string => {
    if (highlightedNutrients.includes(nutrient) && missingNutrients.includes(nutrient)) {
      return 'black';
    }
    switch (classification) {
      case 'none': return 'gray';
      case 'low': return 'red';
      case 'average': return 'yellow';
      case 'high': return 'green';
      default: return 'gray';
    }
  };

  const renderNutrientList = (category: Exclude<NutrientCategory, 'Total'>) => (
    <div className="bg-[#F48668] rounded-lg p-4 w-full mt-4">
      <h3 className="text-lg font-medium mb-2 text-white">
        {category === 'Fatty Acids & Choline' ? 'Fatty Acids' : category}
      </h3>
      <ul className="space-y-1">
        {nutrientCategoryList[category].map((nutrient: string, index: number) => {
          if (nutrient === 'CHOLINE_BREAK') return <li key={index}>&nbsp;</li>;
          const percentage = selectedNutrientData ? selectedNutrientData[nutrient] : undefined;
          const classification = classifyNutrient(percentage);
          const color = getColor(classification, nutrient);

          // Abbreviations
          const displayName = displayNameMap[nutrient] || nutrient;

          return (
            <li key={index} className="text-white opacity-0 nutrient-item group relative">
              <div className="flex items-center">
                <span style={{ color, fontWeight: '500' }} className="flex items-center">
                  {displayName}
                  {classification === 'low' && <span className="ml-1 text-red-500">*</span>}
                  {classification === 'average' && <span className="ml-1 text-yellow-500">*</span>}
                  {classification === 'high' && <span className="ml-1 text-green-500">*</span>}
                </span>
              </div>
              {classification !== 'none' && (
                <div
                className="absolute left-0 right-0 -top-2 mb-2 w-full text-white text-base rounded py-2 px-4 opacity-0 group-hover:opacity-100 transition-opacity duration-300 z-10"

                  style={{ backgroundColor: color }}
                >
                  {percentage?.toFixed(1)}%
                </div>
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );

  const renderOrb = (category: NutrientCategory) => {
    const isFull = nutrientData[category].satisfied === nutrientData[category].total;

    return (
      <div className="flex flex-col items-center relative">
        <div className={`relative w-32 h-32 transition-all duration-500 ${isFull ? 'glow' : ''}`}>
          <div
            className="absolute inset-0 rounded-full orb-container"
            style={{
              background: `conic-gradient(${nutrientData[category].color} var(--fill-percentage, 0%), #e0e0e0 0%)`,
            }}
            ref={(el) => (orbRefs.current[category] = el)}
          ></div>
          <div className="absolute inset-0 flex items-center justify-center rounded-full bg-white bg-opacity-70">
            <span className="text-lg font-bold text-gray-800 text-center px-2">
              {category === 'Fatty Acids & Choline' ? 'Fatty Acids' : category}
            </span>
          </div>

          {/* Particle Effect */}
          <div className="relative w-32 h-32">
            <Particles
              options={{
                particles: {
                  number: { value: 50 },
                  size: { value: 3 },
                  move: { speed: 1, direction: "none", outMode: "out" },
                  line_linked: { enable: false },
                  opacity: { value: 0.5 },
                  color: { value: "#ffffff" },
                },
              }}
              className="absolute inset-0 pointer-events-none"
            />
          </div>

          {/* Fluid Effect */}
          <div className="absolute inset-0 rounded-full pointer-events-none">
            <div className="overlay"></div>
          </div>
        </div>
        <span className="mt-2 text-sm text-gray-700">
          {nutrientData[category].satisfied}/{nutrientData[category].total} 
        </span>
      </div>
    );
  };

  const mainCategories: Exclude<NutrientCategory, 'Total'>[] = ['Minerals', 'Vitamins', 'Amino Acids', 'Fatty Acids & Choline'];

  useEffect(() => {
    // Animate nutrient list items
    gsap.to('.nutrient-item', {
      opacity: 1,
      stagger: 0.1,
      duration: 0.5,
      y: 0,
      ease: 'power2.out',
      delay: 0.5,
    });
  }, [selectedNutrientData]);

  return (
    <div className="flex flex-col items-center">
      <h2 className="text-2xl font-semibold mb-4 text-center text-white">
        Bioessence Extracted from: {selectedIngredient}
      </h2>

      {/* Total Orb */}
      <div className="mb-8 flex justify-center">
        {renderOrb('Total')}
      </div>

      {/* Main Orbs */}
      <div className="flex flex-row justify-center gap-8 w-full">
        {mainCategories.map((category) => (
          <div key={category} className="flex flex-col items-center w-1/4">
            {renderOrb(category)}
            {renderNutrientList(category)}
          </div>
        ))}
      </div>
    </div>
  );
};

export default OrbsPanel;
