import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useTranslation } from 'react-i18next'
import { Car, LogOut, Menu, X, Calendar, MapPin, Settings, Wallet, Ticket } from 'lucide-react'
import { useState } from 'react'
import LanguageSwitcher from './LanguageSwitcher'

const Navbar = () => {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const { t } = useTranslation()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const isDriver = user?.role === 'driver'
  const isOwner = user?.role === 'owner'
  const isAdmin = user?.role === 'admin'

  const driverLinks = [
    { to: '/driver', label: t('nav.dashboard'), icon: Car },
    { to: '/driver/search', label: t('nav.search'), icon: MapPin },
    { to: '/driver/bookings', label: t('nav.myBookings'), icon: Calendar },
    { to: '/driver/balance', label: t('payment.balance.title'), icon: Wallet },
    { to: '/driver/promocodes', label: t('payment.promocode.title'), icon: Ticket },
  ]

  const ownerLinks = [
    { to: '/owner', label: t('nav.dashboard'), icon: Car },
    { to: '/owner/all-parkings', label: t('nav.allParkings'), icon: MapPin },
    { to: '/owner/parkings', label: t('nav.myParkings'), icon: MapPin },
    { to: '/owner/balance', label: t('payment.balance.title'), icon: Wallet },
    { to: '/owner/promocodes', label: t('payment.promocode.title'), icon: Ticket },
  ]

  // Admin links
  const adminLinks = [
    { to: '/admin', label: t('nav.admin'), icon: Settings },
    { to: '/admin/promocodes', label: t('payment.promocode.title'), icon: Ticket },
  ]

  let links = isDriver ? driverLinks : isOwner ? ownerLinks : []
  // Only admins see the admin links
  if (isAdmin) {
    links = [...links, ...adminLinks]
  }

  // Determine home path based on role
  const getHomePath = () => {
    if (isDriver) return '/driver'
    if (isOwner) return '/owner'
    if (isAdmin) return '/admin'
    return '/login'
  }

  return (
    <nav className="bg-white shadow-lg">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <Link to={getHomePath()} className="flex items-center space-x-2">
            <Car className="w-8 h-8 text-primary-600" />
            <span className="text-xl font-bold text-gray-900">{t('app.name')}</span>
          </Link>

          <div className="hidden md:flex items-center space-x-4">
            {links.map((link) => {
              const Icon = link.icon
              return (
                <Link
                  key={link.to}
                  to={link.to}
                  className="flex items-center space-x-1 px-3 py-2 rounded-md text-gray-700 hover:bg-primary-50 hover:text-primary-600 transition-colors"
                >
                  <Icon className="w-4 h-4" />
                  <span>{link.label}</span>
                </Link>
              )
            })}

            <div className="flex items-center space-x-3 ml-4 pl-4 border-l border-gray-200">
              <LanguageSwitcher />
              <div className="text-sm">
                <div className="font-medium text-gray-900">{user?.login}</div>
                <div className="text-gray-500 capitalize">{t(`roles.${user?.role}`)}</div>
              </div>
              <button
                onClick={handleLogout}
                className="flex items-center space-x-1 px-3 py-2 rounded-md text-red-600 hover:bg-red-50 transition-colors"
              >
                <LogOut className="w-4 h-4" />
                <span>{t('auth.logout')}</span>
              </button>
            </div>
          </div>

          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="md:hidden p-2 rounded-md text-gray-700 hover:bg-gray-100"
          >
            {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>

        {mobileMenuOpen && (
          <div className="md:hidden py-4 border-t border-gray-200">
            <div className="flex flex-col space-y-2">
              {links.map((link) => {
                const Icon = link.icon
                return (
                  <Link
                    key={link.to}
                    to={link.to}
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex items-center space-x-2 px-3 py-2 rounded-md text-gray-700 hover:bg-primary-50 hover:text-primary-600 transition-colors"
                  >
                    <Icon className="w-4 h-4" />
                    <span>{link.label}</span>
                  </Link>
                )
              })}

              <div className="pt-4 mt-4 border-t border-gray-200">
                <div className="px-3 pb-3">
                  <LanguageSwitcher />
                </div>
                <div className="px-3 py-2 text-sm">
                  <div className="font-medium text-gray-900">{user?.login}</div>
                  <div className="text-gray-500 capitalize">{t(`roles.${user?.role}`)}</div>
                </div>
                <button
                  onClick={handleLogout}
                  className="w-full flex items-center space-x-2 px-3 py-2 rounded-md text-red-600 hover:bg-red-50 transition-colors"
                >
                  <LogOut className="w-4 h-4" />
                  <span>{t('auth.logout')}</span>
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </nav>
  )
}

export default Navbar
