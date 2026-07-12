<script lang="ts">
	import { fade, slide } from 'svelte/transition';
	import { API_BASE_URL } from '$lib/config';

	// Props using Svelte 5 runes
	let {
		onTabChange,
		onSelectForReupload
	}: {
		onTabChange: (tab: string) => void;
		onSelectForReupload: (activity: {
			name: string;
			category: string;
			organizer: string;
			credits: number;
		}) => void;
	} = $props();

	interface BackendCertificate {
		id: number;
		student_roll_no: string;
		activity_name: string;
		activity_category: string;
		activity_date: string;
		organizer_name: string;
		event_level: string;
		cert_number: string;
		issue_date: string;
		participation_type: string;
		description: string;
		file_name: string;
		file_path: string;
		credits: number;
		status: string;
		rejection_reason: string;
	}

	interface MappedCertificate {
		id: number;
		certNumber: string;
		activity: string;
		category: string;
		date: string;
		credits: number;
		status: string;
		verifiedBy: string;
		file: string;
		feedback: string;
	}

	let token = localStorage.getItem('access_token') || '';
	let certificates = $state<MappedCertificate[]>([]);

	function formatDate(dateStr: string) {
		if (!dateStr) return '';
		const d = new Date(dateStr);
		return d.toLocaleDateString('en-GB', {
			day: '2-digit',
			month: 'short',
			year: 'numeric'
		});
	}

	async function loadCertificates() {
		try {
			const res = await fetch(`${API_BASE_URL}/api/student/certificates`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});

			if (res.ok) {
				const data = await res.json();
				certificates = data.map((cert: BackendCertificate) => ({
					id: cert.id,
					certNumber: cert.cert_number || '',
					activity: cert.activity_name,
					category: cert.activity_category,
					date: formatDate(cert.activity_date),
					credits: cert.credits,
					status: cert.status === 'Approved' ? 'Verified' : cert.status,
					verifiedBy: cert.status === 'Approved' ? 'Faculty Mentor' : '—',
					file: cert.file_name,
					feedback:
						cert.status === 'Rejected'
							? cert.rejection_reason
							: cert.status === 'Approved'
								? 'Credits approved in full.'
								: 'Awaiting faculty review.'
				}));
			}
		} catch (err) {
			console.error('Error fetching certificates:', err);
		}
	}

	$effect(() => {
		loadCertificates();
	});

	// Detail Modal state
	let activeCert = $state<MappedCertificate | null>(null);
	let isModalOpen = $state(false);

	// Toasts state
	interface Toast {
		id: number;
		message: string;
	}
	let toasts = $state<Toast[]>([]);
	let nextToastId = 0;

	function triggerToast(message: string) {
		const id = nextToastId++;
		toasts = [...toasts, { id, message }];
		setTimeout(() => {
			toasts = toasts.filter((t) => t.id !== id);
		}, 3000);
	}

	function openModal(cert: MappedCertificate) {
		activeCert = cert;
		isModalOpen = true;
	}

	function closeModal() {
		isModalOpen = false;
		activeCert = null;
	}

	async function handleDownload(certId: number | string, filename: string) {
		try {
			const res = await fetch(`${API_BASE_URL}/api/student/certificates/${certId}/file`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});

			if (!res.ok) {
				throw new Error('Failed to download certificate');
			}

			const blob = await res.blob();
			const objectUrl = URL.createObjectURL(blob);
			const link = document.createElement('a');
			link.href = objectUrl;
			link.download = filename || `certificate-${certId}`;
			document.body.appendChild(link);
			link.click();
			link.remove();
			URL.revokeObjectURL(objectUrl);

			triggerToast(`Downloading certificate ${certId} (${filename})...`);
		} catch {
			triggerToast('Could not download this certificate. Please try again.');
		}
	}

	function handleReuploadClick(cert: MappedCertificate) {
		// Populate draft and switch tab
		onSelectForReupload({
			name: cert.activity,
			category: cert.category,
			organizer: cert.verifiedBy !== '—' ? 'IIPS DAVV' : '',
			credits: cert.credits
		});
		triggerToast(`Redirecting to Upload Certificate for: ${cert.activity}`);
		setTimeout(() => {
			onTabChange('Upload Certificate');
		}, 800);
	}

	let totalCount = $derived(certificates.length);
	let verifiedCount = $derived(certificates.filter((c) => c.status === 'Verified').length);
	let pendingCount = $derived(certificates.filter((c) => c.status === 'Pending').length);
	let rejectedCount = $derived(certificates.filter((c) => c.status === 'Rejected').length);

	let verifiedPct = $derived(totalCount > 0 ? (verifiedCount / totalCount) * 100 : 0);
	let pendingPct = $derived(totalCount > 0 ? (pendingCount / totalCount) * 100 : 0);
	let rejectedPct = $derived(totalCount > 0 ? (rejectedCount / totalCount) * 100 : 0);
