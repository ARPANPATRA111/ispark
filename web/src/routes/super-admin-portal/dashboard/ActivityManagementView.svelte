<script lang="ts">
	import { fade, slide } from 'svelte/transition';

	// ── Types ───────────────────────────────────────────────────────────────────
	type Track = 'Personality Development' | 'Skill Building';
	type ActivityType = 'Workshop' | 'Seminar' | 'Webinar' | 'Course';
	type ActivityStatus = 'Active' | 'Inactive';

	interface Activity {
		id: string;
		name: string;
		track: Track;
		type: ActivityType;
		credits: number;
		status: ActivityStatus;
	}

	// ── Mock Activity Data ──────────────────────────────────────────────────────
	let activities = $state<Activity[]>([
		{
			id: 'ACT001',
			name: 'Public Speaking Workshop',
			track: 'Personality Development',
			type: 'Workshop',
			credits: 3,
			status: 'Active'
		},
		{
			id: 'ACT002',
			name: 'Python for Beginners',
			track: 'Skill Building',
			type: 'Course',
			credits: 5,
			status: 'Active'
		},
		{
			id: 'ACT003',
			name: 'Leadership Essentials',
			track: 'Personality Development',
			type: 'Seminar',
			credits: 2,
			status: 'Active'
		},
		{
			id: 'ACT004',
			name: 'Data Analysis with Excel',
			track: 'Skill Building',
			type: 'Course',
			credits: 4,
			status: 'Active'
		},
		{
			id: 'ACT005',
			name: 'Emotional Intelligence',
			track: 'Personality Development',
			type: 'Webinar',
			credits: 2,
			status: 'Inactive'
		},
		{
			id: 'ACT006',
			name: 'Web Development Fundamentals',
			track: 'Skill Building',
			type: 'Course',
			credits: 6,
			status: 'Active'
		},
		{
			id: 'ACT007',
			name: 'Time Management Mastery',
			track: 'Personality Development',
			type: 'Workshop',
			credits: 2,
			status: 'Active'
		},
		{
			id: 'ACT008',
			name: 'Machine Learning Basics',
			track: 'Skill Building',
			type: 'Course',
			credits: 5,
			status: 'Inactive'
		},
		{
			id: 'ACT009',
			name: 'Conflict Resolution',
			track: 'Personality Development',
			type: 'Seminar',
			credits: 3,
			status: 'Active'
		},
		{
			id: 'ACT010',
			name: 'Digital Marketing Strategy',
			track: 'Skill Building',
			type: 'Course',
			credits: 4,
			status: 'Active'
		}
	]);

	let activityCounter = 11;

	// ── Credit Rules ────────────────────────────────────────────────────────────
	let creditRules = $state<Record<ActivityType, number>>({
		Workshop: 2,
		Seminar: 2,
		Webinar: 2,
		Course: 4
	});

	// ── Filters / Search / Pagination ───────────────────────────────────────────
	type FilterTab = 'All' | 'Personality Development' | 'Skill Building' | 'Active' | 'Inactive';
	let activeFilter = $state<FilterTab>('All');
	let searchQuery = $state('');
	let currentPage = $state(1);
	const pageSize = 10;

	let filteredActivities = $derived(
		activities.filter((a) => {
			const matchesFilter =
				activeFilter === 'All' || a.track === activeFilter || a.status === activeFilter;
			const matchesSearch =
				a.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				a.type.toLowerCase().includes(searchQuery.toLowerCase()) ||
				a.id.toLowerCase().includes(searchQuery.toLowerCase());
			return matchesFilter && matchesSearch;
		})
	);

	let totalPages = $derived(Math.max(1, Math.ceil(filteredActivities.length / pageSize)));

	let paginatedActivities = $derived(
		filteredActivities.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	function setFilter(tab: FilterTab) {
		activeFilter = tab;
		currentPage = 1;
	}

	function setSearch(value: string) {
		searchQuery = value;
		currentPage = 1;
	}

	function goToPage(page: number) {
		if (page < 1 || page > totalPages) return;
		currentPage = page;
	}

	// ── Stats (Derived) ──────────────────────────────────────────────────────────
	let totalActivitiesCount = $derived(activities.length);
	let activeActivitiesCount = $derived(activities.filter((a) => a.status === 'Active').length);
	let personalityDevCount = $derived(
		activities.filter((a) => a.track === 'Personality Development').length
	);
	let skillBuildingCount = $derived(activities.filter((a) => a.track === 'Skill Building').length);

	// ── Toast Notifications ──────────────────────────────────────────────────────
	interface Toast {
		id: number;
		message: string;
	}
	let toasts = $state<Toast[]>([]);
	let toastCounter = 0;

	function triggerToast(message: string) {
		const id = toastCounter++;
		toasts = [...toasts, { id, message }];
		setTimeout(() => {
			toasts = toasts.filter((t) => t.id !== id);
		}, 3000);
	}

	// ── Add / Edit Activity Modal ────────────────────────────────────────────────
	let isActivityModalOpen = $state(false);
	let modalMode = $state<'add' | 'edit'>('add');
	let editingId = $state<string | null>(null);

	let formName = $state('');
	let formTrack = $state<Track>('Personality Development');
	let formType = $state<ActivityType>('Workshop');
	let formCredits = $state(2);
	let formStatus = $state<ActivityStatus>('Active');

	function openAddActivity() {
		modalMode = 'add';
		editingId = null;
		formName = '';
		formTrack = 'Personality Development';
		formType = 'Workshop';
		formCredits = creditRules['Workshop'];
		formStatus = 'Active';
		isActivityModalOpen = true;
	}

	function openEditActivity(activity: Activity) {
		modalMode = 'edit';
		editingId = activity.id;
		formName = activity.name;
		formTrack = activity.track;
		formType = activity.type;
		formCredits = activity.credits;
		formStatus = activity.status;
		isActivityModalOpen = true;
	}

	function syncCreditsToType() {
		formCredits = creditRules[formType];
	}

	function handleSaveActivity(e: Event) {
		e.preventDefault();
		if (!formName.trim()) return;

		if (modalMode === 'add') {
			const newActivity: Activity = {
				id: `ACT${String(activityCounter++).padStart(3, '0')}`,
				name: formName.trim(),
				track: formTrack,
				type: formType,
				credits: formCredits,
				status: formStatus
			};
			activities = [newActivity, ...activities];
			triggerToast(`Activity "${newActivity.name}" published successfully!`);
		} else if (editingId) {
			activities = activities.map((a) =>
				a.id === editingId
					? {
							...a,
							name: formName.trim(),
							track: formTrack,
							type: formType,
							credits: formCredits,
							status: formStatus
						}
					: a
			);
			triggerToast(`Activity "${formName.trim()}" updated successfully!`);
		}
		isActivityModalOpen = false;
	}

	function handleDeleteActivity(activity: Activity) {
		if (confirm(`Are you sure you want to delete "${activity.name}"?`)) {
			activities = activities.filter((a) => a.id !== activity.id);
			triggerToast(`Activity "${activity.name}" deleted successfully.`);
		}
	}

	function toggleActivityStatus(activity: Activity) {
		activities = activities.map((a) =>
			a.id === activity.id ? { ...a, status: a.status === 'Active' ? 'Inactive' : 'Active' } : a
		);
		triggerToast(
			`"${activity.name}" marked as ${activity.status === 'Active' ? 'Inactive' : 'Active'}.`
		);
	}

	// ── View Activity Modal ───────────────────────────────────────────────────────
	let isViewModalOpen = $state(false);
	let viewingActivity = $state<Activity | null>(null);

	function openViewActivity(activity: Activity) {
		viewingActivity = activity;
		isViewModalOpen = true;
	}

	// ── Edit Credit Rules Modal ──────────────────────────────────────────────────
	let isCreditModalOpen = $state(false);
	let draftCreditRules = $state<Record<ActivityType, number>>({
		Workshop: 0,
		Seminar: 0,
		Webinar: 0,
		Course: 0
	});

	function openCreditRules() {
		draftCreditRules = { ...creditRules };
		isCreditModalOpen = true;
	}

	function handleSaveCreditRules(e: Event) {
		e.preventDefault();
		creditRules = { ...draftCreditRules };
		triggerToast('Credit rules updated successfully!');
		isCreditModalOpen = false;
	}

	// ── Style Helpers ─────────────────────────────────────────────────────────────
	function trackBadgeClass(track: Track): string {
		return track === 'Personality Development'
			? 'bg-purple-50 text-purple-700 border-purple-100'
			: 'bg-rose-50 text-rose-700 border-rose-100';
	}

	function statusTextClass(status: ActivityStatus): string {
		return status === 'Active' ? 'text-emerald-600' : 'text-slate-400';
	}

	function statusDotClass(status: ActivityStatus): string {
		return status === 'Active' ? 'bg-emerald-600' : 'bg-slate-400';
	}
</script>

<div class="space-y-6">
	<!-- ==================== STAT CARDS ==================== -->
	<section
		class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 select-none"
		aria-label="Activity metrics overview"
	>
		<!-- Total Activities -->
		<div
			class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{totalActivitiesCount}</span>
				<div class="p-2.5 rounded-lg bg-slate-100 text-slate-600 border border-slate-200">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
					Total Activities
				</h3>
				<span class="text-[11px] font-bold text-slate-400 mt-1 block">+3 this semester</span>
			</div>
		</div>

		<!-- Active Activities -->
		<div
			class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{activeActivitiesCount}</span>
				<div class="p-2.5 rounded-lg bg-emerald-50 text-emerald-600 border border-emerald-100">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 0 1 5.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
					Active Activities
				</h3>
				<span class="text-[11px] font-bold text-slate-400 mt-1 block">+2 this month</span>
			</div>
		</div>

		<!-- Personality Development -->
		<div
			class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{personalityDevCount}</span>
				<div class="p-2.5 rounded-lg bg-purple-50 text-purple-600 border border-purple-100">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.562.562 0 00-.586 0L6.982 21.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
					Personality Development
				</h3>
				<span class="text-[11px] font-bold text-slate-400 mt-1 block"
					>{personalityDevCount} workshops</span
				>
			</div>
		</div>

		<!-- Skill Building -->
		<div
			class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{skillBuildingCount}</span>
				<div class="p-2.5 rounded-lg bg-rose-50 text-rose-600 border border-rose-100">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
					Skill Building
				</h3>
				<span class="text-[11px] font-bold text-slate-400 mt-1 block"
					>{skillBuildingCount} courses</span
				>
			</div>
		</div>
	</section>

	<!-- ==================== MAIN GRID: OVERVIEW + MANAGE CREDITS ==================== -->
	<div class="grid grid-cols-1 lg:grid-cols-4 gap-6 items-start">
		<!-- Activity Management Overview Table (lg:col-span-2) -->
		<div
			class="lg:col-span-3 bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden flex flex-col"
		>
			<!-- Header -->
			<div
				class="p-5 border-b border-slate-100 flex items-center justify-between gap-3 bg-slate-50/20 select-none"
			>
				<div>
					<h3 class="text-sm font-bold font-serif text-slate-905">Activity Management Overview</h3>
					<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
						{filteredActivities.length} of {activities.length} activities
					</p>
				</div>
				<button
					type="button"
					onclick={openAddActivity}
					class="inline-flex items-center gap-1.5 px-4 py-2 bg-[#C23A3A] hover:bg-[#B03131] text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-colors focus:outline-none shrink-0"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2.5"
						stroke="currentColor"
						class="w-3.5 h-3.5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
					</svg>
					Add Activity
				</button>
			</div>

			<!-- Filters & Search -->
			<div
				class="p-5 border-b border-slate-100 flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white select-none"
			>
				<div class="flex flex-wrap gap-1.5">
					{#each ['All', 'Personality Development', 'Skill Building', 'Active', 'Inactive'] as const as tab}
						<button
							type="button"
							onclick={() => setFilter(tab)}
							class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all
								{activeFilter === tab
								? 'bg-[#C23A3A] text-white shadow-xs'
								: 'bg-slate-50 text-slate-500 hover:bg-slate-100'}"
						>
							{tab}
						</button>
					{/each}
				</div>

				<div class="relative w-full sm:w-64">
					<input
						type="text"
						value={searchQuery}
						oninput={(e) => setSearch((e.target as HTMLInputElement).value)}
						placeholder="Search activity name..."
						class="pl-4 pr-9 py-2 bg-slate-50 rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-slate-350 focus:bg-white w-full transition-all"
					/>
					<span class="absolute right-3 top-2.5 text-slate-400">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="2"
							stroke="currentColor"
							class="w-4 h-4"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
							/>
						</svg>
					</span>
				</div>
			</div>

			<!-- Table -->
			<div class="overflow-x-auto flex-grow no-scrollbar">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr
							class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
						>
							<th class="py-3.5 px-5">Activity Name</th>
							<th class="py-3.5 px-5">Track</th>
							<th class="py-3.5 px-5">Type</th>
							<th class="py-3.5 px-5">Credits</th>
							<th class="py-3.5 px-5">Status</th>
							<th class="py-3.5 px-5 text-center">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 text-xs font-sans">
						{#if paginatedActivities.length === 0}
							<tr>
								<td colspan="6" class="py-8 text-center text-slate-400 font-semibold select-none">
									No activities found matching search filters.
								</td>
							</tr>
						{:else}
							{#each paginatedActivities as activity (activity.id)}
								<tr class="hover:bg-slate-50/30 transition-colors">
									<td class="py-4 px-5">
										<div class="flex flex-col">
											<span class="font-bold text-slate-800">{activity.name}</span>
											<span class="text-[10px] text-slate-400 font-semibold mt-0.5 select-all"
												>{activity.id}</span
											>
										</div>
									</td>
									<td class="py-4 px-5">
										<span
											class="inline-flex items-center px-2 py-0.5 text-[10px] font-bold rounded-full border {trackBadgeClass(
												activity.track
											)}"
										>
											{activity.track === 'Personality Development'
												? 'Personality Dev.'
												: 'Skill Building'}
										</span>
									</td>
									<td class="py-4 px-5 text-slate-500 font-semibold">{activity.type}</td>
									<td class="py-4 px-5">
										<span class="inline-flex items-center gap-1 font-bold text-amber-600">
											<svg
												xmlns="http://www.w3.org/2000/svg"
												fill="none"
												viewBox="0 0 24 24"
												stroke-width="2"
												stroke="currentColor"
												class="w-3.5 h-3.5"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													d="M16.5 18.75h-9m9 0a3 3 0 013 3h-15a3 3 0 013-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497"
												/>
											</svg>
											{activity.credits}
										</span>
									</td>
									<td class="py-4 px-5">
										<button
											type="button"
											onclick={() => toggleActivityStatus(activity)}
											class="inline-flex items-center gap-1.5 font-bold {statusTextClass(
												activity.status
											)} hover:opacity-70 transition-opacity"
										>
											<span
												class="w-1.5 h-1.5 rounded-full shrink-0 {statusDotClass(
													activity.status
												)} {activity.status === 'Active' ? '' : ''}"
											></span>
											{activity.status}
										</button>
									</td>
									<td class="py-4 px-5">
										<div class="flex items-center justify-center gap-3 text-slate-400">
											<button
												type="button"
												onclick={() => openViewActivity(activity)}
												aria-label="View activity"
												class="text-blue-500 hover:text-blue-700 transition-colors p-0.5"
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													fill="none"
													viewBox="0 0 24 24"
													stroke-width="2"
													stroke="currentColor"
													class="w-4 h-4"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
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
												onclick={() => openEditActivity(activity)}
												aria-label="Edit activity"
												class="text-amber-500 hover:text-amber-700 transition-colors p-0.5"
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													fill="none"
													viewBox="0 0 24 24"
													stroke-width="2"
													stroke="currentColor"
													class="w-4 h-4"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
													/>
												</svg>
											</button>
											<button
												type="button"
												onclick={() => handleDeleteActivity(activity)}
												aria-label="Delete activity"
												class="text-rose-500 hover:text-rose-700 transition-colors p-0.5"
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													fill="none"
													viewBox="0 0 24 24"
													stroke-width="2"
													stroke="currentColor"
													class="w-4 h-4"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
													/>
												</svg>
											</button>
										</div>
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>

			<!-- Footer / Pagination -->
			<div
				class="p-4 border-t border-slate-100 bg-slate-50/30 flex items-center justify-between gap-3 select-none"
			>
				<span class="text-slate-500 font-semibold text-[11px]">
					Showing {paginatedActivities.length} of {filteredActivities.length} activities
				</span>
				{#if totalPages > 1}
					<div class="flex items-center gap-1.5">
						{#each Array(totalPages) as _, i}
							<button
								type="button"
								onclick={() => goToPage(i + 1)}
								class="w-7 h-7 flex items-center justify-center rounded-lg text-[11px] font-bold transition-colors
									{currentPage === i + 1
									? 'bg-[#C23A3A] text-white'
									: 'bg-white border border-slate-200 text-slate-500 hover:bg-slate-50'}"
							>
								{i + 1}
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>

		<!-- Manage Credits Panel (lg:col-span-1) -->
		<div class="bg-white border border-slate-200 rounded-xl p-5 shadow-xs space-y-4">
			<div class="flex items-center gap-3">
				<div class="p-2.5 rounded-lg bg-rose-50 text-rose-600 border border-rose-100 shrink-0">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M16.5 18.75h-9m9 0a3 3 0 013 3h-15a3 3 0 013-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497"
						/>
					</svg>
				</div>
				<div>
					<h3 class="text-sm font-bold font-serif text-slate-905">Manage Credits</h3>
					<p class="text-[10px] font-semibold text-slate-400 mt-0.5">Configure credit values</p>
				</div>
			</div>
			<div class="h-px bg-slate-100"></div>

			<div class="space-y-1">
				{#each Object.keys(creditRules) as ActivityType[] as type}
					<div
						class="flex items-center justify-between py-2.5 border-b border-slate-100 last:border-0"
					>
						<span class="text-xs font-semibold text-slate-600">{type}</span>
						<span class="inline-flex items-center gap-1 text-xs font-bold text-slate-800">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2"
								stroke="currentColor"
								class="w-3.5 h-3.5 text-amber-500"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M16.5 18.75h-9m9 0a3 3 0 013 3h-15a3 3 0 013-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497"
								/>
							</svg>
							{creditRules[type]} cr
						</span>
					</div>
				{/each}
			</div>

			<button
				type="button"
				onclick={openCreditRules}
				class="w-full mt-2 inline-flex items-center justify-center gap-1.5 py-2.5 bg-[#C23A3A] hover:bg-[#B03131] text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-colors focus:outline-none"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-4 h-4"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
					/>
				</svg>
				Edit Credit Rules
			</button>
		</div>
	</div>
</div>

<!-- ==================== TOAST NOTIFICATION CONTAINER ==================== -->
<div class="fixed bottom-6 right-6 z-50 flex flex-col gap-2 max-w-sm">
	{#each toasts as toast (toast.id)}
		<div
			transition:slide={{ duration: 150 }}
			class="p-4 bg-slate-800 border border-slate-700 text-white rounded-xl shadow-2xl flex items-center gap-2 text-xs font-semibold font-sans"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="2"
				stroke="currentColor"
				class="w-4 h-4 text-emerald-400"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M9 12.75 11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 01-1.043 3.296 3.745 3.745 0 01-3.296 1.043A3.745 3.745 0 0112 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 01-3.296-1.043 3.745 3.745 0 01-1.043-3.296A3.745 3.745 0 013 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 011.043-3.296 3.746 3.746 0 013.296-1.043A3.746 3.746 0 0112 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 013.296 1.043 3.746 3.746 0 011.043 3.296A3.745 3.745 0 0121 12Z"
				/>
			</svg>
			<span>{toast.message}</span>
		</div>
	{/each}
</div>

<!-- ==================== ADD / EDIT ACTIVITY MODAL ==================== -->
{#if isActivityModalOpen}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={(e) => {
			if (e.target === e.currentTarget) isActivityModalOpen = false;
		}}
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
	>
		<form
			onsubmit={handleSaveActivity}
			class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
		>
			<div class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30">
				<div>
					<h3 class="text-sm font-bold font-serif text-slate-900">
						{modalMode === 'add' ? 'Publish New Activity' : 'Edit Activity'}
					</h3>
					<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
						{modalMode === 'add' ? 'Register Extracurricular Activity' : 'Update activity details'}
					</p>
				</div>
				<button
					type="button"
					onclick={() => (isActivityModalOpen = false)}
					aria-label="Close modal"
					class="p-1 rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-600 transition-colors"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<div class="p-6 space-y-4">
				<div class="flex flex-col gap-1.5">
					<label for="act-name" class="text-[10px] font-extrabold text-slate-650 tracking-wider"
						>ACTIVITY NAME *</label
					>
					<input
						id="act-name"
						type="text"
						bind:value={formName}
						placeholder="e.g. Cyber Security Summit"
						required
						class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
					/>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label for="act-track" class="text-[10px] font-extrabold text-slate-650 tracking-wider"
							>TRACK *</label
						>
						<select
							id="act-track"
							bind:value={formTrack}
							class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-850 bg-white focus:outline-none focus:border-slate-355"
						>
							<option value="Personality Development">Personality Development</option>
							<option value="Skill Building">Skill Building</option>
						</select>
					</div>

					<div class="flex flex-col gap-1.5">
						<label for="act-type" class="text-[10px] font-extrabold text-slate-650 tracking-wider"
							>TYPE *</label
						>
						<select
							id="act-type"
							bind:value={formType}
							onchange={syncCreditsToType}
							class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-850 bg-white focus:outline-none focus:border-slate-355"
						>
							<option value="Workshop">Workshop</option>
							<option value="Seminar">Seminar</option>
							<option value="Webinar">Webinar</option>
							<option value="Course">Course</option>
						</select>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label
							for="act-credits"
							class="text-[10px] font-extrabold text-slate-650 tracking-wider">CREDITS *</label
						>
						<input
							id="act-credits"
							type="number"
							bind:value={formCredits}
							min="1"
							max="50"
							required
							class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
						/>
					</div>

					<div class="flex flex-col gap-1.5">
						<label for="act-status" class="text-[10px] font-extrabold text-slate-650 tracking-wider"
							>STATUS</label
						>
						<select
							id="act-status"
							bind:value={formStatus}
							class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-850 bg-white focus:outline-none focus:border-slate-355"
						>
							<option value="Active">Active</option>
							<option value="Inactive">Inactive</option>
						</select>
					</div>
				</div>
			</div>

			<div
				class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
			>
				<button
					type="button"
					onclick={() => (isActivityModalOpen = false)}
					class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					{modalMode === 'add' ? 'Publish Activity' : 'Save Changes'}
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- ==================== VIEW ACTIVITY MODAL ==================== -->
{#if isViewModalOpen && viewingActivity}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={(e) => {
			if (e.target === e.currentTarget) isViewModalOpen = false;
		}}
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
	>
		<div
			class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
		>
			<div class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30">
				<div>
					<h3 class="text-sm font-bold font-serif text-slate-905">Activity Details</h3>
					<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
						{viewingActivity.id}
					</p>
				</div>
				<button
					type="button"
					onclick={() => (isViewModalOpen = false)}
					aria-label="Close modal"
					class="p-1 rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-600 transition-colors"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<div class="p-6 space-y-1">
				<div class="flex items-center justify-between py-2.5 border-b border-slate-100">
					<span class="text-xs font-semibold text-slate-500">Activity Name</span>
					<span class="text-xs font-bold text-slate-900">{viewingActivity.name}</span>
				</div>
				<div class="flex items-center justify-between py-2.5 border-b border-slate-100">
					<span class="text-xs font-semibold text-slate-500">Track</span>
					<span
						class="inline-flex items-center px-2 py-0.5 text-[10px] font-bold rounded-full border {trackBadgeClass(
							viewingActivity.track
						)}"
					>
						{viewingActivity.track}
					</span>
				</div>
				<div class="flex items-center justify-between py-2.5 border-b border-slate-100">
					<span class="text-xs font-semibold text-slate-500">Type</span>
					<span class="text-xs font-bold text-slate-900">{viewingActivity.type}</span>
				</div>
				<div class="flex items-center justify-between py-2.5 border-b border-slate-100">
					<span class="text-xs font-semibold text-slate-500">Credits</span>
					<span class="text-xs font-bold text-slate-900">{viewingActivity.credits} cr</span>
				</div>
				<div class="flex items-center justify-between py-2.5">
					<span class="text-xs font-semibold text-slate-500">Status</span>
					<span
						class="inline-flex items-center gap-1.5 font-bold {statusTextClass(
							viewingActivity.status
						)}"
					>
						<span class="w-1.5 h-1.5 rounded-full shrink-0 {statusDotClass(viewingActivity.status)}"
						></span>
						{viewingActivity.status}
					</span>
				</div>
			</div>

			<div
				class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
			>
				<button
					type="button"
					onclick={() => (isViewModalOpen = false)}
					class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Close
				</button>
				<button
					type="button"
					onclick={() => {
						isViewModalOpen = false;
						if (viewingActivity) openEditActivity(viewingActivity);
					}}
					class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Edit Activity
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- ==================== EDIT CREDIT RULES MODAL ==================== -->
{#if isCreditModalOpen}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={(e) => {
			if (e.target === e.currentTarget) isCreditModalOpen = false;
		}}
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
	>
		<form
			onsubmit={handleSaveCreditRules}
			class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
		>
			<div class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30">
				<div>
					<h3 class="text-sm font-bold font-serif text-slate-900">Edit Credit Rules</h3>
					<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
						Configure credit values by activity type
					</p>
				</div>
				<button
					type="button"
					onclick={() => (isCreditModalOpen = false)}
					aria-label="Close modal"
					class="p-1 rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-600 transition-colors"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-5 h-5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<div class="p-6 space-y-4">
				{#each Object.keys(draftCreditRules) as ActivityType[] as type}
					<div class="flex items-center justify-between gap-4">
						<label for="credit-{type}" class="text-xs font-bold text-slate-700">{type}</label>
						<input
							id="credit-{type}"
							type="number"
							bind:value={draftCreditRules[type]}
							min="1"
							max="50"
							class="w-24 px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 text-right focus:outline-none focus:border-slate-355"
						/>
					</div>
				{/each}
			</div>

			<div
				class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
			>
				<button
					type="button"
					onclick={() => (isCreditModalOpen = false)}
					class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Save Credit Rules
				</button>
			</div>
		</form>
	</div>
{/if}

<style>
	/* Hide the horizontal scrollbar on the activity table while keeping it scrollable */
	.no-scrollbar {
		scrollbar-width: none; /* Firefox */
		-ms-overflow-style: none; /* IE / Edge legacy */
	}
	.no-scrollbar::-webkit-scrollbar {
		display: none; /* Chrome, Safari, Opera */
		width: 0;
		height: 0;
	}
</style>
