import axios from 'axios'
// import Cookies from 'js-cookie'
import qs from 'qs'
import isPlainObject from 'lodash/isPlainObject'
import {Loading} from 'element-ui'

const http = axios.create({
    baseURL: '/api',
    timeout: 1000 * 180,
    withCredentials: true
})

let loadingObj = null
/**
 * 请求拦截
 */
http.interceptors.request.use(config => {
    // 默认参数
    let defaults = {}
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
    // 由showLoading参数来判断是否需要显示loading
    if (config.showLoading) {
        loadingObj = null
        loadingObj = Loading.service({
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
    loadingObj && loadingObj.close()
    loadingObj = null
    if (response.data.code === 401) {
        return Promise.reject(response.data.msg)
    }
    return response
}, error => {
    loadingObj && loadingObj.close()
    loadingObj = null
    console.error(error)
    alert("请求失败" + error)
    return Promise.reject(error)
})

export default http
