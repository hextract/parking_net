import { useState, useEffect } from 'react'
import { Plus, MapPin, DollarSign, Car as CarIcon, Edit2, Trash2, Calendar } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { parkingService } from '../../services/parkingService'
import { PARKING_TYPES } from '../../config/api'
import { useAuth } from '../../context/AuthContext'
import LoadingSpinner from '../../components/LoadingSpinner'

const MyParkings = () => {
  const { user } = useAuth()
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [parkings, setParkings] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [editingParking, setEditingParking] = useState(null)
  const [deleteLoading, setDeleteLoading] = useState(null)
  const [formData, setFormData] = useState({
    name: '',
    city: '',
    address: '',
    parking_type: PARKING_TYPES.OUTDOOR,
    hourly_rate: '',
    capacity: '',
  })
  const [formLoading, setFormLoading] = useState(false)
  const [success, setSuccess] = useState('')
  const [formErrors, setFormErrors] = useState({})

  useEffect(() => {
    loadParkings()
  }, [])

  const loadParkings = async () => {
    setLoading(true)
    setError('')
    try {
      // Get parkings owned by current user by owner_id
      // Backend expects owner_id to be the Keycloak user_id (UUID)
      const data = await parkingService.getParkings({ owner_id: user?.user_id })
      setParkings(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err.message || 'Failed to load parkings')
      setParkings([])
    } finally {
      setLoading(false)
    }
  }

  const handleOpenModal = (parking = null) => {
    if (parking) {
      setEditingParking(parking)
      setFormData({
        name: parking.name,
        city: parking.city,
        address: parking.address,
        parking_type: parking.parking_type,
        hourly_rate: (parking.hourly_rate / 100).toFixed(2), // Convert cents to dollars for editing
        capacity: parking.capacity,
      })
    } else {
      setEditingParking(null)
      setFormData({
        name: '',
        city: '',
        address: '',
        parking_type: PARKING_TYPES.OUTDOOR,
        hourly_rate: '',
        capacity: '',
      })
    }
    setShowModal(true)
  }

  const handleCloseModal = () => {
    setShowModal(false)
    setEditingParking(null)
    setFormData({
      name: '',
      city: '',
      address: '',
      parking_type: PARKING_TYPES.OUTDOOR,
      hourly_rate: '',
      capacity: '',
    })
  }

  const validateForm = () => {
    const errors = {}

    // Max int64 value: 9,223,372,036,854,775,807
    // Using a safe limit for practical purposes: 2,147,483,647 (max int32)
    const MAX_SAFE_INT = 2147483647

    if (!formData.name.trim()) errors.name = t('validation.required')
    if (!formData.city.trim()) errors.city = t('validation.required')
    if (!formData.address.trim()) errors.address = t('validation.required')

    const hourlyRate = parseInt(formData.hourly_rate)
    if (!formData.hourly_rate || isNaN(hourlyRate) || hourlyRate < 0) {
      errors.hourly_rate = t('validation.validNumberRequired')
    } else if (hourlyRate > MAX_SAFE_INT) {
      errors.hourly_rate = t('validation.numberTooLarge')
    }

    const capacity = parseInt(formData.capacity)
    if (!formData.capacity || isNaN(capacity) || capacity < 1) {
      errors.capacity = t('validation.validCapacityRequired')
    } else if (capacity > MAX_SAFE_INT) {
      errors.capacity = t('validation.numberTooLarge')
    }

    setFormErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    if (!validateForm()) {
      return
    }

    setFormLoading(true)

    try {
      const parkingData = {
        ...formData,
        hourly_rate: Math.round(parseFloat(formData.hourly_rate) * 100), // Convert dollars to cents
        capacity: parseInt(formData.capacity),
      }

      if (editingParking) {
        await parkingService.updateParking(editingParking.id, parkingData)
        setSuccess(t('messages.parkingUpdated'))
      } else {
        await parkingService.createParking(parkingData)
        setSuccess(t('messages.parkingCreated'))
      }

      await loadParkings()
      handleCloseModal()

      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err.message || t('messages.loadFailed'))
    } finally {
      setFormLoading(false)
    }
  }

  const handleDelete = async (parkingId) => {
    if (!confirm(t('parking.confirmDelete'))) {
      return
    }

    setDeleteLoading(parkingId)
    setError('')
    try {
      await parkingService.deleteParking(parkingId)
      setParkings(parkings.filter((p) => p.id !== parkingId))
      setSuccess(t('messages.parkingDeleted'))
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      setError(err.message || t('messages.loadFailed'))
    } finally {
      setDeleteLoading(null)
    }
  }

  const getParkingTypeLabel = (type) => {
    return type
      .split('-')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  }

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <LoadingSpinner size="large" />
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">{t('owner.myParkings')}</h1>
        <button onClick={() => handleOpenModal()} className="btn-primary flex items-center space-x-2">
          <Plus className="w-4 h-4" />
          <span>{t('parking.addParking')}</span>
        </button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-6">
          {success}
        </div>
      )}

      {parkings.length === 0 ? (
        <div className="text-center py-12">
          <MapPin className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t('parking.noParkingsYet')}</h3>
          <p className="text-gray-600 mb-4">{t('parking.createFirst')}</p>
          <button onClick={() => handleOpenModal()} className="btn-primary">
            <Plus className="w-4 h-4 inline mr-2" />
            {t('parking.addParking')}
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {parkings.map((parking) => (
            <div key={parking.id} className="card">
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1 min-w-0 mr-2">
                  <h3 className="text-lg font-semibold text-gray-900 truncate">{parking.name}</h3>
                  <p className="text-sm text-gray-600 flex items-center mt-1">
                    <MapPin className="w-4 h-4 mr-1 flex-shrink-0" />
                    <span className="truncate">{parking.city}</span>
                  </p>
                </div>
                <span className="badge bg-primary-100 text-primary-800 flex-shrink-0">
                  {getParkingTypeLabel(parking.parking_type)}
                </span>
              </div>

              <div className="space-y-2 mb-4">
                <p className="text-sm text-gray-600">
                  <strong>{t('parking.address')}:</strong> {parking.address}
                </p>
                <p className="text-sm text-gray-600 flex items-center">
                  <DollarSign className="w-4 h-4 mr-1" />
                  <strong className="mr-2">{t('parking.rate')}:</strong> ${(parking.hourly_rate / 100).toFixed(2)} {t('parking.perHour')}
                </p>
                <p className="text-sm text-gray-600 flex items-center">
                  <CarIcon className="w-4 h-4 mr-1" />
                  <strong className="mr-2">{t('parking.capacity')}:</strong> {parking.capacity} {t('parking.spots')}
                </p>
              </div>

              <div className="flex space-x-2">
                <button
                  onClick={() => navigate(`/owner/bookings/${parking.id}`)}
                  className="btn-secondary flex-1 flex items-center justify-center space-x-1 text-sm"
                >
                  <Calendar className="w-4 h-4" />
                  <span>{t('actions.bookings')}</span>
                </button>
                <button
                  onClick={() => handleOpenModal(parking)}
                  className="btn-secondary flex items-center justify-center px-3"
                >
                  <Edit2 className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(parking.id)}
                  disabled={deleteLoading === parking.id}
                  className="btn-danger flex items-center justify-center px-3"
                >
                  {deleteLoading === parking.id ? (
                    <LoadingSpinner size="small" />
                  ) : (
                    <Trash2 className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Add/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-2xl w-full p-6 max-h-[90vh] overflow-y-auto">
            <h2 className="text-xl font-bold text-gray-900 mb-4">
              {editingParking ? t('parking.editParking') : t('owner.addParking')}
            </h2>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    {t('parking.name')} *
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
                    {t('parking.city')} *
                  </label>
                  <input
                    type="text"
                    value={formData.city}
                    onChange={(e) => {
                      setFormData({ ...formData, city: e.target.value })
                      if (formErrors.city) setFormErrors({ ...formErrors, city: '' })
                    }}
                    className={`input-field ${formErrors.city ? 'border-red-500' : ''}`}
                    required
                  />
                  {formErrors.city && (
                    <p className="text-red-500 text-sm mt-1">{formErrors.city}</p>
                  )}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  {t('parking.address')} *
                </label>
                <input
                  type="text"
                  value={formData.address}
                  onChange={(e) => {
                    setFormData({ ...formData, address: e.target.value })
                    if (formErrors.address) setFormErrors({ ...formErrors, address: '' })
                  }}
                  className={`input-field ${formErrors.address ? 'border-red-500' : ''}`}
                  required
                />
                {formErrors.address && (
                  <p className="text-red-500 text-sm mt-1">{formErrors.address}</p>
                )}
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    {t('parking.type')} *
                  </label>
                  <select
                    value={formData.parking_type}
                    onChange={(e) => setFormData({ ...formData, parking_type: e.target.value })}
                    className="input-field"
                    required
                  >
                    {Object.values(PARKING_TYPES).map((type) => (
                      <option key={type} value={type}>
                        {getParkingTypeLabel(type)}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    {t('parking.hourlyRate')} (USD) *
                  </label>
                  <div className="relative">
                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                    <input
                      type="number"
                      value={formData.hourly_rate}
                      onChange={(e) => {
                        setFormData({ ...formData, hourly_rate: e.target.value })
                        if (formErrors.hourly_rate) setFormErrors({ ...formErrors, hourly_rate: '' })
                      }}
                      placeholder="10.00"
                      step="0.01"
                      min="0.01"
                      className={`input-field pl-8 ${formErrors.hourly_rate ? 'border-red-500' : ''}`}
                      required
                    />
                  </div>
                  {formErrors.hourly_rate && (
                    <p className="text-red-500 text-sm mt-1">{formErrors.hourly_rate}</p>
                  )}
                  <p className="mt-1 text-xs text-gray-500">
                    {t('parking.hourlyRateHint')}
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    {t('parking.capacity')} *
                  </label>
                  <input
                    type="number"
                    value={formData.capacity}
                    onChange={(e) => {
                      setFormData({ ...formData, capacity: e.target.value })
                      if (formErrors.capacity) setFormErrors({ ...formErrors, capacity: '' })
                    }}
                    className={`input-field ${formErrors.capacity ? 'border-red-500' : ''}`}
                    min="1"
                    max="2147483647"
                    required
                  />
                  {formErrors.capacity && (
                    <p className="text-red-500 text-sm mt-1">{formErrors.capacity}</p>
                  )}
                </div>
              </div>

              <div className="flex space-x-3 pt-4">
                <button
                  type="button"
                  onClick={handleCloseModal}
                  className="btn-secondary flex-1"
                  disabled={formLoading}
                >
                  {t('actions.cancel')}
                </button>
                <button type="submit" className="btn-primary flex-1" disabled={formLoading}>
                  {formLoading ? (
                    <LoadingSpinner size="small" />
                  ) : editingParking ? (
                    t('actions.update')
                  ) : (
                    t('actions.create')
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default MyParkings
