import { describe, it, expect, vi, afterEach, beforeAll, afterAll, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import userEvent from '@testing-library/user-event'
import App from '@/App'

const server = setupServer()

const testUser = {
  id: 1,
  email: 'test@example.com',
  username: 'testuser'
}

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

vi.stubGlobal('import', {
  meta: {
    env: {
      VITE_API_BASE_URL: 'http://localhost:9000'
    }
  }
})

const mockSystemInfo = {
  version: '1.0.0',
  build_id: 'test-build',
  build_time: '2024-01-01T00:00:00Z',
  git_commit: 'abc123'
}

const renderWithRouter = (ui: React.ReactElement) => {
  server.use(
    http.get('http://localhost:9000/api/v1/system/info', () => {
      return HttpResponse.json(mockSystemInfo)
    })
  )
  return render(ui)
}

describe('Dashboard Integration Tests', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  describe('Authentication Guard', () => {
    it('should redirect to login when no token is present', async () => {
      window.history.pushState({}, '', '/dashboard')
      renderWithRouter(<App />)

      await waitFor(() => {
        expect(screen.getByLabelText(/邮箱/i)).toBeInTheDocument()
      })
    })

    it('should show loading state initially when token exists', async () => {
      localStorage.setItem('token', 'fake-jwt-token')
      
      server.use(
        http.get('http://localhost:9000/api/v1/users/info', () => {
          return new Promise(() => {})
        })
      )

      window.history.pushState({}, '', '/dashboard')
      renderWithRouter(<App />)

      expect(screen.getByText(/加载中/i)).toBeInTheDocument()
    })
  })

  describe('Dashboard Content', () => {
    it('should display user information after successful fetch', async () => {
      localStorage.setItem('token', 'fake-jwt-token')

      server.use(
        http.get('http://localhost:9000/api/v1/users/info', () => {
          return HttpResponse.json(testUser)
        })
      )

      window.history.pushState({}, '', '/dashboard')
      renderWithRouter(<App />)

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: /控制面板/i })).toBeInTheDocument()
      })

      expect(screen.getByText(/用户信息/i)).toBeInTheDocument()
      expect(screen.getByText(testUser.email)).toBeInTheDocument()
    })

    it('should display welcome card', async () => {
      localStorage.setItem('token', 'fake-jwt-token')

      server.use(
        http.get('http://localhost:9000/api/v1/users/info', () => {
          return HttpResponse.json(testUser)
        })
      )

      window.history.pushState({}, '', '/dashboard')
      renderWithRouter(<App />)

      await waitFor(() => {
        expect(screen.getByText(/欢迎使用/i)).toBeInTheDocument()
      })

      expect(screen.getByText(/这是您的个人控制面板/i)).toBeInTheDocument()
    })
  })

  describe('Logout Functionality', () => {
    it('should logout and redirect to login page', async () => {
      localStorage.setItem('token', 'fake-jwt-token')

      server.use(
        http.get('http://localhost:9000/api/v1/users/info', () => {
          return HttpResponse.json(testUser)
        })
      )

      window.history.pushState({}, '', '/dashboard')
      renderWithRouter(<App />)

      const user = userEvent.setup()

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: /控制面板/i })).toBeInTheDocument()
      })

      const logoutButton = screen.getByRole('button', { name: /退出登录/i })
      await user.click(logoutButton)

      await waitFor(() => {
        expect(localStorage.getItem('token')).toBeNull()
        expect(screen.getByLabelText(/邮箱/i)).toBeInTheDocument()
      })
    })
  })

  describe('Error Handling', () => {
    it('should redirect to login when API returns 401', async () => {
      localStorage.setItem('token', 'invalid-token')

      server.use(
        http.get('http://localhost:9000/api/v1/users/info', () => {
          return HttpResponse.json(
            { message: 'Unauthorized' },
            { status: 401 }
          )
        })
      )

      window.history.pushState({}, '', '/dashboard')
      renderWithRouter(<App />)

      await waitFor(() => {
        expect(localStorage.getItem('token')).toBeNull()
        expect(screen.getByLabelText(/邮箱/i)).toBeInTheDocument()
      })
    })

    it('should redirect to login when API request fails', async () => {
      localStorage.setItem('token', 'fake-jwt-token')

      server.use(
        http.get('http://localhost:9000/api/v1/users/info', () => {
          return HttpResponse.error()
        })
      )

      window.history.pushState({}, '', '/dashboard')
      renderWithRouter(<App />)

      await waitFor(() => {
        expect(localStorage.getItem('token')).toBeNull()
        expect(screen.getByLabelText(/邮箱/i)).toBeInTheDocument()
      })
    })
  })
})
