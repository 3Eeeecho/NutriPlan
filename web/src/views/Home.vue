<template>
  <div v-if="authStore.isAuthenticated" class="home-container grid-layout">
    <!-- 顶部导航 -->
    <TopNavigation />

    <!-- 主内容区：使用 Tailwind Grid 重构，打破传统死板行结构 -->
    <div class="main-wrapper">
      <div class="dashboard-grid">
        
        <!-- 视觉第一层级：核心状态与高频操作 -->
        
        <!-- 概览卡片（大占比）：问候语 + 热量环形图 -->
        <div class="overview-card grid-card">
          <div class="overview-content">
            <div class="greeting-section">
              <h1 class="greeting-text">
                {{ getGreeting() }}<br/><span class="username">{{ authStore.user?.username }}</span>
              </h1>
              <p class="greeting-subtitle font-sans">{{ getCurrentDateText() }}</p>
            </div>
            
            <div class="calorie-section">
              <div v-if="authStore.hasProfile && authStore.profile" class="calorie-ring-container">
                <!-- 达标完成状态 -->
                <div v-if="isGoalAchieved" class="goal-achieved">
                  <div class="achievement-icon">✓</div>
                  <div class="achievement-text">
                    <div class="achievement-title font-serif">目标达成</div>
                    <div class="achievement-subtitle">{{ calorieTarget }} kcal</div>
                  </div>
                </div>
                <!-- 未达标进度状态 -->
                <div v-else class="progress-container">
                  <n-progress
                    type="circle"
                    :percentage="caloriePercentage"
                    :stroke-width="8"
                    :color="calorieColor"
                    :rail-color="'rgba(0,0,0,0.05)'"
                    :show-indicator="true"
                    :style="{ width: '120px', height: '120px' }"
                  >
                    <div class="ring-inner">
                      <span class="ring-number font-serif">{{ caloriesRemaining }}</span>
                      <span class="ring-label">kcal 剩余</span>
                    </div>
                  </n-progress>
                  <div class="progress-details mt-4 text-sm text-gray-500">
                    已摄入 <span class="font-serif font-semibold text-gray-800">{{ caloriesConsumed }}</span> 
                  </div>
                </div>
              </div>
              <div v-else class="empty-ring">
                <n-button type="primary" size="large" class="organic-btn" @click="goToProfile">
                  完善档案
                </n-button>
              </div>
            </div>
          </div>
        </div>

        <!-- 记录卡片（强调色）：独立卡片，纵向鲜明视觉 -->
        <div class="record-card grid-card highlight-card cursor-pointer group" @click="goToIntake">
          <div class="card-glass-effect"></div>
          <div class="record-content">
            <div class="record-icon group-hover:scale-110 transition-transform duration-500">
               +
            </div>
            <h3 class="record-title font-serif">饮食记录</h3>
            <p class="record-desc font-sans text-sm opacity-80 mt-2">记录每餐摄入，追踪您的健康轨迹</p>
          </div>
        </div>

        <!-- 视觉第二层级：内容推荐与数据统计 -->
        
        <!-- 食谱卡片（宽幅大图）：今日营养食谱 -->
        <div class="recipe-card grid-card has-bg cursor-pointer group" @click="goToRecipes" :class="{ 'opacity-50 pointer-events-none': !authStore.hasProfile }">
          <div class="recipe-bg-image group-hover:scale-105 transition-transform duration-700"></div>
          <div class="recipe-overlay"></div>
          <div class="recipe-content relative z-10 flex flex-col justify-end h-full p-6">
            <div class="recipe-tag mb-4">
              <span class="editorial-tag">Editor's Pick</span>
            </div>
            <h2 class="recipe-title font-serif text-3xl text-white mb-2">今日营养食谱</h2>
            <p class="recipe-description font-sans text-white/80 max-w-xs">
              {{ authStore.hasProfile ? '探索精心调配的地中海风味与全谷物组合。' : '完善档案后为您推荐。' }}
            </p>
          </div>
        </div>

        <!-- 报告卡片（数据微展示）：周报告微型折线图 -->
        <div class="report-card grid-card cursor-pointer group" @click="goToWeekly">
          <div class="report-header flex justify-between items-start mb-4">
            <h3 class="report-title font-serif text-xl">周报告</h3>
            <span class="report-icon text-gray-400 group-hover:text-gray-800 transition-colors">↗</span>
          </div>
          <div class="report-sparkline mt-auto h-16 w-full flex items-end gap-1 opacity-70 group-hover:opacity-100 transition-opacity">
            <div class="bar h-[40%]"></div>
            <div class="bar h-[60%]"></div>
            <div class="bar h-[30%]"></div>
            <div class="bar h-[80%]"></div>
            <div class="bar h-[50%]"></div>
            <div class="bar h-[90%]"></div>
            <div class="bar h-[70%] active"></div>
          </div>
          <p class="mt-4 text-xs text-gray-500 font-sans tracking-wide uppercase">TRENDS & INSIGHTS</p>
        </div>

        <!-- 视觉第三层级：辅助工具 -->
        
        <!-- 辅助功能卡片：小方块并排 -->
        <div class="aux-container grid-card-transparent">
          <div class="aux-grid">
            <div class="aux-card cursor-pointer group" @click="goToFavorites">
              <div class="aux-icon text-amber-600 mb-3 group-hover:-translate-y-1 transition-transform">♥</div>
              <h4 class="font-serif text-lg">我的收藏</h4>
            </div>
            <div class="aux-card cursor-pointer group" @click="goToShopping">
              <div class="aux-icon text-emerald-600 mb-3 group-hover:-translate-y-1 transition-transform">🛒</div>
              <h4 class="font-serif text-lg">购物清单</h4>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
  <HeroSection v-else />
