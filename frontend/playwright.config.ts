import { defineConfig } from '@playwright/test';

export default defineConfig({
	webServer: process.env.CI
		? undefined
		: {
				command: 'pnpm run build && pnpm run preview',
				port: 4173
			},
	testDir: 'e2e',
	reporter: [['html', { open: 'never' }]],
	retries: process.env.CI ? 2 : 0,
	use: {
		baseURL: process.env.BASE_URL || 'http://localhost:4173',
		ignoreHTTPSErrors: true,
		trace: 'retain-on-failure',
		screenshot: 'only-on-failure',
		video: 'retain-on-failure'
	}
});
