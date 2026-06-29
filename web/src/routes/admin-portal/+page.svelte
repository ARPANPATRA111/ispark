<script lang="ts">
	import HeaderSection from './HeaderSection.svelte';
	import AdminPortalSection from './AdminPortalSection.svelte';
	import FooterSection from './FooterSection.svelte';
</script>

<svelte:head>
	<title>Administrative Control Portal | iSPARC IIPS DAVV</title>
	<meta
		name="description"
		content="Secure portal access for iSPARC administrative users and coordinators at IIPS DAVV, Indore."
	/>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="min-h-screen bg-bg-main flex flex-col items-center w-full">
	<HeaderSection />
	<main class="flex-grow w-full flex justify-center items-start">
		<AdminPortalSection />
	</main>
	<FooterSection />
</div>
