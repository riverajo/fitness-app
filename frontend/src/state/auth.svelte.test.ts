import { describe, it, expect, beforeEach, vi } from 'vitest';
import { Auth } from './auth.svelte';

describe('Auth Store', () => {
	let auth: Auth;

	beforeEach(() => {
		auth = new Auth();
		vi.stubGlobal('fetch', vi.fn());
	});

	it('initializes with default state', () => {
		expect(auth.token).toBe(null);
		expect(auth.isRestoring).toBe(true);
	});

	it('sets token correctly', () => {
		auth.setToken('test-token');
		expect(auth.token).toBe('test-token');
	});

	it('clears token correctly', () => {
		auth.setToken('test-token');
		auth.clearToken();
		expect(auth.token).toBe(null);
		expect(auth.isRestoring).toBe(false);
	});

	describe('restoreSession', () => {
		it('updates token on successful refresh', async () => {
			const mockSuccessResponse = {
				ok: true,
				status: 200,
				json: async () => ({ token: 'refreshed-token' })
			};
			vi.mocked(fetch).mockResolvedValue(mockSuccessResponse as Response);

			await auth.restoreSession();

			expect(auth.token).toBe('refreshed-token');
			expect(auth.isRestoring).toBe(false);
			expect(fetch).toHaveBeenCalledWith('/auth/refresh', { method: 'POST' });
		});

		it('clears token on failed refresh', async () => {
			const mockErrorResponse = {
				ok: false,
				status: 401,
				json: async () => ({})
			};
			vi.mocked(fetch).mockResolvedValue(mockErrorResponse as Response);

			await auth.restoreSession();

			expect(auth.token).toBe(null);
			expect(auth.isRestoring).toBe(false);
		});

		it('clears token on network error', async () => {
			vi.mocked(fetch).mockRejectedValue(new Error('Network error'));

			await auth.restoreSession();

			expect(auth.token).toBe(null);
			expect(auth.isRestoring).toBe(false);
		});
	});
});
