<template>
  <div v-if="authStore.isAuthenticated" class="etm-container">
    <div class="etm-main">
          <div class="etm-toolbar">
            <div class="toolbar-left">
              <div class="toggle-group">
                <button class="toggle-btn active">日</button>
                <button class="toggle-btn" @click="goToWeekly">周</button>
              </div>
              <div class="date-nav">
                <n-icon class="nav-icon" @click="goPrevDate"><ChevronBackOutline /></n-icon>
                <n-popover trigger="click" placement="bottom" v-model:show="datePickerVisible">
                  <template #trigger>
                    <n-icon class="nav-icon"><CalendarOutline /></n-icon>
                  </template>
                  <n-date-picker
                    v-model:value="selectedDateTimestamp"
                    type="date"
                    panel
                    :is-date-disabled="disableFutureDate"
                    :clearable="false"
                    @update:value="handleDateChange"
                  />
                </n-popover>
                <n-icon class="nav-icon" @click="goNextDate"><ChevronForwardOutline /></n-icon>
                <span class="date-text">{{ selectedDateLabel }}</span>
              </div>
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
                  <button
                    class="recommend-btn"
                    :class="{ pressed: recommendPressed, loading: isRecommending }"
                    :disabled="isRecommending"
                    @mousedown="recommendPressed = true"
                    @mouseup="recommendPressed = false"
                    @mouseleave="recommendPressed = false"
                    @click="handleRefreshRecommendation"
                  >
                    <n-icon class="recommend-icon" :class="{ 'animate-spin': isRecommending }"><RefreshOutline /></n-icon>
                    <span>{{ isRecommending ? '推荐中...' : '重新推荐' }}</span>
                  </button>
                  <n-dropdown trigger="click" :options="moreActionOptions" @select="handleMoreActionSelect">
                    <n-icon class="action-icon action-trigger"><EllipsisVerticalOutline /></n-icon>
                  </n-dropdown>
                </div>
              </div>
              <div class="daily-calorie-progress">
                <div class="daily-calorie-track">
                  <div class="daily-calorie-fill" :style="{ width: checkedCaloriePercentage + '%' }"></div>
                </div>
                <div class="daily-calorie-text">今日已完成 {{ checkedCaloriePercentage.toFixed(0) }}%</div>
              </div>

              <div class="meal-block" v-for="mealInfo in mealsList" :key="mealInfo.type">
                <div class="meal-header">
                  <div class="meal-title-group">
                    <h3>{{ mealInfo.name }}</h3>
                    <div class="meal-calories text-emerald-600 font-bold ml-2">
                      <div
                        class="meal-macro-pie"
                        :style="{ background: mealInfo.macroGradient }"
                        :title="`已勾选占比：碳水 ${mealInfo.checkedCarbs}g / 蛋白质 ${mealInfo.checkedProtein}g / 脂肪 ${mealInfo.checkedFat}g`"
                      ></div>
                      <span class="meal-checked-calories">{{ mealInfo.checkedCalories }} 千卡</span>
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
                        <n-checkbox
                          :checked="recommendationCompleted[mealInfo.type]"
                          @update:checked="(checked) => handleRecommendationChecked(mealInfo.type, checked)"
                          size="large"
                        />
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
                      <n-checkbox
                        :checked="item.completed"
                        @update:checked="(checked) => handleMealItemChecked(item.id, checked)"
                        size="large"
                      />
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
                  <div class="stats-header compact">
                    <span>营养项</span>
                    <span class="col-right">数值</span>
                  </div>
                  <div class="stats-row compact-row highlight">
                    <span>卡路里</span>
                    <span class="col-right">{{ overviewCalories }} 千卡</span>
                  </div>
                  <div class="stats-row compact-row">
                    <span class="text-amber-500 font-medium">碳水化合物</span>
                    <span class="col-right">{{ overviewCarbs }}克</span>
                  </div>
                  <div class="stats-row compact-row">
                    <span class="text-indigo-400 font-medium">脂肪</span>
                    <span class="col-right">{{ overviewFat }}克</span>
                  </div>
                  <div class="stats-row compact-row">
                    <span class="text-rose-500 font-medium">蛋白质</span>
                    <span class="col-right">{{ overviewProtein }}克</span>
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
</template>

