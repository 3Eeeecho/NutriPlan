<template>
  <div class="nutrition-dashboard">
    <header class="top-header glass-card">
      <div class="header-main">
        <p class="header-greet">{{ 今日问候语 }}</p>
        <h1 class="header-title">{{ 今日看板标题 }}</h1>
        <p class="header-sub">NutriPlan · 基于健康数据的智能饮食推荐</p>
      </div>
      <div class="header-meta">
        <span class="meta-icon">📅</span>
        <div class="meta-text">
          <strong>{{ 今日日期文案 }}</strong>
          <span>{{ 今日星期文案 }}</span>
        </div>
      </div>
    </header>

    <main class="dashboard-grid">
      <section class="center-panel">
        <p v-if="加载中" class="loading-tip">正在同步饮食数据...</p>

        <div class="record-ops glass-card">
          <div class="ops-left">
            <button class="icon-circle soft" aria-label="上一天记录" @click="上一天"><n-icon><ChevronBackOutline /></n-icon></button>
            <button class="icon-circle soft" aria-label="下一天记录" :disabled="是否今天" @click="下一天"><n-icon><ChevronForwardOutline /></n-icon></button>
            <p class="ops-date">{{ 日期标题 }}</p>
          </div>
          <label class="calendar-picker">
            <n-icon><CalendarOutline /></n-icon>
            <input type="date" :value="日期输入值" :max="今日日期输入值" @change="选择日期" />
          </label>
        </div>

        <div class="recipe-action-bar glass-card" v-if="是否今天">
          <div class="action-info">
            <h3 class="plan-status">{{ 今日方案 ? '今日食谱已生成' : '开始规划健康饮食' }}</h3>
            <p class="plan-desc">{{ 今日方案 ? '按计划饮食，保持活力每一天' : '为您量身定制的营养均衡食谱方案' }}</p>
            <div class="diet-mode-chip" :class="`mode-${当前饮食模式.mode || 'normal'}`" @click="去模式设置" title="点击选择模式">
              <span class="mode-title">当前模式：{{ 当前饮食模式文案 }}</span>
              <span class="mode-meta" v-if="当前饮食模式.until">至 {{ 格式化模式日期(当前饮食模式.until) }}</span>
              <span class="mode-meta" v-else>{{ 当前饮食模式说明 }}</span>
            </div>
          </div>
          <div class="recipe-action-buttons">
            <button v-if="false" class="commercial-btn secondary-btn" @click="去受限单餐重构">
              <n-icon><SparklesOutline /></n-icon>
              受限单餐重构
            </button>
            <button class="commercial-btn" @click="去生成食谱">
              <n-icon><SparklesOutline /></n-icon>
              {{ 今日方案 ? '重新生成食谱' : '智能生成今日食谱' }}
            </button>
          </div>
        </div>

        <div class="timeline" :style="时间线样式">
          <article
            v-for="meal in 时间线餐次"
            :key="meal.type"
            class="timeline-row"
            :class="[
              `timeline-${meal.phase}`,
              { 'is-drag-over': 拖拽目标餐次 === meal.type, 'is-active-detail': 详情餐次类型 === meal.type }
            ]"
            :draggable="可拖拽餐次 === meal.type"
            @dragstart="开始拖拽餐次(meal.type)"
            @dragover.prevent="进入拖拽目标(meal.type)"
            @dragleave="离开拖拽目标(meal.type)"
            @drop="放置餐次(meal.type)"
            @dragend="结束拖拽餐次"
          >
            <div class="time-col">
              <span class="dot" />
              <div class="time-editor" @click.stop>
                <button v-if="编辑时间餐次 !== meal.type" class="time-display-btn" @click="开始编辑时间(meal.type)">{{ meal.time }}</button>
                <input
                  v-else
                  type="time"
                  class="time-input"
                  :value="meal.time"
                  @input="(event) => 更新餐次时间(meal.type, event.target.value)"
                  @blur="结束编辑时间"
                  @keyup.enter="结束编辑时间"
                />
              </div>
            </div>

            <div class="timeline-content glass-card meal-card" :class="`meal-${meal.type}`" @click="切换详情餐次(meal.type)">
              <div class="meal-shell">
                <div class="meal-main">
                  <div class="meal-head">
                    <div class="meal-title-wrap">
                      <input
                        v-if="编辑名称餐次 === meal.type"
                        v-model="编辑名称内容"
                        class="meal-name-input"
                        maxlength="16"
                        @blur="保存名称编辑(meal.type)"
                        @keyup.enter="保存名称编辑(meal.type)"
                        @click.stop
                      />
                      <h3
                        v-else
                        class="editable-meal-name"
                        title="点击修改名称"
                        @click.stop="开始编辑名称(meal.type, meal.name)"
                      >
                        {{ meal.name }}
                      </h3>
                      <span class="meal-time-inline">{{ meal.time }}</span>
                    </div>
                    <div class="drag-handle" title="按住拖拽排序" @mousedown.stop="允许拖拽(meal.type)" @mouseup.stop="禁止拖拽" @mouseleave="禁止拖拽" @click.stop>
                      <n-icon><MenuOutline /></n-icon>
                    </div>
                  </div>

                  <div class="meal-food-row">
                    <div
                      v-for="(item, index) in meal.previewItems"
                      :key="`${meal.type}-preview-${index}-${item.name}`"
                      class="food-preview-item"
                    >
                      <div class="food-image-card" :class="{ 'is-link': !!item.recipeId }" @click.stop="打开菜谱详情(item)">
                        <RecipeImage
                          class="food-preview-image"
                          :src="item.image"
                          :name="item.name"
                          :meal-type="meal.type"
                          :alt="item.name || '餐次缩略图'"
                        />
                      </div>
                      <p class="food-preview-name" :class="{ 'is-link': !!item.recipeId }" @click.stop="打开菜谱详情(item)">{{ item.name }}</p>
                    </div>
                  </div>
                </div>

                <div class="meal-side" @click.stop="打开记录悬浮栏(meal.type)">
                  <div class="mini-ring" :style="{ background: `conic-gradient(#8ec662 0 ${meal.progress}%, #e8f1e5 ${meal.progress}% 100%)` }">
                    <div class="mini-ring-inner">{{ meal.progress }}%</div>
                  </div>
                  <p class="meal-summary">{{ meal.display_energy }} kcal</p>
                  <p class="meal-summary-sub">蛋白质 {{ meal.display_protein }}g</p>
                  <button class="record-action-btn" aria-label="立即记录餐次">
                    立即记录
                  </button>
                </div>
              </div>

              <div v-if="meal.variant === 'logged'" class="footer-tip">状态：已记录</div>
            </div>
          </article>
        </div>
      </section>

      <section class="right-panel">
        <div class="overview-card glass-card">
          <div class="ring-wrap">
            <div class="calorie-ring" :style="ringStyle">
              <div class="ring-center">
                <h3>{{ 剩余热量 }}</h3>
                <p>剩余千卡</p>
              </div>
            </div>
            <p class="goal-text">每日目标：{{ 每日目标热量 }} 千卡</p>
          </div>

          <div class="macro-stack">
            <article v-for="macro in 宏量卡片" :key="macro.label" class="macro-card">
              <div class="macro-left">
                <span class="macro-icon" :style="{ background: macro.soft }">{{ macro.icon }}</span>
                <div>
                  <p class="macro-label">{{ macro.label }}</p>
                  <p class="macro-target">目标：{{ macro.target }}</p>
                </div>
              </div>
              <p class="macro-value">{{ macro.value }}</p>
            </article>
          </div>
        </div>

        <div class="details-card details-stage glass-card">
          <Transition name="detail-panel" mode="out-in">
            <div :key="当前详情餐次?.type || 详情餐次类型" class="details-body">
              <div class="details-head">
                <h3>{{ 当前详情标题 }}</h3>
                <p>{{ 当前详情时间 }}</p>
              </div>

              <div v-if="!当前详情食物.length" class="sample-note-wrapper">
                <p class="sample-note">
                  {{ 今日方案 ? '今日暂无记录，以下为今日推荐食谱' : '今日暂无记录，以下为推荐食谱' }}
                </p>
                <button v-if="是否今天 && !今日方案" class="generate-btn" @click="去生成食谱">
                  <n-icon><SparklesOutline /></n-icon> 智能生成今日食谱
                </button>
              </div>

              <div class="micro-grid">
                <div v-for="line in 详情营养条" :key="line.label" class="micro-item">
                  <div class="micro-top">
                    <span>{{ line.label }}</span>
                    <span>{{ line.value }}</span>
                  </div>
                  <div class="micro-mini-bars" aria-hidden="true">
                    <span
                      v-for="(height, index) in line.spark"
                      :key="`${line.label}-${index}`"
                      class="mini-bar"
                      :style="{ height: `${height}%`, background: line.color }"
                    />
                  </div>
                  <div class="micro-track">
                    <div class="micro-fill" :style="{ width: line.progress + '%', background: line.color }" />
                  </div>
                </div>
              </div>

              <div class="food-table">
                <div class="table-head">
                  <span class="check-col"></span>
                  <span>食物项</span>
                  <span>重量(g)</span>
                  <span>热量</span>
                  <span>蛋白质</span>
                  <span>碳水</span>
                  <span>脂肪</span>
                  <span class="action-col"></span>
                </div>

                <div v-for="food in 展示详情食物" :key="food.id" class="food-row">
                  <label class="check-wrap">
                    <input
                      type="checkbox"
                      class="food-check-input"
                      :checked="是否食物已勾选(food)"
                      :disabled="!是否可切换食物勾选(food)"
                      @change="切换食物状态(food)"
                    />
                  </label>
                  <button
                    type="button"
                    class="food-name"
                    :class="{ 'is-link': 可打开菜谱详情(food) }"
                    :disabled="!可打开菜谱详情(food)"
                    @click="打开菜谱详情(food)"
                  >
                    <span class="thumb" :class="{ 'is-link': 可打开菜谱详情(food) }">
                      <RecipeImage
                        class="thumb-image"
                        :src="food.image"
                        :name="food.name"
                        :meal-type="当前详情餐次?.type"
                        alt="食物图"
                      />
                    </span>
                    <span>{{ food.name }}</span>
                  </button>
                  <span class="sub-text">{{ food.weight }}</span>
                  <span class="sub-text">{{ food.kcal }}</span>
                  <span class="sub-text">{{ food.protein }}</span>
                  <span class="sub-text">{{ food.carbs }}</span>
                  <span class="sub-text">{{ food.fat }}</span>
                  <div class="action-wrap">
                    <n-popconfirm
                      v-if="!food.sample && food.realId"
                      @positive-click="删除真实饮食记录(food.realId)"
                      positive-text="删除"
                      negative-text="取消"
                    >
                      <template #trigger>
                        <button class="row-delete-btn" aria-label="删除食物">
                          <n-icon><TrashOutline /></n-icon>
                        </button>
                      </template>
                      确定要删除这道食物记录吗？
                    </n-popconfirm>
                  </div>
                </div>
              </div>

              <div class="add-more">
                <input type="text" placeholder="为这餐添加更多食物..." @click="打开记录悬浮栏(当前详情餐次.type)" readonly />
                <button class="plus-btn" aria-label="添加食物" @click="打开记录悬浮栏(当前详情餐次.type)">+</button>
              </div>
            </div>
          </Transition>
        </div>
      </section>
    </main>

    <n-modal
      v-model:show="显示模式选择弹窗"
      preset="card"
      class="mode-select-modal"
      style="width: 560px"
      title="选择饮食调节模式"
      :mask-closable="!模式保存中"
    >
      <p class="mode-modal-desc">根据当前状态选择模式，系统会按模式含义调节后续推荐。</p>
      <div class="mode-option-list">
        <button
          v-for="item in 饮食模式选项"
          :key="item.mode"
          class="mode-option"
          :class="[`mode-${item.mode}`, { active: item.mode === 待选饮食模式 }]"
          :disabled="模式保存中"
          @click="选择待选模式(item.mode)"
        >
          <div class="mode-option-head">
            <strong>{{ item.label }}</strong>
            <span class="mode-option-days">建议 {{ item.days }} 天</span>
          </div>
          <p>{{ item.summary }}</p>
        </button>
      </div>
      <div class="mode-days-row">
        <span>生效天数</span>
        <n-input-number
          v-model:value="模式生效天数"
          :min="1"
          :max="7"
          :disabled="模式保存中"
          size="small"
        />
      </div>
      <template #footer>
        <div class="mode-modal-footer">
          <span>当前：{{ 当前饮食模式文案 }}</span>
          <div class="mode-modal-actions">
            <n-button secondary :disabled="模式保存中" @click="显示模式选择弹窗 = false">取消</n-button>
            <n-button type="primary" color="#8ec662" :loading="模式保存中" @click="确认保存饮食模式">确认切换</n-button>
          </div>
        </div>
      </template>
    </n-modal>

    <n-modal
      v-model:show="显示单餐重构弹窗"
      preset="card"
      class="mode-select-modal"
      style="width: 680px"
      title="受限单餐重构"
      :mask-closable="!单餐重构中 && !单餐采纳中"
    >
      <div class="regen-grid">
        <label class="regen-label">餐次</label>
        <n-select
          v-model:value="单餐重构表单.meal_type"
          :options="餐次选项"
          placeholder="选择目标餐次"
        />

        <label class="regen-label">食材输入</label>
        <n-input
          v-model:value="单餐重构表单.ingredients_text"
          type="textarea"
          :rows="3"
          placeholder="输入食材，使用逗号分隔，例如：三文鱼, 菠菜, 糙米"
        />

        <div class="regen-upload-row">
          <input type="file" accept="image/png,image/jpeg" @change="选择单餐重构图片" />
          <n-button secondary :loading="食材识别中" @click="识别重构食材">图片识别食材</n-button>
          <n-button secondary @click="解析重构食材文本">加入文本食材</n-button>
        </div>

        <div class="regen-ingredients-list">
          <span
            v-for="item in 单餐重构表单.ingredients"
            :key="item"
            class="regen-chip"
            @click="删除重构食材(item)"
            title="点击移除"
          >
            {{ item }} ×
          </span>
          <span v-if="!单餐重构表单.ingredients.length" class="regen-empty">暂无食材，先输入或识别</span>
        </div>

        <div class="regen-actions">
          <n-button type="primary" color="#8ec662" :loading="单餐重构中" @click="开始单餐重构">开始重构</n-button>
          <n-button type="success" :disabled="!单餐重构结果" :loading="单餐采纳中" @click="采纳单餐重构">采纳该餐次</n-button>
        </div>

        <div v-if="单餐重构中" class="regen-loading-tip">
          <n-spin size="small" />
          <span>正在根据营养缺口和食材重构单餐...</span>
        </div>

        <div v-if="单餐重构结果?.meal" class="regen-result-panel">
          <h4>{{ 单餐重构结果.meal.meal_name }}</h4>
          <p>{{ 单餐重构结果.meal.dietitian_tip }}</p>
          <p v-if="单餐重构结果.supplementary_tip" class="regen-supplementary-tip">{{ 单餐重构结果.supplementary_tip }}</p>
        </div>
      </div>

      <template #footer>
        <div class="mode-modal-footer">
          <span>提示：先锁定今日食谱后再采纳，效果最佳</span>
          <div class="mode-modal-actions">
            <n-button secondary :disabled="单餐重构中 || 单餐采纳中" @click="显示单餐重构弹窗 = false">关闭</n-button>
          </div>
        </div>
      </template>
    </n-modal>

    <!-- 右侧悬浮栏：快速添加饮食记录 -->
    <n-drawer v-model:show="显示添加记录面板" :width="500" placement="right">
      <n-drawer-content title="快速添加饮食记录" closable>
        <template #header>
          <div class="drawer-header-brand">
            <h3 style="margin: 0; color: #235747; font-size: 1.2rem;">添加健康一餐</h3>
          </div>
        </template>
        
        <n-form :model="记录表单" label-placement="top" size="large">
          <n-grid :cols="2" :x-gap="16">
            <n-gi :span="2">
              <n-form-item label="餐点类型">
                <n-select v-model:value="记录表单.meal_type" :options="餐次选项" placeholder="选择餐点" />
              </n-form-item>
            </n-gi>
            <n-gi :span="2">
              <n-form-item label="食物名称 / 描述">
                <n-input-group>
                  <n-input v-model:value="记录表单.food_name" placeholder="例如：一碗牛肉面，少辣" />
                  <n-button type="primary" color="#8ec662" @click="执行AI估算" :loading="分析中" :disabled="!记录表单.food_name">
                    <template #icon><n-icon><SparklesOutline /></n-icon></template>
                    AI估量
                  </n-button>
                  <n-button type="default" @click="triggerAIUpload">
                    <template #icon><n-icon><CameraOutline /></n-icon></template>
                    拍照
                  </n-button>
                  <input type="file" ref="fileInputRef" accept="image/*" style="display:none" @change="处理图片识别" />
                </n-input-group>
              </n-form-item>
            </n-gi>
            <n-gi :span="2">
              <n-form-item label="分量 (g)">
                <n-input-number v-model:value="记录表单.intake_amount" :step="10" placeholder="估算重量" style="width: 100%" />
              </n-form-item>
            </n-gi>
            
            <n-gi :span="2">
              <n-divider dashed style="color:#5f8578; font-size: 0.85rem;">营养成分 (系统可选填/AI自动填)</n-divider>
            </n-gi>

            <n-gi>
              <n-form-item label="热量">
                <n-input-number v-model:value="记录表单.calculated_energy" placeholder="kcal" :show-button="false">
                  <template #suffix>kcal</template>
                </n-input-number>
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="蛋白质">
                <n-input-number v-model:value="记录表单.calculated_protein" placeholder="g" :show-button="false">
                  <template #suffix>g</template>
                </n-input-number>
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="碳水">
                <n-input-number v-model:value="记录表单.calculated_carb" placeholder="g" :show-button="false">
                  <template #suffix>g</template>
                </n-input-number>
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="脂肪">
                <n-input-number v-model:value="记录表单.calculated_fat" placeholder="g" :show-button="false">
                  <template #suffix>g</template>
                </n-input-number>
              </n-form-item>
            </n-gi>
          </n-grid>
        </n-form>

        <template #footer>
          <div style="display: flex; gap: 12px; justify-content: flex-end;">
            <n-button @click="显示添加记录面板 = false">暂不添加</n-button>
            <n-button color="#8ec662" @click="保存饮食记录" :loading="保存中" style="color: white; font-weight: 600;">
              确认记录
            </n-button>
          </div>
        </template>
      </n-drawer-content>
    </n-drawer>
    
    <!-- AI 识别 Loading 遮罩 -->
    <div v-if="识别中" class="ai-loading-mask">
      <div class="ai-loading-inner">
        <n-spin size="large" stroke="#8ec662" />
        <p style="color: #245f4d; font-weight: 600; margin-top: 16px;">正在分析食物营养...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { 
  ref, computed, h, onMounted, onUnmounted, watch, nextTick
} from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NIcon, useMessage, NDrawer, NDrawerContent, NForm, NFormItem, NInput, NInputGroup, NButton, NInputNumber, NSelect, NDivider, NGrid, NGi, NSpin, NPopconfirm, NModal } from 'naive-ui'
import {
  AddCircleOutline,
  CalendarOutline,
  ChevronBackOutline,
  ChevronForwardOutline,
  SparklesOutline,
  MenuOutline,
  CameraOutline,
  FlashOutline,
  TrashOutline
} from '@vicons/ionicons5'
import { getTodayStatus, addIntakeRecord, deleteIntakeRecord } from '@/api/intakeApi'
import { getNutritionRequirements, getDietMode, setDietMode } from '@/api/user'
import {
  getSelectedRecipePlan,
  getRecipeRecommendations,
  selectRecipePlan,
  recognizeMealIngredients,
  regenerateConstrainedMeal,
  adoptRegeneratedMeal
} from '@/api/recipeApi'
import { recognizeFood, analyzeFoodText } from '@/api/foodRecognitionApi'
import { useAuthStore } from '@/store/auth'
import RecipeImage from '@/components/RecipeImage.vue'

