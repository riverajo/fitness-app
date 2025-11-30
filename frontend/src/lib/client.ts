import { Client, cacheExchange, fetchExchange } from '@urql/svelte';

export const client = new Client({
    url: '/query',
    fetchOptions: {
        credentials: 'include',
    },
    exchanges: [cacheExchange, fetchExchange],
});
