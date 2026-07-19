package controllers

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
	"gorm.io/gorm"
)

var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegexp.MatchString(email)
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char) || char == ' ':
			hasSpecial = true
		default:
			if (char >= 32 && char <= 47) || (char >= 58 && char <= 64) || (char >= 91 && char <= 96) || (char >= 123 && char <= 126) {
				hasSpecial = true
			}
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func getAdminStats(admin *models.Admin) fiber.Map {
	var assignedStudents int64
	var verifiedCertificates int64
	var pendingReviews int64
	var supervisedActivities int64

	// 1. Assigned Students Count
	studentQuery := config.DB.Model(&models.Student{})
	if admin.Role != "superadmin" && admin.AssignedBatch != "" {
		studentQuery = studentQuery.Where("roll_no LIKE ?", admin.AssignedBatch+"%")
	}
	studentQuery.Count(&assignedStudents)

	// 2 & 3. Certificates Stats
	certQuery := config.DB.Model(&models.Certificate{})
	if admin.Role != "superadmin" && admin.AssignedBatch != "" {
		certQuery = certQuery.Where("student_roll_no LIKE ?", admin.AssignedBatch+"%")
	}
	certQuery.Where("status = ?", "Approved").Count(&verifiedCertificates)

	certQueryPending := config.DB.Model(&models.Certificate{})
	if admin.Role != "superadmin" && admin.AssignedBatch != "" {
		certQueryPending = certQueryPending.Where("student_roll_no LIKE ?", admin.AssignedBatch+"%")
	}
	certQueryPending.Where("status = ?", "Pending").Count(&pendingReviews)

	// 4. Activities Supervised
	config.DB.Model(&models.Activity{}).Where("coordinator_id = ?", admin.AdminID).Count(&supervisedActivities)

	return fiber.Map{
		"assigned_students":     assignedStudents,
		"verified_certificates": verifiedCertificates,
		"pending_reviews":       pendingReviews,
		"supervised_activities": supervisedActivities,
	}
}

func errJSON(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"error": msg})
}

func AdminLogin(c *fiber.Ctx) error {
	var input models.AdminLoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}
	if input.AdminID == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Admin ID and Password are required"})
	}
	var admin models.Admin
	if err := config.DB.Where("admin_id = ?", input.AdminID).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	if !utils.CheckPasswordHash(input.Password, admin.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	accessToken, err := utils.GenerateAccessToken(admin.AdminID, admin.Email, admin.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}
	return c.JSON(fiber.Map{
		"message":              "Admin logged in successfully",
		"access_token":         accessToken,
		"must_change_password": admin.MustChangePassword,
		"admin": fiber.Map{
			"admin_id": admin.AdminID,
			"name":     admin.Name,
			"role":     admin.Role,
		},
	})
}

