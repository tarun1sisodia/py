package models

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	var err error

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return err
	}

	log.Println("Connected to database")

	// Auto-migrate models
	err = DB.AutoMigrate(
		&User{},
		&AttendanceSession{},
		&AttendanceRecord{},
		&OTPVerification{},
		&AuditLog{},
		&Admin{},
		&AdminAuditLog{},
		&Report{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
