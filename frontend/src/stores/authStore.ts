import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const TOKEN_KEY = 'auth_token';

function createAuthStore() {
	// Initialize from localStorage if available
	const initialToken = browser ? localStorage.getItem(TOKEN_KEY) : null;
	const { subscribe, set } = writable<{ token: string | null }>({ token: initialToken });

	return {
		subscribe,
		setToken: (token: string) => {
			if (browser) {
				localStorage.setItem(TOKEN_KEY, token);
			}
			set({ token });
		},
		clearToken: () => {
			if (browser) {
				localStorage.removeItem(TOKEN_KEY);
			}
			set({ token: null });
		},
		// Helper to get current value synchronously (useful for non-reactive contexts like client.ts)
		getToken: () => {
			if (browser) {
				return localStorage.getItem(TOKEN_KEY);
			}
			return null;
		}
	};
}

export const authStore = createAuthStore();
