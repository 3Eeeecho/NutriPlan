import argparse
import json
import os
from typing import Iterable

import pymysql
import pymysql.cursors


DB_CONFIG = {
    "host": os.getenv("NUTRIPLAN_DB_HOST", "127.0.0.1"),
    "user": os.getenv("NUTRIPLAN_DB_USER", "root"),
    "password": os.getenv("NUTRIPLAN_DB_PASSWORD", "root"),
    "database": os.getenv("NUTRIPLAN_DB_NAME", "nutri_plan"),
    "charset": "utf8mb4",
    "cursorclass": pymysql.cursors.DictCursor,
}


VALID_MEAL_TYPES = ("早餐", "午餐", "晚餐", "加餐")

HIGH_SUGAR_KEYWORDS = [
    "奶茶",
    "果茶",
    "水果茶",
    "糖水",
    "甜汤",
    "杨枝甘露",
    "冰淇淋",
    "奶冻",
    "布丁",
    "提拉米苏",
    "蛋糕",
    "雪花酥",
    "奶昔",
    "可乐",
    "汽水",
    "果汁",
    "沙冰",
    "冰沙",
]

ALCOHOL_KEYWORDS = [
    "酒",
    "鸡尾酒",
    "啤酒",
    "红酒",
    "白酒",
    "黄酒",
    "清酒",
    "威士忌",
    "伏特加",
    "朗姆",
    "金酒",
    "利口酒",
    "莫吉托",
    "mojito",
    "长岛冰茶",
]

FRIED_HEAVY_KEYWORDS = [
    "炸",
    "酥",
    "锅包",
    "红烧肉",
    "扣肉",
    "肥肠",
    "五花肉",
    "奶油",
    "黄油",
    "猪油",
]

CLEAN_PROTEIN_KEYWORDS = [
    "鸡胸",
    "鸡胸肉",
    "牛肉",
    "牛腱",
    "里脊",
    "虾仁",
    "虾",
    "鱼",
    "三文鱼",
    "金枪鱼",
    "鸡蛋",
    "蛋白",
    "豆腐",
    "豆干",
    "牛奶",
    "酸奶",
]

STAPLE_KEYWORDS = [
    "饭",
    "面",
    "粉",
    "粥",
    "燕麦",
    "馒头",
    "包子",
    "饺子",
    "馄饨",
    "米线",
    "意面",
    "土豆",
    "红薯",
    "玉米",
]

BREAKFAST_KEYWORDS = [
    "早餐",
    "燕麦",
    "麦片",
    "牛奶",
    "酸奶",
    "豆浆",
    "豆腐脑",
    "鸡蛋",
    "煎蛋",
    "水煮蛋",
    "蛋饼",
    "吐司",
    "三明治",
    "面包",
    "包子",
    "馒头",
    "花卷",
    "粥",
    "杂粮粥",
    "小米粥",
    "南瓜粥",
]

STRICT_BREAKFAST_KEYWORDS = [
    "粥",
    "小米粥",
    "杂粮粥",
    "燕麦",
    "麦片",
    "豆浆",
    "豆腐脑",
    "牛奶",
    "酸奶",
    "水煮蛋",
    "煎蛋",
    "蒸蛋",
    "茶叶蛋",
    "蛋饼",
    "鸡蛋饼",
    "煎饼",
    "三明治",
    "吐司",
    "面包",
    "包子",
    "馒头",
    "花卷",
    "烧麦",
    "肠粉",
    "馄饨",
    "云吞",
    "汤面",
    "米粉",
]

BREAKFAST_INAPPROPRIATE_KEYWORDS = [
    "小炒",
    "热炒",
    "炒菜",
    "炒蛋",
    "炒鸡蛋",
    "炒饭",
    "盖饭",
    "拌饭",
    "炒面",
    "拌面",
    "宫保",
    "鱼香",
    "麻婆",
    "红烧",
    "白灼",
    "清蒸",
    "爆炒",
    "干锅",
    "彩椒",
    "青椒",
    "芹菜",
    "荷兰豆",
    "菜心",
    "西兰花",
    "毛豆",
    "肉丝",
    "鸡丁",
    "牛柳",
    "滑鸡",
    "虾仁滑蛋",
]

BREAKFAST_SOUP_KEYWORDS = [
    "豆浆",
    "豆腐脑",
    "牛奶",
    "燕麦粥",
    "小米粥",
    "南瓜粥",
    "杂粮粥",
    "银耳羹",
]

SNACK_KEYWORDS = [
    "加餐",
    "酸奶",
    "坚果",
    "水果",
    "果盘",
    "奶昔",
    "沙拉杯",
    "能量棒",
    "蛋白棒",
    "布丁",
]

