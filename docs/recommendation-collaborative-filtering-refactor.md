# 后端推荐算法协同过滤重构方案

## 1. 背景

当前食谱推荐入口是 `RecipeService.RecommendRecipes(user, count)`，由 `RecipeHandler.GetRecommendations` 通过 `GET /api/v1/recipes/recommend` 调用。现有逻辑主要是规则型内容推荐：

- 根据用户健康档案计算目标热量和三大营养素。
- 按餐次加载候选食谱。
- 按健康目标、饮食偏好、过敏/禁忌、最近吃过的食谱过滤。
- 组合早餐、午餐、晚餐、加餐，并用营养匹配度、食材重复惩罚、多样性规则排序。
- 通过 `DailyRecipePlan` 返回每日食谱计划。

这套逻辑适合做安全兜底和营养约束，但不能充分利用用户群体行为，例如“相似用户常吃什么”“收藏过 A 的用户还喜欢 B”“某类用户对哪些组合完成率更高”。因此需要引入协同过滤。

## 2. 重构目标

1. 引入基于隐式反馈的协同过滤推荐能力。
2. 保留现有营养、健康目标、禁忌、餐次、近期去重等安全约束。
3. 推荐结果仍兼容当前 `DailyRecipePlan` DTO，不影响前端调用。
4. 支持冷启动用户和低数据量场景自动回退到现有规则推荐。
5. 将推荐算法拆成可测试、可替换的模块，避免继续膨胀 `recipe_service.go`。

## 3. 推荐策略选择

推荐采用“Item-Based Collaborative Filtering + 内容/规则约束混排”的第一阶段方案。

原因：

- 当前数据规模预计不大，Item-CF 比矩阵分解更容易落地和调试。
- 食谱天然是稳定 item，用户行为较稀疏时，物品相似度比用户相似度更稳定。
- 可以直接利用收藏、选中计划、真实打卡记录构造隐式反馈。
- 推荐结果仍可以接入现有餐次组合器，降低一次性重写风险。

后续数据量增长后，可以在同一接口下替换为 ALS/BPR/Embedding 召回。

## 4. 可用行为数据

现有表可直接使用：

- `user_favorite_recipes`
  - 用户收藏食谱，强正反馈。
- `daily_recipe_plans`
  - 用户选中计划，表示接受推荐。
  - 计划中包含早餐、午餐、晚餐、加餐食谱 ID 和组合 ID 列表。
- `daily_intake_records`
  - 真实饮食记录。
  - `food_source = 1` 表示来自食谱，可作为最强正反馈。
- `users`
  - 健康目标、过敏、饮食偏好、健康条件、TDEE、BMI 等，用于过滤和冷启动。

建议后续新增行为表，统一记录显式/隐式事件：

```sql
CREATE TABLE user_recipe_events (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  recipe_id BIGINT UNSIGNED NOT NULL,
  event_type VARCHAR(32) NOT NULL,
  weight DECIMAL(6,2) NOT NULL,
  event_date DATETIME(3) NOT NULL,
  created_at DATETIME(3),
  updated_at DATETIME(3),
  deleted_at DATETIME(3),
  INDEX idx_user_event (user_id, event_date),
  INDEX idx_recipe_event (recipe_id, event_date),
  INDEX idx_user_recipe (user_id, recipe_id)
);
```

第一阶段可以先不建表，直接从现有三张表聚合；第二阶段再落统一事件表，降低查询复杂度。

## 5. 隐式反馈权重

建议初始权重如下：

| 行为 | 来源 | 权重 |
| --- | --- | --- |
| 实际打卡食谱 | `daily_intake_records.food_source = 1` | 5.0 |
| 选中食谱计划 | `daily_recipe_plans.is_selected = true` | 3.0 |
| 收藏食谱 | `user_favorite_recipes` | 2.0 |
| 最近重复吃过 | `daily_intake_records` 近 5 天 | -1.5 |

最终用户-食谱偏好分：

```text
preference(user, recipe) = sum(event_weight * time_decay)
time_decay = exp(-days_ago / 30)
```

负反馈第一阶段只用于降权，不建议直接删除候选，避免低数据量时召回为空。

## 6. 协同过滤算法设计

### 6.1 Item-CF 相似度

构造用户-食谱偏好矩阵 `R[u][i]`，其中 `u` 是用户，`i` 是食谱。

食谱相似度使用余弦相似度：

```text
sim(i, j) = dot(R[:,i], R[:,j]) / (norm(R[:,i]) * norm(R[:,j]))
```

推荐分：

```text
cf_score(user, candidate) =
  sum(preference(user, history_item) * sim(history_item, candidate))
  / sum(abs(sim(history_item, candidate)))
```

只保留每个食谱 TopK 相似食谱，例如 `K = 50`。

### 6.2 约束过滤

协同过滤只负责“召回和偏好分”，不能绕过健康安全约束。候选仍必须经过：

