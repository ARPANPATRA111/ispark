package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemSetting is a single platform-wide configuration entry managed from the
// super admin System Settings screen. Settings are grouped into categories
// (the UI tabs) and rendered in SortOrder within each category.
type SystemSetting struct {
	Key         string         `gorm:"primaryKey;type:varchar(100)" json:"key"`
	Category    string         `gorm:"type:varchar(50);not null;index" json:"category"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Value       string         `gorm:"type:varchar(255)" json:"value"`
	Status      string         `gorm:"type:varchar(20);not null;default:'Active'" json:"status"` // Active, Enabled, Disabled
	SortOrder   int            `gorm:"not null;default:0" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// UpdateSettingInput is the body accepted when updating a single setting. Both
// fields are pointers so the caller can update value and status independently.
type UpdateSettingInput struct {
	Value  *string `json:"value"`
	Status *string `json:"status"`
}

// BulkSettingUpdate is one entry in a bulk settings update request.
type BulkSettingUpdate struct {
	Key    string  `json:"key"`
	Value  *string `json:"value"`
	Status *string `json:"status"`
}
