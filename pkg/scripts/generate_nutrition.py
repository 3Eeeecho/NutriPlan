import argparse
import json
import os
import re
import time
from concurrent.futures import ThreadPoolExecutor, as_completed

import pymysql
import pymysql.cursors
from openai import OpenAI
from tqdm import tqdm


DB_CONFIG = {
    "host": os.getenv("NUTRIPLAN_DB_HOST", "127.0.0.1"),
    "user": os.getenv("NUTRIPLAN_DB_USER", "root"),
    "password": os.getenv("NUTRIPLAN_DB_PASSWORD", "root"),
    "database": os.getenv("NUTRIPLAN_DB_NAME", "nutri_plan"),
    "charset": "utf8mb4",
    "cursorclass": pymysql.cursors.DictCursor,
}

BASE_URL = os.getenv("NUTRIPLAN_LLM_BASE_URL", "https://api.siliconflow.cn/v1")
MODEL_NAME = os.getenv("NUTRIPLAN_LLM_MODEL", "deepseek-ai/DeepSeek-V3.2")
API_KEY = os.getenv("NUTRIPLAN_LLM_API_KEY", "")

DEFAULT_WORKERS = int(os.getenv("NUTRIPLAN_LLM_WORKERS", "3"))
DEFAULT_LIMIT = int(os.getenv("NUTRIPLAN_LLM_LIMIT", "0"))
DEFAULT_MIN_CONFIDENCE = float(os.getenv("NUTRIPLAN_LLM_MIN_CONFIDENCE", "0.7"))
DEFAULT_RETRY = int(os.getenv("NUTRIPLAN_LLM_RETRY", "3"))


def create_client():
    if not API_KEY:
        raise RuntimeError("缺少环境变量 NUTRIPLAN_LLM_API_KEY")
    return OpenAI(api_key=API_KEY, base_url=BASE_URL, timeout=120.0)


def get_db_connection():
    return pymysql.connect(**DB_CONFIG)


def parse_args():
    parser = argparse.ArgumentParser(description="为 recipes 回填营养数据和饮食目标约束字段")
    parser.add_argument("--limit", type=int, default=DEFAULT_LIMIT, help="最多处理多少条，0 表示不限制")
    parser.add_argument("--workers", type=int, default=DEFAULT_WORKERS, help="并发线程数")
    parser.add_argument("--recipe-id", type=int, default=0, help="只处理指定 recipe id")
    parser.add_argument("--all", action="store_true", help="处理所有 recipes，而不是仅处理缺失营养的记录")
    parser.add_argument("--dry-run", action="store_true", help="只估算不写库")
    parser.add_argument("--retry", type=int, default=DEFAULT_RETRY, help="单条失败重试次数")
    parser.add_argument("--min-confidence", type=float, default=DEFAULT_MIN_CONFIDENCE, help="低于该置信度则跳过")
    parser.add_argument("--keep-invalid", action="store_true", help="回填后不软删除仍无营养数据的菜谱")
    return parser.parse_args()


def load_recipes(args):
    sql = """
    SELECT
        id, name, meal_type, difficulty, cooking_time,
        portion_weight_g, energy, protein, carbohydrate, fat,
        ingredients, cooking_steps
    FROM recipes
    WHERE deleted_at IS NULL
    """
    params = []

    if args.recipe_id > 0:
        sql += " AND id = %s"
        params.append(args.recipe_id)
    elif not args.all:
        sql += """
        AND (
            portion_weight_g IS NULL OR portion_weight_g <= 0
            OR energy IS NULL OR energy <= 0
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


def parse_json_list(value):
    if value is None:
        return []
    if isinstance(value, list):
        return [str(item).strip() for item in value if str(item).strip()]
    if isinstance(value, str):
        value = value.strip()
        if not value:
            return []
        try:
            parsed = json.loads(value)
            if isinstance(parsed, list):
                return [str(item).strip() for item in parsed if str(item).strip()]
        except Exception:
            return [value]
    return [str(value).strip()]


def extract_json_from_text(text):
    if not text:
        return None

    text = text.strip()
    try:
        return json.loads(text)
    except Exception:
        pass

    match = re.search(r"\{.*\}", text, re.DOTALL)
    if not match:
        return None

    try:
        return json.loads(match.group())
    except Exception:
        return None


def build_prompt(recipe):
    ingredients = parse_json_list(recipe.get("ingredients"))
    steps = parse_json_list(recipe.get("cooking_steps"))

    ingredients_text = "\n".join(f"- {item}" for item in ingredients) or "- 无"
    steps_text = "\n".join(f"{idx + 1}. {item}" for idx, item in enumerate(steps)) or "无"

    return f"""
