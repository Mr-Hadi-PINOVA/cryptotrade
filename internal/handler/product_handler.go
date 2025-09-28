package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cryptotrade/internal/domain"
	"cryptotrade/internal/service"
)

// ProductHandler exposes product endpoints.
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler constructs a ProductHandler instance.
func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// RegisterRoutes registers product routes on the provided router group.
func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/products", h.listProducts)
	rg.POST("/products", h.createProduct)
	rg.GET("/products/:id", h.getProduct)
	rg.PUT("/products/:id", h.updateProduct)
	rg.DELETE("/products/:id", h.deleteProduct)
}

type productRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
}

func (h *ProductHandler) createProduct(c *gin.Context) {
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) listProducts(c *gin.Context) {
	products, err := h.service.ListProducts(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) getProduct(c *gin.Context) {
	product, err := h.service.GetProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) updateProduct(c *gin.Context) {
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), c.Param("id"), domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) deleteProduct(c *gin.Context) {
	if err := h.service.DeleteProduct(c.Request.Context(), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
