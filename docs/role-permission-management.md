# 多权限角色管理功能设计文档

## 1. 背景

当前系统已经具备用户注册、登录、JWT 鉴权和登录态路由守卫，但权限模型仍停留在“游客/已登录”两级判断。随着系统功能扩展，需要引入明确的角色体系，用于控制页面访问、接口访问和管理类操作。

本设计引入三种角色：

- 游客：未登录用户，只能访问公开页面和公开数据。
- 注册用户：已注册并登录的普通用户，可使用个人营养、食谱推荐、收藏、饮食记录、购物清单等个人功能。
- 管理员：系统管理人员，除普通用户能力外，可管理用户、食谱、分类、系统数据和运营内容。

## 2. 目标

- 建立统一的角色枚举和权限判断入口。
- 后端接口必须基于角色做强制校验，不能只依赖前端隐藏入口。
- 前端根据角色控制路由、导航菜单和操作按钮展示。
- JWT 中携带角色信息，减少普通请求的数据库查询成本。
- 保持当前注册用户功能不受影响，默认注册用户角色为 `user`。
- 支持后续扩展更多角色或细粒度权限点。

## 3. 角色定义

| 角色 | 标识 | 说明 |
| --- | --- | --- |
| 游客 | `guest` | 未登录访问者，不存入用户表，不签发 JWT |
| 注册用户 | `user` | 注册并登录的普通用户，是系统默认角色 |
| 管理员 | `admin` | 具备后台管理权限的用户 |

角色等级关系：

```text
guest < user < admin
```

说明：

- `guest` 是运行时身份，不是数据库用户角色。
- `user` 是用户注册后的默认角色。
- `admin` 只能由已有管理员或初始化脚本设置，不允许用户注册时指定。

## 4. 权限矩阵

| 功能模块 | 游客 | 注册用户 | 管理员 |
| --- | --- | --- | --- |
| 访问欢迎页 | 允许 | 允许 | 允许 |
| 用户注册 | 允许 | 不需要 | 不需要 |
| 用户登录 | 允许 | 允许 | 允许 |
| 查看公开食谱列表/详情 | 允许 | 允许 | 允许 |
| 维护个人健康档案 | 禁止 | 允许本人 | 允许本人 |
| 获取营养需求 | 禁止 | 允许本人 | 允许本人 |
| 获取食谱推荐 | 禁止 | 允许本人 | 允许本人 |
| 收藏食谱 | 禁止 | 允许本人 | 允许本人 |
| 饮食记录增删查 | 禁止 | 允许本人 | 允许本人 |
| 周报查看 | 禁止 | 允许本人 | 允许本人 |
| 购物清单管理 | 禁止 | 允许本人 | 允许本人 |
| AI 食物识别 | 禁止 | 允许本人 | 允许本人 |
| 用户列表管理 | 禁止 | 禁止 | 允许 |
| 用户角色调整 | 禁止 | 禁止 | 允许 |
| 食谱新增/编辑/删除 | 禁止 | 禁止 | 允许 |
| 分类/标签维护 | 禁止 | 禁止 | 允许 |
| 查看系统统计 | 禁止 | 禁止 | 允许 |

## 5. 后端设计

### 5.1 数据模型

在 `internal/repository/models/user.go` 的 `User` 模型中新增角色字段：

```go
Role UserRole `gorm:"type:varchar(20);not null;default:'user';index;comment:用户角色" json:"role"`
```

新增角色枚举，建议放在 `internal/repository/models/enum.go`：

```go
type UserRole string

const (
    UserRoleGuest UserRole = "guest"
    UserRoleUser  UserRole = "user"
    UserRoleAdmin UserRole = "admin"
)
```

数据库迁移：

```sql
ALTER TABLE users
  ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' COMMENT '用户角色',
  ADD INDEX idx_users_role (role);
```

初始化管理员方式：

- 方案一：提供一次性 SQL，把指定用户升级为管理员。
- 方案二：提供启动脚本或命令行工具创建管理员。
- 不允许公开注册接口传入 `role=admin`。

示例 SQL：

```sql
UPDATE users SET role = 'admin' WHERE username = 'admin';
```

### 5.2 JWT 设计

当前 JWT Claims 包含 `user_id` 和 `username`，需要新增 `role`：

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

登录成功时通过 `jwt.GenerateToken(user.ID, user.Username, string(user.Role))` 写入角色。

鉴权中间件解析后写入 Gin Context：

```go
c.Set("user_id", claims.UserID)
c.Set("username", claims.Username)
c.Set("role", claims.Role)
```

