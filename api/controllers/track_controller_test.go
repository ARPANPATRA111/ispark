package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/controllers"
	"github.com/iips-oss/ispark/api/models"
	"gorm.io/gorm"
)

// setupTrackApp spins up an isolated in-memory database and a Fiber app wired to
// the track handlers directly. Auth middleware is intentionally omitted so the
// tests exercise the controller logic on its own.
func setupTrackApp(t *testing.T) *fiber.App {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory SQLite database: %v", err)
	}
	if err := db.AutoMigrate(&models.Track{}, &models.Activity{}); err != nil {
		t.Fatalf("Failed to migrate track/activity tables: %v", err)
	}
	config.DB = db

	app := fiber.New()
	app.Get("/tracks/stats", controllers.GetTrackStats)
	app.Get("/tracks", controllers.GetTracks)
	app.Post("/tracks", controllers.CreateTrack)
	app.Get("/tracks/:id", controllers.GetTrack)
	app.Put("/tracks/:id", controllers.UpdateTrack)
	app.Delete("/tracks/:id", controllers.DeleteTrack)

	return app
}

// doJSON issues a request with an optional JSON body and returns the response
// together with its decoded body.
func doJSON(t *testing.T, app *fiber.App, method, path string, body any) (*http.Response, map[string]any) {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("%s %s failed: %v", method, path, err)
	}

	decoded := map[string]any{}
	if res.Body != nil {
		_ = json.NewDecoder(res.Body).Decode(&decoded)
	}

	return res, decoded
}

// seedTrack inserts a track straight into the database for test setup.
func seedTrack(t *testing.T, name, description, status string) models.Track {
	t.Helper()
	track := models.Track{Name: name, Description: description, Status: status}
	if err := config.DB.Create(&track).Error; err != nil {
		t.Fatalf("Failed to seed track %q: %v", name, err)
	}
	return track
}

// seedActivityForTrack inserts an activity linked to the given track.
func seedActivityForTrack(t *testing.T, name string, trackID uint) {
	t.Helper()
	activity := models.Activity{Name: name, Category: "TECHNICAL", Credits: 5, Mode: "Online", TrackID: &trackID}
	if err := config.DB.Create(&activity).Error; err != nil {
		t.Fatalf("Failed to seed activity %q: %v", name, err)
	}
}

func TestCreateTrack(t *testing.T) {
	app := setupTrackApp(t)

	t.Run("Success", func(t *testing.T) {
		res, body := doJSON(t, app, http.MethodPost, "/tracks", map[string]any{
			"name":        "Research & Innovation",
			"description": "Research-oriented activities.",
		})
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d (%v)", res.StatusCode, body)
		}
		track, ok := body["track"].(map[string]any)
		if !ok {
			t.Fatalf("expected a track in the response, got %v", body)
		}
		if track["name"] != "Research & Innovation" {
			t.Errorf("expected name to be persisted, got %v", track["name"])
		}
		if track["status"] != "Active" {
			t.Errorf("expected default status Active, got %v", track["status"])
		}
	})

	t.Run("NameRequired", func(t *testing.T) {
		res, _ := doJSON(t, app, http.MethodPost, "/tracks", map[string]any{"name": "   "})
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400 for blank name, got %d", res.StatusCode)
		}
	})

	t.Run("NameTooLongReturns400", func(t *testing.T) {
		res, body := doJSON(t, app, http.MethodPost, "/tracks", map[string]any{
			"name": strings.Repeat("a", 101),
		})
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400 for a 101-character name, got %d (%v)", res.StatusCode, body)
		}
	})

	t.Run("NameAtLimitAllowed", func(t *testing.T) {
		res, _ := doJSON(t, app, http.MethodPost, "/tracks", map[string]any{
			"name": strings.Repeat("b", 100),
		})
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201 for a 100-character name, got %d", res.StatusCode)
		}
	})

	t.Run("DuplicateNameReturns409", func(t *testing.T) {
		seedTrack(t, "Sports", "", "Active")
		res, _ := doJSON(t, app, http.MethodPost, "/tracks", map[string]any{"name": "sports"})
		if res.StatusCode != http.StatusConflict {
			t.Fatalf("expected 409 for a duplicate name, got %d", res.StatusCode)
		}
	})
}