MAIN_MEAL_KEYWORDS = [
    "炒饭",
    "盖饭",
    "焗饭",
    "烩饭",
    "咖喱饭",
    "拌饭",
    "意面",
    "通心粉",
    "炒面",
    "拌面",
    "拉面",
    "汤面",
    "焖面",
    "炒粉",
    "河粉",
    "米线",
    "水饺",
    "饺子",
    "馄饨",
    "煲仔饭",
    "便当",
    "套餐",
    "披萨",
    "汉堡",
]

SOUP_KEYWORDS = [
    "汤",
    "羹",
    "炖",
    "煲",
    "罗宋汤",
]


def parse_args():
    parser = argparse.ArgumentParser(description="Batch classify recipes for suitability and allowed meal slots")
    parser.add_argument("--limit", type=int, default=0, help="Maximum rows to process, 0 means unlimited")
    parser.add_argument("--recipe-id", type=int, default=0, help="Only process one recipe id")
    parser.add_argument("--dry-run", action="store_true", help="Preview results without writing to DB")
    parser.add_argument(
        "--only-unclassified",
        action="store_true",
        help="Only process rows still using default suitability flags",
    )
    parser.add_argument(
        "--skip-suitability",
        action="store_true",
        help="Do not update suitability boolean fields",
    )
    parser.add_argument(
        "--skip-meal-slots",
        action="store_true",
        help="Do not update allowed_meal_types",
    )
    return parser.parse_args()


def get_db_connection():
    return pymysql.connect(**DB_CONFIG)


def get_recipe_columns() -> set[str]:
    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.execute("SHOW COLUMNS FROM recipes")
            return {str(row["Field"]).strip() for row in cursor.fetchall()}
    finally:
        conn.close()


def parse_json_list(value) -> list[str]:
    if value is None:
        return []
    if isinstance(value, list):
        return [str(item).strip() for item in value if str(item).strip()]
    if isinstance(value, str):
        text = value.strip()
        if not text:
            return []
        try:
            parsed = json.loads(text)
        except Exception:
            return [text]
        if isinstance(parsed, list):
            return [str(item).strip() for item in parsed if str(item).strip()]
        return [text]
    return [str(value).strip()]


def to_float(value) -> float:
    if value is None:
        return 0.0
    return float(value)


def contains_any(text: str, keywords: Iterable[str]) -> bool:
    return any(keyword.lower() in text for keyword in keywords if keyword)


def build_content(recipe) -> str:
    name = str(recipe.get("name", "")).strip().lower()
    ingredients = " ".join(parse_json_list(recipe.get("ingredients"))).lower()
    return f"{name} {ingredients}".strip()


def normalize_meal_type(value) -> str:
    text = str(value or "").strip()
    if text in VALID_MEAL_TYPES:
        return text
    return ""


def classify_recipe(recipe):
    content = build_content(recipe)
    energy = to_float(recipe.get("energy"))
    protein = to_float(recipe.get("protein"))
    carbohydrate = to_float(recipe.get("carbohydrate"))
    fat = to_float(recipe.get("fat"))
    dietary_fiber = to_float(recipe.get("dietary_fiber"))

    has_basic_nutrition = energy > 0 and protein >= 0 and carbohydrate >= 0 and fat >= 0
    has_alcohol = contains_any(content, ALCOHOL_KEYWORDS)
    has_high_sugar = contains_any(content, HIGH_SUGAR_KEYWORDS)
    has_fried_heavy = contains_any(content, FRIED_HEAVY_KEYWORDS)
    has_clean_protein = contains_any(content, CLEAN_PROTEIN_KEYWORDS)
    has_staple = contains_any(content, STAPLE_KEYWORDS)

    is_general = has_basic_nutrition and not has_alcohol and not has_high_sugar

    is_weight_loss = (
        is_general
        and not has_fried_heavy
        and energy <= 220
        and fat <= 12
        and carbohydrate <= 25
    )

    is_sugar_control = (
        is_general
        and carbohydrate <= 20
        and energy <= 230
        and (dietary_fiber <= 0 or dietary_fiber >= 1.5 or has_clean_protein)
    )

    is_muscle_gain = (
        is_general
        and energy >= 110
        and protein >= 8
        and (protein >= 12 or has_clean_protein or has_staple)
    )

    return {
        "is_weight_loss_friendly": bool(is_weight_loss),
        "is_muscle_gain_friendly": bool(is_muscle_gain),
        "is_sugar_control_friendly": bool(is_sugar_control),
        "is_general_friendly": bool(is_general),
    }


