<template>
  <div v-if="authStore.isAuthenticated" class="etm-container">
    <div class="etm-wrapper">
      <div class="etm-layout">
        <aside class="etm-sidebar">
          <div class="sidebar-profile">
            <div class="sidebar-avatar">{{ userAvatarText }}</div>
            <div class="sidebar-user">
              <div class="sidebar-name">{{ userDisplayName }}</div>
              <div class="sidebar-sub">NutriPlan 用户</div>
            </div>
          </div>

          <div class="sidebar-menu">
            <button class="sidebar-item" @click="goToProfile">
              <n-icon><PersonOutline /></n-icon>个人档案
            </button>
            <button class="sidebar-item" @click="goToFavorites">
              <n-icon><StarOutline /></n-icon>我的收藏
            </button>
            <button class="sidebar-item active">
              <n-icon><RestaurantOutline /></n-icon>饮食记录
            </button>
            <button class="sidebar-item" @click="goToWeekly">
              <n-icon><BarChartOutline /></n-icon>周报告
            </button>
            <button class="sidebar-item" @click="goToShopping">
              <n-icon><BasketOutline /></n-icon>购物清单
            </button>
          </div>

          <button class="sidebar-logout" @click="handleLogout">
            <n-icon><LogOutOutline /></n-icon>退出登录
          </button>
        </aside>

        <div class="etm-main">
          <div class="etm-toolbar">
            <div class="toggle-group">
              <button class="toggle-btn active">日</button>
              <button class="toggle-btn" @click="goToWeekly">周</button>
            </div>
            <div class="date-nav">
              <n-icon class="nav-icon"><ChevronBackOutline /></n-icon>
              <n-icon class="nav-icon"><CalendarOutline /></n-icon>
              <n-icon class="nav-icon"><ChevronForwardOutline /></n-icon>
              <span class="date-text">今天</span>
            </div>
            <div class="toolbar-actions">
              <n-icon class="action-icon"><EllipsisVerticalOutline /></n-icon>
            </div>
          </div>

          <div class="etm-content">
            <div class="etm-meals">
              <div class="meals-header">
                <h2 class="section-title">饮食记录</h2>
                <div class="meals-calories">
                  <div class="pie-icon-placeholder" :style="{ background: 'conic-gradient(#10b981 ' + caloriePercentage + '%, #e2e8f0 ' + caloriePercentage + '%)' }"></div>
                  <span>{{ caloriesConsumed }} 千卡</span>
                </div>
                <div class="meals-actions">
                  <span class="text-xs text-emerald-600 cursor-pointer flex items-center mr-3 hover:text-emerald-700 transition" @click="() => loadRecommendation(true)">
                    <n-icon class="mr-1" :class="{'animate-spin': isRecommending}"><RefreshOutline /></n-icon>重新推荐
                  </span>
                  <n-icon class="action-icon"><EllipsisVerticalOutline /></n-icon>
                </div>
              </div>

              <div class="meal-block" v-for="mealInfo in mealsList" :key="mealInfo.type">
                <div class="meal-header">
                  <div class="meal-title-group">
                    <h3>{{ mealInfo.name }}</h3>
                    <div class="meal-calories text-emerald-600 font-bold ml-2">
                      <div
                        class="meal-macro-pie"
                        :style="{ background: mealInfo.macroGradient }"
                        :title="`碳水 ${mealInfo.carbs}g / 蛋白质 ${mealInfo.protein}g / 脂肪 ${mealInfo.fat}g`"
                      ></div>
                      <span class="meal-macro-label">碳/蛋/脂</span>
                    </div>
                  </div>
                  <n-icon class="action-icon"><EllipsisVerticalOutline /></n-icon>
                </div>

                <div style="padding: 0.75rem;">
                  <div v-if="dailyRecommendation && dailyRecommendation[mealInfo.type] && dailyRecommendation[mealInfo.type].id" style="margin-bottom: 1rem;">
                    <div style="font-size: 12px; color: #059669; margin-bottom: 8px; font-weight: 500; display: flex; align-items: center;">
                      <n-icon style="margin-right: 4px;"><FlameOutline /></n-icon> 推荐食谱
                    </div>
                    <div class="meal-card" style="border: 1px solid #d1fae5; background: #ecfdf5; cursor: pointer; border-radius: 12px; display: flex; gap: 1rem; padding: 0.75rem; align-items: center;" @click="router.push('/recipes/' + dailyRecommendation[mealInfo.type].id)">
                      <div style="flex-shrink: 0; display: flex; align-items: center; justify-content: center;" @click.stop>
                        <n-checkbox v-model:checked="recommendationCompleted[mealInfo.type]" size="large" />
                      </div>
                      <div class="meal-img" style="background: #ffffff; border: 1px solid #e2e8f0; width: 48px; height: 48px; flex-shrink: 0; border-radius: 8px; overflow: hidden;">
                        <img :src="dailyRecommendation[mealInfo.type].image_url" v-if="dailyRecommendation[mealInfo.type].image_url" style="width: 100%; height: 100%; object-fit: cover;" />
                        <div class="img-placeholder" style="width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #94a3b8;" v-else>
                          <n-icon size="20"><RestaurantOutline /></n-icon>
                        </div>
                      </div>
                      <div class="meal-details" style="flex-grow: 1;">
                        <div style="display: flex; align-items: center; gap: 8px;">
                          <div class="meal-name" style="margin: 0; font-weight: bold; color: #2563eb; text-decoration: underline; text-decoration-color: #93c5fd; text-underline-offset: 4px;">{{ dailyRecommendation[mealInfo.type].name }}</div>
                          <span style="font-size: 10px; padding: 2px 6px; background: #d1fae5; color: #047857; border-radius: 4px; border: 1px solid #a7f3d0; white-space: nowrap;">推荐方案</span>
                        </div>
                        <div style="font-size: 12px; color: #94a3b8; margin-top: 4px;">约 {{ dailyRecommendation[mealInfo.type].energy || 0 }} 千卡</div>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="meal-items" v-if="mealInfo.items.length > 0">
                  <div class="meal-card" style="border-bottom: 1px solid #f1f5f9; cursor: pointer; padding-bottom: 1rem; margin-bottom: 0.5rem; display: flex; align-items: center; gap: 0.75rem;" v-for="item in mealInfo.items" :key="item.id">
                    <div style="flex-shrink: 0; display: flex; align-items: center; justify-content: center;" @click.stop>
                      <n-checkbox v-model:checked="item.completed" size="large" />
                    </div>
                    <div class="meal-img" style="flex-shrink: 0; width: 60px; height: 60px; border-radius: 8px; overflow: hidden;" :style="{ opacity: item.completed ? 0.5 : 1 }">
                      <img v-if="item.image" :src="item.image" :alt="item.foodName" style="width: 100%; height: 100%; object-fit: cover;" />
                      <div v-else class="img-placeholder" style="width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #94a3b8; background: #f1f5f9;">
                        <n-icon size="20"><RestaurantOutline /></n-icon>
                      </div>
                    </div>
                    <div class="meal-details" style="flex-grow: 1;">
                      <div style="display: flex; align-items: center; gap: 8px;">
                        <span class="meal-name" :style="{ margin: 0, color: item.completed ? '#94a3b8' : '#1e293b', fontWeight: 500, textDecoration: item.completed ? 'line-through' : 'none', transition: 'all 0.3s' }">{{ item.foodName }}</span>
                        <span style="font-size: 10px; padding: 2px 6px; background: #f1f5f9; color: #64748b; border-radius: 4px;">自带食物</span>
                      </div>
                      <div style="font-size: 12px; color: #94a3b8; margin-top: 4px;">每 100g 约 {{ item.calories || 0 }} 千卡</div>
                    </div>
                    <div class="meal-serving" style="display: flex; align-items: center; gap: 12px;">
                      <n-input-number v-model:value="item.intakeAmount" size="small" :step="10" :min="1" suffix="克" style="width: 100px;" />
                      <button @click.stop="handleDeleteMealItem(item.id)" style="background: none; border: none; color: #cbd5e1; cursor: pointer; display: flex; align-items: center; justify-content: center; padding: 4px;">
                        <n-icon size="22"><CloseOutline /></n-icon>
                      </button>
                    </div>
                  </div>
                </div>

                <div style="padding: 0.75rem; margin-top: 0.5rem;">
                  <button class="add-food-btn" style="width: 100%;" @click="goToIntake">+ 添加食物到 {{ mealInfo.name }}</button>
                </div>
              </div>
            </div>

            <div class="etm-nutrition">
              <div class="nutrition-header">
                <h2 class="section-title">营养概览</h2>
                <div class="nutrition-actions">
                  <button class="icon-btn active"><n-icon><PieChartOutline /></n-icon></button>
                  <button class="icon-btn"><n-icon><FlameOutline /></n-icon></button>
                </div>
              </div>

              <div class="nutrition-card">
                <div class="chart-container">
                  <div ref="macroChartRef" style="width: 280px; height: 280px;"></div>
                </div>

                <div class="stats-table">
                  <div class="stats-header">
                    <span></span>
                    <span class="col-center">总计</span>
                    <span class="col-right">目标 <n-icon class="edit-icon"><CreateOutline /></n-icon></span>
                  </div>
                  <div class="stats-row highlight">
                    <span>卡路里</span>
                    <span class="col-center">{{ caloriesConsumed }}</span>
                    <span class="col-right">{{ calorieTarget }}</span>
                  </div>

                  <div class="stats-item">
                    <div class="stats-row">
                      <span class="text-amber-500 font-medium">碳水化合物</span>
                      <span class="col-center">{{ totalCarbs }}克</span>
                      <span class="col-right">{{ targetCarbs }}克</span>
                    </div>
                    <div class="mini-progress-container">
                      <n-progress type="line" :percentage="Math.min(100, targetCarbs ? (totalCarbs / targetCarbs * 100) : 0)" color="#fbbf24" :show-indicator="false" :height="6" class="mini-progress" />
                    </div>
                  </div>

                  <div class="stats-item">
                    <div class="stats-row">
                      <span class="text-indigo-400 font-medium">脂肪</span>
                      <span class="col-center">{{ totalFat }}克</span>
                      <span class="col-right">{{ targetFat }}克</span>
                    </div>
                    <div class="mini-progress-container">
                      <n-progress type="line" :percentage="Math.min(100, targetFat ? (totalFat / targetFat * 100) : 0)" color="#818cf8" :show-indicator="false" :height="6" class="mini-progress" />
                    </div>
                  </div>

                  <div class="stats-item">
                    <div class="stats-row">
                      <span class="text-rose-500 font-medium">蛋白质</span>
                      <span class="col-center">{{ totalProtein }}克</span>
                      <span class="col-right">{{ targetProtein }}克</span>
                    </div>
                    <div class="mini-progress-container">
                      <n-progress type="line" :percentage="Math.min(100, targetProtein ? (totalProtein / targetProtein * 100) : 0)" color="#fb7185" :show-indicator="false" :height="6" class="mini-progress" />
                    </div>
                  </div>

                  <div class="stats-divider"></div>

                  <div class="stats-row disabled text-gray-400">
                    <span>纤维</span>
                    <span class="col-center">-</span>
                    <span class="col-right">-</span>
                  </div>
                  <div class="stats-row disabled text-gray-400">
                    <span>钠</span>
                    <span class="col-center">-</span>
                    <span class="col-right">-</span>
                  </div>
                  <div class="stats-row disabled text-gray-400">
                    <span>胆固醇</span>
                    <span class="col-center">-</span>
                    <span class="col-right">-</span>
                  </div>

                  <button class="detailed-btn" @click="goToProfile">详细营养信息</button>
                </div>
              </div>

              <div class="funny-character">
                <span class="avocado-emoji"></span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { onMounted, ref, computed } from "vue";
  import { useRouter } from "vue-router";
  import { NIcon, useMessage, NInputNumber, NCheckbox } from "naive-ui";
  import {
    ChevronBackOutline,
    CalendarOutline,
    ChevronForwardOutline,
    EllipsisVerticalOutline,
    RefreshOutline,
    PieChartOutline,
    FlameOutline,
    CreateOutline,
    CloseOutline,
    RestaurantOutline,
    PersonOutline,
    StarOutline,
    BarChartOutline,
    BasketOutline,
    LogOutOutline
  } from "@vicons/ionicons5";
