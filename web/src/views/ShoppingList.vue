<template>
  <div class="shopping-container">
    <div class="header">
      <h2>🛒 食材采购清单</h2>
      <el-button type="primary" @click="openCreateDialog">
        <el-icon><Plus /></el-icon> 新建清单
      </el-button>
    </div>

    <!-- 清单列表 -->
    <div class="list-grid" v-loading="loading">
      <el-empty v-if="!loading && lists.length === 0" description="暂无购物清单" />
      
      <div v-for="list in lists" :key="list.ID" class="list-card" @click="openDetail(list)">
        <div class="card-header">
          <span class="list-name">{{ list.name }}</span>
          <el-tag :type="list.status === 'completed' ? 'success' : 'warning'" size="small">
            {{ list.status === 'completed' ? '已完成' : '进行中' }}
          </el-tag>
        </div>
        <div class="card-info">
          <p>包含 {{ list.totalItems }} 种食材</p>
          <p class="date">{{ formatDate(list.CreatedAt) }}</p>
        </div>
        <div class="card-actions">
          <el-button link type="danger" @click.stop="handleDelete(list.ID)">删除</el-button>
        </div>
      </div>
    </div>

    <!-- 创建对话框 -->
    <el-dialog v-model="createVisible" title="新建采购清单" width="500px">
      <el-form :model="createForm" label-width="80px">
        <el-form-item label="清单名称">
          <el-input v-model="createForm.name" placeholder="例如：本周采购" />
        </el-form-item>
        <el-form-item label="选择来源">
          <el-checkbox-group v-model="createForm.sources">
            <el-checkbox label="plan">当前食谱计划</el-checkbox>
            <el-checkbox label="favorites">我的收藏</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <div v-if="sourceRecipes.length > 0" class="recipe-preview">
          <p>将包含以下食谱的食材：</p>
          <el-tag v-for="r in sourceRecipes" :key="r.id" class="mr-2 mb-2">{{ r.name }}</el-tag>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">生成清单</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" :title="currentList?.name" width="800px" destroy-on-close>
      <div class="detail-actions">
        <el-button type="success" plain @click="exportImage">
          <el-icon><Picture /></el-icon> 导出图片
        </el-button>
        <el-button type="danger" plain @click="exportPDF">
          <el-icon><Document /></el-icon> 导出 PDF
        </el-button>
        <el-button v-if="currentList?.status !== 'completed'" type="primary" @click="handleComplete">
          标记为完成
        </el-button>
      </div>

      <div id="print-area" class="shopping-detail">
        <div class="detail-header-print">
          <h3>{{ currentList?.name }}</h3>
          <p>创建时间: {{ formatDate(currentList?.CreatedAt) }}</p>
        </div>

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
              <el-checkbox v-model="item.checked" @click.stop="toggleItem(item, false)" />
              <span class="item-name">{{ item.name }}</span>
              <span class="item-amount">{{ item.amount }}</span>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Picture, Document } from '@element-plus/icons-vue'
import { getShoppingLists, createShoppingList, deleteShoppingList, updateShoppingList, completeShoppingList, getShoppingListDetail } from '@/api/shoppingApi'
import { getSelectedRecipePlan, getFavoriteList } from '@/api/recipeApi'

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
  console.log('开始加载购物清单...')
  try {
    const res = await getShoppingLists({ page: 1, pageSize: 50 })
    console.log('购物清单加载结果:', res)
    lists.value = res.lists || []
  } catch (error) {
    console.error('加载购物清单失败:', error)
    ElMessage.error('加载失败，请重试')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  console.log('ShoppingList组件已挂载')
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
    if (planRes) {
      ['breakfast', 'lunch', 'dinner', 'snack'].forEach(type => {
        const r = planRes[type + 'Recipe'] || planRes[type] // 兼容不同字段名
        if (r) recipes.push({ ...r, source: 'plan' })
      })
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
    // 避免重复
    const planIds = new Set(result.map(r => r.id))
    const favs = availableRecipes.value.filter(r => r.source === 'favorites' && !planIds.has(r.id))
    result = result.concat(favs)
  }
  return result
})

// 创建清单
const handleCreate = async () => {
  if (sourceRecipes.value.length === 0) {
    ElMessage.warning('没有可用的食谱来源')
    return
  }
  
  creating.value = true
  try {
    const recipeIds = sourceRecipes.value.map(r => r.id)
    await createShoppingList({
      listName: createForm.value.name,
      recipeIds: recipeIds
    })
    ElMessage.success('创建成功')
    createVisible.value = false
    loadLists()
  } catch (error) {
    ElMessage.error('创建失败')
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
    ElMessage.error('加载详情失败')
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
  try {
    await ElMessageBox.confirm('确定删除该清单吗？', '提示', { type: 'warning' })
    await deleteShoppingList(id)
    ElMessage.success('已删除')
    loadLists()
  } catch (e) {
    // cancel
  }
}

// 完成
const handleComplete = async () => {
  try {
    await completeShoppingList(currentList.value.ID)
    currentList.value.status = 'completed'
    ElMessage.success('清单已完成')
    loadLists()
  } catch (error) {
    ElMessage.error('操作失败')
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
    ElMessage.error('导出图片失败')
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
    ElMessage.error('导出 PDF 失败')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}
</script>

<style scoped>
.shopping-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.list-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.list-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.05);
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #eee;
}

.list-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
  border-color: #10b981;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.list-name {
  font-weight: 600;
  font-size: 16px;
}

.card-info {
  color: #666;
  font-size: 14px;
  margin-bottom: 16px;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
}

.detail-actions {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.shopping-detail {
  background: white;
  padding: 20px;
  border-radius: 8px;
}

.detail-header-print {
  display: none; /* 仅在打印/导出时显示，或者一直显示 */
  text-align: center;
  margin-bottom: 20px;
  border-bottom: 2px solid #eee;
  padding-bottom: 10px;
}

.category-group {
  margin-bottom: 24px;
}

.category-title {
  font-size: 16px;
  color: #10b981;
  margin-bottom: 12px;
  border-left: 4px solid #10b981;
  padding-left: 8px;
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.item-row {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s;
}

.item-row:hover {
  background: #f0fdf4;
}

.item-row.checked {
  opacity: 0.6;
  text-decoration: line-through;
}

.item-name {
  margin-left: 8px;
  flex: 1;
  font-weight: 500;
}

.item-amount {
  color: #666;
  font-size: 13px;
}

.recipe-preview {
  margin-top: 10px;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 4px;
}

.mr-2 { margin-right: 8px; }
.mb-2 { margin-bottom: 8px; }
</style>
