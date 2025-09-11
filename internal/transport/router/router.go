package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/BekzatS8/buhpro/internal/middleware"
	"github.com/BekzatS8/buhpro/internal/repository"
	httpHandlers "github.com/BekzatS8/buhpro/internal/transport/http"
	"github.com/BekzatS8/buhpro/internal/usecase"
	"github.com/BekzatS8/buhpro/pkg/config"
)

// AppDeps carries minimal app dependencies
type AppDeps struct {
	DB  *pgxpool.Pool
	Cfg *config.Config
}

// RegisterRoutes wires repos -> services -> handlers and registers routes to gin.Engine
func RegisterRoutes(deps *AppDeps, r *gin.Engine) {
	// repositories
	userRepo := repository.NewUserRepo(deps.DB)
	orderRepo := repository.NewOrderRepo(deps.DB)
	bidRepo := repository.NewBidRepo(deps.DB)
	paymentRepo := repository.NewPaymentRepo(deps.DB)

	// usecases / services
	userUC := usecase.NewUserUsecase(userRepo, deps.Cfg.JWTSecret, deps.Cfg.JTTTLMin)
	orderSvc := usecase.NewOrderService(orderRepo, paymentRepo)
	bidSvc := usecase.NewBidService(bidRepo, paymentRepo)

	// handlers
	userHandler := httpHandlers.NewUserHandler(userUC)
	orderHandler := httpHandlers.NewOrderHandler(orderSvc)
	bidHandler := httpHandlers.NewBidHandler(bidSvc)

	// middleware
	authMw := middleware.AuthMiddleware(deps.Cfg.JWTSecret)

	api := r.Group("/api/v1")

	// Auth endpoints (we reuse userHandler.RegisterRoutes that registers register/login/count)
	auth := api.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		// If you later implement refresh/logout/email/phone flows - add them here.
	}

	// Users endpoints (protected)
	users := api.Group("/users")
	users.Use(authMw)
	{
		users.GET("/me", userHandler.Me)         // implement Me() in user handler if not yet
		users.PATCH("/me", userHandler.UpdateMe) // implement UpdateMe()
		users.GET("/count", userHandler.Count)   // you already have Count
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

	// Bids routes: create/list under orders, others top-level
	orderBids := api.Group("/orders/:id/bids")
	orderBids.Use(authMw)
	{
		orderBids.POST("", bidHandler.CreateBid)
		orderBids.GET("", bidHandler.ListByOrder)
	}

	bids := api.Group("/bids")
	bids.Use(authMw)
	{
		bids.GET("/:id", bidHandler.GetByID)
		bids.DELETE("/:id", bidHandler.Delete)
		bids.POST("/:id/pay", bidHandler.Pay)
		// add shortlist/win/lose endpoints later
	}
}
