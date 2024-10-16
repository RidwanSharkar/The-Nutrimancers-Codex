// src/components/SuggestionPanel.tsx

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
        if (secondPartLower.includes('pass') || secondPartLower.includes('region') || secondPartLower.includes('store') || secondPartLower.includes('other')) {
          return titleCase(parts[0]);
      } else {
        // Else, display the first two parts joined by a comma
        return titleCase(parts.slice(0, 2).join(', '));
      }
    } else {
      // If there's only one part, display it
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
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">Essence Analysis:</h2>
      {missingNutrients && missingNutrients.length > 0 ? (
        <>
          <h3 className="text-lg font-medium mb-2 text-white">Deficient Bio-Parameters:</h3>
          <ul className="list-disc list-inside mb-4 space-y-1">
            {missingNutrients.map((nutrient, index) => (
              <li key={index} className="text-white">
                {nutrient}
              </li>
            ))}
          </ul>
          {processedSuggestions.length > 0 && (
            <>
              <h3 className="text-lg font-medium mb-2 text-white">Consider Harvesting:</h3>
              <div className="flex flex-wrap gap-2">
                {processedSuggestions.map((suggestion, index) => (
                  <button
                    key={index}
                    onClick={() => onRecommendationClick(suggestion)}
                    className="bg-[#FFC09F] hover:bg-[#EF8354] text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
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
        <p className="text-[#CEF7A0]">You have Ascended.</p>
      )}
    </div>
  );
};

export default SuggestionPanel;