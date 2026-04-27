<template>
  <main class="admin-page">
    <header class="admin-header">
      <div>
        <p class="eyebrow">Admin Console</p>
        <h1>管理员工作台</h1>
      </div>
      <div class="admin-account">
        <span>{{ authStore.user?.username || "-" }}</span>
        <n-tag :type="statusType" size="small">{{ statusMessage }}</n-tag>
      </div>
    </header>

    <section class="metric-grid">
      <article v-for="metric in metrics" :key="metric.label" class="metric-card">
        <span class="metric-label">{{ metric.label }}</span>
        <strong>{{ metric.value }}</strong>
        <small>{{ metric.sub }}</small>
      </article>
    </section>

    <section class="admin-workspace">
      <n-tabs v-model:value="activeTab" type="line" animated>
        <n-tab-pane name="recipes" tab="食谱管理">
          <div class="toolbar">
            <div class="toolbar-title">
              <h2>食谱库</h2>
              <span>{{ recipeTotal }} 条记录</span>
            </div>
            <n-space>
              <n-input v-model:value="recipeKeyword" clearable placeholder="搜索食谱" @keyup.enter="loadRecipes" />
              <n-button secondary :loading="recipesLoading" @click="loadRecipes">刷新</n-button>
              <n-button type="primary" @click="openRecipeModal()">新增食谱</n-button>
            </n-space>
          </div>

          <n-data-table
            :columns="recipeColumns"
            :data="recipes"
            :loading="recipesLoading"
            :pagination="recipePagination"
            :bordered="false"
            remote
          />
        </n-tab-pane>

        <n-tab-pane name="health" tab="健康数据">
          <div class="toolbar">
            <div class="toolbar-title">
              <h2>用户健康概览</h2>
              <span>平均达标率 {{ averageCompliance }}%</span>
            </div>
            <n-space>
              <n-select v-model:value="healthFilter" :options="healthFilterOptions" class="toolbar-select" />
              <n-button secondary :loading="healthLoading" @click="loadHealthStats">刷新</n-button>
            </n-space>
          </div>

          <n-data-table
            :columns="healthColumns"
            :data="filteredHealthRows"
            :loading="healthLoading"
            :pagination="{ pageSize: 8 }"
            :bordered="false"
          />
        </n-tab-pane>

        <n-tab-pane name="completion" tab="完成情况">
          <div class="toolbar">
            <div class="toolbar-title">
              <h2>食谱执行</h2>
              <span>平均完成率 {{ averageCompletion }}%</span>
            </div>
            <n-space>
              <n-select v-model:value="completionFilter" :options="completionFilterOptions" class="toolbar-select" />
              <n-button secondary :loading="completionLoading" @click="loadCompletionStats">刷新</n-button>
            </n-space>
          </div>

          <n-data-table
            :columns="completionColumns"
            :data="filteredCompletionRows"
            :loading="completionLoading"
            :pagination="{ pageSize: 8 }"
            :bordered="false"
          />
        </n-tab-pane>
      </n-tabs>
    </section>

    <n-modal v-model:show="recipeModalVisible" preset="card" :title="editingRecipeId ? '修改食谱' : '新增食谱'" class="recipe-modal">
      <n-form :model="recipeForm" label-placement="top">
        <div class="form-grid">
          <n-form-item label="食谱名称">
            <n-input v-model:value="recipeForm.name" placeholder="例如：鸡胸藜麦能量碗" />
          </n-form-item>
          <n-form-item label="餐次">
            <n-select v-model:value="recipeForm.mealType" :options="mealTypeOptions" />
          </n-form-item>
          <n-form-item label="难度">
            <n-select v-model:value="recipeForm.difficulty" :options="difficultyOptions" />
          </n-form-item>
          <n-form-item label="烹饪时间">
            <n-input-number v-model:value="recipeForm.cookingTime" :min="1" :max="240" class="full-width">
              <template #suffix>分钟</template>
            </n-input-number>
          </n-form-item>
          <n-form-item label="份量">
            <n-input-number v-model:value="recipeForm.portionWeightG" :min="1" class="full-width">
              <template #suffix>g</template>
            </n-input-number>
          </n-form-item>
          <n-form-item label="热量">
            <n-input-number v-model:value="recipeForm.energy" :min="0" class="full-width">
              <template #suffix>kcal</template>
            </n-input-number>
          </n-form-item>
          <n-form-item label="蛋白质">
            <n-input-number v-model:value="recipeForm.protein" :min="0" class="full-width">
              <template #suffix>g</template>
            </n-input-number>
          </n-form-item>
          <n-form-item label="碳水">
            <n-input-number v-model:value="recipeForm.carbohydrate" :min="0" class="full-width">
              <template #suffix>g</template>
            </n-input-number>
          </n-form-item>
          <n-form-item label="脂肪">
            <n-input-number v-model:value="recipeForm.fat" :min="0" class="full-width">
              <template #suffix>g</template>
            </n-input-number>
          </n-form-item>
          <n-form-item label="图片地址">
            <n-input v-model:value="recipeForm.imageUrl" placeholder="https://..." />
          </n-form-item>
        </div>

        <div class="flag-grid">
          <n-checkbox v-model:checked="recipeForm.isWeightLossFriendly">减脂友好</n-checkbox>
          <n-checkbox v-model:checked="recipeForm.isMuscleGainFriendly">增肌友好</n-checkbox>
          <n-checkbox v-model:checked="recipeForm.isSugarControlFriendly">控糖友好</n-checkbox>
          <n-checkbox v-model:checked="recipeForm.isGeneralFriendly">通用推荐</n-checkbox>
        </div>

        <n-form-item label="食材">
          <n-input v-model:value="recipeForm.ingredientsText" type="textarea" placeholder="每行一个食材" />
        </n-form-item>
        <n-form-item label="步骤">
          <n-input v-model:value="recipeForm.stepsText" type="textarea" placeholder="每行一个步骤" />
        </n-form-item>
        <n-form-item label="适用人群">
          <n-input v-model:value="recipeForm.targetUsersText" type="textarea" placeholder="每行一个标签" />
        </n-form-item>
        <n-form-item label="禁忌人群">
          <n-input v-model:value="recipeForm.forbiddenUsersText" type="textarea" placeholder="每行一个标签" />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="modal-actions">
          <n-button @click="recipeModalVisible = false">取消</n-button>
          <n-button type="primary" :loading="savingRecipe" @click="saveRecipe">保存</n-button>
        </div>
      </template>
    </n-modal>
  </main>
