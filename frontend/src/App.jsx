import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './context/AuthContext'
import Layout from './components/Layout'
import LoginPage from './pages/auth/LoginPage'
import RegisterPage from './pages/auth/RegisterPage'
import DriverDashboard from './pages/driver/DriverDashboard'
import SearchParking from './pages/driver/SearchParking'
import MyBookings from './pages/driver/MyBookings'
import BalancePage from './pages/driver/BalancePage'
import PromocodesPage from './pages/driver/PromocodesPage'
import OwnerDashboard from './pages/owner/OwnerDashboard'
import AllParkings from './pages/owner/AllParkings'
import MyParkings from './pages/owner/MyParkings'
import ParkingBookings from './pages/owner/ParkingBookings'
import AdminPage from './pages/admin/AdminPage'
import AdminPromocodesPage from './pages/admin/AdminPromocodesPage'
import ProtectedRoute from './components/ProtectedRoute'
import LoadingSpinner from './components/LoadingSpinner'

function App() {
  const { loading } = useAuth()

  if (loading) {
    return <LoadingSpinner fullScreen />
  }

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      <Route element={<Layout />}>
        {/* Driver Routes */}
        <Route
          path="/driver"
          element={
            <ProtectedRoute requiredRole="driver">
              <DriverDashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/driver/search"
          element={
            <ProtectedRoute requiredRole="driver">
              <SearchParking />
            </ProtectedRoute>
          }
        />
        <Route
          path="/driver/bookings"
          element={
            <ProtectedRoute requiredRole="driver">
              <MyBookings />
            </ProtectedRoute>
          }
        />
        <Route
          path="/driver/balance"
          element={
            <ProtectedRoute requiredRole="driver">
              <BalancePage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/driver/promocodes"
          element={
            <ProtectedRoute requiredRole="driver">
              <PromocodesPage />
            </ProtectedRoute>
          }
        />

        {/* Owner Routes */}
        <Route
          path="/owner"
          element={
            <ProtectedRoute requiredRole="owner">
              <OwnerDashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/owner/all-parkings"
          element={
            <ProtectedRoute requiredRole="owner">
              <AllParkings />
            </ProtectedRoute>
          }
        />
        <Route
          path="/owner/parkings"
          element={
            <ProtectedRoute requiredRole="owner">
              <MyParkings />
            </ProtectedRoute>
          }
        />
        <Route
          path="/owner/bookings/:parkingId"
          element={
            <ProtectedRoute requiredRole="owner">
              <ParkingBookings />
            </ProtectedRoute>
          }
        />
        <Route
          path="/owner/balance"
          element={
            <ProtectedRoute requiredRole="owner">
              <BalancePage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/owner/promocodes"
          element={
            <ProtectedRoute requiredRole="owner">
              <PromocodesPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/admin"
          element={
            <ProtectedRoute requiredRole="admin">
              <AdminPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/promocodes"
          element={
            <ProtectedRoute requiredRole="admin">
              <AdminPromocodesPage />
            </ProtectedRoute>
          }
        />
      </Route>

      <Route path="/" element={<Navigate to="/login" replace />} />
      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  )
}

export default App
