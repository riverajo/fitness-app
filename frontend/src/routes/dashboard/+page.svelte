<script lang="ts">
    import { gql, getContextClient, queryStore } from '@urql/svelte';

    const client = getContextClient();

    let limit = 5;
    let offset = 0;

    $: workoutsQuery = queryStore({
        client,
        query: gql`
            query ListWorkoutLogs($limit: Int, $offset: Int) {
                listWorkoutLogs(limit: $limit, offset: $offset) {
                    id
                    name
                    startTime
                    endTime
                    exerciseLogs {
                        uniqueExercise {
                            id
                        }
                        sets {
                            reps
                            weight
                        }
                    }
                }
            }
        `,
        variables: { limit, offset },
        requestPolicy: 'cache-and-network'
    });

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

    function nextPage() {
        offset += limit;
    }

    function prevPage() {
        if (offset >= limit) {
            offset -= limit;
        }
    }
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
            <div class="flex gap-2">
                <a href="/workouts/new" class="rounded-md bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-green-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600">Log Workout</a>
                <a href="/exercises/new" class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Create Exercise</a>
            </div>
        </div>

        <div class="mt-8">
            <h2 class="text-2xl font-bold tracking-tight text-gray-900">Past Workouts</h2>
            {#if $workoutsQuery.fetching}
                <p>Loading workouts...</p>
            {:else if $workoutsQuery.error}
                <p class="text-red-600">Error loading workouts: {$workoutsQuery.error.message}</p>
            {:else if $workoutsQuery.data?.listWorkoutLogs}
                <div class="mt-4 space-y-4">
                    {#each $workoutsQuery.data.listWorkoutLogs as workout}
                        <div class="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
                            <div class="flex items-center justify-between">
                                <h3 class="text-lg font-medium text-gray-900">
                                    <a href="/workouts/{workout.id}" class="hover:underline focus:outline-none">
                                        {workout.name}
                                    </a>
                                </h3>
                                <p class="text-sm text-gray-500">{new Date(workout.startTime).toLocaleDateString()}</p>
                            </div>
                            <p class="mt-2 text-sm text-gray-600">Exercises: {workout.exerciseLogs.length}</p>
                        </div>
                    {/each}
                    {#if $workoutsQuery.data.listWorkoutLogs.length === 0}
                        <p class="text-gray-500">No workouts found.</p>
                    {/if}
                </div>

                <div class="mt-6 flex justify-between">
                    <button
                        on:click={prevPage}
                        disabled={offset === 0}
                        class="rounded-md bg-gray-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-gray-500 disabled:opacity-50"
                    >
                        Previous
                    </button>
                    <button
                        on:click={nextPage}
                        disabled={$workoutsQuery.data.listWorkoutLogs.length < limit}
                        class="rounded-md bg-gray-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-gray-500 disabled:opacity-50"
                    >
                        Next
                    </button>
                </div>
            {/if}
        </div>
    {:else}
        <p>Please log in.</p>
    {/if}
</div>
