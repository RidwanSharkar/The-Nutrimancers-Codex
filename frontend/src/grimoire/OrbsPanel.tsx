// src/grimoire/OrbsPanel.tsx

import React, { useEffect } from 'react';
import { gsap } from 'gsap';
import LiquidGauge from 'react-liquid-gauge';
import NutrientItem from './NutrientItem';

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
}) => {

  const classifyNutrient = (
    percentage: number | undefined
  ): 'none' | 'low' | 'average' | 'high' => {
    if (percentage === undefined || percentage === 0) return 'none';
    if (percentage > 0 && percentage <= 3) return 'low';
    if (percentage > 3 && percentage <= 20) return 'average';
    return 'high';
  };

  // Define a mapping from classification to Tailwind color classes
  const classificationTextColorMap: { 
    [key in 'none' | 'low' | 'average' | 'high']: string 
  } = {
    none: 'text-gray-500',
    low: 'text-red-500',
    average: 'text-yellow-500',
    high: 'text-green-500',
  };
  const renderNutrientList = (category: Exclude<NutrientCategory, 'Total'>) => (
    <div className="parchment rounded-lg p-4 w-full mt-1 fade-in-up max-h-[415px] max-w-[125px] overflow-auto">
      <h3 className="text-lg font-semibold mb-2 text-[#5d473a] whitespace-nowrap">
        {category === 'Fatty Acids & Choline' ? 'Fatty Acids' : category}
      </h3>
      <ul className="scroll-container text-lg leading-tight mb-1">
        {nutrientCategoryList[category].map((nutrient: string, index: number) => {
          if (nutrient === 'CHOLINE_BREAK') return <li key={index}>&nbsp;</li>;
          const percentage = selectedNutrientData ? selectedNutrientData[nutrient] : undefined;
          const classification = classifyNutrient(percentage);
          const displayName = displayNameMap[nutrient] || nutrient;

          return (
            <li
              key={index}
              className={`relative py-2 nutrient-list-item`}
            >
              <div className="flex items-center">
                <NutrientItem
                  displayName={displayName}
                  percentage={percentage}
                  classificationColor={classificationTextColorMap[classification]}
                />
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
  


  // Type for waveRenderer
  interface WaveRendererProps {
    path: string;
  }

  const renderOrb = (category: NutrientCategory) => {
    const percentageFilled =
      (nutrientData[category].satisfied / nutrientData[category].total) * 100;

    const orbSize = category === 'Total' ? 144 : 128;

    return (
      <div className="flex flex-col items-center relative">
        <div
          className={`relative parchment-orb ${percentageFilled >= 100 ? 'glow' : ''}`}
          style={{
            filter: `drop-shadow(0 0 10px ${nutrientData[category].color})`,
            '--orb-color': `${nutrientData[category].color}`,
            width: `${orbSize}px`,
            height: `${orbSize}px`,
          } as React.CSSProperties & { [key: string]: string | number }}
        >
          <LiquidGauge
            style={{ margin: '0 auto' }}
            width={orbSize}
            height={orbSize}
            value={percentageFilled}
            textRenderer={() => null}
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
              fill: 'none', 
            }}
            waveStyle={{
              fill: `url(#waveGradient-${category})`,
            }}

            
            waveRenderer={(props: WaveRendererProps) => {
              const { path } = props;
              return (
                <>
                  <defs>
                    {/* Gradient for wave */}
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
                        stopOpacity="0.2"
                      />
                      <stop
                        offset="100%"
                        stopColor={`${nutrientData[category].color}`}
                        stopOpacity="0.6"
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
            width={orbSize}
            height={orbSize}
            style={{ position: 'absolute', top: 0, left: 0 }}
          >
            {/* Glass outline */}
            <circle
              cx={orbSize / 2}
              cy={orbSize / 2}
              r={(orbSize / 2) - 4}
              fill="none"
              stroke="rgba(255,255,255,0.5)"
              strokeWidth="2"
            />
            {/* Top reflection */}
            <ellipse
              cx={orbSize / 2}
              cy={orbSize * 0.3}
              rx={orbSize * 0.3}
              ry={orbSize * 0.12}
              fill="rgba(255,255,255,0.2)"
            />
            {/* Bottom reflection */}
            <ellipse
              cx={orbSize / 2}
              cy={orbSize * 0.7}
              rx={orbSize * 0.2}
              ry={orbSize * 0.08}
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
    'Vitamins',
    'Minerals',
    'Amino Acids',
    'Fatty Acids & Choline',
  ];

  useEffect(() => {
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

 

      <div className="flex flex-row justify-center gap-0 w-full">
        {mainCategories.map((category) => (
          <div key={category} className="flex flex-col items-center w-1/4 min-w-[155px]">
            {renderOrb(category)}
            {renderNutrientList(category)}
          </div>
        ))}
      </div>

      <div className=" flex justify-center">
        {renderOrb('Total')}
      </div>


      
    </div>
  );
};

export default OrbsPanel;