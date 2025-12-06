<template>
  <div class="recipe-recommend-container">
    <div class="header">
      <div class="back-btn-wrapper">
        <el-button link @click="router.push('/home')"
          class="back-btn">
          <el-icon>
            <ArrowLeft />
          </el-icon> 返回主页
        </el-button>
      </div>
      <h1>🍽️ 智能食谱推荐</h1>
      <p class="subtitle">基于您的营养需求，为您定制专属每日食谱</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>正在为您生成推荐方案...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-message">
      <p>{{ error }}</p>
      <button @click="fetchRecommendations"
        class="retry-btn">重新获取</button>
    </div>

    <!-- Recommendations -->
    <div v-else class="content">
      <!-- Selected Plan View -->
      <div v-if="currentSelectedPlan"
        class="selected-plan-view">
        <div class="selected-header">
          <h3>✅ 您已选择今日食谱</h3>
          <button @click="reselectPlan"
            class="reselect-btn">重新选择</button>
        </div>

        <div class="plan-card selected-display">
          <div class="plan-header">
            <h4>{{ currentSelectedPlan.plan_name || '今日推荐方案'
              }}</h4>
            <div class="match-score"
              :class="getScoreClass(currentSelectedPlan.match_score)">
              <span class="score-value">{{
                currentSelectedPlan.match_score }}%</span>
              <span class="score-label">匹配度</span>
            </div>
          </div>

          <div class="meals">
            <RecipeCard title="早餐"
              :recipe="currentSelectedPlan.breakfast"
              icon="🌅" />
            <RecipeCard title="午餐"
              :recipe="currentSelectedPlan.lunch"
              icon="☀️" />
            <RecipeCard title="晚餐"
              :recipe="currentSelectedPlan.dinner"
              icon="🌙" />
            <RecipeCard v-if="currentSelectedPlan.snack"
              title="加餐" :recipe="currentSelectedPlan.snack"
              icon="🍎" />
          </div>

          <div class="plan-nutrition">
            <h5>营养总计</h5>
            <div class="nutrition-bars">
              <div class="bar-item">
                <span class="bar-label">热量</span>
                <div class="bar-container">
                  <div class="bar-fill energy"
                    :style="{ width: getPercentage(currentSelectedPlan.total_energy, currentSelectedPlan.target_energy) + '%' }">
                  </div>
                </div>
                <span class="bar-value">{{
                  Math.floor(currentSelectedPlan.total_energy)
                  }} / {{
                    Math.floor(currentSelectedPlan.target_energy)
                  }}</span>
              </div>
              <div class="bar-item">
                <span class="bar-label">蛋白质</span>
                <div class="bar-container">
                  <div class="bar-fill protein"
                    :style="{ width: getPercentage(currentSelectedPlan.total_protein, currentSelectedPlan.target_protein) + '%' }">
                  </div>
                </div>
                <span class="bar-value">{{
                  Math.floor(currentSelectedPlan.total_protein)
                  }}g / {{
                    Math.floor(currentSelectedPlan.target_protein)
                  }}g</span>
              </div>
              <div class="bar-item">
                <span class="bar-label">碳水</span>
                <div class="bar-container">
                  <div class="bar-fill carb"
                    :style="{ width: getPercentage(currentSelectedPlan.total_carbohydrate, currentSelectedPlan.target_carbohydrate) + '%' }">
                  </div>
                </div>
                <span class="bar-value">{{
                  Math.floor(currentSelectedPlan.total_carbohydrate)
                  }}g / {{
                    Math.floor(currentSelectedPlan.target_carbohydrate)
                  }}g</span>
              </div>
              <div class="bar-item">
                <span class="bar-label">脂肪</span>
                <div class="bar-container">
                  <div class="bar-fill fat"
                    :style="{ width: getPercentage(currentSelectedPlan.total_fat, currentSelectedPlan.target_fat) + '%' }">
                  </div>
                </div>
                <span class="bar-value">{{
                  Math.floor(currentSelectedPlan.total_fat)
                  }}g / {{
                    Math.floor(currentSelectedPlan.target_fat)
                  }}g</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recommendation List View -->
      <div v-else>
        <!-- Nutrition Summary -->
        <div class="nutrition-summary">
          <h3>您的每日营养目标</h3>
          <div class="nutrition-grid">
            <div class="nutrition-item energy">
              <span class="label">热量</span>
              <span class="value">{{ targetNutrition.energy
                }} kcal</span>
            </div>
            <div class="nutrition-item protein">
              <span class="label">蛋白质</span>
              <span class="value">{{ targetNutrition.protein
                }} g</span>
            </div>
            <div class="nutrition-item carb">
              <span class="label">碳水</span>
              <span class="value">{{ targetNutrition.carb }}
                g</span>
            </div>
            <div class="nutrition-item fat">
              <span class="label">脂肪</span>
              <span class="value">{{ targetNutrition.fat }}
                g</span>
            </div>
          </div>
        </div>

        <!-- Recipe Plans -->
        <div class="plans-section">
          <h3>为您推荐 {{ plans.length }} 套每日食谱方案</h3>

          <div class="plans-grid">
            <div v-for="(plan, index) in plans"
              :key="plan.id || index" class="plan-card"
              :class="{ selected: selectedPlanIndex === index }"
              @click="selectPlan(index)">
              <div class="plan-header">
                <h4>方案 {{ index + 1 }}</h4>
                <div class="match-score"
                  :class="getScoreClass(plan.match_score)">
                  <span class="score-value">{{
                    plan.match_score }}%</span>
                  <span class="score-label">匹配度</span>
                </div>
              </div>

              <!-- Meals -->
              <div class="meals">
                <RecipeCard title="早餐"
                  :recipe="plan.breakfast" icon="🌅" />
                <RecipeCard title="午餐" :recipe="plan.lunch"
                  icon="☀️" />
                <RecipeCard title="晚餐" :recipe="plan.dinner"
                  icon="🌙" />
                <RecipeCard v-if="plan.snack" title="加餐"
                  :recipe="plan.snack" icon="🍎" />
              </div>

              <!-- Total Nutrition -->
              <div class="plan-nutrition">
                <h5>营养总计</h5>
                <div class="nutrition-bars">
                  <div class="bar-item">
                    <span class="bar-label">热量</span>
                    <div class="bar-container">
                      <div class="bar-fill energy"
                        :style="{ width: getPercentage(plan.total_energy, plan.target_energy) + '%' }">
                      </div>
                    </div>
                    <span class="bar-value">{{
                      Math.floor(plan.total_energy) }} / {{
                        Math.floor(plan.target_energy)
                      }}</span>
                  </div>
                  <div class="bar-item">
                    <span class="bar-label">蛋白质</span>
                    <div class="bar-container">
                      <div class="bar-fill protein"
                        :style="{ width: getPercentage(plan.total_protein, plan.target_protein) + '%' }">
                      </div>
                    </div>
                    <span class="bar-value">{{
                      Math.floor(plan.total_protein) }}g /
                      {{ Math.floor(plan.target_protein)
                      }}g</span>
                  </div>
                  <div class="bar-item">
                    <span class="bar-label">碳水</span>
                    <div class="bar-container">
                      <div class="bar-fill carb"
                        :style="{ width: getPercentage(plan.total_carbohydrate, plan.target_carbohydrate) + '%' }">
                      </div>
                    </div>
                    <span class="bar-value">{{
                      Math.floor(plan.total_carbohydrate)
                      }}g / {{
                        Math.floor(plan.target_carbohydrate)
                      }}g</span>
                  </div>
                  <div class="bar-item">
                    <span class="bar-label">脂肪</span>
                    <div class="bar-container">
                      <div class="bar-fill fat"
                        :style="{ width: getPercentage(plan.total_fat, plan.target_fat) + '%' }">
                      </div>
                    </div>
                    <span class="bar-value">{{
                      Math.floor(plan.total_fat) }}g / {{
                        Math.floor(plan.target_fat) }}g</span>
                  </div>
                </div>
              </div>

              <button v-if="selectedPlanIndex === index"
                @click.stop="confirmSelection"
                class="select-btn selected">
                ✓ 已选择此方案
              </button>
              <button v-else @click.stop="selectPlan(index)"
                class="select-btn">
                选择此方案
              </button>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="actions">
          <button @click="fetchRecommendations"
            class="action-btn secondary">
            🔄 重新推荐
          </button>
          <button @click="confirmSelection"
            class="action-btn primary"
            :disabled="selectedPlanIndex === null">
            确认选择
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { ArrowLeft } from '@element-plus/icons-vue';
import { getRecipeRecommendations, selectRecipePlan, getSelectedRecipePlan } from '../api/recipeApi';
import { getNutritionRequirements } from '../api/user';
import RecipeCard from '../components/RecipeCard.vue';
import Swal from 'sweetalert2';

