import request from './request';

// 获取食谱推荐
export const getRecipeRecommendations = (count = 3) => {
  return request({
    url: '/recipes/recommend',
    method: 'get',
    params: { count }
  });
};

// 选择食谱计划
export const selectRecipePlan = (plan) => {
  return request({
    url: `/recipes/plan/select`,
    method: 'post',
    data: plan
  });
};

// 保存推荐结果
export const saveRecommendations = (plans) => {
  return request({
    url: '/recipes/recommend',
    method: 'post',
    data: plans
  });
};

// 获取已选食谱计划
export const getSelectedRecipePlan = () => {
  return request({
    url: '/recipes/selected',
    method: 'get'
  });
};

// 获取食谱详情
export const getRecipeDetail = (id) => {
  return request({
    url: `/recipes/${id}`,
    method: 'get'
  });
};

// 添加收藏
export const addFavorite = (id) => {
  return request({
    url: `/recipes/${id}/favorite`,
    method: 'post'
  });
};

// 取消收藏
export const removeFavorite = (id) => {
  return request({
    url: `/recipes/${id}/favorite`,
    method: 'delete'
  });
};

// 获取收藏列表
export const getFavoriteList = () => {
  return request({
    url: '/recipes/favorites',
    method: 'get'
  });
};
