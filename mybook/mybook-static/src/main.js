import Vue from 'vue'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import ECharts from 'vue-echarts'

import App from './App.vue'
import http from './util/request';

Vue.component('v-chart', ECharts)

Vue.config.productionTip = false
// 挂载全局
Vue.prototype.$http = http

Vue.use(ElementUI);

new Vue({
  render: h => h(App),
}).$mount('#app')
