<script lang="ts">
    import { page } from '$app/stores';
    import { gql, getContextClient, queryStore } from '@urql/svelte';

    const client = getContextClient();
    const workoutId = $page.params.id;

    $: workoutQuery = queryStore({
        client,
        query: gql`
            query GetWorkoutLog($id: ID!) {
                getWorkoutLog(id: $id) {
                    id
                    name
                    startTime
                    endTime
                    locationName
                    generalNotes
                    exerciseLogs {
                        uniqueExercise {
                            name
                        }
                        sets {
                            reps
                            weight
                            rpe
                            toFailure
                        }
                        notes
                    }
                }
            }
        `,
        variables: { id: workoutId }
    });
</script>

<div class="mx-auto max-w-3xl px-4 py-6 sm:px-6 lg:px-8">
    <div class="mb-6">
        <a href="/dashboard" class="text-sm font-medium text-indigo-600 hover:text-indigo-500">
            &larr; Back to Dashboard
        </a>
    </div>

    {#if $workoutQuery.fetching}
        <div class="text-center py-12">
            <p class="text-gray-500">Loading workout details...</p>
        </div>
    {:else if $workoutQuery.error}
        <div class="rounded-md bg-red-50 p-4">
            <div class="flex">
                <div class="ml-3">
                    <h3 class="text-sm font-medium text-red-800">Error loading workout</h3>
                    <div class="mt-2 text-sm text-red-700">
                        <p>{$workoutQuery.error.message}</p>
                    </div>
                </div>
            </div>
        </div>
    {:else if $workoutQuery.data?.getWorkoutLog}
        {@const workout = $workoutQuery.data.getWorkoutLog}
        
        <div class="bg-white shadow overflow-hidden sm:rounded-lg mb-6">
            <div class="px-4 py-5 sm:px-6">
                <h1 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
                    {workout.name}
                </h1>
                <p class="mt-1 max-w-2xl text-sm text-gray-500">
                    {new Date(workout.startTime).toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
                    &bull;
                    {new Date(workout.startTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                </p>
            </div>
            {#if workout.locationName || workout.generalNotes}
                <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
                    <dl class="grid grid-cols-1 gap-x-4 gap-y-8 sm:grid-cols-2">
                        {#if workout.locationName}
                            <div class="sm:col-span-1">
                                <dt class="text-sm font-medium text-gray-500">Location</dt>
                                <dd class="mt-1 text-sm text-gray-900">{workout.locationName}</dd>
                            </div>
                        {/if}
                        {#if workout.generalNotes}
                            <div class="sm:col-span-2">
                                <dt class="text-sm font-medium text-gray-500">Notes</dt>
                                <dd class="mt-1 text-sm text-gray-900">{workout.generalNotes}</dd>
                            </div>
                        {/if}
                    </dl>
                </div>
            {/if}
        </div>

        <div class="space-y-6">
            <h2 class="text-lg font-medium leading-6 text-gray-900">Exercises</h2>
            
            {#each workout.exerciseLogs as exerciseLog}
                <div class="bg-white shadow sm:rounded-lg overflow-hidden">
                    <div class="px-4 py-5 sm:px-6 bg-gray-50 border-b border-gray-200">
                        <h3 class="text-lg font-medium leading-6 text-gray-900">
                            {exerciseLog.uniqueExercise.name}
                        </h3>
                        {#if exerciseLog.notes}
                            <p class="mt-1 text-sm text-gray-500">{exerciseLog.notes}</p>
                        {/if}
                    </div>
                    <div class="px-4 py-5 sm:p-6">
                        <div class="flex flex-col">
                            <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
                                    <table class="min-w-full divide-y divide-gray-300">
                                        <thead>
                                            <tr>
                                                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6 md:pl-0">Set</th>
                                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Weight (kg)</th>
                                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Reps</th>
                                                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">RPE</th>
                                            </tr>
                                        </thead>
                                        <tbody class="divide-y divide-gray-200">
                                            {#each exerciseLog.sets as set, i}
                                                <tr>
                                                    <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6 md:pl-0">{i + 1}</td>
                                                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{set.weight}</td>
                                                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{set.reps}</td>
                                                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{set.rpe ?? '-'}</td>
                                                </tr>
                                            {/each}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            {/each}
            
            {#if workout.exerciseLogs.length === 0}
                <div class="text-center py-6 bg-white shadow sm:rounded-lg">
                    <p class="text-gray-500">No exercises logged for this workout.</p>
                </div>
            {/if}
        </div>
    {:else}
        <div class="text-center py-12">
            <p class="text-gray-500">Workout not found.</p>
        </div>
    {/if}
</div>
