<template>
  <PageLayout>
    <div class="shopping-container">
      <div class="page-header">
        <div class="header-left">
          <BackButton />
        </div>
        <div class="header-content">
          <h1 class="page-title">食材采购清单</h1>
          <p class="page-subtitle">管理您的购物计划</p>
        </div>
        <div class="header-right">
          <n-button type="primary" @click="openCreateDialog">
            <template #icon>
              <n-icon><Add /></n-icon>
            </template>
            新建清单
          </n-button>
        </div>
      </div>

      <n-spin :show="loading">
        <div v-if="!loading && lists.length === 0" class="empty-state">
          <n-empty description="暂无购物清单" size="large">
            <template #extra>
              <n-button type="primary" @click="openCreateDialog">
                创建第一个清单
              </n-button>
            </template>
          </n-empty>
        </div>

        <div v-else class="list-grid">
          <div 
            v-for="list in lists" 
            :key="list.ID" 
            class="list-card" 
            @click="openDetail(list)"
          >
            <!-- Left Side: Icon & Text -->
            <div class="card-left">
              <div class="icon-anchor">
                <n-icon size="24" color="white"><CartOutline /></n-icon>
              </div>
              <div class="text-block">
                <h3 class="list-title">{{ list.name }}</h3>
                <span class="list-subtitle">共 {{ list.totalItems }} 种食材</span>
              </div>
            </div>

            <!-- Right Side: Status & Action -->
            <div class="card-right">
              <n-tag 
                :type="list.status === 'completed' ? 'success' : 'warning'" 
                size="small" 
                round 
                :bordered="false"
                class="status-tag"
              >
                {{ list.status === 'completed' ? '已完成' : '进行中' }}
              </n-tag>
              
              <div class="delete-btn" @click.stop="handleDelete(list.ID)">
                <n-icon size="18"><TrashOutline /></n-icon>
              </div>
            </div>
          </div>
        </div>
      </n-spin>

      <!-- 创建对话框 -->
      <n-modal v-model:show="createVisible" preset="card" title="新建采购清单" style="width: 500px">
        <n-form :model="createForm" label-placement="left" label-width="80">
          <n-form-item label="清单名称">
            <n-input v-model:value="createForm.name" placeholder="例如：本周采购" />
          </n-form-item>
          <n-form-item label="选择来源">
            <n-checkbox-group v-model:value="createForm.sources">
              <n-space item-style="display: flex;">
                <n-checkbox value="plan" label="当前食谱计划" />
                <n-checkbox value="favorites" label="我的收藏" />
              </n-space>
            </n-checkbox-group>
          </n-form-item>
          <div v-if="sourceRecipes.length > 0" class="recipe-preview">
            <p class="preview-title">将包含以下食谱的食材：</p>
            <div class="tags-wrapper">
              <n-tag v-for="r in sourceRecipes" :key="r.id" size="small" :bordered="false" type="info">
                {{ r.name }}
              </n-tag>
            </div>
          </div>
        </n-form>
        <template #footer>
          <div class="dialog-footer">
            <n-button @click="createVisible = false">取消</n-button>
            <n-button type="primary" @click="handleCreate" :loading="creating">生成清单</n-button>
          </div>
        </template>
      </n-modal>

      <!-- 详情对话框 -->
      <n-modal 
        v-model:show="detailVisible" 
        preset="card" 
        :title="currentList?.name" 
        style="width: 800px; max-width: 95vw;"
        content-style="padding: 0;"
      >
        <template #header-extra>
          <div class="header-actions">
            <n-button size="small" secondary type="success" @click="exportImage">
              <template #icon><n-icon><ImageOutline /></n-icon></template>
              图片
            </n-button>
            <n-button size="small" secondary type="error" @click="exportPDF">
              <template #icon><n-icon><DocumentTextOutline /></n-icon></template>
              PDF
            </n-button>
            <n-button 
              v-if="currentList?.status !== 'completed'" 
              size="small" 
              type="primary" 
              @click="handleComplete"
            >
              <template #icon><n-icon><CheckmarkDoneOutline /></n-icon></template>
              标记完成
            </n-button>
          </div>
        </template>

        <div id="print-area" class="shopping-detail">
          <div class="detail-header-print">
            <h3>{{ currentList?.name }}</h3>
            <p>创建时间: {{ formatDate(currentList?.CreatedAt) }}</p>
          </div>

          <n-scrollbar style="max-height: 60vh">
            <div class="detail-content">
              <div v-for="(group, index) in groupedItems" :key="index" class="category-group">
                <h4 class="category-title">{{ group.category }}</h4>
                <div class="items-grid">
                  <div 
                    v-for="(item, idx) in group.items" 
                    :key="idx" 
                    class="item-row"
                    :class="{ 'checked': item.checked }"
                    @click="toggleItem(item, true)"
                  >
                    <n-checkbox 
                      :checked="item.checked" 
                      @update:checked="(v) => { item.checked = v; toggleItem(item, false) }" 
                      @click.stop
                    />
                    <div class="item-info">
                      <span class="item-name">{{ item.name }}</span>
                      <span class="item-amount">{{ item.amount }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </n-scrollbar>
        </div>
      </n-modal>
    </div>
  </PageLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  NButton, NIcon, NSpin, NEmpty, NModal, NForm, NFormItem, NInput, 
  NCheckboxGroup, NCheckbox, NSpace, NTag, NScrollbar, useMessage, useDialog
} from 'naive-ui'
import { 
  Add, CartOutline, TimeOutline, TrashOutline, ImageOutline, 
  DocumentTextOutline, CheckmarkDoneOutline 
} from '@vicons/ionicons5'
import { getShoppingLists, createShoppingList, deleteShoppingList, updateShoppingList, completeShoppingList, getShoppingListDetail } from '@/api/shoppingApi'
import { getSelectedRecipePlan, getFavoriteList } from '@/api/recipeApi'
import { useAuthStore } from '@/store/auth'
import { formatDateKey } from '@/utils/completedRecipes'
import { getRecommendationByDate } from '@/utils/recommendationCache'
import PageLayout from '@/components/layout/PageLayout.vue'
import BackButton from '@/components/layout/BackButton.vue'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const authStore = useAuthStore()

