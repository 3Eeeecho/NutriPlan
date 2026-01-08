<template>
  <div class="home-container">
    <!-- 顶部导航 -->
    <div class="top-bar">
      <div class="top-bar-content">
        <div class="brand">
          <span class="brand-icon">🥗</span>
          <span class="brand-name">NutriPlan</span>
        </div>
        <el-dropdown @command="handleCommand" trigger="click">
          <div class="user-menu">
            <el-avatar
              :size="36"
              :icon="UserFilled"
              style="background: #10b981"
            />
            <span class="user-name">{{ authStore.user?.username }}</span>
            <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人档案</el-dropdown-item>
              <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
              <el-dropdown-item command="intake">饮食记录</el-dropdown-item>
              <el-dropdown-item command="weekly">周报告</el-dropdown-item>
              <el-dropdown-item divided command="logout"
                >退出登录</el-dropdown-item
              >
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-wrapper">
      <!-- 欢迎区域 -->
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">
            {{ getGreeting() }}，<span class="highlight">{{
              authStore.user?.username
            }}</span>
            👋
          </h1>
          <p class="hero-subtitle">今天想吃点什么呢？</p>
        </div>
        <div
          class="quick-stats"
          v-if="authStore.hasProfile && authStore.profile"
        >
          <div class="stat-item">
            <div class="stat-icon">🔥</div>
            <div class="stat-info">
              <div class="stat-value">
                {{
                  authStore.profile.tdee
                    ? Math.floor(authStore.profile.tdee)
                    : "--"
                }}
              </div>
              <div class="stat-label">每日目标热量</div>
            </div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-icon">⚖️</div>
            <div class="stat-info">
              <div class="stat-value">
                {{ authStore.profile.weight || "--" }}
              </div>
              <div class="stat-label">当前体重 (kg)</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 功能卡片区 -->
      <div class="content-section">
        <div class="features-grid">
          <!-- 主要功能卡片区 (顶层双巨头) -->
          <div class="big-cards-row">
            <!-- 1. 今日食谱 (Big Card) -->
            <div
              class="main-card big-card recipe-card"
              @click="goToRecipes"
              :class="{ disabled: !authStore.hasProfile }"
            >
              <div class="card-content">
                <div class="card-header">
                  <span class="card-emoji">🍳</span>
                  <h2 class="card-title-big">今日食谱</h2>
                </div>
                <p class="card-desc">
                  {{
                    authStore.hasProfile
                      ? "获取推荐的营养均衡食谱"
                      : "完善档案后开启智能推荐"
                  }}
                </p>
                <div class="card-action">
                  <span class="action-text"
                    >查看今日推荐 <el-icon><ArrowRight /></el-icon
                  ></span>
                </div>
              </div>
            </div>

            <!-- 2. 营养分析 (Big Card - 原个人档案升级) -->
            <div
              class="main-card big-card nutrition-card"
              @click="goToNutrition"
            >
              <div class="card-content">
                <div class="card-header">
                  <span class="card-emoji">📊</span>
                  <div class="header-text">
                    <h2 class="card-title-big">营养分析</h2>
                    <span class="card-subtitle">目标追踪 & 档案管理</span>
                  </div>
                </div>

                <div
                  v-if="authStore.hasProfile && nutritionData"
                  class="nutrition-preview"
                >
                  <div class="target-cal">
                    <span class="label">今日目标</span>
                    <span class="value">{{
                      nutritionData.target_calorie || "--"
                    }}</span>
                    <span class="unit">kcal</span>
                  </div>
                  <div class="macros-mini">
                    <div class="macro-mini-item">
                      <span class="macro-label">碳水</span>
                      <span class="macro-val"
                        >{{ nutritionData.carb_ratio?.toFixed(0) }}%</span
                      >
                    </div>
                    <div class="macro-mini-item">
                      <span class="macro-label">蛋白</span>
                      <span class="macro-val"
                        >{{ nutritionData.protein_ratio?.toFixed(0) }}%</span
                      >
                    </div>
                    <div class="macro-mini-item">
                      <span class="macro-label">脂肪</span>
                      <span class="macro-val"
                        >{{ nutritionData.fat_ratio?.toFixed(0) }}%</span
                      >
                    </div>
                  </div>
                </div>
                <p v-else class="card-desc">完善档案，获取专属营养分析</p>

                <div class="card-status" v-if="!authStore.hasProfile">
                  <el-tag type="warning" size="small" effect="plain"
                    >待完善</el-tag
                  >
                </div>
              </div>
            </div>
          </div>

          <!-- 次级功能卡片区 (小卡片网格) -->
          <div class="secondary-grid">
            <!-- 饮食记录 -->
            <div class="main-card small-card" @click="goToIntake">
              <div class="card-content-center">
                <span class="card-emoji-small">📝</span>
                <h3 class="card-title-small">饮食记录</h3>
                <p class="card-desc-mini">记录每餐摄入</p>
              </div>
            </div>

            <!-- 周报告 -->
            <div class="main-card small-card" @click="goToWeekly">
              <div class="card-content-center">
                <span class="card-emoji-small">📈</span>
                <h3 class="card-title-small">周报告</h3>
                <p class="card-desc-mini">查看长期趋势</p>
              </div>
            </div>

            <!-- 收藏夹 -->
            <div class="main-card small-card" @click="goToFavorites">
              <div class="card-content-center">
                <span class="card-emoji-small">⭐</span>
                <h3 class="card-title-small">我的收藏</h3>
                <p class="card-desc-mini">喜爱的食谱</p>
              </div>
            </div>

            <!-- 购物清单 -->
            <div class="main-card small-card" @click="goToShopping">
              <div class="card-content-center">
                <span class="card-emoji-small">🛒</span>
                <h3 class="card-title-small">购物清单</h3>
                <p class="card-desc-mini">食材采购助手</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 提示信息 (当没有档案时显示) -->
        <div class="tip-card" v-if="!authStore.hasProfile">
          <div class="tip-icon">💡</div>
          <div class="tip-content">
            <h4 class="tip-title">开始您的健康之旅</h4>
            <p class="tip-text">完善个人档案，获取专属营养方案</p>
            <el-button type="primary" size="small" @click="goToProfile" round>
              立即完善
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  UserFilled,
  User,
  ArrowDown,
  ArrowRight,
  SwitchButton,
  Document,
  KnifeFork,
  TrendCharts,
  DataLine,
  Medal,
  PieChart,
  Star,
} from "@element-plus/icons-vue";
import { useAuthStore } from "@/store/auth";
import { getNutritionRequirements } from "@/api/user";

