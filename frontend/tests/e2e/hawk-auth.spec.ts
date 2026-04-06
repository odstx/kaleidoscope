import { test, expect } from '@playwright/test';
import { 
  navigateToRegister, 
  fillRegistrationForm, 
  fillLoginForm,
  completeLoginForm,
  mockApiResponse,
  testUser 
} from './helpers';

test.describe('Hawk Authentication E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the application before each test
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test.describe('Hawk Setup Flow', () => {
    test('should display Hawk setup button when not configured', async ({ page }) => {
      // Mock successful login first
      await mockApiResponse(page, '**/api/v1/users/login', {
        status: 200,
        body: {
          success: true,
          message: '登录成功',
          token: 'fake-jwt-token',
          user: { id: 1, username: testUser.username, email: testUser.email, hawk_enabled: false }
        }
      });
      
      // First login to access user profile
      await completeLoginForm(page, testUser);
      
      // Wait for success dialog and close it
      const dialogTitle = page.locator('[role="dialog"] [data-radix-collection-item]').first();
      await expect(dialogTitle).toBeVisible();
      await expect(dialogTitle).toContainText('登录成功');
      
      // Close dialog by clicking outside or OK button
      await page.mouse.click(0, 0);
      
      // Navigate to user profile page
      await page.getByRole('link', { name: /用户信息/i }).click();
      await page.waitForLoadState('networkidle');
      
      // Check that Hawk setup button is visible when not configured
      await expect(page.getByRole('button', { name: /设置 Hawk/i })).toBeVisible();
      await expect(page.getByRole('button', { name: /启用 Hawk/i })).toBeHidden();
      await expect(page.getByRole('button', { name: /禁用 Hawk/i })).toBeHidden();
    });

    test('should setup Hawk key successfully', async ({ page }) => {
      // Mock successful login first
      await mockApiResponse(page, '**/api/v1/users/login', {
        status: 200,
        body: {
          success: true,
          message: '登录成功',
          token: 'fake-jwt-token',
          user: { id: 1, username: testUser.username, email: testUser.email, hawk_enabled: false }
        }
      });
      
      // Login first
      await completeLoginForm(page, testUser);
      
      // Navigate to user profile
      await page.getByRole('link', { name: /用户信息/i }).click();
      await page.waitForLoadState('networkidle');
      
      // Mock Hawk setup API response
      await mockApiResponse(page, '**/api/v1/users/hawk/setup', {
        status: 200,
        body: {
          success: true,
          message: 'Hawk密钥已生成',
          hawk_key: 'werxhqm98rp3ngv998sjpsj9s98qjxh'
        }
      });
      
      // Click setup button
      await page.getByRole('button', { name: /设置 Hawk/i }).click();
      
      // Check that Hawk key is displayed and enable button appears
      await expect(page.getByText(/Hawk密钥:/)).toBeVisible();
      await expect(page.getByText(/werxhqm98rp3ngv998sjpsj9s98qjxh/)).toBeVisible();
      await expect(page.getByRole('button', { name: /启用 Hawk/i })).toBeVisible();
    });
  });

  test.describe('Hawk Enable/Disable Flow', () => {
    test('should enable Hawk authentication successfully', async ({ page }) => {
      // Mock successful login first
      await mockApiResponse(page, '**/api/v1/users/login', {
        status: 200,
        body: {
          success: true,
          message: '登录成功',
          token: 'fake-jwt-token',
          user: { id: 1, username: testUser.username, email: testUser.email, hawk_enabled: false }
        }
      });
      
      // Login first
      await completeLoginForm(page, testUser);
      
      // Navigate to user profile
      await page.getByRole('link', { name: /用户信息/i }).click();
      await page.waitForLoadState('networkidle');
      
      // Setup Hawk first
      await mockApiResponse(page, '**/api/v1/users/hawk/setup', {
        status: 200,
        body: {
          success: true,
          message: 'Hawk密钥已生成',
          hawk_key: 'werxhqm98rp3ngv998sjpsj9s98qjxh'
        }
      });
      
      await page.getByRole('button', { name: /设置 Hawk/i }).click();
      
      // Mock enable Hawk API response
      await mockApiResponse(page, '**/api/v1/users/hawk/enable', {
        status: 200,
        body: {
          success: true,
          message: 'Hawk认证已启用'
        }
      });
      
      // Click enable button
      await page.getByRole('button', { name: /启用 Hawk/i }).click();
      
      // Check success message and verify UI changes
      const enableDialog = page.locator('[role="dialog"]');
      await expect(enableDialog).toBeVisible();
      await expect(enableDialog).toContainText('Hawk认证已启用');
      
      // After enabling, disable button should be visible
      await page.mouse.click(0, 0); // Close dialog
      await expect(page.getByRole('button', { name: /禁用 Hawk/i })).toBeVisible();
      await expect(page.getByRole('button', { name: /启用 Hawk/i })).toBeHidden();
    });

    test('should disable Hawk authentication successfully', async ({ page }) => {
      // Mock successful login with Hawk already enabled
      await mockApiResponse(page, '**/api/v1/users/login', {
        status: 200,
        body: {
          success: true,
          message: '登录成功',
          token: 'fake-jwt-token',
          user: { id: 1, username: testUser.username, email: testUser.email, hawk_enabled: true }
        }
      });
      
      // Login with Hawk already enabled
      await completeLoginForm(page, testUser);
      
      // Navigate to user profile
      await page.getByRole('link', { name: /用户信息/i }).click();
      await page.waitForLoadState('networkidle');
      
      // Verify disable button is visible when Hawk is enabled
      await expect(page.getByRole('button', { name: /禁用 Hawk/i })).toBeVisible();
      
      // Mock disable Hawk API response
      await mockApiResponse(page, '**/api/v1/users/hawk/disable', {
        status: 200,
        body: {
          success: true,
          message: 'Hawk认证已禁用'
        }
      });
      
      // Click disable button
      await page.getByRole('button', { name: /禁用 Hawk/i }).click();
      
      // Check success message
      const disableDialog = page.locator('[role="dialog"]');
      await expect(disableDialog).toBeVisible();
      await expect(disableDialog).toContainText('Hawk认证已禁用');
      
      // After disabling, setup button should be visible again
      await page.mouse.click(0, 0); // Close dialog
      await expect(page.getByRole('button', { name: /设置 Hawk/i })).toBeVisible();
      await expect(page.getByRole('button', { name: /禁用 Hawk/i })).toBeHidden();
    });
  });

  test.describe('Hawk Error Handling', () => {
    test('should handle Hawk setup API errors', async ({ page }) => {
      // Mock successful login first
      await mockApiResponse(page, '**/api/v1/users/login', {
        status: 200,
        body: {
          success: true,
          message: '登录成功',
          token: 'fake-jwt-token',
          user: { id: 1, username: testUser.username, email: testUser.email, hawk_enabled: false }
        }
      });
      
      // Login first
      await completeLoginForm(page, testUser);
      
      // Navigate to user profile
      await page.getByRole('link', { name: /用户信息/i }).click();
      await page.waitForLoadState('networkidle');
      
      // Mock failed Hawk setup API response
      await mockApiResponse(page, '**/api/v1/users/hawk/setup', {
        status: 500,
        body: {
          success: false,
          message: '服务器错误'
        }
      });
      
      // Click setup button
      await page.getByRole('button', { name: /设置 Hawk/i }).click();
      
      // Check error message
      await expect(page.getByText(/服务器错误/i)).toBeVisible();
    });

    test('should handle Hawk enable API errors', async ({ page }) => {
      // Mock successful login first
      await mockApiResponse(page, '**/api/v1/users/login', {
        status: 200,
        body: {
          success: true,
          message: '登录成功',
          token: 'fake-jwt-token',
          user: { id: 1, username: testUser.username, email: testUser.email, hawk_enabled: false }
        }
      });
      
      // Login first
      await completeLoginForm(page, testUser);
      
      // Navigate to user profile
      await page.getByRole('link', { name: /用户信息/i }).click();
      await page.waitForLoadState('networkidle');
      
      // Setup Hawk first
      await mockApiResponse(page, '**/api/v1/users/hawk/setup', {
        status: 200,
        body: {
          success: true,
          message: 'Hawk密钥已生成',
          hawk_key: 'werxhqm98rp3ngv998sjpsj9s98qjxh'
        }
      });
      
      await page.getByRole('button', { name: /设置 Hawk/i }).click();
      
      // Mock failed enable Hawk API response
      await mockApiResponse(page, '**/api/v1/users/hawk/enable', {
        status: 400,
        body: {
          success: false,
          message: '无效的请求'
        }
      });
      
      // Click enable button
      await page.getByRole('button', { name: /启用 Hawk/i }).click();
      
      // Check error message
      await expect(page.getByText(/无效的请求/i)).toBeVisible();
    });
  });
});