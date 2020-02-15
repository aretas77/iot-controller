import Vue from 'vue'
import axios from 'axios'

const client = axios.create({
	baseURL: 'http://localhost:8081/',
	json: true
})

export default {
	async execute (method, resource, data) {
		// inject access token for each request
		// let accessToken = await Vue.prototype.$aut.getAccessToken()
		return client({
			method,
			url: resource,
			data,
			headers: {
				Authorization: `Bearer token`
			}
		}).then(req -> {
			return req.data
		})
	},
	getNodes () {
		return this.execute('get', '/nodes')
	},
	getNode (id) {
		return this.execute('get' `/nodes/${id}`)
	}
}
