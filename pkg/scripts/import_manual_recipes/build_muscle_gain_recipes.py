from __future__ import annotations

import json
from pathlib import Path
from urllib.parse import quote


OUTPUT_PATH = Path(__file__).with_name("muscle_gain_recipes_100.json")


NUTRITION = {
    "鸡胸肉": {"energy": 165, "protein": 31.0, "carbohydrate": 0.0, "fat": 3.6},
    "鸡腿肉(去皮)": {"energy": 175, "protein": 24.0, "carbohydrate": 0.0, "fat": 8.0},
    "牛里脊": {"energy": 155, "protein": 22.0, "carbohydrate": 0.0, "fat": 7.0},
    "牛腩": {"energy": 190, "protein": 20.0, "carbohydrate": 0.0, "fat": 12.0},
    "瘦猪里脊": {"energy": 143, "protein": 21.0, "carbohydrate": 0.0, "fat": 6.0},
    "猪肉末": {"energy": 170, "protein": 20.0, "carbohydrate": 0.0, "fat": 10.0},
    "虾仁": {"energy": 99, "protein": 24.0, "carbohydrate": 0.2, "fat": 0.3},
    "鲈鱼": {"energy": 120, "protein": 20.0, "carbohydrate": 0.0, "fat": 4.0},
    "鳕鱼": {"energy": 105, "protein": 23.0, "carbohydrate": 0.0, "fat": 1.0},
    "鱼片": {"energy": 110, "protein": 21.0, "carbohydrate": 0.0, "fat": 2.0},
    "金枪鱼罐头": {"energy": 116, "protein": 25.0, "carbohydrate": 0.0, "fat": 1.0},
    "鸡蛋": {"energy": 144, "protein": 12.6, "carbohydrate": 1.1, "fat": 9.9},
    "北豆腐": {"energy": 81, "protein": 8.5, "carbohydrate": 3.0, "fat": 4.8},
    "豆干": {"energy": 140, "protein": 18.0, "carbohydrate": 4.0, "fat": 7.0},
    "香干": {"energy": 150, "protein": 16.0, "carbohydrate": 4.0, "fat": 8.0},
    "毛豆": {"energy": 131, "protein": 13.0, "carbohydrate": 11.0, "fat": 6.0},
    "无糖酸奶": {"energy": 60, "protein": 4.0, "carbohydrate": 5.0, "fat": 3.0},
    "无糖豆浆": {"energy": 33, "protein": 3.0, "carbohydrate": 1.8, "fat": 1.6},
    "燕麦": {"energy": 376, "protein": 13.0, "carbohydrate": 67.0, "fat": 7.0},
    "面粉": {"energy": 344, "protein": 10.0, "carbohydrate": 73.0, "fat": 1.5},
    "全麦卷饼": {"energy": 290, "protein": 9.0, "carbohydrate": 50.0, "fat": 6.0},
    "荞麦面(干)": {"energy": 335, "protein": 12.0, "carbohydrate": 65.0, "fat": 2.0},
    "挂面(干)": {"energy": 340, "protein": 10.0, "carbohydrate": 70.0, "fat": 1.0},
    "米粉(干)": {"energy": 356, "protein": 7.0, "carbohydrate": 79.0, "fat": 1.0},
    "馄饨皮": {"energy": 280, "protein": 8.0, "carbohydrate": 58.0, "fat": 2.0},
    "烧麦皮": {"energy": 285, "protein": 8.5, "carbohydrate": 59.0, "fat": 2.0},
    "白米饭": {"energy": 116, "protein": 2.6, "carbohydrate": 25.7, "fat": 0.3},
    "小米(干)": {"energy": 361, "protein": 9.0, "carbohydrate": 75.0, "fat": 3.0},
    "玉米": {"energy": 112, "protein": 3.2, "carbohydrate": 22.8, "fat": 1.2},
    "南瓜": {"energy": 45, "protein": 1.0, "carbohydrate": 10.0, "fat": 0.1},
    "红薯": {"energy": 86, "protein": 1.3, "carbohydrate": 20.1, "fat": 0.1},
    "紫薯": {"energy": 82, "protein": 1.5, "carbohydrate": 19.0, "fat": 0.1},
    "山药": {"energy": 57, "protein": 1.5, "carbohydrate": 13.0, "fat": 0.1},
    "香蕉": {"energy": 89, "protein": 1.1, "carbohydrate": 23.0, "fat": 0.3},
    "苹果": {"energy": 52, "protein": 0.3, "carbohydrate": 14.0, "fat": 0.2},
    "蓝莓": {"energy": 57, "protein": 0.7, "carbohydrate": 14.5, "fat": 0.3},
    "草莓": {"energy": 32, "protein": 0.8, "carbohydrate": 7.1, "fat": 0.3},
    "菠菜": {"energy": 23, "protein": 2.9, "carbohydrate": 3.6, "fat": 0.4},
    "香菇": {"energy": 34, "protein": 2.2, "carbohydrate": 6.9, "fat": 0.5},
    "番茄": {"energy": 18, "protein": 0.9, "carbohydrate": 3.9, "fat": 0.2},
    "黄瓜": {"energy": 16, "protein": 0.7, "carbohydrate": 3.6, "fat": 0.1},
    "洋葱": {"energy": 40, "protein": 1.1, "carbohydrate": 9.0, "fat": 0.1},
    "青椒": {"energy": 22, "protein": 1.0, "carbohydrate": 5.0, "fat": 0.2},
    "西兰花": {"energy": 34, "protein": 2.8, "carbohydrate": 7.0, "fat": 0.4},
    "木耳": {"energy": 25, "protein": 1.5, "carbohydrate": 6.0, "fat": 0.2},
    "芹菜": {"energy": 16, "protein": 0.8, "carbohydrate": 3.5, "fat": 0.2},
    "冬瓜": {"energy": 13, "protein": 0.4, "carbohydrate": 2.6, "fat": 0.2},
    "胡萝卜": {"energy": 41, "protein": 0.9, "carbohydrate": 10.0, "fat": 0.2},
    "西葫芦": {"energy": 19, "protein": 1.2, "carbohydrate": 3.8, "fat": 0.2},
    "芦笋": {"energy": 20, "protein": 2.2, "carbohydrate": 3.7, "fat": 0.1},
    "丝瓜": {"energy": 20, "protein": 1.0, "carbohydrate": 4.0, "fat": 0.2},
    "生菜": {"energy": 15, "protein": 1.4, "carbohydrate": 2.9, "fat": 0.2},
    "海苔": {"energy": 270, "protein": 35.0, "carbohydrate": 30.0, "fat": 2.0},
    "食用油": {"energy": 900, "protein": 0.0, "carbohydrate": 0.0, "fat": 100.0},
    "清水": {"energy": 0, "protein": 0.0, "carbohydrate": 0.0, "fat": 0.0},
}


