<template>
  <div class="intake-container">
    <!-- 顶部导航栏 -->
    <div class="top-nav">
      <el-button @click="goBack" class="back-button" circle>
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <div class="user-info-nav">
        <el-dropdown @command="handleCommand">
          <span class="user-info-display">
            <el-avatar :size="32" :icon="UserFilled" />
            <span class="username">{{ authStore.user?.username || '用户' }}</span>
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="home">
                <el-icon><HomeFilled /></el-icon>
                返回首页
              </el-dropdown-item>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>
                个人档案
              </el-dropdown-item>
              <el-dropdown-item command="weekly">
                <el-icon><TrendCharts /></el-icon>
                周报告
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <div class="page-header">
      <h2>饮食记录</h2>
      <p class="date-info">{{ currentDate }}</p>
    </div>

    <!-- 营养达标率图表 -->
    <el-card class="chart-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>今日营养达标率</span>
        </div>
      </template>
      <div ref="chartRef" class="chart-container"></div>
      <div class="nutrition-summary">
        <div class="summary-item">
          <span class="label">热量:</span>
          <span class="value">{{ nutritionStatus.total_energy?.toFixed(0) || 0 }} / {{ nutritionStatus.target_energy?.toFixed(0) || 0 }} kcal</span>
        </div>
        <div class="summary-item">
          <span class="label">蛋白质:</span>
          <span class="value">{{ nutritionStatus.total_protein?.toFixed(1) || 0 }} / {{ nutritionStatus.target_protein?.toFixed(1) || 0 }} g</span>
        </div>
        <div class="summary-item">
          <span class="label">碳水:</span>
          <span class="value">{{ nutritionStatus.total_carbohydrate?.toFixed(1) || 0 }} / {{ nutritionStatus.target_carbohydrate?.toFixed(1) || 0 }} g</span>
        </div>
        <div class="summary-item">
          <span class="label">脂肪:</span>
          <span class="value">{{ nutritionStatus.total_fat?.toFixed(1) || 0 }} / {{ nutritionStatus.target_fat?.toFixed(1) || 0 }} g</span>
        </div>
      </div>
    </el-card>

    <!-- 添加记录表单 -->
    <el-card class="form-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>添加饮食记录</span>
        </div>
      </template>
      <el-form :model="recordForm" label-width="100px">
        <el-form-item label="餐点类型" required>
          <el-select v-model="recordForm.meal_type" placeholder="请选择餐点类型" style="width: 100%">
            <el-option label="🌅 早餐" value="breakfast"></el-option>
            <el-option label="🌞 午餐" value="lunch"></el-option>
            <el-option label="🌙 晚餐" value="dinner"></el-option>
            <el-option label="🍎 加餐" value="snack"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="食物名称" required>
          <el-input v-model="recordForm.food_name" placeholder="例如：米饭、鸡胸肉、苹果等"></el-input>
        </el-form-item>
        <el-form-item label="摄入量(克)" required>
          <el-input-number 
            v-model="recordForm.intake_amount" 
            :min="1" 
            :max="9999"
            :step="10" 
            :precision="0"
            controls-position="right"
            class="full-width-number"
          ></el-input-number>
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">大约的重量即可</span>
        </el-form-item>
        
        <el-divider content-position="left">营养信息（选填）</el-divider>
        
        <el-form-item label="热量(千卡)">
          <el-input-number 
            v-model="recordForm.calculated_energy" 
            :min="0"
            :max="9999"
            :step="10" 
            :precision="0" 
            controls-position="right"
            class="full-width-number"
          ></el-input-number>
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="蛋白质(克)" label-width="90px">
              <el-input-number 
                v-model="recordForm.calculated_protein" 
                :min="0"
                :max="999"
                :step="0.1" 
                :precision="1"
                controls-position="right"
                class="full-width-number"
              ></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="碳水(克)" label-width="80px">
              <el-input-number 
                v-model="recordForm.calculated_carb" 
                :min="0"
                :max="999"
                :step="0.1" 
                :precision="1"
                controls-position="right"
                class="full-width-number"
              ></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="脂肪(克)" label-width="80px">
              <el-input-number 
                v-model="recordForm.calculated_fat" 
                :min="0"
                :max="999"
                :step="0.1" 
                :precision="1"
                controls-position="right"
                class="full-width-number"
              ></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <el-alert
          type="info"
          :closable="false"
          show-icon
          style="margin-bottom: 20px"
        >
          <template #title>
            <span style="font-size: 12px;">💡 提示：如果不清楚营养信息，可以只填写食物名称和重量，其他留空即可</span>
          </template>
        </el-alert>
        <el-form-item>
          <el-button type="primary" @click="handleAddRecord" :loading="adding" size="large">
            <el-icon style="margin-right: 5px"><Plus /></el-icon>
            添加记录
          </el-button>
          <el-button @click="resetForm" size="large">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 今日记录列表 -->
    <el-card class="list-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>今日记录</span>
          <el-tag>共 {{ nutritionStatus.records?.length || 0 }} 条</el-tag>
        </div>
      </template>
      <el-table :data="nutritionStatus.records" stripe v-loading="loading">
        <el-table-column prop="meal_type" label="餐点" width="100">
          <template #default="{ row }">
            {{ getMealTypeLabel(row.mealType) }}
          </template>
        </el-table-column>
        <el-table-column prop="food_name" label="食物名称" min-width="150">
          <template #default="{ row }">
            {{ row.foodName }}
          </template>
        </el-table-column>
        <el-table-column prop="intake_amount" label="摄入量(g)" width="100">
          <template #default="{ row }">
            {{ row.intakeAmount?.toFixed(0) }}
          </template>
        </el-table-column>
        <el-table-column prop="calculated_energy" label="热量(kcal)" width="110">
          <template #default="{ row }">
            {{ row.calculatedEnergy?.toFixed(0) }}
          </template>
        </el-table-column>
        <el-table-column prop="calculated_protein" label="蛋白质(g)" width="100">
          <template #default="{ row }">
            {{ row.calculatedProtein?.toFixed(1) }}
          </template>
        </el-table-column>
        <el-table-column prop="calculated_carb" label="碳水(g)" width="100">
          <template #default="{ row }">
            {{ row.calculatedCarb?.toFixed(1) }}
          </template>
        </el-table-column>
        <el-table-column prop="calculated_fat" label="脂肪(g)" width="100">
          <template #default="{ row }">
            {{ row.calculatedFat?.toFixed(1) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" size="small" @click="handleDelete(row.ID)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, ArrowLeft, UserFilled, User, ArrowDown, HomeFilled, TrendCharts } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { useAuthStore } from '@/store/auth'
import { getTodayStatus, addIntakeRecord, deleteIntakeRecord } from '@/api/intakeApi'

const router = useRouter()
const authStore = useAuthStore()

const chartRef = ref(null)
let chartInstance = null

const loading = ref(false)
const adding = ref(false)
const currentDate = ref('')
const nutritionStatus = ref({})

const recordForm = ref({
  meal_type: '',
  food_name: '',
  intake_amount: 100,
  calculated_energy: 0,
  calculated_protein: 0,
  calculated_carb: 0,
  calculated_fat: 0
})

onMounted(() => {
  currentDate.value = new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
  loadTodayStatus()
})

// 加载今日营养状态
async function loadTodayStatus() {
  loading.value = true
  try {
    const data = await getTodayStatus()
    nutritionStatus.value = data
    await nextTick()
    initChart()
  } catch (error) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 初始化图表
function initChart() {
  if (!chartRef.value) return

  if (chartInstance) {
    chartInstance.dispose()
  }

  chartInstance = echarts.init(chartRef.value)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      backgroundColor: 'rgba(50, 50, 50, 0.9)',
      borderColor: '#333',
      borderWidth: 1,
      textStyle: {
        color: '#fff',
        fontSize: 14
      },
      formatter: function(params) {
        const value = Math.round(params[0].value)
        let status = '🔴 未达标'
        if (value >= 90 && value <= 110) status = '✅ 已达标'
        else if (value >= 80 && value <= 120) status = '🟡 接近'
        return params[0].name + '<br/>' +
          params[0].marker + params[0].seriesName + ': ' + value + '%<br/>' +
          '<span style="font-size: 12px;">' + status + '</span>'
      }
    },
    grid: {
      left: '5%',
      right: '5%',
      bottom: '5%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['热量', '蛋白质', '碳水化合物', '脂肪'],
      axisLabel: {
        fontSize: 14,
        fontWeight: 500,
        color: '#606266'
      },
      axisLine: {
        lineStyle: {
          color: '#dcdfe6'
        }
      }
    },
    yAxis: {
      type: 'value',
      name: '达标率(%)',
      max: 150,
      axisLabel: {
        formatter: '{value}%',
        color: '#909399'
      },
      splitLine: {
        lineStyle: {
          color: '#ebeef5',
          type: 'dashed'
        }
      },
      axisLine: {
        show: false
      }
    },
    series: [
      {
        name: '达标率',
        type: 'bar',
        data: [
          Math.round(nutritionStatus.value.energy_rate || 0),
          Math.round(nutritionStatus.value.protein_rate || 0),
          Math.round(nutritionStatus.value.carbohydrate_rate || 0),
          Math.round(nutritionStatus.value.fat_rate || 0)
        ],
        itemStyle: {
          borderRadius: [8, 8, 0, 0],
          color: function(params) {
            const value = params.value
            if (value >= 90 && value <= 110) {
              return new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#95de64' },
                { offset: 1, color: '#52c41a' }
              ])
            }
            if (value >= 80 && value <= 120) {
              return new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#ffd666' },
                { offset: 1, color: '#faad14' }
              ])
            }
            return new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#ff7875' },
              { offset: 1, color: '#f5222d' }
            ])
          }
        },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}%',
          fontSize: 14,
          fontWeight: 'bold',
          color: '#303133'
        },
        barWidth: '55%',
        markLine: {
          silent: true,
          symbol: 'none',
          label: {
            position: 'end',
            formatter: '目标线',
            color: '#67C23A'
          },
          lineStyle: {
            color: '#67C23A',
            type: 'solid',
            width: 2
          },
          data: [{ yAxis: 100 }]
        }
      }
    ]
  }

  chartInstance.setOption(option)
}

