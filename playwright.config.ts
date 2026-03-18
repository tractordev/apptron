import { defineConfig, devices } from '@playwright/test';
import * as dotenv from 'dotenv';

// Wrangler loads .env.local automatically, but Node (and Playwright) don't.
// We load it here so MAILSLURP_API_KEY is available in tests.
dotenv.config({ path: '.env.local' });

export default defineConfig({
  testDir: './tests',
  retries: 1,
  use: {
    baseURL: 'http://localhost:8788',
    actionTimeout: 30000,
  },
  // Virtual authenticators (for passkey/WebAuthn testing) only work in Chromium
  projects: [
    // Runs once before all other projects — creates the test account and saves
    // the session to tests/.auth/user.json so other tests can reuse the login.
    {
      name: 'setup',
      testMatch: /.*setup\.ts/,
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        // Load the saved session — every test starts already logged in.
        storageState: 'tests/.auth/user.json',
      },
      dependencies: ['setup'],
    },
  ],
  webServer: {
    command: 'wrangler dev --port=8788',
    url: 'http://localhost:8788',
    reuseExistingServer: !process.env.CI,
  },
});
