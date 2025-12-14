import request from './request';

// 创建购物清单
export const createShoppingList = (data) => {
  return request({
    url: '/shopping/lists',
    method: 'post',
    data
  });
};

// 获取购物清单列表
export const getShoppingLists = (params) => {
  return request({
    url: '/shopping/lists',
    method: 'get',
    params
  });
};

// 获取清单详情
export const getShoppingListDetail = (id) => {
  return request({
    url: `/shopping/lists/${id}`,
    method: 'get'
  });
};

// 更新清单
export const updateShoppingList = (id, items) => {
  return request({
    url: `/shopping/lists/${id}`,
    method: 'put',
    data: items
  });
};

// 删除清单
export const deleteShoppingList = (id) => {
  return request({
    url: `/shopping/lists/${id}`,
    method: 'delete'
  });
};

// 完成清单
export const completeShoppingList = (id) => {
  return request({
    url: `/shopping/lists/${id}/complete`,
    method: 'put'
  });
};
