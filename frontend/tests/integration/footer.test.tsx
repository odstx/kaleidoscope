import { describe, it, expect, vi, beforeAll, afterAll, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { Footer } from '@/components/Footer'

const server = setupServer()

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

describe('Footer Component', () => {
  it('should display loading state initially', () => {
    render(<Footer />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('should fetch and display system info', async () => {
    const mockSystemInfo = {
      version: '1.0.0',
      build_id: 'build-123',
      build_time: '2024-01-01T00:00:00Z',
      git_commit: 'abc123'
    }

    server.use(
      http.get('http://localhost:9000/api/v1/system/info', () => {
        return HttpResponse.json(mockSystemInfo)
      })
    )

    render(<Footer />)

    await waitFor(() => {
      expect(screen.getByText(/Version:/)).toBeInTheDocument()
      expect(screen.getByText('1.0.0')).toBeInTheDocument()
      expect(screen.getByText(/Build ID:/)).toBeInTheDocument()
      expect(screen.getByText('build-123')).toBeInTheDocument()
    })
  })

  it('should display error message when API fails', async () => {
    server.use(
      http.get('http://localhost:9000/api/v1/system/info', () => {
        return new HttpResponse(null, { status: 500 })
      })
    )

    render(<Footer />)

    await waitFor(() => {
      expect(screen.getByText('Unable to load version info')).toBeInTheDocument()
    })
  })
})