const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

// ---- 侧边栏/添加记录相关的状态与逻辑 ----
const 显示添加记录面板 = ref(false)
const 分析中 = ref(false)
const 识别中 = ref(false)
const 保存中 = ref(false)
const fileInputRef = ref(null)
const ai数据缓存 = ref(null)

const 餐次选项 = [
  { label: '🌅 早餐', value: 'breakfast' },
  { label: '☀️ 午餐', value: 'lunch' },
  { label: '🌙 晚餐', value: 'dinner' },
  { label: '🍎 加餐', value: 'snack' }
]

const 记录表单 = ref({
  meal_type: 'breakfast',
  food_name: '',
  intake_amount: 100,
  calculated_energy: 0,
  calculated_protein: 0,
  calculated_carb: 0,
  calculated_fat: 0
})

const 打开记录悬浮栏 = (mealType) => {
  let mappedType = mealType
  if (!['breakfast', 'lunch', 'dinner', 'snack'].includes(mealType)) {
    mappedType = 'breakfast'
  }
  
  记录表单.value = {
    meal_type: mappedType,
    food_name: '',
    intake_amount: 100,
    calculated_energy: 0,
    calculated_protein: 0,
    calculated_carb: 0,
    calculated_fat: 0
  }
  ai数据缓存.value = null
  显示添加记录面板.value = true
}

const 执行AI估算 = async () => {
  if (!记录表单.value.food_name) return
  分析中.value = true
  try {
    const res = await analyzeFoodText(记录表单.value.food_name)
    记录表单.value.food_name = res.dish_name
    记录表单.value.intake_amount = res.estimated_weight
    记录表单.value.calculated_energy = Math.floor(res.total_calories)
    记录表单.value.calculated_protein = Math.floor(res.total_protein)
    记录表单.value.calculated_carb = Math.floor(res.total_carbs)
    记录表单.value.calculated_fat = Math.floor(res.total_fat)

    ai数据缓存.value = {
      calories: res.calories_per_100g,
      protein: res.protein_per_100g,
      carb: res.carbs_per_100g,
      fat: res.fat_per_100g
    }
    message.success('AI文字估算完成')
  } catch (error) {
    message.error('AI估算失败，请手工填写。')
  } finally {
    分析中.value = false
  }
}

const triggerAIUpload = () => {
  fileInputRef.value?.click()
}

const 处理图片识别 = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return
  
  识别中.value = true
  try {
    const res = await recognizeFood(file)
    记录表单.value.food_name = res.dish_name
    记录表单.value.intake_amount = res.estimated_weight
    记录表单.value.calculated_energy = Math.floor(res.total_calories)
    记录表单.value.calculated_protein = Math.floor(res.total_protein)
    记录表单.value.calculated_carb = Math.floor(res.total_carbs)
    记录表单.value.calculated_fat = Math.floor(res.total_fat)

    ai数据缓存.value = {
      calories: res.calories_per_100g,
      protein: res.protein_per_100g,
      carb: res.carbs_per_100g,
      fat: res.fat_per_100g
    }
    
    if (!显示添加记录面板.value) {
      显示添加记录面板.value = true
    }
    message.success('AI视觉分析完成')
  } catch (e) {
    message.error('图片分析失败，请重试')
  } finally {
    识别中.value = false
    if (fileInputRef.value) {
      fileInputRef.value.value = ''
    }
  }
}

