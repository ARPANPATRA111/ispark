package controllers_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/controllers"
	"github.com/iips-oss/ispark/api/models"
	"gorm.io/gorm"
)

func setupAnnouncementApp(t *testing.T) *fiber.App {
	t.Helper()

	SetupTestDB(t)
	if err := config.DB.AutoMigrate(&models.Announcement{}); err != nil {
		t.Fatalf("Failed to migrate announcements: %v", err)
	}
	if err := config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Announcement{}).Error; err != nil {
		t.Fatalf("Failed to clear announcements: %v", err)
	}

	app := fiber.New()
	app.Get("/announcements/stats", controllers.GetAnnouncementStats)
	app.Get("/announcements", controllers.GetAnnouncements)
	app.Post("/announcements", controllers.CreateAnnouncement)
	app.Get("/announcements/:id", controllers.GetAnnouncement)
	app.Put("/announcements/:id", controllers.UpdateAnnouncement)
	app.Delete("/announcements/:id", controllers.DeleteAnnouncement)
	app.Post("/announcements/:id/publish", controllers.PublishAnnouncement)
	return app
}

func announcementDate(days int) string {
	now := time.Now().In(time.FixedZone("Asia/Kolkata", 5*60*60+30*60))
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return today.AddDate(0, 0, days).Format("2006-01-02")
}

func validAnnouncementBody() map[string]any {
	return map[string]any{
		"title":        "Runtime Review Notice",
		"description":  "A persistent announcement used by the controller tests.",
		"category":     "Academic",
		"audience":     "Students",
		"priority":     "High",
		"publish_date": announcementDate(5),
		"expiry_date":  announcementDate(20),
		"status":       "scheduled",
	}
}

func TestAnnouncementCRUDAndPublish(t *testing.T) {
	app := setupAnnouncementApp(t)

	createdResponse, createdBody := doJSON(t, app, http.MethodPost, "/announcements", validAnnouncementBody())
	if createdResponse.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d (%v)", createdResponse.StatusCode, createdBody)
	}
	created := createdBody["announcement"].(map[string]any)
	id := uint(created["id"].(float64))
	if created["status"] != "scheduled" {
		t.Fatalf("expected scheduled status, got %v", created["status"])
	}

	search := url.QueryEscape("academic")
	getResponse, getBody := doJSON(t, app, http.MethodGet, "/announcements?q="+search+"&status=scheduled", nil)
	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResponse.StatusCode)
	}
	announcements := getBody["announcements"].([]any)
	if len(announcements) != 1 {
		t.Fatalf("expected the fresh GET to return one persisted record, got %d", len(announcements))
	}

	updateResponse, updateBody := doJSON(t, app, http.MethodPut, fmt.Sprintf("/announcements/%d", id), map[string]any{
		"title":    "Updated Runtime Review Notice",
		"priority": "Medium",
	})
	if updateResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d (%v)", updateResponse.StatusCode, updateBody)
	}
	updated := updateBody["announcement"].(map[string]any)
	if updated["title"] != "Updated Runtime Review Notice" || updated["priority"] != "Medium" {
		t.Fatalf("expected the update to persist, got %v", updated)
	}
	if updated["expiry_date"] != announcementDate(20) {
		t.Fatalf("expected omitted dates to be preserved, got %v", updated["expiry_date"])
	}

	publishResponse, publishBody := doJSON(t, app, http.MethodPost, fmt.Sprintf("/announcements/%d/publish", id), nil)
	if publishResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d (%v)", publishResponse.StatusCode, publishBody)
	}
	published := publishBody["announcement"].(map[string]any)
	if published["status"] != "active" || published["publish_date"] != announcementDate(0) {
		t.Fatalf("expected publishing to make the notice active today, got %v", published)
	}

	statsResponse, stats := doJSON(t, app, http.MethodGet, "/announcements/stats", nil)
	if statsResponse.StatusCode != http.StatusOK || stats["total_announcements"] != float64(1) || stats["active_announcements"] != float64(1) {
		t.Fatalf("unexpected statistics response: status=%d body=%v", statsResponse.StatusCode, stats)
	}

	deleteResponse, _ := doJSON(t, app, http.MethodDelete, fmt.Sprintf("/announcements/%d", id), nil)
	if deleteResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", deleteResponse.StatusCode)
	}
	notFoundResponse, _ := doJSON(t, app, http.MethodGet, fmt.Sprintf("/announcements/%d", id), nil)
	if notFoundResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("expected deleted announcement to return 404, got %d", notFoundResponse.StatusCode)
	}
}

