package persistence

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseRepository defines common CRUD operations
type BaseRepository[T any] interface {
	GetByID(id uuid.UUID) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id uuid.UUID) error
	List(page, pageSize int) ([]T, int64, error)
}

// BaseRepositoryImpl implements BaseRepository
type BaseRepositoryImpl[T any] struct {
	DB *gorm.DB
}

// GetByID retrieves an entity by its ID
func (r *BaseRepositoryImpl[T]) GetByID(id uuid.UUID) (*T, error) {
	var entity T
	err := r.DB.Where("id = ? AND is_marked_to_delete = ?", id, false).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Create creates a new entity
func (r *BaseRepositoryImpl[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

// Update updates an existing entity
func (r *BaseRepositoryImpl[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

// Delete soft deletes an entity
func (r *BaseRepositoryImpl[T]) Delete(id uuid.UUID) error {
	return r.DB.Model(new(T)).Where("id = ?", id).Update("is_marked_to_delete", true).Error
}

// List retrieves a paginated list of entities
func (r *BaseRepositoryImpl[T]) List(page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	query := r.DB.Model(new(T)).Where("is_marked_to_delete = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
