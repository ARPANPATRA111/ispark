package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/controllers"
	"github.com/iips-oss/ispark/api/models"
	"gorm.io/gorm"
)

func setupSuperadminApp(t *testing.T) *fiber.App {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory SQLite database: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Student{},
		&models.Admin{},
		&models.Activity{},
		&models.Certificate{},
		&models.Track{},
	); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	config.DB = db

	app := fiber.New()
	app.Get("/api/admin/platform/stats", controllers.GetPlatformStats)
	return app
}

func TestGetPlatformStats_CountsActiveTracksFromTracksTable(t *testing.T) {
	app := setupSuperadminApp(t)

	// Seed two active track records
	t1 := models.Track{Name: "Track 1", Status: "Active"}
	t2 := models.Track{Name: "Track 2", Status: "Active"}
	t3 := models.Track{Name: "Track 3", Status: "Inactive"}
	if err := config.DB.Create(&t1).Error; err != nil {
		t.Fatalf("failed to create track 1: %v", err)
	}
	if err := config.DB.Create(&t2).Error; err != nil {
		t.Fatalf("failed to create track 2: %v", err)
	}
	if err := config.DB.Create(&t3).Error; err != nil {
		t.Fatalf("failed to create track 3: %v", err)
	}

	// Seed one activity category
	act := models.Activity{Name: "Activity 1", Category: "TECHNICAL", Track: t1, Credits: 2, Status: "Active"}
	if err := config.DB.Create(&act).Error; err != nil {
		t.Fatalf("failed to create activity: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/admin/platform/stats", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("GET /api/admin/platform/stats failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode json response: %v", err)
	}

	activeTracks, ok := body["active_tracks"].(float64)
	if !ok {
		t.Fatalf("expected active_tracks in response body, got %v", body)
	}

	if int64(activeTracks) != 2 {
		t.Errorf("expected active_tracks to be 2 (counted from tracks table), got %v", activeTracks)
	}
}
