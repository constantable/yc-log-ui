import {createApp} from 'vue'
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

import 'vue-skeletor/dist/vue-skeletor.css'
import { Skeletor }  from 'vue-skeletor'

/**
 * init app
 */
createApp(App)
  .use(pinia)
  .use(i18n)
  .use(router)
  .use(head)
  .component(Skeletor.name, Skeletor)
  .directive('highlightjs', highlightjs)
  .mount('#app')
