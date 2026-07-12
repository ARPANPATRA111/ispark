<script lang="ts">
	import { fade, slide } from 'svelte/transition';

	// ── Types ────────────────────────────────────────────────────────────────
	type TrackStatus = 'Active' | 'Inactive';

	interface Track {
		id: number;
		name: string;
		description: string;
		totalActivities: number;
		status: TrackStatus;
	}

	// ── Track Registry (mock data matching Personality Development / Skill Building) ──
	let tracks = $state<Track[]>([
		{
			id: 1,
			name: 'Personality Development',
			description: 'Activities focused on personal growth, communication, and leadership skills.',
			totalActivities: 48,
			status: 'Active'
		},
		{
			id: 2,
			name: 'Skill Building',
			description: 'Technical and vocational activities that develop practical competencies.',
			totalActivities: 36,
			status: 'Active'
		}
	]);

	let trackFilter = $state('All');
	let trackSearch = $state('');

	let filteredTracks = $derived(
		tracks.filter((t) => {
			const matchesFilter = trackFilter === 'All' || t.status === trackFilter;
			const matchesSearch =
				t.name.toLowerCase().includes(trackSearch.toLowerCase()) ||
				t.description.toLowerCase().includes(trackSearch.toLowerCase());
			return matchesFilter && matchesSearch;
		})
	);

	// ── Stat card derivations ────────────────────────────────────────────────
	let totalTracksCount = $derived(tracks.length);
	let activeTracksCount = $derived(tracks.filter((t) => t.status === 'Active').length);
	let personalityDevActivities = $derived(
		tracks.find((t) => t.name === 'Personality Development')?.totalActivities ?? 0
	);
	let skillBuildingActivities = $derived(
		tracks.find((t) => t.name === 'Skill Building')?.totalActivities ?? 0
	);

	// ── Add Track modal ──────────────────────────────────────────────────────
	let isAddTrackModalOpen = $state(false);
	let newTrackName = $state('');
	let newTrackDescription = $state('');
	let newTrackActivities = $state(0);

	function handleAddTrack(e: Event) {
		e.preventDefault();
		if (!newTrackName.trim()) return;
		const newId = tracks.length ? Math.max(...tracks.map((t) => t.id)) + 1 : 1;
		tracks = [
			...tracks,
			{
				id: newId,
				name: newTrackName.trim(),
				description: newTrackDescription.trim() || 'No description provided.',
				totalActivities: newTrackActivities || 0,
				status: 'Active'
			}
		];
		triggerToast(`Track "${newTrackName.trim()}" created successfully!`);
		newTrackName = '';
		newTrackDescription = '';
		newTrackActivities = 0;
		isAddTrackModalOpen = false;
	}

	// ── Edit Track modal ─────────────────────────────────────────────────────
	let isEditTrackModalOpen = $state(false);
	let editTrackId = $state(-1);
	let editTrackName = $state('');
	let editTrackDescription = $state('');
	let editTrackActivities = $state(0);
	let editTrackStatus = $state<TrackStatus>('Active');

	function openEditTrack(track: Track) {
		editTrackId = track.id;
		editTrackName = track.name;
		editTrackDescription = track.description;
		editTrackActivities = track.totalActivities;
		editTrackStatus = track.status;
		isEditTrackModalOpen = true;
	}

	function handleSaveTrack(e: Event) {
		e.preventDefault();
		if (editTrackId < 0 || !editTrackName.trim()) return;
		tracks = tracks.map((t) =>
			t.id === editTrackId
				? {
						...t,
						name: editTrackName.trim(),
						description: editTrackDescription.trim(),
						totalActivities: editTrackActivities,
						status: editTrackStatus
					}
				: t
		);
		triggerToast(`Track "${editTrackName.trim()}" updated successfully!`);
		isEditTrackModalOpen = false;
		editTrackId = -1;
	}

	// ── View Track modal ─────────────────────────────────────────────────────
	let isViewTrackModalOpen = $state(false);
	let viewTrack = $state<Track | null>(null);

	function openViewTrack(track: Track) {
		viewTrack = track;
		isViewTrackModalOpen = true;
	}

	// ── Toast notifications ──────────────────────────────────────────────────
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

	function handleDeleteTrack(track: Track) {
		if (confirm(`Are you sure you want to delete the "${track.name}" track?`)) {
			tracks = tracks.filter((t) => t.id !== track.id);
			triggerToast(`Track "${track.name}" removed successfully.`);
		}
	}

	function trackStatusClass(status: TrackStatus): string {
		return status === 'Active'
			? 'bg-emerald-50 text-emerald-700 border-emerald-100'
			: 'bg-slate-100 text-slate-500 border-slate-200';
	}

	function trackStatusDot(status: TrackStatus): string {
		return status === 'Active' ? 'bg-emerald-500' : 'bg-slate-400';
	}
