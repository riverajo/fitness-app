import { test, expect } from '@playwright/test';

test.describe('Edit Workout', () => {
    test('should allow user to edit an existing workout', async ({ page }) => {
        // 1. Register a new user
        const email = `edit-test-${Date.now()}@example.com`;
        const password = 'Password123!';

        await page.goto('/register');
        await page.fill('input[name="email"]', email);
        await page.fill('input[name="password"]', password);
        await page.click('button[type="submit"]');
        await expect(page).toHaveURL('/dashboard');

        // 2. Create a workout to edit
        await page.click('a[href="/workouts/new"]');
        await page.fill('input[id="name"]', 'Original Workout');
        await page.fill('input[id="location"]', 'Original Gym');

        // Add an exercise
        await page.fill('input[placeholder="Search exercises..."]', 'Bench Press');
        await page.click('button:has-text("Search")');
        await page.click('button:has-text("Bench Press")');

        // Add a set
        await page.fill('input[id="reps"]', '10');
        await page.fill('input[id="weight"]', '100');
        await page.click('button:has-text("Add Set")');
        await page.click('button:has-text("Done Adding Sets")');

        // Save
        await page.click('button:has-text("Save Workout")');
        await expect(page).toHaveURL('/dashboard');

        // 3. Navigate to the created workout
        await page.click('text=Original Workout');

        // 4. Click Edit
        await page.click('a:has-text("Edit Workout")');
        await expect(page).toHaveURL(/.*\/edit/);

        // 5. Modify the workout
        await page.fill('input[id="name"]', 'Updated Workout');
        await page.fill('input[id="location"]', 'Updated Gym');
        await page.fill('textarea[id="notes"]', 'Updated Notes');

        // Add another exercise (Squat)
        await page.fill('input[placeholder="Search exercises..."]', 'Squat');
        await page.click('button:has-text("Search")');
        await page.click('button:has-text("Squat")');

        await page.fill('input[id="reps"]', '5');
        await page.fill('input[id="weight"]', '120');
        await page.click('button:has-text("Add Set")');
        await page.click('button:has-text("Done Adding Sets")');

        // 6. Save changes
        await page.click('button:has-text("Update Workout")');

        // 7. Verify changes on details page
        await expect(page).toHaveURL(/\/workouts\/.*/);
        await expect(page.locator('h1')).toHaveText('Updated Workout');
        await expect(page.locator('text=Updated Gym')).toBeVisible();
        await expect(page.locator('text=Updated Notes')).toBeVisible();
        await expect(page.locator('text=Bench Press')).toBeVisible();
        await expect(page.locator('text=Squat')).toBeVisible();
    });
});