export default {
  name: 'RecipeRecommend',
  components: {
    RecipeCard,
    ArrowLeft
  },
  setup() {
    const router = useRouter();
    const loading = ref(true);
    const error = ref(null);
    const plans = ref([]);
    const selectedPlanIndex = ref(null);
    const currentSelectedPlan = ref(null);
    const targetNutrition = ref({
      energy: 0,
      protein: 0,
      carb: 0,
      fat: 0
    });

    const fetchRecommendations = async (ignoreSelected = false) => {
      loading.value = true;
      error.value = null;

      try {
        // 1. 先尝试获取已选方案 (除非强制忽略)
        if (!ignoreSelected) {
          try {
            const selected = await getSelectedRecipePlan();
            if (selected && selected.id) {
              currentSelectedPlan.value = selected;
              loading.value = false;
              return; // 如果有已选方案，直接显示，不获取推荐
            }
          } catch (e) {
            // 忽略 404 或其他错误，继续获取推荐
            console.log('No selected plan found or error:', e);
          }
        }

        // 2. 获取营养需求
        const nutritionResp = await getNutritionRequirements();
        targetNutrition.value = {
          energy: Math.floor(nutritionResp.target_calorie || 0),
          protein: Math.floor(nutritionResp.protein_gram || 0),
          carb: Math.floor(nutritionResp.carb_gram || 0),
          fat: Math.floor(nutritionResp.fat_gram || 0)
        };

        // 3. 获取食谱推荐
        const recipeResp = await getRecipeRecommendations(3);
        plans.value = recipeResp.plans || [];
        selectedPlanIndex.value = null;
      } catch (err) {
        error.value = err.response?.data?.error || '获取推荐失败，请重试';
        console.error(err);
      } finally {
        loading.value = false;
      }
    };

    const reselectPlan = async () => {
      currentSelectedPlan.value = null;
      await fetchRecommendations(true);
    };

    const selectPlan = (index) => {
      selectedPlanIndex.value = index;
    };

    const confirmSelection = async () => {
      if (selectedPlanIndex.value === null) return;

      const plan = plans.value[selectedPlanIndex.value];
      try {
        await selectRecipePlan(plan.id);
        // 更新当前已选方案并刷新视图
        currentSelectedPlan.value = plan;

        // 保存成功提示
        Swal.fire({
          title: '成功',
          text: '食谱方案已保存！',
          icon: 'success',
          confirmButtonText: '确定'
        });
      } catch (err) {
        // 保存失败错误提示
        Swal.fire({
          title: '保存失败',
          text: err.response?.data?.error || '未知错误',
          icon: 'error',
          confirmButtonText: '确定'
        });
      }
    };

    const getPercentage = (actual, target) => {
      if (!target) return 0;
      return Math.min(100, Math.round((actual / target) * 100));
    };

    const getScoreClass = (score) => {
      if (score >= 90) return 'excellent';
      if (score >= 80) return 'good';
      return 'fair';
    };

    onMounted(() => {
      fetchRecommendations();
    });

    return {
      loading,
      error,
      plans,
      selectedPlanIndex,
      currentSelectedPlan,
      targetNutrition,
      fetchRecommendations,
      reselectPlan,
      selectPlan,
      confirmSelection,
      getPercentage,
      getPercentage,
      getScoreClass,
      router
    };
  }
};
</script>