你是一名严谨的注册营养师。现在不要生成新菜谱，也不要改写菜名，只对现有菜谱做营养估算。

请根据下面这道菜的食材、步骤、餐次信息，估算：
1. 整道菜最终可食用总重量 total_weight_g
2. 单人推荐食用重量 recommended_portion_weight_g
3. 整道菜总热量 total_energy_kcal
4. 整道菜总蛋白 total_protein_g
5. 整道菜总碳水 total_carbohydrate_g
6. 整道菜总脂肪 total_fat_g
7. confidence，范围 0 到 1
8. reason，简短说明估算依据

要求：
- 只能估算营养，不要补全菜谱文案
- 结合常见中餐烹饪损耗、吸水、出水、调味品用量做合理估算
- 如果信息明显不足，也必须尽量给出保守估值
- 输出必须是单个 JSON 对象，不要输出 markdown，不要输出解释文字
- 所有数值字段必须是数字
- confidence 低于 0.7 代表你认为这道菜信息明显不足
- 单人推荐食用重量必须按餐次给出合理范围：
  - 早餐: 180-350g
  - 午餐: 250-450g
  - 晚餐: 250-450g
  - 加餐: 80-220g
- recommended_portion_weight_g 不得大于 total_weight_g

菜谱信息：
name: {recipe.get("name", "")}
meal_type: {recipe.get("meal_type", "")}
difficulty: {recipe.get("difficulty", "")}
cooking_time: {recipe.get("cooking_time", 0)}

ingredients:
{ingredients_text}

steps:
{steps_text}

