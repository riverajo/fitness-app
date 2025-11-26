import { test, expect } from '@playwright/test';

test('login flow redirects to dashboard and auto-redirects if logged in', async ({ page, context }) => {
    // 1. Register a new user (this will auto-login)
    await page.goto('/register');
    const email = `test-${Date.now()}@example.com`;
    await page.fill('input[name="email"]', email);
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');

    // Verify we are on dashboard (registration auto-redirects? or check register test logic)
    // The register test says it expects redirect to dashboard.
    await expect(page).toHaveURL('/dashboard');

    // 2. Clear cookies to simulate logout
    await context.clearCookies();
    await page.reload(); // Reload to ensure state is cleared
    await page.goto('/'); // Go back to login page

    // 3. Login
    await page.fill('input[name="email"]', email);
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');

    // 4. Verify redirect to dashboard (This should fail currently as it stays on login page)
    await expect(page).toHaveURL('/dashboard');

    // 5. Verify auto-redirect
    // Go back to login page while logged in
    await page.goto('/');
    // Should be redirected back to dashboard
    await expect(page).toHaveURL('/dashboard');
});
