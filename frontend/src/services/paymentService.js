import { paymentApi } from './api'

/**
 * Payment Service
 * Handles all payment-related API calls
 */

// Get user balance
export const getBalance = async () => {
  const response = await paymentApi.get('/payment/balance')
  return response.data
}

// Get user transactions
export const getTransactions = async (limit = 50, offset = 0) => {
  const response = await paymentApi.get('/payment/transactions', {
    params: { limit, offset }
  })
  return response.data
}

// Activate promocode
export const activatePromocode = async (code) => {
  const response = await paymentApi.post('/payment/promocode/activate', { code })
  return response.data
}

// Generate promocode from balance
export const generatePromocode = async (amount) => {
  const response = await paymentApi.post('/payment/promocode/generate', { amount })
  return response.data
}

// Create promocode (admin only)
export const createPromocode = async (data) => {
  const response = await paymentApi.post('/payment/promocode/create', data)
  return response.data
}

// Get promocode info
export const getPromocodeInfo = async (code) => {
  const response = await paymentApi.get(`/payment/promocode/${code}`)
  return response.data
}
