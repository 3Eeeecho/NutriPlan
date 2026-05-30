import argparse
import json
import os
import random
import re
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import dataclass
from typing import Any

import pymysql
import pymysql.cursors
from openai import OpenAI
from tqdm import tqdm

from classify_recipe_suitability import classify_allowed_meal_types, get_recipe_columns


DB_CONFIG = {
    "host": os.getenv("NUTRIPLAN_DB_HOST", "127.0.0.1"),
    "user": os.getenv("NUTRIPLAN_DB_USER", "root"),
    "password": os.getenv("NUTRIPLAN_DB_PASSWORD", "root"),
    "database": os.getenv("NUTRIPLAN_DB_NAME", "nutri_plan"),
    "charset": "utf8mb4",
    "cursorclass": pymysql.cursors.DictCursor,
}

BASE_URL = os.getenv("NUTRIPLAN_LLM_BASE_URL", "https://api.siliconflow.cn/v1")
MODEL_NAME = os.getenv("NUTRIPLAN_LLM_MODEL", "deepseek-ai/DeepSeek-V4-Flash")
API_KEY = os.getenv("NUTRIPLAN_LLM_API_KEY", "sk-hvzxenbdsosgapxauxlnehlgkigyvgmhjaesvulfnwcxeqjp")
LLM_TIMEOUT = float(os.getenv("NUTRIPLAN_LLM_TIMEOUT", "60"))

MEAL_TYPES = ("早餐", "午餐", "晚餐", "加餐")
HEALTH_GOALS = ("减脂", "增肌", "控糖", "维持健康")

GOAL_FLAG_COLUMNS = {
    "减脂": "is_weight_loss_friendly",
    "增肌": "is_muscle_gain_friendly",
    "控糖": "is_sugar_control_friendly",
    "维持健康": "is_general_friendly",
}

GOAL_TARGET_USERS = {
    "减脂": "减脂",
    "增肌": "增肌",
    "控糖": "控糖",
    "维持健康": "维持健康",
}

DEFAULT_TARGET_PER_BUCKET = {
    "早餐": 50,
    "午餐": 80,
    "晚餐": 80,
    "加餐": 45,
}

PORTION_RANGES = {
    "早餐": (180.0, 350.0),
    "午餐": (120.0, 320.0),
    "晚餐": (120.0, 320.0),
    "加餐": (60.0, 200.0),
}

# Single recipe serving energy. Lunch and dinner recipes are intentionally dish-sized
# because the recommendation engine combines 2-3 recipes into one meal.
SERVING_ENERGY_RANGES = {
    ("早餐", "减脂"): (180, 360),
    ("午餐", "减脂"): (120, 320),
    ("晚餐", "减脂"): (120, 300),
    ("加餐", "减脂"): (70, 180),
    ("早餐", "增肌"): (250, 460),
    ("午餐", "增肌"): (180, 420),
    ("晚餐", "增肌"): (160, 400),
    ("加餐", "增肌"): (120, 260),
    ("早餐", "控糖"): (180, 380),
    ("午餐", "控糖"): (120, 340),
    ("晚餐", "控糖"): (120, 320),
    ("加餐", "控糖"): (70, 180),
    ("早餐", "维持健康"): (220, 420),
    ("午餐", "维持健康"): (150, 380),
    ("晚餐", "维持健康"): (140, 360),
    ("加餐", "维持健康"): (80, 220),
}

PER_100G_NUTRITION_LIMITS = {
    "energy": (10.0, 380.0),
    "snack_energy": (10.0, 500.0),
    "protein": (0.0, 40.0),
    "carbohydrate": (0.0, 80.0),
    "fat": (0.0, 40.0),
    "snack_fat": (0.0, 55.0),
}

CUISINE_THEMES = [
    "家常中式", "粤式清淡", "川渝轻辣", "江浙清鲜", "东北杂粮", "西北粗粮",
    "日式定食", "韩式均衡", "地中海", "墨西哥轻食", "东南亚香草", "法式简餐",
    "中东香料", "印度低油", "素食高蛋白", "海鲜低脂", "禽肉高蛋白", "牛肉铁元素",
]

PROTEIN_ROTATION = [
    "鸡胸肉", "鸡腿肉", "火鸡胸", "瘦牛肉", "牛腱", "猪里脊", "鳕鱼", "三文鱼",
    "金枪鱼", "虾仁", "扇贝", "鸡蛋", "豆腐", "豆干", "鹰嘴豆", "毛豆", "低脂奶酪",
]

STAPLE_ROTATION = [
    "糙米", "燕麦", "藜麦", "荞麦面", "全麦意面", "玉米", "红薯", "紫薯",
    "小米", "黑米", "鹰嘴豆泥", "全麦卷饼", "杂粮馒头", "魔芋丝",
]

VEG_ROTATION = [
    "西兰花", "菠菜", "羽衣甘蓝", "芦笋", "秋葵", "番茄", "彩椒", "南瓜",
    "蘑菇", "生菜", "紫甘蓝", "黄瓜", "胡萝卜", "荷兰豆", "菜心", "茄子",
]

