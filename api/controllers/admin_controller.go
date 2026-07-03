package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
)

func AdminLogin(c *fiber.Ctx) error {
	// A simple hardcoded response just to get your server running
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Admin login endpoint hit successfully!",
		"status":  "success",
	})
}

// 1. GET /api/admin/students -> View assigned students
func GetAllStudents(c *fiber.Ctx) error {
	var students []models.Student
	if err := config.DB.Select("roll_no, name, course_name, semester, email_id").Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve students"})
	}
	return c.JSON(fiber.Map{"students": students})
}

// 2. GET /api/admin/students/:roll -> One student's detail
func GetStudentDetail(c *fiber.Ctx) error {
	rollNo := c.Params("roll")
	var student models.Student

	if err := config.DB.Where("roll_no = ?", rollNo).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}
	return c.JSON(fiber.Map{"student": student})
}
