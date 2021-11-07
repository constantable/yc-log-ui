<script setup lang="ts">
import DropdownMenu from '@/components/_common/Buttons/DropdownMenu.vue'
import {ref, Ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useAuth} from '@/stores/auth'
const { t } = useI18n()
const store = useAuth()
const { logout } = store
const props = defineProps({
  isAdmin: {
    type: Boolean,
    default: false,
  },
})
const show: Ref<boolean> = ref(false)
</script>

<template>
  <DropdownMenu
    v-model="show"
    :interactive="true"
  >
    <button
      class="button dropdown-toggle"
    >
      <icon-carbon-user/>
    </button>
    <template #dropdown>
      <template v-if="props.isAdmin">
        <a class="dropdown-item" href="#">{{ t('users') }}</a>
      </template>
      <a class="dropdown-item" href="#">{{ t('change-password') }}</a>
      <a class="dropdown-item" href="#" @click="logout">{{ t('logout') }}</a>
    </template>
  </DropdownMenu>
</template>

