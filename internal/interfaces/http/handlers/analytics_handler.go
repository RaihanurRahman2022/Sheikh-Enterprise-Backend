package handlers

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	usecases "Sheikh-Enterprise-Backend/internal/usecases/impl"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService usecases.AnalyticsService
}

func NewAnalyticsHandler(analyticsService usecases.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetStockTransferAnalytics godoc
// @Summary Get stock transfer analytics
// @Description Get analytics for stock transfers within a date range
// @Tags analytics
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param shop_id query string false "Shop ID"
// @Success 200 {object} map[string]interface{}
// @Router /analytics/stock-transfers [get]
func (h *AnalyticsHandler) GetStockTransferAnalytics(c *gin.Context) {
	startDate, err := time.Parse("2006-01-02", c.Query("start_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	var shopID *uuid.UUID
	if shopIDStr := c.Query("shop_id"); shopIDStr != "" {
		parsedID, err := uuid.Parse(shopIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID format"})
			return
		}
		shopID = &parsedID
	}

	analytics, err := h.analyticsService.GetStockTransferAnalytics(startDate, endDate, shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetInventoryAnalytics godoc
// @Summary Get inventory analytics
// @Description Get analytics for inventory levels
// @Tags analytics
// @Accept json
// @Produce json
// @Param shop_id query string false "Shop ID"
// @Success 200 {object} map[string]interface{}
// @Router /analytics/inventory [get]
func (h *AnalyticsHandler) GetInventoryAnalytics(c *gin.Context) {
	var shopID *uuid.UUID
	if shopIDStr := c.Query("shop_id"); shopIDStr != "" {
		parsedID, err := uuid.Parse(shopIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID format"})
			return
		}
		shopID = &parsedID
	}

	analytics, err := h.analyticsService.GetInventoryAnalytics(shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// Helper functions for grouping data
func groupTransfersByDate(transfers []entities.StockTransfer) map[string]int {
	result := make(map[string]int)
	for _, transfer := range transfers {
		date := transfer.TransferDateTime.Format("2006-01-02")
		result[date] += transfer.Quantity
	}
	return result
}

func groupTransfersByMonth(transfers []entities.StockTransfer) map[string]int {
	result := make(map[string]int)
	for _, transfer := range transfers {
		month := transfer.TransferDateTime.Format("2006-01")
		result[month] += transfer.Quantity
	}
	return result
}

func groupInventoryByShop(inventories []entities.Inventory) map[uuid.UUID]int {
	result := make(map[uuid.UUID]int)
	for _, inventory := range inventories {
		result[inventory.ShopID] += inventory.Quantity
	}
	return result
}