const 保存饮食记录 = async () => {
  if (!记录表单.value.food_name) {
    message.warning('请输入食物名称或描述')
    return
  }
  
  保存中.value = true
  try {
    const isTodaySyncing = 是否今天.value // 记录一下是不是今天
    await addIntakeRecord(记录表单.value)
    message.success('成功记录饮食')
    显示添加记录面板.value = false
    
    // 如果记录的是正在查看的那天，重新拉取看板数据刷新
    if (isTodaySyncing) {
      await 同步看板数据()
    }
  } catch (error) {
    message.error('保存饮食记录失败')
  } finally {
    保存中.value = false
  }
}

const 删除真实饮食记录 = async (recordId) => {
  if (!recordId) return
  
  try {
    加载中.value = true
    const isTodaySyncing = 是否今天.value
    await deleteIntakeRecord(recordId)
    message.success('已删除记录')
    
    if (isTodaySyncing) {
      await 同步看板数据()
    }
  } catch (error) {
    console.error('删除饮食记录失败', error)
    message.error('删除失败，请稍后重试')
  } finally {
    加载中.value = false
  }
}

watch(() => 记录表单.value.intake_amount, (newVal) => {
  if (ai数据缓存.value && newVal > 0) {
    const factor = newVal / 100
    记录表单.value.calculated_energy = Math.round(ai数据缓存.value.calories * factor)
    记录表单.value.calculated_protein = Number((ai数据缓存.value.protein * factor).toFixed(1))
    记录表单.value.calculated_carb = Number((ai数据缓存.value.carb * factor).toFixed(1))
    记录表单.value.calculated_fat = Number((ai数据缓存.value.fat * factor).toFixed(1))
  }
})
// ---------------------------------------------

const 加载中 = ref(false)
const 今日方案 = ref(null)
const 当前饮食模式 = ref({ mode: 'normal', source: 'auto', reason: '', until: null })
const 显示模式选择弹窗 = ref(false)
const 模式保存中 = ref(false)
const 显示单餐重构弹窗 = ref(false)
const 单餐重构中 = ref(false)
const 单餐采纳中 = ref(false)
const 食材识别中 = ref(false)
const 单餐重构图片文件 = ref(null)
const 单餐重构结果 = ref(null)
const 单餐重构表单 = ref({
  meal_type: 'lunch',
  ingredients: [],
  ingredients_text: ''
})
const 待选饮食模式 = ref('normal')
const 模式生效天数 = ref(1)
const 饮食模式选项 = [
  { mode: 'normal', label: '标准模式', summary: '按常规目标推荐，营养结构保持平衡。', days: 1, reason: '手动恢复标准' },
  { mode: 'light_adjust', label: '轻调模式', summary: '适度降低油脂与热量，平稳回到目标轨道。', days: 2, reason: '手动轻调' },
  { mode: 'bland', label: '清淡模式', summary: '更清淡、易消化，减少刺激性与高负担食物。', days: 3, reason: '生病/肠胃不适' },
  { mode: 'heavy_adjust', label: '重调模式', summary: '短期强化控制，快速修正昨日明显偏离。', days: 2, reason: '昨日重度偏离' }
]
const 已选日期 = ref(new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate()))
const 编辑时间餐次 = ref('')
const 编辑名称餐次 = ref('')
const 编辑名称内容 = ref('')
const 拖拽源餐次 = ref('')
const 拖拽目标餐次 = ref('')
const 可拖拽餐次 = ref('')
const 每日状态 = ref({
  records: [],
  totalEnergy: 0,
  totalProtein: 0,
  totalCarbohydrate: 0,
  totalFat: 0,
  targetEnergy: 2500,
  targetProtein: 150,
  targetCarbohydrate: 220,
  targetFat: 60
})
const 详情餐次类型 = ref('breakfast')

const 餐次时间本地键 = 'nutriplan_home_meal_times'
const 餐次顺序本地键 = 'nutriplan_home_meal_order'
const 餐次名称本地键 = 'nutriplan_home_meal_names'
const 默认餐次时间 = {
  breakfast: '08:15',
  lunch: '12:30',
  snack: '16:00',
  dinner: '19:30'
}
const 餐次时间映射 = ref({ ...默认餐次时间 })

const 餐次配置 = [
  { type: 'breakfast', name: '早餐', time: '08:15', icon: '🍳', hero: '🍳🍞', hints: ['🍓', '🥣', '🥛'], ratio: 0.16, emoji: '🍳', bg: 'linear-gradient(135deg, #fdf5d3, #f5e094)' },
  { type: 'lunch', name: '午餐', time: '12:30', icon: '🥗', hero: '🥗🍽️', hints: ['🥦', '🍗', '🍚'], ratio: 0.32, emoji: '🥗', bg: 'linear-gradient(135deg, #e5f5e5, #b8e0b8)' },
  { type: 'snack', name: '下午加餐', time: '16:00', icon: '🫐', hero: '🥜🍎', hints: ['🥜', '🍓', '🍵'], ratio: 0.2, emoji: '🍎', bg: 'linear-gradient(135deg, #fdf5d3, #f5e094)' },
  { type: 'dinner', name: '晚餐', time: '19:30', icon: '🥩', hero: '🥩🍷', hints: ['🥔', '🥬', '🍷'], ratio: 0.32, emoji: '🍽️', bg: 'linear-gradient(135deg, #dff0de, #a7d7a7)' }
]

const 默认餐次顺序 = 餐次配置.map((item) => item.type)
const 默认餐次名称 = Object.fromEntries(餐次配置.map((item) => [item.type, item.name]))
const 餐次顺序 = ref([...默认餐次顺序])
const 餐次名称映射 = ref({ ...默认餐次名称 })

const 安全数值 = (value) => {
  const number = Number(value)
  return Number.isFinite(number) ? number : 0
}

const 规范化食物名 = (value) => String(value || '').trim().toLowerCase()

const 提取数字 = (value) => {
  const matched = String(value ?? '').match(/-?\d+(\.\d+)?/)
  return matched ? Number(matched[0]) : 0
}

const 解析菜谱ID = (item, options = {}) => {
  const { allowEntityId = false } = options
  const directId = Number(item?.recipeId || item?.recipe_id || item?.RecipeID)
  if (Number.isFinite(directId) && directId > 0) return directId

  const foodSource = Number(item?.foodSource ?? item?.food_source ?? item?.FoodSource)
  const sourceId = Number(item?.sourceId ?? item?.source_id ?? item?.SourceID)
  if (foodSource === 1 && Number.isFinite(sourceId) && sourceId > 0) return sourceId

  if (allowEntityId) {
    const entityId = Number(item?.id || item?.ID)
    if (Number.isFinite(entityId) && entityId > 0) return entityId
  }

  return null
}

const 可打开菜谱详情 = (food) => !!解析菜谱ID(food)

const 打开菜谱详情 = (food) => {
  const recipeId = 解析菜谱ID(food)
  if (!recipeId) return
  router.push({ name: 'RecipeDetail', params: { id: recipeId } })
}

const 转日期参数 = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const 格式化时间 = (value) => {
  if (!value) return '记录于 --:--'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '记录于 --:--'
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  return `记录于 ${hour}:${minute}`
}

const 归一化餐次 = (value) => {
  if (value === 1 || value === '1' || value === '早餐' || value === 'breakfast') return 'breakfast'
  if (value === 2 || value === '2' || value === '午餐' || value === 'lunch') return 'lunch'
  if (value === 3 || value === '3' || value === '晚餐' || value === 'dinner') return 'dinner'
  if (value === 4 || value === '4' || value === '加餐' || value === '下午加餐' || value === 'snack') return 'snack'
  return ''
}

const 获取餐次记录 = (type) => {
  return (每日状态.value.records || []).filter((item) => 归一化餐次(item.mealType || item.meal_type) === type)
}

const 获取计划餐次食谱列表 = (plan, mealType) => {
  if (!plan || !mealType) return []

  const itemsKey = `${mealType}_items`
  if (Array.isArray(plan[itemsKey]) && plan[itemsKey].length > 0) {
    return plan[itemsKey].filter((item) => item && item.name)
  }

  const singleRecipe = plan[mealType]
  if (singleRecipe && singleRecipe.name) {
    return [singleRecipe]
  }

  return []
}

const 获取餐次卡片预览项 = (config, records = [], planItems = []) => {
  const actualPreview = records.map((item, recordIndex) => {
    const foodName = item.foodName || item.food_name || `${config.name}食物${recordIndex + 1}`
    return {
      name: foodName,
      image: item.imageUrl || item.image_url || '',
      recipeId: 解析菜谱ID(item)
    }
  })

  const fallbackPreview = planItems.length > 0
    ? planItems.map((item, planIndex) => ({
        name: item.name || `${config.name}推荐${planIndex + 1}`,
        image: item.imageUrl || item.image_url || '',
        recipeId: 解析菜谱ID(item, { allowEntityId: true })
      }))
    : (示例食谱映射[config.type] || []).map((item) => ({
        name: item.name,
        image: item.image || '',
        recipeId: item.recipeId || null
      }))

  const existingNames = new Set(
    actualPreview
      .map((item) => 规范化食物名(item.name))
      .filter(Boolean)
  )

  const mergedPreview = [
    ...actualPreview,
    ...fallbackPreview.filter((item) => {
      const normalizedName = 规范化食物名(item.name)
      if (!normalizedName) return true
      return !existingNames.has(normalizedName)
    })
  ].slice(0, 3)

  if (mergedPreview.length > 0) {
    return mergedPreview
  }

  return [{
    name: `${config.name}待添加`,
    image: '',
    emoji: config.icon
  }]
}

const 今日方案菜谱索引 = computed(() => {
  const byId = new Map()
  const byName = new Map()

  if (!今日方案.value) {
    return { byId, byName }
  }

  默认餐次顺序.forEach((mealType) => {
    const planItems = 获取计划餐次食谱列表(今日方案.value, mealType)
    planItems.forEach((recipe) => {
      const recipeId = 解析菜谱ID(recipe, { allowEntityId: true })
      if (!Number.isFinite(recipeId) || recipeId <= 0) return

      const meta = {
        recipeId,
        image: recipe.imageUrl || recipe.image_url || recipe.image || ''
      }

      byId.set(recipeId, meta)

      const normalizedName = 规范化食物名(recipe.name)
      if (normalizedName && !byName.has(normalizedName)) {
        byName.set(normalizedName, meta)
      }
    })
  })

  return { byId, byName }
})

const 加载餐次时间配置 = () => {
  try {
    const raw = localStorage.getItem(餐次时间本地键)
    if (!raw) return
    const parsed = JSON.parse(raw)
    餐次时间映射.value = { ...默认餐次时间, ...parsed }
  } catch (error) {
    console.warn('读取餐次时间配置失败', error)
  }
}

const 保存餐次时间配置 = () => {
  localStorage.setItem(餐次时间本地键, JSON.stringify(餐次时间映射.value))
}

const 加载餐次个性化配置 = () => {
  try {
    const orderRaw = localStorage.getItem(餐次顺序本地键)
    if (orderRaw) {
      const parsedOrder = JSON.parse(orderRaw)
      const validOrder = parsedOrder.filter((type) => 默认餐次顺序.includes(type))
      const missing = 默认餐次顺序.filter((type) => !validOrder.includes(type))
      餐次顺序.value = [...validOrder, ...missing]
    }

    const nameRaw = localStorage.getItem(餐次名称本地键)
    if (nameRaw) {
      const parsedName = JSON.parse(nameRaw)
      餐次名称映射.value = { ...默认餐次名称, ...parsedName }
    }
  } catch (error) {
    console.warn('读取餐次个性化配置失败', error)
  }
}

const 保存餐次顺序配置 = () => {
  localStorage.setItem(餐次顺序本地键, JSON.stringify(餐次顺序.value))
}

const 保存餐次名称配置 = () => {
  localStorage.setItem(餐次名称本地键, JSON.stringify(餐次名称映射.value))
}

