import { page } from 'vitest/browser';
import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import Page from './+page.svelte';

vi.mock('@urql/svelte', () => {
	return {
		getContextClient: () => ({
			query: vi.fn(),
			mutation: vi.fn(),
		}),
		gql: (t: any) => t,
	};
});

describe('/+page.svelte', () => {
	it('should render h2', async () => {
		render(Page);

		const heading = page.getByRole('heading', { level: 2 });
		await expect.element(heading).toBeInTheDocument();
	});
});
