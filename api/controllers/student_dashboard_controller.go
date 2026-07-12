package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
)

const certificateUploadDir = "./uploads/certificates"

// keeping it like this soo that adding more file types in future is easy.
// If we want to allow more file types, we can just add them here.
var allowedCertificateTypes = map[string]string{
	".pdf":  "application/pdf",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
}

// max certificate file size is 5 MB but i highly doublt that we should keep it 5 MB.
const maxCertificateSize = 5 * 1024 * 1024

// GetCertificates returns student's uploaded certificates
func GetCertificates(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	var certificates []models.Certificate
	if err := config.DB.Where("student_roll_no = ?", rollNo).Order("created_at desc").Find(&certificates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch certificates",
		})
	}

	return c.JSON(certificates)
}

// UploadCertificate uploads a new certificate (handles file upload + details)
func UploadCertificate(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	activityName := c.FormValue("activity_name")
	activityCategory := c.FormValue("activity_category")
	activityDateStr := c.FormValue("activity_date")
	organizerName := c.FormValue("organizer_name")
	eventLevel := c.FormValue("event_level")
	certNumber := c.FormValue("cert_number")
	issueDateStr := c.FormValue("issue_date")
	participationType := c.FormValue("participation_type")
	description := c.FormValue("description")

	if activityName == "" || participationType == "" || activityDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Required fields are missing",
		})
	}

	// Parse dates
	activityDate, err := time.Parse("2006-01-02", activityDateStr)
	if err != nil || activityDate.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid activity date format. Expected YYYY-MM-DD",
		})
	}
	var issueDate *time.Time
	if issueDateStr != "" {
		parsedDate, parseErr := time.Parse("2006-01-02", issueDateStr)
		if parseErr != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid issue date format. Expected YYYY-MM-DD",
			})
		}
		issueDate = &parsedDate
	}

	credits := config.CreditsForCertificate(participationType, eventLevel)

	// Handle File Upload
	file, err := c.FormFile("certificate_file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Certificate file is required",
		})
	}

	// Validate file size (max 5 MB)
	if file.Size > maxCertificateSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size exceeds 5 MB limit",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	expectedType, extAllowed := allowedCertificateTypes[ext]
	if !extAllowed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported file format. Please upload PDF, PNG, or JPG",
		})
	}

	detectedType, err := detectFileType(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read uploaded file",
		})
	}
	if detectedType != expectedType {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File contents do not match its extension. Please upload a valid PDF, PNG, or JPG",
		})
	}

	// Ensure directory exists
	if err := os.MkdirAll(certificateUploadDir, 0o755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	fileName := fmt.Sprintf("%s_%d_%s", rollNo, time.Now().UnixNano(), filepath.Base(file.Filename))
	filePath := filepath.Join(certificateUploadDir, fileName)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	cert := models.Certificate{
		StudentRollNo:     rollNo,
		ActivityName:      activityName,
		ActivityCategory:  activityCategory,
		ActivityDate:      activityDate,
		OrganizerName:     organizerName,
		EventLevel:        eventLevel,
		CertNumber:        certNumber,
		IssueDate:         issueDate,
		ParticipationType: participationType,
		Description:       description,
		FileName:          fileName,
		FilePath:          filePath,
		Credits:           credits,
		Status:            "Pending",
	}

	if err := config.DB.Create(&cert).Error; err != nil {
		_ = os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save certificate record",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Certificate uploaded successfully",
		"certificate": cert,
	})
}

func detectFileType(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	head := make([]byte, 512)
	n, err := f.Read(head)
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.SplitN(http.DetectContentType(head[:n]), ";", 2)[0], nil
}

// now a student can only download their own certificate files.
func DownloadCertificate(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	var cert models.Certificate
	if err := config.DB.Where("id = ? AND student_roll_no = ?", c.Params("id"), rollNo).First(&cert).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Certificate not found",
		})
	}

	filePath := filepath.Join(certificateUploadDir, filepath.Base(cert.FileName))
	if _, err := os.Stat(filePath); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Certificate file is no longer available",
		})
	}

	return c.Download(filePath, cert.FileName)
}

// getAcademicYearDateRange returns start and end date of the academic year
func getAcademicYearDateRange(yearStr string) (time.Time, time.Time, error) {
	if len(yearStr) != 7 || yearStr[4] != '-' {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid year format, expected YYYY-YY")
	}

	startYear, err := strconv.Atoi(yearStr[0:4])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start year: %v", err)
	}

	endYearPart, err := strconv.Atoi(yearStr[5:7])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end year suffix: %v", err)
	}

	expectedEndYearPart := (startYear + 1) % 100
	if endYearPart != expectedEndYearPart {
		return time.Time{}, time.Time{}, fmt.Errorf("academic year mismatch: suffix must be %02d", expectedEndYearPart)
	}

	endYear := startYear + 1
	startDate := time.Date(startYear, time.July, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(endYear, time.June, 30, 23, 59, 59, 999999999, time.UTC)

	return startDate, endDate, nil
}

