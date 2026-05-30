import argparse
import json
import os
import re
from dataclasses import dataclass
from pathlib import Path
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

DEFAULT_SOURCE_DIR = ".cache/howtocook_zip/HowToCook-master"

SKIP_CATEGORIES = {
    "condiment",
    "dessert",
    "drink",
    "semi-finished",
    "template",
}

CATEGORY_MEAL_TYPES = {
    "breakfast": ("早餐", ["早餐"]),
    "staple": ("午餐", ["午餐", "晚餐"]),
    "soup": ("晚餐", ["午餐", "晚餐"]),
    "aquatic": ("午餐", ["午餐", "晚餐"]),
    "meat_dish": ("午餐", ["午餐", "晚餐"]),
    "vegetable_dish": ("晚餐", ["午餐", "晚餐"]),
}

TOOL_WORDS = (
    "锅",
    "碗",
    "盘",
    "刀",
    "筷",
    "勺",
    "铲",
    "案板",
    "砧板",
    "烤箱",
    "空气炸锅",
    "微波炉",
    "电饭煲",
    "高压锅",
    "料理机",
    "保鲜膜",
)

BAD_NAME_WORDS = (
    "奶茶",
    "可乐",
    "酒",
    "鸡尾酒",
    "冰淇淋",
    "蛋糕",
    "布丁",
    "甜品",
    "阿根廷",
    "罗氏",
    "黑鳕",
    "黄油",
    "芥末",
    "海参",
    "咖喱",
    "意式",
    "美式",
    "法式",
)

ALLERGY_KEYWORDS = {
    "虾": "海鲜过敏",
    "蟹": "海鲜过敏",
    "鱼": "海鲜过敏",
    "贝": "海鲜过敏",
    "生蚝": "海鲜过敏",
    "牛奶": "乳制品过敏",
    "奶酪": "乳制品过敏",
    "酸奶": "乳制品过敏",
    "芝士": "乳制品过敏",
    "花生": "坚果过敏",
    "核桃": "坚果过敏",
    "杏仁": "坚果过敏",
    "腰果": "坚果过敏",
    "面粉": "麸质过敏",
    "面条": "麸质过敏",
    "全麦": "麸质过敏",
    "鸡蛋": "鸡蛋过敏",
    "豆腐": "大豆过敏",
    "豆干": "大豆过敏",
    "豆浆": "大豆过敏",
}


@dataclass
class ParsedRecipe:
    source_path: Path
    category: str
    name: str
    meal_type: str
    allowed_meal_types: list[str]
    difficulty: str
    cooking_time: int
    portion_weight_g: float
    energy: float
    protein: float
    carbohydrate: float
    fat: float
    ingredients: list[str]
    cooking_steps: list[str]
    target_users: list[str]
    forbidden_users: list[str]
    is_weight_loss_friendly: bool
    is_muscle_gain_friendly: bool
    is_sugar_control_friendly: bool
    is_general_friendly: bool


def parse_args():
    parser = argparse.ArgumentParser(description="Import common Chinese recipes from Anduin2017/HowToCook.")
    parser.add_argument("--source-dir", default=DEFAULT_SOURCE_DIR, help="Local HowToCook repository directory.")
    parser.add_argument("--limit", type=int, default=0, help="Maximum recipes to import after filtering.")
    parser.add_argument("--dry-run", action="store_true", help="Preview without writing to MySQL.")
    parser.add_argument("--upsert", action="store_true", help="Update rows with the same recipe name.")
    parser.add_argument(
        "--include-category",
        action="append",
        default=[],
        help="Import only these HowToCook category directories. Can be repeated.",
    )
    parser.add_argument(
        "--skip-category",
        action="append",
        default=[],
        help="Additional category directories to skip. Can be repeated.",
    )
    parser.add_argument("--output-json", default="", help="Write parsed recipes to a JSON file instead of MySQL only.")
    return parser.parse_args()


