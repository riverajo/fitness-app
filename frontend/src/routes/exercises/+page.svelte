<script lang="ts">
	import { getContextClient } from '@urql/svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { Heading, Button, Input, Card, Badge, Spinner } from 'flowbite-svelte';
	import { exerciseStore } from '../../state/exercise.svelte';

	const client = getContextClient();

	// Reactive parameters from URL
	let searchQuery = $derived($page.url.searchParams.get('q') || '');
	let currentPage = $derived(parseInt($page.url.searchParams.get('page') || '1'));
	let limit = 50;

	// Use the store for data
	let exercises = $derived(exerciseStore.search(searchQuery));
	let paginatedExercises = $derived(
		exercises.slice((currentPage - 1) * limit, currentPage * limit)
	);

	// Sync on mount
	$effect(() => {
		if (client) {
			exerciseStore.sync(client);
		}
	});

	// eslint-disable-next-line svelte/prefer-writable-derived
	let searchInput = $state($page.url.searchParams.get('q') || '');

	const updateUrl = (newParams: { page?: string; q?: string }) => {
		const existingParams = Object.fromEntries($page.url.searchParams.entries());

		const updatedParams = {
			...existingParams,
			...newParams
		};

		const query = new URLSearchParams(updatedParams);
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(resolve('/exercises') + `?${query.toString()}`, { replaceState: true });
	};

	const handleSearch = () => {
		updateUrl({ page: '1', q: searchInput });
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

	<div class="mb-6">
		<Input bind:value={searchInput} placeholder="Search exercises..." oninput={handleSearch} />
	</div>

	{#if exerciseStore.loading && !exerciseStore.all.length}
		<div class="py-8 text-center"><Spinner /></div>
	{:else if exercises.length === 0}
		<div class="py-8 text-center text-gray-400">No exercises found.</div>
	{:else}
		<div class="space-y-2">
			{#each paginatedExercises as exercise (exercise.id)}
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
			<Button
				disabled={paginatedExercises.length < limit && currentPage * limit >= exercises.length}
				onclick={() => handlePageChange(currentPage + 1)}
				color="dark"
			>
				Next
			</Button>
		</div>
	{/if}
</div>
