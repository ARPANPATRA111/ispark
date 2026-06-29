<script lang="ts">
	// Props configuration
	let {
		badgeText = 'Departmental Affiliated Cell',
		titleText = 'Skill, Personality Advancement & Refinement Cell',
		subtitleText = '',
		descriptionText = 'Welcome to ISPARC, an official cell of the Institute of Infrastructure and Professional Studies (IIPS), DAVV. We facilitate structured extracurricular activities, personality enhancement tracks, and grading systems to foster student growth and leadership.',
		primaryCtaText = 'Grading Matrix',
		primaryCtaHref = '#grading-scheme',
		secondaryCtaText = 'Explore Tracks',
		secondaryCtaHref = '#tracks'
	} = $props();

	// Grading matrix data for the right-hand panel
	const grades = [
		{
			range: '60+ Points',
			grade: 'Grade O (Outstanding)',
			colorClass: 'bg-emerald-50 text-emerald-800 border-emerald-200'
		},
		{
			range: '45 - 59 Points',
			grade: 'Grade A (Excellent)',
			colorClass: 'bg-blue-50 text-blue-800 border-blue-200'
		},
		{
			range: '30 - 44 Points',
			grade: 'Grade B (Very Good)',
			colorClass: 'bg-amber-50 text-amber-850 border-amber-200'
		},
		{
			range: '15 - 29 Points',
			grade: 'Grade C (Good)',
			colorClass: 'bg-purple-50 text-purple-800 border-purple-200'
		},
		{
			range: 'Under 15 Points',
			grade: 'Grade D (Satisfactory)',
			colorClass: 'bg-slate-50 text-slate-800 border-slate-200'
		}
	];
</script>

<section class="relative py-16 sm:py-24 bg-white">
	<div class="max-w-7xl mx-auto px-6 lg:px-8">
		<div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-start">
			<!-- Left Column: Headline and Context -->
			<div class="lg:col-span-7 space-y-6">
				<!-- Academic Affiliation Badge -->
				<div
					class="inline-flex items-center px-3 py-1 bg-[#284B63] rounded-sm text-[10px] font-bold tracking-widest uppercase text-white"
				>
					{badgeText}
				</div>

				<!-- Main Institutional Headline -->
				<h1
					class="text-4xl sm:text-5xl lg:text-6xl font-bold text-[#881B1B] font-serif tracking-tight leading-none"
				>
					{titleText}
				</h1>

				{#if subtitleText}
					<!-- Subtitle -->
					<h2 class="text-lg sm:text-xl font-semibold text-[#1e3a5f] font-sans tracking-tight">
						{subtitleText}
					</h2>
				{/if}

				<!-- Subtext explanations -->
				<p class="text-sm text-slate-650 leading-relaxed font-sans font-normal max-w-2xl">
					{descriptionText}
				</p>

				<!-- Buttons -->
				<div class="flex flex-col sm:flex-row gap-3 pt-2">
					<a
						href={primaryCtaHref}
						class="inline-flex items-center justify-center border-2 border-[#284B63] text-[#284B63] font-semibold text-xs tracking-wider uppercase px-6 py-3 hover:bg-[#284B63] hover:text-white transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#284B63] focus-visible:ring-offset-2"
					>
						{primaryCtaText}
					</a>
					<a
						href={secondaryCtaHref}
						class="inline-flex items-center justify-center border-2 border-[#284B63] text-[#284B63] font-semibold text-xs tracking-wider uppercase px-6 py-3 hover:bg-[#284B63] hover:text-white transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#284B63] focus-visible:ring-offset-2"
					>
						{secondaryCtaText}
					</a>
				</div>
			</div>

			<!-- Right Column: Activity Grading Matrix (Clean Grid Panel) -->
			<div class="lg:col-span-5 bg-white border border-slate-200 p-6 rounded-sm">
				<h2 class="text-lg font-bold font-serif text-[#1e3a5f] mb-3">Activity Grading Matrix</h2>
				<p
					class="text-xs text-slate-500 font-sans leading-relaxed mb-6 border-b border-slate-100 pb-4"
				>
					Extracurricular points are mapped directly to grades in semester results under the DAVV
					guidelines.
				</p>

				<div class="space-y-3.5">
					{#each grades as item}
						<div
							class="flex items-center justify-between py-2.5 border-b border-slate-100 last:border-b-0"
						>
							<span class="text-sm font-semibold text-slate-700 font-sans">{item.range}</span>
							<span
								class="inline-flex items-center px-3 py-1 rounded-sm text-xs font-bold border {item.colorClass} tracking-wide font-sans"
							>
								{item.grade}
							</span>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
</section>
