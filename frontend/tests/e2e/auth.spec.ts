import { test, expect } from '@playwright/test';

const testUser = {
  username: 'testuser',
  email: 'test@example.com',
  password: 'password123'
};

test.describe('Authentication E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the application before each test
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test.describe('Registration Flow', () => {
    test('should display registration form with all required fields', async ({ page }) => {
      // Navigate to registration page
      await page.getByRole('link', { name: '注册' }).click();
      
      // Check form fields are present
      await expect(page.getByLabel(/用户名/i)).toBeVisible();
      await expect(page.getByLabel(/邮箱/i)).toBeVisible();
      await expect(page.getByLabel(/密码/i)).toBeVisible();
      await expect(page.getByRole('button', { name: /注册/i })).toBeVisible();
    });

    test('should validate registration form inputs', async ({ page }) => {
      // Navigate to registration page
      await page.getByRole('link', { name: '注册' }).click();
      
      // Test too short username
      const usernameInput = page.getByLabel(/用户名/i);
      await usernameInput.fill('ab');
      
      // Submit form to trigger validation
      await page.getByRole('button', { name: /注册/i }).click();
      await expect(page.getByText(/用户名至少需要3个字符/i)).toBeVisible();
      
      // Clear previous input and test invalid email
      await usernameInput.clear();
      await usernameInput.fill('validuser');
      const emailInput = page.getByLabel(/邮箱/i);
      await emailInput.clear();
      await emailInput.fill('invalid-email');
      
      // Submit form to trigger validation
      await page.getByRole('button', { name: /注册/i }).click();
      await expect(page.getByText(/请输入有效的邮箱地址/i)).toBeVisible();
      
      // Clear previous inputs and test too short password
      await emailInput.clear();
      await emailInput.fill('test@example.com');
      const passwordInput = page.getByLabel(/密码/i);
      await passwordInput.clear();
      await passwordInput.fill('123');
      
      // Submit form to trigger validation
      await page.getByRole('button', { name: /注册/i }).click();
      await expect(page.getByText(/密码至少需要6个字符/i)).toBeVisible();
    });

    test('should submit registration form with valid data', async ({ page }) => {
      // Navigate to registration page
      await page.getByRole('link', { name: '注册' }).click();
      
      // Fill the form
      await page.getByLabel(/用户名/i).fill(testUser.username);
      await page.getByLabel(/邮箱/i).fill(testUser.email);
      await page.getByLabel(/密码/i).fill(testUser.password);
      
      // Mock API response for successful registration
      await page.route('**/api/v1/users/register', async route => {
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            message: '注册成功',
            user: { id: 1, username: testUser.username, email: testUser.email }
          })
        });
      });
      
      // Submit the form
      await page.getByRole('button', { name: /注册/i }).click();
      
      // Wait for success dialog to appear and check its content
      const dialogTitle = page.locator('[role="dialog"] [data-radix-collection-item]').first();
      await expect(dialogTitle).toBeVisible();
      await expect(dialogTitle).toContainText('注册成功');
      
      // Check dialog description
      const dialogDescription = page.locator('[role="dialog"] [data-radix-collection-item]').nth(1);
      await expect(dialogDescription).toBeVisible();
      await expect(dialogDescription).toContainText('您的账户已成功创建！');
    });

    test('should handle registration API errors', async ({ page }) => {
      // Navigate to registration page
      await page.getByRole('link', { name: '注册' }).click();
      
      // Fill the form
      await page.getByLabel(/用户名/i).fill(testUser.username);
      await page.getByLabel(/邮箱/i).fill(testUser.email);
      await page.getByLabel(/密码/i).fill(testUser.password);
      
      // Mock API response for failed registration
      await page.route('**/api/v1/users/register', async route => {
        await route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({
            success: false,
            message: '邮箱已被注册'
          })
        });
      });
      
      // Submit the form
      await page.getByRole('button', { name: /注册/i }).click();
      
      // Wait for error message
      await expect(page.getByText(/邮箱已被注册/i)).toBeVisible();
    });
  });

  test.describe('Login Flow', () => {
    test('should display login form with all required fields', async ({ page }) => {
      // Default page should show login form
      await expect(page.getByLabel(/邮箱/i)).toBeVisible();
      await expect(page.getByLabel(/密码/i)).toBeVisible();
      await expect(page.getByRole('button', { name: /登录/i })).toBeVisible();
    });

    test('should validate login form inputs', async ({ page }) => {
      // Test invalid email
      const emailInput = page.getByLabel(/邮箱/i);
      await emailInput.fill('invalid-email');
      await emailInput.blur();
      
      await expect(page.getByText(/请输入有效的邮箱地址/i)).toBeVisible();
      
      // Test empty password
      const passwordInput = page.getByLabel(/密码/i);
      await passwordInput.clear();
      await passwordInput.blur();
      
      await expect(page.getByText(/请输入密码/i)).toBeVisible();
    });

    test('should submit login form with valid credentials', async ({ page }) => {
      // Fill the form
      await page.getByLabel(/邮箱/i).fill(testUser.email);
      await page.getByLabel(/密码/i).fill(testUser.password);
      
      // Mock API response for successful login
      await page.route('**/api/v1/users/login', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            message: '登录成功',
            token: 'fake-jwt-token',
            user: { id: 1, username: testUser.username, email: testUser.email }
          })
        });
      });
      
      // Submit the form
      await page.getByRole('button', { name: /登录/i }).click();
      
      // Wait for success dialog to appear and check its content
      const dialogTitle = page.locator('[role="dialog"] [data-radix-collection-item]').first();
      await expect(dialogTitle).toBeVisible();
      await expect(dialogTitle).toContainText('登录成功');
      
      // Check dialog description
      const dialogDescription = page.locator('[role="dialog"] [data-radix-collection-item]').nth(1);
      await expect(dialogDescription).toBeVisible();
      await expect(dialogDescription).toContainText('欢迎回来！');
    });

    test('should handle login API errors', async ({ page }) => {
      // Fill the form
      await page.getByLabel(/邮箱/i).fill(testUser.email);
      await page.getByLabel(/密码/i).fill(testUser.password);
      
      // Mock API response for failed login
      await page.route('**/api/v1/users/login', async route => {
        await route.fulfill({
          status: 401,
          contentType: 'application/json',
          body: JSON.stringify({
            success: false,
            message: '邮箱或密码错误'
          })
        });
      });
      
      // Submit the form
      await page.getByRole('button', { name: /登录/i }).click();
      
      // Wait for error message
      await expect(page.getByText(/邮箱或密码错误/i)).toBeVisible();
    });
  });

  test.describe('Navigation between auth pages', () => {
    test('should navigate from login to register page', async ({ page }) => {
      // Click link to navigate to register page
      await page.getByRole('link', { name: '注册' }).click();
      
      // Confirm switched to registration page
      await expect(page.getByText(/创建账户/i)).toBeVisible();
    });

    test('should navigate from register to login page', async ({ page }) => {
      // Navigate to registration page first
      await page.getByRole('link', { name: '注册' }).click();
      await expect(page.getByText(/创建账户/i)).toBeVisible();
      
      // Click link to navigate back to login page
      await page.getByRole('link', { name: '登录' }).click();
      
      // Confirm switched to login page
      await expect(page.getByText(/登录到您的账户/i)).toBeVisible();
    });
  });
});