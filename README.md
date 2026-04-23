我已经通读了当前工作区的核心代码。这个项目是一个典型的 `Vue 3 + Vite` 前端、`Go + Gin + Gorm` 后端、`MySQL` 持久化的饮食推荐系统，核心目标是基于用户健康档案、营养目标和饮食约束生成每日饮食方案，并配合饮食记录、周报、收藏、采购清单和 AI 食物识别形成闭环。

**项目结构**
后端入口在 [main.go](/e:/vscodes/go/src/NutriPlan/cmd/server/main.go)，路由集中在 [router.go](/e:/vscodes/go/src/NutriPlan/internal/api/router/router.go)。前端路由在 [index.js](/e:/vscodes/go/src/NutriPlan/web/src/router/index.js)。推荐核心在 [recipe_service.go](/e:/vscodes/go/src/NutriPlan/internal/service/recipe_service.go)，营养模型在 [nutri_service.go](/e:/vscodes/go/src/NutriPlan/internal/service/nutri_service.go)。

**前端功能**
前端页面基本覆盖了完整用户流程：

- 欢迎页、登录、注册：负责用户认证和进入系统，[index.js](/e:/vscodes/go/src/NutriPlan/web/src/router/index.js)、[auth.js](/e:/vscodes/go/src/NutriPlan/web/src/store/auth.js)
- 健康档案录入与查看：用户填写性别、年龄、身高、体重、活动水平、健康目标、过敏源、饮食偏好、健康问题等，这些数据直接作为推荐输入，[Profile.vue](/e:/vscodes/go/src/NutriPlan/web/src/views/Profile.vue)
- 首页总览：展示当天热量、宏量营养素完成度、当前饮食模式、今日推荐方案入口，并支持快速添加饮食记录、切换餐次、查看推荐餐单，[Home.vue](/e:/vscodes/go/src/NutriPlan/web/src/views/Home.vue)
- 食谱推荐页：拉取多套推荐方案，展示每套方案的早餐/午餐/晚餐/加餐、总热量、蛋白质/碳水/脂肪，并支持锁定某个方案；还支持“受限单餐重构”，即限定食材后重生成某一餐，[RecipeRecommend.vue](/e:/vscodes/go/src/NutriPlan/web/src/views/RecipeRecommend.vue)
- 食谱详情页、收藏页：查看单个菜谱详情，加入/取消收藏
- 饮食记录页：支持手动添加、图片识别、文本 AI 估算，展示当天记录和按周汇总的摄入情况，[IntakeRecord.vue](/e:/vscodes/go/src/NutriPlan/web/src/views/IntakeRecord.vue)
- 周报页/周计划页：根据每日摄入生成周期性统计视图
- 购物清单页：把选定食谱转成采购清单
- 前后端通信：统一走 Axios 封装和 JWT 鉴权，[request.js](/e:/vscodes/go/src/NutriPlan/web/src/api/request.js)、[recipeApi.js](/e:/vscodes/go/src/NutriPlan/web/src/api/recipeApi.js)、[intakeApi.js](/e:/vscodes/go/src/NutriPlan/web/src/api/intakeApi.js)

前端的特点不是只“展示推荐结果”，而是把推荐方案和“实际吃了什么”联动起来，形成推荐 -> 选定 -> 记录 -> 反馈的流程。

**后端功能**
后端按典型分层组织：

- `handler`：接收 HTTP 请求，做参数绑定和响应 DTO 转换，[recipe_handler.go](/e:/vscodes/go/src/NutriPlan/internal/api/handler/recipe_handler.go)
- `service`：业务逻辑核心，包括用户、营养、推荐、摄入、购物清单、AI 食物识别
- `dao/repository`：数据库访问
- `models`：用户、菜谱、每日计划、饮食记录、收藏、购物清单等实体

主要能力包括：

- 用户注册、登录、JWT 鉴权
- 健康档案维护与营养需求计算
- 每日食谱推荐与方案持久化
- 已选方案查询与重新选择
- 菜谱详情、收藏管理
- 饮食记录增删查、当天营养状态统计、周报统计
- 购物清单生成
- AI 图片识别食物、AI 文本分析食物
- 受限食材单餐重构和采纳

路由总表可直接看 [router.go](/e:/vscodes/go/src/NutriPlan/internal/api/router/router.go)。

**推荐算法原理**
当前生效的推荐流程在 [recipe_service.go](/e:/vscodes/go/src/NutriPlan/internal/service/recipe_service.go) 的 `RecommendRecipes` 中。它不是旧版的“多路召回”，而是一个“营养目标建模 + 分阶段过滤 + 单餐组合生成 + 全日组合评分 + 兜底放宽”的流程。

1. 用户营养目标建模  
在 [nutri_service.go](/e:/vscodes/go/src/NutriPlan/internal/service/nutri_service.go) 中，系统先根据用户档案计算 `BMR`、`TDEE`、`BMI`，再按健康目标生成每日目标热量与三大营养素目标：