const loading = ref(false)
const lists = ref([])
const createVisible = ref(false)
const creating = ref(false)
const detailVisible = ref(false)
const currentList = ref(null)
const currentItems = ref([])

const createForm = ref({
  name: '',
  sources: ['plan']
})

const availableRecipes = ref([])

// 加载列表
const loadLists = async () => {
  loading.value = true
  try {
    const res = await getShoppingLists({ page: 1, pageSize: 50 })
    lists.value = res.lists || []
  } catch (error) {
    message.error('加载购物清单失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadLists()
})

// 打开创建对话框
const openCreateDialog = async () => {
  createForm.value.name = `采购清单 ${new Date().toLocaleDateString()}`
  createForm.value.sources = ['plan']
  createVisible.value = true
  await loadSourceRecipes()
}

// 加载源食谱
const loadSourceRecipes = async () => {
  try {
    const [planRes, favRes] = await Promise.all([
      getSelectedRecipePlan().catch(() => null),
      getFavoriteList().catch(() => ({ recipes: [] }))
    ])

    const recipes = []

    const appendPlanRecipes = (planLike) => {
      if (!planLike || typeof planLike !== 'object') return
      ;['breakfast', 'lunch', 'dinner', 'snack'].forEach((type) => {
        const recipe = planLike[`${type}Recipe`] || planLike[type]
        if (!recipe?.id) return
        recipes.push({ ...recipe, source: 'plan' })
      })
    }

    if (planRes) {
      appendPlanRecipes(planRes)
    } else {
      const todayRecommendation = getRecommendationByDate(authStore.user?.id, formatDateKey(new Date()))
      appendPlanRecipes(todayRecommendation)
    }
    
    if (favRes && favRes.recipes) {
      favRes.recipes.forEach(r => recipes.push({ ...r, source: 'favorites' }))
    }
    
    availableRecipes.value = recipes
  } catch (error) {
    console.error(error)
  }
}

// 计算选中的食谱
const sourceRecipes = computed(() => {
  const sources = createForm.value.sources
  let result = []
  if (sources.includes('plan')) {
    result = result.concat(availableRecipes.value.filter(r => r.source === 'plan'))
  }
  if (sources.includes('favorites')) {
    const planIds = new Set(result.map(r => r.id))
    const favs = availableRecipes.value.filter(r => r.source === 'favorites' && !planIds.has(r.id))
    result = result.concat(favs)
  }
  return result
})

// 创建清单
const handleCreate = async () => {
  if (sourceRecipes.value.length === 0) {
    message.warning('没有可用的食谱来源，请先生成今日推荐或添加收藏')
    return
  }
  
  creating.value = true
  try {
    const recipeIds = sourceRecipes.value.map(r => r.id)
    await createShoppingList({
      listName: createForm.value.name,
      recipeIds: recipeIds
    })
    message.success('创建成功')
    createVisible.value = false
    loadLists()
  } catch (error) {
    message.error('创建失败')
  } finally {
    creating.value = false
  }
}

// 打开详情
const openDetail = async (list) => {
  try {
    const res = await getShoppingListDetail(list.ID)
    currentList.value = res.list
    currentItems.value = res.items || []
    detailVisible.value = true
  } catch (error) {
    message.error('加载详情失败')
  }
}

// 分组显示
const groupedItems = computed(() => {
  const groups = {}
  const order = ['蔬菜类', '肉类海鲜', '豆制品', '蛋奶类', '主食类', '水果类', '调料类', '其他']
  
  currentItems.value.forEach(item => {
    if (!groups[item.category]) {
      groups[item.category] = []
    }
    groups[item.category].push(item)
  })

  return order.map(cat => ({
    category: cat,
    items: groups[cat] || []
  })).filter(g => g.items.length > 0)
})

// 切换选中状态
const toggleItem = async (item, fromRow = false) => {
  if (fromRow) {
    item.checked = !item.checked
  }
  try {
    await updateShoppingList(currentList.value.ID, currentItems.value)
  } catch (error) {
    console.error('保存状态失败')
  }
}

// 删除
const handleDelete = async (id) => {
  dialog.warning({
    title: '删除确认',
    content: '确定删除该清单吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteShoppingList(id)
        message.success('已删除')
        loadLists()
      } catch (e) {
        message.error('删除失败')
      }
    }
  })
}

