# 🪬 The Nutrimancer's Codex - Vol. II
*An AI-powered nutrition analysis system that transforms natural language food descriptions into comprehensive nutrient profiles, identifies dietary deficiencies, and recommends foods to balance your intake.*

## 📜 System Overview

The Nutrimancer's Codex is a full-stack application that combines natural language processing, real-time nutritional analysis, and machine learning to provide personalized dietary insights. Users input a meal description in plain English, and the system extracts ingredients, computes nutrient content against RDA standards, detects deficiencies, and recommends optimal foods using cosine similarity.

### Vol. II:<br>
![Vol  II](https://github.com/user-attachments/assets/23c0f1a1-51d3-4898-b564-c90495477d4b)

### Vol. I:<br>
![Vol  I](https://github.com/user-attachments/assets/af91009a-d7f3-4c40-94fc-d8ace8988c8d)

---

## 🏛️ Architecture &  Data Flow 🔄

### **1. User Input → Ingredient Extraction**
```
User Input: "I ate a chicken caesar salad with parmesan"
    ↓
Gemini Flash 1.5 LLM (NLP)
    ↓
Extracted: ["chicken", "romaine lettuce", "parmesan cheese", "caesar dressing"]
```
- **Frontend** (`App.tsx`) sends food description to backend
- **Backend** (`main.go`) forwards request to Gemini API
- **Gemini Service** (`geminiService.go`) uses prompt engineering to extract core ingredients
- System cleans output (removes bullets, formatting) → returns ingredient list

### **2. Nutrient Data Retrieval**
```
For each ingredient → Nutritionix API
    ↓
Returns: 40+ nutrients with serving quantities
    ↓
Mapped to 37 essential nutrients (minerals, vitamins, amino acids, fatty acids)
```
- **Nutritionix Service** (`nutritionixService.go`) queries each ingredient individually
- Converts API response using `nutrientMapping` (attr_id → nutrient name)
- Handles unit conversions (mg, µg, g, IU) to standardize values

### **3. RDA Percentage Calculation**
```
Raw Nutrient Amount ÷ Daily RDA × 100 = Nutrient Percentage
    ↓
Aggregated per ingredient
    ↓
Total nutrient profile capped at 100% per nutrient
```
- **Backend Logic** (`main.go`) calculates percentages using `nutrientRDA` map
- Adjusts units via `adjustUnits()` and `convertIUtoMg()` functions
- Combines all ingredient nutrients → produces total meal profile

### **4. Deficiency Detection**
```
For each of 37 tracked nutrients:
    if percentage < 3.5% → flagged as deficient
    ↓
Deficiency Vector: [0,1,0,1,1,0,0,1,...] (binary representation)
```
- **Threshold**: 3.5% of RDA marks deficiency
- Creates binary deficiency vector for ML algorithm input

### **5. Food Recommendation (Machine Learning)**
```
USDA Dataset (10,000+ foods) with pre-computed nutrient vectors
    ↓
Cosine Similarity Algorithm
    ↓
Similarity Score = (Food Vector · Deficiency Vector) / (||Food|| × ||Deficiency||)
    ↓
Top 5 foods ranked by similarity score
```
- **Data Loader** (`dataLoader.go`) loads USDA food dataset (`dataset.csv`) at server startup
- **Recommendation Engine** (`recommendTron.go`) compares deficiency vector against all foods
- **Cosine Similarity** (`cosineSimilarity.go`) measures vector alignment (0 to 1 scale)
- Deduplicates similar foods, returns top matches

### **6. Interactive Visualization**
```
Frontend receives:
    • Ingredients list
    • Per-ingredient nutrient breakdown
    • Total nutrient percentages (37 nutrients)
    • Deficiencies array
    • Top 5 food recommendations
    ↓
React Components render:
    • OrbsPanel: Animated nutrient gauges (GSAP)
    • IngredientsPanel: Clickable ingredient breakdown
    • SuggestionPanel: Interactive food recommendations
```
- **Frontend State Management** (`App.tsx`) handles real-time nutrient updates
- Users can click recommendations → system fetches new nutrient data → updates totals dynamically
- Highlights changed nutrients after adding recommended foods

---

## 🛠️ Tech Stack

### **Frontend**
- **React 18** (TypeScript) - Component-based UI
- **GSAP** - Fluid animations for nutrient orbs
- **Tailwind CSS** - Utility-first styling
- **Axios** - HTTP client for backend communication

### **Backend**
- **Go 1.22** - High-performance HTTP server
- **net/http** - Native HTTP routing (`/process-food`, `/fetch-nutrient-data`)
- **CORS** - Cross-origin middleware for frontend integration
- **godotenv** - Environment variable management

### **Data & APIs**
- **Gemini Flash 1.5** - Google's LLM for ingredient extraction
- **Nutritionix API** - Real-time nutrient data (40+ nutrients per food)
- **USDA FoodData Central** - 10,000+ food nutrient vectors for ML dataset

### **Machine Learning**
- **Cosine Similarity** - Custom Go implementation for vector comparison
- **Binary Vector Encoding** - Efficient deficiency representation
- **CSV Dataset** - Preprocessed nutrient matrix for fast lookups

### **Deployment**
- **AWS Elastic Beanstalk** - Backend hosting (EC2-based)
- **AWS Amplify** - Frontend hosting with CI/CD
- **Nginx** - Reverse proxy for production routing

---

## 🧬 Key Features

### **1. Multi-Source Nutrient Tracking (37 Nutrients)**
- **Minerals** (10): Potassium, Sodium, Calcium, Phosphorus, Magnesium, Iron, Zinc, Manganese, Copper, Selenium
- **Vitamins** (12): A, B1-B6, B9, B12, C, D, E, K
- **Essential Amino Acids** (9): Histidine, Isoleucine, Leucine, Lysine, Methionine, Phenylalanine, Threonine, Tryptophan, Valine
- **Essential Fatty Acids** (4): Omega-3 (ALA, EPA, DHA), Omega-6 (Linoleic Acid)
- **Choline** (1): Critical for brain function

### **2. Real-Time Nutrient Calculation**
- Converts all units to RDA-standardized percentages
- Accounts for serving sizes and food preparation methods
- Caps nutrients at 100% (prevents over-supplementation bias)

### **3. Intelligent Food Recommendations**
- Cosine similarity ranges from 0 (no similarity) to 1 (perfect match)
- Recommends foods that collectively address multiple deficiencies
- Deduplicates similar foods (e.g., "raw spinach" vs "cooked spinach")

### **4. Interactive Nutrient Exploration**
- Click individual ingredients → see their specific nutrient contributions
- Click recommendations → preview how they'd affect your totals
- Animated nutrient orbs grouped by category (Minerals, Vitamins, Amino Acids, Fatty Acids)

---

## 📊 Algorithm Choice: Cosine Similarity

### **Why Cosine Similarity?**
Unlike Euclidean distance, cosine similarity measures **direction** rather than magnitude, making it ideal for comparing nutrient profiles where relative composition matters more than absolute amounts.

### **Mathematical Formula**
```
similarity(A, B) = (A · B) / (||A|| × ||B||)

Where:
  A · B = Σ(Aᵢ × Bᵢ)           [dot product]
  ||A|| = √(Σ Aᵢ²)              [magnitude of A]
  ||B|| = √(Σ Bᵢ²)              [magnitude of B]
```

### **Example Calculation**
```
User Deficiencies: [Vitamin C=1, Iron=1, Calcium=0, Zinc=1]
Food (Spinach):    [Vitamin C=0.4, Iron=0.8, Calcium=0.6, Zinc=0.2]

Dot Product = (1×0.4) + (1×0.8) + (0×0.6) + (1×0.2) = 1.4
||Deficiency|| = √(1² + 1² + 0² + 1²) = √3 ≈ 1.732
||Food|| = √(0.4² + 0.8² + 0.6² + 0.2²) ≈ 1.077

Similarity = 1.4 / (1.732 × 1.077) ≈ 0.75 (high match!)
```

Foods with scores > 0.7 are excellent matches for addressing multiple deficiencies simultaneously.

---

## 📁 Project Structure

```
The-Nutrimancers-Codex/
│
├── frontend/
│   ├── src/
│   │   ├── App.tsx                    # Main app logic & state management
│   │   ├── grimoire/
│   │   │   ├── IngredientsPanel.tsx   # Ingredient list UI
│   │   │   ├── OrbsPanel.tsx          # Animated nutrient visualization
│   │   │   └── SuggestionPanel.tsx    # Food recommendations UI
│   │   └── services/
│   │       └── backendService.ts      # API client
│   └── ...
│
├── amplify/backend/
│   ├── main.go                        # HTTP server & request routing
│   ├── services/
│   │   ├── geminiService.go           # LLM ingredient extraction
│   │   └── nutritionixService.go      # Nutrient data fetching
│   ├── machinist/
│   │   ├── dataLoader.go              # USDA dataset loader
│   │   ├── recommendTron.go           # ML recommendation engine
│   │   ├── cosineSimilarity.go        # Similarity algorithm
│   │   └── dataset.csv                # 10K+ food nutrient vectors
│   └── models/
│       ├── food.go                    # Data structures
│       └── model.go
│
└── data/
    ├── dataprocessor.py               # Dataset preprocessing scripts
    ├── normalizer.py
    └── *.csv                          # Raw USDA data
```

---

## 🚀 Setup & Installation

### **Prerequisites**
```bash
Node.js 18+
Go 1.22+
Git
```

### **Backend Setup**
```bash
cd amplify/backend

# Install Go dependencies
go mod download

# Create .env file
cat > .env << EOF
API_KEY=your_gemini_api_key
NUTRITIONIX_APP_ID=your_nutritionix_app_id
NUTRITIONIX_APP_KEY=your_nutritionix_app_key
PORT=5000
EOF

# Run backend
go run main.go
```

### **Frontend Setup**
```bash
cd frontend

# Install dependencies
npm install

# Update backend URL in backendService.ts (if needed)
# Run development server
npm run dev
```

### **Access Application**
```
Frontend: http://localhost:5173
Backend:  http://localhost:5000
```

---

## ⚙️ Potential Future Enhancements
- User accounts with meal history tracking
- Weekly nutrition trend analysis
- Recipe generation from recommended foods
- Mobile app (React Native)
- Integration with fitness trackers (calories, macros)
- Support for dietary restrictions (vegan, keto, etc.)

---

## 📝 License
MIT License - Feel free to use this project for learning and development. 