// getCurrentAcademicYearRange returns the current academic year dates
func getCurrentAcademicYearRange() (time.Time, time.Time) {
	now := time.Now()
	var startYear int
	if now.Month() >= time.July {
		startYear = now.Year()
	} else {
		startYear = now.Year() - 1
	}
	startDate := time.Date(startYear, time.July, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(startYear+1, time.June, 30, 23, 59, 59, 999999999, time.UTC)
	return startDate, endDate
}

// GetLeaderboard returns the leaderboard sorted by total credits for a given academic year
func GetLeaderboard(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)
	year := c.Query("year")

	var startDate, endDate time.Time
	var err error
	if year != "" {
		startDate, endDate, err = getAcademicYearDateRange(year)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid year format. Expected YYYY-YY (e.g., 2025-26)",
			})
		}
	} else {
		startDate, endDate = getCurrentAcademicYearRange()
	}

	type LeaderboardEntry struct {
		RollNo     string `json:"roll_no"`
		Name       string `json:"name"`
		CourseName string `json:"course_name"`
		Semester   int    `json:"semester"`
		Points     int    `json:"points"`
		IsSelf     bool   `json:"is_self"`
	}

	var entries []LeaderboardEntry

	err = config.DB.Raw(`
		SELECT
			s.roll_no,
			s.name,
			s.course_name,
			s.semester,
			COALESCE(SUM(c.credits), 0) as points
		FROM students s
		LEFT JOIN certificates c ON c.student_roll_no = s.roll_no
			AND c.status = 'Approved'
			AND c.activity_date >= ?
			AND c.activity_date <= ?
		GROUP BY s.roll_no, s.name, s.course_name, s.semester
		ORDER BY points DESC, s.name ASC
	`, startDate, endDate).Scan(&entries).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch leaderboard",
		})
	}

	// Mark the authenticated student
	for i := range entries {
		if entries[i].RollNo == rollNo {
			entries[i].IsSelf = true
		}
	}

	return c.JSON(entries)
}

// GetCategoryChampions returns the top student per activity category for a given academic year
func GetCategoryChampions(c *fiber.Ctx) error {
	year := c.Query("year")

	var startDate, endDate time.Time
	var err error
	if year != "" {
		startDate, endDate, err = getAcademicYearDateRange(year)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid year format. Expected YYYY-YY (e.g., 2025-26)",
			})
		}
	} else {
		startDate, endDate = getCurrentAcademicYearRange()
	}

	type ChampionResult struct {
		Track   string `json:"track"`
		RollNo  string `json:"roll_no"`
		Name    string `json:"name"`
		Credits int    `json:"credits"`
	}

	var champions []ChampionResult

	err = config.DB.Raw(`
		WITH CategoryCredits AS (
			SELECT
				UPPER(c.activity_category) AS track,
				s.roll_no,
				s.name,
				SUM(c.credits) AS total_credits
			FROM students s
			JOIN certificates c ON c.student_roll_no = s.roll_no
			WHERE c.status = 'Approved'
			  AND c.activity_date >= ?
			  AND c.activity_date <= ?
			GROUP BY UPPER(c.activity_category), s.roll_no, s.name
		),
		RankedCategoryCredits AS (
			SELECT
				track,
				roll_no,
				name,
				total_credits,
				ROW_NUMBER() OVER (PARTITION BY track ORDER BY total_credits DESC, name ASC) as rn
			FROM CategoryCredits
		)
		SELECT track, roll_no, name, total_credits as credits
		FROM RankedCategoryCredits
		WHERE rn = 1
	`, startDate, endDate).Scan(&champions).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch category champions",
		})
	}

	return c.JSON(champions)
}

// GetActivities returns a list of activities
func GetActivities(c *fiber.Ctx) error {
	category := c.Query("category")
	status := c.Query("status")
	search := c.Query("search")

	query := config.DB.Model(&models.Activity{})
	if category != "" {
		query = query.Where("UPPER(category) = UPPER(?)", category)
	}
	if status != "" {
		query = query.Where("UPPER(status) = UPPER(?)", status)
	}
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR coordinator ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var activities []models.Activity
	if err := query.Order("activity_date asc").Find(&activities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch activities",
		})
	}

	return c.JSON(activities)
}

