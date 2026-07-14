<script lang="ts">
	import { API_BASE_URL } from '$lib/config';

	// Types
	interface CategoryProgress {
		name: string;
		credits: number;
		percentage: number;
	}

	interface SemesterRecord {
		name: string;
		semester: string;
		credits: number;
		activitiesCount: number;
		grade: string;
	}

	interface GradeScale {
		name: string;
		range: string;
		isActive: boolean;
	}

	interface Insight {
		text: string;
		type: 'success' | 'danger' | 'warning' | 'info';
	}

	// States
	let categoryProgresses = $state<CategoryProgress[]>([]);
	let semesterRecords = $state<SemesterRecord[]>([]);
	let gradeScales = $state<GradeScale[]>([]);
	let insights = $state<Insight[]>([]);

	let totalCredits = $state(0);
	let targetCredits = $state(200);
	let finalGrade = $state('...');
	let neededToNext = $state(0);
	let academicYear = $state('2025-26');

	let progressPercentage = $derived(
		Math.min(Math.round((totalCredits / targetCredits) * 100), 100)
	);

	let activeScale = $derived(
		gradeScales.find((s) => s.isActive) || { name: 'Grade D', range: 'Below 40' }
	);

	const insightStyles = {
		success: {
			bg: 'bg-emerald-50/70',
			border: 'border-emerald-200',
			text: 'text-emerald-800'
		},
		danger: {
			bg: 'bg-rose-50/70',
			border: 'border-rose-200',
			text: 'text-rose-800'
		},
		warning: {
			bg: 'bg-amber-50/60',
			border: 'border-amber-200',
			text: 'text-amber-800'
		},
		info: {
			bg: 'bg-blue-50/60',
			border: 'border-blue-200',
			text: 'text-blue-800'
		}
	};

	function getGradeColorClass(grade: string): string {
		switch (grade) {
			case 'Grade O':
				return 'text-emerald-600';
			case 'Grade A':
				return 'text-[#881B1B]';
			case 'Grade B':
				return 'text-blue-600';
			case 'Grade C':
				return 'text-amber-500';
			default:
				return 'text-slate-500';
		}
	}

	async function loadProgressData() {
		try {
			const token = localStorage.getItem('access_token') || '';
			const res = await fetch(`${API_BASE_URL}/api/student/marksheet`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
			if (res.ok) {
				const data = await res.json();
				totalCredits = data.total_credits;
				targetCredits = data.target_credits || 200;
				finalGrade = data.final_grade;
				neededToNext = data.credits_needed_to_next_grade;
				gradeScales = data.grade_scales || [];
				insights = data.insights || [];
				academicYear = data.student_info.academicYear || '2025-26';

				// Map category progress
				categoryProgresses = (data.credit_categories || []).map(
					(cat: { category: string; credits: number; contribution: string }) => {
						const pct = parseInt(cat.contribution) || 0;
						return {
							name: cat.category,
							credits: cat.credits,
							percentage: pct
						};
					}
				);

				// Map semester records
				semesterRecords = data.semester_summary || [];
			}
		} catch (err) {
			console.error(err);
		}
	}

	$effect(() => {
		loadProgressData();
	});
</script>

<div class="space-y-6">
	<!-- ==================== 1. OVERVIEW STAT BOX ==================== -->
	<section class="bg-white border border-slate-200 p-6 rounded-xl shadow-xs space-y-6">
		<div class="flex flex-col md:flex-row md:items-start justify-between gap-6">
			<div>
				<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
					>Academic Year {academicYear}</span
				>
				<div class="flex items-baseline gap-2 mt-1">
					<span class="text-5xl font-extrabold font-serif text-[#0B1535] leading-none"
						>{totalCredits}</span
					>
					<span class="text-xs font-bold text-slate-400 uppercase tracking-widest"
						>Credits Earned</span
					>
				</div>
			</div>

			<!-- Grid of Sub-stats -->
			<div class="grid grid-cols-3 gap-6 md:gap-12 text-slate-800 shrink-0">
				<div class="flex flex-col">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider">Target</span>
					<span class="text-lg font-bold font-serif text-[#0B1535] mt-1">{targetCredits}</span>
					<span class="text-[8px] font-bold text-slate-400 uppercase tracking-wide"
						>Credits Required</span
					>
				</div>
				<div class="flex flex-col">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider"
						>Current Grade</span
					>
					<span class="text-lg font-bold font-serif text-[#881B1B] mt-1"
						>{finalGrade.replace('Grade ', '')}</span
					>
					<span class="text-[8px] font-bold text-slate-400 uppercase tracking-wide">&nbsp;</span>
				</div>
				<div class="flex flex-col">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider"
						>To Next Grade</span
					>
					<span class="text-lg font-bold font-serif text-[#0B1535] mt-1">{neededToNext}</span>
					<span class="text-[8px] font-bold text-slate-400 uppercase tracking-wide"
						>Credits Needed</span
					>
				</div>
			</div>
		</div>

		<!-- Progress bar block -->
		<div class="space-y-2">
			<div class="flex justify-between text-xs font-bold text-slate-700">
				<span>{totalCredits} / {targetCredits} Credits Completed</span>
				<span>{progressPercentage}%</span>
			</div>

			<div class="h-3 w-full bg-slate-100 rounded-full overflow-hidden relative">
				<div class="h-full bg-[#881B1B] rounded-full" style="width: {progressPercentage}%"></div>
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
									<td class="py-3.5 px-4 font-bold text-[#0B1535]">{sem.semester}</td>
									<td class="py-3.5 px-4 text-slate-700 font-semibold">{sem.credits} credits</td>
									<td class="py-3.5 px-4 text-slate-500">{sem.activitiesCount} Activities</td>
									<td class="py-3.5 px-4 text-right font-bold {getGradeColorClass(sem.grade)}"
										>{sem.grade}</td
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
						{activeScale.name.replace('Grade ', '')}
					</div>
					<div class="flex flex-col">
						<span class="text-sm font-bold text-slate-900 leading-tight">{activeScale.name}</span>
						<span class="text-xs text-slate-400 font-medium block mt-0.5"
							>{activeScale.range} Credits</span
						>
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
						{@const style = insightStyles[insight.type] || insightStyles.info}
						<div class="p-3.5 border rounded-lg {style.bg} {style.border} flex items-start gap-2.5">
							<div class="mt-0.5 shrink-0">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2.5"
									stroke="currentColor"
									class="w-3.5 h-3.5 {style.text}"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M11.25 11.25l.041-.02a.75.75 0 111.063.852l-.708 2.836a.75.75 0 001.063.852l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z"
									/>
								</svg>
							</div>
							<p class="text-[11.5px] font-semibold leading-relaxed {style.text}">
								{insight.text}
							</p>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</section>
</div>
