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

type ShopHandler struct {
	shopService services.ShopService
}

func NewShopHandler(shopService services.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
	}
}

// GetShops godoc
// @Summary List shops
// @Description Get a paginated list of shops with optional filters
// @Tags shops
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param company_id query string false "Filter by company ID"
// @Success 200 {object} map[string]interface{}
// @Router /shops [get]
// @Security BearerAuth
func (h *ShopHandler) GetShops(c *gin.Context) {
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
	if managerName := c.Query("manager_name"); managerName != "" {
		filters["manager_name"] = managerName
	}
	if managerPhone := c.Query("manager_phone"); managerPhone != "" {
		filters["manager_phone"] = managerPhone
	}
	if companyID := c.Query("company_id"); companyID != "" {
		filters["company_id"] = companyID
	}

	// Get sort parameters
	var sorts []string
	if sort := c.Query("sort"); sort != "" {
		sorts = append(sorts, sort)
	}

	shops, total, err := h.shopService.GetShops(page, pageSize, filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch shops"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": shops,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetShop godoc
// @Summary Get a shop by ID
// @Description Get detailed information about a shop
// @Tags shops
// @Accept json
// @Produce json
// @Param id path string true "Shop ID"
// @Success 200 {object} entities.Shop
// @Router /shops/{id} [get]
func (h *ShopHandler) GetShop(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	shop, err := h.shopService.GetShopByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// GetShopsByCompany godoc
// @Summary Get shops by company ID
// @Description Get all shops belonging to a specific company
// @Tags shops
// @Accept json
// @Produce json
// @Param company_id path string true "Company ID"
// @Success 200 {array} entities.Shop
// @Router /companies/{company_id}/shops [get]
func (h *ShopHandler) GetShopsByCompany(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	shops, err := h.shopService.GetShopsByCompanyID(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch shops"})
		return
	}

	c.JSON(http.StatusOK, shops)
}

// CreateShop godoc
// @Summary Create shop
// @Description Create a new shop
// @Tags shops
// @Accept json
// @Produce json
// @Param shop body entities.CreateShopRequest true "Shop details"
// @Success 201 {object} entities.Shop
// @Failure 400 {object} validator.ValidationErrors
// @Router /shops [post]
// @Security BearerAuth
func (h *ShopHandler) CreateShop(c *gin.Context) {
	var req entities.CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	shop := &entities.Shop{
		CompanyID:    companyID,
		Name:         req.Name,
		Address:      req.Address,
		Phone:        req.Phone,
		Email:        req.Email,
		ManagerName:  req.ManagerName,
		ManagerPhone: req.ManagerPhone,
		Remarks:      req.Remarks,
	}

	if err := h.shopService.CreateShop(shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create shop"})
		return
	}

	c.JSON(http.StatusCreated, shop)
}

// UpdateShop godoc
// @Summary Update shop
// @Description Update an existing shop
// @Tags shops
// @Accept json
// @Produce json
// @Param id path string true "Shop ID"
// @Param shop body entities.CreateShopRequest true "Shop details"
// @Success 200 {object} entities.Shop
// @Failure 400 {object} validator.ValidationErrors
// @Router /shops/{id} [put]
// @Security BearerAuth
func (h *ShopHandler) UpdateShop(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	var req entities.CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	// First get the existing shop
	existingShop, err := h.shopService.GetShopByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}

	// Update the fields
	existingShop.CompanyID = companyID
	existingShop.Name = req.Name
	existingShop.Address = req.Address
	existingShop.Phone = req.Phone
	existingShop.Email = req.Email
	existingShop.ManagerName = req.ManagerName
	existingShop.ManagerPhone = req.ManagerPhone
	existingShop.Remarks = req.Remarks

	if err := h.shopService.UpdateShop(existingShop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update shop"})
		return
	}

	c.JSON(http.StatusOK, existingShop)
}

// DeleteShop godoc
// @Summary Delete a shop
// @Description Mark a shop as deleted
// @Tags shops
// @Accept json
// @Produce json
// @Param id path string true "Shop ID"
// @Success 200 {object} map[string]string
// @Router /shops/{id} [delete]
func (h *ShopHandler) DeleteShop(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	if err := h.shopService.DeleteShop(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shop deleted successfully"})
}