const router = useRouter();
const authStore = useAuthStore();
const nutritionData = ref(null);

onMounted(async () => {
  if (authStore.isAuthenticated) {
    if (!authStore.user || !authStore.profile) {
      try {
        await authStore.loadProfile();
      } catch (error) {
        console.error("加载档案失败:", error);
      }
    }

    // 加载营养需求数据
    if (authStore.hasProfile) {
      try {
        const response = await getNutritionRequirements();
        nutritionData.value = response;
      } catch (error) {
        console.log("营养需求数据不可用:", error);
      }
    }
  }
});

const goToProfile = () => {
  router.push("/profile/view");
};

const goToRecipes = () => {
  if (!authStore.hasProfile) {
    ElMessage.warning("请先完善健康档案");
    return;
  }
  router.push("/recipes");
};

const goToIntake = () => {
  router.push("/intake");
};

const goToWeekly = () => {
  router.push("/weekly");
};

const goToNutrition = () => {
  if (!authStore.hasProfile) {
    ElMessage.warning("请先完善健康档案");
    return;
  }
  // 可以创建一个专门的营养需求展示页面
  router.push("/profile/view");
};

const goToFavorites = () => {
  router.push("/favorites");
};

const goToShopping = () => {
  console.log("点击购物清单卡片");
  router.push("/shopping");
};

// 获取问候语
const getGreeting = () => {
  const hour = new Date().getHours();
  if (hour < 6) return "夜深了";
  if (hour < 9) return "早上好";
  if (hour < 12) return "上午好";
  if (hour < 14) return "中午好";
  if (hour < 18) return "下午好";
  if (hour < 22) return "晚上好";
  return "夜深了";
};

const handleCommand = async (command) => {
  if (command === "profile") {
    router.push("/profile/view");
  } else if (command === "favorites") {
    goToFavorites();
  } else if (command === "intake") {
    goToIntake();
  } else if (command === "weekly") {
    goToWeekly();
  } else if (command === "logout") {
    try {
      await ElMessageBox.confirm("确定要退出登录吗？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      });
      authStore.logout();
      ElMessage.success("已退出登录");
      router.push("/login");
    } catch {
      // 用户取消
    }
  }
};
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: transparent;
  position: relative;
}

/* 顶部导航栏 */
.top-bar {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid #e8eaed;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.top-bar-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.brand-icon {
  font-size: 24px;
}

.brand-name {
  font-size: 20px;
  color: #1a1a1a;
  letter-spacing: -0.5px;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s;
}

