<script setup lang="ts">
import { isDark, toggleDark } from '@/plugins/dark-mode'
</script>

<template>
  <button
    @click="toggleDark()"
  >
    <icon-carbon-moon v-if="isDark" />
    <icon-carbon-sun v-else />
  </button>
</template>