ALLERGY_KEYWORDS = {
    "虾": "海鲜过敏",
    "鱼": "海鲜过敏",
    "贝": "海鲜过敏",
    "蟹": "海鲜过敏",
    "奶": "乳制品过敏",
    "酸奶": "乳制品过敏",
    "奶酪": "乳制品过敏",
    "花生": "坚果过敏",
    "腰果": "坚果过敏",
    "杏仁": "坚果过敏",
    "核桃": "坚果过敏",
    "全麦": "麸质过敏",
    "意面": "麸质过敏",
    "荞麦": "麸质过敏",
}

CHINESE_DISH_NAME_EXAMPLES = (
    "宫保鸡丁", "鱼香肉丝", "番茄牛腩", "香菇滑鸡", "青椒牛柳", "虾仁滑蛋", "葱爆羊肉", "芹菜香干",
    "清蒸鲈鱼", "白灼虾", "西芹百合", "木须肉", "麻婆豆腐", "蚝油生菜", "荷塘小炒", "冬瓜丸子汤",
)

CHINESE_BREAKFAST_NAME_EXAMPLES = (
    "小米南瓜粥", "杂粮燕麦粥", "豆浆鸡蛋饼", "鲜肉小笼包", "香菇菜包", "玉米馒头",
    "虾仁蒸蛋", "茶叶蛋", "豆腐脑", "鸡丝汤面", "鲜虾馄饨", "牛奶燕麦杯",
)

NAME_BANNED_WORDS = (
    "配", "搭配", "组合", "套餐", "拼盘", "佐", "拌饭", "便当", "营养餐", "健身餐", "减脂餐", "增肌餐",
)

NAME_BANNED_PATTERNS = (
    re.compile(r".+配.+"),
    re.compile(r".+搭配.+"),
    re.compile(r".+佐.+"),
)


@dataclass(frozen=True)
class GenerationBucket:
    meal_type: str
    goal: str
    missing: int


def get_db_connection():
    return pymysql.connect(**DB_CONFIG)


def create_client():
    if not API_KEY:
        raise RuntimeError("缺少环境变量 NUTRIPLAN_LLM_API_KEY")
    return OpenAI(api_key=API_KEY, base_url=BASE_URL, timeout=LLM_TIMEOUT)


def parse_args():
    parser = argparse.ArgumentParser(
        description="Generate NutriPlan recipes with AI, or backfill nutrition for existing recipes."
    )
    parser.add_argument("--task", choices=("generate", "backfill", "coverage"), default="generate")
    parser.add_argument("--dry-run", action="store_true", help="Preview changes without writing to DB")
    parser.add_argument("--batch-size", type=int, default=8, help="Recipes requested per AI call")
    parser.add_argument("--max-new", type=int, default=520, help="Maximum generated recipes in one run")
    parser.add_argument("--target-breakfast", type=int, default=DEFAULT_TARGET_PER_BUCKET["早餐"])
    parser.add_argument("--target-lunch", type=int, default=DEFAULT_TARGET_PER_BUCKET["午餐"])
    parser.add_argument("--target-dinner", type=int, default=DEFAULT_TARGET_PER_BUCKET["晚餐"])
    parser.add_argument("--target-snack", type=int, default=DEFAULT_TARGET_PER_BUCKET["加餐"])
    parser.add_argument("--meal-type", choices=MEAL_TYPES, default="", help="Only generate this meal type")
    parser.add_argument("--goal", choices=HEALTH_GOALS, default="", help="Only generate this health goal")
    parser.add_argument("--retry", type=int, default=3)
    parser.add_argument("--sleep", type=float, default=0.8, help="Seconds to sleep between AI calls")
    parser.add_argument("--workers", type=int, default=3, help="Backfill workers")
    parser.add_argument("--limit", type=int, default=0, help="Backfill row limit")
    parser.add_argument("--recipe-id", type=int, default=0, help="Backfill one recipe id")
    parser.add_argument("--all", action="store_true", help="Backfill all rows instead of missing rows")
    parser.add_argument("--min-confidence", type=float, default=0.7)
    parser.add_argument("--keep-invalid", action="store_true")
    parser.add_argument("--seed", type=int, default=int(time.time()))
    return parser.parse_args()


def json_dumps(value: Any) -> str:
    return json.dumps(value, ensure_ascii=False, separators=(",", ":"))


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
    return json.loads(match.group())


def to_float(value) -> float:
    if value is None:
        return 0.0
    if isinstance(value, (int, float)):
        return float(value)
    text = str(value).strip().replace(",", "")
    match = re.search(r"-?\d+(?:\.\d+)?", text)
    if not match:
        raise ValueError(f"cannot parse number: {value}")
    return float(match.group())


def to_int(value) -> int:
    return int(round(to_float(value)))


def normalize_text(value: str) -> str:
    text = str(value or "").lower()
    text = re.sub(r"[\s·•\-_/（）()【】\[\]：:，,。.!！?？]+", "", text)
    return text


def recipe_signature(recipe: dict) -> str:
    ingredients = parse_json_list(recipe.get("ingredients"))
    cleaned = []
    for item in ingredients[:5]:
        text = re.sub(r"\d+(?:\.\d+)?\s*(g|克|ml|毫升|个|只|片|勺|份)?", "", str(item), flags=re.I)
        text = normalize_text(text)
        if text:
            cleaned.append(text)
    return "|".join(sorted(cleaned))