安全要求：

- JWT 中的角色只代表签发时状态。
- 当管理员调整用户角色后，该用户旧 token 可能仍在有效期内。
- 如果要立即生效，需要增加 token 版本号、用户状态版本号，或缩短 token 有效期。

本期建议：

- 先采用 JWT 携带角色。
- 管理员变更用户角色后提示用户重新登录。
- 后续再引入 token 版本校验。

### 5.3 中间件设计

保留现有 `AuthMiddleware()`，新增两个权限中间件：

```go
func OptionalAuthMiddleware() gin.HandlerFunc
func RequireRole(roles ...models.UserRole) gin.HandlerFunc
```

用途：

- `AuthMiddleware()`：必须登录，适用于普通用户和管理员功能。
- `OptionalAuthMiddleware()`：可登录可不登录，适用于公开食谱详情等游客可看、登录后可返回更多个人状态的接口。
- `RequireRole()`：角色校验，适用于管理员接口。

`RequireRole()` 行为：

- 未登录或缺少角色：返回 `401 Unauthorized`。
- 已登录但角色不满足：返回 `403 Forbidden`。
- 角色满足：继续请求。

示例：

```go
admin := v1.Group("/admin")
admin.Use(jwt.AuthMiddleware(), jwt.RequireRole(models.UserRoleAdmin))
{
    admin.GET("/users", adminHandler.ListUsers)
    admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
}
```

### 5.4 路由分层

建议将 API 分成三类：

公开接口：

```text
POST /api/v1/register
POST /api/v1/login
GET  /api/v1/public/recipes
GET  /api/v1/public/recipes/:id
```

用户接口：

```text
GET  /api/v1/user/profile
PUT  /api/v1/user/profile
GET  /api/v1/user/nutrition
GET  /api/v1/user/diet-mode
PUT  /api/v1/user/diet-mode
GET  /api/v1/recipes/recommend
POST /api/v1/recipes/recommend
POST /api/v1/recipes/:id/favorite
DELETE /api/v1/recipes/:id/favorite
GET  /api/v1/recipes/favorites
POST /api/v1/intake/records
DELETE /api/v1/intake/records/:id
GET  /api/v1/intake/today
GET  /api/v1/intake/weekly
POST /api/v1/shopping/lists
GET  /api/v1/shopping/lists
POST /api/v1/food/recognize
POST /api/v1/food/analyze-text
```

管理员接口：

```text
GET    /api/v1/admin/users
GET    /api/v1/admin/users/:id
PUT    /api/v1/admin/users/:id/role
PATCH  /api/v1/admin/users/:id/status
GET    /api/v1/admin/recipes
POST   /api/v1/admin/recipes
PUT    /api/v1/admin/recipes/:id
DELETE /api/v1/admin/recipes/:id
GET    /api/v1/admin/categories
POST   /api/v1/admin/categories
PUT    /api/v1/admin/categories/:id
DELETE /api/v1/admin/categories/:id
GET    /api/v1/admin/statistics
```

### 5.5 管理员功能

管理员一期建议包含：

- 用户管理：用户列表、搜索、查看详情、禁用/启用用户、调整角色。
- 食谱管理：新增、编辑、删除、上下架、按分类筛选。
- 分类管理：维护食谱分类、标签、适用人群。
- 系统统计：用户数、食谱数、饮食记录数、收藏数、AI 识别调用量。

用户角色调整规则：

- 管理员可以把 `user` 升级为 `admin`。
- 管理员可以把其他管理员降级为 `user`。
- 不允许管理员降级自己，避免系统没有可用管理员。
- 不允许把任何数据库用户设置为 `guest`。

### 5.6 接口响应规范

登录响应新增 `role`：

```json
{
  "token": "jwt-token",
  "user": {
    "id": 1,
    "username": "demo",
    "email": "demo@example.com",
    "role": "user"
  },
  "message": "登录成功"
}
```

无权限响应：

```json
{
  "error": "权限不足"
}
```

推荐状态码：

- `401 Unauthorized`：未登录、token 缺失、token 无效。
- `403 Forbidden`：已登录但角色权限不足。
- `404 Not Found`：资源不存在，或为避免越权枚举而隐藏资源存在性。

## 6. 前端设计

### 6.1 登录态 Store

`web/src/store/auth.js` 中需要持久化用户角色：

```js
const role = computed(() => user.value?.role || 'guest')
const isAdmin = computed(() => role.value === 'admin')
const isUser = computed(() => role.value === 'user' || role.value === 'admin')
```

