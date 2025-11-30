import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Calendar, MapPin, DollarSign, Trash2 } from 'lucide-react'
import { bookingService } from '../../services/bookingService'
import { parkingService } from '../../services/parkingService'
import { BOOKING_STATUSES } from '../../config/api'
import { useAuth } from '../../context/AuthContext'
import { useTranslation } from 'react-i18next'
import LoadingSpinner from '../../components/LoadingSpinner'
import { formatDateTime } from '../../utils/dateUtils'

const MyBookings = () => {
  const { user } = useAuth()
  const { t, i18n } = useTranslation()
  const navigate = useNavigate()
  const [bookings, setBookings] = useState([])
  const [parkingDetails, setParkingDetails] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [deleteLoading, setDeleteLoading] = useState(null)

  useEffect(() => {
    loadBookings()
  }, [])

  const loadBookings = async () => {
    setLoading(true)
    setError('')
    try {
      // Get bookings for current user by user_id
      // Backend expects user_id (Keycloak UUID from /auth/me)
      const data = await bookingService.getBookings({ user_id: user?.user_id })
      const bookingsArray = Array.isArray(data) ? data : []
      setBookings(bookingsArray)

      // Fetch parking details for each booking
      const details = {}
      for (const booking of bookingsArray) {
        try {
          const parking = await parkingService.getParkingById(booking.parking_place_id)
          details[booking.parking_place_id] = parking
        } catch (err) {
        }
      }
      setParkingDetails(details)
    } catch (err) {
      setError(err.message || 'Failed to load bookings')
      setBookings([])
    } finally {
      setLoading(false)
    }
  }

  const handleCancelBooking = async (bookingId) => {
    if (!confirm(t('booking.confirmCancel'))) {
      return
    }

    setDeleteLoading(bookingId)
    setError('')
    try {
      await bookingService.deleteBooking(bookingId)
      setBookings(bookings.filter((b) => b.booking_id !== bookingId))
    } catch (err) {
      setError(err.message || 'Failed to cancel booking')
    } finally {
      setDeleteLoading(null)
    }
  }

  const getStatusBadgeClass = (status) => {
    switch (status) {
      case BOOKING_STATUSES.CONFIRMED:
        return 'badge-confirmed'
      case BOOKING_STATUSES.WAITING:
        return 'badge-waiting'
      case BOOKING_STATUSES.CANCELED:
        return 'badge-canceled'
      default:
        return 'badge bg-gray-100 text-gray-800'
    }
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
      <h1 className="text-3xl font-bold text-gray-900 mb-8">{t('driver.myBookings')}</h1>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
          {error}
        </div>
      )}

      {bookings.length === 0 ? (
        <div className="text-center py-12">
          <Calendar className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t('booking.noBookings')}</h3>
          <p className="text-gray-600 mb-4">{t('booking.noBookingsDesc')}</p>
          <button
            onClick={() => navigate('/driver/search')}
            className="btn-primary"
          >
            {t('driver.searchParking')}
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {bookings.map((booking) => {
            const parking = parkingDetails[booking.parking_place_id]
            return (
              <div key={booking.booking_id} className="card">
                <div className="flex flex-col md:flex-row md:items-center md:justify-between">
                  <div className="flex-1">
                    <div className="flex items-start justify-between mb-3">
                      <div>
                        <h3 className="text-lg font-semibold text-gray-900">
                          {parking ? parking.name : `Parking #${booking.parking_place_id}`}
                        </h3>
                        {parking && (
                          <p className="text-sm text-gray-600 flex items-center mt-1">
                            <MapPin className="w-4 h-4 mr-1" />
                            {parking.city}, {parking.address}
                          </p>
                        )}
                      </div>
                      <span className={`badge ${getStatusBadgeClass(booking.status)}`}>
                        {t(`bookingStatus.${booking.status}`)}
                      </span>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-3 text-sm text-gray-600">
                      <div className="flex items-center">
                        <Calendar className="w-4 h-4 mr-2 text-gray-400" />
                        <div>
                          <p className="font-medium">{t('common.from')}:</p>
                          <p>{formatDateTime(booking.date_from, i18n.language)}</p>
                        </div>
                      </div>

                      <div className="flex items-center">
                        <Calendar className="w-4 h-4 mr-2 text-gray-400" />
                        <div>
                          <p className="font-medium">{t('common.to')}:</p>
                          <p>{formatDateTime(booking.date_to, i18n.language)}</p>
                        </div>
                      </div>

                      <div className="flex items-center">
                        <DollarSign className="w-4 h-4 mr-2 text-gray-400" />
                        <div>
                          <p className="font-medium">{t('booking.totalCost')}:</p>
                          <p className="text-lg font-bold text-primary-600">
                            {booking.full_cost ? `$${(booking.full_cost / 100).toFixed(2)}` : t('booking.calculating')}
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {booking.status !== BOOKING_STATUSES.CANCELED && (
                    <div className="mt-4 md:mt-0 md:ml-4">
                      <button
                        onClick={() => handleCancelBooking(booking.booking_id)}
                        disabled={deleteLoading === booking.booking_id}
                        className="btn-danger w-full md:w-auto flex items-center justify-center space-x-2"
                      >
                        {deleteLoading === booking.booking_id ? (
                          <LoadingSpinner size="small" />
                        ) : (
                          <>
                            <Trash2 className="w-4 h-4" />
                            <span>{t('actions.cancel')}</span>
                          </>
                        )}
                      </button>
                    </div>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}

export default MyBookings
