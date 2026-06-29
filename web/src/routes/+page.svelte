<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import Hero from '$lib/components/Hero.svelte';
	import About from '$lib/components/About.svelte';
	import Tracks from '$lib/components/Tracks.svelte';
	import Timeline from '$lib/components/Timeline.svelte';
	import Outcomes from '$lib/components/Outcomes.svelte';
	import Footer from '$lib/components/Footer.svelte';
</script>

<svelte:head>
	<title>ISPARC | IIPS Skill, Personality Advancement & Refinement Cell</title>
	<meta
		name="description"
		content="Official student development and academic portal of the Institute of Infrastructure and Professional Studies (IIPS), DAVV, Indore."
	/>
</svelte:head>

<div
	class="min-h-screen bg-bg-warm flex flex-col font-sans selection:bg-accent-gold/25 selection:text-primary-dark"
>
	<!-- Main Navigation Bar -->
	<Navbar />

	<main class="flex-grow">
		<!-- Hero Section & Activity Grading Matrix -->
		<Hero />

		<!-- Foundation Pillars (About Section) -->
		<About />

		<!-- Academic Curriculum Tracks -->
		<Tracks />

		<!-- Progression Event Journey Timeline -->
		<Timeline />

		<!-- Accreditation, Quote & Student Registry Portal Banner -->
		<Outcomes />
	</main>

	<!-- Footer -->
	<Footer />
</div>
