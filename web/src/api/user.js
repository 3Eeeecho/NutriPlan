import request from './request'

// 用户注册
export const register = (data) => {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

// 用户登录
export const login = (data) => {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// 获取用户档案
export const getUserProfile = () => {
  return request({
    url: '/user/profile',
    method: 'get'
  })
}

// 更新用户档案
export const updateUserProfile = (data) => {
  return request({
    url: '/user/profile',
    method: 'put',
    data
  })
}

// 获取营养需求
export const getNutritionRequirements = () => {
  return request({
    url: '/user/nutrition',
    method: 'get'
  })
}

