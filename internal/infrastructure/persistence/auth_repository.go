package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByID(id uuid.UUID) (*entities.User, error)
	UpdatePassword(id uuid.UUID, hashedPassword string) error
	Create(user *entities.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) GetUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("username = ? AND is_marked_to_delete = ?", username, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("id = ? AND is_marked_to_delete = ?", id, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdatePassword(id uuid.UUID, hashedPassword string) error {
	return r.db.Model(&entities.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func (r *authRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}
