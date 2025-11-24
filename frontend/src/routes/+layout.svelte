<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { setContextClient } from '@urql/svelte';
	import { client } from '$lib/client';
	import { Navbar, NavBrand, NavHamburger, NavUl, NavLi } from 'flowbite-svelte';
	import { page } from '$app/stores';

	setContextClient(client);

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Navbar fluid={true} class="border-b border-gray-200 bg-white px-4 py-2.5 dark:border-gray-700 dark:bg-gray-800 sm:px-4">
	<NavBrand href="/dashboard">
		<img src={favicon} class="mr-3 h-6 sm:h-9" alt="Fitness App Logo" />
		<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Fitness App</span>
	</NavBrand>
	<NavHamburger />
	<NavUl>
		<NavLi href="/dashboard" active={$page.url.pathname === '/dashboard'}>Dashboard</NavLi>
		<NavLi href="/exercises" active={$page.url.pathname.startsWith('/exercises')}>Exercises</NavLi>
	</NavUl>
</Navbar>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900">
	{@render children()}
</div>
