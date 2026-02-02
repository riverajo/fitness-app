<script lang="ts">
	import { gql, getContextClient } from '@urql/svelte';
	import {
		Heading,
		Button,
		Card,
		Label,
		Input,
		Textarea,
		Checkbox,
		Alert,
		Listgroup,
		ListgroupItem
	} from 'flowbite-svelte';
	import { createEventDispatcher } from 'svelte';
	import type { UniqueExercise } from '../../lib/gql/graphql';
	import { workoutStore, type WorkoutState } from '../../state/workout.svelte';
	import WeightInput from './WeightInput.svelte';

	// Define a type for the exercises as they are being edited in the scratchpad
	interface ScratchpadSet {
		reps: number;
		weight: number;
		unit: string;
		rpe?: number | null;
		toFailure?: boolean | null;
		order: number;
	}

	export let submitLabel: string = 'Save Workout';
	export let error: string = '';

	const dispatch = createEventDispatcher<{
		submit: WorkoutState;
	}>();
	const client = getContextClient();

	let searchQuery = '';
	let searchResults: UniqueExercise[] = [];

	// Local state for the "Add Exercise" scratchpad
	let currentExercise: UniqueExercise | null = null;
	let currentSets: ScratchpadSet[] = [];

	// Scratchpad input fields
	let reps = 0;
	let weight = 0;
	let rpe = 0;
	let toFailure = false;

	const searchExercisesQuery = gql`
		query SearchUniqueExercises($query: String) {
			uniqueExercises(query: $query, limit: 10) {
				id
				name
				description
				isCustom
			}
		}
	`;

	async function handleSearch() {
		if (!searchQuery) return;
		const result = await client.query(searchExercisesQuery, { query: searchQuery }).toPromise();
		if (result.data) {
			searchResults = result.data.uniqueExercises;
		}
	}

	function selectExercise(exercise: UniqueExercise) {
		currentExercise = exercise;
		currentSets = [];
		searchResults = [];
		searchQuery = '';
	}

	function addSet() {
		currentSets = [
			...currentSets,
			{
				reps,
				weight,
				rpe,
				toFailure,
				order: currentSets.length + 1,
				unit: 'KILOGRAMS'
			}
		];
		// Reset scratchpad fields
		reps = 0;
		weight = 0;
		rpe = 0;
		toFailure = false;
	}

	function confirmExercise() {
		if (currentExercise && currentSets.length > 0) {
			// Commit to global store
			// Reassigning exerciseLogs array triggers reactivity in Svelte 5 state
			workoutStore.state.exerciseLogs = [
				...workoutStore.state.exerciseLogs,
				{
					uniqueExerciseId: currentExercise!.id,
					name: currentExercise!.name,
					sets: currentSets,
					notes: ''
				}
			];

			currentExercise = null;
			currentSets = [];
		}
	}

	function removeExercise(index: number) {
		workoutStore.state.exerciseLogs = workoutStore.state.exerciseLogs.filter((_, i) => i !== index);
	}

	function handleSubmit() {
		dispatch('submit', workoutStore.state);
	}
</script>

