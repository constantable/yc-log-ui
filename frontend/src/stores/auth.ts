import { defineStore } from 'pinia'
import axios from '@/plugins/axios'
import jwt_decode from 'jwt-decode'

const REFRESH_LIFETIME = 3600000*24*7

export interface IAuthForm {
  username: string
  password: string
}

export interface IAuthState {
  access_token: string
  refresh_token: string
  last_refresh: number | null,
}
interface ITokenData{
  admin: boolean
  exp: number
  username: string
}
const checkValidRefresh = (milliseconds: number | null): boolean => {
  if (!milliseconds) return false
  const now = Date.now()
  return milliseconds + REFRESH_LIFETIME > now
}

export const useAuth = defineStore('auth', {
  state: (): IAuthState => ({
    access_token: localStorage.getItem('access_token') || '',
    refresh_token: localStorage.getItem('refresh_token') || '',
    last_refresh: localStorage.getItem('last_refresh')
      ? Number(localStorage.getItem('last_refresh'))
      : null,
  }),
  getters: {
    isAuthenticated(): boolean {
      return this.access_token.trim().length > 0 && checkValidRefresh(this.last_refresh)
    },
    isAdmin(): boolean {
      const decoded :ITokenData = jwt_decode(this.access_token.trim())
      return decoded.admin
    },
    getAuthHeaders(): any {
      return {
        headers: { Authorization: `Bearer ${this.access_token}` },
      }
    },
    getRefreshToken(): string {
      return this.refresh_token
    },
  },
  actions: {
    setState(data: any): void {
      const { access_token, refresh_token } = data
      this.access_token = access_token
      this.refresh_token = refresh_token
      this.last_refresh = Date.now()
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)
      localStorage.setItem('last_refresh', String(this.last_refresh))
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
    logout(): void {
      this.access_token = ''
      this.refresh_token = ''
      this.last_refresh = 0
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('last_refresh')
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
