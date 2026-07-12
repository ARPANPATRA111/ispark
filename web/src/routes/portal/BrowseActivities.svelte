<script lang="ts">
	import { slide } from 'svelte/transition';
	import { SvelteDate } from 'svelte/reactivity';

	// Define interfaces for activities
	interface Activity {
		id: number;
		name: string;
		category: string;
		description: string;
		credits: number;
		mode: string;
		regDeadline: string;
		regDeadlineRaw: string;
		activityDate: string;
		venue: string;
		coordinator: string;
		status: 'Open' | 'Closed' | 'Closing Soon';
		daysLeft?: string;
	}

	import { API_BASE_URL } from '$lib/config';

	interface BackendActivity {
		id: number;
		name: string;
		category: string;
		description: string;
		credits: number;
		mode: string;
		reg_deadline: string;
		activity_date: string;
		venue: string;
		coordinator: string;
		status: 'Open' | 'Closed' | 'Closing Soon';
	}

	interface BackendEnrollment {
		activity_id: number;
	}

	let token = localStorage.getItem('access_token') || '';
	let allActivities = $state<Activity[]>([]);
	let enrolledIds = $state<number[]>([]);
	let stats = $state({
		activities_participated: 0,
		certificates_uploaded: 0,
		pending_certificates: 0,
		credits_earned: 0
	});

	function formatDate(dateStr: string) {
		if (!dateStr) return '';
		const d = new Date(dateStr);
		return d.toLocaleDateString('en-GB', {
			day: '2-digit',
			month: 'short',
			year: 'numeric'
		});
	}

	async function loadActivitiesData() {
		try {
			// Fetch activities
			const actRes = await fetch(`${API_BASE_URL}/api/student/activities`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
			if (actRes.ok) {
				const data = (await actRes.json()) as BackendActivity[];
				allActivities = (data || []).map((item) => ({
					id: item.id,
					name: item.name,
					category: item.category,
					description: item.description,
					credits: item.credits,
					mode: item.mode,
					regDeadline: formatDate(item.reg_deadline),
					regDeadlineRaw: item.reg_deadline,
					activityDate: formatDate(item.activity_date),
					venue: item.venue,
					coordinator: item.coordinator,
					status: item.status
				}));
			}

			// Fetch enrollments
			const enrollRes = await fetch(`${API_BASE_URL}/api/student/enrollments`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
			if (enrollRes.ok) {
				const enrollmentsData = (await enrollRes.json()) as BackendEnrollment[];
				enrolledIds = (enrollmentsData || []).map((e) => e.activity_id);
			}

			// Fetch stats
			const statsRes = await fetch(`${API_BASE_URL}/api/student/dashboard/stats`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
			if (statsRes.ok) {
				stats = await statsRes.json();
			}
		} catch (err) {
			console.error('Error fetching activities:', err);
		}
	}

	$effect(() => {
		loadActivitiesData();
	});

	// Derived stats responsive to enrollments
	let enrolledCount = $derived(enrolledIds.length);
	let creditsEarned = $derived(stats.credits_earned);
	let pendingCount = $derived(stats.pending_certificates);
	let completedCount = $derived(stats.activities_participated);

	// Dynamic registration deadlines (Open / Closing Soon activities sorted by deadline)
	let deadlines = $derived.by(() => {
		const today = new SvelteDate();
		today.setHours(0, 0, 0, 0);

		return allActivities
			.filter((a) => a.status === 'Open' || a.status === 'Closing Soon')
			.map((a) => {
				const deadlineDate = new SvelteDate(a.regDeadlineRaw);
				deadlineDate.setHours(0, 0, 0, 0);
				const diffTime = deadlineDate.getTime() - today.getTime();
				const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

				let daysStr: string;
				let colorClass = 'bg-slate-50 text-slate-700 border-slate-100';

				if (diffDays < 0) {
					daysStr = 'Closed';
				} else if (diffDays === 0) {
					daysStr = 'Today';
					colorClass = 'bg-rose-50 text-rose-700 border-rose-100';
				} else {
					daysStr = `${diffDays}d left`;
					if (diffDays <= 3) {
						colorClass = 'bg-rose-50 text-rose-700 border-rose-100';
					} else if (diffDays <= 7) {
						colorClass = 'bg-amber-50 text-amber-700 border-amber-100';
					} else {
						colorClass = 'bg-emerald-50 text-emerald-700 border-emerald-100';
					}
				}

				return {
					name: a.name,
					days: daysStr,
					color: colorClass,
					rawDays: diffDays
				};
			})
			.filter((a) => a.rawDays >= 0)
			.sort((a, b) => a.rawDays - b.rawDays)
			.slice(0, 4);
	});

	// Categories dot color mappings
	const categoryColors: Record<string, string> = {
		TECHNICAL: 'bg-blue-600',
		RESEARCH: 'bg-purple-600',
		SPORTS: 'bg-emerald-600',
		CULTURAL: 'bg-pink-600',
		LEADERSHIP: 'bg-teal-600',
		'SOCIAL SERVICE': 'bg-rose-600',
		'PUBLIC SPEAKING': 'bg-amber-600'
	};

	// Dynamic category counts from allActivities
	let categoriesCounts = $derived.by(() => {
		const counts: Record<string, number> = {};
		allActivities.forEach((a) => {
			const catUpper = a.category.toUpperCase();
			counts[catUpper] = (counts[catUpper] || 0) + 1;
		});

		const order = [
			'TECHNICAL',
			'RESEARCH',
			'SPORTS',
			'CULTURAL',
			'LEADERSHIP',
			'SOCIAL SERVICE',
			'PUBLIC SPEAKING'
		];
		const extra = Object.keys(counts).filter((c) => !order.includes(c));
		const finalOrder = [...order, ...extra];

		return finalOrder.map((catUpper) => {
			const name = catUpper
				.split(' ')
				.map((w) => w.charAt(0) + w.slice(1).toLowerCase())
				.join(' ');
			return {
				name,
				count: counts[catUpper] || 0,
				dotColor: categoryColors[catUpper] || 'bg-slate-600'
			};
		});
	});

	// Filters State
	let showFilters = $state(true);
	let tempFilters = $state({
		search: '',
		category: '',
		status: '',
		credits: '',
		mode: ''
	});
	let activeFilters = $state({
		search: '',
		category: '',
		status: '',
		credits: '',
		mode: ''
	});

	// Pagination State
	let currentPage = $state(1);
	const itemsPerPage = 6;

	// Filter logic
	let filteredActivities = $derived.by(() => {
		return allActivities.filter((activity) => {
			// 1. Search Query
			if (activeFilters.search) {
				const query = activeFilters.search.toLowerCase();
				const matchesName = activity.name.toLowerCase().includes(query);
				const matchesDesc = activity.description.toLowerCase().includes(query);
				if (!matchesName && !matchesDesc) return false;
			}

			// 2. Category
			if (activeFilters.category && activeFilters.category !== 'ALL') {
				if (activity.category !== activeFilters.category) return false;
			}

			// 3. Status
			if (activeFilters.status) {
				if (activity.status !== activeFilters.status) return false;
			}

			// 4. Credits Range
			if (activeFilters.credits) {
				if (activeFilters.credits === '1-10' && (activity.credits < 1 || activity.credits > 10))
					return false;
				if (activeFilters.credits === '11-15' && (activity.credits < 11 || activity.credits > 15))
					return false;
				if (activeFilters.credits === '16-20' && (activity.credits < 16 || activity.credits > 20))
					return false;
			}

			// 5. Mode
			if (activeFilters.mode) {
				if (activity.mode !== activeFilters.mode) return false;
			}

			return true;
		});
	});

	// Paginated activities
	let paginatedActivities = $derived.by(() => {
		const start = (currentPage - 1) * itemsPerPage;
		return filteredActivities.slice(start, start + itemsPerPage);
	});

	let totalPages = $derived(Math.ceil(filteredActivities.length / itemsPerPage) || 1);

	// Actions
	function handleApplyFilters() {
		activeFilters = { ...tempFilters };
		currentPage = 1;
	}

	function handleClearFilters() {
		tempFilters = { search: '', category: '', status: '', credits: '', mode: '' };
		activeFilters = { search: '', category: '', status: '', credits: '', mode: '' };
		currentPage = 1;
	}

	function handleRecommendedClick(category: string) {
		tempFilters.category = category.toUpperCase();
		handleApplyFilters();
	}

	function handleCategoryWidgetClick(category: string) {
		tempFilters.category = category.toUpperCase();
		handleApplyFilters();
	}

	async function toggleEnrollment(id: number) {
		if (enrolledIds.includes(id)) {
			// Already enrolled, no action needed in this mode
			return;
		}

		try {
			const res = await fetch(`${API_BASE_URL}/api/student/activities/${id}/enroll`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${token}`
				}
			});

			if (res.ok) {
				enrolledIds = [...enrolledIds, id];
				// Refresh statistics
				const statsRes = await fetch(`${API_BASE_URL}/api/student/dashboard/stats`, {
					headers: {
						Authorization: `Bearer ${token}`
					}
				});
				if (statsRes.ok) {
					stats = await statsRes.json();
				}
			} else {
				const errData = await res.json();
				alert(errData.error || 'Failed to enroll');
			}
		} catch (err) {
			console.error(err);
			alert('Error enrolling in activity');
		}
	}

	function getCategoryBadgeColor(category: string) {
		switch (category) {
			case 'TECHNICAL':
				return 'bg-slate-100 text-slate-700 border-slate-205 text-[10px]';
			case 'SPORTS':
				return 'bg-emerald-50 text-emerald-700 border-emerald-200 text-[10px]';
			case 'RESEARCH':
				return 'bg-sky-50 text-sky-700 border-sky-200 text-[10px]';
			case 'PUBLIC SPEAKING':
				return 'bg-amber-50 text-amber-700 border-amber-200 text-[10px]';
			case 'SOCIAL SERVICE':
				return 'bg-rose-50 text-rose-700 border-rose-200 text-[10px]';
			case 'CULTURAL':
				return 'bg-purple-50 text-purple-700 border-purple-200 text-[10px]';
			default:
				return 'bg-slate-100 text-slate-700 border-slate-200 text-[10px]';
		}
	}
</script>

<div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start font-sans">
	<!-- LEFT COLUMN (Main browse grid & search) -->
	<div class="lg:col-span-9 space-y-6">
		<!-- Find Activities Panel -->
		<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs">
			<div class="flex items-center justify-between pb-3 border-b border-slate-100 mb-4">
				<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
					Find Activities
				</h3>
				<button
					onclick={() => (showFilters = !showFilters)}
					class="text-xs font-bold text-[#881B1B] hover:underline"
				>
					{showFilters ? 'Hide Filters' : 'Show Filters'}
				</button>
			</div>

			{#if showFilters}
				<div transition:slide={{ duration: 200 }} class="space-y-4">
					<!-- Text search -->
					<div class="relative">
						<input
							type="text"
							bind:value={tempFilters.search}
							placeholder="Search activities, competitions, workshops..."
							class="w-full pl-3 pr-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
						/>
					</div>

					<!-- Dropdowns grid -->
					<div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs">
						<!-- Category -->
						<div class="flex flex-col gap-1.5">
							<label
								for="filter-category"
								class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
								>Category</label
							>
							<select
								id="filter-category"
								bind:value={tempFilters.category}
								class="w-full px-2.5 py-2 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all cursor-pointer"
							>
								<option value="">All Categories</option>
								<option value="TECHNICAL">Technical</option>
								<option value="RESEARCH">Research</option>
								<option value="SPORTS">Sports</option>
								<option value="CULTURAL">Cultural</option>
								<option value="SOCIAL SERVICE">Social Service</option>
								<option value="PUBLIC SPEAKING">Public Speaking</option>
								<option value="LEADERSHIP">Leadership</option>
							</select>
						</div>

						<!-- Status -->
						<div class="flex flex-col gap-1.5">
							<label
								for="filter-status"
								class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Status</label
							>
							<select
								id="filter-status"
								bind:value={tempFilters.status}
								class="w-full px-2.5 py-2 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all cursor-pointer"
							>
								<option value="">All</option>
								<option value="Open">Open</option>
								<option value="Closing Soon">Closing Soon</option>
								<option value="Closed">Closed</option>
							</select>
						</div>

						<!-- Credits Range -->
						<div class="flex flex-col gap-1.5">
							<label
								for="filter-credits"
								class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
								>Credits Range</label
							>
							<select
								id="filter-credits"
								bind:value={tempFilters.credits}
								class="w-full px-2.5 py-2 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all cursor-pointer"
							>
								<option value="">All Credits</option>
								<option value="1-10">1 - 10 Points</option>
								<option value="11-15">11 - 15 Points</option>
								<option value="16-20">16 - 20 Points</option>
							</select>
						</div>

						<!-- Mode -->
						<div class="flex flex-col gap-1.5">
							<label
								for="filter-mode"
								class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Mode</label
							>
							<select
								id="filter-mode"
								bind:value={tempFilters.mode}
								class="w-full px-2.5 py-2 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all cursor-pointer"
							>
								<option value="">All Modes</option>
								<option value="Offline">Offline</option>
								<option value="Online">Online</option>
								<option value="Hybrid">Hybrid</option>
							</select>
						</div>
					</div>

					<!-- Buttons -->
					<div class="flex items-center gap-3 pt-2">
						<button
							onclick={handleApplyFilters}
							class="px-5 py-2 bg-[#0B1535] hover:bg-[#1a2b5e] text-white rounded-lg text-xs font-bold transition duration-200 cursor-pointer"
						>
							Apply Filters
						</button>
						<button
							onclick={handleClearFilters}
							class="px-5 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 rounded-lg text-xs font-bold bg-white transition duration-200 cursor-pointer"
						>
							Clear Filters
						</button>
					</div>
				</div>
			{/if}
		</div>

		<!-- Recommended For You Row -->
		<div class="space-y-2.5">
			<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
				Recommended For You <span class="text-[10px] text-slate-400 normal-case font-medium"
					>Based on your enrolled tracks</span
				>
			</h3>
			<div class="flex flex-wrap gap-2.5">
				<button
					onclick={() => handleRecommendedClick('TECHNICAL')}
					class="flex items-center gap-2 px-4 py-2 border border-[#881B1B]/20 bg-red-50 text-[#881B1B] rounded-lg text-xs font-bold hover:bg-[#881B1B]/10 transition cursor-pointer"
				>
					<!-- Sparkles Icon -->
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
							d="M9.813 15.904L9 21l-.813-5.096L3 15l5.096-.813L9 9l.813 5.096L15 15l-5.187.904zM18 5.25L17.25 8 16.5 5.25 13.75 4.5l2.75-.75L17.25 1l.75 2.75 2.75.75-2.75.75z"
						/>
					</svg>
					Technical Skill Track
					<span class="bg-[#881B1B] text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full"
						>6</span
					>
				</button>
				<button
					onclick={() => handleRecommendedClick('RESEARCH')}
					class="flex items-center gap-2 px-4 py-2 border border-purple-200 bg-purple-50 text-purple-700 rounded-lg text-xs font-bold hover:bg-purple-100 transition cursor-pointer"
				>
					<!-- Beaker/Science Icon -->
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
							d="M9.75 3.104v11.896m5.25-11.896v11.896M2.25 21h19.5M10 3.75h4"
						/>
					</svg>
					Research & Innovation
					<span class="bg-purple-700 text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full"
						>4</span
					>
				</button>
				<button
					onclick={() => handleRecommendedClick('PUBLIC SPEAKING')}
					class="flex items-center gap-2 px-4 py-2 border border-amber-200 bg-amber-50 text-amber-700 rounded-lg text-xs font-bold hover:bg-amber-100 transition cursor-pointer"
				>
					<!-- Mic Icon -->
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
							d="M12 18.75a6 6 0 0 0 6-6v-1.5m-6 7.5a6 6 0 0 1-6-6v-1.5m6 7.5v3.75m-3.75 0h7.5M12 15.75a3 3 0 0 0 3-3V4.5a3 3 0 0 0-6 0v8.25a3 3 0 0 0 3 3z"
						/>
					</svg>
					Public Speaking
					<span class="bg-amber-700 text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full"
						>3</span
					>
				</button>
			</div>
		</div>

		<!-- Activities Grid Section -->
		<div class="space-y-4">
			<div class="flex items-baseline justify-between border-b border-slate-200 pb-2">
				<h2 class="text-base font-bold font-serif text-[#0B1535]">
					All Activities <span class="text-xs font-semibold text-slate-400 normal-case ml-2"
						>{filteredActivities.length} found</span
					>
				</h2>
			</div>

			{#if paginatedActivities.length > 0}
				<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
					{#each paginatedActivities as activity (activity.id)}
						<div
							class="bg-white rounded-2xl border border-slate-200 overflow-hidden flex flex-col justify-between shadow-xs hover:shadow-md transition-all duration-200 relative group"
						>
							<!-- Top Banner / Color strip depending on category -->
							<div class="h-1.5 bg-[#881B1B]/80"></div>

							<div class="p-5 flex-grow space-y-4">
								<!-- Category Badge -->
								<div class="flex justify-between items-center">
									<span
										class="inline-flex px-2 py-0.5 rounded border text-[9px] font-bold uppercase tracking-wider {getCategoryBadgeColor(
											activity.category
										)}"
									>
										{activity.category}
									</span>

									{#if activity.status === 'Closing Soon'}
										<span
											class="text-[9px] font-bold text-rose-600 bg-rose-50 border border-rose-100 px-2 py-0.5 rounded-full uppercase tracking-wider"
										>
											Closing Soon
										</span>
									{/if}
								</div>

								<!-- Title & Description -->
								<div>
									<h4
										class="text-xs font-bold text-[#0B1535] leading-snug group-hover:text-[#881B1B] transition-colors"
									>
										{activity.name}
									</h4>
									<p class="text-[11px] text-slate-500 mt-2 leading-relaxed">
										{activity.description}
									</p>
								</div>

								<!-- Details List -->
								<div
									class="space-y-2 pt-3 border-t border-slate-100 text-[11px] font-medium text-slate-500"
								>
									<!-- Credits -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M11.48 3.499c.15-.358.682-.358.832 0l1.97 3.99a.75.75 0 0 0 .58.58l4.42 1.405c.4.082.56.58.268.88l-3.2 3.12a.75.75 0 0 0-.214.655l.8 4.417c.07.417-.37.74-.75.54L12 18.002l-3.95 2.06c-.38.2-.85-.122-.75-.54l.8-4.417a.75.75 0 0 0-.214-.655L4.69 11.235c-.29-.3-.13-.8.268-.88l4.42-1.405a.75.75 0 0 0 .58-.58l1.97-3.99z"
											/>
										</svg>
										<span
											>Credits: <strong class="text-slate-800">{activity.credits} Points</strong
											></span
										>
									</div>

									<!-- Mode -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M9.75 17L9 20l-1 3h4l-.5-3 .5-3m6-12V15a2.25 2.25 0 0 1-2.25 2.25H5.25A2.25 2.25 0 0 1 3 15V5.25m18 0A2.25 2.25 0 0 0 18.75 3H5.25A2.25 2.25 0 0 0 3 5.25"
											/>
										</svg>
										<span>Mode: <strong class="text-slate-800">{activity.mode}</strong></span>
									</div>

									<!-- Registration Deadline -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0z"
											/>
										</svg>
										<span
											>Registration Deadline: <strong class="text-slate-800"
												>{activity.regDeadline}</strong
											></span
										>
									</div>

									<!-- Activity Date -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5"
											/>
										</svg>
										<span
											>Activity Date: <strong class="text-slate-800">{activity.activityDate}</strong
											></span
										>
									</div>

									<!-- Venue -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
											/>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25s-7.5-4.108-7.5-11.25A7.5 7.5 0 1119.5 10.5z"
											/>
										</svg>
										<span>Venue: <strong class="text-slate-800">{activity.venue}</strong></span>
									</div>

									<!-- Coordinator -->
									<div class="flex items-center gap-2">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2"
											stroke="currentColor"
											class="w-4 h-4 text-slate-400"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0zM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632z"
											/>
										</svg>
										<span
											>Coordinator: <strong class="text-slate-800">{activity.coordinator}</strong
											></span
										>
									</div>
								</div>
							</div>

							<!-- Button -->
							<div class="px-5 pb-5">
								{#if enrolledIds.includes(activity.id)}
									<button
										onclick={() => toggleEnrollment(activity.id)}
										class="w-full py-2.5 bg-[#881B1B] hover:bg-[#731717] text-white rounded-lg text-xs font-bold transition duration-200 cursor-pointer flex items-center justify-center gap-1.5 shadow-sm font-sans uppercase tracking-wider"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2.5"
											stroke="currentColor"
											class="w-4 h-4"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M4.5 12.75l6 6 9-13.5"
											/>
										</svg>
										Enrolled
									</button>
								{:else}
									<button
										onclick={() => toggleEnrollment(activity.id)}
										disabled={activity.status === 'Closed'}
										class="w-full py-2.5 bg-[#881B1B] hover:bg-[#731717] disabled:bg-slate-200 disabled:text-slate-400 text-white rounded-lg text-xs font-bold transition duration-200 cursor-pointer shadow-xs focus:outline-none font-sans uppercase tracking-wider"
									>
										{activity.status === 'Closed' ? 'Registration Closed' : 'Enroll Now'}
									</button>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				<!-- Pagination Controls -->
				{#if totalPages > 1}
					<div
						class="flex items-center justify-center gap-2 pt-6 border-t border-slate-200 mt-6 font-sans"
					>
						<button
							disabled={currentPage === 1}
							onclick={() => (currentPage -= 1)}
							class="px-3 py-1.5 border border-slate-200 rounded-lg text-xs font-bold bg-white text-slate-600 hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
						>
							&lt; Previous
						</button>

						{#each Array(totalPages) as _, index}
							<button
								onclick={() => (currentPage = index + 1)}
								class="w-8 h-8 rounded-lg text-xs font-bold transition flex items-center justify-center border {currentPage ===
								index + 1
									? 'bg-[#0B1535] border-[#0B1535] text-white'
									: 'bg-white border-slate-200 text-slate-600 hover:bg-slate-50'}"
							>
								{index + 1}
							</button>
						{/each}

						<button
							disabled={currentPage === totalPages}
							onclick={() => (currentPage += 1)}
							class="px-3 py-1.5 border border-slate-200 rounded-lg text-xs font-bold bg-white text-slate-600 hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
						>
							Next &gt;
						</button>
					</div>
				{/if}
			{:else}
				<div class="bg-white border border-slate-200 rounded-xl p-12 text-center shadow-xs">
					<p class="text-slate-400 font-semibold text-xs">
						No activities found matching the filters.
					</p>
					<button
						onclick={handleClearFilters}
						class="text-[#881B1B] font-bold text-xs mt-3 hover:underline"
					>
						Clear Filters & Reset
					</button>
				</div>
			{/if}
		</div>
	</div>

	<!-- RIGHT COLUMN (Widgets sidebar) -->
	<div class="lg:col-span-3 space-y-6">
		<!-- Registration Deadlines Widget -->
		<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-2.5">
				<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
					Registration Deadlines
				</h3>
			</div>
			<div class="space-y-3 pt-1">
				{#each deadlines as item}
					<div class="flex items-center justify-between text-xs font-semibold text-slate-800">
						<span class="truncate pr-2 font-serif">{item.name}</span>
						<span
							class="px-2 py-0.5 rounded text-[10px] font-bold border shrink-0 font-sans {item.color}"
						>
							{item.days}
						</span>
					</div>
				{/each}
			</div>
		</div>

		<!-- My Enrollment Status Widget -->
		<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-2.5">
				<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
					My Enrollment Status
				</h3>
			</div>
			<div class="space-y-3 pt-1 text-xs font-semibold text-slate-500">
				<div class="flex justify-between items-center">
					<span>Enrolled</span>
					<span class="text-slate-900 font-bold">{enrolledCount}</span>
				</div>
				<div class="flex justify-between items-center">
					<span>Pending Approval</span>
					<span class="text-slate-900 font-bold">{pendingCount}</span>
				</div>
				<div class="flex justify-between items-center">
					<span>Completed</span>
					<span class="text-slate-900 font-bold">{completedCount}</span>
				</div>
				<div class="flex justify-between items-center pt-2 border-t border-slate-100">
					<span>Credits Earned</span>
					<span class="text-[#881B1B] font-extrabold text-sm">{creditsEarned}</span>
				</div>
			</div>
		</div>

		<!-- Activity Categories Widget -->
		<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-2.5">
				<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider font-sans">
					Activity Categories
				</h3>
			</div>
			<div class="space-y-2.5 pt-1">
				{#each categoriesCounts as cat}
					<button
						onclick={() => handleCategoryWidgetClick(cat.name)}
						class="w-full flex items-center justify-between text-xs font-semibold text-slate-600 hover:text-slate-900 hover:bg-slate-50 p-1.5 rounded transition group cursor-pointer"
					>
						<div class="flex items-center gap-2">
							<span class="w-1.5 h-1.5 rounded-full {cat.dotColor}"></span>
							<span class="group-hover:font-bold">{cat.name}</span>
						</div>
						<span
							class="bg-slate-100 text-slate-700 font-bold text-[10px] px-2 py-0.5 rounded-full font-sans"
						>
							{cat.count}
						</span>
					</button>
				{/each}
			</div>
		</div>
	</div>
</div>
