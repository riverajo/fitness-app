import { test, expect } from '@playwright/test';

// Generate unique credentials for each test run to avoid collisions
const generateUser = () => {
	const timestamp = Date.now();
	return {
		email: `edit-user-${timestamp}@example.com`,
		password: 'Password123!'
	};
};

test.describe('Full Edit Mode', () => {
	test.beforeEach(async ({ page }) => {
		// Register and login
		const { email, password } = generateUser();
		await page.goto('/register');
		await page.fill('input[name="email"]', email);
		await page.fill('input[name="password"]', password);
		await page.click('button[type="submit"]');
		await expect(page).toHaveURL('/dashboard');
	});

	test('should allow creating, editing, and verifying a workout', async ({ page }) => {
		// 1. Create a "readonly" workout first
		await page.click('text=Log Workout');
		await expect(page).toHaveURL('/workouts/new');

		await page.fill('input[id="name"]', 'Chest Day');
		await page.fill('input[id="location"]', 'Gold Gym');

		// Add Exercise
		await page.fill('input[placeholder="Search exercises..."]', 'Bench Press');
		await page.click('button:has-text("Search")');
		// Wait for results and click
		await page.click('button:has-text("Bench Press")');

		// Add Set
		await page.fill('input[id="reps"]', '10');
		await page.fill('input[id="weight"]', '100');
		await page.click('button:has-text("Add Set")');
		await page.click('button:has-text("Done Adding Sets")');

		await page.click('button:has-text("Save Workout")');
		await expect(page).toHaveURL('/dashboard');

		// 2. Go to Detail View
		await page.click('text=Chest Day');
		await expect(page).toHaveURL(/\/workouts\/\w+/);

		// 3. Enter Edit Mode (Verify logic: button should be visible for new workout)
		await expect(page.locator('a:has-text("Edit Workout")')).toBeVisible();

		await page.click('a:has-text("Edit Workout")');
		await expect(page).toHaveURL(/\/edit$/);

		// 4. Modify Field (hydrated state check)
		await expect(page.locator('input[id="name"]')).toHaveValue('Chest Day');
		await page.fill('input[id="name"]', 'Chest Day - Pro');

		// 5. Modify Existing Set (Bound input check)
		// Locate the input for reps in the first added exercise.
		// Structure: Card > div > div > input[type=number]
		// We'll look for value "10" and "100" to be safe.
		const repsInput = page.locator('input[type="number"]').nth(0); // This assumes it's the first input in the list
		// Wait, "Add Exercise" inputs are also type=number.
		// But "Added Exercises" list is above "Add Exercise" card.
		// So nth(0) should be the first set's reps.

		await expect(repsInput).toHaveValue('10');
		await repsInput.fill('12'); // Change reps to 12

		// 6. Save Updates
		await page.click('button:has-text("Update Workout")');

		// 7. Verify Redirect and Content
		await expect(page).toHaveURL(/\/workouts\/\w+$/); // Should be back at detail, not edit
		await expect(page.locator('h1')).toHaveText('Chest Day - Pro');

		// Check the table for updated values
		// Detail view uses a Table, so we look for cell text.
		// The table cell for Reps is the 3rd column.
		await expect(page.locator('td', { hasText: '12' })).toBeVisible();

		// Verify original weight preserved
		await expect(page.locator('td', { hasText: '100' })).toBeVisible();
	});
});