COMMONS_FILE_NAMES = {
    "pancake": "Chive egg pancake.jpg",
    "chicken_noodle": "Chicken Noodle Soup.jpg",
    "congee": "Congee.JPG",
    "onigiri": "Japanese rice balls (onigiri).jpg",
    "steamed_egg": "Chawanmushi.JPG",
    "meat_pie": "Curry beef pie (cropped).jpg",
    "buns": "Steamed Buns (4650154698).jpg",
    "wrap": "Egg Tortilla Wrap (8728191788).jpg",
    "oatmeal": "Cooked oatmeal in bowl.jpg",
    "shrimp_soup": "Shrimp soup at home.jpg",
    "wonton": "Wonton Soup.jpg",
    "shumai": "Shanxi Shumai.jpg",
    "yogurt": "Yogurt fruit bowl.jpg",
    "beef_stir": "Beef with broccoli.JPG",
    "chicken_stir": "Chicken stir fry (1126047571).jpg",
    "shrimp_stir": "Shrimp and broccoli stir fry (8374790146).jpg",
    "steamed_fish": "Steamed rice with fish.jpg",
    "pork_stir": "Teriyaki Pork Stir Fry.jpg",
    "braised_tofu": "Dubu Jorim (Braised tofu).jpg",
    "tofu_stew": "Tofu stew, pork.jpg",
    "winter_melon_soup": "Winter melon soup.jpg",
    "fish_soup": "Fish soup.jpg",
    "tofu_soup": "Spinach-tofu soup.jpg",
    "boiled_eggs": "Boiled Eggs.jpg",
    "soy_milk": "Soy milk.jpg",
    "steamed_corn": "Steamed Corn.jpg",
    "sweet_potato": "Steamed sweet potato.jpg",
    "edamame": "Bowl of Edamame.jpg",
}