def get_existing_recipe_index():
    sql = """
    SELECT name, ingredients
    FROM recipes
    WHERE deleted_at IS NULL
    """
    names = set()
    signatures = set()
    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.execute(sql)
            for row in cursor.fetchall():
                names.add(normalize_text(row.get("name", "")))
                sig = recipe_signature(row)
                if sig:
                    signatures.add(sig)
    finally:
        conn.close()
    return names, signatures


def get_coverage(columns: set[str]) -> dict[tuple[str, str], int]:
    coverage = {}
    has_allowed_meal_types = "allowed_meal_types" in columns
    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            for meal_type in MEAL_TYPES:
                for goal in HEALTH_GOALS:
                    flag = GOAL_FLAG_COLUMNS[goal]
                    if flag not in columns:
                        coverage[(meal_type, goal)] = 0
                        continue
                    meal_filter = "meal_type = %s"
                    params = [meal_type]
                    if has_allowed_meal_types:
                        meal_filter = (
                            "(JSON_CONTAINS(COALESCE(NULLIF(allowed_meal_types, ''), '[]'), JSON_QUOTE(%s)) "
                            "OR ((allowed_meal_types IS NULL OR allowed_meal_types = '' OR allowed_meal_types = '[]') "
                            "AND meal_type = %s))"
                        )
                        params = [meal_type, meal_type]
                    sql = f"""
                    SELECT COUNT(*) AS count
                    FROM recipes
                    WHERE deleted_at IS NULL
                      AND energy > 0
                      AND portion_weight_g > 0
                      AND ({flag} = 1)
                      AND {meal_filter}
                    """
                    cursor.execute(sql, params)
                    coverage[(meal_type, goal)] = int((cursor.fetchone() or {}).get("count", 0))
    finally:
        conn.close()
    return coverage


def target_counts(args) -> dict[str, int]:
    return {
        "早餐": args.target_breakfast,
        "午餐": args.target_lunch,
        "晚餐": args.target_dinner,
        "加餐": args.target_snack,
    }


def build_generation_plan(args, columns: set[str]) -> list[GenerationBucket]:
    coverage = get_coverage(columns)
    targets = target_counts(args)
    buckets = []
    for meal_type in MEAL_TYPES:
        if args.meal_type and meal_type != args.meal_type:
            continue
        for goal in HEALTH_GOALS:
            if args.goal and goal != args.goal:
                continue
            target = targets[meal_type]
            current = coverage.get((meal_type, goal), 0)
            if current < target:
                buckets.append(GenerationBucket(meal_type, goal, target - current))
    buckets.sort(key=lambda item: item.missing, reverse=True)
    return buckets


def print_coverage(args, columns: set[str]):
    coverage = get_coverage(columns)
    targets = target_counts(args)
    print("=== recipe coverage by meal and health goal ===")
    for meal_type in MEAL_TYPES:
        line = [meal_type]
        for goal in HEALTH_GOALS:
            current = coverage.get((meal_type, goal), 0)
            line.append(f"{goal}:{current}/{targets[meal_type]}")
        print(" | ".join(line))


def choose_prompt_seeds(meal_type: str, goal: str, batch_size: int) -> dict[str, list[str]]:
    random.shuffle(CUISINE_THEMES)
    random.shuffle(PROTEIN_ROTATION)
    random.shuffle(STAPLE_ROTATION)
    random.shuffle(VEG_ROTATION)
    return {
        "cuisines": CUISINE_THEMES[: max(4, min(batch_size, 8))],
        "proteins": PROTEIN_ROTATION[: max(5, min(batch_size + 2, 10))],
        "staples": STAPLE_ROTATION[: max(4, min(batch_size + 1, 9))],
        "vegetables": VEG_ROTATION[: max(5, min(batch_size + 2, 10))],
    }


