/**
 * Format an ISO date string to a localized date and time string
 * @param {string} isoString - ISO 8601 date string (e.g., "2025-11-30T15:42:00.000Z")
 * @param {string} locale - Locale string (e.g., "en-US", "ru-RU")
 * @returns {string} Formatted date and time in user's timezone
 */
export const formatDateTime = (isoString, locale = 'en-US') => {
  if (!isoString) return ''

  try {
    const date = new Date(isoString)

    // Check if date is valid
    if (isNaN(date.getTime())) {
      return isoString
    }

    // Format: "Nov 30, 2025, 3:42 PM" (for en-US) or "30 нояб. 2025 г., 15:42" (for ru-RU)
    return date.toLocaleString(locale, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch (error) {
    console.error('Error formatting date:', error)
    return isoString
  }
}

/**
 * Format an ISO date string to a localized date string (without time)
 * @param {string} isoString - ISO 8601 date string
 * @param {string} locale - Locale string
 * @returns {string} Formatted date in user's timezone
 */
export const formatDate = (isoString, locale = 'en-US') => {
  if (!isoString) return ''

  try {
    const date = new Date(isoString)

    if (isNaN(date.getTime())) {
      return isoString
    }

    // Format: "Nov 30, 2025" (for en-US) or "30 нояб. 2025 г." (for ru-RU)
    return date.toLocaleDateString(locale, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  } catch (error) {
    console.error('Error formatting date:', error)
    return isoString
  }
}

/**
 * Format an ISO date string to a localized time string (without date)
 * @param {string} isoString - ISO 8601 date string
 * @param {string} locale - Locale string
 * @returns {string} Formatted time in user's timezone
 */
export const formatTime = (isoString, locale = 'en-US') => {
  if (!isoString) return ''

  try {
    const date = new Date(isoString)

    if (isNaN(date.getTime())) {
      return isoString
    }

    // Format: "3:42 PM" (for en-US) or "15:42" (for ru-RU)
    return date.toLocaleTimeString(locale, {
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch (error) {
    console.error('Error formatting time:', error)
    return isoString
  }
}
