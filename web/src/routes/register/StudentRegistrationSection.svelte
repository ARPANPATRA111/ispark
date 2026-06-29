<script lang="ts">
	import { fade, slide } from 'svelte/transition';

	const formId = 'student-reg';

	// Reactive state for the registration form fields
	let formData = $state({
		rollNumber: '',
		enrollmentNumber: '',
		course: '',
		semester: '',
		fullName: '',
		email: '',
		mobileNumber: '',
		password: '',
		confirmPassword: '',
		consent: false
	});

	// Visbility toggles for password fields
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);

	// OTP flow simulation states
	let otpSent = $state(false);
	let otpSending = $state(false);
	let enteredOtp = $state('');
	let otpVerified = $state(false);
	let otpError = $state('');

	// Final Registration state
	let registered = $state(false);
	let submitting = $state(false);

	// Password Validation Runes
	let hasMinLength = $derived(formData.password.length >= 8);
	let hasUppercase = $derived(/[A-Z]/.test(formData.password));
	let hasLowercase = $derived(/[a-z]/.test(formData.password));
	let hasNumber = $derived(/[0-9]/.test(formData.password));
	let hasSpecialChar = $derived(/[^A-Za-z0-9]/.test(formData.password));

	let allPasswordRequirementsMet = $derived(
		hasMinLength && hasUppercase && hasLowercase && hasNumber && hasSpecialChar
	);

	// Check if confirmation match
	let passwordsMatch = $derived(
		formData.password !== '' && formData.password === formData.confirmPassword
	);

	// Send OTP validation state
	let isOtpDisabled = $derived(
		formData.email.trim() === '' ||
			!formData.email.includes('@') ||
			!formData.email.endsWith('.com')
	);

	const featureItems = [
		{ text: 'Participate in approved activities' },
		{ text: 'Upload and verify certificates' },
		{ text: 'Earn extracurricular credits' },
		{ text: 'View institutional leaderboard' },
		{ text: 'Download extracurricular marksheet' },
		{ text: 'Track semester-wise progress' }
	];

	const guidelineItems = [
		'Roll Number and Enrollment Number must match institute records.',
		'Only one account is allowed per student.',
		'Email verification is mandatory.',
		'Password must meet institutional security standards.'
	];

	const passwordRequirementsList = $derived([
		{ text: 'Minimum 8 characters', met: hasMinLength },
		{ text: 'One uppercase letter', met: hasUppercase },
		{ text: 'One lowercase letter', met: hasLowercase },
		{ text: 'One number', met: hasNumber },
		{ text: 'One special character', met: hasSpecialChar }
	]);

	// Form handlers
	function handleReset() {
		formData = {
			rollNumber: '',
			enrollmentNumber: '',
			course: '',
			semester: '',
			fullName: '',
			email: '',
			mobileNumber: '',
			password: '',
			confirmPassword: '',
			consent: false
		};
		showPassword = false;
		showConfirmPassword = false;
		otpSent = false;
		otpSending = false;
		enteredOtp = '';
		otpVerified = false;
		otpError = '';
		registered = false;
	}

	function handleSendOtp() {
		if (isOtpDisabled) return;
		otpSending = true;
		otpError = '';

		// Simulate API call for sending OTP
		setTimeout(() => {
			otpSending = false;
			otpSent = true;
		}, 1500);
	}

	function handleVerifyOtp() {
		otpError = '';
		// Accept '1234' or '123456' as simulated valid verification OTP code
		if (enteredOtp === '1234' || enteredOtp === '123456') {
			otpVerified = true;
		} else {
			otpError = "Invalid verification code. Enter '1234' to verify.";
		}
	}

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		if (!allPasswordRequirementsMet) {
			alert('Please ensure your password meets all requirements.');
			return;
		}

		if (!passwordsMatch) {
			alert('Passwords do not match!');
			return;
		}

		if (formData.email && !otpVerified) {
			alert('Please verify your email address by confirming the OTP.');
			return;
		}

		submitting = true;

		// Simulate API Submit
		setTimeout(() => {
			submitting = false;
			registered = true;
		}, 1800);
	}
