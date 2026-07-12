<script lang="ts">
	import { fade, slide } from 'svelte/transition';
	import ActivityManagementView from './ActivityManagementView.svelte';
	import TrackManagementView from './TrackManagementView.svelte';

	// Sidebar menu items list for Super Admin Portal
	const menuItems = [
		{
			name: 'Dashboard',
			icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6'
		},
		{
			name: 'User Management',
			icon: 'M15 19.128a9.38 9.38 0 0 0 2.625.372 9.337 9.337 0 0 0 4.121-.952 4.125 4.125 0 0 0-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.109A11.386 11.386 0 0 1 10.089 21c-2.316 0-4.445-.69-6.22-1.879v-.003a4.125 4.125 0 0 1 7.533-2.493M15 19.128v-.003c0-1.112-.285-2.16-.786-3.07M14.214 16.058A9.396 9.396 0 0 0 10.089 15c-1.47 0-2.854.34-4.082.945M14.214 16.058a9.386 9.386 0 0 1 0 3.07M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5'
		},
		{
			name: 'Activity Management',
			icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z'
		},
		{
			name: 'Track Management',
			icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
		},
		{
			name: 'Announcement Management',
			icon: 'M19.114 5.636a9 9 0 010 12.728M16.463 8.288a5.25 5.25 0 010 7.424M6.75 8.25l4.72-4.72a.75.75 0 011.28.53v15.88a.75.75 0 01-1.28.53l-4.72-4.72H4.51c-.88 0-1.704-.507-1.938-1.354A9.01 9.01 0 012.25 12c0-.83.112-1.633.322-2.396C2.806 8.756 3.63 8.25 4.51 8.25H6.75z'
		},
		{
			name: 'Reports Center',
			icon: 'M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z'
		},
		{
			name: 'System Settings',
			icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z'
		}
	];

	// Selected menu state
	let currentTab = $state('Dashboard');
	let isMobileSidebarOpen = $state(false);

	// Search & Notification states
	let isNotificationsOpen = $state(false);
	let searchQuery = $state('');

	// Mock Notifications for Super Admin
	const notifications = [
		{
			id: 1,
			text: 'System health warning: Background DB sync delayed by 3m.',
			time: '15 mins ago',
			unread: true
		},
		{
			id: 2,
			text: 'Security Alert: Failed SSH attempt from 192.168.1.104.',
			time: '1 hour ago',
			unread: true
		},
		{
			id: 3,
			text: 'Daily audit backup successfully uploaded to S3.',
			time: '6 hours ago',
			unread: false
		}
	];

	// Mock Recent Users data (Step 3) - using Svelte 5 reactive states
	let recentUsers = $state([
		{
			name: 'Rahul Sharma',
			id: 'ENR024001',
			role: 'Student',
			dept: 'Computer Science',
			status: 'Active'
		},
		{
			name: 'Dr. Priya Patel',
			id: 'ADM024010',
			role: 'Admin',
			dept: 'Electronics & Comm.',
			status: 'Active'
		},
		{
			name: 'Arjun Desai',
			id: 'ENR024008',
			role: 'Student',
			dept: 'Data Science',
			status: 'Pending'
		},
		{
			name: 'Sneha Kumar',
			id: 'ENR024015',
			role: 'Student',
			dept: 'Mechanical Engg.',
			status: 'Active'
		},
		{
			name: 'Dr. Vikram Singh',
			id: 'ADM024003',
			role: 'Admin',
			dept: 'Civil Engineering',
			status: 'Active'
		}
	]);

	// Mock User Registry for User Management Portal (Step 3 & 4)
	let userRegistry = $state([
		{
			name: 'Rahul Sharma',
			id: 'ENG024001',
			role: 'Student',
			displayRole: 'Student',
			dept: 'Mtech CS',
			status: 'Active',
			avatarBg: 'bg-[#802D2D] text-white'
		},
		{
			name: 'Dr. Priya Patel',
			id: 'MNT024001',
			role: 'Admin',
			displayRole: 'Mentor',
			dept: 'Mtech IT',
			status: 'Active',
			avatarBg: 'bg-[#C07050] text-white'
		},
		{
			name: 'Arjun Desai',
			id: 'ENG024032',
			role: 'Student',
			displayRole: 'Student',
			dept: 'MBA',
			status: 'Pending',
			avatarBg: 'bg-[#2A4D69] text-white'
		},
		{
			name: 'Sneha Kumar',
			id: 'ENG024075',
			role: 'Student',
			displayRole: 'Student',
			dept: 'MCA',
			status: 'Active',
			avatarBg: 'bg-[#7D4427] text-white'
		},
		{
			name: 'Dr. Vikram Singh',
			id: 'MNT024008',
			role: 'Admin',
			displayRole: 'Mentor',
			dept: 'Mtech CS',
			status: 'Active',
			avatarBg: 'bg-[#3A6B88] text-white'
		},
		{
			name: 'Kavya Krishnan',
			id: 'ENG024099',
			role: 'Student',
			displayRole: 'Student',
			dept: 'Mtech IT',
			status: 'Inactive',
			avatarBg: 'bg-[#7A6B3A] text-white'
		},
		{
			name: 'Admin Coordinator',
			id: 'ADM024001',
			role: 'Admin',
			displayRole: 'Admin',
			dept: 'MBA',
			status: 'Active',
			avatarBg: 'bg-[#244E3F] text-white'
		},
		{
			name: 'Dev Mehta',
			id: 'ENG024112',
			role: 'Student',
			displayRole: 'Student',
			dept: 'Mtech CS',
			status: 'Active',
			avatarBg: 'bg-[#B04A3A] text-white'
		},
		{
			name: 'Dr. Anita Rao',
			id: 'MNT024012',
			role: 'Admin',
			displayRole: 'Mentor',
			dept: 'MCA',
			status: 'Pending',
			avatarBg: 'bg-[#4E6B4E] text-white'
		},
		{
			name: 'Rohan Nair',
			id: 'ENG024055',
			role: 'Student',
			displayRole: 'Student',
			dept: 'MBA',
			status: 'Active',
			avatarBg: 'bg-[#3C4A8E] text-white'
		},
		{
			name: 'Meera Joshi',
			id: 'ENG024143',
			role: 'Student',
			displayRole: 'Student',
			dept: 'Mtech IT',
			status: 'Active',
			avatarBg: 'bg-[#A63A50] text-white'
		},
		{
			name: 'Dr. Suresh Iyer',
			id: 'MNT024019',
			role: 'Admin',
			displayRole: 'Mentor',
			dept: 'MCA',
			status: 'Active',
			avatarBg: 'bg-[#6B5E33] text-white'
		}
	]);

	let userFilter = $state('All');
	let userSearch = $state('');

	let filteredUsers = $derived(
		userRegistry.filter((user) => {
			const matchesFilter =
				userFilter === 'All' ||
				(userFilter === 'Student' && user.role === 'Student') ||
				(userFilter === 'Mentor' && user.displayRole === 'Mentor') ||
				(userFilter === 'Admin' && user.displayRole === 'Admin');

			const matchesSearch =
				user.name.toLowerCase().includes(userSearch.toLowerCase()) ||
				user.id.toLowerCase().includes(userSearch.toLowerCase()) ||
				user.dept.toLowerCase().includes(userSearch.toLowerCase());

			return matchesFilter && matchesSearch;
		})
	);

	let totalStudentsCount = $derived(userRegistry.filter((u) => u.role === 'Student').length + 1236);
	let totalMentorsCount = $derived(
		userRegistry.filter((u) => u.displayRole === 'Mentor').length + 83
	);
	let totalAdminsCount = $derived(userRegistry.filter((u) => u.displayRole === 'Admin').length + 5);
	let activeUsersCount = $derived(userRegistry.filter((u) => u.status === 'Active').length + 1179);
	let pendingUsersCount = $derived(userRegistry.filter((u) => u.status === 'Pending').length + 46);
	let inactiveUsersCount = $derived(
		userRegistry.filter((u) => u.status === 'Inactive').length + 105
	);

	// Mock Recent System Activities (Step 5)
	let recentLogs = $state([
		{
			activity: 'New student account created',
			type: 'User Added',
			performedBy: 'Super Admin',
			date: 'Jun 27, 2026',
			status: 'Completed'
		},
		{
			activity: 'Activity "Blood Donation Camp" published',
			type: 'Activity Created',
			performedBy: 'Super Admin',
			date: 'Jun 25, 2026',
			status: 'Completed'
		},
		{
			activity: 'Track "Interdisciplinary" updated',
			type: 'Track Updated',
			performedBy: 'Super Admin',
			date: 'Jun 24, 2026',
			status: 'Completed'
		},
		{
			activity: 'Platform-wide announcement published',
			type: 'Announcement',
			performedBy: 'Super Admin',
			date: 'Jun 22, 2026',
			status: 'Completed'
		},
		{
			activity: 'Monthly activity report generated',
			type: 'Report',
			performedBy: 'Super Admin',
			date: 'Jun 20, 2026',
			status: 'Completed'
		},
		{
			activity: 'Admin assigned to Data Science track',
			type: 'Track Updated',
			performedBy: 'Super Admin',
			date: 'Jun 18, 2026',
			status: 'Completed'
		}
	]);

	// Action modals states (Step 6)
	let isCreateUserModalOpen = $state(false);
	let isCreateActivityModalOpen = $state(false);
	let isCreateTrackModalOpen = $state(false);
	let isGenerateReportModalOpen = $state(false);

	// Action form states
	let newUserName = $state('');
	let newUserRole = $state('Student');
	let newUserDept = $state('');
	let newActivityName = $state('');
	let newActivityCredits = $state(5);
	let newTrackName = $state('');

	// ── System Settings (Step 7) ────────────────────────────────────────────────
	type SettingStatus = 'Active' | 'Enabled' | 'Disabled';
	interface PlatformSetting {
		name: string;
		description: string;
		value: string;
		status: SettingStatus;
	}

	const settingsTabs = [
		'Academic Year',
		'Credit Policy',
		'Activity Rules',
		'Notifications',
		'Platform',
		'Security'
	];

	let activeSettingsTab = $state('Academic Year');

	let settingsSections = $state<Record<string, PlatformSetting[]>>({
		'Academic Year': [
			{
				name: 'Current Academic Year',
				description: 'Active academic cycle label displayed platform-wide',
				value: '2025-2026',
				status: 'Active'
			},
			{
				name: 'Academic Year Start Date',
				description: 'Official start date of the current academic year',
				value: 'Aug 1, 2025',
				status: 'Active'
			},
			{
				name: 'Academic Year End Date',
				description: 'Official end date of the current academic year',
				value: 'May 31, 2026',
				status: 'Active'
			},
			{
				name: 'Enrollment Deadline',
				description: 'Last date for activity enrollment submissions',
				value: 'Sep 15, 2025',
				status: 'Active'
			}
		],
		'Credit Policy': [
			{
				name: 'Minimum Credits Required',
				description: 'Total credits a student must earn to graduate',
				value: '100',
				status: 'Active'
			},
			{
				name: 'Maximum Credits per Activity',
				description: 'Upper cap on credits from a single activity',
				value: '20',
				status: 'Active'
			},
			{
				name: 'Minimum Credits per Semester',
				description: 'Credits a student must earn each semester',
				value: '12',
				status: 'Active'
			},
			{
				name: 'Credit Rollover',
				description: 'Carry unused credits to the next academic year',
				value: 'Enabled',
				status: 'Enabled'
			},
			{
				name: 'Grace Credit Buffer',
				description: 'Extra credits allowed beyond the target',
				value: '5',
				status: 'Active'
			}
		],
		'Activity Rules': [
			{
				name: 'Max Active Enrollments',
				description: 'Activities a student can be enrolled in at once',
				value: '5',
				status: 'Active'
			},
			{
				name: 'Mentor Approval Required',
				description: 'Require mentor sign-off before credit is granted',
				value: 'Enabled',
				status: 'Enabled'
			},
			{
				name: 'Auto-Verification Threshold',
				description: 'Score above which certificates auto-verify',
				value: '90%',
				status: 'Active'
			},
			{
				name: 'Self-Reported Activities',
				description: 'Allow students to log their own activities',
				value: 'Disabled',
				status: 'Disabled'
			},
			{
				name: 'Resubmission Window',
				description: 'Days allowed to resubmit a rejected certificate',
				value: '7 days',
				status: 'Active'
			}
		],
		Notifications: [
			{
				name: 'Email Notifications',
				description: 'Send system emails for key events',
				value: 'Enabled',
				status: 'Enabled'
			},
			{
				name: 'OTP Expiry Duration',
				description: 'Validity window for verification codes',
				value: '15 min',
				status: 'Active'
			},
			{
				name: 'Reminder Frequency',
				description: 'How often pending-task reminders are sent',
				value: 'Weekly',
				status: 'Active'
			},
			{
				name: 'Announcement Broadcasts',
				description: 'Push platform-wide announcements to all users',
				value: 'Enabled',
				status: 'Enabled'
			}
		],
		Platform: [
			{
				name: 'Platform Name',
				description: 'Display name shown across the portal',
				value: 'iSPARC',
				status: 'Active'
			},
			{
				name: 'Maintenance Mode',
				description: 'Temporarily disable access for non-admins',
				value: 'Disabled',
				status: 'Disabled'
			},
			{
				name: 'Default Time Zone',
				description: 'Base time zone for schedules and logs',
				value: 'IST (UTC+5:30)',
				status: 'Active'
			},
			{
				name: 'Default Language',
				description: 'Default interface language for new users',
				value: 'English',
				status: 'Active'
			}
		],
		Security: [
			{
				name: 'Minimum Password Length',
				description: 'Required characters for user passwords',
				value: '8',
				status: 'Active'
			},
			{
				name: 'Session Timeout',
				description: 'Idle minutes before automatic logout',
				value: '30 min',
				status: 'Active'
			},
			{
				name: 'Two-Factor Authentication',
				description: 'Require 2FA for admin accounts',
				value: 'Enabled',
				status: 'Enabled'
			},
			{
				name: 'Max Login Attempts',
				description: 'Failed logins before temporary lockout',
				value: '5',
				status: 'Active'
			}
		]
	});

	let currentSettings = $derived(settingsSections[activeSettingsTab] ?? []);

	// Edit setting modal state
	let isEditSettingOpen = $state(false);
	let editSettingIndex = $state(-1);
	let editSettingName = $state('');
	let editSettingValue = $state('');

	function openEditSetting(index: number) {
		const s = settingsSections[activeSettingsTab][index];
		editSettingIndex = index;
		editSettingName = s.name;
		editSettingValue = s.value;
		isEditSettingOpen = true;
	}

	function handleSaveSetting(e: Event) {
		e.preventDefault();
		if (editSettingIndex < 0 || !editSettingValue.trim()) return;
		settingsSections[activeSettingsTab][editSettingIndex].value = editSettingValue.trim();
		triggerToast(`"${editSettingName}" updated successfully!`);
		isEditSettingOpen = false;
		editSettingIndex = -1;
	}

	function handleSaveAllSettings() {
		triggerToast('All system settings saved successfully!');
	}

	function handleResetSettings() {
		triggerToast(`"${activeSettingsTab}" settings reset to default values.`);
	}

	function settingStatusClass(status: SettingStatus): string {
		return status === 'Disabled'
			? 'bg-slate-100 text-slate-500 border-slate-200'
			: 'bg-emerald-50 text-emerald-700 border-emerald-100';
	}

	function settingStatusDot(status: SettingStatus): string {
		return status === 'Disabled' ? 'bg-slate-400' : 'bg-emerald-500';
	}

	// Toast states
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

	function handleCreateUser(e: Event) {
		e.preventDefault();
		if (!newUserName.trim()) return;
		const generatedId =
			newUserRole === 'Student'
				? `ENR024${Math.floor(100 + Math.random() * 900)}`
				: `ADM024${Math.floor(100 + Math.random() * 900)}`;
		recentUsers = [
			{
				name: newUserName,
				id: generatedId,
				role: newUserRole,
				dept: newUserDept || 'Computer Science',
				status: 'Active'
			},
			...recentUsers
		];

		const avatarColors = [
			'bg-[#802D2D] text-white',
			'bg-[#C07050] text-white',
			'bg-[#2A4D69] text-white',
			'bg-[#7D4427] text-white',
			'bg-[#3A6B88] text-white',
			'bg-[#7A6B3A] text-white',
			'bg-[#244E3F] text-white',
			'bg-[#B04A3A] text-white',
			'bg-[#4E6B4E] text-white',
			'bg-[#3C4A8E] text-white',
			'bg-[#A63A50] text-white',
			'bg-[#6B5E33] text-white'
		];
		const randomColor = avatarColors[Math.floor(Math.random() * avatarColors.length)];

		userRegistry = [
			{
				name: newUserName,
				id: generatedId,
				role: newUserRole,
				displayRole: newUserRole === 'Admin' ? 'Admin' : newUserRole,
				dept: newUserDept || 'Computer Science',
				status: 'Active',
				avatarBg: randomColor
			},
			...userRegistry
		];

		recentLogs = [
			{
				activity: `New ${newUserRole.toLowerCase()} account created for ${newUserName}`,
				type: 'User Added',
				performedBy: 'Super Admin',
				date: 'Jun 28, 2026',
				status: 'Completed'
			},
			...recentLogs
		];
		triggerToast(`User "${newUserName}" created successfully!`);
		newUserName = '';
		newUserDept = '';
		isCreateUserModalOpen = false;
	}

	function handleCreateActivity(e: Event) {
		e.preventDefault();
		if (!newActivityName.trim()) return;
		recentLogs = [
			{
				activity: `Activity "${newActivityName}" published`,
				type: 'Activity Created',
				performedBy: 'Super Admin',
				date: 'Jun 28, 2026',
				status: 'Completed'
			},
			...recentLogs
		];
		triggerToast(`Activity "${newActivityName}" published successfully!`);
		newActivityName = '';
		isCreateActivityModalOpen = false;
	}

	function handleCreateTrack(e: Event) {
		e.preventDefault();
		if (!newTrackName.trim()) return;
		recentLogs = [
			{
				activity: `Track "${newTrackName}" updated`,
				type: 'Track Updated',
				performedBy: 'Super Admin',
				date: 'Jun 28, 2026',
				status: 'Completed'
			},
			...recentLogs
		];
		triggerToast(`Track "${newTrackName}" created successfully!`);
		newTrackName = '';
		isCreateTrackModalOpen = false;
	}

	function handleGenerateReport() {
		triggerToast('System activity report generation initiated. Downloading...');
		isGenerateReportModalOpen = false;
	}

	// ── Reports Center (Step 7) ─────────────────────────────────────────────────
	type ReportStatus = 'Ready' | 'Processing' | 'Failed';
	type ReportFormat = 'PDF' | 'Excel' | 'CSV';

	interface GeneratedReport {
		name: string;
		by: string;
		date: string;
		format: ReportFormat;
		status: ReportStatus;
	}

	// Generate form state
	const reportTypes = [
		'Student Performance',
		'Mentor Analytics',
		'Activity Participation',
		'Credit Distribution',
		'Department Analytics',
		'Certificate Verification',
		'Leaderboard Rankings',
		'Semester Summary'
	];
	const reportCourses = ['All Courses', 'Mtech CS', 'Mtech IT', 'MBA', 'MCA'];
	const reportSemesters = ['All Semesters', 'Semester 1', 'Semester 2', 'Semester 3', 'Semester 4'];
	const reportFormats: ReportFormat[] = ['PDF', 'Excel', 'CSV'];

	let reportType = $state('Student Performance');
	let reportCourse = $state('All Courses');
	let reportSemester = $state('All Semesters');
	let reportFrom = $state('2025-08-01');
	let reportTo = $state('2026-05-31');
	let reportFormat = $state<ReportFormat>('PDF');

	// Recent reports
	let recentReports = $state<GeneratedReport[]>([
		{
			name: 'Student Performance Report',
			by: 'Super Admin',
			date: 'Jun 28, 2026',
			format: 'PDF',
			status: 'Ready'
		},
		{
			name: 'Mentor Analytics Report',
			by: 'Super Admin',
			date: 'Jun 27, 2026',
			format: 'Excel',
			status: 'Ready'
		},
		{
			name: 'Activity Participation Report',
			by: 'Super Admin',
			date: 'Jun 26, 2026',
			format: 'CSV',
			status: 'Processing'
		},
		{
			name: 'Credit Distribution Report',
			by: 'Super Admin',
			date: 'Jun 24, 2026',
			format: 'PDF',
			status: 'Ready'
		},
		{
			name: 'Department-wide Report',
			by: 'Super Admin',
			date: 'Jun 22, 2026',
			format: 'Excel',
			status: 'Ready'
		},
		{
			name: 'Semester Summary – Even',
			by: 'Super Admin',
			date: 'Jun 20, 2026',
			format: 'PDF',
			status: 'Ready'
		},
		{
			name: 'Leaderboard Rankings Q2',
			by: 'Super Admin',
			date: 'Jun 18, 2026',
			format: 'Excel',
			status: 'Ready'
		},
		{
			name: 'Certificate Verification Log',
			by: 'Super Admin',
			date: 'Jun 15, 2026',
			format: 'PDF',
			status: 'Failed'
		}
	]);

	// Quick templates
	const quickTemplates = [
		{ name: 'Student Credit Summary', desc: 'Credits per student' },
		{ name: 'Certificate Verification', desc: 'Approval status log' },
		{ name: 'Activity Participation', desc: 'Enrollment vs completion' },
		{ name: 'Batch Performance', desc: 'Batch-wise comparison' },
		{ name: 'Mentor Activity', desc: 'Mentor engagement' },
		{ name: 'Department Performance', desc: 'Dept-wise metrics' },
		{ name: 'Semester Summary', desc: 'Per-semester rollup' },
		{ name: 'System Usage Report', desc: 'Platform activity' }
	];

	// Scheduled reports
	const scheduledReports = [
		{ name: 'Monthly Student Report', freq: 'Every 1st of the month', format: 'PDF' },
		{ name: 'Semester Summary', freq: 'End of each semester', format: 'Excel' },
		{ name: 'Mentor Performance Report', freq: 'Weekly · Monday', format: 'PDF' },
		{ name: 'Department Analytics', freq: 'Bi-weekly', format: 'Excel' }
	];

	// Export center
	const exportOptions = [
		{ name: 'Export Student Database', desc: 'Full student records', count: '1,248 records' },
		{ name: 'Export Activity Records', desc: 'All logged activities', count: '3,412 records' },
		{ name: 'Export Credit Reports', desc: 'Credit distribution data', count: '892 records' },
		{ name: 'Export Certificate Data', desc: 'Verified certificates', count: '2,076 records' }
	];

	// Report activity audit log
	const reportAuditLog = [
		{
			action: 'Student Performance Report downloaded',
			type: 'PDF',
			user: 'Dr. Mehta',
			time: 'Today, 09:45 AM'
		},
		{
			action: 'Semester Report exported to Excel',
			type: 'Excel',
			user: 'Super Admin',
			time: 'Today, 08:12 AM'
		},
		{
			action: 'Leaderboard Report generated',
			type: 'CSV',
			user: 'Admin Sharma',
			time: 'Yesterday, 05:30 PM'
		},
		{
			action: 'Certificate Report shared with NAAC',
			type: 'PDF',
			user: 'Super Admin',
			time: 'Yesterday, 02:15 PM'
		},
		{
			action: 'Department-wise Report scheduled',
			type: 'Auto',
			user: 'Admin Verma',
			time: '25 Jun 2026, 11:30 AM'
		},
		{
			action: 'Extracurricular Marksheet downloaded',
			type: 'Excel',
			user: 'Dr. Mehta',
			time: '24 Jun 2026, 03:20 PM'
		}
	];

	// Institutional overview snapshot
	const institutionalStats = [
		{ value: '3,842', label: 'Registered Students', color: 'text-indigo-600' },
		{ value: '3,614', label: 'Active Students', color: 'text-emerald-600' },
		{ value: '284', label: 'Activities Conducted', color: 'text-purple-600' },
		{ value: '18.4', label: 'Avg Credits Earned', color: 'text-amber-600' }
	];

	// Avg score by department (bar chart)
	const deptScores = [
		{ dept: 'MCA', score: 88 },
		{ dept: 'MBA', score: 76 },
		{ dept: 'BCA', score: 71 },
		{ dept: 'B.Com', score: 63 },
		{ dept: 'BBA', score: 82 },
		{ dept: 'MSc CS', score: 85 }
	];

	// Monthly participation trend (dual line)
	const participationMonths = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'];
	const activitiesTrend = [40, 46, 50, 45, 62, 70];
	const certificatesTrend = [36, 41, 47, 43, 54, 66];
	const trendMax = 80;
	const trendMin = 20;
	const trendLine = (data: number[]) =>
		data
			.map((v, i) => {
				const x = (i / (data.length - 1)) * 100;
				const y = 40 - ((v - trendMin) / (trendMax - trendMin)) * 36 - 2;
				return `${x},${y}`;
			})
			.join(' ');
	const activitiesPoints = trendLine(activitiesTrend);
	const certificatesPoints = trendLine(certificatesTrend);

	// Handlers
	function runReportGeneration() {
		recentReports = [
			{
				name: `${reportType} Report`,
				by: 'Super Admin',
				date: 'Jun 28, 2026',
				format: reportFormat,
				status: 'Processing'
			},
			...recentReports
		];
		triggerToast(`Generating "${reportType} Report" (${reportFormat})...`);
	}

	function resetReportForm() {
		reportType = 'Student Performance';
		reportCourse = 'All Courses';
		reportSemester = 'All Semesters';
		reportFrom = '2025-08-01';
		reportTo = '2026-05-31';
		reportFormat = 'PDF';
		triggerToast('Report filters reset.');
	}

	function downloadReport(report: GeneratedReport) {
		triggerToast(`Downloading "${report.name}"...`);
	}

	function retryReport(index: number) {
		recentReports[index].status = 'Processing';
		triggerToast(`Retrying "${recentReports[index].name}"...`);
	}

	function useTemplate(name: string) {
		triggerToast(`Template "${name}" loaded into the generator.`);
	}

	function scheduleNewReport() {
		triggerToast('Opening report scheduler...');
	}

	function runExport(name: string) {
		triggerToast(`${name} started. Preparing download...`);
	}

	function reportTypeBadge(type: string): string {
		switch (type) {
			case 'PDF':
				return 'bg-rose-50 text-rose-700 border-rose-100';
			case 'Excel':
				return 'bg-emerald-50 text-emerald-700 border-emerald-100';
			case 'CSV':
				return 'bg-blue-50 text-blue-700 border-blue-100';
			default:
				return 'bg-purple-50 text-purple-700 border-purple-100';
		}
	}

	function reportStatusClass(status: ReportStatus | string): string {
		switch (status) {
			case 'Ready':
			case 'Success':
				return 'bg-emerald-50 text-emerald-700 border-emerald-100';
			case 'Processing':
			case 'Scheduled':
				return 'bg-amber-50 text-amber-700 border-amber-100';
			default:
				return 'bg-rose-50 text-rose-700 border-rose-100';
		}
	}

	function reportStatusDot(status: ReportStatus | string): string {
		switch (status) {
			case 'Ready':
			case 'Success':
				return 'bg-emerald-500';
			case 'Processing':
			case 'Scheduled':
				return 'bg-amber-500';
			default:
				return 'bg-rose-500';
		}
	}

	function reportFormatClass(format: ReportFormat): string {
		switch (format) {
			case 'PDF':
				return 'bg-rose-50 text-rose-700 border-rose-100';
			case 'Excel':
				return 'bg-emerald-50 text-emerald-700 border-emerald-100';
			case 'CSV':
				return 'bg-blue-50 text-blue-700 border-blue-100';
		}
	}

	function toggleMobileSidebar() {
		isMobileSidebarOpen = !isMobileSidebarOpen;
	}

	function handleLogout() {
		window.location.href = '/';
	}
