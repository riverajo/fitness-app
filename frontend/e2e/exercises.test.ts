import { test, expect } from '@playwright/test';

test.describe('Exercises', () => {
    test('should list exercises, search, paginate, and create new exercise', async ({ page }) => {
        // 1. Login (assuming we need to be logged in, though the current implementation might allow public access or auto-login in dev)
        // For now, let's assume we need to register/login or just go to the page if it's protected.
        // Given previous context, we might need to register first.
        await page.goto('/register');
        await page.fill('input[name="email"]', `test-${Date.now()}@example.com`);
        await page.fill('input[name="password"]', 'password123');
        await page.click('button[type="submit"]');
        await expect(page).toHaveURL('/dashboard');

        // 2. Navigate to Exercises
        await page.click('a[href="/exercises/new"]'); // Use the link from dashboard to go to create, or navigate to list
        // Wait, dashboard has "Create Exercise" link to /exercises/new.
        // Let's go to /exercises directly first to check listing.
        await page.goto('/exercises');
        await expect(page.locator('h1')).toHaveText('Exercises');

        // 3. Search (should be empty initially or have system exercises)
        // Let's assume there are some system exercises or we create one.
        // Let's create one first to be sure.
        await page.click('text=Create Exercise');
        await expect(page).toHaveURL('/exercises/new');

        const exerciseName = `Test Exercise ${Date.now()}`;
        await page.fill('input[name="name"]', exerciseName);
        await page.fill('textarea[name="description"]', 'A test description');
        await page.click('button[type="submit"]');

        // 4. Verify Redirection and Listing
        await expect(page).toHaveURL('/exercises');
        await expect(page.locator(`text=${exerciseName}`)).toBeVisible();

        // 5. Search
        await page.fill('input[placeholder="Search exercises..."]', exerciseName);
        await page.click('text=Search');
        await expect(page.locator(`text=${exerciseName}`)).toBeVisible();

        // 6. Pagination (Check buttons exist while we have results)
        await expect(page.locator('button:has-text("Previous")')).toBeVisible();
        await expect(page.locator('button:has-text("Next")')).toBeVisible();

        // 7. Search for something non-existent
        await page.fill('input[placeholder="Search exercises..."]', 'NonExistentExerciseXYZ');
        await page.click('text=Search');
        await expect(page.locator('text=No exercises found')).toBeVisible();
    });
});
