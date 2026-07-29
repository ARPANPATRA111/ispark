package config

import (
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/iips-oss/ispark/api/models"
	"gorm.io/gorm"
)

func TestNormalizeStudentCourses(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Student{}); err != nil {
		t.Fatalf("migrate students: %v", err)
	}

	previousDB := DB
	DB = db
	t.Cleanup(func() { DB = previousDB })

	tests := []struct {
		legacy    string
		canonical string
	}{
		{"M.Tech (Computer Science - CS)", models.CourseMTechCS},
		{"M.Tech (Information Technology - IT)", models.CourseMTechIT},
		{"MCA (Master of Computer Applications)", models.CourseMCA5Yr},
		{"MBA (Management Science - MS)", models.CourseMBAMS5Yr},
		{"MBA (Management Science)", models.CourseMBAMS2Yr},
		{"MBA (Advertising and Public Relations - APR)", models.CourseMBAAPR},
		{"MBA (Entrepreneurship)", models.CourseMBAEnt},
		{"B.Com. (Hons.)", models.CourseBComHons},
		{"  mca  ", models.CourseMCA5Yr},
	}

	for i, test := range tests {
		student := models.Student{
			RollNo:       fmt.Sprintf("R%02d", i),
			Name:         "Migration Test",
			CourseName:   test.legacy,
			Semester:     1,
			EmailID:      fmt.Sprintf("student-%d@example.com", i),
			EnrollmentNo: fmt.Sprintf("E%02d", i),
			Password:     "test",
		}
		if err := db.Create(&student).Error; err != nil {
			t.Fatalf("seed %q: %v", test.legacy, err)
		}
	}

	if err := normalizeStudentCourses(); err != nil {
		t.Fatalf("normalize courses: %v", err)
	}
	if err := normalizeStudentCourses(); err != nil {
		t.Fatalf("normalize courses a second time: %v", err)
	}

	for i, test := range tests {
		var student models.Student
		if err := db.First(&student, "roll_no = ?", fmt.Sprintf("R%02d", i)).Error; err != nil {
			t.Fatalf("load %q: %v", test.legacy, err)
		}
		if student.CourseName != test.canonical {
			t.Errorf("%q: expected %q, got %q", test.legacy, test.canonical, student.CourseName)
		}
	}
}
