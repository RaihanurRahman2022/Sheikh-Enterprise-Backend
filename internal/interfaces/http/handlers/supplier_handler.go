
package handlers

import (
	"net/http"
	"strconv"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	validator "Sheikh-Enterprise-Backend/internal/infrastructure/validation"
	services "Sheikh-Enterprise-Backend/internal/usecases/impl"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SupplierHandler struct {
	supplierService services.SupplierService
}

func NewSupplierHandler(supplierService services.SupplierService) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
	}
}

func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if phone := c.Query("phone"); phone != "" {
		filters["phone"] = phone
	}
	if email := c.Query("email"); email != "" {
		filters["email"] = email
	}

	var sorts []string
	if sort := c.Query("sort"); sort != "" {
		sorts = append(sorts, sort)
	}

	suppliers, total, err := h.supplierService.GetSuppliers(page, pageSize, filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch suppliers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": suppliers,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier ID"})
		return
	}

	supplier, err := h.supplierService.GetSupplierByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var supplier entities.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.supplierService.CreateSupplier(&supplier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create supplier"})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier ID"})
		return
	}

	if err := h.supplierService.DeleteSupplier(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "supplier deleted successfully"})
}