输出格式：
{{
  "total_weight_g": 0,
  "recommended_portion_weight_g": 0,
  "total_energy_kcal": 0,
  "total_protein_g": 0,
  "total_carbohydrate_g": 0,
  "total_fat_g": 0,
  "confidence": 0,
  "reason": ""
}}
""".strip()


def request_estimate(client, recipe):
    payload = {
        "model": MODEL_NAME,
        "messages": [{"role": "user", "content": build_prompt(recipe)}],
        "temperature": 0.2,
        "max_tokens": 1200,
    }

    try:
        response = client.chat.completions.create(
            **payload,
            response_format={"type": "json_object"},
        )
    except Exception as exc:
        error_text = str(exc).lower()
        if "response_format" not in error_text and "json_object" not in error_text:
            raise
        response = client.chat.completions.create(**payload)

    content = response.choices[0].message.content
    return extract_json_from_text(content)


def to_float(value):
    if value is None:
        return 0.0
    if isinstance(value, (int, float)):
        return float(value)
    if isinstance(value, str):
        cleaned = value.strip().replace(",", "")
        match = re.search(r"-?\d+(?:\.\d+)?", cleaned)
        if match:
            return float(match.group())
    raise ValueError(f"无法转成数字: {value}")


def contains_any(text, keywords):
    content = str(text or "").strip().lower()
    return any(keyword.lower() in content for keyword in keywords)


def infer_min_energy_per_100g(recipe):
    content = f"{recipe.get('name', '')} {' '.join(parse_json_list(recipe.get('ingredients')))}"

    milk_tea_keywords = [
        "奶茶", "鸳鸯", "奶绿", "奶咖", "奶昔", "珍珠奶茶", "丝袜奶茶",
        "港式奶茶", "印度奶茶", "奶盖", "可可奶", "焦糖奶茶",
    ]
    cocktail_keywords = [
        "莫吉托", "mojito", "鸡尾酒", "长岛冰茶", "margarita", "玛格丽塔",
        "b52", "特调", "利口酒", "威士忌", "伏特加", "朗姆", "金酒",
    ]
    dessert_drink_keywords = [
        "冰淇淋", "奶冻", "提拉米苏", "奶昔", "布丁", "雪花酥",
    ]

    if contains_any(content, cocktail_keywords):
        return 60.0
    if contains_any(content, milk_tea_keywords):
        return 55.0
    if contains_any(content, dessert_drink_keywords):
        return 120.0
    return 0.0


def apply_reasonableness_rules(recipe, estimate):
    total_weight = estimate["total_weight_g"]
    macro_energy = (
        estimate["total_protein_g"] * 4
        + estimate["total_carbohydrate_g"] * 4
        + estimate["total_fat_g"] * 9
    )
    macro_floor = macro_energy * 0.95
    keyword_floor = infer_min_energy_per_100g(recipe) * total_weight / 100.0
    estimate["total_energy_kcal"] = round(
        max(estimate["total_energy_kcal"], macro_floor, keyword_floor),
        2,
    )
    return estimate


def normalize_estimate(data):
    if not isinstance(data, dict):
        raise ValueError("模型未返回 JSON 对象")

    estimate = {
        "total_weight_g": round(to_float(data.get("total_weight_g")), 2),
        "recommended_portion_weight_g": round(to_float(data.get("recommended_portion_weight_g")), 2),
        "total_energy_kcal": round(to_float(data.get("total_energy_kcal")), 2),
        "total_protein_g": round(to_float(data.get("total_protein_g")), 2),
        "total_carbohydrate_g": round(to_float(data.get("total_carbohydrate_g")), 2),
        "total_fat_g": round(to_float(data.get("total_fat_g")), 2),
        "confidence": round(to_float(data.get("confidence")), 4),
        "reason": str(data.get("reason", "")).strip(),
    }

    if estimate["total_weight_g"] <= 0:
        raise ValueError("total_weight_g <= 0")
    if estimate["recommended_portion_weight_g"] <= 0:
        raise ValueError("recommended_portion_weight_g <= 0")
    if estimate["total_energy_kcal"] <= 0:
        raise ValueError("total_energy_kcal <= 0")
    if estimate["total_protein_g"] < 0 or estimate["total_carbohydrate_g"] < 0 or estimate["total_fat_g"] < 0:
        raise ValueError("宏量营养素存在负值")
    if estimate["confidence"] < 0 or estimate["confidence"] > 1:
        raise ValueError("confidence 不在 0~1 范围内")
    if estimate["recommended_portion_weight_g"] > estimate["total_weight_g"]:
        estimate["recommended_portion_weight_g"] = estimate["total_weight_g"]

    return estimate


def normalize_portion_by_meal_type(meal_type, portion_weight_g, total_weight_g):
    ranges = {
        "早餐": (180.0, 350.0),
        "午餐": (250.0, 450.0),
        "晚餐": (250.0, 450.0),
        "加餐": (80.0, 220.0),
    }
    lower, upper = ranges.get(str(meal_type).strip(), (150.0, 400.0))
    normalized = max(lower, min(portion_weight_g, upper))
    return min(normalized, total_weight_g)


def convert_to_recipe_nutrition(recipe, estimate):
    total_weight = estimate["total_weight_g"]
    portion_weight = normalize_portion_by_meal_type(
        recipe.get("meal_type", ""),
        estimate["recommended_portion_weight_g"],
        total_weight,
    )
    return {
        "portion_weight_g": round(portion_weight, 2),
        "energy": round(estimate["total_energy_kcal"] / total_weight * 100, 2),
        "protein": round(estimate["total_protein_g"] / total_weight * 100, 2),
        "carbohydrate": round(estimate["total_carbohydrate_g"] / total_weight * 100, 2),
        "fat": round(estimate["total_fat_g"] / total_weight * 100, 2),
    }


def classify_recipe(recipe, nutrition):
    name = str(recipe.get("name", "")).strip().lower()
    ingredients = " ".join(parse_json_list(recipe.get("ingredients"))).lower()
    content = f"{name} {ingredients}"

    high_sugar_drink_keywords = [
        "奶茶", "莫吉托", "mojito", "长岛冰茶", "鸡尾酒", "特调", "果茶", "水果茶",
        "冰茶", "奶昔", "糖水", "杨枝甘露", "冰淇淋", "奶冻", "提拉米苏", "雪花酥", "布丁",
    ]
    alcohol_keywords = [
        "酒", "威士忌", "伏特加", "朗姆", "金酒", "利口酒", "啤酒", "红酒", "白酒", "黄酒",
    ]
    fried_heavy_keywords = [
        "炸", "酥", "锅包", "红烧肉", "扣肉", "肥肠", "五花肉", "奶油", "黄油",
    ]
    clean_protein_keywords = [
        "鸡胸", "鸡胸肉", "牛肉", "牛腱", "虾仁", "鱼", "三文鱼", "鸡蛋", "豆腐", "牛奶", "酸奶",
    ]
    staple_keywords = [
        "饭", "面", "粉", "粥", "燕麦", "馒头", "包子", "饺子", "馄饨", "米线", "意面",
    ]

    energy = float(nutrition["energy"])
    protein = float(nutrition["protein"])
    carbohydrate = float(nutrition["carbohydrate"])
    fat = float(nutrition["fat"])

    is_general = True
    if contains_any(content, alcohol_keywords):
        is_general = False

    is_weight_loss = (
        is_general
        and not contains_any(content, high_sugar_drink_keywords + alcohol_keywords + fried_heavy_keywords)
        and energy <= 220
        and fat <= 12
        and carbohydrate <= 25
    )

    is_sugar_control = (
        is_general
        and not contains_any(content, high_sugar_drink_keywords + alcohol_keywords)
        and carbohydrate <= 20
        and energy <= 230
    )

    is_muscle_gain = (
        is_general
        and energy >= 110
        and protein >= 8
        and (
            protein >= 12
            or contains_any(content, clean_protein_keywords)
            or contains_any(content, staple_keywords)
        )
    )

    return {
        "is_weight_loss_friendly": bool(is_weight_loss),
        "is_muscle_gain_friendly": bool(is_muscle_gain),
        "is_sugar_control_friendly": bool(is_sugar_control),
        "is_general_friendly": bool(is_general),
    }


def update_recipe(recipe_id, nutrition, suitability):
    sql = """
    UPDATE recipes
    SET
        portion_weight_g = %s,
        energy = %s,
        protein = %s,
        carbohydrate = %s,
        fat = %s,
        is_weight_loss_friendly = %s,
        is_muscle_gain_friendly = %s,
        is_sugar_control_friendly = %s,
        is_general_friendly = %s,
        updated_at = NOW()
    WHERE id = %s
    """

    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.execute(
                sql,
                (
                    nutrition["portion_weight_g"],
                    nutrition["energy"],
                    nutrition["protein"],
                    nutrition["carbohydrate"],
                    nutrition["fat"],
                    suitability["is_weight_loss_friendly"],
                    suitability["is_muscle_gain_friendly"],
                    suitability["is_sugar_control_friendly"],
                    suitability["is_general_friendly"],
                    recipe_id,
                ),
            )
        conn.commit()
    finally:
        conn.close()


def delete_invalid_recipes(dry_run=False):
    count_sql = """
    SELECT COUNT(*) AS count
    FROM recipes
    WHERE deleted_at IS NULL
      AND (
        portion_weight_g IS NULL OR portion_weight_g <= 0
        OR energy IS NULL OR energy <= 0
      )
    """

    delete_sql = """
    UPDATE recipes
    SET
        deleted_at = NOW(),
        updated_at = NOW()
    WHERE deleted_at IS NULL
      AND (
        portion_weight_g IS NULL OR portion_weight_g <= 0
        OR energy IS NULL OR energy <= 0
      )
    """

    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.execute(count_sql)
            row = cursor.fetchone() or {}
            count = int(row.get("count", 0))
            if dry_run or count == 0:
                return count

            cursor.execute(delete_sql)
        conn.commit()
        return count
    finally:
        conn.close()


def process_recipe(client, recipe, args):
    last_error = None

    for _ in range(args.retry):
        try:
            raw = request_estimate(client, recipe)
            estimate = normalize_estimate(raw)
            estimate = apply_reasonableness_rules(recipe, estimate)
            if estimate["confidence"] < args.min_confidence:
                return {
                    "ok": False,
                    "recipe_id": recipe["id"],
                    "name": recipe["name"],
                    "reason": f"模型置信度过低: {estimate['confidence']}",
                    "estimate": estimate,
                }

            nutrition = convert_to_recipe_nutrition(recipe, estimate)
            suitability = classify_recipe(recipe, nutrition)
            if not args.dry_run:
                update_recipe(recipe["id"], nutrition, suitability)

            return {
                "ok": True,
                "recipe_id": recipe["id"],
                "name": recipe["name"],
                "nutrition": nutrition,
                "estimate": estimate,
                "suitability": suitability,
            }
        except Exception as exc:
            last_error = exc
            time.sleep(1)

    return {
        "ok": False,
        "recipe_id": recipe["id"],
        "name": recipe["name"],
        "reason": str(last_error) if last_error else "未知错误",
    }


def main():
    args = parse_args()
    client = create_client()
    recipes = load_recipes(args)

    print("=== AI 营养回填工具 ===")
    print(f"待处理 recipes: {len(recipes)}")
    if not recipes:
        print("没有需要处理的 recipes。")
        return

    success = 0
    skipped = 0

    with tqdm(total=len(recipes)) as pbar:
        with ThreadPoolExecutor(max_workers=max(1, args.workers)) as executor:
            futures = [executor.submit(process_recipe, client, recipe, args) for recipe in recipes]

            for future in as_completed(futures):
                result = future.result()
                if result["ok"]:
                    nutrition = result["nutrition"]
                    estimate = result["estimate"]
                    suitability = result["suitability"]
                    action = "估算" if args.dry_run else "回填"
                    tqdm.write(
                        f"✅ [{result['recipe_id']}] {result['name']} {action}成功 | "
                        f"推荐份量={nutrition['portion_weight_g']}g | "
                        f"每100g={nutrition['energy']}kcal P{nutrition['protein']} "
                        f"C{nutrition['carbohydrate']} F{nutrition['fat']} | "
                        f"WL={int(suitability['is_weight_loss_friendly'])} "
                        f"MG={int(suitability['is_muscle_gain_friendly'])} "
                        f"SC={int(suitability['is_sugar_control_friendly'])} "
                        f"GEN={int(suitability['is_general_friendly'])} | "
                        f"conf={estimate['confidence']}"
                    )
                    success += 1
                else:
                    tqdm.write(
                        f"⚠️ [{result['recipe_id']}] {result['name']} 跳过: {result['reason']}"
                    )
                    skipped += 1
                pbar.update(1)

    deleted = 0
    if args.keep_invalid:
        print("\n已跳过无营养数据清理。")
    else:
        deleted = delete_invalid_recipes(dry_run=args.dry_run)
        action = "将删除" if args.dry_run else "已删除"
        print(f"{action}仍无营养数据的菜谱: {deleted}")

    print(f"\n完成: 成功={success}, 跳过={skipped}, 删除无效={deleted}, dry_run={args.dry_run}")


if __name__ == "__main__":
    main()
