import { defineStore } from 'pinia'
import axios from '@/plugins/axios'

export interface IAuthForm {
  username: string
  password: string
}

export interface IAuthState {
  access_token: string
  refresh_token: string
}

export const useAuth = defineStore('auth', {
  state: (): IAuthState => ({
    access_token: '',
    refresh_token: '',
  }),
  getters: {
    isAuthenticated(): boolean {
      return this.access_token.trim().length > 0
    },
  },
  actions: {
    login(data: IAuthForm): Promise<any> {
      return new Promise((resolve, reject) => {
        axios.post(import.meta.env.VITE_BACKEND_URL + '/api/login', data)
          .then(response => {
            // @ts-ignore
            const { access_token, refresh_token } = response.data
            this.access_token = access_token
            this.refresh_token = refresh_token
            resolve(response)
          })
          .catch(err => reject(err))
      })
    },
  },
})
