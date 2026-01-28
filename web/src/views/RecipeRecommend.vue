<template>
  <div class="page-container">
    <TopNavigation />

    <div class="main-wrapper">
      <!-- 页面头部: 标题与操作 -->
      <div class="page-header">
        <div class="header-left">
          <n-button quaternary circle size="large" @click="router.push('/home')">
            <template #icon>
              <n-icon><ArrowBack /></n-icon>
            </template>
          </n-button>
          <div>
            <h1 class="page-title">智能食谱推荐</h1>
            <p class="page-subtitle">基于您的营养需求，为您定制专属每日食谱</p>
          </div>
        </div>
        
        <div class="header-right">
          <n-button 
            v-if="currentSelectedPlan" 
            type="error" 
            secondary
            class="action-btn"
            @click="reselectPlan"
          >
            <template #icon>
              <n-icon><Refresh /></n-icon>
            </template>
            重新选择
          </n-button>

          <n-button 
            v-else
            secondary 
            type="primary" 
            :loading="loading" 
            @click="fetchRecommendations(true)"
          >
            <template #icon>
              <n-icon><Refresh /></n-icon>
            </template>
            换一批推荐
          </n-button>
        </div>
      </div>

      <!-- Loading 状态 -->
      <div v-if="loading" class="loading-state">
        <n-spin size="large">
          <template #description>正在为您生成个性化营养方案...</template>
        </n-spin>
      </div>

      <!-- Error 状态 -->
      <div v-else-if="error" class="error-state">
        <n-result status="500" title="获取推荐失败" :description="error">
          <template #footer>
            <n-button type="primary" @click="fetchRecommendations(true)">重试</n-button>
          </template>
        </n-result>
      </div>

      <!-- 内容区域 -->
      <div v-else class="content-area">
        
        <!-- 场景1: 已选择方案 (Selected Plan View) -->
        <div v-if="currentSelectedPlan" class="selected-view">
          <n-alert title="今日食谱已锁定" type="success" class="mb-6">
            您已选择今日的饮食计划。按照此计划执行，助您达成健康目标！
            <template #action>
              <n-button size="small" type="error" ghost @click="reselectPlan">
                重新选择
              </n-button>
            </template>
          </n-alert>

          <n-grid :x-gap="24" :y-gap="24" cols="1 l:3" responsive="screen">
            <!-- 左侧：营养概览 -->
            <n-gi span="1">
              <n-card title="今日营养目标" :bordered="false" class="nutrition-card">
                <div class="nutrition-summary">
                  <div class="macro-item">
                    <n-progress type="circle" :percentage="getPercentage(currentSelectedPlan.total_energy, currentSelectedPlan.target_energy)" color="#10b981">
                      <div class="progress-text">
                        <span class="value">{{ Math.floor(currentSelectedPlan.total_energy) }}</span>
                        <span class="label">kcal</span>
                      </div>
                    </n-progress>
                    <div class="macro-label">热量 ({{ Math.floor(currentSelectedPlan.target_energy) }})</div>
                  </div>
                  
                  <div class="macros-bars">
                    <div class="macro-bar">
                      <div class="mb-header">
                        <div class="mb-label-group">
                          <n-icon color="#ef4444"><FitnessOutline /></n-icon>
                          <span class="mb-label protein">蛋白质</span>
                        </div>
                        <span class="mb-val">{{ Math.floor(currentSelectedPlan.total_protein) }} / {{ Math.floor(currentSelectedPlan.target_protein) }}g</span>
                      </div>
                      <n-progress 
                        type="line" 
                        :percentage="getPercentage(currentSelectedPlan.total_protein, currentSelectedPlan.target_protein)" 
                        color="#ef4444" 
                        :height="12"
                        :show-indicator="false" 
                        class="thick-progress"
                      />
                    </div>
                    
                    <div class="macro-bar">
                      <div class="mb-header">
                        <div class="mb-label-group">
                          <n-icon color="#f59e0b"><LeafOutline /></n-icon>
                          <span class="mb-label carb">碳水</span>
                        </div>
                        <span class="mb-val">{{ Math.floor(currentSelectedPlan.total_carbohydrate) }} / {{ Math.floor(currentSelectedPlan.target_carbohydrate) }}g</span>
                      </div>
                      <n-progress 
                        type="line" 
                        :percentage="getPercentage(currentSelectedPlan.total_carbohydrate, currentSelectedPlan.target_carbohydrate)" 
                        color="#f59e0b" 
                        :height="12"
                        :show-indicator="false" 
                        class="thick-progress"
                      />
                    </div>
                    
                    <div class="macro-bar">
                      <div class="mb-header">
                        <div class="mb-label-group">
                          <n-icon color="#8b5cf6"><WaterOutline /></n-icon>
                          <span class="mb-label fat">脂肪</span>
                        </div>
                        <span class="mb-val">{{ Math.floor(currentSelectedPlan.total_fat) }} / {{ Math.floor(currentSelectedPlan.target_fat) }}g</span>
                      </div>
                      <n-progress 
                        type="line" 
                        :percentage="getPercentage(currentSelectedPlan.total_fat, currentSelectedPlan.target_fat)" 
                        color="#8b5cf6" 
                        :height="12"
                        :show-indicator="false" 
                        class="thick-progress"
                      />
                    </div>
                  </div>
                </div>
              </n-card>
            </n-gi>

            <!-- 右侧：餐单详情 -->
            <n-gi span="2">
              <n-card title="每日菜单" :bordered="false">
                <div class="meal-timeline">
                  <MealItem title="早餐" icon="🌅" :recipe="currentSelectedPlan.breakfast" />
                  <MealItem title="午餐" icon="☀️" :recipe="currentSelectedPlan.lunch" />
                  <MealItem title="晚餐" icon="🌙" :recipe="currentSelectedPlan.dinner" />
                  <MealItem v-if="currentSelectedPlan.snack" title="加餐" icon="🍎" :recipe="currentSelectedPlan.snack" />
                </div>
              </n-card>
            </n-gi>
          </n-grid>
        </div>

        <!-- 场景2: 推荐列表 (Recommendations List) -->
        <div v-else class="recommendations-view">
          <div class="section-title">
            <h3>为您生成的 {{ plans.length }} 套方案</h3>
            <span class="subtitle">点击查看详情并选择</span>
          </div>

          <n-grid :x-gap="24" :y-gap="24" cols="1 m:2 l:3" responsive="screen">
            <n-gi v-for="(plan, index) in plans" :key="index">
              <n-card 
                class="plan-card" 
                :class="{ 'active': selectedPlanIndex === index }"
                @click="selectPlan(index)"
                content-style="padding: 0;"
              >
                <div class="card-content">
                  <!-- Header Section -->
                  <div class="card-header-section">
                    <div class="plan-info">
                      <h3 class="plan-name">方案 {{ index + 1 }}</h3>
                      <div class="match-badge">
                        <n-icon size="14"><CheckmarkCircle /></n-icon>
                        {{ plan.match_score }}% 匹配
                      </div>
                    </div>
                    <div class="calories-display">
                      <span class="cal-val">{{ Math.floor(plan.total_energy) }}</span>
                      <span class="cal-unit">kcal</span>
                    </div>
                  </div>

                  <!-- Meal List (Compact) -->
                  <div class="compact-meal-list">
                    <div class="compact-meal-item">
                      <div class="cmi-icon">🌅</div>
                      <div class="cmi-content">
                        <span class="cmi-name">{{ plan.breakfast?.name || '未安排' }}</span>
                        <span class="cmi-cal">{{ Math.floor(plan.breakfast?.energy || 0) }} kcal</span>
                      </div>
                    </div>
                    <div class="compact-meal-item">
                      <div class="cmi-icon">☀️</div>
                      <div class="cmi-content">
                        <span class="cmi-name">{{ plan.lunch?.name || '未安排' }}</span>
                        <span class="cmi-cal">{{ Math.floor(plan.lunch?.energy || 0) }} kcal</span>
                      </div>
                    </div>
                    <div class="compact-meal-item">
                      <div class="cmi-icon">🌙</div>
                      <div class="cmi-content">
                        <span class="cmi-name">{{ plan.dinner?.name || '未安排' }}</span>
                        <span class="cmi-cal">{{ Math.floor(plan.dinner?.energy || 0) }} kcal</span>
                      </div>
                    </div>
                  </div>

                  <!-- Macro Pills -->
                  <div class="macro-pills">
                    <div class="macro-pill protein">
                      <span class="pill-dot"></span>
                      <span>{{ Math.floor(plan.total_protein) }}g 蛋白</span>
                    </div>
                    <div class="macro-pill carb">
                      <span class="pill-dot"></span>
                      <span>{{ Math.floor(plan.total_carbohydrate) }}g 碳水</span>
                    </div>
                    <div class="macro-pill fat">
                      <span class="pill-dot"></span>
                      <span>{{ Math.floor(plan.total_fat) }}g 脂肪</span>
                    </div>
                  </div>

                  <!-- Action Area -->
                  <div class="card-action-area">
                    <n-button 
                      v-if="selectedPlanIndex === index" 
                      type="primary" 
                      block 
                      size="large"
                      class="confirm-btn"
                      @click.stop="confirmSelection"
                    >
                      确认选择
                    </n-button>
                    <div v-else class="select-hint">点击查看详情</div>
                  </div>
                </div>
              </n-card>
            </n-gi>
          </n-grid>

          <!-- 选中方案的详情弹窗/展开视图 (Optional, for now handled by selection state visually) -->
           <!-- 为了简化交互，点击卡片即选中高亮，再次点击确认或点击下方按钮确认 -->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { 
  NButton, NIcon, NSpin, NResult, NGrid, NGi, NCard, NTag, NProgress, NAlert, useMessage, NImage 
} from 'naive-ui';
import { 
  ArrowBack, Refresh, FitnessOutline, LeafOutline, WaterOutline, CheckmarkCircle 
} from '@vicons/ionicons5';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import MealItem from '@/components/MealItem.vue';
import { getRecipeRecommendations, selectRecipePlan, getSelectedRecipePlan } from '@/api/recipeApi';
import { getNutritionRequirements } from '@/api/user';
import { useAuthStore } from '@/store/auth';

