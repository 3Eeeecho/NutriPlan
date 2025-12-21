<template>
  <div class="profile-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button text @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <h1 class="page-title">健康档案管理</h1>
          <div></div>
        </div>
      </el-header>

      <el-main class="main-content">
        <el-card shadow="hover" class="profile-card">
          <el-steps :active="currentStep" finish-status="success" align-center>
            <el-step title="基础信息" icon="User" />
            <el-step title="健康目标" icon="Aim" />
            <el-step title="偏好设置" icon="Setting" />
            <el-step title="完成" icon="CircleCheck" />
          </el-steps>

          <div class="form-container">
            <!-- 步骤1: 基础信息 -->
            <div v-show="currentStep === 0" class="step-content">
              <h3 class="step-title">基础生理信息</h3>
              <el-form
                ref="formRef1"
                :model="profileForm"
                :rules="rules"
                label-position="top"
              >
                <!-- 性别选择卡片 -->
                <el-form-item label="性别" prop="gender">
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
                </el-form-item>

                <div class="input-grid">
                  <el-form-item label="年龄" prop="age">
                    <el-input-number
                      v-model="profileForm.age"
                      :min="1"
                      :max="150"
                      placeholder="岁"
                      controls-position="right"
                    />
                  </el-form-item>

                  <el-form-item label="身高 (cm)" prop="height">
                    <el-input-number
                      v-model="profileForm.height"
                      :min="50"
                      :max="300"
                      :precision="1"
                      placeholder="厘米"
                      controls-position="right"
                    />
                  </el-form-item>

                  <el-form-item label="体重 (kg)" prop="weight">
                    <el-input-number
                      v-model="profileForm.weight"
                      :min="20"
                      :max="500"
                      :precision="1"
                      placeholder="公斤"
                      controls-position="right"
                    />
                  </el-form-item>
                </div>
              </el-form>
            </div>

            <!-- 步骤2: 健康目标 -->
            <div v-show="currentStep === 1" class="step-content">
              <h3 class="step-title">健康目标与活动水平</h3>
              <el-form
                ref="formRef2"
                :model="profileForm"
                :rules="rules"
                label-position="top"
              >
                <el-form-item label="您的健康目标是什么？" prop="health_goal">
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
                </el-form-item>

                <el-form-item label="目标体重 (kg)" prop="target_weight">
                  <el-input-number
                    v-model="profileForm.target_weight"
                    :min="20"
                    :max="500"
                    :precision="1"
                    placeholder="请输入目标体重（可选）"
                    style="width: 100%"
                  />
                </el-form-item>

                <el-form-item
                  label="您平时的活动量如何？"
                  prop="activity_level"
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
                </el-form-item>
              </el-form>
            </div>

            <!-- 步骤3: 偏好设置 -->
            <div v-show="currentStep === 2" class="step-content">
              <h3 class="step-title">个性化偏好与禁忌</h3>
              <el-form
                ref="formRef3"
                :model="profileForm"
                label-width="120px"
                label-position="left"
              >
                <el-form-item label="过敏源">
                  <el-input
                    v-model="profileForm.allergies"
                    type="textarea"
                    :rows="3"
                    placeholder="请输入过敏源，多个用逗号分隔（如：海鲜,花生）"
                  />
                </el-form-item>

                <el-form-item label="饮食偏好">
                  <el-input
                    v-model="profileForm.dietary_prefs"
                    type="textarea"
                    :rows="3"
                    placeholder="请输入饮食偏好（如：喜辣,低脂）"
                  />
                </el-form-item>

                <el-form-item label="健康问题">
                  <el-input
                    v-model="profileForm.health_conditions"
                    type="textarea"
                    :rows="3"
                    placeholder="请输入健康问题（如：糖尿病,乳糖不耐）"
                  />
                </el-form-item>

                <el-form-item label="每日用餐次数" prop="meal_times_per_day">
                  <el-input-number
                    v-model="profileForm.meal_times_per_day"
                    :min="1"
                    :max="6"
                    placeholder="请输入每日用餐次数"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-form>
            </div>

            <!-- 步骤4: 完成 -->
            <div
              v-show="currentStep === 3"
              class="step-content success-content"
            >
              <el-result
                icon="success"
                title="档案保存成功！"
                sub-title="您的健康档案已更新"
              >
                <template #extra>
                  <el-button type="primary" @click="goToProfileView"
                    >查看档案</el-button
                  >
                  <el-button @click="goToHome">返回首页</el-button>
                  <el-button @click="resetForm">继续编辑</el-button>
                </template>
              </el-result>
            </div>
          </div>

          <div class="step-actions">
            <el-button
              v-if="currentStep > 0 && currentStep < 3"
              @click="prevStep"
              >上一步</el-button
            >
            <el-button v-if="currentStep < 2" type="primary" @click="nextStep">
              下一步
            </el-button>
            <el-button
              v-if="currentStep === 2"
              type="primary"
              :loading="loading"
              @click="submitForm"
            >
              保存档案
            </el-button>
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { ArrowLeft } from "@element-plus/icons-vue";
import { useAuthStore } from "@/store/auth";
import { updateUserProfile } from "@/api/user";

