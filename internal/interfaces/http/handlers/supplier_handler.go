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

// GetSuppliers godoc
// @Summary List suppliers
// @Description Get a paginated list of suppliers with optional filters
// @Tags suppliers
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /suppliers [get]
// @Security BearerAuth
func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Get filters from query parameters
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

	// Get sort parameters
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

// GetSupplier godoc
// @Summary Get a supplier by ID
// @Description Get detailed information about a supplier
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Success 200 {object} entities.Supplier
// @Router /suppliers/{id} [get]
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

// CreateSupplier godoc
// @Summary Create supplier
// @Description Create a new supplier
// @Tags suppliers
// @Accept json
// @Produce json
// @Param supplier body entities.CreateSupplierRequest true "Supplier details"
// @Success 201 {object} entities.Supplier
// @Failure 400 {object} validator.ValidationErrors
// @Router /suppliers [post]
// @Security BearerAuth
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req entities.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supplier := &entities.Supplier{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
		Remarks: req.Remarks,
	}

	if err := h.supplierService.CreateSupplier(supplier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create supplier"})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

// UpdateSupplier godoc
// @Summary Update supplier
// @Description Update an existing supplier
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Param supplier body entities.CreateSupplierRequest true "Supplier details"
// @Success 200 {object} entities.Supplier
// @Failure 400 {object} validator.ValidationErrors
// @Router /suppliers/{id} [put]
// @Security BearerAuth
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier ID"})
		return
	}

	var req entities.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First get the existing supplier
	existingSupplier, err := h.supplierService.GetSupplierByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return
	}

	// Update the fields
	existingSupplier.Name = req.Name
	existingSupplier.Address = req.Address
	existingSupplier.Phone = req.Phone
	existingSupplier.Email = req.Email
	existingSupplier.Remarks = req.Remarks

	if err := h.supplierService.UpdateSupplier(existingSupplier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update supplier"})
		return
	}

	c.JSON(http.StatusOK, existingSupplier)
}

// DeleteSupplier godoc
// @Summary Delete a supplier
// @Description Mark a supplier as deleted
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Success 200 {object} map[string]string
// @Router /suppliers/{id} [delete]
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
