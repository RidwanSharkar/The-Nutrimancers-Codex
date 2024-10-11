import pandas as pd
import os

nutrient_df = pd.read_csv(r'C:\Users\Lenovo\Desktop\Bioessence\data\Nutrient.csv')
print(nutrient_df[['id', 'name']])

Nutrients = [
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
    "Alpha-Linolenic Acid",
    "Linoleic Acid",
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
    "Choline",
]

nutrientMapping = {
    "Potassium": "Potassium, K",
    "Sodium": "Sodium, Na",
    "Calcium": "Calcium, Ca",
    "Phosphorus": "Phosphorus, P",
    "Magnesium": "Magnesium, Mg",
    "Iron": "Iron, Fe",
    "Zinc": "Zinc, Zn",
    "Manganese": "Manganese, Mn",
    "Copper": "Copper, Cu",
    "Selenium": "Selenium, Se",
    "Histidine": "Histidine",
    "Isoleucine": "Isoleucine",
    "Leucine": "Leucine",
    "Lysine": "Lysine",
    "Methionine": "Methionine",
    "Phenylalanine": "Phenylalanine",
    "Threonine": "Threonine",
    "Tryptophan": "Tryptophan",
    "Valine": "Valine",
    "Alpha-Linolenic Acid": "18:3 n-3 c,c,c (ALA)",
    "Linoleic Acid": "18:2 n-6 c,c",
    "Vitamin A": "Vitamin A, RAE",
    "Vitamin B1": "Thiamin",
    "Vitamin B2": "Riboflavin",
    "Vitamin B3": "Niacin",
    "Vitamin B5": "Pantothenic acid",
    "Vitamin B6": "Vitamin B-6",
    "Vitamin B9": "Folate, total",
    "Vitamin B12": "Vitamin B-12",
    "Vitamin C": "Vitamin C, total ascorbic acid",
    "Vitamin D": "Vitamin D (D2 + D3)",
    "Vitamin E": "Vitamin E (alpha-tocopherol)",
    "Vitamin K": "Vitamin K (phylloquinone)",
    "Choline": "Choline, total",
}

usda_nutrient_names = nutrient_df['name'].tolist()
missing_nutrients = [usda_name for usda_name in nutrientMapping.values() if usda_name not in usda_nutrient_names]

if missing_nutrients:
    print("These USDA nutrient names are not found in the dataset:")
    for nutrient in missing_nutrients:
        print(nutrient)
else:
    print("All nutrients are found in the dataset.")

#-----------------------------------------------------

# Load
data_dir = r'C:\Users\Lenovo\Desktop\Bioessence\data'
food_df = pd.read_csv(os.path.join(data_dir, 'Food.csv'))
nutrient_df = pd.read_csv(os.path.join(data_dir, 'Nutrient.csv'))
food_nutrient_df = pd.read_csv(os.path.join(data_dir, 'FoodNutrient.csv'))

#-----------------------------------------------------

# Get USDA names/IDs
usda_nutrient_names = list(nutrientMapping.values())
filtered_nutrient_df = nutrient_df[nutrient_df['name'].isin(usda_nutrient_names)]
nutrient_ids = filtered_nutrient_df['id'].tolist()

# Filter for our nutrients
filtered_food_nutrient_df = food_nutrient_df[food_nutrient_df['nutrient_id'].isin(nutrient_ids)]

# Merge all
merged_df = filtered_food_nutrient_df.merge(
    filtered_nutrient_df[['id', 'name']],
    left_on='nutrient_id',
    right_on='id',
    suffixes=('_fn', '_nutrient')
)
merged_df = merged_df.merge(food_df[['fdc_id', 'description']], on='fdc_id')

#-----------------------------------------------------

# Pivot to have nutrients as columns
pivot_df = merged_df.pivot_table(
    index=['fdc_id', 'description'],
    columns='name',
    values='amount',
    aggfunc='first'
)
pivot_df.reset_index(inplace=True)

#-----------------------------------------------------

output_file = 'filtered_nutrient_data.csv'
pivot_df.to_csv(output_file, index=False)
print(f"Filtered data saved to {output_file}")
