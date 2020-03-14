import ApiService from '@/common/api.service'
import JwtService from '@/common/jwt.service'
import {
  LOGIN,
  LOGOUT,
  CHECK_AUTH
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
    console.log('hello ' + credentials.email)
    return new Promise(resolve => {
      ApiService.post('users/login', { user: credentials })
        .then(({ data }) => {
          context.commit(SET_AUTH, data.user)
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
  [CHECK_AUTH] (context) {
    if (JwtService.getToken()) {
      ApiService.setHeader()
      ApiService.get('user')
        .then(({ data }) => {
          context.commit(SET_AUTH, data.user)
        })
        .catch(({ response }) => {
          context.commit(SET_ERROR, response.data.errors)
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
    state.isAuthenticated = true
    state.user = user
    state.errors = {}
    JwtService.saveToken(state.user.token)
  },
  [PURGE_AUTH] (state) {
    state.isAuthenticated = false
    state.user = {}
    state.errors = {}
    JwtService.destroyToken()
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
