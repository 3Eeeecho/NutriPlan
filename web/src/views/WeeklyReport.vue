<template>
  <PageLayout>
    <div class="weekly-container">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-left">
          <BackButton />
        </div>
        <div class="header-content">
          <h2 class="page-title">本周饮食趋势报告</h2>
          <p class="date-range" v-if="report.start_date">
            {{ formatDate(report.start_date) }} - {{ formatDate(report.end_date) }}
          </p>
        </div>
        <div class="header-right"></div>
      </div>

      <n-spin :show="loading">
        <!-- 周平均营养统计 -->
        <div class="section-title">本周平均摄入</div>
        <div class="stats-grid">
          <StatCard
            icon="⚡"
            :iconBg="'linear-gradient(135deg, #10b981 0%, #059669 100%)'"
            :value="report.avg_energy?.toFixed(0) || 0"
            label="平均热量"
            subtitle="kcal"
          />
          <StatCard
            icon="🍗"
            :iconBg="'linear-gradient(135deg, #ef4444 0%, #b91c1c 100%)'"
            :value="report.avg_protein?.toFixed(1) || 0"
            label="平均蛋白质"
            subtitle="g"
          />
          <StatCard
            icon="🍚"
            :iconBg="'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)'"
            :value="report.avg_carb?.toFixed(1) || 0"
            label="平均碳水"
            subtitle="g"
          />
          <StatCard
            icon="🥑"
            :iconBg="'linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%)'"
            :value="report.avg_fat?.toFixed(1) || 0"
            label="平均脂肪"
            subtitle="g"
          />
        </div>

        <!-- 趋势图表 -->
        <n-card class="chart-card" title="每日营养摄入趋势" :bordered="false">
          <div ref="chartRef" class="chart-container"></div>
        </n-card>

        <!-- 对比图表 -->
        <n-card class="chart-card" title="实际摄入 vs 目标摄入" :bordered="false">
          <div ref="compareChartRef" class="chart-container"></div>
        </n-card>
      </n-spin>
    </div>
  </PageLayout>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, NCard, NSpin } from 'naive-ui'
import * as echarts from 'echarts'
import { getWeeklyReport } from '@/api/intakeApi'
import PageLayout from '@/components/layout/PageLayout.vue'
import BackButton from '@/components/layout/BackButton.vue'
import StatCard from '@/components/ui/StatCard.vue'

const router = useRouter()
const message = useMessage()

const chartRef = ref(null)
const compareChartRef = ref(null)
let chartInstance = null
let compareChartInstance = null

const loading = ref(false)
const report = ref({})

// 设计系统颜色
const colors = {
  energy: '#10b981',
  protein: '#ef4444',
  carb: '#f59e0b',
  fat: '#8b5cf6'
}

onMounted(() => {
  loadWeeklyReport()
  window.addEventListener('resize', handleResize)
})

function handleResize() {
  chartInstance?.resize()
  compareChartInstance?.resize()
}

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
    message.error('加载周报告失败')
    console.error(error)
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
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      borderColor: '#eee',
      textStyle: { color: '#333' },
      axisPointer: { type: 'cross', lineStyle: { color: '#999' } }
    },
    legend: {
      data: ['热量', '蛋白质', '碳水化合物', '脂肪'],
      bottom: 0,
      icon: 'circle'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: { lineStyle: { color: '#e5e7eb' } },
      axisLabel: { color: '#6b7280' }
    },
    yAxis: [
      {
        type: 'value',
        name: '热量(kcal)',
        position: 'left',
        splitLine: { lineStyle: { type: 'dashed', color: '#f3f4f6' } },
        axisLabel: { color: '#9ca3af' },
        nameTextStyle: { color: colors.energy, fontWeight: 'bold' }
      },
      {
        type: 'value',
        name: '营养素(g)',
        position: 'right',
        splitLine: { show: false },
        axisLabel: { color: '#9ca3af' },
        nameTextStyle: { color: '#6b7280', fontWeight: 'bold' }
      }
    ],
    series: [
      {
        name: '热量',
        type: 'line',
        data: energyData,
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 3, color: colors.energy },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(16, 185, 129, 0.2)' },
            { offset: 1, color: 'rgba(16, 185, 129, 0.01)' }
          ])
        },
        yAxisIndex: 0
      },
      {
        name: '蛋白质',
        type: 'line',
        data: proteinData,
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: colors.protein },
        yAxisIndex: 1
      },
      {
        name: '碳水化合物',
        type: 'line',
        data: carbData,
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: colors.carb },
        yAxisIndex: 1
      },
      {
        name: '脂肪',
        type: 'line',
        data: fatData,
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2, color: colors.fat },
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

  const firstDay = report.value.daily_data[0]
  
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
      axisPointer: { type: 'shadow' },
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      textStyle: { color: '#333' },
      formatter: function(params) {
        const dataIndex = params[0].dataIndex
        let result = params[0].name + '<br/>'
        params.forEach((item, idx) => {
          const rawValue = idx === 0 ? actualRawData[dataIndex] : targetRawData[dataIndex]
          result += item.marker + item.seriesName + ': ' + Math.round(rawValue) + '<br/>'
        })
        return result
      }
    },
    legend: {
      data: ['实际摄入', '目标摄入'],
      bottom: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['热量(/10)', '蛋白质', '碳水', '脂肪'],
      axisLine: { lineStyle: { color: '#e5e7eb' } },
      axisLabel: { color: '#6b7280', fontWeight: 'bold' }
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { type: 'dashed', color: '#f3f4f6' } },
      axisLabel: { color: '#9ca3af' }
    },
    series: [
      {
        name: '实际摄入',
        type: 'bar',
        data: actualData.map(v => Math.round(v)),
        itemStyle: { 
          color: function(params) {
            const colorMap = [colors.energy, colors.protein, colors.carb, colors.fat]
            return colorMap[params.dataIndex]
          },
          borderRadius: [4, 4, 0, 0]
        },
        barGap: '20%',
        barCategoryGap: '40%'
      },
      {
        name: '目标摄入',
        type: 'bar',
        data: targetData.map(v => Math.round(v)),
        itemStyle: { 
          color: '#e5e7eb', // 灰色作为目标背景
          borderRadius: [4, 4, 0, 0]
        },
        label: {
          show: true,
          position: 'top',
          color: '#9ca3af',
          formatter: function(params) {
             return Math.round(targetRawData[params.dataIndex])
          }
        }
      }
    ]
  }

  compareChartInstance.setOption(option)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    month: 'long',
    day: 'numeric'
  })
}
</script>

<style scoped>
.weekly-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  align-items: center;
  margin-bottom: 32px;
  position: relative;
}

.header-left, .header-right {
  width: 80px; /* 占位，保证标题居中 */
}

.header-content {
  flex: 1;
  text-align: center;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 4px;
}

.date-range {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 16px;
  margin-left: 4px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.chart-card {
  margin-bottom: 24px;
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
}

.chart-container {
  width: 100%;
  height: 400px;
}

@media (max-width: 768px) {
  .weekly-container {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-left {
    width: 100%;
    text-align: left;
  }
  
  .chart-container {
    height: 300px;
  }
}
</style>
