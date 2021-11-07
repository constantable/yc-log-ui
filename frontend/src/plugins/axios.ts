import axios from 'axios'
import {useAuth} from '@/stores/auth'
import createAuthRefreshInterceptor from 'axios-auth-refresh'

if (import.meta.env.NODE_ENV === 'development') {
  console.log('🦕 VITE ENV', import.meta.env)
}

const axiosInstance = axios.create({
  headers: {
    Accept: 'application/json',
  },
})
interface ITokenRefreshResponse {
  data: ITokenRefreshResponseData
}
interface ITokenRefreshResponseData {
  access_token: string,
  refresh_token: string
}
const refreshAuthLogic = (failedRequest: any) => axiosInstance.post<any, ITokenRefreshResponse>(import.meta.env.VITE_BACKEND_URL + '/api/token', {
  refresh_token: useAuth().refresh_token,
}).then((tokenRefreshResponse: ITokenRefreshResponse) => {
  const data:ITokenRefreshResponseData = tokenRefreshResponse.data
  useAuth().setState(data)
  failedRequest.response.config.headers['Authorization'] = 'Bearer ' + data.access_token
  return Promise.resolve()
})

createAuthRefreshInterceptor(axiosInstance, refreshAuthLogic)

export default axiosInstance
