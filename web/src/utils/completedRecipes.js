const COMPLETED_RECIPE_KEY_PREFIX = 'nutriplan_completed_recipes'

export const formatDateKey = (date = new Date()) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const getStorageKey = (userId) => `${COMPLETED_RECIPE_KEY_PREFIX}_${userId || 'guest'}`

export const getCompletedRecipeMap = (userId) => {
  try {
    const raw = localStorage.getItem(getStorageKey(userId))
    if (!raw) return {}
    const parsed = JSON.parse(raw)
    if (!parsed || typeof parsed !== 'object') return {}
    return parsed
  } catch (error) {
    console.warn('读取完成食谱记录失败', error)
    return {}
  }
}

const saveCompletedRecipeMap = (userId, data) => {
  localStorage.setItem(getStorageKey(userId), JSON.stringify(data || {}))
}

export const updateCompletedRecipe = ({ userId, dateKey, mealType, recipe, checked }) => {
  if (!dateKey || !mealType) return

  const map = getCompletedRecipeMap(userId)
  const dayList = Array.isArray(map[dateKey]) ? map[dateKey] : []
  const recipeId = recipe?.id || `${mealType}_${recipe?.name || 'unknown'}`
  const filtered = dayList.filter((item) => !(item.mealType === mealType && String(item.recipeId) === String(recipeId)))

  if (checked && recipe) {
    filtered.push({
      recipeId,
      mealType,
      name: recipe.name || '未命名食谱',
      energy: Number(recipe.energy || recipe.calories || 0),
      protein: Number(recipe.protein || recipe.total_protein || 0),
      carbohydrate: Number(recipe.carbohydrate || recipe.carbs || recipe.total_carbohydrate || 0),
      fat: Number(recipe.fat || recipe.total_fat || 0)
    })
  }

  map[dateKey] = filtered
  saveCompletedRecipeMap(userId, map)
}

export const buildCompletedRecordsForDate = (userId, dateKey) => {
  const map = getCompletedRecipeMap(userId)
  const dayList = Array.isArray(map[dateKey]) ? map[dateKey] : []

  return dayList.map((item) => ({
    id: `completed-${dateKey}-${item.mealType}-${item.recipeId}`,
    mealType: item.mealType,
    foodName: item.name,
    intakeAmount: 1,
    intakeUnit: '份',
    calculatedEnergy: Number(item.energy || 0),
    calculatedProtein: Number(item.protein || 0),
    calculatedCarb: Number(item.carbohydrate || 0),
    calculatedFat: Number(item.fat || 0),
    isCompletedRecipe: true
  }))
}