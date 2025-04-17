package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"gorm.io/gorm"
)

type SupplierRepository interface {
	BaseRepository[entities.Supplier]
	GetSuppliersWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.Supplier, int64, error)
}

type supplierRepository struct {
	BaseRepositoryImpl[entities.Supplier]
}

func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.Supplier]{DB: db},
	}
}

func (r *supplierRepository) GetSuppliersWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.Supplier, int64, error) {
	var suppliers []entities.Supplier
	var total int64

	query := r.DB.Model(&entities.Supplier{}).Where("is_marked_to_delete = ?", false)

	// Apply filters
	for field, value := range filters {
		switch field {
		case "name", "phone", "email":
			query = query.Where(field+" LIKE ?", "%"+value.(string)+"%")
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
	err = query.Offset(offset).Limit(pageSize).Find(&suppliers).Error
	if err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}
