package usecases

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock in source shop")
	ErrSameShopTransfer  = errors.New("cannot transfer stock to the same shop")
)

type StockTransferService interface {
	GetStockTransfers(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.StockTransfer, int64, error)
	GetStockTransferByID(id uuid.UUID) (*entities.StockTransfer, error)
	GetStockTransfersByShopID(shopID uuid.UUID) ([]entities.StockTransfer, error)
	CreateStockTransfer(transfer *entities.StockTransfer) error
	UpdateStockTransfer(transfer *entities.StockTransfer) error
	DeleteStockTransfer(id uuid.UUID) error
}

type stockTransferService struct {
	stockTransferRepo repository.StockTransferRepository
	inventoryRepo     repository.InventoryRepository
}

func NewStockTransferService(stockTransferRepo repository.StockTransferRepository, inventoryRepo repository.InventoryRepository) StockTransferService {
	return &stockTransferService{
		stockTransferRepo: stockTransferRepo,
		inventoryRepo:     inventoryRepo,
	}
}

func (s *stockTransferService) GetStockTransfers(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.StockTransfer, int64, error) {
	return s.stockTransferRepo.GetStockTransfersWithFilters(filters, sorts, page, pageSize)
}

func (s *stockTransferService) GetStockTransferByID(id uuid.UUID) (*entities.StockTransfer, error) {
	return s.stockTransferRepo.GetByID(id)
}

func (s *stockTransferService) GetStockTransfersByShopID(shopID uuid.UUID) ([]entities.StockTransfer, error) {
	// Use current time as end date to get all transfers up to now
	return s.stockTransferRepo.GetTransfersByShop(shopID, time.Time{}, time.Now())
}

func (s *stockTransferService) CreateStockTransfer(transfer *entities.StockTransfer) error {
	// Check if source and destination shops are different
	if transfer.FromShopID == transfer.ToShopID {
		return ErrSameShopTransfer
	}

	// Check if source shop has sufficient stock
	sourceInventory, err := s.inventoryRepo.GetByProductAndShop(transfer.ProductID, transfer.FromShopID)
	if err != nil {
		return err
	}

	if sourceInventory.Quantity < transfer.Quantity {
		return ErrInsufficientStock
	}

	// Create the transfer record
	if err := s.stockTransferRepo.Create(transfer); err != nil {
		return err
	}

	// Update source shop inventory
	sourceInventory.Quantity -= transfer.Quantity
	if err := s.inventoryRepo.UpdateStock(transfer.ProductID, transfer.FromShopID, -transfer.Quantity); err != nil {
		return err
	}

	// Update destination shop inventory
	if err := s.inventoryRepo.UpdateStock(transfer.ProductID, transfer.ToShopID, transfer.Quantity); err != nil {
		return err
	}

	return nil
}

func (s *stockTransferService) UpdateStockTransfer(transfer *entities.StockTransfer) error {
	// Get existing transfer to compare changes
	existing, err := s.stockTransferRepo.GetByID(transfer.ID)
	if err != nil {
		return err
	}

	// If quantity changed, update inventories
	if existing.Quantity != transfer.Quantity {
		// Update source shop inventory
		if err := s.inventoryRepo.UpdateStock(transfer.ProductID, transfer.FromShopID, existing.Quantity-transfer.Quantity); err != nil {
			return err
		}

		// Update destination shop inventory
		if err := s.inventoryRepo.UpdateStock(transfer.ProductID, transfer.ToShopID, transfer.Quantity-existing.Quantity); err != nil {
			return err
		}
	}

	return s.stockTransferRepo.Create(transfer) // Use Create for update as it handles both insert and update
}

func (s *stockTransferService) DeleteStockTransfer(id uuid.UUID) error {
	// Get transfer to update inventories
	transfer, err := s.stockTransferRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Update source shop inventory
	if err := s.inventoryRepo.UpdateStock(transfer.ProductID, transfer.FromShopID, transfer.Quantity); err != nil {
		return err
	}

	// Update destination shop inventory
	if err := s.inventoryRepo.UpdateStock(transfer.ProductID, transfer.ToShopID, -transfer.Quantity); err != nil {
		return err
	}

	return s.stockTransferRepo.Create(transfer) // Use Create for soft delete
}
