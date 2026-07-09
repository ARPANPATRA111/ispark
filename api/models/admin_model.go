package models

import (
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	AdminID            string         `gorm:"primaryKey;type:varchar(50)" json:"admin_id"`
	Name               string         `gorm:"type:varchar(100);not null" json:"name"`
	Email              string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password           string         `gorm:"type:varchar(255);not null" json:"-"`
	Role               string         `gorm:"type:varchar(20);not null;default:'admin'" json:"role"`
	MustChangePassword bool           `gorm:"default:false" json:"must_change_password"` // NEW
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

type AdminLoginInput struct {
	AdminID  string `json:"admin_id"`
	Password string `json:"password"`
}
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
