# NutriPlan 推荐算法（当前实现总结）

本文总结当前分阶段落地后的推荐算法实现，重点覆盖“每餐多食物组合”和“多样性控制”。

---

## 1. 目标与输入

- 输入：用户画像（`TDEE`、`HealthGoal`、过敏/偏好/健康状况）。
- 输出：每日推荐计划 `N` 套（默认 3~5），每套包含早餐/午餐/晚餐/加餐。
- 约束：
   - 过滤过敏与禁忌人群；
   - 过滤近期已吃过菜谱；
   - 保持营养目标匹配度（热量+三大营养素）。

---

## 2. 阶段一：多路召回（已实现）

在 `RecipeService.RecommendRecipes` 中对每餐并发执行召回：

1. 传统热量接近召回（`FindByMealType` + `getTopCandidates`）
2. 低热量召回（`FindLowCalorieRecipes`）
3. 高热量召回（`FindHighCalorieRecipes`）
4. 偏好标签召回（`FindByIngredientTags`）
5. 随机探索召回（`FindRandomRecipes`）

随后统一去重、做用户偏好过滤和历史去重，形成候选池。

---

## 3. 阶段二：单餐组合生成（已实现）

### 3.1 早餐组合（`buildBreakfastCombos`）

- 从候选集中先取 Top-K。
- 组合生成支持 1 道或 2 道菜。
- 稳定优先策略：若命中“酸奶”和“鸡蛋”，优先注入该组合并给额外加权。
- 输出 `mealCombo`：`Items`（组合全量）+ `Primary`（兼容旧 ID 字段）+ 汇总营养。

### 3.2 午餐/晚餐/加餐组合（`buildMealCombos`）

- 通用生成器支持：
   - 午/晚餐：最多 3 道菜；
   - 加餐：最多 2 道菜；
   - 单菜方案仍可保留（兼容老逻辑）。
- 通过组合 key（排序后 ID 串）去重，避免重复组合。
- 通过 `rankMealCombos` 按单餐目标评分后截取 Top-N。

---

## 4. 阶段三：全日组合与多样性约束（已实现）

`generateRankedPlans` 由“单菜笛卡尔积”升级为“组合笛卡尔积”：

- 早餐组合 × 午餐组合 × 晚餐组合 × 加餐组合。
- 硬约束：
   - 跨餐菜品去重（`hasRecipeOverlap`）；
   - 热量上限剪枝（超过目标上限直接跳过）。
- 评分：
   - 基础分：`calculateMatchScore`（热量+三宏量）；
   - 惩罚分：`ingredientDiversityPenalty`（跨餐食材重叠扣分）。

即：

$$
Score_{final}=Score_{nutrition}-Penalty_{ingredient\ overlap}
$$

最终保留 `matchScore >= 60` 的方案，排序后再走加权洗牌与最终多样性过滤。

---

## 5. DTO 与兼容策略（已实现）

已扩展返回 DTO，同时保持旧接口兼容：

- 保留旧字段：`breakfast/lunch/dinner/snack`（单菜视图）
- 新增字段：
   - `breakfast_items`
   - `lunch_items`
   - `dinner_items`
   - `snack_items`

模型侧增加非持久化字段（`gorm:"-"`），用于推荐阶段承载组合明细，不影响数据库结构与“选中计划”逻辑。

---

## 6. 当前验证结果（已通过）

服务层关键测试均通过：

- `TestBuildBreakfastCombos_PreferYogurtAndEgg`
- `TestBuildMealCombos_GeneratesMultiItemCombo`
- `TestIngredientDiversityPenalty_WithOverlap`

验证结论：

1. 早餐可稳定优先产出“酸奶 + 鸡蛋”组合；
2. 午/晚餐可生成多菜组合；
3. 跨餐食材重叠会触发惩罚。

---

## 7. 后续建议（下一阶段）

1. 将组合明细持久化（JSON 字段或新表）。
2. 升级“选中计划”接口，支持保存完整组合。
3. 引入用户反馈闭环，动态学习偏好权重。