import { useAuthStore } from "@/store/auth";
import { getNutritionRequirements } from "@/api/user";
import { getTodayStatus } from "@/api/intakeApi";
import { getRecipeRecommendations } from "@/api/recipeApi";

const router = useRouter();
const message = useMessage();
const authStore = useAuthStore();
const nutritionData = ref(null);
const todayIntake = ref(null);
const userDisplayName = computed(() => authStore.user?.username || "用户");
const userAvatarText = computed(() => (userDisplayName.value || "用").slice(0, 1).toUpperCase());

  const dailyRecommendation = ref(null);
  const recommendationCompleted = ref({
    breakfast: false,
    lunch: false,
    dinner: false,
    snack: false
  });
  const isRecommending = ref(false);

  const CACHE_KEY = 'nutriplan_daily_recommendation';

  const loadRecommendation = async (force = false) => {
    isRecommending.value = true;
    try {
      const now = new Date();
      const today = `${now.getFullYear()}-${now.getMonth()+1}-${now.getDate()}`;

      // Check cache if not forcing a refresh
      if (!force) {
        const cached = localStorage.getItem(CACHE_KEY);
        if (cached) {
          try {
            const parsed = JSON.parse(cached);
            if (parsed.date === today && parsed.data) {
              dailyRecommendation.value = parsed.data;
              return; // Use cached data
            }
          } catch (e) {
            console.warn("解析推荐缓存失败:", e);
          }
        }
      }

      const res = await getRecipeRecommendations(1);
      if (res && res.plans && res.plans.length > 0) {
        dailyRecommendation.value = {
          breakfast: res.plans[0].breakfast,
          lunch: res.plans[0].lunch,
          dinner: res.plans[0].dinner,
          snack: res.plans[0].snack
        };

        // Update local storage cache
        localStorage.setItem(CACHE_KEY, JSON.stringify({
          date: today,
          data: dailyRecommendation.value
        }));

        if (force) {
          message.success("食谱推荐已更新");
        }
      }
    } catch (e) {
      console.error("加载食谱推荐失败", e);
      if (force) {
        message.error("加载推荐失败");
      }
    } finally {
      isRecommending.value = false;
    }
  };

  onMounted(async () => {
  if (authStore.isAuthenticated) {
    if (!authStore.user || !authStore.profile) {
      try {
        await authStore.loadProfile();
      } catch (error) {
        console.error("加载档案失败:", error);
      }
    }

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
      
      loadRecommendation();
    }
  }
});

