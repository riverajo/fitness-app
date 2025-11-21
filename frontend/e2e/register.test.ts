import { test, expect } from '@playwright/test';

test('user can register successfully', async ({ page }) => {
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    page.on('pageerror', err => console.log('PAGE ERROR:', err));
    page.on('requestfailed', req => console.log('REQUEST FAILED:', req.url(), req.failure()?.errorText));

    // 1. Navigate to register page
    await page.goto('/register');

    // 2. Fill form
    const email = `test-${Date.now()}@example.com`;
    await page.fill('input[name="email"]', email);
    await page.fill('input[name="password"]', 'password123');

    // 3. Submit
    await page.click('button[type="submit"]');

    // 4. Verify redirection to home (or login)
    try {
        await expect(page).toHaveURL('/', { timeout: 5000 });
    } catch (e) {
        // If timeout, check if there's an error message on the page
        const errorMessage = await page.textContent('.text-red-800').catch(() => 'No error message found');
        console.log('Registration failed with error:', errorMessage);
        throw new Error(`Registration failed. URL stayed at /register. UI Error: ${errorMessage}`);
    }
});