// 1. POST /admin/change-password
func AdminChangePassword(c *fiber.Ctx) error {
	var input models.ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return errJSON(c, fiber.StatusBadRequest, "Cannot parse request body")
	}
	if input.CurrentPassword == "" || input.NewPassword == "" || input.ConfirmPassword == "" {
		return errJSON(c, fiber.StatusBadRequest, "All fields are required")
	}
	if input.NewPassword != input.ConfirmPassword {
		return errJSON(c, fiber.StatusBadRequest, "Passwords do not match")
	}

	if !isStrongPassword(input.NewPassword) {
		return errJSON(c, fiber.StatusBadRequest, "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.")
	}

	if len(input.NewPassword) > 72 {
		return errJSON(c, fiber.StatusBadRequest, "Password cannot exceed 72 characters")
	}

	adminID, ok := c.Locals("roll_no").(string)
	if !ok || adminID == "" {
		return errJSON(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var admin models.Admin
	if err := config.DB.Where("admin_id = ?", adminID).First(&admin).Error; err != nil {
		return errJSON(c, fiber.StatusNotFound, "Admin not found")
	}

	if !utils.CheckPasswordHash(input.CurrentPassword, admin.Password) {
		return errJSON(c, fiber.StatusUnauthorized, "Current password is incorrect")
	}

	newHash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return errJSON(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	admin.Password = newHash
	admin.MustChangePassword = false

	if err := config.DB.Save(&admin).Error; err != nil {
		return errJSON(c, fiber.StatusInternalServerError, "Failed to update password")
	}

	return c.JSON(fiber.Map{"message": "Password changed successfully"})
}

func getAuthenticatedAdmin(c *fiber.Ctx) (*models.Admin, error) {
	adminID, ok := c.Locals("roll_no").(string)
	if !ok || adminID == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	var admin models.Admin
	if err := config.DB.Where("admin_id = ?", adminID).First(&admin).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	return &admin, nil
}

func scopeToAssignedBatch(query *gorm.DB, admin *models.Admin) (*gorm.DB, bool) {
	if admin.Role != "admin" {
		return query, true
	}
	if admin.AssignedBatch == "" {
		return query, false
	}

	return query.Where("roll_no LIKE ?", admin.AssignedBatch+"%"), true
}

func applyStudentStats(student *models.Student) {
	if len(student.RollNo) >= 6 {
		student.Batch = student.RollNo[:6]
	} else {
		student.Batch = "Unknown"
	}

	credits := 0
	pending := 0
	for _, cert := range student.Certificates {
		if cert.Status == "Approved" {
			credits += cert.Credits
		}
		if cert.Status == "Pending" {
			pending++
		}
	}

	student.CreditsEarned = credits
	student.PendingCertificates = pending
	student.ActivityCount = len(student.Enrollments)

	switch {
	case len(student.Enrollments) == 0 && len(student.Certificates) == 0:
		student.EngagementStatus = "Inactive"
	case pending > 0:
		student.EngagementStatus = "Pending Review"
	default:
		student.EngagementStatus = "Active"
	}
}

// 2. GET /api/admin/students -> View assigned students
func GetAllStudents(c *fiber.Ctx) error {
	currentUser, err := getAuthenticatedAdmin(c)
	if err != nil {
		return err
	}

	dbQuery, scoped := scopeToAssignedBatch(
		config.DB.Preload("Certificates").Preload("Enrollments"),
		currentUser,
	)
	if !scoped {
		return c.JSON(fiber.Map{"students": []models.Student{}})
	}

	var students []models.Student
	if err := dbQuery.Find(&students).Error; err != nil {
		return errJSON(c, fiber.StatusInternalServerError, "Failed to retrieve students")
	}

	for i := range students {
		applyStudentStats(&students[i])
	}

	return c.JSON(fiber.Map{"students": students})
}

// 3. GET /api/admin/students/:roll -> One student's detail
func GetStudentDetail(c *fiber.Ctx) error {
	roll := c.Params("roll")

	currentUser, err := getAuthenticatedAdmin(c)
	if err != nil {
		return err
	}

	query, scoped := scopeToAssignedBatch(
		config.DB.Preload("Enrollments.Activity").
			Preload("Certificates").
			Where("roll_no = ?", roll),
		currentUser,
	)
	if !scoped {
		return errJSON(c, fiber.StatusNotFound, "Student not found")
	}

	var student models.Student
	if err := query.First(&student).Error; err != nil {
		return errJSON(c, fiber.StatusNotFound, "Student not found")
	}
	applyStudentStats(&student)

	return c.JSON(fiber.Map{"student": student})
}

// GET /api/admin/profile -> Retrieve authenticated admin details
func GetAdminProfile(c *fiber.Ctx) error {
	admin, err := getAuthenticatedAdmin(c)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"admin": fiber.Map{
			"admin_id":       admin.AdminID,
			"name":           admin.Name,
			"email":          admin.Email,
			"role":           admin.Role,
			"assigned_batch": admin.AssignedBatch,
			"created_at":     admin.CreatedAt,
			"updated_at":     admin.UpdatedAt,
		},
		"stats": getAdminStats(admin),
	})
}

// PUT /api/admin/profile -> Update authenticated admin details
func UpdateAdminProfile(c *fiber.Ctx) error {
	admin, err := getAuthenticatedAdmin(c)
	if err != nil {
		return err
	}

	type ProfileUpdate struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var input ProfileUpdate
	if err := c.BodyParser(&input); err != nil {
		return errJSON(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	// Trim inputs
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if input.Name == "" {
		return errJSON(c, fiber.StatusBadRequest, "Name is required")
	}

	if len(input.Name) > 100 {
		return errJSON(c, fiber.StatusBadRequest, "Name cannot exceed 100 characters")
	}

	if input.Email == "" {
		return errJSON(c, fiber.StatusBadRequest, "Email is required")
	}

	if len(input.Email) > 100 {
		return errJSON(c, fiber.StatusBadRequest, "Email cannot exceed 100 characters")
	}

	if !isValidEmail(input.Email) {
		return errJSON(c, fiber.StatusBadRequest, "Invalid email format")
	}

	// Normalize email to lowercase
	input.Email = strings.ToLower(input.Email)

	// Check if email already exists for another admin case-insensitively
	var existing models.Admin
	if err := config.DB.Where("LOWER(email) = ? AND admin_id <> ?", input.Email, admin.AdminID).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists.",
		})
	}

	admin.Name = input.Name
	admin.Email = input.Email

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(admin).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Activity{}).
			Where("coordinator_id = ?", admin.AdminID).
			Update("coordinator", admin.Name).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errJSON(c, fiber.StatusInternalServerError, "Failed to update profile")
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"admin": fiber.Map{
			"admin_id":       admin.AdminID,
			"name":           admin.Name,
			"email":          admin.Email,
			"role":           admin.Role,
			"assigned_batch": admin.AssignedBatch,
			"created_at":     admin.CreatedAt,
			"updated_at":     admin.UpdatedAt,
		},
		"stats": getAdminStats(admin),
	})
}
