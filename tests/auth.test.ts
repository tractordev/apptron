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

//   test('can create a project and see the editor loading', async ({ page }) => {
//     await page.goto('/dashboard');

//     // wait for session and html-component custom element to be registered
//     await page.locator('#header-bar.signedin').waitFor({ state: 'visible' });
//     await page.waitForFunction(() => customElements.get('html-component') !== undefined);

//     // click create project button and wait for modal
//     await page.getByRole('button', { name: 'Create Project' }).click();
//     await expect(page.getByRole('heading', { name: 'New Project' })).toBeVisible({ timeout: 30000 });

//     // set a unique name per project
//     const projectName = `project${Date.now()}`;
//     await page.getByLabel('Project name').fill(projectName);

//     // wait for duplicate name check
//     await expect(page.locator('#project-submit')).toBeEnabled({ timeout: 15000 });
//     await page.locator('#project-submit').click();

//     // wait for editor redirect
//     await page.waitForURL('**/edit/**', { timeout: 15000 });

//     // pierce the outer iframe first, then locate the inner one
//     await page.frameLocator('iframe').locator('#vscode-frame').waitFor({ state: 'attached' });
//   });

});
