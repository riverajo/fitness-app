<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { setContextClient, gql, queryStore } from '@urql/svelte';
	import { client } from '$lib/client';
	import { Navbar, NavBrand, NavHamburger, NavUl, NavLi, DarkMode } from 'flowbite-svelte';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { initializeFaro, getWebInstrumentations, LogLevel } from '@grafana/faro-web-sdk';
	import { TracingInstrumentation } from '@grafana/faro-web-tracing';

	if (browser) {
		initializeFaro({
			url: '/faro/collect',
			app: {
				name: 'fitness-app',
				version: '0.0.1',
				environment: 'production'
			},
			instrumentations: [
				...getWebInstrumentations(),
				new TracingInstrumentation(),
			],
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
		query: ME_QUERY
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
		// Force a full page reload to clear the urql cache and ensure a fresh state
		window.location.href = '/';
	}

	let { children } = $props();
	let activeUrl = $derived($page.url.pathname);
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
		<NavLi href="/dashboard">Dashboard</NavLi>
		<NavLi href="/exercises">Exercises</NavLi>
		{#if $meQuery.data?.me}
			<NavLi class="cursor-pointer" onclick={handleLogout}>Logout</NavLi>
		{/if}
	</NavUl>
</Navbar>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900">
	{@render children()}
</div>
