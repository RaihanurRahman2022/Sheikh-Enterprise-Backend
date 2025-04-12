package middleware

import (
	"net/http"

	validator "Sheikh-Enterprise-Backend/internal/infrastructure/validation"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

// ValidationMiddleware handles validation errors consistently across the application
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				// Check if it's a validation error
				if validationErrors, ok := err.Err.(val.ValidationErrors); ok {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"errors": validator.FormatError(validationErrors),
					})
					return
				}
			}
		}
	}
}

// ValidateRequest is a helper function to validate request bodies
func ValidateRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		if validationErrors, ok := err.(val.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": validator.FormatError(validationErrors),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
		}
		return false
	}
	return true
}
