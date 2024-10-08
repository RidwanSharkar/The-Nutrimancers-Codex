// src/components/SuggestionPanel.tsx

import React from 'react';

interface SuggestionPanelProps {
  missingNutrients: string[] | null;
  suggestions: string[] | null;
}

const SuggestionPanel: React.FC<SuggestionPanelProps> = ({ missingNutrients, suggestions }) => {
  return (
    <div className="bg-[#F48668] rounded-lg p-4 flex-1">
      <h2 className="text-xl font-semibold mb-4 text-white">Suggestions</h2>
      {missingNutrients && missingNutrients.length > 0 ? (
        <>
          <h3 className="text-lg font-medium mb-2 text-white">You're Missing:</h3>
          <ul className="list-disc list-inside mb-4 space-y-1">
            {missingNutrients.map((nutrient, index) => (
              <li key={index} className="text-white">
                {nutrient}
              </li>
            ))}
          </ul>
          {suggestions && suggestions.length > 0 && (
            <>
              <h3 className="text-lg font-medium mb-2 text-white">Next Meal Suggestions:</h3>
              <ul className="list-disc list-inside space-y-1">
                {suggestions.map((suggestion, index) => (
                  <li key={index} className="text-[#CEF7A0]">
                    {suggestion}
                  </li>
                ))}
              </ul>
            </>
          )}
        </>
      ) : (
        <p className="text-[#CEF7A0]">Great job! You're meeting all your essential nutrient needs.</p>
      )}
    </div>
  );
};

export default SuggestionPanel;
