<script lang="ts">
	import { graphql } from '$lib/gql';
	import { queryStore } from '@urql/svelte';
	import { client } from '$lib/client';

	const workoutsQuery = graphql(`
		query ListWorkouts {
			listWorkoutLogs {
				id
				name
			}
		}
	`);

	const workouts = queryStore({
		client,
		query: workoutsQuery
	});
</script>

<h1>Welcome to SvelteKit</h1>

{#if $workouts.fetching}
	<p>Loading...</p>
{:else if $workouts.error}
	<p>Error: {$workouts.error.message}</p>
{:else}
	<ul>
		{#each $workouts.data?.listWorkoutLogs ?? [] as workout}
			<li>{workout.name}</li>
		{/each}
	</ul>
{/if}

<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
