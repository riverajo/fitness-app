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
						// Attempt to refresh the token using the HTTP-only cookie
						const response = await fetch('/auth/refresh', {
							method: 'POST' // Using POST as it mutates state (rotates token)
						});

						if (response.ok) {
							const data = await response.json();
							if (data.token) {
								authStore.setToken(data.token);
								return; // success, urql will retry
							}
						}
					} catch (e) {
						console.error('Failed to refresh token', e);
					}

					// If refresh fails, logout
					authStore.clearToken();
					window.location.href = '/';
				}
			};
		}),
		fetchExchange
	]
});
