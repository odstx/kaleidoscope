import { describe, it, expect, vi, afterEach, beforeAll, afterAll } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import userEvent from '@testing-library/user-event'
import App from '@/App'

// 创建模拟服务器
const server = setupServer()

// 定义测试用户数据
const testUser = {
  username: 'testuser',
  email: 'test@example.com',
  password: 'password123'
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

describe('Authentication Integration Tests', () => {
  describe('Registration Flow', () => {
    it('should display registration form with all required fields', () => {
      renderWithRouter(<App />)
      
      // 导航到注册页面
      const registerLink = screen.getByText('注册')
      fireEvent.click(registerLink)
      
      // 检查表单字段
      expect(screen.getByLabelText(/用户名/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/邮箱/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/密码/i)).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /注册/i })).toBeInTheDocument()
    })

    it('should validate registration form inputs', async () => {
      renderWithRouter(<App />)
      
      // 导航到注册页面
      const registerLink = screen.getByText('注册')
      fireEvent.click(registerLink)
      
      const user = userEvent.setup()
      
      // 测试太短的用户名
      const usernameInput = screen.getByLabelText(/用户名/i)
      await user.type(usernameInput, 'ab')
      await user.tab()
      
      // 等待验证错误出现
      await waitFor(() => {
        expect(screen.getByText(/用户名至少需要3个字符/i)).toBeInTheDocument()
      })
      
      // 测试无效邮箱
      const emailInput = screen.getByLabelText(/邮箱/i)
      await user.clear(emailInput)
      await user.type(emailInput, 'invalid-email')
      await user.tab()
      
      await waitFor(() => {
        expect(screen.getByText(/请输入有效的邮箱地址/i)).toBeInTheDocument()
      })
      
      // 测试太短的密码
      const passwordInput = screen.getByLabelText(/密码/i)
      await user.clear(passwordInput)
      await user.type(passwordInput, '123')
      await user.tab()
      
      await waitFor(() => {
        expect(screen.getByText(/密码至少需要6个字符/i)).toBeInTheDocument()
      })
    })

    it('should submit registration form with valid data', async () => {
      // 模拟成功的注册API响应
      server.use(
        http.post('http://localhost:3000/api/v1/users/register', async () => {
          return HttpResponse.json(
            { 
              success: true, 
              message: '注册成功',
              user: { id: 1, username: testUser.username, email: testUser.email }
            },
            { status: 201 }
          )
        })
      )

      renderWithRouter(<App />)
      
      // 导航到注册页面
      const registerLink = screen.getByText('注册')
      fireEvent.click(registerLink)
      
      const user = userEvent.setup()
      
      // 填写表单
      await user.type(screen.getByLabelText(/用户名/i), testUser.username)
      await user.type(screen.getByLabelText(/邮箱/i), testUser.email)
      await user.type(screen.getByLabelText(/密码/i), testUser.password)
      
      // 提交表单
      const submitButton = screen.getByRole('button', { name: /注册/i })
      await user.click(submitButton)
      
      // 等待API调用成功
      await waitFor(() => {
        expect(screen.getByText(/注册成功/i)).toBeInTheDocument()
      })
    })

    it('should handle registration API errors', async () => {
      // 模拟失败的注册API响应
      server.use(
        http.post('http://localhost:3000/api/v1/users/register', async () => {
          return HttpResponse.json(
            { 
              success: false, 
              message: '邮箱已被注册'
            },
            { status: 400 }
          )
        })
      )

      renderWithRouter(<App />)
      
      // 导航到注册页面
      const registerLink = screen.getByText('注册')
      fireEvent.click(registerLink)
      
      const user = userEvent.setup()
      
      // 填写表单
      await user.type(screen.getByLabelText(/用户名/i), testUser.username)
      await user.type(screen.getByLabelText(/邮箱/i), testUser.email)
      await user.type(screen.getByLabelText(/密码/i), testUser.password)
      
      // 提交表单
      const submitButton = screen.getByRole('button', { name: /注册/i })
      await user.click(submitButton)
      
      // 等待错误消息显示
      await waitFor(() => {
        expect(screen.getByText(/邮箱已被注册/i)).toBeInTheDocument()
      })
    })
  })

  describe('Login Flow', () => {
    it('should display login form with all required fields', () => {
      renderWithRouter(<App />)
      
      // 默认应该显示登录页面
      expect(screen.getByLabelText(/邮箱/i)).toBeInTheDocument()
      expect(screen.getByLabelText(/密码/i)).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /登录/i })).toBeInTheDocument()
    })

    it('should validate login form inputs', async () => {
      renderWithRouter(<App />)
      
      const user = userEvent.setup()
      
      // 测试无效邮箱
      const emailInput = screen.getByLabelText(/邮箱/i)
      await user.type(emailInput, 'invalid-email')
      await user.tab()
      
      await waitFor(() => {
        expect(screen.getByText(/请输入有效的邮箱地址/i)).toBeInTheDocument()
      })
      
      // 测试空密码
      const passwordInput = screen.getByLabelText(/密码/i)
      await user.clear(passwordInput)
      await user.tab()
      
      await waitFor(() => {
        expect(screen.getByText(/请输入密码/i)).toBeInTheDocument()
      })
    })

    it('should submit login form with valid credentials', async () => {
      // 模拟成功的登录API响应
      server.use(
        http.post('http://localhost:3000/api/v1/users/login', async () => {
          return HttpResponse.json(
            { 
              success: true, 
              message: '登录成功',
              token: 'fake-jwt-token',
              user: { id: 1, username: testUser.username, email: testUser.email }
            },
            { status: 200 }
          )
        })
      )

      renderWithRouter(<App />)
      
      const user = userEvent.setup()
      
      // 填写表单
      await user.type(screen.getByLabelText(/邮箱/i), testUser.email)
      await user.type(screen.getByLabelText(/密码/i), testUser.password)
      
      // 提交表单
      const submitButton = screen.getByRole('button', { name: /登录/i })
      await user.click(submitButton)
      
      // 等待API调用成功
      await waitFor(() => {
        expect(screen.getByText(/登录成功/i)).toBeInTheDocument()
      })
    })

    it('should handle login API errors', async () => {
      // 模拟失败的登录API响应
      server.use(
        http.post('http://localhost:3000/api/v1/users/login', async () => {
          return HttpResponse.json(
            { 
              success: false, 
              message: '邮箱或密码错误'
            },
            { status: 401 }
          )
        })
      )

      renderWithRouter(<App />)
      
      const user = userEvent.setup()
      
      // 填写表单
      await user.type(screen.getByLabelText(/邮箱/i), testUser.email)
      await user.type(screen.getByLabelText(/密码/i), testUser.password)
      
      // 提交表单
      const submitButton = screen.getByRole('button', { name: /登录/i })
      await user.click(submitButton)
      
      // 等待错误消息显示
      await waitFor(() => {
        expect(screen.getByText(/邮箱或密码错误/i)).toBeInTheDocument()
      })
    })
  })

  describe('Navigation between auth pages', () => {
    it('should navigate from login to register page', async () => {
      renderWithRouter(<App />)
      
      const user = userEvent.setup()
      
      // 点击导航到注册页面
      const registerLink = screen.getByText('注册')
      await user.click(registerLink)
      
      // 确认已切换到注册页面
      await waitFor(() => {
        expect(screen.getByText(/创建账户/i)).toBeInTheDocument()
      })
    })

    it('should navigate from register to login page', async () => {
      renderWithRouter(<App />)
      
      const user = userEvent.setup()
      
      // 先导航到注册页面
      const registerLink = screen.getByText('注册')
      await user.click(registerLink)
      
      // 等待页面切换
      await waitFor(() => {
        expect(screen.getByText(/创建账户/i)).toBeInTheDocument()
      })
      
      // 点击返回登录链接
      const loginLink = screen.getByText('登录')
      await user.click(loginLink)
      
      // 确认已切换到登录页面
      await waitFor(() => {
        expect(screen.getByText(/登录到您的账户/i)).toBeInTheDocument()
      })
    })
  })
})