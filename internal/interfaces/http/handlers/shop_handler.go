
package handlers

import (
	"net/http"
	"strconv"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"Sheikh-Enterprise-Backend/internal/usecases/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShopHandler struct {
	shopService interfaces.ShopService
}

func NewShopHandler(shopService interfaces.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

func (h *ShopHandler) GetShops(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if companyID := c.Query("company_id"); companyID != "" {
		filters["company_id"] = companyID
	}

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

func (h *ShopHandler) CreateShop(c *gin.Context) {
	var shop entities.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.shopService.CreateShop(&shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create shop"})
		return
	}

	c.JSON(http.StatusCreated, shop)
}

func (h *ShopHandler) UpdateShop(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	var shop entities.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shop.ID = id
	if err := h.shopService.UpdateShop(&shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update shop"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

func (h *ShopHandler) DeleteShop(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	if err := h.shopService.DeleteShop(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete shop"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shop deleted successfully"})
}
