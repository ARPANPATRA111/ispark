<script lang="ts">
	// Optional callback so quick-actions inside this view can switch tabs,
	// same pattern as AdminDashboardView's onTabChange prop.
	let { onTabChange = (tab: string) => {} }: { onTabChange?: (tab: string) => void } = $props();

	// ---- Mock data (swap for real API calls) ----
	const summary = {
		assignedBatches: 3,
		totalStudents: 342,
		compliantStudents: 278,
		defaulters: 64
	};

	const batches = [
		{
			name: 'IT 2025',
			students: 120,
			pd: 218,
			sb: 246,
			compliance: 92,
			defaulters: 8,
			pendingCerts: 4,
			status: 'Excellent'
		},
		{
			name: 'IT 2024',
			students: 118,
			pd: 192,
			sb: 201,
			compliance: 84,
			defaulters: 19,
			pendingCerts: 9,
			status: 'Good'
		},
		{
			name: 'MBA 2025',
			students: 104,
			pd: 143,
			sb: 129,
			compliance: 68,
			defaulters: 37,
			pendingCerts: 13,
			status: 'At Risk'
		}
	];

	const alerts = [
		{
			tone: 'rose',
			title: 'Defaulters',
			description: 'Students who have not met compliance requirements',
			value: '64 Students',
			icon: 'M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z'
		},
		{
			tone: 'amber',
			title: 'Inactive Students',
			description: 'No activity recorded in the current semester',
			value: '18 Students',
			icon: 'M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z'
		},
		{
			tone: 'blue',
			title: 'Pending Certificate Reviews',
			description: 'Awaiting mentor approval and verification',
			value: '25 Submissions',
			icon: 'M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z'
		}
	];

	const requirements = [
		{ tone: 'emerald', label: 'Completed Both Tracks', count: 278 },
		{ tone: 'amber', label: 'Only Personality Development Completed', count: 24 },
		{ tone: 'amber', label: 'Only Skill Building Completed', count: 18 },
		{ tone: 'rose', label: 'Neither Track Completed', count: 22 }
	];

	const reports = [
		{
			title: 'Batch Report',
			icon: 'M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z'
		},
		{
			title: 'Student Progress Report',
			icon: 'M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.109A11.386 11.386 0 0110.089 21c-2.316 0-4.445-.69-6.22-1.879v-.003a4.125 4.125 0 017.533-2.493M15 19.128v-.003c0-1.112-.285-2.16-.786-3.07M14.214 16.058A9.396 9.396 0 0010.089 15c-1.47 0-2.854.34-4.082.945M14.214 16.058a9.386 9.386 0 010 3.07M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5'
		},
		{
			title: 'Compliance Report',
			icon: 'M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z'
		},
		{
			title: 'Certificate Verification Report',
			icon: 'M4.26 10.147a60.436 60.436 0 00-.491 6.347A48.627 48.627 0 0112 20.904a48.627 48.627 0 018.232-4.41 60.46 60.46 0 00-.491-6.347m-15.482 0a50.57 50.57 0 00-2.658-.813A59.905 59.905 0 0112 3.493a59.902 59.902 0 0110.399 5.84c-.896.248-1.783.52-2.658.814m-15.482 0A50.697 50.697 0 0112 13.489a50.702 50.702 0 017.74-3.342M6.75 15a.75.75 0 100-1.5.75.75 0 000 1.5zm0 0v-3.675A55.378 55.378 0 0112 8.443m-7.007 11.55A5.981 5.981 0 006.75 15.75v-1.5'
		}
	];

	let searchQuery = $state('');
	let filteredBatches = $derived(
		batches.filter((b) => b.name.toLowerCase().includes(searchQuery.trim().toLowerCase()))
	);

	function statusStyles(status: string) {
		if (status === 'Excellent')
			return {
				pill: 'bg-emerald-50 text-emerald-700',
				dot: 'bg-emerald-500',
				bar: 'bg-emerald-500'
			};
		if (status === 'Good')
			return { pill: 'bg-blue-50 text-blue-700', dot: 'bg-blue-500', bar: 'bg-blue-500' };
		return { pill: 'bg-amber-50 text-amber-700', dot: 'bg-amber-500', bar: 'bg-amber-500' };
	}

	function toneClasses(tone: string) {
		const map: Record<string, { iconBg: string; iconText: string; valueText: string }> = {
			rose: { iconBg: 'bg-rose-50', iconText: 'text-[#881B1B]', valueText: 'text-[#881B1B]' },
			amber: { iconBg: 'bg-amber-50', iconText: 'text-amber-600', valueText: 'text-amber-600' },
			blue: { iconBg: 'bg-blue-50', iconText: 'text-blue-600', valueText: 'text-blue-600' },
			emerald: {
				iconBg: 'bg-emerald-50',
				iconText: 'text-emerald-600',
				valueText: 'text-emerald-600'
			}
		};
		return map[tone] ?? map.blue;
	}
</script>

