<script lang="ts">
	import { gql, getContextClient } from '@urql/svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { Card, Button, Label, Input, Alert } from 'flowbite-svelte';
	import { authStore } from '../stores/authStore';

	const client = getContextClient();

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	const loginMutation = gql`
		mutation Login($input: LoginInput!) {
			login(input: $input) {
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

	onMount(async () => {
		// Check if we have a token in the store
		const token = authStore.getToken();
		if (token) {
			await goto(resolve('/dashboard'));
		}
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';

		try {
			const result = await client
				.mutation(loginMutation, { input: { email, password } })
				.toPromise();

			if (result.error) {
				error = result.error.message;
			} else if (result.data?.login?.success) {
				// Store the token
				if (result.data.login.token) {
					authStore.setToken(result.data.login.token);
				}
				// Refresh the Me query to ensure the layout and dashboard are updated
				// Force a hard navigation to reset the URQL client and pick up the new token
				window.location.href = resolve('/dashboard');
			} else {
				error = result.data?.login?.message || 'Login failed';
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
			Sign in to your account
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
				{#if loading}Signing in...{:else}Sign in{/if}
			</Button>

			<div class="text-center text-sm font-medium text-gray-500 dark:text-gray-300">
				Don't have an account? <a
					href={resolve('/register')}
					class="text-primary-700 dark:text-primary-500 hover:underline">Register</a
				>
			</div>
		</form>
	</Card>
</div>
