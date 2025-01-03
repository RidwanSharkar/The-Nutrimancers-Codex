# The Nutrimancer's Codex - Vol. II
An AI-powered application that analyzes food descriptions to extract ingredients, compute nutrient deficiencies, and recommend foods to balance your diet.

## Overview:
The Nutrimancer's Codex is a full-stack application that leverages AI and machine learning to help users understand their nutrient intake and make informed dietary choices. By inputting a food description, users receive an analysis of nutrient content, identify deficiencies, and get personalized food recommendations.

---

## Features:
• Ingredient Extraction: utilizes the Gemini Language Model to parse natural language food descriptions.<br>
• Nutrient Analysis: calculates nutrient percentages based on recommended daily allowances using data from Nutritionix and USDA. <br>
• Deficiency Detection: identifies low or missing essential nutrients in the user's diet. <br>
• Recommendation: cosine similarity algorithm is applied across dataset to display the foods most capable of alleviating the current active deficiencies. 

---

## Vol. II:<br>
![Vol  II](https://github.com/user-attachments/assets/23c0f1a1-51d3-4898-b564-c90495477d4b)

---

## Vol. I:<br>
![Vol  I](https://github.com/user-attachments/assets/af91009a-d7f3-4c40-94fc-d8ace8988c8d)

---

## Tech Stack:

• **Frontend:** React (TypeScript), GSAP, Tailwind (CSS), Axios

• **Backend:** Go (GoLang), net/http, CORS, JSON processing 

• **Dataset & APIs:** USDA FoodData Central, Google Generative Language, Nutritionix

• **Natural Language Processing:** Gemini Flash 1.5

• **Machine Learning:** Cosine Similarity (GoLang)

• **Deployment:** AWS Amplify, AWS Elastic Beanstack via EC2, Nginx 
