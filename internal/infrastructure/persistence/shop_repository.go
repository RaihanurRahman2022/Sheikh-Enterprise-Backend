package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"gorm.io/gorm"
)

type ShopRepository interface {
	BaseRepository[entities.Shop]
	GetShopsWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.Shop, int64, error)
	GetShopsByCompanyID(companyID string) ([]entities.Shop, error)
}

type shopRepository struct {
	BaseRepositoryImpl[entities.Shop]
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.Shop]{DB: db},
	}
}

func (r *shopRepository) GetShopsWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.Shop, int64, error) {
	var shops []entities.Shop
	var total int64

	query := r.DB.Model(&entities.Shop{}).Where("is_marked_to_delete = ?", false)

	// Apply filters
	for field, value := range filters {
		switch field {
		case "name", "phone", "email", "manager_name", "manager_phone":
			query = query.Where(field+" LIKE ?", "%"+value.(string)+"%")
		case "company_id":
			query = query.Where("company_id = ?", value)
		}
	}

	// Count total before pagination
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting
	for _, sort := range sorts {
		query = query.Order(sort)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&shops).Error
	if err != nil {
		return nil, 0, err
	}

	return shops, total, nil
}

func (r *shopRepository) GetShopsByCompanyID(companyID string) ([]entities.Shop, error) {
	var shops []entities.Shop
	err := r.DB.Where("company_id = ? AND is_marked_to_delete = ?", companyID, false).Find(&shops).Error
	if err != nil {
		return nil, err
	}
	return shops, nil
}
