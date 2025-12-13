<template>
  <div class="weekly-container">
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
              <el-dropdown-item command="intake">
                <el-icon><DataLine /></el-icon>
                饮食记录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <div class="page-header">
      <h2>本周饮食趋势报告</h2>
      <p class="date-range" v-if="report.start_date">
        {{ formatDate(report.start_date) }} - {{ formatDate(report.end_date) }}
      </p>
    </div>

    <!-- 周平均营养 -->
    <el-card class="average-card" shadow="hover" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>本周平均营养摄入</span>
        </div>
      </template>
      <div class="average-grid">
        <div class="average-item">
          <div class="icon-wrapper energy">
            <i class="el-icon-lightning"></i>
          </div>
          <div class="info">
            <div class="label">平均热量</div>
            <div class="value">{{ report.avg_energy?.toFixed(0) || 0 }} kcal</div>
          </div>
        </div>
        <div class="average-item">
          <div class="icon-wrapper protein">
            <i class="el-icon-food"></i>
          </div>
          <div class="info">
            <div class="label">平均蛋白质</div>
            <div class="value">{{ report.avg_protein?.toFixed(1) || 0 }} g</div>
          </div>
        </div>
        <div class="average-item">
          <div class="icon-wrapper carb">
            <i class="el-icon-dessert"></i>
          </div>
          <div class="info">
            <div class="label">平均碳水</div>
            <div class="value">{{ report.avg_carb?.toFixed(1) || 0 }} g</div>
          </div>
        </div>
        <div class="average-item">
          <div class="icon-wrapper fat">
            <i class="el-icon-ice-cream"></i>
          </div>
          <div class="info">
            <div class="label">平均脂肪</div>
            <div class="value">{{ report.avg_fat?.toFixed(1) || 0 }} g</div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 趋势图表 -->
    <el-card class="chart-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>每日营养摄入趋势</span>
        </div>
      </template>
      <div ref="chartRef" class="chart-container"></div>
    </el-card>

    <!-- 对比图表 -->
    <el-card class="compare-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>实际摄入 vs 目标摄入对比</span>
        </div>
      </template>
      <div ref="compareChartRef" class="chart-container"></div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, UserFilled, User, ArrowDown, HomeFilled, DataLine } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { useAuthStore } from '@/store/auth'
import { getWeeklyReport } from '@/api/intakeApi'

const router = useRouter()
const authStore = useAuthStore()

const chartRef = ref(null)
const compareChartRef = ref(null)
let chartInstance = null
let compareChartInstance = null

const loading = ref(false)
const report = ref({})

onMounted(() => {
  loadWeeklyReport()
})

// 加载周报告
async function loadWeeklyReport() {
  loading.value = true
  try {
    const data = await getWeeklyReport()
    report.value = data
    await nextTick()
    initTrendChart()
    initCompareChart()
  } catch (error) {
    ElMessage.error('加载周报告失败')
  } finally {
    loading.value = false
  }
}

// 初始化趋势图表
function initTrendChart() {
  if (!chartRef.value || !report.value.daily_data) return

  if (chartInstance) {
    chartInstance.dispose()
  }

  chartInstance = echarts.init(chartRef.value)

  const dates = report.value.daily_data.map(d => {
    const date = new Date(d.date)
    return `${date.getMonth() + 1}/${date.getDate()}`
  })

  const energyData = report.value.daily_data.map(d => d.total_energy)
  const proteinData = report.value.daily_data.map(d => d.total_protein)
  const carbData = report.value.daily_data.map(d => d.total_carbohydrate)
  const fatData = report.value.daily_data.map(d => d.total_fat)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        crossStyle: {
          color: '#999'
        }
      },
      backgroundColor: 'rgba(50, 50, 50, 0.9)',
      borderColor: '#333',
      textStyle: {
        color: '#fff'
      }
    },
    legend: {
      data: ['热量', '蛋白质', '碳水化合物', '脂肪'],
      top: 10,
      textStyle: {
        fontSize: 13,
        fontWeight: 500
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
      boundaryGap: false,
      data: dates,
      axisLabel: {
        color: '#606266',
        fontSize: 12
      },
      axisLine: {
        lineStyle: {
          color: '#dcdfe6'
        }
      }
    },
    yAxis: [
      {
        type: 'value',
        name: '热量(kcal)',
        position: 'left',
        axisLabel: {
          formatter: '{value}',
          color: '#909399'
        },
        splitLine: {
          lineStyle: {
            color: '#ebeef5',
            type: 'dashed'
          }
        },
        nameTextStyle: {
          color: '#E6A23C',
          fontWeight: 'bold'
        }
      },
      {
        type: 'value',
        name: '营养素(g)',
        position: 'right',
        axisLabel: {
          formatter: '{value}',
          color: '#909399'
        },
        splitLine: {
          show: false
        },
        nameTextStyle: {
          color: '#67C23A',
          fontWeight: 'bold'
        }
      }
    ],
    series: [
      {
        name: '热量',
        type: 'line',
        data: energyData,
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        lineStyle: {
          width: 3,
          color: '#E6A23C'
        },
        itemStyle: { 
          color: '#E6A23C',
          borderWidth: 2,
          borderColor: '#fff'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(230, 162, 60, 0.3)' },
            { offset: 1, color: 'rgba(230, 162, 60, 0.05)' }
          ])
        },
        yAxisIndex: 0
      },
      {
        name: '蛋白质',
        type: 'line',
        data: proteinData,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: {
          width: 2,
          color: '#67C23A'
        },
        itemStyle: { 
          color: '#67C23A'
        },
        yAxisIndex: 1
      },
      {
        name: '碳水化合物',
        type: 'line',
        data: carbData,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: {
          width: 2,
          color: '#409EFF'
        },
        itemStyle: { 
          color: '#409EFF'
        },
        yAxisIndex: 1
      },
      {
        name: '脂肪',
        type: 'line',
        data: fatData,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: {
          width: 2,
          color: '#F56C6C'
        },
        itemStyle: { 
          color: '#F56C6C'
        },
        yAxisIndex: 1
      }
    ]
  }

  chartInstance.setOption(option)
}

