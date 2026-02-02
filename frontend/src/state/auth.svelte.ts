class Auth {
	#token = $state<string | null>(null);
	#isRestoring = $state(true);

	// Public read-only accessors
	get token() {
		return this.#token;
	}

	get isRestoring() {
		return this.#isRestoring;
	}

	setToken(token: string) {
		this.#token = token;
	}

	clearToken() {
		this.#token = null;
		this.#isRestoring = false;
	}

	async restoreSession() {
		console.log('[authStore] restoreSession: starting');
		this.#isRestoring = true;
		try {
			const response = await fetch('/auth/refresh', {
				method: 'POST'
			});
			console.log('[authStore] restoreSession: fetch complete', response.status);

			if (response.ok) {
				const data = await response.json();
				console.log('[authStore] restoreSession: data received', data);
				if (data.token) {
					this.#token = data.token;
					this.#isRestoring = false;
					return;
				}
			}
		} catch (e) {
			console.error('[authStore] Failed to restore session', e);
		}
		// If we fail, we are just not authenticated
		console.log('[authStore] restoreSession: finished (failed)');
		this.#token = null;
		this.#isRestoring = false;
	}
}

export const authStore = new Auth();
export { Auth };
