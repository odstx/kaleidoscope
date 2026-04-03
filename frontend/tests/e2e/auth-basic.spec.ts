import { test, expect } from '@playwright/test';

test.describe('Authentication - Basic UI Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display login page by default', async ({ page }) => {
    await expect(page.getByText('用户登录')).toBeVisible();
    await expect(page.getByRole('heading', { name: '用户登录' })).toBeVisible();
    await expect(page.getByRole('button', { name: '登录' })).toBeVisible();
    await expect(page.getByRole('link', { name: '去注册' })).toBeVisible();
  });

  test('should navigate to registration page', async ({ page }) => {
    await page.getByRole('link', { name: '去注册' }).click();
    await expect(page.getByText('用户注册')).toBeVisible();
    await expect(page.getByRole('heading', { name: '用户注册' })).toBeVisible();
    await expect(page.getByRole('button', { name: '注册' })).toBeVisible();
    await expect(page.getByRole('link', { name: '去登录' })).toBeVisible();
  });

  test('should navigate back to login page from registration', async ({ page }) => {
    await page.goto('/register');
    await expect(page.getByText('用户注册')).toBeVisible();
    
    await page.getByRole('link', { name: '去登录' }).click();
    await expect(page.getByText('用户登录')).toBeVisible();
  });

  test('registration form validation - username too short', async ({ page }) => {
    await page.goto('/register');
    
    const usernameInput = page.getByLabel('用户名');
    await usernameInput.fill('ab');
    await usernameInput.blur();
    
    // Submit form to trigger validation
    await page.getByRole('button', { name: '注册' }).click();
    await expect(page.getByText('用户名至少需要3个字符')).toBeVisible();
  });

  test('registration form validation - invalid email', async ({ page }) => {
    await page.goto('/register');
    
    const emailInput = page.getByLabel('邮箱');
    await emailInput.fill('invalid-email');
    await emailInput.blur();
    
    // Submit form to trigger validation
    await page.getByRole('button', { name: '注册' }).click();
    await expect(page.getByText('请输入有效的邮箱地址')).toBeVisible();
  });

  test('registration form validation - password too short', async ({ page }) => {
    await page.goto('/register');
    
    const passwordInput = page.getByLabel('密码');
    await passwordInput.fill('123');
    await passwordInput.blur();
    
    // Submit form to trigger validation
    await page.getByRole('button', { name: '注册' }).click();
    await expect(page.getByText('密码至少需要6个字符')).toBeVisible();
  });

  test('login form validation - invalid email', async ({ page }) => {
    const emailInput = page.getByLabel('邮箱');
    await emailInput.fill('invalid-email');
    await emailInput.blur();
    
    // Submit form to trigger validation
    await page.getByRole('button', { name: '登录' }).click();
    await expect(page.getByText('请输入有效的邮箱地址')).toBeVisible();
  });

  test('login form validation - empty password', async ({ page }) => {
    const passwordInput = page.getByLabel('密码');
    await passwordInput.fill('');
    await passwordInput.blur();
    
    // Submit form to trigger validation
    await page.getByRole('button', { name: '登录' }).click();
    await expect(page.getByText('请输入密码')).toBeVisible();
  });

  test('registration form - valid username format shows no errors', async ({ page }) => {
    await page.goto('/register');
    
    const usernameInput = page.getByLabel('用户名');
    await usernameInput.fill('valid_user123');
    await usernameInput.blur();
    
    // Submit form - should not show username errors
    await page.getByRole('button', { name: '注册' }).click();
    await expect(page.getByText('用户名至少需要3个字符')).not.toBeVisible();
    await expect(page.getByText('用户名最多20个字符')).not.toBeVisible();
    await expect(page.getByText('用户名只能包含字母、数字和下划线')).not.toBeVisible();
    // Should show email/password required errors instead
    await expect(page.getByText('请输入有效的邮箱地址')).toBeVisible();
  });

  test('registration form - invalid username characters', async ({ page }) => {
    await page.goto('/register');
    
    const usernameInput = page.getByLabel('用户名');
    await usernameInput.fill('invalid@user');
    await usernameInput.blur();
    
    // Submit form to trigger validation
    await page.getByRole('button', { name: '注册' }).click();
    await expect(page.getByText('用户名只能包含字母、数字和下划线')).toBeVisible();
  });

  test('registration form - username too long', async ({ page }) => {
    await page.goto('/register');
    
    const usernameInput = page.getByLabel('用户名');
    await usernameInput.fill('thisusernameistoolongforvalidation');
    await usernameInput.blur();
    
    // Submit form to trigger validation
    await page.getByRole('button', { name: '注册' }).click();
    await expect(page.getByText('用户名最多20个字符')).toBeVisible();
  });
});