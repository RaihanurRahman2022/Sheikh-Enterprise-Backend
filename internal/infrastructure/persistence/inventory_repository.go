package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetByProductAndShop(productID, shopID uuid.UUID) (*entities.Inventory, error)
	GetInventoryByShopID(shopID uuid.UUID) ([]entities.Inventory, error)
	UpdateStock(productID, shopID uuid.UUID, quantity int) error
	TransferStock(fromShopID, toShopID, productID uuid.UUID, quantity int) error
	GetLowStockItems(shopID uuid.UUID, threshold int) ([]entities.Inventory, error)
	GetInventoryWithFilters(filters map[string]interface{}, sortBy, sortOrder string, page, pageSize int) ([]entities.Inventory, int64, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{
		db: db,
	}
}

func (r *inventoryRepository) GetByProductAndShop(productID, shopID uuid.UUID) (*entities.Inventory, error) {
	var inventory entities.Inventory
	err := r.db.Where("product_id = ? AND shop_id = ?", productID, shopID).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetInventoryByShopID(shopID uuid.UUID) ([]entities.Inventory, error) {
	var inventory []entities.Inventory
	err := r.db.Where("shop_id = ?", shopID).Find(&inventory).Error
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *inventoryRepository) UpdateStock(productID, shopID uuid.UUID, quantity int) error {
	return r.db.Exec(`
		INSERT INTO inventory (product_id, shop_id, quantity)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		quantity = quantity + ?
	`, productID, shopID, quantity, quantity).Error
}

func (r *inventoryRepository) TransferStock(fromShopID, toShopID, productID uuid.UUID, quantity int) error {
	tx := r.db.Begin()

	// Decrease stock in source shop
	if err := tx.Exec(`
		UPDATE inventory
		SET quantity = quantity - ?
		WHERE shop_id = ? AND product_id = ? AND quantity >= ?
	`, quantity, fromShopID, productID, quantity).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Increase stock in destination shop
	if err := tx.Exec(`
		INSERT INTO inventory (product_id, shop_id, quantity)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		quantity = quantity + ?
	`, productID, toShopID, quantity, quantity).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *inventoryRepository) GetLowStockItems(shopID uuid.UUID, threshold int) ([]entities.Inventory, error) {
	var inventory []entities.Inventory
	err := r.db.
		Joins("JOIN products ON inventory.product_id = products.product_id").
		Where("inventory.shop_id = ? AND inventory.quantity <= ?", shopID, threshold).
		Find(&inventory).Error
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *inventoryRepository) GetInventoryWithFilters(filters map[string]interface{}, sortBy, sortOrder string, page, pageSize int) ([]entities.Inventory, int64, error) {
	var inventories []entities.Inventory
	var total int64

	query := r.db.Model(&entities.Inventory{})

	// Apply filters
	for key, value := range filters {
		if value != nil {
			query = query.Where(key, value)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sortBy != "" {
		query = query.Order(sortBy + " " + sortOrder)
	}

	// Apply pagination
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	// Execute query
	if err := query.Find(&inventories).Error; err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}