def build_generation_prompt(meal_type: str, goal: str, count: int, existing_names: set[str]) -> str:
    seeds = choose_prompt_seeds(meal_type, goal, count)
    low, high = SERVING_ENERGY_RANGES[(meal_type, goal)]
    portion_low, portion_high = PORTION_RANGES[meal_type]
    sampled_existing = list(existing_names)
    random.shuffle(sampled_existing)
    sampled_existing = sampled_existing[:80]

    goal_rules = {
        "减脂": "高蛋白、低油、低能量密度；避免油炸、奶油、含糖饮料；加餐碳水不超过18g/100g。",
        "增肌": "高蛋白、主食充足、优质脂肪适量；每份蛋白质优先达到25g以上，早餐/加餐也要有明确蛋白来源。",
        "控糖": "低GI主食、蔬菜充足、避免精制糖和甜饮；每100g碳水早餐/正餐尽量不超过18g，加餐尽量不超过12g。",
        "维持健康": "营养均衡、适合日常长期执行；蛋白质、蔬菜、主食搭配完整，油盐适中。",
    }[goal]
    meal_rules = {
        "早餐": (
            "必须符合普通中国早餐习惯，只生成粥、豆浆、豆腐脑、牛奶、酸奶、燕麦、包子、馒头、花卷、"
            "烧麦、鸡蛋饼、蛋饼、煎饼、吐司、三明治、水煮蛋、煎蛋、蒸蛋、茶叶蛋、馄饨、汤面、米粉等早餐型食物。"
            "禁止把家常热炒菜放到早餐，例如彩椒毛豆炒蛋、青椒牛柳、宫保鸡丁、鱼香肉丝、麻婆豆腐、清炒时蔬、白灼虾、清蒸鱼等。"
            "早餐蛋类只能是水煮蛋、煎蛋、蒸蛋、茶叶蛋、蛋饼/鸡蛋饼，不要生成“某某炒蛋/炒鸡蛋”。"
        ),
        "午餐": "必须符合中国午餐习惯，可生成荤素菜、豆腐菜、鱼虾、汤、主食类菜；allowed_meal_types 通常应为 [\"午餐\", \"晚餐\"]。",
        "晚餐": "必须符合中国晚餐习惯，可生成清淡荤素菜、豆腐菜、鱼虾、汤、少油主食类菜；allowed_meal_types 通常应为 [\"午餐\", \"晚餐\"]。",
        "加餐": "必须是加餐习惯食物，如酸奶、水果、坚果、豆浆、牛奶、燕麦杯、蛋白棒、小份蒸蛋等；不要生成正餐热炒菜。",
    }[meal_type]
    dish_examples = CHINESE_BREAKFAST_NAME_EXAMPLES if meal_type == "早餐" else CHINESE_DISH_NAME_EXAMPLES

    return f"""
你是 NutriPlan 的专业食谱数据生成器。请生成 {count} 条全新的、可直接入库的食谱 JSON 数据。

硬性目标：
- 餐次必须是：{meal_type}
- 主要目标人群必须覆盖：{goal}
- 每份热量必须在 {low}-{high} kcal 之间
- 推荐份量 portion_weight_g 必须在 {portion_low}-{portion_high} g 之间
- 字段 energy/protein/carbohydrate/fat 必须是“每100g营养值”，不是整份营养
- 必须按公式自检：serving_energy_kcal = energy * portion_weight_g / 100，且 serving_energy_kcal 必须落在 {low}-{high} kcal
- 不要把整份热量直接填到 energy；例如一份 250g、整份 300kcal 时，energy 应填 120，而不是 300
- ingredients 必须写明每种食材的克数或毫升数，例如“鸡胸肉 120g”
- cooking_steps 必须可执行，3-6步
- 每条食谱必须有独特菜名，避免和已有名称、常见模板高度相似
- 菜名必须像当前餐次中常见的短菜名，优先使用这些命名风格：{", ".join(dish_examples)}
- 菜名长度控制在 2-8 个汉字，不要写成长标题，不要加入“高蛋白、低脂、减脂、增肌、营养、健康”等功能描述
- 菜名禁止出现“配、搭配、佐、组合、套餐、拼盘”等连接词，尤其禁止“xx配xx”“xx搭配xx”的字样
- 不要生成酒精饮品、奶茶、甜品、油炸高糖食谱
- 不要生成“测试-”开头的名称
- 不要输出 markdown，只输出单个 JSON 对象

餐次习惯规则：
{meal_rules}

目标规则：
{goal_rules}

多样性要求：
- 菜系可参考：{", ".join(seeds["cuisines"])}
- 常见菜名风格可参考：{", ".join(dish_examples)}
- 蛋白来源可参考：{", ".join(seeds["proteins"])}
- 主食可参考：{", ".join(seeds["staples"])}
- 蔬菜可参考：{", ".join(seeds["vegetables"])}
- 同一批内不要重复主要蛋白、主食和烹饪方式

已有名称归一化列表，禁止重复或近似：
{json_dumps(sampled_existing)}

输出格式：
{{
  "recipes": [
    {{
      "name": "菜名",
      "meal_type": "{meal_type}",
      "allowed_meal_types": ["{meal_type}"],
      "difficulty": "简单|中等|困难",
      "cooking_time": 15,
      "portion_weight_g": 300,
      "serving_energy_kcal": 360,
      "energy": 120,
      "protein": 12,
      "carbohydrate": 15,
      "fat": 4,
      "ingredients": ["食材 100g"],
      "cooking_steps": ["步骤1"],
      "target_users": ["{goal}", "维持健康"],
      "forbidden_users": ["海鲜过敏"],
      "is_weight_loss_friendly": false,
      "is_muscle_gain_friendly": false,
      "is_sugar_control_friendly": false,
      "is_general_friendly": true
    }}
  ]
}}
""".strip()


def request_json(client, prompt: str) -> dict:
    payload = {
        "model": MODEL_NAME,
        "messages": [{"role": "user", "content": prompt}],
        "temperature": 0.75,
        "max_tokens": 6000,
    }
    try:
        response = client.chat.completions.create(**payload, response_format={"type": "json_object"})
    except Exception as exc:
        error_text = str(exc).lower()
        if "response_format" not in error_text and "json_object" not in error_text:
            raise
        response = client.chat.completions.create(**payload)
    return extract_json_from_text(response.choices[0].message.content)


def normalize_bool(value) -> bool:
    if isinstance(value, bool):
        return value
    if isinstance(value, (int, float)):
        return value != 0
    return str(value).strip().lower() in {"1", "true", "yes", "是"}


def infer_forbidden_users(recipe: dict) -> list[str]:
    text = f"{recipe.get('name', '')} {' '.join(parse_json_list(recipe.get('ingredients')))}"
    forbidden = []
    for keyword, tag in ALLERGY_KEYWORDS.items():
        if keyword in text and tag not in forbidden:
            forbidden.append(tag)
    return forbidden


