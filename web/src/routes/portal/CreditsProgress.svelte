<script lang="ts">
	// Types
	interface CategoryProgress {
		name: string;
		credits: number;
		percentage: number;
	}

	interface SemesterRecord {
		name: string;
		credits: number;
		activitiesCount: number;
		grade: string;
		gradeColorClass: string;
	}

	interface GradeScale {
		name: string;
		range: string;
		isActive: boolean;
	}

	interface Insight {
		text: string;
		bgClass: string;
		textClass: string;
		borderClass: string;
	}

	// Mock data based on screenshot
	const categoryProgresses: CategoryProgress[] = [
		{ name: 'Technical Skills', credits: 48, percentage: 80 },
		{ name: 'Public Speaking', credits: 24, percentage: 60 },
		{ name: 'Research', credits: 18, percentage: 45 },
		{ name: 'Social Service', credits: 12, percentage: 30 },
		{ name: 'Sports', credits: 10, percentage: 25 },
		{ name: 'Cultural Activities', credits: 4, percentage: 10 },
		{ name: 'Leadership', credits: 2, percentage: 5 }
	];

	const semesterRecords: SemesterRecord[] = [
		{
			name: 'Semester 1',
			credits: 12,
			activitiesCount: 3,
			grade: 'Grade C',
			gradeColorClass: 'text-amber-500'
		},
		{
			name: 'Semester 2',
			credits: 18,
			activitiesCount: 5,
			grade: 'Grade B',
			gradeColorClass: 'text-blue-600'
		},
		{
			name: 'Semester 3',
			credits: 22,
			activitiesCount: 6,
			grade: 'Grade A',
			gradeColorClass: 'text-[#881B1B]'
		},
		{
			name: 'Semester 4',
			credits: 28,
			activitiesCount: 7,
			grade: 'Grade A',
			gradeColorClass: 'text-[#881B1B]'
		},
		{
			name: 'Semester 5',
			credits: 38,
			activitiesCount: 9,
			grade: 'Grade O',
			gradeColorClass: 'text-emerald-600'
		}
	];

	const gradeScales: GradeScale[] = [
		{ name: 'Grade O', range: '140+', isActive: false },
		{ name: 'Grade A', range: '100-139', isActive: true },
		{ name: 'Grade B', range: '70-99', isActive: false },
		{ name: 'Grade C', range: '40-69', isActive: false },
		{ name: 'Grade D', range: 'Below 40', isActive: false }
	];

	const insights: Insight[] = [
		{
			text: 'Your strongest contribution area is Technical Activities.',
			bgClass: 'bg-emerald-50/70',
			textClass: 'text-emerald-800',
			borderClass: 'border-emerald-150'
		},
		{
			text: 'You need 22 more credits to reach Grade O.',
			bgClass: 'bg-rose-50/70',
			textClass: 'text-rose-800',
			borderClass: 'border-rose-150'
		},
		{
			text: 'Research activities have the lowest contribution.',
			bgClass: 'bg-amber-50/60',
			textClass: 'text-amber-800',
			borderClass: 'border-amber-150'
		},
		{
			text: 'Participating in upcoming hackathons may improve your score faster.',
			bgClass: 'bg-blue-50/60',
			textClass: 'text-blue-800',
			borderClass: 'border-blue-150'
		}
	];
</script>