// 完成
const handleComplete = async () => {
  try {
    await completeShoppingList(currentList.value.ID)
    currentList.value.status = 'completed'
    message.success('清单已完成')
    loadLists()
  } catch (error) {
    message.error('操作失败')
  }
}

// 导出图片
const exportImage = async () => {
  const el = document.getElementById('print-area')
  if (!el) return
  
  try {
    const html2canvas = (await import('html2canvas')).default
    const canvas = await html2canvas(el, {
      useCORS: true,
      scale: 2,
      backgroundColor: '#ffffff'
    })
    const link = document.createElement('a')
    link.download = `${currentList.value.name}.png`
    link.href = canvas.toDataURL()
    link.click()
  } catch (error) {
    console.error(error)
    message.error('导出图片失败')
  }
}

// 导出 PDF
const exportPDF = async () => {
  const el = document.getElementById('print-area')
  if (!el) return

  try {
    const html2canvas = (await import('html2canvas')).default
    const { jsPDF } = await import('jspdf')
    
    const canvas = await html2canvas(el, {
      scale: 2,
      backgroundColor: '#ffffff'
    })
    const imgData = canvas.toDataURL('image/png')
    
    const pdf = new jsPDF('p', 'mm', 'a4')
    const pdfWidth = pdf.internal.pageSize.getWidth()
    const pdfHeight = pdf.internal.pageSize.getHeight()
    const imgWidth = pdfWidth
    const imgHeight = (canvas.height * pdfWidth) / canvas.width

    let heightLeft = imgHeight
    let position = 0

    pdf.addImage(imgData, 'PNG', 0, position, imgWidth, imgHeight)
    heightLeft -= pdfHeight

    while (heightLeft >= 0) {
      position = heightLeft - imgHeight
      pdf.addPage()
      pdf.addImage(imgData, 'PNG', 0, position, imgWidth, imgHeight)
      heightLeft -= pdfHeight
    }

    pdf.save(`${currentList.value.name}.pdf`)
  } catch (error) {
    message.error('导出 PDF 失败')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}
</script>

<style scoped>
.shopping-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32px;
}

