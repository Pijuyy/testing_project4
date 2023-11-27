package config

import (
	"fmt"
	"os"

	"github.com/Pijuyy/testing_project4/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var (
// 	host     = os.Getenv("PGHOST")
// 	user     = os.Getenv("PGUSER")
// 	password = os.Getenv("PGPASSWORD")
// 	dbPort   = os.Getenv("PGPORT")
// 	dbname   = os.Getenv("PGDATABASE")
// 	// host       = "localhost"
// 	// user       = "postgres"
// 	// password   = "postgres890"
// 	// dbPort     = "5432"
// 	// dbname     = "mygram"
// 	// debug_mode = os.Getenv("DEBUG_MODE")
// 	DB  *gorm.DB
// 	err error
// )

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres dbname=baru password=Hafidzurr1 sslmode=disable" // Lokal DB config
		// dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, dbPort)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database. Error:", err)
		return nil, err
	}

	// Call Migrate function to perform database migrations
	if err := Migrate(db); err != nil {
		fmt.Println("Failed to perform database migrations. Error:", err)
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.TransactionHistory{})
}