退出登录后角色恢复为 `guest`。

### 6.2 路由元信息

`web/src/router/index.js` 中建议扩展路由 meta：

```js
meta: {
  requiresAuth: true,
  roles: ['user', 'admin']
}
```

管理员路由示例：

```js
{
  path: '/admin',
  component: () => import('@/layouts/AdminShell.vue'),
  meta: { requiresAuth: true, roles: ['admin'] },
  children: [
    { path: 'users', name: 'AdminUsers', component: () => import('@/views/admin/UserList.vue') },
    { path: 'recipes', name: 'AdminRecipes', component: () => import('@/views/admin/RecipeList.vue') },
    { path: 'statistics', name: 'AdminStatistics', component: () => import('@/views/admin/Statistics.vue') }
  ]
}
```

路由守卫逻辑：

- 未登录访问 `requiresAuth=true`：跳转登录页。
- 已登录但角色不在 `roles`：跳转 `403` 页面。
- 管理员访问普通用户页面：允许。
- 普通用户访问管理员页面：禁止。

### 6.3 菜单和按钮

顶部导航或侧边栏根据角色展示入口：

- 游客：欢迎页、登录、注册、公开食谱。
- 注册用户：首页、个人档案、推荐食谱、收藏、饮食记录、周报、购物清单。
- 管理员：普通用户菜单 + 管理后台入口。

页面按钮也要按角色控制：

- 食谱“编辑/删除/上下架”只对管理员展示。
- 用户“调整角色/禁用账号”只对管理员展示。
- 收藏、饮食记录、购物清单操作只对登录用户展示。

注意：前端隐藏按钮只改善体验，真实权限必须以后端校验为准。

## 7. 数据安全与越权控制

个人数据接口必须同时满足：

- 请求方已登录。
- 数据归属用户 ID 等于当前 `user_id`，或当前角色为 `admin` 且接口明确支持管理员查看。

普通用户禁止通过传入其他用户 ID 访问以下数据：

- 健康档案
- 营养需求
- 饮食记录
- 周报
- 收藏列表
- 购物清单
- 已选食谱计划

管理员接口必须使用 `/api/v1/admin` 路由前缀，避免和普通用户接口混用。

## 8. 实施步骤

1. 数据层：新增 `UserRole` 枚举和 `users.role` 字段，执行数据库迁移。
2. 登录鉴权：JWT Claims 增加 `role`，登录响应返回角色。
3. 中间件：新增 `RequireRole`，统一处理 `401` 和 `403`。
4. 路由调整：划分公开接口、用户接口和管理员接口。
5. 管理接口：实现用户管理、角色调整、食谱管理等管理员功能。
6. 前端 Store：持久化并计算当前用户角色。
7. 前端路由：增加 `roles` meta 和 `403` 页面。
8. 前端菜单：根据角色展示不同导航和操作按钮。
9. 测试：覆盖游客、普通用户、管理员三类访问路径。

## 9. 测试用例

| 场景 | 预期 |
| --- | --- |
| 游客访问欢迎页 | 成功 |
| 游客访问个人档案 | 跳转登录页，接口返回 401 |
| 游客访问公开食谱详情 | 成功 |
| 普通用户访问首页 | 成功 |
| 普通用户访问自己的饮食记录 | 成功 |
| 普通用户访问管理员用户列表 | 前端跳转 403，接口返回 403 |
| 管理员访问用户列表 | 成功 |
| 管理员调整普通用户为管理员 | 成功 |
| 管理员尝试降级自己 | 失败，返回 400 或 403 |
| 用户注册时传入 role=admin | 忽略或拒绝，最终角色仍为 user |
| 旧 token 在用户被降级后访问管理员接口 | 本期可能仍有效，需记录为已知限制 |

## 10. 验收标准

- 新注册用户默认角色为 `user`。
- 未登录用户不再能访问任何个人数据接口。
- 普通用户无法访问 `/api/v1/admin/**`。
- 管理员可以访问后台管理页面和管理员接口。
- 登录响应、用户信息、本地 Store 均能正确识别角色。
- 前端导航能根据游客、注册用户、管理员展示不同入口。
- 后端权限校验覆盖所有管理接口。
- 权限不足时返回明确的 `403 Forbidden`。

## 11. 后续扩展

后续如果角色继续增加，可以从 RBAC 简化模型升级为“角色 + 权限点”模型：

```text
roles
permissions
role_permissions
user_roles
```

当前三角色场景中，单字段 `users.role` 已能满足需求，复杂 RBAC 可以等管理功能继续扩展后再引入。
