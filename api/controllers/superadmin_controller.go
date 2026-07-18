package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
	"gorm.io/gorm"
)

// allowedSettingStatuses is the set of statuses a system setting may hold.
var allowedSettingStatuses = map[string]bool{
	"Active":   true,
	"Enabled":  true,
	"Disabled": true,
}

// platformUser is the flattened shape the super admin user registry renders.
// Students and admins are different tables, so they are normalised here.
type platformUser struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Role   string `json:"role"`
	Dept   string `json:"dept"`
	Status string `json:"status"`
}

// createPlatformUserInput is what the super admin "Create User" form submits.
type createPlatformUserInput struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	ID       string `json:"id"`
	Email    string `json:"email"`
	Dept     string `json:"dept"`
	Semester int    `json:"semester"`
}

// generateTemporaryPassword returns a random password for a newly created
// account. The account holder is expected to change it on first login.
func generateTemporaryPassword() (string, error) {
	buf := make([]byte, 9)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return "iSPARC-" + base64.RawURLEncoding.EncodeToString(buf), nil
}

// GetPlatformStats returns platform-wide counts for the super admin dashboard.
func GetPlatformStats(c *fiber.Ctx) error {
	var (
		totalStudents       int64
		totalAdmins         int64
		totalActivities     int64
		totalCertificates   int64
		pendingCertificates int64
		activeTracks        int64
	)

	counts := []struct {
		model any
		where string
		args  []any
		into  *int64
	}{
		{model: &models.Student{}, into: &totalStudents},
		{model: &models.Admin{}, into: &totalAdmins},
		{model: &models.Activity{}, into: &totalActivities},
		{model: &models.Certificate{}, into: &totalCertificates},
		{model: &models.Certificate{}, where: "status = ?", args: []any{"Pending"}, into: &pendingCertificates},
	}

	for _, count := range counts {
		query := config.DB.Model(count.model)
		if count.where != "" {
			query = query.Where(count.where, count.args...)
		}
		if err := query.Count(count.into).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load platform statistics",
			})
		}
	}

	// A track is an activity category, so the number of distinct categories in
	// use is the number of active tracks.
	if err := config.DB.Model(&models.Activity{}).
		Distinct("category").
		Count(&activeTracks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load platform statistics",
		})
	}

	return c.JSON(fiber.Map{
		"total_students":       totalStudents,
		"total_admins":         totalAdmins,
		"total_users":          totalStudents + totalAdmins,
		"total_activities":     totalActivities,
		"total_certificates":   totalCertificates,
		"pending_certificates": pendingCertificates,
		"active_tracks":        activeTracks,
	})
}

// GetPlatformUsers returns every student and admin account for the super admin
// user registry.
func GetPlatformUsers(c *fiber.Ctx) error {
	var students []models.Student
	if err := config.DB.
		Select("roll_no", "name", "course_name", "is_verified").
		Order("created_at desc").
		Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load users",
		})
	}

	var admins []models.Admin
	if err := config.DB.
		Select("admin_id", "name", "role", "assigned_batch").
		Order("created_at desc").
		Find(&admins).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load users",
		})
	}

	users := make([]platformUser, 0, len(students)+len(admins))

	for _, admin := range admins {
		dept := admin.AssignedBatch
		if dept == "" {
			dept = "All Batches"
		}

		role := "Admin"
		if admin.Role == "superadmin" {
			role = "Super Admin"
		}

		users = append(users, platformUser{
			Name:   admin.Name,
			ID:     admin.AdminID,
			Role:   role,
			Dept:   dept,
			Status: "Active",
		})
	}

	for _, student := range students {
		status := "Pending"
		if student.IsVerified {
			status = "Active"
		}

		users = append(users, platformUser{
			Name:   student.Name,
			ID:     student.RollNo,
			Role:   "Student",
			Dept:   student.CourseName,
			Status: status,
		})
	}

	return c.JSON(fiber.Map{"users": users})
}

// CreatePlatformUser registers a new student or admin account from the super
// admin user registry. The account is created with a generated temporary
// password, which is returned once so the super admin can pass it on.
func CreatePlatformUser(c *fiber.Ctx) error {
	var input createPlatformUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	input.Name = strings.TrimSpace(input.Name)
	input.ID = strings.TrimSpace(input.ID)
	input.Email = strings.TrimSpace(input.Email)
	input.Dept = strings.TrimSpace(input.Dept)

	if input.Name == "" || input.ID == "" || input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, ID and email are required",
		})
	}

	if input.Role != "Student" && input.Role != "Admin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Role must be either Student or Admin",
		})
	}

	tempPassword, err := generateTemporaryPassword()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create account",
		})
	}

	hashed, err := utils.HashPassword(tempPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create account",
		})
	}

	var created platformUser

	if input.Role == "Student" {
		semester := input.Semester
		if semester <= 0 {
			semester = 1
		}

		student := models.Student{
			RollNo:       input.ID,
			Name:         input.Name,
			CourseName:   input.Dept,
			Semester:     semester,
			EmailID:      input.Email,
			EnrollmentNo: "EN-" + input.ID,
			Password:     hashed,
			// The account still has to verify its email, so it shows as Pending
			// in the registry until the student completes OTP verification.
			IsVerified: false,
		}

		if err := config.DB.Create(&student).Error; err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "A student with this roll number, email or enrollment number already exists",
			})
		}

		created = platformUser{
			Name:   student.Name,
			ID:     student.RollNo,
			Role:   "Student",
			Dept:   student.CourseName,
			Status: "Pending",
		}
	} else {
		admin := models.Admin{
			AdminID:            input.ID,
			Name:               input.Name,
			Email:              input.Email,
			Password:           hashed,
			Role:               "admin",
			AssignedBatch:      input.Dept,
			MustChangePassword: true,
		}

		if err := config.DB.Create(&admin).Error; err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "An admin with this ID or email already exists",
			})
		}

		dept := admin.AssignedBatch
		if dept == "" {
			dept = "All Batches"
		}

		created = platformUser{
			Name:   admin.Name,
			ID:     admin.AdminID,
			Role:   "Admin",
			Dept:   dept,
			Status: "Active",
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":               created,
		"temporary_password": tempPassword,
	})
}

