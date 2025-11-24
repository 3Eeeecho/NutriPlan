<template>
  <div class="recipe-card">
    <div class="card-header">
      <span class="meal-icon">{{ icon }}</span>
      <h4>{{ title }}</h4>
    </div>
    <div class="card-body">
      <h5 class="recipe-name">{{ recipe.name }}</h5>
      <p v-if="recipe.description" class="recipe-desc">{{ recipe.description }}</p>
      
      <div class="nutrition-info">
        <div class="nutrition-badge">
          <span class="badge-label">热量</span>
          <span class="badge-value">{{ Math.floor(recipe.energy) }} kcal</span>
        </div>
        <div class="nutrition-badge">
          <span class="badge-label">蛋白质</span>
          <span class="badge-value">{{ Math.floor(recipe.protein) }}g</span>
        </div>
        <div class="nutrition-badge">
          <span class="badge-label">碳水</span>
          <span class="badge-value">{{ Math.floor(recipe.carbohydrate) }}g</span>
        </div>
        <div class="nutrition-badge">
          <span class="badge-label">脂肪</span>
          <span class="badge-value">{{ Math.floor(recipe.fat) }}g</span>
        </div>
      </div>

      <div v-if="recipe.ingredients && recipe.ingredients.length > 0" class="ingredients">
        <span class="ingredients-label">食材：</span>
        <span class="ingredients-list">{{ formatIngredients(recipe.ingredients) }}</span>
      </div>

      <div v-if="recipe.cooking_time" class="cooking-info">
        <span>⏱️ {{ recipe.cooking_time }}分钟</span>
        <span v-if="recipe.difficulty">{{ getDifficultyEmoji(recipe.difficulty) }} {{ recipe.difficulty }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'RecipeCard',
  props: {
    title: {
      type: String,
      required: true
    },
    recipe: {
      type: Object,
      required: true
    },
    icon: {
      type: String,
      default: '🍽️'
    }
  },
  methods: {
    formatIngredients(ingredients) {
      if (Array.isArray(ingredients)) {
        return ingredients.join('、');
      }
      return ingredients;
    },
    getDifficultyEmoji(difficulty) {
      const emojiMap = {
        '简单': '⭐',
        '中等': '⭐⭐',
        '困难': '⭐⭐⭐'
      };
      return emojiMap[difficulty] || '⭐';
    }
  }
};
</script>

<style scoped>
.recipe-card {
  background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.recipe-card:hover {
  transform: scale(1.02);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.meal-icon {
  font-size: 1.5rem;
}

.card-header h4 {
  margin: 0;
  font-size: 1rem;
  color: #333;
  font-weight: 600;
}

.card-body {
  padding: 1rem;
  background: rgba(255, 255, 255, 0.6);
}

.recipe-name {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  color: #2c3e50;
  font-weight: 700;
}

.recipe-desc {
  margin: 0 0 1rem 0;
  font-size: 0.85rem;
  color: #666;
}

.nutrition-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.nutrition-badge {
  background: white;
  padding: 0.4rem 0.6rem;
  border-radius: 6px;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.badge-label {
  display: block;
  font-size: 0.7rem;
  color: #888;
  margin-bottom: 0.2rem;
}

.badge-value {
  display: block;
  font-size: 0.9rem;
  font-weight: 700;
  color: #333;
}

.ingredients {
  margin-bottom: 0.75rem;
  font-size: 0.85rem;
  color: #555;
}

.ingredients-label {
  font-weight: 600;
}

.ingredients-list {
  color: #666;
}

.cooking-info {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  color: #777;
  padding-top: 0.5rem;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}
</style>
