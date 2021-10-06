<script setup lang="ts">
import { useOnline } from '@vueuse/core'
const isOnline = useOnline()
</script>

<template>
  <button>
    <icon-carbon:wifi v-if="isOnline" />
    <icon-carbon:wifi-off v-else />
  </button>
</template>