// DeletePlatformUser removes a student or admin account. A super admin cannot
// delete their own account, and super admin accounts cannot be deleted here.
func DeletePlatformUser(c *fiber.Ctx) error {
	id := c.Params("id")

	callerID, _ := c.Locals("roll_no").(string)
	if id == callerID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "You cannot delete your own account",
		})
	}

	var admin models.Admin
	err := config.DB.Where("admin_id = ?", id).First(&admin).Error
	switch {
	case err == nil:
		if admin.Role == "superadmin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Super admin accounts cannot be deleted here",
			})
		}

		if err := config.DB.Delete(&admin).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete user",
			})
		}

		return c.JSON(fiber.Map{"message": "User deleted successfully"})

	case errors.Is(err, gorm.ErrRecordNotFound):
		// Not an admin, fall through and try students.

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	var student models.Student
	if err := config.DB.Where("roll_no = ?", id).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := config.DB.Delete(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

// ---------------------------------------------------------------------------
// System Settings
// ---------------------------------------------------------------------------

// groupSettingsByCategory loads every setting ordered by SortOrder and groups
// them under their category, which is the shape the super admin settings screen
// renders (one list per tab).
func groupSettingsByCategory() (map[string][]models.SystemSetting, error) {
	var settings []models.SystemSetting
	if err := config.DB.Order("sort_order asc").Find(&settings).Error; err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.SystemSetting)
	for _, setting := range settings {
		grouped[setting.Category] = append(grouped[setting.Category], setting)
	}

	return grouped, nil
}

// GetPlatformSettings returns every platform setting grouped by category for the
// super admin System Settings screen.
func GetPlatformSettings(c *fiber.Ctx) error {
	grouped, err := groupSettingsByCategory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load settings",
		})
	}

	return c.JSON(fiber.Map{"settings": grouped})
}

// UpdatePlatformSetting updates a single setting's value and/or status.
func UpdatePlatformSetting(c *fiber.Ctx) error {
	key := c.Params("key")

	var input models.UpdateSettingInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if input.Value == nil && input.Status == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nothing to update: provide a value or a status",
		})
	}

	if input.Status != nil && !allowedSettingStatuses[*input.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status must be one of Active, Enabled or Disabled",
		})
	}

	var setting models.SystemSetting
	if err := config.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Setting not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load setting"})
	}

	if input.Value != nil {
		setting.Value = strings.TrimSpace(*input.Value)
	}
	if input.Status != nil {
		setting.Status = *input.Status
	}

	if err := config.DB.Save(&setting).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update setting"})
	}

	return c.JSON(fiber.Map{"setting": setting})
}

// UpdatePlatformSettings applies several setting updates in one request. It is
// what the settings screen's "Save" action uses so a whole tab can be persisted
// at once. All updates run in a transaction: if any key is unknown or invalid,
// nothing is written.
func UpdatePlatformSettings(c *fiber.Ctx) error {
	var input struct {
		Settings []models.BulkSettingUpdate `json:"settings"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if len(input.Settings) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No settings provided"})
	}

	// Validate everything up front so the transaction never partially applies.
	for _, item := range input.Settings {
		if strings.TrimSpace(item.Key) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Each setting must include a key"})
		}
		if item.Value == nil && item.Status == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Setting %q has nothing to update", item.Key),
			})
		}
		if item.Status != nil && !allowedSettingStatuses[*item.Status] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid status %q for setting %q", *item.Status, item.Key),
			})
		}
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range input.Settings {
			updates := map[string]any{}
			if item.Value != nil {
				updates["value"] = strings.TrimSpace(*item.Value)
			}
			if item.Status != nil {
				updates["status"] = *item.Status
			}

			res := tx.Model(&models.SystemSetting{}).Where("key = ?", item.Key).Updates(updates)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("setting %q not found", item.Key)
			}
		}
		return nil
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	grouped, err := groupSettingsByCategory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reload settings"})
	}

	return c.JSON(fiber.Map{"settings": grouped})
}
