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

var DBAkademik *gorm.DB
var SqlDBAkademik *sql.DB

var DBDigiclass *gorm.DB
var SqlDBDigiclass *sql.DB

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

	cfgakademik := config.LoadConfigAkademik()

	dsnakademik := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfgakademik.DBHost, cfgakademik.DBUser, cfgakademik.DBPassword, cfgakademik.DBName, cfgakademik.DBPort, cfgakademik.SSLMode)

	var errakademik error
	DBAkademik, errakademik = gorm.Open(postgres.Open(dsnakademik), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "akademik.", // default schema
			SingularTable: true,
		},
	})
	if errakademik != nil {
		log.Fatal("Failed to connect to database:", errakademik)
	}

	// Create schema if it doesn't exist
	if errakademik := DB.Exec("CREATE SCHEMA IF NOT EXISTS akademik").Error; errakademik != nil {
		log.Fatal("Failed to create schema:", errakademik)
	}

	// Set search path to include the akademik schema
	if errakademik := DB.Exec("SET search_path TO akademik, public").Error; errakademik != nil {
		log.Fatal("Failed to set search path:", errakademik)
	}

	// Get the underlying *sql.DB from *gorm.DB
	SqlDBAkademik, err = DBAkademik.DB()
	if err != nil {
		log.Fatal("Failed to get *sql.DB from *gorm.DB:", errakademik)
	}

	// You can now use SqlDB if needed, for example, to configure connection pool settings
	SqlDBAkademik.SetMaxIdleConns(10)
	SqlDBAkademik.SetMaxOpenConns(100)
	SqlDBAkademik.SetConnMaxLifetime(time.Hour)

	log.Println("Database Akademik connection successful")
}

func ConnectDigiclass() {
	cfg := config.LoadConfigDigiclass()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.SSLMode)

	var err error
	DBDigiclass, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
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
	SqlDBDigiclass, err = DBDigiclass.DB()
	if err != nil {
		log.Printf("⚠️ Warning: Failed to get *sql.DB from *gorm.DB: %v\n", err)
		return
	}

	// Set database connection pool settings
	SqlDBDigiclass.SetMaxIdleConns(10)
	SqlDBDigiclass.SetMaxOpenConns(100)
	SqlDBDigiclass.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database Digiclass connection successful")

}
