import { createToast } from 'mosha-vue-toastify'

type TToast = 'default' | 'info' | 'warning' | 'success' | 'danger'

const timeout = 4000
// const timeout = 999999999

export const toast = (msg: string, type: TToast = 'default'): void => {
  createToast(msg, {
    type,
    position: 'bottom-center',
    transition: 'slide',
    timeout,
  })
}

export const errorToast = (title: string, description: undefined | string = undefined): void => {
  createToast({
    title,
    description,
  },
  {
    type: 'danger',
    position: 'bottom-center',
    transition: 'slide',
    timeout,
  })
}