</template>

<script setup>
import { computed, h, onMounted, reactive, ref } from "vue";
import {
  NButton,
  NCheckbox,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NPopconfirm,
  NProgress,
  NSelect,
  NSpace,
  NTabPane,
  NTabs,
  NTag,
  useMessage,
} from "naive-ui";
import {
  createAdminRecipe,
  deleteAdminRecipe,
  getAdminHealthStats,
  getAdminRecipeCompletionStats,
  getAdminRecipes,
  getAdminStatus,
  updateAdminRecipe,
} from "@/api/admin";
import { useAuthStore } from "@/store/auth";

const authStore = useAuthStore();
const message = useMessage();
const activeTab = ref("recipes");
const statusMessage = ref("校验中");
const statusType = ref("warning");

const recipes = ref([]);
const recipeTotal = ref(0);
const recipeKeyword = ref("");
const recipesLoading = ref(false);
const savingRecipe = ref(false);
const recipeModalVisible = ref(false);
const editingRecipeId = ref(null);
const recipePage = ref(1);
const recipePageSize = ref(10);

const healthRows = ref([]);
const healthSummary = ref({});
const healthFilter = ref("all");
const healthLoading = ref(false);

const completionRows = ref([]);
const completionSummary = ref({});
const completionFilter = ref("all");
const completionLoading = ref(false);

const mealTypeOptions = [
  { label: "早餐", value: "早餐" },
  { label: "午餐", value: "午餐" },
  { label: "晚餐", value: "晚餐" },
  { label: "加餐", value: "加餐" },
];

const difficultyOptions = [
  { label: "简单", value: "简单" },
  { label: "中等", value: "中等" },
  { label: "进阶", value: "进阶" },
];

const healthFilterOptions = [
  { label: "全部用户", value: "all" },
  { label: "需关注", value: "attention" },
  { label: "达标稳定", value: "stable" },
];

