package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/routes"
	"github.com/iips-oss/ispark/api/utils"
)

func seedCertFixtures(t *testing.T) {
	t.Helper()
	hashed, _ := utils.HashPassword("Password123")

	config.DB.Create(&models.Admin{AdminID: "cadmin", Name: "Batch Admin", Email: "c@a.dev", Password: hashed, Role: "admin", AssignedBatch: "IT2K24"})
	config.DB.Create(&models.Admin{AdminID: "csuper", Name: "Super", Email: "s@a.dev", Password: hashed, Role: "superadmin"})
	config.DB.Create(&models.Student{RollNo: "IT2K24001", Name: "In Batch", EmailID: "in@s.dev", EnrollmentNo: "EN-IN-001", Password: hashed, IsVerified: true})
	config.DB.Create(&models.Student{RollNo: "IT2K25001", Name: "Other Batch", EmailID: "out@s.dev", EnrollmentNo: "EN-OUT-001", Password: hashed, IsVerified: true})

	config.DB.Create(&models.Certificate{StudentRollNo: "IT2K24001", ActivityName: "In-Batch Cert", ActivityCategory: "TECHNICAL", ActivityDate: time.Now(), ParticipationType: "Participant", Status: "Pending", Credits: 3})
	config.DB.Create(&models.Certificate{StudentRollNo: "IT2K25001", ActivityName: "Other-Batch Cert", ActivityCategory: "TECHNICAL", ActivityDate: time.Now(), ParticipationType: "Participant", Status: "Pending", Credits: 3})
}

func certToken(t *testing.T, adminID, role string) string {
	t.Helper()
	tok, err := utils.GenerateAccessToken(adminID, adminID+"@a.dev", role)
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	return tok
}

func TestAdminCertificateVerification(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("test-jwt-", 4))
	t.Setenv("JWT_REFRESH_SECRET", strings.Repeat("test-refresh-jwt-", 4))
	SetupTestDB(t)
	seedCertFixtures(t)

	app := fiber.New()
	routes.SetupRoutes(app)

	batchTok := certToken(t, "cadmin", "admin")
	superTok := certToken(t, "csuper", "superadmin")

	// ids of the two seeded certs
	var inBatch, otherBatch models.Certificate
	config.DB.Where("student_roll_no = ?", "IT2K24001").First(&inBatch)
	config.DB.Where("student_roll_no = ?", "IT2K25001").First(&otherBatch)

	do := func(method, path, token, body string) *http.Response {
		var r *http.Request
		if body != "" {
			r = httptest.NewRequest(method, path, strings.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
		} else {
			r = httptest.NewRequest(method, path, nil)
		}
		r.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(r, -1)
		if err != nil {
			t.Fatalf("request %s %s: %v", method, path, err)
		}
		return resp
	}

	t.Run("batch admin lists only their batch", func(t *testing.T) {
		resp := do("GET", "/api/admin/certificates", batchTok, "")
		var out struct {
			Certificates []map[string]any `json:"certificates"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&out)
		if len(out.Certificates) != 1 {
			t.Fatalf("want 1 cert for batch admin, got %d", len(out.Certificates))
		}
		if out.Certificates[0]["student_roll_no"] != "IT2K24001" {
			t.Fatalf("wrong cert leaked: %v", out.Certificates[0]["student_roll_no"])
		}
	})

	t.Run("super admin sees all", func(t *testing.T) {
		resp := do("GET", "/api/admin/certificates", superTok, "")
		var out struct {
			Certificates []map[string]any `json:"certificates"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&out)
		if len(out.Certificates) != 2 {
			t.Fatalf("want 2 certs for super admin, got %d", len(out.Certificates))
		}
	})

	t.Run("batch admin cannot approve out-of-batch cert", func(t *testing.T) {
		resp := do("POST", "/api/admin/certificates/"+itoa(otherBatch.ID)+"/approve", batchTok, "")
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("want 404 for cross-batch approve, got %d", resp.StatusCode)
		}
	})

	t.Run("approve persists status and reviewer", func(t *testing.T) {
		resp := do("POST", "/api/admin/certificates/"+itoa(inBatch.ID)+"/approve", batchTok, "")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("want 200, got %d", resp.StatusCode)
		}
		var reloaded models.Certificate
		config.DB.First(&reloaded, inBatch.ID)
		if reloaded.Status != "Approved" || reloaded.ReviewedBy != "cadmin" || reloaded.ReviewedAt == nil {
			t.Fatalf("approve did not persist: status=%s by=%s at=%v", reloaded.Status, reloaded.ReviewedBy, reloaded.ReviewedAt)
		}
	})

	t.Run("reject requires a reason", func(t *testing.T) {
		resp := do("POST", "/api/admin/certificates/"+itoa(inBatch.ID)+"/reject", batchTok, `{"reason":""}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("want 400 for empty reason, got %d", resp.StatusCode)
		}
	})

	t.Run("reject persists reason", func(t *testing.T) {
		resp := do("POST", "/api/admin/certificates/"+itoa(inBatch.ID)+"/reject", batchTok, `{"reason":"Blurry scan"}`)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("want 200, got %d", resp.StatusCode)
		}
		var reloaded models.Certificate
		config.DB.First(&reloaded, inBatch.ID)
		if reloaded.Status != "Rejected" || reloaded.RejectionReason != "Blurry scan" {
			t.Fatalf("reject did not persist: status=%s reason=%s", reloaded.Status, reloaded.RejectionReason)
		}
	})

	t.Run("student token is rejected", func(t *testing.T) {
		studentTok := certToken(t, "IT2K24001", "student")
		resp := do("GET", "/api/admin/certificates", studentTok, "")
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("want 403 for student on admin route, got %d", resp.StatusCode)
		}
	})
}

func itoa(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
