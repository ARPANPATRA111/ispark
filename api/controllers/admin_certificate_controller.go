package controllers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/storage"
)

// certWithStudent is a certificate enriched with the few student fields the
// review queue needs, so the frontend does not have to make a second call.
type certWithStudent struct {
	models.Certificate
	StudentName   string `json:"student_name"`
	StudentCourse string `json:"student_course"`
}

// loadCertificateForAdmin fetches a certificate by id and confirms the calling
// admin is allowed to act on it: a batch admin only reaches certificates whose
// student roll number falls in their assigned batch; a super admin reaches all.
func loadCertificateForAdmin(c *fiber.Ctx) (*models.Certificate, *models.Admin, error) {
	admin, err := getAuthenticatedAdmin(c)
	if err != nil {
		return nil, nil, err
	}

	var cert models.Certificate
	if err := config.DB.Where("id = ?", c.Params("id")).First(&cert).Error; err != nil {
		return nil, nil, fiber.NewError(fiber.StatusNotFound, "Certificate not found")
	}

	if admin.Role == "admin" {
		if admin.AssignedBatch == "" || !strings.HasPrefix(cert.StudentRollNo, admin.AssignedBatch) {
			// Do not reveal that the certificate exists outside the batch.
			return nil, nil, fiber.NewError(fiber.StatusNotFound, "Certificate not found")
		}
	}

	return &cert, admin, nil
}

// GetAdminCertificates returns certificates for review, scoped to the admin's
// batch (all for a super admin), newest first. Supports ?status=Pending filter.
func GetAdminCertificates(c *fiber.Ctx) error {
	admin, err := getAuthenticatedAdmin(c)
	if err != nil {
		return err
	}

	query := config.DB.Model(&models.Certificate{}).
		Select("certificates.*, students.name AS student_name, students.course_name AS student_course").
		Joins("JOIN students ON students.roll_no = certificates.student_roll_no").
		Order("certificates.created_at DESC")

	if admin.Role == "admin" {
		if admin.AssignedBatch == "" {
			return c.JSON(fiber.Map{"certificates": []certWithStudent{}})
		}
		query = query.Where("certificates.student_roll_no LIKE ?", admin.AssignedBatch+"%")
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" && !strings.EqualFold(status, "All") {
		query = query.Where("certificates.status = ?", status)
	}

	certs := []certWithStudent{}
	if err := query.Scan(&certs).Error; err != nil {
		return errJSON(c, fiber.StatusInternalServerError, "Failed to load certificates")
	}

	return c.JSON(fiber.Map{"certificates": certs})
}

// DownloadAdminCertificate streams a student's certificate file to an
// authorised admin, reusing the same storage backend as the student download.
func DownloadAdminCertificate(c *fiber.Ctx) error {
	cert, _, err := loadCertificateForAdmin(c)
	if err != nil {
		return err
	}

	ref := cert.FilePath
	if ref == "" {
		ref = cert.FileName
	}
	f, err := storage.Default().Open(ref)
	if err != nil {
		return errJSON(c, fiber.StatusNotFound, "Certificate file is no longer available")
	}

	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", cert.FileName))
	if contentType, ok := allowedCertificateTypes[strings.ToLower(filepath.Ext(cert.FileName))]; ok {
		c.Set(fiber.HeaderContentType, contentType)
	}
	return c.SendStream(f)
}

// ApproveCertificate marks a certificate approved. Credits already sit on the
// row from upload time; approval is what makes them count towards dashboards
// and the leaderboard, so this completes the upload -> verify -> credit loop.
func ApproveCertificate(c *fiber.Ctx) error {
	cert, admin, err := loadCertificateForAdmin(c)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := map[string]any{
		"status":           "Approved",
		"rejection_reason": "",
		"reviewed_by":      admin.AdminID,
		"reviewed_at":      now,
	}
	if err := config.DB.Model(cert).Updates(updates).Error; err != nil {
		return errJSON(c, fiber.StatusInternalServerError, "Failed to approve certificate")
	}

	return c.JSON(fiber.Map{"message": "Certificate approved", "certificate": cert})
}

// RejectCertificate marks a certificate rejected with a required reason so the
// student knows what to fix before re-uploading.
func RejectCertificate(c *fiber.Ctx) error {
	var input struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&input); err != nil {
		return errJSON(c, fiber.StatusBadRequest, "Cannot parse request body")
	}
	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		return errJSON(c, fiber.StatusBadRequest, "A rejection reason is required")
	}

	cert, admin, err := loadCertificateForAdmin(c)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := map[string]any{
		"status":           "Rejected",
		"rejection_reason": reason,
		"reviewed_by":      admin.AdminID,
		"reviewed_at":      now,
	}
	if err := config.DB.Model(cert).Updates(updates).Error; err != nil {
		return errJSON(c, fiber.StatusInternalServerError, "Failed to reject certificate")
	}

	return c.JSON(fiber.Map{"message": "Certificate rejected", "certificate": cert})
}
