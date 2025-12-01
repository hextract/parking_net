/**
 * Environment configuration
 * All environment variables must be prefixed with VITE_ to be exposed to the client.
 *
 * Default values are taken from ../.env-example
 * To customize, set VITE_* environment variables at build time.
 */

// Get base host from environment or default to localhost
const BASE_HOST = import.meta.env.VITE_BASE_HOST || 'localhost'

// Service ports from root .env (matching .env-example)
const NGINX_PORT = import.meta.env.VITE_NGINX_PORT || '80'
const AUTH_REST_PORT = import.meta.env.VITE_AUTH_REST_PORT || '8800'
const PARKING_REST_PORT = import.meta.env.VITE_PARKING_REST_PORT || '8888'
const BOOKING_REST_PORT = import.meta.env.VITE_BOOKING_REST_PORT || '8880'
const PAYMENT_REST_PORT = import.meta.env.VITE_PAYMENT_REST_PORT || '8890'
const KEYCLOAK_PORT = import.meta.env.VITE_KEYCLOAK_PORT || '8080'
const JAEGER_PORT = import.meta.env.VITE_JAEGER_PORT || '16686'
const PROMETHEUS_PORT = import.meta.env.VITE_PROMETHEUS_PORT || '9090'
const GRAFANA_PORT = import.meta.env.VITE_GRAFANA_PORT || '3000'

// API Configuration (via Nginx Gateway)
// Use HTTPS if VITE_API_BASE_URL is not set and we're in production (not localhost)
const isProduction = BASE_HOST !== 'localhost' && BASE_HOST !== '127.0.0.1'
const defaultProtocol = isProduction ? 'https' : 'http'
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || `${defaultProtocol}://${BASE_HOST}${NGINX_PORT === '80' ? '' : `:${NGINX_PORT}`}`

// Monitoring Tools - use HTTPS subdomains in production
const monitoringProtocol = isProduction ? 'https' : 'http'
const jaegerSubdomain = import.meta.env.VITE_JAEGER_SUBDOMAIN || (isProduction ? 'jaeger.backend.parking-net.space' : null)
const prometheusSubdomain = import.meta.env.VITE_PROMETHEUS_SUBDOMAIN || (isProduction ? 'prometheus.backend.parking-net.space' : null)
const grafanaSubdomain = import.meta.env.VITE_GRAFANA_SUBDOMAIN || (isProduction ? 'grafana.backend.parking-net.space' : null)
const keycloakSubdomain = import.meta.env.VITE_KEYCLOAK_SUBDOMAIN || (isProduction ? 'keycloak.backend.parking-net.space' : null)

const jaegerHost = jaegerSubdomain || `${BASE_HOST}:${JAEGER_PORT}`
const prometheusHost = prometheusSubdomain || `${BASE_HOST}:${PROMETHEUS_PORT}`
const grafanaHost = grafanaSubdomain || `${BASE_HOST}:${GRAFANA_PORT}`
const keycloakHost = keycloakSubdomain || `${BASE_HOST}:${KEYCLOAK_PORT}`

export const JAEGER_URL = import.meta.env.VITE_JAEGER_URL || `${monitoringProtocol}://${jaegerHost}`
export const PROMETHEUS_URL = import.meta.env.VITE_PROMETHEUS_URL || `${monitoringProtocol}://${prometheusHost}`
export const GRAFANA_URL = import.meta.env.VITE_GRAFANA_URL || `${monitoringProtocol}://${grafanaHost}`

// Backend Services - use API_BASE_URL paths in production, direct ports in development
export const AUTH_SERVICE_URL = import.meta.env.VITE_AUTH_SERVICE_URL || (isProduction ? `${API_BASE_URL}/auth` : `http://${BASE_HOST}:${AUTH_REST_PORT}`)
export const PARKING_SERVICE_URL = import.meta.env.VITE_PARKING_SERVICE_URL || (isProduction ? `${API_BASE_URL}/parking` : `http://${BASE_HOST}:${PARKING_REST_PORT}`)
export const BOOKING_SERVICE_URL = import.meta.env.VITE_BOOKING_SERVICE_URL || (isProduction ? `${API_BASE_URL}/booking` : `http://${BASE_HOST}:${BOOKING_REST_PORT}`)
export const PAYMENT_SERVICE_URL = import.meta.env.VITE_PAYMENT_SERVICE_URL || (isProduction ? `${API_BASE_URL}/payment` : `http://${BASE_HOST}:${PAYMENT_REST_PORT}`)
export const KEYCLOAK_URL = import.meta.env.VITE_KEYCLOAK_URL || `${monitoringProtocol}://${keycloakHost}`

// Metrics Endpoints (via Nginx Gateway)
export const AUTH_METRICS_URL = import.meta.env.VITE_AUTH_METRICS_URL || `${API_BASE_URL}/auth/metrics`
export const PARKING_METRICS_URL = import.meta.env.VITE_PARKING_METRICS_URL || `${API_BASE_URL}/parking/metrics`
export const BOOKING_METRICS_URL = import.meta.env.VITE_BOOKING_METRICS_URL || `${API_BASE_URL}/booking/metrics`
export const PAYMENT_METRICS_URL = import.meta.env.VITE_PAYMENT_METRICS_URL || `${API_BASE_URL}/payment/metrics`

// Export all as a single object for easier access
export const ENV = {
  API_BASE_URL,
  AUTH_SERVICE_URL,
  PARKING_SERVICE_URL,
  BOOKING_SERVICE_URL,
  PAYMENT_SERVICE_URL,
  KEYCLOAK_URL,
  AUTH_METRICS_URL,
  PARKING_METRICS_URL,
  BOOKING_METRICS_URL,
  PAYMENT_METRICS_URL,
  JAEGER_URL,
  PROMETHEUS_URL,
  GRAFANA_URL,
}

export default ENV