IMAGE_KEYS = [
    "pancake",
    "chicken_noodle",
    "congee",
    "onigiri",
    "steamed_egg",
    "meat_pie",
    "pancake",
    "buns",
    "chicken_noodle",
    "congee",
    "oatmeal",
    "buns",
    "pancake",
    "shrimp_soup",
    "wrap",
    "steamed_egg",
    "pancake",
    "chicken_noodle",
    "wonton",
    "shumai",
    "pancake",
    "onigiri",
    "yogurt",
    "pancake",
    "onigiri",
    "beef_stir",
    "chicken_stir",
    "shrimp_stir",
    "steamed_fish",
    "chicken_stir",
    "pork_stir",
    "beef_stir",
    "braised_tofu",
    "chicken_stir",
    "tofu_stew",
    "beef_stir",
    "pork_stir",
    "chicken_stir",
    "winter_melon_soup",
    "fish_soup",
    "chicken_stir",
    "beef_stir",
    "chicken_stir",
    "beef_stir",
    "braised_tofu",
    "steamed_egg",
    "beef_stir",
    "chicken_stir",
    "braised_tofu",
    "pork_stir",
    "steamed_fish",
    "chicken_stir",
    "shrimp_soup",
    "chicken_stir",
    "beef_stir",
    "steamed_egg",
    "winter_melon_soup",
    "shrimp_stir",
    "chicken_stir",
    "fish_soup",
    "beef_stir",
    "tofu_stew",
    "pork_stir",
    "beef_stir",
    "chicken_stir",
    "tofu_stew",
    "chicken_stir",
    "fish_soup",
    "shrimp_stir",
    "beef_stir",
    "winter_melon_soup",
    "steamed_egg",
    "beef_stir",
    "shrimp_stir",
    "tofu_soup",
    "boiled_eggs",
    "boiled_eggs",
    "yogurt",
    "soy_milk",
    "steamed_corn",
    "sweet_potato",
    "sweet_potato",
    "onigiri",
    "onigiri",
    "yogurt",
    "yogurt",
    "yogurt",
    "yogurt",
    "edamame",
    "braised_tofu",
    "steamed_corn",
    "yogurt",
    "sweet_potato",
    "steamed_egg",
    "wrap",
    "oatmeal",
    "wrap",
    "yogurt",
    "edamame",
    "oatmeal",
]


def ing(name: str, grams: float) -> dict:
    return {"name": name, "grams": grams}


def format_grams(value: float) -> str:
    if int(value) == value:
        return str(int(value))
    return f"{value:.1f}".rstrip("0").rstrip(".")


def ingredient_line(item: dict) -> str:
    return f"{item['name']} {format_grams(item['grams'])}g"


def nutrition_totals(ingredients: list[dict]) -> tuple[float, float, float, float, float]:
    total_weight = energy = protein = carbohydrate = fat = 0.0
    for item in ingredients:
        grams = item["grams"]
        total_weight += grams
        base = NUTRITION[item["name"]]
        energy += base["energy"] * grams / 100.0
        protein += base["protein"] * grams / 100.0
        carbohydrate += base["carbohydrate"] * grams / 100.0
        fat += base["fat"] * grams / 100.0
    return (
        round(total_weight, 1),
        round(energy, 1),
        round(protein, 1),
        round(carbohydrate, 1),
        round(fat, 1),
    )


def forbidden_users(ingredients: list[dict]) -> list[str]:
    names = {item["name"] for item in ingredients}
    tags: list[str] = []
    if names & {"虾仁", "鲈鱼", "鳕鱼", "鱼片", "金枪鱼罐头"}:
        tags.append("海鲜过敏")
    if names & {"无糖酸奶"}:
        tags.append("乳制品过敏")
    if names & {"北豆腐", "豆干", "香干", "毛豆", "无糖豆浆"}:
        tags.append("大豆过敏")
    if names & {"面粉", "全麦卷饼", "荞麦面(干)", "挂面(干)", "米粉(干)", "馄饨皮", "烧麦皮"}:
        tags.append("麸质过敏")
    if "鸡蛋" in names:
        tags.append("蛋类过敏")
    if names & {"牛里脊", "牛腩"}:
        tags.append("牛肉忌口")
    if names & {"瘦猪里脊", "猪肉末"}:
        tags.append("猪肉忌口")
    return tags


def build_recipe(name: str, meal_type: str, allowed_meal_types: list[str], difficulty: str, cooking_time: int, ingredients: list[dict], cooking_steps: list[str]) -> dict:
    portion_weight_g, energy, protein, carbohydrate, fat = nutrition_totals(ingredients)
    return {
        "name": name,
        "image_url": "",
        "meal_type": meal_type,
        "allowed_meal_types": allowed_meal_types,
        "difficulty": difficulty,
        "cooking_time": cooking_time,
        "portion_weight_g": portion_weight_g,
        "energy": energy,
        "protein": protein,
        "carbohydrate": carbohydrate,
        "fat": fat,
        "ingredients": [ingredient_line(item) for item in ingredients],
        "cooking_steps": cooking_steps,
        "target_users": ["增肌", "大众"],
        "forbidden_users": forbidden_users(ingredients),
        "is_weight_loss_friendly": False,
        "is_muscle_gain_friendly": True,
        "is_sugar_control_friendly": False,
        "is_general_friendly": True,
    }


def commons_file_url(file_name: str) -> str:
    return "https://commons.wikimedia.org/wiki/Special:FilePath/" + quote(file_name, safe="")


def apply_image_urls(recipes: list[dict]) -> None:
    assert len(recipes) == len(IMAGE_KEYS), (len(recipes), len(IMAGE_KEYS))
    for recipe, key in zip(recipes, IMAGE_KEYS):
        recipe["image_url"] = commons_file_url(COMMONS_FILE_NAMES[key])


