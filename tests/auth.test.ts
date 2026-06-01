import { test, expect } from '@playwright/test';
import * as fs from 'fs';
import { addFileViaTerminal } from './helpers';

function getProjectUrl(): string {
  return JSON.parse(fs.readFileSync('.auth/test-project.json', 'utf-8')).url;
}

function getProjectName(): string {
  const url = getProjectUrl();
  return new URL(url).pathname.split('/edit/')[1];
}

// By the time these tests run, setup.ts has already created a test account
// and saved the session to tests/.auth/user.json. The chromium project in
// playwright.config.ts loads that file automatically, so every test here
// starts fully logged in.

test.describe('Authenticated user', () => {

  test('lands on the dashboard after sign in', async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page.getByRole('heading', { name: 'Projects' })).toBeVisible();
  });

  test('can run a terminal command and see the result in the file explorer', async ({ page }) => {
    await page.goto(getProjectUrl());

    const editor = page.frameLocator('iframe');

    // wait for VS Code to mount and the explorer toolbar to be ready
    await editor.locator('.monaco-workbench').waitFor({ state: 'attached', timeout: 60000 });
    await editor.locator('[aria-label="Refresh Explorer"]').first().waitFor({ state: 'visible', timeout: 30000 });

    await addFileViaTerminal(page, editor, 'hello.txt');
  });

  test('can change project visibility to public via the dashboard context menu', async ({ page, browser }) => {
    await page.goto('/dashboard');

    // wait for session and projects to load
    await page.locator('#header-bar.signedin').waitFor({ state: 'visible', timeout: 15000 });
    await page.waitForFunction(() => customElements.get('html-component') !== undefined);

    // find the project card and open its context menu
    const card = page.locator('.project-card', { hasText: getProjectName() });
    await card.waitFor({ state: 'visible', timeout: 30000 });
    await card.locator('.ellipsis-btn').click();

    // click Edit Settings
    await card.getByRole('menuitem', { name: 'Edit Settings' }).click();

    // select Public visibility and save
    await page.locator('#vis-public').check();
    await expect(page.locator('#project-submit')).toBeEnabled({ timeout: 10000 });
    await page.locator('#project-submit').click();

    // click the project card to open the editor
    await card.click();
    await page.waitForURL('**/edit/**', { timeout: 15000 });

    // wait for VS Code and the explorer toolbar to be ready
    const env = page.frameLocator('iframe');
    await env.locator('.monaco-workbench').waitFor({ state: 'attached', timeout: 60000 });
    await env.locator('[aria-label="Refresh Explorer"]').first().waitFor({ state: 'visible', timeout: 30000 });

    // create a text file via the terminal
    await addFileViaTerminal(page, env, 'hello.txt');

    // open the Share modal — #share is inside the env iframe
    await env.locator('#share').click();
    await expect(env.locator('#share-url')).toBeVisible({ timeout: 10000 });

    // read the share URL from the input
    const shareUrl = await env.locator('#share-url').inputValue();

    // open a fresh browser session (no auth) and navigate to the share URL
    console.log('[test] opening public browser session, share URL:', shareUrl);
    const publicContext = await browser.newContext();
    const publicPage = await publicContext.newPage();
    await publicPage.goto(shareUrl);
    await publicPage.waitForURL(shareUrl, { timeout: 15000 });

    // wait for VS Code and the file explorer to load in the public session
    const publicEnv = publicPage.frameLocator('iframe');
    await publicEnv.locator('.monaco-workbench').waitFor({ state: 'attached', timeout: 60000 });
    await publicEnv.locator('[aria-label="Refresh Explorer"]').first().waitFor({ state: 'visible', timeout: 30000 });
    await publicPage.waitForTimeout(5000);

    await publicContext.close();
  });

  test('private project returns 404 for anonymous users', async ({ page, browser }) => {
    await page.goto('/dashboard');
    await page.locator('#header-bar.signedin').waitFor({ state: 'visible', timeout: 15000 });
    await page.waitForFunction(() => customElements.get('html-component') !== undefined);

    // open settings and set visibility to private
    const card = page.locator('.project-card', { hasText: getProjectName() });
    await card.waitFor({ state: 'visible', timeout: 30000 });
    await card.locator('.ellipsis-btn').click();
    await card.getByRole('menuitem', { name: 'Edit Settings' }).click();
    await page.locator('#vis-private').check();
    await expect(page.locator('#vis-private')).toBeChecked();
    await expect(page.locator('#project-submit')).toBeEnabled({ timeout: 10000 });
    await page.locator('#project-submit').click();

    // open the editor and verify the share URL is hidden when project is private
    await card.click();
    await page.waitForURL('**/edit/**', { timeout: 15000 });
    const env = page.frameLocator('iframe');
    await env.locator('.monaco-workbench').waitFor({ state: 'attached', timeout: 60000 });
    await env.locator('[aria-label="Refresh Explorer"]').first().waitFor({ state: 'visible', timeout: 30000 });
    await env.locator('#share').click();
    await expect(env.locator('#share-block')).toHaveAttribute('hidden');
    await page.goBack();
    await page.waitForURL('**/dashboard**', { timeout: 15000 });

    // open the project URL in a fresh anonymous session and expect a 404
    const editUrl = new URL(getProjectUrl());
    editUrl.search = '';
    console.log('[test] checking anonymous access to private project:', editUrl.toString());
    const publicContext = await browser.newContext();
    const publicPage = await publicContext.newPage();
    const response = await publicPage.goto(editUrl.toString());
    expect(response?.status()).toBe(404);
    await publicContext.close();
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

//     // wait for VS Code to mount — wb.create() injects .monaco-workbench into <vscode-workbench>
//     await page.frameLocator('iframe').locator('.monaco-workbench').waitFor({ state: 'attached', timeout: 60000 });
//   });

});
