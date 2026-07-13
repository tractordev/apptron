import { expect, Page, FrameLocator } from '@playwright/test';

export async function createProject(page: Page): Promise<string> {
  await page.goto('/dashboard');
  await page.locator('#header-bar.signedin').waitFor({ state: 'visible', timeout: 15000 });
  await page.waitForFunction(() => customElements.get('html-component') !== undefined);
  await page.getByRole('button', { name: 'Create Project' }).click();
  await page.getByRole('heading', { name: 'New Project' }).waitFor({ state: 'visible', timeout: 30000 });
  const projectName = `project${Date.now()}`;
  await page.getByLabel('Project name').fill(projectName);
  await expect(page.locator('#project-submit')).toBeEnabled({ timeout: 15000 });
  await page.locator('#project-submit').click();
  await page.waitForURL('**/edit/**', { timeout: 15000 });
  return page.url();
}

export async function addFileViaTerminal(page: Page, env: FrameLocator, filename: string) {
  await env.locator('.xterm-helper-textarea').click();
  await page.keyboard.type(`echo "hello world" > ${filename}`);
  await page.keyboard.press('Enter');
  await page.waitForTimeout(2000);
  await env.locator('[aria-label="Refresh Explorer"]').first().click();
  await expect(env.locator(`.monaco-list-row[aria-label*="${filename}"]`)).toBeVisible({ timeout: 15000 });
}
