package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the global DB variable using the DSN.
func InitDB() error {
	// Build DSN from environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println(dsn)

	var err error
	for i := 0; i < 10; i++ { // retry 10 times
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return nil
		}
		fmt.Println("Waiting for MySQL to be ready...")
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("failed to connect to database after retries: %w", err)

	// Connect using GORM
	// database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to database: %w", err)
	// }

	// // Set the global DB variable
	// DB = database
	// return nil
}
