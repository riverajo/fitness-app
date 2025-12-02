import { test, expect } from '@playwright/test';

test('user can create a workout log', async ({ page }) => {
	// 1. Register a new user
	const email = `test-${Date.now()}@example.com`;
	const password = 'Password123!';

	await page.goto('/register');
	await page.fill('input[name="email"]', email);
	await page.fill('input[name="password"]', password);
	await page.click('button[type="submit"]');

	// Wait for redirection to dashboard
	await expect(page).toHaveURL('/dashboard');

	// 2. Create a unique exercise
	await page.click('text=Create Exercise');
	await expect(page).toHaveURL('/exercises/new');

	const exerciseName = `Test Press ${Date.now()}`;
	await page.fill('input[name="name"]', exerciseName);
	await page.fill('textarea[name="description"]', 'A test exercise');
	await page.click('button:has-text("Create Exercise")');

	// Wait for redirection to dashboard
	// Wait for redirection to exercises list
	await expect(page).toHaveURL(/\/exercises/);

	// Navigate back to dashboard to log workout
	await page.goto('/dashboard');

	// 3. Create a workout log
	await page.click('text=Log Workout');
	await expect(page).toHaveURL('/workouts/new');

	const workoutName = 'My E2E Workout';
	await page.fill('input[id="name"]', workoutName);

	// Search and select exercise
	await page.fill('input[placeholder="Search exercises..."]', exerciseName);
	await page.click('button:has-text("Search")');
	await page.click(`button:has-text("${exerciseName}")`);

	// Add a set
	await page.fill('input[id="reps"]', '10');
	await page.fill('input[id="weight"]', '50');
	await page.fill('input[id="rpe"]', '8');
	await page.click('button:has-text("Add Set")');

	// Confirm exercise
	await page.click('button:has-text("Done Adding Sets")');

	// Save workout
	await page.click('button:has-text("Save Workout")');

	// Verify redirection to dashboard
	await expect(page).toHaveURL('/dashboard');
});