def iter_markdown_files(source_dir: Path, include_categories: set[str], skip_categories: set[str]) -> Iterable[Path]:
    dishes_dir = source_dir / "dishes"
    if not dishes_dir.exists():
        raise FileNotFoundError(f"HowToCook dishes directory not found: {dishes_dir}")

    for category_dir in sorted(path for path in dishes_dir.iterdir() if path.is_dir()):
        category = category_dir.name
        if include_categories and category not in include_categories:
            continue
        if category in skip_categories:
            continue
        if category not in CATEGORY_MEAL_TYPES:
            continue
        yield from sorted(category_dir.rglob("*.md"))


def clean_title(raw: str, path: Path) -> str:
    title = raw.strip().lstrip("#").strip()
    if not title:
        title = path.stem
    title = re.sub(r"的?做法$", "", title).strip()
    title = re.sub(r"怎么做$", "", title).strip()
    title = re.sub(r"[（(].*?[）)]", "", title).strip()
    return title or path.stem


def split_sections(text: str) -> dict[str, list[str]]:
    sections: dict[str, list[str]] = {}
    current = ""
    for line in text.splitlines():
        match = re.match(r"^##+\s+(.+?)\s*$", line)
        if match:
            current = match.group(1).strip()
            sections.setdefault(current, [])
            continue
        if current:
            sections[current].append(line.rstrip())
    return sections


def section_by_keyword(sections: dict[str, list[str]], keywords: tuple[str, ...]) -> list[str]:
    for title, lines in sections.items():
        if any(keyword in title for keyword in keywords):
            return lines
    return []


def strip_markdown_item(line: str) -> str:
    text = re.sub(r"^\s*[-*+]\s+", "", line).strip()
    text = re.sub(r"^\s*\d+[.)、]\s+", "", text).strip()
    text = text.replace("`", "").strip()
    return text


def extract_list_items(lines: list[str]) -> list[str]:
    items = []
    for line in lines:
        if not re.match(r"^\s*(?:[-*+]|\d+[.)、])\s+", line):
            continue
        item = strip_markdown_item(line)
        if item:
            items.append(item)
    return items


def clean_ingredients(items: list[str]) -> list[str]:
    result = []
    seen = set()
    for item in items:
        item = re.sub(r"\s+", " ", item).strip()
        if not item:
            continue
        if any(word in item for word in TOOL_WORDS):
            continue
        if item in seen:
            continue
        seen.add(item)
        result.append(item)
    return result[:12]


def clean_steps(items: list[str]) -> list[str]:
    result = []
    for item in items:
        item = re.sub(r"\s+", " ", item).strip()
        if not item or item.startswith("注意"):
            continue
        result.append(item)
    return result[:8]


def parse_difficulty(text: str) -> str:
    match = re.search(r"预估烹饪难度[：:]\s*([★☆]+)", text)
    stars = match.group(1).count("★") if match else 2
    if stars <= 2:
        return "简单"
    if stars <= 4:
        return "中等"
    return "困难"


def estimate_cooking_time(name: str, category: str, difficulty: str) -> int:
    if category == "breakfast":
        base = 12
    elif category == "soup":
        base = 35
    elif category == "staple":
        base = 25
    else:
        base = 18

    if any(word in name for word in ("炖", "煲", "卤", "红烧", "烤")):
        base = max(base, 40)
    if any(word in name for word in ("凉拌", "白灼", "清炒", "滑蛋", "煎蛋")):
        base = min(base, 15)
    if difficulty == "困难":
        base += 20
    elif difficulty == "中等":
        base += 8
    return min(base, 90)


