package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
)

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