func TestGetTracksWithActivityCounts(t *testing.T) {
	app := setupTrackApp(t)

	pd := seedTrack(t, "Personality Development", "", "Active")
	sb := seedTrack(t, "Skill Building", "", "Active")
	seedActivityForTrack(t, "Debate", pd.ID)
	seedActivityForTrack(t, "Leadership", pd.ID)
	seedActivityForTrack(t, "Python", sb.ID)

	res, body := doJSON(t, app, http.MethodGet, "/tracks", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	tracks, ok := body["tracks"].([]any)
	if !ok {
		t.Fatalf("expected a tracks array, got %v", body["tracks"])
	}

	counts := map[string]float64{}
	for _, raw := range tracks {
		track := raw.(map[string]any)
		counts[track["name"].(string)] = track["total_activities"].(float64)
	}

	if counts["Personality Development"] != 2 {
		t.Errorf("expected 2 activities for Personality Development, got %v", counts["Personality Development"])
	}
	if counts["Skill Building"] != 1 {
		t.Errorf("expected 1 activity for Skill Building, got %v", counts["Skill Building"])
	}
}

func TestGetTrackNotFound(t *testing.T) {
	app := setupTrackApp(t)
	res, _ := doJSON(t, app, http.MethodGet, "/tracks/999", nil)
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

func TestUpdateTrackPartial(t *testing.T) {
	t.Run("StatusOnlyPreservesDescription", func(t *testing.T) {
		app := setupTrackApp(t)
		track := seedTrack(t, "Community Service", "Volunteering and outreach.", "Active")

		// A status-only PUT must not wipe the existing description.
		res, body := doJSON(t, app, http.MethodPut, fmt.Sprintf("/tracks/%d", track.ID), map[string]any{
			"status": "Inactive",
		})
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d (%v)", res.StatusCode, body)
		}

		updated := body["track"].(map[string]any)
		if updated["description"] != "Volunteering and outreach." {
			t.Errorf("expected description to be preserved, got %v", updated["description"])
		}
		if updated["status"] != "Inactive" {
			t.Errorf("expected status Inactive, got %v", updated["status"])
		}
	})

	t.Run("EmptyDescriptionClearsIt", func(t *testing.T) {
		app := setupTrackApp(t)
		track := seedTrack(t, "Arts", "Creative activities.", "Active")

		res, body := doJSON(t, app, http.MethodPut, fmt.Sprintf("/tracks/%d", track.ID), map[string]any{
			"description": "",
		})
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		if desc := body["track"].(map[string]any)["description"]; desc != "" {
			t.Errorf("expected an explicit empty description to clear it, got %v", desc)
		}
	})

	t.Run("NameTooLongReturns400", func(t *testing.T) {
		app := setupTrackApp(t)
		track := seedTrack(t, "Coding", "", "Active")

		res, _ := doJSON(t, app, http.MethodPut, fmt.Sprintf("/tracks/%d", track.ID), map[string]any{
			"name": strings.Repeat("x", 101),
		})
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400 for a 101-character name, got %d", res.StatusCode)
		}
	})

	t.Run("InvalidStatusReturns400", func(t *testing.T) {
		app := setupTrackApp(t)
		track := seedTrack(t, "Music", "", "Active")

		res, _ := doJSON(t, app, http.MethodPut, fmt.Sprintf("/tracks/%d", track.ID), map[string]any{
			"status": "Archived",
		})
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400 for an invalid status, got %d", res.StatusCode)
		}
	})
}

func TestDeleteTrack(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app := setupTrackApp(t)
		track := seedTrack(t, "Empty Track", "", "Active")

		res, _ := doJSON(t, app, http.MethodDelete, fmt.Sprintf("/tracks/%d", track.ID), nil)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}

		var count int64
		config.DB.Model(&models.Track{}).Where("id = ?", track.ID).Count(&count)
		if count != 0 {
			t.Errorf("expected the track to be deleted, still found %d", count)
		}
	})

	t.Run("BlockedWhenActivitiesLinked", func(t *testing.T) {
		app := setupTrackApp(t)
		track := seedTrack(t, "Busy Track", "", "Active")
		seedActivityForTrack(t, "Hackathon", track.ID)

		res, _ := doJSON(t, app, http.MethodDelete, fmt.Sprintf("/tracks/%d", track.ID), nil)
		if res.StatusCode != http.StatusConflict {
			t.Fatalf("expected 409 when activities are still linked, got %d", res.StatusCode)
		}

		var count int64
		config.DB.Model(&models.Track{}).Where("id = ?", track.ID).Count(&count)
		if count != 1 {
			t.Errorf("expected the track to survive the blocked delete, found %d", count)
		}
	})
}

