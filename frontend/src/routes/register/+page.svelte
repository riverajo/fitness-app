<script lang="ts">
	import { gql, getContextClient } from '@urql/svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { Card, Button, Label, Input, Alert } from 'flowbite-svelte';
	import { authStore } from '../../stores/authStore';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	const registerMutation = gql`
		mutation Register($input: RegisterInput!) {
			register(input: $input) {
				success
				message
				user {
					id
					email
				}
				token
			}
		}
	`;

	import { onMount } from 'svelte';

	const client = getContextClient();

	onMount(async () => {
		const token = authStore.getToken();
		if (token) {
			await goto(resolve('/dashboard'));
		}
	});

	async function handleSubmit(e: Event) {
		console.log('Submitting form...');
		e.preventDefault();
		loading = true;
		error = '';

		try {
			const result = await client
				.mutation(registerMutation, { input: { email, password } })
				.toPromise();

			if (result.error) {
				error = result.error.message;
			} else if (result.data?.register?.success) {
				// Store the token
				if (result.data.register.token) {
					authStore.setToken(result.data.register.token);
				}
				// Refresh the Me query to ensure the layout and dashboard are updated
				// Force a hard navigation to reset the URQL client and pick up the new token
				window.location.href = resolve('/dashboard');
			} else {
				error = result.data?.register?.message || 'Registration failed';
			}
		} catch (e) {
			error = 'An unexpected error occurred';
			console.error(e);
		} finally {
			loading = false;
		}
	}
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8 dark:bg-gray-900"
>
	<Card class="w-full max-w-md">
		<h2 class="text-center text-2xl font-bold text-gray-900 dark:text-white">
			Create your account
		</h2>
		<form class="mt-8 space-y-6" onsubmit={handleSubmit}>
			<div class="space-y-4">
				<div>
					<Label for="email" class="mb-2">Email address</Label>
					<Input
						id="email"
						name="email"
						type="email"
						placeholder="name@company.com"
						required
						bind:value={email}
					/>
				</div>
				<div>
					<Label for="password" class="mb-2">Password</Label>
					<Input
						id="password"
						name="password"
						type="password"
						placeholder="••••••••"
						required
						bind:value={password}
					/>
				</div>
			</div>

			{#if error}
				<Alert color="red" class="mt-4">
					<span class="font-medium">Error!</span>
					{error}
				</Alert>
			{/if}

			<Button type="submit" class="w-full" disabled={loading}>
				{#if loading}Signing up...{:else}Sign up{/if}
			</Button>
		</form>
	</Card>
</div>