def pancake(name: str, protein_items: list[dict], veg_items: list[dict], *, flour_g: float = 55, egg_g: float = 55, meal_type: str = "早餐", time: int = 15) -> dict:
    ingredients = [ing("面粉", flour_g)]
    if egg_g > 0:
        ingredients.append(ing("鸡蛋", egg_g))
    ingredients.extend(protein_items)
    ingredients.extend(veg_items)
    ingredients.append(ing("食用油", 4))
    steps = [
        "把主料和蔬菜切碎，和面粉、鸡蛋一起拌成面糊。",
        "平底锅刷少量油，倒入面糊摊平。",
        "小火煎到两面熟透后出锅切块。",
    ]
    allowed = ["早餐"] if meal_type == "早餐" else ["加餐", "早餐"]
    return build_recipe(name, meal_type, allowed, "简单", time, ingredients, steps)


def steamed_cake(name: str, ingredients: list[dict], *, meal_type: str = "早餐", time: int = 25) -> dict:
    steps = [
        "把粉类和液体食材拌成细腻面糊，再加入主料拌匀。",
        "面糊倒入模具或耐热碗中，静置片刻。",
        "上锅蒸到完全定型熟透，放温后切块食用。",
    ]
    return build_recipe(name, meal_type, ["早餐"], "中等", time, ingredients, steps)


def steamed_bun(name: str, dough_items: list[dict], filling_items: list[dict], *, meal_type: str = "早餐", time: int = 30) -> dict:
    ingredients = [*dough_items, *filling_items]
    steps = [
        "把面团材料揉匀后静置，馅料切碎拌好备用。",
        "将面团分剂子包入馅料，收口后略微醒发。",
        "上锅蒸到外皮熟透松软即可。",
    ]
    return build_recipe(name, meal_type, ["早餐"], "中等", time, ingredients, steps)


def noodle(name: str, protein_items: list[dict], veg_items: list[dict], noodle_item: dict, *, meal_type: str, time: int = 16) -> dict:
    ingredients = [noodle_item, *protein_items, *veg_items, ing("清水", 220), ing("食用油", 3)]
    steps = [
        "锅中加水烧开，先把面或米粉煮至八成熟。",
        "另起锅把主料和蔬菜煮熟或炒香，再倒回面中。",
        "调到咸淡合适，煮透后即可出锅。",
    ]
    allowed = ["早餐"] if meal_type == "早餐" else ["午餐", "晚餐"]
    return build_recipe(name, meal_type, allowed, "简单", time, ingredients, steps)


def porridge(name: str, base_items: list[dict], extra_items: list[dict], *, meal_type: str = "早餐", time: int = 22) -> dict:
    ingredients = [*base_items, *extra_items, ing("清水", 260)]
    allowed = ["早餐", "加餐"] if meal_type == "早餐" else ["晚餐"]
    steps = [
        "将谷物或主食洗净后入锅，加清水煮开。",
        "转小火煮到软烂，再加入其他食材继续煮熟。",
        "收至喜欢的稠度后盛出即可。",
    ]
    return build_recipe(name, meal_type, allowed, "简单", time, ingredients, steps)


def rice_ball(name: str, fillings: list[dict], *, meal_type: str = "早餐", time: int = 12) -> dict:
    ingredients = [ing("白米饭", 160), *fillings, ing("海苔", 5)]
    allowed = ["早餐", "加餐"] if meal_type == "早餐" else ["加餐"]
    steps = [
        "把主料提前做熟切碎，和米饭拌匀。",
        "加入海苔碎后分成两份。",
        "戴手套捏成饭团，压实后即可食用。",
    ]
    return build_recipe(name, meal_type, allowed, "简单", time, ingredients, steps)


def steamed_egg(name: str, extra_items: list[dict], *, meal_type: str, time: int = 18) -> dict:
    ingredients = [ing("鸡蛋", 110), *extra_items, ing("清水", 120)]
    allowed = ["早餐", "加餐"] if meal_type in {"早餐", "加餐"} else ["午餐", "晚餐"]
    steps = [
        "鸡蛋打散后加入温水，搅匀过筛。",
        "放入配料后盖上保鲜膜或盘子。",
        "上锅蒸到蛋液凝固熟透即可。",
    ]
    return build_recipe(name, meal_type, allowed, "简单", time, ingredients, steps)


def dough_wrap(name: str, fillings: list[dict], *, wrapper: str = "全麦卷饼", meal_type: str = "早餐", time: int = 14) -> dict:
    ingredients = [ing(wrapper, 85), *fillings, ing("食用油", 3)]
    allowed = ["早餐"] if meal_type == "早餐" else ["加餐"]
    steps = [
        "把夹馅食材提前处理熟透，切成小块或细丝。",
        "将饼皮加热到柔软后铺开，放入馅料。",
        "卷紧后切段，趁热食用口感更好。",
    ]
    return build_recipe(name, meal_type, allowed, "简单", time, ingredients, steps)