const router = useRouter();
const message = useMessage();
const authStore = useAuthStore();

const loading = ref(true);
const error = ref(null);
const plans = ref([]);
const selectedPlanIndex = ref(null);
const currentSelectedPlan = ref(null);
const targetNutrition = ref({});

// 初始化
onMounted(async () => {
  await fetchRecommendations();
});

const fetchRecommendations = async (forceRefresh = false) => {
  loading.value = true;
  error.value = null;
  selectedPlanIndex.value = null;

  try {
    // 1. 如果不是强制刷新，先检查是否有已选方案
    if (!forceRefresh) {
      try {
        const selected = await getSelectedRecipePlan();
        if (selected && selected.id) {
          currentSelectedPlan.value = selected;
          loading.value = false;
          return;
        }
      } catch (e) {
        // 忽略 404
      }
    }

    // 2. 获取推荐
    // 并行获取营养目标和推荐列表
    const [nutritionResp, recipeResp] = await Promise.all([
      getNutritionRequirements(),
      getRecipeRecommendations(3)
    ]);

    targetNutrition.value = nutritionResp;
    plans.value = recipeResp.plans || [];

  } catch (err) {
    console.error(err);
    error.value = err.response?.data?.error || '获取推荐失败，请稍后重试';
    message.error(error.value);
  } finally {
    loading.value = false;
  }
};

