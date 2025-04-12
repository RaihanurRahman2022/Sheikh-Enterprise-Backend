package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	services "Sheikh-Enterprise-Backend/internal/usecases/impl"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SalesHandler struct {
	salesService services.SalesService
}

func NewSalesHandler(salesService services.SalesService) *SalesHandler {
	return &SalesHandler{
		salesService: salesService,
	}
}

// GetSales godoc
// @Summary Get list of sales
// @Description Get paginated list of sales with filtering and sorting
// @Tags sales
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param sort query string false "Sort fields (comma-separated, prefix with - for desc)" example("sale_datetime,-total")
// @Success 200 {object} map[string]interface{}
// @Router /sales [get]
func (h *SalesHandler) GetSales(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Parse filters
	filters := make(map[string]interface{})
	for key, value := range c.Request.URL.Query() {
		if !strings.HasPrefix(key, "filter_") {
			continue
		}
		filterKey := strings.TrimPrefix(key, "filter_")
		filters[filterKey] = value[0]
	}

	// Parse sort parameters
	var sorts []string
	if sortParam := c.Query("sort"); sortParam != "" {
		sorts = strings.Split(sortParam, ",")
	}

	sales, total, err := h.salesService.GetSales(page, pageSize, filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": sales,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetSale godoc
// @Summary Get a sale by ID
// @Description Get detailed information about a sale
// @Tags sales
// @Accept json
// @Produce json
// @Param id path string true "Sale ID"
// @Success 200 {object} entities.SalesInvoice
// @Router /sales/{id} [get]
func (h *SalesHandler) GetSale(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sale ID"})
		return
	}

	sale, err := h.salesService.GetSaleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sale not found"})
		return
	}

	c.JSON(http.StatusOK, sale)
}

// CreateSale godoc
// @Summary Create a new sale
// @Description Create a new sale with the provided information
// @Tags sales
// @Accept json
// @Produce json
// @Param sale body entities.SalesInvoice true "Sale information"
// @Success 201 {object} entities.SalesInvoice
// @Router /sales [post]
func (h *SalesHandler) CreateSale(c *gin.Context) {
	var sale entities.SalesInvoice
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set current user as sales_by
	if userID, exists := c.Get("user_id"); exists {
		sale.SalesByID = userID.(uuid.UUID)
	}

	// Set current shop if not admin
	if shopID, exists := c.Get("shop_id"); exists && shopID != nil {
		sale.ShopID = shopID.(uuid.UUID)
	}

	sale.SaleDateTime = time.Now()

	if err := h.salesService.CreateSale(&sale); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sale)
}

// DeleteSale godoc
// @Summary Delete a sale
// @Description Mark a sale as deleted
// @Tags sales
// @Accept json
// @Produce json
// @Param id path string true "Sale ID"
// @Success 200 {object} map[string]string
// @Router /sales/{id} [delete]
func (h *SalesHandler) DeleteSale(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sale ID"})
		return
	}

	if err := h.salesService.DeleteSale(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sale deleted successfully"})
}

// ExportToExcel godoc
// @Summary Export sales to Excel
// @Description Export filtered sales to Excel file
// @Tags sales
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param filter query string false "Filter parameters"
// @Param sort query string false "Sort parameters"
// @Success 200 {file} file
// @Router /sales/export [get]
func (h *SalesHandler) ExportToExcel(c *gin.Context) {
	filters := make(map[string]interface{})
	for key, value := range c.Request.URL.Query() {
		if !strings.HasPrefix(key, "filter_") {
			continue
		}
		filterKey := strings.TrimPrefix(key, "filter_")
		filters[filterKey] = value[0]
	}

	var sorts []string
	if sortParam := c.Query("sort"); sortParam != "" {
		sorts = strings.Split(sortParam, ",")
	}

	// Add shop filter for non-admin users
	if shopID, exists := c.Get("shop_id"); exists && shopID != nil {
		filters["shop_id"] = shopID
	}

	file, err := h.salesService.ExportToExcel(filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=sales_%s.xlsx", time.Now().Format("20060102150405")))

	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
		return
	}
}

// GetAnalytics godoc
// @Summary Get sales analytics
// @Description Get sales analytics including today's, monthly, and yearly sales
// @Tags sales
// @Accept json
// @Produce json
// @Success 200 {object} repository.SalesAnalytics
// @Router /sales/analytics [get]
func (h *SalesHandler) GetAnalytics(c *gin.Context) {
	var shopID *uuid.UUID
	if id, exists := c.Get("shop_id"); exists && id != nil {
		sid := id.(uuid.UUID)
		shopID = &sid
	}

	analytics, err := h.salesService.GetAnalytics(shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetLast7DaysSales godoc
// @Summary Get last 7 days sales
// @Description Get daily sales data for the last 7 days
// @Tags sales
// @Accept json
// @Produce json
// @Success 200 {array} repository.DailySales
// @Router /sales/last-7-days [get]
func (h *SalesHandler) GetLast7DaysSales(c *gin.Context) {
	var shopID *uuid.UUID
	if id, exists := c.Get("shop_id"); exists && id != nil {
		sid := id.(uuid.UUID)
		shopID = &sid
	}

	sales, err := h.salesService.GetLast7DaysSales(shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sales)
}