// UpdateProfile updates student's editable profile info
func UpdateProfile(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	type ProfileUpdate struct {
		EmailID   string `json:"email_id"`
		ContactNo string `json:"contact_no"`
		DOB       string `json:"dob"`
		Gender    string `json:"gender"`
	}

	var input ProfileUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var student models.Student
	if err := config.DB.Where("roll_no = ?", rollNo).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	// Update fields if provided
	if input.EmailID != "" {
		student.EmailID = input.EmailID
	}
	if input.ContactNo != "" {
		student.ContactNo = input.ContactNo
	}
	if input.DOB != "" {
		student.DOB = input.DOB
	}
	if input.Gender != "" {
		student.Gender = input.Gender
	}

	if err := config.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"student": fiber.Map{
			"roll_no":       student.RollNo,
			"name":          student.Name,
			"email_id":      student.EmailID,
			"course_name":   student.CourseName,
			"semester":      student.Semester,
			"contact_no":    student.ContactNo,
			"dob":           student.DOB,
			"gender":        student.Gender,
			"enrollment_no": student.EnrollmentNo,
			"is_verified":   student.IsVerified,
		},
	})
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangePassword changes the authenticated student's password
func ChangePassword(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	var input ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON input",
		})
	}

	if len(input.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password must be at least 6 characters long",
		})
	}

	var student models.Student
	if err := config.DB.Where("roll_no = ?", rollNo).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	// Verify current password
	if !utils.CheckPasswordHash(input.CurrentPassword, student.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incorrect current password",
		})
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	student.Password = hashedPassword
	if err := config.DB.Save(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

// EnrollActivity enrolls a student in an activity
func EnrollActivity(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)
	activityID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid activity ID",
		})
	}

	// Check if activity exists
	var activity models.Activity
	if err := config.DB.First(&activity, activityID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Activity not found",
		})
	}

	// Check if already enrolled
	var existing models.Enrollment
	if err := config.DB.Where("student_roll_no = ? AND activity_id = ?", rollNo, activityID).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "You are already enrolled in this activity",
		})
	}

	enrollment := models.Enrollment{
		StudentRollNo: rollNo,
		ActivityID:    uint(activityID),
		Status:        "Enrolled",
	}

	if err := config.DB.Create(&enrollment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to enroll in activity",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Enrolled successfully",
		"enrollment": enrollment,
	})
}

// GetEnrollments returns enrollments for the student
func GetEnrollments(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	var enrollments []models.Enrollment
	if err := config.DB.Preload("Activity").Where("student_roll_no = ?", rollNo).Order("created_at desc").Find(&enrollments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch enrollments",
		})
	}

	return c.JSON(enrollments)
}

// GetDashboardStats returns stats for the student dashboard home page
func GetDashboardStats(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	// 1. Total activities participated (count of enrollments)
	var activitiesCount int64
	if err := config.DB.Model(&models.Enrollment{}).Where("student_roll_no = ? AND status IN ('Enrolled', 'Completed', 'Registered', 'Verified', 'Pending Verification')", rollNo).Count(&activitiesCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate activities count",
		})
	}

	// 2. Total certificates uploaded and pending/approved/rejected
	var certificatesCount int64
	if err := config.DB.Model(&models.Certificate{}).Where("student_roll_no = ?", rollNo).Count(&certificatesCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate certificates count",
		})
	}

	var pendingCertificatesCount int64
	if err := config.DB.Model(&models.Certificate{}).Where("student_roll_no = ? AND status = 'Pending'", rollNo).Count(&pendingCertificatesCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate pending certificates count",
		})
	}

	var approvedCertificatesCount int64
	if err := config.DB.Model(&models.Certificate{}).Where("student_roll_no = ? AND status = 'Approved'", rollNo).Count(&approvedCertificatesCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate approved certificates count",
		})
	}

	var rejectedCertificatesCount int64
	if err := config.DB.Model(&models.Certificate{}).Where("student_roll_no = ? AND status = 'Rejected'", rollNo).Count(&rejectedCertificatesCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate rejected certificates count",
		})
	}

	// 3. Credits earned from approved certificates
	type SumResult struct {
		Total int
	}
	var sumResult SumResult
	if err := config.DB.Raw("SELECT COALESCE(SUM(credits), 0) as total FROM certificates WHERE student_roll_no = ? AND status = 'Approved'", rollNo).Scan(&sumResult).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate total credits",
		})
	}

	// 4. Current Rank on the Leaderboard
	type StudentRank struct {
		RollNo       string
		TotalCredits int
	}
	var ranks []StudentRank
	if err := config.DB.Raw(`
		SELECT s.roll_no, COALESCE(SUM(c.credits), 0) as total_credits
		FROM students s
		LEFT JOIN certificates c ON c.student_roll_no = s.roll_no AND c.status = 'Approved'
		GROUP BY s.roll_no
		ORDER BY total_credits DESC, s.roll_no ASC
	`).Scan(&ranks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate student ranks",
		})
	}

	rank := len(ranks)
	totalStudents := len(ranks)
	for i, r := range ranks {
		if r.RollNo == rollNo {
			rank = i + 1
			break
		}
	}

	// 5. Recent extracurricular activities list (approved/pending certificates)
	var recentActivities []models.Certificate
	if err := config.DB.Where("student_roll_no = ?", rollNo).Order("created_at desc").Limit(5).Find(&recentActivities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch recent activities",
		})
	}
	return c.JSON(fiber.Map{
		"activities_participated": activitiesCount,
		"certificates_uploaded":   certificatesCount,
		"pending_certificates":    pendingCertificatesCount,
		"approved_certificates":   approvedCertificatesCount,
		"rejected_certificates":   rejectedCertificatesCount,
		"credits_earned":          sumResult.Total,
		"current_rank":            rank,
		"total_students":          totalStudents,
		"recent_activities":       recentActivities,
	})
}
