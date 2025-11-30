import { describe, it, expect } from 'vitest'
import { formatDateTime, formatDate, formatTime } from '../dateUtils'

describe('dateUtils', () => {
  describe('formatDateTime', () => {
    it('should format ISO date string to localized date and time', () => {
      const isoString = '2025-11-30T15:42:00.000Z'
      const result = formatDateTime(isoString, 'en-US')

      // Should include date components (checking for some part of the date)
      expect(result).toMatch(/2025/)
      expect(result).toMatch(/30/)
    })

    it('should handle different locales', () => {
      const isoString = '2025-11-30T15:42:00.000Z'
      const enResult = formatDateTime(isoString, 'en-US')
      const ruResult = formatDateTime(isoString, 'ru-RU')

      // Both should contain the year
      expect(enResult).toMatch(/2025/)
      expect(ruResult).toMatch(/2025/)

      // Results should be different due to locale
      expect(enResult).not.toBe(ruResult)
    })

    it('should return empty string for empty input', () => {
      expect(formatDateTime('')).toBe('')
      expect(formatDateTime(null)).toBe('')
      expect(formatDateTime(undefined)).toBe('')
    })

    it('should handle invalid date gracefully', () => {
      const result = formatDateTime('invalid-date')
      expect(result).toBe('invalid-date')
    })

    it('should use en-US as default locale', () => {
      const isoString = '2025-11-30T15:42:00.000Z'
      const result = formatDateTime(isoString)

      expect(result).toMatch(/2025/)
    })
  })

  describe('formatDate', () => {
    it('should format ISO date string to localized date without time', () => {
      const isoString = '2025-11-30T15:42:00.000Z'
      const result = formatDate(isoString, 'en-US')

      // Should include date but not specific time like 15:42
      expect(result).toMatch(/2025/)
      expect(result).toMatch(/30/)
      expect(result).not.toMatch(/15:42/)
    })

    it('should handle different locales', () => {
      const isoString = '2025-11-30T15:42:00.000Z'
      const enResult = formatDate(isoString, 'en-US')
      const ruResult = formatDate(isoString, 'ru-RU')

      expect(enResult).toMatch(/2025/)
      expect(ruResult).toMatch(/2025/)
      expect(enResult).not.toBe(ruResult)
    })

    it('should return empty string for empty input', () => {
      expect(formatDate('')).toBe('')
      expect(formatDate(null)).toBe('')
      expect(formatDate(undefined)).toBe('')
    })

    it('should handle invalid date gracefully', () => {
      const result = formatDate('not-a-date')
      expect(result).toBe('not-a-date')
    })
  })

  describe('formatTime', () => {
    it('should format ISO date string to localized time without date', () => {
      const isoString = '2025-11-30T15:42:00.000Z'
      const result = formatTime(isoString, 'en-US')

      // Should include time components
      expect(result).toMatch(/\d+:\d+/)
    })

    it('should handle different locales', () => {
      const isoString = '2025-11-30T15:42:00.000Z'
      const enResult = formatTime(isoString, 'en-US')
      const ruResult = formatTime(isoString, 'ru-RU')

      // Both should contain time components
      expect(enResult).toMatch(/\d+:\d+/)
      expect(ruResult).toMatch(/\d+:\d+/)
    })

    it('should return empty string for empty input', () => {
      expect(formatTime('')).toBe('')
      expect(formatTime(null)).toBe('')
      expect(formatTime(undefined)).toBe('')
    })

    it('should handle invalid date gracefully', () => {
      const result = formatTime('invalid')
      expect(result).toBe('invalid')
    })

    it('should use en-US as default locale', () => {
      const isoString = '2025-11-30T15:42:00.000Z'
      const result = formatTime(isoString)

      expect(result).toMatch(/\d+:\d+/)
    })
  })

  describe('timezone handling', () => {
    it('should convert UTC time to local timezone', () => {
      const isoString = '2025-11-30T00:00:00.000Z' // Midnight UTC
      const result = formatDateTime(isoString, 'en-US')

      // Result should contain a date, timezone conversion will vary by test environment
      expect(result).toBeTruthy()
      expect(result).toMatch(/202\d/)
    })
  })
})
