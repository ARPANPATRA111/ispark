<script lang="ts">
	import { fade, slide } from 'svelte/transition';

	// Mock Profile State
	let profile = $state({
		name: 'Rahul Verma',
		rollNo: 'CSE-2K23-78',
		enrollmentNo: 'DX7890543',
		email: 'rahul.iips@gmail.com',
		mobile: '+91 98765 43210',
		dob: '-',
		gender: 'Male',
		course: 'B.Tech',
		program: 'Computer Science & Engg.',
		semester: 'Semester 6',
		batch: '2021 - 2025',
		department: 'Computer Science',
		institute: 'IIPS, DAVV',
		status: 'Active Student',
		verified: true,
		photoUrl: '', // Empty for now, missing photo
		emailVerified: true,
		phoneVerified: true,
		lastLogin: 'Today, 09:42 AM',
		lastPasswordChange: '12 Jun 2025'
	});

	// UI Interactive States
	let isEditing = $state(false);
	let showPasswordModal = $state(false);

	// Temp Edit Form Fields
	let editForm = $state({
		email: '',
		mobile: '',
		dob: '',
		gender: ''
	});

	// Temp Password Form
	let passwordForm = $state({
		currentPassword: '',
		newPassword: '',
		confirmPassword: ''
	});
	let passwordError = $state('');
	let passwordSuccess = $state(false);

	// Derived Completion Percentage & Missing Info
	let missingInfo = $derived.by(() => {
		const list = [];
		if (!profile.dob || profile.dob === '-') list.push({ key: 'dob', name: 'Date of Birth' });
		if (!profile.photoUrl) list.push({ key: 'photo', name: 'Profile Photo' });
		return list;
	});

	let completionRate = $derived.by(() => {
		let rate = 100;
		if (!profile.dob || profile.dob === '-') rate -= 10;
		if (!profile.photoUrl) rate -= 5;
		return rate;
	});

	// Start edit mode
	function handleStartEdit() {
		editForm.email = profile.email;
		editForm.mobile = profile.mobile;
		editForm.dob = profile.dob === '-' ? '' : profile.dob;
		editForm.gender = profile.gender;
		isEditing = true;
	}

	// Save edit form
	function handleSaveEdit(e: Event) {
		e.preventDefault();
		profile.email = editForm.email;
		profile.mobile = editForm.mobile;
		profile.dob = editForm.dob ? editForm.dob : '-';
		profile.gender = editForm.gender;
		isEditing = false;
	}

	// Cancel edit mode
	function handleCancelEdit() {
		isEditing = false;
	}

	// Handle password submission
	function handlePasswordSubmit(e: Event) {
		e.preventDefault();
		passwordError = '';
		if (passwordForm.newPassword !== passwordForm.confirmPassword) {
			passwordError = 'Passwords do not match.';
			return;
		}
		if (passwordForm.newPassword.length < 6) {
			passwordError = 'Password must be at least 6 characters.';
			return;
		}

		// Success flow
		passwordSuccess = true;
		profile.lastPasswordChange = new Date().toLocaleDateString('en-GB', {
			day: 'numeric',
			month: 'short',
			year: 'numeric'
		});

		setTimeout(() => {
			showPasswordModal = false;
			passwordSuccess = false;
			passwordForm = { currentPassword: '', newPassword: '', confirmPassword: '' };
		}, 1500);
	}

	// Complete profile automatically / quick fix
	function handleCompleteProfile() {
		if (profile.dob === '-') {
			profile.dob = '2004-10-15';
		}
		if (!profile.photoUrl) {
			profile.photoUrl =
				'https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?auto=format&fit=crop&q=80&w=120'; // Mock avatar
		}
	}
</script>

