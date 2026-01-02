import { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { serviceService } from '../services/serviceService'
import LoadingSpinner from './LoadingSpinner'

const ParkingServicesModal = ({ parking, isOpen, onClose }) => {
  const { t } = useTranslation()
  const [services, setServices] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [showServiceForm, setShowServiceForm] = useState(false)
  const [editingService, setEditingService] = useState(null)
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    price: '',
  })
  const [formErrors, setFormErrors] = useState({})
  const [formLoading, setFormLoading] = useState(false)
  const [deleteLoading, setDeleteLoading] = useState(null)

  useEffect(() => {
    if (isOpen && parking) {
      loadServices()
    }
  }, [isOpen, parking])

  const loadServices = async () => {
    setLoading(true)
    setError('')
    try {
      const data = await serviceService.getServices(parking.id)
      setServices(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err.message || t('messages.loadFailed'))
    } finally {
      setLoading(false)
    }
  }

  const handleOpenServiceForm = (service = null) => {
    if (service) {
      setEditingService(service)
      setFormData({
        name: service.name || '',
        description: service.description || '',
        price: ((service.price || 0) / 100).toFixed(2),
      })
    } else {
      setEditingService(null)
      setFormData({
        name: '',
        description: '',
        price: '',
      })
    }
    setFormErrors({})
    setShowServiceForm(true)
  }

  const handleCloseServiceForm = () => {
    setShowServiceForm(false)
    setEditingService(null)
    setFormData({
      name: '',
      description: '',
      price: '',
    })
    setFormErrors({})
  }

  const validateForm = () => {
    const errors = {}
    if (!formData.name.trim()) errors.name = t('validation.required')
    const price = parseFloat(formData.price)
    if (!formData.price || isNaN(price) || price < 0) {
      errors.price = t('validation.validNumberRequired')
    }
    setFormErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmitService = async (e) => {
    e.preventDefault()
    setError('')
    if (!validateForm()) return

    setFormLoading(true)
    try {
      const serviceData = {
        name: formData.name.trim(),
        description: formData.description.trim(),
        price: Math.round(parseFloat(formData.price) * 100),
      }

      if (editingService) {
        await serviceService.updateService(parking.id, editingService.id, serviceData)
      } else {
        await serviceService.createService(parking.id, serviceData)
      }

      await loadServices()
      handleCloseServiceForm()
    } catch (err) {
      setError(err.message || t('messages.loadFailed'))
    } finally {
      setFormLoading(false)
    }
  }

  const handleDeleteService = async (serviceId) => {
    if (!confirm(t('parking.confirmDeleteService'))) return

    setDeleteLoading(serviceId)
    setError('')
    try {
      await serviceService.deleteService(parking.id, serviceId)
      await loadServices()
    } catch (err) {
      setError(err.message || t('messages.loadFailed'))
    } finally {
      setDeleteLoading(null)
    }
  }

  if (!isOpen || !parking) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-center p-6 border-b">
          <h2 className="text-xl font-bold text-gray-900">
            {t('parking.manageServices')} - {parking.name}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        <div className="p-6">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
              {error}
            </div>
          )}

          {!showServiceForm ? (
            <>
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-semibold text-gray-900">{t('parking.services')}</h3>
                <button
                  onClick={() => handleOpenServiceForm()}
                  className="btn-primary flex items-center space-x-2"
                >
                  <Plus className="w-4 h-4" />
                  <span>{t('parking.addService')}</span>
                </button>
              </div>

              {loading ? (
                <div className="flex justify-center py-12">
                  <LoadingSpinner size="large" />
                </div>
              ) : services.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-gray-600 mb-4">{t('parking.noServices')}</p>
                  <button
                    onClick={() => handleOpenServiceForm()}
                    className="btn-primary"
                  >
                    <Plus className="w-4 h-4 inline mr-2" />
                    {t('parking.addFirstService')}
                  </button>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {services.map((service) => (
                    <div key={service.id} className="card">
                      <div className="flex justify-between items-start mb-2">
                        <h4 className="text-lg font-semibold text-gray-900">{service.name}</h4>
                        <div className="flex space-x-2">
                          <button
                            onClick={() => handleOpenServiceForm(service)}
                            className="text-blue-600 hover:text-blue-800"
                          >
                            <Edit2 className="w-4 h-4" />
                          </button>
                          <button
                            onClick={() => handleDeleteService(service.id)}
                            disabled={deleteLoading === service.id}
                            className="text-red-600 hover:text-red-800"
                          >
                            {deleteLoading === service.id ? (
                              <LoadingSpinner size="small" />
                            ) : (
                              <Trash2 className="w-4 h-4" />
                            )}
                          </button>
                        </div>
                      </div>
                      {service.description && (
                        <p className="text-sm text-gray-600 mb-2">{service.description}</p>
                      )}
                      <p className="text-lg font-bold text-blue-600">
                        ${((service.price || 0) / 100).toFixed(2)}
                      </p>
                    </div>
                  ))}
                </div>
              )}
            </>
          ) : (
            <div>
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                {editingService ? t('parking.editService') : t('parking.addService')}
              </h3>
              <form onSubmit={handleSubmitService} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    {t('parking.serviceName')} *
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => {
                      setFormData({ ...formData, name: e.target.value })
                      if (formErrors.name) setFormErrors({ ...formErrors, name: '' })
                    }}
                    className={`input-field ${formErrors.name ? 'border-red-500' : ''}`}
                    required
                  />
                  {formErrors.name && (
                    <p className="text-red-500 text-sm mt-1">{formErrors.name}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    {t('parking.serviceDescription')}
                  </label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    className="input-field"
                    rows="3"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    {t('parking.servicePrice')} (USD) *
                  </label>
                  <div className="relative">
                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                    <input
                      type="number"
                      value={formData.price}
                      onChange={(e) => {
                        setFormData({ ...formData, price: e.target.value })
                        if (formErrors.price) setFormErrors({ ...formErrors, price: '' })
                      }}
                      placeholder="15.00"
                      step="0.01"
                      min="0.01"
                      className={`input-field pl-8 ${formErrors.price ? 'border-red-500' : ''}`}
                      required
                    />
                  </div>
                  {formErrors.price && (
                    <p className="text-red-500 text-sm mt-1">{formErrors.price}</p>
                  )}
                  <p className="mt-1 text-xs text-gray-500">
                    {t('parking.servicePriceHint')}
                  </p>
                </div>

                <div className="flex space-x-3 pt-4">
                  <button
                    type="button"
                    onClick={handleCloseServiceForm}
                    className="btn-secondary flex-1"
                    disabled={formLoading}
                  >
                    {t('actions.cancel')}
                  </button>
                  <button
                    type="submit"
                    className="btn-primary flex-1"
                    disabled={formLoading}
                  >
                    {formLoading ? (
                      <LoadingSpinner size="small" />
                    ) : editingService ? (
                      t('actions.update')
                    ) : (
                      t('actions.create')
                    )}
                  </button>
                </div>
              </form>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default ParkingServicesModal

