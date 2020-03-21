import { NetworkService } from '@/common/api.service'
import NetworkStorageService from '@/common/network.service'
import {
  FETCH_NETWORK,
  FETCH_NETWORKS,
  UPDATE_CURRENT_NETWORK,
  NETWORK_RESET_STATE
} from './actions.type'
import {
  SET_NETWORKS,
  SET_NETWORK,
  PURGE_NETWORK
} from './mutations.type'

const initialState = {
  // This should be a global network for a user
  network: NetworkStorageService.getNetwork(),
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

  async [FETCH_NETWORKS] (context, userID) {
    if (userID === undefined) {
      return null
    }

    const { data } = await NetworkService.getByUser(userID)
    context.commit(SET_NETWORKS, data)
    return data
  },

  async [UPDATE_CURRENT_NETWORK] (context, network) {
    console.log('UPDATE_CURRENT_NETWORK start')
    if (network === null || network === undefined || network.name === '') {
      return null
    }
    context.commit(SET_NETWORK, network)
  },

  async [NETWORK_RESET_STATE] (context) {
    console.log('NETWORK_RESET_STATE start')
    context.commit(PURGE_NETWORK)
  }
}

/* eslint no-param-reassign: ["error", { "props": false }] */
export const mutations = {
  [SET_NETWORK] (state, network) {
    console.log('SET_NETWORK')
    state.network = network
    NetworkStorageService.saveNetwork(network)
  },
  [SET_NETWORKS] (state, networks) {
    console.log('SET_NETWORKS')
    state.networks = networks
  },
  [PURGE_NETWORK] (state) {
    state.network = {}
    state.networks = {}
    state.networksCount = 0
  }
}

const getters = {
  currentNetwork (state) {
    if (state.network === null) {
      return {}
    }
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
