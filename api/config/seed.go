package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
)

const DevPassword = "Pass@123"

const certificateSeedDir = "./uploads/certificates"

func SeedDevData() {
	enabled, err := strconv.ParseBool(os.Getenv("SEED_DEV_DATA"))
	if err != nil || !enabled {
		log.Println("SEED_DEV_DATA is not enabled, skipping demo data seeding")
		return
	}

	hashed, err := utils.HashPassword(DevPassword)
	if err != nil {
		log.Printf("Seeding failed: could not hash the demo password: %v", err)
		return
	}

	if err := seedAdmins(hashed); err != nil {
		log.Printf("Seeding admins failed: %v", err)
		return
	}

	students, err := seedStudents(hashed)
	if err != nil {
		log.Printf("Seeding students failed: %v", err)
		return
	}

	activities, err := seedActivities()
	if err != nil {
		log.Printf("Seeding activities failed: %v", err)
		return
	}

	if err := seedEnrollments(students, activities); err != nil {
		log.Printf("Seeding enrollments failed: %v", err)
		return
	}

	if err := seedCertificates(); err != nil {
		log.Printf("Seeding certificates failed: %v", err)
		return
	}

	if err := seedSystemSettings(); err != nil {
		log.Printf("Seeding system settings failed: %v", err)
		return
	}

	log.Printf("Demo data ready: %d students, %d activities. Every account's password is %q.",
		len(students), len(activities), DevPassword)
}

// ---------------------------------------------------------------------------
// System settings
// ---------------------------------------------------------------------------

