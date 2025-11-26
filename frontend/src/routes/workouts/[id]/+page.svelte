<script lang="ts">
    import { page } from '$app/stores';
    import { gql, getContextClient, queryStore } from '@urql/svelte';
    import { Card, Button, Heading, Table, TableHead, TableBody, TableBodyRow, TableHeadCell, TableBodyCell, Spinner, Alert } from 'flowbite-svelte';

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
    <div class="mb-6 flex justify-between items-center">
        <Button color="light" href="/dashboard" size="xs">&larr; Back to Dashboard</Button>
        <Button color="alternative" href="/workouts/{workoutId}/edit" size="xs">Edit Workout</Button>
    </div>

    {#if $workoutQuery.fetching}
        <div class="text-center py-12">
            <Spinner />
            <p class="mt-2 text-gray-500">Loading workout details...</p>
        </div>
    {:else if $workoutQuery.error}
        <Alert color="red">
            <span class="font-medium">Error loading workout:</span> {$workoutQuery.error.message}
        </Alert>
    {:else if $workoutQuery.data?.getWorkoutLog}
        {@const workout = $workoutQuery.data.getWorkoutLog}
        
        <Card class="w-full max-w-none mb-6">
            <div class="mb-4">
                <Heading tag="h1" class="text-2xl sm:text-3xl mb-2">{workout.name}</Heading>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                    {new Date(workout.startTime).toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
                    &bull;
                    {new Date(workout.startTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                </p>
            </div>
            {#if workout.locationName || workout.generalNotes}
                <div class="border-t border-gray-200 pt-4 dark:border-gray-700">
                    <dl class="grid grid-cols-1 gap-x-4 gap-y-8 sm:grid-cols-2">
                        {#if workout.locationName}
                            <div class="sm:col-span-1">
                                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Location</dt>
                                <dd class="mt-1 text-sm text-gray-900 dark:text-white">{workout.locationName}</dd>
                            </div>
                        {/if}
                        {#if workout.generalNotes}
                            <div class="sm:col-span-2">
                                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Notes</dt>
                                <dd class="mt-1 text-sm text-gray-900 dark:text-white">{workout.generalNotes}</dd>
                            </div>
                        {/if}
                    </dl>
                </div>
            {/if}
        </Card>

        <div class="space-y-6">
            <Heading tag="h2" class="text-lg">Exercises</Heading>
            
            {#each workout.exerciseLogs as exerciseLog}
                <Card class="w-full max-w-none p-0 overflow-hidden">
                    <div class="px-4 py-5 sm:px-6 bg-gray-50 border-b border-gray-200 dark:bg-gray-700 dark:border-gray-600">
                        <Heading tag="h3" class="text-lg">{exerciseLog.uniqueExercise.name}</Heading>
                        {#if exerciseLog.notes}
                            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{exerciseLog.notes}</p>
                        {/if}
                    </div>
                    <div class="p-0">
                        <Table shadow>
                            <TableHead>
                                <TableHeadCell>Set</TableHeadCell>
                                <TableHeadCell>Weight (kg)</TableHeadCell>
                                <TableHeadCell>Reps</TableHeadCell>
                                <TableHeadCell>RPE</TableHeadCell>
                            </TableHead>
                            <TableBody>
                                {#each exerciseLog.sets as set, i}
                                    <TableBodyRow>
                                        <TableBodyCell>{i + 1}</TableBodyCell>
                                        <TableBodyCell>{set.weight}</TableBodyCell>
                                        <TableBodyCell>{set.reps}</TableBodyCell>
                                        <TableBodyCell>{set.rpe ?? '-'}</TableBodyCell>
                                    </TableBodyRow>
                                {/each}
                            </TableBody>
                        </Table>
                    </div>
                </Card>
            {/each}
            
            {#if workout.exerciseLogs.length === 0}
                <div class="text-center py-6">
                    <p class="text-gray-500 dark:text-gray-400">No exercises logged for this workout.</p>
                </div>
            {/if}
        </div>
    {:else}
        <div class="text-center py-12">
            <p class="text-gray-500 dark:text-gray-400">Workout not found.</p>
        </div>
    {/if}
</div>
