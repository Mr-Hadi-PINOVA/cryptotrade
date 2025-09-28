package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "cryptotrade/internal/domain"
    "cryptotrade/internal/service"
)

// UserHandler exposes user endpoints.
type UserHandler struct {
    service *service.UserService
}

// NewUserHandler constructs a UserHandler instance.
func NewUserHandler(service *service.UserService) *UserHandler {
    return &UserHandler{service: service}
}

// RegisterRoutes registers user routes on the provided router group.
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
    rg.GET("/users", h.listUsers)
    rg.POST("/users", h.createUser)
    rg.GET("/users/:id", h.getUser)
}

type userRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

func (h *UserHandler) createUser(c *gin.Context) {
    var req userRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.service.CreateUser(c.Request.Context(), domain.User{
        Name:  req.Name,
        Email: req.Email,
    })
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) listUsers(c *gin.Context) {
    users, err := h.service.ListUsers(c.Request.Context())
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusOK, users)
}

func (h *UserHandler) getUser(c *gin.Context) {
    user, err := h.service.GetUser(c.Request.Context(), c.Param("id"))
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusOK, user)
}
