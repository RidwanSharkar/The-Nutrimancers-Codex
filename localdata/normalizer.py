import pandas as pd

column_mapping = {
    "Potassium, K": "Potassium",
    "Sodium, Na": "Sodium",
    "Calcium, Ca": "Calcium",
    "Phosphorus, P": "Phosphorus",
    "Magnesium, Mg": "Magnesium",
    "Iron, Fe": "Iron",
    "Zinc, Zn": "Zinc",
    "Manganese, Mn": "Manganese",
    "Copper, Cu": "Copper",
    "Selenium, Se": "Selenium",

    "Histidine": "Histidine",
    "Isoleucine": "Isoleucine",
    "Leucine": "Leucine",
    "Lysine": "Lysine",
    "Methionine": "Methionine",
    "Phenylalanine": "Phenylalanine",
    "Threonine": "Threonine",
    "Tryptophan": "Tryptophan",
    "Valine": "Valine",

    "Vitamin A, RAE": "Vitamin A",
    "Thiamin": "Vitamin B1",
    "Riboflavin": "Vitamin B2",
    "Niacin": "Vitamin B3",
    "Pantothenic acid": "Vitamin B5",
    "Vitamin B-6": "Vitamin B6",
    "Folate, total": "Vitamin B9",
    "Vitamin B-12": "Vitamin B12",
    "Vitamin C, total ascorbic acid": "Vitamin C",
    "Vitamin D (D2 + D3)": "Vitamin D",
    "Vitamin E (alpha-tocopherol)": "Vitamin E",
    "Vitamin K (phylloquinone)": "Vitamin K",
    
    "Choline, total": "Choline"
}

df = pd.read_csv('dataset.csv')
df.rename(columns=column_mapping, inplace=True)

new_column_order = [
    "Potassium", 
    "Sodium", 
    "Calcium", 
    "Phosphorus", 
    "Magnesium", 
    "Iron", 
    "Zinc", 
    "Manganese",
    "Copper", 
    "Selenium", 

    "Histidine", 
    "Isoleucine", 
    "Leucine", 
    "Lysine", 
    "Methionine", 
    "Phenylalanine", 
    "Threonine", 
    "Tryptophan", 
    "Valine", 

    "Vitamin A", 
    "Vitamin B1", 
    "Vitamin B2", 
    "Vitamin B3", 
    "Vitamin B5", 
    "Vitamin B6", 
    "Vitamin B9", 
    "Vitamin B12", 
    "Vitamin C", 
    "Vitamin D", 
    "Vitamin E", 
    "Vitamin K", 

    "Choline"
]

df = df[new_column_order]
df.to_csv('ArrangedDataset.csv', index=False)
