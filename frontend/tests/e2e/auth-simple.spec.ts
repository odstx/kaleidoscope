import { test, expect } from '@playwright/test';
import { navigateToRegister, fillRegistrationForm, fillLoginForm, mockApiResponse, testUser } from './helpers';

test.describe('Authentication E2E Tests (Simple)', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test('registration with valid data shows success message', async ({ page }) => {
    await navigateToRegister(page);
    await fillRegistrationForm(page, testUser);
    
    await mockApiResponse(page, '**/api/v1/users/register', {
      status: 201,
      body: {
        success: true,
        message: '注册成功',
        user: { id: 1, username: testUser.username, email: testUser.email }
      }
    });

    await page.getByRole('button', { name: /注册/i }).click();
    await expect(page.getByText(/注册成功/i)).toBeVisible();
  });

  test('registration with existing email shows error', async ({ page }) => {
    await navigateToRegister(page);
    await fillRegistrationForm(page, testUser);
    
    await mockApiResponse(page, '**/api/v1/users/register', {
      status: 400,
      body: {
        success: false,
        message: '邮箱已被注册'
      }
    });

    await page.getByRole('button', { name: /注册/i }).click();
    await expect(page.getByText(/邮箱已被注册/i)).toBeVisible();
  });

  test('login with valid credentials shows success message', async ({ page }) => {
    await fillLoginForm(page, testUser);
    
    await mockApiResponse(page, '**/api/v1/users/login', {
      status: 200,
      body: {
        success: true,
        message: '登录成功',
        token: 'fake-jwt-token',
        user: { id: 1, username: testUser.username, email: testUser.email }
      }
    });

    await page.getByRole('button', { name: /登录/i }).click();
    await expect(page.getByText(/登录成功/i)).toBeVisible();
  });

  test('login with invalid credentials shows error', async ({ page }) => {
    await fillLoginForm(page, testUser);
    
    await mockApiResponse(page, '**/api/v1/users/login', {
      status: 401,
      body: {
        success: false,
        message: '邮箱或密码错误'
      }
    });

    await page.getByRole('button', { name: /登录/i }).click();
    await expect(page.getByText(/邮箱或密码错误/i)).toBeVisible();
  });

  test('form validation shows appropriate error messages', async ({ page }) => {
    await navigateToRegister(page);
    
    // Test username validation
    const usernameInput = page.getByLabel(/用户名/i);
    await usernameInput.fill('ab');
    await usernameInput.blur();
    await expect(page.getByText(/用户名至少需要3个字符/i)).toBeVisible();
    
    // Test email validation
    const emailInput = page.getByLabel(/邮箱/i);
    await emailInput.clear();
    await emailInput.fill('invalid-email');
    await emailInput.blur();
    await expect(page.getByText(/请输入有效的邮箱地址/i)).toBeVisible();
    
    // Test password validation
    const passwordInput = page.getByLabel(/密码/i);
    await passwordInput.clear();
    await passwordInput.fill('123');
    await passwordInput.blur();
    await expect(page.getByText(/密码至少需要6个字符/i)).toBeVisible();
  });

  test('navigation between login and register pages works', async ({ page }) => {
    // From login to register
    await page.getByRole('link', { name: '注册' }).click();
    await expect(page.getByText(/用户注册/i)).toBeVisible();
    
    // From register back to login
    await page.getByRole('link', { name: '登录' }).click();
    await expect(page.getByText(/用户登录/i)).toBeVisible();
  });
});