def dumpling_like(name: str, wrapper_name: str, fillings: list[dict], *, meal_type: str = "早餐", time: int = 22) -> dict:
    wrapper_g = 90 if wrapper_name == "馄饨皮" else 95
    ingredients = [ing(wrapper_name, wrapper_g), *fillings, ing("清水", 120)]
    steps = [
        "把馅料切碎后拌匀，包入皮中。",
        "锅中水开后下锅煮到浮起，再煮2到3分钟。",
        "捞出后即可直接食用或配少量汤汁。",
    ]
    return build_recipe(name, meal_type, ["早餐"], "中等", time, ingredients, steps)


def stir_fry(name: str, protein_items: list[dict], veg_items: list[dict], *, meal_type: str, time: int = 15) -> dict:
    ingredients = [*protein_items, *veg_items, ing("食用油", 5)]
    steps = [
        "把主料切片或切丝，蔬菜洗净切好。",
        "热锅少油先下肉类或海鲜炒到变色，再下蔬菜翻炒。",
        "炒到断生入味即可出锅。",
    ]
    allowed = ["午餐", "晚餐"] if meal_type in {"午餐", "晚餐"} else ["早餐"]
    return build_recipe(name, meal_type, allowed, "简单", time, ingredients, steps)


def steamed(name: str, main_items: list[dict], veg_items: list[dict], *, meal_type: str, time: int = 18) -> dict:
    ingredients = [*main_items, *veg_items]
    steps = [
        "把主料和配料处理干净后装盘。",
        "水开后上锅蒸到完全熟透。",
        "出锅后根据口味少量调味即可。",
    ]
    return build_recipe(name, meal_type, ["午餐", "晚餐"], "简单", time, ingredients, steps)


def soup(name: str, main_items: list[dict], veg_items: list[dict], *, meal_type: str, time: int = 20) -> dict:
    ingredients = [*main_items, *veg_items, ing("清水", 320)]
    steps = [
        "锅中加水烧开，先下比较耐煮的食材。",
        "再放入主料和易熟配料，小火煮到食材成熟。",
        "最后调味，保持汤体清爽即可。",
    ]
    return build_recipe(name, meal_type, ["午餐", "晚餐"], "简单", time, ingredients, steps)


def stew(name: str, main_items: list[dict], veg_items: list[dict], *, meal_type: str, time: int = 22) -> dict:
    ingredients = [*main_items, *veg_items, ing("清水", 180), ing("食用油", 4)]
    steps = [
        "先把主料煎香或焯水，蔬菜切块备用。",
        "所有食材入锅后加少量清水焖煮到入味。",
        "收汁或保留少量汤汁后即可出锅。",
    ]
    return build_recipe(name, meal_type, ["午餐", "晚餐"], "中等", time, ingredients, steps)


def yogurt_cup(name: str, items: list[dict], *, meal_type: str = "加餐", time: int = 4) -> dict:
    steps = [
        "把需要蒸熟的食材提前处理好并放凉。",
        "将食材按层次放入杯中。",
        "吃前轻轻拌匀即可。",
    ]
    return build_recipe(name, meal_type, ["加餐", "早餐"], "简单", time, items, steps)


def simple_snack(name: str, items: list[dict], *, time: int = 3) -> dict:
    steps = [
        "把食材准备好或加热到适口状态。",
        "装盘后即可直接食用。",
        "适合作为训练前后的小份加餐。",
    ]
    return build_recipe(name, "加餐", ["加餐"], "简单", time, items, steps)