</script>

<section
	class="max-w-6xl w-full mx-auto grid grid-cols-1 lg:grid-cols-12 gap-10 px-4 sm:px-6 lg:px-8 py-12 items-start font-sans"
>
	<!-- Left Column: Marketing / Guidelines -->
	<aside class="lg:col-span-5 flex flex-col gap-6 animate-fade-in">
		<div class="flex flex-col gap-4">
			<div>
				<h1 class="text-4xl font-extrabold text-acad-red font-serif leading-tight tracking-tight">
					Create Your<br />
					<span class="text-inst-navy">iSPARC Account</span>
				</h1>
			</div>

			<p class="text-slate-500 text-sm leading-relaxed">
				Register using your official institute credentials to participate in extracurricular
				activities, upload certificates, track credits, and monitor your overall iSPARC progress.
			</p>
		</div>

		<!-- Feature list -->
		<div class="flex flex-col gap-3.5 bg-white p-6 rounded-xl border border-slate-200 shadow-xs">
			<h3 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-1">
				Account Features
			</h3>
			{#each featureItems as item}
				<div class="flex items-start gap-3 group">
					<div
						class="flex w-5 h-5 items-center justify-center shrink-0 bg-emerald-50 rounded-full border border-emerald-250 transition-colors group-hover:bg-emerald-100"
					>
						<!-- Green checkmark SVG -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="3"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="w-3.5 h-3.5 text-emerald-600"
						>
							<polyline points="20 6 9 17 4 12" />
						</svg>
					</div>
					<span
						class="text-[13px] text-slate-700 font-medium group-hover:text-slate-900 transition-colors"
					>
						{item.text}
					</span>
				</div>
			{/each}
		</div>

		<!-- Security notice -->
		<div class="flex items-start gap-3.5 p-4 bg-amber-50/50 rounded-xl border border-amber-200/60">
			<div class="shrink-0 p-1 bg-amber-100/80 rounded-lg text-amber-800">
				<!-- Shield-alert SVG -->
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
					<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
					<line x1="12" y1="8" x2="12" y2="12" />
					<line x1="12" y1="16" x2="12.01" y2="16" />
				</svg>
			</div>
			<div>
				<h4 class="text-xs font-semibold text-slate-800 leading-tight">Security Notice</h4>
				<p class="text-[11px] text-slate-600 mt-1 leading-relaxed">
					Registration is restricted to officially enrolled IIPS students. Your information is
					verified before account activation.
				</p>
			</div>
		</div>

		<!-- Guidelines notice -->
		<section class="flex flex-col gap-3 p-5 bg-white rounded-xl border border-slate-200 shadow-xs">
			<h2 class="text-xs font-bold text-slate-400 tracking-widest uppercase">
				REGISTRATION GUIDELINES
			</h2>
			<ul class="flex flex-col gap-2.5 list-none">
				{#each guidelineItems as item}
					<li class="flex items-start gap-2.5 text-xs text-slate-600">
						<span class="text-acad-red font-bold text-sm leading-none">•</span>
						<p class="leading-relaxed">{item}</p>
					</li>
				{/each}
			</ul>
		</section>
	</aside>

	<!-- Right Column: Registration Card -->
	<main
		class="lg:col-span-7 bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden animate-fade-up"
	>
		{#if registered}
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
					<h2 class="text-2xl font-bold text-slate-900 font-serif">Registration Submitted!</h2>
					<p class="text-slate-500 text-sm mt-2 max-w-md mx-auto">
						Thank you, <span class="font-semibold text-slate-800"
							>{formData.fullName || 'Student'}</span
						>. Your iSPARC account registration details are successfully recorded and are now
						undergoing verification with IIPS administrative records.
					</p>
				</div>

				<div
					class="w-full bg-slate-50 p-5 rounded-xl border border-slate-100 text-left flex flex-col gap-2.5 max-w-md mx-auto text-xs text-slate-600"
				>
					<div class="flex justify-between border-b border-slate-200 pb-1.5">
						<span class="font-medium text-slate-400">Roll Number:</span>
						<span class="font-semibold text-slate-800">{formData.rollNumber}</span>
					</div>
					<div class="flex justify-between border-b border-slate-200 pb-1.5">
						<span class="font-medium text-slate-400">Enrollment:</span>
						<span class="font-semibold text-slate-800">{formData.enrollmentNumber}</span>
					</div>
					<div class="flex justify-between border-b border-slate-200 pb-1.5">
						<span class="font-medium text-slate-400">Course / Sem:</span>
						<span class="font-semibold text-slate-800 uppercase"
							>{formData.course} (Sem {formData.semester})</span
						>
					</div>
					<div class="flex justify-between">
						<span class="font-medium text-slate-400">Email:</span>
						<span class="font-semibold text-slate-800">{formData.email}</span>
					</div>
				</div>

				<div class="flex flex-col sm:flex-row gap-3 w-full max-w-md">
					<a
						href="/"
						class="flex-1 py-3 text-center bg-inst-navy hover:bg-inst-navy/90 text-white font-semibold text-xs tracking-wider uppercase rounded-xl transition duration-200"
					>
						Back to Home
					</a>
					<button
						type="button"
						onclick={handleReset}
						class="flex-1 py-3 border border-slate-200 hover:bg-slate-50 text-slate-700 font-semibold text-xs tracking-wider uppercase rounded-xl transition duration-200"
					>
						Register Another
					</button>
				</div>
			</div>
		{:else}
			<!-- Form View -->
			<div class="p-6 sm:p-8 border-b border-slate-200 bg-slate-50/50">
				<h2 class="text-2xl font-bold text-inst-navy font-serif leading-tight">
					Student Registration
				</h2>
				<p class="text-slate-500 text-xs mt-1">
					Complete the form below using verified credentials to set up your profile.
				</p>
			</div>

			<form onsubmit={handleSubmit} class="p-6 sm:p-8 flex flex-col gap-6">
				<!-- SECTION 1: ACADEMIC INFORMATION -->
				<div class="flex flex-col gap-3">
					<div class="text-[10px] font-bold tracking-widest text-slate-400 uppercase">
						ACADEMIC INFORMATION
					</div>

					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<!-- Roll Number -->
						<div class="flex flex-col gap-1.5">
							<label
								for="{formId}-rollNumber"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								ROLL NUMBER <span class="text-acad-red">*</span>
							</label>
							<input
								id="{formId}-rollNumber"
								type="text"
								bind:value={formData.rollNumber}
								placeholder="Enter Roll Number"
								class="w-full px-3.5 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200"
								required
							/>
						</div>

						<!-- Enrollment Number -->
						<div class="flex flex-col gap-1.5">
							<label
								for="{formId}-enrollmentNumber"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								ENROLLMENT NUMBER <span class="text-acad-red">*</span>
							</label>
							<input
								id="{formId}-enrollmentNumber"
								type="text"
								bind:value={formData.enrollmentNumber}
								placeholder="Enter Enrollment Number"
								class="w-full px-3.5 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200"
								required
							/>
						</div>

						<!-- Course -->
						<div class="flex flex-col gap-1.5">
							<label
								for="{formId}-course"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								COURSE <span class="text-acad-red">*</span>
							</label>
							<div class="relative">
								<select
									id="{formId}-course"
									bind:value={formData.course}
									class="w-full px-3.5 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200 appearance-none"
									required
								>
									<option value="" disabled class="text-slate-450">Select Course</option>
									<option value="MBA (Management Science - MS)"
										>MBA (Management Science - MS)</option
									>
									<option value="MCA (Master of Computer Applications)"
										>MCA (Master of Computer Applications)</option
									>
									<option value="M.Tech (Information Technology - IT)"
										>M.Tech (Information Technology - IT)</option
									>
									<option value="M.Tech (Computer Science - CS)"
										>M.Tech (Computer Science - CS)</option
									>
									<option value="MBA (Management Science)">MBA (Management Science)</option>
									<option value="MBA (Advertising and Public Relations - APR)"
										>MBA (Advertising and Public Relations - APR)</option
									>
									<option value="MBA (Entrepreneurship)">MBA (Entrepreneurship)</option>
									<option value="B.Com. (Hons.)">B.Com. (Hons.)</option>
								</select>
								<div
									class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-3.5 text-slate-400"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 20 20"
										fill="currentColor"
										class="w-4 h-4"
									>
										<path
											fill-rule="evenodd"
											d="M5.22 8.22a.75.75 0 0 1 1.06 0L10 11.94l3.72-3.72a.75.75 0 1 1 1.06 1.06l-4.25 4.25a.75.75 0 0 1-1.06 0L5.22 9.28a.75.75 0 0 1 0-1.06Z"
											clip-rule="evenodd"
										/>
									</svg>
								</div>
							</div>
						</div>

						<!-- Semester -->
						<div class="flex flex-col gap-1.5">
							<label
								for="{formId}-semester"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								SEMESTER <span class="text-acad-red">*</span>
							</label>
							<div class="relative">
								<select
									id="{formId}-semester"
									bind:value={formData.semester}
									class="w-full px-3.5 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200 appearance-none"
									required
								>
									<option value="" disabled class="text-slate-450">Select Semester</option>
									{#each Array(10) as _, i}
										<option value={i + 1}>Semester {i + 1}</option>
									{/each}
								</select>
								<div
									class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-3.5 text-slate-400"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 20 20"
										fill="currentColor"
										class="w-4 h-4"
									>
										<path
											fill-rule="evenodd"
											d="M5.22 8.22a.75.75 0 0 1 1.06 0L10 11.94l3.72-3.72a.75.75 0 1 1 1.06 1.06l-4.25 4.25a.75.75 0 0 1-1.06 0L5.22 9.28a.75.75 0 0 1 0-1.06Z"
											clip-rule="evenodd"
										/>
									</svg>
								</div>
							</div>
						</div>
					</div>
				</div>

				<div class="h-px bg-slate-250"></div>

				<!-- SECTION 2: PERSONAL INFORMATION -->
				<div class="flex flex-col gap-4">
					<div class="text-[10px] font-bold tracking-widest text-slate-400 uppercase">
						PERSONAL INFORMATION
					</div>

					<div class="flex flex-col gap-4">
						<!-- Full Name -->
						<div class="flex flex-col gap-1.5">
							<label
								for="{formId}-fullName"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								FULL NAME <span class="text-acad-red">*</span>
							</label>
							<input
								id="{formId}-fullName"
								type="text"
								bind:value={formData.fullName}
								placeholder="Enter your full name"
								class="w-full px-3.5 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200"
								required
							/>
						</div>

						<!-- Email Address -->
						<div class="flex flex-col gap-1.5">
							<label
								for="{formId}-email"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								EMAIL <span class="text-acad-red">*</span>
							</label>
							<div class="flex gap-2">
								<input
									id="{formId}-email"
									type="email"
									bind:value={formData.email}
									disabled={otpVerified}
									placeholder="student@iips.com"
									class="flex-1 px-3.5 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200 disabled:bg-slate-100 disabled:text-slate-500"
									required
								/>

								{#if !otpVerified}
									<button
										type="button"
										onclick={handleSendOtp}
										disabled={isOtpDisabled || otpSending}
										class="shrink-0 px-4 py-2.5 bg-inst-navy hover:bg-inst-navy/90 text-white font-semibold text-xs tracking-wider uppercase rounded-lg disabled:opacity-40 transition duration-200 min-w-[100px] flex items-center justify-center"
									>
										{#if otpSending}
											<!-- Spinner SVG -->
											<svg
												class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
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
										{:else if otpSent}
											RESEND OTP
										{:else}
											SEND OTP
										{/if}
									</button>
								{/if}
							</div>
							<p class="text-[10px] text-slate-400 leading-normal">
								An activation code will be sent to this email address to verify your identity.
							</p>
						</div>

						<!-- OTP Verification Code Entry (Conditional) -->
						{#if otpSent && !otpVerified}
							<div
								class="p-4 bg-slate-50 border border-slate-200 rounded-xl flex flex-col gap-2.5"
								transition:slide={{ duration: 250 }}
							>
								<div class="text-[11px] font-bold text-slate-700 tracking-wider">
									ENTER VERIFICATION CODE
								</div>
								<div class="flex gap-2">
									<input
										type="text"
										bind:value={enteredOtp}
										placeholder="Enter Code (Use '1234' for testing)"
										class="flex-grow px-3.5 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200"
									/>
									<button
										type="button"
										onclick={handleVerifyOtp}
										class="px-4 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold text-xs tracking-wider uppercase rounded-lg transition duration-200"
									>
										VERIFY CODE
									</button>
								</div>
								{#if otpError}
									<div class="text-xs text-rose-600 font-medium">{otpError}</div>
								{/if}
							</div>
						{/if}

						{#if otpVerified}
							<!-- OTP Verified indicator -->
							<div
								class="flex items-center gap-2 text-emerald-600 text-xs font-semibold p-2 bg-emerald-50 rounded-lg border border-emerald-150"
								transition:slide={{ duration: 200 }}
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="3"
									stroke-linecap="round"
									stroke-linejoin="round"
									class="w-4 h-4"
								>
									<polyline points="20 6 9 17 4 12" />
								</svg>
								<span>Email Address Verified Successfully</span>
							</div>
						{/if}

						<!-- Mobile Number -->
						<div class="flex flex-col gap-1.5">
							<label
								for="{formId}-mobileNumber"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								MOBILE NUMBER <span class="text-acad-red">*</span>
							</label>
							<input
								id="{formId}-mobileNumber"
								type="tel"
								bind:value={formData.mobileNumber}
								placeholder="+91 XXXXX XXXXX"
								class="w-full px-3.5 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200"
								required
							/>
						</div>
					</div>
				</div>

				<div class="h-px bg-slate-255"></div>

				<!-- SECTION 3: ACCOUNT SETUP -->
				<div class="flex flex-col gap-4">
					<div class="text-[10px] font-bold tracking-widest text-slate-400 uppercase">
						ACCOUNT SETUP
					</div>

					<div class="flex flex-col gap-4.5">
						<!-- Password -->
						<div class="flex flex-col gap-1.5 relative">
							<label
								for="{formId}-password"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								PASSWORD <span class="text-acad-red">*</span>
							</label>
							<div class="relative">
								<input
									id="{formId}-password"
									type={showPassword ? 'text' : 'password'}
									bind:value={formData.password}
									placeholder="Create a strong password"
									class="w-full pl-3.5 pr-12 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200"
									required
								/>
								<button
									type="button"
									aria-label={showPassword ? 'Hide password' : 'Show password'}
									onclick={() => (showPassword = !showPassword)}
									class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-650"
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

							<!-- Password strength guidelines - reactive checker -->
							<div class="mt-2.5 p-3.5 bg-slate-50 border border-slate-150 rounded-lg">
								<div class="text-[10px] font-bold text-slate-400 uppercase mb-2">
									Password Requirements
								</div>
								<div class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-2">
									{#each passwordRequirementsList as requirement}
										<div class="flex items-center gap-2">
											<div
												class="w-4 h-4 flex items-center justify-center rounded-full border transition-all duration-200 {requirement.met
													? 'bg-emerald-50 border-emerald-250 text-emerald-600'
													: 'bg-slate-100 border-slate-200 text-slate-350'}"
											>
												<!-- Tiny checkmark -->
												<svg
													xmlns="http://www.w3.org/2000/svg"
													viewBox="0 0 24 24"
													fill="none"
													stroke="currentColor"
													stroke-width="4"
													stroke-linecap="round"
													stroke-linejoin="round"
													class="w-2.5 h-2.5 {requirement.met ? 'opacity-100' : 'opacity-20'}"
												>
													<polyline points="20 6 9 17 4 12" />
												</svg>
											</div>
											<span
												class="text-[11px] font-medium transition-all duration-200 {requirement.met
													? 'text-slate-650'
													: 'text-slate-400'}"
											>
												{requirement.text}
											</span>
										</div>
									{/each}
								</div>
							</div>
						</div>

						<!-- Confirm Password -->
						<div class="flex flex-col gap-1.5 relative">
							<label
								for="{formId}-confirmPassword"
								class="text-[11px] font-bold text-slate-700 tracking-wider"
							>
								CONFIRM PASSWORD <span class="text-acad-red">*</span>
							</label>
							<div class="relative">
								<input
									id="{formId}-confirmPassword"
									type={showConfirmPassword ? 'text' : 'password'}
									bind:value={formData.confirmPassword}
									placeholder="Re-enter your password"
									class="w-full pl-3.5 pr-12 py-2.5 bg-white rounded-lg border border-slate-250 text-[13px] text-slate-800 placeholder:text-slate-400 focus:outline-none focus:border-acad-red focus:ring-2 focus:ring-acad-red/10 transition-all duration-200"
									required
								/>
								<button
									type="button"
									aria-label={showConfirmPassword
										? 'Hide confirm password'
										: 'Show confirm password'}
									onclick={() => (showConfirmPassword = !showConfirmPassword)}
									class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-650"
								>
									{#if showConfirmPassword}
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

							<!-- Password mismatch warning -->
							{#if formData.confirmPassword !== '' && !passwordsMatch}
								<div
									class="text-xs text-rose-600 font-medium mt-1.5"
									transition:slide={{ duration: 150 }}
								>
									Passwords do not match.
								</div>
							{:else if formData.confirmPassword !== '' && passwordsMatch}
								<div
									class="text-xs text-emerald-600 font-medium mt-1.5 flex items-center gap-1"
									transition:slide={{ duration: 150 }}
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="3"
										stroke-linecap="round"
										stroke-linejoin="round"
										class="w-3 h-3"
									>
										<polyline points="20 6 9 17 4 12" />
									</svg>
									<span>Passwords match!</span>
								</div>
							{/if}
						</div>
					</div>
				</div>

				<div class="h-px bg-slate-260"></div>

				<!-- SECTION 4: CONSENT & SUBMISSION -->
				<div class="flex flex-col gap-5">
					<label class="flex gap-3 cursor-pointer select-none group items-start">
						<div class="pt-0.5 relative shrink-0">
							<input
								type="checkbox"
								bind:checked={formData.consent}
								class="w-4 h-4 rounded border-slate-250 text-acad-red focus:ring-acad-red focus:ring-opacity-20 cursor-pointer accent-acad-red"
								required
							/>
						</div>
						<span
							class="text-xs text-slate-600 group-hover:text-slate-800 transition-colors leading-relaxed"
						>
							I certify that the information provided is accurate and belongs to me. I agree to
							abide by the iSPARC participation guidelines and institutional policies.
						</span>
					</label>

					<!-- Buttons -->
					<div class="flex flex-col gap-3 w-full">
						<button
							type="submit"
							disabled={submitting ||
								!formData.consent ||
								!allPasswordRequirementsMet ||
								!passwordsMatch ||
								(formData.email !== '' && !otpVerified)}
							class="w-full py-3.5 bg-acad-red hover:bg-acad-red/90 text-white font-bold text-sm tracking-widest uppercase rounded-xl transition duration-200 shadow-sm disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center"
						>
							{#if submitting}
								<svg
									class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
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
								REGISTERING ACCOUNT...
							{:else}
								REGISTER ACCOUNT
							{/if}
						</button>

						<button
							type="button"
							onclick={handleReset}
							class="w-full py-3 bg-white border border-slate-250 text-slate-700 font-semibold text-xs tracking-wider uppercase rounded-xl hover:bg-slate-50 transition duration-200"
						>
							Reset Form
						</button>
					</div>

					<!-- Login Redirect link -->
					<div class="text-center text-xs text-slate-500 mt-2">
						Already have an account?
						<a
							href="/login"
							class="font-bold text-acad-red hover:underline hover:text-acad-red/90 transition-colors ml-1"
						>
							Login Here
						</a>
					</div>
				</div>
			</form>
		{/if}
	</main>
</section>

<style>
	/* Rich, Sleek Animations */
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
