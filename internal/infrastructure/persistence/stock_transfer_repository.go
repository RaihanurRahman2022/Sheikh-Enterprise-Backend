package persistence

import (
	"errors"
	"time"

	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StockTransferRepository defines the interface for stock transfer operations
type StockTransferRepository interface {
	Create(stockTransfer *entities.StockTransfer) error
	Update(stockTransfer *entities.StockTransfer) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*entities.StockTransfer, error)
	GetByShopID(shopID uuid.UUID) ([]entities.StockTransfer, error)
	GetByProductID(productID uuid.UUID) ([]entities.StockTransfer, error)
	GetByStatus(status entities.StockTransferStatus) ([]entities.StockTransfer, error)
	GetHistory(stockTransferID uuid.UUID) ([]entities.StockTransferHistory, error)
	AddHistory(history *entities.StockTransferHistory) error
	GetStockTransfersWithFilters(filters map[string]interface{}, orderBy []string, page, pageSize int) ([]entities.StockTransfer, int64, error)
	GetTransfersByShop(shopID uuid.UUID, startDate, endDate time.Time) ([]entities.StockTransfer, error)
}

type stockTransferRepository struct {
	db *gorm.DB
}

// NewStockTransferRepository creates a new instance of StockTransferRepository
func NewStockTransferRepository(db *gorm.DB) StockTransferRepository {
	return &stockTransferRepository{db: db}
}

func (r *stockTransferRepository) Create(stockTransfer *entities.StockTransfer) error {
	return r.db.Create(stockTransfer).Error
}

func (r *stockTransferRepository) Update(stockTransfer *entities.StockTransfer) error {
	return r.db.Save(stockTransfer).Error
}

func (r *stockTransferRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.StockTransfer{}, id).Error
}

func (r *stockTransferRepository) GetByID(id uuid.UUID) (*entities.StockTransfer, error) {
	var stockTransfer entities.StockTransfer
	err := r.db.Preload("Product").
		Preload("ToShop").
		Preload("FromShop").
		Preload("TransferredByUser").
		First(&stockTransfer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stockTransfer, nil
}

func (r *stockTransferRepository) GetByShopID(shopID uuid.UUID) ([]entities.StockTransfer, error) {
	var stockTransfers []entities.StockTransfer
	err := r.db.Preload("Product").
		Preload("ToShop").
		Preload("FromShop").
		Preload("TransferredByUser").
		Where("to_shop_id = ? OR from_shop_id = ?", shopID, shopID).
		Find(&stockTransfers).Error
	return stockTransfers, err
}

func (r *stockTransferRepository) GetByProductID(productID uuid.UUID) ([]entities.StockTransfer, error) {
	var stockTransfers []entities.StockTransfer
	err := r.db.Preload("Product").
		Preload("ToShop").
		Preload("FromShop").
		Preload("TransferredByUser").
		Where("product_id = ?", productID).
		Find(&stockTransfers).Error
	return stockTransfers, err
}

func (r *stockTransferRepository) GetByStatus(status entities.StockTransferStatus) ([]entities.StockTransfer, error) {
	var stockTransfers []entities.StockTransfer
	err := r.db.Preload("Product").
		Preload("ToShop").
		Preload("FromShop").
		Preload("TransferredByUser").
		Where("status = ?", status).
		Find(&stockTransfers).Error
	return stockTransfers, err
}

func (r *stockTransferRepository) GetHistory(stockTransferID uuid.UUID) ([]entities.StockTransferHistory, error) {
	var history []entities.StockTransferHistory
	err := r.db.Where("stock_transfer_id = ?", stockTransferID).
		Order("created_at DESC").
		Find(&history).Error
	return history, err
}

func (r *stockTransferRepository) AddHistory(history *entities.StockTransferHistory) error {
	return r.db.Create(history).Error
}

func (r *stockTransferRepository) GetStockTransfersWithFilters(filters map[string]interface{}, orderBy []string, page, pageSize int) ([]entities.StockTransfer, int64, error) {
	var stockTransfers []entities.StockTransfer
	var total int64

	query := r.db.Model(&entities.StockTransfer{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key, value)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	for _, order := range orderBy {
		query = query.Order(order)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Preload("Product").
		Preload("ToShop").
		Preload("FromShop").
		Preload("TransferredByUser").
		Offset(offset).
		Limit(pageSize).
		Find(&stockTransfers).Error

	return stockTransfers, total, err
}

func (r *stockTransferRepository) GetTransfersByShop(shopID uuid.UUID, startDate, endDate time.Time) ([]entities.StockTransfer, error) {
	var stockTransfers []entities.StockTransfer
	query := r.db.Where("(from_shop_id = ? OR to_shop_id = ?)", shopID, shopID)

	if !startDate.IsZero() {
		query = query.Where("transfer_datetime >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("transfer_datetime <= ?", endDate)
	}

	err := query.Preload("Product").
		Preload("ToShop").
		Preload("FromShop").
		Preload("TransferredByUser").
		Find(&stockTransfers).Error
	return stockTransfers, err
}