def breakfast_recipes() -> list[dict]:
    return [
        pancake("蔬菜鸡蛋饼", [], [ing("菠菜", 35), ing("胡萝卜", 30)], flour_g=58),
        noodle("鸡丝荞麦面", [ing("鸡胸肉", 85)], [ing("黄瓜", 50), ing("香菇", 40)], ing("荞麦面(干)", 70), meal_type="早餐"),
        porridge("南瓜小米粥", [ing("小米(干)", 48)], [ing("南瓜", 120)], meal_type="早餐"),
        rice_ball("金枪鱼饭团", [ing("金枪鱼罐头", 75), ing("黄瓜", 35)]),
        steamed_egg("虾仁蒸蛋", [ing("虾仁", 85)], meal_type="早餐"),
        pancake("牛肉馅饼", [ing("牛里脊", 90)], [ing("洋葱", 45)], flour_g=65, egg_g=0),
        pancake("菠菜鸡蛋卷", [], [ing("菠菜", 70)], flour_g=42, egg_g=100),
        steamed_cake("玉米发糕", [ing("面粉", 62), ing("玉米", 100), ing("无糖豆浆", 90), ing("鸡蛋", 45)], meal_type="早餐", time=25),
        noodle("家常番茄鸡蛋面", [ing("鸡蛋", 55)], [ing("番茄", 120)], ing("挂面(干)", 72), meal_type="早餐"),
        porridge("鸡肉蔬菜粥", [ing("白米饭", 140)], [ing("鸡胸肉", 80), ing("胡萝卜", 35), ing("菠菜", 40)], meal_type="早餐"),
        porridge("紫薯豆浆燕麦粥", [ing("燕麦", 42)], [ing("紫薯", 110), ing("无糖豆浆", 180)], meal_type="早餐"),
        steamed_bun("香菇鸡肉包", [ing("面粉", 95), ing("清水", 55)], [ing("鸡胸肉", 85), ing("香菇", 45)], meal_type="早餐", time=28),
        pancake("山药鸡蛋饼", [], [ing("山药", 100)], flour_g=40, egg_g=100),
        soup("虾仁豆腐羹", [ing("虾仁", 90), ing("北豆腐", 120)], [ing("番茄", 70)], meal_type="早餐", time=16),
        dough_wrap("葱香牛肉卷饼", [ing("牛里脊", 85), ing("洋葱", 40), ing("青椒", 35)], meal_type="早餐"),
        steamed_egg("南瓜鸡蛋羹", [ing("南瓜", 110)], meal_type="早餐"),
        pancake("鸡蛋灌饼", [ing("鸡胸肉", 70)], [ing("生菜", 45)], flour_g=62, egg_g=55, meal_type="早餐"),
        noodle("鸡丝米粉", [ing("鸡胸肉", 85)], [ing("菠菜", 50), ing("香菇", 40)], ing("米粉(干)", 68), meal_type="早餐"),
        dumpling_like("虾仁馄饨", "馄饨皮", [ing("虾仁", 90), ing("菠菜", 35)], meal_type="早餐"),
        dumpling_like("牛肉烧麦", "烧麦皮", [ing("牛里脊", 95), ing("香菇", 35)], meal_type="早餐"),
        pancake("豆腐蔬菜饼", [ing("北豆腐", 100)], [ing("胡萝卜", 35), ing("菠菜", 35)], flour_g=38, egg_g=55),
        rice_ball("海苔鸡蛋饭团", [ing("鸡蛋", 60), ing("黄瓜", 30)]),
        yogurt_cup("蓝莓燕麦酸奶杯", [ing("无糖酸奶", 180), ing("燕麦", 28), ing("蓝莓", 60)], meal_type="早餐"),
        pancake("玉米鸡蛋饼", [], [ing("玉米", 85)], flour_g=46, egg_g=90),
        rice_ball("鸡肉饭团", [ing("鸡胸肉", 85), ing("玉米", 40)]),
    ]


def lunch_recipes() -> list[dict]:
    return [
        stew("番茄牛肉", [ing("牛里脊", 145)], [ing("番茄", 140), ing("洋葱", 40)], meal_type="午餐"),
        stir_fry("鲜香香菇滑鸡", [ing("鸡腿肉(去皮)", 145)], [ing("香菇", 80)], meal_type="午餐"),
        stir_fry("西兰花炒虾仁", [ing("虾仁", 150)], [ing("西兰花", 130)], meal_type="午餐"),
        steamed("葱丝清蒸鲈鱼", [ing("鲈鱼", 180)], [ing("番茄", 40)], meal_type="午餐"),
        stir_fry("家常青椒鸡丁", [ing("鸡胸肉", 145)], [ing("青椒", 90), ing("洋葱", 40)], meal_type="午餐"),
        stir_fry("木耳炒里脊", [ing("瘦猪里脊", 145)], [ing("木耳", 70), ing("青椒", 45)], meal_type="午餐"),
        stir_fry("小炒芹菜牛肉丝", [ing("牛里脊", 140)], [ing("芹菜", 110), ing("胡萝卜", 35)], meal_type="午餐"),
        stew("番茄豆腐", [ing("北豆腐", 180)], [ing("番茄", 150), ing("香菇", 50)], meal_type="午餐"),
        stir_fry("黄瓜鸡丝", [ing("鸡胸肉", 140)], [ing("黄瓜", 120), ing("胡萝卜", 30)], meal_type="午餐"),
        stew("虾仁豆腐煲", [ing("虾仁", 130), ing("北豆腐", 150)], [ing("番茄", 90)], meal_type="午餐"),
        stir_fry("洋葱炒牛肉", [ing("牛里脊", 145)], [ing("洋葱", 120), ing("青椒", 45)], meal_type="午餐"),
        stir_fry("香干炒肉丝", [ing("瘦猪里脊", 125), ing("香干", 80)], [ing("青椒", 60)], meal_type="午餐"),
        stir_fry("毛豆鸡丁", [ing("鸡胸肉", 135), ing("毛豆", 70)], [ing("胡萝卜", 35)], meal_type="午餐"),
        soup("清鲜冬瓜虾仁", [ing("虾仁", 140)], [ing("冬瓜", 180), ing("香菇", 35)], meal_type="午餐"),
        stew("番茄鱼片", [ing("鱼片", 150)], [ing("番茄", 150), ing("北豆腐", 80)], meal_type="午餐"),
        steamed("香菇蒸鸡", [ing("鸡腿肉(去皮)", 150)], [ing("香菇", 90)], meal_type="午餐", time=20),
        stir_fry("西葫芦炒牛肉", [ing("牛里脊", 140)], [ing("西葫芦", 120), ing("胡萝卜", 35)], meal_type="午餐"),
        stir_fry("芦笋鸡丁", [ing("鸡胸肉", 140)], [ing("芦笋", 120), ing("胡萝卜", 30)], meal_type="午餐"),
        stir_fry("青椒牛柳", [ing("牛里脊", 145)], [ing("青椒", 100), ing("洋葱", 40)], meal_type="午餐"),
        steamed("嫩豆腐蒸虾仁", [ing("虾仁", 130), ing("北豆腐", 140)], [ing("番茄", 50)], meal_type="午餐", time=20),
        steamed_egg("肉末蒸蛋", [ing("猪肉末", 80)], meal_type="午餐"),
        stew("番茄牛腩", [ing("牛腩", 150)], [ing("番茄", 150), ing("洋葱", 35)], meal_type="午餐", time=28),
        stir_fry("胡萝卜炒鸡丝", [ing("鸡胸肉", 140)], [ing("胡萝卜", 100), ing("青椒", 35)], meal_type="午餐"),
        stew("香菇豆腐", [ing("北豆腐", 180)], [ing("香菇", 100), ing("青椒", 40)], meal_type="午餐"),
        stir_fry("洋葱里脊丝", [ing("瘦猪里脊", 145)], [ing("洋葱", 120), ing("青椒", 40)], meal_type="午餐"),
    ]


