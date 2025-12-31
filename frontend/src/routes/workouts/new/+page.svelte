<script lang="ts">
	import { gql, getContextClient } from '@urql/svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { Heading } from 'flowbite-svelte';
	import WorkoutForm from '../../../components/workouts/WorkoutForm.svelte';
	import { workoutStore } from '../../../stores/workoutStore';
	import { onMount, onDestroy } from 'svelte';

	const client = getContextClient();
	let error = '';

	// Reset store on mount to ensure fresh state for new workout
	onMount(() => {
		workoutStore.reset();
	});

	// Also reset on destroy to be clean
	onDestroy(() => {
		workoutStore.reset();
	});

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
		// The form now dispatches the full store state
		const storeState = event.detail;

		const input = {
			name: storeState.name,
			startTime: new Date().toISOString(),
			endTime: new Date().toISOString(),
			locationName: storeState.locationName,
			generalNotes: storeState.generalNotes,
			exerciseLogs: storeState.exerciseLogs.map(
				(log: {
					uniqueExerciseId: string;
					sets: {
						reps: number;
						weight: number;
						unit: string;
						rpe?: number | null;
						toFailure?: boolean | null;
						order: number;
					}[];
					notes: string;
				}) => ({
					uniqueExerciseId: log.uniqueExerciseId,
					sets: log.sets.map(
						(s: {
							reps: number;
							weight: number;
							unit: string;
							rpe?: number | null;
							toFailure?: boolean | null;
							order: number;
						}) => ({
							reps: s.reps,
							weight: s.weight,
							unit: s.unit,
							rpe: s.rpe,
							toFailure: s.toFailure,
							order: s.order
						})
					),
					notes: log.notes
				})
			)
		};

		const result = await client.mutation(createWorkoutLogMutation, { input }).toPromise();

		if (result.error) {
			error = result.error.message;
		} else {
			await goto(resolve('/dashboard'));
		}
	}
</script>

<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
	<Heading tag="h1" class="mb-6">Log Workout</Heading>
	<WorkoutForm on:submit={handleSubmit} {error} />
</div>
