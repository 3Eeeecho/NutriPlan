<template>
  <div class="home-container">
    <!-- 顶部导航 - 使用新组件 -->
    <TopNavigation />

    <!-- 主内容区 -->
    <div class="main-wrapper">
      <!-- 欢迎区域 -->
      <n-card class="hero-card" :bordered="false">
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
          <StatCard
            icon="🔥"
            :value="authStore.profile.tdee ? Math.floor(authStore.profile.tdee).toString() : '--'"
            label="每日目标热量"
            subtitle="kcal"
          />
          <div class="stat-divider"></div>
          <StatCard
            icon="⚖️"
            :value="authStore.profile.weight?.toString() || '--'"
            label="当前体重"
            subtitle="kg"
          />
        </div>
      </n-card>

      <!-- 功能卡片区 -->
      <div class="content-section">
        <div class="features-grid">
          <!-- 主要功能卡片区 (顶层双巨头) -->
          <div class="big-cards-row">
            <!-- 1. 今日食谱 (Big Card) -->
            <n-card
              class="feature-card big-card recipe-card"
              hoverable
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
                  <span class="action-text">
                    查看今日推荐 
                    <n-icon><ArrowForwardOutline /></n-icon>
                  </span>
                </div>
              </div>
            </n-card>

            <!-- 2. 营养分析 (Big Card - 原个人档案升级) -->
            <n-card
              class="feature-card big-card nutrition-card"
              hoverable
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
                    <n-tag round size="small" :bordered="false" class="macro-tag">
                      碳水 {{ nutritionData.carb_ratio?.toFixed(0) }}%
                    </n-tag>
                    <n-tag round size="small" :bordered="false" class="macro-tag">
                      蛋白 {{ nutritionData.protein_ratio?.toFixed(0) }}%
                    </n-tag>
                    <n-tag round size="small" :bordered="false" class="macro-tag">
                      脂肪 {{ nutritionData.fat_ratio?.toFixed(0) }}%
                    </n-tag>
                  </div>
                </div>
                <p v-else class="card-desc">完善档案，获取专属营养分析</p>

                <div class="card-status" v-if="!authStore.hasProfile">
                  <n-tag type="warning" size="small" round>待完善</n-tag>
                </div>
              </div>
            </n-card>
          </div>

          <!-- 次级功能卡片区 (小卡片网格) -->
          <div class="secondary-grid">
            <!-- 饮食记录 -->
            <n-card class="feature-card small-card" hoverable @click="goToIntake">
              <div class="card-content-center">
                <span class="card-emoji-small">📝</span>
                <h3 class="card-title-small">饮食记录</h3>
                <p class="card-desc-mini">记录每餐摄入</p>
              </div>
            </n-card>

            <!-- 周报告 -->
            <n-card class="feature-card small-card" hoverable @click="goToWeekly">
              <div class="card-content-center">
                <span class="card-emoji-small">📈</span>
                <h3 class="card-title-small">周报告</h3>
                <p class="card-desc-mini">查看长期趋势</p>
              </div>
            </n-card>

            <!-- 收藏夹 -->
            <n-card class="feature-card small-card" hoverable @click="goToFavorites">
              <div class="card-content-center">
                <span class="card-emoji-small">⭐</span>
                <h3 class="card-title-small">我的收藏</h3>
                <p class="card-desc-mini">喜爱的食谱</p>
              </div>
            </n-card>

            <!-- 购物清单 -->
            <n-card class="feature-card small-card" hoverable @click="goToShopping">
              <div class="card-content-center">
                <span class="card-emoji-small">🛒</span>
                <h3 class="card-title-small">购物清单</h3>
                <p class="card-desc-mini">食材采购助手</p>
              </div>
            </n-card>
          </div>
        </div>

        <!-- 提示信息 (当没有档案时显示) -->
        <n-card v-if="!authStore.hasProfile" class="tip-card" :bordered="false">
          <div class="tip-icon">💡</div>
          <div class="tip-content">
            <h4 class="tip-title">开始您的健康之旅</h4>
            <p class="tip-text">完善个人档案，获取专属营养方案</p>
            <n-button type="primary" size="small" round @click="goToProfile">
              立即完善
            </n-button>
          </div>
        </n-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { NCard, NTag, NButton, NIcon } from "naive-ui";
import { ArrowForwardOutline } from "@vicons/ionicons5";
import { ElMessage } from "element-plus";
import { useAuthStore } from "@/store/auth";
import { getNutritionRequirements } from "@/api/user";
import TopNavigation from "@/components/layout/TopNavigation.vue";
import StatCard from "@/components/ui/StatCard.vue";

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
  background: linear-gradient(135deg, #f0fdf4 0%, #e0f2fe 50%, #fef3c7 100%);
}

