
package usecases

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	"github.com/google/uuid"
)

type StockTransferService interface {
	GetStockTransfers(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.StockTransfer, int64, error)
	GetStockTransferByID(id uuid.UUID) (*entities.StockTransfer, error)
	CreateStockTransfer(transfer *entities.StockTransfer) error
	DeleteStockTransfer(id uuid.UUID) error
}

type stockTransferService struct {
	stockTransferRepo repository.StockTransferRepository
}

func NewStockTransferService(stockTransferRepo repository.StockTransferRepository) StockTransferService {
	return &stockTransferService{
		stockTransferRepo: stockTransferRepo,
	}
}

func (s *stockTransferService) GetStockTransfers(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.StockTransfer, int64, error) {
	return s.stockTransferRepo.GetStockTransfersWithFilters(filters, sorts, page, pageSize)
}

func (s *stockTransferService) GetStockTransferByID(id uuid.UUID) (*entities.StockTransfer, error) {
	return s.stockTransferRepo.GetByID(id)
}

func (s *stockTransferService) CreateStockTransfer(transfer *entities.StockTransfer) error {
	return s.stockTransferRepo.Create(transfer)
}

func (s *stockTransferService) DeleteStockTransfer(id uuid.UUID) error {
	return s.stockTransferRepo.Delete(id)
}