.header-left, .header-right {
  width: 100px; /* 保持标题居中 */
}

.header-right {
  display: flex;
  justify-content: flex-end;
}

.header-content {
  text-align: center;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
}

.page-subtitle {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 4px 0 0;
}

.empty-state {
  padding: 60px 0;
  display: flex;
  justify-content: center;
}

.list-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.list-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  border: 1px solid #f3f4f6;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  height: 90px; /* Fixed height for consistency */
}

.list-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  border-color: #10b981;
}

/* Left Side */
.card-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.icon-anchor {
  width: 48px;
  height: 48px;
  background-color: #10b981;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 6px -1px rgba(16, 185, 129, 0.2);
}

.text-block {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.list-title {
  font-size: 18px;
  font-weight: 700;
  color: #1f2937;
  margin: 0 0 4px 0;
  line-height: 1.2;
}

.list-subtitle {
  font-size: 14px;
  color: #6b7280;
}

/* Right Side */
.card-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: space-between;
  height: 100%;
  padding: 2px 0;
}

.status-tag {
  font-weight: 600;
}

.delete-btn {
  color: #9ca3af;
  cursor: pointer;
  transition: color 0.2s;
  padding: 4px;
  display: flex;
  align-items: center;
}

.delete-btn:hover {
  color: #ef4444;
}

/* Dialog Styles */
.recipe-preview {
  margin-top: 16px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.preview-title {
  font-size: 13px;
  color: #6b7280;
  margin: 0 0 8px 0;
}

.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.shopping-detail {
  padding: 24px;
}

.detail-header-print {
  display: none;
  text-align: center;
  margin-bottom: 24px;
  border-bottom: 2px solid #eee;
  padding-bottom: 16px;
}

.category-group {
  margin-bottom: 24px;
}

.category-title {
  font-size: 15px;
  font-weight: 600;
  color: #10b981;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
}

.category-title::before {
  content: '';
  display: block;
  width: 4px;
  height: 16px;
  background: #10b981;
  margin-right: 8px;
  border-radius: 2px;
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
}

.item-row {
  display: flex;
  align-items: center;
  padding: 10px;
  background: #f9fafb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.item-row:hover {
  background: #f0fdf4;
  border-color: #bbf7d0;
}

.item-row.checked {
  opacity: 0.6;
  background: #f3f4f6;
}

.item-row.checked .item-name {
  text-decoration: line-through;
  color: #9ca3af;
}

.item-info {
  margin-left: 12px;
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-name {
  font-weight: 500;
  color: #374151;
}

.item-amount {
  font-size: 12px;
  color: #6b7280;
  background: white;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
}

@media (max-width: 768px) {
  .shopping-container {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-left, .header-right {
    width: 100%;
    justify-content: center;
  }
  
  .items-grid {
    grid-template-columns: 1fr;
  }
}
</style>