</script>

<svelte:head>
	<title>Super Admin Console | iSPARC</title>
	<meta
		name="description"
		content="iSPARC Super Admin Portal. Configure tracks, view activity logs, manage administration settings, and oversee platform status."
	/>
</svelte:head>

<div class="min-h-screen bg-[#F7F6F3] text-slate-800 flex font-sans animate-fade-in">
	<!-- ==================== SIDEBAR ==================== -->
	<!-- Desktop Sidebar -->
	<aside
		class="hidden lg:flex flex-col w-64 bg-white border-r border-slate-200 h-screen sticky top-0 shrink-0 select-none"
	>
		<!-- Logo area -->
		<div class="h-[72px] flex items-center px-6 gap-3.5 border-b border-slate-100">
			<!-- Red square badge with white SA -->
			<div
				class="w-9 h-9 bg-[#881B1B] rounded-lg flex items-center justify-center text-white font-bold font-serif text-[11px] shrink-0"
			>
				SA
			</div>
			<div class="flex flex-col">
				<span class="text-xl font-bold tracking-tight text-slate-900 font-serif">
					i<span class="text-[#881B1B]">SPARC</span>
				</span>
				<span class="text-[9px] font-bold text-slate-600 tracking-wider uppercase -mt-1">
					IIPS DAVV CELL
				</span>
			</div>
		</div>

		<!-- Menu Items -->
		<nav class="flex-grow py-6 px-4 space-y-1 overflow-y-auto">
			{#each menuItems as item}
				<button
					onclick={() => (currentTab = item.name)}
					class="w-full flex items-center gap-3.5 px-4 py-3 rounded-lg text-[13px] font-bold tracking-wide transition-all duration-200 {currentTab ===
					item.name
						? 'bg-[#881B1B]/10 text-[#881B1B] border-l-[3px] border-[#881B1B] rounded-l-none'
						: 'text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
						class="w-5 h-5"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d={item.icon} />
					</svg>
					{item.name}
				</button>
			{/each}
		</nav>

		<!-- Logout -->
		<div class="p-4 border-t border-slate-100">
			<button
				onclick={handleLogout}
				class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-[13px] font-bold text-slate-800 hover:bg-rose-50 hover:text-rose-600 transition-colors"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
					class="w-5 h-5"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
					/>
				</svg>
				Logout
			</button>
		</div>
	</aside>

	<!-- Mobile Sidebar Overlay Menu (Drawer) -->
	{#if isMobileSidebarOpen}
		<!-- Backdrop -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div
			onclick={toggleMobileSidebar}
			transition:fade={{ duration: 150 }}
			class="lg:hidden fixed inset-0 bg-black/40 z-40"
		></div>

		<!-- Drawer Content -->
		<aside
			transition:slide={{ axis: 'x', duration: 250 }}
			class="lg:hidden fixed left-0 top-0 bottom-0 w-64 bg-white z-50 flex flex-col h-full shadow-2xl"
		>
			<div class="h-[72px] flex items-center justify-between px-6 border-b border-slate-100">
				<div class="flex items-center gap-3">
					<div
						class="w-9 h-9 bg-[#881B1B] rounded-lg flex items-center justify-center text-white font-bold font-serif text-[11px] shrink-0"
					>
						SA
					</div>
					<div class="flex flex-col">
						<span class="text-xl font-bold tracking-tight text-slate-900 font-serif">
							i<span class="text-[#881B1B]">SPARC</span>
						</span>
						<span class="text-[9px] font-bold text-slate-600 tracking-wider uppercase -mt-1">
							IIPS DAVV CELL
						</span>
					</div>
				</div>
				<button
					onclick={toggleMobileSidebar}
					aria-label="Close sidebar"
					class="p-1 rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-600"
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

			<nav class="flex-grow py-6 px-4 space-y-1 overflow-y-auto">
				{#each menuItems as item}
					<button
						onclick={() => {
							currentTab = item.name;
							toggleMobileSidebar();
						}}
						class="w-full flex items-center gap-3.5 px-4 py-3 rounded-lg text-[13px] font-bold tracking-wide transition-all duration-200 {currentTab ===
						item.name
							? 'bg-[#881B1B]/10 text-[#881B1B] border-l-[3px] border-[#881B1B] rounded-l-none'
							: 'text-slate-800 hover:bg-slate-50 hover:text-slate-900'}"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
							class="w-5 h-5"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d={item.icon} />
						</svg>
						{item.name}
					</button>
				{/each}
			</nav>

			<div class="p-4 border-t border-slate-100">
				<button
					onclick={() => {
						handleLogout();
						toggleMobileSidebar();
					}}
					class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-[13px] font-bold text-slate-800 hover:bg-rose-50 hover:text-rose-600 transition-colors"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
						class="w-5 h-5"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
						/>
					</svg>
					Logout
				</button>
			</div>
		</aside>
	{/if}

	<!-- ==================== MAIN WORKSPACE ==================== -->
	<div class="flex-grow flex flex-col min-w-0">
		<!-- Top Navigation Header -->
		<header
			class="bg-white border-b border-slate-200 h-[72px] flex items-center justify-between px-4 sm:px-6 lg:px-8 sticky top-0 z-30 select-none"
		>
			<!-- Mobile sidebar toggle -->
			<div class="flex items-center gap-3">
				<button
					onclick={toggleMobileSidebar}
					aria-label="Open sidebar"
					class="lg:hidden p-2 rounded-lg text-slate-500 hover:bg-slate-100 hover:text-slate-900 focus:outline-none"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
						class="w-6 h-6"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
					</svg>
				</button>
				<div>
					<h1 class="text-lg sm:text-xl font-bold text-slate-900 font-serif leading-tight">
						{currentTab === 'Dashboard' ? 'Welcome Back, Super Admin!' : currentTab}
					</h1>
					<p
						class="text-[10px] sm:text-xs font-semibold text-slate-400 uppercase tracking-wider mt-0.5"
					>
						Tuesday, 23 June 2026
					</p>
				</div>
			</div>

			<!-- Action items -->
			<div class="flex items-center gap-3.5 relative">
				<!-- Search box -->
				<div class="hidden md:flex items-center relative">
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search students..."
						class="pl-4 pr-9 py-2 bg-slate-50 rounded-lg border border-slate-200 text-xs text-slate-800 focus:outline-none focus:border-slate-350 focus:bg-white w-52 transition-all focus:w-60"
					/>
					<span class="absolute right-3 text-slate-400">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
							stroke-width="2"
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

				<!-- Notification button -->
				<button
					onclick={() => {
						isNotificationsOpen = !isNotificationsOpen;
					}}
					class="w-10 h-10 border border-slate-250 bg-white rounded-lg flex items-center justify-center text-slate-500 hover:bg-slate-50 transition-colors relative"
				>
					<span class="sr-only">Notifications</span>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
						class="w-5 h-5"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
						/>
					</svg>
					<span class="absolute top-2 right-2.5 w-2 h-2 bg-[#881B1B] rounded-full"></span>
				</button>

				{#if isNotificationsOpen}
					<div
						transition:slide={{ duration: 150 }}
						class="absolute right-0 top-12 w-80 bg-white border border-slate-200 rounded-xl shadow-lg py-2.5 z-40 mt-1 animate-fade-in"
					>
						<h3
							class="px-4 py-1.5 text-xs font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100"
						>
							System Alerts
						</h3>
						<div class="max-h-64 overflow-y-auto">
							{#each notifications as notice}
								<div
									class="px-4 py-3 hover:bg-slate-50 border-b border-slate-50 last:border-b-0 flex gap-2.5"
								>
									{#if notice.unread}
										<span class="w-2.5 h-2.5 bg-[#881B1B] rounded-full mt-1.5 shrink-0"></span>
									{/if}
									<div class="flex-grow">
										<p class="text-xs text-slate-700 font-semibold">{notice.text}</p>
										<span class="text-[10px] text-slate-405 block mt-1">{notice.time}</span>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}

				<!-- Profile Pill -->
				<div class="flex items-center gap-3 pl-2 border-l border-slate-200 select-none">
					<div
						class="w-10 h-10 rounded-full bg-[#881B1B] text-white flex items-center justify-center font-bold text-sm border-2 border-white shadow-sm shrink-0 font-serif"
					>
						SA
					</div>
					<div class="hidden sm:flex flex-col">
						<span class="text-xs font-bold text-slate-900 leading-none">Super Admin</span>
						<span class="text-[9px] text-slate-400 font-bold tracking-wider block mt-1 uppercase">
							Platform Administrator
						</span>
					</div>
				</div>
			</div>
		</header>

		<!-- Main Content Body -->
		<main class="flex-grow p-4 sm:p-6 lg:p-8 space-y-6 overflow-y-auto max-w-7xl mx-auto w-full">
			{#if currentTab === 'Dashboard'}
				<!-- Dashboard Statistics Cards Grid (Step 2) -->
				<section
					class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6"
					aria-label="Super Admin metrics overview"
				>
					<!-- Card 1: Total Students -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
					>
						<div class="flex items-center justify-between">
							<span class="text-2xl font-bold font-serif text-slate-900">1,248</span>
							<div
								class="p-2.5 rounded-lg bg-blue-50 text-blue-600 border border-blue-100 animate-pulse"
							>
								<!-- People icon -->
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
										d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.109A11.386 11.386 0 0 1 10.089 21c-2.316 0-4.445-.69-6.22-1.879v-.003a4.125 4.125 0 0 1 7.533-2.493M15 19.128v-.003c0-1.112-.285-2.16-.786-3.07M14.214 16.058A9.396 9.396 0 0 0 10.089 15c-1.47 0-2.854.34-4.082.945M14.214 16.058a9.386 9.386 0 0 1 0 3.07"
									/>
								</svg>
							</div>
						</div>
						<div class="mt-4">
							<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
								Total Students
							</h3>
							<span class="text-[11px] font-bold text-slate-400 mt-1 block">+32 this semester</span>
						</div>
					</div>

					<!-- Card 2: Total Admins -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
					>
						<div class="flex items-center justify-between">
							<span class="text-2xl font-bold font-serif text-slate-900"
								>{recentUsers.filter((u) => u.role === 'Admin').length + 85}</span
							>
							<div class="p-2.5 rounded-lg bg-blue-50 text-blue-600 border border-blue-100">
								<!-- Single user admin icon -->
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
										d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z"
									/>
								</svg>
							</div>
						</div>
						<div class="mt-4">
							<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
								Total Mentors
							</h3>
							<span class="text-[11px] font-bold text-slate-400 mt-1 block">+5 this month</span>
						</div>
					</div>

					<!-- Card 3: Total Activities -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
					>
						<div class="flex items-center justify-between">
							<span class="text-2xl font-bold font-serif text-slate-900">3,412</span>
							<div class="p-2.5 rounded-lg bg-blue-50 text-blue-600 border border-blue-100">
								<!-- Graph icon -->
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
							<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
								Total Activities
							</h3>
							<span class="text-[11px] font-bold text-slate-400 mt-1 block">+110 this week</span>
						</div>
					</div>

					<!-- Card 4: Active Tracks -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
					>
						<div class="flex items-center justify-between">
							<span class="text-2xl font-bold font-serif text-slate-900">6</span>
							<div class="p-2.5 rounded-lg bg-blue-50 text-blue-600 border border-blue-100">
								<!-- Stack icon -->
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
										d="M6.429 9.75L2.25 12l4.179 2.25m0-4.5l5.571 3 5.571-3m-11.142 0L21.75 12l-4.179 2.25m0 0L21.75 16.5l-9.75 5.25-9.75-5.25L6.429 14.25m11.142 0L12 16.5m0-13.5L21.75 7.5 12 12.75 2.25 7.5 12 3z"
									/>
								</svg>
							</div>
						</div>
						<div class="mt-4">
							<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
								Active Tracks
							</h3>
							<span class="text-[11px] font-bold text-slate-400 mt-1 block">2 new this year</span>
						</div>
					</div>
				</section>

				<!-- Middle Grid section (Step 3 & 4) -->
				<div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
					<!-- User Management Overview (lg:col-span-2) -->
					<div
						class="lg:col-span-2 bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden"
					>
						<div
							class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/20"
						>
							<h3 class="text-sm font-bold font-serif text-slate-905">User Management Overview</h3>
							<button
								onclick={() => (currentTab = 'User Management')}
								class="text-[#881B1B] hover:underline text-xs font-bold uppercase tracking-wider"
							>
								View All
							</button>
						</div>
						<div class="overflow-x-auto">
							<table class="w-full text-left border-collapse">
								<thead>
									<tr
										class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
									>
										<th class="py-3.5 px-5">User Name</th>
										<th class="py-3.5 px-5">Role</th>
										<th class="py-3.5 px-5">Department</th>
										<th class="py-3.5 px-5">Status</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-slate-100 text-xs font-sans">
									{#each recentUsers as user}
										<tr class="hover:bg-slate-50/30 transition-colors">
											<td class="py-4 px-5">
												<div class="flex flex-col">
													<span class="font-bold text-slate-800">{user.name}</span>
													<span class="text-[10px] text-slate-400 font-semibold mt-0.5 select-all"
														>{user.id}</span
													>
												</div>
											</td>
											<td class="py-4 px-5">
												<span
													class="inline-flex items-center px-2 py-0.5 text-[10px] font-bold uppercase rounded-full border
													{user.role === 'Student'
														? 'bg-slate-50 text-slate-650 border-slate-200'
														: 'bg-blue-50 text-blue-700 border-blue-100'}"
												>
													{user.role === 'Admin' ? 'Mentor' : user.role}
												</span>
											</td>
											<td class="py-4 px-5 text-slate-500 font-semibold">{user.dept}</td>
											<td class="py-4 px-5">
												<span
													class="inline-flex items-center gap-1.5 font-bold
													{user.status === 'Active' ? 'text-emerald-600' : 'text-amber-600'}"
												>
													<span
														class="w-1.5 h-1.5 rounded-full shrink-0
														{user.status === 'Active' ? 'bg-emerald-600' : 'bg-amber-500 animate-pulse'}"
													></span>
													{user.status}
												</span>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</div>

					<!-- Quick Actions Panel (Step 4) -->
					<div class="bg-white border border-slate-200 rounded-xl p-5 shadow-xs space-y-4">
						<h3 class="text-sm font-bold font-serif text-slate-905">Quick Actions</h3>
						<div class="h-px bg-slate-100 my-2"></div>

						<div class="grid grid-cols-2 gap-4 select-none">
							<!-- Button 1: Create User (Solid Red) -->
							<button
								onclick={() => (isCreateUserModalOpen = true)}
								class="p-5 bg-[#C23A3A] hover:bg-[#B03131] text-white rounded-[20px] shadow-xs flex flex-col items-center justify-center text-center space-y-2.5 h-[160px] w-full transition duration-200 focus:outline-none"
							>
								<!-- User plus icon -->
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2"
									stroke="currentColor"
									class="w-6 h-6 text-white"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235A8.902 8.902 0 0110 18a8.902 8.902 0 016 1.235c.19.115.3.322.3.54v.725c0 .19-.153.344-.344.344H4.344A.344.344 0 014 20.5v-.725c0-.218.11-.425.3-.54z"
									/>
								</svg>
								<div class="space-y-1">
									<span class="font-extrabold text-xs block font-sans">Create User</span>
									<span class="text-[9px] text-white/80 font-medium block font-sans leading-tight"
										>Add student or mentor</span
									>
								</div>
							</button>

							<!-- Button 2: Create Activity -->
							<button
								onclick={() => (isCreateActivityModalOpen = true)}
								class="p-5 bg-[#EFECE9] hover:bg-[#E4DFDB] text-slate-800 rounded-[20px] flex flex-col items-center justify-center text-center space-y-2.5 h-[160px] w-full transition duration-200 focus:outline-none"
							>
								<!-- Activity graph icon -->
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2"
									stroke="currentColor"
									class="w-6 h-6 text-slate-700"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 015.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941"
									/>
								</svg>
								<div class="space-y-1">
									<span class="font-extrabold text-xs block font-sans text-slate-905"
										>Create Activity</span
									>
									<span
										class="text-[9px] text-slate-500 font-semibold block font-sans leading-tight"
										>Publish new activity</span
									>
								</div>
							</button>

							<!-- Button 3: Create Track -->
							<button
								onclick={() => (isCreateTrackModalOpen = true)}
								class="p-5 bg-[#EFECE9] hover:bg-[#E4DFDB] text-slate-800 rounded-[20px] flex flex-col items-center justify-center text-center space-y-2.5 h-[160px] w-full transition duration-200 focus:outline-none"
							>
								<!-- Track icon -->
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2"
									stroke="currentColor"
									class="w-6 h-6 text-slate-700"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M6.429 9.75L2.25 12l4.179 2.25m0-4.5l5.571 3 5.571-3m-11.142 0L21.75 12l-4.179 2.25m0 0L21.75 16.5l-9.75 5.25-9.75-5.25L6.429 14.25m11.142 0L12 16.5m0-13.5L21.75 7.5 12 12.75 2.25 7.5 12 3z"
									/>
								</svg>
								<div class="space-y-1">
									<span class="font-extrabold text-xs block font-sans text-slate-905"
										>Create Track</span
									>
									<span
										class="text-[9px] text-slate-500 font-semibold block font-sans leading-tight"
										>Set up a new track</span
									>
								</div>
							</button>

							<!-- Button 4: Generate Report -->
							<button
								onclick={() => (isGenerateReportModalOpen = true)}
								class="p-5 bg-[#EFECE9] hover:bg-[#E4DFDB] text-slate-800 rounded-[20px] flex flex-col items-center justify-center text-center space-y-2.5 h-[160px] w-full transition duration-200 focus:outline-none"
							>
								<!-- Document download icon -->
								<svg
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
									stroke-width="2"
									stroke="currentColor"
									class="w-6 h-6 text-slate-700"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
									/>
								</svg>
								<div class="space-y-1">
									<span class="font-extrabold text-xs block font-sans text-slate-905"
										>Generate Report</span
									>
									<span
										class="text-[9px] text-slate-500 font-semibold block font-sans leading-tight"
										>Download activity data</span
									>
								</div>
							</button>
						</div>
					</div>
				</div>

				<!-- Recent System Activities Table Section (Step 5) -->
				<section class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
					<div
						class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/20 select-none"
					>
						<div class="flex items-center gap-2">
							<h3 class="text-sm font-bold font-serif text-slate-905">Recent System Activities</h3>
							<span
								class="px-2 py-0.5 bg-rose-50 text-rose-705 font-extrabold text-[10px] uppercase rounded-md border border-rose-100"
							>
								+ 6 this week
							</span>
						</div>
						<button
							onclick={() => triggerToast('Navigating to full audit log...')}
							class="text-[#881B1B] hover:underline text-xs font-bold uppercase tracking-wider"
						>
							View Full Log
						</button>
					</div>

					<div class="overflow-x-auto">
						<table class="w-full text-left border-collapse">
							<thead>
								<tr
									class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
								>
									<th class="py-3.5 px-5">Activity</th>
									<th class="py-3.5 px-5">Type</th>
									<th class="py-3.5 px-5">Performed By</th>
									<th class="py-3.5 px-5">Date</th>
									<th class="py-3.5 px-5">Status</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-slate-100 text-xs font-sans">
								{#each recentLogs as log}
									<tr class="hover:bg-slate-50/30 transition-colors">
										<td class="py-4 px-5 font-bold text-slate-800">{log.activity}</td>
										<td class="py-4 px-5">
											<span
												class="inline-flex items-center px-2.5 py-0.5 text-[9px] font-bold uppercase rounded-md border
												{log.type === 'User Added'
													? 'bg-slate-100 text-slate-650 border-slate-200'
													: log.type === 'Activity Created'
														? 'bg-blue-50 text-blue-700 border-blue-100'
														: log.type === 'Track Updated'
															? 'bg-purple-50 text-purple-700 border-purple-100'
															: log.type === 'Announcement'
																? 'bg-amber-50 text-amber-705 border-amber-100'
																: 'bg-emerald-50 text-emerald-700 border-emerald-100'}"
											>
												{log.type}
											</span>
										</td>
										<td class="py-4 px-5 text-slate-600 font-semibold">{log.performedBy}</td>
										<td class="py-4 px-5 text-slate-500 font-semibold">{log.date}</td>
										<td class="py-4 px-5">
											<span class="inline-flex items-center gap-1.5 font-bold text-emerald-600">
												<span class="w-1.5 h-1.5 rounded-full bg-emerald-600 shrink-0"></span>
												{log.status}
											</span>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>

					<div
						class="p-4 border-t border-slate-100 bg-slate-50/30 text-slate-500 font-semibold text-[11px] select-none"
					>
						<span>Showing {recentLogs.length} of {recentLogs.length} recent system activities</span>
					</div>
				</section>
			{:else if currentTab === 'Activity Management'}
				<ActivityManagementView />
			{:else if currentTab === 'Track Management'}
				<TrackManagementView />
			{:else if currentTab === 'User Management'}
				<!-- User Statistics Cards Grid (Step 2) -->
				<section
					class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 select-none"
					aria-label="User metrics overview"
				>
					<!-- Card 1: Total Students -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
					>
						<div class="flex items-center justify-between">
							<span class="text-2xl font-bold font-serif text-slate-900">{totalStudentsCount}</span>
							<div class="p-2.5 rounded-lg bg-rose-50 text-rose-600 border border-rose-100">
								<!-- Book open icon -->
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
										d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25"
									/>
								</svg>
							</div>
						</div>
						<div class="mt-4">
							<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
								Total Students
							</h3>
							<span class="text-[11px] font-bold text-slate-400 mt-1 block">+12 this semester</span>
						</div>
					</div>

					<!-- Card 2: Total Mentors -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
					>
						<div class="flex items-center justify-between">
							<span class="text-2xl font-bold font-serif text-slate-900">{totalMentorsCount}</span>
							<div class="p-2.5 rounded-lg bg-blue-50 text-blue-600 border border-blue-100">
								<!-- Users icon -->
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
										d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.109A11.386 11.386 0 0110.089 21c-2.316 0-4.445-.69-6.22-1.879v-.003a4.125 4.125 0 017.533-2.493M15 19.128v-.003c0-1.112-.285-2.16-.786-3.07M14.214 16.058A9.396 9.396 0 0010.089 15c-1.47 0-2.854.34-4.082.945M14.214 16.058a9.386 9.386 0 010 3.07"
									/>
								</svg>
							</div>
						</div>
						<div class="mt-4">
							<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
								Total Mentors
							</h3>
							<span class="text-[11px] font-bold text-slate-400 mt-1 block">+3 this month</span>
						</div>
					</div>

					<!-- Card 3: Total Admins -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
					>
						<div class="flex items-center justify-between">
							<span class="text-2xl font-bold font-serif text-slate-900">{totalAdminsCount}</span>
							<div class="p-2.5 rounded-lg bg-red-50 text-red-655 border border-red-100">
								<!-- Ribbon icon -->
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
										d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z"
									/>
								</svg>
							</div>
						</div>
						<div class="mt-4">
							<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
								Total Admins
							</h3>
							<span class="text-[11px] font-bold text-slate-400 mt-1 block">2 new this year</span>
						</div>
					</div>

					<!-- Card 4: Active Users -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex flex-col justify-between hover:shadow-md transition-shadow"
					>
						<div class="flex items-center justify-between">
							<span class="text-2xl font-bold font-serif text-slate-900">{activeUsersCount}</span>
							<div class="p-2.5 rounded-lg bg-[#E8F5E9] text-[#2E7D32] border border-[#C8E6C9]">
								<!-- Graph trending up icon -->
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
							<h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
								Active Users
							</h3>
							<span class="text-[11px] font-bold text-slate-400 mt-1 block">94.7% active rate</span>
						</div>
					</div>
				</section>

				<!-- Main 2-Column Grid (Step 3 & 4 Layout) -->
				<div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
					<!-- Left Column: User Management Table (lg:col-span-2) -->
					<div
						class="lg:col-span-2 bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden flex flex-col"
					>
						<!-- Table Header -->
						<div
							class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/20 select-none"
						>
							<h3 class="text-sm font-bold font-serif text-slate-905">User Management Overview</h3>
							<button
								type="button"
								onclick={() => {
									userFilter = 'All';
									userSearch = '';
								}}
								class="text-[#881B1B] hover:underline text-xs font-bold uppercase tracking-wider"
							>
								View All
							</button>
						</div>

						<!-- Table Filters & Search Controls -->
						<div
							class="p-5 border-b border-slate-100 flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white select-none"
						>
							<!-- Filter Tabs -->
							<div class="flex flex-wrap gap-1.5">
								{#each ['All', 'Student', 'Mentor', 'Admin'] as roleType}
									<button
										type="button"
										onclick={() => (userFilter = roleType)}
										class="px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all
											{userFilter === roleType
											? 'bg-[#C23A3A] text-white shadow-xs'
											: 'bg-slate-50 text-slate-500 hover:bg-slate-100'}"
									>
										{roleType === 'All'
											? 'All Users'
											: roleType === 'Student'
												? 'Students'
												: roleType === 'Mentor'
													? 'Mentors'
													: 'Admins'}
									</button>
								{/each}
							</div>

							<!-- Search Bar -->
							<div class="relative w-full sm:w-64">
								<input
									type="text"
									bind:value={userSearch}
									placeholder="Search by name, roll no., or email..."
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

						<!-- Table Grid -->
						<div class="overflow-x-auto">
							<table class="w-full text-left border-collapse">
								<thead>
									<tr
										class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider select-none"
									>
										<th class="py-3.5 px-5">User Name</th>
										<th class="py-3.5 px-5">Role</th>
										<th class="py-3.5 px-5">Department</th>
										<th class="py-3.5 px-5">Status</th>
										<th class="py-3.5 px-5 text-right sm:text-left">Actions</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-slate-100 text-xs font-sans">
									{#if filteredUsers.length === 0}
										<tr>
											<td
												colspan="5"
												class="py-8 text-center text-slate-400 font-semibold select-none"
											>
												No users found matching search filters.
											</td>
										</tr>
									{:else}
										{#each filteredUsers as user}
											<tr class="hover:bg-slate-50/30 transition-colors">
												<!-- User Info (Initials Avatar + Name & ID) -->
												<td class="py-4 px-5">
													<div class="flex items-center gap-3">
														<div
															class="w-8 h-8 rounded-full flex items-center justify-center font-bold text-xs uppercase shrink-0 select-none {user.avatarBg}"
														>
															{user.name
																.split(' ')
																.map((n) => n[0])
																.join('')
																.slice(0, 2)}
														</div>
														<div class="flex flex-col">
															<span class="font-bold text-slate-800">{user.name}</span>
															<span
																class="text-[9px] text-slate-400 font-semibold mt-0.5 select-all"
																>{user.id}</span
															>
														</div>
													</div>
												</td>

												<!-- Role Badge -->
												<td class="py-4 px-5 select-none">
													<span
														class="inline-flex items-center px-2 py-0.5 text-[9px] font-bold uppercase rounded-full border
														{user.displayRole === 'Student'
															? 'bg-blue-50 text-blue-700 border-blue-100'
															: user.displayRole === 'Mentor'
																? 'bg-emerald-50 text-emerald-700 border-emerald-100'
																: 'bg-rose-50 text-rose-700 border-rose-100'}"
													>
														{user.displayRole}
													</span>
												</td>

												<!-- Department -->
												<td class="py-4 px-5 text-slate-500 font-semibold">{user.dept}</td>

												<!-- Status Dot -->
												<td class="py-4 px-5 select-none">
													<span
														class="inline-flex items-center gap-1.5 font-bold
														{user.status === 'Active'
															? 'text-emerald-600'
															: user.status === 'Pending'
																? 'text-amber-600'
																: 'text-slate-400'}"
													>
														<span
															class="w-1.5 h-1.5 rounded-full shrink-0
															{user.status === 'Active'
																? 'bg-emerald-600'
																: user.status === 'Pending'
																	? 'bg-amber-500 animate-pulse'
																	: 'bg-slate-400'}"
														></span>
														{user.status}
													</span>
												</td>

												<!-- Action Buttons -->
												<td class="py-4 px-5">
													<div class="flex items-center gap-3 text-slate-400">
														<button
															type="button"
															onclick={() => triggerToast(`Viewing stats for ${user.name}...`)}
															aria-label="View user profile"
															class="hover:text-[#881B1B] transition-colors p-0.5"
														>
															<!-- Eye icon -->
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
															onclick={() =>
																triggerToast(`Edit functionality for ${user.name} is coming soon.`)}
															aria-label="Edit user details"
															class="hover:text-blue-600 transition-colors p-0.5"
														>
															<!-- Pencil edit icon -->
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
																	d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
																/>
															</svg>
														</button>
														<button
															type="button"
															onclick={() => {
																if (confirm(`Are you sure you want to delete ${user.name}?`)) {
																	userRegistry = userRegistry.filter((u) => u.id !== user.id);
																	triggerToast(`User "${user.name}" removed successfully.`);
																}
															}}
															aria-label="Delete user"
															class="hover:text-red-650 transition-colors p-0.5"
														>
															<!-- Trash delete icon -->
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

						<!-- Table Footer Summary -->
						<div
							class="p-4 border-t border-slate-100 bg-slate-50/30 text-slate-500 font-semibold text-[11px] select-none"
						>
							<span>Showing {filteredUsers.length} of {userRegistry.length} users</span>
						</div>
					</div>

					<!-- Right Column Sidebar Panels (Step 4 Cards) -->
					<div class="space-y-6">
						<!-- Quick Actions Card -->
						<div class="bg-white border border-slate-200 rounded-xl p-5 shadow-xs space-y-4">
							<h3 class="text-sm font-bold font-serif text-slate-905">Quick Actions</h3>
							<div class="h-px bg-slate-100 my-2"></div>

							<div class="grid grid-cols-2 gap-4 select-none">
								<!-- Button 1: Add User (Solid Red) -->
								<button
									type="button"
									onclick={() => (isCreateUserModalOpen = true)}
									class="p-4 bg-[#C23A3A] hover:bg-[#B03131] text-white rounded-2xl shadow-xs flex flex-col items-center justify-center text-center space-y-2 h-[100px] w-full transition duration-200 focus:outline-none"
								>
									<!-- User plus icon -->
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2"
										stroke="currentColor"
										class="w-5 h-5 text-white"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235A8.902 8.902 0 0110 18a8.902 8.902 0 016 1.235c.19.115.3.322.3.54v.725c0 .19-.153.344-.344.344H4.344A.344.344 0 014 20.5v-.725c0-.218.11-.425.3-.54z"
										/>
									</svg>
									<span class="font-extrabold text-[11px] block font-sans">Add User</span>
								</button>

								<!-- Button 2: Assign Mentor -->
								<button
									type="button"
									onclick={() => triggerToast('Assign Mentor flow coming soon...')}
									class="p-4 bg-[#F5F2F0] hover:bg-[#EAE5DF] text-slate-800 rounded-2xl flex flex-col items-center justify-center text-center space-y-2 h-[100px] w-full transition duration-200 focus:outline-none"
								>
									<!-- Link icon -->
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2"
										stroke="currentColor"
										class="w-5 h-5 text-[#C23A3A]"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244"
										/>
									</svg>
									<span class="font-extrabold text-[11px] block font-sans text-slate-805"
										>Assign Mentor</span
									>
								</button>

								<!-- Button 3: Assign Role -->
								<button
									type="button"
									onclick={() => triggerToast('Assign Role flow coming soon...')}
									class="p-4 bg-[#F5F2F0] hover:bg-[#EAE5DF] text-slate-800 rounded-2xl flex flex-col items-center justify-center text-center space-y-2 h-[100px] w-full transition duration-200 focus:outline-none"
								>
									<!-- Shield / user cog icon -->
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2"
										stroke="currentColor"
										class="w-5 h-5 text-[#C23A3A]"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z"
										/>
									</svg>
									<span class="font-extrabold text-[11px] block font-sans text-slate-805"
										>Assign Role</span
									>
								</button>

								<!-- Button 4: User Stats -->
								<button
									type="button"
									onclick={() => triggerToast('User statistics view selected.')}
									class="p-4 bg-[#F5F2F0] hover:bg-[#EAE5DF] text-slate-800 rounded-2xl flex flex-col items-center justify-center text-center space-y-2 h-[100px] w-full transition duration-200 focus:outline-none"
								>
									<!-- Speedometer / chart icon -->
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2"
										stroke="currentColor"
										class="w-5 h-5 text-[#C23A3A]"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M10.5 6a7.5 7.5 0 107.5 7.5h-7.5V6z"
										/>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="M13.5 10.5H21A7.5 7.5 0 0013.5 3v7.5z"
										/>
									</svg>
									<span class="font-extrabold text-[11px] block font-sans text-slate-805"
										>User Stats</span
									>
								</button>
							</div>
						</div>

						<!-- Account Status Summary Card -->
						<div
							class="bg-white border border-slate-200 rounded-xl p-5 shadow-xs space-y-4 select-none"
						>
							<h3 class="text-sm font-bold font-serif text-slate-905">Account Status</h3>
							<div class="h-px bg-slate-100 my-2"></div>

							<div class="space-y-3 font-sans text-xs">
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2 font-semibold text-slate-600">
										<span class="w-2 h-2 rounded-full bg-emerald-500"></span>
										Active
									</div>
									<span class="font-extrabold text-slate-800">{activeUsersCount}</span>
								</div>
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2 font-semibold text-slate-600">
										<span class="w-2 h-2 rounded-full bg-amber-500"></span>
										Pending
									</div>
									<span class="font-extrabold text-[#C06014]">{pendingUsersCount}</span>
								</div>
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2 font-semibold text-slate-600">
										<span class="w-2 h-2 rounded-full bg-slate-400"></span>
										Inactive
									</div>
									<span class="font-extrabold text-slate-400">{inactiveUsersCount}</span>
								</div>
							</div>
						</div>
					</div>
				</div>
			{:else if currentTab === 'Reports Center'}
				<!-- Reports Center (Step 7) -->
				<!-- Overview Stat Cards -->
				<section
					class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 select-none"
					aria-label="Reports overview"
				>
					<!-- Total Reports -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex items-center justify-between hover:shadow-md transition-shadow"
					>
						<div>
							<span class="text-2xl font-bold font-serif text-slate-900">156</span>
							<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1.5">Total Reports</h3>
							<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
								Generated all-time
							</p>
						</div>
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
									d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
								/>
							</svg>
						</div>
					</div>

					<!-- Scheduled -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex items-center justify-between hover:shadow-md transition-shadow"
					>
						<div>
							<span class="text-2xl font-bold font-serif text-slate-900"
								>{scheduledReports.length + 8}</span
							>
							<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1.5">
								Scheduled Reports
							</h3>
							<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
								Auto-generated
							</p>
						</div>
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
									d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
								/>
							</svg>
						</div>
					</div>

					<!-- Downloads -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex items-center justify-between hover:shadow-md transition-shadow"
					>
						<div>
							<span class="text-2xl font-bold font-serif text-slate-900">84</span>
							<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1.5">Downloads</h3>
							<p class="text-[10px] font-bold text-emerald-500 uppercase tracking-wider mt-0.5">
								+18 this month
							</p>
						</div>
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
									d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
								/>
							</svg>
						</div>
					</div>

					<!-- Last Generated -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex items-center justify-between hover:shadow-md transition-shadow"
					>
						<div>
							<span class="text-2xl font-bold font-serif text-slate-900">Today</span>
							<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1.5">Last Generated</h3>
							<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
								Jun 28, 2026
							</p>
						</div>
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
									d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5"
								/>
							</svg>
						</div>
					</div>
				</section>

				<!-- Generate Report + Quick Templates -->
				<div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
					<!-- Generate Report (col-span-2) -->
					<div
						class="lg:col-span-2 bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden"
					>
						<div class="p-5 border-b border-slate-100 bg-slate-50/20 select-none">
							<h3 class="text-sm font-bold font-serif text-slate-905">Generate Report</h3>
							<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
								Configure parameters and export your institutional report.
							</p>
						</div>
						<div class="p-5 space-y-5">
							<!-- Report Type -->
							<div class="flex flex-col gap-1.5">
								<label
									for="report-type"
									class="text-[10px] font-extrabold text-slate-650 tracking-wider"
									>REPORT TYPE</label
								>
								<select
									id="report-type"
									bind:value={reportType}
									class="px-3 py-2 border border-slate-200 rounded-lg text-xs font-semibold text-slate-800 bg-white focus:outline-none focus:border-slate-355"
								>
									{#each reportTypes as type}
										<option value={type}>{type}</option>
									{/each}
								</select>
							</div>

							<!-- Course / Semester -->
							<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
								<div class="flex flex-col gap-1.5">
									<label
										for="report-course"
										class="text-[10px] font-extrabold text-slate-650 tracking-wider"
										>COURSE / BATCH</label
									>
									<select
										id="report-course"
										bind:value={reportCourse}
										class="px-3 py-2 border border-slate-200 rounded-lg text-xs font-semibold text-slate-800 bg-white focus:outline-none focus:border-slate-355"
									>
										{#each reportCourses as course}
											<option value={course}>{course}</option>
										{/each}
									</select>
								</div>
								<div class="flex flex-col gap-1.5">
									<label
										for="report-semester"
										class="text-[10px] font-extrabold text-slate-650 tracking-wider">SEMESTER</label
									>
									<select
										id="report-semester"
										bind:value={reportSemester}
										class="px-3 py-2 border border-slate-200 rounded-lg text-xs font-semibold text-slate-800 bg-white focus:outline-none focus:border-slate-355"
									>
										{#each reportSemesters as sem}
											<option value={sem}>{sem}</option>
										{/each}
									</select>
								</div>
							</div>

							<!-- From / To -->
							<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
								<div class="flex flex-col gap-1.5">
									<label
										for="report-from"
										class="text-[10px] font-extrabold text-slate-650 tracking-wider">FROM</label
									>
									<input
										id="report-from"
										type="date"
										bind:value={reportFrom}
										class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
									/>
								</div>
								<div class="flex flex-col gap-1.5">
									<label
										for="report-to"
										class="text-[10px] font-extrabold text-slate-650 tracking-wider">TO</label
									>
									<input
										id="report-to"
										type="date"
										bind:value={reportTo}
										class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
									/>
								</div>
							</div>

							<!-- Format -->
							<div class="flex flex-col gap-2">
								<span class="text-[10px] font-extrabold text-slate-650 tracking-wider">FORMAT</span>
								<div class="flex flex-wrap gap-2">
									{#each reportFormats as fmt}
										<button
											type="button"
											onclick={() => (reportFormat = fmt)}
											class="px-4 py-2 rounded-lg text-xs font-bold border transition-all
											{reportFormat === fmt
												? 'bg-[#881B1B]/10 text-[#881B1B] border-[#881B1B]/30'
												: 'bg-white text-slate-500 border-slate-200 hover:bg-slate-50'}"
										>
											{fmt === 'Excel' ? 'Excel (.xlsx)' : fmt}
										</button>
									{/each}
								</div>
							</div>

							<!-- Actions -->
							<div class="flex items-center gap-3 pt-1">
								<button
									type="button"
									onclick={runReportGeneration}
									class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-[#C23A3A] hover:bg-[#B03131] text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-colors focus:outline-none"
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
											d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
										/>
									</svg>
									Generate Report
								</button>
								<button
									type="button"
									onclick={resetReportForm}
									class="px-4 py-2.5 border border-slate-200 hover:bg-slate-50 text-slate-600 font-bold text-xs uppercase tracking-wider rounded-lg transition-colors focus:outline-none"
								>
									Reset
								</button>
							</div>
						</div>
					</div>

					<!-- Quick Templates (col-span-1) -->
					<div class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
						<div class="p-5 border-b border-slate-100 bg-slate-50/20 select-none">
							<h3 class="text-sm font-bold font-serif text-slate-905">Quick Templates</h3>
							<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
								One-click report presets
							</p>
						</div>
						<div class="p-3 space-y-1 max-h-[420px] overflow-y-auto">
							{#each quickTemplates as tpl}
								<button
									type="button"
									onclick={() => useTemplate(tpl.name)}
									class="w-full flex items-center gap-3 p-2.5 rounded-lg hover:bg-slate-50 transition-colors text-left group"
								>
									<div
										class="w-8 h-8 rounded-lg bg-[#881B1B]/10 text-[#881B1B] flex items-center justify-center shrink-0 group-hover:bg-[#881B1B]/15 transition-colors"
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
												d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z"
											/>
										</svg>
									</div>
									<div class="min-w-0 flex-grow">
										<span class="block text-xs font-bold text-slate-800 truncate">{tpl.name}</span>
										<span class="block text-[10px] text-slate-400 font-semibold truncate"
											>{tpl.desc}</span
										>
									</div>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										fill="none"
										viewBox="0 0 24 24"
										stroke-width="2.5"
										stroke="currentColor"
										class="w-3.5 h-3.5 text-slate-300 group-hover:text-[#881B1B] transition-colors shrink-0"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											d="m8.25 4.5 7.5 7.5-7.5 7.5"
										/>
									</svg>
								</button>
							{/each}
						</div>
					</div>
				</div>

				<!-- Recent Reports -->
				<section class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
					<div
						class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/20 select-none"
					>
						<div class="flex items-center gap-2">
							<h3 class="text-sm font-bold font-serif text-slate-905">Recent Reports</h3>
							<span
								class="px-2 py-0.5 bg-rose-50 text-rose-705 font-extrabold text-[10px] uppercase rounded-md border border-rose-100"
							>
								{recentReports.length} total
							</span>
						</div>
						<button
							onclick={() => triggerToast('Loading full report archive...')}
							class="text-[#881B1B] hover:underline text-xs font-bold uppercase tracking-wider"
						>
							View All Reports
						</button>
					</div>
					<div class="overflow-x-auto">
						<table class="w-full text-left border-collapse">
							<thead>
								<tr
									class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
								>
									<th class="py-3.5 px-5">Report Name</th>
									<th class="py-3.5 px-5">Generated By</th>
									<th class="py-3.5 px-5">Date</th>
									<th class="py-3.5 px-5">Format</th>
									<th class="py-3.5 px-5">Status</th>
									<th class="py-3.5 px-5 text-center">Actions</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-slate-100 text-xs font-sans">
								{#each recentReports as report, index (report.name + index)}
									<tr class="hover:bg-slate-50/30 transition-colors">
										<td class="py-4 px-5 font-bold text-slate-800">{report.name}</td>
										<td class="py-4 px-5 text-slate-600 font-semibold">{report.by}</td>
										<td class="py-4 px-5 text-slate-500 font-semibold whitespace-nowrap"
											>{report.date}</td
										>
										<td class="py-4 px-5">
											<span
												class="inline-flex px-2 py-0.5 text-[10px] font-bold uppercase rounded-md border {reportFormatClass(
													report.format
												)}"
											>
												{report.format}
											</span>
										</td>
										<td class="py-4 px-5">
											<span
												class="inline-flex items-center gap-1.5 px-2 py-0.5 text-[10px] font-bold uppercase rounded-full border {reportStatusClass(
													report.status
												)}"
											>
												<span
													class="w-1.5 h-1.5 rounded-full {reportStatusDot(
														report.status
													)} {report.status === 'Processing' ? 'animate-pulse' : ''}"
												></span>
												{report.status}
											</span>
										</td>
										<td class="py-4 px-5">
											<div class="flex items-center justify-center">
												{#if report.status === 'Ready'}
													<button
														type="button"
														onclick={() => downloadReport(report)}
														class="inline-flex items-center gap-1 px-2.5 py-1.5 text-[10px] font-bold text-white bg-[#C23A3A] hover:bg-[#B03131] rounded-lg transition-colors"
													>
														<svg
															xmlns="http://www.w3.org/2000/svg"
															fill="none"
															viewBox="0 0 24 24"
															stroke-width="2"
															stroke="currentColor"
															class="w-3 h-3"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
															/>
														</svg>
														Download
													</button>
												{:else if report.status === 'Processing'}
													<span class="text-[10px] font-bold text-amber-600 uppercase tracking-wide"
														>Processing…</span
													>
												{:else}
													<button
														type="button"
														onclick={() => retryReport(index)}
														class="inline-flex items-center gap-1 px-2.5 py-1.5 text-[10px] font-bold text-rose-600 border border-rose-200 bg-rose-50 hover:bg-rose-100 rounded-lg transition-colors"
													>
														<svg
															xmlns="http://www.w3.org/2000/svg"
															fill="none"
															viewBox="0 0 24 24"
															stroke-width="2"
															stroke="currentColor"
															class="w-3 h-3"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99"
															/>
														</svg>
														Retry
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

				<!-- Scheduled Reports + Export Center -->
				<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
					<!-- Scheduled Reports -->
					<div class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
						<div
							class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/20 select-none"
						>
							<h3 class="text-sm font-bold font-serif text-slate-905">Scheduled Reports</h3>
							<button
								type="button"
								onclick={scheduleNewReport}
								class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-[#C23A3A] hover:bg-[#B03131] text-white font-bold text-[10px] uppercase tracking-wider rounded-lg transition-colors focus:outline-none"
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
								Schedule New
							</button>
						</div>
						<div class="p-3 space-y-1">
							{#each scheduledReports as sched}
								<div
									class="flex items-center gap-3 p-3 rounded-lg hover:bg-slate-50 transition-colors"
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
												d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
											/>
										</svg>
									</div>
									<div class="min-w-0 flex-grow">
										<span class="block text-xs font-bold text-slate-800 truncate">{sched.name}</span
										>
										<span class="block text-[10px] text-slate-400 font-semibold truncate"
											>{sched.freq}</span
										>
									</div>
									<span
										class="inline-flex items-center gap-1.5 px-2 py-0.5 text-[10px] font-bold uppercase rounded-full border bg-emerald-50 text-emerald-700 border-emerald-100 shrink-0"
									>
										<span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
										Scheduled
									</span>
								</div>
							{/each}
						</div>
					</div>

					<!-- Export Center -->
					<div class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
						<div class="p-5 border-b border-slate-100 bg-slate-50/20 select-none">
							<h3 class="text-sm font-bold font-serif text-slate-905">Export Center</h3>
							<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
								Bulk export platform datasets
							</p>
						</div>
						<div class="p-4 grid grid-cols-1 sm:grid-cols-2 gap-3">
							{#each exportOptions as opt}
								<button
									type="button"
									onclick={() => runExport(opt.name)}
									class="p-4 rounded-xl border border-slate-150 hover:border-[#881B1B]/30 hover:bg-slate-50/60 transition-colors text-left flex flex-col gap-2 group"
								>
									<div
										class="w-9 h-9 rounded-lg bg-blue-50 text-blue-600 border border-blue-100 flex items-center justify-center"
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
												d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3"
											/>
										</svg>
									</div>
									<div>
										<span class="block text-xs font-bold text-slate-800">{opt.name}</span>
										<span class="block text-[10px] text-slate-400 font-semibold mt-0.5"
											>{opt.desc}</span
										>
										<span class="block text-[10px] text-[#881B1B] font-bold mt-1">{opt.count}</span>
									</div>
								</button>
							{/each}
						</div>
					</div>
				</div>

				<!-- Institutional Overview + Report Activity Log -->
				<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
					<!-- Institutional Overview -->
					<div class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
						<div class="p-5 border-b border-slate-100 bg-slate-50/20 select-none">
							<h3 class="text-sm font-bold font-serif text-slate-905">Institutional Overview</h3>
							<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
								Academic year 2025–26 at a glance
							</p>
						</div>
						<div class="p-5 space-y-6">
							<!-- Stat cards -->
							<div class="grid grid-cols-2 sm:grid-cols-4 gap-3 select-none">
								{#each institutionalStats as stat}
									<div
										class="text-center rounded-xl bg-slate-50 border border-slate-150 py-3.5 px-2"
									>
										<div class="text-xl font-bold font-serif {stat.color}">{stat.value}</div>
										<div
											class="text-[9px] font-bold text-slate-400 uppercase tracking-wider mt-1 leading-tight"
										>
											{stat.label}
										</div>
									</div>
								{/each}
							</div>

							<!-- Avg score by department (bar chart) -->
							<div>
								<span class="text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
									>Avg Score by Department</span
								>
								<div class="mt-3 flex gap-2">
									<div
										class="flex flex-col justify-between h-36 text-[9px] font-bold text-slate-300 text-right shrink-0"
									>
										<span>100</span><span>75</span><span>50</span><span>25</span><span>0</span>
									</div>
									<div class="flex-1 min-w-0">
										<div class="relative h-36">
											<div class="absolute inset-0 flex flex-col justify-between">
												{#each ['a', 'b', 'c', 'd', 'e'] as line (line)}
													<div class="border-t border-dashed border-slate-100"></div>
												{/each}
											</div>
											<div class="relative flex items-end justify-between gap-2 h-full">
												{#each deptScores as d}
													<div class="flex-1 flex items-end justify-center h-full">
														<div
															class="w-full max-w-[26px] rounded-t bg-gradient-to-t from-[#881B1B] to-[#C23A3A] transition-all"
															style="height: {d.score}%"
															title="{d.dept}: {d.score}"
														></div>
													</div>
												{/each}
											</div>
										</div>
										<div class="flex justify-between gap-2 mt-1.5">
											{#each deptScores as d}
												<span class="flex-1 text-center text-[9px] font-bold text-slate-400"
													>{d.dept}</span
												>
											{/each}
										</div>
									</div>
								</div>
							</div>

							<!-- Monthly participation trend (dual line) -->
							<div>
								<div class="flex items-center justify-between">
									<span class="text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
										>Monthly Participation Trend</span
									>
									<div class="flex items-center gap-3">
										<span class="flex items-center gap-1 text-[9px] font-bold text-slate-500">
											<span class="w-2.5 h-0.5 rounded-full bg-[#C23A3A]"></span>Activities
										</span>
										<span class="flex items-center gap-1 text-[9px] font-bold text-slate-500">
											<span class="w-2.5 h-0.5 rounded-full bg-emerald-500"></span>Certificates
										</span>
									</div>
								</div>
								<div class="mt-3 flex gap-2">
									<div
										class="flex flex-col justify-between h-24 text-[9px] font-bold text-slate-300 text-right shrink-0"
									>
										<span>80</span><span>60</span><span>40</span><span>20</span>
									</div>
									<div class="flex-1 min-w-0">
										<div class="relative h-24">
											<div class="absolute inset-0 flex flex-col justify-between">
												{#each ['a', 'b', 'c', 'd'] as line (line)}
													<div class="border-t border-dashed border-slate-100"></div>
												{/each}
											</div>
											<svg
												viewBox="0 0 100 40"
												preserveAspectRatio="none"
												class="relative w-full h-full overflow-visible"
											>
												<polyline
													points={activitiesPoints}
													fill="none"
													stroke="#C23A3A"
													stroke-width="1.5"
													stroke-linecap="round"
													stroke-linejoin="round"
													vector-effect="non-scaling-stroke"
												/>
												<polyline
													points={certificatesPoints}
													fill="none"
													stroke="#10b981"
													stroke-width="1.5"
													stroke-linecap="round"
													stroke-linejoin="round"
													vector-effect="non-scaling-stroke"
												/>
											</svg>
										</div>
										<div class="flex justify-between gap-2 mt-1.5">
											{#each participationMonths as m (m)}
												<span class="flex-1 text-center text-[9px] font-bold text-slate-400"
													>{m}</span
												>
											{/each}
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>

					<!-- Report Activity Log -->
					<div class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
						<div
							class="p-5 border-b border-slate-100 flex items-center justify-between bg-slate-50/20 select-none"
						>
							<div>
								<h3 class="text-sm font-bold font-serif text-slate-905">Report Activity Log</h3>
								<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
									Audit trail of recent report actions
								</p>
							</div>
							<button
								type="button"
								onclick={() => triggerToast('Loading full activity log...')}
								class="px-3 py-1.5 border border-slate-200 hover:bg-slate-50 text-slate-600 font-bold text-[10px] uppercase tracking-wider rounded-lg transition-colors focus:outline-none"
							>
								View All
							</button>
						</div>
						<div class="overflow-x-auto">
							<table class="w-full text-left border-collapse">
								<thead>
									<tr
										class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
									>
										<th class="py-3 px-5">Action</th>
										<th class="py-3 px-4">Type</th>
										<th class="py-3 px-4">User</th>
										<th class="py-3 px-5">Timestamp</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-slate-100 text-xs font-sans">
									{#each reportAuditLog as entry (entry.action)}
										<tr class="hover:bg-slate-50/30 transition-colors">
											<td class="py-3.5 px-5">
												<div class="flex items-center gap-2.5">
													<span class="w-4 h-4 rounded-full border-2 border-slate-200 shrink-0"
													></span>
													<span class="font-bold text-slate-800">{entry.action}</span>
												</div>
											</td>
											<td class="py-3.5 px-4">
												<span
													class="inline-flex px-2 py-0.5 text-[10px] font-bold uppercase rounded-md border {reportTypeBadge(
														entry.type
													)}"
												>
													{entry.type}
												</span>
											</td>
											<td class="py-3.5 px-4 text-slate-600 font-semibold whitespace-nowrap"
												>{entry.user}</td
											>
											<td class="py-3.5 px-5 text-slate-500 font-semibold whitespace-nowrap"
												>{entry.time}</td
											>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</div>
				</div>
			{:else if currentTab === 'System Settings'}
				<!-- System Settings (Step 7) -->
				<!-- Overview Stat Cards -->
				<section
					class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 select-none"
					aria-label="System settings overview"
				>
					<!-- Academic Year -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex items-center justify-between hover:shadow-md transition-shadow"
					>
						<div>
							<span class="text-2xl font-bold font-serif text-slate-900">2025–26</span>
							<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1.5">Academic Year</h3>
							<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
								Current cycle
							</p>
						</div>
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
									d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5"
								/>
							</svg>
						</div>
					</div>

					<!-- Credit Policies -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex items-center justify-between hover:shadow-md transition-shadow"
					>
						<div>
							<span class="text-2xl font-bold font-serif text-slate-900"
								>{settingsSections['Credit Policy'].length + 3}</span
							>
							<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1.5">Credit Policies</h3>
							<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
								Rules active
							</p>
						</div>
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
									d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M7.73 9.728a6.726 6.726 0 0 0 2.748 1.35m8.272-6.842V4.5c0 2.108-.966 3.99-2.48 5.228m2.48-5.492a46.32 46.32 0 0 1 2.916.52 6.003 6.003 0 0 1-5.395 4.972m0 0a6.726 6.726 0 0 1-2.749 1.35m0 0a6.772 6.772 0 0 1-3.044 0"
								/>
							</svg>
						</div>
					</div>

					<!-- Required Activities -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex items-center justify-between hover:shadow-md transition-shadow"
					>
						<div>
							<span class="text-2xl font-bold font-serif text-slate-900">12</span>
							<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1.5">
								Required Activities
							</h3>
							<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
								Configured
							</p>
						</div>
						<div class="p-2.5 rounded-lg bg-rose-50 text-rose-600 border border-rose-100">
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
						</div>
					</div>

					<!-- System Status -->
					<div
						class="bg-white border border-slate-200 rounded-xl p-6 shadow-xs flex items-center justify-between hover:shadow-md transition-shadow"
					>
						<div>
							<span class="text-2xl font-bold font-serif text-emerald-600">Active</span>
							<h3 class="text-xs font-bold text-slate-800 tracking-wide mt-1.5">System Status</h3>
							<p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-0.5">
								All systems normal
							</p>
						</div>
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
									d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z"
								/>
							</svg>
						</div>
					</div>
				</section>

				<!-- Overview + Platform Summary grid -->
				<div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
					<!-- System Settings Overview (col-span-2) -->
					<div
						class="lg:col-span-2 bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden flex flex-col"
					>
						<!-- Header -->
						<div
							class="p-5 border-b border-slate-100 flex items-center justify-between gap-3 bg-slate-50/20 select-none"
						>
							<div>
								<h3 class="text-sm font-bold font-serif text-slate-905">
									System Settings Overview
								</h3>
								<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
									Manage platform-wide configuration settings.
								</p>
							</div>
							<button
								type="button"
								onclick={handleSaveAllSettings}
								class="px-4 py-2 bg-[#C23A3A] hover:bg-[#B03131] text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-colors focus:outline-none shrink-0"
							>
								Save Changes
							</button>
						</div>

						<!-- Section Tabs -->
						<div
							class="px-5 pt-3 border-b border-slate-100 flex flex-wrap gap-1.5 bg-white select-none"
						>
							{#each settingsTabs as tab}
								<button
									type="button"
									onclick={() => (activeSettingsTab = tab)}
									class="px-3.5 py-2 rounded-t-lg text-xs font-bold transition-all border-b-2 -mb-px
									{activeSettingsTab === tab
										? 'text-[#881B1B] border-[#881B1B]'
										: 'text-slate-500 border-transparent hover:text-slate-800'}"
								>
									{tab}
								</button>
							{/each}
						</div>

						<!-- Settings Table -->
						<div class="overflow-x-auto flex-grow">
							<table class="w-full text-left border-collapse">
								<thead>
									<tr
										class="border-b border-slate-150 bg-slate-50/50 text-[10px] font-extrabold text-slate-405 uppercase tracking-wider"
									>
										<th class="py-3.5 px-5">Setting Name</th>
										<th class="py-3.5 px-5">Description</th>
										<th class="py-3.5 px-5">Value</th>
										<th class="py-3.5 px-5">Status</th>
										<th class="py-3.5 px-5 text-center">Actions</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-slate-100 text-xs font-sans">
									{#each currentSettings as setting, index (setting.name)}
										<tr class="hover:bg-slate-50/30 transition-colors">
											<td class="py-4 px-5 font-bold text-slate-800 align-top">{setting.name}</td>
											<td class="py-4 px-5 text-slate-500 font-semibold align-top max-w-xs">
												{setting.description}
											</td>
											<td class="py-4 px-5 font-bold text-slate-800 align-top whitespace-nowrap">
												{setting.value}
											</td>
											<td class="py-4 px-5 align-top">
												<span
													class="inline-flex items-center gap-1.5 px-2 py-0.5 text-[10px] font-bold uppercase rounded-full border {settingStatusClass(
														setting.status
													)}"
												>
													<span class="w-1.5 h-1.5 rounded-full {settingStatusDot(setting.status)}"
													></span>
													{setting.status}
												</span>
											</td>
											<td class="py-4 px-5 align-top">
												<div class="flex items-center justify-center">
													<button
														type="button"
														onclick={() => openEditSetting(index)}
														aria-label="Edit {setting.name}"
														class="p-1.5 border border-slate-200 text-slate-500 hover:text-[#881B1B] hover:bg-slate-100 rounded-lg transition-colors focus:outline-none"
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
																d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Z"
															/>
														</svg>
													</button>
												</div>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>

						<!-- Footer -->
						<div
							class="p-4 border-t border-slate-100 bg-slate-50/30 flex items-center justify-between gap-3 select-none"
						>
							<span class="text-slate-500 font-semibold text-[11px]">
								{currentSettings.length} settings in this section
							</span>
							<button
								type="button"
								onclick={handleResetSettings}
								class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-600 font-bold text-[11px] uppercase tracking-wider rounded-lg transition-colors focus:outline-none"
							>
								Reset to Default
							</button>
						</div>
					</div>

					<!-- Platform Summary (col-span-1) -->
					<div class="bg-white border border-slate-200 rounded-xl shadow-xs overflow-hidden">
						<div class="p-5 border-b border-slate-100 bg-slate-50/20 select-none">
							<h3 class="text-sm font-bold font-serif text-slate-905">Platform Summary</h3>
							<p class="text-[11px] text-slate-500 font-semibold mt-0.5">
								Global configuration overview
							</p>
						</div>
						<div class="p-5 space-y-1">
							<div
								class="flex items-center justify-between py-2.5 border-b border-slate-100 last:border-0"
							>
								<span class="text-xs font-semibold text-slate-500">Academic Year</span>
								<span class="text-xs font-bold text-slate-900">2025–26</span>
							</div>
							<div
								class="flex items-center justify-between py-2.5 border-b border-slate-100 last:border-0"
							>
								<span class="text-xs font-semibold text-slate-500">Credit Policies</span>
								<span class="text-xs font-bold text-slate-900"
									>{settingsSections['Credit Policy'].length + 3}</span
								>
							</div>
							<div
								class="flex items-center justify-between py-2.5 border-b border-slate-100 last:border-0"
							>
								<span class="text-xs font-semibold text-slate-500">Activity Rules</span>
								<span class="text-xs font-bold text-slate-900">12</span>
							</div>
							<div
								class="flex items-center justify-between py-2.5 border-b border-slate-100 last:border-0"
							>
								<span class="text-xs font-semibold text-slate-500">Notif. Templates</span>
								<span class="text-xs font-bold text-slate-900">5</span>
							</div>
							<div class="flex items-center justify-between py-2.5">
								<span class="text-xs font-semibold text-slate-500">System Status</span>
								<span class="text-xs font-bold text-emerald-600">Active</span>
							</div>

							<button
								type="button"
								onclick={() => triggerToast('Opening global settings editor...')}
								class="w-full mt-3 py-2.5 bg-[#C23A3A] hover:bg-[#B03131] text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-colors focus:outline-none"
							>
								Edit Global Settings
							</button>
						</div>
					</div>
				</div>
			{:else}
				<!-- Under Construction placeholder for other tabs -->
				<div
					class="bg-white border border-slate-200 rounded-xl p-8 sm:p-12 text-center max-w-2xl mx-auto my-12 shadow-xs space-y-6"
				>
					<div
						class="w-16 h-16 bg-[#881B1B]/10 text-[#881B1B] flex items-center justify-center rounded-full mx-auto"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="2"
							stroke="currentColor"
							class="w-8 h-8 animate-bounce"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
							/>
						</svg>
					</div>
					<div class="space-y-2">
						<h2 class="text-xl font-bold font-serif text-[#0b1535]">
							{currentTab} Under Construction
						</h2>
						<p class="text-xs text-slate-500 max-w-md mx-auto leading-relaxed">
							The features for the <strong>{currentTab}</strong> module are scheduled for development
							by other contributors.
						</p>
					</div>
					<button
						onclick={() => (currentTab = 'Dashboard')}
						class="px-5 py-2.5 bg-[#881B1B] hover:bg-[#881B1B]/95 text-white font-bold text-xs uppercase tracking-wider rounded-lg transition-colors"
					>
						Back to Dashboard
					</button>
				</div>
			{/if}
		</main>

		<!-- ==================== TOAST NOTIFICATION CONTAINER (Step 6) ==================== -->
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

		<!-- ==================== MODALS LIST (Step 6) ==================== -->
		<!-- Create User Modal -->
		{#if isCreateUserModalOpen}
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				onclick={(e) => {
					if (e.target === e.currentTarget) isCreateUserModalOpen = false;
				}}
				transition:fade={{ duration: 150 }}
				class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
			>
				<form
					onsubmit={handleCreateUser}
					class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
				>
					<div
						class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30"
					>
						<div>
							<h3 class="text-sm font-bold font-serif text-slate-900">Create New User</h3>
							<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
								Register Student or Mentor
							</p>
						</div>
						<button
							type="button"
							onclick={() => (isCreateUserModalOpen = false)}
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
								for="new-username"
								class="text-[10px] font-extrabold text-slate-650 tracking-wider">FULL NAME *</label
							>
							<input
								id="new-username"
								type="text"
								bind:value={newUserName}
								placeholder="e.g. Amit Kumar"
								required
								class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
							/>
						</div>

						<div class="grid grid-cols-2 gap-4">
							<div class="flex flex-col gap-1.5">
								<label
									for="new-user-role"
									class="text-[10px] font-extrabold text-slate-650 tracking-wider">ROLE *</label
								>
								<select
									id="new-user-role"
									bind:value={newUserRole}
									class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-850 bg-white focus:outline-none focus:border-slate-355"
								>
									<option value="Student">Student</option>
									<option value="Admin">Mentor</option>
								</select>
							</div>

							<div class="flex flex-col gap-1.5">
								<label
									for="new-user-dept"
									class="text-[10px] font-extrabold text-slate-650 tracking-wider">DEPARTMENT</label
								>
								<input
									id="new-user-dept"
									type="text"
									bind:value={newUserDept}
									placeholder="e.g. Computer Science"
									class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
								/>
							</div>
						</div>
					</div>

					<div
						class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
					>
						<button
							type="button"
							onclick={() => (isCreateUserModalOpen = false)}
							class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Create User
						</button>
					</div>
				</form>
			</div>
		{/if}

		<!-- Create Activity Modal -->
		{#if isCreateActivityModalOpen}
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				onclick={(e) => {
					if (e.target === e.currentTarget) isCreateActivityModalOpen = false;
				}}
				transition:fade={{ duration: 150 }}
				class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
			>
				<form
					onsubmit={handleCreateActivity}
					class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
				>
					<div
						class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30"
					>
						<div>
							<h3 class="text-sm font-bold font-serif text-slate-900">Publish New Activity</h3>
							<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
								Register Extracurricular Activity
							</p>
						</div>
						<button
							type="button"
							onclick={() => (isCreateActivityModalOpen = false)}
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
								for="new-act-name"
								class="text-[10px] font-extrabold text-slate-655 tracking-wider"
								>ACTIVITY NAME *</label
							>
							<input
								id="new-act-name"
								type="text"
								bind:value={newActivityName}
								placeholder="e.g. Cyber Security Summit"
								required
								class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
							/>
						</div>

						<div class="flex flex-col gap-1.5">
							<label
								for="new-act-credits"
								class="text-[10px] font-extrabold text-slate-655 tracking-wider"
								>CREDITS OFFERED *</label
							>
							<input
								id="new-act-credits"
								type="number"
								bind:value={newActivityCredits}
								min="1"
								max="50"
								required
								class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
							/>
						</div>
					</div>

					<div
						class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
					>
						<button
							type="button"
							onclick={() => (isCreateActivityModalOpen = false)}
							class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Publish Activity
						</button>
					</div>
				</form>
			</div>
		{/if}

		<!-- Create Track Modal -->
		{#if isCreateTrackModalOpen}
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				onclick={(e) => {
					if (e.target === e.currentTarget) isCreateTrackModalOpen = false;
				}}
				transition:fade={{ duration: 150 }}
				class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
			>
				<form
					onsubmit={handleCreateTrack}
					class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
				>
					<div
						class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30"
					>
						<div>
							<h3 class="text-sm font-bold font-serif text-slate-900">Create New Track</h3>
							<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
								Add specialized extracurricular track
							</p>
						</div>
						<button
							type="button"
							onclick={() => (isCreateTrackModalOpen = false)}
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
								class="text-[10px] font-extrabold text-slate-655 tracking-wider"
								>TRACK TITLE *</label
							>
							<input
								id="new-track-name"
								type="text"
								bind:value={newTrackName}
								placeholder="e.g. Social Entrepreneurship & Innovations"
								required
								class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
							/>
						</div>
					</div>

					<div
						class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
					>
						<button
							type="button"
							onclick={() => (isCreateTrackModalOpen = false)}
							class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Create Track
						</button>
					</div>
				</form>
			</div>
		{/if}

		<!-- Generate Report Modal -->
		{#if isGenerateReportModalOpen}
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				onclick={(e) => {
					if (e.target === e.currentTarget) isGenerateReportModalOpen = false;
				}}
				transition:fade={{ duration: 150 }}
				class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
			>
				<div
					class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
				>
					<div
						class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30"
					>
						<div>
							<h3 class="text-sm font-bold font-serif text-slate-905">Generate System Report</h3>
							<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
								Export extracurricular summaries
							</p>
						</div>
						<button
							type="button"
							onclick={() => (isGenerateReportModalOpen = false)}
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
								for="report-type"
								class="text-[10px] font-extrabold text-slate-655 tracking-wider">REPORT TYPE</label
							>
							<select
								id="report-type"
								class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 bg-white focus:outline-none focus:border-slate-355"
							>
								<option value="activities">Student Activities Summary</option>
								<option value="credits">Credit Distribution Ledger</option>
								<option value="audits">Administrative Audit Trails</option>
							</select>
						</div>

						<div class="grid grid-cols-2 gap-4">
							<div class="flex flex-col gap-1.5">
								<label
									for="report-from"
									class="text-[10px] font-extrabold text-slate-655 tracking-wider">START DATE</label
								>
								<input
									id="report-from"
									type="date"
									value="2026-06-01"
									class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
								/>
							</div>
							<div class="flex flex-col gap-1.5">
								<label
									for="report-to"
									class="text-[10px] font-extrabold text-slate-655 tracking-wider">END DATE</label
								>
								<input
									id="report-to"
									type="date"
									value="2026-06-30"
									class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
								/>
							</div>
						</div>
					</div>

					<div
						class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
					>
						<button
							type="button"
							onclick={() => (isGenerateReportModalOpen = false)}
							class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Cancel
						</button>
						<button
							type="button"
							onclick={handleGenerateReport}
							class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Export Data
						</button>
					</div>
				</div>
			</div>
		{/if}

		<!-- Edit Setting Modal (Step 7) -->
		{#if isEditSettingOpen}
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				onclick={(e) => {
					if (e.target === e.currentTarget) isEditSettingOpen = false;
				}}
				transition:fade={{ duration: 150 }}
				class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs"
			>
				<form
					onsubmit={handleSaveSetting}
					class="w-full max-w-md bg-white border border-slate-200 rounded-2xl shadow-2xl overflow-hidden flex flex-col font-sans"
				>
					<div
						class="p-5 border-b border-slate-150 flex items-center justify-between bg-slate-50/30"
					>
						<div>
							<h3 class="text-sm font-bold font-serif text-slate-900">Edit Setting</h3>
							<p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
								{activeSettingsTab}
							</p>
						</div>
						<button
							type="button"
							onclick={() => (isEditSettingOpen = false)}
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
							<span class="text-[10px] font-extrabold text-slate-650 tracking-wider"
								>SETTING NAME</span
							>
							<span class="text-sm font-bold text-slate-800">{editSettingName}</span>
						</div>
						<div class="flex flex-col gap-1.5">
							<label
								for="edit-setting-value"
								class="text-[10px] font-extrabold text-slate-650 tracking-wider">VALUE *</label
							>
							<input
								id="edit-setting-value"
								type="text"
								bind:value={editSettingValue}
								required
								class="px-3 py-2 border border-slate-200 rounded-lg text-xs text-slate-800 focus:outline-none focus:border-slate-355"
							/>
						</div>
					</div>

					<div
						class="p-5 border-t border-slate-150 flex items-center justify-end gap-2.5 bg-slate-50/30"
					>
						<button
							type="button"
							onclick={() => (isEditSettingOpen = false)}
							class="px-4 py-2 border border-slate-200 hover:bg-slate-50 text-slate-700 font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-[#881B1B] hover:bg-[#881B1B]/90 text-white font-bold text-xs uppercase rounded-lg transition-colors focus:outline-none"
						>
							Save Setting
						</button>
					</div>
				</form>
			</div>
		{/if}
	</div>
</div>
