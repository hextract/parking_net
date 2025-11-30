import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../../context/AuthContext'
import { Wallet, TrendingDown, TrendingUp, Clock, RefreshCw } from 'lucide-react'
import LoadingSpinner from '../../components/LoadingSpinner'
import { getBalance, getTransactions } from '../../services/paymentService'

const BalancePage = () => {
  const { t } = useTranslation()
  const { user } = useAuth()
  const [balance, setBalance] = useState(null)
  const [transactions, setTransactions] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [refreshing, setRefreshing] = useState(false)

  const loadData = async (showLoader = true) => {
    try {
      if (showLoader) {
        setLoading(true)
      } else {
        setRefreshing(true)
      }
      setError('')

      const [balanceData, transactionsData] = await Promise.all([
        getBalance(),
        getTransactions(50, 0)
      ])

      setBalance(balanceData)
      setTransactions(transactionsData || [])
    } catch (err) {
      console.error('Failed to load payment data:', err)
      setError(err.message || t('messages.loadError'))
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  const formatCurrency = (amountInCents) => {
    const amount = (amountInCents / 100).toFixed(2)
    return `$${amount}`
  }

  const formatDate = (dateString) => {
    if (!dateString) return '-'
    const date = new Date(dateString)
    return date.toLocaleString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const getTransactionIcon = (type) => {
    switch (type) {
      case 'charge':
      case 'promocode_generate':
        return <TrendingDown className="w-5 h-5 text-red-500" />
      case 'payment':
      case 'refund':
      case 'promocode_activate':
        return <TrendingUp className="w-5 h-5 text-green-500" />
      default:
        return <Clock className="w-5 h-5 text-gray-400" />
    }
  }

  const getTransactionColor = (type) => {
    switch (type) {
      case 'charge':
      case 'promocode_generate':
        return 'text-red-600'
      case 'payment':
      case 'refund':
      case 'promocode_activate':
        return 'text-green-600'
      default:
        return 'text-gray-600'
    }
  }

  const getStatusBadge = (status) => {
    const colors = {
      completed: 'bg-green-100 text-green-800',
      pending: 'bg-yellow-100 text-yellow-800',
      failed: 'bg-red-100 text-red-800',
      canceled: 'bg-gray-100 text-gray-800',
    }

    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${colors[status] || colors.pending}`}>
        {t(`payment.transactionStatus.${status}`) || status}
      </span>
    )
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <LoadingSpinner />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      )}

      {/* Balance Card */}
      <div className="bg-gradient-to-br from-primary-500 to-primary-700 rounded-xl shadow-lg p-8 text-white">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-3">
            <Wallet className="w-8 h-8" />
            <h2 className="text-2xl font-bold">{t('payment.balance.title')}</h2>
          </div>
          <button
            onClick={() => loadData(false)}
            disabled={refreshing}
            className="p-2 hover:bg-white/20 rounded-lg transition-colors disabled:opacity-50"
            title={t('actions.refresh')}
          >
            <RefreshCw className={`w-5 h-5 ${refreshing ? 'animate-spin' : ''}`} />
          </button>
        </div>

        {balance && (
          <div>
            <p className="text-white/80 text-sm mb-2">{t('payment.balance.current')}</p>
            <p className="text-5xl font-bold mb-2">
              {formatCurrency(balance.balance)}
            </p>
            <p className="text-white/60 text-sm">{balance.currency || 'USD'}</p>
          </div>
        )}
      </div>

      {/* Transactions */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-xl font-semibold text-gray-900">{t('payment.transactions.title')}</h3>
          <span className="text-sm text-gray-500">
            {t('payment.transactions.total', { count: transactions.length })}
          </span>
        </div>

        {transactions.length === 0 ? (
          <div className="text-center py-12">
            <Clock className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-500">{t('payment.transactions.empty')}</p>
          </div>
        ) : (
          <div className="space-y-4">
            {transactions.map((transaction) => (
              <div
                key={transaction.id}
                className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-center space-x-4 flex-1">
                  <div className="flex-shrink-0">
                    {getTransactionIcon(transaction.transaction_type)}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-2 mb-1">
                      <p className="font-medium text-gray-900">
                        {t(`payment.transactionType.${transaction.transaction_type}`) || transaction.transaction_type}
                      </p>
                      {getStatusBadge(transaction.status)}
                    </div>
                    {transaction.description && (
                      <p className="text-sm text-gray-600 truncate">{transaction.description}</p>
                    )}
                    <p className="text-xs text-gray-400">{formatDate(transaction.created_at)}</p>
                  </div>
                </div>
                <div className={`font-bold text-lg ${getTransactionColor(transaction.transaction_type)}`}>
                  {transaction.amount > 0 ? '+' : ''}{formatCurrency(transaction.amount)}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default BalancePage
