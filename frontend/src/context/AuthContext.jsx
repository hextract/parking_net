import { createContext, useContext, useState, useEffect } from 'react'
import { authService } from '../services/authService'
import { isTokenExpired } from '../utils/jwt'
import { getAuthToken } from '../utils/cookies'

const AuthContext = createContext(null)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null)
  const [token, setToken] = useState(null)
  const [loading, setLoading] = useState(true)

  // Fetch user info from backend on mount
  useEffect(() => {
    const initializeAuth = async () => {
      const storedToken = authService.getStoredToken()

      if (storedToken && !isTokenExpired(storedToken)) {
        setToken(storedToken)

        // Fetch user info from backend
        try {
          const userInfo = await authService.getUserInfo()
          if (userInfo) {
            setUser(userInfo)
          } else {
            // Token invalid, clear it
            authService.logout()
            setToken(null)
          }
        } catch (error) {
          // Token invalid or expired
          authService.logout()
          setToken(null)
        }
      } else if (storedToken) {
        // Token expired
        authService.logout()
      }

      setLoading(false)
    }

    initializeAuth()
  }, [])

  const login = async (credentials) => {
    try {
      const response = await authService.login(credentials)
      const { token: newToken } = response

      if (!newToken) {
        throw new Error('No token received')
      }

      // Store token in cookie
      authService.setAuthData(newToken)
      setToken(newToken)

      // Fetch user info from backend
      const userInfo = await authService.getUserInfo()
      if (userInfo) {
        setUser(userInfo)
        return { success: true, user: userInfo }
      } else {
        throw new Error('Failed to fetch user info')
      }
    } catch (error) {
      return { success: false, error: error.message || 'Login failed' }
    }
  }

  const register = async (userData) => {
    try {
      const response = await authService.register(userData)

      if (!response || !response.token) {
        throw new Error(response?.error_message || response?.message || 'Registration failed: No token received')
      }

      const { token: newToken } = response

      // Store token in cookie
      authService.setAuthData(newToken)
      setToken(newToken)

      // Fetch user info from backend
      const userInfo = await authService.getUserInfo()
      if (userInfo) {
        setUser(userInfo)
        return { success: true, user: userInfo }
      } else {
        throw new Error('Failed to fetch user info')
      }
    } catch (error) {
      const errorMessage = error.response?.data?.error_message ||
                          error.response?.data?.message ||
                          error.message ||
                          'Registration failed'
      return { success: false, error: errorMessage }
    }
  }

  const logout = () => {
    authService.logout()
    setToken(null)
    setUser(null)
  }

  const changePassword = async (passwordData) => {
    try {
      const response = await authService.changePassword(passwordData)
      const { token: newToken } = response

      if (newToken) {
        authService.setAuthData(newToken)
        setToken(newToken)
      }

      return { success: true }
    } catch (error) {
      return { success: false, error: error.message || 'Password change failed' }
    }
  }

  const refreshUserInfo = async () => {
    try {
      const userInfo = await authService.getUserInfo()
      if (userInfo) {
        setUser(userInfo)
        return true
      }
      return false
    } catch (error) {
      return false
    }
  }

  const value = {
    user,
    token,
    loading,
    isAuthenticated: !!token && !!user,
    login,
    register,
    logout,
    changePassword,
    refreshUserInfo,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