def validate_recipe_name(name: str):
    if not name or name.startswith("测试-"):
        raise ValueError("invalid recipe name")
    if len(name) > 8:
        raise ValueError(f"recipe name too long: {name}")
    if any(word in name for word in NAME_BANNED_WORDS):
        raise ValueError(f"recipe name contains banned word: {name}")
    if any(pattern.search(name) for pattern in NAME_BANNED_PATTERNS):
        raise ValueError(f"recipe name uses banned pairing pattern: {name}")
    if re.search(r"[，,、/＋+&和与]", name):
        raise ValueError(f"recipe name looks like a combined title: {name}")


def maybe_convert_serving_nutrition_to_per_100g(recipe: dict, meal_type: str, goal: str):
    factor = recipe["portion_weight_g"] / 100.0
    if factor <= 1:
        return

    energy_low, energy_high = SERVING_ENERGY_RANGES[(meal_type, goal)]
    calculated_serving_energy = recipe["energy"] * factor
    energy_looks_like_serving_total = energy_low * 0.85 <= recipe["energy"] <= energy_high * 1.15
    calculated_is_too_high = calculated_serving_energy > energy_high * 1.25
    if not (energy_looks_like_serving_total and calculated_is_too_high):
        return

    for field in ("energy", "protein", "carbohydrate", "fat"):
        recipe[field] = round(recipe[field] / factor, 2)


def validate_per_100g_nutrition(recipe: dict, meal_type: str):
    energy_min, energy_max = PER_100G_NUTRITION_LIMITS["snack_energy" if meal_type == "加餐" else "energy"]
    fat_min, fat_max = PER_100G_NUTRITION_LIMITS["snack_fat" if meal_type == "加餐" else "fat"]
    protein_min, protein_max = PER_100G_NUTRITION_LIMITS["protein"]
    carb_min, carb_max = PER_100G_NUTRITION_LIMITS["carbohydrate"]

    if not (energy_min <= recipe["energy"] <= energy_max):
        raise ValueError(f"energy per 100g out of plausible range: {recipe['energy']}")
    if not (protein_min <= recipe["protein"] <= protein_max):
        raise ValueError(f"protein per 100g out of plausible range: {recipe['protein']}")
    if not (carb_min <= recipe["carbohydrate"] <= carb_max):
        raise ValueError(f"carbohydrate per 100g out of plausible range: {recipe['carbohydrate']}")
    if not (fat_min <= recipe["fat"] <= fat_max):
        raise ValueError(f"fat per 100g out of plausible range: {recipe['fat']}")

    macro_energy = recipe["protein"] * 4.0 + recipe["carbohydrate"] * 4.0 + recipe["fat"] * 9.0
    if macro_energy > 0 and recipe["energy"] < macro_energy * 0.65:
        raise ValueError(f"energy too low for macros: energy={recipe['energy']} macro_energy={macro_energy:.1f}")
    if macro_energy > 0 and recipe["energy"] > macro_energy * 1.6 + 30:
        raise ValueError(f"energy too high for macros: energy={recipe['energy']} macro_energy={macro_energy:.1f}")


def normalize_generated_recipe(raw: dict, meal_type: str, goal: str) -> dict:
    if not isinstance(raw, dict):
        raise ValueError("recipe is not an object")
    recipe = {
        "name": str(raw.get("name", "")).strip(),
        "image_url": str(raw.get("image_url", "") or "").strip(),
        "meal_type": str(raw.get("meal_type", meal_type)).strip(),
        "allowed_meal_types": parse_json_list(raw.get("allowed_meal_types")) or [meal_type],
        "difficulty": str(raw.get("difficulty", "中等")).strip(),
        "cooking_time": to_int(raw.get("cooking_time", 20)),
        "portion_weight_g": round(to_float(raw.get("portion_weight_g")), 2),
        "energy": round(to_float(raw.get("energy")), 2),
        "protein": round(to_float(raw.get("protein")), 2),
        "carbohydrate": round(to_float(raw.get("carbohydrate")), 2),
        "fat": round(to_float(raw.get("fat")), 2),
        "ingredients": parse_json_list(raw.get("ingredients")),
        "cooking_steps": parse_json_list(raw.get("cooking_steps")),
        "target_users": parse_json_list(raw.get("target_users")),
        "forbidden_users": parse_json_list(raw.get("forbidden_users")),
        "is_weight_loss_friendly": normalize_bool(raw.get("is_weight_loss_friendly")),
        "is_muscle_gain_friendly": normalize_bool(raw.get("is_muscle_gain_friendly")),
        "is_sugar_control_friendly": normalize_bool(raw.get("is_sugar_control_friendly")),
        "is_general_friendly": normalize_bool(raw.get("is_general_friendly", True)),
    }

    validate_recipe_name(recipe["name"])
    if recipe["meal_type"] != meal_type:
        raise ValueError(f"meal_type mismatch: {recipe['meal_type']} != {meal_type}")
    if recipe["difficulty"] not in {"简单", "中等", "困难"}:
        recipe["difficulty"] = "中等"
    if recipe["cooking_time"] <= 0 or recipe["cooking_time"] > 120:
        raise ValueError("invalid cooking_time")
    if len(recipe["ingredients"]) < 3:
        raise ValueError("ingredients too short")
    if len(recipe["cooking_steps"]) < 2:
        raise ValueError("cooking_steps too short")

    portion_low, portion_high = PORTION_RANGES[meal_type]
    if not (portion_low <= recipe["portion_weight_g"] <= portion_high):
        raise ValueError(f"portion out of range: {recipe['portion_weight_g']}")
    if recipe["energy"] <= 0 or recipe["protein"] < 0 or recipe["carbohydrate"] < 0 or recipe["fat"] < 0:
        raise ValueError("invalid nutrition")
    maybe_convert_serving_nutrition_to_per_100g(recipe, meal_type, goal)
    validate_per_100g_nutrition(recipe, meal_type)

    serving_energy = recipe["energy"] * recipe["portion_weight_g"] / 100.0
    energy_low, energy_high = SERVING_ENERGY_RANGES[(meal_type, goal)]
    if serving_energy < energy_low:
        raise ValueError(f"serving energy below range: {serving_energy:.1f} < {energy_low}")
    if serving_energy > energy_high:
        raise ValueError(f"serving energy above range: {serving_energy:.1f} > {energy_high}")

    flag = GOAL_FLAG_COLUMNS[goal]
    recipe[flag] = True
    if goal not in recipe["target_users"]:
        recipe["target_users"].insert(0, goal)
    if recipe["is_general_friendly"] and "维持健康" not in recipe["target_users"]:
        recipe["target_users"].append("维持健康")

    inferred_forbidden = infer_forbidden_users(recipe)
    for tag in inferred_forbidden:
        if tag not in recipe["forbidden_users"]:
            recipe["forbidden_users"].append(tag)

    recipe_for_slots = {
        **recipe,
        "ingredients": recipe["ingredients"],
    }
    recipe["allowed_meal_types"] = classify_allowed_meal_types(recipe_for_slots)
    if meal_type not in recipe["allowed_meal_types"]:
        raise ValueError(
            f"recipe is not suitable for requested meal_type: {recipe['name']} "
            f"requested={meal_type} allowed={recipe['allowed_meal_types']}"
        )

    return recipe


