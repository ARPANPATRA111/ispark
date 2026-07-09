package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
)

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

// GetDashboardStats returns stats for the student dashboard home page
func GetDashboardStats(c *fiber.Ctx) error {
	rollNo := c.Locals("roll_no").(string)

	// 1. Total activities participated (count of enrollments)
	var activitiesCount int64
	if err := config.DB.Model(&models.Enrollment{}).Where("student_roll_no = ? AND status IN ('Enrolled', 'Completed')", rollNo).Count(&activitiesCount).Error; err != nil {
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
