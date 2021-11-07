<script setup lang="ts">
import {onMounted, Ref, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useAuth} from '@/stores/auth'
import axios from '@/plugins/axios'
import {errorToast} from '@/use/useToast'
import dayjs from 'dayjs'
import LocalizedFormat from 'dayjs/plugin/localizedFormat'

dayjs.locale('ru')
dayjs.extend(LocalizedFormat)

interface LogTime{
  nanos: number;
  seconds: number;
}
interface Log {
  ingested_at: LogTime;
  json_payload: any;
  level: number;
  message: string;
  resource: any;
  saved_at: LogTime;
  timestamp: LogTime;
  uid: string;
}

const { t } = useI18n()
const store = useAuth()
const levels: Ref<Array<string>> = ref(['', 'TRACE', 'DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL'])
const logs: Ref<Array<Log>> = ref([])
const isLoading: Ref<boolean> = ref(false)
const isLoadingBefore: Ref<boolean> = ref(false)
const isLoadingAfter: Ref<boolean> = ref(false)

const filter: Ref<string> = ref('')
const seconds: Ref<number> = ref(3600)
const levelSelected: Ref<string> = ref('')
const nextPageToken: Ref<string> = ref('')
const previousPageToken: Ref<string> = ref('')
const limit: Ref<number> = ref(50)

const getData = () => {
  isLoading.value = true

  axios.post(import.meta.env.VITE_BACKEND_URL + '/api/logs', {
    'levels': Number.isInteger(parseInt(levelSelected.value)) ? [parseInt(levelSelected.value)] : [],
    'filter': filter.value,
    'since': Math.floor(Date.now() / 1000) - seconds.value,
    'until': Math.floor(Date.now() / 1000),
    'limit': limit.value,
  }, store.getAuthHeaders)
    .then((response) => {
      // @ts-ignore
      logs.value = response.data.entries
      // @ts-ignore
      nextPageToken.value = response.data.next_page_token
      // @ts-ignore
      previousPageToken.value = response.data.previous_page_token
      isLoading.value = false
    })
    .catch((err) => {
      console.error('🦕', err)
      errorToast(`🦕 ${t('error')}`)
      isLoading.value = false
    })
}
const loadPrevious = () => {
  isLoadingBefore.value = true
  axios.post(import.meta.env.VITE_BACKEND_URL + '/api/logs/page', {
    'previous_page_token': previousPageToken.value,
  }, store.getAuthHeaders)
    .then((response) => {
      // @ts-ignore
      const entries = response.data.entries
      logs.value = entries.concat(logs.value)
      // @ts-ignore
      previousPageToken.value = response.data.previous_page_token
      isLoadingBefore.value = false
    })
    .catch((err) => {
      console.error('🦕', err)
      errorToast(`🦕 ${t('error')}`)
      isLoadingBefore.value = false
    })
}
const loadNext = () => {
  isLoadingAfter.value = true
  axios.post(import.meta.env.VITE_BACKEND_URL + '/api/logs/page', {
    'next_page_token': nextPageToken.value,
  }, store.getAuthHeaders)
    .then((response) => {
      // @ts-ignore
      const entries = response.data.entries
      logs.value = logs.value.concat(entries)
      // @ts-ignore
      nextPageToken.value = response.data.next_page_token
      isLoadingAfter.value = false
    })
    .catch((err) => {
      console.error('🦕', err)
      errorToast(`🦕 ${t('error')}`)
      isLoadingAfter.value = false
    })
}
const filterSource = (label: string, labelVal: string) => {
  filter.value = 'json_payload.kubernetes.labels.' + label + ': "' + labelVal + '"'
  getData()
}
const isOpenJsonModal: Ref<boolean> = ref(false)
const currentJson: Ref<any> = ref(null)
const handleJsonModal = (json: any) => {
  currentJson.value = json
  isOpenJsonModal.value = true
}
onMounted(() => {
  getData()
})
</script>

