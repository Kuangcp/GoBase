/**
 * 配置参考: https://cli.vuejs.org/zh/config/
 */
module.exports = {
  publicPath: './',
  productionSourceMap: false,
  devServer: {
    // host: '',
    open: true,
    port: 8081,
    overlay: {
      errors: true,
      warnings: true
    },
    proxy: {
      // 配置跨域, 普通请求
      '/api': {
        target: 'http://localhost:9090',
        ws: false,
        changOrigin: true
      }
    }
  }
}