// Routing overrides
  const goToProfile = () => router.push("/profile/view");
  const goToIntake = () => router.push("/intake");
  const goToWeekly = () => router.push("/weekly");
  const goToFavorites = () => router.push("/favorites");
  const goToShopping = () => router.push("/shopping");
  const handleLogout = () => {
    authStore.logout();
    router.push("/login");
  };

  const handleDeleteMealItem = (id) => {
    message.success("已删除该记录（演示使用）");
    // In real app, call delete API and refresh
  };// Calorie calcs
const calorieTarget = computed(() => {
  if (nutritionData.value && nutritionData.value.target_calorie) return Math.floor(nutritionData.value.target_calorie);
  return Math.floor(authStore.profile?.tdee || 0);
});

const caloriesConsumed = computed(() => Math.floor(todayIntake.value?.total_energy || 0));

const caloriePercentage = computed(() => {
  if (!calorieTarget.value) return 0;
  return Math.min(100, Math.max(0, (caloriesConsumed.value / calorieTarget.value) * 100));
});

// Macro calcs
const totalCarbs = computed(() => Math.floor(todayIntake.value?.total_carbohydrate || 0));
const totalFat = computed(() => Math.floor(todayIntake.value?.total_fat || 0));
const totalProtein = computed(() => Math.floor(todayIntake.value?.total_protein || 0));

