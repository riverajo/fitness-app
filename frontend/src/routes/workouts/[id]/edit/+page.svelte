<script lang="ts">
	import { page } from '$app/stores';
	import { getContextClient } from '@urql/svelte';
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { Heading, Spinner, Alert } from 'flowbite-svelte';
	import WorkoutForm from '../../../../components/workouts/WorkoutForm.svelte';
	import { workoutStore } from '../../../../stores/workoutStore';

	const client = getContextClient();
	const workoutId = $page.params.id ?? '';

	let loading = true;
	let error = '';

	import { gql } from '@urql/svelte';

	const UPDATE_WORKOUT_MUTATION = gql`
		mutation UpdateWorkoutLog($input: UpdateWorkoutLogInput!) {
			updateWorkoutLog(input: $input) {
				id
				name
			}
		}
	`;

	onMount(async () => {
		try {
			await workoutStore.loadWorkoutForEditing(client, workoutId);
			loading = false;
		} catch (e) {
			if (e instanceof Error) {
				error = e.message;
			} else {
				error = 'An unexpected error occurred';
			}
			loading = false;
		}
	});

	interface StoreStateLog {
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
	}

	async function handleUpdate(event: CustomEvent) {
		const storeState = event.detail;

		const input = {
			id: workoutId,
			name: storeState.name,
			locationName: storeState.locationName,
			generalNotes: storeState.generalNotes,
			startTime: storeState.startTime,
			endTime: storeState.endTime, // Now present in store

			exerciseLogs: storeState.exerciseLogs.map((log: StoreStateLog) => ({
				uniqueExerciseId: log.uniqueExerciseId,
				sets: log.sets.map((s) => ({
					reps: s.reps,
					weight: s.weight,
					unit: s.unit,
					rpe: s.rpe,
					toFailure: s.toFailure,
					order: s.order
				})),
				notes: log.notes
			}))
		};

		const result = await client.mutation(UPDATE_WORKOUT_MUTATION, { input }).toPromise();

		if (result.error) {
			error = result.error.message;
		} else {
			await goto(resolve(`/workouts/${workoutId}`));
		}
	}
</script>

<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
	<Heading tag="h1" class="mb-6">Edit Workout</Heading>

	{#if loading}
		<div class="flex justify-center p-8">
			<Spinner />
		</div>
	{:else if error}
		<Alert color="red">
			<span class="font-medium">Error:</span>
			{error}
		</Alert>
		<div class="mt-4">
			<a href={resolve(`/workouts/${workoutId}`)} class="text-blue-600 hover:underline"
				>&larr; Back to Workout</a
			>
		</div>
	{:else}
		<WorkoutForm on:submit={handleUpdate} submitLabel="Update Workout" {error} />
	{/if}
</div>
