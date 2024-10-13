// src/components/SuggestionPanel.tsx

import React from 'react';

interface SuggestionPanelProps {
  missingNutrients: string[] | null;
  suggestions: string[] | null;
  onRecommendationClick: (suggestion: string) => void;
}

const SuggestionPanel: React.FC<SuggestionPanelProps> = ({
  missingNutrients,
  suggestions,
  onRecommendationClick,
}) => {
  return (
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">Essence-Profile Analysis:</h2>
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
          {suggestions && suggestions.length > 0 && (
            <>
              <h3 className="text-lg font-medium mb-2 text-white">Consider Harvesting:</h3>
              <div className="flex flex-wrap gap-2">
                {suggestions.map((suggestion, index) => {
                  // Split the suggestion by commas and take the first two elements
                  const displayedText = suggestion
                    .split(',')
                    .slice(0, 2)
                    .map(part => part.trim())
                    .join(', ');

                  return (
                    <button
                      key={index}
                      onClick={() => onRecommendationClick(suggestion)}
                      className="bg-[#FFC09F] hover:bg-[#EF8354] text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
                      title={suggestion} // Optional: Show full suggestion on hover
                    >
                      {displayedText}
                    </button>
                  );
                })}
              </div>
            </>
          )}
        </>
      ) : (
        <p className="text-[#CEF7A0]">You have Ascended</p>
      )}
    </div>
  );
};

export default SuggestionPanel;
