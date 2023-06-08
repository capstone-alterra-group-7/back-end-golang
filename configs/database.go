package configs

import (
	"back-end-golang/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load the Asia/Jakarta location
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle the error
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	dbConn = dbConn.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return time.Now().In(location)
		},
	})

	return dbConn, nil
}
func AccountSeeder(db *gorm.DB) error {
	password := "$2a$10$QXBNiEWub5z3TX5LFewSy.atj0iARk1vCZDgzRQTDp5xOQopj4WRW"
	users := []models.User{
		{
			FullName:       "Admin",
			Email:          "admin@gmail.com",
			Password:       password,
			PhoneNumber:    "08523884322",
			ProfilePicture: "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg",
			Citizen:        "Indonesia",
			Role:           "admin",
		},
		{
			FullName:       "User",
			Email:          "user@gmail.com",
			Password:       password,
			PhoneNumber:    "08523884322",
			ProfilePicture: "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg",
			Citizen:        "Indonesia",
			Role:           "user",
		},
	}

	for _, user := range users {

		// Check if data already exists
		var count int64
		if err := db.Model(&models.User{}).Where(&user).Count(&count).Error; err != nil {
			return err
		}

		// If data exists, skip seeding
		if count > 0 {
			continue
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func TrainSeatSeeder(db *gorm.DB) error {
	classTypes := []string{"Ekonomi", "Bisnis", "Eksekutif"}

	for _, class := range classTypes {
		prefixes := []string{"A", "B", "C", "D"}

		if class == "Ekonomi" {
			prefixes = append(prefixes, "E")
		}

		for _, prefix := range prefixes {
			startIndex := 1
			endIndex := 24

			if prefix == "D" && class == "Ekonomi" || prefix == "E" && class == "Ekonomi" {
				startIndex = 3
			}

			if prefix == "C" && class == "Ekonomi" {
				startIndex = 4
				endIndex = 21
			}

			if prefix == "A" && class == "Ekonomi" || prefix == "B" && class == "Ekonomi" {
				endIndex = 22
			}

			if prefix == "C" && class == "Bisnis" || prefix == "D" && class == "Bisnis" {
				startIndex = 2
			}

			if prefix == "A" && class == "Bisnis" || prefix == "B" && class == "Bisnis" {
				endIndex = 16
			}

			if prefix == "C" && class == "Bisnis" || prefix == "D" && class == "Bisnis" {
				startIndex = 2
				endIndex = 17
			}

			if prefix == "A" && class == "Eksekutif" {
				endIndex = 12
			}

			if prefix == "B" && class == "Eksekutif" || prefix == "C" && class == "Eksekutif" || prefix == "D" && class == "Eksekutif" {
				endIndex = 13
			}

			if prefix == "D" && class == "Eksekutif" {
				startIndex = 2
			}

			for i := startIndex; i <= endIndex; i++ {
				seat := models.TrainSeat{
					Class: class,
					Name:  fmt.Sprintf("%s%d", prefix, i),
				}

				// Check if data already exists
				var count int64
				if err := db.Model(&models.TrainSeat{}).Where(&seat).Count(&count).Error; err != nil {
					return err
				}

				// If data exists, skip seeding
				if count > 0 {
					continue
				}

				if err := db.Create(&seat).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Station{},
		&models.Train{},
		&models.TrainStation{},
		&models.TrainCarriage{},
		&models.TrainSeat{},
		&models.TravelerDetail{},
		&models.TicketOrder{},
		&models.TicketTravelerDetail{},
		&models.Article{},
		&models.HistorySearch{},
		&models.Payment{},
		&models.Hotel{},
		&models.HotelImage{},
		&models.HotelFacilities{},
		&models.HotelPolicies{},
		&models.HotelRoom{},
		&models.HotelRoomImage{},
		&models.HotelRoomFacilities{},
	)
}
