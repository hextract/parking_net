import { bookingApi } from './api'
import { API_ENDPOINTS } from '../config/api'

export const bookingService = {
  getBookings: async (filters = {}) => {
    const params = new URLSearchParams()
    if (filters.parking_place_id) params.append('parking_place_id', filters.parking_place_id)
    if (filters.user_id) params.append('user_id', filters.user_id)

    const response = await bookingApi.get(
      `${API_ENDPOINTS.BOOKING.LIST}${params.toString() ? `?${params.toString()}` : ''}`
    )
    return response.data
  },

  getBookingById: async (id) => {
    const response = await bookingApi.get(API_ENDPOINTS.BOOKING.DETAIL(id))
    return response.data
  },

  createBooking: async (bookingData) => {
    const response = await bookingApi.post(API_ENDPOINTS.BOOKING.CREATE, bookingData)
    return response.data
  },

  updateBooking: async (id, bookingData) => {
    const response = await bookingApi.put(API_ENDPOINTS.BOOKING.UPDATE(id), bookingData)
    return response.data
  },

  deleteBooking: async (id) => {
    const response = await bookingApi.delete(API_ENDPOINTS.BOOKING.DELETE(id))
    return response.data
  },
}
