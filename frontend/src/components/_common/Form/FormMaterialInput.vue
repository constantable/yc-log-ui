<script setup lang="ts" >
import { ref, Ref, computed, ComputedRef } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: '',
  },
  id: {
    type: String,
    default: Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15),
  },
  type: {
    type: String,
    default: 'text',
  },
  placeholder: {
    type: String,
    default: '',
  },
  autocomplete: {
    type: String,
    default: 'off',
  },
  step: {
    type: String,
    default: undefined,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  required: {
    type: Boolean,
    default: undefined,
  },
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void,
}>()

const isFocus: Ref<boolean> = ref(false)

const valueIsNotEmpty: ComputedRef<boolean> = computed(() =>
  props.modelValue
    ? props.modelValue.length > 0
    : false,
)

const onFocus = (): void => { isFocus.value = true }
const onBlur = (): void => { isFocus.value = false }
const onInput = (ev: any): void => { emit('update:modelValue', ev.target.value) }
</script>

<template>
  <label
    :for="id"
    class="relative inline-block w-full rounded-[4px] bg-$secondary"
    :style="{
      boxShadow: isFocus && !disabled ? '0 0 0 .2rem var(--primary)' : 'none',
    }"
  >
    <input
      :id="id"
      :type="type"
      :autocomplete="autocomplete"
      :value="modelValue"
      :disabled="disabled"
      :required="required"
      :step="step"
      placeholder=""
      class="clear w-full min-h-[34px] px-[10px] pt-[16px] pb-[2px] text-[16px] focus:outline-none bg-transparent"
      @focus="onFocus"
      @blur="onBlur"
      @input="onInput"
    >
    <span
      class="absolute inline-block top-0 left-[10px] text-[16px] transform text-$typography-secondary"
      :class="isFocus || valueIsNotEmpty
        ? 'translate-y-[2px] scale-70'
        : 'translate-y-[10px] scale-100'
      "
      :style="{
        transition: '.35s cubic-bezier(0.175,0.885,0.32,1)',
        transformOrigin: 'top left',
      }"
    >
      {{ placeholder }}
    </span>
  </label>
</template>



