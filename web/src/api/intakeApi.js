import request from './request'

/**
 * 添加饮食记录
 * @param {Object} data - 饮食记录数据
 * @returns {Promise}
 */
export function addIntakeRecord(data) {
  return request({
    url: '/intake/records',
    method: 'post',
    data
  })
}

/**
 * 删除饮食记录
 * @param {Number} id - 记录ID
 * @returns {Promise}
 */
export function deleteIntakeRecord(id) {
  return request({
    url: `/intake/records/${id}`,
    method: 'delete'
  })
}

/**
 * 获取当日营养状态
 * @returns {Promise}
 */
export function getTodayStatus() {
  return request({
    url: '/intake/today',
    method: 'get'
  })
}

/**
 * 获取周报告
 * @returns {Promise}
 */
export function getWeeklyReport() {
  return request({
    url: '/intake/weekly',
    method: 'get'
  })
}
