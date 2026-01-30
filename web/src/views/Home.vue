<template>
  <div class="home-container">
    <!-- 顶部导航 -->
    <TopNavigation />

    <!-- 主内容区 -->
    <div class="main-wrapper">
      <n-grid :x-gap="16" :y-gap="16" cols="1">
        <!-- Hero Card: Greeting + Calorie Ring -->
        <n-gi>
          <n-card :bordered="false" class="hero-card">
            <n-grid :x-gap="24" :cols="2" responsive="screen" :collapsed-rows="1" item-responsive>
              <!-- Left: Greeting -->
              <n-gi span="2 m:1">
                <div class="greeting-section">
                  <h1 class="greeting-text">
                    {{ getGreeting() }}，<span class="username">{{ authStore.user?.username }}</span> 👋
                  </h1>
                  <p class="greeting-subtitle">{{ getCurrentDateText() }}</p>
                </div>
              </n-gi>

              <!-- Right: Calorie Ring -->
              <n-gi span="2 m:1">
                <div class="calorie-ring-container" v-if="authStore.hasProfile && authStore.profile">
                  <n-progress
                    type="circle"
                    :percentage="caloriePercentage"
                    :stroke-width="12"
                    :color="calorieColor"
                    :rail-color="'#f1f5f9'"
                    :show-indicator="true"
                    :style="{ width: '130px', height: '130px' }"
                  >
                    <div class="ring-inner">
                      <span class="ring-number" :style="{ color: calorieColor }">{{ caloriesRemaining }}</span>
                      <span class="ring-label">kcal 剩余</span>
                    </div>
                  </n-progress>
                  <div class="ring-footer">
                    已摄入 {{ caloriesConsumed }} / 目标 {{ calorieTarget }}
                  </div>
                </div>
                <div v-else class="empty-ring">
                  <n-button type="primary" size="large" @click="goToProfile">
                    完善档案开始
                  </n-button>
                </div>
              </n-gi>
            </n-grid>
          </n-card>
        </n-gi>

        <!-- Today's Recipe Card: Magazine Cover Style -->
        <n-gi>
          <n-card 
            :bordered="false" 
            class="recipe-magazine-card" 
            hoverable
            @click="goToRecipes"
            :class="{ disabled: !authStore.hasProfile }"
          >
            <div class="recipe-bg-image"></div>
            <div class="recipe-overlay"></div>
            <div class="recipe-content">
              <div class="recipe-tag">
                <n-tag type="success" round size="small">今日推荐</n-tag>
              </div>
              <h2 class="recipe-title">今日营养食谱</h2>
              <p class="recipe-description">
                {{ authStore.hasProfile ? '为您精心挑选的营养均衡食谱' : '完善档案后开启智能推荐' }}
              </p>
              <div class="recipe-action">
                <!-- Arrow removed -->
              </div>
            </div>
          </n-card>
        </n-gi>

        <!-- Action Cards Grid: 4 columns -->
        <n-gi>
          <n-grid :x-gap="16" :y-gap="16" :cols="'2 m:4'" responsive="screen">
            <!-- 饮食记录 -->
            <n-gi>
              <n-card 
                :bordered="false" 
                class="action-card" 
                hoverable
                @click="goToIntake"
              >
                <div class="action-icon-wrapper bg-red-100">
                  <n-icon size="36" color="#dc2626">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M21 6.5c0 1.38-1.12 2.5-2.5 2.5S16 7.88 16 6.5 17.12 4 18.5 4 21 5.12 21 6.5zM9 4h6v2H9V4zm11 6h-2c0 1.1-.9 2-2 2h-5v11c0 .55-.45 1-1 1s-1-.45-1-1V12H4c-1.1 0-2-.9-2-2H0V6h2c0-1.1.9-2 2-2h1V2c0-.55.45-1 1-1s1 .45 1 1v2h1c1.1 0 2 .9 2 2v4h8V6c0-1.1.9-2 2-2h1V2c0-.55.45-1 1-1s1 .45 1 1v2h1c1.1 0 2 .9 2 2v4z"/>
                    </svg>
                  </n-icon>
                </div>
                <h3 class="action-title">饮食记录</h3>
                <p class="action-desc">记录每餐摄入</p>
              </n-card>
            </n-gi>

            <!-- 周报告 -->
            <n-gi>
              <n-card 
                :bordered="false" 
                class="action-card" 
                hoverable
                @click="goToWeekly"
              >
                <div class="action-icon-wrapper bg-blue-100">
                  <n-icon size="36" color="#2563eb">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M3 13h2v-2H3v2zm0 4h2v-2H3v2zm0-8h2V7H3v2zm4 4h14v-2H7v2zm0 4h14v-2H7v2zM7 7v2h14V7H7z"/>
                    </svg>
                  </n-icon>
                </div>
                <h3 class="action-title">周报告</h3>
                <p class="action-desc">查看长期趋势</p>
              </n-card>
            </n-gi>

            <!-- 我的收藏 -->
            <n-gi>
              <n-card 
                :bordered="false" 
                class="action-card" 
                hoverable
                @click="goToFavorites"
              >
                <div class="action-icon-wrapper bg-yellow-100">
                  <n-icon size="36" color="#ca8a04">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"/>
                    </svg>
                  </n-icon>
                </div>
                <h3 class="action-title">我的收藏</h3>
                <p class="action-desc">喜爱的食谱</p>
              </n-card>
            </n-gi>

            <!-- 购物清单 -->
            <n-gi>
              <n-card 
                :bordered="false" 
                class="action-card" 
                hoverable
                @click="goToShopping"
              >
                <div class="action-icon-wrapper bg-green-100">
                  <n-icon size="36" color="#16a34a">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M7 18c-1.1 0-1.99.9-1.99 2S5.9 22 7 22s2-.9 2-2-.9-2-2-2zM1 2v2h2l3.6 7.59-1.35 2.45c-.16.28-.25.61-.25.96 0 1.1.9 2 2 2h12v-2H7.42c-.14 0-.25-.11-.25-.25l.03-.12.9-1.63h7.45c.75 0 1.41-.41 1.75-1.03l3.58-6.49c.08-.14.12-.31.12-.48 0-.55-.45-1-1-1H5.21l-.94-2H1zm16 16c-1.1 0-1.99.9-1.99 2s.89 2 1.99 2 2-.9 2-2-.9-2-2-2z"/>
                    </svg>
                  </n-icon>
                </div>
                <h3 class="action-title">购物清单</h3>
                <p class="action-desc">食材采购助手</p>
              </n-card>
            </n-gi>
          </n-grid>
        </n-gi>

        <!-- Tip Card (if no profile) -->
        <n-gi v-if="!authStore.hasProfile">
          <n-card :bordered="false" class="tip-card">
            <div class="tip-content">
              <n-icon size="48" color="#f59e0b">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M9 21c0 .55.45 1 1 1h4c.55 0 1-.45 1-1v-1H9v1zm3-19C8.14 2 5 5.14 5 9c0 2.38 1.19 4.47 3 5.74V17c0 .55.45 1 1 1h6c.55 0 1-.45 1-1v-2.26c1.81-1.27 3-3.36 3-5.74 0-3.86-3.14-7-7-7zm2.85 11.1l-.85.6V16h-4v-2.3l-.85-.6C7.8 12.16 7 10.63 7 9c0-2.76 2.24-5 5-5s5 2.24 5 5c0 1.63-.8 3.16-2.15 4.1z"/>
                </svg>
              </n-icon>
              <h3 class="tip-title">开始您的健康之旅</h3>
              <p class="tip-text">完善个人档案，获取专属营养方案</p>
              <n-button type="primary" size="large" @click="goToProfile">
                立即完善档案
              </n-button>
            </div>
          </n-card>
        </n-gi>
      </n-grid>
    </div>
  </div>
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

