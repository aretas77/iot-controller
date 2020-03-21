import Vue from 'vue'
import { NodesService } from '@/common/api.service'
import {
  NODE_GET_ALL,
  FETCH_NODE,
  FETCH_NODES,
  NODE_EDIT,
  NODE_REMOVE,
  NODE_ADD,
  NODE_RESET_STATE
} from './actions.type'
import {
  SET_NODE,
  SET_NODES,
  FETCH_START,
  FETCH_END,
  RESET_STATE
} from './mutations.type'

// initialState
const initialState = {
  node: {
    name: ''
  },
  nodes: [],
  nodesCount: 0,
  isLoading: true
}

export const state = { ...initialState }

export const actions = {
  // FETCH_NODE get node by its ID.
  async [FETCH_NODE] (context, nodeSlug, prevNode) {
    console.log('FETCH_NODE start')
    // avoid additional network call if node exists
    if (prevNode !== undefined) {
      return context.commit(SET_NODE, prevNode)
    }
    const { data } = await NodesService.get(nodeSlug)
    context.commit(SET_NODE, data)
    return data
  },

  // NODE_GET_ALL fetches all Node devices.
  async [NODE_GET_ALL] (context, nodeSlug) {
    console.log('NODE_GET_ALL start')
    const { data } = await NodesService.get(nodeSlug)
    context.commit(SET_NODES, data)
    return data
  },

  // FETCH_NODES fetches Node devices by given filters.
  [FETCH_NODES] ({ commit }, params) {
    console.log('FETCH_NODES start')
    commit(FETCH_START)
    return NodesService.query(params.type, params.filters)
      .then(({ data }) => {
        commit(FETCH_END, data)
      })
      .catch(error => {
        throw new Error(error)
      })
  },
  [NODE_ADD] ({ state }) {
    console.log('NODE_ADD call')
    return NodesService.create(state.node)
  },
  [NODE_EDIT] ({ state }) {
    return NodesService.update(state.node.slug, state.node)
  },
  [NODE_REMOVE] (context, slug) {
    return NodesService.destroy(slug)
  },
  [NODE_RESET_STATE] ({ commit }) {
    console.log('NODE_RESET_STATE start')
    commit(RESET_STATE)
  }

}

/* eslint no-param-reassign: ["error", { "props": false }] */
export const mutations = {
  [SET_NODE] (state, node) {
    console.log('SET_NODE' + node)
    state.node = node
  },
  [SET_NODES] (state, nodes) {
    state.nodes = nodes
  },
  [FETCH_START] (state) {
    state.isLoading = true
  },
  [FETCH_END] (state, { nodes }) {
    state.nodes = nodes
    state.nodesCount = nodes.length
    state.isLoading = false
  },
  [RESET_STATE] () {
    for (const f in state) {
      Vue.set(state, f, initialState[f])
    }
  }
}

const getters = {
  node (state) {
    return state.node
  },
  nodes (state) {
    return state.nodes
  },
  nodesCount (state) {
    return state.nodesCount
  },
  isLoading (state) {
    return state.isLoading
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
