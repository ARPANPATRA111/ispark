package config

import (
	"log"
	"time"

	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
)

func SeedDefaultStudents() {
	hashedPassword, _ := utils.HashPassword("student123")

	student := models.Student{
		RollNo:       "IT2K24011",
		Name:         "Rahul Sharma",
		CourseName:   "Computer Science",
		Semester:     6,
		ContactNo:    "9876543210",
		EmailID:      "rahul.sharma@iips.edu",
		EnrollmentNo: "EN2026001",
		Password:     hashedPassword,
		IsVerified:   true,
	}

	// Seed student
	DB.Where("roll_no = ?", student.RollNo).FirstOrCreate(&student)

	// ---------------- Activities ----------------

	activities := []models.Activity{
		{
			Name:         "Leadership Workshop",
			Category:     "Leadership",
			Description:  "Leadership development workshop.",
			Credits:      10,
			Mode:         "Offline",
			RegDeadline:  time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC),
			ActivityDate: time.Date(2026, 6, 23, 10, 0, 0, 0, time.UTC),
			Venue:        "Seminar Hall",
			Coordinator:  "Dr. Mehta",
			Status:       "Open",
		},
		{
			Name:         "Beach Cleanup",
			Category:     "Social Service",
			Description:  "Community service event.",
			Credits:      5,
			Mode:         "Offline",
			RegDeadline:  time.Date(2026, 6, 24, 0, 0, 0, 0, time.UTC),
			ActivityDate: time.Date(2026, 6, 25, 9, 0, 0, 0, time.UTC),
			Venue:        "City Beach",
			Coordinator:  "NSS Cell",
			Status:       "Open",
		},
	}

	for i := range activities {
		DB.Where("name = ?", activities[i].Name).FirstOrCreate(&activities[i])
		DB.Where("student_roll_no = ? AND activity_id = ?", student.RollNo, activities[i].ID).
			FirstOrCreate(&models.Enrollment{
				StudentRollNo: student.RollNo,
				ActivityID:    activities[i].ID,
				Status:        "Enrolled",
			})
	}

	// ---------------- Certificates ----------------

	pythonIssueDate := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
	debateIssueDate := time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)

	certificates := []models.Certificate{
		{
			StudentRollNo:     student.RollNo,
			ActivityName:      "Python Programming Course",
			ActivityCategory:  "Technical",
			ActivityDate:      time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC),
			OrganizerName:     "Coursera",
			EventLevel:        "National",
			CertNumber:        "CERT-PY-001",
			IssueDate:         &pythonIssueDate,
			ParticipationType: "Participant",
			Description:       "Completed Python Programming Course.",
			FileName:          "python_certificate.pdf",
			FilePath:          "/certificates/python_certificate.pdf",
			Credits:           CreditsForCertificate("Participant", "National"),
			Status:            "Approved",
		},
		{
			StudentRollNo:     student.RollNo,
			ActivityName:      "State Debate Championship",
			ActivityCategory:  "Public Speaking",
			ActivityDate:      time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC),
			OrganizerName:     "Govt. of MP",
			EventLevel:        "State",
			CertNumber:        "CERT-DB-002",
			IssueDate:         &debateIssueDate,
			ParticipationType: "Runner Up",
			Description:       "Secured Runner Up position.",
			FileName:          "debate_certificate.pdf",
			FilePath:          "/certificates/debate_certificate.pdf",
			Credits:           CreditsForCertificate("Runner Up", "State"),
			Status:            "Pending",
		},
	}

	for _, cert := range certificates {
		DB.Where("cert_number = ?", cert.CertNumber).
			FirstOrCreate(&cert)
	}

	log.Println("Seeded default student, activities and certificates.")
}
