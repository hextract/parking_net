import axios from 'axios'
import { API_ENDPOINTS } from '../config/api'
import { getAuthToken, removeAuthToken } from '../utils/cookies'

const createApiClient = (baseURL) => {
  const client = axios.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  client.interceptors.request.use(
    (config) => {
      const token = getAuthToken()
      if (token) {
        config.headers['api_key'] = token
      }
      return config
    },
    (error) => {
      return Promise.reject(error)
    }
  )

  client.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response) {
        if (error.response.status === 401) {
          removeAuthToken()
          if (window.location.pathname !== '/login') {
            window.location.href = '/login'
          }
        }
        const errorMessage = error.response.data?.error_message ||
                            error.response.data?.message ||
                            error.response.data?.error ||
                            `Request failed with status ${error.response.status}`
        return Promise.reject({
          message: errorMessage,
          status: error.response.status,
          data: error.response.data,
          response: error.response
        })
      }
      return Promise.reject({
        message: error.message || 'Network error',
        status: 0,
        data: null
      })
    }
  )

  return client
}

export const authApi = createApiClient(API_ENDPOINTS.AUTH.BASE)
export const parkingApi = createApiClient(API_ENDPOINTS.PARKING.BASE)
export const bookingApi = createApiClient(API_ENDPOINTS.BOOKING.BASE)
export const paymentApi = createApiClient(API_ENDPOINTS.PAYMENT.BASE)
