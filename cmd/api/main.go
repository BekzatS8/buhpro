package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BekzatS8/buhpro/internal/repository"
	httpHandlers "github.com/BekzatS8/buhpro/internal/transport/http"
	"github.com/BekzatS8/buhpro/internal/usecase"
	"github.com/BekzatS8/buhpro/pkg/config"
	"github.com/BekzatS8/buhpro/pkg/db"
	"github.com/gin-gonic/gin"
	ginhttp "github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	router := ginhttp.Default()

	// build layers
	userRepo := repository.NewUserRepo(pool)
	userUC := usecase.NewUserUsecase(userRepo, cfg.JWTSecret, cfg.JTTTLMin)
	userHandler := httpHandlers.NewUserHandler(userUC)

	api := router.Group("/api/v1")
	users := api.Group("/users")
	userHandler.RegisterRoutes(users)

	router.GET("/healthz", func(c *ginhttp.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	srv := &http.Server{Addr: cfg.AppAddr, Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	fmt.Printf("Server started on %s\n", cfg.AppAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	fmt.Println("Server stopped")
}