<div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start font-sans">
	<!-- LEFT COLUMN (cols 1-7 on desktop) -->
	<div class="lg:col-span-7 space-y-6">
		<!-- 1. STUDENT AVATAR / SUMMARY CARD -->
		<div
			class="bg-white rounded-xl border border-slate-200 p-6 flex flex-col items-center relative overflow-hidden shadow-xs"
		>
			<!-- Decor Top Line -->
			<div class="absolute top-0 left-0 right-0 h-1 bg-[#881B1B]"></div>

			<!-- Circular Avatar -->
			<div
				class="w-24 h-24 rounded-full bg-[#881B1B] text-white flex items-center justify-center font-bold text-3xl border-4 border-white shadow-md shrink-0 relative overflow-hidden"
			>
				{#if profile.photoUrl}
					<img src={profile.photoUrl} alt={profile.name} class="w-full h-full object-cover" />
				{:else}
					RV
				{/if}
			</div>

			<h2 class="text-xl font-bold text-[#0B1535] mt-4 font-serif">{profile.name}</h2>
			<p class="text-xs text-slate-500 font-semibold mt-1">Roll No: {profile.rollNo}</p>

			<!-- Verified Badge -->
			<div
				class="mt-3 flex items-center gap-1.5 px-3 py-1 bg-emerald-50 text-emerald-700 border border-emerald-200 rounded-full text-[11px] font-bold"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 20 20"
					fill="currentColor"
					class="w-3.5 h-3.5"
				>
					<path
						fill-rule="evenodd"
						d="M16.403 12.652a3 3 0 000-5.304 3 3 0 00-3.75-3.751 3 3 0 00-5.305 0 3 3 0 00-3.751 3.75 3 3 0 000 5.305 3 3 0 003.75 3.751 3 3 0 005.305 0 3 3 0 003.751-3.75zm-2.546-3.073a.75.75 0 00-1.06-1.06L9 11.293 7.72 10.01a.75.75 0 00-1.06 1.06l1.8 1.8a.75.75 0 001.06 0l3.8-3.8z"
						clip-rule="evenodd"
					/>
				</svg>
				Verified Student
			</div>
		</div>

		<!-- 2. PERSONAL INFORMATION CARD -->
		<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
			<div class="flex items-center justify-between pb-3 border-b border-slate-100 mb-5">
				<h3 class="text-xs font-bold text-slate-400 tracking-wider uppercase font-sans">
					Personal Information
				</h3>
			</div>

			{#if !isEditing}
				<!-- VIEW MODE -->
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-y-4 gap-x-6 text-xs">
					<div>
						<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
							>Full Name</span
						>
						<span class="font-bold text-[#0B1535] mt-1 block">{profile.name}</span>
					</div>
					<div>
						<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
							>Roll Number</span
						>
						<span class="font-bold text-[#0B1535] mt-1 block">{profile.rollNo}</span>
					</div>
					<div>
						<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
							>Enrollment No.</span
						>
						<span class="font-bold text-[#0B1535] mt-1 block">{profile.enrollmentNo}</span>
					</div>
					<div>
						<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
							>Email Address</span
						>
						<span class="font-bold text-[#0B1535] mt-1 block break-all">{profile.email}</span>
					</div>
					<div>
						<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
							>Mobile Number</span
						>
						<span class="font-bold text-[#0B1535] mt-1 block">{profile.mobile}</span>
					</div>
					<div>
						<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
							>Date of Birth</span
						>
						<span class="font-bold text-[#0B1535] mt-1 block">
							{profile.dob === '-'
								? '-'
								: new Date(profile.dob).toLocaleDateString('en-GB', {
										day: 'numeric',
										month: 'short',
										year: 'numeric'
									})}
						</span>
					</div>
					<div>
						<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
							>Gender</span
						>
						<span class="font-bold text-[#0B1535] mt-1 block">{profile.gender}</span>
					</div>
				</div>

				<!-- ACTIONS ROW -->
				<div class="grid grid-cols-2 gap-3.5 mt-8 pt-6 border-t border-slate-100">
					<button
						onclick={handleStartEdit}
						class="flex items-center justify-center gap-2 py-3 bg-[#0B1535] hover:bg-[#1a2b5e] text-white rounded-lg text-xs font-bold transition duration-200 shadow-xs focus:outline-none cursor-pointer"
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
								d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.83 21.82a.75.75 0 0 1-.34.201L3 22.887l.859-3.542a.75.75 0 0 1 .202-.34l11.758-11.76H16.862z"
							/>
						</svg>
						Edit Profile
					</button>

					<button
						onclick={() => (showPasswordModal = true)}
						class="flex items-center justify-center gap-2 py-3 border border-slate-200 hover:bg-slate-50 text-[#0B1535] bg-white rounded-lg text-xs font-bold transition duration-200 shadow-xs focus:outline-none cursor-pointer"
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
								d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25z"
							/>
						</svg>
						Change Password
					</button>
				</div>
			{:else}
				<!-- EDIT MODE FORM -->
				<form onsubmit={handleSaveEdit} class="space-y-4 text-xs">
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<!-- Read-only items -->
						<div>
							<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
								>Full Name</span
							>
							<span
								class="font-bold text-slate-500 mt-1.5 block px-3 py-2 bg-slate-50 rounded-lg border border-slate-100"
								>{profile.name}</span
							>
						</div>
						<div>
							<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
								>Roll Number</span
							>
							<span
								class="font-bold text-slate-500 mt-1.5 block px-3 py-2 bg-slate-50 rounded-lg border border-slate-100"
								>{profile.rollNo}</span
							>
						</div>
						<div>
							<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
								>Enrollment No.</span
							>
							<span
								class="font-bold text-slate-500 mt-1.5 block px-3 py-2 bg-slate-50 rounded-lg border border-slate-100"
								>{profile.enrollmentNo}</span
							>
						</div>

						<!-- Editable Items -->
						<div class="flex flex-col gap-1.5">
							<label
								for="edit-email"
								class="text-[10px] font-bold text-slate-700 tracking-wider uppercase"
								>Email Address *</label
							>
							<input
								id="edit-email"
								type="email"
								required
								bind:value={editForm.email}
								class="w-full px-3 py-2 bg-white rounded-lg border border-slate-250 text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
							/>
						</div>

						<div class="flex flex-col gap-1.5">
							<label
								for="edit-mobile"
								class="text-[10px] font-bold text-slate-700 tracking-wider uppercase"
								>Mobile Number *</label
							>
							<input
								id="edit-mobile"
								type="text"
								required
								bind:value={editForm.mobile}
								class="w-full px-3 py-2 bg-white rounded-lg border border-slate-250 text-xs text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
							/>
						</div>

						<div class="flex flex-col gap-1.5">
							<label
								for="edit-dob"
								class="text-[10px] font-bold text-slate-700 tracking-wider uppercase"
								>Date of Birth</label
							>
							<input
								id="edit-dob"
								type="date"
								bind:value={editForm.dob}
								class="w-full px-3 py-2 bg-white rounded-lg border border-slate-250 text-xs text-slate-800 placeholder:text-slate-450 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
							/>
						</div>

						<div class="flex flex-col gap-1.5">
							<label
								for="edit-gender"
								class="text-[10px] font-bold text-slate-700 tracking-wider uppercase"
								>Gender *</label
							>
							<select
								id="edit-gender"
								bind:value={editForm.gender}
								class="w-full px-3 py-2 bg-white rounded-lg border border-slate-250 text-xs text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
							>
								<option value="Male">Male</option>
								<option value="Female">Female</option>
								<option value="Other">Other</option>
							</select>
						</div>
					</div>

					<!-- EDIT ACTIONS -->
					<div class="flex items-center gap-3.5 mt-8 pt-6 border-t border-slate-100">
						<button
							type="button"
							onclick={handleCancelEdit}
							class="flex-1 py-3 border border-slate-200 hover:bg-slate-50 text-slate-700 bg-white rounded-lg font-bold text-xs tracking-wider uppercase transition duration-200 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="flex-grow py-3 bg-[#881B1B] hover:bg-[#731717] text-white rounded-lg font-bold text-xs tracking-wider uppercase transition duration-200 cursor-pointer"
						>
							Save Changes
						</button>
					</div>
				</form>
			{/if}
		</div>

		<!-- 3. ACADEMIC INFORMATION CARD -->
		<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
			<div class="flex items-center justify-between pb-3 border-b border-slate-100 mb-5">
				<h3 class="text-xs font-bold text-slate-400 tracking-wider uppercase font-sans">
					Academic Information
				</h3>
			</div>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-y-4 gap-x-6 text-xs">
				<div>
					<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
						>Course</span
					>
					<span class="font-bold text-[#0B1535] mt-1 block">{profile.course}</span>
				</div>
				<div>
					<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
						>Program</span
					>
					<span class="font-bold text-[#0B1535] mt-1 block">{profile.program}</span>
				</div>
				<div>
					<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
						>Semester</span
					>
					<span class="font-bold text-[#0B1535] mt-1 block">{profile.semester}</span>
				</div>
				<div>
					<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
						>Batch</span
					>
					<span class="font-bold text-[#0B1535] mt-1 block">{profile.batch}</span>
				</div>
				<div>
					<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
						>Department</span
					>
					<span class="font-bold text-[#0B1535] mt-1 block">{profile.department}</span>
				</div>
				<div>
					<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
						>Institute</span
					>
					<span class="font-bold text-[#0B1535] mt-1 block">{profile.institute}</span>
				</div>
				<div class="sm:col-span-2">
					<span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block"
						>Current Status</span
					>
					<div class="flex items-center gap-1.5 mt-1.5">
						<span class="w-2 h-2 rounded-full bg-emerald-500"></span>
						<span class="font-bold text-emerald-600">{profile.status}</span>
					</div>
				</div>
			</div>
		</div>

		<!-- 4. SECURITY & ACCOUNT CARD -->
		<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
			<div class="flex items-center justify-between pb-3 border-b border-slate-100 mb-5">
				<h3 class="text-xs font-bold text-slate-400 tracking-wider uppercase font-sans">
					Security & Account
				</h3>
			</div>

			<div class="space-y-4 text-xs font-semibold text-slate-700">
				<div class="flex items-center justify-between py-1">
					<div class="flex items-center gap-3">
						<!-- Icon -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="2"
							stroke="currentColor"
							class="w-5 h-5 text-slate-400"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25H4.5A2.25 2.25 0 0 1 2.25 17.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5H4.5a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75"
							/>
						</svg>
						<span class="text-slate-800 font-bold">Email Verification</span>
					</div>
					<div class="flex items-center gap-1 text-emerald-600 font-bold">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 20 20"
							fill="currentColor"
							class="w-4 h-4"
						>
							<path
								fill-rule="evenodd"
								d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5z"
								clip-rule="evenodd"
							/>
						</svg>
						Verified
					</div>
				</div>

				<div class="flex items-center justify-between py-1">
					<div class="flex items-center gap-3">
						<!-- Icon -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="2"
							stroke="currentColor"
							class="w-5 h-5 text-slate-400"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 0 0 2.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-2.824-1.502-5.184-3.864-6.686-6.686l1.294-.97c.362-.272.528-.734.417-1.173L6.963 3.102a1.125 1.125 0 0 0-1.091-.852H4.5A2.25 2.25 0 0 0 2.25 4.5v2.25z"
							/>
						</svg>
						<span class="text-slate-800 font-bold">Phone Verification</span>
					</div>
					<div class="flex items-center gap-1 text-emerald-600 font-bold">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 20 20"
							fill="currentColor"
							class="w-4 h-4"
						>
							<path
								fill-rule="evenodd"
								d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5z"
								clip-rule="evenodd"
							/>
						</svg>
						Verified
					</div>
				</div>

				<div class="flex items-center justify-between py-1">
					<div class="flex items-center gap-3">
						<!-- Icon -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="2"
							stroke="currentColor"
							class="w-5 h-5 text-slate-400"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0z"
							/>
						</svg>
						<span class="text-slate-800 font-bold">Last Login</span>
					</div>
					<span class="text-slate-600 font-bold">{profile.lastLogin}</span>
				</div>

				<div class="flex items-center justify-between py-1">
					<div class="flex items-center gap-3">
						<!-- Icon -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="2"
							stroke="currentColor"
							class="w-5 h-5 text-slate-400"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25z"
							/>
						</svg>
						<span class="text-slate-800 font-bold">Last Password Change</span>
					</div>
					<span class="text-slate-600 font-bold">{profile.lastPasswordChange}</span>
				</div>

				<div class="flex items-center justify-between py-1">
					<div class="flex items-center gap-3">
						<!-- Icon -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="2"
							stroke="currentColor"
							class="w-5 h-5 text-slate-400"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M9 12.75 11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 0 1-1.043 3.296 3.745 3.745 0 0 1-3.296 1.043A3.745 3.745 0 0 1 12 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 0 1-3.296-1.043 3.745 3.745 0 0 1-1.043-3.296A3.745 3.745 0 0 1 3 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 0 1 1.043-3.296 3.745 3.745 0 0 1 3.296-1.043A3.746 3.746 0 0 1 12 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 0 1 3.296 1.043 3.746 3.746 0 0 1 1.043 3.296A3.745 3.745 0 0 1 21 12Z"
							/>
						</svg>
						<span class="text-slate-800 font-bold">Account Status</span>
					</div>
					<span class="text-emerald-600 font-bold">Active</span>
				</div>
			</div>
		</div>
	</div>

	<!-- RIGHT COLUMN (cols 8-12 on desktop) -->
	<div class="lg:col-span-5 space-y-6">
		<!-- 1. PARTICIPATION OVERVIEW GRID -->
		<div class="space-y-3.5">
			<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase font-sans">
				Participation Overview
			</h3>

			<div class="grid grid-cols-2 gap-4">
				<!-- Activities Participated -->
				<div
					class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs flex flex-col justify-between"
				>
					<span class="text-3xl font-extrabold text-slate-900 font-serif leading-none">24</span>
					<span class="text-[11px] font-bold text-slate-500 mt-3 block leading-tight"
						>Activities Participated</span
					>
				</div>

				<!-- Certificates Uploaded -->
				<div
					class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs flex flex-col justify-between"
				>
					<span class="text-3xl font-extrabold text-slate-900 font-serif leading-none">18</span>
					<span class="text-[11px] font-bold text-slate-500 mt-3 block leading-tight"
						>Certificates Uploaded</span
					>
				</div>

				<!-- Credits Earned -->
				<div
					class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs flex flex-col justify-between"
				>
					<span class="text-3xl font-extrabold text-slate-900 font-serif leading-none">118</span>
					<span class="text-[11px] font-bold text-slate-500 mt-3 block leading-tight"
						>Credits Earned</span
					>
				</div>

				<!-- Leaderboard Rank -->
				<div
					class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs flex flex-col justify-between"
				>
					<span class="text-3xl font-extrabold text-slate-900 font-serif leading-none">#4</span>
					<span class="text-[11px] font-bold text-slate-500 mt-3 block leading-tight"
						>Leaderboard Rank</span
					>
				</div>
			</div>
		</div>

		<!-- 2. CREDIT PROGRESS -->
		<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-3 mb-2">
				<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase font-sans">
					Credit Progress
				</h3>
			</div>

			<div class="flex items-baseline justify-between">
				<div class="flex items-baseline gap-1">
					<span class="text-3xl font-extrabold text-[#0B1535]">118</span>
					<span class="text-xs text-slate-400 font-bold uppercase tracking-wider"
						>/ 200 Credits</span
					>
				</div>
				<span class="text-xs font-bold text-[#881B1B] uppercase tracking-wide">59% Complete</span>
			</div>

			<!-- Custom Progress Bar -->
			<div class="h-2.5 w-full bg-slate-105 rounded-full overflow-hidden">
				<div class="h-full bg-[#881B1B] rounded-full" style="width: 59%"></div>
			</div>

			<!-- Detail counters -->
			<div class="flex items-center justify-between text-[11px] font-bold text-slate-500">
				<span>Credits Earned: <strong class="text-[#0B1535]">118</strong></span>
				<span>Remaining: <strong class="text-slate-800">82</strong></span>
			</div>
		</div>

		<!-- 3. REGISTERED SKILL TRACKS -->
		<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-3 mb-2">
				<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase font-sans">
					Registered Skill Tracks
				</h3>
			</div>

			<div class="flex flex-wrap gap-2.5 pt-1">
				{#each [{ name: 'Technical Skills', color: 'bg-red-50 text-[#881B1B] border-red-150' }, { name: 'Public Speaking', color: 'bg-amber-50 text-[#C89B3C] border-amber-150' }, { name: 'Social Service (NSS)', color: 'bg-emerald-50 text-emerald-700 border-emerald-150' }, { name: 'Athletics', color: 'bg-blue-50 text-blue-700 border-blue-150' }, { name: 'Research & Innovation', color: 'bg-purple-50 text-purple-700 border-purple-150' }] as track}
					<span
						class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold border {track.color} transition-all duration-200 shadow-3xs"
					>
						<span class="w-1.5 h-1.5 rounded-full bg-current shrink-0"></span>
						{track.name}
					</span>
				{/each}
			</div>
		</div>

		<!-- 4. ACHIEVEMENTS & RECOGNITION -->
		<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b border-slate-100 pb-3 mb-2">
				<h3 class="text-xs font-bold text-slate-455 tracking-wider uppercase font-sans">
					Achievements & Recognition
				</h3>
			</div>

			<div class="space-y-3.5 pt-1">
				{#each [{ name: 'National Science Olympiad', points: '20 Credits' }, { name: 'Inter College Debate Championship', points: '15 Credits' }, { name: 'Hackathon Finalist', points: '12 Credits' }] as award}
					<div
						class="flex items-center justify-between p-3.5 border border-slate-150 bg-slate-50/40 hover:bg-slate-50 rounded-xl transition duration-150 gap-4"
					>
						<div class="flex items-center gap-3">
							<!-- Star icon in circle -->
							<div
								class="w-8 h-8 rounded-full bg-amber-50 text-amber-500 flex items-center justify-center border border-amber-100 shrink-0"
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
										d="M11.48 3.499c.15-.358.682-.358.832 0l1.97 3.99a.75.75 0 0 0 .58.58l4.42 1.405c.4.082.56.58.268.88l-3.2 3.12a.75.75 0 0 0-.214.655l.8 4.417c.07.417-.37.74-.75.54L12 18.002l-3.95 2.06c-.38.2-.85-.122-.75-.54l.8-4.417a.75.75 0 0 0-.214-.655L4.69 11.235c-.29-.3-.13-.8.268-.88l4.42-1.405a.75.75 0 0 0 .58-.58l1.97-3.99z"
									/>
								</svg>
							</div>
							<div>
								<h4 class="text-xs font-bold text-[#0B1535] leading-tight">{award.name}</h4>
								<p class="text-[10px] text-slate-400 font-bold mt-1 uppercase tracking-wider">
									{award.points}
								</p>
							</div>
						</div>

						<div
							class="flex items-center gap-1 text-emerald-600 text-[10px] font-extrabold uppercase tracking-wide shrink-0"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 20 20"
								fill="currentColor"
								class="w-3.5 h-3.5"
							>
								<path
									fill-rule="evenodd"
									d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5z"
									clip-rule="evenodd"
								/>
							</svg>
							Verified
						</div>
					</div>
				{/each}
			</div>
		</div>

		<!-- 5. PROFILE COMPLETION -->
		<div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs space-y-5">
			<div class="flex items-center justify-between border-b border-slate-100 pb-3 mb-2">
				<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase font-sans">
					Profile Completion
				</h3>
			</div>

			<div class="flex items-baseline justify-between">
				<span class="text-3xl font-extrabold text-[#0B1535]">{completionRate}%</span>
				<span class="text-xs text-slate-500 font-bold"
					>{completionRate === 100 ? 'Fully Complete' : 'Almost complete'}</span
				>
			</div>

			<!-- Completion Progress Bar -->
			<div class="h-2.5 w-full bg-slate-100 rounded-full overflow-hidden">
				<div
					class="h-full bg-[#0B1535] rounded-full transition-all duration-500"
					style="width: {completionRate}%"
				></div>
			</div>

			{#if missingInfo.length > 0}
				<div class="space-y-2.5">
					<span class="text-[10px] font-bold text-[#C89B3C] uppercase tracking-wider"
						>Missing Information</span
					>
					<div class="space-y-2">
						{#each missingInfo as missing}
							<div class="flex items-center gap-2 text-xs font-semibold text-slate-655">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2.5"
									stroke="currentColor"
									class="w-4 h-4 text-[#C89B3C]"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0zm-9 3.75h.008v.008H12v-.008z"
									/>
								</svg>
								{missing.name}
							</div>
						{/each}
					</div>
				</div>
			{/if}

			<!-- Complete Profile Action Button -->
			{#if completionRate < 100}
				<button
					onclick={handleCompleteProfile}
					class="w-full text-center py-3.5 bg-[#881B1B] hover:bg-[#731717] text-white text-xs font-bold tracking-wider uppercase rounded-xl transition duration-200 shadow-xs focus:outline-none cursor-pointer"
				>
					Complete Profile
				</button>
			{:else}
				<div
					class="w-full text-center py-3.5 bg-emerald-50 border border-emerald-250 text-emerald-700 text-xs font-bold tracking-wider uppercase rounded-xl flex items-center justify-center gap-2"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 20 20"
						fill="currentColor"
						class="w-4 h-4"
					>
						<path
							fill-rule="evenodd"
							d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5z"
							clip-rule="evenodd"
						/>
					</svg>
					Profile Complete !
				</div>
			{/if}
		</div>
	</div>
</div>

<!-- ==================== MODALS ==================== -->

<!-- 1. CHANGE PASSWORD MODAL -->
{#if showPasswordModal}
	<div
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
	>
		<div
			transition:slide={{ duration: 250 }}
			class="bg-white border border-slate-200 rounded-2xl w-full max-w-md shadow-2xl overflow-hidden font-sans"
			role="dialog"
			aria-modal="true"
		>
			<div class="h-1 bg-[#881B1B]"></div>

			<!-- Header -->
			<div class="p-6 border-b border-slate-100 flex items-center justify-between">
				<div>
					<h3 class="text-lg font-bold font-serif text-[#0B1535]">Change Account Password</h3>
					<p class="text-[10px] font-semibold text-[#6B7280] uppercase tracking-widest mt-0.5">
						student security portal
					</p>
				</div>
				<button
					onclick={() => (showPasswordModal = false)}
					aria-label="Close modal"
					class="p-1 rounded-lg text-slate-400 hover:bg-slate-50 hover:text-slate-600 transition-colors"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
						class="w-6 h-6"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<!-- Form -->
			<form onsubmit={handlePasswordSubmit} class="p-6 space-y-4 text-xs">
				{#if passwordSuccess}
					<div transition:fade={{ duration: 250 }} class="text-center py-8 space-y-4">
						<div
							class="w-14 h-14 bg-emerald-100 text-emerald-800 flex items-center justify-center rounded-full mx-auto border-2 border-emerald-250 animate-bounce"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="3"
								class="w-7 h-7"
							>
								<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
							</svg>
						</div>
						<div>
							<h4 class="text-base font-bold text-slate-900">Password Changed Successfully!</h4>
							<p class="text-xs text-slate-500 mt-1">
								Your security credentials have been updated.
							</p>
						</div>
					</div>
				{:else}
					{#if passwordError}
						<div
							class="p-3 bg-rose-50 text-rose-700 border border-rose-100 rounded-lg text-xs font-bold"
						>
							{passwordError}
						</div>
					{/if}

					<!-- Current Password -->
					<div class="flex flex-col gap-1.5">
						<label
							for="current-password"
							class="text-[11px] font-bold text-slate-700 tracking-wider">CURRENT PASSWORD *</label
						>
						<input
							id="current-password"
							type="password"
							required
							bind:value={passwordForm.currentPassword}
							class="w-full px-3.5 py-2.5 bg-slate-50/50 rounded-lg border border-slate-250 text-xs text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
						/>
					</div>

					<!-- New Password -->
					<div class="flex flex-col gap-1.5">
						<label for="new-password" class="text-[11px] font-bold text-slate-700 tracking-wider"
							>NEW PASSWORD *</label
						>
						<input
							id="new-password"
							type="password"
							required
							bind:value={passwordForm.newPassword}
							class="w-full px-3.5 py-2.5 bg-slate-50/50 rounded-lg border border-slate-250 text-xs text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
						/>
					</div>

					<!-- Confirm Password -->
					<div class="flex flex-col gap-1.5">
						<label
							for="confirm-password"
							class="text-[11px] font-bold text-slate-700 tracking-wider"
							>CONFIRM NEW PASSWORD *</label
						>
						<input
							id="confirm-password"
							type="password"
							required
							bind:value={passwordForm.confirmPassword}
							class="w-full px-3.5 py-2.5 bg-slate-50/50 rounded-lg border border-slate-250 text-xs text-slate-800 focus:outline-none focus:border-[#881B1B] focus:ring-1 focus:ring-[#881B1B] transition-all"
						/>
					</div>

					<!-- Buttons -->
					<div class="flex items-center gap-3 pt-3">
						<button
							type="button"
							onclick={() => (showPasswordModal = false)}
							class="flex-1 py-3 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs tracking-wider uppercase rounded-xl transition duration-200 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="flex-grow py-3 bg-[#881B1B] hover:bg-[#731717] text-white font-bold text-xs tracking-wider uppercase rounded-xl transition duration-200 flex items-center justify-center cursor-pointer"
						>
							Update Password
						</button>
					</div>
				{/if}
			</form>
		</div>
	</div>
{/if}
