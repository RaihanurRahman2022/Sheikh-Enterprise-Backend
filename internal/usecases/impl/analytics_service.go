package usecases

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"time"

	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetStockTransferAnalytics(startDate, endDate time.Time, shopID *uuid.UUID) (map[string]interface{}, error)
	GetInventoryAnalytics(shopID *uuid.UUID) (map[string]interface{}, error)
}

type analyticsService struct {
	stockTransferRepo repository.StockTransferRepository
	inventoryRepo     repository.InventoryRepository
}

func NewAnalyticsService(stockTransferRepo repository.StockTransferRepository, inventoryRepo repository.InventoryRepository) AnalyticsService {
	return &analyticsService{
		stockTransferRepo: stockTransferRepo,
		inventoryRepo:     inventoryRepo,
	}
}

func (s *analyticsService) GetStockTransferAnalytics(startDate, endDate time.Time, shopID *uuid.UUID) (map[string]interface{}, error) {
	filters := map[string]interface{}{
		"transfer_datetime >= ?": startDate,
		"transfer_datetime <= ?": endDate,
	}

	if shopID != nil {
		filters["(to_shop_id = ? OR from_shop_id = ?)"] = []interface{}{*shopID, *shopID}
	}

	transfers, _, err := s.stockTransferRepo.GetStockTransfersWithFilters(filters, []string{"transfer_datetime DESC"}, 1, 1000)
	if err != nil {
		return nil, err
	}

	// Calculate analytics
	totalTransfers := len(transfers)
	totalQuantity := 0
	shopTransfers := make(map[uuid.UUID]int)
	productTransfers := make(map[uuid.UUID]int)

	for _, transfer := range transfers {
		totalQuantity += transfer.Quantity
		shopTransfers[transfer.FromShopID]++
		shopTransfers[transfer.ToShopID]++
		productTransfers[transfer.ProductID]++
	}

	return map[string]interface{}{
		"total_transfers":    totalTransfers,
		"total_quantity":     totalQuantity,
		"shop_transfers":     shopTransfers,
		"product_transfers":  productTransfers,
		"transfers_by_date":  s.groupTransfersByDate(transfers),
		"transfers_by_month": s.groupTransfersByMonth(transfers),
	}, nil
}

func (s *analyticsService) GetInventoryAnalytics(shopID *uuid.UUID) (map[string]interface{}, error) {
	var inventories []entities.Inventory
	var err error

	if shopID != nil {
		inventories, err = s.inventoryRepo.GetInventoryByShopID(*shopID)
	} else {
		inventories, _, err = s.inventoryRepo.GetInventoryWithFilters(nil, "quantity", "DESC", 1, 1000)
	}

	if err != nil {
		return nil, err
	}

	// Calculate analytics
	totalProducts := len(inventories)
	totalQuantity := 0
	lowStockProducts := 0
	productDistribution := make(map[string]int)

	for _, inventory := range inventories {
		totalQuantity += inventory.Quantity
		if inventory.Quantity < 10 { // Assuming 10 is the low stock threshold
			lowStockProducts++
		}
		productDistribution[inventory.Product.MasterCategory]++
	}

	return map[string]interface{}{
		"total_products":       totalProducts,
		"total_quantity":       totalQuantity,
		"low_stock_products":   lowStockProducts,
		"product_distribution": productDistribution,
		"inventory_by_shop":    s.groupInventoryByShop(inventories),
	}, nil
}

// Helper functions for grouping data
func (s *analyticsService) groupTransfersByDate(transfers []entities.StockTransfer) map[string]int {
	result := make(map[string]int)
	for _, transfer := range transfers {
		date := transfer.TransferDateTime.Format("2006-01-02")
		result[date] += transfer.Quantity
	}
	return result
}

func (s *analyticsService) groupTransfersByMonth(transfers []entities.StockTransfer) map[string]int {
	result := make(map[string]int)
	for _, transfer := range transfers {
		month := transfer.TransferDateTime.Format("2006-01")
		result[month] += transfer.Quantity
	}
	return result
}

func (s *analyticsService) groupInventoryByShop(inventories []entities.Inventory) map[uuid.UUID]int {
	result := make(map[uuid.UUID]int)
	for _, inventory := range inventories {
		result[inventory.ShopID] += inventory.Quantity
	}
	return result
}
