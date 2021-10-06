import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

/**
 * Plugins
 */
import head from '@/plugins/vueuse-head'
import i18n from '@/plugins/vue-i18n'
import pinia from '@/plugins/pinia'

/**
 * Directives
 */
import highlightjs from '@/directives/v-highlightjs'

/**
 * Styles
 */
import 'highlight.js/styles/github-dark-dimmed.css'
import 'mosha-vue-toastify/dist/style.css'
import 'virtual:windi.css' // windicss demon
import '@/assets/sass/main.sass'

/**
 * init app
 */
createApp(App)
  .use(pinia)
  .use(i18n)
  .use(router)
  .use(head)
  .directive('highlightjs', highlightjs)
  .mount('#app')

/**
 * Middleware
 */
import { useAuth }  from '@/stores/auth'
const store = useAuth()
router.beforeEach((to, from, next) => {
  if (to.meta.auth && !store.isAuthenticated) {
    next({ name: 'Login' })
  } else if (to.name === 'Login' && store.isAuthenticated) {
    next({ name: 'Home' })
  }
  next()
})
