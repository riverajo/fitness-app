<script lang="ts">
    import { page } from '$app/stores';
    import { gql, getContextClient, queryStore } from '@urql/svelte';
    import { goto } from '$app/navigation';
    import { Heading, Spinner, Alert } from 'flowbite-svelte';
    import WorkoutForm from '../../../../components/workouts/WorkoutForm.svelte';

    const client = getContextClient();
    const workoutId = $page.params.id;
    let error = '';

    $: workoutQuery = queryStore({
        client,
        query: gql`
            query GetWorkoutLogForEdit($id: ID!) {
                getWorkoutLog(id: $id) {
                    id
                    name
                    startTime
                    endTime
                    locationName
                    generalNotes
                    exerciseLogs {
                        uniqueExercise {
                            id
                            name
                        }
                        sets {
                            reps
                            weight
                            rpe
                            toFailure
                            order
                        }
                        notes
                    }
                }
            }
        `,
        variables: { id: workoutId }
    });

    const updateWorkoutLogMutation = gql`
        mutation UpdateWorkoutLog($input: UpdateWorkoutLogInput!) {
            updateWorkoutLog(input: $input) {
                id
                name
            }
        }
    `;

    async function handleSubmit(event: CustomEvent) {
        error = '';
        const formData = event.detail;

        const input = {
            id: workoutId,
            name: formData.name,
            locationName: formData.locationName,
            generalNotes: formData.generalNotes,
            exerciseLogs: formData.exerciseLogs
        };

        const result = await client.mutation(updateWorkoutLogMutation, { input }).toPromise();

        if (result.error) {
            error = result.error.message;
        } else {
            goto(`/workouts/${workoutId}`);
        }
    }
</script>

<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
    <Heading tag="h1" class="mb-6">Edit Workout</Heading>

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
        <WorkoutForm 
            initialData={$workoutQuery.data.getWorkoutLog} 
            submitLabel="Update Workout"
            on:submit={handleSubmit} 
            {error} 
        />
    {:else}
        <div class="text-center py-12">
            <p class="text-gray-500 dark:text-gray-400">Workout not found.</p>
        </div>
    {/if}
</div>
