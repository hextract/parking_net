import { useNavigate } from 'react-router-dom'
import { MapPin, Calendar, Plus, TrendingUp, Wallet, Ticket } from 'lucide-react'
import { useAuth } from '../../context/AuthContext'
import { useTranslation } from 'react-i18next'
import { useState, useEffect } from 'react'
import { parkingService } from '../../services/parkingService'
import { bookingService } from '../../services/bookingService'
import { getBalance } from '../../services/paymentService'

const OwnerDashboard = () => {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { t } = useTranslation()
  const [totalParkings, setTotalParkings] = useState(0)
  const [activeBookings, setActiveBookings] = useState(0)
  const [totalRevenue, setTotalRevenue] = useState(0)
  const [balance, setBalance] = useState(0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadDashboardData()
  }, [user])

  const loadDashboardData = async () => {
    if (!user?.user_id) return

    try {
      // Load parkings count and balance
      const [parkings, balanceData] = await Promise.all([
        parkingService.getParkings({ owner_id: user.user_id }),
        getBalance().catch(() => ({ balance: 0 }))
      ])

      // Set balance
      if (balanceData?.balance !== undefined) {
        setBalance((balanceData.balance / 100).toFixed(2))
      }

      setTotalParkings(Array.isArray(parkings) ? parkings.length : 0)

      // Load bookings for all owner's parkings and calculate revenue
      if (Array.isArray(parkings) && parkings.length > 0) {
        let totalActive = 0
        let revenue = 0

        for (const parking of parkings) {
          try {
            const bookings = await bookingService.getBookings({ parking_place_id: parking.id })
            if (Array.isArray(bookings)) {
              // Count active bookings
              const active = bookings.filter(b => b.status === 'Confirmed' || b.status === 'Waiting').length
              totalActive += active

              // Calculate revenue from confirmed bookings
              const confirmedBookings = bookings.filter(b => b.status === 'Confirmed')
              for (const booking of confirmedBookings) {
                if (booking.full_cost) {
                  revenue += booking.full_cost
                }
              }
            }
          } catch (err) {
            // Ignore errors for individual parkings
          }
        }
        setActiveBookings(totalActive)
        setTotalRevenue((revenue / 100).toFixed(2))
      }
    } catch (err) {
      console.error('Failed to load dashboard data:', err)
    } finally {
      setLoading(false)
    }
  }

  const quickActions = [
    {
      title: t('nav.allParkings'),
      description: t('owner.allParkingsDesc'),
      icon: MapPin,
      color: 'bg-blue-500',
      action: () => navigate('/owner/all-parkings'),
    },
    {
      title: t('owner.myParkings'),
      description: t('owner.myParkingsDesc'),
      icon: MapPin,
      color: 'bg-green-500',
      action: () => navigate('/owner/parkings'),
    },
    {
      title: t('payment.balance.title'),
      description: t('payment.transactions.title'),
      icon: Wallet,
      color: 'bg-purple-500',
      action: () => navigate('/owner/balance'),
    },
    {
      title: t('payment.promocode.title'),
      description: t('payment.promocode.generateDescription'),
      icon: Ticket,
      color: 'bg-yellow-500',
      action: () => navigate('/owner/promocodes'),
    },
  ]

  const stats = [
    {
      label: t('owner.totalParkings'),
      value: loading ? '-' : totalParkings,
      icon: MapPin,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
    },
    {
      label: t('owner.activeBookings'),
      value: loading ? '-' : activeBookings,
      icon: Calendar,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
    },
    {
      label: t('payment.balance.current'),
      value: loading ? '-' : `$${balance}`,
      icon: Wallet,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
    },
    {
      label: t('owner.totalRevenue'),
      value: loading ? '-' : `$${totalRevenue}`,
      icon: TrendingUp,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
    },
  ]

  return (
    <div className="max-w-7xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">{t('owner.welcome', { name: user?.login })}</h1>
        <p className="text-gray-600 mt-2">{t('owner.tagline')}</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat, index) => {
          const Icon = stat.icon
          return (
            <div key={index} className="card">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600 mb-1">{stat.label}</p>
                  <p className="text-3xl font-bold text-gray-900">{stat.value}</p>
                </div>
                <div className={`${stat.bgColor} p-3 rounded-lg`}>
                  <Icon className={`w-6 h-6 ${stat.color}`} />
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* Quick Actions */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 mb-4">{t('common.quickActions')}</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {quickActions.map((action, index) => {
            const Icon = action.icon
            return (
              <button
                key={index}
                onClick={action.action}
                className="card text-left hover:shadow-xl transition-all transform hover:-translate-y-1"
              >
                <div className="flex items-start space-x-4">
                  <div className={`${action.color} p-3 rounded-lg`}>
                    <Icon className="w-6 h-6 text-white" />
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-1">
                      {action.title}
                    </h3>
                    <p className="text-gray-600">{action.description}</p>
                  </div>
                </div>
              </button>
            )
          })}
        </div>
      </div>

      {/* Info Section */}
      <div className="mt-8 bg-primary-50 border border-primary-200 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-primary-900 mb-2">
          {t('owner.gettingStarted')}
        </h3>
        <ul className="space-y-2 text-primary-800">
          <li className="flex items-start space-x-2">
            <span className="font-bold">1.</span>
            <span>{t('owner.ownerStep1')}</span>
          </li>
          <li className="flex items-start space-x-2">
            <span className="font-bold">2.</span>
            <span>{t('owner.ownerStep2')}</span>
          </li>
          <li className="flex items-start space-x-2">
            <span className="font-bold">3.</span>
            <span>{t('owner.ownerStep3')}</span>
          </li>
        </ul>
      </div>
    </div>
  )
}

export default OwnerDashboard