<style scoped>
.recipe-recommend-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.header {
  text-align: center;
  color: white;
  margin-bottom: 3rem;
  position: relative;
}

.back-btn-wrapper {
  position: absolute;
  left: 0;
  top: 0;
}

.back-btn {
  color: white !important;
  font-size: 1rem;
  font-weight: 500;
}

.back-btn:hover {
  opacity: 0.8;
}

.header h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  font-weight: 700;
}

.subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
}

.loading {
  text-align: center;
  color: white;
  padding: 4rem 0;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-message {
  background: rgba(255, 77, 77, 0.9);
  color: white;
  padding: 2rem;
  border-radius: 12px;
  text-align: center;
}

.retry-btn {
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background: white;
  color: #ff4d4d;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
}

.content {
  background: white;
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.nutrition-summary {
  margin-bottom: 3rem;
  padding: 1.5rem;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 15px;
  color: white;
}

.nutrition-summary h3 {
  margin-bottom: 1rem;
  font-size: 1.3rem;
}

.nutrition-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.nutrition-item {
  background: rgba(255, 255, 255, 0.2);
  padding: 1rem;
  border-radius: 10px;
  text-align: center;
  backdrop-filter: blur(10px);
}

.nutrition-item .label {
  display: block;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
  opacity: 0.9;
}

.nutrition-item .value {
  display: block;
  font-size: 1.5rem;
  font-weight: 700;
}

.plans-section h3 {
  margin-bottom: 2rem;
  color: #333;
  font-size: 1.5rem;
}

.plans-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.plan-card {
  border: 3px solid transparent;
  border-radius: 15px;
  padding: 1.5rem;
  background: linear-gradient(white, white) padding-box,
    linear-gradient(135deg, #667eea, #764ba2) border-box;
  cursor: pointer;
  transition: all 0.3s ease;
}

.plan-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.plan-card.selected {
  border-color: #667eea;
  background: linear-gradient(white, white) padding-box,
    linear-gradient(135deg, #667eea, #764ba2) border-box;
  box-shadow: 0 15px 40px rgba(102, 126, 234, 0.4);
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.plan-header h4 {
  font-size: 1.3rem;
  color: #333;
}

.match-score {
  text-align: center;
  padding: 0.5rem 1rem;
  border-radius: 10px;
  background: #f0f0f0;
}

.match-score.excellent {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  color: white;
}

.match-score.good {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.match-score.fair {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  color: white;
}

.score-value {
  display: block;
  font-size: 1.5rem;
  font-weight: 700;
}

.score-label {
  display: block;
  font-size: 0.75rem;
  opacity: 0.9;
}

.meals {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.plan-nutrition {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 10px;
  margin-bottom: 1rem;
}

.plan-nutrition h5 {
  margin-bottom: 1rem;
  color: #333;
}

.nutrition-bars {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.bar-item {
  display: grid;
  grid-template-columns: 60px 1fr 100px;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
}

.bar-label {
  color: #666;
  font-weight: 500;
}

.bar-container {
  height: 8px;
  background: #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.bar-fill.energy {
  background: linear-gradient(90deg, #f093fb 0%, #f5576c 100%);
}

.bar-fill.protein {
  background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
}

.bar-fill.carb {
  background: linear-gradient(90deg, #43e97b 0%, #38f9d7 100%);
}

.bar-fill.fat {
  background: linear-gradient(90deg, #fa709a 0%, #fee140 100%);
}

.bar-value {
  color: #666;
  font-size: 0.8rem;
  text-align: right;
}

.select-btn {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #667eea;
  background: white;
  color: #667eea;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.select-btn:hover {
  background: #667eea;
  color: white;
}

.select-btn.selected {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 1rem;
  margin-top: 2rem;
}

.action-btn {
  padding: 1rem 2rem;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.action-btn.primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 20px rgba(102, 126, 234, 0.4);
}

.action-btn.primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.secondary {
  background: white;
  color: #667eea;
  border: 2px solid #667eea;
}

.action-btn.secondary:hover {
  background: #667eea;
  color: white;
}

.selected-plan-view {
  animation: fadeIn 0.5s ease;
}

.selected-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.selected-header h3 {
  font-size: 1.5rem;
  color: #333;
  margin: 0;
}

.reselect-btn {
  padding: 0.5rem 1rem;
  border: 1px solid #ff4d4d;
  background: white;
  color: #ff4d4d;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
}

.reselect-btn:hover {
  background: #ff4d4d;
  color: white;
}

.selected-display {
  border-color: #667eea;
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.2);
  cursor: default;
}

.selected-display:hover {
  transform: none;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
