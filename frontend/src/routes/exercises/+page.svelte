<script lang="ts">
    import { gql, getContextClient, queryStore } from '@urql/svelte';
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';

    const client = getContextClient();

    // Reactive parameters from URL
    let searchQuery = $derived($page.url.searchParams.get('q') || '');
    let currentPage = $derived(parseInt($page.url.searchParams.get('page') || '1'));
    let limit = 50;
    let offset = $derived((currentPage - 1) * limit);

    // Query store
    let exercisesQuery = $derived(queryStore({
        client,
        query: gql`
            query UniqueExercises($query: String, $limit: Int, $offset: Int) {
                uniqueExercises(query: $query, limit: $limit, offset: $offset) {
                    id
                    name
                    description
                    isCustom
                }
            }
        `,
        variables: { query: searchQuery, limit, offset },
        requestPolicy: 'cache-and-network'
    }));

    let searchInput = $state($page.url.searchParams.get('q') || '');

    function handleSearch() {
        goto(`/exercises?q=${searchInput}&page=1`);
    }

    function handlePageChange(newPage: number) {
        goto(`/exercises?q=${searchQuery}&page=${newPage}`);
    }

    // Update search input when URL changes (e.g. back button)
    $effect(() => {
        searchInput = searchQuery;
    });
</script>

<div class="container mx-auto p-4 max-w-2xl">
    <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold">Exercises</h1>
        <a href="/exercises/new" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition">
            Create Exercise
        </a>
    </div>

    <div class="mb-6 flex gap-2">
        <input
            type="text"
            bind:value={searchInput}
            placeholder="Search exercises..."
            class="flex-1 p-2 border rounded bg-gray-800 border-gray-700 text-white"
            onkeydown={(e) => e.key === 'Enter' && handleSearch()}
        />
        <button onclick={handleSearch} class="bg-gray-700 text-white px-4 py-2 rounded hover:bg-gray-600">
            Search
        </button>
    </div>

    {#if $exercisesQuery.fetching}
        <div class="text-center py-8">Loading...</div>
    {:else if $exercisesQuery.error}
        <div class="text-center py-8 text-red-500">Error: {$exercisesQuery.error.message}</div>
    {:else if $exercisesQuery.data?.uniqueExercises.length === 0}
        <div class="text-center py-8 text-gray-400">No exercises found.</div>
    {:else}
        <div class="space-y-2">
            {#each $exercisesQuery.data?.uniqueExercises || [] as exercise}
                <div class="p-4 bg-gray-800 rounded border border-gray-700 flex justify-between items-center">
                    <div>
                        <h3 class="font-semibold text-lg">{exercise.name}</h3>
                        {#if exercise.description}
                            <p class="text-sm text-gray-400">{exercise.description}</p>
                        {/if}
                    </div>
                    {#if exercise.isCustom}
                        <span class="text-xs bg-purple-900 text-purple-200 px-2 py-1 rounded">Custom</span>
                    {:else}
                        <span class="text-xs bg-gray-700 text-gray-300 px-2 py-1 rounded">System</span>
                    {/if}
                </div>
            {/each}
        </div>

        <div class="mt-6 flex justify-center gap-4">
            <button
                disabled={currentPage <= 1}
                onclick={() => handlePageChange(currentPage - 1)}
                class="px-4 py-2 bg-gray-700 rounded disabled:opacity-50 hover:bg-gray-600"
            >
                Previous
            </button>
            <span class="py-2">Page {currentPage}</span>
            <button
                onclick={() => handlePageChange(currentPage + 1)}
                class="px-4 py-2 bg-gray-700 rounded hover:bg-gray-600"
            >
                Next
            </button>
        </div>
    {/if}
</div>
