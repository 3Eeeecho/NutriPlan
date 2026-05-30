import argparse
import csv
import json
import os
import re
from pathlib import Path

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

DEFAULT_CSV = ".cache/yunyoujun_cook_zip/cook-main/app/data/recipe.csv"

BAD_NAME_WORDS = (
    "蛋糕",
    "布丁",
    "奶茶",
    "可乐",
    "冰淇淋",
    "酒",
    "鸡尾酒",
    "甜品",
    "废手",
    "开店",
    "国宴",
    "同款",
    "进阶吃法",
    "特调",
    "饮料",
)

BAD_TAGS = {"零食", "甜品", "饮品"}

ALLERGY_KEYWORDS = {
    "虾": "海鲜过敏",
    "蟹": "海鲜过敏",
    "鱼": "海鲜过敏",
    "贝": "海鲜过敏",
    "牛奶": "乳制品过敏",
    "奶酪": "乳制品过敏",
    "芝士": "乳制品过敏",
    "面": "麸质过敏",
    "面包": "麸质过敏",
    "吐司": "麸质过敏",
    "鸡蛋": "鸡蛋过敏",
    "豆腐": "大豆过敏",
    "豆浆": "大豆过敏",
    "香干": "大豆过敏",
    "花生": "坚果过敏",
    "核桃": "坚果过敏",
    "杏仁": "坚果过敏",
}


def parse_args():
    parser = argparse.ArgumentParser(description="Import recipes from YunYouJun/cook recipe.csv.")
    parser.add_argument("--csv", default=DEFAULT_CSV)
    parser.add_argument("--dry-run", action="store_true")
    parser.add_argument("--upsert", action="store_true")
    parser.add_argument("--limit", type=int, default=0)
    parser.add_argument("--output-json", default="")
    return parser.parse_args()


def split_values(value: str) -> list[str]:
    return [item.strip() for item in re.split(r"[、,，/]+", value or "") if item.strip()]


def clean_name(name: str) -> str:
    name = re.sub(r"[（(].*?[）)]", "", name or "").strip()
    name = re.sub(r"^(电饭煲版|空气炸锅版|微波炉版)", "", name).strip()
    return name


def should_skip(row: dict, name: str, tags: list[str]) -> bool:
    if not name or len(name) > 14:
        return True
    if any(word in name for word in BAD_NAME_WORDS):
        return True
    if BAD_TAGS.intersection(tags):
        return True
    if not split_values(row.get("stuff", "")):
        return True
    return False


def infer_meal_type(name: str, tags: list[str], stuff: list[str], method: str) -> tuple[str, list[str]]:
    text = name + " " + " ".join(tags) + " " + " ".join(stuff)
    if any(tag in tags for tag in ("早饭", "早餐")) or any(word in name for word in ("粥", "蛋饼", "吐司", "三明治", "花卷", "茶叶蛋")):
        return "早餐", ["早餐"]
    if any(word in text for word in ("饭", "面", "粉", "饺", "馄饨", "包", "烧卖", "饼")):
        return "午餐", ["午餐", "晚餐"]
    if any(word in name for word in ("汤", "羹", "粥")) or method in {"煮", "炖", "煲"}:
        return "晚餐", ["午餐", "晚餐"]
    return "午餐", ["午餐", "晚餐"]


def normalize_difficulty(value: str) -> str:
    if value == "简单":
        return "简单"
    if value in {"困难", "复杂"}:
        return "困难"
    return "中等"


def cooking_time(method: str, difficulty: str) -> int:
    base = {
        "炒": 15,
        "煎": 15,
        "拌": 10,
        "凉拌": 10,
        "蒸": 25,
        "煮": 25,
        "烧": 30,
        "炖": 45,
        "煲": 45,
        "烤": 35,
        "炸": 25,
        "卤": 60,
        "焖": 35,
    }.get(method, 20)
    if difficulty == "困难":
        base += 20
    elif difficulty == "中等":
        base += 8
    return min(base, 90)


