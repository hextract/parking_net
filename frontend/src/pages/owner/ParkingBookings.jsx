import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Calendar, MapPin, User, DollarSign, ArrowLeft } from 'lucide-react'
import { bookingService } from '../../services/bookingService'
import { parkingService } from '../../services/parkingService'
import { BOOKING_STATUSES } from '../../config/api'
import LoadingSpinner from '../../components/LoadingSpinner'
import { formatDateTime } from '../../utils/dateUtils'

const ParkingBookings = () => {
  const { t, i18n } = useTranslation()
  const { parkingId } = useParams()
  const navigate = useNavigate()
  const [parking, setParking] = useState(null)
  const [bookings, setBookings] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [updateLoading, setUpdateLoading] = useState(null)

  useEffect(() => {
    loadData()
  }, [parkingId])

  const loadData = async () => {
    setLoading(true)
    setError('')
    try {
      // Load parking details
      const parkingData = await parkingService.getParkingById(parkingId)
      setParking(parkingData)

      // Load bookings for this parking
      const bookingsData = await bookingService.getBookings({
        parking_place_id: parkingId,
      })
      setBookings(Array.isArray(bookingsData) ? bookingsData : [])
    } catch (err) {
      setError(err.message || 'Failed to load data')
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateStatus = async (bookingId, newStatus) => {
    setUpdateLoading(bookingId)
    setError('')
    try {
      const booking = bookings.find((b) => b.booking_id === bookingId)
      await bookingService.updateBooking(bookingId, {
        ...booking,
        status: newStatus,
      })

      // Update local state
      setBookings(
        bookings.map((b) =>
          b.booking_id === bookingId ? { ...b, status: newStatus } : b
        )
      )
    } catch (err) {
      setError(err.message || 'Failed to update booking status')
    } finally {
      setUpdateLoading(null)
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
      <button
        onClick={() => navigate('/owner/parkings')}
        className="flex items-center space-x-2 text-gray-600 hover:text-gray-900 mb-6"
      >
        <ArrowLeft className="w-4 h-4" />
        <span>{t('actions.backTo', { page: t('owner.myParkings') })}</span>
      </button>

      {parking && (
        <div className="card mb-8">
          <div className="flex items-start justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">{parking.name}</h1>
              <p className="text-gray-600 flex items-center">
                <MapPin className="w-4 h-4 mr-1" />
                {parking.city}, {parking.address}
              </p>
            </div>
            <div className="text-right">
              <p className="text-sm text-gray-600">{t('booking.hourlyRate')}</p>
              <p className="text-2xl font-bold text-primary-600">${(parking.hourly_rate / 100).toFixed(2)}</p>
            </div>
          </div>
        </div>
      )}

      <h2 className="text-xl font-semibold text-gray-900 mb-4">
        {t('actions.bookings')} ({bookings.length})
      </h2>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
          {error}
        </div>
      )}

      {bookings.length === 0 ? (
        <div className="text-center py-12">
          <Calendar className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t('booking.noBookingsForParking')}</h3>
          <p className="text-gray-600">{t('booking.bookingsWillAppear')}</p>
        </div>
      ) : (
        <div className="space-y-4">
          {bookings.map((booking) => (
            <div key={booking.booking_id} className="card">
              <div className="flex flex-col md:flex-row md:items-center md:justify-between">
                <div className="flex-1">
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <div className="flex items-center space-x-2 mb-2">
                        <span className={`badge ${getStatusBadgeClass(booking.status)}`}>
                          {t(`bookingStatus.${booking.status}`)}
                        </span>
                        <span className="text-sm text-gray-600">
                          {t('booking.bookingId')} #{booking.booking_id}
                        </span>
                      </div>
                      <p className="text-sm text-gray-600 flex items-center">
                        <User className="w-4 h-4 mr-1" />
                        {t('booking.customer')}: {booking.user_id}
                      </p>
                    </div>
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

                {booking.status === BOOKING_STATUSES.WAITING && (
                  <div className="mt-4 md:mt-0 md:ml-4 flex flex-col space-y-2">
                    <button
                      onClick={() =>
                        handleUpdateStatus(booking.booking_id, BOOKING_STATUSES.CONFIRMED)
                      }
                      disabled={updateLoading === booking.booking_id}
                      className="btn-primary whitespace-nowrap"
                    >
                      {updateLoading === booking.booking_id ? (
                        <LoadingSpinner size="small" />
                      ) : (
                        t('booking.confirm')
                      )}
                    </button>
                    <button
                      onClick={() =>
                        handleUpdateStatus(booking.booking_id, BOOKING_STATUSES.CANCELED)
                      }
                      disabled={updateLoading === booking.booking_id}
                      className="btn-danger whitespace-nowrap"
                    >
                      {updateLoading === booking.booking_id ? (
                        <LoadingSpinner size="small" />
                      ) : (
                        t('booking.cancel')
                      )}
                    </button>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default ParkingBookings