// Remove duplicated function definitions below this line if any


</script>

<style scoped>
/* Container */
.home-container {
  min-height: 100vh;
  background: #F5F7FA;
}

.main-wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
}

/* Hero Card */
.hero-card {
  padding: 32px;
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.greeting-section {
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
}

.greeting-text {
  font-size: 2.25rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 8px 0;
  line-height: 1.2;
}

.username {
  color: #10b981;
}

.greeting-subtitle {
  font-size: 0.95rem;
  color: #6b7280;
  margin: 0;
}

/* Calorie Ring */
.calorie-ring-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 12px 0;
}

.ring-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

.ring-number {
  font-size: 2rem;
  font-weight: 800;
  line-height: 1;
  font-family: 'Inter', -apple-system, sans-serif;
}

.ring-label {
  font-size: 0.75rem;
  color: #9ca3af;
  margin-top: 4px;
}

.ring-footer {
  margin-top: 12px;
  font-size: 0.8rem;
  color: #9ca3af;
  font-weight: 500;
}

.empty-ring {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.calorie-footer {
  text-align: center;
  margin-top: -16px; /* Pull up closer to ring */
  position: relative;
  z-index: 10;
}

.calorie-detail {
  font-size: 0.875rem;
  color: #6b7280;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.8);
  padding: 4px 12px;
  border-radius: 20px;
  display: inline-block;
}

.empty-ring {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

/* Recipe Magazine Card */
.recipe-magazine-card {
  position: relative;
  min-height: 320px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.recipe-magazine-card:hover {
  transform: translateY(-4px);
}

.recipe-magazine-card.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.recipe-bg-image {
  position: absolute;
  inset: 0;
  background-image: url('https://images.unsplash.com/photo-1498837167922-ddd27525d352?w=1200&q=80');
  background-size: cover;
  background-position: center;
  z-index: 0;
}

.recipe-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(0,0,0,0.6) 0%, rgba(0,0,0,0.3) 100%);
  z-index: 1;
}

.recipe-content {
  position: relative;
  z-index: 2;
  height: 100%;
  min-height: 280px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 32px 32px 32px 40px;
  color: white;
}

.recipe-tag {
  margin-bottom: 16px;
}

.recipe-title {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 12px 0;
  color: white;
  text-shadow: 0 2px 8px rgba(0,0,0,0.3);
}

.recipe-description {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 24px 0;
  line-height: 1.6;
}

.recipe-action {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

/* Action Cards */
.action-card {
  padding: 24px;
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
}

.action-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.action-icon-wrapper {
  width: 72px;
  height: 72px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.bg-red-100 {
  background-color: #fee2e2;
}

.bg-blue-100 {
  background-color: #dbeafe;
}

.bg-yellow-100 {
  background-color: #fef3c7;
}

.bg-green-100 {
  background-color: #d1fae5;
}

.action-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.action-desc {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0;
}

/* Tip Card */
.tip-card {
  padding: 32px;
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.tip-content {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.tip-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #78350f;
  margin: 0;
}

.tip-text {
  font-size: 1rem;
  color: #92400e;
  margin: 0;
  max-width: 500px;
}

/* Responsive */
@media (max-width: 1024px) {
  .calorie-ring-container :deep(.n-progress) {
    width: 140px !important;
    height: 140px !important;
  }
  
  .calorie-number {
    font-size: 2.25rem;
  }
}

@media (max-width: 768px) {
  .main-wrapper {
    padding: 16px;
  }

  .hero-card {
    padding: 24px;
  }

  /* Hero卡片移动端居中 */
  .greeting-section {
    text-align: center;
  }

  .greeting-text {
    font-size: 1.75rem;
  }

  .greeting-subtitle {
    font-size: 0.875rem;
  }

  /* 卡路里圆环居中 */
  .calorie-ring-container {
    justify-content: center;
    padding: 24px 0;
  }

  .calorie-ring-container :deep(.n-progress) {
    width: 160px !important;
    height: 160px !important;
  }

  .calorie-number {
    font-size: 2.25rem;
  }

  .calorie-detail {
    font-size: 0.7rem;
  }

  /* Banner高度增加 */
  .recipe-magazine-card {
    min-height: 400px;
  }

  .recipe-content {
    min-height: 360px;
    padding: 28px;
  }

  .recipe-title {
    font-size: 1.75rem;
  }

  .recipe-description {
    font-size: 1rem;
  }

  /* Action卡片调整 */
  .action-card {
    padding: 20px;
  }

  .action-icon-wrapper {
    width: 64px;
    height: 64px;
  }

  .action-title {
    font-size: 1rem;
  }

  .action-desc {
    font-size: 0.8125rem;
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
