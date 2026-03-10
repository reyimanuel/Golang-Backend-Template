package user

import (
	"errors"

	"github.com/reyimanuel/template/internal/migration"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetByEmail(email string) (*migration.User, error) {
	var user migration.User
	if err := r.DB.Preload("Roles").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateAccount(user *migration.User) error {
	return r.DB.Create(user).Error
}
