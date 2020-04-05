import ApiService from '@/common/api.service'
import JwtService from '@/common/jwt.service'
import NetworkService from '@/common/network.service'
import {
  LOGIN,
  LOGOUT,
  CHECK_AUTH,
  FETCH_CURRENT_USER
} from './actions.type'
import { SET_AUTH, PURGE_AUTH, SET_ERROR } from './mutations.type'

const state = {
  errors: null,
  user: {},
  isAuthenticated: !!JwtService.getToken()
}

const getters = {
  currentUser (state) {
    return state.user
  },
  isAuthenticated (state) {
    return state.isAuthenticated
  }
}

// actions for auth should provide some way of:
//  - Login
//  - Logout
//  - Check Authentication header
//  - Register a new user
const actions = {
  [LOGIN] (context, credentials) {
    const email = credentials.email
    const password = credentials.password
    // clear state
    context.commit(PURGE_AUTH)
    return new Promise(resolve => {
      ApiService.post('login', { email, password })
        .then(({ data }) => {
          context.commit(SET_AUTH, data)
          resolve(data)
        })
        .catch(({ response }) => {
          context.commit(SET_ERROR, response.data.errors)
        })
    })
  },
  [LOGOUT] (context) {
    context.commit(PURGE_AUTH)
  },

  // Is called each time we load a new View.
  [CHECK_AUTH] (context) {
    console.log('Calling CHECK_AUTH')

    // Getting token from JwtService - currently local storage.
    if (JwtService.getToken()) {
      ApiService.setHeader()
      ApiService.get('users/check')
        .then(({ data }) => {
          context.commit(SET_AUTH, data)
        })
        .catch(({ error }) => {
          context.commit(SET_ERROR, error)
          context.commit(PURGE_AUTH)
        })
    } else {
      context.commit(PURGE_AUTH)
    }
  }
}

const mutations = {
  [SET_ERROR] (state, error) {
    console.log('Setting SET_ERROR: ' + error)
    state.errors = error
  },
  [SET_AUTH] (state, user) {
    console.log('Setting SET_AUTH')
    state.isAuthenticated = true
    state.user = user
    state.errors = {}
    JwtService.saveToken(state.user.token)
  },
  [PURGE_AUTH] (state) {
    console.log('Calling PURGE_AUTH')
    state.isAuthenticated = false
    state.user = {}
    state.errors = {}

    // We don't need to see them - destroy
    JwtService.destroyToken()
    NetworkService.destroyNetwork()
  },
  [FETCH_CURRENT_USER] (state) {
    console.log('Calling FETCH_CURRENT_USER')
    return state.user
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
