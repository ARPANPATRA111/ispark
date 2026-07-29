package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/iips-oss/ispark/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// env reads a variable and strips surrounding whitespace. Values pasted into a
// hosting dashboard routinely pick up a trailing newline or space, and in the
// libpq key=value DSN format whitespace terminates a value — an untrimmed
// password is silently truncated and surfaces as an authentication failure.
func env(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

// quoteDSN renders a value for the libpq key=value DSN format, single-quoting
// it when it contains characters that would otherwise break parsing.
func quoteDSN(value string) string {
	if value == "" {
		return "''"
	}
	if !strings.ContainsAny(value, " '\\") {
		return value
	}

	var b strings.Builder
	b.WriteByte('\'')
	for _, r := range value {
		if r == '\'' || r == '\\' {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	b.WriteByte('\'')
	return b.String()
}

func ConnectDB() {
	var err error

	// A full connection string (as handed out by Supabase or Render) takes
	// precedence, since it is a single value and cannot be assembled wrongly.
	dsn := env("DATABASE_URL")
	if dsn != "" {
		log.Println("Connecting to database using DATABASE_URL...")
	} else {
		dbHost := env("DB_HOST")
		dbPort := env("DB_PORT")
		dbUser := env("DB_USER")
		dbPassword := env("DB_PASSWORD")
		dbName := env("DB_NAME")
		dbSSLMode := env("DB_SSLMODE")

		if dbSSLMode == "" {
			dbSSLMode = "disable"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			quoteDSN(dbHost), quoteDSN(dbUser), quoteDSN(dbPassword),
			quoteDSN(dbName), quoteDSN(dbPort), quoteDSN(dbSSLMode))

		// Password length only — never the value. Makes a truncated or missing
		// credential obvious in the deploy logs without leaking the secret.
		log.Printf("Connecting to database at %s:%s (user=%s db=%s sslmode=%s, password=%d chars)...",
			dbHost, dbPort, dbUser, dbName, dbSSLMode, len(dbPassword))
	}

	// Transaction-mode poolers (Supabase Supavisor/PgBouncer on port 6543)
	// don't support prepared statements, which pgx uses by default. Set
	// DB_PREFER_SIMPLE_PROTOCOL=true when connecting through one.
	preferSimple, _ := strconv.ParseBool(os.Getenv("DB_PREFER_SIMPLE_PROTOCOL"))

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: preferSimple,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Auto Migration
	log.Println("Running AutoMigration...")
	err = DB.AutoMigrate(&models.Student{}, &models.OTP{}, &models.Admin{}, &models.Activity{}, &models.Certificate{}, &models.Enrollment{}, &models.SystemSetting{}, &models.AdminNote{}, &models.Track{}, &models.Announcement{}, &models.GeneratedReport{}, &models.ScheduledReport{}, &models.ReportAuditLog{}, &models.BatchOverride{})

	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migration completed.")

	// Safe data migration: backfill coordinator_id from matching admin Name if it is empty/null
	var admins []models.Admin
	if err := DB.Find(&admins).Error; err == nil {
		for _, admin := range admins {
			if err := DB.Model(&models.Activity{}).
				Where("coordinator = ? AND (coordinator_id = ? OR coordinator_id IS NULL)", admin.Name, "").
				Update("coordinator_id", admin.AdminID).Error; err != nil {
				log.Printf("Warning: Failed to backfill coordinator_id for admin %s: %v", admin.AdminID, err)
			}
		}
	}

	// Safe data migration: backfill track_id for existing activities if null or 0
	var unassignedCount int64
	if err := DB.Model(&models.Activity{}).Where("track_id IS NULL OR track_id = 0").Count(&unassignedCount).Error; err == nil && unassignedCount > 0 {
		var skillTrack, personalityTrack models.Track
		DB.Where("LOWER(name) = ?", "skill building").FirstOrCreate(&skillTrack, models.Track{
			Name: "Skill Building", Description: "Technical and vocational activities that develop practical competencies.", Status: "Active",
		})
		DB.Where("LOWER(name) = ?", "personality development").FirstOrCreate(&personalityTrack, models.Track{
			Name: "Personality Development", Description: "Activities focused on personal growth, communication, and leadership skills.", Status: "Active",
		})

		if skillTrack.ID != 0 {
			DB.Model(&models.Activity{}).
				Where("(track_id IS NULL OR track_id = 0) AND UPPER(category) IN ?", []string{"TECHNICAL", "RESEARCH", "SPORTS", "CULTURAL"}).
				Update("track_id", skillTrack.ID)
		}
		if personalityTrack.ID != 0 {
			DB.Model(&models.Activity{}).
				Where("track_id IS NULL OR track_id = 0").
				Update("track_id", personalityTrack.ID)
		}
	}

	// Safe data migration: existing verified students predate the status column
	// (whose default is "Pending"), so mark them Active. This keeps the super
	// admin registry from reporting already-active accounts as Pending.
	DB.Model(&models.Student{}).
		Where("is_verified = ? AND status = ?", true, "Pending").
		Update("status", "Active")
}
