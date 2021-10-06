import { defineStore } from 'pinia'
import axios from '@/plugins/axios'

const MIN15 = 900000 // 900000 milliseconds === 15 min

export interface IAuthForm {
  username: string
  password: string
}

export interface IAuthState {
  access_token: string
  refresh_token: string
  last_refresh: number | null,
}

const check15Min = (milliseconds: number | null): boolean => {
  if (!milliseconds) return false
  const now = Date.now()
  return milliseconds + MIN15 > now //
}

export const useAuth = defineStore('auth', {
  state: (): IAuthState => ({
    access_token: '',
    refresh_token:  localStorage.getItem('refresh_token') || '',
    last_refresh: localStorage.getItem('last_refresh')
      ? Number(localStorage.getItem('last_refresh'))
      : null,
  }),
  getters: {
    isAuthenticated(): boolean {
      return this.refresh_token.trim().length > 0 && check15Min(this.last_refresh)
    },
    getAuthHeaders(): any {
      return {
        headers: { Authorization: `Bearer ${this.refresh_token}` },
      }
    },
  },
  actions: {
    setState(data: any): void {
      const { access_token, refresh_token } = data
      this.access_token = access_token
      this.refresh_token = refresh_token
      this.last_refresh = Date.now()
      localStorage.setItem('refresh_token', refresh_token)
      localStorage.setItem('last_refresh', String(this.last_refresh))
      setTimeout(() => {
        this.refresh()
      }, MIN15)
    },
    login(data: IAuthForm): Promise<any> {
      return new Promise((resolve, reject) => {
        axios.post(import.meta.env.VITE_BACKEND_URL + '/api/login', data)
          .then(response => {
            this.setState(response.data)
            resolve(response)
          })
          .catch(err => reject(err))
      })
    },
    refresh(): Promise<any> {
      return new Promise((resolve, reject) => {
        axios.post(import.meta.env.VITE_BACKEND_URL + '/api/token', {
          refresh_token: this.refresh_token,
        })
          .then(response => {
            this.setState(response.data)
            resolve(response)
          })
          .catch(err => reject(err))
      })
    },
  },
})