def classify_allowed_meal_types(recipe) -> list[str]:
    content = build_content(recipe)
    primary = normalize_meal_type(recipe.get("meal_type"))
    energy = to_float(recipe.get("energy"))

    has_breakfast = contains_any(content, BREAKFAST_KEYWORDS)
    has_breakfast_soup = contains_any(content, BREAKFAST_SOUP_KEYWORDS)
    has_snack = contains_any(content, SNACK_KEYWORDS)
    has_main_meal = contains_any(content, MAIN_MEAL_KEYWORDS)
    has_soup = contains_any(content, SOUP_KEYWORDS)

    allowed: set[str] = set()
    if primary:
        allowed.add(primary)

    if primary in {"午餐", "晚餐"}:
        allowed.update({"午餐", "晚餐"})

    if has_breakfast:
        allowed.add("早餐")

    if has_main_meal or has_soup:
        allowed.update({"午餐", "晚餐"})

    if has_snack and energy <= 280:
        allowed.add("加餐")

    if primary == "早餐" and energy <= 280:
        allowed.add("早餐")
        if has_snack or "酸奶" in content or "燕麦" in content or "水果" in content:
            allowed.add("加餐")

    if primary == "加餐":
        allowed.add("加餐")
        if energy <= 260 and (has_breakfast or "酸奶" in content or "燕麦" in content):
            allowed.add("早餐")

    if has_soup and not has_breakfast_soup:
        allowed.discard("早餐")

    if has_main_meal:
        allowed.discard("早餐")
        if energy > 320:
            allowed.discard("加餐")

    if primary in {"午餐", "晚餐"} and not has_breakfast:
        allowed.discard("早餐")

    if not allowed:
        allowed.add(primary or "午餐")

    return [meal_type for meal_type in VALID_MEAL_TYPES if meal_type in allowed]


def classify_allowed_meal_types_v2(recipe) -> list[str]:
    content = build_content(recipe)
    primary = normalize_meal_type(recipe.get("meal_type"))
    breakfast = VALID_MEAL_TYPES[0]
    lunch = VALID_MEAL_TYPES[1]
    dinner = VALID_MEAL_TYPES[2]
    snack = VALID_MEAL_TYPES[3]

    has_breakfast = contains_any(content, STRICT_BREAKFAST_KEYWORDS)
    breakfast_inappropriate = contains_any(content, BREAKFAST_INAPPROPRIATE_KEYWORDS)
    has_snack = contains_any(content, SNACK_KEYWORDS)
    has_main_meal = contains_any(content, MAIN_MEAL_KEYWORDS)
    has_soup = contains_any(content, SOUP_KEYWORDS)
    energy = to_float(recipe.get("energy"))

    allowed: set[str] = set()

    if has_breakfast and not breakfast_inappropriate:
        allowed.add(breakfast)
        if has_snack or ("酸奶" in content or "燕麦" in content or "水果" in content) and energy <= 280:
            allowed.add(snack)

    if primary == snack:
        allowed.add(snack)
        if has_breakfast and not breakfast_inappropriate and energy <= 260:
            allowed.add(breakfast)

    if primary in {lunch, dinner} or has_main_meal or has_soup or breakfast_inappropriate:
        allowed.update({lunch, dinner})

    if has_snack and not has_main_meal and not breakfast_inappropriate and energy <= 280:
        allowed.add(snack)

    if primary == breakfast and not breakfast_inappropriate and has_breakfast:
        allowed.add(breakfast)

    if not allowed:
        if primary == breakfast:
            allowed.update({lunch, dinner})
        elif primary in {lunch, dinner}:
            allowed.update({lunch, dinner})
        elif primary == snack:
            allowed.add(snack)
        else:
            allowed.add(lunch)

    if breakfast_inappropriate:
        allowed.discard(breakfast)

    return [meal_type for meal_type in VALID_MEAL_TYPES if meal_type in allowed]


classify_allowed_meal_types = classify_allowed_meal_types_v2


