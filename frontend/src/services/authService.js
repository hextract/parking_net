import { authApi } from './api'
import { API_ENDPOINTS } from '../config/api'
import { setAuthToken, removeAuthToken, getAuthToken } from '../utils/cookies'

export const authService = {
  login: async (credentials) => {
    const response = await authApi.post(API_ENDPOINTS.AUTH.LOGIN, credentials)
    return response.data
  },

  register: async (userData) => {
    const response = await authApi.post(API_ENDPOINTS.AUTH.REGISTER, userData)
    return response.data
  },

  changePassword: async (passwordData) => {
    const response = await authApi.post(API_ENDPOINTS.AUTH.CHANGE_PASSWORD, passwordData)
    return response.data
  },

  getUserInfo: async () => {
    try {
      const response = await authApi.get(API_ENDPOINTS.AUTH.ME || '/auth/me')
      return response.data
    } catch (error) {
      return null
    }
  },

  logout: () => {
    removeAuthToken()
  },

  getStoredToken: () => {
    return getAuthToken()
  },

  setAuthData: (token) => {
    setAuthToken(token)
  },
}
