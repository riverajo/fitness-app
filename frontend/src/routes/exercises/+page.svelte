<script lang="ts">
	import { gql, getContextClient, queryStore } from '@urql/svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { Heading, Button, Input, Card, Badge, Spinner, Alert } from 'flowbite-svelte';

	const client = getContextClient();

	// Reactive parameters from URL
	let searchQuery = $derived($page.url.searchParams.get('q') || '');
	let currentPage = $derived(parseInt($page.url.searchParams.get('page') || '1'));
	let limit = 50;
	let offset = $derived((currentPage - 1) * limit);

	// Query store
	let exercisesQuery = $derived(
		queryStore({
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
		})
	);

	// eslint-disable-next-line svelte/prefer-writable-derived
	let searchInput = $state($page.url.searchParams.get('q') || '');

	const updateUrl = (newParams: { page?: string }) => {
		const currentQueryValue = String(searchInput);
		const existingParams = Object.fromEntries($page.url.searchParams.entries());

		const updatedParams = {
			...existingParams,
			// Accesses the current reactive value of the input field
			q: currentQueryValue,
			...newParams
		};

		const query = new URLSearchParams(updatedParams);
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(resolve('/exercises') + `?${query.toString()}`);
	};

	const handleSearch = () => {
		updateUrl({ page: '1' });
	};

	const handlePageChange = (newPage: number) => {
		updateUrl({ page: newPage.toString() });
	};

	$effect(() => {
		searchInput = searchQuery;
	});
</script>

<div class="container mx-auto max-w-2xl p-4">
	<div class="mb-6 flex items-center justify-between">
		<Heading tag="h1">Exercises</Heading>
		<Button href="/exercises/new" color="blue">Create Exercise</Button>
	</div>

	<div class="mb-6 flex gap-2">
		<div class="flex-1">
			<Input
				bind:value={searchInput}
				placeholder="Search exercises..."
				onkeydown={(e) => e.key === 'Enter' && handleSearch()}
			/>
		</div>
		<Button onclick={handleSearch} color="dark">Search</Button>
	</div>

	{#if $exercisesQuery.fetching}
		<div class="py-8 text-center"><Spinner /></div>
	{:else if $exercisesQuery.error}
		<Alert color="red">Error: {$exercisesQuery.error.message}</Alert>
	{:else if $exercisesQuery.data?.uniqueExercises.length === 0}
		<div class="py-8 text-center text-gray-400">No exercises found.</div>
	{:else}
		<div class="space-y-2">
			{#each $exercisesQuery.data?.uniqueExercises || [] as exercise (exercise.id)}
				<Card class="w-full max-w-none flex-row items-center justify-between p-4">
					<div>
						<h3 class="text-lg font-semibold text-gray-900 dark:text-white">{exercise.name}</h3>
						{#if exercise.description}
							<p class="text-sm text-gray-500 dark:text-gray-400">{exercise.description}</p>
						{/if}
					</div>
					{#if exercise.isCustom}
						<Badge color="purple">Custom</Badge>
					{:else}
						<Badge color="gray">System</Badge>
					{/if}
				</Card>
			{/each}
		</div>

		<div class="mt-6 flex items-center justify-center gap-4">
			<Button
				disabled={currentPage <= 1}
				onclick={() => handlePageChange(currentPage - 1)}
				color="dark"
			>
				Previous
			</Button>
			<span class="py-2 text-gray-700 dark:text-gray-300">Page {currentPage}</span>
			<Button onclick={() => handlePageChange(currentPage + 1)} color="dark">Next</Button>
		</div>
	{/if}
</div>
