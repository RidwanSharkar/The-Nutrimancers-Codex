// src/grimoire/OrbsPanel.tsx

import React, { useEffect, useRef } from 'react';
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
}

const OrbsPanel: React.FC<OrbsPanelProps> = ({ nutrientData }) => {
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

  return (
    <div className="flex flex-wrap justify-center gap-8 mt-8">
      {(Object.keys(nutrientData) as NutrientCategory[]).map((category) => (
        <div key={category} className="flex flex-col items-center">
          <div className="relative w-32 h-32">
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
  );
};

export default OrbsPanel;




/*
TO TRY:


<img src="/decorative-border.svg" alt="Decorative Border" className="w-full h-full object-cover" />


<div
  className="absolute inset-0 rounded-full bg-gray-300 animate-pulse"
  style={{
    background: `conic-gradient(${data.color} var(--fill-percentage, 0%), #e0e0e0 0%)`,
  }}
  ref={(el) => (orbRefs.current[category] = el)}
></div>

.orb-container:hover {
  animation: rotate 5s linear infinite;
  box-shadow: 0 0 20px rgba(255, 255, 255, 0.5);
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
<div className="absolute inset-0 rounded-full bg-gray-300 orb-container" ... ></div>

*/ 