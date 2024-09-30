package dao

import (
	"errors"
	"gcalsync/gophers/clients/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

var DB *gorm.DB

// TODO: Move migrations and seed to a seperate command

// InitDB ...
func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	log := logger.GetInstance()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		log.Error(nil, "Failed to connect to the database: %v", err)
		return err
	}
	err = DB.AutoMigrate(&User{}, &Calendar{}, &Event{}, &UserToken{}, &Watch{})

	if err != nil {
		log.Error(nil, "Failed to migrate database schema: %v", err)
		return err
	}

	err = seedUsers(DB)
	if err != nil {
		log.Error(nil, "Failed to migrate database schema: %v", err)
		return err
	}

	log.Println("Database migration completed.")
	return nil
}

func seedUsers(DB *gorm.DB) error {

	return DB.Transaction(func(tx *gorm.DB) error {
		users := []User{
			{Email: "jane@test.com", PasswordHash: "test"},
			{Email: "john@test.com", PasswordHash: "test"},
		}

		for _, user := range users {
			if err := upsertUser(DB, &user); err != nil {
				logger.GetInstance().Info(nil, "Error seeding user: %v\n", err)
				return err
			}
		}

		return nil
	})
}

func upsertUser(DB *gorm.DB, user *User) error {
	err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(user).Error

	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return err
	}
	return nil
}
