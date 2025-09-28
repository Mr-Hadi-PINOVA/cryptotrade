package router

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "cryptotrade/internal/config"
    "cryptotrade/internal/handler"
)

// SetupRouter configures the HTTP routes and middleware stack.
func SetupRouter(cfg config.Config, productHandler *handler.ProductHandler, userHandler *handler.UserHandler, orderHandler *handler.OrderHandler) *gin.Engine {
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":      "ok",
            "environment": cfg.Environment,
        })
    })

    api := r.Group("/api/v1")
    productHandler.RegisterRoutes(api)
    userHandler.RegisterRoutes(api)
    orderHandler.RegisterRoutes(api)

    return r
}
