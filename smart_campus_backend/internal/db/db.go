package db

import (
	"database/sql"
	"fmt"
	"smart_campus_backend/config"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &Database{db}, nil
}

func (db *Database) NewUserRepository() *UserRepository {
	return &UserRepository{db: db}
}

func (db *Database) NewSessionRepository() *SessionRepository {
	return &SessionRepository{db: db}
}

func (db *Database) NewCourseRepository() *CourseRepository {
	return &CourseRepository{db: db}
}

func (db *Database) NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (db *Database) NewAttendanceRepository() *AttendanceRepository {
	return &AttendanceRepository{db: db}
}
