import { Client, cacheExchange, fetchExchange } from '@urql/svelte';

export const client = new Client({
    url: 'http://localhost:8080/query',
    fetchOptions: {
        credentials: 'include',
    },
    exchanges: [cacheExchange, fetchExchange],
});
