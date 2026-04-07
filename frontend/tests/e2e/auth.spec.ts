import { test, expect } from '@playwright/test';

test('should register a new user successfully', async ({ page }) => {
  // Navigate to registration page
  await page.goto('/register');
  
  // Generate unique test data (keep username under 20 chars)
  const uniqueId = Date.now().toString().slice(-6); // Last 6 digits of timestamp
  const username = `test${uniqueId}`;
  const email = `test${uniqueId}@example.com`;
  
  // Fill registration form with valid data
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'password123');
  
  // Submit registration form
  await page.click('button[type="submit"]');
  
  // Wait for success dialog to appear
  await expect(page.getByText('Registration Successful')).toBeVisible();
  
  // Wait for automatic redirect to login page (5 second countdown)
  await page.waitForURL('/login');
  
  // Verify we're on the login page
  await expect(page).toHaveURL('/login');
  await expect(page.getByText('Enter your account information')).toBeVisible();
});

test('should login with registered user successfully', async ({ page }) => {
  // Generate unique test data (keep username under 20 chars)
  const uniqueId = Date.now().toString().slice(-6); // Last 6 digits of timestamp
  const username = `login${uniqueId}`;
  const email = `login${uniqueId}@example.com`;
  
  // First register a user
  await page.goto('/register');
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  await page.waitForURL('/login');
  
  // Now login with the registered user
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  
  // Should be redirected to dashboard after successful login
  await page.waitForURL('/dashboard');
  await expect(page).toHaveURL('/dashboard');
});

test('should show validation errors for invalid registration data', async ({ page }) => {
  await page.goto('/register');
  
  // Test username too short
  await page.fill('input[name="username"]', 'ab');
  await page.fill('input[name="email"]', 'test@example.com');
  await page.fill('input[name="password"]', 'pass');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('Username must be at least 3 characters')).toBeVisible();
  
  // Test invalid email
  await page.fill('input[name="username"]', 'validuser');
  await page.fill('input[name="email"]', 'invalid-email');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('Please enter a valid email address')).toBeVisible();
  
  // Test password too short
  await page.fill('input[name="email"]', 'test@example.com');
  await page.fill('input[name="password"]', '123');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('Password must be at least 6 characters')).toBeVisible();
});

test('should show error for duplicate registration', async ({ page }) => {
  // Generate unique test data (keep username under 20 chars)
  const uniqueId = Date.now().toString().slice(-6); // Last 6 digits of timestamp
  const username = `dupuser${uniqueId}`;
  const email = `duplicate${uniqueId}@example.com`;
  
  // Register first user
  await page.goto('/register');
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  await page.waitForURL('/login');
  
  // Try to register same user again
  await page.goto('/register');
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  
  // Should show registration error
  await expect(page.locator('.bg-destructive/10')).toBeVisible();
});

test('should show validation errors for invalid login data', async ({ page }) => {
  await page.goto('/login');
  
  // Test invalid email format
  await page.fill('input[name="email"]', 'invalid-email');
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('Please enter a valid email address')).toBeVisible();
  
  // Test empty password
  await page.fill('input[name="email"]', 'test@example.com');
  await page.fill('input[name="password"]', '');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('Please enter password')).toBeVisible();
});

test('should navigate between login and register pages', async ({ page }) => {
  // Start on login page
  await page.goto('/login');
  await expect(page).toHaveURL('/login');
  
  // Navigate to register page
  await page.click('text=Go to Register');
  await page.waitForURL('/register');
  await expect(page).toHaveURL('/register');
  
  // Navigate back to login page
  await page.click('text=Go to login');
  await page.waitForURL('/login');
  await expect(page).toHaveURL('/login');
});

test('should send forgot password email successfully', async ({ page }) => {
  await page.goto('/forgot-password');
  
  await page.fill('input[name="email"]', 'test@example.com');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('If an account exists with this email, you will receive a password reset link shortly.')).toBeVisible();
  
  await expect(page.getByRole('link', { name: 'Back to Login' })).toBeVisible();
});

test('should show validation error for invalid email in forgot password', async ({ page }) => {
  await page.goto('/forgot-password');
  
  await page.fill('input[name="email"]', 'invalid-email');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('Please enter a valid email address')).toBeVisible();
});

test('should show error when accessing reset password without token', async ({ page }) => {
  await page.goto('/reset-password');
  
  await expect(page.getByText('Invalid or expired reset token')).toBeVisible();
  
  await expect(page.getByRole('button', { name: 'Request New Link' })).toBeVisible();
});

test('should navigate to forgot password from reset password without token', async ({ page }) => {
  await page.goto('/reset-password');
  
  await page.click('button:has-text("Request New Link")');
  await page.waitForURL('/forgot-password');
  await expect(page).toHaveURL('/forgot-password');
});

test('should show validation error for password mismatch in reset password', async ({ page }) => {
  await page.goto('/reset-password?token=valid-test-token');
  
  await page.fill('input[name="password"]', 'password123');
  await page.fill('input[name="confirmPassword"]', 'password456');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('Passwords do not match')).toBeVisible();
});

test('should show validation error for short password in reset password', async ({ page }) => {
  await page.goto('/reset-password?token=valid-test-token');
  
  await page.fill('input[name="password"]', '123');
  await page.fill('input[name="confirmPassword"]', '123');
  await page.click('button[type="submit"]');
  
  await expect(page.getByText('Password must be at least 6 characters')).toBeVisible();
});

test('should show error for wrong password on login', async ({ page }) => {
  const uniqueId = Date.now().toString().slice(-6);
  const email = `wrongpass${uniqueId}@example.com`;
  
  await page.goto('/register');
  await page.fill('input[name="username"]', `wrongpass${uniqueId}`);
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  await page.waitForURL('/login');
  
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'wrongpassword');
  await page.click('button[type="submit"]');
  
  await expect(page.locator('.bg-destructive\\/10')).toBeVisible();
});

test('should show error for non-existent user on login', async ({ page }) => {
  await page.goto('/login');
  
  await page.fill('input[name="email"]', 'nonexistent@example.com');
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  
  await expect(page.locator('.bg-destructive\\/10')).toBeVisible();
});

test('should navigate from login to forgot password', async ({ page }) => {
  await page.goto('/login');
  
  await page.click('text=Forgot Password?');
  await page.waitForURL('/forgot-password');
  await expect(page).toHaveURL('/forgot-password');
});

test('should logout successfully', async ({ page }) => {
  const uniqueId = Date.now().toString().slice(-6);
  const username = `logout${uniqueId}`;
  const email = `logout${uniqueId}@example.com`;
  
  await page.goto('/register');
  await page.fill('input[name="username"]', username);
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  await page.waitForURL('/login');
  
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');
  await page.waitForURL('/dashboard');
  
  await page.click('button:has-text("Menu")');
  await page.click('text=Logout');
  
  await page.waitForURL('/login');
  await expect(page).toHaveURL('/login');
});