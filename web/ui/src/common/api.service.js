import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import { API_URL } from '@/common/config'

const ApiService = {
  init () {
    Vue.use(VueAxios, axios)
    Vue.axios.defaults.baseURL = API_URL
  },
  setHeader () {
    Vue.axios.defaults.headers.common.Authorization = 'Token ad'
  },
  query (resource, params) {
    return Vue.axios.get(resource, params).catch(error => {
      throw new Error(`[Iotctl] ApiService ${error}`)
    })
  },
  get (resource, slug = '') {
    return Vue.axios.get(`${resource}/${slug}`).catch(error => {
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
    return Vue.axios.delete(`${resource}/${resource}`).catch(error => {
      throw new Error(`[Iotctl] ApiService ${error}`)
    })
  }
}

export default ApiService

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
  create (params) {
    return ApiService.post('nodes', { node: params })
  },
  update (slug, params) {
    return ApiService.update('nodes', slug, { node: params })
  },
  destroy (slug) {
    return ApiService.update(`nodes/${slug}`)
  }
}