const selectPlan = (index) => {
  selectedPlanIndex.value = index;
};

const confirmSelection = async () => {
  if (selectedPlanIndex.value === null) return;
  
  const plan = plans.value[selectedPlanIndex.value];
  try {
    loading.value = true;
    await selectRecipePlan(plan);
    currentSelectedPlan.value = plan;
    message.success('已成功锁定今日食谱！');
    // Scroll to top
    window.scrollTo({ top: 0, behavior: 'smooth' });
  } catch (err) {
    message.error(err.response?.data?.error || '保存失败');
    loading.value = false;
  } finally {
    loading.value = false;
  }
};

const reselectPlan = () => {
  currentSelectedPlan.value = null;
  // 此时 plans 还有数据，可以直接显示推荐列表，或者重新获取
  if (plans.value.length === 0) {
    fetchRecommendations(true);
  }
};

const getPercentage = (val, target) => {
  if (!target) return 0;
  return Math.min(100, Math.round((val / target) * 100));
};
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background-color: #F5F7FA;
}

.main-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

/* Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
  line-height: 1.2;
}

.page-subtitle {
  font-size: 14px;
  color: #6b7280;
  margin: 4px 0 0 0;
}

/* Loading & Error */
.loading-state, .error-state {
  display: flex;
  justify-content: center;
  padding: 60px 0;
}

