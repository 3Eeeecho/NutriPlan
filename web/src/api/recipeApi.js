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
export const selectRecipePlan = (planId) => {
  return request({
    url: `/recipes/plan/${planId}/select`,
    method: 'post'
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
