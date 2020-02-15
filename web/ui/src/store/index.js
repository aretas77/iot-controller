import Vue from 'vue'
import Vuex from 'vuex'

import node from './node.module'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    node
  }
})
