package main

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/config"
	"ssl-custom-api/internal/providers/postgres"
)

func Run(db *gorm.DB, password string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var user models.User

		err := tx.
			Where("email = ?", "admin@ssl.com").
			First(&user).Error

		if err == nil {
			fmt.Println("admin already exists")
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check admin user: %w", err)
		}

		user = models.User{
			Username: "admin",
			Email:    "admin@ssl.com",
			Password: password,
		}

		if err := auth.NewUserRepository(tx, nil).CreateUser(&user); err != nil {
			return fmt.Errorf("create admin: %w", err)
		}

		fmt.Println("admin created successfully")

		return nil
	})
}
func readPassword() (string, error) {
	var password string
	fmt.Print("Enter admin password: ")
	_, err := fmt.Scanln(&password)
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}

	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	return password, nil
}

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := postgres.Config{
		Host:     config.PGHost,
		Port:     config.PGPort,
		User:     config.PGUser,
		Password: config.PGPassword,
		Database: config.PGDatabase,
		SSLMode:  config.PGSSLMode,
	}

	provider, err := postgres.NewProvider(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer provider.Close()

	db := provider.DB()

	password, err := readPassword()
	if err != nil {
		log.Fatalf("failed to read password: %v", err)
	}

	if err := Run(db, password); err != nil {
		log.Fatalf("failed to run seed: %v", err)
	}
}
