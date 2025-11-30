/**
 * Frontend Validation Utilities
 * Matches backend validation rules from auth/internal/utils/validation.go
 */

// Validation constants (matching backend)
export const VALIDATION_RULES = {
  LOGIN: {
    MIN_LENGTH: 3,
    MAX_LENGTH: 50,
    PATTERN: /^[a-zA-Z0-9_]+$/,
  },
  PASSWORD: {
    MIN_LENGTH: 8,
    MAX_LENGTH: 128,
  },
  EMAIL: {
    MIN_LENGTH: 5,
    MAX_LENGTH: 254,
    PATTERN: /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/,
  },
  TELEGRAM_ID: {
    MIN: 1,
    MAX: 9223372036854775807,
  },
  PARKING: {
    NAME: {
      MIN_LENGTH: 3,
      MAX_LENGTH: 100,
    },
    CITY: {
      MIN_LENGTH: 2,
      MAX_LENGTH: 50,
    },
    ADDRESS: {
      MIN_LENGTH: 5,
      MAX_LENGTH: 200,
    },
    CAPACITY: {
      MIN: 1,
      MAX: 2147483647,
    },
    HOURLY_RATE: {
      MIN: 1,
      MAX: 2147483647,
    },
  },
  PROMOCODE: {
    CODE: {
      MIN_LENGTH: 4,
      MAX_LENGTH: 20,
      PATTERN: /^[A-Za-z0-9_-]+$/, // Allow both uppercase and lowercase
    },
    AMOUNT: {
      MIN: 1,
      MAX: 2147483647,
    },
    MAX_USES: {
      MIN: 1,
      MAX: 2147483647,
    },
  },
}

/**
 * Validate login
 * @param {string} login
 * @returns {string|null} Error message or null if valid
 */
export const validateLogin = (login) => {
  if (!login || login.trim() === '') {
    return 'validation.loginRequired'
  }

  const length = login.length
  if (length < VALIDATION_RULES.LOGIN.MIN_LENGTH || length > VALIDATION_RULES.LOGIN.MAX_LENGTH) {
    return `validation.loginLength`
  }

  if (!VALIDATION_RULES.LOGIN.PATTERN.test(login)) {
    return 'validation.loginPattern'
  }

  return null
}

/**
 * Validate password
 * @param {string} password
 * @returns {string|null} Error message or null if valid
 */
export const validatePassword = (password) => {
  if (!password || password.trim() === '') {
    return 'validation.passwordRequired'
  }

  const length = password.length
  if (length < VALIDATION_RULES.PASSWORD.MIN_LENGTH || length > VALIDATION_RULES.PASSWORD.MAX_LENGTH) {
    return 'validation.passwordLength'
  }

  // Check for uppercase, lowercase, and number
  const hasUpper = /[A-Z]/.test(password)
  const hasLower = /[a-z]/.test(password)
  const hasNumber = /[0-9]/.test(password)

  if (!hasUpper || !hasLower || !hasNumber) {
    return 'validation.passwordComplexity'
  }

  return null
}

/**
 * Validate email
 * @param {string} email
 * @returns {string|null} Error message or null if valid
 */
export const validateEmail = (email) => {
  if (!email || email.trim() === '') {
    return 'validation.emailRequired'
  }

  const length = email.length
  if (length < VALIDATION_RULES.EMAIL.MIN_LENGTH || length > VALIDATION_RULES.EMAIL.MAX_LENGTH) {
    return 'validation.emailLength'
  }

  if (!VALIDATION_RULES.EMAIL.PATTERN.test(email)) {
    return 'validation.emailInvalid'
  }

  return null
}

/**
 * Validate Telegram ID
 * @param {number} telegramId
 * @returns {string|null} Error message or null if valid
 */
export const validateTelegramId = (telegramId) => {
  // Telegram ID is optional - allow empty/null/undefined/0
  if (telegramId === undefined || telegramId === null || telegramId === '' || telegramId === 0) {
    return null // Optional field - user doesn't have Telegram
  }

  const id = parseInt(telegramId)
  if (isNaN(id)) {
    return 'validation.telegramIdInvalid'
  }

  // Only validate if a value was actually provided
  if (id < VALIDATION_RULES.TELEGRAM_ID.MIN) {
    return 'validation.telegramIdPositive'
  }

  if (id > VALIDATION_RULES.TELEGRAM_ID.MAX) {
    return 'validation.telegramIdTooLarge'
  }

  return null
}

/**
 * Validate parking name
 * @param {string} name
 * @returns {string|null} Error message or null if valid
 */
