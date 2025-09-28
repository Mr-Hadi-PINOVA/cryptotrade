package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "cryptotrade/internal/domain"
    "cryptotrade/internal/service"
)

// OrderHandler exposes order endpoints.
type OrderHandler struct {
    service *service.OrderService
}

// NewOrderHandler constructs a new OrderHandler.
func NewOrderHandler(service *service.OrderService) *OrderHandler {
    return &OrderHandler{service: service}
}

// RegisterRoutes registers order routes on the provided router group.
func (h *OrderHandler) RegisterRoutes(rg *gin.RouterGroup) {
    rg.GET("/orders", h.listOrders)
    rg.GET("/orders/:id", h.getOrder)
    rg.POST("/orders", h.createOrder)
}

type orderItemRequest struct {
    ProductID string `json:"product_id" binding:"required"`
    Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type orderRequest struct {
    UserID string             `json:"user_id" binding:"required"`
    Items  []orderItemRequest `json:"items" binding:"required,dive"`
}

func (h *OrderHandler) createOrder(c *gin.Context) {
    var req orderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    items := make([]domain.OrderItem, 0, len(req.Items))
    for _, item := range req.Items {
        items = append(items, domain.OrderItem{ProductID: item.ProductID, Quantity: item.Quantity})
    }

    order, err := h.service.CreateOrder(c.Request.Context(), req.UserID, items)
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) listOrders(c *gin.Context) {
    orders, err := h.service.ListOrders(c.Request.Context())
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) getOrder(c *gin.Context) {
    order, err := h.service.GetOrder(c.Request.Context(), c.Param("id"))
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusOK, order)
}
