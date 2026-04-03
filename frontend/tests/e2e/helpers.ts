import { Page, expect } from '@playwright/test';

export async function navigateToRegister(page: Page) {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  await page.getByRole('link', { name: '注册' }).click();
  await page.waitForLoadState('networkidle');
  await expect(page.getByText(/用户注册/i)).toBeVisible();
}

export async function fillRegistrationForm(page: Page, user: { username: string; email: string; password: string }) {
  await page.getByLabel(/用户名/i).fill(user.username);
  await page.getByLabel(/邮箱/i).fill(user.email);
  await page.getByLabel(/密码/i).fill(user.password);
}

export async function fillLoginForm(page: Page, user: { email: string; password: string }) {
  await page.getByLabel(/邮箱/i).fill(user.email);
  await page.getByLabel(/密码/i).fill(user.password);
}

export async function mockApiResponse(
  page: Page,
  urlPattern: string,
  response: {
    status: number;
    body: any;
    contentType?: string;
  }
) {
  await page.route(urlPattern, async route => {
    await route.fulfill({
      status: response.status,
      contentType: response.contentType || 'application/json',
      body: JSON.stringify(response.body)
    });
  });
}

export const testUser = {
  username: 'testuser',
  email: 'test@example.com',
  password: 'password123'
};