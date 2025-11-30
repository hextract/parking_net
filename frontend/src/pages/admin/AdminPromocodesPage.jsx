import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Ticket, Plus, AlertCircle, CheckCircle } from 'lucide-react'
import LoadingSpinner from '../../components/LoadingSpinner'
import { createPromocode } from '../../services/paymentService'
import { validatePromocodeCode, validatePromocodeAmount, validatePromocodeMaxUses } from '../../utils/validation'

const AdminPromocodesPage = () => {
  const { t } = useTranslation()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [createdPromocodes, setCreatedPromocodes] = useState([])

  const [formData, setFormData] = useState({
    code: '',
    amount: '',
    maxUses: '1',
    expiresAt: '',
  })

  const [validationErrors, setValidationErrors] = useState({})

  const validateField = (name, value) => {
    let error = null

    switch (name) {
      case 'code':
        error = validatePromocodeCode(value)
        break
      case 'amount':
        error = validatePromocodeAmount(value)
        break
      case 'maxUses':
        error = validatePromocodeMaxUses(value)
        break
      default:
        break
    }

    return error
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
    setError('')
    setSuccess('')

    // Validate field on change
    const error = validateField(name, value)
    setValidationErrors(prev => ({
      ...prev,
      [name]: error ? t(error) : null
    }))
  }

  const validateForm = () => {
    const errors = {}

    // Validate code (optional)
    if (formData.code) {
      const codeError = validatePromocodeCode(formData.code)
      if (codeError) errors.code = t(codeError)
    }

    // Validate amount (required)
    if (!formData.amount) {
      errors.amount = t('validation.required')
    } else {
      const amountError = validatePromocodeAmount(formData.amount)
      if (amountError) errors.amount = t(amountError)
    }

    // Validate max uses (required)
    if (!formData.maxUses) {
      errors.maxUses = t('validation.required')
    } else {
      const maxUsesError = validatePromocodeMaxUses(formData.maxUses)
      if (maxUsesError) errors.maxUses = t(maxUsesError)
    }

    setValidationErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (!validateForm()) {
      setError(t('validation.pleaseFixErrors'))
      return
    }

    try {
      setLoading(true)
      setError('')
      setSuccess('')

      const data = {
        amount: Math.round(parseFloat(formData.amount) * 100), // Convert dollars to cents
        max_uses: parseInt(formData.maxUses),
      }

      // Add optional fields
      if (formData.code.trim()) {
        data.code = formData.code.trim() // Send exactly as typed (case-sensitive)
      }

      if (formData.expiresAt) {
        data.expires_at = new Date(formData.expiresAt).toISOString()
      }

      const result = await createPromocode(data)

      // Add to created list
      setCreatedPromocodes(prev => [result, ...prev])

      setSuccess(t('admin.promocode.createSuccess', { code: result.code }))

      // Reset form
      setFormData({
        code: '',
        amount: '',
        maxUses: '1',
        expiresAt: '',
      })
      setValidationErrors({})
    } catch (err) {
      console.error('Failed to create promocode:', err)
      setError(err.message || t('admin.promocode.createError'))
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (amountInCents) => {
    return `$${(amountInCents / 100).toFixed(2)}`
  }

  const formatDate = (dateString) => {
    if (!dateString) return t('common.never')
    const date = new Date(dateString)
    return date.toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center space-x-3">
        <Ticket className="w-8 h-8 text-primary-600" />
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{t('admin.promocode.title')}</h1>
          <p className="text-gray-600 mt-1">{t('admin.promocode.subtitle')}</p>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg flex items-start space-x-2">
          <AlertCircle className="w-5 h-5 flex-shrink-0 mt-0.5" />
          <span>{error}</span>
        </div>
      )}

      {success && (
        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg flex items-start space-x-2">
          <CheckCircle className="w-5 h-5 flex-shrink-0 mt-0.5" />
          <span>{success}</span>
        </div>
      )}

      {/* Create Promocode Form */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">
          {t('admin.promocode.createNew')}
        </h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Custom Code (Optional) */}
            <div>
              <label htmlFor="code" className="block text-sm font-medium text-gray-700 mb-2">
                {t('payment.promocode.code')} ({t('common.optional')})
              </label>
                  <input
                    type="text"
                    id="code"
                    name="code"
                    value={formData.code}
                    onChange={handleChange}
                    placeholder={t('admin.promocode.codeHint')}
                    className={`input-field ${validationErrors.code ? 'border-red-500' : ''}`}
                    disabled={loading}
                  />
              {validationErrors.code && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.code}</p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                {t('admin.promocode.codeDescription')}
              </p>
            </div>

            {/* Amount */}
            <div>
              <label htmlFor="amount" className="block text-sm font-medium text-gray-700 mb-2">
                {t('payment.promocode.amount')} <span className="text-red-500">*</span>
              </label>
              <div className="relative">
                <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  type="number"
                  id="amount"
                  name="amount"
                  value={formData.amount}
                  onChange={handleChange}
                  placeholder="10.00"
                  min="0.01"
                  step="0.01"
                  className={`input-field pl-8 ${validationErrors.amount ? 'border-red-500' : ''}`}
                  disabled={loading}
                  required
                />
              </div>
              {validationErrors.amount && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.amount}</p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                {t('payment.promocode.amountInCents', {
                  cents: formData.amount ? (parseFloat(formData.amount) * 100).toFixed(0) : '0'
                })}
              </p>
            </div>

            {/* Max Uses */}
            <div>
              <label htmlFor="maxUses" className="block text-sm font-medium text-gray-700 mb-2">
                {t('admin.promocode.maxUses')} <span className="text-red-500">*</span>
              </label>
              <input
                type="number"
                id="maxUses"
                name="maxUses"
                value={formData.maxUses}
                onChange={handleChange}
                min="1"
                className={`input-field ${validationErrors.maxUses ? 'border-red-500' : ''}`}
                disabled={loading}
                required
              />
              {validationErrors.maxUses && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.maxUses}</p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                {t('admin.promocode.maxUsesDescription')}
              </p>
            </div>

            {/* Expiration Date (Optional) */}
            <div>
              <label htmlFor="expiresAt" className="block text-sm font-medium text-gray-700 mb-2">
                {t('admin.promocode.expiresAt')} ({t('common.optional')})
              </label>
              <input
                type="datetime-local"
                id="expiresAt"
                name="expiresAt"
                value={formData.expiresAt}
                onChange={handleChange}
                className="input-field"
                disabled={loading}
              />
              <p className="mt-1 text-xs text-gray-500">
                {t('admin.promocode.expiresAtDescription')}
              </p>
            </div>
          </div>

          <button
            type="submit"
            disabled={loading || Object.keys(validationErrors).some(key => validationErrors[key])}
            className="btn-primary flex items-center space-x-2"
          >
            {loading ? (
              <LoadingSpinner size="small" />
            ) : (
              <>
                <Plus className="w-4 h-4" />
                <span>{t('admin.promocode.create')}</span>
              </>
            )}
          </button>
        </form>
      </div>

      {/* Created Promocodes List */}
      {createdPromocodes.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            {t('admin.promocode.recentlyCreated')}
          </h2>
          <div className="space-y-3">
            {createdPromocodes.map((promo, index) => (
              <div
                key={index}
                className="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center space-x-2 mb-2">
                      <code className="px-3 py-1 bg-primary-100 text-primary-900 rounded font-mono font-bold text-lg">
                        {promo.code}
                      </code>
                      <button
                        onClick={() => {
                          navigator.clipboard.writeText(promo.code)
                          setSuccess(t('common.copiedToClipboard'))
                          setTimeout(() => setSuccess(''), 2000)
                        }}
                        className="text-xs text-primary-600 hover:text-primary-700 px-2 py-1 hover:bg-primary-50 rounded"
                      >
                        {t('actions.copy')}
                      </button>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-2 text-sm">
                      <div>
                        <span className="text-gray-600">{t('payment.promocode.value')}:</span>
                        <span className="ml-1 font-medium">{formatCurrency(promo.amount)}</span>
                      </div>
                      <div>
                        <span className="text-gray-600">{t('admin.promocode.maxUses')}:</span>
                        <span className="ml-1 font-medium">{promo.max_uses}</span>
                      </div>
                      <div>
                        <span className="text-gray-600">{t('payment.promocode.remainingUses')}:</span>
                        <span className="ml-1 font-medium">{promo.remaining_uses}</span>
                      </div>
                      <div>
                        <span className="text-gray-600">{t('payment.promocode.expiresAt')}:</span>
                        <span className="ml-1 font-medium text-xs">{formatDate(promo.expires_at)}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default AdminPromocodesPage
