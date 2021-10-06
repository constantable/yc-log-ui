/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import Home from '@/views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    meta: {
      auth: true,
    },
    component: Home,
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('@/views/Error404.vue'),
  },
]

// console.log('🚥 routes', routes)

export default routes
