package handlers

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	usecases "Sheikh-Enterprise-Backend/internal/usecases/impl"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockTransferHandler struct {
	stockTransferService usecases.StockTransferService
}

func NewStockTransferHandler(stockTransferService usecases.StockTransferService) *StockTransferHandler {
	return &StockTransferHandler{
		stockTransferService: stockTransferService,
	}
}

// GetStockTransfers godoc
// @Summary Get all stock transfers
// @Description Get a list of stock transfers with pagination and filtering
// @Tags stock-transfers
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param filters query string false "Filters in JSON format"
// @Param sorts query string false "Sort fields"
// @Success 200 {object} map[string]interface{}
// @Router /stock-transfers [get]
func (h *StockTransferHandler) GetStockTransfers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	filters := make(map[string]interface{})
	sorts := make([]string, 0)

	transfers, total, err := h.stockTransferService.GetStockTransfers(page, pageSize, filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transfers,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetStockTransfer godoc
// @Summary Get a stock transfer by ID
// @Description Get a specific stock transfer by its ID
// @Tags stock-transfers
// @Accept json
// @Produce json
// @Param id path string true "Stock Transfer ID"
// @Success 200 {object} entities.StockTransfer
// @Router /stock-transfers/{id} [get]
func (h *StockTransferHandler) GetStockTransfer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	transfer, err := h.stockTransferService.GetStockTransferByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock transfer not found"})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

// GetStockTransfersByShopID godoc
// @Summary Get stock transfers by shop ID
// @Description Get all stock transfers for a specific shop
// @Tags stock-transfers
// @Accept json
// @Produce json
// @Param shop_id path string true "Shop ID"
// @Success 200 {array} entities.StockTransfer
// @Router /shops/{shop_id}/stock-transfers [get]
func (h *StockTransferHandler) GetStockTransfersByShopID(c *gin.Context) {
	shopID, err := uuid.Parse(c.Param("shop_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shop ID format"})
		return
	}

	transfers, err := h.stockTransferService.GetStockTransfersByShopID(shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transfers)
}

// CreateStockTransfer godoc
// @Summary Create a new stock transfer
// @Description Create a new stock transfer between shops
// @Tags stock-transfers
// @Accept json
// @Produce json
// @Param transfer body entities.CreateStockTransferRequest true "Stock Transfer"
// @Success 201 {object} entities.StockTransfer
// @Router /stock-transfers [post]
func (h *StockTransferHandler) CreateStockTransfer(c *gin.Context) {
	var request entities.CreateStockTransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(request.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	toShopID, err := uuid.Parse(request.ToShopID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to shop ID format"})
		return
	}

	fromShopID, err := uuid.Parse(request.FromShopID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from shop ID format"})
		return
	}

	transfer := &entities.StockTransfer{
		ProductID:        productID,
		ToShopID:         toShopID,
		FromShopID:       fromShopID,
		Quantity:         request.Quantity,
		TransferDateTime: request.TransferDateTime,
		TransferredBy:    uuid.MustParse(c.GetString("user_id")), // Assuming user ID is stored in context
	}

	if err := h.stockTransferService.CreateStockTransfer(transfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transfer)
}

// UpdateStockTransfer godoc
// @Summary Update a stock transfer
// @Description Update an existing stock transfer
// @Tags stock-transfers
// @Accept json
// @Produce json
// @Param id path string true "Stock Transfer ID"
// @Param transfer body entities.StockTransfer true "Stock Transfer"
// @Success 200 {object} entities.StockTransfer
// @Router /stock-transfers/{id} [put]
func (h *StockTransferHandler) UpdateStockTransfer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var transfer entities.StockTransfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transfer.ID = id
	if err := h.stockTransferService.UpdateStockTransfer(&transfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

// DeleteStockTransfer godoc
// @Summary Delete a stock transfer
// @Description Delete an existing stock transfer
// @Tags stock-transfers
// @Accept json
// @Produce json
// @Param id path string true "Stock Transfer ID"
// @Success 204 "No Content"
// @Router /stock-transfers/{id} [delete]
func (h *StockTransferHandler) DeleteStockTransfer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.stockTransferService.DeleteStockTransfer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