<div class="space-y-6">
	<!-- ==================== 1. OVERVIEW STAT BOX ==================== -->
	<section class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs space-y-6">
		<div class="flex flex-col md:flex-row md:items-start justify-between gap-6">
			<div>
				<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
					>Academic Year 2024-25</span
				>
				<div class="flex items-baseline gap-2 mt-1">
					<span class="text-5xl font-extrabold font-serif text-[#0B1535] leading-none">118</span>
					<span class="text-xs font-bold text-slate-400 uppercase tracking-widest"
						>Credits Earned</span
					>
				</div>
			</div>

			<!-- Grid of Sub-stats -->
			<div class="grid grid-cols-3 gap-6 md:gap-12 text-slate-800 shrink-0">
				<div class="flex flex-col">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider">Target</span>
					<span class="text-lg font-bold font-serif text-[#0B1535] mt-1">200</span>
					<span class="text-[8px] font-bold text-slate-400 uppercase tracking-wide"
						>Credits Required</span
					>
				</div>
				<div class="flex flex-col">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider"
						>Current Grade</span
					>
					<span class="text-lg font-bold font-serif text-[#881B1B] mt-1">A</span>
					<span class="text-[8px] font-bold text-slate-400 uppercase tracking-wide">&nbsp;</span>
				</div>
				<div class="flex flex-col">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider"
						>To Next Grade</span
					>
					<span class="text-lg font-bold font-serif text-[#0B1535] mt-1">22</span>
					<span class="text-[8px] font-bold text-slate-400 uppercase tracking-wide"
						>Credits Needed</span
					>
				</div>
			</div>
		</div>

		<!-- Progress bar block -->
		<div class="space-y-2">
			<div class="flex justify-between text-xs font-bold text-slate-700">
				<span>118 / 200 Credits Completed</span>
				<span>59%</span>
			</div>

			<div class="h-3 w-full bg-slate-100 rounded-full overflow-hidden relative">
				<div class="h-full bg-[#881B1B] rounded-full" style="width: 59%"></div>
			</div>

			<!-- Labels scale -->
			<div class="flex justify-between text-[10px] font-bold text-slate-400 font-sans px-1">
				<span>0</span>
				<span>50</span>
				<span>100</span>
				<span>150</span>
				<span>200</span>
			</div>
		</div>
	</section>

	<!-- ==================== 2. DOUBLE COLUMN CONTENT ==================== -->
	<section class="grid grid-cols-1 lg:grid-cols-12 gap-6">
		<!-- LEFT COLUMN -->
		<div class="lg:col-span-8 space-y-6">
			<!-- Credit Distribution Card -->
			<div class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs">
				<h2 class="text-base font-bold font-serif text-[#0B1535] pb-4 border-b border-slate-100">
					Credit Distribution by Activity Category
				</h2>

				<div class="space-y-4.5 mt-5">
					{#each categoryProgresses as cat}
						<div class="space-y-1.5">
							<div class="flex justify-between text-xs font-bold text-slate-800">
								<span class="text-[11px] font-bold text-slate-655">{cat.name}</span>
								<span>{cat.credits} credits</span>
							</div>
							<div class="h-2 w-full bg-slate-100 rounded-full overflow-hidden">
								<div
									class="h-full bg-[#0B1535] rounded-full"
									style="width: {cat.percentage}%"
								></div>
							</div>
						</div>
					{/each}
				</div>
			</div>

			<!-- Semester Progress Card -->
			<div class="bg-white border border-slate-200 p-5 rounded-xl shadow-xs">
				<h2 class="text-base font-bold font-serif text-[#0B1535] pb-4 border-b border-slate-100">
					Semester Progress
				</h2>

				<div class="overflow-x-auto mt-4">
					<table class="w-full text-left border-collapse text-xs">
						<thead>
							<tr
								class="text-[10px] font-bold text-[#6B7280] uppercase tracking-widest border-b border-slate-100 bg-slate-50/50"
							>
								<th class="py-3 px-4">Semester</th>
								<th class="py-3 px-4">Credits Earned</th>
								<th class="py-3 px-4">Activities</th>
								<th class="py-3 px-4 text-right">Grade Contribution</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100">
							{#each semesterRecords as sem}
								<tr class="hover:bg-slate-50/50 transition-colors">
									<td class="py-3.5 px-4 font-bold text-[#0B1535]">{sem.name}</td>
									<td class="py-3.5 px-4 text-slate-700 font-semibold">{sem.credits} credits</td>
									<td class="py-3.5 px-4 text-slate-500">{sem.activitiesCount} Activities</td>
									<td class="py-3.5 px-4 text-right font-bold {sem.gradeColorClass}">{sem.grade}</td
									>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		</div>

		<!-- RIGHT COLUMN -->
		<div class="lg:col-span-4 space-y-6">
			<!-- Current Grade Standing Card -->
			<div class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs">
				<h2 class="text-base font-bold font-serif text-[#0B1535] pb-4 border-b border-slate-100">
					Current Grade Standing
				</h2>

				<!-- Active Grade Banner Block -->
				<div
					class="mt-4 p-4.5 bg-[#881B1B]/5 border border-[#881B1B]/15 rounded-xl flex items-center gap-4"
				>
					<div
						class="w-12 h-12 bg-[#881B1B] text-white font-serif font-bold text-xl flex items-center justify-center rounded-lg shadow-sm"
					>
						A
					</div>
					<div class="flex flex-col">
						<span class="text-sm font-bold text-slate-900 leading-tight">Grade A</span>
						<span class="text-xs text-slate-400 font-medium block mt-0.5">100 - 139 Credits</span>
					</div>
				</div>

				<!-- Grade list table -->
				<div class="mt-5 space-y-1">
					<div
						class="grid grid-cols-12 text-[10px] font-bold text-slate-400 uppercase tracking-wider pb-2 px-2 border-b border-slate-100"
					>
						<span class="col-span-6">Grade</span>
						<span class="col-span-6 text-right">Credit Requirement</span>
					</div>

					<div class="divide-y divide-slate-50">
						{#each gradeScales as scale}
							<div
								class="grid grid-cols-12 items-center py-3 px-2 text-xs transition-colors duration-150
								{scale.isActive
									? 'bg-[#881B1B]/5 text-[#881B1B] font-extrabold border-l-2 border-[#881B1B] -ml-2 pl-2.5'
									: 'text-slate-700 font-semibold'}"
							>
								<div class="col-span-6 flex items-center gap-1.5">
									{#if scale.isActive}
										<span class="w-1.5 h-1.5 rounded-full bg-[#881B1B] shrink-0"></span>
									{/if}
									{scale.name}
								</div>
								<div class="col-span-6 text-right font-bold">
									{scale.range}
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>

			<!-- Performance Insights Card -->
			<div class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs space-y-4">
				<h2 class="text-base font-bold font-serif text-[#0B1535] pb-4 border-b border-slate-100">
					Performance Insights
				</h2>

				<div class="space-y-3 pt-1">
					{#each insights as insight}
						<div
							class="p-3.5 border rounded-lg {insight.bgClass} {insight.borderClass} flex items-start gap-2.5"
						>
							<div class="mt-0.5 shrink-0">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2.5"
									stroke="currentColor"
									class="w-3.5 h-3.5 {insight.textClass}"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M11.25 11.25l.041-.02a.75.75 0 111.063.852l-.708 2.836a.75.75 0 001.063.852l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z"
									/>
								</svg>
							</div>
							<p class="text-[11.5px] font-semibold leading-relaxed {insight.textClass}">
								{insight.text}
							</p>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</section>
</div>
