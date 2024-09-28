package dao

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	err = DB.AutoMigrate(&User{}, &Calendar{}, &Event{}, &UserTokens{})

	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	err = seedUsers(DB)
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	log.Println("Database migration completed.")
}

func seedUsers(DB *gorm.DB) error {
	users := []User{
		{Email: "jane@test.com", PasswordHash: "test"},
		{Email: "john@test.com", PasswordHash: "test"},
	}

	for _, user := range users {
		if err := saveUserIgnoringUniqueConstraint(DB, &user); err != nil {
			log.Printf("Error seeding user: %v\n", err)
			return err
		}
	}

	return nil
}

func saveUserIgnoringUniqueConstraint(DB *gorm.DB, user *User) error {
	err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(user).Error

	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return err
	}
	return nil
}
