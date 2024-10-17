// src/grimoire/SuggestionPanel.tsx

import React, { useMemo } from 'react';

interface SuggestionPanelProps {
  missingNutrients: string[] | null;
  suggestions: string[] | null;
  onRecommendationClick: (suggestion: string) => void;
}

const titleCase = (str: string): string => {
  return str
    .toLowerCase()
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
};

/*=============================================================================================*/

const processSuggestion = (suggestion: string): string => {
  const parts = suggestion.split(',').map(part => part.trim());

  if (parts.length === 0) return '';

  const firstPartLower = parts[0].toLowerCase();

  if (firstPartLower === 'fatty acids' || firstPartLower === 'amino acids') {
    if (parts.length >= 2) {
      return titleCase(parts[1]);
    } else {
      return titleCase(parts[0]);
    }
  } else {
    if (parts.length >= 2) {
      const secondPartLower = parts[1].toLowerCase(); // PLACEHOLDER TILL FIND MORE for logic swap
      if (
        secondPartLower.includes('pass') ||
        secondPartLower.includes('region') ||
        secondPartLower.includes('store') ||
        secondPartLower.includes('other')
      ) {
        return titleCase(parts[0]);
      } else {
        return titleCase(parts.slice(0, 2).join(', '));
      }
    } else {
      return titleCase(parts[0]);
    }
  }
};

/*=============================================================================================*/

const SuggestionPanel: React.FC<SuggestionPanelProps> = ({
  missingNutrients,
  suggestions,
  onRecommendationClick,
}) => {
  const processedSuggestions = useMemo(() => {
    if (!suggestions) return [];
    const processed = suggestions.map(processSuggestion).filter(s => s !== '');
    const unique = Array.from(new Set(processed));
    return unique;
  }, [suggestions]);

  return (
    <div className="parchment rounded-lg p-4 fade-in-up flex-1">
      <h2 className="text-xl font-semibold mb-4 text-[#5d473a]" style={{ whiteSpace: 'nowrap' }}>Essence Analysis:</h2>
      {missingNutrients && missingNutrients.length > 0 ? (
        <>
          <h3 className="text-lg font-medium mb-2 text-[#5d473a]">You're low on:</h3>
          <ul className="list-disc list-inside mb-4 space-y-1 scroll-container">
            {missingNutrients.map((nutrient, index) => (
              <li key={index} className="text-[#5d473a]">
                {nutrient}
              </li>
            ))}
          </ul>
          {processedSuggestions.length > 0 && (
            <>
              <h3 className="text-lg font-medium mb-2 text-[#5d473a]">Consider:</h3>
              <div className="flex flex-wrap gap-2 scroll-container">
                {processedSuggestions.map((suggestion, index) => (
                  <button
                    key={index}
                    onClick={() => onRecommendationClick(suggestion)}
                    className="button-magical bg-[#fff8e1] hover:bg-[#c9a66b] text-[#5d473a] font-semibold py-2 px-4 rounded-lg transition duration-300"
                    title={suggestion}
                  >
                    {suggestion}
                  </button>
                ))}
              </div>
            </>
          )}
        </>
      ) : (
        <p className="text-[#a6784a]">You have Ascended.</p>
      )}
    </div>
  );
};

export default SuggestionPanel;
