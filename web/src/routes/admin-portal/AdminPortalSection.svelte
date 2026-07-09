<script lang="ts">
	import { goto } from '$app/navigation';
	import { onDestroy } from 'svelte';
	import { fade, slide } from 'svelte/transition';

	import { API_BASE_URL } from '$lib/config';

	const formId = 'admin-login';

	// Reactive state using Svelte 5 runes
	let adminId = $state('');
	let password = $state('');
	let showPassword = $state(false);

	// Simulated states
	let submitting = $state(false);
	let loginSuccess = $state(false);
	let errorMsg = $state('');

	// Rate Limiting & Lockout states
	let failedAttempts = $state(0);
	let lockoutTimeLeft = $state(0);
	let lockoutInterval: ReturnType<typeof setInterval> | undefined;

	function startLockout() {
		lockoutTimeLeft = 60;
		errorMsg = 'Too many failed login attempts. Please try again in 60 seconds.';

		if (lockoutInterval) clearInterval(lockoutInterval);

		lockoutInterval = setInterval(() => {
			lockoutTimeLeft -= 1;
			if (lockoutTimeLeft <= 0) {
				clearInterval(lockoutInterval);
				failedAttempts = 0;
				errorMsg = '';
			} else {
				errorMsg = `Too many failed login attempts. Please try again in ${lockoutTimeLeft} seconds.`;
			}
		}, 1000);
	}

	onDestroy(() => {
		if (lockoutInterval) clearInterval(lockoutInterval);
	});

	// Form validity derived state
	let isFormValid = $derived(adminId.trim() !== '' && password.trim() !== '');

	const featureItems = [
		{
			title: 'Verify Activities',
			description: 'Review and verify student activity submissions and certificates.',
			iconPath: 'verify'
		},
		{
			title: 'Manage Grading Records',
			description: 'Monitor extracurricular points and grading allocations.',
			iconPath: 'grade'
		},
		{
			title: 'Generate Reports',
			description: 'Access semester-wise participation reports and statistics.',
			iconPath: 'report'
		},
		{
			title: 'Institutional Oversight',
			description: 'Ensure compliance with iSPARC participation guidelines.',
			iconPath: 'oversight'
		}
	];

	// Check if admin password is updated or not.
	function routeAfterLogin(data: { must_change_password?: boolean }) {
		loginSuccess = true;

		// Persist the flag so the dashboard can guard itself against direct URL access
		localStorage.setItem('admin_must_change_password', String(!!data.must_change_password));

		setTimeout(() => {
			if (data.must_change_password) {
				goto('/admin-portal/update');
			} else {
				goto('/admin-portal/dashboard');
			}
		}, 1000);
	}

	// Submit Handler
	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (!isFormValid || lockoutTimeLeft > 0) return;

		submitting = true;
		errorMsg = '';

		try {
			const response = await fetch(`${API_BASE_URL}/api/admin/auth/login`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ admin_id: adminId, password: password })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Invalid Administrator ID or Password.');
			}

			localStorage.setItem('admin_token', data.access_token);
			failedAttempts = 0;
			routeAfterLogin(data);
		} catch (err) {
			failedAttempts += 1;
			if (failedAttempts >= 5) {
				startLockout();
			} else {
				errorMsg =
					err instanceof Error
						? err.message
						: `Invalid Administrator ID or Password. (Attempt ${failedAttempts} of 5)`;
			}
		} finally {
			submitting = false;
		}
	}
</script>

<section
	class="max-w-6xl w-full mx-auto grid grid-cols-1 lg:grid-cols-12 gap-10 px-4 sm:px-6 lg:px-8 py-12 items-start font-sans"
