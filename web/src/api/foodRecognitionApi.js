import request from './request'

/**
 * 上传图片进行食物识别
 * @param {File} imageFile - 用户选择的图片文件
 * @returns {Promise<any>} - 包含识别结果的 Promise
 */
export const recognizeFood = (imageFile) => {
  // 创建 FormData 对象来封装图片数据
  const formData = new FormData()
  // 'image' 是后端 Gin 框架中 c.FormFile("image") 需要的字段名
  formData.append('image', imageFile)

  // 发送 POST 请求
  // 需要重写 headers，将 Content-Type 设置为 'multipart/form-data'
  // axios 会自动处理这种类型
  return request({
    url: '/food/recognize',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    // 可以为耗时长的请求设置更长的超时时间
    timeout: 60000 // 60秒
  })
}
