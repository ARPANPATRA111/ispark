<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { API_BASE_URL } from '$lib/config';

	// Types
	interface AdminProfile {
		admin_id: string;
		name: string;
		email: string;
		role: string;
		assigned_batch: string;
		created_at: string;
		updated_at: string;
	}

	interface BatchInfo {
		batch: string;
		course: string;
		semester: string;
		studentCount: number;
	}

	// Props using Svelte 5 runes
	let {
		admin,
		loading,
		error,
		onEditProfile
	}: {
		admin: AdminProfile | null;
		loading: boolean;
		error: string | null;
		onEditProfile: () => void;
	} = $props();

	// Stats state using Svelte 5 runes
	let statsLoading = $state(true);
	let statsError = $state<string | null>(null);
	let assignedStudentsCount = $state(0);
	let verifiedCertificatesCount = $state(0);
	let pendingReviewsCount = $state(0);
	let activitiesSupervisedCount = $state(0);
	let assignedBatches = $state<BatchInfo[]>([]);

	// Change Password form state
	let isChangePasswordOpen = $state(false);
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let isPasswordSubmitting = $state(false);
	let passwordError = $state<string | null>(null);
	let passwordSuccess = $state<string | null>(null);

	function getInitials(name: string): string {
		if (!name) return 'A';
		const parts = name.split(' ').filter((part) => {
			const lower = part.toLowerCase();
			return (
				lower !== 'dr.' &&
				lower !== 'prof.' &&
				lower !== 'mr.' &&
				lower !== 'ms.' &&
				lower !== 'mrs.'
			);
		});
		if (parts.length === 0) return 'A';
		if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
		return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
	}

	function openChangePassword() {
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
		passwordError = null;
		passwordSuccess = null;
		isChangePasswordOpen = true;
	}

	async function handleChangePassword(e: SubmitEvent) {
		e.preventDefault();
		if (!currentPassword || !newPassword || !confirmPassword) {
			passwordError = 'All fields are required';
			return;
		}
		if (newPassword !== confirmPassword) {
			passwordError = 'New passwords do not match';
			return;
		}

		const isMinLength = newPassword.length >= 8;
		const hasUppercase = /[A-Z]/.test(newPassword);
		const hasLowercase = /[a-z]/.test(newPassword);
		const hasNumber = /\d/.test(newPassword);
		const hasSpecial = /[^A-Za-z0-9]/.test(newPassword);

		if (!isMinLength || !hasUppercase || !hasLowercase || !hasNumber || !hasSpecial) {
			passwordError =
				'Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.';
			return;
		}

		isPasswordSubmitting = true;
		passwordError = null;
		passwordSuccess = null;

		const token = localStorage.getItem('admin_token');
		try {
			const res = await fetch(`${API_BASE_URL}/api/admin/change-password`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify({
					current_password: currentPassword,
					new_password: newPassword,
					confirm_password: confirmPassword
				})
			});

			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || 'Failed to change password');
			}

			passwordSuccess = 'Password changed successfully!';
			setTimeout(() => {
				isChangePasswordOpen = false;
			}, 1500);
		} catch (err) {
			passwordError = err instanceof Error ? err.message : 'An error occurred';
		} finally {
			isPasswordSubmitting = false;
		}
	}

	// Fetch dynamic stats from students list
	onMount(async () => {
		const token = localStorage.getItem('admin_token');
		if (!token) {
			statsLoading = false;
			return;
		}

		try {
			const res = await fetch(`${API_BASE_URL}/api/admin/students`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});

			if (!res.ok) {
				throw new Error('Failed to fetch students stats');
			}

			const data = await res.json();
			const studentsList = data.students || [];

			assignedStudentsCount = studentsList.length;

			let verified = 0;
			let pending = 0;
			const uniqueActivityIds: number[] = [];
			const grouped: BatchInfo[] = [];

			for (const s of studentsList) {
				// 1. Count certificates
				if (s.certificates) {
					for (const cert of s.certificates) {
						if (cert.status === 'Approved') {
							verified++;
						}
						if (cert.status === 'Pending') {
							pending++;
						}
					}
				}

				// 2. Count supervised activities
				if (s.enrollments) {
					for (const enrollment of s.enrollments) {
						if (enrollment.activity_id && !uniqueActivityIds.includes(enrollment.activity_id)) {
							uniqueActivityIds.push(enrollment.activity_id);
						}
					}
				}

				// 3. Group by Batch / Course / Semester
				if (admin) {
					const batch = s.roll_no?.match(/^[A-Z]+2K\d+/)?.[0] ?? admin.assigned_batch ?? 'Unknown';
					const course = s.course_name || '—';
					const semester = s.semester ? `Semester ${s.semester}` : '—';

					const existing = grouped.find(
						(g) => g.batch === batch && g.course === course && g.semester === semester
					);
					if (existing) {
						existing.studentCount++;
					} else {
						grouped.push({
							batch,
							course,
							semester,
							studentCount: 1
						});
					}
				}
			}

			verifiedCertificatesCount = verified;
			pendingReviewsCount = pending;
			activitiesSupervisedCount = uniqueActivityIds.length;
			assignedBatches = grouped;
		} catch (err) {
			console.error('Error fetching admin profile stats:', err);
			statsError = err instanceof Error ? err.message : 'Error loading overview data';
		} finally {
			statsLoading = false;
		}
	});