- 减脂：`TDEE - 热量赤字`
- 增肌：`TDEE + 热量盈余`
- 控糖：`TDEE - 轻微热量赤字`
- 维持健康：约等于 `TDEE`

宏量营养素配比同样随目标变化，例如增肌更强调能量和蛋白质，减脂更强调热量与脂肪控制。最终推荐的基准不是固定菜单，而是用户专属的 `目标热量 + 目标蛋白质 + 目标碳水 + 目标脂肪`。

2. 候选菜谱构建与分阶段过滤  
系统会先按餐次分别取出早餐、午餐、晚餐、加餐的基础候选池，然后依次做四层过滤：

- 允许餐次过滤：只保留适合当前餐次的菜谱
- 用户偏好过滤：结合 `TargetUsers` / `ForbiddenUsers`、健康标签、饮食偏好、过敏源做精筛
- 最近历史过滤：尽量剔除最近 5 天计划里已经出现过的食谱
- 饮食模式过滤：根据 `normal`、`light_adjust`、`heavy_adjust` 等模式控制禁用关键词、热量和脂肪上限

之后还会叠加健康目标过滤，例如减脂和控糖场景下会进一步剔除高糖饮品和不合适的高碳水加餐。

3. 逐级放宽的推荐阶段  
当前算法不是“一次过滤到底”，而是按阶段尝试生成方案：

- `strict`：最近历史过滤 + 饮食模式过滤 + 健康目标过滤
- `relax_health_goal`：放宽健康目标过滤，只保留最近历史与饮食模式约束
- `relax_diet_mode`：再放宽饮食模式约束
- `relax_recent_history`：减脂场景下如果还无解，最后连近期重复限制也放宽

也就是说，系统优先返回最严格、最符合当前状态的方案；只有严格阶段无结果时，才逐层回退，保证“尽量准”与“不要空结果”之间的平衡。

4. 单餐组合生成  
算法不是“一餐一道菜”，而是先把每一餐生成候选组合，再进入整日评分：

- 早餐：最多 2 个菜品组合
- 午餐、晚餐：2 到 3 个菜品组合，并要求优先包含主食
- 加餐：1 到 2 个菜品组合

组合生成时会先做候选裁剪和多样性抽样，再去重。早餐还有额外偏好，例如对“酸奶 + 鸡蛋”这类稳定早餐搭配做加分。每个组合都会先按该餐目标营养进行一次局部评分和截断，避免组合空间爆炸。

5. 全日组合评分  
系统把早餐、午餐、晚餐、加餐组合成一整天方案，并按 `25% / 35% / 30% / 10%` 分配单餐目标。评分分为两层：

- 单菜/单餐评分：看该食谱或组合与单餐目标的接近度
- 全日评分：看整天总热量、蛋白质、碳水、脂肪与日目标的偏差

全日匹配分是一个 `0-100` 的分数，本质上是“带权营养距离”的映射。不同目标的权重不同：

- 增肌：能量 `0.4`，蛋白质 `0.4`，碳水 `0.1`，脂肪 `0.1`
- 减脂：能量 `0.5`，蛋白质 `0.1`，碳水 `0.1`，脂肪 `0.3`
- 控糖：能量 `0.4`，蛋白质 `0.2`，碳水 `0.3`，脂肪 `0.1`
- 维持健康：能量 `0.4`，蛋白质 `0.2`，碳水 `0.2`，脂肪 `0.2`

6. 多样性控制、阈值截断与加权选优  
在整日组合阶段，系统还会显式控制“不要太重复”：

- 不允许同一道菜跨餐重复出现
- 食材交集会产生惩罚分，鼓励一日菜单更丰富
- 候选推荐之间如果午餐和晚餐主菜都重复，或者四餐里有 3 餐相同，会被判为过于相似
- 组合池较小时会自动降低最低分阈值，避免小样本下全部被筛空

最终所有候选方案先按分数排序，再对前 `Top 100` 做加权洗牌 `shuffleByWeights`，高分方案更容易被选中，但不会永远固定返回同一套。之后再结合最低得分门槛和多样性检查，选出最终返回的 `3-5` 套推荐。

**数据模型如何支撑算法**
- 用户输入在 [user.go](/e:/vscodes/go/src/NutriPlan/internal/repository/models/user.go)
- 菜谱营养字段、适用人群、禁忌人群在 [recipe.go](/e:/vscodes/go/src/NutriPlan/internal/repository/models/recipe.go)
- 每日计划既保存主菜，也保存组合内多个菜的 `ItemIDs`，说明系统已经支持“多食物单餐”的持久化，[daily_recipe_plan.go](/e:/vscodes/go/src/NutriPlan/internal/repository/models/daily_recipe_plan.go)
