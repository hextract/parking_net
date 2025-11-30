import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Ticket, Plus, Search, Gift, Info } from 'lucide-react'
import LoadingSpinner from '../../components/LoadingSpinner'
import { activatePromocode, generatePromocode, getPromocodeInfo } from '../../services/paymentService'

const PromocodesPage = () => {
  const { t } = useTranslation()
  const [activeTab, setActiveTab] = useState('activate') // 'activate', 'generate', 'check'
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  // Activate tab state
  const [activateCode, setActivateCode] = useState('')
  const [activateResult, setActivateResult] = useState(null)

  // Generate tab state
  const [generateAmount, setGenerateAmount] = useState('')
  const [generateResult, setGenerateResult] = useState(null)

  // Check tab state
  const [checkCode, setCheckCode] = useState('')
  const [checkResult, setCheckResult] = useState(null)

  const handleActivate = async (e) => {
    e.preventDefault()
    if (!activateCode.trim()) return

    try {
      setLoading(true)
      setError('')
      setSuccess('')

      // Send exactly as user typed (case-sensitive)
      const result = await activatePromocode(activateCode.trim())
      setActivateResult(result)
      setSuccess(t('payment.promocode.activateSuccess', {
        balance: `$${(result.balance / 100).toFixed(2)}`
      }))
      setActivateCode('')
    } catch (err) {
      console.error('Failed to activate promocode:', err)
      setError(err.message || t('payment.promocode.activateError'))
    } finally {
      setLoading(false)
    }
  }

  const handleGenerate = async (e) => {
    e.preventDefault()
    const amountInDollars = parseFloat(generateAmount)
    if (!amountInDollars || amountInDollars <= 0) {
      setError(t('payment.promocode.invalidAmount'))
      return
    }

    try {
      setLoading(true)
      setError('')
      setSuccess('')

      // Convert dollars to cents for backend
      const amountInCents = Math.round(amountInDollars * 100)
      const result = await generatePromocode(amountInCents)
      setGenerateResult(result)
      setSuccess(t('payment.promocode.generateSuccess', { code: result.code }))
      setGenerateAmount('')
    } catch (err) {
      console.error('Failed to generate promocode:', err)
      setError(err.message || t('payment.promocode.generateError'))
    } finally {
      setLoading(false)
    }
  }

  const handleCheck = async (e) => {
    e.preventDefault()
    if (!checkCode.trim()) return

    try {
      setLoading(true)
      setError('')
      setSuccess('')

      // Send exactly as user typed (case-sensitive)
      const result = await getPromocodeInfo(checkCode.trim())
      setCheckResult(result)
    } catch (err) {
      console.error('Failed to check promocode:', err)
      setError(err.message || t('payment.promocode.checkError'))
      setCheckResult(null)
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
    <div className="space-y-6">
      <div className="flex items-center space-x-3">
        <Ticket className="w-8 h-8 text-primary-600" />
        <h1 className="text-3xl font-bold text-gray-900">{t('payment.promocode.title')}</h1>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg">
          {success}
        </div>
      )}

      {/* Tabs */}
      <div className="bg-white rounded-lg shadow">
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8 px-6" aria-label="Tabs">
            <button
              onClick={() => {
                setActiveTab('activate')
                setError('')
                setSuccess('')
              }}
              className={`${
                activeTab === 'activate'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center space-x-2`}
            >
              <Plus className="w-4 h-4" />
              <span>{t('payment.promocode.activate')}</span>
            </button>
            <button
              onClick={() => {
                setActiveTab('generate')
                setError('')
                setSuccess('')
              }}
              className={`${
                activeTab === 'generate'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center space-x-2`}
            >
              <Gift className="w-4 h-4" />
              <span>{t('payment.promocode.generate')}</span>
            </button>
            <button
              onClick={() => {
                setActiveTab('check')
                setError('')
                setSuccess('')
              }}
              className={`${
                activeTab === 'check'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center space-x-2`}
            >
              <Search className="w-4 h-4" />
              <span>{t('payment.promocode.check')}</span>
            </button>
          </nav>
        </div>

        <div className="p-6">
          {/* Activate Tab */}
          {activeTab === 'activate' && (
            <div className="max-w-md">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                {t('payment.promocode.activateTitle')}
              </h3>
              <p className="text-gray-600 mb-6">{t('payment.promocode.activateDescription')}</p>

              <form onSubmit={handleActivate} className="space-y-4">
                <div>
                  <label htmlFor="activate-code" className="block text-sm font-medium text-gray-700 mb-2">
                    {t('payment.promocode.code')}
                  </label>
                  <input
                    type="text"
                    id="activate-code"
                    value={activateCode}
                    onChange={(e) => setActivateCode(e.target.value)}
                    placeholder={t('payment.promocode.codePlaceholder')}
                    className="input-field"
                    disabled={loading}
                  />
                </div>

                <button
                  type="submit"
                  disabled={loading || !activateCode.trim()}
                  className="btn-primary w-full"
                >
                  {loading ? <LoadingSpinner size="small" /> : t('payment.promocode.activate')}
                </button>
              </form>
            </div>
          )}

          {/* Generate Tab */}
          {activeTab === 'generate' && (
            <div className="max-w-md">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                {t('payment.promocode.generateTitle')}
              </h3>
              <p className="text-gray-600 mb-6">{t('payment.promocode.generateDescription')}</p>

              <form onSubmit={handleGenerate} className="space-y-4">
                <div>
                  <label htmlFor="generate-amount" className="block text-sm font-medium text-gray-700 mb-2">
                    {t('payment.promocode.amount')}
                  </label>
                  <div className="relative">
                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                    <input
                      type="number"
                      id="generate-amount"
                      value={generateAmount}
                      onChange={(e) => setGenerateAmount(e.target.value)}
                      placeholder="10.00"
                      min="1"
                      step="0.01"
                      className="input-field pl-8"
                      disabled={loading}
                    />
                  </div>
                  <p className="mt-1 text-xs text-gray-500">
                    {t('payment.promocode.amountInCents', {
                      cents: generateAmount ? (parseFloat(generateAmount) * 100).toFixed(0) : '0'
                    })}
                  </p>
                </div>

                <button
                  type="submit"
                  disabled={loading || !generateAmount}
                  className="btn-primary w-full"
                >
                  {loading ? <LoadingSpinner size="small" /> : t('payment.promocode.generate')}
                </button>
              </form>

              {generateResult && (
                <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg">
                  <p className="text-sm text-green-800 mb-2">{t('payment.promocode.generatedCode')}:</p>
                  <div className="flex items-center space-x-2">
                    <code className="flex-1 px-3 py-2 bg-white border border-green-300 rounded text-lg font-mono font-bold text-green-900">
                      {generateResult.code}
                    </code>
                    <button
                      onClick={() => navigator.clipboard.writeText(generateResult.code)}
                      className="btn-secondary text-sm"
                    >
                      {t('actions.copy')}
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Check Tab */}
          {activeTab === 'check' && (
            <div className="max-w-md">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                {t('payment.promocode.checkTitle')}
              </h3>
              <p className="text-gray-600 mb-6">{t('payment.promocode.checkDescription')}</p>

              <form onSubmit={handleCheck} className="space-y-4">
                <div>
                  <label htmlFor="check-code" className="block text-sm font-medium text-gray-700 mb-2">
                    {t('payment.promocode.code')}
                  </label>
                  <input
                    type="text"
                    id="check-code"
                    value={checkCode}
                    onChange={(e) => setCheckCode(e.target.value)}
                    placeholder={t('payment.promocode.codePlaceholder')}
                    className="input-field"
                    disabled={loading}
                  />
                </div>

                <button
                  type="submit"
                  disabled={loading || !checkCode.trim()}
                  className="btn-primary w-full"
                >
                  {loading ? <LoadingSpinner size="small" /> : t('payment.promocode.check')}
                </button>
              </form>

              {checkResult && (
                <div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg space-y-3">
                  <div className="flex items-start space-x-2">
                    <Info className="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
                    <div className="flex-1">
                      <h4 className="font-semibold text-blue-900 mb-3">{t('payment.promocode.info')}</h4>

                      <dl className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <dt className="text-blue-800 font-medium">{t('payment.promocode.code')}:</dt>
                          <dd className="font-mono text-blue-900">{checkResult.code}</dd>
                        </div>
                        <div className="flex justify-between">
                          <dt className="text-blue-800 font-medium">{t('payment.promocode.value')}:</dt>
                          <dd className="text-blue-900">{formatCurrency(checkResult.amount)}</dd>
                        </div>
                        <div className="flex justify-between">
                          <dt className="text-blue-800 font-medium">{t('payment.promocode.remainingUses')}:</dt>
                          <dd className="text-blue-900">{checkResult.remaining_uses} / {checkResult.max_uses}</dd>
                        </div>
                        <div className="flex justify-between">
                          <dt className="text-blue-800 font-medium">{t('payment.promocode.expiresAt')}:</dt>
                          <dd className="text-blue-900">{formatDate(checkResult.expires_at)}</dd>
                        </div>
                        <div className="flex justify-between">
                          <dt className="text-blue-800 font-medium">{t('payment.promocode.status')}:</dt>
                          <dd>
                            <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                              checkResult.is_active
                                ? 'bg-green-100 text-green-800'
                                : 'bg-red-100 text-red-800'
                            }`}>
                              {checkResult.is_active ? t('payment.promocode.active') : t('payment.promocode.inactive')}
                            </span>
                          </dd>
                        </div>
                      </dl>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default PromocodesPage
