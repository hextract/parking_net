/**
 * Secure cookie utilities for token storage
 * Uses httpOnly-like behavior with secure, sameSite attributes
 *
 * SECURITY NOTE:
 * - Cookies set via JavaScript CANNOT be HttpOnly
 * - This makes them vulnerable to XSS attacks
 * - RECOMMENDATION: Move token storage to backend-only HttpOnly cookies
 * - For now, we mitigate with: Secure flag (HTTPS), SameSite=Strict (CSRF protection)
 *
 * PRODUCTION RECOMMENDATIONS:
 * 1. Backend should set tokens in HttpOnly cookies
 * 2. Use HTTPS only (enforce with HSTS headers)
 * 3. Implement Content-Security-Policy headers
 * 4. Consider using short-lived tokens with refresh tokens
 */

const COOKIE_NAME = 'auth_token'
const COOKIE_OPTIONS = {
  secure: window.location.protocol === 'https:', // Secure in production
  sameSite: 'Strict', // CSRF protection
  path: '/',
  maxAge: 7 * 24 * 60 * 60, // 7 days (consider shorter for production)
}

/**
 * Set authentication token in cookie
 *
 * WARNING: Not HttpOnly! Accessible to JavaScript (XSS vulnerability)
 * In production, backend should set tokens in HttpOnly cookies via Set-Cookie header
 */
export const setAuthToken = (token) => {
  if (!token) {
    removeAuthToken()
    return
  }

  const expires = new Date()
  expires.setTime(expires.getTime() + COOKIE_OPTIONS.maxAge * 1000)

  let cookieString = `${COOKIE_NAME}=${encodeURIComponent(token)}`
  cookieString += `; expires=${expires.toUTCString()}`
  cookieString += `; path=${COOKIE_OPTIONS.path}`
  cookieString += `; SameSite=${COOKIE_OPTIONS.sameSite}`

  if (COOKIE_OPTIONS.secure) {
    cookieString += '; Secure'
  }

  document.cookie = cookieString
}

/**
 * Get authentication token from cookie
 */
export const getAuthToken = () => {
  const name = COOKIE_NAME + '='
  const decodedCookie = decodeURIComponent(document.cookie)
  const cookieArray = decodedCookie.split(';')

  for (let i = 0; i < cookieArray.length; i++) {
    let cookie = cookieArray[i]
    while (cookie.charAt(0) === ' ') {
      cookie = cookie.substring(1)
    }
    if (cookie.indexOf(name) === 0) {
      return decodeURIComponent(cookie.substring(name.length, cookie.length))
    }
  }
  return null
}

/**
 * Remove authentication token from cookie
 */
export const removeAuthToken = () => {
  document.cookie = `${COOKIE_NAME}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=${COOKIE_OPTIONS.path}; SameSite=${COOKIE_OPTIONS.sameSite}`
  if (COOKIE_OPTIONS.secure) {
    document.cookie = `${COOKIE_NAME}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=${COOKIE_OPTIONS.path}; SameSite=${COOKIE_OPTIONS.sameSite}; Secure`
  }
}
