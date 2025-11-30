import { useState, useEffect } from 'react'
import { Search, MapPin, DollarSign, Car as CarIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { parkingService } from '../../services/parkingService'
import { PARKING_TYPES } from '../../config/api'
import LoadingSpinner from '../../components/LoadingSpinner'

const AllParkings = () => {
  const { t } = useTranslation()
  const [filters, setFilters] = useState({
    city: '',
    name: '',
    parking_type: '',
  })
  const [parkings, setParkings] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    loadAllParkings()
  }, [])

  const loadAllParkings = async () => {
    setLoading(true)
    setError('')
    try {
      const data = await parkingService.getParkings({})
      setParkings(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err.message || t('messages.loadFailed'))
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

  const handleSearch = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const data = await parkingService.getParkings(filters)
      setParkings(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err.message || t('messages.loadFailed'))
      setParkings([])
    } finally {
      setLoading(false)
    }
  }

  const getParkingTypeLabel = (type) => {
    return t(`parkingTypes.${type}`)
  }

  return (
    <div className="max-w-7xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">{t('nav.allParkings')}</h1>

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

      {/* Error Message */}
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

              <div className="space-y-2">
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
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default AllParkings
