import { test, expect } from '@playwright/test';

test('WeightInput handles dual units and conversion correctly', async ({ page }) => {
	// 1. Setup: Register user
	const email = `weight-test-${Date.now()}@example.com`;
	const password = 'Password123!';

	await page.goto('/register');
	await page.fill('input[name="email"]', email);
	await page.fill('input[name="password"]', password);
	await page.click('button[type="submit"]');
	await expect(page).toHaveURL('/dashboard');

	// 2. Setup: Create Exercise
	await page.goto('/exercises/new');
	const exerciseName = `Weight Test ${Date.now()}`;
	await page.fill('input[name="name"]', exerciseName);
	await page.click('button:has-text("Create Exercise")');

	// 3. Go to Workout Log
	await page.goto('/workouts/new');
	await page.fill('input[id="name"]', 'Weight Input Test');

	// Select exercise
	await page.fill('input[placeholder="Search exercises..."]', exerciseName);
	await page.click('button:has-text("Search")');
	await page.click(`button:has-text("${exerciseName}")`);

	// 4. Test Inputs
	// Enter 100 Lbs
	await page.getByTestId('weight-input-lbs').fill('100');
	// Enter 10 Kgs
	await page.getByTestId('weight-input-kgs').fill('10');

	// Add Set
	await page.fill('input[id="reps"]', '5');
	await page.click('button:has-text("Add Set")');

	// 5. Verify Added Set
	// 100 lbs = 45.3592... kg
	// Total = 45.36 + 10 = 55.36 kg (rounded)
	await expect(page.locator('text=55.36kg')).toBeVisible();

	// 6. Verify Reset
	// Inputs should be back to 0 or empty (placeholder '0')
	await expect(page.getByTestId('weight-input-lbs')).toHaveValue('0');
	await expect(page.getByTestId('weight-input-kgs')).toHaveValue((0).toString()); // Because parent resets 'weight' to 0, and kgs = value
});
