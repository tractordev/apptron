import { test as setup } from '@playwright/test';
import { MailSlurp } from 'mailslurp-client';
import * as fs from 'fs';

// Where the authenticated browser state (cookies + localStorage) will be saved.
// The chromium project in playwright.config.ts loads this at the start of every test,
// so tests don't need to sign in themselves.
const authFile = 'tests/.auth/user.json';

// A virtual authenticator simulates a hardware passkey device (like Touch ID or a USB key).
// This lets Playwright handle WebAuthn ceremonies automatically, with no real biometrics needed.
async function setupVirtualAuthenticator(page: any, context: any) {
  const cdp = await context.newCDPSession(page);
  await cdp.send('WebAuthn.enable', { enableUI: false });
  await cdp.send('WebAuthn.addVirtualAuthenticator', {
    options: {
      protocol: 'ctap2',                // CTAP2 is the protocol used by modern passkeys
      transport: 'internal',            // 'internal' = platform authenticator (e.g. Face ID / Touch ID)
      hasResidentKey: true,             // allows the credential to be stored on the authenticator
      hasUserVerification: true,        // supports biometric/PIN verification
      isUserVerified: true,             // auto-passes verification (no real biometrics required)
      automaticPresenceSimulation: true // auto-responds to "tap your key" prompts
    },
  });
  return cdp;
}

setup('create test account', async ({ page, context }) => {
  // Log browser console messages and failed network requests so we can diagnose
  // issues with the external Hanko API without staring at a blank spinner.
  page.on('console', (msg: any) => console.log('[browser]', msg.text()));
  page.on('requestfailed', (req: any) =>
    console.log('[failed request]', req.url(), req.failure()?.errorText)
  );

  await setupVirtualAuthenticator(page, context);
  await page.goto('/signin');

  // Wait for Hanko to finish initialising before interacting with it.
  const hankoAuth = page.locator('hanko-auth');
  await hankoAuth.getByRole('button', { name: 'Create account' }).waitFor({ state: 'visible', timeout: 30000 });
  await hankoAuth.getByRole('button', { name: 'Create account' }).click();

  // Create a real MailSlurp inbox so we can receive the Hanko verification email.
  const mailslurp = new MailSlurp({ apiKey: process.env.MAILSLURP_API_KEY!, basePath: 'https://api.mailslurp.com' });
  const inbox = await mailslurp.createInbox();

  await hankoAuth.getByLabel('Username').fill(`testuser${Date.now()}`);
  await hankoAuth.getByLabel('Email').fill(inbox.emailAddress!);
  await hankoAuth.getByRole('button', { name: 'Continue' }).click();

  // Hanko occasionally returns a transient "technical error" after form submission.
  // If that happens, throw so Playwright's retry (retries: 1) re-runs setup.
  const hankoError = hankoAuth.locator('#errorMessage');
  if (await hankoError.isVisible({ timeout: 3000 }).catch(() => false)) {
    throw new Error('Hanko returned a technical error — retrying');
  }

  // Wait for the passcode email and extract the 6-digit code.
  const email = await mailslurp.waitForLatestEmail(inbox.id!, 30000);
  const code = email.body!.match(/\d{6}/)![0];

  await hankoAuth.locator('input').first().click();
  await page.keyboard.type(code);

  // Hanko may auto-submit when all 6 digits are entered, making Continue disappear.
  // Wait for whichever comes first, then only click Continue if it's still there.
  const createPasskeyHeading = hankoAuth.locator('h1').filter({ hasText: 'Create a passkey' });
  const continueBtn = hankoAuth.getByRole('button', { name: 'Continue' });
  await Promise.race([
    createPasskeyHeading.waitFor({ state: 'visible' }),
    continueBtn.waitFor({ state: 'visible' }),
  ]);
  if (await continueBtn.isVisible()) {
    await continueBtn.click();
  }

  // The virtual authenticator intercepts navigator.credentials.create() automatically.
  await hankoAuth.getByRole('button', { name: 'Create a passkey' }).click();
  await page.waitForURL('**/dashboard**', { timeout: 30000 });

  // Save cookies + localStorage so all other tests start already logged in.
  fs.mkdirSync('tests/.auth', { recursive: true });
  await page.context().storageState({ path: authFile });
});