def estimate_nutrition(name: str, category: str) -> tuple[float, float, float, float, float]:
    if category == "breakfast":
        portion, energy, protein, carb, fat = 250.0, 120.0, 6.0, 16.0, 4.0
    elif category == "staple":
        portion, energy, protein, carb, fat = 220.0, 150.0, 5.0, 28.0, 2.5
    elif category == "soup":
        portion, energy, protein, carb, fat = 280.0, 55.0, 4.0, 4.5, 2.5
    elif category == "vegetable_dish":
        portion, energy, protein, carb, fat = 220.0, 85.0, 4.0, 8.0, 4.5
    elif category == "aquatic":
        portion, energy, protein, carb, fat = 180.0, 120.0, 16.0, 4.0, 5.0
    else:
        portion, energy, protein, carb, fat = 180.0, 165.0, 16.0, 6.0, 8.0

    if any(word in name for word in ("炸", "干煸", "油焖", "红烧", "糖醋")):
        energy += 45
        fat += 5
    if any(word in name for word in ("鸡胸", "牛肉", "鱼", "虾", "蛋", "豆腐", "豆干")):
        protein += 4
    if any(word in name for word in ("饭", "面", "粉", "饼", "馒头", "包")):
        carb += 10
        energy += 25
    if any(word in name for word in ("凉拌", "清蒸", "白灼", "水煮")):
        fat = max(1.5, fat - 2)

    return portion, energy, protein, carb, fat


def infer_forbidden_users(name: str, ingredients: list[str]) -> list[str]:
    text = name + " " + " ".join(ingredients)
    result = []
    for keyword, tag in ALLERGY_KEYWORDS.items():
        if keyword in text and tag not in result:
            result.append(tag)
    return result


def build_target_flags(name: str, category: str, energy: float, protein: float, carb: float, fat: float):
    high_oil = any(word in name for word in ("炸", "油焖", "干煸", "糖醋", "红烧"))
    is_weight_loss = not high_oil and energy <= 150 and fat <= 8
    is_muscle_gain = protein >= 10 and category in {"aquatic", "meat_dish", "breakfast", "vegetable_dish"}
    is_sugar_control = carb <= 18 and not any(word in name for word in ("糖", "甜", "饭", "面", "粉", "饼"))
    is_general = True

    target_users = []
    if is_weight_loss:
        target_users.append("减脂")
    if is_muscle_gain:
        target_users.append("增肌")
    if is_sugar_control:
        target_users.append("控糖")
    target_users.append("维持健康")
    return target_users, is_weight_loss, is_muscle_gain, is_sugar_control, is_general


def should_skip_name(name: str) -> bool:
    if not name or len(name) > 16:
        return True
    return any(word in name for word in BAD_NAME_WORDS)


def parse_recipe(path: Path, source_dir: Path) -> ParsedRecipe | None:
    text = path.read_text(encoding="utf-8", errors="ignore")
    title_line = next((line for line in text.splitlines() if line.startswith("# ")), "")
    category = path.relative_to(source_dir / "dishes").parts[0]
    name = clean_title(title_line, path)
    if should_skip_name(name):
        return None

    sections = split_sections(text)
    ingredient_lines = section_by_keyword(sections, ("必备原料", "原料"))
    step_lines = section_by_keyword(sections, ("操作", "步骤", "制作"))
    ingredients = clean_ingredients(extract_list_items(ingredient_lines))
    steps = clean_steps(extract_list_items(step_lines))
    if len(ingredients) < 2 or len(steps) < 2:
        return None

    meal_type, allowed_meal_types = CATEGORY_MEAL_TYPES[category]
    difficulty = parse_difficulty(text)
    cooking_time = estimate_cooking_time(name, category, difficulty)
    portion, energy, protein, carb, fat = estimate_nutrition(name, category)
    target_users, is_loss, is_gain, is_sugar, is_general = build_target_flags(name, category, energy, protein, carb, fat)

    return ParsedRecipe(
        source_path=path,
        category=category,
        name=name,
        meal_type=meal_type,
        allowed_meal_types=allowed_meal_types,
        difficulty=difficulty,
        cooking_time=cooking_time,
        portion_weight_g=portion,
        energy=energy,
        protein=protein,
        carbohydrate=carb,
        fat=fat,
        ingredients=ingredients,
        cooking_steps=steps,
        target_users=target_users,
        forbidden_users=infer_forbidden_users(name, ingredients),
        is_weight_loss_friendly=is_loss,
        is_muscle_gain_friendly=is_gain,
        is_sugar_control_friendly=is_sugar,
        is_general_friendly=is_general,
    )


