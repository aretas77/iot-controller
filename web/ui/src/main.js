import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import { CHECK_AUTH } from './store/actions.type'
import ApiService from './common/api.service'
import DateFilter from './common/date.filter'
import TimeFilter from './common/time.filter'
import ErrorFilter from './common/error.filter'

Vue.use(BootstrapVue)

Vue.config.productionTip = false
Vue.filter('date', DateFilter)
Vue.filter('error', ErrorFilter)
Vue.filter('time', TimeFilter)

ApiService.init()

router.beforeEach((to, from, next) =>
  Promise.all([store.dispatch(CHECK_AUTH)]).then(next)
)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
