package router

import (
	httpHandlers "github.com/BekzatS8/buhpro/internal/transport/http"
	"github.com/gin-gonic/gin"
)

type RouteDeps struct {
	UserHandler  *httpHandlers.UserHandler
	OrderHandler *httpHandlers.OrderHandler
	BidHandler   *httpHandlers.BidHandler

	AuthMW gin.HandlerFunc
}

func RegisterRoutes(r *gin.Engine, deps *RouteDeps) {
	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", deps.UserHandler.Register)
		auth.POST("/login", deps.UserHandler.Login)
		auth.POST("/refresh", deps.UserHandler.Refresh)

		authProtected := auth.Group("")
		authProtected.Use(deps.AuthMW)
		{
			authProtected.POST("/logout", deps.UserHandler.Logout)
		}
	}
	users := api.Group("/users")
	users.Use(deps.AuthMW)
	{
		users.GET("/me", deps.UserHandler.Me)
		users.PATCH("/me", deps.UserHandler.UpdateMe)
		users.GET("count", deps.UserHandler.Count)
	}
	orders := api.Group("/orders")
	{
		orders.GET("", deps.OrderHandler.List)
		orders.GET("/:id", deps.OrderHandler.GetByID)

		orderAuth := orders.Group("")
		orderAuth.Use(deps.AuthMW)
		{
			orderAuth.POST("", deps.OrderHandler.Create)
			orderAuth.PATCH("/:id", deps.OrderHandler.Update)
			orderAuth.DELETE("/:id", deps.OrderHandler.Delete)

			orderAuth.POST("/:id/publish", deps.OrderHandler.Publish)
			orderAuth.POST("/:id/select-executor", deps.OrderHandler.SelectExecutor)
			orderAuth.POST("/:id/start", deps.OrderHandler.Start)
			orderAuth.POST("/:id/complete", deps.OrderHandler.Complete)
			orderAuth.POST("/:id/cancel", deps.OrderHandler.Cancel)
			orderAuth.GET("/:id/history", deps.OrderHandler.History)
		}
	}
	orderBids := api.Group("/bids")
	orderBids.Use(deps.AuthMW)
	{
		orderBids.POST("", deps.BidHandler.CreateBid)
		orderBids.GET("", deps.BidHandler.ListByOrder)
	}
	bids := api.Group("/bids")
	bids.Use(deps.AuthMW)
	{
		bids.GET("/:id", deps.BidHandler.GetByID)
		bids.DELETE("/:id", deps.BidHandler.Delete)
		bids.POST("/:id/pay", deps.BidHandler.Pay)
	}
}
