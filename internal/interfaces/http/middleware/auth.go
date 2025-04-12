package middleware

import (
	"net/http"
	"strings"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	"Sheikh-Enterprise-Backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("shop_id", claims.ShopID)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role not found in token"})
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

func ShopMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		shopID, exists := c.Get("shop_id")

		// Admin can access all shops
		if role == string(entities.RoleAdmin) {
			c.Next()
			return
		}

		// Other roles must have a shop_id
		if !exists || shopID == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "shop access not configured for user"})
			return
		}

		// Check if the requested shop matches user's assigned shop
		if requestedShopID := c.Param("shop_id"); requestedShopID != "" {
			if requestedShopID != shopID.(string) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied to requested shop"})
				return
			}
		}

		c.Next()
	}
}
