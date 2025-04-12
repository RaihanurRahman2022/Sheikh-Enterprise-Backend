
package handlers

import (
	"net/http"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"Sheikh-Enterprise-Backend/internal/usecases/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	customerService interfaces.CustomerService
}

func NewCustomerHandler(customerService interfaces.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	customers, err := h.customerService.GetAllCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch customers"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	customer, err := h.customerService.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var customer entities.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.customerService.CreateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	var customer entities.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer.ID = id
	if err := h.customerService.UpdateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	if err := h.customerService.DeleteCustomer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted successfully"})
}
