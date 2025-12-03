import { Client, cacheExchange, fetchExchange, mapExchange } from '@urql/svelte';
import { authExchange } from '@urql/exchange-auth';

export const client = new Client({
	url: '/query',
	exchanges: [
		cacheExchange,
		authExchange(async (utils) => {
			return {
				getAuth: async () => {
					// We don't manage tokens client-side, so we just return null.
					// This tells authExchange we are "ready" to make requests.
					return null;
				},
				addAuthToOperation(operation) {
					// We rely on HttpOnly cookies, so we just ensure credentials are included.
					// fetchOptions handles this, but we can double check or just return operation.
					return operation;
				},
				didAuthError(error, operation) {
					// Ignore auth errors for the 'Me' query to prevent infinite loops,
					// as it's expected to fail when not logged in.
					if (operation.kind === 'query' && operation.query.definitions.some((d: any) => d.name?.value === 'Me')) {
						return false;
					}
					return error.graphQLErrors.some((e) => e.message.includes('unauthorized') || e.extensions?.code === 'UNAUTHENTICATED');
				},
				async refreshAuth() {
					// If we get an auth error (401), it means our cookie is invalid/expired.
					// We should redirect to login.
					window.location.href = '/';
				}
			};
		}),
		fetchExchange
	],
	fetchOptions: {
		credentials: 'include'
	}
});
