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
	"ssl-custom-api/internal/utils"
)

func Run(db *gorm.DB, password string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var role models.Role

		err := tx.
			Where("name = ?", "admin").
			First(&role).Error

		if err == nil {
			fmt.Println("admin role already exists")
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			role = models.Role{
				ID:      utils.NewV7ID(),
				Name:    "admin",
				IsAdmin: true,
			}

			if err := tx.Create(&role).Error; err != nil {
				return fmt.Errorf("create admin role: %w", err)
			}

			fmt.Println("admin role created successfully")
		} else {
			return fmt.Errorf("check admin role: %w", err)
		}
		var user models.User

		err_user := tx.
			Where("email = ?", "admin@ssl.com").
			First(&user).Error

		if err_user == nil {
			fmt.Println("admin already exists")
			return nil
		}

		if !errors.Is(err_user, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check admin user: %w", err_user)
		}

		user = models.User{
			ID:       utils.NewV7ID(),
			Username: "admin",
			Email:    "admin@ssl.com",
			Password: password,
			RoleID:   role.ID,
		}

		if err_user := auth.NewUserRepository(tx, nil).CreateUser(&user); err_user != nil {
			return fmt.Errorf("create admin: %w", err_user)
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