// 初始化对比图表
function initCompareChart() {
  if (!compareChartRef.value || !report.value.daily_data || report.value.daily_data.length === 0) return

  if (compareChartInstance) {
    compareChartInstance.dispose()
  }

  compareChartInstance = echarts.init(compareChartRef.value)

  // 使用本周的平均数据与目标对比
  const firstDay = report.value.daily_data[0]
  
  // 原始数据
  const actualRawData = [
    report.value.avg_energy || 0,
    report.value.avg_protein || 0,
    report.value.avg_carb || 0,
    report.value.avg_fat || 0
  ]

  const targetRawData = [
    firstDay.target_energy || 0,
    firstDay.target_protein || 0,
    firstDay.target_carbohydrate || 0,
    firstDay.target_fat || 0
  ]
  
  // 缩放后用于显示的数据（热量除以10）
  const energyScale = 10
  const actualData = [
    actualRawData[0] / energyScale,
    actualRawData[1],
    actualRawData[2],
    actualRawData[3]
  ]

  const targetData = [
    targetRawData[0] / energyScale,
    targetRawData[1],
    targetRawData[2],
    targetRawData[3]
  ]

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      backgroundColor: 'rgba(50, 50, 50, 0.9)',
      borderColor: '#333',
      textStyle: {
        color: '#fff'
      },
      formatter: function(params) {
        const dataIndex = params[0].dataIndex
        let result = params[0].name + '<br/>'
        params.forEach((item, idx) => {
          const rawValue = idx === 0 ? actualRawData[dataIndex] : targetRawData[dataIndex]
          const percent = idx === 0 && params.length > 1
            ? Math.round((actualRawData[dataIndex] / targetRawData[dataIndex]) * 100) + '%'
            : ''
          result += item.marker + item.seriesName + ': ' + Math.round(rawValue) + 
                   (percent ? ' (' + percent + ')' : '') + '<br/>'
        })
        return result
      }
    },
    legend: {
      data: ['实际摄入', '目标摄入'],
      top: 10,
      textStyle: {
        fontSize: 13,
        fontWeight: 500
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
      data: ['热量(kcal)×10', '蛋白质(g)', '碳水化合物(g)', '脂肪(g)'],
      axisLabel: {
        fontSize: 12,
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
      name: '数值',
      axisLabel: {
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
        name: '实际摄入',
        type: 'bar',
        data: actualData.map(v => Math.round(v)),
        itemStyle: { 
          color: function(params) {
            // 热量使用橙色系
            if (params.dataIndex === 0) {
              return new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#E6A23C' },
                { offset: 1, color: '#faad14' }
              ])
            }
            return new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#409EFF' },
              { offset: 1, color: '#1890ff' }
            ])
          },
          borderRadius: [8, 8, 0, 0]
        },
        label: {
          show: true,
          position: 'top',
          formatter: function(params) {
            return Math.round(actualRawData[params.dataIndex])
          },
          fontSize: 13,
          fontWeight: 'bold',
          color: function(params) {
            return params.dataIndex === 0 ? '#E6A23C' : '#409EFF'
          }
        },
        barWidth: '35%'
      },
      {
        name: '目标摄入',
        type: 'bar',
        data: targetData.map(v => Math.round(v)),
        itemStyle: { 
          color: function(params) {
            // 热量使用淡绿色系
            if (params.dataIndex === 0) {
              return new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#95de64' },
                { offset: 1, color: '#73d13d' }
              ])
            }
            return new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#67C23A' },
              { offset: 1, color: '#52c41a' }
            ])
          },
          borderRadius: [8, 8, 0, 0]
        },
        label: {
          show: true,
          position: 'top',
          formatter: function(params) {
            return Math.round(targetRawData[params.dataIndex])
          },
          fontSize: 13,
          fontWeight: 'bold',
          color: function(params) {
            return params.dataIndex === 0 ? '#73d13d' : '#67C23A'
          }
        },
        barWidth: '35%'
      }
    ]
  }

  compareChartInstance.setOption(option)
}

// 格式化日期
function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    month: 'long',
    day: 'numeric'
  })
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
    case 'intake':
      router.push('/intake')
      break
  }
}
</script>

<style scoped>
.weekly-container {
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

.date-range {
  font-size: 16px;
  opacity: 0.9;
}

.average-card,
.chart-card,
.compare-card {
  margin-bottom: 20px;
  border-radius: 12px;
}

.card-header {
  font-weight: 600;
  font-size: 18px;
}

.average-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.average-item {
  display: flex;
  align-items: center;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 12px;
  transition: transform 0.3s;
}

.average-item:hover {
  transform: translateY(-5px);
}

.icon-wrapper {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: white;
  margin-right: 15px;
}

.icon-wrapper.energy {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.icon-wrapper.protein {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.icon-wrapper.carb {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.icon-wrapper.fat {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.info {
  flex: 1;
}

.info .label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 5px;
}

.info .value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.chart-container {
  width: 100%;
  height: 450px;
}
</style>
