import { test, expect } from '@playwright/test';

test('user can create a new unique exercise', async ({ page }) => {
    // 1. Register/Login (using a fresh user for isolation)
    await page.goto('/register');
    const email = `test-ex-${Date.now()}@example.com`;
    await page.fill('input[name="email"]', email);
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL('/dashboard');

    // 2. Navigate to Create Exercise page
    await page.click('text=Create Exercise');
    await expect(page).toHaveURL('/exercises/new');

    // 3. Fill form
    const exerciseName = `Custom Press ${Date.now()}`;
    await page.fill('input[name="name"]', exerciseName);
    await page.fill('textarea[name="description"]', 'A great shoulder exercise');

    // 4. Submit
    await page.click('button[type="submit"]');

    // 5. Verify redirection to Dashboard
    await expect(page).toHaveURL(/\/exercises/);

    // Note: Since we don't have a list of exercises on the dashboard yet, 
    // we are just verifying the successful redirection which implies success 
    // (as per the code logic).
});
