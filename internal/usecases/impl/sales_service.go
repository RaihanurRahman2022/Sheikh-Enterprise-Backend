package usecases

import (
	"fmt"
	"time"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type SalesService interface {
	GetSales(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.SalesInvoice, int64, error)
	GetSaleByID(id uuid.UUID) (*entities.SalesInvoice, error)
	CreateSale(sale *entities.SalesInvoice) error
	DeleteSale(id uuid.UUID) error
	ExportToExcel(filters map[string]interface{}, sorts []string) (*excelize.File, error)
	GetAnalytics(shopID *uuid.UUID) (*repository.SalesAnalytics, error)
	GetLast7DaysSales(shopID *uuid.UUID) ([]repository.DailySales, error)
}

type salesService struct {
	salesRepo repository.SalesRepository
}

func NewSalesService(salesRepo repository.SalesRepository) SalesService {
	return &salesService{
		salesRepo: salesRepo,
	}
}

func (s *salesService) GetSales(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.SalesInvoice, int64, error) {
	return s.salesRepo.GetSalesWithFilters(filters, sorts, page, pageSize)
}

func (s *salesService) GetSaleByID(id uuid.UUID) (*entities.SalesInvoice, error) {
	return s.salesRepo.GetByID(id)
}

func (s *salesService) CreateSale(sale *entities.SalesInvoice) error {
	// Calculate totals
	var total float64
	for i := range sale.SalesDetails {
		detail := &sale.SalesDetails[i]
		detail.Subtotal = detail.SalesPrice * float64(detail.Quantity)
		total += detail.Subtotal
	}
	sale.Total = total - sale.Discount

	return s.salesRepo.Create(sale)
}

func (s *salesService) DeleteSale(id uuid.UUID) error {
	return s.salesRepo.Delete(id)
}

func (s *salesService) ExportToExcel(filters map[string]interface{}, sorts []string) (*excelize.File, error) {
	// Get all sales with filters
	sales, _, err := s.salesRepo.GetSalesWithFilters(filters, sorts, 1, 1000000) // Large page size to get all
	if err != nil {
		return nil, err
	}

	// Create new Excel file
	f := excelize.NewFile()

	// Create headers
	headers := []string{
		"Invoice ID", "Shop", "Customer", "Sales By", "Sale Date",
		"Total", "Discount", "Final Total", "Remarks",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Add data
	for i, sale := range sales {
		row := i + 2 // Start from row 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), sale.ID)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), sale.Shop.Name)
		customerName := "Cash Sale"
		if sale.Customer != nil {
			customerName = sale.Customer.Name
		}
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), customerName)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), fmt.Sprintf("%s %s", sale.SalesBy.FirstName, sale.SalesBy.LastName))
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), sale.SaleDateTime.Format("2006-01-02 15:04:05"))
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), sale.Total+sale.Discount)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), sale.Discount)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), sale.Total)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", row), sale.Remarks)
	}

	return f, nil
}

func (s *salesService) GetAnalytics(shopID *uuid.UUID) (*repository.SalesAnalytics, error) {
	return s.salesRepo.GetSalesAnalytics(shopID, time.Now(), time.Now())
}

func (s *salesService) GetLast7DaysSales(shopID *uuid.UUID) ([]repository.DailySales, error) {
	return s.salesRepo.GetLast7DaysSales(shopID)
}