const completionFilterOptions = [
  { label: "全部状态", value: "all" },
  { label: "已完成", value: "done" },
  { label: "进行中", value: "active" },
  { label: "未完成", value: "missed" },
];

const recipeForm = reactive({
  name: "",
  imageUrl: "",
  mealType: "午餐",
  difficulty: "简单",
  cookingTime: 20,
  portionWeightG: 100,
  energy: 0,
  protein: 0,
  carbohydrate: 0,
  fat: 0,
  ingredientsText: "",
  stepsText: "",
  targetUsersText: "",
  forbiddenUsersText: "",
  isWeightLossFriendly: false,
  isMuscleGainFriendly: false,
  isSugarControlFriendly: false,
  isGeneralFriendly: true,
});

const metrics = computed(() => [
  { label: "食谱总数", value: recipeTotal.value, sub: `${recipes.value.length} 条当前页记录` },
  { label: "健康达标率", value: `${averageCompliance.value}%`, sub: `${healthSummary.value?.attentionUsers || 0} 个用户需关注` },
  { label: "食谱完成率", value: `${averageCompletion.value}%`, sub: `${completionSummary.value?.completedPlans || 0} 个计划已完成` },
  { label: "当前角色", value: authStore.role, sub: statusMessage.value },
]);

const averageCompliance = computed(() => Math.round(healthSummary.value?.averageCompliance || 0));
const averageCompletion = computed(() => Math.round(completionSummary.value?.averageCompletion || 0));

const filteredHealthRows = computed(() => {
  if (healthFilter.value === "all") return healthRows.value;
  return healthRows.value.filter((row) => row.risk === healthFilter.value);
});

const filteredCompletionRows = computed(() => {
  if (completionFilter.value === "all") return completionRows.value;
  return completionRows.value.filter((row) => row.status === completionFilter.value);
});

const recipePagination = computed(() => ({
  page: recipePage.value,
  pageSize: recipePageSize.value,
  itemCount: recipeTotal.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page) => {
    recipePage.value = page;
    loadRecipes();
  },
  onUpdatePageSize: (pageSize) => {
    recipePageSize.value = pageSize;
    recipePage.value = 1;
    loadRecipes();
  },
}));

const recipeColumns = [
  { title: "食谱", key: "name", minWidth: 180 },
  { title: "餐次", key: "mealType", width: 90 },
  {
    title: "营养",
    key: "nutrition",
    minWidth: 210,
    render(row) {
      return `${row.energy || 0} kcal / 蛋白 ${row.protein || 0}g / 碳水 ${row.carbohydrate || 0}g / 脂肪 ${row.fat || 0}g`;
    },
  },
  {
    title: "标签",
    key: "tags",
    minWidth: 220,
    render(row) {
      const tags = [];
      if (row.isWeightLossFriendly) tags.push(["减脂", "success"]);
      if (row.isMuscleGainFriendly) tags.push(["增肌", "info"]);
      if (row.isSugarControlFriendly) tags.push(["控糖", "warning"]);
      if (row.isGeneralFriendly) tags.push(["通用", "default"]);
      if (!tags.length) tags.push(["未分类", "default"]);
      return h(NSpace, { size: 6 }, () =>
        tags.map(([label, type]) => h(NTag, { size: "small", type }, () => label)),
      );
    },
  },
  {
    title: "操作",
    key: "actions",
    width: 220,
    render(row) {
      return h(NSpace, { size: 8 }, () => [
        h(NButton, { size: "small", onClick: () => openRecipeModal(row) }, () => "编辑"),
        h(NButton, { size: "small", secondary: true, onClick: () => duplicateRecipe(row) }, () => "复制"),
        h(
          NPopconfirm,
          { onPositiveClick: () => removeRecipe(row) },
          {
            trigger: () => h(NButton, { size: "small", type: "error", secondary: true }, () => "删除"),
            default: () => "确认删除该食谱？",
          },
        ),
      ]);
    },
  },
];

