package models

import (
	"time"

	"gorm.io/gorm"
)

// GeneratedReport is a single report a super admin has produced from the
// Reports Center. The heavy aggregation runs at generation time and the result
// is written to a CSV file under the uploads directory; this row tracks the
// report's metadata, the file it produced and its processing status.
type GeneratedReport struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Type        string         `gorm:"type:varchar(100);not null" json:"type"`
	Course      string         `gorm:"type:varchar(100)" json:"course"`
	Semester    string         `gorm:"type:varchar(50)" json:"semester"`
	DateFrom    *time.Time     `json:"date_from"`
	DateTo      *time.Time     `json:"date_to"`
	Format      string         `gorm:"type:varchar(20);not null" json:"format"`                      // CSV, PDF, Excel
	Status      string         `gorm:"type:varchar(20);not null;default:'Processing'" json:"status"` // Processing, Ready, Failed
	GeneratedBy string         `gorm:"type:varchar(100)" json:"generated_by"`
	FileName    string         `gorm:"type:varchar(255)" json:"file_name"`
	FilePath    string         `gorm:"type:varchar(255)" json:"-"`
	FileSize    int64          `gorm:"default:0" json:"file_size"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// ScheduledReport is a report the platform generates automatically on a
// recurring schedule. It is managed from the Reports Center's scheduled
// reports panel.
type ScheduledReport struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Type      string         `gorm:"type:varchar(100);not null" json:"type"`
	Frequency string         `gorm:"type:varchar(100);not null" json:"frequency"` // e.g. "Weekly · Monday", "Every 1st of the month"
	Format    string         `gorm:"type:varchar(20);not null" json:"format"`
	Enabled   bool           `gorm:"not null;default:true" json:"enabled"`
	LastRunAt *time.Time     `json:"last_run_at"`
	NextRunAt *time.Time     `json:"next_run_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ReportAuditLog records a single report-related action for the Reports Center
// activity log. Category is the machine-readable action bucket used for
// filtering and counting (Generate, Download, Export, Schedule); Format is the
// human badge the UI renders (PDF, Excel, CSV, Auto).
type ReportAuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Action    string    `gorm:"type:varchar(255);not null" json:"action"`
	Category  string    `gorm:"type:varchar(50);not null;index" json:"category"` // Generate, Download, Export, Schedule
	Format    string    `gorm:"type:varchar(20)" json:"type"`                    // PDF, Excel, CSV, Auto
	User      string    `gorm:"type:varchar(100)" json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

// GenerateReportInput is the body accepted by the "Generate Report" form.
type GenerateReportInput struct {
	Type     string `json:"type"`
	Course   string `json:"course"`
	Semester string `json:"semester"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Format   string `json:"format"`
}

// ScheduledReportInput is the body accepted when creating or updating a
// scheduled report. Enabled is a pointer so an update can toggle it on its own
// without resending the other fields.
type ScheduledReportInput struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Frequency string `json:"frequency"`
	Format    string `json:"format"`
	Enabled   *bool  `json:"enabled"`
}