>
	<!-- Left Column: Admin Portal Description & Administrative Operations -->
	<aside class="lg:col-span-5 flex flex-col gap-6 animate-fade-in">
		<div class="flex flex-col gap-4">
			<div>
				<h1 class="text-4xl font-extrabold text-acad-red font-serif leading-tight tracking-tight">
					Administrative<br />
					<span class="text-inst-navy">Control Portal</span>
				</h1>
			</div>

			<p class="text-slate-500 text-sm leading-relaxed">
				Secure access for authorized iSPARC administrators responsible for monitoring student
				participation records, activity verification, grading workflows, and institutional
				reporting.
			</p>
		</div>

		<!-- Administrative Operations Details -->
		<div class="flex flex-col gap-4 bg-white p-6 rounded-xl border border-border-base shadow-xs">
			<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">
				System Privileges
			</h3>

			<div class="flex flex-col gap-5">
				{#each featureItems as item}
					<div class="flex items-start gap-4">
						<div
							class="flex w-9 h-9 items-center justify-center shrink-0 bg-slate-50 rounded-lg border border-border-base text-inst-navy"
						>
							{#if item.iconPath === 'verify'}
								<!-- Verify SVG -->
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
										d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
									/>
								</svg>
							{:else if item.iconPath === 'grade'}
								<!-- Grade SVG -->
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
										d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582m15.686 0A11.953 11.953 0 0 1 12 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0 1 21 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0 1 12 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 0 1 3 12c0-.778.099-1.533.284-2.253"
									/>
								</svg>
							{:else if item.iconPath === 'report'}
								<!-- Report SVG -->
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
										d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
									/>
								</svg>
							{:else}
								<!-- Oversight SVG -->
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
										d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582"
									/>
								</svg>
							{/if}
						</div>

						<div class="flex flex-col gap-0.5">
							<h4 class="text-sm font-bold text-slate-800">{item.title}</h4>
							<p class="text-xs text-slate-500 leading-normal">{item.description}</p>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</aside>

	<!-- Right Column: Administrative Login Card -->
	<main
		class="lg:col-span-7 bg-white rounded-2xl border border-border-base shadow-sm overflow-hidden animate-fade-up"
	>
		{#if loginSuccess}
			<!-- Success View -->
			<div
				class="p-8 sm:p-12 text-center flex flex-col items-center gap-6"
				in:fade={{ duration: 400 }}
			>
				<div
					class="w-16 h-16 bg-emerald-100 text-emerald-800 flex items-center justify-center rounded-full border-2 border-emerald-250 animate-bounce"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="3"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="w-8 h-8"
					>
						<polyline points="20 6 9 17 4 12" />
					</svg>
				</div>

				<div>
					<h2 class="text-2xl font-bold text-slate-900 font-serif">Authenticated Successfully!</h2>
					<p class="text-slate-500 text-sm mt-2 max-w-md mx-auto">
						Welcome, Administrator <span class="font-semibold text-slate-800">{adminId}</span>. You
						have logged into the iSPARC Administrative Panel. Redirecting to administrative
						console...
					</p>
				</div>

				<div class="flex gap-3 w-full max-w-sm pt-4">
					<a
						href="/admin-portal/dashboard"
						class="flex-1 py-3 text-center bg-inst-navy hover:bg-inst-navy/90 text-white font-semibold text-xs tracking-wider uppercase rounded-xl transition duration-200"
					>
						Go to Console
					</a>
				</div>
			</div>
		{:else}
			<!-- Admin Login Form -->
			<div class="p-6 sm:p-8 border-b border-border-base bg-slate-50/50">
				<div class="text-[10px] font-bold tracking-widest text-slate-400 uppercase">
					ADMIN LOGIN
				</div>
				<h2 class="text-2xl font-bold text-inst-navy font-serif leading-tight mt-1">
					Welcome Admin
				</h2>
				<p class="text-slate-500 text-xs mt-1">Login to access the iSPARC Administration Portal</p>
			</div>

			<form onsubmit={handleSubmit} class="p-6 sm:p-8 flex flex-col gap-6">
				<!-- Admin ID Input -->
				<div class="flex flex-col gap-1.5">
					<label
						for="{formId}-admin-id"
						class="text-[11px] font-bold text-slate-700 tracking-wider"
					>
						ADMIN ID
					</label>
					<input
						id="{formId}-admin-id"
						type="text"
						bind:value={adminId}
						disabled={lockoutTimeLeft > 0}
						placeholder="Enter your administrator ID"
						autoComplete="username"
						class="w-full px-3.5 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
						required
					/>
				</div>

				<!-- Password Input -->
				<div class="flex flex-col gap-1.5">
					<label
						for="{formId}-password"
						class="text-[11px] font-bold text-slate-700 tracking-wider"
					>
						PASSWORD
					</label>
					<div class="relative flex items-center">
						<input
							id="{formId}-password"
							type={showPassword ? 'text' : 'password'}
							bind:value={password}
							disabled={lockoutTimeLeft > 0}
							placeholder="Enter your password"
							autoComplete="current-password"
							class="w-full pl-3.5 pr-12 py-2.5 bg-white rounded-lg border border-border-base text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-inst-navy focus:ring-2 focus:ring-inst-navy/10 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
							required
						/>
						<button
							type="button"
							onclick={() => (showPassword = !showPassword)}
							class="absolute right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-650"
							aria-label={showPassword ? 'Hide password' : 'Show password'}
						>
							{#if showPassword}
								<!-- Eye-off SVG -->
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-5 h-5"
								>
									<line x1="2" y1="2" x2="22" y2="22" /><path
										d="M17.547 17.547a8.553 8.553 0 0 1-5.547 1.953 8.8 8.8 0 0 1-8.547-5.5 10.87 10.87 0 0 1 1.761-3.239"
									/><path
										d="M9.88 4.22a8.8 8.8 0 0 1 1.62-.22 8.82 8.82 0 0 1 8.547 5.5 10.64 10.64 0 0 1-1.341 2.871"
									/><circle cx="12" cy="12" r="3" />
								</svg>
							{:else}
								<!-- Eye SVG -->
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-5 h-5"
								>
									<path
										d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0z"
									/><circle cx="12" cy="12" r="3" />
								</svg>
							{/if}
						</button>
					</div>
				</div>

				<!-- Checkbox Options & Forgot Password -->
				<div class="flex items-center justify-between text-xs pt-1">
					<div></div>

					<button
						type="button"
						class="font-bold text-acad-red hover:underline transition-colors focus:outline-none"
					>
						Forgot Password?
					</button>
				</div>

				{#if errorMsg}
					<div
						class="p-3 bg-red-50 border border-red-200 text-red-700 text-xs rounded-lg font-semibold animate-fade-in"
						transition:slide
					>
						{errorMsg}
					</div>
				{/if}

				<!-- Submit Button -->
				<button
					type="submit"
					disabled={submitting || !isFormValid || lockoutTimeLeft > 0}
					class="w-full py-3.5 bg-inst-navy hover:bg-inst-navy/90 text-white font-bold text-sm tracking-widest uppercase rounded-xl transition duration-200 shadow-sm disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center"
				>
					{#if submitting}
						<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
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
						LOGGING IN...
					{:else}
						LOGIN
					{/if}
				</button>
			</form>
		{/if}
	</main>
</section>

<style>
	.animate-fade-in {
		opacity: 0;
		transform: translateY(-10px);
		animation: fadeIn 0.8s ease-out forwards;
	}

	.animate-fade-up {
		opacity: 0;
		transform: translateY(25px);
		animation: fadeUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) forwards;
		animation-delay: 0.15s;
	}

	@keyframes fadeIn {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes fadeUp {
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
