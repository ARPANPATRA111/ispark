package models

import (
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	AdminID   string         `gorm:"primaryKey;type:varchar(50)" json:"admin_id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Role      string         `gorm:"type:varchar(20);not null;default:'admin'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Certificate to be verified by Admin/Mentor
type Certificate struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StudentRollNo string    `gorm:"type:varchar(50);not null" json:"student_roll_no"`
	ActivityName  string    `gorm:"type:varchar(200);not null" json:"activity_name"`
	FileURL       string    `gorm:"type:text;not null" json:"file_url"`
	Status        string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, rejected
	Feedback      string    `gorm:"type:text" json:"feedback"`
	VerifiedBy    string    `gorm:"type:varchar(50)" json:"verified_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Recommendation represents a course suggested by a mentor to a student
type Recommendation struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StudentRollNo string    `gorm:"type:varchar(50);not null" json:"student_roll_no"`
	MentorID      string    `gorm:"type:varchar(50);not null" json:"mentor_id"`
	CourseName    string    `gorm:"type:varchar(100);not null" json:"course_name"`
	Message       string    `gorm:"type:text" json:"message"`
	CreatedAt     time.Time `json:"created_at"`
}