const healthColumns = [
  { title: "用户", key: "username", width: 120 },
  { title: "目标", key: "healthGoal", width: 110 },
  { title: "BMI", key: "bmi", width: 90 },
  { title: "TDEE", key: "tdee", width: 100 },
  {
    title: "平均摄入",
    key: "averageIntake",
    width: 120,
    render(row) {
      return `${row.averageIntake || 0} kcal`;
    },
  },
  {
    title: "达标率",
    key: "compliance",
    minWidth: 160,
    render(row) {
      return h(NProgress, {
        type: "line",
        percentage: Math.min(Math.round(row.compliance || 0), 100),
        indicatorPlacement: "inside",
        processing: row.risk === "attention",
        status: row.risk === "attention" ? "warning" : "success",
      });
    },
  },
  {
    title: "状态",
    key: "risk",
    width: 100,
    render(row) {
      return h(NTag, { type: row.risk === "attention" ? "warning" : "success", size: "small" }, () =>
        row.risk === "attention" ? "需关注" : "稳定",
      );
    },
  },
];

const completionColumns = [
  {
    title: "日期",
    key: "planDate",
    width: 120,
    render(row) {
      return formatDate(row.planDate);
    },
  },
  { title: "用户", key: "username", width: 110 },
  { title: "计划", key: "planName", minWidth: 160 },
  {
    title: "餐次完成",
    key: "completedMeals",
    width: 130,
    render(row) {
      return `${row.completedMeals}/${row.totalMeals}`;
    },
  },
  {
    title: "完成率",
    key: "completionRate",
    minWidth: 150,
    render(row) {
      return h(NProgress, {
        type: "line",
        percentage: Math.min(Math.round(row.completionRate || 0), 100),
        indicatorPlacement: "inside",
        status: row.status === "missed" ? "error" : row.status === "done" ? "success" : "info",
      });
    },
  },
  {
    title: "热量",
    key: "calories",
    width: 170,
    render(row) {
      return `${row.calories || 0}/${row.targetCalories || 0} kcal`;
    },
  },
  {
    title: "状态",
    key: "status",
    width: 100,
    render(row) {
      const map = {
        done: { label: "已完成", type: "success" },
        active: { label: "进行中", type: "info" },
        missed: { label: "未完成", type: "error" },
      };
      const current = map[row.status] || map.active;
      return h(NTag, { type: current.type, size: "small" }, () => current.label);
    },
  },
];

const parseLines = (value) =>
  String(value || "")
    .split("\n")
    .map((item) => item.trim())
    .filter(Boolean);

const normalizeRecipe = (recipe) => ({
  ...recipe,
  id: recipe.id || recipe.ID,
  imageUrl: recipe.imageUrl || recipe.image_url || "",
  mealType: recipe.mealType || recipe.meal_type || "",
  cookingTime: recipe.cookingTime || recipe.cooking_time || 0,
  portionWeightG: recipe.portionWeightG || recipe.portion_weight_g || 100,
  cookingSteps: recipe.cookingSteps || recipe.cooking_steps || [],
  targetUsers: recipe.targetUsers || [],
  forbiddenUsers: recipe.forbiddenUsers || [],
  isWeightLossFriendly: !!recipe.isWeightLossFriendly,
  isMuscleGainFriendly: !!recipe.isMuscleGainFriendly,
  isSugarControlFriendly: !!recipe.isSugarControlFriendly,
  isGeneralFriendly: recipe.isGeneralFriendly !== false,
});

const buildRecipePayload = (overrides = {}) => ({
  name: recipeForm.name.trim(),
  imageUrl: recipeForm.imageUrl.trim(),
  mealType: recipeForm.mealType,
  allowedMealTypes: [recipeForm.mealType],
  difficulty: recipeForm.difficulty,
  cookingTime: recipeForm.cookingTime || 0,
  portionWeightG: recipeForm.portionWeightG || 100,
  energy: recipeForm.energy || 0,
  protein: recipeForm.protein || 0,
  carbohydrate: recipeForm.carbohydrate || 0,
  fat: recipeForm.fat || 0,
  ingredients: parseLines(recipeForm.ingredientsText),
  cookingSteps: parseLines(recipeForm.stepsText),
  targetUsers: parseLines(recipeForm.targetUsersText),
  forbiddenUsers: parseLines(recipeForm.forbiddenUsersText),
  isWeightLossFriendly: recipeForm.isWeightLossFriendly,
  isMuscleGainFriendly: recipeForm.isMuscleGainFriendly,
  isSugarControlFriendly: recipeForm.isSugarControlFriendly,
  isGeneralFriendly: recipeForm.isGeneralFriendly,
  ...overrides,
});

