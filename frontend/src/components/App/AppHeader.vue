<script setup lang="ts">
import {computed, ComputedRef, ref, Ref} from 'vue'
import {useAuth} from '@/stores/auth'
const store = useAuth()
const isAdmin: Ref<boolean> = ref(true)
</script>
<template>
  <header class="relative w-full flex justify-between items-center px-[16px] py-[8px] space-x-[16px]">
    <div class="flex-1 flex flex-col space-y-[4px]">
      <p class="font-medium text-[18px]">
        YC Log UI
      </p>
    </div>
    <div class="flex items-center space-x-[8px]">
      <ProfileButton
        v-if="store.isAuthenticated"
        :is-admin="store.isAdmin"
      />
      <LocaleButton />
      <DarkModeButton />
    </div>
  </header>
</template>
