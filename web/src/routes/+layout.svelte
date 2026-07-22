<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { dev } from '$app/environment';
	import { injectAnalytics } from '@vercel/analytics/sveltekit';

	// Vercel Web Analytics: counts page views and visitors on the deployment.
	// A no-op outside Vercel, so local and Docker runs are unaffected.
	injectAnalytics({ mode: dev ? 'development' : 'production' });

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
{@render children()}
