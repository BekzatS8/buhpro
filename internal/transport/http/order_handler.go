package http

import (
	"net/http"
	"strconv"

	"github.com/BekzatS8/buhpro/internal/domain"
	"github.com/BekzatS8/buhpro/internal/usecase"
	"github.com/gin-gonic/gin"
)

// structure depends on your wiring
type OrderHandler struct {
	svc *usecase.OrderService
}

func NewOrderHandler(s *usecase.OrderService) *OrderHandler { return &OrderHandler{svc: s} }

func (h *OrderHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.Create)
	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)
	rg.PATCH("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)

	rg.POST("/:id/publish", h.Publish)
	rg.POST("/:id/select-executor", h.SelectExecutor)
	rg.POST("/:id/start", h.Start)
	rg.POST("/:id/complete", h.Complete)
	rg.POST("/:id/cancel", h.Cancel)
	rg.GET("/:id/history", h.History)
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req domain.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// set client_user_id from auth context if available
	if uid, ok := c.Get("user_id"); ok {
		req.ClientUserID = uid.(string)
	}
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *OrderHandler) List(c *gin.Context) {
	filters := map[string]string{
		"status":     c.Query("status"),
		"category":   c.Query("category"),
		"region":     c.Query("region"),
		"min_budget": c.Query("min_budget"),
		"max_budget": c.Query("max_budget"),
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	list, total, err := h.svc.List(c.Request.Context(), filters, page, per)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": list, "total": total, "page": page, "per_page": per})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	o, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, o)
}

func (h *OrderHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req domain.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, req)
}

func (h *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

// Publish: client initiates payment for publish — returns created Payment (mock)
type publishReq struct {
	Amount int64 `json:"amount"`
}

func (h *OrderHandler) Publish(c *gin.Context) {
	id := c.Param("id")
	var req publishReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	uid, _ := c.Get("user_id")
	p, err := h.svc.Publish(c.Request.Context(), id, uid.(string), req.Amount)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, p)
}

type selectReq struct {
	BidID string `json:"bid_id" binding:"required"`
}

func (h *OrderHandler) SelectExecutor(c *gin.Context) {
	id := c.Param("id")
	var req selectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	actor := ""
	if uid, ok := c.Get("user_id"); ok {
		actor = uid.(string)
	}
	if err := h.svc.SelectExecutor(c.Request.Context(), id, req.BidID, actor); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func (h *OrderHandler) Start(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Start(c.Request.Context(), id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func (h *OrderHandler) Complete(c *gin.Context) {
	id := c.Param("id")
	actor := ""
	if uid, ok := c.Get("user_id"); ok {
		actor = uid.(string)
	}
	if err := h.svc.Complete(c.Request.Context(), id, actor); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	actor := ""
	if uid, ok := c.Get("user_id"); ok {
		actor = uid.(string)
	}
	if err := h.svc.Cancel(c.Request.Context(), id, actor); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func (h *OrderHandler) History(c *gin.Context) {
	// minimal: delegate to repo via service (not implemented returning payload)
	c.JSON(200, gin.H{"message": "history endpoint - implement reading audit_logs in repo/service"})
}