const router = useRouter();
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
  age: [{ required: true, message: "请输入年龄", trigger: "blur" }],
  height: [{ required: true, message: "请输入身高", trigger: "blur" }],
  weight: [{ required: true, message: "请输入体重", trigger: "blur" }],
  health_goal: [
    { required: true, message: "请选择健康目标", trigger: "change" },
  ],
  activity_level: [
    { required: true, message: "请选择活动水平", trigger: "change" },
  ],
};

const genderOptions = [
  { label: "男", value: "男", icon: "👨", desc: "Male" },
  { label: "女", value: "女", icon: "👩", desc: "Female" },
];

const goalOptions = [
  {
    label: "减脂",
    value: "减脂",
    icon: "🔥",
    desc: "降低体脂，减轻体重",
    color: "#e6a23c",
  },
  {
    label: "增肌",
    value: "增肌",
    icon: "💪",
    desc: "增加肌肉，塑造体型",
    color: "#f56c6c",
  },
  {
    label: "控糖",
    value: "控糖",
    icon: "🥗",
    desc: "稳定血糖，健康饮食",
    color: "#67c23a",
  },
  {
    label: "维持健康",
    value: "维持健康",
    icon: "❤️",
    desc: "保持当前良好状态",
    color: "#409eff",
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
  // 加载现有档案数据
  try {
    const profile = await authStore.loadProfile();
    if (profile) {
      Object.assign(profileForm, {
        gender: profile.gender || "",
        age: profile.age || null,
        height: profile.height || null,
        weight: profile.weight || null,
        health_goal: profile.health_goal || "",
        target_weight: profile.target_weight || null,
        activity_level: profile.activity_level || "",
        allergies: profile.allergies || "",
        dietary_prefs: profile.dietary_prefs || "",
        health_conditions: profile.health_conditions || "",
        meal_times_per_day: profile.meal_times_per_day || 3,
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
    await formRef.validate((valid) => {
      if (valid) {
        currentStep.value++;
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

  await formRef3.value.validate(async (valid) => {
    if (valid) {
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
        ElMessage.success("档案更新成功");
        currentStep.value = 3;
      } catch (error) {
        console.error("更新档案失败:", error);
      } finally {
        loading.value = false;
      }
    }
  });
};

const resetForm = () => {
  currentStep.value = 0;
};

const goBack = () => {
  router.back();
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
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 30px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.main-content {
  padding: 30px;
  max-width: 900px;
  margin: 0 auto;
}

.profile-card {
  border-radius: 12px;
  padding: 30px;
}

.form-container {
  margin: 40px 0;
  min-height: 400px;
}

.step-content {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.step-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid #e4e7ed;
}

.success-content {
  text-align: center;
  padding: 40px 0;
}

.step-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

:deep(.el-step__title) {
  font-size: 14px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-input-number) {
  width: 100%;
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
}

.selection-grid.two-cols {
  grid-template-columns: repeat(2, 1fr);
}

/* 选择卡片样式 */
.selection-card {
  background: #f8fafc;
  border: 2px solid #eef2f6;
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 10px;
}

.selection-card.horizontal {
  flex-direction: row;
  text-align: left;
  padding: 12px 20px;
}

.selection-card:hover {
  border-color: #cbd5e1;
  background: #f1f5f9;
}

.selection-card.active {
  border-color: #10b981;
  background: #ecfdf5;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.1);
}

.selection-card .card-icon {
  font-size: 24px;
}

.selection-card .card-icon-large {
  font-size: 36px;
}

.selection-card .card-label {
  font-weight: 600;
  color: #334155;
  font-size: 15px;
}

.selection-card.active .card-label {
  color: #059669;
}

.selection-card .card-desc {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.selection-card.active .card-desc {
  color: #10b981;
}

/* 响应式调整 */
@media (max-width: 640px) {
  .input-grid {
    grid-template-columns: 1fr;
    gap: 0;
  }
  
  .profile-container {
    padding: 0;
  }
  
  .main-content {
    padding: 16px;
  }
  
  .profile-card {
    padding: 20px 16px;
  }
  
  .selection-grid.two-cols {
    grid-template-columns: 1fr;
  }
  
  .step-actions {
    flex-direction: column-reverse;
    gap: 12px;
  }
  
  .step-actions .el-button {
    width: 100%;
    margin-left: 0 !important;
  }
}
</style>