def insert_recipe(recipe: dict, columns: set[str], dry_run: bool):
    insertable = {
        "created_at": "NOW()",
        "updated_at": "NOW()",
        "name": recipe["name"],
        "image_url": recipe["image_url"],
        "meal_type": recipe["meal_type"],
        "allowed_meal_types": json_dumps(recipe["allowed_meal_types"]),
        "difficulty": recipe["difficulty"],
        "cooking_time": recipe["cooking_time"],
        "portion_weight_g": recipe["portion_weight_g"],
        "energy": recipe["energy"],
        "protein": recipe["protein"],
        "carbohydrate": recipe["carbohydrate"],
        "fat": recipe["fat"],
        "ingredients": json_dumps(recipe["ingredients"]),
        "cooking_steps": json_dumps(recipe["cooking_steps"]),
        "target_users": json_dumps(recipe["target_users"]),
        "forbidden_users": json_dumps(recipe["forbidden_users"]),
        "is_weight_loss_friendly": recipe["is_weight_loss_friendly"],
        "is_muscle_gain_friendly": recipe["is_muscle_gain_friendly"],
        "is_sugar_control_friendly": recipe["is_sugar_control_friendly"],
        "is_general_friendly": recipe["is_general_friendly"],
    }
    if "health_score" in columns:
        insertable["health_score"] = 8

    field_names = [field for field in insertable if field in columns]
    sql_fields = []
    placeholders = []
    values = []
    for field in field_names:
        sql_fields.append(f"`{field}`")
        if insertable[field] == "NOW()":
            placeholders.append("NOW()")
        else:
            placeholders.append("%s")
            values.append(insertable[field])

    sql = f"INSERT INTO recipes ({', '.join(sql_fields)}) VALUES ({', '.join(placeholders)})"
    if dry_run:
        return

    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.execute(sql, values)
        conn.commit()
    finally:
        conn.close()


def rotate_buckets(buckets: list[GenerationBucket]):
    active = [GenerationBucket(item.meal_type, item.goal, item.missing) for item in buckets if item.missing > 0]
    while active:
        next_active = []
        for item in active:
            yield item
            remaining = item.missing - 1
            if remaining > 0:
                next_active.append(GenerationBucket(item.meal_type, item.goal, remaining))
        active = next_active


