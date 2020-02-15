import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import NodesManager from '@/components/NodesManager'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('../views/About.vue')
  },
  {
    path: '/nodes-manager',
    name: 'NodesManager',
    component: NodesManager,
    meta: {
    }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