<script setup>
  import { onMounted, ref, computed } from "vue";
  import { useRouter } from "vue-router";
  import { NIcon, NDropdown, NDatePicker, NPopover, useMessage, NInputNumber, NCheckbox } from "naive-ui";
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
    RestaurantOutline
  } from "@vicons/ionicons5";
import { useAuthStore } from "@/store/auth";
import { getNutritionRequirements } from "@/api/user";
import { getTodayStatus } from "@/api/intakeApi";
import { getRecipeRecommendations } from "@/api/recipeApi";
import { formatDateKey, updateCompletedRecipe } from "@/utils/completedRecipes";

const router = useRouter();
const message = useMessage();
const authStore = useAuthStore();
const nutritionData = ref(null);
const todayIntake = ref(null);
const mealCompletionMap = ref({});
const selectedDate = ref(new Date());
const selectedDateTimestamp = ref(new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate()).getTime());
const datePickerVisible = ref(false);

  const dailyRecommendation = ref(null);
  const recommendationCompleted = ref({
    breakfast: false,
    lunch: false,
    dinner: false,
    snack: false
  });
  const isRecommending = ref(false);
  const recommendPressed = ref(false);

  const moreActionOptions = [
    { label: '立即重新推荐', key: 'refresh' },
    { label: '清空推荐缓存', key: 'clear-cache' },
    { label: '重置今日勾选', key: 'reset-checklist' },
    { label: '前往周计划', key: 'goto-week' }
  ];

  const CHECKLIST_KEY_PREFIX = 'nutriplan_daily_checklist';
  const CACHE_KEY = 'nutriplan_daily_recommendation';

  const getDateTag = () => {
    return formatDateKey(selectedDate.value);
  };

  const disableFutureDate = (timestamp) => {
    const target = new Date(timestamp);
    const today = new Date();
    const endOfToday = new Date(today.getFullYear(), today.getMonth(), today.getDate(), 23, 59, 59, 999);
    return target.getTime() > endOfToday.getTime();
  };

  const isSelectedToday = computed(() => formatDateKey(selectedDate.value) === formatDateKey(new Date()));

  const selectedDateLabel = computed(() => {
    if (isSelectedToday.value) return '今天';
    return selectedDate.value.toLocaleDateString('zh-CN', { month: 'numeric', day: 'numeric', weekday: 'short' });
  });

  const setSelectedDate = (date) => {
    selectedDate.value = new Date(date.getFullYear(), date.getMonth(), date.getDate());
    selectedDateTimestamp.value = selectedDate.value.getTime();
  };

  const goPrevDate = async () => {
    setSelectedDate(new Date(selectedDate.value.getFullYear(), selectedDate.value.getMonth(), selectedDate.value.getDate() - 1));
    await loadSelectedDayData();
  };

  const goNextDate = async () => {
    const next = new Date(selectedDate.value.getFullYear(), selectedDate.value.getMonth(), selectedDate.value.getDate() + 1);
    if (disableFutureDate(next.getTime())) return;
    setSelectedDate(next);
    await loadSelectedDayData();
  };

  const handleDateChange = async (value) => {
    if (!value) return;
    setSelectedDate(new Date(value));
    datePickerVisible.value = false;
    await loadSelectedDayData();
  };

  const getChecklistStorageKey = () => {
    const userId = authStore.user?.id || 'guest';
    return `${CHECKLIST_KEY_PREFIX}_${userId}_${getDateTag()}`;
  };

  const loadChecklistState = () => {
    try {
      const raw = localStorage.getItem(getChecklistStorageKey());
      if (!raw) {
        mealCompletionMap.value = {};
        recommendationCompleted.value = {
          breakfast: false,
          lunch: false,
          dinner: false,
          snack: false
        };
        return;
      }
      const parsed = JSON.parse(raw);
      mealCompletionMap.value = parsed.mealCompletionMap || {};
      recommendationCompleted.value = {
        breakfast: !!parsed.recommendationCompleted?.breakfast,
        lunch: !!parsed.recommendationCompleted?.lunch,
        dinner: !!parsed.recommendationCompleted?.dinner,
        snack: !!parsed.recommendationCompleted?.snack
      };
    } catch (error) {
      console.warn('加载勾选状态失败', error);
    }
  };

  const saveChecklistState = () => {
    localStorage.setItem(
      getChecklistStorageKey(),
      JSON.stringify({
        mealCompletionMap: mealCompletionMap.value,
        recommendationCompleted: recommendationCompleted.value
      })
    );
  };

  const handleMealItemChecked = (itemId, checked) => {
    mealCompletionMap.value = {
      ...mealCompletionMap.value,
      [itemId]: !!checked
    };
    saveChecklistState();
  };

  const handleRecommendationChecked = (mealType, checked) => {
    recommendationCompleted.value = {
      ...recommendationCompleted.value,
      [mealType]: !!checked
    };

    const recipe = dailyRecommendation.value?.[mealType];
    if (recipe) {
      updateCompletedRecipe({
        userId: authStore.user?.id,
        dateKey: formatDateKey(selectedDate.value),
        mealType,
        recipe,
        checked: !!checked
      });
    }

    saveChecklistState();
  };

  const handleRefreshRecommendation = async () => {
    if (!isSelectedToday.value) {
      message.info('仅支持当天重新推荐');
      return;
    }
    recommendPressed.value = false;
    await loadRecommendation(true);
  };

  const resetTodayChecklist = () => {
    const recommendationTypes = ['breakfast', 'lunch', 'dinner', 'snack'];
    recommendationTypes.forEach((mealType) => {
      const recipe = dailyRecommendation.value?.[mealType];
      if (recipe) {
        updateCompletedRecipe({
          userId: authStore.user?.id,
          dateKey: formatDateKey(selectedDate.value),
          mealType,
          recipe,
          checked: false
        });
      }
    });

    mealCompletionMap.value = {};
    recommendationCompleted.value = {
      breakfast: false,
      lunch: false,
      dinner: false,
      snack: false
    };
    saveChecklistState();
  };

  const handleMoreActionSelect = async (key) => {
    if (key === 'refresh') {
      await handleRefreshRecommendation();
      return;
    }

    if (key === 'clear-cache') {
      localStorage.removeItem(CACHE_KEY);
      dailyRecommendation.value = null;
      message.success('已清空推荐缓存');
      return;
    }

    if (key === 'reset-checklist') {
      resetTodayChecklist();
      message.success('已重置今日勾选状态');
      return;
    }

    if (key === 'goto-week') {
      goToWeekly();
    }
  };

  const loadRecommendation = async (force = false) => {
    isRecommending.value = true;
    try {
      const selectedDateKey = formatDateKey(selectedDate.value);
      const today = formatDateKey(new Date());

      if (selectedDateKey !== today) {
        dailyRecommendation.value = null;
        return;
      }

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

    await loadSelectedDayData();

    try {
      if (authStore.hasProfile) {
        nutritionData.value = await getNutritionRequirements();
      } else {
        nutritionData.value = null;
      }
    } catch (error) {
      console.log("获取营养目标失败:", error);
      nutritionData.value = null;
    }

  }
});

  const loadSelectedDayData = async () => {
    try {
      const statusRes = await getTodayStatus(formatDateKey(selectedDate.value));
      todayIntake.value = {
        ...statusRes,
        total_energy: statusRes?.total_energy ?? statusRes?.totalEnergy ?? 0,
        total_protein: statusRes?.total_protein ?? statusRes?.totalProtein ?? 0,
        total_carbohydrate: statusRes?.total_carbohydrate ?? statusRes?.totalCarbohydrate ?? 0,
        total_fat: statusRes?.total_fat ?? statusRes?.totalFat ?? 0,
        records: statusRes?.records || []
      };
      loadChecklistState();
      await loadRecommendation();
    } catch (error) {
      console.log("获取日期营养状态失败:", error);
      todayIntake.value = {
        total_energy: 0,
        total_protein: 0,
        total_carbohydrate: 0,
        total_fat: 0,
        records: []
      };
      dailyRecommendation.value = null;
      loadChecklistState();
    }
  };

