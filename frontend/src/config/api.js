import { API_BASE_URL } from './env'

export const API_ENDPOINTS = {
  AUTH: {
    BASE: API_BASE_URL,
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    CHANGE_PASSWORD: '/auth/change-password',
    ME: '/auth/me',
  },
  PARKING: {
    BASE: API_BASE_URL,
    LIST: '/parking',
    DETAIL: (id) => `/parking/${id}`,
    CREATE: '/parking',
    UPDATE: (id) => `/parking/${id}`,
    DELETE: (id) => `/parking/${id}`,
  },
  BOOKING: {
    BASE: API_BASE_URL,
    LIST: '/booking',
    DETAIL: (id) => `/booking/${id}`,
    CREATE: '/booking',
    UPDATE: (id) => `/booking/${id}`,
    DELETE: (id) => `/booking/${id}`,
  },
  PAYMENT: {
    BASE: API_BASE_URL,
    BALANCE: '/payment/balance',
    TRANSACTIONS: '/payment/transactions',
    PROMOCODE_ACTIVATE: '/payment/promocode/activate',
    PROMOCODE_GENERATE: '/payment/promocode/generate',
    PROMOCODE_CREATE: '/payment/promocode/create',
    PROMOCODE_INFO: (code) => `/payment/promocode/${code}`,
  },
}

export const PARKING_TYPES = {
  OUTDOOR: 'outdoor',
  COVERED: 'covered',
  UNDERGROUND: 'underground',
  MULTI_LEVEL: 'multi-level',
}

export const BOOKING_STATUSES = {
  WAITING: 'Waiting',
  CONFIRMED: 'Confirmed',
  CANCELED: 'Canceled',
}

export const USER_ROLES = {
  DRIVER: 'driver',
  OWNER: 'owner',
  ADMIN: 'admin',
}
