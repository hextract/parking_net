import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import LoadingSpinner from '../../components/LoadingSpinner'

const CallbackPage = () => {
  const [searchParams] = useSearchParams()
  const { setAuthData } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    const handleOAuthCallback = async () => {
      const token = searchParams.get('token')
      if (token) {
        console.log('OAuth callback: received token')
        const result = await setAuthData(token)
        console.log('OAuth callback: setAuthData result', result)
        
        if (result.success && result.user) {
          const redirectPath = result.user.role === 'driver' ? '/driver' : '/owner'
          console.log('OAuth callback: redirecting to', redirectPath)
          navigate(redirectPath, { replace: true })
        } else {
          console.error('OAuth callback: authentication failed', result.error)
          navigate('/login?error=oauth_failed', { replace: true })
        }
      } else {
        console.error('OAuth callback: no token in URL')
        navigate('/login?error=no_token', { replace: true })
      }
    }
    
    handleOAuthCallback()
  }, [searchParams, setAuthData, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100">
      <div className="text-center">
        <LoadingSpinner size="large" />
        <p className="mt-4 text-gray-600">Completing authentication...</p>
      </div>
    </div>
  )
}

export default CallbackPage

