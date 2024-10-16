// src/grimoire/OrbsPanel.tsx

import React, { useEffect } from 'react';
import { gsap } from 'gsap';
import LiquidGauge from 'react-liquid-gauge';

type NutrientCategory =
  | 'Vitamins'
  | 'Minerals'
  | 'Amino Acids'
  | 'Fatty Acids & Choline'
  | 'Total';

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

const nutrientCategoryList: {
  [key in Exclude<NutrientCategory, 'Total'>]: string[];
} = {
  Vitamins: [
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
  ],
  Minerals: [
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
  ],
  'Amino Acids': [
    'Histidine',
    'Isoleucine',
    'Leucine',
    'Lysine',
    'Methionine',
    'Phenylalanine',
    'Threonine',
    'Tryptophan',
    'Valine',
  ],
  'Fatty Acids & Choline': [
    'Linoleic Acid',
    'Alpha-Linolenic Acid',
    'EPA',
    'DHA',
    'CHOLINE_BREAK',
    'Choline',
  ],
};

const OrbsPanel: React.FC<OrbsPanelProps> = ({
  nutrientData,
  selectedIngredient,
  selectedNutrientData,
  highlightedNutrients,
  missingNutrients,
}) => {
  const classifyNutrient = (
    percentage: number | undefined
  ): 'none' | 'low' | 'average' | 'high' => {
    if (percentage === undefined || percentage === 0) return 'none';
    if (percentage > 0 && percentage <= 4) return 'low';
    if (percentage > 4 && percentage <= 15) return 'average';
    return 'high';
  };

  const getColor = (
    classification: 'none' | 'low' | 'average' | 'high',
    nutrient: string
  ): string => {
    if (highlightedNutrients.includes(nutrient) && missingNutrients.includes(nutrient)) {
      return '#5d473a'; // Dark brown for highlighted and missing
    }
    switch (classification) {
      case 'none':
        return '#7d7d7d'; // Gray
      case 'low':
        return '#d9534f'; // Red
      case 'average':
        return '#f0ad4e'; // Yellow
      case 'high':
        return '#5cb85c'; // Green
      default:
        return '#7d7d7d'; // Default Gray
    }
  };

  const renderNutrientList = (category: Exclude<NutrientCategory, 'Total'>) => (
    <div className="parchment rounded-lg p-4 w-full mt-4 fade-in-up">
      <h3 className="text-lg font-medium mb-2 text-[#5d473a]">
        {category === 'Fatty Acids & Choline' ? 'Fatty Acids' : category}
      </h3>
      <ul className="space-y-1 scroll-container">
        {nutrientCategoryList[category].map((nutrient: string, index: number) => {
          if (nutrient === 'CHOLINE_BREAK') return <li key={index}>&nbsp;</li>;
          const percentage = selectedNutrientData
            ? selectedNutrientData[nutrient]
            : undefined;
          const classification = classifyNutrient(percentage);
          const color = getColor(classification, nutrient);

          // Abbreviations
          const displayName = displayNameMap[nutrient] || nutrient;

          return (
            <li
              key={index}
              className="text-[#5d473a] opacity-0 nutrient-item group relative"
            >
              <div className="flex items-center">
                <span
                  style={{ color, fontWeight: 500 }}
                  className="flex items-center"
                >
                  {displayName}
                  {classification === 'low' && (
                    <span className="ml-1 text-red-500">*</span>
                  )}
                  {classification === 'average' && (
                    <span className="ml-1 text-yellow-500">*</span>
                  )}
                  {classification === 'high' && (
                    <span className="ml-1 text-green-500">*</span>
                  )}
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

  // Define type for waveRenderer
  interface WaveRendererProps {
    path: string;
  }

  const renderOrb = (category: NutrientCategory) => {
    const percentageFilled =
      (nutrientData[category].satisfied / nutrientData[category].total) * 100;

    return (
      <div className="flex flex-col items-center relative">
        <div
          className="relative parchment-orb"
          style={{
            filter: `drop-shadow(0 0 10px ${nutrientData[category].color})`,
          }}
        >
          <LiquidGauge
            style={{ margin: '0 auto' }}
            width={128}
            height={128}
            value={percentageFilled}
            textRenderer={() => null} // Hide default text
            riseAnimation
            waveAnimation
            waveFrequency={1}
            waveAmplitude={5}
            gradient
            gradientStops={[
              {
                key: '0%',
                stopColor: `${nutrientData[category].color}`,
                stopOpacity: 1,
                offset: '0%',
              },
              {
                key: '100%',
                stopColor: `${nutrientData[category].color}`,
                stopOpacity: 0.7,
                offset: '100%',
              },
            ]}
            circleStyle={{
              fill: 'none', //outer circle transparent
            }}
            waveStyle={{
              fill: `url(#waveGradient-${category})`,
            }}
            waveRenderer={(props: WaveRendererProps) => {
              const { path } = props;
              return (
                <>
                  <defs>
                    {/* Gradient for the wave */}
                    <linearGradient
                      id={`waveGradient-${category}`}
                      x1="0%"
                      y1="0%"
                      x2="0%"
                      y2="100%"
                    >
                      <stop
                        offset="0%"
                        stopColor={`${nutrientData[category].color}`}
                        stopOpacity="0.8"
                      />
                      <stop
                        offset="100%"
                        stopColor={`${nutrientData[category].color}`}
                        stopOpacity="0.4"
                      />
                    </linearGradient>
                  </defs>
                  <path d={path} fill={`url(#waveGradient-${category})`} />
                </>
              );
            }}
          />
          {/* Orb label */}
          <div className="absolute inset-0 flex items-center justify-center">
            <span className="text-lg font-bold text-[#5d473a] text-center px-2">
              {category === 'Fatty Acids & Choline' ? 'Fatty Acids' : category}
            </span>
          </div>
          {/* Glass effect overlay */}
          <svg
            width={128}
            height={128}
            style={{ position: 'absolute', top: 0, left: 0 }}
          >
            {/* Glass outline */}
            <circle
              cx={64}
              cy={64}
              r={60}
              fill="none"
              stroke="rgba(255,255,255,0.5)"
              strokeWidth="2"
            />
            {/* Top reflection */}
            <ellipse
              cx={64}
              cy={40}
              rx={35}
              ry={15}
              fill="rgba(255,255,255,0.2)"
            />
            {/* Bottom reflection */}
            <ellipse
              cx={64}
              cy={88}
              rx={25}
              ry={10}
              fill="rgba(255,255,255,0.1)"
            />
          </svg>
          {/* Shine effect */}
          <div className="orb-shine"></div>
        </div>
        {/* Satisfied / Total display */}
        <span className="mt-2 text-sm text-[#5d473a]">
          {nutrientData[category].satisfied}/{nutrientData[category].total}
        </span>
      </div>
    );
  };

  const mainCategories: Exclude<NutrientCategory, 'Total'>[] = [
    'Minerals',
    'Vitamins',
    'Amino Acids',
    'Fatty Acids & Choline',
  ];

  useEffect(() => {
    // Animate nutrient list items with fade-in-up
    gsap.from('.nutrient-item', {
      opacity: 0,
      y: 20,
      stagger: 0.1,
      duration: 0.6,
      ease: 'power2.out',
    });
  }, [selectedNutrientData]);

  return (
    <div className="flex flex-col items-center parchment fade-in-up p-4">
      <h2 className="text-2xl font-semibold mb-4 text-[#5d473a] text-center">
        Bioessence Extracted from: {selectedIngredient}
      </h2>

      {/* Main Orbs */}
      <div className="flex flex-row justify-center gap-8 w-full">
        {mainCategories.map((category) => (
          <div key={category} className="flex flex-col items-center w-1/4">
            {renderOrb(category)}
            {renderNutrientList(category)}
          </div>
        ))}
      </div>

      {/* Total Orb */}
      <div className="mb-8 flex justify-center">
        {renderOrb('Total')}
      </div>
    </div>
  );
};

export default OrbsPanel;
