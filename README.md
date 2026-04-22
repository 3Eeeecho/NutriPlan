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
这个项目的推荐算法不是简单“按热量排序取前几名”，而是一个“营养目标驱动 + 多路召回 + 组合生成 + 全日评分 + 多样性控制”的流程。

1. 用户营养目标建模  
在 [nutri_service.go](/e:/vscodes/go/src/NutriPlan/internal/service/nutri_service.go) 中，系统先根据用户档案计算：
- `BMR`
- `TDEE`
- `BMI`

然后依据健康目标决定每日目标热量：
- 减脂：`TDEE - 热量缺口`
- 增肌：`TDEE + 热量盈余`
- 控糖：`TDEE - 较小缺口`
- 维持：约等于 `TDEE`

再根据目标分配宏量营养素比例：
- 减脂更重视蛋白和脂肪控制
- 增肌更重视蛋白和总能量
- 控糖更重视碳水控制
- 维持则更均衡

也就是说，推荐算法的“标准答案”不是固定菜单，而是用户专属的 `目标热量 + 目标蛋白质 + 目标碳水 + 目标脂肪`。

2. 多路召回候选菜谱  
在 [recipe_service.go](/e:/vscodes/go/src/NutriPlan/internal/service/recipe_service.go) 的 `multiChannelRecall` 中，每个餐次都会并发从多个渠道找候选菜谱，而不是只靠一种规则：

- 按餐次和目标热量找接近项
- 低热量召回
- 高热量召回
- 按用户偏好标签召回
- 随机探索召回

这样做的作用是同时兼顾：
- 命中营养目标
- 保留多样性
- 避免推荐总是一样
- 给组合搭配留下空间

3. 用户约束过滤  
候选菜谱在进入组合前会经过几轮过滤：

- 过敏源、健康禁忌过滤
- `TargetUsers` / `ForbiddenUsers` 人群标签过滤，[recipe.go](/e:/vscodes/go/src/NutriPlan/internal/repository/models/recipe.go)
- 最近几天吃过的菜谱过滤，减少重复
- 饮食模式过滤，比如清淡模式、重调模式等，[user.go](/e:/vscodes/go/src/NutriPlan/internal/repository/models/user.go)

其中“饮食模式”比较有意思。系统不仅支持用户手动设置模式，也会根据前一天摄入偏差自动切换到 `light_adjust`、`heavy_adjust` 等模式，再用不同的热量/脂肪约束修正推荐。

4. 单餐组合生成  
这部分是当前算法的关键升级点。不是“一餐只推荐一道菜”，而是支持一个餐次由多个菜组成。

从 [algorithm.md](/e:/vscodes/go/src/NutriPlan/algorithm.md) 和 [recipe_service.go](/e:/vscodes/go/src/NutriPlan/internal/service/recipe_service.go) 可以确认：
- 早餐支持 1 到 2 个菜品组合
- 午餐、晚餐支持最多 3 个菜品组合
- 加餐支持最多 2 个菜品组合
- 组合会去重，并按单餐目标营养进行排序
- 早餐还对某些优先组合做加权，例如酸奶 + 鸡蛋这类稳定搭配

这说明算法从“单菜推荐”升级成了“单餐搭配推荐”。

5. 全日组合与评分  
之后系统会把早餐/午餐/晚餐/加餐的候选组合做全日组合，再计算整天方案的总营养。

核心评分逻辑是：
- 计算实际摄入和目标值的偏差
- 按健康目标设置不同权重
- 生成一个 `0-100` 的匹配分

在 [recipe_service.go](/e:/vscodes/go/src/NutriPlan/internal/service/recipe_service.go) 中能看到两层评分思想：
- 单菜或单餐评分：看局部营养接近度
- 全日评分：看整天总热量、蛋白、碳水、脂肪与目标差距

不同目标的权重大致是：
- 增肌：更重视总能量和蛋白质
- 减脂：更重视热量和脂肪控制
- 控糖：更重视碳水控制
- 维持：相对均衡

本质上这是一个“带权营养距离最小化”的方案排序问题。

6. 多样性与重叠惩罚  
系统不仅追求“营养最接近”，还控制“吃得别太重复”。

做法包括：
- 跨餐不能大量重复同一道菜
- 原料重叠会被惩罚
- 午餐/晚餐过度相似的方案会被剔除
- 最近几天吃过的菜也尽量不再推荐

因此最终得分可以理解成：

`最终得分 = 营养匹配得分 - 食材/菜谱重复惩罚`

这使得推荐结果更像真实的一日菜单，而不是数学上最接近但体验很差的组合。

7. 加权随机而不是死板 Top1  
生成所有方案后，系统不是永远返回分数最高的固定方案，而是先排序，再对高分方案做加权洗牌 `shuffleByWeights`。  
高分更容易被选中，但不是绝对固定。

这样做的意义：
- 保留质量
- 增加刷新后的变化性
- 避免用户每次看到完全相同的结果

**数据模型如何支撑算法**
- 用户输入在 [user.go](/e:/vscodes/go/src/NutriPlan/internal/repository/models/user.go)
- 菜谱营养字段、适用人群、禁忌人群在 [recipe.go](/e:/vscodes/go/src/NutriPlan/internal/repository/models/recipe.go)
- 每日计划既保存主菜，也保存组合内多个菜的 `ItemIDs`，说明系统已经支持“多食物单餐”的持久化，[daily_recipe_plan.go](/e:/vscodes/go/src/NutriPlan/internal/repository/models/daily_recipe_plan.go)