</script>

<!-- ==================== TRACK MANAGEMENT ==================== -->
<div class="space-y-6">
	<!-- Overview Stat Cards -->
	<section
		class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 select-none"
		aria-label="Track management overview"
	>
		<!-- Total Tracks -->
		<div
			class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{totalTracksCount}</span>
				<div class="p-2.5 rounded-lg bg-slate-100 text-slate-600 border border-slate-200">
					<!-- Layers icon -->
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
							d="M12 3 2.25 8.25 12 13.5l9.75-5.25L12 3Z"
						/>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="m2.25 12.75 9.75 5.25 9.75-5.25"
						/>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="m2.25 17.25 9.75 5.25 9.75-5.25"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Total Tracks</h3>
				<span class="text-[11px] font-bold text-slate-400 mt-1 block">+0 this semester</span>
			</div>
		</div>

		<!-- Active Tracks -->
		<div
			class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{activeTracksCount}</span>
				<div class="p-2.5 rounded-lg bg-emerald-50 text-emerald-600 border border-emerald-100">
					<!-- Trending up icon -->
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
							d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 015.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Active Tracks</h3>
				<span class="text-[11px] font-bold text-slate-400 mt-1 block">2 new this year</span>
			</div>
		</div>

		<!-- Personality Development Activities -->
		<div
			class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{personalityDevActivities}</span>
				<div class="p-2.5 rounded-lg bg-purple-50 text-purple-600 border border-purple-100">
					<!-- Star icon -->
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
							d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.562.562 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.562.562 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
					Personality Development
				</h3>
				<span class="text-[11px] font-bold text-slate-400 mt-1 block">+6 this month</span>
			</div>
		</div>

		<!-- Skill Building Activities -->
		<div
			class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{skillBuildingActivities}</span>
				<div class="p-2.5 rounded-lg bg-rose-50 text-rose-600 border border-rose-100">
					<!-- Lightning bolt icon -->
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
				<span class="text-[11px] font-bold text-slate-400 mt-1 block">+4 this month</span>
			</div>
		</div>
	</section>

	<!-- Track Management Overview -->
	<section class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
		<!-- Header -->
		<div
			class="p-5 border-b border-slate-100 flex items-center justify-between gap-3 bg-slate-50/20 select-none"
		>
			<div>
				<h3 class="text-sm font-bold font-serif text-slate-905">Track Management Overview</h3>
				<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
					All registered tracks on the iSPARC platform
				</p>
			</div>
			<button
				type="button"
				onclick={() => (isAddTrackModalOpen = true)}
				class="inline-flex items-center gap-1.5 px-4 py-2 bg-[#C23A3A] hover:bg-[#B03131] text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-colors focus:outline-none shrink-0"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-4 h-4"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
				</svg>
				Add Track
			</button>
		</div>

		<!-- Filters & Search -->
		<div
			class="p-5 border-b border-slate-100 flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white select-none"
		>
			<div class="flex flex-wrap gap-1.5">
				{#each ['All', 'Active', 'Inactive'] as filterOption}
					<button
						type="button"
						onclick={() => (trackFilter = filterOption)}
						class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all
							{trackFilter === filterOption
							? 'bg-[#C23A3A] text-white shadow-xs'
							: 'bg-slate-50 text-slate-500 hover:bg-slate-100'}"
					>
						{filterOption === 'All' ? 'All Tracks' : filterOption}
					</button>
				{/each}
			</div>

			<div class="relative w-full sm:w-64">
				<input
					type="text"
					bind:value={trackSearch}
					placeholder="Search track name"
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
		<div class="overflow-x-auto">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr
						class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
					>
						<th class="py-3.5 px-5">Track Name</th>
						<th class="py-3.5 px-5">Description</th>
						<th class="py-3.5 px-5">Total Activities</th>
						<th class="py-3.5 px-5">Status</th>
						<th class="py-3.5 px-5">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-xs font-sans">
					{#if filteredTracks.length === 0}
						<tr>
							<td colspan="5" class="py-8 text-center text-slate-400 font-semibold select-none">
								No tracks found matching search filters.
							</td>
						</tr>
					{:else}
						{#each filteredTracks as track (track.id)}
							<tr class="hover:bg-slate-50/30 transition-colors">
								<td class="py-4 px-5 font-bold text-slate-800 align-top whitespace-nowrap">
									{track.name}
								</td>
								<td class="py-4 px-5 text-slate-500 font-semibold align-top max-w-sm">
									{track.description}
								</td>
								<td class="py-4 px-5 font-bold text-slate-800 align-top">
									{track.totalActivities}
								</td>
								<td class="py-4 px-5 align-top">
									<span
										class="inline-flex items-center gap-1.5 px-2 py-0.5 text-[10px] font-bold uppercase rounded-full border {trackStatusClass(
											track.status
										)}"
									>
										<span class="w-1.5 h-1.5 rounded-full {trackStatusDot(track.status)}"></span>
										{track.status}
									</span>
								</td>
								<td class="py-4 px-5 align-top">
									<div class="flex items-center gap-3">
										<button
											type="button"
											onclick={() => openViewTrack(track)}
											aria-label="View track"
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
											onclick={() => openEditTrack(track)}
											aria-label="Edit track"
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
													d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125"
												/>
											</svg>
										</button>
										<button
											type="button"
											onclick={() => handleDeleteTrack(track)}
											aria-label="Delete track"
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

		<!-- Footer -->
		<div
			class="p-4 border-t border-slate-100 bg-slate-50/30 text-slate-500 font-semibold text-[11px] select-none"
		>
			<span>Showing {filteredTracks.length} of {tracks.length} tracks</span>
		</div>
	</section>
</div>

<!-- ==================== TOAST NOTIFICATIONS ==================== -->
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

<!-- ==================== MODALS ==================== -->

<!-- Add Track Modal -->
{#if isAddTrackModalOpen}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={(e) => {
			if (e.target === e.currentTarget) isAddTrackModalOpen = false;
		}}
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
	>
		<form
			onsubmit={handleAddTrack}
			class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
		>
			<div class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30">
				<div>
					<h3 class="text-sm font-bold font-serif text-slate-900">Add New Track</h3>
					<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
						Register a platform track
					</p>
				</div>
				<button
					type="button"
					onclick={() => (isAddTrackModalOpen = false)}
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
					<label
						for="new-track-name"
						class="text-[10px] font-extrabold text-slate-650 tracking-wider">TRACK NAME *</label
					>
					<input
						id="new-track-name"
						type="text"
						bind:value={newTrackName}
						placeholder="e.g. Social Entrepreneurship & Innovation"
						required
						class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label
						for="new-track-description"
						class="text-[10px] font-extrabold text-slate-650 tracking-wider">DESCRIPTION</label
					>
					<textarea
						id="new-track-description"
						bind:value={newTrackDescription}
						rows="3"
						placeholder="Briefly describe the focus of this track"
						class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355 resize-none"
					></textarea>
				</div>

				<div class="flex flex-col gap-1.5">
					<label
						for="new-track-activities"
						class="text-[10px] font-extrabold text-slate-650 tracking-wider">TOTAL ACTIVITIES</label
					>
					<input
						id="new-track-activities"
						type="number"
						bind:value={newTrackActivities}
						min="0"
						class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
					/>
				</div>
			</div>

			<div
				class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
			>
				<button
					type="button"
					onclick={() => (isAddTrackModalOpen = false)}
					class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Add Track
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- Edit Track Modal -->
{#if isEditTrackModalOpen}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={(e) => {
			if (e.target === e.currentTarget) isEditTrackModalOpen = false;
		}}
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
	>
		<form
			onsubmit={handleSaveTrack}
			class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
		>
			<div class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30">
				<div>
					<h3 class="text-sm font-bold font-serif text-slate-900">Edit Track</h3>
					<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
						Update track details
					</p>
				</div>
				<button
					type="button"
					onclick={() => (isEditTrackModalOpen = false)}
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
					<label
						for="edit-track-name"
						class="text-[10px] font-extrabold text-slate-650 tracking-wider">TRACK NAME *</label
					>
					<input
						id="edit-track-name"
						type="text"
						bind:value={editTrackName}
						required
						class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
					/>
				</div>

				<div class="flex flex-col gap-1.5">
					<label
						for="edit-track-description"
						class="text-[10px] font-extrabold text-slate-650 tracking-wider">DESCRIPTION</label
					>
					<textarea
						id="edit-track-description"
						bind:value={editTrackDescription}
						rows="3"
						class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355 resize-none"
					></textarea>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label
							for="edit-track-activities"
							class="text-[10px] font-extrabold text-slate-650 tracking-wider"
							>TOTAL ACTIVITIES</label
						>
						<input
							id="edit-track-activities"
							type="number"
							bind:value={editTrackActivities}
							min="0"
							class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
						/>
					</div>
					<div class="flex flex-col gap-1.5">
						<label
							for="edit-track-status"
							class="text-[10px] font-extrabold text-slate-650 tracking-wider">STATUS</label
						>
						<select
							id="edit-track-status"
							bind:value={editTrackStatus}
							class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 bg-white focus:outline-none focus:border-slate-355"
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
					onclick={() => (isEditTrackModalOpen = false)}
					class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Save Changes
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- View Track Modal -->
{#if isViewTrackModalOpen && viewTrack}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={(e) => {
			if (e.target === e.currentTarget) isViewTrackModalOpen = false;
		}}
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
	>
		<div
			class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
		>
			<div class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30">
				<div>
					<h3 class="text-sm font-bold font-serif text-slate-900">{viewTrack.name}</h3>
					<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
						Track Details
					</p>
				</div>
				<button
					type="button"
					onclick={() => (isViewTrackModalOpen = false)}
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

			<div class="p-6 space-y-4 text-xs font-sans">
				<div class="flex flex-col gap-1">
					<span class="text-[10px] font-extrabold text-slate-650 tracking-wider">DESCRIPTION</span>
					<p class="text-slate-700 font-semibold leading-relaxed">{viewTrack.description}</p>
				</div>
				<div class="grid grid-cols-2 gap-4 pt-2">
					<div class="flex flex-col gap-1">
						<span class="text-[10px] font-extrabold text-slate-650 tracking-wider"
							>TOTAL ACTIVITIES</span
						>
						<span class="text-sm font-bold text-slate-900">{viewTrack.totalActivities}</span>
					</div>
					<div class="flex flex-col gap-1">
						<span class="text-[10px] font-extrabold text-slate-650 tracking-wider">STATUS</span>
						<span
							class="inline-flex items-center gap-1.5 px-2 py-0.5 w-fit text-[10px] font-bold uppercase rounded-full border {trackStatusClass(
								viewTrack.status
							)}"
						>
							<span class="w-1.5 h-1.5 rounded-full {trackStatusDot(viewTrack.status)}"></span>
							{viewTrack.status}
						</span>
					</div>
				</div>
			</div>

			<div
				class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
			>
				<button
					type="button"
					onclick={() => (isViewTrackModalOpen = false)}
					class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Close
				</button>
				<button
					type="button"
					onclick={() => {
						isViewTrackModalOpen = false;
						if (viewTrack) openEditTrack(viewTrack);
					}}
					class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
				>
					Edit Track
				</button>
			</div>
		</div>
	</div>
{/if}