.user-menu:hover {
  background: #f0f2f5;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
}

.dropdown-icon {
  font-size: 14px;
  color: #8b95a5;
}

/* 主内容包装 */
.main-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

/* Hero区域 */
.hero-section {
  background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
  border-radius: 20px;
  padding: 40px;
  margin-bottom: 32px;
}

.hero-content {
  margin-bottom: 24px;
}

.hero-title {
  font-size: 32px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px 0;
  line-height: 1.2;
}

.highlight {
  color: #10b981;
}

.hero-subtitle {
  font-size: 16px;
  color: #6b7280;
  margin: 0;
}

/* 快速统计 */
.quick-stats {
  display: flex;
  gap: 24px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon {
  font-size: 32px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #6b7280;
  margin-top: 4px;
}

.stat-divider {
  width: 1px;
  background: #e5e7eb;
}

/* 功能区布局 */
.content-section {
  animation: fadeIn 0.5s ease-out;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.features-grid {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 顶部大卡片行 */
.big-cards-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

/* 次级小卡片网格 */
.secondary-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

/* 卡片通用样式 */
.main-card {
  background: #ffffff;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #eef2f6;
  position: relative;
  overflow: hidden;
}

.main-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.08);
  border-color: #10b981;
}

.main-card:active {
  transform: scale(0.98);
}

/* 大卡片样式 */
.big-card {
  height: 220px;
  padding: 28px;
  /* background: linear-gradient(145deg, #ffffff 0%, #f9fafb 100%); */
}

.recipe-card {
  background: linear-gradient(135deg, #e6fffa 0%, #ffffff 100%);
  border-color: #b2f5ea;
}

.nutrition-card {
  background: linear-gradient(135deg, #ebf8ff 0%, #ffffff 100%);
  border-color: #bee3f8;
}

.nutrition-card:hover {
  border-color: #3182ce;
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 16px;
}

.header-text {
  display: flex;
  flex-direction: column;
}

.card-subtitle {
  font-size: 13px;
  color: #64748b;
  margin-top: 4px;
}

.card-emoji {
  font-size: 42px;
  line-height: 1;
}

.card-title-big {
  font-size: 24px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  line-height: 1.2;
}

.card-desc {
  font-size: 15px;
  color: #64748b;
  line-height: 1.6;
  margin-bottom: auto;
}

.card-action {
  margin-top: auto;
  color: #10b981;
  font-weight: 600;
  font-size: 14px;
  display: flex;
  align-items: center;
}

.action-text {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 营养数据预览 */
.nutrition-preview {
  margin-top: 10px;
}

.target-cal {
  display: flex;
  align-items: baseline;
  gap: 6px;
  margin-bottom: 12px;
}

.target-cal .label {
  font-size: 13px;
  color: #64748b;
}

.target-cal .value {
  font-size: 28px;
  font-weight: 800;
  color: #2c5282;
  font-family: "DIN Alternate", sans-serif;
}

.target-cal .unit {
  font-size: 13px;
  color: #64748b;
}

.macros-mini {
  display: flex;
  gap: 12px;
}

.macro-mini-item {
  background: rgba(255, 255, 255, 0.6);
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.macro-label {
  color: #94a3b8;
  margin-right: 4px;
}

.macro-val {
  color: #334155;
  font-weight: 600;
}

/* 小卡片样式 */
.small-card {
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  height: 140px;
}

.card-content-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.card-emoji-small {
  font-size: 32px;
  margin-bottom: 4px;
}

.card-title-small {
  font-size: 16px;
  font-weight: 600;
  color: #334155;
  margin: 0;
}

.card-desc-mini {
  font-size: 12px;
  color: #94a3b8;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .secondary-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .top-bar-content {
    padding: 12px 16px;
  }

  .main-wrapper {
    padding: 16px;
  }

  .hero-section {
    padding: 24px;
    margin-bottom: 24px;
  }

  .hero-title {
    font-size: 24px;
  }

  /* 移动端改为单列 */
  .big-cards-row {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .secondary-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .big-card {
    height: auto;
    min-height: 180px;
  }

  .small-card {
    height: 120px;
    padding: 16px;
  }

  .quick-stats {
    flex-direction: column;
    gap: 16px;
  }

  .stat-divider {
    display: none;
  }

  .tip-card {
    flex-direction: column;
    text-align: center;
  }
}

@media (max-width: 480px) {
  /* 极窄屏幕 */
  .hero-title {
    font-size: 20px;
  }

  .stat-value {
    font-size: 20px;
  }
}
</style>
