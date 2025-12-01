import { useState, useEffect } from 'react'
import { Search, MapPin, DollarSign, Car as CarIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { parkingService } from '../../services/parkingService'
import { bookingService } from '../../services/bookingService'
import { PARKING_TYPES } from '../../config/api'
import LoadingSpinner from '../../components/LoadingSpinner'
import { format } from 'date-fns'

const SearchParking = () => {
  const { t } = useTranslation()
  const [filters, setFilters] = useState({
    city: '',
    name: '',
    parking_type: '',
  })
  const [parkings, setParkings] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [selectedParking, setSelectedParking] = useState(null)
  const [bookingData, setBookingData] = useState({
    date_from: '',
    date_to: '',
  })
  const [bookingLoading, setBookingLoading] = useState(false)
  const [bookingSuccess, setBookingSuccess] = useState(false)
  const [bookingError, setBookingError] = useState('')

  useEffect(() => {
    searchParkings()
  }, [])

  const searchParkings = async () => {
    setLoading(true)
    setError('')
    try {
      const data = await parkingService.getParkings(filters)
      setParkings(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err.message || 'Failed to load parkings')
      setParkings([])
    } finally {
      setLoading(false)
    }
  }

  const handleFilterChange = (e) => {
    setFilters({
      ...filters,
      [e.target.name]: e.target.value,
    })
  }

  const handleSearch = (e) => {
    e.preventDefault()
    searchParkings()
  }

  const handleBooking = async (e) => {
    e.preventDefault()
    if (!selectedParking) return

    // Validate dates and times
    const dateFrom = new Date(bookingData.date_from)
    const dateTo = new Date(bookingData.date_to)
    const now = new Date()

    if (dateFrom < now) {
      setBookingError(t('booking.errorPastDate'))
      return
    }

    if (dateTo <= dateFrom) {
      setBookingError(t('booking.errorEndBeforeStart'))
      return
    }

    const durationInHours = (dateTo - dateFrom) / (1000 * 60 * 60)
    if (durationInHours < 1) {
      setBookingError(t('booking.errorMinimumDuration'))
      return
    }

    setBookingLoading(true)
    setBookingError('')
    setBookingSuccess(false)

    try {
      // Convert dates to ISO 8601 format (RFC3339) as expected by the API
      // Format: 2024-12-31T10:00:00Z
      const formatDateToISO = (dateString) => {
        const date = new Date(dateString)
        return date.toISOString()
      }

      const formattedData = {
        parking_place_id: selectedParking.id,
        date_from: formatDateToISO(bookingData.date_from),
        date_to: formatDateToISO(bookingData.date_to),
      }

      await bookingService.createBooking(formattedData)
      setBookingSuccess(true)
      setSelectedParking(null)
      setBookingData({ date_from: '', date_to: '' })
      setBookingError('')

      setTimeout(() => {
        setBookingSuccess(false)
      }, 3000)
    } catch (err) {
      const errorMessage = err.message || err.data?.error_message || ''
      const lowerMessage = errorMessage.toLowerCase()
      
      if (lowerMessage.includes('insufficient funds') || 
          lowerMessage.includes('payment processing failed') ||
          lowerMessage.includes('insufficient balance')) {
        setBookingError(t('messages.insufficientFunds'))
      } else if (lowerMessage.includes('payment')) {
        setBookingError(t('messages.paymentFailed'))
      } else {
        setBookingError(errorMessage || t('messages.loadFailed'))
      }
    } finally {
      setBookingLoading(false)
    }
  }

  const getParkingTypeLabel = (type) => {
    return t(`parkingTypes.${type}`)
  }

  return (
    <div className="max-w-7xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">{t('driver.searchParking')}</h1>

      {/* Search Form */}
      <form onSubmit={handleSearch} className="card mb-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {t('parking.city')}
            </label>
            <input
              type="text"
              name="city"
              value={filters.city}
              onChange={handleFilterChange}
              className="input-field"
              placeholder={t('parking.cityPlaceholder')}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {t('parking.name')}
            </label>
            <input
              type="text"
              name="name"
              value={filters.name}
              onChange={handleFilterChange}
              className="input-field"
              placeholder={t('parking.namePlaceholder')}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {t('parking.type')}
            </label>
            <select
              name="parking_type"
              value={filters.parking_type}
              onChange={handleFilterChange}
              className="input-field"
            >
              <option value="">{t('parking.allTypes')}</option>
              {Object.values(PARKING_TYPES).map((type) => (
                <option key={type} value={type}>
                  {getParkingTypeLabel(type)}
                </option>
              ))}
            </select>
          </div>

          <div className="flex items-end">
            <button type="submit" disabled={loading} className="btn-primary w-full">
              <Search className="w-4 h-4 inline mr-2" />
              {t('actions.search')}
            </button>
          </div>
        </div>
      </form>

      {/* Success Message */}
      {bookingSuccess && (
        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-6">
          {t('messages.bookingCreated')}
        </div>
      )}

      {/* Error Message - only for search errors */}
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
          {error}
        </div>
      )}

      {/* Results */}
      {loading ? (
        <div className="flex justify-center py-12">
          <LoadingSpinner size="large" />
        </div>
      ) : parkings.length === 0 ? (
        <div className="text-center py-12">
          <MapPin className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t('parking.noParking')}</h3>
          <p className="text-gray-600">{t('parking.noParkingDesc')}</p>
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
                <p className="text-sm text-gray-600 truncate">
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

              <button
                onClick={() => setSelectedParking(parking)}
                className="btn-primary w-full"
              >
                {t('booking.bookNow')}
              </button>
            </div>
          ))}
        </div>
      )}

      {/* Booking Modal */}
      {selectedParking && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">
              {t('booking.bookParking', { name: selectedParking.name })}
            </h2>

            {bookingError && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
                {bookingError}
              </div>
            )}

            <form onSubmit={handleBooking} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  {t('booking.startDateTime')}
                </label>
                <input
                  type="datetime-local"
                  value={bookingData.date_from}
                  onChange={(e) =>
                    setBookingData({ ...bookingData, date_from: e.target.value })
                  }
                  className="input-field"
                  required
                  min={format(new Date(), "yyyy-MM-dd'T'HH:mm")}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  {t('booking.endDateTime')}
                </label>
                <input
                  type="datetime-local"
                  value={bookingData.date_to}
                  onChange={(e) =>
                    setBookingData({ ...bookingData, date_to: e.target.value })
                  }
                  className="input-field"
                  required
                  min={bookingData.date_from || format(new Date(), "yyyy-MM-dd'T'HH:mm")}
                />
              </div>

              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700">{t('booking.hourlyRate')}:</span>
                  <span className="text-lg font-bold text-blue-600">${(selectedParking.hourly_rate / 100).toFixed(2)}/{t('booking.hour')}</span>
                </div>
                {bookingData.date_from && bookingData.date_to && (() => {
                  const start = new Date(bookingData.date_from)
                  const end = new Date(bookingData.date_to)
                  const hours = Math.max(0, (end - start) / (1000 * 60 * 60))
                  const estimatedCost = (Math.ceil(hours) * selectedParking.hourly_rate / 100).toFixed(2)

                  if (hours > 0) {
                    return (
                      <>
                        <div className="flex items-center justify-between text-sm text-gray-600 mt-1">
                          <span>{t('booking.duration')}:</span>
                          <span className="font-medium">{Math.ceil(hours)} {Math.ceil(hours) === 1 ? t('booking.hour') : t('booking.hours')}</span>
                        </div>
                        <div className="flex items-center justify-between mt-2 pt-2 border-t border-blue-300">
                          <span className="text-sm font-semibold text-gray-700">{t('booking.estimatedCost')}:</span>
                          <span className="text-xl font-bold text-blue-700">${estimatedCost}</span>
                        </div>
                      </>
                    )
                  }
                  return null
                })()}
                <p className="text-xs text-gray-500 mt-2">
                  * {t('booking.minimumDuration')}
                </p>
              </div>

              <div className="flex space-x-3">
                <button
                  type="button"
                  onClick={() => {
                    setSelectedParking(null)
                    setBookingData({ date_from: '', date_to: '' })
                    setBookingError('')
                  }}
                  className="btn-secondary flex-1"
                  disabled={bookingLoading}
                >
                  {t('actions.cancel')}
                </button>
                <button
                  type="submit"
                  className="btn-primary flex-1"
                  disabled={bookingLoading}
                >
                  {bookingLoading ? <LoadingSpinner size="small" /> : t('booking.confirmBooking')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default SearchParking