- 餐次匹配：早餐、午餐、晚餐、加餐。
- 用户过敏/禁忌过滤。
- 健康目标过滤或降权。
- 饮食模式过滤。
- 最近吃过的食谱降权或过滤。
- 营养目标校验。

### 6.3 混合排序

最终食谱级别分数：

```text
recipe_score =
  0.45 * cf_score
  + 0.30 * nutrition_fit_score
  + 0.15 * content_match_score
  + 0.10 * diversity_score
```

对于冷启动或 CF 数据不足：

```text
recipe_score =
  0.55 * nutrition_fit_score
  + 0.30 * content_match_score
  + 0.15 * diversity_score
```

计划级别仍复用现有 `calculateMatchScore`、食材重复惩罚和多样性检查，但组合前每个餐次候选池应按 `recipe_score` 截断到 TopN，减少组合爆炸。

## 7. 代码结构重构

建议新增以下模块：

```text
internal/service/recommendation/
  engine.go
  collaborative_filter.go
  behavior_repository.go
  candidate_builder.go
  scorer.go
  constraints.go
  plan_builder.go
```

### 7.1 Engine

```go
type RecommendationEngine interface {
    Recommend(ctx context.Context, user *models.User, count int) ([]*models.DailyRecipePlan, error)
}
```

`RecipeServiceImpl.RecommendRecipes` 改为调用 `RecommendationEngine`，保留方法签名，避免影响 handler。

### 7.2 BehaviorRepository

负责读取用户行为：

```go
type BehaviorRepository interface {
    GetUserRecipePreferences(ctx context.Context, userID uint, days int) ([]UserRecipePreference, error)
    GetAllRecipePreferences(ctx context.Context, days int) ([]UserRecipePreference, error)
}
```

第一阶段从现有表聚合：

- 收藏表。
- 选中计划表。
- 饮食记录表。

### 7.3 CollaborativeFilter

负责构建或加载 item 相似度，并返回用户候选：

```go
type CollaborativeFilter interface {
    RecommendRecipes(ctx context.Context, userID uint, mealType models.MealType, limit int) ([]RecipeScore, error)
}
```

初期可以每次请求按最近 90 天行为在线计算，数据量上来后改成定时任务预计算。

### 7.4 ConstraintFilter

把当前分散在 `recipe_service.go` 的过滤逻辑收敛：

- `filterRecipesByAllowedMealType`
- `filterRecipesByDietMode`
- `filterRecipesByHealthGoal`
- `filterRecipesByUserPreferences`
- `filterRecentRecipes`

### 7.5 PlanBuilder

保留现有组合能力，但输入改成已经带混合分的候选池。

## 8. 数据流

```text
GET /api/v1/recipes/recommend
  -> RecipeHandler.GetRecommendations
  -> RecipeService.RecommendRecipes
  -> RecommendationEngine.Recommend
      -> load user nutrition target
      -> load user behavior
      -> collaborative filtering recall
      -> rule/content fallback recall
      -> constraint filtering
      -> hybrid scoring
      -> meal combo generation
      -> daily plan ranking
  -> DailyPlanDTO
```

## 9. 接口兼容性

外部接口保持不变：

```text
GET /api/v1/recipes/recommend?count=3
```

响应仍为：

```json
{
  "plans": []
}
```

可选增强：新增调试字段，但默认不返回。

```go
type RecommendRequest struct {
    Count int `form:"count" json:"count"`
    Debug bool `form:"debug" json:"debug"`
}
```

调试字段可包括：

- `algorithm`: `collaborative_filtering` / `fallback_rules`
- `stage`: `cf_strict` / `cf_relaxed` / `cold_start`
- `candidate_count`

## 10. 冷启动策略

### 新用户

无行为历史时，使用现有规则推荐：

- 健康目标。
- TDEE 和三大营养素目标。
- 过敏和饮食偏好。
- 热门食谱。

热门食谱可由近 30 天行为聚合：

```text
popular_score = 5 * intake_count + 3 * selected_count + 2 * favorite_count
```

### 新食谱

没有协同过滤相似度时，使用内容相似：

- 餐次相同。
- 健康目标标签相同。
- 食材 Jaccard 相似度。
- 营养接近度。

## 11. 性能方案

第一阶段：

- 最近 90 天行为在线计算。
- 限制用户数和行为数，例如最多读取最近 50,000 条行为。
- 每个用户历史最多取最近/最高分 50 个食谱。
- 每个餐次候选池截断到 Top 30。

第二阶段：

新增相似度缓存表：

```sql
CREATE TABLE recipe_similarities (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  recipe_id BIGINT UNSIGNED NOT NULL,
  similar_recipe_id BIGINT UNSIGNED NOT NULL,
  score DECIMAL(8,4) NOT NULL,
  computed_at DATETIME(3) NOT NULL,
  UNIQUE KEY uk_recipe_similar (recipe_id, similar_recipe_id),
  INDEX idx_recipe_score (recipe_id, score)
);
```

通过后台任务每天或每小时重算。