const resetRecipeForm = () => {
  editingRecipeId.value = null;
  Object.assign(recipeForm, {
    name: "",
    imageUrl: "",
    mealType: "午餐",
    difficulty: "简单",
    cookingTime: 20,
    portionWeightG: 100,
    energy: 0,
    protein: 0,
    carbohydrate: 0,
    fat: 0,
    ingredientsText: "",
    stepsText: "",
    targetUsersText: "",
    forbiddenUsersText: "",
    isWeightLossFriendly: false,
    isMuscleGainFriendly: false,
    isSugarControlFriendly: false,
    isGeneralFriendly: true,
  });
};

const openRecipeModal = (rawRecipe) => {
  if (!rawRecipe) {
    resetRecipeForm();
    recipeModalVisible.value = true;
    return;
  }

  const recipe = normalizeRecipe(rawRecipe);
  editingRecipeId.value = recipe.id;
  Object.assign(recipeForm, {
    name: recipe.name || "",
    imageUrl: recipe.imageUrl || "",
    mealType: recipe.mealType || "午餐",
    difficulty: recipe.difficulty || "简单",
    cookingTime: recipe.cookingTime || 20,
    portionWeightG: recipe.portionWeightG || 100,
    energy: recipe.energy || 0,
    protein: recipe.protein || 0,
    carbohydrate: recipe.carbohydrate || 0,
    fat: recipe.fat || 0,
    ingredientsText: (recipe.ingredients || []).join("\n"),
    stepsText: (recipe.cookingSteps || []).join("\n"),
    targetUsersText: (recipe.targetUsers || []).join("\n"),
    forbiddenUsersText: (recipe.forbiddenUsers || []).join("\n"),
    isWeightLossFriendly: recipe.isWeightLossFriendly,
    isMuscleGainFriendly: recipe.isMuscleGainFriendly,
    isSugarControlFriendly: recipe.isSugarControlFriendly,
    isGeneralFriendly: recipe.isGeneralFriendly,
  });
  recipeModalVisible.value = true;
};

const loadStatus = async () => {
  try {
    const response = await getAdminStatus();
    statusMessage.value = response.message || "已通过";
    statusType.value = "success";
  } catch {
    statusMessage.value = "校验失败";
    statusType.value = "error";
  }
};

const loadRecipes = async () => {
  recipesLoading.value = true;
  try {
    const response = await getAdminRecipes({
      page: recipePage.value,
      pageSize: recipePageSize.value,
      keyword: recipeKeyword.value.trim() || undefined,
    });
    recipes.value = (response.recipes || []).map(normalizeRecipe);
    recipeTotal.value = response.total || recipes.value.length;
  } catch {
    recipes.value = [];
    recipeTotal.value = 0;
  } finally {
    recipesLoading.value = false;
  }
};

const loadHealthStats = async () => {
  healthLoading.value = true;
  try {
    const response = await getAdminHealthStats({ limit: 200, days: 7 });
    healthSummary.value = response.summary || {};
    healthRows.value = response.users || [];
  } finally {
    healthLoading.value = false;
  }
};

const loadCompletionStats = async () => {
  completionLoading.value = true;
  try {
    const response = await getAdminRecipeCompletionStats({ limit: 200, days: 7 });
    completionSummary.value = response.summary || {};
    completionRows.value = response.plans || [];
  } finally {
    completionLoading.value = false;
  }
};

const saveRecipe = async () => {
  if (!recipeForm.name.trim()) {
    message.warning("请填写食谱名称");
    return;
  }
  if (!recipeForm.mealType) {
    message.warning("请选择餐次");
    return;
  }
  if (!recipeForm.energy || recipeForm.energy <= 0) {
    message.warning("请填写热量");
    return;
  }

  savingRecipe.value = true;
  try {
    const payload = buildRecipePayload();
    if (editingRecipeId.value) {
      await updateAdminRecipe(editingRecipeId.value, payload);
      message.success("食谱已更新");
    } else {
      await createAdminRecipe(payload);
      message.success("食谱已新增");
    }
    recipeModalVisible.value = false;
    await loadRecipes();
  } finally {
    savingRecipe.value = false;
  }
};