export const validateParkingName = (name) => {
  if (!name || name.trim() === '') {
    return 'validation.parkingNameRequired'
  }

  const length = name.trim().length
  if (length < VALIDATION_RULES.PARKING.NAME.MIN_LENGTH || length > VALIDATION_RULES.PARKING.NAME.MAX_LENGTH) {
    return 'validation.parkingNameLength'
  }

  return null
}

/**
 * Validate city name
 * @param {string} city
 * @returns {string|null} Error message or null if valid
 */
export const validateCity = (city) => {
  if (!city || city.trim() === '') {
    return 'validation.cityRequired'
  }

  const length = city.trim().length
  if (length < VALIDATION_RULES.PARKING.CITY.MIN_LENGTH || length > VALIDATION_RULES.PARKING.CITY.MAX_LENGTH) {
    return 'validation.cityLength'
  }

  return null
}

/**
 * Validate address
 * @param {string} address
 * @returns {string|null} Error message or null if valid
 */
export const validateAddress = (address) => {
  if (!address || address.trim() === '') {
    return 'validation.addressRequired'
  }

  const length = address.trim().length
  if (length < VALIDATION_RULES.PARKING.ADDRESS.MIN_LENGTH || length > VALIDATION_RULES.PARKING.ADDRESS.MAX_LENGTH) {
    return 'validation.addressLength'
  }

  return null
}

/**
 * Validate capacity
 * @param {number} capacity
 * @returns {string|null} Error message or null if valid
 */
export const validateCapacity = (capacity) => {
  const num = parseInt(capacity)
  if (isNaN(num)) {
    return 'validation.validNumberRequired'
  }

  if (num < VALIDATION_RULES.PARKING.CAPACITY.MIN || num > VALIDATION_RULES.PARKING.CAPACITY.MAX) {
    return 'validation.validCapacityRequired'
  }

  return null
}

/**
 * Validate hourly rate (expects dollars, will be converted to cents)
 * @param {number} rate - Rate in dollars
 * @returns {string|null} Error message or null if valid
 */
export const validateHourlyRate = (rate) => {
  const num = parseFloat(rate)
  if (isNaN(num)) {
    return 'validation.validNumberRequired'
  }

  if (num <= 0) {
    return 'validation.hourlyRateInvalid'
  }

  // Convert to cents for backend validation
  const rateInCents = Math.round(num * 100)
  if (rateInCents > VALIDATION_RULES.PARKING.HOURLY_RATE.MAX) {
    return 'validation.numberTooLarge'
  }

  return null
}

/**
 * Validate promocode code format
 * @param {string} code
 * @returns {string|null} Error message or null if valid
 */
export const validatePromocodeCode = (code) => {
  if (!code || code.trim() === '') {
    return null // Optional, will be generated if not provided
  }

  const length = code.length
  if (length < VALIDATION_RULES.PROMOCODE.CODE.MIN_LENGTH || length > VALIDATION_RULES.PROMOCODE.CODE.MAX_LENGTH) {
    return 'validation.promocodeCodeLength'
  }

  if (!VALIDATION_RULES.PROMOCODE.CODE.PATTERN.test(code)) {
    return 'validation.promocodeCodePattern'
  }

  return null
}

/**
 * Validate promocode amount (expects dollars, will be converted to cents)
 * @param {number} amount - Amount in dollars
 * @returns {string|null} Error message or null if valid
 */
export const validatePromocodeAmount = (amount) => {
  const num = parseFloat(amount)
  if (isNaN(num)) {
    return 'validation.validNumberRequired'
  }

  if (num <= 0) {
    return 'validation.promocodeAmountInvalid'
  }

  // Convert to cents for backend validation
  const amountInCents = Math.round(num * 100)
  if (amountInCents > VALIDATION_RULES.PROMOCODE.AMOUNT.MAX) {
    return 'validation.numberTooLarge'
  }

  return null
}

/**
 * Validate promocode max uses
 * @param {number} maxUses
 * @returns {string|null} Error message or null if valid
 */
export const validatePromocodeMaxUses = (maxUses) => {
  const num = parseInt(maxUses)
  if (isNaN(num)) {
    return 'validation.validNumberRequired'
  }

  if (num < VALIDATION_RULES.PROMOCODE.MAX_USES.MIN || num > VALIDATION_RULES.PROMOCODE.MAX_USES.MAX) {
    return 'validation.promocodeMaxUsesInvalid'
  }

  return null
}

/**
 * Format validation error for display
 * @param {string} errorKey - Translation key
 * @param {Object} t - Translation function
 * @returns {string} Formatted error message
 */
export const formatValidationError = (errorKey, t) => {
  return t(errorKey)
}
