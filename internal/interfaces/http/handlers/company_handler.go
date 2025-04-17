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

type CompanyHandler struct {
	companyService services.CompanyService
}

func NewCompanyHandler(companyService services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

// GetCompanies godoc
// @Summary List companies
// @Description Get a paginated list of companies with optional filters
// @Tags companies
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /companies [get]
// @Security BearerAuth
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
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

	companies, total, err := h.companyService.GetCompanies(page, pageSize, filters, sorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch companies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": companies,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetCompany godoc
// @Summary Get a company by ID
// @Description Get detailed information about a company
// @Tags companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} entities.Company
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	company, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// CreateCompany godoc
// @Summary Create company
// @Description Create a new company
// @Tags companies
// @Accept json
// @Produce json
// @Param company body entities.CreateCompanyRequest true "Company details"
// @Success 201 {object} entities.Company
// @Failure 400 {object} validator.ValidationErrors
// @Router /companies [post]
// @Security BearerAuth
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req entities.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company := &entities.Company{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
		Slogan:  req.Slogan,
		Remarks: req.Remarks,
	}

	if err := h.companyService.CreateCompany(company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, company)
}

// UpdateCompany godoc
// @Summary Update company
// @Description Update an existing company
// @Tags companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param company body entities.CreateCompanyRequest true "Company details"
// @Success 200 {object} entities.Company
// @Failure 400 {object} validator.ValidationErrors
// @Router /companies/{id} [put]
// @Security BearerAuth
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	var req entities.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First get the existing company
	existingCompany, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	// Update the fields
	existingCompany.Name = req.Name
	existingCompany.Address = req.Address
	existingCompany.Phone = req.Phone
	existingCompany.Email = req.Email
	existingCompany.Slogan = req.Slogan
	existingCompany.Remarks = req.Remarks

	if err := h.companyService.UpdateCompany(existingCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update company"})
		return
	}

	c.JSON(http.StatusOK, existingCompany)
}

// DeleteCompany godoc
// @Summary Delete a company
// @Description Mark a company as deleted
// @Tags companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} map[string]string
// @Router /companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	if err := h.companyService.DeleteCompany(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "company deleted successfully"})
}
