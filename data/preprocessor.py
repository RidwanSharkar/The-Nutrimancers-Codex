import pandas as pd
import os

# Load datasets
data_dir = r'C:\Users\Lenovo\Desktop\Bioessence\data'
food_df = pd.read_csv(os.path.join(data_dir, 'FoundationalFood.csv'))
nutrient_df = pd.read_csv(os.path.join(data_dir, 'FoundationalNutrient.csv'))
food_nutrient_df = pd.read_csv(os.path.join(data_dir, 'FoundationalFoodNutrient.csv'), low_memory=False)

# Inspect
print(food_nutrient_df.dtypes)
print(food_nutrient_df.iloc[:, 9].unique())


nutrientMapping = {
    # Essential Ions
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

    # Essential Amino Acids
    "Histidine": "Histidine",
    "Isoleucine": "Isoleucine",
    "Leucine": "Leucine",
    "Lysine": "Lysine",
    "Methionine": "Methionine",
    "Phenylalanine": "Phenylalanine",
    "Threonine": "Threonine",
    "Tryptophan": "Tryptophan",
    "Valine": "Valine",

    # Essential Omega Fatty Acids
    "Alpha-Linolenic Acid": "PUFA 18:3 n-3 c,c,c (ALA)",
    "Linoleic Acid": "PUFA 18:2 n-6 c,c",
    "EPA": "PUFA 20:5 n-3 (EPA)",
    "DHA": "PUFA 22:6 n-3 (DHA)",

    # Vitamins
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

    "Choline": "Choline, total"
}

# Filter
usda_nutrient_names = list(nutrientMapping.values())
filtered_nutrient_df = nutrient_df[nutrient_df['name'].isin(usda_nutrient_names)]
nutrient_ids = filtered_nutrient_df['id'].tolist()
filtered_food_nutrient_df = food_nutrient_df[food_nutrient_df['nutrient_id'].isin(nutrient_ids)]

# Merge
merged_df = filtered_food_nutrient_df.merge(
    filtered_nutrient_df[['id', 'name']],
    left_on='nutrient_id',
    right_on='id',
    suffixes=('_fn', '_nutrient')
)
merged_df = merged_df.merge(food_df[['fdc_id', 'description']], on='fdc_id')

# Pivot
pivot_df = merged_df.pivot_table(
    index=['fdc_id', 'description'],
    columns='name',
    values='amount',
    aggfunc='first'
)
pivot_df.reset_index(inplace=True)

# Rename
reverse_mapping = {v: k for k, v in nutrientMapping.items()}
pivot_df.rename(columns=reverse_mapping, inplace=True)

# Reorder
final_columns = [
    "fdc_id", "description", "Potassium", "Sodium", "Calcium", "Phosphorus", "Magnesium",
    "Iron", "Zinc", "Manganese", "Copper", "Selenium", "Histidine", "Isoleucine", 
    "Leucine", "Lysine", "Methionine", "Phenylalanine", "Threonine", "Tryptophan", "Valine", 
    "Alpha-Linolenic Acid", "Linoleic Acid", "EPA", "DHA", "Vitamin A", "Vitamin B1", 
    "Vitamin B2", "Vitamin B3", "Vitamin B5", "Vitamin B6", "Vitamin B9", "Vitamin B12", 
    "Vitamin C", "Vitamin D", "Vitamin E", "Vitamin K", "Choline"
]

# Fill
pivot_df = pivot_df.reindex(columns=final_columns)

# Output
output_file = 'dataset.csv'
pivot_df.to_csv(output_file, index=False)

