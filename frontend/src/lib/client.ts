import { Client, cacheExchange, fetchExchange } from '@urql/svelte';
import { PUBLIC_GRAPHQL_API } from '$env/static/private';

export const client = new Client({
    url: PUBLIC_GRAPHQL_API || 'http://localhost:8080/query',
    exchanges: [cacheExchange, fetchExchange],
});