/* 主内容包装 */
.main-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--spacing-2xl) var(--spacing-lg);
}

/* Hero卡片 */
.hero-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-2xl);
  margin-bottom: var(--spacing-2xl);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.8);
}

.hero-content {
  margin-bottom: var(--spacing-lg);
}

.hero-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-xs) 0;
  line-height: 1.2;
}

.highlight {
  color: var(--color-primary);
}

.hero-subtitle {
  font-size: var(--font-size-base);
  color: var(--text-secondary);
  margin: 0;
}

/* 快速统计 */
.quick-stats {
  display: flex;
  gap: var(--spacing-lg);
  padding-top: var(--spacing-lg);
  border-top: 1px solid var(--border-color);
  align-items: center;
}

.stat-divider {
  width: 1px;
  height: 60px;
  background: var(--border-color);
}

/* 功能区布局 */
.content-section {
  animation: fadeIn 0.5s ease-out;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.features-grid {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

/* 顶部大卡片行 */
.big-cards-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-lg);
}

/* 次级小卡片网格 */
.secondary-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--spacing-md);
}

/* 卡片通用样式 */
.feature-card {
  cursor: pointer;
  transition: all var(--transition-normal);
  border-radius: var(--radius-xl);
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

.feature-card:active {
  transform: scale(0.98);
}

/* 大卡片样式 */
.big-card {
  min-height: 220px;
}

.big-card :deep(.n-card__content) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.recipe-card {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.recipe-card:hover {
  background: rgba(255, 255, 255, 0.75);
  border-color: var(--color-primary);
}

.nutrition-card {
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.nutrition-card:hover {
  background: rgba(255, 255, 255, 0.75);
  border-color: #3b82f6;
}

.card-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.header-text {
  display: flex;
  flex-direction: column;
}

.card-subtitle {
  font-size: var(--font-size-sm);
  color: var(--text-tertiary);
  margin-top: var(--spacing-xs);
}

.card-emoji {
  font-size: 42px;
  line-height: 1;
}

.card-title-big {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0;
  line-height: 1.2;
}

.card-desc {
  font-size: var(--font-size-base);
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: auto;
}

.card-action {
  margin-top: auto;
  color: var(--color-primary);
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  display: flex;
  align-items: center;
}

.action-text {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

/* 营养数据预览 */
.nutrition-preview {
  margin-top: var(--spacing-sm);
}

.target-cal {
  display: flex;
  align-items: baseline;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-md);
}

.target-cal .label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.target-cal .value {
  font-size: 28px;
  font-weight: var(--font-weight-bold);
  color: #2c5282;
  font-family: "DIN Alternate", sans-serif;
}

.target-cal .unit {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.macros-mini {
  display: flex;
  gap: var(--spacing-sm);
}

.macro-tag {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

/* 小卡片样式 */
.small-card {
  min-height: 140px;
  background: rgba(255, 255, 255, 0.5);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.6);
}

.small-card:hover {
  background: rgba(255, 255, 255, 0.7);
  border-color: var(--color-primary);
}

.small-card :deep(.n-card__content) {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-content-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-xs);
  text-align: center;
}

.card-emoji-small {
  font-size: 32px;
  margin-bottom: var(--spacing-xs);
}

.card-title-small {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.card-desc-mini {
  font-size: var(--font-size-sm);
  color: var(--text-tertiary);
  margin: 0;
}

/* 提示卡片 */
.tip-card {
  background: rgba(254, 243, 199, 0.5);
  backdrop-filter: blur(8px);
  border-radius: var(--radius-xl);
  box-shadow: 0 4px 20px rgba(251, 191, 36, 0.15);
  border: 1px solid rgba(251, 191, 36, 0.2);
}

.tip-card :deep(.n-card__content) {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.tip-icon {
  font-size: 48px;
}

.tip-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.tip-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.tip-text {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin: 0;
}

.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .secondary-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .main-wrapper {
    padding: var(--spacing-md);
  }

  .hero-title {
    font-size: var(--font-size-2xl);
  }

  /* 移动端改为单列 */
  .big-cards-row {
    grid-template-columns: 1fr;
    gap: var(--spacing-md);
  }

  .secondary-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-sm);
  }

  .big-card {
    min-height: 180px;
  }

  .small-card {
    min-height: 120px;
  }

  .quick-stats {
    flex-direction: column;
    gap: var(--spacing-md);
    align-items: stretch;
  }

  .stat-divider {
    display: none;
  }

  .tip-card :deep(.n-card__content) {
    flex-direction: column;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .hero-title {
    font-size: var(--font-size-xl);
  }
}
</style>
