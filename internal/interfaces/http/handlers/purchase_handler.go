package handlers

import (
	"net/http"
	"strconv"
	"time"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	validator "Sheikh-Enterprise-Backend/internal/infrastructure/validation"
	services "Sheikh-Enterprise-Backend/internal/usecases/impl"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PurchaseHandler struct {
	purchaseService services.PurchaseService
}

func NewPurchaseHandler(purchaseService services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService: purchaseService,
	}
}

func (h *PurchaseHandler) GetPurchases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if supplierID := c.Query("supplier_id"); supplierID != "" {
		filters["supplier_id"] = supplierID
	}
	if paymentType := c.Query("payment_type"); paymentType != "" {
		filters["payment_type"] = paymentType
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

	purchases, total, err := h.purchaseService.GetPurchases(page, pageSize, filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch purchases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": purchases,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *PurchaseHandler) GetPurchase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase ID"})
		return
	}

	purchase, err := h.purchaseService.GetPurchaseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *PurchaseHandler) CreatePurchase(c *gin.Context) {
	var purchase entities.PurchaseInvoice
	if err := c.ShouldBindJSON(&purchase); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set entry timestamp
	purchase.PurchaseDateTime = time.Now()

	// Get user from context for entry_by
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	purchase.EntryByID = user.(entities.User).ID

	if err := h.purchaseService.CreatePurchase(&purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create purchase"})
		return
	}

	c.JSON(http.StatusCreated, purchase)
}

func (h *PurchaseHandler) DeletePurchase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase ID"})
		return
	}

	if err := h.purchaseService.DeletePurchase(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "purchase deleted successfully"})
}