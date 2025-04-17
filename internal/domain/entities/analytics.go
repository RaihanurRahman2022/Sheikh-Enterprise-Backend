package entities

import "time"

// SalesSummary represents a summary of sales metrics
type SalesSummary struct {
	TotalSales        int     `json:"total_sales"`
	TotalAmount       float64 `json:"total_amount"`
	TotalDiscount     float64 `json:"total_discount"`
	AverageSaleAmount float64 `json:"average_sale_amount"`
}

// DailySales represents daily sales metrics
type DailySales struct {
	Date        time.Time `json:"date"`
	TotalSales  int       `json:"total_sales"`
	TotalAmount float64   `json:"total_amount"`
}

// TopProduct represents a product with its sales metrics
type TopProduct struct {
	ProductID     string  `json:"product_id"`
	Name          string  `json:"name"`
	TotalQuantity int     `json:"total_quantity"`
	TotalAmount   float64 `json:"total_amount"`
}

// SalesTrend represents sales trend data over a period
type SalesTrend struct {
	Period        string  `json:"period"`
	TotalSales    int     `json:"total_sales"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}