// seedSystemSettings loads the default platform settings shown on the super
// admin System Settings screen. Value and Status are seeded via Attrs so they
// are only written on first create — a super admin's later edits are never
// overwritten on re-seed. The structural fields (category, name, description,
// order) are kept in sync via Assign so code changes propagate.
func seedSystemSettings() error {
	settings := []models.SystemSetting{
		// Academic Year
		{Key: "academic_year_current", Category: "Academic Year", Name: "Current Academic Year", Description: "Active academic cycle label displayed platform-wide", Value: "2025-2026", Status: "Active"},
		{Key: "academic_year_start", Category: "Academic Year", Name: "Academic Year Start Date", Description: "Official start date of the current academic year", Value: "Aug 1, 2025", Status: "Active"},
		{Key: "academic_year_end", Category: "Academic Year", Name: "Academic Year End Date", Description: "Official end date of the current academic year", Value: "May 31, 2026", Status: "Active"},
		{Key: "enrollment_deadline", Category: "Academic Year", Name: "Enrollment Deadline", Description: "Last date for activity enrollment submissions", Value: "Sep 15, 2025", Status: "Active"},

		// Credit Policy
		{Key: "credit_min_required", Category: "Credit Policy", Name: "Minimum Credits Required", Description: "Total credits a student must earn to graduate", Value: "100", Status: "Active"},
		{Key: "credit_max_per_activity", Category: "Credit Policy", Name: "Maximum Credits per Activity", Description: "Upper cap on credits from a single activity", Value: "20", Status: "Active"},
		{Key: "credit_min_per_semester", Category: "Credit Policy", Name: "Minimum Credits per Semester", Description: "Credits a student must earn each semester", Value: "12", Status: "Active"},
		{Key: "credit_rollover", Category: "Credit Policy", Name: "Credit Rollover", Description: "Carry unused credits to the next academic year", Value: "Enabled", Status: "Enabled"},
		{Key: "credit_grace_buffer", Category: "Credit Policy", Name: "Grace Credit Buffer", Description: "Extra credits allowed beyond the target", Value: "5", Status: "Active"},

		// Activity Rules
		{Key: "activity_max_enrollments", Category: "Activity Rules", Name: "Max Active Enrollments", Description: "Activities a student can be enrolled in at once", Value: "5", Status: "Active"},
		{Key: "activity_mentor_approval", Category: "Activity Rules", Name: "Mentor Approval Required", Description: "Require mentor sign-off before credit is granted", Value: "Enabled", Status: "Enabled"},
		{Key: "activity_auto_verify_threshold", Category: "Activity Rules", Name: "Auto-Verification Threshold", Description: "Score above which certificates auto-verify", Value: "90%", Status: "Active"},
		{Key: "activity_self_reported", Category: "Activity Rules", Name: "Self-Reported Activities", Description: "Allow students to log their own activities", Value: "Disabled", Status: "Disabled"},
		{Key: "activity_resubmission_window", Category: "Activity Rules", Name: "Resubmission Window", Description: "Days allowed to resubmit a rejected certificate", Value: "7 days", Status: "Active"},

		// Notifications
		{Key: "notify_email", Category: "Notifications", Name: "Email Notifications", Description: "Send system emails for key events", Value: "Enabled", Status: "Enabled"},
		{Key: "notify_otp_expiry", Category: "Notifications", Name: "OTP Expiry Duration", Description: "Validity window for verification codes", Value: "15 min", Status: "Active"},
		{Key: "notify_reminder_frequency", Category: "Notifications", Name: "Reminder Frequency", Description: "How often pending-task reminders are sent", Value: "Weekly", Status: "Active"},
		{Key: "notify_announcement_broadcast", Category: "Notifications", Name: "Announcement Broadcasts", Description: "Push platform-wide announcements to all users", Value: "Enabled", Status: "Enabled"},

		// Platform
		{Key: "platform_name", Category: "Platform", Name: "Platform Name", Description: "Display name shown across the portal", Value: "iSPARC", Status: "Active"},
		{Key: "platform_maintenance_mode", Category: "Platform", Name: "Maintenance Mode", Description: "Temporarily disable access for non-admins", Value: "Disabled", Status: "Disabled"},
		{Key: "platform_time_zone", Category: "Platform", Name: "Default Time Zone", Description: "Base time zone for schedules and logs", Value: "IST (UTC+5:30)", Status: "Active"},
		{Key: "platform_language", Category: "Platform", Name: "Default Language", Description: "Default interface language for new users", Value: "English", Status: "Active"},

		// Security
		{Key: "security_min_password_length", Category: "Security", Name: "Minimum Password Length", Description: "Required characters for user passwords", Value: "8", Status: "Active"},
		{Key: "security_session_timeout", Category: "Security", Name: "Session Timeout", Description: "Idle minutes before automatic logout", Value: "30 min", Status: "Active"},
		{Key: "security_two_factor", Category: "Security", Name: "Two-Factor Authentication", Description: "Require 2FA for admin accounts", Value: "Enabled", Status: "Enabled"},
		{Key: "security_max_login_attempts", Category: "Security", Name: "Max Login Attempts", Description: "Failed logins before temporary lockout", Value: "5", Status: "Active"},
	}

	for i, setting := range settings {
		var existing models.SystemSetting
		if err := DB.Where(models.SystemSetting{Key: setting.Key}).
			Attrs(models.SystemSetting{Value: setting.Value, Status: setting.Status}).
			Assign(models.SystemSetting{
				Category:    setting.Category,
				Name:        setting.Name,
				Description: setting.Description,
				SortOrder:   i,
			}).
			FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Admins
// ---------------------------------------------------------------------------

func seedAdmins(hashedPassword string) error {
	admins := []models.Admin{
		{
			AdminID:       "superadmin",
			Name:          "Ananya Deshmukh",
			Email:         "superadmin@isparc.dev",
			Role:          "superadmin",
			AssignedBatch: "", // A super admin is not scoped to a batch.
		},
		{
			AdminID:       "admin",
			Name:          "Dr. Rajesh Kumar",
			Email:         "rajesh.kumar@isparc.dev",
			Role:          "admin",
			AssignedBatch: "IT2K24",
		},
		{
			AdminID:       "admin2",
			Name:          "Dr. Priya Patel",
			Email:         "priya.patel@isparc.dev",
			Role:          "admin",
			AssignedBatch: "IT2K25",
		},
	}

	for _, admin := range admins {
		var existing models.Admin
		// MustChangePassword stays false for demo accounts so that logging in
		// lands straight on the dashboard instead of the reset screen.
		if err := DB.Where(models.Admin{AdminID: admin.AdminID}).
			Attrs(models.Admin{Password: hashedPassword, MustChangePassword: false}).
			Assign(models.Admin{
				Name:          admin.Name,
				Email:         admin.Email,
				Role:          admin.Role,
				AssignedBatch: admin.AssignedBatch,
			}).
			FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Students
// ---------------------------------------------------------------------------

func seedStudents(hashedPassword string) ([]models.Student, error) {
	students := []models.Student{
		{RollNo: "IT2K24011", Name: "Rahul Sharma", CourseName: "Computer Science", Semester: 6, ContactNo: "9876543210", EmailID: "rahul.sharma@iips.edu", EnrollmentNo: "EN-IT2K24011"},
		{RollNo: "IT2K24012", Name: "Sneha Kumar", CourseName: "Computer Science", Semester: 6, ContactNo: "9876543211", EmailID: "sneha.kumar@iips.edu", EnrollmentNo: "EN-IT2K24012"},
		{RollNo: "IT2K24013", Name: "Arjun Desai", CourseName: "Information Technology", Semester: 4, ContactNo: "9876543212", EmailID: "arjun.desai@iips.edu", EnrollmentNo: "EN-IT2K24013"},
		{RollNo: "IT2K24014", Name: "Kavya Krishnan", CourseName: "Computer Science", Semester: 4, ContactNo: "9876543213", EmailID: "kavya.krishnan@iips.edu", EnrollmentNo: "EN-IT2K24014"},
		{RollNo: "IT2K24015", Name: "Vikram Singh", CourseName: "Information Technology", Semester: 2, ContactNo: "9876543214", EmailID: "vikram.singh@iips.edu", EnrollmentNo: "EN-IT2K24015"},
		{RollNo: "IT2K25001", Name: "Priya Nair", CourseName: "MCA", Semester: 2, ContactNo: "9876543215", EmailID: "priya.nair@iips.edu", EnrollmentNo: "EN-IT2K25001"},
		{RollNo: "IT2K25002", Name: "Rohan Verma", CourseName: "MCA", Semester: 2, ContactNo: "9876543216", EmailID: "rohan.verma@iips.edu", EnrollmentNo: "EN-IT2K25002"},
		{RollNo: "IT2K25003", Name: "Meera Iyer", CourseName: "MCA", Semester: 4, ContactNo: "9876543217", EmailID: "meera.iyer@iips.edu", EnrollmentNo: "EN-IT2K25003"},
	}

	seeded := make([]models.Student, 0, len(students))

	for _, student := range students {
		var existing models.Student
		if err := DB.Where(models.Student{RollNo: student.RollNo}).
			Attrs(models.Student{Password: hashedPassword}).
			Assign(models.Student{
				Name:         student.Name,
				CourseName:   student.CourseName,
				Semester:     student.Semester,
				ContactNo:    student.ContactNo,
				EmailID:      student.EmailID,
				EnrollmentNo: student.EnrollmentNo,
				// Seeded students skip OTP: they are ready to log in.
				IsVerified: true,
			}).
			FirstOrCreate(&existing).Error; err != nil {
			return nil, err
		}

		seeded = append(seeded, existing)
	}

	return seeded, nil
}

// ---------------------------------------------------------------------------
// Activities
// ---------------------------------------------------------------------------

// Categories are upper case because that is what the activity catalogue in the
// student portal groups and filters on.
func seedActivities() ([]models.Activity, error) {
	now := time.Now()

	activities := []models.Activity{
		{
			Name: "National Hackathon 2026", Category: "TECHNICAL",
			Description: "A 36-hour coding challenge open to all students. Build solutions for real-world problems.",
			Credits:     15, Mode: "Offline", Venue: "IIPS Auditorium", Coordinator: "Dr. Rajesh Kumar", CoordinatorID: "admin",
			RegDeadline: now.AddDate(0, 0, 3), ActivityDate: now.AddDate(0, 0, 10), Status: "Closing Soon",
		},
		{
			Name: "National Science Olympiad", Category: "RESEARCH",
			Description: "National-level science competition covering physics, chemistry and biology.",
			Credits:     20, Mode: "Hybrid", Venue: "IIPS Seminar Hall", Coordinator: "Dr. Priya Patel", CoordinatorID: "admin2",
			RegDeadline: now.AddDate(0, 0, 14), ActivityDate: now.AddDate(0, 0, 21), Status: "Open",
		},
		{
			Name: "Inter-College Athletics Meet", Category: "SPORTS",
			Description: "Annual inter-college athletics championship. Represent IIPS in track and field events.",
			Credits:     10, Mode: "Offline", Venue: "DAVV Sports Ground", Coordinator: "Prof. Anjali Sharma",
			RegDeadline: now.AddDate(0, 0, 7), ActivityDate: now.AddDate(0, 0, 12), Status: "Open",
		},
		{
			Name: "Cultural Fest - Rangmanch", Category: "CULTURAL",
			Description: "Annual cultural festival with music, dance and theatre performances.",
			Credits:     10, Mode: "Offline", Venue: "IIPS Open Air Theatre", Coordinator: "Prof. Anjali Sharma",
			RegDeadline: now.AddDate(0, 0, 9), ActivityDate: now.AddDate(0, 0, 16), Status: "Open",
		},
		{
			Name: "Student Leadership Workshop", Category: "LEADERSHIP",
			Description: "Leadership development workshop covering team building and decision making.",
			Credits:     10, Mode: "Offline", Venue: "IIPS Seminar Hall", Coordinator: "Dr. Mehta",
			RegDeadline: now.AddDate(0, 0, 5), ActivityDate: now.AddDate(0, 0, 11), Status: "Open",
		},
		{
			Name: "Inter College Debate Championship", Category: "PUBLIC SPEAKING",
			Description: "Parliamentary-style debate on contemporary socio-political topics.",
			Credits:     12, Mode: "Offline", Venue: "IIPS Conference Hall", Coordinator: "Dr. Rajesh Kumar", CoordinatorID: "admin",
			RegDeadline: now.AddDate(0, 0, 4), ActivityDate: now.AddDate(0, 0, 9), Status: "Closing Soon",
		},
		{
			Name: "Blood Donation Camp", Category: "SOCIAL SERVICE",
			Description: "Community health initiative with District Hospital Indore. Volunteers earn social service credit.",
			Credits:     8, Mode: "Offline", Venue: "IIPS Main Ground", Coordinator: "NSS Cell",
			RegDeadline: now.AddDate(0, 0, -2), ActivityDate: now.AddDate(0, 0, 1), Status: "Closed",
		},
	}

	seeded := make([]models.Activity, 0, len(activities))

	for _, activity := range activities {
		var existing models.Activity
		if err := DB.Where(models.Activity{Name: activity.Name}).
			Assign(models.Activity{
				Category:      activity.Category,
				Description:   activity.Description,
				Credits:       activity.Credits,
				Mode:          activity.Mode,
				Venue:         activity.Venue,
				Coordinator:   activity.Coordinator,
				CoordinatorID: activity.CoordinatorID,
				RegDeadline:   activity.RegDeadline,
				ActivityDate:  activity.ActivityDate,
				Status:        activity.Status,
			}).
			FirstOrCreate(&existing).Error; err != nil {
			return nil, err
		}

		seeded = append(seeded, existing)
	}

	return seeded, nil
}

// ---------------------------------------------------------------------------
// Enrollments
// ---------------------------------------------------------------------------

func seedEnrollments(students []models.Student, activities []models.Activity) error {
	if len(students) == 0 || len(activities) == 0 {
		return nil
	}

	// Roll number -> the activities that student is enrolled in, by name.
	plan := map[string][]struct {
		activity string
		status   string
	}{
		"IT2K24011": {{"National Hackathon 2026", "Completed"}, {"Student Leadership Workshop", "Enrolled"}, {"Blood Donation Camp", "Completed"}},
		"IT2K24012": {{"National Science Olympiad", "Completed"}, {"Cultural Fest - Rangmanch", "Enrolled"}},
		"IT2K24013": {{"Inter-College Athletics Meet", "Enrolled"}, {"National Hackathon 2026", "Enrolled"}},
		"IT2K24014": {{"Inter College Debate Championship", "Completed"}},
		"IT2K25001": {{"National Hackathon 2026", "Enrolled"}, {"Cultural Fest - Rangmanch", "Enrolled"}},
		"IT2K25002": {{"Blood Donation Camp", "Completed"}},
		// IT2K24015 and IT2K25003 stay unenrolled on purpose: the admin views need
		// students with no activity to show as "Inactive".
	}

	byName := make(map[string]models.Activity, len(activities))
	for _, activity := range activities {
		byName[activity.Name] = activity
	}

	for rollNo, entries := range plan {
		for _, entry := range entries {
			activity, ok := byName[entry.activity]
			if !ok {
				continue
			}

			var existing models.Enrollment
			if err := DB.Where(models.Enrollment{StudentRollNo: rollNo, ActivityID: activity.ID}).
				Assign(models.Enrollment{Status: entry.status}).
				FirstOrCreate(&existing).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Certificates
// ---------------------------------------------------------------------------

func currentAcademicYear(now time.Time) (time.Time, time.Time) {
	startYear := now.Year()
	if now.Month() < time.July {
		startYear--
	}

	start := time.Date(startYear, time.July, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(startYear+1, time.June, 30, 23, 59, 59, 0, time.UTC)

	return start, end
}

func seedCertificates() error {
	now := time.Now().UTC()
	yearStart, _ := currentAcademicYear(now)

	type certSeed struct {
		rollNo            string
		activityName      string
		category          string
		organizer         string
		eventLevel        string
		participationType string
		status            string
		rejectionReason   string
	}

	certs := []certSeed{
		{"IT2K24011", "National Hackathon 2026", "TECHNICAL", "IIPS DAVV", "National", "Winner", "Approved", ""},
		{"IT2K24011", "Blood Donation Camp", "SOCIAL SERVICE", "District Hospital Indore", "State", "Volunteer", "Approved", ""},
		{"IT2K24011", "Python Programming Course", "TECHNICAL", "Coursera", "International", "Participant", "Pending", ""},
		{"IT2K24012", "National Science Olympiad", "RESEARCH", "Govt. of MP", "National", "Runner Up", "Approved", ""},
		{"IT2K24012", "Cultural Fest - Rangmanch", "CULTURAL", "IIPS DAVV", "State", "Participant", "Pending", ""},
		{"IT2K24013", "Inter-College Athletics Meet", "SPORTS", "DAVV", "State", "Winner", "Approved", ""},
		{"IT2K24013", "Web Development Bootcamp", "TECHNICAL", "Udemy", "International", "Participant", "Rejected", "Certificate image is blurred. Please upload a clearer copy."},
		{"IT2K24014", "Inter College Debate Championship", "PUBLIC SPEAKING", "IIPS DAVV", "National", "Winner", "Approved", ""},
		{"IT2K25001", "Student Leadership Workshop", "LEADERSHIP", "IIPS DAVV", "State", "Coordinator", "Pending", ""},
		{"IT2K25002", "Blood Donation Camp", "SOCIAL SERVICE", "District Hospital Indore", "State", "Volunteer", "Approved", ""},
	}

	span := now.Sub(yearStart)
	step := span / time.Duration(len(certs)+1)

	for i, cert := range certs {
		activityDate := yearStart.Add(step * time.Duration(i+1))
		issueDate := activityDate.AddDate(0, 0, 3)
		if issueDate.After(now) {
			issueDate = now
		}

		certNumber := fmt.Sprintf("CERT-%s-%03d", cert.rollNo, i+1)
		fileName := fmt.Sprintf("%s_%s.pdf", cert.rollNo, certNumber)

		if err := writeSeedCertificateFile(fileName, cert.activityName); err != nil {
			return err
		}

		var existing models.Certificate
		if err := DB.Where(models.Certificate{CertNumber: certNumber}).
			Assign(models.Certificate{
				StudentRollNo:     cert.rollNo,
				ActivityName:      cert.activityName,
				ActivityCategory:  cert.category,
				ActivityDate:      activityDate,
				OrganizerName:     cert.organizer,
				EventLevel:        cert.eventLevel,
				IssueDate:         &issueDate,
				ParticipationType: cert.participationType,
				Description:       fmt.Sprintf("%s - %s", cert.activityName, cert.participationType),
				FileName:          fileName,
				FilePath:          filepath.Join(certificateSeedDir, fileName),
				Credits:           CreditsForCertificate(cert.participationType, cert.eventLevel),
				Status:            cert.status,
				RejectionReason:   cert.rejectionReason,
			}).
			FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}

	return nil
}

func writeSeedCertificateFile(fileName, activityName string) error {
	if err := os.MkdirAll(certificateSeedDir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(certificateSeedDir, fileName)
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	text := fmt.Sprintf("iSPARC demo certificate - %s", activityName)
	pdf := fmt.Sprintf(`%%PDF-1.4
1 0 obj << /Type /Catalog /Pages 2 0 R >> endobj
2 0 obj << /Type /Pages /Kids [3 0 R] /Count 1 >> endobj
3 0 obj << /Type /Page /Parent 2 0 R /MediaBox [0 0 300 120] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >> endobj
4 0 obj << /Length %d >> stream
BT /F1 10 Tf 20 60 Td (%s) Tj ET
endstream endobj
5 0 obj << /Type /Font /Subtype /Type1 /BaseFont /Helvetica >> endobj
trailer << /Root 1 0 R >>
%%%%EOF
`, len(text)+30, text)

	return os.WriteFile(path, []byte(pdf), 0o600)
}