// Routing overrides
  const goToProfile = () => router.push("/profile/view");
  const goToIntake = () => router.push("/intake");
  const goToWeekly = () => router.push("/weekly-planner");

  const handleDeleteMealItem = (id) => {
    message.success("已删除该记录（演示使用）");
    // In real app, call delete API and refresh
  };// Calorie calcs
const calorieTarget = computed(() => {
  if (nutritionData.value && nutritionData.value.target_calorie) return Math.floor(nutritionData.value.target_calorie);
  return Math.floor(authStore.profile?.tdee || 0);
});

const caloriesConsumed = computed(() => Math.floor(todayIntake.value?.total_energy || 0));

const checkedCalories = computed(() => {
  const records = todayIntake.value?.records || [];
  const checkedFromRecords = records
    .filter((record) => !!mealCompletionMap.value[record.ID])
    .reduce((sum, record) => sum + (record.calculatedEnergy || 0), 0);

  const recommendationTypes = ["breakfast", "lunch", "dinner", "snack"];
  const checkedFromRecommendations = recommendationTypes.reduce((sum, type) => {
    if (!recommendationCompleted.value?.[type]) return sum;
    const recommendation = dailyRecommendation.value?.[type];
    if (!recommendation) return sum;
    const energy = Number(recommendation.energy ?? recommendation.calories ?? 0);
    return sum + (Number.isFinite(energy) ? energy : 0);
  }, 0);

  return Math.floor(checkedFromRecords + checkedFromRecommendations);
});

