package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BekzatS8/buhpro/internal/middleware"
	"github.com/BekzatS8/buhpro/internal/repository"
	"github.com/BekzatS8/buhpro/internal/services"
	httpHandlers "github.com/BekzatS8/buhpro/internal/transport/http"
	"github.com/BekzatS8/buhpro/pkg/config"
)

// AppDeps carries minimal app dependencies
type AppDeps struct {
	DB  *pgxpool.Pool
	Cfg *config.Config
}

func RegisterRoutes(deps *AppDeps, r *gin.Engine) {
	// repos
	userRepo := repository.NewUserRepo(deps.DB)
	refreshRepo := repository.NewRefreshRepo(deps.DB)
	orderRepo := repository.NewOrderRepo(deps.DB)
	bidRepo := repository.NewBidRepo(deps.DB)
	paymentRepo := repository.NewPaymentRepo(deps.DB)

	// usecases
	userUC := services.NewUserUsecase(userRepo, refreshRepo, deps.Cfg.JWTSecret, deps.Cfg.JTTTLMin, deps.Cfg.RefreshTTLDays)
	orderSvc := services.NewOrderService(orderRepo, paymentRepo)
	bidSvc := services.NewBidService(bidRepo, paymentRepo)

	// handlers
	userHandler := httpHandlers.NewUserHandler(userUC)
	orderHandler := httpHandlers.NewOrderHandler(orderSvc)
	bidHandler := httpHandlers.NewBidHandler(bidSvc)

	// middleware
	authMw := middleware.AuthMiddleware(deps.Cfg.JWTSecret)

	api := r.Group("/api/v1")

	// Auth endpoints (register/login/refresh are public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.Refresh)
		// logout - require access token -> protect with authMw
		authProtected := auth.Group("")
		authProtected.Use(authMw)
		{
			authProtected.POST("/logout", userHandler.Logout)
		}
	}

	// Users endpoints (protected)
	users := api.Group("/users")
	users.Use(authMw)
	{
		users.GET("/me", userHandler.Me)
		users.PATCH("/me", userHandler.UpdateMe)
		users.GET("/count", userHandler.Count)
	}

	// Orders
	orders := api.Group("/orders")
	{
		// public list/get
		orders.GET("", orderHandler.List)
		orders.GET("/:id", orderHandler.GetByID)

		// protected actions
		ordersAuth := orders.Group("")
		ordersAuth.Use(authMw)
		{
			ordersAuth.POST("", orderHandler.Create)
			ordersAuth.PATCH("/:id", orderHandler.Update)
			ordersAuth.DELETE("/:id", orderHandler.Delete)

			ordersAuth.POST("/:id/publish", orderHandler.Publish)
			ordersAuth.POST("/:id/select-executor", orderHandler.SelectExecutor)
			ordersAuth.POST("/:id/start", orderHandler.Start)
			ordersAuth.POST("/:id/complete", orderHandler.Complete)
			ordersAuth.POST("/:id/cancel", orderHandler.Cancel)
			ordersAuth.GET("/:id/history", orderHandler.History)
		}
	}

	// Bids under orders (protected)
	orderBids := api.Group("/orders/:id/bids")
	orderBids.Use(authMw)
	{
		orderBids.POST("", bidHandler.CreateBid)
		orderBids.GET("", bidHandler.ListByOrder)
	}

	// Top-level bids (protected)
	bids := api.Group("/bids")
	bids.Use(authMw)
	{
		bids.GET("/:id", bidHandler.GetByID)
		bids.DELETE("/:id", bidHandler.Delete)
		bids.POST("/:id/pay", bidHandler.Pay)
	}
}
