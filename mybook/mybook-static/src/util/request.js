import axios from 'axios'
// import Cookies from 'js-cookie'
import qs from 'qs'
import isPlainObject from 'lodash/isPlainObject'
// import { Loading } from '@baturu/heiner-ui'
import { Loading } from 'element-ui'

const http = axios.create({
  baseURL: 'http://localhost:9090',
  timeout: 1000 * 180,
  withCredentials: true
})

let loadinginstance = null
/**
 * 请求拦截
 */
http.interceptors.request.use(config => {
  // config.headers['Accept-Language'] = Cookies.get('language') || 'zh-CN'
  // config.headers['token'] = Cookies.get('token') || ''
  // config.headers['CSRFTOKEN'] = Cookies.get('CSRFTOKEN') || ''
  // 默认参数
  var defaults = {}
  // 防止缓存，GET请求默认带_t参数
  if (config.method === 'get') {
    config.params = {
      ...config.params,
      // ...{ '_t': new Date().getTime() }
    }
  }
  if (isPlainObject(config.params)) {
    config.params = {
      ...defaults,
      ...config.params
    }
  }
  if (isPlainObject(config.data)) {
    config.data = {
      ...defaults,
      ...config.data
    }
    if (/^application\/x-www-form-urlencoded/.test(config.headers['content-type'])) {
      config.data = qs.stringify(config.data)
    }
  }
  // 由showLoding参数来判断是否需要显示loading
  if (config.showLoding) {
    loadinginstance = null
    loadinginstance = Loading.service({
      fullscreen: true,
      lock: true,
      text: 'Loading',
      background: 'rgba(0, 0, 0, 0.5)'
    })
  }
  return config
}, error => {
  return Promise.reject(error)
})

/**
 * 响应拦截
 */
http.interceptors.response.use(response => {
  loadinginstance && loadinginstance.close()
  loadinginstance = null
  if (response.data.code === 401) {
    return Promise.reject(response.data.msg)
  }
  return response
}, error => {
  loadinginstance && loadinginstance.close()
  loadinginstance = null
  console.error(error)
  return Promise.reject(error)
})

export default http
