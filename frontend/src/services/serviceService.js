import { parkingApi } from './api'
import { API_ENDPOINTS } from '../config/api'

export const serviceService = {
  getServices: async (parkingId) => {
    const response = await parkingApi.get(API_ENDPOINTS.PARKING.SERVICES.LIST(parkingId))
    return response.data
  },

  getServiceById: async (parkingId, serviceId) => {
    const response = await parkingApi.get(API_ENDPOINTS.PARKING.SERVICES.DETAIL(parkingId, serviceId))
    return response.data
  },

  createService: async (parkingId, serviceData) => {
    const response = await parkingApi.post(API_ENDPOINTS.PARKING.SERVICES.CREATE(parkingId), serviceData)
    return response.data
  },

  updateService: async (parkingId, serviceId, serviceData) => {
    const response = await parkingApi.put(API_ENDPOINTS.PARKING.SERVICES.UPDATE(parkingId, serviceId), serviceData)
    return response.data
  },

  deleteService: async (parkingId, serviceId) => {
    const response = await parkingApi.delete(API_ENDPOINTS.PARKING.SERVICES.DELETE(parkingId, serviceId))
    return response.data
  },
}

