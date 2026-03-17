<template>
  <div class="profile-container">
    <div class="header">
      <div class="header-content">
        <n-button text @click="goBack" class="back-btn">
          <template #icon>
            <n-icon><ArrowBack /></n-icon>
          </template>
          返回
        </n-button>
        <h1 class="page-title">健康档案管理</h1>
        <div style="width: 60px"></div>
      </div>
    </div>

    <div class="main-content">
      <n-card :bordered="false" class="profile-card">
        <n-steps :current="currentStep + 1" status="process" class="mb-8">
          <n-step title="基础信息" description="生理特征" />
          <n-step title="健康目标" description="制定计划" />
          <n-step title="偏好设置" description="个性化" />
          <n-step title="完成" description="开始旅程" />
        </n-steps>

        <div class="form-container">
          <!-- 步骤1: 基础信息 -->
          <div v-if="currentStep === 0" class="step-content">
            <h3 class="step-title">基础生理信息</h3>
            <n-form
              ref="formRef1"
              :model="profileForm"
              :rules="rules"
              label-placement="top"
              size="large"
            >
              <!-- 性别选择卡片 -->
              <n-form-item label="性别" path="gender">
                <div class="selection-grid two-cols">
                  <div
                    v-for="item in genderOptions"
                    :key="item.value"
                    class="selection-card"
                    :class="{ active: profileForm.gender === item.value }"
                    @click="profileForm.gender = item.value"
                  >
                    <div class="card-icon-large">{{ item.icon }}</div>
                    <div class="card-label">{{ item.label }}</div>
                  </div>
                </div>
              </n-form-item>

              <div class="input-grid">
                <n-form-item label="年龄" path="age">
                  <n-input-number
                    v-model:value="profileForm.age"
                    :min="1"
                    :max="150"
                    placeholder="岁"
                    style="width: 100%"
                  />
                </n-form-item>

                <n-form-item label="身高 (cm)" path="height">
                  <n-input-number
                    v-model:value="profileForm.height"
                    :min="50"
                    :max="300"
                    :precision="1"
                    placeholder="厘米"
                    style="width: 100%"
                  />
                </n-form-item>

                <n-form-item label="体重 (kg)" path="weight">
                  <n-input-number
                    v-model:value="profileForm.weight"
                    :min="20"
                    :max="500"
                    :precision="1"
                    placeholder="公斤"
                    style="width: 100%"
                  />
                </n-form-item>
              </div>
            </n-form>
          </div>

          <!-- 步骤2: 健康目标 -->
          <div v-if="currentStep === 1" class="step-content">
            <h3 class="step-title">健康目标与活动水平</h3>
            <n-form
              ref="formRef2"
              :model="profileForm"
              :rules="rules"
              label-placement="top"
              size="large"
            >
              <n-form-item label="您的健康目标是什么？" path="health_goal">
                <div class="selection-grid two-cols">
                  <div
                    v-for="item in goalOptions"
                    :key="item.value"
                    class="selection-card"
                    :class="{
                      active: profileForm.health_goal === item.value,
                    }"
                    @click="profileForm.health_goal = item.value"
                  >
                    <div class="card-icon">{{ item.icon }}</div>
                    <div class="card-info">
                      <div class="card-label">{{ item.label }}</div>
                      <div class="card-desc">{{ item.desc }}</div>
                    </div>
                  </div>
                </div>
              </n-form-item>

              <n-form-item label="目标体重 (kg)" path="target_weight">
                <n-input-number
                  v-model:value="profileForm.target_weight"
                  :min="20"
                  :max="500"
                  :precision="1"
                  placeholder="请输入目标体重（可选）"
                  style="width: 100%"
                />
              </n-form-item>

              <n-form-item
                label="您平时的活动量如何？"
                path="activity_level"
              >
                <div class="selection-grid">
                  <div
                    v-for="item in activityOptions"
                    :key="item.value"
                    class="selection-card horizontal"
                    :class="{
                      active: profileForm.activity_level === item.value,
                    }"
                    @click="profileForm.activity_level = item.value"
                  >
                    <div class="card-icon">{{ item.icon }}</div>
                    <div class="card-info">
                      <div class="card-label">{{ item.label }}</div>
                      <div class="card-desc">{{ item.desc }}</div>
                    </div>
                  </div>
                </div>
              </n-form-item>
            </n-form>
          </div>

          <!-- 步骤3: 偏好设置 -->
          <div v-if="currentStep === 2" class="step-content">
            <h3 class="step-title">个性化偏好与禁忌</h3>
            <n-form
              ref="formRef3"
              :model="profileForm"
              label-placement="left"
              label-width="120"
              size="large"
            >
              <n-form-item label="过敏源">
                <n-input
                  v-model:value="profileForm.allergies"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入过敏源，多个用逗号分隔（如：海鲜,花生）"
                />
              </n-form-item>

              <n-form-item label="饮食偏好">
                <n-input
                  v-model:value="profileForm.dietary_prefs"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入饮食偏好（如：喜辣,低脂）"
                />
              </n-form-item>

              <n-form-item label="健康问题">
                <n-input
                  v-model:value="profileForm.health_conditions"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入健康问题（如：糖尿病,乳糖不耐）"
                />
              </n-form-item>

              <n-form-item label="每日用餐次数" path="meal_times_per_day">
                <n-input-number
                  v-model:value="profileForm.meal_times_per_day"
                  :min="1"
                  :max="6"
                  placeholder="请输入每日用餐次数"
                  style="width: 100%"
                />
              </n-form-item>
            </n-form>
          </div>

          <!-- 步骤4: 完成 -->
          <div
            v-if="currentStep === 3"
            class="step-content success-content"
          >
            <n-result
              status="success"
              title="档案保存成功！"
              description="您的健康档案已更新，现在可以开始您的健康饮食之旅了。"
            >
              <template #footer>
                <div class="success-actions">
                  <n-button type="primary" size="large" @click="goToProfileView">
                    查看档案
                  </n-button>
                  <n-button size="large" @click="goToHome">返回首页</n-button>
                  <n-button text @click="resetForm">继续编辑</n-button>
                </div>
              </template>
            </n-result>
          </div>
        </div>

        <div class="step-actions" v-if="currentStep < 3">
          <n-button
            v-if="currentStep > 0"
            size="large"
            @click="prevStep"
          >
            上一步
          </n-button>
          <div style="flex: 1"></div>
          <n-button 
            v-if="currentStep < 2" 
            type="primary" 
            size="large"
            @click="nextStep"
          >
            下一步
          </n-button>
          <n-button
            v-if="currentStep === 2"
            type="primary"
            size="large"
            :loading="loading"
            @click="submitForm"
          >
            保存档案
          </n-button>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { 
  NButton, NIcon, NSteps, NStep, NCard, NForm, NFormItem, 
  NInput, NInputNumber, NResult, useMessage 
} from "naive-ui";
import { ArrowBack } from "@vicons/ionicons5";
import { useAuthStore } from "@/store/auth";
import { updateUserProfile } from "@/api/user";

