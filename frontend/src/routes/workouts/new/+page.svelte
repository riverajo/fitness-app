<script lang="ts">
    import { gql, getContextClient } from '@urql/svelte';
    import { goto } from '$app/navigation';
    import { Heading } from 'flowbite-svelte';
    import WorkoutForm from '../../../components/workouts/WorkoutForm.svelte';

    const client = getContextClient();
    let error = '';

    const createWorkoutLogMutation = gql`
        mutation CreateWorkoutLog($input: CreateWorkoutLogInput!) {
            createWorkoutLog(input: $input) {
                id
                name
            }
        }
    `;

    async function handleSubmit(event: CustomEvent) {
        error = '';
        const formData = event.detail;

        const input = {
            name: formData.name,
            startTime: new Date().toISOString(),
            endTime: new Date().toISOString(),
            locationName: formData.locationName,
            generalNotes: formData.generalNotes,
            exerciseLogs: formData.exerciseLogs
        };

        const result = await client.mutation(createWorkoutLogMutation, { input }).toPromise();

        if (result.error) {
            error = result.error.message;
        } else {
            goto('/dashboard');
        }
    }
</script>

<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
    <Heading tag="h1" class="mb-6">Log Workout</Heading>
    <WorkoutForm on:submit={handleSubmit} {error} />
</div>
