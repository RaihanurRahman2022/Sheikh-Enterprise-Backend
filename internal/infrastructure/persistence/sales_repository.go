package persistence

import (
	"time"

	"Sheikh-Enterprise-Backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SalesRepository interface {
	BaseRepository[entities.SalesInvoice]
	GetSalesWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.SalesInvoice, int64, error)
	GetSalesAnalytics(shopID *uuid.UUID, startDate, endDate time.Time) (*SalesAnalytics, error)
	GetLast7DaysSales(shopID *uuid.UUID) ([]DailySales, error)
}

type SalesAnalytics struct {
	TodaySales        float64 `json:"today_sales"`
	MonthlySales      float64 `json:"monthly_sales"`
	YearlySales       float64 `json:"yearly_sales"`
	ProductsSoldToday int     `json:"products_sold_today"`
}

type DailySales struct {
	Date  time.Time `json:"date"`
	Total float64   `json:"total"`
}

type salesRepository struct {
	BaseRepositoryImpl[entities.SalesInvoice]
}

func NewSalesRepository(db *gorm.DB) SalesRepository {
	return &salesRepository{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.SalesInvoice]{DB: db},
	}
}

func (r *salesRepository) GetSalesWithFilters(filters map[string]interface{}, sorts []string, page, pageSize int) ([]entities.SalesInvoice, int64, error) {
	var sales []entities.SalesInvoice
	var total int64

	query := r.DB.Model(&entities.SalesInvoice{}).
		Preload("Shop").
		Preload("Customer").
		Preload("SalesBy").
		Preload("SalesDetails").
		Preload("SalesDetails.Product").
		Where("is_marked_to_delete = ?", false)

	// Apply filters
	for field, value := range filters {
		switch field {
		case "shop_id", "customer_id", "sales_by_id":
			query = query.Where(field+" = ?", value)
		case "min_total":
			query = query.Where("total >= ?", value)
		case "max_total":
			query = query.Where("total <= ?", value)
		case "date_from":
			query = query.Where("sale_datetime >= ?", value)
		case "date_to":
			query = query.Where("sale_datetime <= ?", value)
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
	err = query.Offset(offset).Limit(pageSize).Find(&sales).Error
	if err != nil {
		return nil, 0, err
	}

	return sales, total, nil
}

func (r *salesRepository) GetSalesAnalytics(shopID *uuid.UUID, startDate, endDate time.Time) (*SalesAnalytics, error) {
	var analytics SalesAnalytics
	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC)
	startOfYear := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	// Base query with shop filter if provided
	baseQuery := r.DB.Model(&entities.SalesInvoice{}).Where("is_marked_to_delete = ?", false)
	if shopID != nil {
		baseQuery = baseQuery.Where("shop_id = ?", shopID)
	}

	// Get today's sales
	err := baseQuery.Where("sale_datetime >= ?", startOfDay).
		Select("COALESCE(SUM(total), 0)").
		Scan(&analytics.TodaySales).Error
	if err != nil {
		return nil, err
	}

	// Get monthly sales
	err = baseQuery.Where("sale_datetime >= ?", startOfMonth).
		Select("COALESCE(SUM(total), 0)").
		Scan(&analytics.MonthlySales).Error
	if err != nil {
		return nil, err
	}

	// Get yearly sales
	err = baseQuery.Where("sale_datetime >= ?", startOfYear).
		Select("COALESCE(SUM(total), 0)").
		Scan(&analytics.YearlySales).Error
	if err != nil {
		return nil, err
	}

	// Get products sold today
	err = r.DB.Model(&entities.SalesDetail{}).
		Joins("JOIN sales_invoices ON sales_details.invoice_id = sales_invoices.id").
		Where("sales_invoices.sale_datetime >= ?", startOfDay).
		Where("sales_invoices.is_marked_to_delete = ?", false).
		Select("COALESCE(SUM(sales_details.quantity), 0)").
		Scan(&analytics.ProductsSoldToday).Error

	return &analytics, err
}

func (r *salesRepository) GetLast7DaysSales(shopID *uuid.UUID) ([]DailySales, error) {
	var sales []DailySales
	today := time.Now().UTC()
	startDate := today.AddDate(0, 0, -6) // 7 days ago

	query := `
		WITH RECURSIVE dates AS (
			SELECT date_trunc('day', ?) AS date
			UNION ALL
			SELECT date + interval '1 day'
			FROM dates
			WHERE date < date_trunc('day', ?)
		)
		SELECT d.date,
			COALESCE(SUM(si.total), 0) as total
		FROM dates d
		LEFT JOIN sales_invoices si ON date_trunc('day', si.sale_datetime) = d.date
			AND si.is_marked_to_delete = false
	`

	if shopID != nil {
		query += " AND si.shop_id = ?"
		query += " GROUP BY d.date ORDER BY d.date"
		err := r.DB.Raw(query, startDate, today, shopID).Scan(&sales).Error
		return sales, err
	}

	query += " GROUP BY d.date ORDER BY d.date"
	err := r.DB.Raw(query, startDate, today).Scan(&sales).Error
	return sales, err
}