const router = useRouter();
const message = useMessage();
const authStore = useAuthStore();

const currentStep = ref(0);
const loading = ref(false);
const formRef1 = ref(null);
const formRef2 = ref(null);
const formRef3 = ref(null);

const profileForm = reactive({
  gender: "",
  age: null,
  height: null,
  weight: null,
  health_goal: "",
  target_weight: null,
  activity_level: "",
  allergies: "",
  dietary_prefs: "",
  health_conditions: "",
  meal_times_per_day: 3,
});

const rules = {
  gender: [{ required: true, message: "请选择性别", trigger: "change" }],
  age: [{ required: true, type: 'number', message: "请输入年龄", trigger: "blur" }],
  height: [{ required: true, type: 'number', message: "请输入身高", trigger: "blur" }],
  weight: [{ required: true, type: 'number', message: "请输入体重", trigger: "blur" }],
  health_goal: [
    { required: true, message: "请选择健康目标", trigger: "change" },
  ],
  activity_level: [
    { required: true, message: "请选择活动水平", trigger: "change" },
  ],
};

const genderOptions = [
  { label: "男", value: "男", icon: "👨" },
  { label: "女", value: "女", icon: "👩" },
];

const goalOptions = [
  {
    label: "减脂",
    value: "减脂",
    icon: "🔥",
    desc: "降低体脂，减轻体重",
  },
  {
    label: "增肌",
    value: "增肌",
    icon: "💪",
    desc: "增加肌肉，塑造体型",
  },
  {
    label: "控糖",
    value: "控糖",
    icon: "🥗",
    desc: "稳定血糖，健康饮食",
  },
  {
    label: "维持健康",
    value: "维持健康",
    icon: "❤️",
    desc: "保持当前良好状态",
  },
];

