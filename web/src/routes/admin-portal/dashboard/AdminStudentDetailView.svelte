<script lang="ts">
	import { fade, slide } from 'svelte/transition';

	// ── Types (structural match with AdminStudentManagementView's Student) ──────
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

	// ── Props ───────────────────────────────────────────────────────────────────
	let {
		student,
		rank,
		cohortSize,
		onBack,
		onToast
	}: {
		student: Student;
		rank: number;
		cohortSize: number;
		onBack: () => void;
		onToast?: (message: string, type?: 'success' | 'danger') => void;
	} = $props();

	function toast(message: string, type: 'success' | 'danger' = 'success') {
		onToast?.(message, type);
	}

	// ── Derived profile stats ────────────────────────────────────────────────────
	const creditPct = $derived(
		Math.min(100, Math.round((student.creditsEarned / student.creditsTarget) * 100))
	);
	const creditsRemaining = $derived(Math.max(0, student.creditsTarget - student.creditsEarned));
	const participationScore = $derived(Math.min(98, 50 + student.activityCount * 2));
	const pendingCerts = 2;

	// Static mock meta (not present in Student, shown for demo fidelity)
	const course = 'B.Tech';
	const track = 'Skill Building';

	// ── Mock Activity History ────────────────────────────────────────────────────
	type ActivityStatus = 'Completed' | 'Pending' | 'Rejected';
	interface Activity {
		id: number;
		name: string;
		category: string;
		date: string;
		credits: number;
		status: ActivityStatus;
	}

	const activities: Activity[] = [
		{
			id: 1,
			name: 'NPTEL Online Certification – Python',
			category: 'Technical',
			date: '15 Jun 2025',
			credits: 20,
			status: 'Completed'
		},
		{
			id: 2,
			name: 'National Science Day Seminar',
			category: 'Academic',
			date: '28 Feb 2025',
			credits: 8,
			status: 'Completed'
		},
		{
			id: 3,
			name: 'Leadership Development Workshop',
			category: 'Leadership',
			date: '10 May 2024',
			credits: 12,
			status: 'Completed'
		},
		{
			id: 4,
			name: 'Internship – Tech Startup',
			category: 'Professional',
			date: '22 Apr 2025',
			credits: 15,
			status: 'Completed'
		},
		{
			id: 5,
			name: 'Community Service – Beach Cleanup',
			category: 'Social',
			date: '5 Mar 2025',
			credits: 8,
			status: 'Completed'
		},
		{
			id: 6,
			name: 'Research Project – AI Ethics',
			category: 'Research',
			date: '18 Jan 2025',
			credits: 10,
			status: 'Pending'
		},
		{
			id: 7,
			name: 'Sports Event Coordination',
			category: 'Sports',
			date: '30 Dec 2024',
			credits: 6,
			status: 'Completed'
		},
		{
			id: 8,
			name: 'Cultural Fest Organization',
			category: 'Cultural',
			date: '20 Nov 2025',
			credits: 8,
			status: 'Rejected'
		}
	];

	// ── Mock Certificate History ─────────────────────────────────────────────────
	interface Certificate {
		id: number;
		name: string;
		issuer: string;
		date: string;
		credits: number;
		status: ActivityStatus;
	}

	const certificates: Certificate[] = [
		{
			id: 1,
			name: 'Python for Data Science',
			issuer: 'NPTEL',
			date: '15 Jun 2025',
			credits: 20,
			status: 'Completed'
		},
		{
			id: 2,
			name: 'Full Stack Web Development',
			issuer: 'Coursera',
			date: '2 Apr 2025',
			credits: 15,
			status: 'Completed'
		},
		{
			id: 3,
			name: 'Machine Learning Specialization',
			issuer: 'DeepLearning.AI',
			date: '18 Jan 2025',
			credits: 18,
			status: 'Pending'
		},
		{
			id: 4,
			name: 'Cloud Practitioner Essentials',
			issuer: 'AWS',
			date: '9 Nov 2024',
			credits: 12,
			status: 'Completed'
		}
	];

	// ── Table State ──────────────────────────────────────────────────────────────
	let activeTab = $state<'activity' | 'certificate'>('activity');
	let searchQuery = $state('');
	let filterCategory = $state('All');
	let filterStatus = $state<'All' | ActivityStatus>('All');

	const categories = ['All', ...Array.from(new Set(activities.map((a) => a.category)))];

	const filteredActivities = $derived(
		activities.filter((a) => {
			const matchSearch =
				searchQuery === '' ||
				a.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				a.category.toLowerCase().includes(searchQuery.toLowerCase());
			const matchCat = filterCategory === 'All' || a.category === filterCategory;
			const matchStatus = filterStatus === 'All' || a.status === filterStatus;
			return matchSearch && matchCat && matchStatus;
		})
	);

	// ── Mentor Observations ──────────────────────────────────────────────────────
	interface Note {
		id: number;
		author: string;
		role: string;
		badge: string;
		text: string;
	}

	let notes = $state<Note[]>([
		{
			id: 1,
			author: 'Dr. Rajesh Kumar',
			role: 'CSE Dept · Semester 6',
			badge: 'Academic Advisor',
			text: `Student actively participates in technical activities and consistently submits certificates on time. Academic performance in extracurricular domain is commendable — particularly strong engagement with NPTEL certifications and research-oriented activities. Leadership activity participation can be improved; recommend encouraging enrolment in upcoming Student Council or Workshop facilitation programmes. Overall, this student demonstrates initiative and discipline appropriate for the Skill Building track.`
		}
	]);

	let nextNoteId = 2;
	let composerOpen = $state(false);
	let composerText = $state('');
	let editingId = $state<number | null>(null);

	function openAddNote() {
		editingId = null;
		composerText = '';
		composerOpen = true;
	}

	function openEditNote(note: Note) {
		editingId = note.id;
		composerText = note.text;
		composerOpen = true;
	}

	function saveNote() {
		const text = composerText.trim();
		if (text === '') return;
		if (editingId !== null) {
			notes = notes.map((n) => (n.id === editingId ? { ...n, text } : n));
			toast('Observation updated.');
		} else {
			notes = [
				...notes,
				{
					id: nextNoteId++,
					author: 'Dr. Rajesh Kumar',
					role: 'CSE Dept · Semester 6',
					badge: 'Academic Advisor',
					text
				}
			];
			toast('Observation added.');
		}
		composerOpen = false;
		composerText = '';
		editingId = null;
	}

	// ── Helpers ──────────────────────────────────────────────────────────────────
	function initials(name: string): string {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.slice(0, 2)
			.toUpperCase();
	}

	function statusChip(status: ActivityStatus): string {
		switch (status) {
			case 'Completed':
				return 'bg-emerald-50 text-emerald-700 border border-emerald-100';
			case 'Pending':
				return 'bg-amber-50 text-amber-700 border border-amber-100';
			case 'Rejected':
				return 'bg-rose-50 text-rose-700 border border-rose-100';
		}
	}

	function statusDot(status: ActivityStatus): string {
		switch (status) {
			case 'Completed':
				return 'bg-emerald-500';
			case 'Pending':
				return 'bg-amber-500';
			case 'Rejected':
				return 'bg-rose-500';
		}
	}

	function studentStatusChip(status: Status): string {
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
</script>

<div class="space-y-6" transition:fade={{ duration: 150 }}>
	<!-- ── Back link ─────────────────────────────────────────────────────────── -->
	<button
		onclick={onBack}
		class="inline-flex items-center gap-1.5 text-[11px] font-bold text-slate-500 hover:text-inst-navy transition-colors uppercase tracking-wider"
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
		Back to Students
	</button>

	<!-- ── Page Heading ──────────────────────────────────────────────────────── -->
	<div class="flex flex-col gap-1">
		<h1 class="text-2xl font-bold font-serif text-slate-900">Student Details</h1>
		<p class="text-xs text-slate-400 font-semibold tracking-wide">
			View complete student extracurricular profile and performance.
		</p>
	</div>

	<!-- ── Profile Header Card ───────────────────────────────────────────────── -->
	<section class="bg-white border border-slate-200 rounded-xl shadow-xs p-5 sm:p-6">
		<div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-5">
			<!-- Identity -->
			<div class="flex items-start gap-4 min-w-0">
				<div
					class="w-16 h-16 rounded-full bg-[#881B1B] text-white flex items-center justify-center font-bold text-xl border-2 border-white shadow-md font-serif shrink-0"
				>
					{initials(student.name)}
				</div>
				<div class="min-w-0">
					<div class="flex items-center gap-2.5 flex-wrap">
						<h2 class="text-lg font-bold font-serif text-slate-900">{student.name}</h2>
						<span
							class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md text-[10px] font-extrabold uppercase tracking-wide {studentStatusChip(
								student.status
							)}"
						>
							<span class="w-1.5 h-1.5 rounded-full bg-current opacity-70"></span>
							{student.status}
						</span>
					</div>
					<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
						{student.regNo}
					</p>

					<!-- Meta grid -->
					<div class="grid grid-cols-2 sm:grid-cols-4 gap-x-8 gap-y-3 mt-4">
						<div>
							<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
								>Department</span
							>
							<span class="text-xs font-bold text-slate-800 block mt-0.5">{student.department}</span
							>
						</div>
						<div>
							<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
								>Course</span
							>
							<span class="text-xs font-bold text-slate-800 block mt-0.5">{course}</span>
						</div>
						<div>
							<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
								>Semester</span
							>
							<span class="text-xs font-bold text-slate-800 block mt-0.5"
								>{student.semester}th Semester</span
							>
						</div>
						<div>
							<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
								>Track</span
							>
							<span class="text-xs font-bold text-slate-800 block mt-0.5">{track}</span>
						</div>
					</div>
				</div>
			</div>

			<!-- Actions -->
			<div class="flex flex-wrap items-center gap-2 shrink-0">
				<button
					onclick={() => (activeTab = 'activity')}
					class="inline-flex items-center gap-1.5 px-3 py-2 text-[11px] font-bold text-slate-600 border border-slate-200 rounded-lg hover:bg-slate-50 transition-colors"
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
							d="M3.75 3v11.25A2.25 2.25 0 0 0 6 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0 1 18 16.5h-2.25m-7.5 0h7.5m-7.5 0-1 3m8.5-3 1 3m0 0 .5 1.5m-.5-1.5h-9.5m0 0-.5 1.5M9 11.25v1.5M12 9v3.75m3-6v6"
						/>
					</svg>
					View Activities
				</button>
				<button
					onclick={() => (activeTab = 'certificate')}
					class="inline-flex items-center gap-1.5 px-3 py-2 text-[11px] font-bold text-slate-600 border border-slate-200 rounded-lg hover:bg-slate-50 transition-colors"
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
							d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M7.73 9.728a6.726 6.726 0 0 0 2.748 1.35m8.272-6.842V4.5c0 2.108-.966 3.99-2.48 5.228m2.48-5.492a46.32 46.32 0 0 1 2.916.52 6.003 6.003 0 0 1-5.395 4.972m0 0a6.726 6.726 0 0 1-2.749 1.35m0 0a6.772 6.772 0 0 1-3.044 0"
						/>
					</svg>
					View Certificates
				</button>
				<button
					onclick={() => toast(`Generating report for ${student.name}...`)}
					class="inline-flex items-center gap-1.5 px-3 py-2 text-[11px] font-bold text-white bg-[#881B1B] hover:bg-[#881B1B]/90 rounded-lg transition-colors shadow-xs"
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
							d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5m-13.5-9L12 3m0 0 4.5 4.5M12 3v13.5"
						/>
					</svg>
					Download Report
				</button>
			</div>
		</div>
	</section>

	<!-- ── Stat Cards ────────────────────────────────────────────────────────── -->
	<section class="grid grid-cols-2 lg:grid-cols-4 gap-4">
		<!-- Total Credits Earned -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<div class="p-2 rounded-lg bg-rose-50 text-rose-600 border border-rose-100">
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
							d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-3">
				<span class="text-2xl font-bold font-serif text-slate-900">{student.creditsEarned}</span>
				<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1">Total Credits Earned</h3>
				<p class="text-[10px] font-bold text-slate-400 mt-0.5 uppercase tracking-wider">
					{creditPct}% Complete
				</p>
			</div>
		</div>

		<!-- Activities Logged -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<div class="p-2 rounded-lg bg-blue-50 text-blue-600 border border-blue-100">
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
							d="M3.75 3v11.25A2.25 2.25 0 0 0 6 16.5h12M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-16.5 0v11.25m4.5-8.25v6m4.5-9v9m4.5-6v6"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-3">
				<span class="text-2xl font-bold font-serif text-slate-900">{student.activityCount}</span>
				<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1">Activities Logged</h3>
				<p class="text-[10px] font-bold text-slate-400 mt-0.5 uppercase tracking-wider">
					This academic year
				</p>
			</div>
		</div>

		<!-- Certificates Verified -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<div class="p-2 rounded-lg bg-emerald-50 text-emerald-600 border border-emerald-100">
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
							d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-3">
				<span class="text-2xl font-bold font-serif text-slate-900">{student.certificates}</span>
				<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1">Certificates Verified</h3>
				<p class="text-[10px] font-bold text-amber-500 mt-0.5 uppercase tracking-wider">
					{pendingCerts} pending review
				</p>
			</div>
		</div>

		<!-- Current Rank -->
		<div
			class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs hover:shadow-md transition-shadow"
		>
			<div class="flex items-center justify-between">
				<div class="p-2 rounded-lg bg-amber-50 text-amber-600 border border-amber-100">
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
							d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M7.73 9.728a6.726 6.726 0 0 0 2.748 1.35m8.272-6.842V4.5c0 2.108-.966 3.99-2.48 5.228m2.48-5.492a46.32 46.32 0 0 1 2.916.52 6.003 6.003 0 0 1-5.395 4.972m0 0a6.726 6.726 0 0 1-2.749 1.35m0 0a6.772 6.772 0 0 1-3.044 0"
						/>
					</svg>
				</div>
			</div>
			<div class="mt-3">
				<span class="text-2xl font-bold font-serif text-slate-900">#{rank}</span>
				<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1">Current Rank</h3>
				<p class="text-[10px] font-bold text-slate-400 mt-0.5 uppercase tracking-wider">
					Among {cohortSize} students
				</p>
			</div>
		</div>
	</section>

	<!-- ── Credit Progress ───────────────────────────────────────────────────── -->
	<section class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
		<div class="p-5 border-b border-slate-100 bg-slate-50/20">
			<h2 class="text-sm font-bold font-serif text-inst-navy">Credit Progress</h2>
		</div>
		<div class="p-5 space-y-4">
			<div class="grid grid-cols-3 gap-4">
				<div class="text-center rounded-xl bg-emerald-50/60 border border-emerald-100 py-4">
					<div class="text-2xl font-bold font-serif text-emerald-600">{student.creditsEarned}</div>
					<div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
						Credits Earned
					</div>
				</div>
				<div class="text-center rounded-xl bg-slate-50 border border-slate-150 py-4">
					<div class="text-2xl font-bold font-serif text-slate-800">{student.creditsTarget}</div>
					<div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
						Credits Required
					</div>
				</div>
				<div class="text-center rounded-xl bg-amber-50/60 border border-amber-100 py-4">
					<div class="text-2xl font-bold font-serif text-amber-600">{creditsRemaining}</div>
					<div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
						Remaining
					</div>
				</div>
			</div>

			<div>
				<div class="flex items-center justify-between mb-1.5">
					<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
						>{student.creditsEarned} / {student.creditsTarget} Credits</span
					>
					<span class="text-[10px] font-extrabold text-slate-600">{creditPct}% Complete</span>
				</div>
				<div class="h-2.5 w-full bg-slate-100 rounded-full overflow-hidden">
					<div
						class="h-full rounded-full bg-gradient-to-r from-[#881B1B] to-rose-500 transition-all duration-500"
						style="width: {creditPct}%"
					></div>
				</div>
				<p class="text-[10px] text-slate-400 font-semibold mt-2">
					{student.name.split(' ')[0]} has completed {creditPct}% of the required {student.creditsTarget}
					credits. Maintaining the current activity pace, on track to complete before the end of the semester.
				</p>
			</div>
		</div>
	</section>

	<!-- ── History (tabs) ────────────────────────────────────────────────────── -->
	<section class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
		<!-- Tabs -->
		<div class="px-5 pt-4 border-b border-slate-100 flex items-center gap-6">
			<button
				onclick={() => (activeTab = 'activity')}
				class="pb-3 text-xs font-bold border-b-2 transition-colors {activeTab === 'activity'
					? 'border-[#881B1B] text-inst-navy'
					: 'border-transparent text-slate-400 hover:text-slate-600'}"
			>
				Activity Participation History
			</button>
			<button
				onclick={() => (activeTab = 'certificate')}
				class="pb-3 text-xs font-bold border-b-2 transition-colors {activeTab === 'certificate'
					? 'border-[#881B1B] text-inst-navy'
					: 'border-transparent text-slate-400 hover:text-slate-600'}"
			>
				Certificate History
			</button>
		</div>

		{#if activeTab === 'activity'}
			<!-- Filter bar -->
			<div class="px-5 py-3.5 border-b border-slate-100 flex flex-wrap items-center gap-3">
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
						placeholder="Search activity..."
						bind:value={searchQuery}
						class="pl-8 pr-4 py-2 bg-slate-50 rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-slate-350 focus:bg-white w-48 transition-all focus:w-56"
					/>
				</div>

				<select
					bind:value={filterCategory}
					class="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-[11px] font-bold text-slate-600 focus:outline-none focus:border-slate-350"
				>
					{#each categories as cat}
						<option value={cat}>{cat === 'All' ? 'Category' : cat}</option>
					{/each}
				</select>

				<select
					bind:value={filterStatus}
					class="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-[11px] font-bold text-slate-600 focus:outline-none focus:border-slate-350"
				>
					<option value="All">Status</option>
					<option value="Completed">Completed</option>
					<option value="Pending">Pending</option>
					<option value="Rejected">Rejected</option>
				</select>
			</div>

			<!-- Activity table -->
			<div class="overflow-x-auto">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr
							class="border-b border-slate-100 bg-slate-50/50 text-[10px] font-extrabold text-slate-400 uppercase tracking-wider"
						>
							<th class="py-3 px-5">Activity Name</th>
							<th class="py-3 px-5">Category</th>
							<th class="py-3 px-5">Date</th>
							<th class="py-3 px-5">Credits</th>
							<th class="py-3 px-5">Status</th>
							<th class="py-3 px-5 text-center">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 text-xs font-sans">
						{#if filteredActivities.length === 0}
							<tr>
								<td colspan="6" class="py-14 text-center text-slate-400 font-semibold text-xs">
									No activities match your filters.
								</td>
							</tr>
						{:else}
							{#each filteredActivities as activity (activity.id)}
								<tr class="hover:bg-slate-50/50 transition-colors">
									<td class="py-3.5 px-5 font-bold text-slate-800">{activity.name}</td>
									<td class="py-3.5 px-5">
										<span
											class="text-[10px] font-bold text-slate-500 bg-slate-100 border border-slate-200 px-2 py-0.5 rounded-md"
											>{activity.category}</span
										>
									</td>
									<td class="py-3.5 px-5 text-slate-500 font-semibold">{activity.date}</td>
									<td class="py-3.5 px-5 font-extrabold text-slate-800">{activity.credits}</td>
									<td class="py-3.5 px-5">
										<span
											class="inline-flex items-center gap-1.5 px-2 py-1 rounded-md text-[10px] font-extrabold uppercase tracking-wide {statusChip(
												activity.status
											)}"
										>
											<span class="w-1.5 h-1.5 rounded-full {statusDot(activity.status)}"></span>
											{activity.status}
										</span>
									</td>
									<td class="py-3.5 px-5 text-center">
										<button
											onclick={() => toast(`Opening “${activity.name}”`)}
											class="text-[11px] font-bold text-inst-navy hover:underline"
										>
											View Activity
										</button>
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>
		{:else}
			<!-- Certificate table -->
			<div class="overflow-x-auto">
				<table class="w-full text-left border-collapse">
					<thead>
						<tr
							class="border-b border-slate-100 bg-slate-50/50 text-[10px] font-extrabold text-slate-400 uppercase tracking-wider"
						>
							<th class="py-3 px-5">Certificate</th>
							<th class="py-3 px-5">Issuer</th>
							<th class="py-3 px-5">Date</th>
							<th class="py-3 px-5">Credits</th>
							<th class="py-3 px-5">Status</th>
							<th class="py-3 px-5 text-center">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 text-xs font-sans">
						{#each certificates as cert (cert.id)}
							<tr class="hover:bg-slate-50/50 transition-colors">
								<td class="py-3.5 px-5 font-bold text-slate-800">{cert.name}</td>
								<td class="py-3.5 px-5 text-slate-600 font-semibold">{cert.issuer}</td>
								<td class="py-3.5 px-5 text-slate-500 font-semibold">{cert.date}</td>
								<td class="py-3.5 px-5 font-extrabold text-slate-800">{cert.credits}</td>
								<td class="py-3.5 px-5">
									<span
										class="inline-flex items-center gap-1.5 px-2 py-1 rounded-md text-[10px] font-extrabold uppercase tracking-wide {statusChip(
											cert.status
										)}"
									>
										<span class="w-1.5 h-1.5 rounded-full {statusDot(cert.status)}"></span>
										{cert.status}
									</span>
								</td>
								<td class="py-3.5 px-5 text-center">
									<button
										onclick={() => toast(`Opening certificate “${cert.name}”`)}
										class="text-[11px] font-bold text-inst-navy hover:underline"
									>
										View
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</section>

	<!-- ── Performance Summary + Quick Insights ──────────────────────────────── -->
	<section class="grid grid-cols-1 lg:grid-cols-2 gap-5 items-stretch">
		<!-- Performance Summary -->
		<div
			class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden flex flex-col"
		>
			<div class="p-5 border-b border-slate-100 bg-slate-50/20">
				<h2 class="text-sm font-bold font-serif text-inst-navy">Performance Summary</h2>
			</div>
			<div class="p-5 grid grid-cols-1 sm:grid-cols-2 gap-4 flex-grow">
				<div class="flex items-center gap-3 bg-slate-50 rounded-xl border border-slate-150 p-4">
					<div
						class="w-9 h-9 rounded-lg bg-blue-100 text-blue-600 border border-blue-200 flex items-center justify-center shrink-0"
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
								d="M3.75 3v11.25A2.25 2.25 0 0 0 6 16.5h12M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-16.5 0v11.25m4.5-8.25v6m4.5-9v9m4.5-6v6"
							/>
						</svg>
					</div>
					<div class="min-w-0">
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Most Active Category</span
						>
						<span class="text-sm font-bold text-slate-900 block">Technical</span>
					</div>
				</div>

				<div class="flex items-center gap-3 bg-slate-50 rounded-xl border border-slate-150 p-4">
					<div
						class="w-9 h-9 rounded-lg bg-yellow-100 text-yellow-600 border border-yellow-200 flex items-center justify-center shrink-0"
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
					<div class="min-w-0">
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Highest Credit Activity</span
						>
						<span class="text-sm font-bold text-slate-900 block truncate">NPTEL Certification</span>
						<span class="text-[10px] text-slate-500 font-semibold">20 Credits</span>
					</div>
				</div>

				<div class="flex items-center gap-3 bg-slate-50 rounded-xl border border-slate-150 p-4">
					<div
						class="w-9 h-9 rounded-lg bg-emerald-100 text-emerald-600 border border-emerald-200 flex items-center justify-center shrink-0"
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
					<div class="min-w-0">
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Participation Score</span
						>
						<span class="text-sm font-bold text-slate-900 block">{participationScore}%</span>
					</div>
				</div>

				<div class="flex items-center gap-3 bg-slate-50 rounded-xl border border-slate-150 p-4">
					<div
						class="w-9 h-9 rounded-lg bg-purple-100 text-purple-600 border border-purple-200 flex items-center justify-center shrink-0"
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
								d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5"
							/>
						</svg>
					</div>
					<div class="min-w-0">
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Last Activity Date</span
						>
						<span class="text-sm font-bold text-slate-900 block">24 Jun 2025</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Quick Insights -->
		<div
			class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden flex flex-col"
		>
			<div class="p-5 border-b border-slate-100 bg-slate-50/20">
				<h2 class="text-sm font-bold font-serif text-inst-navy">Quick Insights</h2>
			</div>
			<div class="p-4 space-y-2.5 flex-grow flex flex-col">
				<div
					class="flex items-center gap-2.5 p-3 rounded-lg bg-amber-50/60 border border-amber-100"
				>
					<span class="w-2 h-2 rounded-full bg-amber-500 shrink-0"></span>
					<p class="text-[11px] font-bold text-slate-700">
						{pendingCerts} Certificates Pending Verification
					</p>
				</div>
				<div
					class="flex items-center gap-2.5 p-3 rounded-lg bg-emerald-50/60 border border-emerald-100"
				>
					<span class="w-2 h-2 rounded-full bg-emerald-500 shrink-0"></span>
					<p class="text-[11px] font-bold text-slate-700">Excellent Participation Rate</p>
				</div>
				<div class="flex items-center gap-2.5 p-3 rounded-lg bg-blue-50/60 border border-blue-100">
					<span class="w-2 h-2 rounded-full bg-blue-500 shrink-0"></span>
					<p class="text-[11px] font-bold text-slate-700">Top 10% in Assigned Batch</p>
				</div>
				<div
					class="flex items-center gap-2.5 p-3 rounded-lg bg-purple-50/60 border border-purple-100"
				>
					<span class="w-2 h-2 rounded-full bg-purple-500 shrink-0"></span>
					<p class="text-[11px] font-bold text-slate-700">Consistent Activity Submission</p>
				</div>

				<div class="pt-2 mt-auto">
					<div class="flex items-center justify-between mb-1.5">
						<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
							>Overall Progress</span
						>
						<span class="text-[10px] font-extrabold text-slate-600">{creditPct}%</span>
					</div>
					<div class="h-2 w-full bg-slate-100 rounded-full overflow-hidden">
						<div
							class="h-full rounded-full bg-gradient-to-r from-[#881B1B] to-rose-500"
							style="width: {creditPct}%"
						></div>
					</div>
					<button
						onclick={() => toast(`Preparing full activity report for ${student.name}...`)}
						class="w-full mt-3 py-2 text-[11px] font-bold text-[#881B1B] border border-[#881B1B]/20 bg-[#881B1B]/5 hover:bg-[#881B1B]/10 rounded-lg transition-colors tracking-wide uppercase"
					>
						View Full Activity Report →
					</button>
				</div>
			</div>
		</div>
	</section>

	<!-- ── Mentor Observations ───────────────────────────────────────────────── -->
	<section class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
		<div
			class="p-5 border-b border-slate-100 bg-slate-50/20 flex items-center justify-between gap-3"
		>
			<h2 class="text-sm font-bold font-serif text-inst-navy flex items-center gap-2">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-4 h-4 text-[#881B1B]"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z"
					/>
				</svg>
				Mentor Observations
			</h2>
			<div class="flex items-center gap-2 shrink-0">
				<button
					onclick={() => openEditNote(notes[notes.length - 1])}
					disabled={notes.length === 0}
					class="inline-flex items-center gap-1.5 px-3 py-1.5 text-[11px] font-bold text-slate-600 border border-slate-200 rounded-lg hover:bg-slate-50 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
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
							d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Z"
						/>
					</svg>
					Edit Note
				</button>
				<button
					onclick={openAddNote}
					class="inline-flex items-center gap-1.5 px-3 py-1.5 text-[11px] font-bold text-white bg-[#881B1B] hover:bg-[#881B1B]/90 rounded-lg transition-colors shadow-xs"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						class="w-3.5 h-3.5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
					</svg>
					Add Note
				</button>
			</div>
		</div>

		<div class="p-5 space-y-4">
			{#if composerOpen}
				<div
					transition:slide={{ duration: 150 }}
					class="rounded-xl border border-slate-200 p-4 bg-slate-50/40"
				>
					<textarea
						bind:value={composerText}
						rows="3"
						placeholder="Write an observation about this student..."
						class="w-full text-xs text-slate-700 bg-white border border-slate-200 rounded-lg p-3 focus:outline-none focus:border-slate-350 resize-none"
					></textarea>
					<div class="flex items-center justify-end gap-2 mt-3">
						<button
							onclick={() => {
								composerOpen = false;
								composerText = '';
								editingId = null;
							}}
							class="px-3 py-1.5 text-[11px] font-bold text-slate-600 border border-slate-200 rounded-lg hover:bg-slate-100 transition-colors"
						>
							Cancel
						</button>
						<button
							onclick={saveNote}
							class="px-3 py-1.5 text-[11px] font-bold text-white bg-[#881B1B] hover:bg-[#881B1B]/90 rounded-lg transition-colors shadow-xs"
						>
							{editingId !== null ? 'Update' : 'Save'} Observation
						</button>
					</div>
				</div>
			{/if}

			{#each notes as note (note.id)}
				<div class="rounded-xl border border-slate-150 p-4">
					<div class="flex items-start gap-3">
						<div
							class="w-9 h-9 rounded-full bg-[#881B1B] text-white flex items-center justify-center font-bold text-xs border-2 border-white shadow-sm font-serif shrink-0"
						>
							{initials(note.author)}
						</div>
						<div class="min-w-0 flex-grow">
							<div class="flex items-center gap-2 flex-wrap">
								<span class="text-xs font-bold text-slate-900">{note.author}</span>
								<span
									class="text-[9px] font-extrabold text-emerald-700 bg-emerald-50 border border-emerald-100 px-1.5 py-0.5 rounded uppercase tracking-wide"
									>{note.badge}</span
								>
							</div>
							<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider"
								>{note.role}</span
							>
							<p class="text-xs text-slate-600 font-medium leading-relaxed mt-2">{note.text}</p>
						</div>
					</div>
				</div>
			{/each}

			<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
				Showing {notes.length} of {notes.length} observations
			</p>
		</div>
	</section>
</div>