</script>

<!-- Toast Container -->
<div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-sm">
	{#each toasts as toast (toast.id)}
		<div
			transition:slide={{ duration: 150 }}
			class="p-4 rounded-xl border shadow-lg flex items-center gap-2 text-xs font-semibold font-sans bg-slate-800 border-slate-700 text-white"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="2"
				stroke="currentColor"
				class="w-4 h-4 text-accent-gold"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M12 9v3.75m0-10.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.75c0 5.592 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.57-.598-3.75h-.152c-3.196 0-6.1-1.249-8.25-3.286Zm0 13.036h.008v.008H12v-.008Z"
				/>
			</svg>
			<span>{toast.message}</span>
		</div>
	{/each}
</div>

<!-- ==================== SUMMARY CARDS ==================== -->
<section
	class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5"
	aria-label="Certificate overview statistics"
>
	<!-- Card 1: Total Uploaded -->
	<div
		class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between shadow-xs hover:shadow-md transition-shadow duration-200"
	>
		<div class="flex items-center justify-between">
			<span class="text-2xl font-bold font-serif text-slate-900">{totalCount}</span>
			<div class="p-2.5 rounded-lg bg-blue-50 text-blue-600 border border-blue-100">
				<!-- File icon -->
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
			<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase">Total Uploaded</h3>
		</div>
	</div>

	<!-- Card 2: Verified Certificates -->
	<div
		class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between shadow-xs hover:shadow-md transition-shadow duration-200"
	>
		<div class="flex items-center justify-between">
			<span class="text-2xl font-bold font-serif text-slate-900">{verifiedCount}</span>
			<div class="p-2.5 rounded-lg bg-emerald-50 text-emerald-600 border border-emerald-100">
				<!-- Check circle icon -->
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
						d="M9 12.75 11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 0 1-1.043 3.296 3.745 3.745 0 0 1-3.296 1.043A3.745 3.745 0 0 1 12 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 0 1-3.296-1.043 3.745 3.745 0 0 1-1.043-3.296A3.745 3.745 0 0 1 3 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 0 1 1.043-3.296 3.746 3.746 0 0 1 3.296-1.043A3.746 3.746 0 0 1 12 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 0 1 3.296 1.043 3.746 3.746 0 0 1 1.043 3.296A3.745 3.745 0 0 1 21 12Z"
					/>
				</svg>
			</div>
		</div>
		<div class="mt-4">
			<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase">
				Verified Certificates
			</h3>
		</div>
	</div>

	<!-- Card 3: Pending Verification -->
	<div
		class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between shadow-xs hover:shadow-md transition-shadow duration-200"
	>
		<div class="flex items-center justify-between">
			<span class="text-2xl font-bold font-serif text-slate-900">{pendingCount}</span>
			<div class="p-2.5 rounded-lg bg-amber-50 text-amber-600 border border-amber-100">
				<!-- Clock Icon -->
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
						d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
					/>
				</svg>
			</div>
		</div>
		<div class="mt-4">
			<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase">
				Pending Verification
			</h3>
		</div>
	</div>

	<!-- Card 4: Rejected Certificates -->
	<div
		class="bg-white p-5 rounded-xl border border-slate-200 flex flex-col justify-between shadow-xs hover:shadow-md transition-shadow duration-200"
	>
		<div class="flex items-center justify-between">
			<span class="text-2xl font-bold font-serif text-slate-900">{rejectedCount}</span>
			<div class="p-2.5 rounded-lg bg-rose-50 text-rose-600 border border-rose-100">
				<!-- X circle icon -->
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
						d="m9.75 9.75 4.5 4.5m0-4.5-4.5 4.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
					/>
				</svg>
			</div>
		</div>
		<div class="mt-4">
			<h3 class="text-xs font-bold text-slate-405 tracking-wider uppercase">
				Rejected Certificates
			</h3>
		</div>
	</div>
</section>

<!-- ==================== STATUS OVERVIEW SECTION ==================== -->
<section class="bg-white border border-slate-200 rounded-xl p-5 shadow-xs space-y-4">
	<div>
		<h2 class="text-xs font-bold font-serif text-inst-navy uppercase tracking-wider">
			Certificate Status Overview
		</h2>
	</div>

	<div class="space-y-4">
		<!-- Horizontal Proportion Bar -->
		<div class="h-3 w-full bg-slate-100 rounded-full overflow-hidden flex">
			<!-- Verified -->
			<div
				class="h-full bg-emerald-600 rounded-l-full"
				style="width: {verifiedPct}%"
				title="Verified: {verifiedPct.toFixed(1)}%"
			></div>
			<!-- Pending -->
			<div
				class="h-full bg-amber-500"
				style="width: {pendingPct}%"
				title="Pending: {pendingPct.toFixed(1)}%"
			></div>
			<!-- Rejected -->
			<div
				class="h-full bg-rose-600 rounded-r-full"
				style="width: {rejectedPct}%"
				title="Rejected: {rejectedPct.toFixed(1)}%"
			></div>
		</div>

		<!-- Status Legend -->
		<div
			class="flex flex-wrap items-center gap-6 text-xs font-semibold text-slate-600 font-sans pl-1"
		>
			<span class="flex items-center gap-2">
				<span class="w-2.5 h-2.5 bg-emerald-600 rounded-full"></span>
				Verified
			</span>
			<span class="flex items-center gap-2">
				<span class="w-2.5 h-2.5 bg-amber-500 rounded-full"></span>
				Pending
			</span>
			<span class="flex items-center gap-2">
				<span class="w-2.5 h-2.5 bg-rose-600 rounded-full"></span>
				Rejected
			</span>
		</div>
	</div>
</section>

<!-- ==================== CERTIFICATE REGISTRY SECTION ==================== -->
<section class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
	<div class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/20">
		<div>
			<h2 class="text-xs font-bold font-serif text-inst-navy uppercase tracking-wider">
				Certificate Registry
			</h2>
		</div>
		<span
			class="bg-slate-100 border border-slate-200 text-slate-600 text-[10px] font-extrabold uppercase px-2.5 py-1 rounded"
		>
			8 records
		</span>
	</div>

	<div class="overflow-x-auto">
		<table class="w-full text-left border-collapse">
			<thead>
				<tr
					class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
				>
					<th class="py-3.5 px-5">Certificate ID</th>
					<th class="py-3.5 px-5">Activity Name</th>
					<th class="py-3.5 px-5">Category</th>
					<th class="py-3.5 px-5">Submission Date</th>
					<th class="py-3.5 px-5">Credits</th>
					<th class="py-3.5 px-5">Status</th>
					<th class="py-3.5 px-5">Verified By</th>
					<th class="py-3.5 px-5 text-center">Actions</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-slate-100 text-xs font-sans">
				{#each certificates as cert}
					<tr class="hover:bg-slate-50/50 transition-colors">
						<td class="py-4 px-5 font-bold text-slate-800 select-all">{cert.certNumber || '—'}</td>
						<td class="py-4 px-5 font-bold text-slate-850">{cert.activity}</td>
						<td class="py-4 px-5 text-slate-500 font-semibold">{cert.category}</td>
						<td class="py-4 px-5 text-slate-500 font-medium">{cert.date}</td>
						<td class="py-4 px-5 font-extrabold text-slate-900 text-sm">{cert.credits}</td>
						<td class="py-4 px-5">
							<span
								class="inline-flex items-center gap-1 px-2.5 py-1 rounded-md text-[10px] font-extrabold uppercase tracking-wide
								{cert.status === 'Verified'
									? 'bg-emerald-50 text-emerald-700 border border-emerald-105'
									: cert.status === 'Pending'
										? 'bg-amber-50 text-amber-705 border border-amber-105'
										: 'bg-rose-50 text-rose-700 border border-rose-105'}"
							>
								{#if cert.status === 'Verified'}
									<span class="w-1.5 h-1.5 bg-emerald-600 rounded-full"></span>
								{:else if cert.status === 'Pending'}
									<span class="w-1.5 h-1.5 bg-amber-500 rounded-full"></span>
								{:else}
									<span class="w-1.5 h-1.5 bg-rose-600 rounded-full"></span>
								{/if}
								{cert.status}
							</span>
						</td>
						<td class="py-4 px-5 text-slate-600 font-semibold">{cert.verifiedBy}</td>
						<td class="py-4 px-5">
							<div class="flex items-center justify-center gap-2">
								<!-- View Button -->
								<button
									onclick={() => openModal(cert)}
									class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-[10px] uppercase rounded-lg transition-colors focus:outline-none"
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
									View
								</button>

								<!-- Action depends on status -->
								{#if cert.status === 'Verified'}
									<!-- Download Button -->
									<button
										onclick={() => handleDownload(cert.id, cert.file)}
										class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-emerald-50 text-emerald-700 hover:bg-emerald-100 font-extrabold text-[10px] uppercase rounded-lg transition-colors focus:outline-none"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2.5"
											stroke="currentColor"
											class="w-3.5 h-3.5"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
											/>
										</svg>
										Download
									</button>
								{:else if cert.status === 'Rejected'}
									<!-- Re-upload Button -->
									<button
										onclick={() => handleReuploadClick(cert)}
										class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-accent-red text-white hover:bg-accent-red/90 font-extrabold text-[10px] uppercase rounded-lg transition-colors focus:outline-none"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											fill="none"
											viewBox="0 0 24 24"
											stroke-width="2.5"
											stroke="currentColor"
											class="w-3.5 h-3.5"
										>
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"
											/>
										</svg>
										Re-upload
									</button>
								{/if}
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>

<!-- ==================== DIALOG DETAILED VIEW MODAL ==================== -->
{#if isModalOpen && activeCert}
	<div
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
	>
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div
			onclick={(e) => e.stopPropagation()}
			class="w-full max-w-lg bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
		>
			<!-- Modal Header -->
			<div class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30">
				<div>
					<h3 class="text-sm font-bold font-serif text-inst-navy">Certificate Details</h3>
					{#if activeCert.certNumber}
						<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
							ID: {activeCert.certNumber}
						</p>
					{/if}
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
			<div class="p-6 space-y-5">
				<!-- Detail Info Grid -->
				<div class="grid grid-cols-2 gap-4 bg-slate-50 p-4 rounded-xl border border-slate-150">
					<div>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Activity Name</span
						>
						<span class="text-xs font-bold text-slate-800 block mt-0.5">{activeCert.activity}</span>
					</div>
					<div>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Category</span
						>
						<span class="text-xs font-bold text-slate-800 block mt-0.5">{activeCert.category}</span>
					</div>
					<div>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Submission Date</span
						>
						<span class="text-xs font-bold text-slate-800 block mt-0.5">{activeCert.date}</span>
					</div>
					<div>
						<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
							>Credits Allocated</span
						>
						<span class="text-xs font-bold text-slate-800 block mt-0.5"
							>+{activeCert.credits} Credits</span
						>
					</div>
				</div>

				<!-- Document Info -->
				<div class="space-y-1">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
						>Uploaded File</span
					>
					<div
						class="flex items-center gap-3 p-3.5 border border-slate-200 rounded-xl bg-slate-50/30"
					>
						<div class="p-2 bg-rose-50 text-rose-600 border border-rose-100 rounded-lg shrink-0">
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
						</div>
						<div class="flex-grow min-w-0">
							<div class="text-xs font-bold text-slate-800 truncate">{activeCert.file}</div>
							<div class="text-[10px] text-slate-400 font-semibold">Verification scan</div>
						</div>
					</div>
				</div>

				<!-- Status Verification Status -->
				<div class="space-y-1.5">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block"
						>Verification Review</span
					>
					<div
						class="p-4 rounded-xl border flex flex-col gap-1.5
						{activeCert.status === 'Verified'
							? 'bg-emerald-50/30 border-emerald-100 text-emerald-900'
							: activeCert.status === 'Pending'
								? 'bg-amber-50/30 border-amber-100 text-amber-900'
								: 'bg-rose-50/30 border-rose-100 text-rose-900'}"
					>
						<div class="flex items-center gap-2 text-xs font-bold uppercase tracking-wider">
							<span
								class="w-2 h-2 rounded-full
								{activeCert.status === 'Verified'
									? 'bg-emerald-600'
									: activeCert.status === 'Pending'
										? 'bg-amber-500'
										: 'bg-rose-600'}"
							></span>
							Status: {activeCert.status}
						</div>
						<p class="text-xs leading-relaxed italic text-slate-700 pt-1 font-medium">
							"{activeCert.feedback}"
						</p>
						<div class="text-[10px] text-slate-400 font-bold uppercase mt-1">
							Audited By: <span class="text-slate-700">{activeCert.verifiedBy}</span>
						</div>
					</div>
				</div>
			</div>

			<!-- Modal Footer -->
			<div
				class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
			>
				<button
					onclick={closeModal}
					class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs tracking-wider uppercase rounded-lg transition duration-200 focus:outline-none"
				>
					Close Details
				</button>
				{#if activeCert.status === 'Verified'}
					<button
						onclick={() => {
							handleDownload(activeCert!.id, activeCert!.file);
							closeModal();
						}}
						class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs tracking-wider uppercase rounded-lg transition duration-200 focus:outline-none"
					>
						Download File
					</button>
				{:else if activeCert.status === 'Rejected'}
					<button
						onclick={() => handleReuploadClick(activeCert!)}
						class="px-4 py-2 bg-accent-red hover:bg-accent-red/90 text-white font-bold text-xs tracking-wider uppercase rounded-lg transition duration-200 focus:outline-none"
					>
						Re-upload
					</button>
				{/if}
			</div>
		</div>
	</div>
{/if}
