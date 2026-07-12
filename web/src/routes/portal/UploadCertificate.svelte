<script lang="ts">
	import { slide } from 'svelte/transition';
	import { API_BASE_URL } from '$lib/config';

	let token = localStorage.getItem('access_token') || '';

	// Props using Svelte 5 runes
	let {
		onBackToDashboard
	}: {
		onBackToDashboard: () => void;
	} = $props();

	// Form field states
	let activityName = $state('');
	let activityCategory = $state('');
	let activityDate = $state('');
	let organizerName = $state('');
	let eventLevel = $state('');
	let certNumber = $state('');
	let issueDate = $state('');
	let participationType = $state('');
	let description = $state('');

	// Upload file states
	let fileInput: HTMLInputElement;
	let selectedFile = $state<File | null>(null);
	let uploadProgress = $state(0);
	let isDragging = $state(false);
	let isUploading = $state(false);
	let uploadSuccess = $state(false);

	// Error / Success notification states
	let errorMessage = $state('');
	let successMessage = $state('');

	// Auto-fill category simulation based on common words in activity name
	$effect(() => {
		if (activityName.trim() === '') {
			activityCategory = '';
			return;
		}
		const lowerName = activityName.toLowerCase();
		if (lowerName.includes('debate') || lowerName.includes('speech') || lowerName.includes('mun')) {
			activityCategory = 'Literary';
		} else if (
			lowerName.includes('dance') ||
			lowerName.includes('drama') ||
			lowerName.includes('singing') ||
			lowerName.includes('music')
		) {
			activityCategory = 'Cultural';
		} else if (
			lowerName.includes('hackathon') ||
			lowerName.includes('coding') ||
			lowerName.includes('robotics')
		) {
			activityCategory = 'Technical';
		} else if (
			lowerName.includes('blood') ||
			lowerName.includes('nss') ||
			lowerName.includes('cleanup')
		) {
			activityCategory = 'Social Service';
		} else if (
			lowerName.includes('olympiad') ||
			lowerName.includes('quiz') ||
			lowerName.includes('paper')
		) {
			activityCategory = 'Academic';
		} else if (
			lowerName.includes('cricket') ||
			lowerName.includes('football') ||
			lowerName.includes('sports') ||
			lowerName.includes('athletics')
		) {
			activityCategory = 'Sports & Games';
		} else {
			activityCategory = 'General / Extracurricular';
		}
	});

	// Trigger simulated upload progress
	function startSimulatedUpload(file: File) {
		selectedFile = file;
		isUploading = true;
		uploadSuccess = false;
		uploadProgress = 0;
		errorMessage = '';

		const interval = setInterval(() => {
			uploadProgress += 10;
			if (uploadProgress >= 100) {
				clearInterval(interval);
				isUploading = false;
				uploadSuccess = true;
			}
		}, 100);
	}

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files[0]) {
			const file = target.files[0];
			const validTypes = ['application/pdf', 'image/png', 'image/jpeg', 'image/jpg'];
			if (!validTypes.includes(file.type)) {
				errorMessage = 'Unsupported file format. Please upload PDF, PNG, or JPG.';
				return;
			}
			if (file.size > 5 * 1024 * 1024) {
				errorMessage = 'File size exceeds 5 MB limit.';
				return;
			}
			startSimulatedUpload(file);
		}
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		isDragging = true;
	}

	function handleDragLeave() {
		isDragging = false;
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		isDragging = false;
		if (event.dataTransfer && event.dataTransfer.files && event.dataTransfer.files[0]) {
			const file = event.dataTransfer.files[0];
			const validTypes = ['application/pdf', 'image/png', 'image/jpeg', 'image/jpg'];
			if (!validTypes.includes(file.type)) {
				errorMessage = 'Unsupported file format. Please upload PDF, PNG, or JPG.';
				return;
			}
			if (file.size > 5 * 1024 * 1024) {
				errorMessage = 'File size exceeds 5 MB limit.';
				return;
			}
			startSimulatedUpload(file);
		}
	}

	function triggerFilePicker() {
		fileInput.click();
	}

	function removeFile() {
		selectedFile = null;
		uploadProgress = 0;
		uploadSuccess = false;
		isUploading = false;
		errorMessage = '';
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();
		errorMessage = '';
		successMessage = '';

		// Validation check
		if (!activityName || !activityDate || !organizerName || !eventLevel) {
			errorMessage = 'Please fill in all required Activity Information fields.';
			return;
		}
		if (!selectedFile || !uploadSuccess) {
			errorMessage = 'Please upload a valid certificate file.';
			return;
		}
		if (!participationType) {
			errorMessage = 'Please select a Participation Type.';
			return;
		}

		isUploading = true;

		const formData = new FormData();
		formData.append('activity_name', activityName);
		formData.append('activity_category', activityCategory);
		formData.append('activity_date', activityDate);
		formData.append('organizer_name', organizerName);
		formData.append('event_level', eventLevel);
		formData.append('cert_number', certNumber);
		formData.append('issue_date', issueDate);
		formData.append('participation_type', participationType);
		formData.append('description', description);
		formData.append('certificate_file', selectedFile);

		try {
			const res = await fetch(`${API_BASE_URL}/api/student/certificates`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${token}`
				},
				body: formData
			});

			if (!res.ok) {
				const errorData = await res.json();
				throw new Error(errorData.error || 'Failed to upload certificate.');
			}

			successMessage = 'Certificate submitted successfully for review!';

			// Clear local draft upon successful submission
			localStorage.removeItem('certificate_draft');

			// Reset form
			activityName = '';
			activityCategory = '';
			activityDate = '';
			organizerName = '';
			eventLevel = '';
			certNumber = '';
			issueDate = '';
			participationType = '';
			description = '';
			removeFile();

			setTimeout(() => {
				successMessage = '';
				onBackToDashboard();
			}, 2000);
		} catch (err) {
			console.error(err);
			errorMessage =
				err instanceof Error ? err.message : 'Error uploading certificate. Please try again.';
		} finally {
			isUploading = false;
		}
	}

	// Load draft from localStorage on mount
	$effect(() => {
		const saved = localStorage.getItem('certificate_draft');
		if (saved) {
			try {
				const draft = JSON.parse(saved);
				activityName = draft.activityName || '';
				activityCategory = draft.activityCategory || '';
				activityDate = draft.activityDate || '';
				organizerName = draft.organizerName || '';
				eventLevel = draft.eventLevel || '';
				certNumber = draft.certNumber || '';
				issueDate = draft.issueDate || '';
				participationType = draft.participationType || '';
				description = draft.description || '';
			} catch (e) {
				console.error('Error parsing draft:', e);
			}
		}
	});

	function handleSaveDraft() {
		const draft = {
			activityName,
			activityCategory,
			activityDate,
			organizerName,
			eventLevel,
			certNumber,
			issueDate,
			participationType,
			description
		};
		localStorage.setItem('certificate_draft', JSON.stringify(draft));
		successMessage = 'Draft saved successfully!';
		setTimeout(() => {
			successMessage = '';
		}, 2000);
	}

	function handleCancel() {
		removeFile();
		onBackToDashboard();
	}
</script>

<div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start font-sans">
	<!-- Left Column: Certificate Submission Form -->
	<form
		onsubmit={handleSubmit}
		class="lg:col-span-8 bg-white border border-slate-200 rounded-xl p-6 sm:p-8 shadow-xs space-y-6"
	>
		<div>
			<h2 class="text-xl font-bold font-serif text-[#0B1535] leading-tight">
				Certificate Submission Form
			</h2>
			<div class="h-px bg-slate-100 my-4"></div>
		</div>

		<!-- Notifications -->
		{#if errorMessage}
			<div
				transition:slide={{ duration: 150 }}
				class="p-4 bg-rose-50 border border-rose-200 text-rose-800 text-xs font-bold rounded-lg flex items-center gap-2"
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
				<span>{errorMessage}</span>
			</div>
		{/if}

		{#if successMessage}
			<div
				transition:slide={{ duration: 150 }}
				class="p-4 bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs font-bold rounded-lg flex items-center gap-2"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="2"
					stroke="currentColor"
					class="w-4 h-4"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
				</svg>
				<span>{successMessage}</span>
			</div>
		{/if}

		<!-- SECTION 1: ACTIVITY INFORMATION -->
		<section class="space-y-4">
			<h3 class="text-xs font-bold text-slate-400 uppercase tracking-widest">
				Activity Information
			</h3>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div class="flex flex-col gap-1.5">
					<label for="activity-name" class="text-[11px] font-bold text-slate-700 tracking-wider"
						>ACTIVITY NAME *</label
					>
					<input
						id="activity-name"
						type="text"
						bind:value={activityName}
						placeholder="e.g. National Hackathon, Robotics Workshop"
						class="px-3 py-2 bg-white rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-[#0B1535] focus:ring-2 focus:ring-[#0B1535]/5 transition-all"
						required
					/>
				</div>
				<div class="flex flex-col gap-1.5">
					<label for="activity-category" class="text-[11px] font-bold text-slate-700 tracking-wider"
						>ACTIVITY CATEGORY</label
					>
					<input
						id="activity-category"
						type="text"
						bind:value={activityCategory}
						placeholder="Auto-filled after selection"
						class="px-3 py-2 bg-slate-50 rounded-lg border border-slate-200 text-xs text-slate-400 focus:outline-none select-all"
						disabled
					/>
				</div>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div class="flex flex-col gap-1.5">
					<label for="activity-date" class="text-[11px] font-bold text-slate-700 tracking-wider"
						>ACTIVITY DATE *</label
					>
					<input
						id="activity-date"
						type="date"
						bind:value={activityDate}
						class="px-3 py-2 bg-white rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-[#0B1535] focus:ring-2 focus:ring-[#0B1535]/5 transition-all"
						required
					/>
				</div>
				<div class="flex flex-col gap-1.5">
					<label for="organizer-name" class="text-[11px] font-bold text-slate-700 tracking-wider"
						>ORGANIZER NAME *</label
					>
					<input
						id="organizer-name"
						type="text"
						bind:value={organizerName}
						placeholder="Enter organizing institution"
						class="px-3 py-2 bg-white rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-[#0B1535] focus:ring-2 focus:ring-[#0B1535]/5 transition-all"
						required
					/>
				</div>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="event-level" class="text-[11px] font-bold text-slate-700 tracking-wider"
					>EVENT LEVEL *</label
				>
				<select
					id="event-level"
					bind:value={eventLevel}
					class="px-3 py-2 bg-white rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-[#0B1535] focus:ring-2 focus:ring-[#0B1535]/5 transition-all"
					required
				>
					<option value="" disabled selected>Select level</option>
					<option value="International">International</option>
					<option value="National">National</option>
					<option value="State">State</option>
					<option value="University">University</option>
					<option value="Department">Department</option>
				</select>
			</div>
		</section>

		<div class="h-px bg-slate-100 my-2"></div>

		<!-- SECTION 2: UPLOAD CERTIFICATE -->
		<section class="space-y-4">
			<h3 class="text-xs font-bold text-slate-400 uppercase tracking-widest">Upload Certificate</h3>

			<!-- Hidden native file picker -->
			<input
				type="file"
				accept=".pdf,.png,.jpg,.jpeg"
				bind:this={fileInput}
				onchange={handleFileSelect}
				class="hidden"
			/>

			{#if !selectedFile}
				<!-- Drag & Drop Zone -->
				<!-- svelte-ignore a11y_click_events_have_key_events -->
				<!-- svelte-ignore a11y_no_static_element_interactions -->
				<div
					ondragover={handleDragOver}
					ondragleave={handleDragLeave}
					ondrop={handleDrop}
					onclick={triggerFilePicker}
					class="border-2 border-dashed rounded-xl p-8 flex flex-col items-center justify-center gap-3 transition-colors cursor-pointer select-none
					{isDragging
						? 'border-[#0B1535] bg-[#0B1535]/5'
						: 'border-slate-300 hover:border-[#0B1535] hover:bg-slate-50'}"
				>
					<div class="p-3 bg-slate-50 border border-slate-100 rounded-full text-slate-400">
						<!-- Upload cloud icon -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.8"
							stroke="currentColor"
							class="w-6 h-6"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M12 16.5V9.75m0 0 3 3m-3-3-3 3M6.75 19.5a4.5 4.5 0 0 1-1.41-8.775 5.25 5.25 0 0 1 10.233-2.33 3 3 0 0 1 3.758 3.848A3.752 3.752 0 0 1 18 19.5H6.75Z"
							/>
						</svg>
					</div>
					<div class="text-center">
						<p class="text-xs text-slate-700 font-bold">Drag & Drop your certificate here</p>
						<p class="text-[10px] text-slate-400 font-bold uppercase tracking-widest mt-1">or</p>
					</div>
					<button
						type="button"
						class="px-4 py-2 bg-[#0B1535] hover:bg-[#0B1535]/95 text-white text-xs font-bold rounded-lg transition duration-200 focus:outline-none"
					>
						BROWSE FILES
					</button>
					<p class="text-[9px] font-bold text-slate-400 mt-1 uppercase tracking-wider">
						Supported: PDF, PNG, JPEG, JPG &middot; Max size: 5 MB
					</p>
				</div>
			{:else}
				<!-- File selected details -->
				<div
					transition:slide={{ duration: 150 }}
					class="p-4 border border-slate-200 rounded-xl bg-slate-50/50 flex items-center justify-between gap-3"
				>
					<div class="flex items-center gap-3 min-w-0">
						<div class="p-2.5 rounded-lg bg-rose-50 text-rose-600 border border-rose-100 shrink-0">
							<!-- PDF/Image icon -->
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
						<div class="flex-grow min-w-0">
							<div class="text-xs font-bold text-slate-800 truncate">{selectedFile.name}</div>
							<div class="text-[10px] text-slate-400 font-semibold mt-0.5">
								{(selectedFile.size / (1024 * 1024)).toFixed(2)} MB
							</div>
							{#if isUploading}
								<!-- Progress bar -->
								<div class="mt-2 space-y-1">
									<div class="h-1.5 w-full bg-slate-100 rounded-full overflow-hidden">
										<div
											class="h-full bg-accent-red rounded-full transition-all duration-100"
											style="width: {uploadProgress}%"
										></div>
									</div>
									<div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">
										Uploading {uploadProgress}%
									</div>
								</div>
							{:else if uploadSuccess}
								<div
									class="text-[9px] text-emerald-600 font-bold uppercase tracking-widest mt-1 flex items-center gap-1"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="3"
										stroke="currentColor"
										class="w-3.5 h-3.5"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="m4.5 12.75 6 6 9-13.5"
										/>
									</svg>
									Upload Complete
								</div>
							{/if}
						</div>
					</div>

					<button
						type="button"
						onclick={removeFile}
						class="text-xs font-bold text-rose-600 hover:text-rose-800 transition-colors focus:outline-none"
					>
						Remove
					</button>
				</div>
			{/if}
		</section>

		<div class="h-px bg-slate-100 my-2"></div>

		<!-- SECTION 3: CERTIFICATE INFORMATION -->
		<section class="space-y-4">
			<h3 class="text-xs font-bold text-slate-400 uppercase tracking-widest">
				Certificate Information
			</h3>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div class="flex flex-col gap-1.5">
					<label for="cert-number" class="text-[11px] font-bold text-slate-700 tracking-wider"
						>CERTIFICATE NUMBER</label
					>
					<input
						id="cert-number"
						type="text"
						bind:value={certNumber}
						placeholder="e.g. CERT-2025-001 (Optional)"
						class="px-3 py-2 bg-white rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-[#0B1535] focus:ring-2 focus:ring-[#0B1535]/5 transition-all"
					/>
				</div>
				<div class="flex flex-col gap-1.5">
					<label for="issue-date" class="text-[11px] font-bold text-slate-700 tracking-wider"
						>ISSUE DATE</label
					>
					<input
						id="issue-date"
						type="date"
						bind:value={issueDate}
						class="px-3 py-2 bg-white rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-[#0B1535] focus:ring-2 focus:ring-[#0B1535]/5 transition-all"
					/>
				</div>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="participation-type" class="text-[11px] font-bold text-slate-700 tracking-wider"
					>PARTICIPATION TYPE *</label
				>
				<select
					id="participation-type"
					bind:value={participationType}
					class="px-3 py-2 bg-white rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-[#0B1535] focus:ring-2 focus:ring-[#0B1535]/5 transition-all"
					required
				>
					<option value="" disabled selected>Select participation type</option>
					<option value="Winner">Winner</option>
					<option value="Runner-up">Runner-up</option>
					<option value="Participant">Participant</option>
					<option value="Organizer">Organizer</option>
				</select>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="achievement-desc" class="text-[11px] font-bold text-slate-700 tracking-wider"
					>ACHIEVEMENT DESCRIPTION</label
				>
				<textarea
					id="achievement-desc"
					bind:value={description}
					rows="4"
					placeholder="Provide a brief description of your participation or achievement."
					class="px-3 py-2 bg-white rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-[#0B1535] focus:ring-2 focus:ring-[#0B1535]/5 transition-all resize-y"
				></textarea>
			</div>
		</section>

		<div class="h-px bg-slate-100 my-2"></div>

		<!-- Action Buttons -->
		<div class="flex items-center gap-4 pt-2">
			<button
				type="submit"
				class="px-5 py-2.5 bg-[#0B1535] hover:bg-[#0B1535]/95 text-white font-bold text-xs tracking-wider uppercase rounded-lg transition duration-200 shadow-xs focus:outline-none"
			>
				Submit Certificate
			</button>

			<button
				type="button"
				onclick={handleSaveDraft}
				class="px-5 py-2.5 bg-white border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs tracking-wider uppercase rounded-lg transition duration-200 focus:outline-none"
			>
				Save Draft
			</button>

			<button
				type="button"
				onclick={handleCancel}
				class="text-xs font-bold text-slate-500 hover:text-slate-700 transition-colors focus:outline-none"
			>
				Cancel
			</button>
		</div>
	</form>

	<!-- Right Column: Guidelines and Integrity notice -->
	<aside class="lg:col-span-4 space-y-6">
		<!-- Guidelines Card -->
		<div class="bg-white border border-slate-200 rounded-xl p-5 shadow-xs">
			<h3 class="text-sm font-bold font-serif text-[#0B1535]">Submission Guidelines</h3>
			<div class="h-px bg-slate-100 my-3.5"></div>

			<ul class="space-y-3 font-sans">
				<li class="flex items-start gap-2.5 text-xs text-slate-600">
					<span
						class="w-4.5 h-4.5 bg-emerald-50 text-emerald-600 flex items-center justify-center rounded-full border border-emerald-100 shrink-0 mt-0.5"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="3"
							stroke="currentColor"
							class="w-2.5 h-2.5"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
						</svg>
					</span>
					<span class="leading-relaxed">Upload clear and readable certificates</span>
				</li>
				<li class="flex items-start gap-2.5 text-xs text-slate-600">
					<span
						class="w-4.5 h-4.5 bg-emerald-50 text-emerald-600 flex items-center justify-center rounded-full border border-emerald-100 shrink-0 mt-0.5"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="3"
							stroke="currentColor"
							class="w-2.5 h-2.5"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
						</svg>
					</span>
					<span class="leading-relaxed">Ensure your name is visible</span>
				</li>
				<li class="flex items-start gap-2.5 text-xs text-slate-600">
					<span
						class="w-4.5 h-4.5 bg-emerald-50 text-emerald-600 flex items-center justify-center rounded-full border border-emerald-100 shrink-0 mt-0.5"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="3"
							stroke="currentColor"
							class="w-2.5 h-2.5"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
						</svg>
					</span>
					<span class="leading-relaxed">Certificates must be authentic</span>
				</li>
				<li class="flex items-start gap-2.5 text-xs text-slate-600">
					<span
						class="w-4.5 h-4.5 bg-emerald-50 text-emerald-600 flex items-center justify-center rounded-full border border-emerald-100 shrink-0 mt-0.5"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="3"
							stroke="currentColor"
							class="w-2.5 h-2.5"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
						</svg>
					</span>
					<span class="leading-relaxed">Supported formats: PDF, JPG, PNG</span>
				</li>
				<li class="flex items-start gap-2.5 text-xs text-slate-600">
					<span
						class="w-4.5 h-4.5 bg-emerald-50 text-emerald-600 flex items-center justify-center rounded-full border border-emerald-100 shrink-0 mt-0.5"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="3"
							stroke="currentColor"
							class="w-2.5 h-2.5"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
						</svg>
					</span>
					<span class="leading-relaxed">File size must not exceed 5 MB</span>
				</li>
				<li class="flex items-start gap-2.5 text-xs text-slate-600">
					<span
						class="w-4.5 h-4.5 bg-emerald-50 text-emerald-600 flex items-center justify-center rounded-full border border-emerald-100 shrink-0 mt-0.5"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="3"
							stroke="currentColor"
							class="w-2.5 h-2.5"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
						</svg>
					</span>
					<span class="leading-relaxed">One certificate per submission</span>
				</li>
			</ul>
		</div>

		<!-- Integrity Notice Card -->
		<div
			class="bg-white border-l-4 border-rose-600 border border-slate-200 rounded-xl p-5 shadow-xs"
		>
			<h4 class="text-[11px] font-bold text-rose-700 uppercase tracking-widest font-sans">
				Academic Integrity Notice
			</h4>
			<p class="text-xs text-slate-650 leading-relaxed font-sans mt-2.5">
				Submission of forged or altered certificates may result in disciplinary action and permanent
				rejection of extracurricular credits.
			</p>
		</div>
	</aside>
</div>