/* Selected View */
.nutrition-card {
  height: 100%;
}

.nutrition-summary {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

.macro-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.progress-text {
  text-align: center;
  display: flex;
  flex-direction: column;
}

.progress-text .value {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
}

.progress-text .label {
  font-size: 12px;
  color: #6b7280;
}

.macro-label {
  margin-top: 8px;
  font-size: 14px;
  color: #4b5563;
  font-weight: 500;
}

.macros-bars {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.mb-header {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  margin-bottom: 4px;
}

.mb-label.protein { color: #ef4444; }
.mb-label.carb { color: #f59e0b; }
.mb-label.fat { color: #8b5cf6; }

.mb-val { color: #6b7280; }

.meal-timeline {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Recommendations View */
.section-title {
  margin-bottom: 24px;
}

.section-title h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.section-title .subtitle {
  font-size: 14px;
  color: #6b7280;
}

.mb-label-group {
  display: flex;
  align-items: center;
  gap: 4px;
}

.thick-progress :deep(.n-progress-graph-line-rail),
.thick-progress :deep(.n-progress-graph-line-fill) {
  border-radius: 9999px;
}

/* Plan Card Styles */
.plan-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid transparent;
  cursor: pointer;
  overflow: hidden;
  border-radius: 16px;
  background-color: white;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03); /* shadow-sm */
}

.plan-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05); /* shadow-lg */
}

.plan-card.active {
  border-color: #10b981; /* ring-primary-500 */
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04); /* shadow-xl */
}

.card-content {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.card-header-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.plan-name {
  font-size: 18px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 4px 0;
}

.match-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #059669;
  background-color: #d1fae5;
  padding: 2px 8px;
  border-radius: 999px;
  font-weight: 600;
}

.calories-display {
  text-align: right;
}

.cal-val {
  display: block;
  font-size: 24px;
  font-weight: 800;
  color: #1f2937;
  line-height: 1;
}

.cal-unit {
  font-size: 12px;
  color: #6b7280;
}

.compact-meal-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 4px 0;
}

.compact-meal-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cmi-icon {
  font-size: 18px;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f3f4f6;
  border-radius: 8px;
  flex-shrink: 0;
}

.cmi-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  overflow: hidden;
}

.cmi-name {
  font-size: 14px;
  color: #374151;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-right: 8px;
}

.cmi-cal {
  font-size: 12px;
  color: #9ca3af;
  white-space: nowrap;
  flex-shrink: 0;
}

.macro-pills {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.macro-pill {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 12px;
  padding: 6px 4px;
  border-radius: 999px;
  font-weight: 600;
  white-space: nowrap;
}

.macro-pill.protein { background-color: #fee2e2; color: #991b1b; }
.macro-pill.carb { background-color: #fef3c7; color: #92400e; }
.macro-pill.fat { background-color: #f3e8ff; color: #6b21a8; }

.pill-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: currentColor;
}

.card-action-area {
  margin-top: 8px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.select-hint {
  font-size: 14px;
  color: #9ca3af;
  font-weight: 500;
}

.confirm-btn {
  font-weight: 600;
  letter-spacing: 0.5px;
}

.mb-6 {
  margin-bottom: 24px;
}

.action-btn {
  font-weight: 500;
}
</style>
