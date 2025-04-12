
package handlers

import (
	"net/http"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"Sheikh-Enterprise-Backend/internal/usecases/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CompanyHandler struct {
	companyService interfaces.CompanyService
}

func NewCompanyHandler(companyService interfaces.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	companies, err := h.companyService.GetAllCompanies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch companies"})
		return
	}
	c.JSON(http.StatusOK, companies)
}

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

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company entities.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.companyService.CreateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, company)
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	var company entities.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company.ID = id
	if err := h.companyService.UpdateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update company"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	if err := h.companyService.DeleteCompany(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "company deleted successfully"})
}