<div class="max-w-2xl space-y-6">
	{#if error}
		<Alert color="red">{error}</Alert>
	{/if}

	<!-- Workout Details -->
	<Card class="w-full max-w-none">
		<Heading tag="h2" class="mb-4 text-lg">Details</Heading>
		<div class="space-y-4">
			<div>
				<Label for="name" class="mb-2">Workout Name</Label>
				<Input
					type="text"
					id="name"
					bind:value={workoutStore.state.name}
					placeholder="e.g. Morning Lift"
				/>
			</div>
			<div>
				<Label for="location" class="mb-2">Location</Label>
				<Input type="text" id="location" bind:value={workoutStore.state.locationName} />
			</div>
			<div>
				<Label for="notes" class="mb-2">Notes</Label>
				<Textarea id="notes" bind:value={workoutStore.state.generalNotes} />
			</div>
		</div>
	</Card>

	<!-- Added Exercises List (Bound to Store) -->
	{#if workoutStore.state.exerciseLogs.length > 0}
		<div class="space-y-4">
			<Heading tag="h2" class="text-lg">Exercises</Heading>
			{#each workoutStore.state.exerciseLogs as exercise, exerciseIndex (exerciseIndex)}
				<Card class="w-full max-w-none bg-gray-50 dark:bg-gray-700">
					<div class="flex items-start justify-between">
						<h3 class="font-medium text-gray-900 dark:text-white">{exercise.name}</h3>
						<Button color="red" size="xs" onclick={() => removeExercise(exerciseIndex)}
							>Remove</Button
						>
					</div>
					<div class="mt-4 space-y-2">
						{#each exercise.sets as set (set.order)}
							<div class="flex items-center gap-2 text-sm">
								<span class="w-12 font-semibold text-gray-500">Set {set.order}</span>
								<div class="flex items-center gap-2">
									<Input type="number" size="sm" class="w-20" bind:value={set.reps} />
									<span class="text-gray-500">reps @</span>
									<WeightInput size="sm" bind:value={set.weight} />
								</div>
							</div>
						{/each}
					</div>
				</Card>
			{/each}
		</div>
	{/if}

	<!-- Add Exercise Section -->
	<Card class="w-full max-w-none">
		<Heading tag="h2" class="mb-4 text-lg">Add Exercise</Heading>

		{#if !currentExercise}
			<div class="mt-4 flex gap-2">
				<div class="flex-1">
					<Input
						type="text"
						bind:value={searchQuery}
						placeholder="Search exercises..."
						onkeydown={(e) => e.key === 'Enter' && handleSearch()}
					/>
				</div>
				<Button onclick={handleSearch} color="blue">Search</Button>
			</div>

			{#if searchResults.length > 0}
				<Listgroup class="mt-4">
					{#each searchResults as exercise (exercise.id)}
						<ListgroupItem class="p-0">
							<button
								type="button"
								class="flex w-full items-center px-4 py-2 text-left transition-colors duration-200 hover:bg-gray-100 dark:hover:bg-gray-600"
								onclick={() => selectExercise(exercise)}
							>
								<span class="font-medium">{exercise.name}</span>
								{#if exercise.description}
									<span class="ml-2 text-sm text-gray-500">- {exercise.description}</span>
								{/if}
							</button>
						</ListgroupItem>
					{/each}
				</Listgroup>
			{/if}
		{:else}
			<div class="mt-4 space-y-4">
				<div class="flex items-center justify-between">
					<h3 class="font-medium text-indigo-600 dark:text-indigo-400">{currentExercise.name}</h3>
					<Button color="light" size="xs" onclick={() => (currentExercise = null)}
						>Change Exercise</Button
					>
				</div>

				<!-- Sets List -->
				{#if currentSets.length > 0}
					<div class="space-y-2">
						{#each currentSets as set (set.order)}
							<div class="text-sm text-gray-700 dark:text-gray-300">
								Set {set.order}: {set.reps} reps @ {set.weight}kg
							</div>
						{/each}
					</div>
				{/if}

				<!-- Add Set Form -->
				<div class="grid grid-cols-2 gap-4 rounded-md bg-gray-50 p-3 dark:bg-gray-700">
					<div>
						<Label for="reps" class="mb-1 text-xs">Reps</Label>
						<Input type="number" id="reps" bind:value={reps} size="sm" />
					</div>
					<div>
						<Label for="weight" class="mb-1 text-xs">Weight</Label>
						<WeightInput size="sm" bind:value={weight} />
					</div>
					<div>
						<Label for="rpe" class="mb-1 text-xs">RPE</Label>
						<Input type="number" id="rpe" bind:value={rpe} size="sm" />
					</div>
					<div class="mt-6 flex items-center">
						<Checkbox id="toFailure" bind:checked={toFailure}>To Failure</Checkbox>
					</div>
					<div class="col-span-2">
						<Button onclick={addSet} color="light" class="w-full">Add Set</Button>
					</div>
				</div>

				<Button onclick={confirmExercise} class="w-full">Done Adding Sets</Button>
			</div>
		{/if}
	</Card>

	<div class="flex justify-end gap-4 pt-4">
		<Button color="light" href="/dashboard">Cancel</Button>
		<Button onclick={handleSubmit}>{submitLabel}</Button>
	</div>
</div>