const duplicateRecipe = async (rawRecipe) => {
  const recipe = normalizeRecipe(rawRecipe);
  const payload = {
    name: `${recipe.name} 副本`,
    imageUrl: recipe.imageUrl || "",
    mealType: recipe.mealType || "午餐",
    allowedMealTypes: recipe.allowedMealTypes?.length ? recipe.allowedMealTypes : [recipe.mealType || "午餐"],
    difficulty: recipe.difficulty || "简单",
    cookingTime: recipe.cookingTime || 20,
    portionWeightG: recipe.portionWeightG || 100,
    energy: recipe.energy || 0,
    protein: recipe.protein || 0,
    carbohydrate: recipe.carbohydrate || 0,
    fat: recipe.fat || 0,
    ingredients: recipe.ingredients || [],
    cookingSteps: recipe.cookingSteps || [],
    targetUsers: recipe.targetUsers || [],
    forbiddenUsers: recipe.forbiddenUsers || [],
    isWeightLossFriendly: recipe.isWeightLossFriendly,
    isMuscleGainFriendly: recipe.isMuscleGainFriendly,
    isSugarControlFriendly: recipe.isSugarControlFriendly,
    isGeneralFriendly: recipe.isGeneralFriendly,
  };
  await createAdminRecipe(payload);
  message.success("已复制为新食谱");
  await loadRecipes();
};

const removeRecipe = async (rawRecipe) => {
  const recipe = normalizeRecipe(rawRecipe);
  await deleteAdminRecipe(recipe.id);
  message.success("食谱已删除");
  if (recipes.value.length === 1 && recipePage.value > 1) {
    recipePage.value -= 1;
  }
  await loadRecipes();
};

const formatDate = (value) => {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("zh-CN", { month: "2-digit", day: "2-digit" });
};

onMounted(async () => {
  await Promise.all([loadStatus(), loadRecipes(), loadHealthStats(), loadCompletionStats()]);
});
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  background: #eef5f1;
  padding: 28px;
  color: #173b31;
}

.admin-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.eyebrow {
  margin: 0 0 6px;
  color: #2f7c63;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0;
  text-transform: uppercase;
}

h1,
h2 {
  margin: 0;
  color: #173b31;
}

h1 {
  font-size: 30px;
  line-height: 1.2;
}

h2 {
  font-size: 18px;
}

.admin-account {
  min-height: 36px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #587a6f;
  font-weight: 600;
}

.metric-grid {
  margin-top: 22px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.metric-card {
  min-height: 112px;
  padding: 18px;
  border: 1px solid rgba(31, 86, 68, 0.1);
  border-radius: 8px;
  background: #ffffff;
  display: grid;
  align-content: space-between;
  box-shadow: 0 10px 22px rgba(31, 86, 68, 0.06);
}

.metric-label,
.metric-card small,
.toolbar-title span {
  color: #6d8c81;
  font-size: 13px;
}

.metric-card strong {
  margin-top: 8px;
  color: #173b31;
  font-size: 28px;
  line-height: 1.1;
}

.admin-workspace {
  margin-top: 18px;
  border: 1px solid rgba(31, 86, 68, 0.1);
  border-radius: 8px;
  background: #ffffff;
  padding: 18px;
  box-shadow: 0 12px 28px rgba(31, 86, 68, 0.07);
}

.toolbar {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
}

.toolbar-title {
  display: grid;
  gap: 4px;
}

.toolbar-select {
  width: 160px;
}

.recipe-modal {
  width: min(840px, calc(100vw - 32px));
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
}

.flag-grid {
  margin-bottom: 16px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.full-width {
  width: 100%;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

@media (max-width: 1100px) {
  .metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .form-grid,
  .flag-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .admin-page {
    padding: 18px;
  }

  .admin-header,
  .toolbar {
    display: grid;
  }

  .metric-grid,
  .form-grid,
  .flag-grid {
    grid-template-columns: 1fr;
  }

  .toolbar-select {
    width: 100%;
  }
}
</style>