<div class="space-y-6">
	<!-- ================= Summary cards ================= -->
	<div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
		<div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5">
			<div class="w-10 h-10 rounded-lg bg-[#881B1B]/10 flex items-center justify-center">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
					class="w-5 h-5 text-[#881B1B]"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M20.25 14.15v4.25c0 1.094-.787 2.036-1.872 2.18-2.087.277-4.216.42-6.378.42s-4.291-.143-6.378-.42c-1.085-.144-1.872-1.086-1.872-2.18v-4.25m16.5 0a2.18 2.18 0 00.75-1.661V8.706c0-1.081-.768-2.015-1.837-2.175a48.114 48.114 0 00-3.413-.387m4.5 8.006c-.194.165-.42.295-.673.38A23.978 23.978 0 0112 15.75c-2.648 0-5.195-.429-7.577-1.22a2.016 2.016 0 01-.673-.38m0 0A2.18 2.18 0 013 12.489V8.706c0-1.081.768-2.015 1.837-2.175a48.111 48.111 0 013.413-.387m7.5 0V5.25A2.25 2.25 0 0013.5 3h-3a2.25 2.25 0 00-2.25 2.25v.894m7.5 0a48.667 48.667 0 00-7.5 0"
					/>
				</svg>
			</div>
			<p class="mt-3 text-2xl font-bold font-serif text-slate-900">{summary.assignedBatches}</p>
			<p class="mt-1 text-xs font-bold text-slate-700">Assigned Batches</p>
		</div>

		<div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5">
			<div class="w-10 h-10 rounded-lg bg-blue-50 flex items-center justify-center">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
					class="w-5 h-5 text-blue-600"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.109A11.386 11.386 0 0110.089 21c-2.316 0-4.445-.69-6.22-1.879v-.003a4.125 4.125 0 017.533-2.493M15 19.128v-.003c0-1.112-.285-2.16-.786-3.07M14.214 16.058A9.396 9.396 0 0010.089 15c-1.47 0-2.854.34-4.082.945M14.214 16.058a9.386 9.386 0 010 3.07M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
					/>
				</svg>
			</div>
			<p class="mt-3 text-2xl font-bold font-serif text-slate-900">{summary.totalStudents}</p>
			<p class="mt-1 text-xs font-bold text-slate-700">Total Students</p>
			<p class="text-xs text-slate-400">Across assigned batches</p>
		</div>

		<div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5">
			<div class="w-10 h-10 rounded-lg bg-emerald-50 flex items-center justify-center">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
					class="w-5 h-5 text-emerald-600"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
			</div>
			<p class="mt-3 text-2xl font-bold font-serif text-slate-900">{summary.compliantStudents}</p>
			<p class="mt-1 text-xs font-bold text-slate-700">Compliant Students</p>
			<p class="text-xs text-slate-400">Completed semester requirements</p>
		</div>

		<div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5">
			<div class="w-10 h-10 rounded-lg bg-rose-50 flex items-center justify-center">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
					class="w-5 h-5 text-[#881B1B]"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"
					/>
				</svg>
			</div>
			<p class="mt-3 text-2xl font-bold font-serif text-[#881B1B]">{summary.defaulters}</p>
			<p class="mt-1 text-xs font-bold text-slate-700">Defaulters</p>
			<p class="text-xs text-slate-400">Require mentor intervention</p>
		</div>
	</div>

	<!-- ================= Batch performance table ================= -->
	<div class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
		<div
			class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 px-5 py-4 border-b border-slate-200"
		>
			<div>
				<h2 class="font-serif font-bold text-sm text-slate-900">Batch Performance Overview</h2>
				<p class="text-xs text-slate-400 mt-0.5">Performance tracking for assigned batches</p>
			</div>
			<div
				class="flex items-center gap-2 bg-slate-50 border border-slate-200 rounded-lg px-3 py-1.5 w-full sm:w-56"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
					class="w-4 h-4 text-slate-400 shrink-0"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
					/>
				</svg>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search batch name..."
					class="bg-transparent text-xs text-slate-800 placeholder-slate-400 focus:outline-none w-full"
				/>
			</div>
		</div>

		<div class="overflow-x-auto">
			<table class="w-full min-w-[900px] text-left">
				<thead>
					<tr class="border-b border-slate-200">
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>Batch Name</th
						>
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>Students</th
						>
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>PD Credits</th
						>
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>SB Credits</th
						>
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>Compliance %</th
						>
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>Defaulters</th
						>
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>Pending Certs</th
						>
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>Status</th
						>
						<th class="px-5 py-3 text-[11px] font-bold text-slate-400 uppercase tracking-wide"
							>Actions</th
						>
					</tr>
				</thead>
				<tbody>
					{#each filteredBatches as batch}
						{@const styles = statusStyles(batch.status)}
						<tr class="border-b border-slate-100 last:border-b-0 hover:bg-slate-50/60">
							<td class="px-5 py-3.5 text-sm font-bold text-slate-900">{batch.name}</td>
							<td class="px-5 py-3.5 text-sm text-slate-700">{batch.students}</td>
							<td class="px-5 py-3.5 text-sm text-slate-700">{batch.pd}</td>
							<td class="px-5 py-3.5 text-sm text-slate-700">{batch.sb}</td>
							<td class="px-5 py-3.5">
								<div class="flex items-center gap-2">
									<div class="w-20 h-1.5 rounded-full bg-slate-100 overflow-hidden">
										<div
											class="h-full rounded-full {styles.bar}"
											style="width:{batch.compliance}%"
										></div>
									</div>
									<span class="text-xs font-semibold text-slate-800">{batch.compliance}%</span>
								</div>
							</td>
							<td
								class="px-5 py-3.5 text-sm font-bold {batch.defaulters > 20
									? 'text-[#881B1B]'
									: 'text-slate-900'}"
							>
								{batch.defaulters}
							</td>
							<td class="px-5 py-3.5 text-sm text-slate-500">{batch.pendingCerts}</td>
							<td class="px-5 py-3.5">
								<span
									class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold {styles.pill}"
								>
									<span class="w-1.5 h-1.5 rounded-full {styles.dot}"></span>
									{batch.status}
								</span>
							</td>
							<td class="px-5 py-3.5">
								<div class="flex items-center gap-1.5">
									<button
										type="button"
										aria-label="View {batch.name}"
										class="w-7 h-7 rounded-md bg-[#881B1B]/10 text-[#881B1B] flex items-center justify-center hover:bg-[#881B1B]/20 transition-colors"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke="currentColor"
											stroke-width="2"
											class="w-3.5 h-3.5"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
											/>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
											/>
										</svg>
									</button>
									<button
										type="button"
										aria-label="Edit {batch.name}"
										class="w-7 h-7 rounded-md bg-slate-900 text-white flex items-center justify-center hover:bg-slate-800 transition-colors"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke="currentColor"
											stroke-width="2"
											class="w-3.5 h-3.5"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="m16.862 4.487 1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125"
											/>
										</svg>
									</button>
								</div>
							</td>
						</tr>
					{:else}
						<tr>
							<td colspan="9" class="px-5 py-8 text-center text-sm text-slate-400">
								No batches match &ldquo;{searchQuery}&rdquo;.
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>

	<!-- ================= Mentor alerts ================= -->
	<div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5">
		<h3 class="text-xs font-bold uppercase tracking-wider text-[#881B1B] font-serif">
			Mentor Alerts
		</h3>
		<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
			{#each alerts as alert}
				{@const t = toneClasses(alert.tone)}
				<div class="flex gap-3 p-3 rounded-xl bg-slate-50 border border-slate-100">
					<div class="w-8 h-8 rounded-md {t.iconBg} flex items-center justify-center shrink-0">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
							class="w-4 h-4 {t.iconText}"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d={alert.icon} />
						</svg>
					</div>
					<div>
						<p class="text-xs font-bold text-slate-900">{alert.title}</p>
						<p class="text-xs text-slate-500 mt-0.5">{alert.description}</p>
						<p class="text-xs font-bold mt-1 {t.valueText}">{alert.value}</p>
					</div>
				</div>
			{/each}
		</div>
	</div>

	<!-- ================= Semester requirement status ================= -->
	<div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5">
		<h2 class="font-serif font-bold text-sm text-slate-900">Semester Requirement Status</h2>
		<p class="text-xs text-slate-400 mt-0.5">
			Track student completion across both development tracks
		</p>

		<div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4 mt-4">
			{#each requirements as req}
				{@const tone =
					req.tone === 'emerald'
						? { bg: 'bg-emerald-50', border: 'border-emerald-200', text: 'text-emerald-600' }
						: req.tone === 'amber'
							? { bg: 'bg-amber-50', border: 'border-amber-200', text: 'text-amber-600' }
							: { bg: 'bg-rose-50', border: 'border-rose-200', text: 'text-[#881B1B]' }}
				<div class="rounded-xl border {tone.border} {tone.bg} p-4">
					<p class="text-xs text-slate-600 leading-snug">{req.label}</p>
					<p class="mt-3 text-2xl font-bold font-serif {tone.text}">{req.count}</p>
					<p class="text-xs text-slate-400">Students</p>
				</div>
			{/each}
		</div>
	</div>

	<!-- ================= Reports ================= -->
	<div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5">
		<h2 class="font-serif font-bold text-sm text-slate-900">Reports</h2>
		<p class="text-xs text-slate-400 mt-0.5">Export batch data and compliance records</p>

		<div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4 mt-4">
			{#each reports as report}
				<div class="rounded-xl border border-slate-200 bg-slate-50 p-4">
					<div class="flex items-center gap-2.5">
						<div
							class="w-8 h-8 rounded-md bg-[#881B1B]/10 flex items-center justify-center shrink-0"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
								class="w-4 h-4 text-[#881B1B]"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d={report.icon} />
							</svg>
						</div>
						<span class="text-xs font-bold font-serif text-slate-900">{report.title}</span>
					</div>
					<button
						type="button"
						class="mt-3 w-full flex items-center justify-center gap-2 h-8 rounded-lg bg-slate-900 text-white text-xs font-semibold hover:bg-slate-800 transition-colors"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
							class="w-3.5 h-3.5"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3"
							/>
						</svg>
						Export Excel
					</button>
				</div>
			{/each}
		</div>
	</div>
</div>
