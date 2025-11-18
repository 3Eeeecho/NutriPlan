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
                label-width="120px"
                label-position="left"
              >
                <el-form-item label="性别" prop="gender">
                  <el-radio-group v-model="profileForm.gender">
                    <el-radio label="男">男</el-radio>
                    <el-radio label="女">女</el-radio>
                  </el-radio-group>
                </el-form-item>

                <el-form-item label="年龄" prop="age">
                  <el-input-number
                    v-model="profileForm.age"
                    :min="1"
                    :max="150"
                    placeholder="请输入年龄"
                    style="width: 100%"
                  />
                </el-form-item>

                <el-form-item label="身高 (cm)" prop="height">
                  <el-input-number
                    v-model="profileForm.height"
                    :min="50"
                    :max="300"
                    :precision="1"
                    placeholder="请输入身高"
                    style="width: 100%"
                  />
                </el-form-item>

                <el-form-item label="体重 (kg)" prop="weight">
                  <el-input-number
                    v-model="profileForm.weight"
                    :min="20"
                    :max="500"
                    :precision="1"
                    placeholder="请输入体重"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-form>
            </div>

            <!-- 步骤2: 健康目标 -->
            <div v-show="currentStep === 1" class="step-content">
              <h3 class="step-title">健康目标与活动水平</h3>
              <el-form
                ref="formRef2"
                :model="profileForm"
                :rules="rules"
                label-width="120px"
                label-position="left"
              >
                <el-form-item label="健康目标" prop="health_goal">
                  <el-select
                    v-model="profileForm.health_goal"
                    placeholder="请选择健康目标"
                    style="width: 100%"
                  >
                    <el-option label="减脂" value="减脂" />
                    <el-option label="增肌" value="增肌" />
                    <el-option label="控糖" value="控糖" />
                    <el-option label="维持健康" value="维持健康" />
                  </el-select>
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

                <el-form-item label="活动水平" prop="activity_level">
                  <el-select
                    v-model="profileForm.activity_level"
                    placeholder="请选择活动水平"
                    style="width: 100%"
                  >
                    <el-option label="久坐" value="久坐" />
                    <el-option label="轻度活动" value="轻度活动" />
                    <el-option label="中度活动" value="中度活动" />
                    <el-option label="重度活动" value="重度活动" />
                    <el-option label="极重度活动" value="极重度活动" />
                  </el-select>
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
            <div v-show="currentStep === 3" class="step-content success-content">
              <el-result icon="success" title="档案保存成功！" sub-title="您的健康档案已更新">
                <template #extra>
                  <el-button type="primary" @click="goToProfileView">查看档案</el-button>
                  <el-button @click="goToHome">返回首页</el-button>
                  <el-button @click="resetForm">继续编辑</el-button>
                </template>
              </el-result>
            </div>
          </div>

          <div class="step-actions">
            <el-button v-if="currentStep > 0 && currentStep < 3" @click="prevStep">上一步</el-button>
            <el-button
              v-if="currentStep < 2"
              type="primary"
              @click="nextStep"
            >
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
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'
import { updateUserProfile } from '@/api/user'

const router = useRouter()
const authStore = useAuthStore()

const currentStep = ref(0)
const loading = ref(false)
const formRef1 = ref(null)
const formRef2 = ref(null)
const formRef3 = ref(null)

const profileForm = reactive({
  gender: '',
  age: null,
  height: null,
  weight: null,
  health_goal: '',
  target_weight: null,
  activity_level: '',
  allergies: '',
  dietary_prefs: '',
  health_conditions: '',
  meal_times_per_day: 3
})

const rules = {
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  age: [{ required: true, message: '请输入年龄', trigger: 'blur' }],
  height: [{ required: true, message: '请输入身高', trigger: 'blur' }],
  weight: [{ required: true, message: '请输入体重', trigger: 'blur' }],
  health_goal: [{ required: true, message: '请选择健康目标', trigger: 'change' }],
  activity_level: [{ required: true, message: '请选择活动水平', trigger: 'change' }]
}

onMounted(async () => {
  // 加载现有档案数据
  try {
    const profile = await authStore.loadProfile()
    if (profile) {
      Object.assign(profileForm, {
        gender: profile.gender || '',
        age: profile.age || null,
        height: profile.height || null,
        weight: profile.weight || null,
        health_goal: profile.health_goal || '',
        target_weight: profile.target_weight || null,
        activity_level: profile.activity_level || '',
        allergies: profile.allergies || '',
        dietary_prefs: profile.dietary_prefs || '',
        health_conditions: profile.health_conditions || '',
        meal_times_per_day: profile.meal_times_per_day || 3
      })
    }
  } catch (error) {
    console.error('加载档案失败:', error)
  }
})

const nextStep = async () => {
  let formRef = null
  if (currentStep.value === 0) formRef = formRef1.value
  if (currentStep.value === 1) formRef = formRef2.value

  if (formRef) {
    await formRef.validate((valid) => {
      if (valid) {
        currentStep.value++
      }
    })
  } else {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const submitForm = async () => {
  if (!formRef3.value) return

  await formRef3.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const data = { ...profileForm }
        // 移除空值
        Object.keys(data).forEach((key) => {
          if (data[key] === '' || data[key] === null) {
            delete data[key]
          }
        })

        await updateUserProfile(data)
        await authStore.loadProfile()
        ElMessage.success('档案更新成功')
        currentStep.value = 3
      } catch (error) {
        console.error('更新档案失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

const resetForm = () => {
  currentStep.value = 0
}

const goBack = () => {
  router.back()
}

const goToHome = () => {
  router.push('/home')
}

const goToProfileView = () => {
  router.push('/profile/view')
}
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
</style>