</template>

<script setup>
import { onMounted, ref, computed } from "vue";
import { useRouter } from "vue-router";
import { NCard, NGrid, NGi, NTag, NButton, NIcon, NProgress, useMessage } from "naive-ui";
import { ArrowForwardOutline } from "@vicons/ionicons5";
import { useAuthStore } from "@/store/auth";
import { getNutritionRequirements } from "@/api/user";
import { getTodayStatus } from "@/api/intakeApi";
import TopNavigation from "@/components/layout/TopNavigation.vue";
import HeroSection from "@/components/HeroSection.vue";

const router = useRouter();
const message = useMessage();
const authStore = useAuthStore();
const nutritionData = ref(null);
const todayIntake = ref(null);

onMounted(async () => {
  if (authStore.isAuthenticated) {
    if (!authStore.user || !authStore.profile) {
      try {
        await authStore.loadProfile();
      } catch (error) {
        console.error("加载档案失败:", error);
      }
    }

    // 加载营养需求数据和今日摄入
    if (authStore.hasProfile) {
      try {
        const [reqRes, statusRes] = await Promise.all([
          getNutritionRequirements(),
          getTodayStatus()
        ]);
        nutritionData.value = reqRes;
        todayIntake.value = statusRes;
      } catch (error) {
        console.log("数据加载失败:", error);
      }
    }
  }
});

const goToProfile = () => {
  router.push("/profile/view");
};

const goToRecipes = () => {
  if (!authStore.hasProfile) {
    message.warning("请先完善健康档案");
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
    message.warning("请先完善健康档案");
    return;
  }
  router.push("/profile/view");
};

const goToFavorites = () => {
  router.push("/favorites");
};

const goToShopping = () => {
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

// 获取当前日期文本
const getCurrentDateText = () => {
  const today = new Date();
  const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
  return today.toLocaleDateString('zh-CN', options);
};

// 计算属性
const calorieTarget = computed(() => {
  if (nutritionData.value && nutritionData.value.target_calorie) {
    return Math.floor(nutritionData.value.target_calorie);
  }
  return Math.floor(authStore.profile?.tdee || 0);
});

const caloriesConsumed = computed(() => Math.floor(todayIntake.value?.total_energy || 0));

const caloriesRemaining = computed(() => {
  if (!calorieTarget.value) return 0;
  return Math.max(0, calorieTarget.value - caloriesConsumed.value);
});

const caloriePercentage = computed(() => {
  if (!calorieTarget.value) return 0;
  return Math.min(100, Math.max(0, (caloriesConsumed.value / calorieTarget.value) * 100));
});

const calorieColor = computed(() => {
  const p = caloriePercentage.value;
  if (p >= 100) return '#ef4444'; // 红色（超标）
  if (p > 80) return '#f59e0b'; // 黄色（接近）
  return '#10b981'; // 绿色（健康）
});

// 判断是否达成目标（摄入量在目标的90%-110%之间）
const isGoalAchieved = computed(() => {
  if (!calorieTarget.value) return false;
  const percentage = caloriePercentage.value;
  return percentage >= 90 && percentage <= 110;
});

// Remove duplicated function definitions below this line if any


</script>

<style scoped>
/* =========== Grid 布局体系 =========== */
.grid-layout {
  min-height: 100vh;
  /* 使用全局的有机底色背景，去掉原生渐变 */
  background: transparent;
}

.main-wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  grid-auto-rows: minmax(100px, auto);
  gap: 1.5rem;
}

