package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository[entities.User]
	GetByUsername(username string) (*entities.User, error)
	UpdatePassword(id uuid.UUID, hashedPassword string) error
	UpdateDetails(user *entities.User) error
}

type userRepository struct {
	BaseRepositoryImpl[entities.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.User]{DB: db},
	}
}

func (r *userRepository) GetByUsername(username string) (*entities.User, error) {
	var user entities.User
	err := r.DB.Where("username = ? AND is_marked_to_delete = ?", username, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(id uuid.UUID, hashedPassword string) error {
	return r.DB.Model(&entities.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func (r *userRepository) UpdateDetails(user *entities.User) error {
	return r.DB.Model(&entities.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"phone":      user.Phone,
			"shop_id":    user.ShopID,
		}).Error
}