</script>

<div class="space-y-6 font-sans" transition:fade={{ duration: 150 }}>
	<!-- Page Title Header -->
	<div class="space-y-1">
		<h2 class="text-2xl font-bold font-serif text-slate-900 leading-tight">My Profile</h2>
		<p class="text-xs text-slate-500 font-semibold">
			Manage your professional profile and administrative account
		</p>
	</div>

	<!-- Profile Header Card -->
	<div class="bg-white rounded-xl border border-slate-200 p-6 sm:p-8 shadow-xs relative">
		{#if loading}
			<!-- Loading State Skeleton -->
			<div class="flex flex-col md:flex-row items-center gap-6 animate-pulse">
				<div class="w-24 h-24 rounded-full bg-slate-200 shrink-0"></div>
				<div class="flex-grow space-y-3.5 w-full">
					<div class="h-6 bg-slate-200 rounded w-1/3"></div>
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-y-3 gap-x-8">
						<div class="h-4 bg-slate-200 rounded w-3/4"></div>
						<div class="h-4 bg-slate-200 rounded w-3/4"></div>
						<div class="h-4 bg-slate-200 rounded w-2/3"></div>
						<div class="h-4 bg-slate-200 rounded w-2/3"></div>
						<div class="h-4 bg-slate-200 rounded w-1/2 sm:col-span-2"></div>
					</div>
				</div>
			</div>
		{:else if error}
			<!-- Error State -->
			<div class="p-6 text-center text-rose-600 bg-rose-50 border border-rose-100 rounded-lg">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-8 h-8 mx-auto mb-2 text-rose-500"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"
					/>
				</svg>
				<h4 class="font-bold text-sm">Failed to Load Profile</h4>
				<p class="text-xs text-rose-500 mt-1">{error}</p>
			</div>
		{:else if admin}
			<!-- Main Profile Layout -->
			<div class="flex flex-col md:flex-row items-start gap-6">
				<!-- Left: Circular Avatar -->
				<div
					class="w-24 h-24 rounded-full bg-[#0B1535] text-white flex items-center justify-center font-bold text-3xl border-4 border-slate-100 shadow-md shrink-0 relative overflow-hidden font-serif select-none"
				>
					{getInitials(admin.name)}
				</div>

				<!-- Center: Info Fields -->
				<div class="flex-grow space-y-4 w-full">
					<div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
						<h3 class="text-2xl font-bold text-[#0B1535] font-serif leading-none">
							{admin.name}
						</h3>

						<!-- Right: Actions Buttons -->
						<div class="flex gap-3 shrink-0">
							<!-- Edit Profile button -->
							<button
								type="button"
								onclick={onEditProfile}
								class="inline-flex items-center justify-center gap-1.5 px-4 py-2 border border-slate-250 bg-white hover:bg-slate-50 text-slate-800 rounded-lg text-xs font-bold transition-colors shadow-3xs cursor-pointer focus:outline-none"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2.2"
									stroke="currentColor"
									class="w-3.5 h-3.5 text-slate-500"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="m16.862 4.487 1.687-1.688a1.875 1.875 0 112.652 2.652L6.83 21.82a.75.75 0 01-.34.201L3 22.887l.859-3.542a.75.75 0 01.202-.34l11.758-11.76H16.862z"
									/>
								</svg>
								Edit Profile
							</button>

							<!-- Change Password button -->
							<button
								type="button"
								onclick={openChangePassword}
								class="inline-flex items-center justify-center gap-1.5 px-4 py-2 bg-[#0B1535] hover:bg-[#1a2b5e] text-white rounded-lg text-xs font-bold transition-colors shadow-3xs cursor-pointer focus:outline-none"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2.2"
									stroke="currentColor"
									class="w-3.5 h-3.5 text-white/90"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 00 2.25 2.25z"
									/>
								</svg>
								Change Password
							</button>
						</div>
					</div>

					<!-- Details Grid -->
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-y-3.5 gap-x-8 text-xs leading-normal">
						<!-- Department -->
						<div class="flex items-center gap-3">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2"
								stroke="currentColor"
								class="w-4 h-4 text-slate-400 shrink-0"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M18 18.72a9.094 9.094 0 0 0 3.741-.479 3 3 0 0 0-4.682-2.72m.94 3.198.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0 1 12 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 0 1 6 18.719m0 0a8.984 8.984 0 0 1 3.224-3.357m8.816-2.13A4.9 4.9 0 0 0 16 7.5c0-.985-.29-1.902-.79-2.671M12 18.75h-1.5V1.5h1.5v17.25Z"
								/>
							</svg>
							<div class="flex items-center gap-1.5">
								<span class="text-slate-500 font-semibold">Department:</span>
								<span class="font-bold text-slate-800">-</span>
							</div>
						</div>

						<!-- Employee ID -->
						<div class="flex items-center gap-3">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2"
								stroke="currentColor"
								class="w-4 h-4 text-slate-400 shrink-0"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M15 9h3.75M15 12h3.75M15 15h3.75M4.5 19.5h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5zm6-10.125a1.875 1.875 0 11-3.75 0 1.875 1.875 0 013.75 0zm-1.2 6.477a3.375 3.375 0 00-5.1 0 2.25 2.25 0 002.25 2.25h.6a2.25 2.25 0 002.25-2.25v-.15z"
								/>
							</svg>
							<div class="flex items-center gap-1.5">
								<span class="text-slate-500 font-semibold">Employee ID:</span>
								<span class="font-bold text-slate-800">{admin.admin_id}</span>
							</div>
						</div>

						<!-- Email -->
						<div class="flex items-center gap-3">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2"
								stroke="currentColor"
								class="w-4 h-4 text-slate-400 shrink-0"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25H4.5A2.25 2.25 0 0 1 2.25 17.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5H4.5a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75"
								/>
							</svg>
							<div class="flex items-center gap-1.5 min-w-0">
								<span class="text-slate-500 font-semibold">Email:</span>
								<span class="font-bold text-slate-800 truncate break-all">{admin.email}</span>
							</div>
						</div>

						<!-- Contact -->
						<div class="flex items-center gap-3">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2"
								stroke="currentColor"
								class="w-4 h-4 text-slate-400 shrink-0"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 0 0 2.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-2.824-1.502-5.184-3.864-6.686-6.686l1.294-.97c.362-.272.528-.734.417-1.173L6.963 3.102a1.125 1.125 0 0 0-1.091-.852H4.5A2.25 2.25 0 0 0 2.25 4.5v2.25z"
								/>
							</svg>
							<div class="flex items-center gap-1.5">
								<span class="text-slate-500 font-semibold">Contact:</span>
								<span class="font-bold text-slate-800">-</span>
							</div>
						</div>

						<!-- Office Location -->
						<div class="flex items-center gap-3 sm:col-span-2">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="2"
								stroke="currentColor"
								class="w-4 h-4 text-slate-400 shrink-0"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
								/>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"
								/>
							</svg>
							<div class="flex items-center gap-1.5">
								<span class="text-slate-500 font-semibold">Office Location:</span>
								<span class="font-bold text-slate-800">-</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>

	<!-- Professional Information & Administrative Overview Row -->
	{#if !loading && !error && admin}
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Professional Information Card -->
			<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs flex flex-col">
				<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase font-sans mb-5">
					Professional Information
				</h3>

				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4 flex-grow">
					<!-- Designation -->
					<div
						class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col justify-between"
					>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest leading-none"
							>Designation</span
						>
						<span class="text-xs font-bold text-slate-400 mt-2">—</span>
					</div>

					<!-- Department -->
					<div
						class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col justify-between"
					>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest leading-none"
							>Department</span
						>
						<span class="text-xs font-bold text-slate-800 mt-2">-</span>
					</div>

					<!-- Years of Experience -->
					<div
						class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col justify-between"
					>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest leading-none"
							>Years of Experience</span
						>
						<span class="text-xs font-bold text-slate-400 mt-2">—</span>
					</div>

					<!-- Qualification -->
					<div
						class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col justify-between"
					>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest leading-none"
							>Qualification</span
						>
						<span class="text-xs font-bold text-slate-400 mt-2">—</span>
					</div>

					<!-- Specialization -->
					<div
						class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col justify-between"
					>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest leading-none"
							>Specialization</span
						>
						<span class="text-xs font-bold text-slate-400 mt-2">—</span>
					</div>

					<!-- Mentor Since -->
					<div
						class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col justify-between"
					>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest leading-none"
							>Mentor Since</span
						>
						<span class="text-xs font-bold text-slate-400 mt-2">—</span>
					</div>
				</div>
			</div>

			<!-- Administrative Overview Card -->
			<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs flex flex-col">
				<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase font-sans mb-5">
					Administrative Overview
				</h3>

				{#if statsLoading}
					<!-- Stats Loading State -->
					<div class="grid grid-cols-2 gap-4 flex-grow items-center">
						{#each Array(4) as _}
							<div
								class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col items-center justify-center animate-pulse h-28"
							>
								<div class="w-8 h-8 rounded-full bg-slate-200"></div>
								<div class="h-6 bg-slate-200 rounded w-1/3 mt-3"></div>
								<div class="h-3 bg-slate-200 rounded w-2/3 mt-2"></div>
							</div>
						{/each}
					</div>
				{:else if statsError}
					<!-- Stats Empty State Placeholder -->
					<div class="grid grid-cols-2 gap-4 flex-grow">
						{#each ['Assigned Students', 'Certificates Verified', 'Pending Reviews', 'Activities Supervised'] as label}
							<div
								class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col items-center text-center justify-center"
							>
								<span class="text-2xl font-bold font-serif text-slate-400 mt-3 leading-none">—</span
								>
								<span
									class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-2 leading-tight"
								>
									{label}
								</span>
							</div>
						{/each}
					</div>
				{:else}
					<!-- Stats Content -->
					<div class="grid grid-cols-2 gap-4 flex-grow">
						<!-- Assigned Students -->
						<div
							class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col items-center text-center justify-center hover:shadow-2xs transition-shadow"
						>
							<div
								class="w-9 h-9 rounded-lg bg-blue-50 text-blue-600 border border-blue-100 flex items-center justify-center shrink-0"
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
										d="M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.109A11.386 11.386 0 0 1 10.089 21c-2.316 0-4.445-.69-6.22-1.879v-.003a4.125 4.125 0 0 1 7.533-2.493M15 19.128v-.003c0-1.112-.285-2.16-.786-3.07M14.214 16.058A9.396 9.396 0 0 0 10.089 15c-1.47 0-2.854.34-4.082.945"
									/>
								</svg>
							</div>
							<span class="text-2xl font-bold font-serif text-slate-900 mt-3 leading-none">
								{assignedStudentsCount}
							</span>
							<span
								class="text-[9px] font-bold text-slate-500 uppercase tracking-widest mt-2 leading-tight"
							>
								Assigned Students
							</span>
						</div>

						<!-- Certificates Verified -->
						<div
							class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col items-center text-center justify-center hover:shadow-2xs transition-shadow"
						>
							<div
								class="w-9 h-9 rounded-lg bg-emerald-50 text-emerald-600 border border-emerald-100 flex items-center justify-center shrink-0"
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
										d="M9 12.75 11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 0 1-1.043 3.296 3.745 3.745 0 0 1-3.296 1.043A3.745 3.745 0 0 1 12 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 0 1-3.296-1.043 3.745 3.745 0 0 1-1.043-3.296A3.745 3.745 0 0 1 3 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 0 1 1.043-3.296 3.746 3.746 0 0 1 3.296-1.043A3.746 3.746 0 0 1 12 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 0 1 3.296 1.043 3.746 3.746 0 0 1 1.043 3.296A3.745 3.745 0 0 1 21 12Z"
									/>
								</svg>
							</div>
							<span class="text-2xl font-bold font-serif text-slate-900 mt-3 leading-none">
								{verifiedCertificatesCount}
							</span>
							<span
								class="text-[9px] font-bold text-slate-500 uppercase tracking-widest mt-2 leading-tight"
							>
								Certificates Verified
							</span>
						</div>

						<!-- Pending Reviews -->
						<div
							class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col items-center text-center justify-center hover:shadow-2xs transition-shadow"
						>
							<div
								class="w-9 h-9 rounded-lg bg-amber-50 text-amber-600 border border-amber-100 flex items-center justify-center shrink-0"
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
										d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"
									/>
								</svg>
							</div>
							<span class="text-2xl font-bold font-serif text-slate-900 mt-3 leading-none">
								{pendingReviewsCount}
							</span>
							<span
								class="text-[9px] font-bold text-slate-500 uppercase tracking-widest mt-2 leading-tight"
							>
								Pending Reviews
							</span>
						</div>

						<!-- Activities Supervised -->
						<div
							class="bg-slate-50 border border-slate-150 rounded-xl p-4 flex flex-col items-center text-center justify-center hover:shadow-2xs transition-shadow"
						>
							<div
								class="w-9 h-9 rounded-lg bg-rose-50 text-rose-600 border border-rose-100 flex items-center justify-center shrink-0"
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
							<span class="text-2xl font-bold font-serif text-slate-900 mt-3 leading-none">
								{activitiesSupervisedCount}
							</span>
							<span
								class="text-[9px] font-bold text-slate-500 uppercase tracking-widest mt-2 leading-tight"
							>
								Activities Supervised
							</span>
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Account Security & Assigned Batches Row -->
	{#if !loading && !error && admin}
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Account Security Card -->
			<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs flex flex-col">
				<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase font-sans mb-5">
					Account Security
				</h3>

				<div class="flex-grow space-y-1">
					<!-- Email Verification -->
					<div
						class="flex items-center justify-between py-3 border-b border-slate-100 text-xs font-semibold text-slate-700"
					>
						<div class="flex items-center gap-2.5">
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
									d="M9 12.75 11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068"
								/>
							</svg>
							<span>Email Verification</span>
						</div>
						<span class="text-emerald-600 font-bold">Verified</span>
					</div>

					<!-- Two-Factor Authentication -->
					<div
						class="flex items-center justify-between py-3 border-b border-slate-100 text-xs font-semibold text-slate-700"
					>
						<div class="flex items-center gap-2.5">
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
									d="M9 12.75 11.25 15 15 9.75m-3-3.75h.008v.008H12V9Z"
								/>
							</svg>
							<span>Two-Factor Authentication</span>
						</div>
						<span class="text-slate-400 font-normal">Disabled</span>
					</div>

					<!-- Last Password Change -->
					<div
						class="flex items-center justify-between py-3 border-b border-slate-100 text-xs font-semibold text-slate-700"
					>
						<div class="flex items-center gap-2.5">
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
									d="M15.75 5.25a3 3 0 013 3m3 0a6 6 0 01-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H3.75v-2.25H6V17.25h2.25v-2.25h1.5l1.688-1.688a1.5 1.5 0 011.06-.439c.362 0 .71-.143.963-.398L15.75 5.25Z"
								/>
							</svg>
							<span>Last Password Change</span>
						</div>
						<span class="text-slate-400 font-normal">—</span>
					</div>

					<!-- Last Login -->
					<div
						class="flex items-center justify-between py-3 border-b border-slate-100 text-xs font-semibold text-slate-700"
					>
						<div class="flex items-center gap-2.5">
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
									d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
							<span>Last Login</span>
						</div>
						<span class="text-slate-400 font-normal">—</span>
					</div>

					<!-- Login Devices -->
					<div class="flex items-center justify-between py-3 text-xs font-semibold text-slate-700">
						<div class="flex items-center gap-2.5">
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
									d="M9 17.25v1.007a3 3 0 01-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0115 18.257V17.25m6-12V15a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 15V5.25m18 0A2.25 2.25 0 0018.75 3H5.25A2.25 2.25 0 003 5.25m18 0V12a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 12V5.25"
								/>
							</svg>
							<span>Login Devices</span>
						</div>
						<span class="text-slate-805 font-bold">Current Session</span>
					</div>
				</div>

				<!-- Action buttons -->
				<div class="grid grid-cols-2 gap-3 mt-6 pt-5 border-t border-slate-100">
					<button
						type="button"
						class="px-4 py-2 border border-slate-250 hover:bg-slate-50 text-slate-800 bg-white rounded-lg text-xs font-bold transition-colors cursor-pointer focus:outline-none text-center shadow-3xs"
					>
						Manage Security
					</button>
					<button
						type="button"
						class="px-4 py-2 border border-slate-250 hover:bg-slate-50 text-slate-800 bg-white rounded-lg text-xs font-bold transition-colors cursor-pointer focus:outline-none text-center shadow-3xs"
					>
						View Login History
					</button>
				</div>
			</div>

			<!-- Assigned Batches Card -->
			<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs flex flex-col">
				<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase font-sans mb-5">
					Assigned Batches
				</h3>

				<div class="flex-grow overflow-x-auto">
					<table class="w-full text-left text-xs border-collapse">
						<thead>
							<tr
								class="border-b border-slate-100 text-[10px] font-bold text-slate-400 uppercase tracking-widest"
							>
								<th class="pb-3 font-semibold">Batch</th>
								<th class="pb-3 font-semibold">Course</th>
								<th class="pb-3 font-semibold">Semester</th>
								<th class="pb-3 font-semibold text-right">Students</th>
							</tr>
						</thead>
						<tbody>
							{#if statsLoading}
								<!-- Loading Rows -->
								{#each Array(2) as _}
									<tr class="animate-pulse border-b border-slate-50 last:border-b-0">
										<td class="py-3.5"><div class="h-4 bg-slate-200 rounded w-16"></div></td>
										<td class="py-3.5"><div class="h-4 bg-slate-200 rounded w-12"></div></td>
										<td class="py-3.5"><div class="h-4 bg-slate-200 rounded w-20"></div></td>
										<td class="py-3.5 text-right"
											><div class="h-4 bg-slate-200 rounded w-16 ml-auto"></div></td
										>
									</tr>
								{/each}
							{:else if statsError || assignedBatches.length === 0}
								<!-- Empty state placeholder -->
								<tr>
									<td colspan="4" class="py-12 text-center text-slate-400 font-medium font-sans">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="1.5"
											stroke="currentColor"
											class="w-8 h-8 mx-auto mb-2 text-slate-300"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z"
											/>
										</svg>
										No assigned batches available
									</td>
								</tr>
							{:else}
								{#each assignedBatches as b}
									<tr class="border-b border-slate-50 last:border-b-0 font-semibold text-slate-800">
										<td class="py-3.5 text-slate-900 font-bold">{b.batch}</td>
										<td class="py-3.5 text-slate-500">{b.course}</td>
										<td class="py-3.5 text-slate-500">{b.semester}</td>
										<td class="py-3.5 text-right">
											<span
												class="px-2.5 py-1 bg-slate-100 text-slate-700 border border-slate-200 rounded-md text-[10px] font-bold"
											>
												{b.studentCount} Students
											</span>
										</td>
									</tr>
								{/each}
							{/if}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	{/if}
</div>

<!-- Change Password Modal Overlay -->
{#if isChangePasswordOpen}
	<div
		class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-fade-in"
		transition:fade={{ duration: 100 }}
	>
		<div
			class="bg-white rounded-2xl border border-slate-200 shadow-xl max-w-md w-full p-6 space-y-4"
		>
			<div>
				<h3 class="text-lg font-bold text-slate-900 font-serif leading-tight">Change Password</h3>
				<p class="text-xs text-slate-500 font-medium mt-1">
					Update your account security credential password.
				</p>
			</div>

			<form onsubmit={handleChangePassword} class="space-y-4">
				{#if passwordError}
					<div
						class="p-3 text-xs font-semibold text-rose-650 bg-rose-50 border border-rose-100 rounded-lg"
					>
						{passwordError}
					</div>
				{/if}
				{#if passwordSuccess}
					<div
						class="p-3 text-xs font-semibold text-emerald-650 bg-emerald-50 border border-emerald-100 rounded-lg"
					>
						{passwordSuccess}
					</div>
				{/if}

				<div class="space-y-1.5">
					<label
						for="curr-password"
						class="text-[10px] font-bold text-slate-450 uppercase tracking-wider"
						>Current Password</label
					>
					<input
						id="curr-password"
						type="password"
						bind:value={currentPassword}
						disabled={isPasswordSubmitting}
						class="w-full px-3 py-2 border border-slate-200 focus:border-slate-350 focus:ring-1 focus:ring-slate-350 rounded-lg text-sm focus:outline-none bg-white disabled:bg-slate-50 disabled:text-slate-400 transition-colors"
					/>
				</div>

				<div class="space-y-1.5">
					<label
						for="new-password"
						class="text-[10px] font-bold text-slate-450 uppercase tracking-wider"
						>New Password</label
					>
					<input
						id="new-password"
						type="password"
						bind:value={newPassword}
						disabled={isPasswordSubmitting}
						class="w-full px-3 py-2 border border-slate-200 focus:border-slate-350 focus:ring-1 focus:ring-slate-350 rounded-lg text-sm focus:outline-none bg-white disabled:bg-slate-50 disabled:text-slate-400 transition-colors"
					/>
				</div>

				<div class="space-y-1.5">
					<label
						for="conf-password"
						class="text-[10px] font-bold text-slate-450 uppercase tracking-wider"
						>Confirm New Password</label
					>
					<input
						id="conf-password"
						type="password"
						bind:value={confirmPassword}
						disabled={isPasswordSubmitting}
						class="w-full px-3 py-2 border border-slate-200 focus:border-slate-350 focus:ring-1 focus:ring-slate-350 rounded-lg text-sm focus:outline-none bg-white disabled:bg-slate-50 disabled:text-slate-400 transition-colors"
					/>
				</div>

				<div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100">
					<button
						type="button"
						onclick={() => (isChangePasswordOpen = false)}
						disabled={isPasswordSubmitting}
						class="px-4 py-2 border border-slate-200 hover:bg-slate-50 disabled:opacity-50 text-slate-700 bg-white rounded-lg text-xs font-bold transition-colors cursor-pointer focus:outline-none"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={isPasswordSubmitting}
						class="inline-flex items-center justify-center gap-1.5 px-4 py-2 bg-[#0B1535] hover:bg-[#1a2b5e] disabled:bg-[#0b1535]/50 text-white rounded-lg text-xs font-bold transition-colors cursor-pointer focus:outline-none"
					>
						{#if isPasswordSubmitting}
							<svg
								class="animate-spin -ml-1 mr-1.5 h-3.5 w-3.5 text-white"
								fill="none"
								viewBox="0 0 24 24"
							>
								<circle
									class="opacity-25"
									cx="12"
									cy="12"
									r="10"
									stroke="currentColor"
									stroke-width="4"
								></circle>
								<path
									class="opacity-75"
									fill="currentColor"
									d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
								></path>
							</svg>
							Changing...
						{:else}
							Change Password
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