// 添加记录
async function handleAddRecord() {
  if (!recordForm.value.meal_type) {
    ElMessage.warning('请选择餐点类型')
    return
  }
  if (!recordForm.value.food_name) {
    ElMessage.warning('请输入食物名称')
    return
  }
  if (recordForm.value.intake_amount <= 0) {
    ElMessage.warning('请输入摄入量')
    return
  }

  adding.value = true
  try {
    await addIntakeRecord(recordForm.value)
    ElMessage.success('添加成功')
    resetForm()
    loadTodayStatus()
  } catch (error) {
    ElMessage.error('添加失败')
  } finally {
    adding.value = false
  }
}

// 删除记录
async function handleDelete(id) {
  try {
    await deleteIntakeRecord(id)
    ElMessage.success('删除成功')
    loadTodayStatus()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// 重置表单
function resetForm() {
  recordForm.value = {
    meal_type: '',
    food_name: '',
    intake_amount: 100,
    calculated_energy: 0,
    calculated_protein: 0,
    calculated_carb: 0,
    calculated_fat: 0
  }
}

// 获取餐点类型标签
function getMealTypeLabel(type) {
  const map = {
    breakfast: '早餐',
    lunch: '午餐',
    dinner: '晚餐',
    snack: '加餐'
  }
  return map[type] || type
}

// 返回上一页
function goBack() {
  router.back()
}

// 处理下拉菜单命令
function handleCommand(command) {
  switch(command) {
    case 'home':
      router.push('/home')
      break
    case 'profile':
      router.push('/profile/view')
      break
    case 'weekly':
      router.push('/weekly')
      break
  }
}
</script>

<style scoped>
.intake-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.top-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 0 10px;
}

.back-button {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  font-size: 20px;
  transition: all 0.3s;
}

.back-button:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.1);
}

