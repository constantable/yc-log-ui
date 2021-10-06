<script setup lang="ts">
import { watch, computed, ComputedRef } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void,
}>()

const onToggleModal = (bool: boolean) => {
  emit('update:modelValue', bool)
}

const isShow: ComputedRef<boolean> = computed(() => props.modelValue)
watch(isShow, (val) => {
  if (val) {
    document.body.classList.add('is-blocked')
  } else {
    document.body.classList.remove('is-blocked')
  }
})
</script>

<template>
  <teleport to="body">
    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <section v-if="modelValue" class="fixed z-11 top-0 left-0 w-full h-full flex flex-col justify-center px-[16px] bg-[rgba(0,0,0,.9)]">
        <button class="absolute top-[14px] right-[14px] w-[32px] h-[32px] p-0 grid place-items-center text-[18px]" @click="onToggleModal(false)">
          <icon-carbon:close />
        </button>
        <div class="h-[60px]" />
        <slot />
        <div class="h-[60px]" />
      </section>
    </transition>
  </teleport>
</template>