func TestGetTrackStats(t *testing.T) {
	app := setupTrackApp(t)

	active := seedTrack(t, "Active One", "", "Active")
	seedTrack(t, "Active Two", "", "Active")
	seedTrack(t, "Inactive One", "", "Inactive")
	seedActivityForTrack(t, "Linked Activity", active.ID)
	// An unlinked activity must not count toward assigned_activities.
	if err := config.DB.Create(&models.Activity{Name: "Floating", Category: "SPORTS", Credits: 3, Mode: "Offline"}).Error; err != nil {
		t.Fatalf("Failed to seed unlinked activity: %v", err)
	}

	res, body := doJSON(t, app, http.MethodGet, "/tracks/stats", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	expect := map[string]float64{
		"total_tracks":        3,
		"active_tracks":       2,
		"inactive_tracks":     1,
		"assigned_activities": 1,
	}
	for key, want := range expect {
		if got, ok := body[key].(float64); !ok || got != want {
			t.Errorf("expected %s = %v, got %v", key, want, body[key])
		}
	}
}

func TestDeleteAssignedTrackBlockedWithMultipleActivities(t *testing.T) {
	app := setupTrackApp(t)
	track := seedTrack(t, "Personality Development", "Seeded track", "Active")

	// Seed 5 activities assigned to this track
	for i := 1; i <= 5; i++ {
		seedActivityForTrack(t, fmt.Sprintf("Activity %d", i), track.ID)
	}

	res, body := doJSON(t, app, http.MethodDelete, fmt.Sprintf("/tracks/%d", track.ID), nil)
	if res.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict when 5 activities are linked, got %d (%v)", res.StatusCode, body)
	}

	errMsg, _ := body["error"].(string)
	if !strings.Contains(errMsg, "Cannot delete a track that still has activities assigned to it") {
		t.Errorf("expected deletion blocked message, got %q", errMsg)
	}

	var count int64
	config.DB.Model(&models.Track{}).Where("id = ?", track.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected track to survive deletion attempt, found count = %d", count)
	}
}

func TestCreateAndFilterNewlyCreatedTrack(t *testing.T) {
	app := setupTrackApp(t)

	// 1. Create a new track "Research and Innovation"
	res, body := doJSON(t, app, http.MethodPost, "/tracks", map[string]any{
		"name":        "Research and Innovation",
		"description": "Activities focused on research and creative innovation.",
		"status":      "Active",
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d (%v)", res.StatusCode, body)
	}

	trackData, ok := body["track"].(map[string]any)
	if !ok {
		t.Fatalf("expected track object in response, got %v", body)
	}
	trackID := uint(trackData["id"].(float64))

	// 2. Assign an activity to it
	seedActivityForTrack(t, "Paper Presentation", trackID)

	// 3. GET /tracks?q=Research
	resSearch, bodySearch := doJSON(t, app, http.MethodGet, "/tracks?q=Research", nil)
	if resSearch.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for search query, got %d", resSearch.StatusCode)
	}

	tracksArr, ok := bodySearch["tracks"].([]any)
	if !ok || len(tracksArr) != 1 {
		t.Fatalf("expected 1 track matching 'Research', got %v", bodySearch["tracks"])
	}

	foundTrack := tracksArr[0].(map[string]any)
	if foundTrack["name"] != "Research and Innovation" {
		t.Errorf("expected track name 'Research and Innovation', got %v", foundTrack["name"])
	}
	if foundTrack["total_activities"].(float64) != 1 {
		t.Errorf("expected total_activities = 1, got %v", foundTrack["total_activities"])
	}
}
