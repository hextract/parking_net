import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { useTranslation } from 'react-i18next'
import { Car, UserPlus } from 'lucide-react'
import LoadingSpinner from '../../components/LoadingSpinner'
import LanguageSwitcher from '../../components/LanguageSwitcher'
import { validateLogin, validateEmail, validatePassword, validateTelegramId } from '../../utils/validation'

const RegisterPage = () => {
  const [formData, setFormData] = useState({
    email: '',
    login: '',
    password: '',
    confirmPassword: '',
    role: 'driver',
    telegram_id: 0,
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [validationErrors, setValidationErrors] = useState({})

  const { register, isAuthenticated, user } = useAuth()
  const { t } = useTranslation()
  const navigate = useNavigate()

  if (isAuthenticated) {
    const redirectPath = user?.role === 'driver' ? '/driver' : '/owner'
    navigate(redirectPath, { replace: true })
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    const processedValue = name === 'telegram_id'
      ? parseInt(value) || 0
      : value

    setFormData({
      ...formData,
      [name]: processedValue,
    })
    setError('')

    // Clear validation error for this field
    setValidationErrors(prev => ({ ...prev, [name]: null }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    // Validate all fields
    const errors = {}

    const emailError = validateEmail(formData.email)
    if (emailError) errors.email = t(emailError)

    const loginError = validateLogin(formData.login)
    if (loginError) errors.login = t(loginError)

    const passwordError = validatePassword(formData.password)
    if (passwordError) errors.password = t(passwordError)

    if (formData.password !== formData.confirmPassword) {
      errors.confirmPassword = t('messages.passwordMismatch')
    }

    const telegramError = validateTelegramId(formData.telegram_id)
    if (telegramError) errors.telegram_id = t(telegramError)

    if (Object.keys(errors).length > 0) {
      setValidationErrors(errors)
      setError(t('validation.pleaseFixErrors'))
      return
    }

    setLoading(true)

    const { confirmPassword, ...registrationData } = formData
    const result = await register(registrationData)

    setLoading(false)

    if (result.success) {
      const redirectPath = formData.role === 'driver' ? '/driver' : '/owner'
      navigate(redirectPath, { replace: true })
    } else {
      setError(result.error)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4 py-12">
      <div className="max-w-md w-full">
        <div className="absolute top-4 right-4">
          <LanguageSwitcher />
        </div>

        <div className="text-center mb-8">
          <div className="flex justify-center mb-4">
            <Car className="w-16 h-16 text-primary-600" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900">{t('app.name')}</h1>
          <p className="text-gray-600 mt-2">{t('auth.joinMessage')}</p>
        </div>

        <div className="bg-white rounded-lg shadow-xl p-8">
          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                {t('auth.email')}
              </label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className={`input-field ${validationErrors.email ? 'border-red-500' : ''}`}
                required
                disabled={loading}
              />
              {validationErrors.email && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.email}</p>
              )}
            </div>

            <div>
              <label htmlFor="login" className="block text-sm font-medium text-gray-700 mb-2">
                {t('auth.username')}
              </label>
              <input
                type="text"
                id="login"
                name="login"
                value={formData.login}
                onChange={handleChange}
                className={`input-field ${validationErrors.login ? 'border-red-500' : ''}`}
                required
                disabled={loading}
              />
              {validationErrors.login && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.login}</p>
              )}
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                {t('auth.password')}
              </label>
              <input
                type="password"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                className={`input-field ${validationErrors.password ? 'border-red-500' : ''}`}
                required
                disabled={loading}
              />
              {validationErrors.password && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.password}</p>
              )}
            </div>

            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-2">
                {t('auth.confirmPassword')}
              </label>
              <input
                type="password"
                id="confirmPassword"
                name="confirmPassword"
                value={formData.confirmPassword}
                onChange={handleChange}
                className={`input-field ${validationErrors.confirmPassword ? 'border-red-500' : ''}`}
                required
                disabled={loading}
              />
              {validationErrors.confirmPassword && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.confirmPassword}</p>
              )}
            </div>

            <div>
              <label htmlFor="role" className="block text-sm font-medium text-gray-700 mb-2">
                {t('auth.accountType')}
              </label>
              <select
                id="role"
                name="role"
                value={formData.role}
                onChange={handleChange}
                className="input-field"
                required
                disabled={loading}
              >
                <option value="driver">{t('roles.driver')} - {t('roles.driverDesc')}</option>
                <option value="owner">{t('roles.owner')} - {t('roles.ownerDesc')}</option>
              </select>
            </div>

            <div>
              <label htmlFor="telegram_id" className="block text-sm font-medium text-gray-700 mb-2">
                {t('common.telegramId')} ({t('common.optional')})
              </label>
              <input
                type="number"
                id="telegram_id"
                name="telegram_id"
                value={formData.telegram_id === 0 ? '' : formData.telegram_id}
                onChange={handleChange}
                placeholder={t('auth.telegramIdPlaceholder')}
                className={`input-field ${validationErrors.telegram_id ? 'border-red-500' : ''}`}
                disabled={loading}
              />
              {validationErrors.telegram_id && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.telegram_id}</p>
              )}
              <div className="mt-2 p-3 bg-blue-50 border border-blue-200 rounded-lg">
                <p className="text-xs text-blue-800">
                  ℹ️ {t('auth.telegramBotDescription')}{' '}
                  <a
                    href="https://t.me/ParkingNetRobot"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-medium text-blue-600 hover:text-blue-800 underline"
                  >
                    @ParkingNetRobot
                  </a>
                  {' '}{t('auth.telegramBotInstructions')}
                  {' '}
                  <span className="text-blue-700 font-medium">{t('auth.telegramOptional')}</span>
                </p>
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full btn-primary flex items-center justify-center space-x-2"
            >
              {loading ? (
                <LoadingSpinner size="small" />
              ) : (
                <>
                  <UserPlus className="w-4 h-4" />
                  <span>{t('auth.register')}</span>
                </>
              )}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              {t('auth.haveAccount')}{' '}
              <Link to="/login" className="text-primary-600 hover:text-primary-700 font-medium">
                {t('auth.login')}
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default RegisterPage
