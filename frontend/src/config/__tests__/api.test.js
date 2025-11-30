import { describe, it, expect } from 'vitest'
import { API_ENDPOINTS, PARKING_TYPES, BOOKING_STATUSES, USER_ROLES } from '../api'

describe('API Configuration', () => {
  it('has AUTH endpoints', () => {
    expect(API_ENDPOINTS.AUTH).toBeDefined()
    expect(API_ENDPOINTS.AUTH.BASE).toBeDefined()
    expect(API_ENDPOINTS.AUTH.LOGIN).toBe('/auth/login')
    expect(API_ENDPOINTS.AUTH.REGISTER).toBe('/auth/register')
  })

  it('has PARKING endpoints', () => {
    expect(API_ENDPOINTS.PARKING).toBeDefined()
    expect(API_ENDPOINTS.PARKING.LIST).toBe('/parking')
    expect(API_ENDPOINTS.PARKING.DETAIL(1)).toBe('/parking/1')
    expect(API_ENDPOINTS.PARKING.CREATE).toBe('/parking')
  })

  it('has BOOKING endpoints', () => {
    expect(API_ENDPOINTS.BOOKING).toBeDefined()
    expect(API_ENDPOINTS.BOOKING.LIST).toBe('/booking')
    expect(API_ENDPOINTS.BOOKING.DETAIL(1)).toBe('/booking/1')
    expect(API_ENDPOINTS.BOOKING.CREATE).toBe('/booking')
  })

  it('has PARKING_TYPES constants', () => {
    expect(PARKING_TYPES.OUTDOOR).toBe('outdoor')
    expect(PARKING_TYPES.COVERED).toBe('covered')
    expect(PARKING_TYPES.UNDERGROUND).toBe('underground')
    expect(PARKING_TYPES.MULTI_LEVEL).toBe('multi-level')
  })

  it('has BOOKING_STATUSES constants', () => {
    expect(BOOKING_STATUSES.WAITING).toBe('Waiting')
    expect(BOOKING_STATUSES.CONFIRMED).toBe('Confirmed')
    expect(BOOKING_STATUSES.CANCELED).toBe('Canceled')
  })

  it('has USER_ROLES constants', () => {
    expect(USER_ROLES.DRIVER).toBe('driver')
    expect(USER_ROLES.OWNER).toBe('owner')
    expect(USER_ROLES.ADMIN).toBe('admin')
  })
})
