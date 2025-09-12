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

// AppDeps carries minimal app dependencies (передаём в InitAndRegister)
type AppDeps struct {
	DB  *pgxpool.Pool
	Cfg *config.Config
}

// InitAndRegister — создаёт репозитории/сервисы/хендлеры и регистрирует роуты.
// Здесь — вся инициализация, а регистрация чисто вызывает RegisterRoutes (routes.go).
func InitAndRegister(deps *AppDeps, r *gin.Engine) {
	// repos
	userRepo := repository.NewUserRepo(deps.DB)
	refreshRepo := repository.NewRefreshRepo(deps.DB)
	orderRepo := repository.NewOrderRepo(deps.DB)
	bidRepo := repository.NewBidRepo(deps.DB)
	paymentRepo := repository.NewPaymentRepo(deps.DB)

	// usecases / services
	userUC := services.NewUserUsecase(userRepo, refreshRepo, deps.Cfg.JWTSecret, deps.Cfg.JTTTLMin, deps.Cfg.RefreshTTLDays)
	orderSvc := services.NewOrderService(orderRepo, paymentRepo)
	bidSvc := services.NewBidService(bidRepo, paymentRepo)

	// handlers (готовые для передачи в routes.go)
	userHandler := httpHandlers.NewUserHandler(userUC)
	orderHandler := httpHandlers.NewOrderHandler(orderSvc)
	bidHandler := httpHandlers.NewBidHandler(bidSvc)

	// middleware
	authMw := middleware.AuthMiddleware(deps.Cfg.JWTSecret)

	// собираем RouteDeps и регистрируем маршруты
	routeDeps := &RouteDeps{
		UserHandler:  userHandler,
		OrderHandler: orderHandler,
		BidHandler:   bidHandler,
		AuthMW:       authMw,
	}
	RegisterRoutes(r, routeDeps)
}
