import { EventsService } from '@/common/api.service'
import {
  FETCH_NODE_EVENTS
} from './actions.type'
import {
  FETCH_NODE_EVENTS_START,
  FETCH_NODE_EVENTS_END
} from './mutations.type'

const initialState = {
  node: 0,
  events: [],
  eventsCount: 0,
  isLoadingEvents: false
}

export const state = { ...initialState }

export const actions = {
  async [FETCH_NODE_EVENTS] (context, nodeSlug) {
    console.log('FETCH_NODE_EVENTS start')
    context.commit(FETCH_NODE_EVENTS_START)
    await EventsService.getByNode(nodeSlug)
      .then(({ data }) => {
        context.commit(FETCH_NODE_EVENTS_END, data)
      })
      .catch(error => {
        throw new Error(error)
      })
  }
}

export const mutations = {
  [FETCH_NODE_EVENTS_END] (state, { events }) {
    state.events = events
    state.eventsCount = events.length
    state.isLoadingEvents = false
  },
  [FETCH_NODE_EVENTS_START] (state) {
    state.isLoadingEvents = true
  }
}

const getters = {
  eventNodeId (state) {
    return state.node
  },
  events (state) {
    return state.events
  },
  eventsCount (state) {
    return state.eventsCount
  }
}

export default {
  state,
  actions,
  mutations,
  getters
}
