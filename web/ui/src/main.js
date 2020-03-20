import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import { CHECK_AUTH, CHECK_NETWORK } from './store/actions.type'
import ApiService from './common/api.service'
import DateFilter from './common/date.filter'
import ErrorFilter from './common/error.filter'

Vue.use(BootstrapVue)

Vue.config.productionTip = false
Vue.filter('date', DateFilter)
Vue.filter('error', ErrorFilter)

ApiService.init()

router.beforeEach((to, from, next) =>
  Promise.all([store.dispatch(CHECK_AUTH)]).then(next)
)

router.beforeEach((to, from, next) =>
  Promise.all([store.dispatch(CHECK_NETWORK)]).then(next)
)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
