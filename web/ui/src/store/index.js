import Vue from 'vue'
import Vuex from 'vuex'

import node from './node.module'
import auth from './auth.module'
import netw from './network.module'
import evts from './event.module'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    auth,
    node,
    netw,
    evts
  }
})
