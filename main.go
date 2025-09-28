package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cryptotrade/internal/config"
	"cryptotrade/internal/handler"
	"cryptotrade/internal/repository/memory"
	"cryptotrade/internal/router"
	"cryptotrade/internal/service"
)

func main() {
	cfg := config.Load()

	productRepo := memory.NewProductRepository()
	userRepo := memory.NewUserRepository()
	orderRepo := memory.NewOrderRepository()

	productService := service.NewProductService(productRepo)
	userService := service.NewUserService(userRepo)
	orderService := service.NewOrderService(orderRepo, userRepo, productRepo)

	productHandler := handler.NewProductHandler(productService)
	userHandler := handler.NewUserHandler(userService)
	orderHandler := handler.NewOrderHandler(orderService)

	engine := router.SetupRouter(cfg, productHandler, userHandler, orderHandler)

	srv := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("starting server on %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