def generate_recipes(args, columns: set[str]):
    print_coverage(args, columns)
    buckets = build_generation_plan(args, columns)
    if not buckets:
        print("coverage target already satisfied")
        return
    if args.max_new <= 0:
        print("max-new is 0, nothing to generate")
        return

    client = create_client()
    existing_names, existing_signatures = get_existing_recipe_index()
    generated = 0
    skipped = 0
    failed_rounds: dict[tuple[str, str], int] = {}
    max_failed_rounds_per_bucket = 3

    print("\n=== AI recipe generation ===")
    for bucket in rotate_buckets(buckets):
        if generated >= args.max_new:
            break
        bucket_key = (bucket.meal_type, bucket.goal)
        if failed_rounds.get(bucket_key, 0) >= max_failed_rounds_per_bucket:
            continue

        print(f"\nBucket {bucket.meal_type}/{bucket.goal}: remaining={bucket.missing}")
        request_count = min(args.batch_size, bucket.missing, args.max_new - generated)
        prompt = build_generation_prompt(bucket.meal_type, bucket.goal, request_count, existing_names)
        last_error = None
        data = None
        for attempt in range(args.retry):
            try:
                print(
                    f"  requesting AI batch: count={request_count}, "
                    f"attempt={attempt + 1}/{args.retry}, model={MODEL_NAME}, timeout={LLM_TIMEOUT}s",
                    flush=True,
                )
                data = request_json(client, prompt)
                print("  AI batch received, validating recipes...", flush=True)
                break
            except Exception as exc:
                last_error = exc
                print(f"  request attempt failed: {exc}", flush=True)
                time.sleep(1 + attempt)
        if data is None:
            failed_rounds[bucket_key] = failed_rounds.get(bucket_key, 0) + 1
            print(f"  request failed: {last_error}")
            continue

        raw_recipes = data.get("recipes") if isinstance(data, dict) else None
        if not isinstance(raw_recipes, list):
            failed_rounds[bucket_key] = failed_rounds.get(bucket_key, 0) + 1
            print("  response skipped: recipes is not a list")
            skipped += 1
            continue

        accepted_this_round = 0
        for raw in raw_recipes:
            try:
                recipe = normalize_generated_recipe(raw, bucket.meal_type, bucket.goal)
                name_key = normalize_text(recipe["name"])
                sig = recipe_signature(recipe)
                if name_key in existing_names:
                    raise ValueError("duplicate name")
                if sig and sig in existing_signatures:
                    raise ValueError("duplicate ingredient signature")
                insert_recipe(recipe, columns, args.dry_run)
                existing_names.add(name_key)
                if sig:
                    existing_signatures.add(sig)
                generated += 1
                accepted_this_round += 1
                serving_energy = recipe["energy"] * recipe["portion_weight_g"] / 100.0
                action = "would insert" if args.dry_run else "inserted"
                print(
                    f"  {action}: {recipe['name']} | {bucket.meal_type}/{bucket.goal} | "
                    f"{serving_energy:.0f}kcal/serving | P{recipe['protein']} C{recipe['carbohydrate']} F{recipe['fat']}"
                )
                if generated >= args.max_new:
                    break
            except Exception as exc:
                skipped += 1
                print(f"  skipped generated item: {exc}")

        if accepted_this_round == 0:
            failed_rounds[bucket_key] = failed_rounds.get(bucket_key, 0) + 1
            print(
                f"  no accepted recipe in this AI batch; "
                f"failed_rounds={failed_rounds[bucket_key]}/{max_failed_rounds_per_bucket}"
            )
        else:
            failed_rounds[bucket_key] = 0
        time.sleep(max(0, args.sleep))

    print(f"\nfinished generation: generated={generated}, skipped={skipped}, dry_run={args.dry_run}")
    print_coverage(args, columns)


def load_recipes_for_backfill(args):
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


