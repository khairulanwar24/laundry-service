package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"laundry-service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var SqlDB *sql.DB

func Connect() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "public.", // default schema
			SingularTable: false,
		},
	})
	if err != nil {
		log.Printf("⚠️ Warning: Failed to connect to database: %v\n", err)
		return // Jangan hentikan program, cukup log error
	}

	// Get the underlying *sql.DB from *gorm.DB
	SqlDB, err = DB.DB()
	if err != nil {
		log.Printf("⚠️ Warning: Failed to get *sql.DB from *gorm.DB: %v\n", err)
		return
	}

	// Set database connection pool settings
	SqlDB.SetMaxIdleConns(10)
	SqlDB.SetMaxOpenConns(100)
	SqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection successful")
}
