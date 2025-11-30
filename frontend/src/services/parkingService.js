import { parkingApi } from './api'
import { API_ENDPOINTS } from '../config/api'

export const parkingService = {
  getParkings: async (filters = {}) => {
    const params = new URLSearchParams()
    if (filters.city) params.append('city', filters.city)
    if (filters.name) params.append('name', filters.name)
    if (filters.parking_type) params.append('parking_type', filters.parking_type)
    if (filters.owner_id) params.append('owner_id', filters.owner_id)

    const response = await parkingApi.get(
      `${API_ENDPOINTS.PARKING.LIST}${params.toString() ? `?${params.toString()}` : ''}`
    )
    return response.data
  },

  getParkingById: async (id) => {
    const response = await parkingApi.get(API_ENDPOINTS.PARKING.DETAIL(id))
    return response.data
  },

  createParking: async (parkingData) => {
    const response = await parkingApi.post(API_ENDPOINTS.PARKING.CREATE, parkingData)
    return response.data
  },

  updateParking: async (id, parkingData) => {
    const response = await parkingApi.put(API_ENDPOINTS.PARKING.UPDATE(id), parkingData)
    return response.data
  },

  deleteParking: async (id) => {
    const response = await parkingApi.delete(API_ENDPOINTS.PARKING.DELETE(id))
    return response.data
  },
}
