import { page } from 'vitest/browser';
import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import Page from './+page.svelte';

vi.mock('@urql/svelte', () => {
	const mockClient = {
		query: vi.fn(),
		mutation: vi.fn()
	};

	// Mock chainable methods
	mockClient.query.mockReturnValue({
		toPromise: vi.fn().mockResolvedValue({ data: {} })
	});
	mockClient.mutation.mockReturnValue({
		toPromise: vi.fn().mockResolvedValue({ data: {} })
	});

	return {
		getContextClient: () => mockClient,
		gql: (t: TemplateStringsArray) => t
	};
});

describe('/+page.svelte', () => {
	it('should render h2', async () => {
		render(Page);

		const heading = page.getByRole('heading', { level: 2 });
		await expect.element(heading).toBeInTheDocument();
	});
});
