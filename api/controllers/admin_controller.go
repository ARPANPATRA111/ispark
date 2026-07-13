package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
	"gorm.io/gorm"
)

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
	})
}