package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
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
	// Find the admin using your existing Admin model
	var admin models.Admin
	if err := config.DB.Where("admin_id = ?", input.AdminID).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	// Verify Password
	if !utils.CheckPasswordHash(input.Password, admin.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	// Generate Access Token
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

// 1. POST /admin/change-password -> ChangePassword lets a logged-in admin set a new password (used for both voluntary changes and the forced first-login reset flow)
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
