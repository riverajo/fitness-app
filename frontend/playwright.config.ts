import { defineConfig } from '@playwright/test';

export default defineConfig({
	webServer: process.env.CI ? undefined : {
		command: 'npm run build && npm run preview',
		port: 4173
	},
	testDir: 'e2e',
	use: {
		baseURL: process.env.BASE_URL || 'http://localhost:4173'
	}
});
