<script setup lang="ts">
import { ref, Ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuth, IAuthForm } from '@/stores/auth'
import { toast, errorToast } from '@/use/useToast'

const { t } = useI18n()
const router = useRouter()
const store = useAuth()
const { login } = store

const form: Ref<IAuthForm> = ref({
  username: '',
  password: '',
})

const handleSubmit = async(): Promise<void> => {
  console.log('🦕 handleSubmit')
  try {
    await login(form.value)
    router.push({ name: 'Home' })
    toast(t('success'), 'success')
  } catch (err) {
    console.error('🦕', err)
    errorToast(`🦕 ${t('error')}`)
  }
}
</script>

<template>
  <main class="px-[16px]">
    <section class="flex flex-col items-center justify-center py-[20px] md:py-[40px]">
      <h1 class="mb-[24px] md:mb-[40px] text-[18px] font-medium">
        {{ t('login') }}
      </h1>
      <form
        class="w-full max-w-[320px] space-y-[16px]"
        @submit.prevent="handleSubmit"
      >
        <FormMaterialInput
          v-model="form.username"
          placeholder="Username"
          required
        />
        <FormMaterialInput
          v-model="form.password"
          placeholder="Password"
          type="password"
          required
        />
        <div class="text-center">
          <button
            class="button"
          >
            <span class="px-[24px]">
              {{ t('submit') }}
            </span>
          </button>
        </div>
      </form>
    </section>
  </main>
</template>
