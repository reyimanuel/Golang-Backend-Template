package migration

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/reyimanuel/template/internal/infrastructures/pkg/helpers"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, force bool) error {
	host := os.Getenv("DB_HOST")
	if !force && !strings.EqualFold(host, "localhost") && host != "127.0.0.1" {
		return fmt.Errorf("seeding blocked in production (DB_HOST=%s), use --force", host)
	}

	const (
		defaultUsername = "admin"
		defaultEmail    = "admin@example.com"
		defaultPassword = "admin123"
	)

	var existingUser User
	if err := db.Where("email = ? OR username = ?", defaultEmail, defaultUsername).First(&existingUser).Error; err == nil {
		fmt.Printf("Seeder account already exists: %s\n", existingUser.Email)
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed checking seeder account: %w", err)
	}

	hashedPassword, err := helpers.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed hashing seeder password: %w", err)
	}

	seedUser := User{
		Username: defaultUsername,
		Email:    defaultEmail,
		Password: hashedPassword,
	}

	if err := db.Create(&seedUser).Error; err != nil {
		return fmt.Errorf("failed creating seeder account: %w", err)
	}

	fmt.Printf("Seeder account created: %s\n", seedUser.Email)
	return nil
}
