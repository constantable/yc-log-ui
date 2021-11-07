import {
  createRouter,
  createWebHistory,
  // createWebHashHistory,
} from 'vue-router'
import routes from './routes'
import {useAuth} from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  // history: createWebHashHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  if (to.meta.auth && !useAuth().isAuthenticated) {
    next({ name: 'Login' })
    return
  } else if (to.name === 'Login' && useAuth().isAuthenticated) {
    next({ name: 'Home' })
    return
  }
  next()
})

export default router

