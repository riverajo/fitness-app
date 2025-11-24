import { test, expect } from '@playwright/test';

test('user can view workout details', async ({ page }) => {
    // 1. Register a new user
    const email = `test-details-${Date.now()}@example.com`;
    const password = 'Password123!';

    await page.goto('/register');
    await page.fill('input[name="email"]', email);
    await page.fill('input[name="password"]', password);
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL('/dashboard');

    // 2. Create a unique exercise
    await page.click('text=Create Exercise');
    await expect(page).toHaveURL('/exercises/new');
    const exerciseName = `Details Exercise ${Date.now()}`;
    await page.fill('input[name="name"]', exerciseName);
    await page.click('button:has-text("Create Exercise")');
    await expect(page).toHaveURL('/exercises');
    await page.goto('/dashboard');

    // 3. Create a workout
    await page.click('text=Log Workout');
    await expect(page).toHaveURL('/dashboard/create-workout');

    const workoutName = `Details Workout ${Date.now()}`;
    await page.fill('input[id="name"]', workoutName);

    // Select exercise
    await page.fill('input[placeholder="Search exercises..."]', exerciseName);
    await page.click('button:has-text("Search")');
    await page.click(`button:has-text("${exerciseName}")`);

    // Add set
    await page.fill('input[id="reps"]', '12');
    await page.fill('input[id="weight"]', '60');
    await page.click('button:has-text("Add Set")');
    await page.click('button:has-text("Done Adding Sets")');

    // Save
    await page.click('button:has-text("Save Workout")');
    await expect(page).toHaveURL('/dashboard');

    // 4. Navigate to details page
    // Wait for the workout to appear in the list
    await expect(page.locator(`text=${workoutName}`)).toBeVisible();

    // Click on the workout name (or card) to go to details. 
    // Note: The dashboard implementation currently doesn't link to the details page.
    // I need to update the dashboard to link to the details page first, or just navigate directly for now.
    // Let's navigate directly first to test the page itself, then I'll update the dashboard.
    // Actually, I should update the dashboard to link to it as part of "Implement workouts/[id] page" implies it's accessible.
    // But for this test, I'll find the link if I added it, or just click the card if I make it clickable.
    // Let's assume I'll make the workout name a link in the dashboard.

    // Wait, I haven't updated the dashboard yet! I should do that.
    // For now, let's just try to navigate by clicking, and if it fails, I'll know I need to update the dashboard.
    // But wait, the dashboard currently just displays text.
    // I will update the dashboard in the next step. For this test, I will assume the name is a link.

    // Let's update the test to expect a link.
    await page.click(`text=${workoutName}`);

    // 5. Verify Details Page
    await expect(page.locator(`h1:has-text("${workoutName}")`)).toBeVisible();
    await expect(page.locator(`text=${exerciseName}`)).toBeVisible();
    await expect(page.locator('text=12')).toBeVisible(); // Reps
    await expect(page.locator('text=60')).toBeVisible(); // Weight
});
