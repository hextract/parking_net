import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useState, useEffect } from 'react'
import { Search, Calendar, MapPin, TrendingUp, Wallet, Ticket } from 'lucide-react'
import { useAuth } from '../../context/AuthContext'
import { bookingService } from '../../services/bookingService'
import { getBalance } from '../../services/paymentService'

const DriverDashboard = () => {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { t } = useTranslation()
  const [totalBookings, setTotalBookings] = useState(0)
  const [activeBookings, setActiveBookings] = useState(0)
  const [totalSpent, setTotalSpent] = useState(0)
  const [balance, setBalance] = useState(0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadDashboardData()
  }, [user])

  const loadDashboardData = async () => {
    if (!user?.user_id) return

    try {
      // Load all bookings for the driver
      const [bookings, balanceData] = await Promise.all([
        bookingService.getBookings({ user_id: user.user_id }),
        getBalance().catch(() => ({ balance: 0 }))
      ])

      // Set balance
      if (balanceData?.balance !== undefined) {
        setBalance((balanceData.balance / 100).toFixed(2))
      }

      if (Array.isArray(bookings)) {
        setTotalBookings(bookings.length)

        // Count active bookings (Confirmed or Waiting)
        const active = bookings.filter(b => b.status === 'Confirmed' || b.status === 'Waiting').length
        setActiveBookings(active)

        // Calculate total spent from confirmed bookings
        let spent = 0
        const confirmedBookings = bookings.filter(b => b.status === 'Confirmed')
        for (const booking of confirmedBookings) {
          if (booking.full_cost) {
            spent += booking.full_cost
          }
        }
        setTotalSpent((spent / 100).toFixed(2))
      }
    } catch (err) {
      console.error('Failed to load dashboard data:', err)
    } finally {
      setLoading(false)
    }
  }

  const quickActions = [
    {
      title: t('driver.searchParking'),
      description: t('driver.searchDesc'),
      icon: Search,
      color: 'bg-blue-500',
      action: () => navigate('/driver/search'),
    },
    {
      title: t('driver.myBookings'),
      description: t('driver.myBookingsDesc'),
      icon: Calendar,
      color: 'bg-green-500',
      action: () => navigate('/driver/bookings'),
    },
    {
      title: t('payment.balance.title'),
      description: t('payment.transactions.title'),
      icon: Wallet,
      color: 'bg-purple-500',
      action: () => navigate('/driver/balance'),
    },
    {
      title: t('payment.promocode.title'),
      description: t('payment.promocode.activateDescription'),
      icon: Ticket,
      color: 'bg-yellow-500',
      action: () => navigate('/driver/promocodes'),
    },
  ]

  const stats = [
    {
      label: t('driver.totalBookings'),
      value: loading ? '-' : totalBookings,
      icon: Calendar,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
    },
    {
      label: t('driver.activeBookings'),
      value: loading ? '-' : activeBookings,
      icon: MapPin,
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
      label: t('driver.totalSpent'),
      value: loading ? '-' : `$${totalSpent}`,
      icon: TrendingUp,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
    },
  ]

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">{t('driver.welcome', { name: user?.login })}</h1>
        <p className="text-gray-600 mt-2">{t('driver.tagline')}</p>
      </div>

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

      <div className="mt-8 bg-primary-50 border border-primary-200 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-primary-900 mb-2">
          {t('driver.howToStart')}
        </h3>
        <ul className="space-y-2 text-primary-800">
          <li className="flex items-start space-x-2">
            <span className="font-bold">1.</span>
            <span>{t('driver.step1')}</span>
          </li>
          <li className="flex items-start space-x-2">
            <span className="font-bold">2.</span>
            <span>{t('driver.step2')}</span>
          </li>
          <li className="flex items-start space-x-2">
            <span className="font-bold">3.</span>
            <span>{t('driver.step3')}</span>
          </li>
        </ul>
      </div>
    </div>
  )
}

export default DriverDashboard
