import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login')
  },
  {
    path: '/@:username',
    component: () => import('@/views/Profile'),
    children: [
      {
        path: '/@:username/profile',
        name: 'profile',
        component: () => import('@/views/ProfileNodes')
      }
    ]
  },
  {
    path: '/nodes',
    name: 'home-nodes',
    component: () => import('@/views/HomeNodes')
  },
  {
    path: '/nodes/:slug',
    name: 'node',
    component: () => import('@/views/Node'),
    props: true
  },
  {
    path: '/editor/:slug?',
    name: 'node-register',
    component: () => import('@/views/NodeRegister'),
    props: true
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