def build_backfill_prompt(recipe):
    ingredients = parse_json_list(recipe.get("ingredients"))
    steps = parse_json_list(recipe.get("cooking_steps"))
    ingredients_text = "\n".join(f"- {item}" for item in ingredients) or "- 无"
    steps_text = "\n".join(f"{idx + 1}. {item}" for idx, item in enumerate(steps)) or "无"
    return f"""
你是注册营养师。请只对现有菜谱做营养估算，不要改写菜名，不要生成新菜。
输出必须是单个 JSON 对象，不要 markdown。

菜谱：
name: {recipe.get("name", "")}
meal_type: {recipe.get("meal_type", "")}
difficulty: {recipe.get("difficulty", "")}
cooking_time: {recipe.get("cooking_time", 0)}

ingredients:
{ingredients_text}

steps:
{steps_text}

请估算：
- total_weight_g 整道菜最终可食用重量
- recommended_portion_weight_g 单人推荐份量
- total_energy_kcal 整道菜总热量
- total_protein_g 整道菜总蛋白质
- total_carbohydrate_g 整道菜总碳水
- total_fat_g 整道菜总脂肪
- confidence 0到1
- reason 简短依据

JSON 格式：
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


def request_backfill_estimate(client, recipe):
    return request_json(client, build_backfill_prompt(recipe))


def normalize_estimate(data):
    if not isinstance(data, dict):
        raise ValueError("model did not return a JSON object")
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
    if estimate["recommended_portion_weight_g"] > estimate["total_weight_g"]:
        estimate["recommended_portion_weight_g"] = estimate["total_weight_g"]
    if estimate["confidence"] < 0 or estimate["confidence"] > 1:
        raise ValueError("confidence out of range")
    return estimate


def normalize_portion_by_meal_type(meal_type, portion_weight_g, total_weight_g):
    lower, upper = PORTION_RANGES.get(str(meal_type).strip(), (150.0, 400.0))
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
    meal_type = str(recipe.get("meal_type", ""))
    serving_energy = nutrition["energy"] * nutrition["portion_weight_g"] / 100.0
    protein = nutrition["protein"]
    carbohydrate = nutrition["carbohydrate"]
    fat = nutrition["fat"]
    return {
        "is_weight_loss_friendly": serving_energy <= SERVING_ENERGY_RANGES.get((meal_type, "减脂"), (0, 999))[1] and fat <= 12,
        "is_muscle_gain_friendly": protein >= 8 and serving_energy >= SERVING_ENERGY_RANGES.get((meal_type, "增肌"), (0, 999))[0] * 0.65,
        "is_sugar_control_friendly": carbohydrate <= 20,
        "is_general_friendly": True,
    }


def build_target_users(suitability):
    users = []
    for goal, flag in GOAL_FLAG_COLUMNS.items():
        if suitability.get(flag):
            users.append(GOAL_TARGET_USERS[goal])
    return users or ["维持健康"]


def update_recipe(recipe_id, nutrition, suitability, target_users, allowed_meal_types, columns: set[str], dry_run=False):
    updates = {
        "portion_weight_g": nutrition["portion_weight_g"],
        "energy": nutrition["energy"],
        "protein": nutrition["protein"],
        "carbohydrate": nutrition["carbohydrate"],
        "fat": nutrition["fat"],
        "target_users": json_dumps(target_users),
        "is_weight_loss_friendly": suitability["is_weight_loss_friendly"],
        "is_muscle_gain_friendly": suitability["is_muscle_gain_friendly"],
        "is_sugar_control_friendly": suitability["is_sugar_control_friendly"],
        "is_general_friendly": suitability["is_general_friendly"],
        "allowed_meal_types": json_dumps(allowed_meal_types),
    }
    set_clauses = []
    params = []
    for field, value in updates.items():
        if field in columns:
            set_clauses.append(f"`{field}` = %s")
            params.append(value)
    set_clauses.append("updated_at = NOW()")
    params.append(recipe_id)
    sql = f"UPDATE recipes SET {', '.join(set_clauses)} WHERE id = %s"
    if dry_run:
        return
    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.execute(sql, params)
        conn.commit()
    finally:
        conn.close()


def delete_invalid_recipes(dry_run=False):
    sql = """
    SELECT COUNT(*) AS count
    FROM recipes
    WHERE deleted_at IS NULL
      AND (portion_weight_g IS NULL OR portion_weight_g <= 0 OR energy IS NULL OR energy <= 0)
    """
    delete_sql = """
    UPDATE recipes
    SET deleted_at = NOW(), updated_at = NOW()
    WHERE deleted_at IS NULL
      AND (portion_weight_g IS NULL OR portion_weight_g <= 0 OR energy IS NULL OR energy <= 0)
    """
    conn = get_db_connection()
    try:
        with conn.cursor() as cursor:
            cursor.execute(sql)
            count = int((cursor.fetchone() or {}).get("count", 0))
            if not dry_run and count > 0:
                cursor.execute(delete_sql)
        conn.commit()
        return count
    finally:
        conn.close()


def process_backfill_recipe(client, recipe, args, columns: set[str]):
    last_error = None
    for _ in range(args.retry):
        try:
            estimate = normalize_estimate(request_backfill_estimate(client, recipe))
            if estimate["confidence"] < args.min_confidence:
                return {"ok": False, "recipe": recipe, "reason": f"low confidence {estimate['confidence']}"}
            nutrition = convert_to_recipe_nutrition(recipe, estimate)
            suitability = classify_recipe(recipe, nutrition)
            target_users = build_target_users(suitability)
            allowed_meal_types = classify_allowed_meal_types({**recipe, **nutrition})
            update_recipe(recipe["id"], nutrition, suitability, target_users, allowed_meal_types, columns, args.dry_run)
            return {
                "ok": True,
                "recipe": recipe,
                "nutrition": nutrition,
                "suitability": suitability,
                "target_users": target_users,
                "allowed_meal_types": allowed_meal_types,
            }
        except Exception as exc:
            last_error = exc
            time.sleep(1)
    return {"ok": False, "recipe": recipe, "reason": str(last_error)}


def backfill_recipes(args, columns: set[str]):
    recipes = load_recipes_for_backfill(args)
    print(f"recipes to backfill: {len(recipes)}")
    if not recipes:
        return
    client = create_client()
    success = 0
    skipped = 0
    with tqdm(total=len(recipes)) as pbar:
        with ThreadPoolExecutor(max_workers=max(1, args.workers)) as executor:
            futures = [executor.submit(process_backfill_recipe, client, recipe, args, columns) for recipe in recipes]
            for future in as_completed(futures):
                result = future.result()
                recipe = result["recipe"]
                if result["ok"]:
                    nutrition = result["nutrition"]
                    action = "would update" if args.dry_run else "updated"
                    tqdm.write(
                        f"[{recipe['id']}] {recipe['name']} {action} | "
                        f"portion={nutrition['portion_weight_g']}g | "
                        f"per100g={nutrition['energy']}kcal P{nutrition['protein']} "
                        f"C{nutrition['carbohydrate']} F{nutrition['fat']}"
                    )
                    success += 1
                else:
                    tqdm.write(f"[{recipe['id']}] {recipe['name']} skipped: {result['reason']}")
                    skipped += 1
                pbar.update(1)
    deleted = 0 if args.keep_invalid else delete_invalid_recipes(args.dry_run)
    print(f"finished backfill: success={success}, skipped={skipped}, invalid_deleted={deleted}, dry_run={args.dry_run}")


def main():
    args = parse_args()
    random.seed(args.seed)
    columns = get_recipe_columns()
    if args.task == "coverage":
        print_coverage(args, columns)
        return
    if args.task == "backfill":
        backfill_recipes(args, columns)
        return
    generate_recipes(args, columns)


if __name__ == "__main__":
    main()
