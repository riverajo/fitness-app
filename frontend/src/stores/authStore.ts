import { writable } from 'svelte/store';

interface AuthState {
	token: string | null;
	isRestoring: boolean;
}

function createAuthStore() {
	// Initialize with no token and isRestoring false.
	// We rely on the app (e.g. layout onMount) to call restoreSession().
	const { subscribe, set, update } = writable<AuthState>({
		token: null,
		isRestoring: true
	});

	let currentToken: string | null = null;

	subscribe((state) => {
		currentToken = state.token;
	});

	return {
		subscribe,
		setToken: (token: string) => {
			update((s) => ({ ...s, token }));
		},
		clearToken: () => {
			set({ token: null, isRestoring: false });
		},
		// Helper to get current value synchronously
		getToken: () => {
			return currentToken;
		},
		restoreSession: async () => {
			console.log('[authStore] restoreSession: starting');
			update((s) => ({ ...s, isRestoring: true }));
			try {
				const response = await fetch('/auth/refresh', {
					method: 'POST'
				});
				console.log('[authStore] restoreSession: fetch complete', response.status);

				if (response.ok) {
					const data = await response.json();
					console.log('[authStore] restoreSession: data received', data);
					if (data.token) {
						update((s) => ({ ...s, token: data.token, isRestoring: false }));
						return;
					}
				}
			} catch (e) {
				console.error('[authStore] Failed to restore session', e);
			}
			// If we fail, we are just not authenticated
			console.log('[authStore] restoreSession: finished (failed)');
			update((s) => ({ ...s, token: null, isRestoring: false }));
		}
	};
}

export const authStore = createAuthStore();