<template>
  <main class="relative">
    <header class="sticky top-0 left-0 z-1 w-full bg-$document border-t border-b border-$secondary px-[16px] py-[8px]">
      <h1 class="mb-[8px]">
        {{ t('filter') }}
      </h1>
      <div class="space-y-[8px]">
        <div>
          <label>
            <textarea
              v-model="filter"
              class="w-full rounded-[4px] bg-$secondary border-none p-[8px] focused"
              rows="2"
              @blur="getData()"
            />
          </label>
        </div>
        <div>
          <label class="inline-flex items-center min-h-[27px] md:min-h-[36px] space-x-[8px] rounded-[4px] bg-$secondary px-[8px] py-[4px] mr-[8px] md:mr-[16px] mb-[8px]">
            <span class="text-$typography-secondary">
              {{ t('level') }}:
            </span>
            <select
              v-model="levelSelected"
              class="flex-1 w-full bg-transparent"
              @change="getData()"
              @blur="getData()"
            >
              <option value="">Все</option>
              <option :value="1">TRACE</option>
              <option :value="2">DEBUG</option>
              <option :value="3">INFO</option>
              <option :value="4">WARN</option>
              <option :value="5">ERROR</option>
              <option :value="6">FATAL</option>
              <option :value="0">Нет</option>
            </select>
          </label>
          <label class="inline-flex items-center min-h-[27px]  md:min-h-[36px] space-x-[8px] rounded-[4px] bg-$secondary px-[8px] py-[4px] mr-[8px] md:mr-[16px] mb-[8px]">
            <span class="text-$typography-secondary">
              {{ t('limit') }}:
            </span>
            <select
              v-model="limit"
              class="flex-1 w-full bg-transparent"
              @change="getData()"
              @blur="getData()"
            >
              <option :value="50">50</option>
              <option :value="100">100</option>
            </select>
          </label>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
            :class="seconds === 3600?'active':''"
            @click.prevent="seconds = 3600;getData()"
          >
            1h
          </button>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
            :class="seconds === 3*3600?'active':''"
            @click.prevent="seconds = 3*3600;getData()"
          >
            3h
          </button>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
            :class="seconds === 24*3600?'active':''"
            @click.prevent="seconds = 24*3600;getData()"
          >
            1d
          </button>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
            :class="seconds === 7*24*3600?'active':''"
            @click.prevent="seconds = 7*24*3600;getData()"
          >
            1w
          </button>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
            :class="seconds === 2*7*24*3600?'active':''"
            @click.prevent="seconds = 2*7*24*3600;getData()"
          >
            2w
          </button>

          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
          >
            Другое
          </button>
        </div>
      </div>
      <div class="pt-[16px] hidden md:grid grid-cols-[150px,100px,200px,40px,1fr] gap-x-[10px]">
        <div>
          {{ t('time') }}
        </div>
        <div>
          {{ t('level') }}
        </div>
        <div>
          {{ t('label') }}
        </div>
        <div>
          Json
        </div>
        <div>
          {{ t('message') }}
        </div>
      </div>
    </header>
    <section class="w-full">
      <div v-if="!isLoading" class="px-[16px] py-[8px] text-center">
        <button class="button" @click="loadPrevious">
          {{ t('load-previous') }}
        </button>
      </div>
      <div v-if="isLoading">
        <ul>
          <li
            v-for="n in 10"
            :key="n"
            class="log border-t border-$secondary px-[16px] py-[8px]"
          >
            <div
              class="log--time"
            >
              <Skeletor />
            </div>
            <div
              class="log--level"
            >
              <Skeletor />
            </div>
            <div
              class="log--label"
            >
              <Skeletor />
            </div>
            <div
              class="log--json md:text-center"
            />
            <div
              class="log--msg"
            >
              <Skeletor />
            </div>
          </li>
        </ul>
      </div>
      <ul v-else>
        <li
          v-for="log in logs"
          :key="log.uid"
          class="log border-t border-$secondary px-[16px] py-[8px]"
        >
          <div
            class="log--time"
          >
            {{ dayjs.unix(log.timestamp.seconds).format('DD.MM.YYYY HH:mm:ss') }}
          </div>
          <div
            class="log--level"
          >
            <p :class="'level-' + levels[log.level]">
              {{ levels[log.level] }}
            </p>
          </div>
          <div
            class="log--label"
          >
            <template v-if="log.json_payload.kubernetes.labels.app">
              <a href="#" @click.prevent="filterSource('app',log.json_payload.kubernetes.labels.app)">app: {{ log.json_payload.kubernetes.labels.app }}</a>
            </template>
            <template v-else-if="log.json_payload.kubernetes.labels.run">
              <a href="#" @click.prevent="filterSource('run',log.json_payload.kubernetes.labels.run)">run: {{ log.json_payload.kubernetes.labels.run }}</a>
            </template>
            <template v-else-if="log.json_payload.kubernetes.labels.job">
              <a href="#" @click.prevent="filterSource('job',log.json_payload.kubernetes.labels.job)">job: {{ log.json_payload.kubernetes.labels.job }}</a>
            </template>
          </div>
          <div
            class="log--json md:text-center"
          >
            <button
              class="button"
              @click="handleJsonModal(log.json_payload)"
            >
              <icon-carbon-view />
            </button>
          </div>
          <div
            class="log--msg"
          >
            {{ log.message }}
          </div>
        </li>
      </ul>
      <div v-if="!isLoading" class="px-[16px] py-[8px] text-center">
        <button class="button" @click="loadNext">
          {{ t('load-next') }}
        </button>
      </div>
    </section>
  </main>

  <ModalDefault v-model="isOpenJsonModal">
    <section class="relative rounded-[4px] overflow-scroll max-h-[75%]">
      <CodeBlock lang="json">
        {{ currentJson }}
      </CodeBlock>
    </section>
  </ModalDefault>
</template>