const targetCarbs = computed(() => Math.floor(nutritionData.value?.target_carbohydrate || 0));
const targetFat = computed(() => Math.floor(nutritionData.value?.target_fat || 0));
const targetProtein = computed(() => Math.floor(nutritionData.value?.target_protein || 0));

const macroTotal = computed(() => totalCarbs.value + totalFat.value + totalProtein.value);

const carbPercentage = computed(() => macroTotal.value ? Math.round((totalCarbs.value / macroTotal.value) * 100) : 0);
const fatPercentage = computed(() => macroTotal.value ? Math.round((totalFat.value / macroTotal.value) * 100) : 0);
const proteinPercentage = computed(() => macroTotal.value ? Math.round((totalProtein.value / macroTotal.value) * 100) : 0);

const macroConicGradient = computed(() => {
  if (macroTotal.value === 0) return 'background: conic-gradient(#e2e8f0 100%)';
  const fatEnd = fatPercentage.value;
  const proteinEnd = fatEnd + proteinPercentage.value;
  return 'background: conic-gradient(#818cf8 0% ' + fatEnd + '%, #fb7185 ' + fatEnd + '% ' + proteinEnd + '%, #fbbf24 ' + proteinEnd + '% 100%)';
});

const buildMacroGradient = (fat, protein, carbs) => {
  const total = fat + protein + carbs;
  if (!total) return 'conic-gradient(#e2e8f0 100%)';
  const fatEnd = Math.round((fat / total) * 100);
  const proteinEnd = fatEnd + Math.round((protein / total) * 100);
  return `conic-gradient(#818cf8 0% ${fatEnd}%, #fb7185 ${fatEnd}% ${proteinEnd}%, #fbbf24 ${proteinEnd}% 100%)`;
};

