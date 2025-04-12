package handlers

import (
	"net/http"

	models "Sheikh-Enterprise-Backend/internal/domain/entities"
	validator "Sheikh-Enterprise-Backend/internal/infrastructure/validation"
	"Sheikh-Enterprise-Backend/internal/usecases/interfaces"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService interfaces.AuthService
}

func NewAuthHandler(authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} validator.ValidationErrors
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == interfaces.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		if err == interfaces.ErrUserInactive {
			c.JSON(http.StatusForbidden, gin.H{"error": "user account is inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Register godoc
// @Summary Register new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} validator.ValidationErrors
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Register(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param passwords body models.ChangePasswordRequest true "Old and new passwords"
// @Success 200 {object} map[string]string
// @Failure 400 {object} validator.ValidationErrors
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors := validator.FormatError(err); validationErrors != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.authService.ChangePassword(userID.(string), req.OldPassword, req.NewPassword); err != nil {
		if err == interfaces.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
