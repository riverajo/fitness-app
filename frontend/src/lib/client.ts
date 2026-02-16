import { Client, cacheExchange, fetchExchange } from '@urql/svelte';
import { authExchange } from '@urql/exchange-auth';
import { authStore } from '../state/auth.svelte';

export const client = new Client({
	url: '/query',
	exchanges: [
		cacheExchange,
		authExchange(async (utils) => {
			return {
				addAuthToOperation(operation) {
					const token = authStore.token;
					if (!token) {
						return operation;
					}

					return utils.appendHeaders(operation, {
						Authorization: `Bearer ${token}`
					});
				},
				didAuthError(error, _operation) {
					// Handle 401s from the server (e.g. invalid/expired token)
					const hasAuthError = error.graphQLErrors.some(
						(e) => e.message.includes('unauthorized') || e.extensions?.code === 'UNAUTHENTICATED'
					);
					return hasAuthError;
				},
				async refreshAuth() {
					try {
						const response = await fetch('/auth/refresh', {
							method: 'POST'
						});

						if (response.ok) {
							const data = await response.json();
							if (data.token) {
								authStore.setToken(data.token);
								return;
							}
						}

						// If we are here, response was not OK.
						// This means server rejected refresh (e.g. 401).
						console.log('[client] Refresh failed with status', response.status);
						authStore.clearToken();
						window.location.href = '/';
					} catch (e) {
						console.error('[client] Network error during refresh', e);
						// Do not clear token. Do not redirect.
					}
				}
			};
		}),
		fetchExchange
	]
});
