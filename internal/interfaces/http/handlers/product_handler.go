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

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProducts godoc
// @Summary List products
// @Description Get a paginated list of products with optional filters
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
// @Security BearerAuth
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Get filters from query parameters
	filters := make(map[string]interface{})
	if code := c.Query("code"); code != "" {
		filters["code"] = code
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if style := c.Query("style"); style != "" {
		filters["style"] = style
	}
	if category := c.Query("master_category"); category != "" {
		filters["master_category"] = category
	}
	if subcategory := c.Query("sub_category"); subcategory != "" {
		filters["sub_category"] = subcategory
	}

	// Get sort parameters
	var sorts []string
	if sort := c.Query("sort"); sort != "" {
		sorts = append(sorts, sort)
	}

	products, total, err := h.productService.GetProducts(page, pageSize, filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get detailed information about a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entities.Product
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := h.productService.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body entities.CreateProductRequest true "Product details"
// @Success 201 {object} entities.Product
// @Failure 400 {object} validator.ValidationErrors
// @Router /products [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req entities.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to product model
	product := &entities.Product{
		Code:           req.Code,
		Name:           req.Name,
		Style:          req.Style,
		MasterCategory: req.MasterCategory,
		SubCategory:    req.SubCategory,
		Color:          req.Color,
		Size:           req.Size,
		PurchasePrice:  req.PurchasePrice,
		SalesPrice:     req.SalesPrice,
		SalesType:      entities.SalesType(req.SalesType),
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}
	product.ShopID = shopID

	if err := h.productService.CreateProduct(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Mark a product as deleted
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// BulkImport godoc
// @Summary Bulk import products
// @Description Import multiple products from a CSV file
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing product data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /products/bulk-import [post]
// @Security BearerAuth
func (h *ProductHandler) BulkImport(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty file"})
		return
	}

	// Check file extension
	if ext := file.Filename[len(file.Filename)-4:]; ext != ".csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only CSV files are allowed"})
		return
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	if err := h.productService.BulkImportProducts(f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to import products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "products imported successfully"})
}

// ExportToExcel godoc
// @Summary Export products to Excel
// @Description Export filtered products to an Excel file
// @Tags products
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param filters query string false "JSON encoded filters"
// @Success 200 {file} file
// @Router /products/export [get]
// @Security BearerAuth
func (h *ProductHandler) ExportToExcel(c *gin.Context) {
	// Get filters from query parameters (similar to GetProducts)
	filters := make(map[string]interface{})
	if code := c.Query("code"); code != "" {
		filters["code"] = code
	}
	// ... add other filters ...

	// Get sort parameters
	var sorts []string
	if sort := c.Query("sort"); sort != "" {
		sorts = append(sorts, sort)
	}

	file, err := h.productService.ExportToExcel(filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export products"})
		return
	}

	// Set headers for file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=products.xlsx")

	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
		return
	}
}
