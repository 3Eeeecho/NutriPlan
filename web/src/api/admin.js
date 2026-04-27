import request from './request'

export const getAdminStatus = () => {
  return request({
    url: '/admin/status',
    method: 'get'
  })
}

export const getAdminRecipes = (params = {}) => {
  return request({
    url: '/admin/recipes',
    method: 'get',
    params
  })
}

export const getAdminRecipe = (id) => {
  return request({
    url: `/admin/recipes/${id}`,
    method: 'get'
  })
}

export const createAdminRecipe = (data) => {
  return request({
    url: '/admin/recipes',
    method: 'post',
    data
  })
}

export const updateAdminRecipe = (id, data) => {
  return request({
    url: `/admin/recipes/${id}`,
    method: 'put',
    data
  })
}

export const deleteAdminRecipe = (id) => {
  return request({
    url: `/admin/recipes/${id}`,
    method: 'delete'
  })
}

export const getAdminHealthStats = (params = {}) => {
  return request({
    url: '/admin/stats/health',
    method: 'get',
    params
  })
}

export const getAdminRecipeCompletionStats = (params = {}) => {
  return request({
    url: '/admin/stats/recipe-completion',
    method: 'get',
    params
  })
}
