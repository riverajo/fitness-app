<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { setContextClient, gql, queryStore } from '@urql/svelte';
	import { client } from '$lib/client';
	import { Navbar, NavBrand, NavHamburger, NavUl, NavLi, DarkMode } from 'flowbite-svelte';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { initializeFaro, getWebInstrumentations } from '@grafana/faro-web-sdk';
	import { TracingInstrumentation } from '@grafana/faro-web-tracing';

	import { authStore } from '../stores/authStore';

	if (browser) {
		initializeFaro({
			url: '/faro/collect',
			app: {
				name: 'fitness-app',
				version: '0.0.1',
				environment: 'production'
			},
			instrumentations: [...getWebInstrumentations(), new TracingInstrumentation()]
		});
	}

	setContextClient(client);

	const ME_QUERY = gql`
		query Me {
			me {
				id
				email
			}
		}
	`;

	const meQuery = queryStore({
		client,
		query: ME_QUERY,
		pause: !browser || !$authStore.token
	});

	const logoutMutation = gql`
		mutation Logout {
			logout {
				success
				message
			}
		}
	`;

	async function handleLogout() {
		await client.mutation(logoutMutation, {}).toPromise();
		// Clear the token
		authStore.clearToken();
		// Force a full page reload to clear the urql cache and ensure a fresh state
		window.location.href = '/';
	}

	let { children } = $props();
	let activeUrl = $derived($page.url.pathname);

	const publicRoutes = ['/', '/register'];

	$effect(() => {
		if (browser) {
			// If we have no token, and we are not on a public route, redirect to /
			if (!$authStore.token && !publicRoutes.includes(activeUrl)) {
				goto(resolve('/'));
				return;
			}

			// If we DO have a token, we might wait for Me query
			if ($authStore.token && !$meQuery.fetching && !$meQuery.data?.me) {
				// If Me query failed or returned null (invalid session despite token), we might want to logout?
				// But let's trust the client.ts refreshAuth to handle invalid tokens.
				// This block might be redundant or tricky.
				// Let's just handle the case where we simply have no Me data yet.
			}
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Navbar
	fluid={true}
	class="border-b border-gray-200 bg-white px-4 py-2.5 sm:px-4 dark:border-gray-700 dark:bg-gray-800"
>
	<NavBrand href="/dashboard">
		<img src={favicon} class="mr-3 h-6 sm:h-9" alt="Fitness App Logo" />
		<span class="self-center text-xl font-semibold whitespace-nowrap dark:text-white"
			>Fitness App</span
		>
	</NavBrand>
	<div class="flex items-center md:order-2">
		<DarkMode />
		<NavHamburger class="ml-3" />
	</div>
	<NavUl {activeUrl}>
		{#if $meQuery.data?.me}
			<NavLi href="/dashboard">Dashboard</NavLi>
			<NavLi href="/exercises">Exercises</NavLi>
		{/if}
		{#if $meQuery.data?.me}
			<NavLi class="cursor-pointer" onclick={handleLogout}>Logout</NavLi>
		{/if}
	</NavUl>
</Navbar>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900">
	{#if $meQuery.fetching}
		<div class="flex h-screen items-center justify-center">
			<div
				class="h-8 w-8 animate-spin rounded-full border-b-2 border-gray-900 dark:border-white"
			></div>
		</div>
	{:else if publicRoutes.includes(activeUrl) || $meQuery.data?.me}
		{@render children()}
	{/if}
</div>
