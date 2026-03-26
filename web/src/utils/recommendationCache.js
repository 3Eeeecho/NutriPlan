const RECOMMENDATION_CACHE_KEY_PREFIX = 'nutriplan_daily_recommendation_map'
const LEGACY_RECOMMENDATION_CACHE_KEY = 'nutriplan_daily_recommendation'

const getStorageKey = (userId) => `${RECOMMENDATION_CACHE_KEY_PREFIX}_${userId || 'guest'}`

const normalizeDateKey = (dateKey = '') => {
  const match = String(dateKey).match(/^(\d{4})-(\d{1,2})-(\d{1,2})$/)
  if (!match) return String(dateKey)
  return `${match[1]}-${String(match[2]).padStart(2, '0')}-${String(match[3]).padStart(2, '0')}`
}

export const getRecommendationMap = (userId) => {
  try {
    const raw = localStorage.getItem(getStorageKey(userId))
    const parsed = raw ? JSON.parse(raw) : {}
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      return {}
    }
    return parsed
  } catch (error) {
    console.warn('读取推荐缓存失败', error)
    return {}
  }
}

const saveRecommendationMap = (userId, map) => {
  localStorage.setItem(getStorageKey(userId), JSON.stringify(map || {}))
}

const migrateLegacyRecommendationCache = (userId) => {
  try {
    const legacyRaw = localStorage.getItem(LEGACY_RECOMMENDATION_CACHE_KEY)
    if (!legacyRaw) return getRecommendationMap(userId)

    const legacy = JSON.parse(legacyRaw)
    const legacyDate = normalizeDateKey(legacy?.date || '')
    const legacyData = legacy?.data

    const map = getRecommendationMap(userId)
    if (legacyDate && legacyData && !map[legacyDate]) {
      map[legacyDate] = legacyData
      saveRecommendationMap(userId, map)
    }

    localStorage.removeItem(LEGACY_RECOMMENDATION_CACHE_KEY)
    return map
  } catch (error) {
    console.warn('迁移旧版推荐缓存失败', error)
    return getRecommendationMap(userId)
  }
}

export const loadRecommendationMapWithMigration = (userId) => {
  return migrateLegacyRecommendationCache(userId)
}

export const getRecommendationByDate = (userId, dateKey) => {
  const map = loadRecommendationMapWithMigration(userId)
  const normalizedDate = normalizeDateKey(dateKey)
  return normalizedDate ? map[normalizedDate] || null : null
}

export const setRecommendationByDate = (userId, dateKey, recommendation) => {
  const normalizedDate = normalizeDateKey(dateKey)
  if (!normalizedDate) return

  const map = loadRecommendationMapWithMigration(userId)
  map[normalizedDate] = recommendation || null
  saveRecommendationMap(userId, map)
}

export const removeRecommendationMealByDate = (userId, dateKey, mealType) => {
  const normalizedDate = normalizeDateKey(dateKey)
  if (!normalizedDate || !mealType) return

  const map = loadRecommendationMapWithMigration(userId)
  const dayRecommendation = map[normalizedDate]
  if (!dayRecommendation || typeof dayRecommendation !== 'object') return

  const next = { ...dayRecommendation }
  delete next[mealType]

  if (Object.keys(next).length === 0) {
    delete map[normalizedDate]
  } else {
    map[normalizedDate] = next
  }

  saveRecommendationMap(userId, map)
}

export const clearRecommendationMap = (userId) => {
  localStorage.removeItem(getStorageKey(userId))
  localStorage.removeItem(LEGACY_RECOMMENDATION_CACHE_KEY)
}
