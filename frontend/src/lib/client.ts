import { Client, cacheExchange, fetchExchange } from '@urql/svelte';
import { authExchange } from '@urql/exchange-auth';
import { authStore } from '../stores/authStore';

export const client = new Client({
	url: '/query',
	exchanges: [
		cacheExchange,
		authExchange(async (utils) => {
			return {
				addAuthToOperation(operation) {
					const token = authStore.getToken();
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
					// If we get an auth error (401), it means our token is invalid/expired.
					// We should remove it and redirect.
					localStorage.removeItem('auth_token');
					window.location.href = '/';
				}
			};
		}),
		fetchExchange
	]
});