def dinner_recipes() -> list[dict]:
    return [
        steamed("清蒸鳕鱼", [ing("鳕鱼", 180)], [ing("番茄", 45)], meal_type="晚餐", time=18),
        stir_fry("山药木耳鸡片", [ing("鸡胸肉", 145)], [ing("山药", 120), ing("木耳", 70)], meal_type="晚餐"),
        soup("番茄虾滑汤", [ing("虾仁", 145)], [ing("番茄", 140), ing("北豆腐", 80)], meal_type="晚餐"),
        steamed("葱油鸡腿", [ing("鸡腿肉(去皮)", 160)], [ing("青椒", 40)], meal_type="晚餐", time=20),
        stir_fry("西兰花牛肉", [ing("牛里脊", 145)], [ing("西兰花", 130)], meal_type="晚餐"),
        steamed_egg("虾仁蒸水蛋", [ing("虾仁", 85)], meal_type="晚餐"),
        soup("冬瓜鸡丸汤", [ing("鸡胸肉", 120)], [ing("冬瓜", 180), ing("香菇", 35)], meal_type="晚餐"),
        stir_fry("蒜蓉虾仁", [ing("虾仁", 150)], [ing("西兰花", 90)], meal_type="晚餐"),
        stew("香菇鸡腿煲", [ing("鸡腿肉(去皮)", 155)], [ing("香菇", 90), ing("番茄", 50)], meal_type="晚餐"),
        soup("番茄鱼片汤", [ing("鱼片", 150)], [ing("番茄", 150), ing("西兰花", 50)], meal_type="晚餐"),
        stir_fry("木耳牛肉片", [ing("牛里脊", 145)], [ing("木耳", 80), ing("青椒", 40)], meal_type="晚餐"),
        stew("豆腐鸡肉丸", [ing("鸡胸肉", 130), ing("北豆腐", 130)], [ing("番茄", 70)], meal_type="晚餐"),
        stir_fry("青椒肉丝", [ing("瘦猪里脊", 145)], [ing("青椒", 100), ing("胡萝卜", 35)], meal_type="晚餐"),
        stir_fry("山药牛肉片", [ing("牛里脊", 140)], [ing("山药", 120), ing("青椒", 35)], meal_type="晚餐"),
        stew("家常番茄鸡胸肉", [ing("鸡胸肉", 145)], [ing("番茄", 150), ing("洋葱", 35)], meal_type="晚餐"),
        stew("香菇豆腐煲", [ing("北豆腐", 170)], [ing("香菇", 100), ing("番茄", 70)], meal_type="晚餐"),
        stir_fry("芹菜鸡丝", [ing("鸡胸肉", 140)], [ing("芹菜", 120), ing("胡萝卜", 35)], meal_type="晚餐"),
        soup("鲈鱼豆腐汤", [ing("鲈鱼", 150), ing("北豆腐", 120)], [ing("冬瓜", 120)], meal_type="晚餐"),
        stir_fry("西葫芦虾仁", [ing("虾仁", 145)], [ing("西葫芦", 120), ing("胡萝卜", 35)], meal_type="晚餐"),
        stir_fry("胡萝卜牛肉丝", [ing("牛里脊", 140)], [ing("胡萝卜", 110), ing("青椒", 35)], meal_type="晚餐"),
        soup("虾仁冬瓜羹", [ing("虾仁", 135)], [ing("冬瓜", 180), ing("北豆腐", 80)], meal_type="晚餐"),
        steamed_egg("香菇肉末蒸蛋", [ing("猪肉末", 70), ing("香菇", 45)], meal_type="晚餐"),
        stir_fry("洋葱牛肉片", [ing("牛里脊", 145)], [ing("洋葱", 120), ing("青椒", 35)], meal_type="晚餐"),
        stir_fry("丝瓜虾仁", [ing("虾仁", 145)], [ing("丝瓜", 130), ing("番茄", 40)], meal_type="晚餐"),
        soup("番茄豆腐鸡蛋汤", [ing("鸡蛋", 60), ing("北豆腐", 120)], [ing("番茄", 150), ing("菠菜", 40)], meal_type="晚餐"),
    ]


