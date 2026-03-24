import { test as teardown } from '@playwright/test';
import * as fs from 'fs';

const userFile = 'tests/.auth/test-user.json';

teardown('delete test account', async () => {
  if (!fs.existsSync(userFile)) {
    console.log('[teardown] No test user file found, skipping.');
    return;
  }

  const { email } = JSON.parse(fs.readFileSync(userFile, 'utf-8'));
  const adminApiKey = process.env.HANKO_ADMIN_API_KEY!;
  const adminUrl = `${process.env.AUTH_URL}/admin/users`;;

  const headers = {
    Authorization: `Bearer ${adminApiKey}`,
    'Content-Type': 'application/json',
  };

  // Look up the user by email to get their ID.
  const listRes = await fetch(`${adminUrl}?email=${encodeURIComponent(email)}`, { headers });
  if (!listRes.ok) {
    throw new Error(`Failed to list users: ${listRes.status} ${await listRes.text()}`);
  }

  const users = await listRes.json() as Array<{ id: string }>;
  if (users.length === 0) {
    console.log(`[teardown] No Hanko user found for ${email}, skipping delete.`);
    return;
  }

  const userId = users[0].id;

  // Delete the user.
  const deleteRes = await fetch(`${adminUrl}/${userId}`, { method: 'DELETE', headers });
  if (!deleteRes.ok && deleteRes.status !== 404) {
    throw new Error(`Failed to delete user ${userId}: ${deleteRes.status} ${await deleteRes.text()}`);
  }

  // Verify the user is gone — expect a 404.
  const verifyRes = await fetch(`${adminUrl}/${userId}`, { headers });
  if (verifyRes.status !== 404) {
    throw new Error(`Expected user ${userId} to be deleted but got status ${verifyRes.status}`);
  }

  console.log(`[teardown] Deleted and verified removal of Hanko user ${userId} (${email})`);
  fs.unlinkSync(userFile);
});
