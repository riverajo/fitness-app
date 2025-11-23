import { test, expect } from '@playwright/test';

test('user can view past workouts with pagination', async ({ page }) => {
    // 1. Register a new user
    const email = `test-pagination-${Date.now()}@example.com`;
    const password = 'Password123!';

    await page.goto('/register');
    await page.fill('input[name="email"]', email);
    await page.fill('input[name="password"]', password);
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL('/dashboard');

    // 2. Create a unique exercise
    await page.click('text=Create Exercise');
    await expect(page).toHaveURL('/exercises/new');
    const exerciseName = `Pagination Exercise ${Date.now()}`;
    await page.fill('input[name="name"]', exerciseName);
    await page.click('button:has-text("Create Exercise")');
    await expect(page).toHaveURL('/exercises');
    await page.goto('/dashboard');

    // 3. Create multiple workouts to trigger pagination (limit is 5)
    // We need at least 6 workouts.
    for (let i = 1; i <= 6; i++) {
        await page.click('text=Log Workout');
        await expect(page).toHaveURL('/dashboard/create-workout');

        const workoutName = `Workout ${i}`;
        await page.fill('input[id="name"]', workoutName);

        // Select exercise
        await page.fill('input[placeholder="Search exercises..."]', exerciseName);
        await page.click('button:has-text("Search")');
        await page.click(`button:has-text("${exerciseName}")`);

        // Add set
        await page.fill('input[id="reps"]', '10');
        await page.fill('input[id="weight"]', '50');
        await page.click('button:has-text("Add Set")');
        await page.click('button:has-text("Done Adding Sets")');

        // Save
        await page.click('button:has-text("Save Workout")');
        await expect(page).toHaveURL('/dashboard');
    }

    // 4. Verify Dashboard List (First Page)
    // Wait for loading to finish
    await expect(page.locator('text=Loading workouts...')).not.toBeVisible();

    // Should show the latest 5 workouts (Workout 6 down to Workout 2)
    // Note: Default sort is usually by creation time or start time.
    // Assuming backend sorts by startTime desc (which we added in repository).

    // Check for Workout 6
    await expect(page.locator(`text=Workout 6`)).toBeVisible();
    // Check for Workout 2
    await expect(page.locator(`text=Workout 2`)).toBeVisible();
    // Workout 1 should NOT be visible on first page
    await expect(page.locator(`text=Workout 1`)).not.toBeVisible();

    // 5. Test Pagination (Next Page)
    await page.click('button:has-text("Next")');

    // Check for Workout 1
    await expect(page.locator(`text=Workout 1`)).toBeVisible();
    // Workout 6 should NOT be visible
    await expect(page.locator(`text=Workout 6`)).not.toBeVisible();

    // 6. Test Pagination (Previous Page)
    await page.click('button:has-text("Previous")');
    await expect(page.locator(`text=Workout 6`)).toBeVisible();
});