/* 基础卡片样式: 有机、微边框、大倒角 */
.grid-card {
  background: var(--color-bg-primary);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border-primary);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  position: relative;
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.grid-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

/* ================== 卡片分配 ================== */

/* 1. 概览卡片 (左上，大尺寸) */
.overview-card {
  grid-column: span 12;
  min-height: 240px;
  background: linear-gradient(135deg, var(--color-bg-primary) 0%, var(--color-primary-supperlight) 100%);
}
@media (min-width: 1024px) {
  .overview-card {
    grid-column: span 8;
  }
}

.overview-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 2.5rem;
}
@media (min-width: 768px) {
  .overview-content {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}

.greeting-text {
  font-family: var(--font-family-display);
  font-size: 2.5rem;
  line-height: 1.2;
  color: var(--color-text-primary);
  margin-bottom: 0.5rem;
}

.username {
  font-style: italic;
  color: var(--color-primary);
}

/* 2. 补签/记录卡片 (右上，高亮) */
.record-card {
  grid-column: span 12;
  background: var(--color-accent-warm);
  color: #fff;
  border: none;
  min-height: 240px;
  position: relative;
}
@media (min-width: 1024px) {
  .record-card {
    grid-column: span 4;
  }
}

.card-glass-effect {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: linear-gradient(145deg, rgba(255,255,255,0.1) 0%, rgba(0,0,0,0.05) 100%);
  pointer-events: none;
}

.record-content {
  position: relative;
  z-index: 2;
  padding: 2rem;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
}

.record-icon {
  font-size: 3rem;
  font-weight: 300;
  margin-bottom: 1rem;
  color: rgba(255,255,255,0.9);
}

.record-title {
  font-size: 1.75rem;
}

/* 3. 食谱卡片 (中左，宽幅杂志风) */
.recipe-card {
  grid-column: span 12;
  min-height: 320px;
}
@media (min-width: 1024px) {
  .recipe-card {
    grid-column: span 8;
  }
}

.recipe-bg-image {
  position: absolute;
  inset: 0;
  background-image: url('https://images.unsplash.com/photo-1490645935967-10de6ba17061?q=80&w=2053&auto=format&fit=crop');
  background-size: cover;
  background-position: center;
}

.recipe-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(30,40,30,0.85) 0%, rgba(0,0,0,0.1) 60%, transparent 100%);
}

.editorial-tag {
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(8px);
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #fff;
  border: 1px solid rgba(255,255,255,0.3);
}

/* 4. 报告卡片 (中右，火花线趋势) */
.report-card {
  grid-column: span 12;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-primary);
  border: 1px solid var(--color-border-primary);
}
@media (min-width: 768px) {
  .report-card {
    grid-column: span 6;
  }
}
@media (min-width: 1024px) {
  .report-card {
    grid-column: span 4;
  }
}

.report-sparkline {
  align-items: flex-end;
}

.report-sparkline .bar {
  flex: 1;
  background: var(--color-border-primary);
  border-radius: 2px 2px 0 0;
  transition: height 0.5s ease-in-out, background 0.3s;
}

.report-sparkline .bar:hover {
  background: var(--color-primary-light);
}

.report-sparkline .bar.active {
  background: var(--color-accent-cool);
}

/* 5. 辅助功能 (下右，小巧双格) */
.aux-container {
  grid-column: span 12;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}
@media (min-width: 768px) {
  .aux-container {
    grid-column: span 6;
  }
}
@media (min-width: 1024px) {
  .aux-container {
    grid-column: span 4;
  }
}

.aux-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  height: 100%;
}

.aux-card {
  background: var(--color-bg-primary);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

.aux-card:hover {
  background: var(--color-bg-secondary);
  border-color: var(--color-primary-lighter);
}

.aux-icon {
  font-size: 2rem;
}

/* ==================================================== */

/* Optional responsive cleanups underneath */
@media (max-width: 1024px) {
  /* Removed old ring container scale */
}

@media (max-width: 768px) {
  .main-wrapper {
    padding: 1rem;
  }
}

/* Animation */
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

.home-container > * {
  animation: fadeIn 0.5s ease-out;
}
</style>
