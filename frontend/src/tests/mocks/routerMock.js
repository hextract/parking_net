import { vi } from 'vitest'

export const mockNavigate = vi.fn()
export const mockLocation = { pathname: '/' }

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useLocation: () => mockLocation,
    useParams: () => ({}),
    BrowserRouter: ({ children }) => children,
    Link: ({ children, to, ...props }) => {
      const element = document.createElement('a')
      element.href = to
      return element
    },
    Navigate: () => null,
  }
})
