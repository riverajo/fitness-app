class Auth {
	#token = $state<string | null>(null);
	#isRestoring = $state(true);
	#isOffline = $state(false);

	// Public read-only accessors
	get token() {
		return this.#token;
	}

	get isRestoring() {
		return this.#isRestoring;
	}

	get isOffline() {
		return this.#isOffline;
	}

	setToken(token: string) {
		this.#token = token;
		// If we set a token, we assume we are online or at least have a valid session
		this.#isOffline = false;
	}

	clearToken() {
		this.#token = null;
		this.#isRestoring = false;
		this.#isOffline = false;
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
					this.#isOffline = false;
					return;
				}
			} else {
				// Server responded with an error (e.g. 401, 403, 500)
				// This usually means our session is invalid
				this.#token = null;
			}
		} catch (e) {
			console.error('[authStore] Failed to restore session', e);
			// Check if it's likely a network error
			// In many browsers, network errors are TypeErrors with specific messages,
			// but simply catching the error here (which catches network failures) is distinct from
			// the response.ok check above (which handles HTTP errors).

			// If we entered catch, fetch failed completely (network error, CORS, offline)
			// We should NOT clear the token here to support offline mode.
			this.#isOffline = true;
			this.#isRestoring = false;
			return;
		}
		// If we fail (and it wasn't a network error caught above), we are just not authenticated
		console.log('[authStore] restoreSession: finished (failed)');
		this.#isRestoring = false;
		this.#isOffline = false;
	}
}

export const authStore = new Auth();
export { Auth };
