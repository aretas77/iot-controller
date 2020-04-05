import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import JwtService from '@/common/jwt.service'
import { API_URL } from '@/common/config'

const ApiService = {
  init () {
    Vue.use(VueAxios, axios)
    Vue.axios.defaults.baseURL = API_URL
  },

  setHeader () {
    Vue.axios.defaults.headers.common.Authorization = `${JwtService.getToken()}`
  },

  query (resource, params) {
    return Vue.axios.get(resource, params).catch(error => {
      throw new Error(`[Iotctl] ApiService ${error}`)
    })
  },

  // TODO: maybe transfer more specific calls to upper layer (e.g. NodeService
  // or any other service)?
  get (resource, slug = '') {
    var params = null
    if (slug.filters != null && slug.filters.limit != null) {
      params = {
        offset: slug.filters.offset,
        limit: slug.filters.limit
      }
    }

    var requestUrl
    if (params != null || slug !== '') {
      requestUrl = `${resource}/${slug}`
    } else {
      requestUrl = `${resource}`
    }

    return Vue.axios.get(requestUrl, {
      params
    }).catch(error => {
      throw new Error(`[Iotctl] ApiService ${error}`)
    })
  },

  post (resource, params) {
    return Vue.axios.post(`${resource}`, params)
  },

  update (resource, slug, params) {
    return Vue.axios.put(`${resource}/${slug}`, params)
  },

  delete (resource) {
    return Vue.axios.delete(`${resource}`).catch(error => {
      throw new Error(`[Iotctl] ApiService ${error}`)
    })
  }
}

export default ApiService

// An exported Wrapper of ApiService.
export const NodesService = {
  // Will run a query on nodes.
  query (type, params) {
    return ApiService.query('nodes', {
      params: params
    })
  },
  // Will run /nodes/${id} and return a given Node.
  get (slug) {
    return ApiService.get('nodes', slug)
  },
  getUnregistered (networkName) {
    return ApiService.get(`networks/${networkName}/unregistered`)
  },
  create (params) {
    return ApiService.post('nodes', { node: params })
  },
  createUnregistered (params) {
    return ApiService.post('nodes', params)
  },
  update (slug, params) {
    return ApiService.update('nodes', slug, { node: params })
  },
  destroy (slug) {
    return ApiService.delete(`nodes/${slug}`)
  }
}

export const NetworkService = {
  query (type, params) {
    return ApiService.query('networks', {
      params: params
    })
  },
  get (slug) {
    return ApiService.get('networks', slug)
  },
  getByUser (userId) {
    return ApiService.get(`users/${userId}/networks`)
  },
  create (params) {
    return ApiService.post('networks', { network: params })
  },
  update (slug, params) {
    return ApiService.update('networks', slug, { network: params })
  }
}

export const EventsService = {
  getByNode (nodeId) {
    return ApiService.get(`nodes/${nodeId}/events`)
  }
}

export const StatisticsService = {
  getByNode (nodeId) {
    return ApiService.get(`nodes/${nodeId}/statistics`)
  }
}