const 开始编辑时间 = (mealType) => {
  编辑时间餐次.value = mealType
}

const 更新餐次时间 = (mealType, value) => {
  if (!value) return
  餐次时间映射.value = {
    ...餐次时间映射.value,
    [mealType]: value
  }
}

const 结束编辑时间 = () => {
  if (!编辑时间餐次.value) return
  保存餐次时间配置()
  编辑时间餐次.value = ''
}

const 开始编辑名称 = (mealType, currentName) => {
  编辑名称餐次.value = mealType
  编辑名称内容.value = currentName || 餐次名称映射.value[mealType] || ''
}

const 保存名称编辑 = (mealType) => {
  const nextName = 编辑名称内容.value.trim()
  if (!nextName) {
    编辑名称餐次.value = ''
    编辑名称内容.value = ''
    return
  }
  餐次名称映射.value = {
    ...餐次名称映射.value,
    [mealType]: nextName
  }
  保存餐次名称配置()
  编辑名称餐次.value = ''
  编辑名称内容.value = ''
}

const 开始拖拽餐次 = (mealType) => {
  拖拽源餐次.value = mealType
}

const 允许拖拽 = (type) => {
  可拖拽餐次.value = type
}

const 禁止拖拽 = () => {
  可拖拽餐次.value = ''
}

const 进入拖拽目标 = (mealType) => {
  if (拖拽源餐次.value && 拖拽源餐次.value !== mealType) {
    拖拽目标餐次.value = mealType
  }
}

const 离开拖拽目标 = (mealType) => {
  if (拖拽目标餐次.value === mealType) {
    拖拽目标餐次.value = ''
  }
}

const 放置餐次 = (targetType) => {
  const sourceType = 拖拽源餐次.value
  拖拽目标餐次.value = ''
  if (!sourceType || sourceType === targetType) return

  const nextOrder = [...餐次顺序.value]
  const sourceIndex = nextOrder.indexOf(sourceType)
  const targetIndex = nextOrder.indexOf(targetType)
  if (sourceIndex < 0 || targetIndex < 0) return

  nextOrder.splice(sourceIndex, 1)
  nextOrder.splice(targetIndex, 0, sourceType)
  餐次顺序.value = nextOrder
  保存餐次顺序配置()
}

const 结束拖拽餐次 = () => {
  拖拽源餐次.value = ''
  拖拽目标餐次.value = ''
  可拖拽餐次.value = ''
}

const 从日期字符串创建日期 = (value) => {
  if (!value) return null
  const [year, month, day] = value.split('-').map((item) => Number(item))
  if (!year || !month || !day) return null
  const date = new Date(year, month - 1, day)
  return Number.isNaN(date.getTime()) ? null : date
}

const 今日日期输入值 = computed(() => 转日期参数(new Date()))
const 日期输入值 = computed(() => 转日期参数(已选日期.value))

const 选择日期 = async (event) => {
  const nextDate = 从日期字符串创建日期(event?.target?.value)
  if (!nextDate) return
  已选日期.value = nextDate
  await 同步看板数据()
}

const 日期标题 = computed(() => {
  const now = new Date()
  const isToday = now.toDateString() === 已选日期.value.toDateString()
  const month = 已选日期.value.getMonth() + 1
  const day = 已选日期.value.getDate()
  return isToday ? `今天, ${month}月${day}日` : `${month}月${day}日`
})

const 是否今天 = computed(() => new Date().toDateString() === 已选日期.value.toDateString())

const 今日问候语 = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '🌙 夜深了'
  if (hour < 11) return '👋 早上好'
  if (hour < 14) return '☀️ 中午好'
  if (hour < 18) return '🌤 下午好'
  return '🌆 晚上好'
})

const 今日看板标题 = computed(() => {
  const status = 已摄入热量.value > 0 ? '今日营养计划进行中' : '今日营养计划已就绪'
  return status
})

