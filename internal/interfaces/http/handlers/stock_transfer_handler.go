
package handlers

import (
	"net/http"
	"strconv"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"Sheikh-Enterprise-Backend/internal/usecases/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockTransferHandler struct {
	stockTransferService interfaces.StockTransferService
}

func NewStockTransferHandler(stockTransferService interfaces.StockTransferService) *StockTransferHandler {
	return &StockTransferHandler{stockTransferService: stockTransferService}
}

func (h *StockTransferHandler) GetStockTransfers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if productID := c.Query("product_id"); productID != "" {
		filters["product_id"] = productID
	}
	if toShopID := c.Query("to_shop_id"); toShopID != "" {
		filters["to_shop_id"] = toShopID
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}

	var sorts []string
	if sort := c.Query("sort"); sort != "" {
		sorts = append(sorts, sort)
	}

	transfers, total, err := h.stockTransferService.GetStockTransfers(page, pageSize, filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stock transfers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transfers,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *StockTransferHandler) GetStockTransfer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock transfer ID"})
		return
	}

	transfer, err := h.stockTransferService.GetStockTransferByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock transfer not found"})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (h *StockTransferHandler) CreateStockTransfer(c *gin.Context) {
	var transfer entities.StockTransfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.stockTransferService.CreateStockTransfer(&transfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create stock transfer"})
		return
	}

	c.JSON(http.StatusCreated, transfer)
}

func (h *StockTransferHandler) DeleteStockTransfer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock transfer ID"})
		return
	}

	if err := h.stockTransferService.DeleteStockTransfer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete stock transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock transfer deleted successfully"})
}