const activityOptions = [
  { label: "久坐", value: "久坐", icon: "🪑", desc: "办公室工作，几乎不运动" },
  { label: "轻度活动", value: "轻度活动", icon: "🚶", desc: "每周运动 1-3 天" },
  { label: "中度活动", value: "中度活动", icon: "🏃", desc: "每周运动 3-5 天" },
  { label: "重度活动", value: "重度活动", icon: "🏋️", desc: "每周运动 6-7 天" },
  {
    label: "极重度",
    value: "极重度活动",
    icon: "🏊",
    desc: "高强度体力劳动或专业训练",
  },
];

onMounted(async () => {
  try {
    const profile = await authStore.loadProfile();
    if (profile) {
      Object.assign(profileForm, {
        gender: profile.gender || "",
        age: profile.age || null,
        height: profile.height || null,
        weight: profile.weight || null,
        health_goal: profile.health_goal || profile.healthGoal || "",
        target_weight: profile.target_weight ?? profile.targetWeight ?? null,
        activity_level: profile.activity_level || profile.activityLevel || "",
        allergies: profile.allergies || "",
        dietary_prefs: profile.dietary_prefs || profile.dietaryPrefs || "",
        health_conditions: profile.health_conditions || profile.healthConditions || "",
        meal_times_per_day: profile.meal_times_per_day || profile.mealTimesPerDay || 3,
      });
    }
  } catch (error) {
    console.error("加载档案失败:", error);
  }
});

const nextStep = async () => {
  let formRef = null;
  if (currentStep.value === 0) formRef = formRef1.value;
  if (currentStep.value === 1) formRef = formRef2.value;

  if (formRef) {
    // Naive UI form validation
    formRef.validate((errors) => {
      if (!errors) {
        currentStep.value++;
      } else {
        message.error("请完善当前步骤信息");
      }
    });
  } else {
    currentStep.value++;
  }
};

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--;
  }
};

const submitForm = async () => {
  if (!formRef3.value) return;

  formRef3.value.validate(async (errors) => {
    if (!errors) {
      loading.value = true;
      try {
        const data = { ...profileForm };
        // 移除空值
        Object.keys(data).forEach((key) => {
          if (data[key] === "" || data[key] === null) {
            delete data[key];
          }
        });

        await updateUserProfile(data);
        await authStore.loadProfile();
        message.success("档案更新成功");
        currentStep.value = 3;
      } catch (error) {
        console.error("更新档案失败:", error);
        message.error("更新档案失败，请重试");
      } finally {
        loading.value = false;
      }
    } else {
      message.error("请检查填写信息");
    }
  });
};

const resetForm = () => {
  currentStep.value = 0;
};

const goBack = () => {
  router.push("/home");
};

const goToHome = () => {
  router.push("/home");
};

const goToProfileView = () => {
  router.push("/profile/view");
};
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  height: 64px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.main-content {
  padding: 40px 20px;
  max-width: 900px;
  margin: 0 auto;
}

.profile-card {
  border-radius: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.form-container {
  margin: 40px 0;
  min-height: 400px;
}

.step-content {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.step-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 1px solid #f3f4f6;
}

.success-content {
  padding: 60px 0;
}

.success-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 24px;
}

.step-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid #f3f4f6;
}

.mb-8 {
  margin-bottom: 32px;
}

/* 布局网格 */
.input-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-top: 10px;
}

.selection-grid {
  display: grid;
  gap: 16px;
  width: 100%;
}

.selection-grid.two-cols {
  grid-template-columns: repeat(2, 1fr);
}

/* 选择卡片样式 */
.selection-card {
  background: #f9fafb;
  border: 2px solid transparent;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 12px;
}

.selection-card.horizontal {
  flex-direction: row;
  text-align: left;
  padding: 16px 24px;
}

.selection-card:hover {
  background: #f0fdf4;
  border-color: #bbf7d0;
}

.selection-card.active {
  border-color: #10b981;
  background: #ecfdf5;
  box-shadow: 0 4px 6px -1px rgba(16, 185, 129, 0.1);
}

.selection-card .card-icon {
  font-size: 24px;
}

.selection-card .card-icon-large {
  font-size: 40px;
  margin-bottom: 4px;
}

.selection-card .card-label {
  font-weight: 600;
  color: #374151;
  font-size: 16px;
}

.selection-card.active .card-label {
  color: #059669;
}

.selection-card .card-desc {
  font-size: 13px;
  color: #6b7280;
  margin-top: 4px;
}

.selection-card.active .card-desc {
  color: #059669;
}

/* 响应式调整 */
@media (max-width: 640px) {
  .input-grid {
    grid-template-columns: 1fr;
    gap: 0;
  }
  
  .selection-grid.two-cols {
    grid-template-columns: 1fr;
  }
  
  .main-content {
    padding: 16px;
  }
}
</style>