const caloriePercentage = computed(() => {
  if (!calorieTarget.value) return 0;
  return Math.min(100, Math.max(0, (caloriesConsumed.value / calorieTarget.value) * 100));
});

const checkedCaloriePercentage = computed(() => {
  if (allRecommendedChecked.value) return 100;
  if (!calorieTarget.value) return 0;
  return Math.min(100, Math.max(0, (checkedCalories.value / calorieTarget.value) * 100));
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

const parseMacroValue = (value) => {
  const number = Number(value);
  return Number.isFinite(number) ? number : 0;
};

const getRecommendationMacro = (mealType) => {
  if (!recommendationCompleted.value?.[mealType]) {
    return { fat: 0, protein: 0, carbs: 0 };
  }
  const recommendation = dailyRecommendation.value?.[mealType];
  if (!recommendation) {
    return { fat: 0, protein: 0, carbs: 0 };
  }
  return {
    fat: parseMacroValue(recommendation.fat ?? recommendation.total_fat),
    protein: parseMacroValue(recommendation.protein ?? recommendation.total_protein),
    carbs: parseMacroValue(recommendation.carbohydrate ?? recommendation.carbs ?? recommendation.total_carbohydrate)
  };
};

const getRecommendationCheckedEnergy = (mealType) => {
  if (!recommendationCompleted.value?.[mealType]) {
    return 0;
  }
  const recommendation = dailyRecommendation.value?.[mealType];
  if (!recommendation) {
    return 0;
  }
  return parseMacroValue(recommendation.energy ?? recommendation.calories);
};

const mealTypes = ["breakfast", "lunch", "dinner", "snack"];

const allRecommendedChecked = computed(() => {
  const availableTypes = mealTypes.filter((type) => !!dailyRecommendation.value?.[type]?.id);
  if (availableTypes.length === 0) return false;
  return availableTypes.every((type) => !!recommendationCompleted.value?.[type]);
});

const checkedNutritionTotals = computed(() => {
  const records = todayIntake.value?.records || [];
  const checkedRecords = records.filter((record) => !!mealCompletionMap.value[record.ID]);

  const fromRecords = checkedRecords.reduce(
    (sum, record) => {
      sum.calories += parseMacroValue(record.calculatedEnergy);
      sum.carbs += parseMacroValue(record.calculatedCarb);
      sum.fat += parseMacroValue(record.calculatedFat);
      sum.protein += parseMacroValue(record.calculatedProtein);
      return sum;
    },
    { calories: 0, carbs: 0, fat: 0, protein: 0 }
  );

  const fromRecommendations = mealTypes.reduce(
    (sum, type) => {
      if (!recommendationCompleted.value?.[type]) return sum;
      const recommendation = dailyRecommendation.value?.[type];
      if (!recommendation) return sum;
      sum.calories += parseMacroValue(recommendation.energy ?? recommendation.calories);
      sum.carbs += parseMacroValue(recommendation.carbohydrate ?? recommendation.carbs ?? recommendation.total_carbohydrate);
      sum.fat += parseMacroValue(recommendation.fat ?? recommendation.total_fat);
      sum.protein += parseMacroValue(recommendation.protein ?? recommendation.total_protein);
      return sum;
    },
    { calories: 0, carbs: 0, fat: 0, protein: 0 }
  );

  return {
    calories: fromRecords.calories + fromRecommendations.calories,
    carbs: fromRecords.carbs + fromRecommendations.carbs,
    fat: fromRecords.fat + fromRecommendations.fat,
    protein: fromRecords.protein + fromRecommendations.protein
  };
});

const overviewCalories = computed(() => {
  const checkedValue = Math.floor(checkedNutritionTotals.value.calories);
  if (checkedValue > 0) return checkedValue;
  return Math.floor(todayIntake.value?.total_energy || 0);
});

const overviewCarbs = computed(() => {
  const checkedValue = Math.floor(checkedNutritionTotals.value.carbs);
  if (checkedValue > 0) return checkedValue;
  return Math.floor(todayIntake.value?.total_carbohydrate || 0);
});

const overviewFat = computed(() => {
  const checkedValue = Math.floor(checkedNutritionTotals.value.fat);
  if (checkedValue > 0) return checkedValue;
  return Math.floor(todayIntake.value?.total_fat || 0);
});

const overviewProtein = computed(() => {
  const checkedValue = Math.floor(checkedNutritionTotals.value.protein);
  if (checkedValue > 0) return checkedValue;
  return Math.floor(todayIntake.value?.total_protein || 0);
});

const recommendationChartTotals = computed(() => {
  return mealTypes.reduce(
    (sum, type) => {
      const recommendation = dailyRecommendation.value?.[type];
      if (!recommendation) return sum;
      sum.carbs += parseMacroValue(recommendation.carbohydrate ?? recommendation.carbs ?? recommendation.total_carbohydrate);
      sum.protein += parseMacroValue(recommendation.protein ?? recommendation.total_protein);
      sum.fat += parseMacroValue(recommendation.fat ?? recommendation.total_fat);
      return sum;
    },
    { carbs: 0, protein: 0, fat: 0 }
  );
});

const recommendationChartCarbs = computed(() => Math.floor(recommendationChartTotals.value.carbs));
const recommendationChartProtein = computed(() => Math.floor(recommendationChartTotals.value.protein));
const recommendationChartFat = computed(() => Math.floor(recommendationChartTotals.value.fat));

// Meals List formatting
const generateMealBlock = (label, type, code) => {
    const records = todayIntake.value?.records || [];
    const typeRecords = records.filter(r => r.mealType === type || r.mealType === code || r.mealType === label.toLowerCase());
    const checkedRecords = typeRecords.filter((record) => !!mealCompletionMap.value[record.ID]);
    const recommendationMacro = getRecommendationMacro(type);

    const calories = typeRecords.reduce((sum, r) => sum + (r.calculatedEnergy || 0), 0);
    const fat = typeRecords.reduce((sum, r) => sum + (r.calculatedFat || 0), 0);
    const protein = typeRecords.reduce((sum, r) => sum + (r.calculatedProtein || 0), 0);
    const carbs = typeRecords.reduce((sum, r) => sum + (r.calculatedCarb || 0), 0);
    const checkedCalories = checkedRecords.reduce((sum, r) => sum + (r.calculatedEnergy || 0), 0) + getRecommendationCheckedEnergy(type);
    const checkedFat = checkedRecords.reduce((sum, r) => sum + (r.calculatedFat || 0), 0) + recommendationMacro.fat;
    const checkedProtein = checkedRecords.reduce((sum, r) => sum + (r.calculatedProtein || 0), 0) + recommendationMacro.protein;
    const checkedCarbs = checkedRecords.reduce((sum, r) => sum + (r.calculatedCarb || 0), 0) + recommendationMacro.carbs;

    return {
      name: label,
      type: type,
      calories: Math.floor(calories),
      fat: Math.floor(fat),
      protein: Math.floor(protein),
      carbs: Math.floor(carbs),
      checkedCalories: Math.floor(checkedCalories),
      checkedFat: Math.floor(checkedFat),
      checkedProtein: Math.floor(checkedProtein),
      checkedCarbs: Math.floor(checkedCarbs),
      macroGradient: buildMacroGradient(checkedFat, checkedProtein, checkedCarbs),
      items: typeRecords.map(r => ({
        id: r.ID,
        foodName: r.foodName,
        intakeAmount: r.intakeAmount,
        image: r.imageUrl || null,
        completed: !!mealCompletionMap.value[r.ID]
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
      series: [
        {
          name: '营养素',
          type: 'pie',
          radius: ['0%', '78%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            position: 'inside',
            formatter: '{b}\n{d}%',
            color: '#ffffff',
            fontWeight: 700,
            fontSize: 12
          },
          labelLine: {
            show: false
          },
          data: recommendationChartCarbs.value === 0 && recommendationChartProtein.value === 0 && recommendationChartFat.value === 0
            ? [{ value: 1, name: '暂无数据', itemStyle: { color: '#f1f5f9' } }]
            : [
              { value: recommendationChartCarbs.value, name: '碳水', itemStyle: { color: '#fbbf24' } },
              { value: recommendationChartProtein.value, name: '蛋白', itemStyle: { color: '#fb7185' } },
              { value: recommendationChartFat.value, name: '脂肪', itemStyle: { color: '#818cf8' } }
            ]
        }
      ]
    };
    donutChart.setOption(option);
  }
};

watch([recommendationChartCarbs, recommendationChartProtein, recommendationChartFat], () => {
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
  justify-content: flex-start;
  gap: 1.25rem;
  padding: 1.5rem 0 1rem;
  border-bottom: 1px solid #e2e8f0; /* slate-200 */
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 1rem;
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

.daily-calorie-progress {
  margin: -0.2rem 0 1rem;
}

.daily-calorie-track {
  width: 100%;
  height: 5px;
  border-radius: 999px;
  background: #e2e8f0;
  overflow: hidden;
}

.daily-calorie-fill {
  height: 100%;
  border-radius: 999px;
  background: linear-gradient(90deg, #fb923c, #f97316);
  transition: width 0.3s ease;
}

.daily-calorie-text {
  margin-top: 0.4rem;
  font-size: 0.84rem;
  color: #64748b;
  font-weight: 600;
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
  align-items: center;
  gap: 1rem;
}

.recommend-btn {
  border: 1px solid #a7f3d0;
  background: #ecfdf5;
  color: #047857;
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.78rem;
  font-weight: 600;
  padding: 0.3rem 0.65rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.recommend-btn:hover {
  background: #d1fae5;
  border-color: #6ee7b7;
  transform: translateY(-1px);
}

.recommend-btn.pressed {
  transform: scale(0.96);
  background: #bbf7d0;
}

.recommend-btn.loading,
.recommend-btn:disabled {
  cursor: not-allowed;
  opacity: 0.85;
}

.recommend-icon {
  font-size: 0.95rem;
}

.action-icon {
  color: #64748b;
  font-size: 1.1rem;
  cursor: pointer;
}

.action-trigger {
  border-radius: 8px;
  padding: 4px;
  transition: all 0.2s ease;
}

.action-trigger:hover {
  background: #f1f5f9;
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
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.meal-checked-calories {
  color: #64748b;
  font-size: 0.85rem;
  font-weight: 600;
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

.stats-header.compact {
  grid-template-columns: 2fr 1fr;
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

.stats-row.compact-row {
  grid-template-columns: 2fr 1fr;
  padding: 0.5rem 0;
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