def make_steps(name: str, stuff: list[str], method: str) -> list[str]:
    ingredients = "、".join(stuff)
    method = method or guess_method(name)
    if method in {"炒", "烧", "煎", "焖"}:
        return [
            f"准备{ingredients}，清洗后按需要切块或切片。",
            "热锅少油，先处理肉蛋类或不易熟的食材。",
            f"加入其余食材翻炒，按口味加入盐、生抽等调味。",
            f"继续{method}至食材熟透，收汁后盛出。",
        ]
    if method in {"煮", "炖", "煲"}:
        return [
            f"准备{ingredients}，清洗后切成适口大小。",
            "锅中加水或高汤，先放入耐煮食材。",
            f"转中小火{method}至食材熟软。",
            "按口味调入盐和少量调味料，盛出食用。",
        ]
    if method in {"蒸"}:
        return [
            f"准备{ingredients}，装入耐热盘中。",
            "蒸锅上汽后放入食材。",
            "蒸至熟透后取出。",
            "按口味加入少量调味料。",
        ]
    if method in {"拌", "凉拌"}:
        return [
            f"准备{ingredients}，清洗并切好。",
            "需要焯水的食材先焯熟后沥干。",
            "加入盐、生抽、醋等调味料。",
            "拌匀后装盘。",
        ]
    return [
        f"准备{ingredients}。",
        "按食材特点清洗、切配并预处理。",
        "使用少油少盐方式烹饪至熟透。",
        "调味后装盘。",
    ]


def guess_method(name: str) -> str:
    for method in ("炒", "煎", "蒸", "煮", "炖", "煲", "焖", "烤", "炸", "拌", "卤"):
        if method in name:
            return method
    return "炒"


def estimate_nutrition(name: str, meal_type: str, stuff: list[str], method: str):
    text = name + " " + " ".join(stuff)
    if meal_type == "早餐":
        portion, energy, protein, carb, fat = 220.0, 135.0, 6.0, 18.0, 4.5
    elif any(word in text for word in ("饭", "面", "粉", "饼", "包", "吐司", "面包", "米")):
        portion, energy, protein, carb, fat = 260.0, 165.0, 6.0, 25.0, 4.5
    elif any(word in text for word in ("汤", "粥", "羹")):
        portion, energy, protein, carb, fat = 300.0, 55.0, 3.5, 5.0, 2.0
    elif any(word in text for word in ("牛肉", "猪肉", "鸡肉", "虾", "鱼", "排骨")):
        portion, energy, protein, carb, fat = 190.0, 145.0, 13.0, 5.0, 8.0
    else:
        portion, energy, protein, carb, fat = 200.0, 85.0, 4.0, 8.0, 4.0

    if method in {"炸", "煎"} or any(word in name for word in ("炸", "干煸", "红烧", "糖醋", "油焖")):
        energy += 45
        fat += 5
    if any(word in text for word in ("鸡蛋", "豆腐", "虾", "鱼", "牛肉", "鸡肉")):
        protein += 4
    if any(word in text for word in ("土豆", "玉米", "南瓜", "米", "面", "粉")):
        carb += 8
    return portion, energy, protein, carb, fat


def infer_forbidden(name: str, stuff: list[str]) -> list[str]:
    text = name + " " + " ".join(stuff)
    result = []
    for keyword, tag in ALLERGY_KEYWORDS.items():
        if keyword in text and tag not in result:
            result.append(tag)
    return result


def target_flags(name: str, energy: float, protein: float, carb: float, fat: float, method: str):
    high_oil = method in {"炸"} or any(word in name for word in ("炸", "油焖", "干煸", "糖醋", "红烧"))
    is_loss = not high_oil and energy <= 145 and fat <= 8
    is_gain = protein >= 10
    is_sugar = carb <= 18 and not any(word in name for word in ("糖", "甜", "饭", "面", "粉", "饼", "包"))
    users = []
    if is_loss:
        users.append("减脂")
    if is_gain:
        users.append("增肌")
    if is_sugar:
        users.append("控糖")
    users.append("维持健康")
    return users, is_loss, is_gain, is_sugar, True