// Meals List formatting
const generateMealBlock = (label, type, code) => {
    const records = todayIntake.value?.records || [];
    const typeRecords = records.filter(r => r.mealType === type || r.mealType === code || r.mealType === label.toLowerCase());

    const calories = typeRecords.reduce((sum, r) => sum + (r.calculatedEnergy || 0), 0);
    const fat = typeRecords.reduce((sum, r) => sum + (r.calculatedFat || 0), 0);
    const protein = typeRecords.reduce((sum, r) => sum + (r.calculatedProtein || 0), 0);
    const carbs = typeRecords.reduce((sum, r) => sum + (r.calculatedCarb || 0), 0);

    return {
      name: label,
      type: type,
      calories: Math.floor(calories),
      fat: Math.floor(fat),
      protein: Math.floor(protein),
      carbs: Math.floor(carbs),
      macroGradient: buildMacroGradient(fat, protein, carbs),
      items: typeRecords.map(r => ({
        id: r.ID,
        foodName: r.foodName,
        intakeAmount: r.intakeAmount,
        image: r.imageUrl || null,
        completed: false
      }))
    };
  };const mealsList = computed(() => {
  return [
    generateMealBlock('早餐', 'breakfast', 1),
    generateMealBlock('午餐', 'lunch', 2),
    generateMealBlock('晚餐', 'dinner', 3),
    generateMealBlock('加餐', 'snack', 4)
  ];
});

// ECharts logic
import * as echarts from 'echarts';
import { watch, onUnmounted, nextTick } from 'vue';
import { NProgress } from 'naive-ui'; // Import NProgress for the mini progress bars

const macroChartRef = ref(null);
let donutChart = null;

const initChart = () => {
  if (macroChartRef.value) {
    if (!donutChart) {
      donutChart = echarts.init(macroChartRef.value);
    }
    const option = {
      tooltip: {
        trigger: 'item'
      },
      title: {
        text: caloriesConsumed.value + '\n{small|kcal}',
        left: 'center',
        top: 'center',
        textStyle: {
          fontSize: 28,
          fontWeight: 'bold',
          color: '#1e293b',
          rich: {
            small: {
              fontSize: 14,
              color: '#64748b',
              padding: [4, 0, 0, 0]
            }
          }
        }
      },
      series: [
        {
          name: '营养素',
          type: 'pie',
          radius: ['60%', '80%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: false,
            position: 'center'
          },
          labelLine: {
            show: false
          },
          data: totalCarbs.value === 0 && totalProtein.value === 0 && totalFat.value === 0
            ? [{ value: 1, name: '暂无数据', itemStyle: { color: '#f1f5f9' } }]
            : [
              { value: totalCarbs.value, name: '碳水', itemStyle: { color: '#fbbf24' } },
              { value: totalProtein.value, name: '蛋白', itemStyle: { color: '#fb7185' } },
              { value: totalFat.value, name: '脂肪', itemStyle: { color: '#818cf8' } }
            ]
        }
      ]
    };
    donutChart.setOption(option);
  }
};

watch([totalCarbs, totalProtein, totalFat, caloriesConsumed], () => {
  nextTick(() => {
    initChart();
  });
}, { deep: true });

onMounted(() => {
  nextTick(() => {
    initChart();
  });
  
  // Resize chart on window resize
  window.addEventListener('resize', () => {
    if (donutChart) {
      donutChart.resize();
    }
  });
});

