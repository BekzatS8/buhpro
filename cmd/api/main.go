package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BekzatS8/buhpro/internal/transport/router"
	"github.com/BekzatS8/buhpro/pkg/config"
	"github.com/BekzatS8/buhpro/pkg/db"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	r := gin.Default()

	// register app routes (router will create repos/services/handlers)
	deps := &router.AppDeps{
		DB:  pool,
		Cfg: cfg,
	}
	router.RegisterRoutes(deps, r)

	// healthz (keep)
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    cfg.AppAddr,
		Handler: r,
	}

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
