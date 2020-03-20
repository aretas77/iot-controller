import { NetworkService } from '@/common/api.service'
import {
  CHECK_NETWORK,
  FETCH_NETWORK,
  SET_NETWORK
} from './actions.type'

const initialState = {
  // This should be a global network for a user
  network: {
    name: ''
  },
  networks: [],
  networksCount: 0,
  isLoadingNetworks: true
}

export const state = { ...initialState }

export const actions = {
  // FETCH_NETWORK will get the network by its ID.
  async [FETCH_NETWORK] (context, networkSlug) {
    console.log('FETCH_NETWORK start')
    const { data } = await NetworkService.get(networkSlug)
    context.commit(SET_NETWORK, data)
    return data
  },

  [CHECK_NETWORK] (context) {
    // Get current user ID from ProfileService.
    // const user = context.dispatch(CHECK_AUTH)
  }
}

/* eslint no-param-reassign: ["error", { "props": false }] */
export const mutations = {
  [SET_NETWORK] (state, network) {
    console.log('SET_NETWORK' + network)
    state.network = network
  }
}

const getters = {
  currentNetwork (state) {
    return state.network
  },
  networks (state) {
    return state.networks
  },
  networksCount (state) {
    return state.networksCount
  },
  isLoadingNetworks (state) {
    return state.isLoadingNetworks
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