onUnmounted(() => {
  if (donutChart) {
    window.removeEventListener('resize', () => donutChart.resize());
    donutChart.dispose();
  }
});
</script>

<style scoped>
/* Modern Theme Base */
.etm-container {
  min-height: 100vh;
  background-color: #f8fafc; /* slate-50 */
  color: #1e293b; /* slate-800 */
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.etm-wrapper {
  max-width: 1400px;
  margin: 0 auto;
  padding: 1rem 2rem;
}

.etm-layout {
  display: grid;
  grid-template-columns: 250px minmax(0, 1fr);
  gap: 1.5rem;
  align-items: start;
}

.etm-main {
  min-width: 0;
}

.etm-sidebar {
  position: sticky;
  top: 1rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 1rem;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.sidebar-profile {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f1f5f9;
}

.sidebar-avatar {
  width: 42px;
  height: 42px;
  border-radius: 999px;
  background: #10b981;
  color: #ffffff;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sidebar-user {
  min-width: 0;
}

.sidebar-name {
  font-weight: 700;
  color: #1e293b;
  line-height: 1.1;
}

.sidebar-sub {
  color: #94a3b8;
  font-size: 0.78rem;
  margin-top: 0.2rem;
}

.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
}

.sidebar-item,
.sidebar-logout {
  border: none;
  background: transparent;
  color: #334155;
  border-radius: 0.65rem;
  padding: 0.55rem 0.7rem;
  width: 100%;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s;
  font-size: 0.92rem;
}

.sidebar-item:hover,
.sidebar-logout:hover {
  background: #f8fafc;
}

.sidebar-item.active {
  background: #ecfdf5;
  color: #047857;
  font-weight: 600;
}

.sidebar-logout {
  margin-top: auto;
  color: #b45309;
}

/* Toolbar */
.etm-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 0 1rem;
  border-bottom: 1px solid #e2e8f0; /* slate-200 */
}

.toggle-group {
  display: flex;
  background: #f1f5f9; /* slate-100 */
  border-radius: 6px;
  overflow: hidden;
}

.toggle-btn {
  background: none;
  border: none;
  color: #64748b; /* slate-500 */
  padding: 0.4rem 1rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
}

.toggle-btn.active {
  background: #10b981; /* emerald-500 */
  color: #fff;
}

.date-nav {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1e293b;
}

.nav-icon {
  cursor: pointer;
  color: #64748b;
  transition: color 0.2s;
}

.nav-icon:hover {
  color: #10b981;
}

.date-text {
  min-width: 80px;
  text-align: center;
}

.toolbar-actions .action-icon {
  font-size: 1.2rem;
  color: #64748b;
  cursor: pointer;
}

/* Content Layout */
.etm-content {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 2rem;
  padding: 1.5rem 0;
}

/* Left: Meals */
.meals-header {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.section-title {
  font-size: 1.4rem;
  font-weight: 600;
  color: #1e293b;
  margin-right: 1.5rem;
  margin-bottom: 0;
}

.meals-calories {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #64748b;
  font-size: 0.9rem;
}

.pie-icon-placeholder {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: conic-gradient(#10b981 75%, #e2e8f0 75%);
}

.pie-icon-placeholder.small {
  width: 12px;
  height: 12px;
}

.meals-actions {
  margin-left: auto;
  display: flex;
  gap: 1rem;
}

.action-icon {
  color: #64748b;
  font-size: 1.1rem;
  cursor: pointer;
}

.action-icon:hover {
  color: #1e293b;
}

.meal-block {
  background: #ffffff;
  border-radius: 1rem; /* rounded-2xl */
  margin-bottom: 1.5rem;
  overflow: hidden;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05); /* shadow-sm */
  border: 1px solid #f1f5f9;
}

.meal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #f1f5f9;
}

.meal-title-group {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.meal-title-group h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #1e293b;
}

.meal-calories {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.85rem;
  color: #64748b;
}