const 今日日期文案 = computed(() => {
  const date = 已选日期.value
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`
})

const 今日星期文案 = computed(() => {
  const weekdays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  return weekdays[已选日期.value.getDay()] || ''
})

const 时间转分钟 = (timeValue) => {
  if (!timeValue || !timeValue.includes(':')) return 0
  const [hour, minute] = timeValue.split(':').map((item) => Number(item))
  return (hour || 0) * 60 + (minute || 0)
}

const 当前饮食模式文案 = computed(() => {
  const mode = 当前饮食模式.value?.mode || 'normal'
  const map = {
    normal: '标准模式',
    light_adjust: '轻调模式',
    bland: '清淡模式',
    heavy_adjust: '重调模式'
  }
  return map[mode] || '标准模式'
})

const 当前饮食模式说明 = computed(() => {
  const mode = 当前饮食模式.value?.mode || 'normal'
  return 饮食模式选项.find((item) => item.mode === mode)?.summary || '按常规目标推荐，营养结构保持平衡。'
})

const 格式化模式日期 = (value) => {
  if (!value) return '--'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '--'
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const 本地勾选食物 = ref([])
const 勾选同步中食物ID = ref([])

const 查找匹配推荐食谱 = (food) => {
  if (!food) return null

  const recipeId = 解析菜谱ID(food, { allowEntityId: true })
  const normalizedName = 规范化食物名(food.name)

  return (示例详情食物.value || []).find((item) => {
    const sampleRecipeId = 解析菜谱ID(item, { allowEntityId: true })
    if (recipeId && sampleRecipeId && recipeId === sampleRecipeId) {
      return true
    }
    return normalizedName && normalizedName === 规范化食物名(item.name)
  }) || null
}

const 查找已保存推荐食谱记录 = (food) => {
  const matchedSample = food?.sample ? food : 查找匹配推荐食谱(food)
  if (!matchedSample) return null

  const recipeId = 解析菜谱ID(matchedSample, { allowEntityId: true })
  const normalizedName = 规范化食物名(matchedSample.name)

  return 当前详情食物.value.find((item) => {
    const actualRecipeId = 解析菜谱ID(item)
    if (recipeId && actualRecipeId && recipeId === actualRecipeId) {
      return true
    }
    return normalizedName && normalizedName === 规范化食物名(item.name)
  }) || null
}

const 是否可切换食物勾选 = (food) => {
  if (!是否今天.value || !food || 勾选同步中食物ID.value.includes(food.id)) {
    return false
  }
  if (food.sample) {
    return true
  }
  return !!food.realId && !!查找匹配推荐食谱(food)
}

const 是否食物已勾选 = (food) => {
  if (!food) return false
  if (!food.sample) return true
  if (查找已保存推荐食谱记录(food)) return true
  return 本地勾选食物.value.some((item) => item.id === food.id)
}

const 添加本地勾选食物 = (food, mealType) => {
  const parsedKcal = parseFloat(food.kcal) || 0
  const parsedProtein = parseFloat(food.protein) || 0
  const parsedCarbs = parseFloat(food.carbs) || 0
  const parsedFat = parseFloat(food.fat) || 0
  const exists = 本地勾选食物.value.some((item) => item.id === food.id)
  if (exists) return

  本地勾选食物.value.push({
    id: food.id,
    mealType,
    kcal: parsedKcal,
    protein: parsedProtein,
    carbs: parsedCarbs,
    fat: parsedFat
  })
}

const 移除本地勾选食物 = (foodId) => {
  const idx = 本地勾选食物.value.findIndex((item) => item.id === foodId)
  if (idx !== -1) {
    本地勾选食物.value.splice(idx, 1)
  }
}

const 构建推荐食谱记录 = (food, mealType) => {
  const recipeId = 解析菜谱ID(food, { allowEntityId: true })
  return {
    meal_type: mealType,
    food_source: recipeId ? 1 : 2,
    source_id: recipeId || 0,
    food_name: food.name || '推荐食谱',
    intake_amount: Math.max(1, 提取数字(food.weight) || 100),
    calculated_energy: Math.max(0, 提取数字(food.kcal)),
    calculated_protein: Math.max(0, 提取数字(food.protein)),
    calculated_carb: Math.max(0, 提取数字(food.carbs)),
    calculated_fat: Math.max(0, 提取数字(food.fat))
  }
}

const 完成推荐食谱 = async (food, mealType) => {
  添加本地勾选食物(food, mealType)
  勾选同步中食物ID.value.push(food.id)

  try {
    await addIntakeRecord(构建推荐食谱记录(food, mealType))
    message.success(`已完成今日${当前详情餐次.value?.name || '饮食'}计划：${food.name}`)
    await 同步看板数据()
  } catch (error) {
    console.error('保存推荐食谱勾选失败', error)
    移除本地勾选食物(food.id)
    message.error('更新饮食计划状态失败，请稍后重试')
  } finally {
    勾选同步中食物ID.value = 勾选同步中食物ID.value.filter((id) => id !== food.id)
  }
}

const 取消完成推荐食谱 = async (food) => {
  const savedRecord = food?.realId ? food : 查找已保存推荐食谱记录(food)
  if (!savedRecord?.realId) return

  勾选同步中食物ID.value.push(food.id)

  try {
    await deleteIntakeRecord(savedRecord.realId)
    message.success(`已取消今日${当前详情餐次.value?.name || '饮食'}计划：${savedRecord.name || food.name}`)
    await 同步看板数据()
  } catch (error) {
    console.error('取消推荐食谱完成状态失败', error)
    message.error('更新饮食计划状态失败，请稍后重试')
  } finally {
    勾选同步中食物ID.value = 勾选同步中食物ID.value.filter((id) => id !== food.id)
  }
}

const 切换食物状态 = async (food) => {
  if (!是否可切换食物勾选(food)) return
  if (!是否今天.value) {
    message.warning('仅支持在今天将推荐食谱加入饮食记录')
    return
  }

  const mealType = 当前详情餐次.value?.type || 'breakfast'
  const savedRecord = 查找已保存推荐食谱记录(food)
  if (savedRecord) {
    await 取消完成推荐食谱(food)
    return
  }

  await 完成推荐食谱(food, mealType)
}

const 本地勾选摄入总量 = computed(() => {
  return 本地勾选食物.value.reduce((sum, item) => sum + item.kcal, 0)
})
const 本地勾选碳水总量 = computed(() => 本地勾选食物.value.reduce((sum, item) => sum + item.carbs, 0))
const 本地勾选蛋白质总量 = computed(() => 本地勾选食物.value.reduce((sum, item) => sum + item.protein, 0))
const 本地勾选脂肪总量 = computed(() => 本地勾选食物.value.reduce((sum, item) => sum + item.fat, 0))

const 每日目标热量 = computed(() => Math.round(安全数值(每日状态.value.targetEnergy)))
const 已摄入热量 = computed(() => Math.round(安全数值(每日状态.value.totalEnergy)) + 本地勾选摄入总量.value)

const 真实剩余热量 = computed(() => Math.max(0, 每日目标热量.value - 已摄入热量.value))
const 真实进度条 = computed(() => {
  const goal = Math.max(1, 每日目标热量.value)
  return Math.min(100, Math.round((已摄入热量.value / goal) * 100))
})

const 动画剩余热量 = ref(0)
const 动画进度条 = ref(0)
let 动画标识 = null

watch([真实剩余热量, 真实进度条], ([新热量, 新进度]) => {
  const 起始热量 = 动画剩余热量.value
  const 起始进度 = 动画进度条.value
  const 差值热量 = 新热量 - 起始热量
  const 差值进度 = 新进度 - 起始进度
  
  if (差值热量 === 0 && 差值进度 === 0 && 动画剩余热量.value !== 0) return

  const 开始时间 = performance.now()
  const 持续时间 = 600

  if (动画标识) cancelAnimationFrame(动画标识)

  const 执行动画 = (当前时间) => {
    const 经过时间 = 当前时间 - 开始时间
    const 进度 = Math.min(经过时间 / 持续时间, 1)
    const 缓动 = 1 - Math.pow(1 - 进度, 4) // easeOutQuart 缓动函数
    
    动画剩余热量.value = 起始热量 + 差值热量 * 缓动
    动画进度条.value = 起始进度 + 差值进度 * 缓动

    if (进度 < 1) {
      动画标识 = requestAnimationFrame(执行动画)
    } else {
      动画剩余热量.value = 新热量
      动画进度条.value = 新进度
    }
  }
  动画标识 = requestAnimationFrame(执行动画)
}, { immediate: true })

const 剩余热量 = computed(() => Math.round(动画剩余热量.value).toLocaleString('zh-CN'))

const 时间线餐次 = computed(() => {
  const 配置映射 = Object.fromEntries(餐次配置.map((item) => [item.type, item]))
  const now = new Date()
  const nowMinutes = now.getHours() * 60 + now.getMinutes()
  const selectedIsToday = 是否今天.value
  const mapped = 餐次顺序.value.map((type, index, orderedTypes) => {
    const 基础配置 = 配置映射[type]
    if (!基础配置) return null

    const config = {
      ...基础配置,
      name: 餐次名称映射.value[type] || 基础配置.name
    }
    const records = 获取餐次记录(config.type)
    const consumed = Math.round(records.reduce((sum, item) => sum + 安全数值(item.calculatedEnergy || item.calculated_energy), 0))
    const consumedProtein = Math.round(records.reduce((sum, item) => sum + 安全数值(item.calculatedProtein || item.calculated_protein), 0))
    const checkedMealItems = 本地勾选食物.value.filter((item) => item.mealType === config.type)
    const checkedEnergy = Math.round(checkedMealItems.reduce((sum, item) => sum + 安全数值(item.kcal), 0))
    const checkedProtein = Math.round(checkedMealItems.reduce((sum, item) => sum + 安全数值(item.protein), 0))
    const previewEnergy = consumed + checkedEnergy
    const previewProtein = consumedProtein + checkedProtein
    const planItems = (是否今天.value && 今日方案.value)
      ? 获取计划餐次食谱列表(今日方案.value, config.type)
      : []
    const recommendedEnergy = Math.round(planItems.reduce((sum, item) => sum + 安全数值(item.energy), 0))
    const recommendedProtein = Math.round(planItems.reduce((sum, item) => sum + 安全数值(item.protein), 0))
    const displayEnergy = previewEnergy > 0 ? previewEnergy : recommendedEnergy
    const displayProtein = previewProtein > 0 ? previewProtein : recommendedProtein
    const target = Math.max(1, Math.round(每日目标热量.value * config.ratio))
    const progress = previewEnergy > 0
      ? Math.min(100, Math.round((previewEnergy / target) * 100))
      : 0

    let defaultHints = []
    if (是否今天.value && 今日方案.value) {
      if (planItems.length > 0) {
        defaultHints = planItems.map((item) => item.name).filter(Boolean)
      }
    }

    if (defaultHints.length === 0) {
      defaultHints = (示例食谱映射[config.type] || []).map(item => item.name)
    } else {
      defaultHints = defaultHints.slice(0, 4)
    }

    let hints = []
    if (records.length > 0) {
      const addedNames = records.map(r => r.foodName || r.food_name).filter(Boolean)
      hints = Array.from(new Set([...addedNames, ...defaultHints])).slice(0, 4)
    } else {
      hints = defaultHints
    }

    const currentTime = 餐次时间映射.value[config.type] || config.time
    const currentMinutes = 时间转分钟(currentTime)
    const nextType = orderedTypes[index + 1]
    const nextTime = nextType ? (餐次时间映射.value[nextType] || 配置映射[nextType]?.time) : ''
    const nextMinutes = nextType ? 时间转分钟(nextTime) : 24 * 60

    let phase = 'future'
    if (!selectedIsToday) {
      phase = now > 已选日期.value ? 'past' : 'future'
    } else if (nowMinutes >= currentMinutes && nowMinutes < nextMinutes) {
      phase = 'current'
    } else if (nowMinutes >= nextMinutes) {
      phase = 'past'
    }

    const normalizedPreviewItems = 获取餐次卡片预览项(config, records, planItems)

    return {
      ...config,
      time: 餐次时间映射.value[config.type] || config.time,
      records,
      consumed,
      protein: consumedProtein,
      target,
      hints,
      kcal: consumed > 0 ? `${consumed} / ${target} 千卡` : '',
      progress,
      display_energy: Math.max(0, displayEnergy),
      display_protein: Math.max(0, displayProtein),
      thumb: records[0]?.imageUrl || records[0]?.image_url || planItems[0]?.imageUrl || planItems[0]?.image_url || '',
      previewItems: normalizedPreviewItems,
      variant: consumed > 0 ? (config.type === 'lunch' ? 'logged' : 'filled') : 'empty',
      phase
    }
  }).filter(Boolean)

  return mapped
})

const 时间线样式 = computed(() => {
  const meals = 时间线餐次.value || []
  if (!meals.length) {
    return { '--timeline-progress': '0%' }
  }
  const currentIndex = meals.findIndex((item) => item.phase === 'current')
  const pastCount = meals.filter((item) => item.phase === 'past').length
  const done = currentIndex >= 0 ? currentIndex + 0.5 : pastCount
  const progress = Math.min(100, Math.max(0, Math.round((done / meals.length) * 100)))
  return { '--timeline-progress': `${progress}%` }
})

const 宏量卡片 = computed(() => {
  const totalCarbs = Math.round(安全数值(每日状态.value.totalCarbohydrate) + 本地勾选碳水总量.value)
  const totalProtein = Math.round(安全数值(每日状态.value.totalProtein) + 本地勾选蛋白质总量.value)
  const totalFat = Math.round(安全数值(每日状态.value.totalFat) + 本地勾选脂肪总量.value)

  const targetCarbs = Math.round(安全数值(每日状态.value.targetCarbohydrate))
  const targetProtein = Math.round(安全数值(每日状态.value.targetProtein))
  const targetFat = Math.round(安全数值(每日状态.value.targetFat))

  return [
    { label: '碳水', value: `${totalCarbs}g`, target: `${targetCarbs}g`, icon: '🌾', soft: 'rgba(255, 186, 108, 0.22)' },
    { label: '蛋白质', value: `${totalProtein}g`, target: `${targetProtein}g`, icon: '💧', soft: 'rgba(116, 163, 255, 0.2)' },
    { label: '脂肪', value: `${totalFat}g`, target: `${targetFat}g`, icon: '🟡', soft: 'rgba(255, 215, 83, 0.22)' }
  ]
})

const 当前详情餐次 = computed(() => {
  const current = 时间线餐次.value.find((item) => item.type === 详情餐次类型.value)
  if (current) return current
  return 时间线餐次.value.find((item) => item.records.length > 0) || 时间线餐次.value[1] || 时间线餐次.value[0]
})

const 当前详情标题 = computed(() => `${当前详情餐次.value?.name || '餐次'}详情`)

const 当前详情时间 = computed(() => {
  const record = 当前详情餐次.value?.records?.[0]
  return 格式化时间(record?.createdAt || record?.CreatedAt)
})

const 详情营养条 = computed(() => {
  const records = 当前详情餐次.value?.records || []
  const ratio = 当前详情餐次.value?.ratio || 0.25
  let protein = records.reduce((sum, item) => sum + 安全数值(item.calculatedProtein || item.calculated_protein), 0)
  let carbs = records.reduce((sum, item) => sum + 安全数值(item.calculatedCarb || item.calculated_carb), 0)
  let fat = records.reduce((sum, item) => sum + 安全数值(item.calculatedFat || item.calculated_fat), 0)

  if (展示详情食物.value) {
    const checkedSample = 展示详情食物.value.filter(food => food.sample && 本地勾选食物.value.some(f => f.id === food.id))
    protein += checkedSample.reduce((sum, item) => sum + (parseFloat(item.protein) || 0), 0)
    carbs += checkedSample.reduce((sum, item) => sum + (parseFloat(item.carbs) || 0), 0)
    fat += checkedSample.reduce((sum, item) => sum + (parseFloat(item.fat) || 0), 0)
  }

  const targetProtein = Math.max(1, 安全数值(每日状态.value.targetProtein) * ratio)
  const targetCarbs = Math.max(1, 安全数值(每日状态.value.targetCarbohydrate) * ratio)
  const targetFat = Math.max(1, 安全数值(每日状态.value.targetFat) * ratio)

  return [
    {
      label: '蛋白质',
      value: `${Math.round(protein)}g`,
      progress: Math.min(100, Math.round((protein / targetProtein) * 100)),
      color: '#4e8fff',
      spark: [22, 38, 44, 60, Math.min(100, Math.round((protein / targetProtein) * 100))]
    },
    {
      label: '碳水',
      value: `${Math.round(carbs)}g`,
      progress: Math.min(100, Math.round((carbs / targetCarbs) * 100)),
      color: '#ff8d2f',
      spark: [18, 28, 46, 54, Math.min(100, Math.round((carbs / targetCarbs) * 100))]
    },
    {
      label: '脂肪',
      value: `${Math.round(fat)}g`,
      progress: Math.min(100, Math.round((fat / targetFat) * 100)),
      color: '#f5c84b',
      spark: [24, 30, 36, 48, Math.min(100, Math.round((fat / targetFat) * 100))]
    }
  ]
})

const 当前详情食物 = computed(() => {
  const meal = 当前详情餐次.value
  const fallback = { emoji: meal?.emoji || '🍽️', bg: meal?.bg || 'linear-gradient(135deg, #7aa694, #4f7a69)' }
  const { byId, byName } = 今日方案菜谱索引.value
  return (meal?.records || []).map((item, index) => {
    const foodName = item.foodName || item.food_name || '未知食物'
    const explicitRecipeId = 解析菜谱ID(item)
    const matchedMeta = (explicitRecipeId && byId.get(explicitRecipeId)) || byName.get(规范化食物名(foodName))
    return {
      id: item.id || item.ID || `${meal.type}-${index}`,
      realId: item.id || item.ID, // 提取真实 ID
      name: foodName,
      recipeId: explicitRecipeId || matchedMeta?.recipeId || null,
      weight: `${Math.round(安全数值(item.intakeAmount || item.intake_amount))}g`,
      kcal: `${Math.round(安全数值(item.calculatedEnergy || item.calculated_energy))} 千卡`,
      protein: `${Math.round(安全数值(item.calculatedProtein || item.calculated_protein))}g`,
      carbs: `${Math.round(安全数值(item.calculatedCarb || item.calculated_carb))}g`,
      fat: `${Math.round(安全数值(item.calculatedFat || item.calculated_fat))}g`,
      emoji: fallback.emoji,
      bg: fallback.bg,
      image: item.imageUrl || item.image_url || matchedMeta?.image || '',
      sample: false
    }
  })
})

const 示例食谱映射 = {
  breakfast: [
    { name: '牛油果全麦吐司', weight: '130g', kcal: '285 千卡', protein: '8g', carbs: '32g', fat: '15g', emoji: '🥑' },
    { name: '蓝莓酸奶碗', weight: '180g', kcal: '210 千卡', protein: '12g', carbs: '24g', fat: '7g', emoji: '🫐' }
  ],
  lunch: [
    { name: '鸡胸肉藜麦沙拉', weight: '260g', kcal: '420 千卡', protein: '42g', carbs: '38g', fat: '12g', emoji: '🥗' },
    { name: '南瓜浓汤', weight: '220g', kcal: '165 千卡', protein: '4g', carbs: '28g', fat: '5g', emoji: '🎃' }
  ],
  snack: [
    { name: '混合坚果', weight: '35g', kcal: '195 千卡', protein: '7g', carbs: '9g', fat: '16g', emoji: '🥜' },
    { name: '苹果切片', weight: '120g', kcal: '62 千卡', protein: '0g', carbs: '14g', fat: '0g', emoji: '🍎' }
  ],
  dinner: [
    { name: '香煎三文鱼配芦笋', weight: '240g', kcal: '468 千卡', protein: '38g', carbs: '8g', fat: '28g', emoji: '🐟' },
    { name: '黑椒菌菇意面', weight: '210g', kcal: '338 千卡', protein: '14g', carbs: '52g', fat: '8g', emoji: '🍝' }
  ]
}

const 示例详情食物 = computed(() => {
  const meal = 当前详情餐次.value
  let list = []

  if (是否今天.value && 今日方案.value) {
    const planItems = 获取计划餐次食谱列表(今日方案.value, meal?.type)
    if (planItems.length > 0) {
      list = planItems.map((recipe) => ({
        name: recipe.name,
        recipeId: 解析菜谱ID(recipe, { allowEntityId: true }),
        weight: `${Math.round(Number(recipe.portion_weight_g || recipe.portionWeightG || 100))}g`,
        kcal: `${Math.round(recipe.energy || 0)} 千卡`,
        protein: `${Math.round(recipe.protein || 0)}g`,
        carbs: `${Math.round(recipe.carbohydrate || 0)}g`,
        fat: `${Math.round(recipe.fat || 0)}g`,
        emoji: meal?.emoji || '🍽️',
        image: recipe.imageUrl || recipe.image_url || ''
      }))
    }
  }

  if (!list.length) {
    list = 示例食谱映射[meal?.type] || []
  }

  return list.map((item, index) => ({
    id: `sample-${meal?.type || 'meal'}-${index}`,
    name: item.name,
    recipeId: item.recipeId,
    weight: item.weight || item.portion || '100g',
    kcal: item.kcal,
    protein: item.protein,
    carbs: item.carbs,
    fat: item.fat,
    emoji: item.emoji || meal?.emoji || '🍽️',
    bg: meal?.bg || 'linear-gradient(135deg, #7aa694, #4f7a69)',
    image: item.image || '',
    sample: true
  }))
})

const 展示详情食物 = computed(() => {
  const actual = 当前详情食物.value || []
  const samples = 示例详情食物.value || []
  
  // 保留原有的示例食谱样貌，每次添加的食物追加/叠加在前面，保持 UI 不会空洞
  const actualNames = new Set(actual.map((item) => 规范化食物名(item.name)).filter(Boolean))
  const filteredSamples = samples.filter((item) => {
    const normalizedName = 规范化食物名(item.name)
    if (!normalizedName) return true
    return !actualNames.has(normalizedName)
  })

  return [...actual, ...filteredSamples]
})

const ringStyle = computed(() => {
  const progress = 动画进度条.value
  return {
    background: `conic-gradient(#8ec662 0 ${progress}%, #e7f2da ${progress}% 100%)`
  }
})

