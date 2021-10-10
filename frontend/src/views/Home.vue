<script setup lang="ts">
import {onMounted, Ref, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useAuth} from '@/stores/auth'
import axios from '@/plugins/axios'
import {errorToast} from '@/use/useToast'
import dayjs from 'dayjs'

import LocalizedFormat from 'dayjs/plugin/localizedFormat';
dayjs.locale('ru');
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
const levels: Ref<Array<string>> = ref(["", "TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"])
const logs: Ref<Array<Log>> = ref([])
const isLoading: Ref<boolean> = ref(false)

const filter: Ref<string> = ref("")
const limit: Ref<number> = ref(100)

const getData = async() => {
  console.log('🦕 getData')
  isLoading.value = true
  try {
    const response = await axios.post(import.meta.env.VITE_BACKEND_URL + '/api/logs', {
      'levels': [],
      'filter': filter.value,
      'since': 1632608612,
      'until': 1634612212,
      'limit': limit.value,
    }, store.getAuthHeaders)
    // @ts-ignore
    logs.value = response.data.entries
    isLoading.value = false
  } catch (err) {
    console.error('🦕', err)
    errorToast(`🦕 ${t('error')}`)
    isLoading.value = false
  }
}

const isOpenJsonModal: Ref<boolean> = ref(false)
const currentJson: Ref<any> = ref(null)
const handleJsonModal = (json: any) => {
  currentJson.value = json
  isOpenJsonModal.value = true
}
onMounted(async() => {
  await getData()
})
</script>

<template>
  <main class="relative">
    <header class="sticky top-0 left-0 z-1 w-full bg-$document border-t border-b border-$secondary px-[16px] py-[8px]">
      <h1 class="mb-[8px]">
        {{ t('logs') }}
      </h1>
      <div class="space-y-[8px]">
        <div>
          <label>
            <textarea
              rows="2"
              v-model="filter"
              class="w-full rounded-[4px] bg-$secondary border-none p-[8px] focused"
            />
          </label>
        </div>
        <div>
          <label class="inline-flex items-center min-h-[27px] md:min-h-[36px] space-x-[8px] rounded-[4px] bg-$secondary px-[8px] py-[4px] mr-[8px] md:mr-[16px] mb-[8px]">
            <span class="text-$typography-secondary">
              {{ t('level') }}:
            </span>
            <select
              class="flex-1 w-full bg-transparent"
            >
              <option>Все</option>
            </select>
          </label>
          <label class="inline-flex items-center min-h-[27px]  md:min-h-[36px] space-x-[8px] rounded-[4px] bg-$secondary px-[8px] py-[4px] mr-[8px] md:mr-[16px] mb-[8px]">
            <span class="text-$typography-secondary">
              Кол-во:
            </span>
            <select
              class="flex-1 w-full bg-transparent"
            >
              <option>50</option>
              <option>100</option>
            </select>
          </label>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
          >
            1h
          </button>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
          >
            3h
          </button>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
          >
            1d
          </button>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
          >
            1w
          </button>
          <button
            class="button mr-[8px] md:mr-[16px] mb-[8px]"
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
      <div class="px-[16px] py-[8px] text-center">
        <button class="button">
          Загрузить предыдущие записи
        </button>
      </div>
      <div v-if="isLoading">
        <p>isLoading</p>
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
            app: {{ log.json_payload.kubernetes.labels.app}}
            </template>
            <template v-else-if="log.json_payload.kubernetes.labels.app">
              run: {{ log.json_payload.kubernetes.labels.run}}
            </template>
            <template v-else-if="log.json_payload.kubernetes.labels.app">
              job: {{ log.json_payload.kubernetes.labels.job}}
            </template>
          </div>
          <div
            class="log--json md:text-center"
          >
            <button
              class="button"
              @click="handleJsonModal(log.json_payload)"
            >
              <icon-carbon:view />
            </button>
          </div>
          <div
            class="log--msg"
          >
            {{ log.message }}
          </div>
        </li>
      </ul>
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