.meal-macro-pie {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.meal-macro-label {
  color: #64748b;
  font-weight: 500;
}

.meal-items {
  padding: 0.5rem;
}

.meal-card {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  gap: 1rem;
  border-radius: 0.75rem;
  transition: background 0.2s;
}

.meal-card:hover {
  background: #f8fafc;
}

.meal-checkbox {
  width: 18px;
  height: 18px;
  border: 2px solid #cbd5e1;
  border-radius: 4px;
}

.meal-img {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  overflow: hidden;
  background: #f1f5f9;
  flex-shrink: 0;
}

.meal-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.img-placeholder {
  width: 100%;
  height: 100%;
  background: #e2e8f0;
}

.meal-details {
  flex-grow: 1;
}

.meal-name {
  color: #1e293b;
  font-weight: 500;
  margin-bottom: 0.4rem;
  cursor: pointer;
  transition: color 0.2s;
}

.meal-name:hover {
  color: #10b981;
}

.meal-serving {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.serving-value {
  background: #f1f5f9;
  color: #1e293b;
  padding: 0.2rem 0.6rem;
  border-radius: 4px;
  font-size: 0.85rem;
  border: 1px solid #e2e8f0;
}

.serving-unit {
  font-size: 0.85rem;
  color: #64748b;
  display: flex;
  align-items: center;
  gap: 0.2rem;
  cursor: pointer;
}

.empty-meal {
  padding: 1rem;
}

.add-food-btn {
  background: none;
  border: 1px dashed #a7f3d0; /* emerald-200 */
  color: #059669; /* emerald-600 */
  width: 100%;
  padding: 0.75rem;
  border-radius: 8px;
  cursor: pointer;
  transition: colors 0.2s;
}

.add-food-btn:hover {
  background: #ecfdf5; /* emerald-50 */
}

/* Right: Nutrition */
.nutrition-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.nutrition-actions {
  display: flex;
  background: #f1f5f9;
  border-radius: 6px;
  padding: 2px;
}

.icon-btn {
  background: none;
  border: none;
  color: #64748b;
  padding: 0.4rem 0.6rem;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-btn.active {
  background: #ffffff;
  color: #10b981;
  border-radius: 6px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.nutrition-card {
  background: #ffffff;
  border-radius: 1rem; /* rounded-2xl */
  padding: 1.5rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05); /* shadow-sm */
  border: 1px solid #f1f5f9;
}

/* Chart */
.chart-container {
  display: flex;
  justify-content: center;
  margin: 2rem 0;
  position: relative;
}

/* Stats Table */
.stats-table {
  font-size: 0.9rem;
}

.stats-header {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  color: #64748b;
  margin-bottom: 0.5rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.col-center { text-align: right; margin-right: 1.5rem; }
.col-right { text-align: right; }

.stats-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  padding: 0.4rem 0;
  align-items: center;
  color: #1e293b;
}

.stats-row.highlight {
  font-weight: 600;
  font-size: 0.95rem;
  margin-bottom: 0.5rem;
}

.stats-row.disabled {
  color: #94a3b8;
}

.stats-item {
  margin-bottom: 0.5rem;
}

.stats-item .stats-row {
  padding-bottom: 0.1rem;
}

.mini-progress-container {
  padding-right: 1.5rem; /* match col-right */
  padding-left: 0.2rem;
  margin-bottom: 0.2rem;
}

/* Colors for specific text since Tailwind might not build correctly if dynamically injected */
.text-amber-500 { color: #f59e0b; }
.text-indigo-400 { color: #818cf8; }
.text-rose-500 { color: #f43f5e; }
.text-gray-400 { color: #9ca3af; }
.font-medium { font-weight: 500; }

.stats-divider {
  height: 1px;
  background: #e2e8f0;
  margin: 0.5rem 0;
}

.edit-icon {
  font-size: 0.9rem;
  vertical-align: middle;
  cursor: pointer;
  color: #64748b;
}

.edit-icon:hover {
  color: #10b981;
}

.detailed-btn {
  width: 100%;
  margin-top: 1.5rem;
  padding: 0.6rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #64748b;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
}

.detailed-btn:hover {
  background: #10b981;
  color: #fff;
  border-color: #10b981;
}

.funny-character {
  text-align: right;
  margin-top: 2rem;
  opacity: 0.8;
}
.avocado-emoji {
  font-size: 4rem;
  display: inline-block;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-10px) rotate(5deg); }
  100% { transform: translateY(0px) rotate(0deg); }
}

@media (max-width: 900px) {
  .etm-layout { grid-template-columns: 1fr; }
  .etm-sidebar { position: static; }
  .etm-content { grid-template-columns: 1fr; }
  .chart-container { margin: 1rem 0; }
}
</style>