## 12. 测试计划

### 单元测试

- 行为权重聚合。
- 时间衰减计算。
- Item-CF 相似度计算。
- 用户候选召回。
- 冷启动回退。
- 禁忌食材不会出现在结果中。
- 营养目标分不会被 CF 分绕过。

### 集成测试

- 有收藏和打卡历史的用户，推荐包含相似食谱。
- 无行为用户仍可获得推荐。
- 管理员新增食谱后，新食谱可进入冷启动候选。
- 最近吃过的食谱被降权或过滤。

### 回归测试

- `GET /api/v1/recipes/recommend` 响应结构不变。
- 选中计划、保存推荐、详情收藏逻辑不受影响。

## 13. 迁移步骤

当前状态：阶段 1、阶段 2、阶段 3、阶段 4、阶段 5 收尾已完成。代码已新增 `RecommendationEngine` 抽象和 `RuleBasedRecommendationEngine`，`RecipeServiceImpl.RecommendRecipes` 已改为委托推荐引擎，默认行为仍保留规则约束。行为聚合层已能从收藏、选中计划和真实打卡记录生成用户-食谱偏好分。Item-CF 召回组件已能基于用户行为矩阵计算相似食谱候选，并已接入餐次候选池与候选排序。推荐接口已返回算法元数据，前端食谱推荐页已展示协同过滤混合推荐状态，推荐页打卡会写入 `food_source/source_id` 形成隐式反馈闭环。

### 阶段 1：抽象和无行为回归

1. 新增 `RecommendationEngine` 接口。
2. 将现有 `RecommendRecipes` 逻辑迁移到 `RuleBasedRecommendationEngine`。
3. `RecipeServiceImpl.RecommendRecipes` 改为调用 engine。
4. 保证所有现有测试通过。

### 阶段 2：行为聚合

1. 新增 `BehaviorRepository`。
2. 从收藏、选中计划、饮食记录聚合隐式反馈。
3. 增加行为聚合单元测试。

状态：已完成。当前实现包含 `BehaviorRepository`、`GormBehaviorRepository`、`UserRecipePreference`，并已接入 `RecipeServiceImpl` 依赖，为后续 Item-CF 召回和混排提供行为偏好数据。

### 阶段 3：Item-CF 召回

1. 新增 `CollaborativeFilter`。
2. 在线计算 item 相似度。
3. 对有历史用户返回 CF 候选。
4. 对无历史用户回退规则推荐。

状态：已完成。当前实现包含 `CollaborativeFilter`、`ItemCollaborativeFilter`、`RecipeScore`，支持按餐次返回协同过滤候选。

### 阶段 4：混合排序和计划生成

1. 将 CF 分数注入餐次候选。
2. 改造 `buildRecommendedPlans`，支持候选分数。
3. 保留营养目标和多样性为硬约束。

状态：已完成。当前实现会按餐次加载 CF 候选，将候选合并进原始餐次池，再继续经过用户偏好、近期历史、饮食模式、健康目标等既有过滤。CF 分数被归一化为最多 15 分的候选排序加成，只影响 `selectDiverseCandidates` 的候选截断顺序，不绕过最终营养匹配、组合去重和多样性检查。

### 阶段 5：缓存和运维

1. 增加 `recipe_similarities` 缓存表。
2. 增加定时重算脚本或后台任务。
3. 增加推荐日志和关键指标。

状态：基础收尾已完成。当前版本先采用用户级短 TTL 内存缓存，避免一次推荐请求中按四个餐次重复计算同一用户的 Item-CF 排序；接口返回 `algorithm` 元数据，便于前端展示和后续观测。数据库级 `recipe_similarities` 缓存表和定时重算任务保留为数据量扩大后的运维增强项。

## 14. 风险和应对

| 风险 | 应对 |
| --- | --- |
| 数据稀疏导致 CF 无结果 | 自动回退规则推荐 |
| 热门食谱过度集中 | 加多样性惩罚和同食材去重 |
| 推荐违反过敏/禁忌 | 约束过滤放在最终候选前，作为硬规则 |
| 在线计算变慢 | 限制窗口和候选数，后续引入相似度缓存 |
| 新食谱没有曝光 | 内容相似和热门兜底给新食谱探索流量 |
| 算法不可解释 | 增加 debug 信息和推荐日志 |

## 15. 建议的第一版验收标准

1. 有行为历史用户的推荐结果会优先包含相似食谱。
2. 无行为历史用户仍能稳定返回 `count` 个计划。
3. 推荐结果不包含用户过敏/禁忌食谱。
4. 返回结构与当前前端兼容。
5. `go test ./...` 通过。
6. 推荐耗时在本地开发数据量下不超过现有算法 2 倍。

## 16. 不在第一版处理的内容

- 深度学习推荐模型。
- 实时事件流。
- A/B 测试平台。
- 向量数据库。

这些能力应在 Item-CF 版本稳定、行为数据足够后再设计。