def row_to_recipe(row: dict):
    raw_name = row.get("name", "")
    name = clean_name(raw_name)
    tags = split_values(row.get("tags", ""))
    if should_skip(row, name, tags):
        return None
    stuff = split_values(row.get("stuff", ""))
    method = split_values(row.get("methods", ""))
    method = method[0] if method else guess_method(name)
    meal_type, allowed = infer_meal_type(name, tags, stuff, method)
    difficulty = normalize_difficulty(row.get("difficulty", ""))
    portion, energy, protein, carb, fat = estimate_nutrition(name, meal_type, stuff, method)
    target_users, is_loss, is_gain, is_sugar, is_general = target_flags(name, energy, protein, carb, fat, method)
    return {
        "name": name,
        "image_url": "",
        "meal_type": meal_type,
        "allowed_meal_types": allowed,
        "difficulty": difficulty,
        "cooking_time": cooking_time(method, difficulty),
        "portion_weight_g": portion,
        "energy": energy,
        "protein": protein,
        "carbohydrate": carb,
        "fat": fat,
        "ingredients": [f"{item} 适量" for item in stuff],
        "cooking_steps": make_steps(name, stuff, method),
        "target_users": target_users,
        "forbidden_users": infer_forbidden(name, stuff),
        "is_weight_loss_friendly": is_loss,
        "is_muscle_gain_friendly": is_gain,
        "is_sugar_control_friendly": is_sugar,
        "is_general_friendly": is_general,
    }


def connect_db():
    return pymysql.connect(**DB_CONFIG)


def insert_or_update(conn, recipe: dict, upsert: bool):
    json_fields = {"allowed_meal_types", "ingredients", "cooking_steps", "target_users", "forbidden_users"}
    values = {
        key: json.dumps(value, ensure_ascii=False) if key in json_fields else value
        for key, value in recipe.items()
    }
    with conn.cursor() as cursor:
        cursor.execute("SELECT id FROM recipes WHERE deleted_at IS NULL AND name = %s LIMIT 1", (recipe["name"],))
        existing = cursor.fetchone()
        if existing and not upsert:
            return "skipped"
        if existing:
            set_clause = ", ".join(f"`{key}` = %s" for key in values)
            cursor.execute(
                f"UPDATE recipes SET {set_clause}, updated_at = NOW() WHERE id = %s",
                list(values.values()) + [existing["id"]],
            )
            return "updated"
        fields = ["created_at", "updated_at", *values.keys()]
        placeholders = ["NOW()", "NOW()", *["%s"] * len(values)]
        cursor.execute(
            f"INSERT INTO recipes ({', '.join(f'`{field}`' for field in fields)}) "
            f"VALUES ({', '.join(placeholders)})",
            list(values.values()),
        )
        return "created"


def main():
    args = parse_args()
    rows = list(csv.DictReader(Path(args.csv).open(encoding="utf-8-sig")))
    recipes = []
    seen = set()
    for row in rows:
        recipe = row_to_recipe(row)
        if recipe is None:
            continue
        if recipe["name"] in seen:
            continue
        seen.add(recipe["name"])
        recipes.append(recipe)
        if args.limit > 0 and len(recipes) >= args.limit:
            break

    if args.output_json:
        Path(args.output_json).write_text(json.dumps(recipes, ensure_ascii=False, indent=2), encoding="utf-8")

    print(f"parsed={len(recipes)}, dry_run={args.dry_run}")
    for recipe in recipes[:20]:
        print(f"  {recipe['name']} | {recipe['meal_type']} | {','.join(recipe['ingredients'][:3])}")

    if args.dry_run:
        return

    conn = connect_db()
    created = updated = skipped = 0
    try:
        for recipe in recipes:
            action = insert_or_update(conn, recipe, args.upsert)
            if action == "created":
                created += 1
            elif action == "updated":
                updated += 1
            else:
                skipped += 1
        conn.commit()
    finally:
        conn.close()
    print(f"import finished: created={created}, updated={updated}, skipped={skipped}")


if __name__ == "__main__":
    main()
