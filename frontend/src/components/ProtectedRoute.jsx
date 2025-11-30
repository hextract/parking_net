import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const ProtectedRoute = ({ children, requiredRole }) => {
  const { isAuthenticated, user } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (requiredRole && user?.role !== requiredRole) {
    // Redirect to appropriate dashboard based on actual role
    let redirectPath = '/login'
    if (user?.role === 'driver') {
      redirectPath = '/driver'
    } else if (user?.role === 'owner') {
      redirectPath = '/owner'
    } else if (user?.role === 'admin') {
      redirectPath = '/admin'
    }
    return <Navigate to={redirectPath} replace />
  }

  return children
}

export default ProtectedRoute
