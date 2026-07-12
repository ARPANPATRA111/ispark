package config

import (
	"log"

	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
)

func SeedDefaultAdmin() {
	hashed, _ := utils.HashPassword("admin123")

	var admin models.Admin
	if err := DB.Where(models.Admin{AdminID: "admin"}).
		Attrs(models.Admin{Password: hashed, MustChangePassword: true}).
		Assign(models.Admin{Name: "Platform Admin", Email: "admin@isparc.dev", Role: "admin", AssignedBatch: "IT2K24"}).
		FirstOrCreate(&admin).Error; err != nil {
		log.Printf("Seed admin failed: %v", err)
	}
}
