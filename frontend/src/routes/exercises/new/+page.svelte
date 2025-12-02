<script lang="ts">
	import { gql, getContextClient } from '@urql/svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { Heading, Label, Input, Textarea, Button, Alert } from 'flowbite-svelte';

	const client = getContextClient();

	let name = '';
	let description = '';
	let error = '';

	const createExerciseMutation = gql`
		mutation CreateUniqueExercise($input: CreateUniqueExerciseInput!) {
			createUniqueExercise(input: $input) {
				id
				name
				description
				isCustom
			}
		}
	`;

	async function handleSubmit() {
		error = '';
		const result = await client
			.mutation(createExerciseMutation, {
				input: {
					name,
					description
				}
			})
			.toPromise();

		if (result.error) {
			error = result.error.message;
		} else {
			await goto(resolve('/exercises'));
		}
	}
</script>

<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
	<Heading tag="h1" class="mb-6">Create New Exercise</Heading>

	<form onsubmit={handleSubmit} class="mt-6 max-w-xl space-y-6">
		{#if error}
			<Alert color="red">
				<span class="font-medium">Error creating exercise:</span>
				{error}
			</Alert>
		{/if}

		<div>
			<Label for="name" class="mb-2">Name</Label>
			<Input
				type="text"
				id="name"
				name="name"
				required
				bind:value={name}
				placeholder="e.g. Bench Press"
			/>
		</div>

		<div>
			<Label for="description" class="mb-2">Description</Label>
			<Textarea
				id="description"
				name="description"
				rows={3}
				bind:value={description}
				placeholder="Optional description..."
			/>
		</div>

		<div class="flex items-center justify-end gap-x-6">
			<Button color="light" href="/dashboard">Cancel</Button>
			<Button type="submit">Create Exercise</Button>
		</div>
	</form>
</div>