func TestAnnouncementValidation(t *testing.T) {
	tests := []struct {
		name   string
		change func(map[string]any)
	}{
		{name: "BlankTitle", change: func(body map[string]any) { body["title"] = "   " }},
		{name: "TitleTooLong", change: func(body map[string]any) { body["title"] = strings.Repeat("a", 201) }},
		{name: "InvalidCategory", change: func(body map[string]any) { body["category"] = "Finance" }},
		{name: "InvalidAudience", change: func(body map[string]any) { body["audience"] = "Guests" }},
		{name: "InvalidPriority", change: func(body map[string]any) { body["priority"] = "Urgent" }},
		{name: "InvalidStatus", change: func(body map[string]any) { body["status"] = "hidden" }},
		{name: "InvalidDate", change: func(body map[string]any) { body["publish_date"] = "22-07-2026" }},
		{name: "ExpiryBeforePublish", change: func(body map[string]any) { body["expiry_date"] = announcementDate(1) }},
		{name: "ScheduledToday", change: func(body map[string]any) { body["publish_date"] = announcementDate(0) }},
		{name: "ActiveInFuture", change: func(body map[string]any) { body["status"] = "active" }},
		{name: "ExpiredInFuture", change: func(body map[string]any) { body["status"] = "expired" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := setupAnnouncementApp(t)
			body := validAnnouncementBody()
			test.change(body)
			response, decoded := doJSON(t, app, http.MethodPost, "/announcements", body)
			if response.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d (%v)", response.StatusCode, decoded)
			}
		})
	}
}

func TestAnnouncementLifecycleRefresh(t *testing.T) {
	app := setupAnnouncementApp(t)
	past, _ := time.Parse("2006-01-02", announcementDate(-2))
	yesterday, _ := time.Parse("2006-01-02", announcementDate(-1))
	tomorrow, _ := time.Parse("2006-01-02", announcementDate(1))

	records := []models.Announcement{
		{Title: "Scheduled Now", Category: "General", Audience: "All Users", Priority: "Medium", PublishDate: yesterday, ExpiryDate: tomorrow, Status: "scheduled"},
		{Title: "Expired Now", Category: "Events", Audience: "Students", Priority: "Low", PublishDate: past, ExpiryDate: yesterday, Status: "active"},
		{Title: "Draft Stays Draft", Category: "Academic", Audience: "Mentors", Priority: "High", PublishDate: past, ExpiryDate: yesterday, Status: "draft"},
	}
	if err := config.DB.Create(&records).Error; err != nil {
		t.Fatalf("failed to seed lifecycle records: %v", err)
	}

	response, _ := doJSON(t, app, http.MethodGet, "/announcements", nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}

	statuses := map[string]string{}
	var stored []models.Announcement
	if err := config.DB.Find(&stored).Error; err != nil {
		t.Fatalf("failed to reload lifecycle records: %v", err)
	}
	for _, announcement := range stored {
		statuses[announcement.Title] = announcement.Status
	}
	if statuses["Scheduled Now"] != "active" || statuses["Expired Now"] != "expired" || statuses["Draft Stays Draft"] != "draft" {
		t.Fatalf("unexpected lifecycle transitions: %v", statuses)
	}
}

func TestAnnouncementUpdateAndPublishErrors(t *testing.T) {
	app := setupAnnouncementApp(t)

	response, _ := doJSON(t, app, http.MethodPut, "/announcements/999", map[string]any{"title": "Missing"})
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.StatusCode)
	}

	createdResponse, createdBody := doJSON(t, app, http.MethodPost, "/announcements", validAnnouncementBody())
	if createdResponse.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createdResponse.StatusCode)
	}
	id := uint(createdBody["announcement"].(map[string]any)["id"].(float64))

	emptyResponse, _ := doJSON(t, app, http.MethodPut, fmt.Sprintf("/announcements/%d", id), map[string]any{})
	if emptyResponse.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for an empty update, got %d", emptyResponse.StatusCode)
	}

	if err := config.DB.Model(&models.Announcement{}).Where("id = ?", id).
		Updates(map[string]any{"publish_date": announcementDate(-5), "expiry_date": announcementDate(-1), "status": "expired"}).Error; err != nil {
		t.Fatalf("failed to expire announcement: %v", err)
	}
	publishResponse, _ := doJSON(t, app, http.MethodPost, fmt.Sprintf("/announcements/%d/publish", id), nil)
	if publishResponse.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 when publishing an expired announcement, got %d", publishResponse.StatusCode)
	}
}