const 同步看板数据 = async () => {
  加载中.value = true
  try {
    if (authStore.isAuthenticated && !authStore.profile) {
      await authStore.loadProfile()
    }

    const promises = [
      getTodayStatus(转日期参数(已选日期.value)),
      getNutritionRequirements().catch(() => null),
      getDietMode().catch(() => null)
    ]
    if (是否今天.value) {
      promises.push(getSelectedRecipePlan().catch(() => null))
    }

    const res = await Promise.all(promises)
    const status = res[0]
    const nutrition = res[1]
    const modeResp = res[2]
    const plan = res[3]

    当前饮食模式.value = {
      mode: modeResp?.mode || 'normal',
      source: modeResp?.source || 'auto',
      reason: modeResp?.reason || '',
      until: modeResp?.until || null
    }

    if (是否今天.value && plan && plan.id) {
      今日方案.value = plan
    } else {
      今日方案.value = null
    }

    每日状态.value = {
      records: status?.records || [],
      totalEnergy: 安全数值(status?.total_energy ?? status?.totalEnergy),
      totalProtein: 安全数值(status?.total_protein ?? status?.totalProtein),
      totalCarbohydrate: 安全数值(status?.total_carbohydrate ?? status?.totalCarbohydrate),
      totalFat: 安全数值(status?.total_fat ?? status?.totalFat),
      targetEnergy: 安全数值(status?.target_energy ?? status?.targetEnergy ?? nutrition?.target_calorie ?? nutrition?.targetCalorie) || 2500,
      targetProtein: 安全数值(status?.target_protein ?? status?.targetProtein ?? nutrition?.target_protein ?? nutrition?.targetProtein) || 150,
      targetCarbohydrate: 安全数值(status?.target_carbohydrate ?? status?.targetCarbohydrate ?? nutrition?.target_carbohydrate ?? nutrition?.targetCarbohydrate) || 220,
      targetFat: 安全数值(status?.target_fat ?? status?.targetFat ?? nutrition?.target_fat ?? nutrition?.targetFat) || 60
    }

    const hasCurrent = 时间线餐次.value.some((item) => item.type === 详情餐次类型.value)
    if (!hasCurrent) {
      详情餐次类型.value = 'breakfast'
    }
    本地勾选食物.value = []
  } catch (error) {
    console.error('加载首页看板数据失败', error)
    每日状态.value = {
      ...每日状态.value,
      records: []
    }
  } finally {
    加载中.value = false
  }
}

const 上一天 = async () => {
  已选日期.value = new Date(已选日期.value.getFullYear(), 已选日期.value.getMonth(), 已选日期.value.getDate() - 1)
  await 同步看板数据()
}

const 下一天 = async () => {
  if (是否今天.value) return
  已选日期.value = new Date(已选日期.value.getFullYear(), 已选日期.value.getMonth(), 已选日期.value.getDate() + 1)
  await 同步看板数据()
}

const 切换详情餐次 = (mealType) => {
  if (详情餐次类型.value === mealType) return
  详情餐次类型.value = mealType
}

const 去模式设置 = () => {
  待选饮食模式.value = 当前饮食模式.value?.mode || 'normal'
  const defaultDays = 饮食模式选项.find((item) => item.mode === 待选饮食模式.value)?.days || 1
  模式生效天数.value = Math.max(1, Math.min(7, Number(defaultDays) || 1))
  显示模式选择弹窗.value = true
}

const 选择待选模式 = (mode) => {
  const target = 饮食模式选项.find((item) => item.mode === mode)
  if (!target) return
  待选饮食模式.value = target.mode
  if (!模式生效天数.value || 模式生效天数.value < 1) {
    模式生效天数.value = target.days
  }
}

const 确认保存饮食模式 = async () => {
  const mode = 待选饮食模式.value || 'normal'
  const target = 饮食模式选项.find((item) => item.mode === mode)
  if (!target) return

  const days = Math.max(1, Math.min(7, Number(模式生效天数.value) || target.days || 1))
  模式保存中.value = true
  try {
    const resp = await setDietMode({ mode: target.mode, days, reason: target.reason })
    当前饮食模式.value = {
      mode: resp?.mode || target.mode,
      source: resp?.source || 'manual',
      reason: resp?.reason || target.reason,
      until: resp?.until || null
    }
    message.success(`已切换为${当前饮食模式文案.value}`)
    显示模式选择弹窗.value = false
  } catch (error) {
    message.error(error?.response?.data?.error || '饮食模式更新失败')
  } finally {
    模式保存中.value = false
  }
}

const 去生成食谱 = async () => {
  if (!是否今天.value) return
  加载中.value = true
  try {
    const res = await getRecipeRecommendations(1)
    if (res && res.plans && res.plans.length > 0) {
      await selectRecipePlan(res.plans[0])
      message.success('已为您生成新的今日食谱')
      await 同步看板数据()
    } else {
      message.warning('暂无合适的食谱推荐')
    }
  } catch (error) {
    console.error('智能生成食谱失败', error)
    message.error('生成食谱失败，请稍后重试')
  } finally {
    加载中.value = false
  }
}

const 去受限单餐重构 = () => {
  单餐重构结果.value = null
  单餐重构图片文件.value = null
  单餐重构表单.value = {
    meal_type: 'lunch',
    ingredients: [],
    ingredients_text: ''
  }
  显示单餐重构弹窗.value = true
}

const 解析重构食材文本 = () => {
  const parsed = (单餐重构表单.value.ingredients_text || '')
    .split(/[，,\n]/)
    .map((item) => item.trim())
    .filter(Boolean)

  单餐重构表单.value.ingredients = Array.from(new Set([
    ...单餐重构表单.value.ingredients,
    ...parsed
  ]))
}

const 删除重构食材 = (item) => {
  单餐重构表单.value.ingredients = 单餐重构表单.value.ingredients.filter((value) => value !== item)
}

const 选择单餐重构图片 = (event) => {
  单餐重构图片文件.value = event.target.files?.[0] || null
}

const 识别重构食材 = async () => {
  if (!单餐重构图片文件.value) {
    message.warning('请先选择食材图片')
    return
  }

  食材识别中.value = true
  try {
    const resp = await recognizeMealIngredients(单餐重构图片文件.value)
    const recognized = Array.isArray(resp?.ingredients) ? resp.ingredients : []
    单餐重构表单.value.ingredients = Array.from(new Set([
      ...单餐重构表单.value.ingredients,
      ...recognized
    ]))
    message.success('识别完成，请确认食材列表')
  } catch (error) {
    message.error(error?.response?.data?.error || '食材识别失败，请稍后重试')
  } finally {
    食材识别中.value = false
  }
}

const 开始单餐重构 = async () => {
  解析重构食材文本()
  if (!单餐重构表单.value.meal_type) {
    message.warning('请选择餐次')
    return
  }
  if (!单餐重构表单.value.ingredients.length) {
    message.warning('请至少提供一种食材')
    return
  }

  单餐重构中.value = true
  try {
    const payload = {
      meal_type: 单餐重构表单.value.meal_type,
      ingredients: 单餐重构表单.value.ingredients
    }
    单餐重构结果.value = await regenerateConstrainedMeal(payload)
    message.success('单餐重构完成，可直接采纳')
  } catch (error) {
    message.error(error?.response?.data?.error || '单餐重构失败，请稍后再试')
  } finally {
    单餐重构中.value = false
  }
}

const 采纳单餐重构 = async () => {
  if (!单餐重构结果.value?.meal) return

  单餐采纳中.value = true
  try {
    const resp = await adoptRegeneratedMeal({
      meal_type: 单餐重构表单.value.meal_type,
      meal: 单餐重构结果.value.meal
    })

    if (resp?.plan) {
      今日方案.value = resp.plan
    }
    await 同步看板数据()
    显示单餐重构弹窗.value = false
    message.success('已采纳重构结果，当前页数据已刷新')
  } catch (error) {
    message.error(error?.response?.data?.error || '采纳失败，请先锁定今日食谱')
  } finally {
    单餐采纳中.value = false
  }
}

onMounted(async () => {
  加载餐次时间配置()
  加载餐次个性化配置()
  await 同步看板数据()
})
</script>

<style scoped>
:global(body) {
  font-family: 'Inter', 'Segoe UI', Roboto, Arial, sans-serif;
  background: #EFF4F1;
}

.nutrition-dashboard {
  min-height: 100vh;
  padding: 24px;
  position: relative;
  overflow: hidden;
  background-color: #EFF4F1;
  background-image: 
    radial-gradient(circle at 12% 8%, rgba(255, 255, 255, 0.8) 0%, transparent 40%),
    radial-gradient(circle at 82% 88%, rgba(236, 245, 239, 0.75) 0%, transparent 42%);
  color: #1a2f24;
}

.glass-card {
  background: #FAFCFB;
  border: 1px solid rgba(231, 238, 234, 0.9);
  border-radius: 18px;
  box-shadow: 0px 8px 24px rgba(0, 0, 0, 0.04);
  transition: all 0.2s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.glass-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 24px rgba(35, 73, 57, 0.07);
}

.top-header,
.center-panel,
.timeline-content,
.overview-card,
.details-card,
.food-table {
  border-radius: 16px;
}

