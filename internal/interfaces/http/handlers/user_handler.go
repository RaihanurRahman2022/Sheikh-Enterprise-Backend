package handlers

import (
	"net/http"

	"Sheikh-Enterprise-Backend/internal/domain/entities"
	services "Sheikh-Enterprise-Backend/internal/usecases/impl"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserDetails godoc
// @Summary Get user details
// @Description Get details of the currently logged in user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} entities.User
// @Router /users/me [get]
// @Security BearerAuth
func (h *UserHandler) GetUserDetails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserDetails godoc
// @Summary Update user details
// @Description Update details of the currently logged in user
// @Tags users
// @Accept json
// @Produce json
// @Param user body entities.User true "User details"
// @Success 200 {object} entities.User
// @Router /users/me [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUserDetails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = userID.(uuid.UUID)
	if err := h.userService.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
