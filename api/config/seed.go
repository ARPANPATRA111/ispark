package config

import (
	"log"
	"os"
	"strconv"
)

func SeedDevData() {
	enabled, err := strconv.ParseBool(os.Getenv("SEED_DEV_DATA"))
	if err != nil || !enabled {
		log.Println("SEED_DEV_DATA is not enabled, skipping demo data seeding")
		return
	}

	SeedDefaultAdmin()
	SeedDefaultStudents()
}
