import { test, expect } from '@playwright/test';

test('home page has expected title', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'Sign in to your account' })).toBeVisible();
});
