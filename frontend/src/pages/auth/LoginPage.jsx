import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { useTranslation } from 'react-i18next'
import { Car, LogIn } from 'lucide-react'
import LoadingSpinner from '../../components/LoadingSpinner'
import LanguageSwitcher from '../../components/LanguageSwitcher'
import { validateLogin, validatePassword } from '../../utils/validation'

const LoginPage = () => {
  const [formData, setFormData] = useState({
    login: '',
    password: '',
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [validationErrors, setValidationErrors] = useState({})

  const { login, isAuthenticated, user } = useAuth()
  const { t } = useTranslation()
  const navigate = useNavigate()

  if (isAuthenticated) {
    const redirectPath = user?.role === 'driver' ? '/driver' : '/owner'
    navigate(redirectPath, { replace: true })
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData({
      ...formData,
      [name]: value,
    })
    setError('')
    // Clear validation error for this field
    setValidationErrors(prev => ({ ...prev, [name]: null }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    // Validate fields
    const errors = {}

    const loginError = validateLogin(formData.login)
    if (loginError) errors.login = t(loginError)

    const passwordError = validatePassword(formData.password)
    if (passwordError) errors.password = t(passwordError)

    if (Object.keys(errors).length > 0) {
      setValidationErrors(errors)
      setError(t('validation.pleaseFixErrors'))
      return
    }
    setLoading(true)

    const result = await login(formData)

    setLoading(false)

    if (result.success) {
      const redirectPath = result.user?.role === 'driver' ? '/driver' : '/owner'
      navigate(redirectPath, { replace: true })
    } else {
      setError(result.error)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4">
      <div className="max-w-md w-full">
        <div className="absolute top-4 right-4">
          <LanguageSwitcher />
        </div>

        <div className="text-center mb-8">
          <div className="flex justify-center mb-4">
            <Car className="w-16 h-16 text-primary-600" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900">{t('auth.welcomeBack')}</h1>
          <p className="text-gray-600 mt-2">{t('auth.signInMessage')}</p>
        </div>

        <div className="bg-white rounded-lg shadow-xl p-8">
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

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

            <button
              type="submit"
              disabled={loading}
              className="w-full btn-primary flex items-center justify-center space-x-2"
            >
              {loading ? (
                <LoadingSpinner size="small" />
              ) : (
                <>
                  <LogIn className="w-4 h-4" />
                  <span>{t('auth.login')}</span>
                </>
              )}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              {t('auth.noAccount')}{' '}
              <Link to="/register" className="text-primary-600 hover:text-primary-700 font-medium">
                {t('auth.register')}
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LoginPage