.user-info-nav {
  display: flex;
  align-items: center;
}

.user-info-display {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 24px;
  transition: all 0.3s;
}

.user-info-display:hover {
  background: rgba(255, 255, 255, 0.3);
}

.user-info-display .username {
  color: white;
  font-weight: 500;
  font-size: 14px;
}

.user-info-display .el-icon {
  color: white;
}

.page-header {
  text-align: center;
  color: white;
  margin-bottom: 30px;
}

.page-header h2 {
  font-size: 32px;
  margin-bottom: 10px;
}

.date-info {
  font-size: 16px;
  opacity: 0.9;
}

.chart-card,
.form-card,
.list-card {
  margin-bottom: 20px;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 18px;
}

.chart-container {
  width: 100%;
  height: 400px;
}

.nutrition-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 15px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.summary-item {
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
}

.summary-item .label {
  font-weight: 600;
  color: #606266;
}

.summary-item .value {
  color: #409EFF;
  font-weight: 500;
}

.el-form {
  max-width: 800px;
}

.full-width-number {
  width: 100% !important;
}

:deep(.el-input-number) {
  width: 100% !important;
}

:deep(.el-input-number .el-input__wrapper) {
  width: 100% !important;
  padding-left: 11px;
  padding-right: 11px;
}

:deep(.el-input-number .el-input__inner) {
  text-align: left !important;
  width: 100% !important;
}

:deep(.el-table) {
  font-size: 14px;
}
</style>
