import { test, expect } from '@playwright/test';

test.describe('Authentication Session Restoration', () => {
	// We need to register/login first to get the cookies
	test('should restore session on page reload', async ({ page }) => {
		// 1. Register a new user to ensure we have a valid session
		const email = `test-session-${Date.now()}@example.com`;
		const password = 'password123';

		// Use API to register for speed (if possible), or UI
		// Using UI here to be safe and mimicking real user flow
		await page.goto('/register');
		await page.getByLabel('Email').fill(email);
		await page.getByLabel('Password').fill(password);
		await page.getByRole('button', { name: 'Sign up' }).click();

		// Should be on dashboard
		await expect(page).toHaveURL('/dashboard');
		// Use heading to avoid ambiguity with nav link
		await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();

		// 2. Reload the page
		await page.reload();

		// 3. Verify we are still on dashboard and not redirected to login
		// The "Loading..." spinner might appear briefly
		await expect(page).toHaveURL('/dashboard');
		await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
	});

	test('should redirect to login if refresh token is missing/invalid', async ({ page }) => {
		// 1. Register/Login
		const email = `test-no-session-${Date.now()}@example.com`;
		const password = 'password123';

		await page.goto('/register');
		await page.getByLabel('Email').fill(email);
		await page.getByLabel('Password').fill(password);
		await page.getByRole('button', { name: 'Sign up' }).click();
		await expect(page).toHaveURL('/dashboard');

		// 2. Clear cookies to simulate expired/missing refresh token
		await page.context().clearCookies();

		// 3. Reload
		await page.reload();

		// 4. Should be redirected to login (or root which redirects to login?)
		// Root "/" redirects to dashboard if logged in, but we are not logged in.
		// Protected route "/dashboard" should redirect to "/" if not logged in.
		// Layout logic: if (!token) goto('/');
		// So we expect to end up at '/'
		await expect(page).toHaveURL('/');
	});

	test('should not store access token in localStorage', async ({ page }) => {
		const email = `test-localstorage-${Date.now()}@example.com`;
		const password = 'password123';

		await page.goto('/register');
		await page.getByLabel('Email').fill(email);
		await page.getByLabel('Password').fill(password);
		await page.getByRole('button', { name: 'Sign up' }).click();
		await expect(page).toHaveURL('/dashboard');

		// Check localStorage
		const tokenInStorage = await page.evaluate(() => localStorage.getItem('auth_token'));
		expect(tokenInStorage).toBeNull();
	});
});