def recipe_to_row(recipe: ParsedRecipe) -> dict:
    return {
        "name": recipe.name,
        "image_url": "",
        "meal_type": recipe.meal_type,
        "allowed_meal_types": recipe.allowed_meal_types,
        "difficulty": recipe.difficulty,
        "cooking_time": recipe.cooking_time,
        "portion_weight_g": recipe.portion_weight_g,
        "energy": recipe.energy,
        "protein": recipe.protein,
        "carbohydrate": recipe.carbohydrate,
        "fat": recipe.fat,
        "ingredients": recipe.ingredients,
        "cooking_steps": recipe.cooking_steps,
        "target_users": recipe.target_users,
        "forbidden_users": recipe.forbidden_users,
        "is_weight_loss_friendly": recipe.is_weight_loss_friendly,
        "is_muscle_gain_friendly": recipe.is_muscle_gain_friendly,
        "is_sugar_control_friendly": recipe.is_sugar_control_friendly,
        "is_general_friendly": recipe.is_general_friendly,
    }


def connect_db():
    return pymysql.connect(**DB_CONFIG)


def existing_names(conn) -> set[str]:
    with conn.cursor() as cursor:
        cursor.execute("SELECT name FROM recipes WHERE deleted_at IS NULL")
        return {row["name"] for row in cursor.fetchall()}


def insert_or_update(conn, recipe: ParsedRecipe, upsert: bool):
    row = recipe_to_row(recipe)
    json_fields = {"allowed_meal_types", "ingredients", "cooking_steps", "target_users", "forbidden_users"}
    db_row = {
        key: json.dumps(value, ensure_ascii=False) if key in json_fields else value
        for key, value in row.items()
    }

    with conn.cursor() as cursor:
        cursor.execute("SELECT id FROM recipes WHERE deleted_at IS NULL AND name = %s LIMIT 1", (recipe.name,))
        existing = cursor.fetchone()
        if existing and not upsert:
            return "skipped"
        if existing:
            set_clause = ", ".join(f"`{key}` = %s" for key in db_row)
            values = list(db_row.values()) + [existing["id"]]
            cursor.execute(f"UPDATE recipes SET {set_clause}, updated_at = NOW() WHERE id = %s", values)
            return "updated"

        fields = ["created_at", "updated_at", *db_row.keys()]
        placeholders = ["NOW()", "NOW()", *["%s"] * len(db_row)]
        cursor.execute(
            f"INSERT INTO recipes ({', '.join(f'`{field}`' for field in fields)}) "
            f"VALUES ({', '.join(placeholders)})",
            list(db_row.values()),
        )
        return "created"


def main():
    args = parse_args()
    source_dir = Path(args.source_dir).resolve()
    include_categories = set(args.include_category)
    skip_categories = SKIP_CATEGORIES | set(args.skip_category)

    parsed: list[ParsedRecipe] = []
    failed = 0
    for path in iter_markdown_files(source_dir, include_categories, skip_categories):
        try:
            recipe = parse_recipe(path, source_dir)
        except Exception as exc:
            failed += 1
            print(f"parse failed: {path}: {exc}")
            continue
        if recipe is None:
            continue
        parsed.append(recipe)
        if args.limit > 0 and len(parsed) >= args.limit:
            break

    if args.output_json:
        Path(args.output_json).write_text(
            json.dumps([recipe_to_row(recipe) for recipe in parsed], ensure_ascii=False, indent=2),
            encoding="utf-8",
        )

    print(f"parsed={len(parsed)}, failed={failed}, dry_run={args.dry_run}")
    for recipe in parsed[:20]:
        print(f"  {recipe.name} | {recipe.meal_type} | {recipe.category} | {len(recipe.ingredients)} ingredients")

    if args.dry_run:
        return

    conn = connect_db()
    created = updated = skipped = 0
    try:
        for recipe in parsed:
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
