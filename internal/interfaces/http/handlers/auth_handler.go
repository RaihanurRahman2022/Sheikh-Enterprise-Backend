package handlers

import (
	"Sheikh-Enterprise-Backend/internal/domain/entities"
	_val "Sheikh-Enterprise-Backend/internal/infrastructure/validation"
	services "Sheikh-Enterprise-Backend/internal/usecases/impl"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body entities.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entities.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		if err == services.ErrUserInactive {
			c.JSON(http.StatusForbidden, gin.H{"error": "user account is inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entities.RegisterRequest true "User registration details"
// @Success 201 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	reqType := reflect.TypeOf(entities.RegisterRequest{})
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fmt.Printf("Field: %s, JSON Tag: %s\n", field.Name, field.Tag.Get("json"))
	}
	// Read request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error reading body:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	// Log raw body
	fmt.Println("Raw request body:", string(bodyBytes))

	// Restore body for ShouldBindJSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req entities.RegisterRequest
	fmt.Println("Attempting to bind JSON")
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Binding error:", err.Error())
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := _val.FormatError(validationErrors)
			fmt.Println("Validation errors:", formattedErrors)
			c.JSON(http.StatusBadRequest, gin.H{"errors": formattedErrors})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// Debug bound struct
	fmt.Printf("Bound struct: %+v\n", req)

	if err := h.authService.Register(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for the currently logged in user
// @Tags auth
// @Accept json
// @Produce json
// @Param passwords body entities.ChangePasswordRequest true "Current and new password"
// @Success 200 {object} map[string]string
// @Router /auth/password [put]
// @Security BearerAuth
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var req entities.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ChangePassword(userID.(string), req.OldPassword, req.NewPassword); err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

// RefreshToken godoc
// @Summary Refresh JWT tokens
// @Description Get new access and refresh tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		if err == services.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
