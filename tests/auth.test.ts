import { test, expect } from '@playwright/test';

// By the time these tests run, setup.ts has already created a test account
// and saved the session to tests/.auth/user.json. The chromium project in
// playwright.config.ts loads that file automatically, so every test here
// starts fully logged in.

test.describe('Authenticated user', () => {

  test('lands on the dashboard after sign in', async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page.getByRole('heading', { name: 'Projects' })).toBeVisible();
  });

});