.top-header {
  min-height: 78px;
  padding: 14px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
  position: relative;
  z-index: 2;
}

.header-main {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.header-greet {
  margin: 0;
  font-size: 0.92rem;
  color: #6a8a7a;
  font-weight: 700;
}

.header-title {
  margin: 0;
  font-size: 1.35rem;
  color: #184c3d;
  font-weight: 800;
}

.header-sub {
  margin: 0;
  font-size: 0.82rem;
  color: #84a195;
  font-weight: 700;
}

.header-meta {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  background: #f4faf6;
  border: 1px solid #e1ece6;
  border-radius: 14px;
  padding: 8px 12px;
}

.meta-icon {
  width: 30px;
  height: 30px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  background: #ffffff;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.05);
}

.meta-text {
  display: flex;
  flex-direction: column;
}

.meta-text strong {
  color: #2b5d4c;
  font-size: 0.86rem;
}

.meta-text span {
  color: #7b988c;
  font-size: 0.78rem;
}

.top-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
}

.icon-circle {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 50%;
  background: linear-gradient(145deg, #ffffff 0%, #f2f7f2 100%);
  color: #2b5643; /* 深叶绿 */
  display: grid;
  place-items: center;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.2, 0.8, 0.2, 1);
  box-shadow: 0 4px 14px rgba(32, 59, 45, 0.08); /* 柔和植物影 */
}

.icon-circle:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 22px rgba(32, 59, 45, 0.12); /* 偏光效果 */
}

.icon-circle:active {
  transform: translateY(1px) scale(0.96); /* 下沉动效 */
  box-shadow: 0 2px 8px rgba(32, 59, 45, 0.06);
}

.icon-circle:disabled {
  opacity: 0.45;
  cursor: not-allowed;
  transform: none;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: minmax(540px, 1.1fr) minmax(360px, 0.9fr);
  align-items: start;
  gap: 24px;
  position: relative;
  z-index: 1;
}

.center-panel {
  background: transparent;
  border: none;
  padding: 0;
  backdrop-filter: none;
}

.loading-tip {
  margin: 0 0 12px;
  color: #5d8778;
  font-size: 1rem;
  font-weight: 600;
}

.record-ops {
  margin-bottom: 24px;
  padding: 12px 18px;
  border: 1px solid #e9f0eb;
  border-radius: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.ops-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ops-date {
  margin: 0 0 0 6px;
  color: #2b5745;
  font-size: 1rem;
  font-weight: 700;
}

.ops-subtitle {
  margin: 0 0 0 6px;
  color: #7b9f8f;
  font-size: 0.95rem;
}

.calendar-picker {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #2a6552;
  font-size: 1rem;
  font-weight: 700;
}

.calendar-picker input {
  border: 1px solid rgba(45, 112, 90, 0.2);
  background: #f5fbf8;
  color: #2f6554;
  height: 34px;
  padding: 0 8px;
}

.recipe-action-bar {
  margin-bottom: 24px;
  padding: 20px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-left: 6px solid #8ec662;
}

.action-info .plan-status {
  margin: 0 0 4px;
  font-size: 1.1rem;
  color: #2c3e35;
  font-weight: 700;
}

.action-info .plan-desc {
  margin: 0;
  font-size: 0.85rem;
  color: #6c8a7b;
}

.diet-mode-chip {
  margin-top: 10px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  border-radius: 999px;
  padding: 6px 12px;
  border: 1px solid #d5e7db;
  background: #f6fbf7;
  cursor: pointer;
  transition: all 0.2s ease;
}

.diet-mode-chip:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
}

.diet-mode-chip .mode-title {
  font-size: 0.88rem;
  font-weight: 700;
  color: #111111;
}

.diet-mode-chip .mode-meta {
  font-size: 0.8rem;
  color: #111111;
  opacity: 0.82;
}

.diet-mode-chip.mode-normal {
  background: #eef7ff;
  border-color: #cfe4ff;
}

.diet-mode-chip.mode-light_adjust {
  background: #fff8ea;
  border-color: #ffe1a6;
}

.diet-mode-chip.mode-bland {
  background: #edf9ef;
  border-color: #bfe5c4;
}

.diet-mode-chip.mode-heavy_adjust {
  background: #fff1f1;
  border-color: #ffc8c8;
}

.mode-select-modal :deep(.n-card-header__main) {
  color: #184c3d;
  font-weight: 700;
}

.mode-modal-desc {
  margin: 0 0 12px;
  font-size: 14px;
  color: #2f5d4f;
}

.mode-option-list {
  display: grid;
  gap: 10px;
}

.mode-option {
  text-align: left;
  border-radius: 14px;
  border: 1px solid #d8e8dd;
  background: #f8fbf9;
  padding: 12px;
  cursor: pointer;
  transition: all 0.18s ease;
}

.mode-option:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 20px rgba(32, 59, 45, 0.08);
}

.mode-option:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.mode-option.active {
  border-color: #86bd5a;
  box-shadow: 0 0 0 2px rgba(141, 199, 93, 0.2);
}

.mode-option-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 4px;
}

.mode-option-head strong {
  color: #184c3d;
  font-size: 15px;
}

.mode-option-days {
  color: #2f5d4f;
  font-size: 12px;
}

.mode-option p {
  margin: 0;
  color: #111111;
  font-size: 13px;
  line-height: 1.5;
}

.mode-option.mode-normal {
  background: #eef7ff;
}

.mode-option.mode-light_adjust {
  background: #fff8ea;
}

.mode-option.mode-bland {
  background: #edf9ef;
}

.mode-option.mode-heavy_adjust {
  background: #fff1f1;
}

.mode-modal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  color: #2f5d4f;
}

.mode-modal-actions {
  display: flex;
  gap: 10px;
}

.mode-days-row {
  margin-top: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #184c3d;
  font-size: 14px;
}

.regen-grid {
  display: grid;
  gap: 10px;
}

.regen-label {
  font-size: 13px;
  color: #2f5d4f;
  font-weight: 600;
}

.regen-upload-row {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.regen-ingredients-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  min-height: 32px;
}

.regen-chip {
  background: #edf7e8;
  color: #2c5b3f;
  border: 1px solid #cfe6c4;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
  cursor: pointer;
}

.regen-empty {
  color: #8aa497;
  font-size: 12px;
}

.regen-actions {
  display: flex;
  gap: 10px;
  margin-top: 4px;
}

.regen-loading-tip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #2f5d4f;
  font-size: 13px;
}

.regen-result-panel {
  margin-top: 4px;
  padding: 10px 12px;
  border-radius: 12px;
  background: #f6fbf5;
  border: 1px solid #dcead7;
}

.regen-result-panel h4 {
  margin: 0 0 6px;
  color: #184c3d;
}

.regen-result-panel p {
  margin: 0;
  color: #2f5d4f;
  font-size: 13px;
}

.regen-supplementary-tip {
  margin-top: 6px !important;
  color: #2b8a3e !important;
  font-weight: 600;
}