def load_recipes(args, columns: set[str]):
    select_fields = [
        "id",
        "name",
        "meal_type",
        "ingredients",
        "energy",
        "protein",
        "carbohydrate",
        "fat",
        "is_weight_loss_friendly",
        "is_muscle_gain_friendly",
        "is_sugar_control_friendly",
        "is_general_friendly",
    ]
    if "dietary_fiber" in columns:
        select_fields.append("dietary_fiber")
    if "allowed_meal_types" in columns:
        select_fields.append("allowed_meal_types")

    sql = f"""
    SELECT
        {", ".join(select_fields)}
    FROM recipes
    WHERE deleted_at IS NULL
    """
    params = []

    if args.recipe_id > 0:
        sql += " AND id = %s"
        params.append(args.recipe_id)

    if args.only_unclassified:
        sql += """
        AND (
            is_weight_loss_friendly = 0
            AND is_muscle_gain_friendly = 0
            AND is_sugar_control_friendly = 0
            AND is_general_friendly = 1
        )
        """

    sql += " ORDER BY id ASC"
    if args.limit > 0:
        sql += " LIMIT %s"
        params.append(args.limit)

    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.execute(sql, params)
            return cursor.fetchall()
    finally:
        conn.close()


def batch_update(updates, update_suitability: bool, update_meal_slots: bool):
    set_clauses = []
    if update_suitability:
        set_clauses.extend(
            [
                "is_weight_loss_friendly = %s",
                "is_muscle_gain_friendly = %s",
                "is_sugar_control_friendly = %s",
                "is_general_friendly = %s",
            ]
        )
    if update_meal_slots:
        set_clauses.append("allowed_meal_types = %s")
    set_clauses.append("updated_at = NOW()")

    sql = f"""
    UPDATE recipes
    SET
        {", ".join(set_clauses)}
    WHERE id = %s
    """

    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.executemany(sql, updates)
        conn.commit()
    finally:
        conn.close()


def main():
    args = parse_args()
    columns = get_recipe_columns()
    update_suitability = not args.skip_suitability
    update_meal_slots = not args.skip_meal_slots and "allowed_meal_types" in columns

    if not update_suitability and not update_meal_slots:
        print("nothing to update")
        return

    if not args.skip_meal_slots and "allowed_meal_types" not in columns:
        print("allowed_meal_types column not found, meal slot backfill will be skipped")

    recipes = load_recipes(args, columns)
    if not recipes:
        print("no recipes to process")
        return

    updates = []
    counts = {
        "weight_loss": 0,
        "muscle_gain": 0,
        "sugar_control": 0,
        "general": 0,
        "all_false": 0,
        "breakfast": 0,
        "lunch": 0,
        "dinner": 0,
        "snack": 0,
    }

    for recipe in recipes:
        suitability = classify_recipe(recipe)
        allowed_meal_types = classify_allowed_meal_types(recipe)

        row = []
        if update_suitability:
            row.extend(
                [
                    suitability["is_weight_loss_friendly"],
                    suitability["is_muscle_gain_friendly"],
                    suitability["is_sugar_control_friendly"],
                    suitability["is_general_friendly"],
                ]
            )
        if update_meal_slots:
            row.append(json.dumps(allowed_meal_types, ensure_ascii=False))
        row.append(recipe["id"])
        updates.append(tuple(row))

        counts["weight_loss"] += int(suitability["is_weight_loss_friendly"])
        counts["muscle_gain"] += int(suitability["is_muscle_gain_friendly"])
        counts["sugar_control"] += int(suitability["is_sugar_control_friendly"])
        counts["general"] += int(suitability["is_general_friendly"])
        counts["all_false"] += int(not any(suitability.values()))
        counts["breakfast"] += int("早餐" in allowed_meal_types)
        counts["lunch"] += int("午餐" in allowed_meal_types)
        counts["dinner"] += int("晚餐" in allowed_meal_types)
        counts["snack"] += int("加餐" in allowed_meal_types)

    print(f"to process: {len(recipes)}")
    if update_suitability:
        print(f"weight_loss: {counts['weight_loss']}")
        print(f"muscle_gain: {counts['muscle_gain']}")
        print(f"sugar_control: {counts['sugar_control']}")
        print(f"general: {counts['general']}")
        print(f"all_false: {counts['all_false']}")
    if update_meal_slots:
        print(f"breakfast_allowed: {counts['breakfast']}")
        print(f"lunch_allowed: {counts['lunch']}")
        print(f"dinner_allowed: {counts['dinner']}")
        print(f"snack_allowed: {counts['snack']}")

    if args.dry_run:
        for recipe, update in zip(recipes[:20], updates[:20]):
            payload = {
                "id": recipe["id"],
                "name": recipe["name"],
                "meal_type": recipe.get("meal_type"),
                "allowed_meal_types": classify_allowed_meal_types(recipe),
            }
            if update_suitability:
                payload.update(classify_recipe(recipe))
            print("sample:", payload)
        print("dry-run complete")
        return

    batch_update(updates, update_suitability, update_meal_slots)
    print(f"updated: {len(updates)}")


if __name__ == "__main__":
    main()