def snack_recipes() -> list[dict]:
    return [
        simple_snack("白水煮蛋", [ing("鸡蛋", 60)], time=8),
        simple_snack("五香茶叶蛋", [ing("鸡蛋", 65)], time=12),
        simple_snack("无糖酸奶", [ing("无糖酸奶", 200)], time=2),
        simple_snack("现磨无糖豆浆", [ing("无糖豆浆", 260)], time=3),
        simple_snack("蒸玉米", [ing("玉米", 180)], time=12),
        simple_snack("烤红薯", [ing("红薯", 180)], time=20),
        simple_snack("软糯蒸紫薯", [ing("紫薯", 170)], time=18),
        rice_ball("海苔金枪鱼饭团", [ing("金枪鱼罐头", 70), ing("黄瓜", 35)], meal_type="加餐"),
        rice_ball("香煎鸡肉饭团", [ing("鸡胸肉", 80), ing("玉米", 35)], meal_type="加餐"),
        yogurt_cup("燕麦酸奶杯", [ing("无糖酸奶", 180), ing("燕麦", 30), ing("香蕉", 70)]),
        yogurt_cup("香蕉酸奶杯", [ing("无糖酸奶", 180), ing("香蕉", 110)]),
        yogurt_cup("蓝莓酸奶杯", [ing("无糖酸奶", 180), ing("蓝莓", 80)]),
        yogurt_cup("草莓酸奶杯", [ing("无糖酸奶", 180), ing("草莓", 110)]),
        simple_snack("毛豆杯", [ing("毛豆", 120)], time=6),
        simple_snack("卤香豆干", [ing("豆干", 90)], time=4),
        yogurt_cup("玉米鸡蛋杯", [ing("玉米", 120), ing("鸡蛋", 55), ing("无糖酸奶", 100)]),
        yogurt_cup("南瓜酸奶杯", [ing("南瓜", 140), ing("无糖酸奶", 150)]),
        simple_snack("山药泥", [ing("山药", 180)], time=12),
        steamed_egg("迷你虾仁蒸蛋", [ing("虾仁", 60)], meal_type="加餐", time=15),
        dough_wrap("牛肉饭卷", [ing("牛里脊", 75), ing("白米饭", 90), ing("黄瓜", 25)], meal_type="加餐", time=12),
        yogurt_cup("玉米燕麦杯", [ing("玉米", 90), ing("燕麦", 28), ing("无糖酸奶", 120)]),
        dough_wrap("鸡丝卷饼", [ing("鸡胸肉", 75), ing("生菜", 35), ing("黄瓜", 25)], meal_type="加餐", time=12),
        yogurt_cup("苹果酸奶杯", [ing("苹果", 120), ing("无糖酸奶", 160)]),
        simple_snack("毛豆玉米杯", [ing("毛豆", 90), ing("玉米", 100)], time=6),
        yogurt_cup("香蕉燕麦杯", [ing("香蕉", 100), ing("燕麦", 28), ing("无糖豆浆", 180)]),
    ]


def validate_counts(recipes: list[dict]) -> None:
    counts = {"早餐": 0, "午餐": 0, "晚餐": 0, "加餐": 0}
    for recipe in recipes:
        counts[recipe["meal_type"]] += 1
        assert "配" not in recipe["name"], recipe["name"]
    assert len(recipes) == 100, len(recipes)
    assert all(count == 25 for count in counts.values()), counts


def main() -> None:
    recipes = breakfast_recipes() + lunch_recipes() + dinner_recipes() + snack_recipes()
    validate_counts(recipes)
    apply_image_urls(recipes)
    OUTPUT_PATH.write_text(json.dumps(recipes, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"generated {len(recipes)} recipes -> {OUTPUT_PATH}")


if __name__ == "__main__":
    main()
