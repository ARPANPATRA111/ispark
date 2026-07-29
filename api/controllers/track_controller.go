package controllers

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"gorm.io/gorm"
)

// allowedTrackStatuses is the set of statuses a track may hold.
var allowedTrackStatuses = map[string]bool{
	"Active":   true,
	"Inactive": true,
}

// maxTrackNameLength mirrors the varchar(100) limit on Track.Name so an
// over-long name is rejected with a 400 instead of surfacing as a database
// error (500) once it hits the column constraint.
const maxTrackNameLength = 100

// trackActivityCounts returns the number of activities assigned to each track,
// keyed by track ID. Tracks with no activities are simply absent from the map.
func trackActivityCounts() (map[uint]int64, error) {
	type row struct {
		TrackID uint
		Count   int64
	}
	var rows []row
	if err := config.DB.Model(&models.Activity{}).
		Select("track_id, COUNT(*) AS count").
		Where("track_id IS NOT NULL").
		Group("track_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	counts := make(map[uint]int64, len(rows))
	for _, r := range rows {
		counts[r.TrackID] = r.Count
	}
	return counts, nil
}

// attachActivityCounts fills in TotalActivities for each track from the counts map.
func attachActivityCounts(tracks []models.Track, counts map[uint]int64) {
	for i := range tracks {
		tracks[i].TotalActivities = counts[tracks[i].ID]
	}
}

// trackNameTaken reports whether an active track already uses the given name,
// case-insensitively. excludeID skips a track (used when renaming in place).
func trackNameTaken(name string, excludeID uint) (bool, error) {
	query := config.DB.Model(&models.Track{}).Where("LOWER(name) = LOWER(?)", name)
	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetTracks returns every track with its activity count, supporting a status
// filter (?status=Active|Inactive) and a search query (?q= over name and
// description).
func GetTracks(c *fiber.Ctx) error {
	query := config.DB.Model(&models.Track{}).Order("created_at asc")

	if status := strings.TrimSpace(c.Query("status")); status != "" && !strings.EqualFold(status, "All") {
		if !allowedTrackStatuses[status] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status must be Active or Inactive"})
		}
		query = query.Where("status = ?", status)
	}

	if q := strings.TrimSpace(c.Query("q")); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", like, like)
	}

	var tracks []models.Track
	if err := query.Find(&tracks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load tracks"})
	}

	counts, err := trackActivityCounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load tracks"})
	}
	attachActivityCounts(tracks, counts)

	return c.JSON(fiber.Map{"tracks": tracks})
}

// GetTrack returns a single track with its activity count.
func GetTrack(c *fiber.Ctx) error {
	var track models.Track
	if err := config.DB.Where("id = ?", c.Params("id")).First(&track).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Track not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load track"})
	}

	var count int64
	if err := config.DB.Model(&models.Activity{}).Where("track_id = ?", track.ID).Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load track"})
	}
	track.TotalActivities = count

	return c.JSON(fiber.Map{"track": track})
}

// CreateTrack registers a new track. Names must be unique among active tracks.
func CreateTrack(c *fiber.Ctx) error {
	var input models.TrackInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Track name is required"})
	}
	if utf8.RuneCountInString(name) > maxTrackNameLength {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Track name must be 100 characters or fewer"})
	}

	description := ""
	if input.Description != nil {
		description = strings.TrimSpace(*input.Description)
	}

	status := "Active"
	if input.Status != nil {
		if !allowedTrackStatuses[*input.Status] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status must be Active or Inactive"})
		}
		status = *input.Status
	}

	taken, err := trackNameTaken(name, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create track"})
	}
	if taken {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A track with this name already exists"})
	}

	track := models.Track{
		Name:        name,
		Description: description,
		Status:      status,
	}
	if err := config.DB.Create(&track).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create track"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"track": track})
}

// UpdateTrack edits a track's name, description and/or status. Only fields
// present in the request are changed.
func UpdateTrack(c *fiber.Ctx) error {
	var track models.Track
	if err := config.DB.Where("id = ?", c.Params("id")).First(&track).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Track not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load track"})
	}

	var input models.TrackInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if name := strings.TrimSpace(input.Name); name != "" {
		if utf8.RuneCountInString(name) > maxTrackNameLength {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Track name must be 100 characters or fewer"})
		}
		if name != track.Name {
			taken, err := trackNameTaken(name, track.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update track"})
			}
			if taken {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A track with this name already exists"})
			}
			track.Name = name
		}
	}

	if input.Status != nil {
		if !allowedTrackStatuses[*input.Status] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status must be Active or Inactive"})
		}
		track.Status = *input.Status
	}

	// Only touch the description when the caller actually sent the field. A nil
	// pointer means it was omitted (e.g. a status-only update), so the existing
	// description is preserved; a non-nil pointer may explicitly clear it.
	if input.Description != nil {
		track.Description = strings.TrimSpace(*input.Description)
	}

	if err := config.DB.Save(&track).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update track"})
	}

	var count int64
	if err := config.DB.Model(&models.Activity{}).Where("track_id = ?", track.ID).Count(&count).Error; err == nil {
		track.TotalActivities = count
	}

	return c.JSON(fiber.Map{"track": track})
}

// DeleteTrack removes a track. Deletion is blocked while activities are still
// assigned to it, so credits and activities are never orphaned.
func DeleteTrack(c *fiber.Ctx) error {
	var track models.Track
	if err := config.DB.Where("id = ?", c.Params("id")).First(&track).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Track not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load track"})
	}

	var assigned int64
	if err := config.DB.Model(&models.Activity{}).Where("track_id = ?", track.ID).Count(&assigned).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete track"})
	}
	if assigned > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Cannot delete a track that still has activities assigned to it",
		})
	}

	if err := config.DB.Delete(&track).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete track"})
	}

	return c.JSON(fiber.Map{"message": "Track deleted successfully"})
}

// GetTrackStats returns the Track Management overview counts.
func GetTrackStats(c *fiber.Ctx) error {
	var total, active, inactive, assignedActivities int64

	if err := config.DB.Model(&models.Track{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load track stats"})
	}
	if err := config.DB.Model(&models.Track{}).Where("status = ?", "Active").Count(&active).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load track stats"})
	}
	if err := config.DB.Model(&models.Track{}).Where("status = ?", "Inactive").Count(&inactive).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load track stats"})
	}
	if err := config.DB.Model(&models.Activity{}).Where("track_id IS NOT NULL").Count(&assignedActivities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load track stats"})
	}

	return c.JSON(fiber.Map{
		"total_tracks":        total,
		"active_tracks":       active,
		"inactive_tracks":     inactive,
		"assigned_activities": assignedActivities,
	})
}
