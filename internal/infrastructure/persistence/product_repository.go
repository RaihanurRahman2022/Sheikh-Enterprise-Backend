package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"gorm.io/gorm"
)

type ProductRepository interface {
	BaseRepository[entities.Product]
	GetProductsWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.Product, int64, error)
	BulkCreate(products []entities.Product) error
}

type productRepository struct {
	BaseRepositoryImpl[entities.Product]
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.Product]{DB: db},
	}
}

func (r *productRepository) GetProductsWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.Product, int64, error) {
	var products []entities.Product
	var total int64

	query := r.DB.Model(&entities.Product{}).Where("is_marked_to_delete = ?", false)

	// Apply filters
	for field, value := range filters {
		switch field {
		case "code", "name", "style", "master_category", "sub_category", "color", "size", "sales_type":
			query = query.Where(field+" = ?", value)
		case "min_purchase_price":
			query = query.Where("purchase_price >= ?", value)
		case "max_purchase_price":
			query = query.Where("purchase_price <= ?", value)
		case "min_sales_price":
			query = query.Where("sales_price >= ?", value)
		case "max_sales_price":
			query = query.Where("sales_price <= ?", value)
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
	err = query.Offset(offset).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) BulkCreate(products []entities.Product) error {
	return r.DB.Create(&products).Error
}
