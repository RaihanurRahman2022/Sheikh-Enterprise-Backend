package usecases

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type ProductService interface {
	GetProducts(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.Product, int64, error)
	GetProductByID(id uuid.UUID) (*entities.Product, error)
	CreateProduct(product *entities.Product) error
	UpdateProduct(product *entities.Product) error
	DeleteProduct(id uuid.UUID) error
	BulkImportProducts(reader io.Reader) error
	ExportToExcel(filters map[string]interface{}, sorts []string) (*excelize.File, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) GetProducts(page, pageSize int, filters map[string]interface{}, sorts []string) ([]entities.Product, int64, error) {
	return s.productRepo.GetProductsWithFilters(filters, sorts, page, pageSize)
}

func (s *productService) GetProductByID(id uuid.UUID) (*entities.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *productService) CreateProduct(product *entities.Product) error {
	return s.productRepo.Create(product)
}

func (s *productService) UpdateProduct(product *entities.Product) error {
	return s.productRepo.Update(product)
}

func (s *productService) DeleteProduct(id uuid.UUID) error {
	return s.productRepo.Delete(id)
}

func (s *productService) BulkImportProducts(reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	// Skip header row
	_, err := csvReader.Read()
	if err != nil {
		return err
	}

	var products []entities.Product
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		purchasePrice, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return fmt.Errorf("invalid purchase price in row: %v", record)
		}

		salesPrice, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return fmt.Errorf("invalid sales price in row: %v", record)
		}

		product := entities.Product{
			Code:           record[0],
			Name:           record[1],
			Style:          record[2],
			MasterCategory: record[3],
			SubCategory:    record[4],
			Color:          record[5],
			Size:           record[6],
			PurchasePrice:  purchasePrice,
			SalesPrice:     salesPrice,
			SalesType:      entities.SalesType(record[9]),
			Remarks:        record[10],
		}

		products = append(products, product)
	}

	return s.productRepo.BulkCreate(products)
}

func (s *productService) ExportToExcel(filters map[string]interface{}, sorts []string) (*excelize.File, error) {
	// Get all products with filters
	products, _, err := s.productRepo.GetProductsWithFilters(filters, sorts, 1, 1000000) // Large page size to get all
	if err != nil {
		return nil, err
	}

	// Create new Excel file
	f := excelize.NewFile()

	// Create headers
	headers := []string{
		"Code", "Name", "Style", "Master Category", "Sub Category",
		"Color", "Size", "Purchase Price", "Sales Price", "Sales Type", "Remarks",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Add data
	for i, product := range products {
		row := i + 2 // Start from row 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), product.Code)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), product.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), product.Style)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), product.MasterCategory)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), product.SubCategory)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), product.Color)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), product.Size)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), product.PurchasePrice)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", row), product.SalesPrice)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", row), product.SalesType)
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", row), product.Remarks)
	}

	return f, nil
}
