package connection

import (
	"os"
	// TODO: init driver?
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDevConnection() (*gorm.DB, error) {
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASS") + " dbname=" + os.Getenv("DB_NAME") +
		" port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