.commercial-btn {
  border: none;
  /* 亲自然：叶绿自然渐变 */
  background: linear-gradient(135deg, #8dc75d 0%, #72aa46 100%);
  color: #ffffff;
  padding: 10px 20px;
  border-radius: 20px; /* 有机圆角 */
  font-size: 1rem;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  box-shadow: 0 6px 16px rgba(122, 180, 77, 0.25);
  transition: all 0.18s cubic-bezier(0.2, 0.8, 0.2, 1); /* 在 150-220ms 范围内 */
}

.recipe-action-buttons {
  display: flex;
  gap: 10px;
}

.secondary-btn {
  background: #ffffff;
  color: #5a7f3b;
  border: 1px solid #9fc07f;
  box-shadow: none;
}

.commercial-btn:hover {
  background: linear-gradient(135deg, #95cc68 0%, #79b04d 100%); /* 亮度微升 */
  transform: translateY(-2px); /* 悬浮感 */
  box-shadow: 0 10px 24px rgba(122, 180, 77, 0.4); 
}

.commercial-btn:active {
  transform: translateY(1px) scale(0.97); /* Active 轻微下沉 */
  box-shadow: 0 2px 8px rgba(122, 180, 77, 0.15); /* 阴影收拢 */
}

.timeline {
  padding-left: 22px;
  border-left: none;
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 3px;
  top: 0;
  bottom: 0;
  width: 4px;
  border-radius: 999px;
  background: #d6dfd8;
}

.timeline::after {
  content: '';
  position: absolute;
  left: 3px;
  top: 0;
  width: 4px;
  height: var(--timeline-progress, 0%);
  border-radius: 999px;
  background: linear-gradient(180deg, #8ec662, #78b456);
  transition: height 0.35s ease;
}

.timeline-row {
  display: grid;
  grid-template-columns: 98px minmax(280px, 1fr);
  gap: 12px;
  position: relative;
  transition: transform 0.22s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.timeline-row:hover {
  transform: translateX(4px) scale(1.005);
}

.timeline-row.is-drag-over .timeline-content {
  border: 2px dashed #8dc75d;
  background: #f4faf0;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(122, 180, 77, 0.15);
}

.timeline-row.is-active-detail .timeline-content {
  box-shadow: 0 0 0 2px rgba(132, 189, 87, 0.4); /* 植物环绕感 */
}

.time-col {
  padding-top: 10px;
  position: relative;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #ffffff;
  position: absolute;
  left: -24px;
  top: 14px;
  border: 2px solid #8ec662;
  box-shadow: 0 0 0 5px #EFF4F1;
}

.timeline-row.timeline-past .dot {
  border-color: #a9b7b0;
  background: #e9efec;
}

.timeline-row.timeline-current .dot {
  border-color: #74b848;
  box-shadow: 0 0 0 5px #ecf6e7;
}

.timeline-row.timeline-future .dot {
  border-color: #d0dbd5;
  background: #ffffff;
}

.time-editor {
  min-height: 24px;
  display: flex;
  align-items: center;
}

.time-display-btn {
  border: none;
  background: transparent;
  color: #5f8a7b;
  font-size: 0.92rem;
  font-weight: 700;
  padding: 0;
  cursor: pointer;
}

.time-input {
  width: 92px;
  border: 1px solid rgba(47, 124, 99, 0.25);
  border-radius: 8px;
  background: #f4faf6;
  color: #3f6f60;
  padding: 3px 6px;
  font-size: 0.78rem;
}

.timeline-content {
  padding: 12px 14px;
  border-radius: 14px;
  position: relative;
  overflow: hidden;
  cursor: pointer;
}

.timeline-content::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: radial-gradient(circle at 16% 16%, rgba(142, 198, 98, 0.08), transparent 32%),
    radial-gradient(circle at 88% 86%, rgba(114, 170, 70, 0.06), transparent 30%);
  pointer-events: none;
}

.timeline-row.timeline-past .timeline-content {
  border-left-color: #bac8c0;
}

.timeline-row.timeline-current .timeline-content {
  box-shadow: 0 0 0 2px rgba(134, 193, 89, 0.25), 0 10px 24px rgba(73, 117, 80, 0.08);
}

.meal-card::before,
.meal-card::after {
  display: none;
}

.meal-breakfast {
  background: #FAFCFB;
  border-left: 6px solid #f5e094;
}

.meal-lunch {
  background: #FAFCFB;
  border-left: 6px solid #8ec662;
}

.meal-snack {
  background: #FAFCFB;
  border-left: 6px solid #f5e094;
}

.meal-dinner {
  background: #FAFCFB;
  border-left: 6px solid #8ec662;
}

.meal-shell {
  display: grid;
  grid-template-columns: minmax(0, 7fr) minmax(120px, 3fr);
  gap: 10px;
  align-items: stretch;
  position: relative;
  z-index: 1;
}

.meal-main {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.meal-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.meal-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.meal-head h3 {
  margin: 0;
  font-size: 1.24rem;
  color: #204c3f;
}

.meal-time-inline {
  font-size: 0.9rem;
  color: #6b8d7e;
  font-weight: 700;
}

.editable-meal-name {
  cursor: pointer;
  transition: color 0.2s ease;
}

.meal-name-input {
  width: 150px;
  border: 1px solid rgba(47, 124, 99, 0.28);
  border-radius: 8px;
  background: #f4faf6;
  color: #2c6251;
  padding: 5px 8px;
  font-size: 1.1rem;
  font-weight: 700;
}

.drag-handle {
  cursor: grab;
  font-size: 1rem;
  color: #a0bcae;
  display: flex;
  padding: 2px;
  border-radius: 4px;
  transition: background 0.2s, color 0.2s;
}

.drag-handle:hover {
  background: #f0f5f2;
  color: #5f8a7b;
}

.drag-handle:active {
  cursor: grabbing;
}

.meal-food-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(118px, 1fr));
  gap: 16px;
  min-width: 0;
  width: 100%;
}

.food-preview-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  gap: 8px;
  min-width: 0;
}

.food-image-card {
  width: 100%;
  max-width: 168px;
  height: 112px;
  border-radius: 20px;
  border: 1px solid rgba(214, 230, 218, 0.88);
  background: #f8fbf8;
  box-shadow: 0 10px 18px rgba(58, 90, 72, 0.1);
  overflow: hidden;
  display: grid;
  place-items: center;
  padding: 6px;
  box-sizing: border-box;
  color: #446f5f;
}

.food-image-card.secondary {
  opacity: 0.96;
}

.food-image-card.is-link {
  cursor: pointer;
}

.food-image-card.is-link:hover {
  transform: translateY(-1px);
  box-shadow: 0 14px 24px rgba(58, 90, 72, 0.14);
}

.food-image-card img {
  width: auto;
  height: auto;
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  object-position: center;
  background: #f8fbf8;
}

.food-preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.food-preview-image.recipe-image--placeholder {
  object-fit: contain;
  object-position: center;
  padding: 6px;
  box-sizing: border-box;
  background: #f8fbf8;
}

.food-image-fallback {
  font-size: 1.7rem;
}

.food-preview-name {
  width: 100%;
  margin: 0;
  color: #204c3f;
  font-size: 0.92rem;
  font-weight: 700;
  line-height: 1.35;
  text-align: center;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.food-preview-name.is-link {
  cursor: pointer;
}

.food-preview-name.is-link:hover {
  color: #2f7a59;
}

.meal-side {
  border-left: 1px dashed #d8e5da;
  padding-left: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.record-action-btn {
  border: none;
  border-radius: 999px;
  background: linear-gradient(135deg, #8dc75d 0%, #72aa46 100%);
  color: #ffffff;
  font-size: 0.84rem;
  font-weight: 700;
  line-height: 1;
  padding: 9px 16px;
  cursor: pointer;
  box-shadow: 0 8px 18px rgba(114, 170, 70, 0.24);
  transition: transform 0.18s ease, box-shadow 0.18s ease, filter 0.18s ease;
}

.record-action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 22px rgba(114, 170, 70, 0.32);
  filter: brightness(1.03);
}

.record-action-btn:active {
  transform: translateY(0);
}

.meal-primary-name,
.meal-secondary-line,
.food-text-container,
.quick-add-icon {
  display: none;
}

.meal-side {
  min-width: 0;
  flex-shrink: 0;
}

.mini-ring {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  display: grid;
  place-items: center;
}

.mini-ring-inner {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: #f7fbf6;
  color: #2c5f4c;
  font-size: 0.82rem;
  font-weight: 800;
  display: grid;
  place-items: center;
}

.meal-summary {
  margin: 0;
  color: #2a5d4b;
  font-size: 0.96rem;
  font-weight: 800;
}

.meal-summary-sub {
  margin: 0;
  color: #5f8474;
  font-size: 0.86rem;
  font-weight: 700;
}

.footer-tip {
  margin-top: 8px;
  font-size: 0.76rem;
  color: #346b58;
  font-weight: 700;
  text-align: right;
}

.right-panel {
  display: flex;
  flex-direction: column;
  gap: 18px;
  min-height: 0;
  margin-top: 8px;
}

.overview-card {
  padding: 16px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
  position: relative;
}

.ring-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.calorie-ring {
  width: 190px;
  height: 190px;
  border-radius: 50%;
  display: grid;
  place-items: center;
}

.ring-center {
  width: 142px;
  height: 142px;
  border-radius: 50%;
  background: #f6fbf8;
  display: grid;
  place-items: center;
  text-align: center;
}

.ring-center h3 {
  margin: 0;
  font-size: 2.2rem;
  color: #1f5a46;
}

.ring-center p {
  margin: 2px 0 0;
  font-size: 0.95rem;
  color: #668f80;
  font-weight: 700;
}

.goal-text {
  margin: 8px 0 0;
  font-size: 1.1rem;
  color: #5a8375;
  font-weight: 600;
}

.macro-stack {
  display: flex;
  flex-direction: column;
  gap: 10px;
  justify-content: center;
}

.macro-card {
  border-radius: 12px;
  background: #FAFCFB;
  border: 1px solid #e9f0eb;
  padding: 12px 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.macro-left {
  display: flex;
  gap: 10px;
  align-items: center;
}

.macro-icon {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: grid;
  place-items: center;
}

.macro-label {
  margin: 0;
  color: #245946;
  font-weight: 700;
  font-size: 1.1rem;
}

.macro-target {
     margin: 2px 0 0;
  color: #6d9384;
  font-size: 0.9rem;
}

.macro-value {
  margin: 0;
  color: #1f5947;
  font-weight: 800;
  font-size: 1.3rem;
}

.details-card {
  padding: 16px;
  position: relative;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 520px;
}

.details-stage {
  background: linear-gradient(180deg, #f4f8f6 0%, #f1f6f3 100%);
  box-shadow: inset 0 0 0 1px rgba(138, 168, 151, 0.16), inset 0 8px 20px rgba(108, 141, 123, 0.05), 0 8px 24px rgba(0, 0, 0, 0.04);
}

.details-body {
  display: flex;
  flex-direction: column;
  min-height: 100%;
}

.detail-panel-enter-active,
.detail-panel-leave-active {
  transition: opacity 0.24s ease, transform 0.24s ease;
}

.detail-panel-enter-from,
.detail-panel-leave-to {
  opacity: 0;
  transform: translateY(10px) scale(0.992);
}

.detail-panel-enter-to,
.detail-panel-leave-from {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.details-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.details-head h3 {
  margin: 0;
  color: #235747;
  font-size: 1.8rem;
}

.details-head p {
  margin: 0;
  color: #111111;
  font-size: 1rem;
}

.sample-note-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  color: #111111;
  font-size: 0.95rem;
}

.sample-note {
  margin: 0;
  color: #111111;
  font-size: 0.82rem;
  font-weight: 700;
}

.generate-btn {
  border: none;
  background: #8ec662;
  color: white;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.82rem;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 600;
  box-shadow: 0 4px 8px rgba(142, 198, 98, 0.2);
}

.generate-btn:hover {
  transform: translateY(-1px);
  background: #7bb550;
  box-shadow: 0 6px 12px rgba(142, 198, 98, 0.3);
}

.micro-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.micro-item {
  padding: 10px 10px 8px;
  border-radius: 12px;
  background: #f8fbf9;
  border: 1px solid #e8efea;
}

.micro-top {
  display: flex;
  justify-content: space-between;
  color: #111111;
  font-size: 0.95rem;
  font-weight: 700;
  margin-bottom: 6px;
}

.micro-track {
  height: 6px;
  border-radius: 999px;
  overflow: hidden;
  background: #f0f5f2;
}

.micro-fill {
  height: 100%;
}

.micro-mini-bars {
  height: 22px;
  display: flex;
  align-items: flex-end;
  gap: 4px;
  margin-bottom: 8px;
}

.mini-bar {
  flex: 1;
  border-radius: 8px 8px 2px 2px;
  opacity: 0.75;
  min-height: 6px;
}

.meal-empty-guide {
  margin: 4px 0 12px;
  padding: 12px;
  border: 1px dashed #bad8c0;
  border-radius: 12px;
  background: #f3faf5;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.meal-empty-guide h4 {
  margin: 0 0 2px;
  color: #2f694c;
  font-size: 0.95rem;
}

.meal-empty-guide p {
  margin: 0;
  color: #5d8970;
  font-size: 0.82rem;
}

.plate-illustration {
  width: 56px;
  height: 56px;
  position: relative;
  flex-shrink: 0;
  display: grid;
  place-items: center;
}

.plate-ring {
  width: 46px;
  height: 46px;
  border: 2px solid #8ec662;
  border-radius: 50%;
}

.plate-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #8ec662;
  position: absolute;
}

.food-table {
  border-radius: 12px;
  background: #FAFCFB;
  border: 1px solid #e9f0eb;
  padding: 10px 12px;
  flex: 1;
}

.table-head,
.food-row {
  display: grid;
  grid-template-columns: 28px 1.9fr 0.6fr 0.8fr 0.6fr 0.6fr 0.6fr 32px;
  align-items: center;
  gap: 8px;
}

.check-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}

.food-check-input {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: #8ec662;
}

.table-head {
  color: #111111;
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 6px;
  font-weight: 700;
}

.food-row {
  height: 72px;
  border-top: 1px solid #f0f5f2;
  font-size: 1.05rem;
  color: #2c3e35;
}

.action-col {
  text-align: center;
}

.action-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}

.row-delete-btn {
  background: transparent;
  border: none;
  color: #111111;
  font-size: 1.3rem;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.row-delete-btn:hover {
  background: #fdf5f5;
  color: #e12f3c;
}

.food-empty {
  color: #111111;
}

.food-name {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 600;
  font-size: 1.05rem;
  width: 100%;
  padding: 0;
  border: none;
  background: transparent;
  text-align: left;
  color: #2c3e35;
  font-family: inherit;
  line-height: inherit;
}

.food-name.is-link {
  cursor: pointer;
}

.food-name:disabled {
  cursor: default;
  opacity: 1;
  color: #2c3e35;
}

.food-name.is-link:hover {
  color: #2f7a59;
}

.food-name.is-link:hover .thumb {
  transform: translateY(-1px);
  box-shadow: 0 8px 18px rgba(47, 122, 89, 0.18);
}

.thumb {
  width: 68px;
  height: 68px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  color: #f4fffb;
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.thumb.is-link {
  cursor: pointer;
}

.thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumb-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.thumb-image.recipe-image--placeholder {
  object-fit: contain;
  object-position: center;
  padding: 6px;
  box-sizing: border-box;
  background: #f8fbf8;
}

.sub-text {
  font-size: 0.9rem;
  color: #111111;
}

.unit-text {
  font-size: 0.85rem;
  color: #111111;
  margin-left: 2px;
}

.add-more {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.add-more input {
  flex: 1;
  width: 100%;
  height: 40px;
  border: 1px solid #dce8e1;
  border-radius: 10px;
  background: #FAFCFB;
  padding: 0 12px;
  color: #2e5346;
  font-size: 0.95rem;
  outline: none;
}

.add-more input::placeholder {
  color: #7f9a8f;
}

.plus-btn {
  width: 40px;
  height: 40px;
  border: 1px solid #cfe1d7;
  border-radius: 10px;
  background: #eef7f1;
  color: #2f694c;
  font-size: 1.25rem;
  line-height: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.plus-btn:hover {
  background: #e5f3ea;
  transform: translateY(-1px);
}

/* AI 识别 Loading 遮罩 */
.ai-loading-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(244, 246, 243, 0.85); /* 亲自然半透 */
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(6px);
}

.ai-loading-inner {
  text-align: center;
  background: #FAFCFB;
  padding: 32px 48px;
  border-radius: 20px;
  box-shadow: 0 16px 48px rgba(32, 59, 45, 0.1);
}
</style>
