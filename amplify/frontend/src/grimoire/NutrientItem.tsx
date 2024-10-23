// src/grimoire/NutrientItem.tsx

import React from 'react';

interface NutrientItemProps {
  displayName: string;
  percentage?: number;
  classificationColor: string;
}

const NutrientItem: React.FC<NutrientItemProps> = ({
  displayName,
  percentage,
  classificationColor,
}) => {
  const isPresent = percentage !== undefined && percentage > 0;
  const percentageText = isPresent ? `${percentage.toFixed(2)}%` : '';

  // Empty/Missing Nutrient Color
  const textColorClass = isPresent ? classificationColor : 'text-gray-400';

  return (
    <span
      className={`font-medium text-lg ${textColorClass} cursor-pointer relative inline-block group w-full`}
    >
      {/* Nutrient*/}
      <span
        className={`absolute inset-0 flex items-center justify-start transition-opacity duration-200 ${
          isPresent ? 'opacity-100 group-hover:opacity-0' : 'opacity-100'
        }`}
      >
        {displayName}
      </span>
      {/* Percent */}
      {isPresent && (
        <span
          className="absolute inset-0 text-xl flex items-center justify-start opacity-0 transition-opacity duration-200 group-hover:opacity-100"
        >
          {percentageText}
        </span>
      )}
    </span>
  );
};

export default NutrientItem;
