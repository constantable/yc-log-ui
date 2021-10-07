<script setup lang="ts">
import { ref, Ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuth } from '@/stores/auth'
import axios from '@/plugins/axios'
import { errorToast } from '@/use/useToast'

const { t } = useI18n()
const store = useAuth()
const levels: Ref<Array<string>> = ref(["", "TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"])
const logs: Ref<Array<any>> = ref([])
const isLoading: Ref<boolean> = ref(false)
const getData = async() => {
  console.log('🦕 getData')
  isLoading.value = true
  try {
    const response = await axios.post(import.meta.env.VITE_BACKEND_URL + '/api/logs', {
      'levels': [],
      'filter': '',
      'since': 1632608612,
      'until': 1634612212,
      'limit': 100,
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
              class="w-full rounded-[4px] bg-$secondary border-none p-[8px] focused"
            />
          </label>
        </div>
        <div>
          <label class="inline-flex items-center min-h-[27px] md:min-h-[36px] space-x-[8px] rounded-[4px] bg-$secondary px-[8px] py-[4px] mr-[8px] md:mr-[16px] mb-[8px]">
            <p class="text-$typography-secondary">
              {{ t('level') }}:
            </p>
            <select
              class="flex-1 w-full bg-transparent"
            >
              <option>Все</option>
            </select>
          </label>
          <label class="inline-flex items-center min-h-[27px]  md:min-h-[36px] space-x-[8px] rounded-[4px] bg-$secondary px-[8px] py-[4px] mr-[8px] md:mr-[16px] mb-[8px]">
            <p class="text-$typography-secondary">
              Кол-во:
            </p>
            <select
              class="flex-1 w-full bg-transparent"
            >
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
      <div class="pt-[16px] hidden md:grid grid-cols-[150px,100px,40px,1fr] gap-x-[16px]">
        <div>
          {{ t('time') }}
        </div>
        <div>
          {{ t('level') }}
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
          Загрузить 100 предыдущих записей
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
            :title="t('time')"
            class="log--time"
          >
            {{ log.timestamp.seconds }}
          </div>
          <div
            :title="t('level')"
            class="log--level"
          >
            <p :class="'level-' + levels[log.level]">
              {{ levels[log.level] }}
            </p>
          </div>
          <div
            title="Json"
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
            :title="t('message')"
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
