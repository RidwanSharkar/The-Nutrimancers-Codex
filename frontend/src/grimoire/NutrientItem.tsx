// src/grimoire/NutrientItem.tsx

import React, { useState } from 'react';

interface NutrientItemProps {
  displayName: string;
  percentage?: number;
  classificationColor: string;
}

const NutrientItem: React.FC<NutrientItemProps> = ({ displayName, percentage, classificationColor }) => {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <span
      className={`font-medium ${classificationColor} cursor-pointer`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {isHovered && percentage !== undefined ? `${percentage.toFixed(2)}%` : displayName}
      {percentage !== undefined && (
        <span className="ml-1"></span>
      )}
    </span>
  );
};

export default NutrientItem;
