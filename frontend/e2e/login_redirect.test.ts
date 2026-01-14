import { test, expect } from '@playwright/test';

test('login flow redirects to dashboard and auto-redirects if logged in', async ({ page }) => {
	// 1. Register a new user (this will auto-login)
	await page.goto('/register');
	const email = `test-${Date.now()}@example.com`;
	await page.fill('input[name="email"]', email);
	await page.fill('input[name="password"]', 'password123');
	await page.click('button[type="submit"]');

	// Verify we are on dashboard (registration auto-redirects? or check register test logic)
	// The register test says it expects redirect to dashboard.
	await expect(page).toHaveURL('/dashboard');

	// 2. Clear localStorage to simulate logout (cookies are not used anymore)
	await page.evaluate(() => localStorage.removeItem('auth_token'));
	await page.reload(); // Reload to ensure state is cleared
	await page.goto('/'); // Go back to login page

	// 3. Login
	await page.fill('input[name="email"]', email);
	await page.fill('input[name="password"]', 'password123');
	await page.click('button[type="submit"]');

	// 4. Verify redirect to dashboard
	await expect(page).toHaveURL('/dashboard');

	// 5. Verify auto-redirect
	// Go back to login page while logged in
	await page.goto('/');
	// Should be redirected back to dashboard
	await expect(page).toHaveURL('/dashboard');
});

test('accessing dashboard without login redirects to login', async ({ page }) => {
	await page.goto('/dashboard');
	// Ensure we don't see dashboard content while verifying
	await expect(page.getByRole('heading', { name: 'Dashboard' })).not.toBeVisible();
	await expect(page).toHaveURL('/');
});
