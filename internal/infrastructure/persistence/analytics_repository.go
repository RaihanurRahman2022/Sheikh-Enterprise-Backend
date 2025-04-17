package persistence

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetSalesSummary(shopID *string, startDate, endDate time.Time) (*entities.SalesSummary, error)
	GetDailySales(shopID *string, startDate, endDate time.Time) ([]entities.DailySales, error)
	GetTopProducts(shopID *string, limit int, startDate, endDate time.Time) ([]entities.TopProduct, error)
	GetSalesTrend(shopID *string, period string) ([]entities.SalesTrend, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{
		db: db,
	}
}

func (r *analyticsRepository) GetSalesSummary(shopID *string, startDate, endDate time.Time) (*entities.SalesSummary, error) {
	var summary entities.SalesSummary
	query := r.db.Table("sales_invoices")

	if shopID != nil {
		query = query.Where("shop_id = ?", *shopID)
	}

	err := query.
		Select(`
			COUNT(*) as total_sales,
			SUM(total) as total_amount,
			SUM(discount) as total_discount,
			AVG(total) as average_sale_amount
		`).
		Where("sale_datetime BETWEEN ? AND ?", startDate, endDate).
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (r *analyticsRepository) GetDailySales(shopID *string, startDate, endDate time.Time) ([]entities.DailySales, error) {
	var dailySales []entities.DailySales
	query := r.db.Table("sales_invoices")

	if shopID != nil {
		query = query.Where("shop_id = ?", *shopID)
	}

	err := query.
		Select(`
			DATE(sale_datetime) as date,
			COUNT(*) as total_sales,
			SUM(total) as total_amount
		`).
		Where("sale_datetime BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(sale_datetime)").
		Order("date").
		Scan(&dailySales).Error

	if err != nil {
		return nil, err
	}

	return dailySales, nil
}

func (r *analyticsRepository) GetTopProducts(shopID *string, limit int, startDate, endDate time.Time) ([]entities.TopProduct, error) {
	var topProducts []entities.TopProduct
	query := r.db.Table("sales_details").
		Joins("JOIN sales_invoices ON sales_details.invoice_id = sales_invoices.invoice_id").
		Joins("JOIN products ON sales_details.product_id = products.product_id")

	if shopID != nil {
		query = query.Where("sales_invoices.shop_id = ?", *shopID)
	}

	err := query.
		Select(`
			products.product_id,
			products.name,
			SUM(sales_details.quantity) as total_quantity,
			SUM(sales_details.subtotal) as total_amount
		`).
		Where("sales_invoices.sale_datetime BETWEEN ? AND ?", startDate, endDate).
		Group("products.product_id, products.name").
		Order("total_quantity DESC").
		Limit(limit).
		Scan(&topProducts).Error

	if err != nil {
		return nil, err
	}

	return topProducts, nil
}

func (r *analyticsRepository) GetSalesTrend(shopID *string, period string) ([]entities.SalesTrend, error) {
	var trend []entities.SalesTrend
	query := r.db.Table("sales_invoices")

	if shopID != nil {
		query = query.Where("shop_id = ?", *shopID)
	}

	groupBy := "DATE(sale_datetime)"
	switch period {
	case "monthly":
		groupBy = "DATE_FORMAT(sale_datetime, '%Y-%m')"
	case "yearly":
		groupBy = "YEAR(sale_datetime)"
	}

	err := query.
		Select(`
			${groupBy} as period,
			COUNT(*) as total_sales,
			SUM(total) as total_amount,
			AVG(total) as average_amount
		`).
		Group(groupBy).
		Order("period").
		Scan(&trend).Error

	if err != nil {
		return nil, err
	}

	return trend, nil
}
