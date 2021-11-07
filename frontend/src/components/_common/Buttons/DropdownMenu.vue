<template>
  <div
    class="dropdown"
    :class="{ dropup: top }"
    @mouseleave="mouseLeave"
    @mouseover="mouseOver"
    @mouseenter="mouseEnter"
    @click="toggleMenu"
  >
    <slot />
    <transition :name="transition">
      <div
        v-show="modelValue"
        ref="dropdown"
        class="dropdown-menu show"
        :class="{ 'dropdown-menu-right': right}"
        :style="styles"
        @mouseleave="startTimer"
        @mouseenter="stopTimer"
        @click.stop
      >
        <slot name="dropdown" />
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  props: {
    modelValue: Boolean,
    right: Boolean,
    hover: Boolean,
    hoverTime: {
      type: Number,
      default: 100,
    },
    hoverTimeout: {
      type: Number,
      default: 500,
    },
    styles: {
      type: Object,
      default() {
        return {}
      },
    },
    interactive: { //whether should stay open until clicked outside
      type: Boolean,
      default: false,
    },
    transition: {
      type: String,
      default: '',
    },
    closeOnClickOutside: {
      type: Boolean,
      default: true,
    },
  },
  emits: ['update:modelValue'],
  data() {
    return {
      hovering: false,
      top: false,
    }
  },
  watch: {
    modelValue(v) {
      if (v) {
        let vm = this
        this.top = false
        this.$nextTick(() => {
            let rect = vm.$refs.dropdown.getBoundingClientRect()
            let window_height = (window.innerHeight || document.documentElement.clientHeight)
            this.top = (rect.bottom > window_height) && (rect.top >= rect.height)
          },
        )
      }
    },
    interactive: {
      handler(value) {
        if (typeof document === 'object')
          value ? document.body.addEventListener('click', this.closeMenu) : document.body.removeEventListener('click', this.closeMenu)
      },
      immediate: true,
    },
  },
  unmounted() {
    document.body.removeEventListener('click', this.closeMenu)
  },
  methods: {
    mouseEnter() {
      // console.log('mouseEnter', $event.target)
      this.stopTimer()
      if (this.hover && this.hoverTime > 0 && !this.modelValue) {
        // console.log('start open timer', this.hoverTime)
        this.hoverOpenTimer = setTimeout(() => {
          this.$emit('update:modelValue', true)

          //disable for a moment
          this.hovering = true
          setTimeout(() => {
            this.hovering = false
          }, this.hoverTimeout)
        }, this.hoverTime)
      }

      if (this.hover && !this.modelValue && this.hoverTime === 0) {
        this.hovering = true
        setTimeout(() => {
          this.hovering = false
        }, this.hoverTimeout)
        this.$emit('update:modelValue', true)
      }

    },
    mouseLeave() {
      // console.log('mouseLeave', $event.target)
      if (!this.hoverTimer) { //left the link and no time active
        this.startTimer()
      }

      if (this.hoverTime > 0 && this.hover) {
        // console.log('clear hover timer')
        clearTimeout(this.hoverOpenTimer)
      }
    },
    mouseOver() {
      this.stopTimer()
      // console.log('mouseOver')
    },
    closeMenu($event) {
      if (!$event || !this.$el.contains($event.target)) {
        if (this.modelValue && this.closeOnClickOutside) {
          this.$emit('update:modelValue', false)
        }
      }
    },
    toggleMenu() {
      if (this.hovering)
        return
      if (this.modelValue && this.hover)
        return

      this.$emit('update:modelValue', !this.modelValue)
    },
    stopTimer() {
      // console.log('stop timer')
      clearTimeout(this.hoverTimer)
      this.hoverTimer = null
    },
    startTimer() {
      // console.log('start timer')
      if (!this.interactive)
        this.hoverTimer = setTimeout(this.closeMenu, this.hoverTimeout)
    },
  },
}
</script>
