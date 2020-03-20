import ApiService from '@/common/api.service'
import {
  FETCH_PROFILE
} from './actions.type'
import { SET_PROFILE } from './mutations.type'

const state = {
  errors: {},
  profile: {}
}

const getters = {
  profile (state) {
    return state.profile
  }
}

const actions = {
  [FETCH_PROFILE] (context, payload) {
    const { username } = payload
    return ApiService.get('user', username)
      .then(({ data }) => {
        context.commit(SET_PROFILE, data.profile)
        return data
      })
      .catch(() => {

      })
  }
}

const mutations = {
  [SET_PROFILE] (state, profile) {
    state.profile = profile
    state.errors = {}
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
