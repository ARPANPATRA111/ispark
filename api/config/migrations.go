package config

import (
	"log"

	"github.com/iips-oss/ispark/api/models"
)

// RunMigrations applies data migrations that must run on every boot, regardless
// of whether demo seeding is enabled. AutoMigrate (in ConnectDB) handles the
// schema; this handles data that needs to be reshaped to match current
// expectations — for example canonicalising legacy course names so production
// records line up with the report filters' canonical program list.
//
// Every migration here must be idempotent: it runs on each start and must be a
// no-op once the data is already in its target shape.
func RunMigrations() {
	if err := normalizeStudentCourses(); err != nil {
		// A migration failure should not stop the server from booting, but it
		// must be visible in the logs so it can be investigated.
		log.Printf("Migration warning: normalising student courses failed: %v", err)
	}
}

// normalizeStudentCourses rewrites any legacy course name to its canonical
// equivalent. It is idempotent: rows already on a canonical name match nothing.
// It runs as a migration (see RunMigrations) so production records are
// normalised without depending on development seeding.
func normalizeStudentCourses() error {
	for legacy, canonical := range models.CourseNameAliases() {
		if err := DB.Model(&models.Student{}).
			Where("LOWER(TRIM(course_name)) = ?", legacy).
			Update("course_name", canonical).Error; err != nil {
			return err
		}
	}
	return nil
}
