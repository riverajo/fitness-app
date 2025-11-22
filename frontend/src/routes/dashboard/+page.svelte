<script lang="ts">
    import { gql, getContextClient, queryStore } from '@urql/svelte';

    const client = getContextClient();

    const meQuery = queryStore({
        client,
        query: gql`
            query Me {
                me {
                    id
                    email
                }
            }
        `
    });
</script>

<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
    {#if $meQuery.fetching}
        <p>Loading...</p>
    {:else if $meQuery.error}
        <p class="text-red-600">Error: {$meQuery.error.message}</p>
    {:else if $meQuery.data?.me}
        <h1 class="text-3xl font-bold tracking-tight text-gray-900">Dashboard</h1>
        <div class="mt-6 flex gap-4">
            <p class="text-lg text-gray-700">Welcome, <span class="font-semibold">{$meQuery.data.me.email}</span>!</p>
            <a href="/exercises/new" class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Create Exercise</a>
        </div>
    {:else}
        <p>Please log in.</p>
    {/if}
</div>
