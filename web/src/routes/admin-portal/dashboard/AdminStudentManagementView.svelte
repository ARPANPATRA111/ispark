<script lang="ts">
	import { fade, slide } from 'svelte/transition';
	import AdminStudentDetailView from './AdminStudentDetailView.svelte';

	// ── Types ──────────────────────────────────────────────────────────────────
	type Status = 'Active' | 'At Risk' | 'Pending Review' | 'Inactive';

	interface Student {
		id: string;
		name: string;
		regNo: string;
		department: string;
		semester: number;
		creditsEarned: number;
		creditsTarget: number;
		certificates: number;
		activityCount: number;
		status: Status;
		email: string;
		batch: string;
	}

	// ── Mock Data ──────────────────────────────────────────────────────────────
	const allStudents: Student[] = [
		{
			id: 'S001',
			name: 'Rahul Sharma',
			regNo: 'EN2022001',
			department: 'Computer Science',
			semester: 6,
			creditsEarned: 162,
			creditsTarget: 200,
			certificates: 4,
			activityCount: 19,
			status: 'Active',
			email: 'rahul.sharma@iips.edu',
			batch: '2022'
		},
		{
			id: 'S002',
			name: 'Priya Patel',
			regNo: 'EN2023012',
			department: 'Electronics',
			semester: 4,
			creditsEarned: 98,
			creditsTarget: 200,
			certificates: 3,
			activityCount: 7,
			status: 'At Risk',
			email: 'priya.patel@iips.edu',
			batch: '2023'
		},
		{
			id: 'S003',
			name: 'Arjun Patel',
			regNo: 'EN2022018',
			department: 'Mechanical',
			semester: 5,
			creditsEarned: 125,
			creditsTarget: 200,
			certificates: 4,
			activityCount: 18,
			status: 'Active',
			email: 'arjun.patel@iips.edu',
			batch: '2022'
		},
		{
			id: 'S004',
			name: 'Sneha Kumar',
			regNo: 'EN2022015',
			department: 'Computer Science',
			semester: 6,
			creditsEarned: 180,
			creditsTarget: 200,
			certificates: 7,
			activityCount: 16,
			status: 'Active',
			email: 'sneha.kumar@iips.edu',
			batch: '2022'
		},
		{
			id: 'S005',
			name: 'Vikram Singh',
			regNo: 'EN2023003',
			department: 'Civil',
			semester: 4,
			creditsEarned: 72,
			creditsTarget: 200,
			certificates: 1,
			activityCount: 3,
			status: 'Pending Review',
			email: 'vikram.singh@iips.edu',
			batch: '2023'
		},
		{
			id: 'S006',
			name: 'Kavya Krishnan',
			regNo: 'EN2022008',
			department: 'Electronics',
			semester: 5,
			creditsEarned: 110,
			creditsTarget: 200,
			certificates: 4,
			activityCount: 9,
			status: 'Active',
			email: 'kavya.krishnan@iips.edu',
			batch: '2022'
		},
		{
			id: 'S007',
			name: 'Dev Mehra',
			regNo: 'EN2023014',
			department: 'Mechanical',
			semester: 3,
			creditsEarned: 55,
			creditsTarget: 200,
			certificates: 0,
			activityCount: 5,
			status: 'Pending Review',
			email: 'dev.mehra@iips.edu',
			batch: '2023'
		},
		{
			id: 'S008',
			name: 'Anita Nair',
			regNo: 'EN2022004',
			department: 'Computer Science',
			semester: 6,
			creditsEarned: 138,
			creditsTarget: 200,
			certificates: 4,
			activityCount: 14,
			status: 'Active',
			email: 'anita.nair@iips.edu',
			batch: '2022'
		},
		{
			id: 'S009',
			name: 'Rohan Verma',
			regNo: 'EN2023009',
			department: 'Civil',
			semester: 5,
			creditsEarned: 89,
			creditsTarget: 200,
			certificates: 2,
			activityCount: 6,
			status: 'At Risk',
			email: 'rohan.verma@iips.edu',
			batch: '2023'
		},
		{
			id: 'S010',
			name: 'Meera Iyer',
			regNo: 'EN2023021',
			department: 'Computer Science',
			semester: 4,
			creditsEarned: 115,
			creditsTarget: 200,
			certificates: 3,
			activityCount: 8,
			status: 'Active',
			email: 'meera.iyer@iips.edu',
			batch: '2023'
		},
		{
			id: 'S011',
			name: 'Karthik Rao',
			regNo: 'EN2022011',
			department: 'Electronics',
			semester: 6,
			creditsEarned: 145,
			creditsTarget: 200,
			certificates: 5,
			activityCount: 12,
			status: 'Active',
			email: 'karthik.rao@iips.edu',
			batch: '2022'
		},
		{
			id: 'S012',
			name: 'Divya Menon',
			regNo: 'EN2023017',
			department: 'Mechanical',
			semester: 3,
			creditsEarned: 44,
			creditsTarget: 200,
			certificates: 0,
			activityCount: 2,
			status: 'Inactive',
			email: 'divya.menon@iips.edu',
			batch: '2023'
		},
		{
			id: 'S013',
			name: 'Palak Rai',
			regNo: 'EN2022007',
			department: 'Computer Science',
			semester: 6,
			creditsEarned: 185,
			creditsTarget: 200,
			certificates: 2,
			activityCount: 11,
			status: 'Active',
			email: 'palak.rai@iips.edu',
			batch: '2022'
		}
	];

	// ── Derived Stats ──────────────────────────────────────────────────────────
	const totalStudents = allStudents.length;
	const activeStudents = allStudents.filter((s) => s.status === 'Active').length;
	const pendingCertReviews = 7;
	const avgCredits = Math.round(
		allStudents.reduce((sum, s) => sum + s.creditsEarned, 0) / allStudents.length
	);

	// Student Overview highlights
	const perfScore = (s: Student) => s.creditsEarned + s.certificates * 15 + s.activityCount * 2;
	const topPerformer = [...allStudents].sort((a, b) => perfScore(b) - perfScore(a))[0];
	const highestCredits = [...allStudents].sort((a, b) => b.creditsEarned - a.creditsEarned)[0];
	const mostActive = [...allStudents].sort((a, b) => b.activityCount - a.activityCount)[0];
	const pendingAttention = allStudents.filter(
		(s) => s.status === 'At Risk' || s.status === 'Pending Review'
	);

	// ── Table State ────────────────────────────────────────────────────────────
	let searchQuery = $state('');
	let filterStatus = $state<Status | 'All'>('All');
	let filterDept = $state('All');
	let currentPage = $state(1);
	const pageSize = 10;
	let showFilters = $state(false);

	const departments = ['All', ...Array.from(new Set(allStudents.map((s) => s.department)))];

	const filteredStudents = $derived(
		allStudents.filter((s) => {
			const matchSearch =
				searchQuery === '' ||
				s.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				s.regNo.toLowerCase().includes(searchQuery.toLowerCase()) ||
				s.department.toLowerCase().includes(searchQuery.toLowerCase());
			const matchStatus = filterStatus === 'All' || s.status === filterStatus;
			const matchDept = filterDept === 'All' || s.department === filterDept;
			return matchSearch && matchStatus && matchDept;
		})
	);

	const totalPages = $derived(Math.max(1, Math.ceil(filteredStudents.length / pageSize)));

	const pagedStudents = $derived(
		filteredStudents.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	function resetPage() {
		currentPage = 1;
	}

	// ── Modal ──────────────────────────────────────────────────────────────────
	let activeStudent = $state<Student | null>(null);
	let isModalOpen = $state(false);

	function openStudentModal(student: Student) {
		activeStudent = student;
		isModalOpen = true;
	}

	function closeModal() {
		isModalOpen = false;
		activeStudent = null;
	}

	// ── Student Detail View ──────────────────────────────────────────────────────
	let detailStudent = $state<Student | null>(null);

	// Rank by credits earned within the cohort (1 = highest)
	function studentRank(student: Student): number {
		return allStudents.filter((s) => s.creditsEarned > student.creditsEarned).length + 1;
	}

	function openStudentDetail(student: Student) {
		detailStudent = student;
		isModalOpen = false;
		activeStudent = null;
	}

	function closeStudentDetail() {
		detailStudent = null;
	}

	// ── Toast ──────────────────────────────────────────────────────────────────
	interface Toast {
		id: number;
		message: string;
		type: 'success' | 'danger';
	}
	let toasts = $state<Toast[]>([]);
	let nextToastId = 0;

	function triggerToast(message: string, type: 'success' | 'danger' = 'success') {
		const id = nextToastId++;
		toasts = [...toasts, { id, message, type }];
		setTimeout(() => {
			toasts = toasts.filter((t) => t.id !== id);
		}, 3000);
	}

	// ── Helpers ────────────────────────────────────────────────────────────────
	function statusClass(status: Status): string {
		switch (status) {
			case 'Active':
				return 'bg-emerald-50 text-emerald-700 border border-emerald-100';
			case 'At Risk':
				return 'bg-rose-50 text-rose-700 border border-rose-100';
			case 'Pending Review':
				return 'bg-amber-50 text-amber-700 border border-amber-100';
			case 'Inactive':
				return 'bg-slate-100 text-slate-500 border border-slate-200';
		}
	}

	function statusDot(status: Status): string {
		switch (status) {
			case 'Active':
				return 'bg-emerald-500';
			case 'At Risk':
				return 'bg-rose-500';
			case 'Pending Review':
				return 'bg-amber-500';
			case 'Inactive':
				return 'bg-slate-400';
		}
	}

	function initials(name: string): string {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.slice(0, 2)
			.toUpperCase();
	}
</script>

<!-- ── Toast Container ─────────────────────────────────────────────────────── -->
<div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-sm pointer-events-none">
	{#each toasts as toast (toast.id)}
		<div
			transition:slide={{ duration: 150 }}
			class="p-4 rounded-xl border shadow-lg flex items-center gap-3 text-xs font-semibold font-sans pointer-events-auto {toast.type ===
			'success'
				? 'bg-emerald-50 border-emerald-200 text-emerald-800'
				: 'bg-rose-50 border-rose-200 text-rose-800'}"
		>
			{#if toast.type === 'success'}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2.5"
					stroke="currentColor"
					class="w-4 h-4 text-emerald-600 shrink-0"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
				</svg>
			{:else}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2.5"
					stroke="currentColor"
					class="w-4 h-4 text-rose-600 shrink-0"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
					/>
				</svg>
			{/if}
			<span>{toast.message}</span>
		</div>
	{/each}
</div>

{#if detailStudent}
	<AdminStudentDetailView
		student={detailStudent}
		rank={studentRank(detailStudent)}
		cohortSize={totalStudents}
		onBack={closeStudentDetail}
		onToast={triggerToast}
	/>
{:else}
	<!-- ── Stat Cards ──────────────────────────────────────────────────────────── -->
	<section class="grid grid-cols-2 lg:grid-cols-4 gap-4">
		<!-- Total Students -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between shadow-xs hover:shadow-md transition-shadow duration-200"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{totalStudents}</span>
				<div class="p-2.5 rounded-lg bg-blue-50 text-blue-600 border border-blue-100">
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
							d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.109A11.386 11.386 0 0 1 10.089 21c-2.316 0-4.445-.69-6.22-1.879v-.003a4.125 4.125 0 0 1 7.533-2.493M15 19.128v-.003c0-1.112-.285-2.16-.786-3.07M14.214 16.058A9.396 9.396 0 0 0 10.089 15c-1.47 0-2.854.34-4.082.945M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-xs font-bold text-slate-800 tracking-wide">Total Students</h3>
				<p class="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-wider">
					+2 this semester
				</p>
			</div>
		</div>

		<!-- Active Students -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between shadow-xs hover:shadow-md transition-shadow duration-200"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{activeStudents}</span>
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
							d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-xs font-bold text-slate-800 tracking-wide">Active Students</h3>
				<p class="text-[10px] font-bold text-emerald-500 mt-1 uppercase tracking-wider">
					{activeStudents} on engagement
				</p>
			</div>
		</div>

		<!-- Pending Certificate Reviews -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between shadow-xs hover:shadow-md transition-shadow duration-200"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{pendingCertReviews}</span>
				<div class="p-2.5 rounded-lg bg-amber-50 text-amber-600 border border-amber-100">
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
							d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-xs font-bold text-slate-800 tracking-wide">Pending Certificate Reviews</h3>
				<p class="text-[10px] font-bold text-amber-500 mt-1 uppercase tracking-wider">
					3 marked urgent
				</p>
			</div>
		</div>

		<!-- Average Credits Earned -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between shadow-xs hover:shadow-md transition-shadow duration-200"
		>
			<div class="flex items-center justify-between">
				<span class="text-2xl font-bold font-serif text-slate-900">{avgCredits}</span>
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
							d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 0 1 5.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-4">
				<h3 class="text-xs font-bold text-slate-800 tracking-wide">Average Credits Earned</h3>
				<p class="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-wider">
					57.5% avg from last batch
				</p>
			</div>
		</div>
	</section>

	<!-- ── Student Overview + Quick Insights ──────────────────────────────────── -->
	<section class="grid grid-cols-1 lg:grid-cols-12 gap-5 items-stretch">
		<!-- Student Overview -->
		<div
			class="lg:col-span-8 bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden flex flex-col"
		>
			<div class="p-5 border-b border-slate-100 bg-slate-50/20">
				<h2 class="text-sm font-bold font-serif text-inst-navy">Student Overview</h2>
				<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
					Key performance profiles of enrolled students
				</p>
			</div>
			<div class="p-5 grid grid-cols-1 sm:grid-cols-2 gap-4 flex-grow auto-rows-fr">
				<!-- Top Performer -->
				<div
					class="bg-slate-50 rounded-xl border border-slate-150 p-4 flex items-center gap-3 hover:shadow-sm transition-shadow"
				>
					<div
						class="w-9 h-9 rounded-lg bg-yellow-100 text-yellow-600 flex items-center justify-center border border-yellow-200 shrink-0"
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
								d="M11.48 3.499a.562.562 0 0 1 1.04 0l2.125 5.111a.563.563 0 0 0 .475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 0 0-.182.557l1.285 5.385a.562.562 0 0 1-.84.61l-4.725-2.885a.562.562 0 0 0-.586 0L6.982 20.54a.562.562 0 0 1-.84-.61l1.285-5.386a.562.562 0 0 0-.182-.557l-4.204-3.602a.562.562 0 0 1 .321-.988l5.518-.442a.563.563 0 0 0 .475-.345L11.48 3.5Z"
							/>
						</svg>
					</div>
					<div class="flex flex-col gap-0.5 min-w-0">
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider"
							>Top Performer</span
						>
						<span class="font-bold text-sm text-slate-900 truncate">{topPerformer.name}</span>
						<span class="text-[10px] text-slate-500 font-semibold">
							{topPerformer.creditsEarned} credits · {topPerformer.department}
						</span>
					</div>
				</div>

				<!-- Highest Credits Earned -->
				<div
					class="bg-slate-50 rounded-xl border border-slate-150 p-4 flex items-center gap-3 hover:shadow-sm transition-shadow"
				>
					<div
						class="w-9 h-9 rounded-lg bg-emerald-100 text-emerald-600 flex items-center justify-center border border-emerald-200 shrink-0"
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
								d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 0 1 5.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
							/>
						</svg>
					</div>
					<div class="flex flex-col gap-0.5 min-w-0">
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider"
							>Highest Credits Earned</span
						>
						<span class="font-bold text-sm text-slate-900 truncate">{highestCredits.name}</span>
						<span class="text-[10px] text-slate-500 font-semibold">
							{highestCredits.creditsEarned} credits earned this batch
						</span>
					</div>
				</div>

				<!-- Most Active Student -->
				<div
					class="bg-slate-50 rounded-xl border border-slate-150 p-4 flex items-center gap-3 hover:shadow-sm transition-shadow"
				>
					<div
						class="w-9 h-9 rounded-lg bg-blue-100 text-blue-600 flex items-center justify-center border border-blue-200 shrink-0"
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
								d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z"
							/>
						</svg>
					</div>
					<div class="flex flex-col gap-0.5 min-w-0">
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider"
							>Most Active Student</span
						>
						<span class="font-bold text-sm text-slate-900 truncate">{mostActive.name}</span>
						<span class="text-[10px] text-slate-500 font-semibold">
							{mostActive.activityCount} activities logged
						</span>
					</div>
				</div>

				<!-- Requiring Attention -->
				<div
					class="bg-rose-50/60 rounded-xl border border-rose-100 p-4 flex items-center gap-3 hover:shadow-sm transition-shadow"
				>
					<div
						class="w-9 h-9 rounded-lg bg-rose-100 text-rose-600 flex items-center justify-center border border-rose-200 shrink-0"
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
								d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
							/>
						</svg>
					</div>
					<div class="flex flex-col gap-0.5 min-w-0">
						<span class="text-[9px] font-bold text-rose-400 uppercase tracking-wider"
							>Requiring Attention</span
						>
						<span class="font-bold text-sm text-rose-700">{pendingAttention.length} students</span>
						<span class="text-[10px] text-rose-500 font-semibold">
							{pendingAttention
								.slice(0, 4)
								.map((s) => s.name.split(' ')[0])
								.join(', ')}
						</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Quick Insights -->
		<div
			class="lg:col-span-4 bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden"
		>
			<div class="p-5 border-b border-slate-100 bg-slate-50/20">
				<h2 class="text-sm font-bold font-serif text-inst-navy">Quick Insights</h2>
				<p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
					Administrative actions
				</p>
			</div>
			<div class="p-4 space-y-3">
				<!-- Pending certificates -->
				<div
					class="flex items-center justify-between p-3 rounded-lg hover:bg-slate-50 border border-slate-100 transition-colors cursor-pointer"
				>
					<div class="flex items-center gap-2.5">
						<div
							class="w-7 h-7 rounded-lg bg-amber-50 text-amber-600 border border-amber-100 flex items-center justify-center"
						>
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
									d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m6.75 12H9m1.5-12H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
								/>
							</svg>
						</div>
						<div>
							<p class="text-[11px] font-bold text-slate-700">Pending certificates</p>
							<p class="text-[9px] text-slate-400 font-semibold">Need review</p>
						</div>
					</div>
					<span
						class="text-xs font-extrabold text-amber-600 bg-amber-50 border border-amber-100 px-2 py-0.5 rounded-md"
						>{pendingCertReviews}</span
					>
				</div>

				<!-- Review credit target -->
				<div
					class="flex items-center justify-between p-3 rounded-lg hover:bg-slate-50 border border-slate-100 transition-colors cursor-pointer"
				>
					<div class="flex items-center gap-2.5">
						<div
							class="w-7 h-7 rounded-lg bg-blue-50 text-blue-600 border border-blue-100 flex items-center justify-center"
						>
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
									d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 0 1 5.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
								/>
							</svg>
						</div>
						<div>
							<p class="text-[11px] font-bold text-slate-700">Review credit target</p>
							<p class="text-[9px] text-slate-400 font-semibold">Below threshold</p>
						</div>
					</div>
					<span
						class="text-xs font-extrabold text-blue-600 bg-blue-50 border border-blue-100 px-2 py-0.5 rounded-md"
						>4</span
					>
				</div>

				<!-- Inactive students -->
				<div
					class="flex items-center justify-between p-3 rounded-lg hover:bg-slate-50 border border-slate-100 transition-colors cursor-pointer"
				>
					<div class="flex items-center gap-2.5">
						<div
							class="w-7 h-7 rounded-lg bg-slate-100 text-slate-500 border border-slate-200 flex items-center justify-center"
						>
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
									d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z"
								/>
							</svg>
						</div>
						<div>
							<p class="text-[11px] font-bold text-slate-700">Inactive students</p>
							<p class="text-[9px] text-slate-400 font-semibold">30 days</p>
						</div>
					</div>
					<span
						class="text-xs font-extrabold text-slate-600 bg-slate-100 border border-slate-200 px-2 py-0.5 rounded-md"
						>3</span
					>
				</div>

				<!-- Pending review -->
				<div
					class="flex items-center justify-between p-3 rounded-lg hover:bg-slate-50 border border-slate-100 transition-colors cursor-pointer"
				>
					<div class="flex items-center gap-2.5">
						<div
							class="w-7 h-7 rounded-lg bg-rose-50 text-rose-600 border border-rose-100 flex items-center justify-center"
						>
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
									d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
								/>
							</svg>
						</div>
						<div>
							<p class="text-[11px] font-bold text-slate-700">Pending review</p>
							<p class="text-[9px] text-slate-400 font-semibold">Action required</p>
						</div>
					</div>
					<span
						class="text-xs font-extrabold text-rose-600 bg-rose-50 border border-rose-100 px-2 py-0.5 rounded-md"
						>{pendingAttention.length}</span
					>
				</div>

				<div class="pt-1">
					<button
						onclick={() => {
							filterStatus = 'At Risk';
							resetPage();
						}}
						class="w-full py-2 text-[11px] font-bold text-[#881B1B] border border-[#881B1B]/20 bg-[#881B1B]/5 hover:bg-[#881B1B]/10 rounded-lg transition-colors tracking-wide uppercase"
					>
						View All Flagged Students
					</button>
				</div>
			</div>
		</div>
	</section>

	<!-- ── Student Management Table ───────────────────────────────────────────── -->
	<section class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
		<!-- Table Header -->
		<div
			class="p-5 border-b border-slate-100 flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-slate-50/20"
		>
			<h2 class="text-sm font-bold font-serif text-inst-navy">Student Management</h2>
			<div
				class="flex items-center gap-2 text-[10px] font-extrabold text-slate-400 uppercase tracking-wider"
			>
				Total Students: <span class="text-slate-700 ml-1">{filteredStudents.length}</span>
			</div>
		</div>

		<!-- Search & Filter Bar -->
		<div class="px-5 py-3.5 border-b border-slate-100 flex flex-wrap items-center gap-3">
			<!-- Search -->
			<div class="relative flex items-center">
				<span class="absolute left-3 text-slate-400 pointer-events-none">
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
							d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
						/>
					</svg>
				</span>
				<input
					type="text"
					placeholder="Search students..."
					bind:value={searchQuery}
					oninput={resetPage}
					class="pl-8 pr-4 py-2 bg-slate-50 rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-slate-350 focus:bg-white w-48 transition-all focus:w-56"
				/>
			</div>

			<!-- More Filters toggle -->
			<button
				onclick={() => (showFilters = !showFilters)}
				class="flex items-center gap-1.5 px-3 py-2 text-[11px] font-bold text-slate-500 border border-slate-200 rounded-lg hover:bg-slate-50 transition-colors"
			>
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
						d="M10.5 6h9.75M10.5 6a1.5 1.5 0 1 1-3 0m3 0a1.5 1.5 0 1 0-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m-9.75 0h9.75"
					/>
				</svg>
				More Filters
			</button>

			{#if showFilters}
				<div transition:slide={{ duration: 150 }} class="flex flex-wrap items-center gap-2">
					<select
						bind:value={filterStatus}
						onchange={resetPage}
						class="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-[11px] font-bold text-slate-600 focus:outline-none focus:border-slate-350"
					>
						<option value="All">All Status</option>
						<option value="Active">Active</option>
						<option value="At Risk">At Risk</option>
						<option value="Pending Review">Pending Review</option>
						<option value="Inactive">Inactive</option>
					</select>

					<select
						bind:value={filterDept}
						onchange={resetPage}
						class="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-[11px] font-bold text-slate-600 focus:outline-none focus:border-slate-350"
					>
						{#each departments as dept}
							<option value={dept}>{dept === 'All' ? 'All Departments' : dept}</option>
						{/each}
					</select>

					{#if filterStatus !== 'All' || filterDept !== 'All' || searchQuery !== ''}
						<button
							onclick={() => {
								filterStatus = 'All';
								filterDept = 'All';
								searchQuery = '';
								resetPage();
							}}
							class="px-3 py-2 text-[11px] font-bold text-rose-600 border border-rose-200 bg-rose-50 rounded-lg hover:bg-rose-100 transition-colors"
						>
							Clear
						</button>
					{/if}
				</div>
			{/if}
		</div>

		<!-- Table -->
		<div class="overflow-x-auto">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr
						class="border-b border-slate-100 bg-slate-50/50 text-[10px] font-extrabold text-slate-400 uppercase tracking-wider"
					>
						<th class="py-3 px-5">Name</th>
						<th class="py-3 px-5">Department</th>
						<th class="py-3 px-5">Sem</th>
						<th class="py-3 px-5">Credits Earned</th>
						<th class="py-3 px-5">Certificates</th>
						<th class="py-3 px-5">Activity Count</th>
						<th class="py-3 px-5">Status</th>
						<th class="py-3 px-5 text-center">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-xs font-sans">
					{#if pagedStudents.length === 0}
						<tr>
							<td colspan="8" class="py-16 text-center text-slate-400 font-semibold text-xs">
								No students match your search criteria.
							</td>
						</tr>
					{:else}
						{#each pagedStudents as student (student.id)}
							<tr class="hover:bg-slate-50/50 transition-colors">
								<!-- Name -->
								<td class="py-3.5 px-5">
									<div class="font-bold text-slate-800">{student.name}</div>
									<div class="text-[10px] text-slate-400 font-semibold uppercase">
										{student.regNo}
									</div>
								</td>
								<!-- Department -->
								<td class="py-3.5 px-5 text-slate-600 font-semibold">{student.department}</td>
								<!-- Semester -->
								<td class="py-3.5 px-5 text-slate-600 font-extrabold">{student.semester}th</td>
								<!-- Credits Earned -->
								<td class="py-3.5 px-5">
									<div class="flex items-baseline gap-1">
										<span class="font-extrabold text-slate-800">{student.creditsEarned}</span>
										<span class="text-[9px] text-slate-400 font-bold">/{student.creditsTarget}</span
										>
									</div>
									<div class="mt-1.5 h-1 w-20 bg-slate-100 rounded-full overflow-hidden">
										<div
											class="h-full rounded-full {student.creditsEarned >= 150
												? 'bg-emerald-400'
												: student.creditsEarned >= 100
													? 'bg-amber-400'
													: 'bg-rose-400'}"
											style="width: {Math.min(
												100,
												(student.creditsEarned / student.creditsTarget) * 100
											)}%"
										></div>
									</div>
								</td>
								<!-- Certificates -->
								<td class="py-3.5 px-5">
									<span class="font-extrabold text-[#881B1B]">{student.certificates}</span>
								</td>
								<!-- Activity Count -->
								<td class="py-3.5 px-5 font-extrabold text-slate-700">{student.activityCount}</td>
								<!-- Status -->
								<td class="py-3.5 px-5">
									<span
										class="inline-flex items-center gap-1.5 px-2 py-1 rounded-md text-[10px] font-extrabold uppercase tracking-wide {statusClass(
											student.status
										)}"
									>
										<span class="w-1.5 h-1.5 rounded-full {statusDot(student.status)}"></span>
										{student.status}
									</span>
								</td>
								<!-- Actions -->
								<td class="py-3.5 px-5">
									<div class="flex items-center justify-center">
										<button
											onclick={() => openStudentModal(student)}
											aria-label="View student"
											class="p-1.5 border border-slate-200 text-slate-500 hover:text-inst-navy hover:bg-slate-100 rounded-lg transition-colors focus:outline-none"
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
													d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
												/>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
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

		<!-- Table Footer: info + pagination -->
		<div
			class="px-5 py-3.5 border-t border-slate-100 flex flex-col sm:flex-row items-center justify-between gap-3 bg-slate-50/20"
		>
			<span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">
				Showing {filteredStudents.length === 0
					? 0
					: Math.min((currentPage - 1) * pageSize + 1, filteredStudents.length)}–{Math.min(
					currentPage * pageSize,
					filteredStudents.length
				)} of {filteredStudents.length} students
			</span>

			<div class="flex items-center gap-1.5">
				<button
					onclick={() => {
						if (currentPage > 1) currentPage -= 1;
					}}
					disabled={currentPage === 1}
					class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-500 hover:bg-slate-100 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
					aria-label="Previous page"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2.5"
						stroke="currentColor"
						class="w-3.5 h-3.5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
					</svg>
				</button>

				{#each Array.from({ length: totalPages }, (_, i) => i + 1) as page}
					<button
						onclick={() => (currentPage = page)}
						class="w-8 h-8 flex items-center justify-center rounded-lg border text-xs font-extrabold transition-colors {currentPage ===
						page
							? 'bg-[#881B1B] text-white border-[#881B1B]'
							: 'border-slate-200 text-slate-600 hover:bg-slate-100'}"
					>
						{page}
					</button>
				{/each}

				<button
					onclick={() => {
						if (currentPage < totalPages) currentPage += 1;
					}}
					disabled={currentPage === totalPages}
					class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-500 hover:bg-slate-100 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
					aria-label="Next page"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2.5"
						stroke="currentColor"
						class="w-3.5 h-3.5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
					</svg>
				</button>
			</div>
		</div>
	</section>
{/if}

<!-- ── Student Detail Modal ───────────────────────────────────────────────── -->
{#if isModalOpen && activeStudent}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
		onclick={closeModal}
	>
		<div
			onclick={(e) => e.stopPropagation()}
			class="w-full max-w-lg bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
		>
			<!-- Modal Header -->
			<div class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30">
				<div>
					<h3 class="text-sm font-bold font-serif text-inst-navy">Student Profile</h3>
					<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
						ID: {activeStudent.id} · Batch {activeStudent.batch}
					</p>
				</div>
				<button
					onclick={closeModal}
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

			<!-- Modal Body -->
			<div class="p-6 space-y-5 overflow-y-auto max-h-[65vh]">
				<!-- Avatar + Name -->
				<div class="flex items-center gap-4">
					<div
						class="w-14 h-14 rounded-full bg-[#881B1B] text-white flex items-center justify-center font-bold text-lg border-2 border-white shadow-md font-serif shrink-0"
					>
						{initials(activeStudent.name)}
					</div>
					<div class="flex-grow">
						<div class="font-bold text-lg text-slate-900 font-serif">{activeStudent.name}</div>
						<div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
							{activeStudent.regNo}
						</div>
						<div class="text-[10px] text-slate-500 font-semibold">{activeStudent.email}</div>
					</div>
					<span
						class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-[10px] font-extrabold uppercase tracking-wide {statusClass(
							activeStudent.status
						)}"
					>
						<span class="w-1.5 h-1.5 rounded-full {statusDot(activeStudent.status)}"></span>
						{activeStudent.status}
					</span>
				</div>

				<!-- Info Grid -->
				<div class="grid grid-cols-2 gap-4 bg-slate-50 p-4 rounded-xl border border-slate-150">
					<div>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Department</span
						>
						<span class="text-xs font-bold text-slate-800 block mt-0.5"
							>{activeStudent.department}</span
						>
					</div>
					<div>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Semester</span
						>
						<span class="text-xs font-bold text-slate-800 block mt-0.5"
							>{activeStudent.semester}th Semester</span
						>
					</div>
					<div>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Certificates</span
						>
						<span class="text-xs font-bold text-[#881B1B] block mt-0.5"
							>{activeStudent.certificates} uploaded</span
						>
					</div>
					<div>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Activities</span
						>
						<span class="text-xs font-bold text-slate-800 block mt-0.5"
							>{activeStudent.activityCount} logged</span
						>
					</div>
				</div>

				<!-- Credits Progress -->
				<div class="space-y-2">
					<div class="flex items-center justify-between">
						<span class="text-[10px] font-bold text-slate-500 uppercase tracking-wider"
							>Credit Progress</span
						>
						<span class="text-[10px] font-extrabold text-slate-700">
							{activeStudent.creditsEarned} / {activeStudent.creditsTarget}
						</span>
					</div>
					<div class="h-2.5 w-full bg-slate-100 rounded-full overflow-hidden">
						<div
							class="h-full rounded-full transition-all duration-500 {activeStudent.creditsEarned >=
							150
								? 'bg-emerald-500'
								: activeStudent.creditsEarned >= 100
									? 'bg-amber-500'
									: 'bg-rose-500'}"
							style="width: {Math.min(
								100,
								(activeStudent.creditsEarned / activeStudent.creditsTarget) * 100
							)}%"
						></div>
					</div>
					<div class="flex justify-between text-[9px] font-bold text-slate-400 px-0.5">
						<span>0</span>
						<span>50</span>
						<span>100</span>
						<span>150</span>
						<span>200 (Target)</span>
					</div>
				</div>
			</div>

			<!-- Modal Footer -->
			<div class="p-4 border-t border-slate-100 flex items-center justify-end gap-3 bg-slate-50/20">
				<button
					onclick={closeModal}
					class="px-4 py-2 text-xs font-bold text-slate-600 border border-slate-200 rounded-lg hover:bg-slate-100 transition-colors focus:outline-none"
				>
					Close
				</button>
				<button
					onclick={() => {
						triggerToast(`Sending notice to ${activeStudent?.name}...`);
						closeModal();
					}}
					class="inline-flex items-center gap-1.5 px-4 py-2 text-xs font-bold text-white bg-[#881B1B] hover:bg-[#881B1B]/90 rounded-lg transition-colors focus:outline-none shadow-xs"
				>
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
							d="M14.857 17.082a23.848 23.848 0 0 0 5.454-1.31A8.967 8.967 0 0 1 18 9.75V9A6 6 0 0 0 6 9v.75a8.967 8.967 0 0 1-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 0 1-5.714 0m5.714 0a3 3 0 1 1-5.714 0"
						/>
					</svg>
					Send Notice
				</button>
				<button
					onclick={() => activeStudent && openStudentDetail(activeStudent)}
					class="inline-flex items-center gap-1.5 px-4 py-2 text-xs font-bold text-white bg-[#881B1B] hover:bg-[#881B1B]/90 rounded-lg transition-colors focus:outline-none shadow-xs"
				>
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
							d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"
						/>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
						/>
					</svg>
					View Details
				</button>
			</div>
		</div>
	</div>
{/if}
