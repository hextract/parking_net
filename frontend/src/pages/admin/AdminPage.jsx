import React from 'react'
import { useTranslation } from 'react-i18next'
import { Activity, BarChart3, Gauge, Shield, MapPin, Calendar, Database, Wallet } from 'lucide-react'
import ENV from '../../config/env'

const AdminPage = () => {
  const { t } = useTranslation()

  const monitoringTools = [
    {
      key: 'jaeger',
      url: ENV.JAEGER_URL,
      icon: Activity,
      color: 'bg-blue-500',
    },
    {
      key: 'prometheus',
      url: ENV.PROMETHEUS_URL,
      icon: BarChart3,
      color: 'bg-orange-500',
    },
    {
      key: 'grafana',
      url: ENV.GRAFANA_URL,
      icon: Gauge,
      color: 'bg-yellow-500',
    },
  ]

  const backendServices = [
    {
      key: 'auth',
      url: ENV.AUTH_SERVICE_URL,
      metricsUrl: ENV.AUTH_METRICS_URL,
      icon: Shield,
      color: 'bg-indigo-500',
    },
    {
      key: 'parking',
      url: ENV.PARKING_SERVICE_URL,
      metricsUrl: ENV.PARKING_METRICS_URL,
      icon: MapPin,
      color: 'bg-emerald-500',
    },
    {
      key: 'booking',
      url: ENV.BOOKING_SERVICE_URL,
      metricsUrl: ENV.BOOKING_METRICS_URL,
      icon: Calendar,
      color: 'bg-rose-500',
    },
    {
      key: 'payment',
      url: ENV.PAYMENT_SERVICE_URL,
      metricsUrl: ENV.PAYMENT_METRICS_URL,
      icon: Wallet,
      color: 'bg-purple-500',
    },
    {
      key: 'keycloak',
      url: ENV.KEYCLOAK_URL,
      metricsUrl: null,
      icon: Database,
      color: 'bg-cyan-500',
    },
  ]

  const handleToolClick = (url) => {
    window.open(url, '_blank', 'noopener,noreferrer')
  }


  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">{t('admin.title')}</h1>
        <p className="text-gray-600 mt-2">{t('admin.subtitle')}</p>
      </div>

      <div className="mb-8">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">{t('admin.monitoringTools')}</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {monitoringTools.map((tool) => {
            const Icon = tool.icon
            return (
              <button
                key={tool.key}
                onClick={() => handleToolClick(tool.url)}
                className="card text-left hover:shadow-xl transition-all transform hover:-translate-y-1 cursor-pointer"
              >
                <div className="flex items-start space-x-4">
                  <div className={`${tool.color} p-3 rounded-lg`}>
                    <Icon className="w-6 h-6 text-white" />
                  </div>
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-gray-900 mb-1">
                      {t(`admin.${tool.key}.name`)}
                    </h3>
                    <p className="text-sm text-primary-600 font-medium mb-2">
                      {t(`admin.${tool.key}.description`)}
                    </p>
                    <p className="text-xs text-gray-600 mb-3">{t(`admin.${tool.key}.details`)}</p>
                    <p className="text-xs text-gray-500 font-mono bg-gray-100 px-2 py-1 rounded">
                      {tool.url}
                    </p>
                  </div>
                </div>
              </button>
            )
          })}
        </div>
      </div>

      {/* Backend Services */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900">{t('admin.backendServices')}</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {backendServices.map((service) => {
            const Icon = service.icon

            return (
              <div key={service.key} className="card relative">
                <div className="flex items-start space-x-3">
                  <div className={`${service.color} p-2.5 rounded-lg`}>
                    <Icon className="w-5 h-5 text-white" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center space-x-2 mb-1">
                      <h3 className="text-base font-semibold text-gray-900">
                        {t(`admin.services.${service.key}.name`)}
                      </h3>
                    </div>
                    <p className="text-xs text-gray-600 mb-2">
                      {t(`admin.services.${service.key}.description`)}
                    </p>
                    <div className="space-y-1">
                      <button
                        onClick={() => handleToolClick(service.url)}
                        className="text-xs text-primary-600 hover:text-primary-700 font-mono bg-gray-100 px-2 py-1 rounded block w-full text-left hover:bg-primary-50 transition-colors"
                      >
                        {service.url}
                      </button>
                      {service.metricsUrl && (
                        <button
                          onClick={() => handleToolClick(service.metricsUrl)}
                          className="text-xs text-green-600 hover:text-green-700 font-mono bg-gray-100 px-2 py-1 rounded block w-full text-left hover:bg-green-50 transition-colors"
                        >
                          /metrics
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      </div>

      <div className="mt-8 bg-gray-50 border border-gray-200 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">{t('admin.quickInfo')}</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-sm">
          <div>
            <h4 className="font-semibold text-gray-700 mb-2">{t('admin.monitoringTools')}</h4>
            <ul className="space-y-1 text-gray-600">
              <li>• {t('admin.info.jaeger')}</li>
              <li>• {t('admin.info.prometheus')}</li>
              <li>• {t('admin.info.grafana')}</li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold text-gray-700 mb-2">{t('admin.backendServices')}</h4>
            <ul className="space-y-1 text-gray-600">
              <li>• {t('admin.services.auth.name')} - {ENV.AUTH_SERVICE_URL}</li>
              <li>• {t('admin.services.parking.name')} - {ENV.PARKING_SERVICE_URL}</li>
              <li>• {t('admin.services.booking.name')} - {ENV.BOOKING_SERVICE_URL}</li>
              <li>• {t('admin.services.payment.name')} - {ENV.PAYMENT_SERVICE_URL}</li>
              <li>• {t('admin.services.keycloak.name')} - {ENV.KEYCLOAK_URL}</li>
            </ul>
          </div>
        </div>

      </div>
    </div>
  )
}

export default AdminPage
