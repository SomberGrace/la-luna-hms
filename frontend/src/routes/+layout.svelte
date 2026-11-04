<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { APP_ROLE, isPathAllowedForApp, ROLE_ROUTE_MAP } from '#lib/mobile';
	import './layout.css';
	import './fonts.css';
	import '#lib/assets/phosphor-icons/regular/style.css';
	import '#lib/assets/phosphor-icons/light/style.css';
	import ThemeToggle from '#lib/components/ThemeToggle.svelte';
	import favicon from '#lib/assets/favicon.svg';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();
	let guestAppHome = $derived(APP_ROLE === 'guest' && page.url.pathname === '/');
	let routeAllowed = $derived(isPathAllowedForApp(page.url.pathname) && !guestAppHome);

	$effect(() => {
		if (!routeAllowed && APP_ROLE !== 'shared') {
			void goto(ROLE_ROUTE_MAP[APP_ROLE], { replaceState: true });
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{#if routeAllowed}
	{@render children()}
{/if}
<ThemeToggle />
