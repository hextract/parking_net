import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProtectedRoute from '../ProtectedRoute'

const mockUseAuth = vi.fn()
const mockNavigate = vi.fn()

vi.mock('../../context/AuthContext', () => ({
  useAuth: () => mockUseAuth(),
}))

vi.mock('react-router-dom', () => ({
  Navigate: ({ to }) => {
    mockNavigate(to)
    return null
  },
}))

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('redirects to login when not authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      user: null,
    })

    render(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>
    )

    expect(mockNavigate).toHaveBeenCalledWith('/login')
  })

  it('renders children when authenticated without role requirement', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      user: { login: 'testuser', role: 'driver' },
    })

    render(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>
    )

    expect(screen.getByText('Protected Content')).toBeInTheDocument()
  })

  it('renders children when authenticated with correct role', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      user: { login: 'testuser', role: 'driver' },
    })

    render(
      <ProtectedRoute requiredRole="driver">
        <div>Driver Content</div>
      </ProtectedRoute>
    )

    expect(screen.getByText('Driver Content')).toBeInTheDocument()
  })

  it('redirects when authenticated with wrong role', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      user: { login: 'testuser', role: 'owner' },
    })

    render(
      <ProtectedRoute requiredRole="driver">
        <div>Driver Content</div>
      </ProtectedRoute>
    )

    expect(mockNavigate).toHaveBeenCalled()
  })
})
