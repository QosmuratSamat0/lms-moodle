package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/user"
	userUC "github.com/ap1-final-mini-moodle/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *userUC.Service
}

func NewUserHandler(service *userUC.Service) *UserHandler {
	return &UserHandler{service: service}
}

// Register registers a new user
// @Summary Register user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body struct{Email string `json:"email" binding:"required,email"`; Password string `json:"password" binding:"required,min=6"`; FirstName string `json:"first_name" binding:"required"`; LastName string `json:"last_name" binding:"required"`; Role string `json:"role" binding:"required"`} true "Registration Request"
// @Success 201 {object} user.User "Created user"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=6"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Role      string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.service.Register(&user.CreateUserInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      user.Role(req.Role),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}

// GetByID returns a user by ID
// @Summary Get user by ID
// @Description Returns user details by their unique ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.User "User details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	u, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// List returns a list of users
// @Summary List users
// @Description Returns a paginated list of users (Admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} user.User "Users list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := h.service.List(req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Update updates user details
// @Summary Update user
// @Description Update user first name, last name or active status
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body user.UpdateUserInput true "Update Request"
// @Success 200 {object} user.User "Updated user"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/users/{id} [patch]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Active    *bool   `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.service.Update(id, &user.UpdateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Active:    req.Active,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

// Login authenticates a user
// @Summary Login user (Old)
// @Description Authenticate user and return token (Legacy endpoint)
// @Tags users
// @Accept json
// @Produce json
// @Param request body struct{Email string `json:"email" binding:"required,email"`; Password string `json:"password" binding:"required"`} true "Login Request"
// @Success 200 {object} map[string]string "Token"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Delete deletes a user
// @Summary Delete user
// @Description Deletes a user by ID (Admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
