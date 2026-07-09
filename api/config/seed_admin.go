package config

import (
	"log"

	"github.com/iips-oss/ispark/api/models"
	"github.com/iips-oss/ispark/api/utils"
)

func SeedDefaultAdmin() {
	hashed, _ := utils.HashPassword("admin123")

	admin := models.Admin{
		AdminID:            "admin",
		Name:               "Platform Admin",
		Email:              "admin@isparc.dev",
		Password:           hashed,
		Role:               "admin",
		MustChangePassword: true,
	}

	DB.Where("admin_id = ?", "admin").FirstOrCreate(&admin)
	log.Println("Seed check done -> admin_id: admin | password: admin123")
}
