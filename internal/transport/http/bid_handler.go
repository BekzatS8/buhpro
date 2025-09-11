package http

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"github.com/BekzatS8/buhpro/internal/domain"
	"github.com/BekzatS8/buhpro/internal/usecase"
	"github.com/gin-gonic/gin"
)

type BidHandler struct {
	svc *usecase.BidService
}

func NewBidHandler(s *usecase.BidService) *BidHandler { return &BidHandler{svc: s} }

func (h *BidHandler) RegisterRoutes(rg *gin.RouterGroup) {
	// Note: when router groups include params (orders/:id/bids) ensure group matches router.RegisterRoutes
	rg.POST("/orders/:id/bids", h.CreateBid)
	rg.GET("/orders/:id/bids", h.ListByOrder)
	rg.GET("/bids/:id", h.GetByID)
	rg.DELETE("/bids/:id", h.Delete)
	rg.POST("/bids/:id/pay", h.Pay)
}

func (h *BidHandler) CreateBid(c *gin.Context) {
	orderID := c.Param("id")

	// validate orderID
	if _, err := uuid.Parse(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	executorID, ok := uidVal.(string)
	if !ok || executorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in context"})
		return
	}
	fmt.Printf("CreateBid: executorID from ctx = %s\n", executorID)

	var req domain.Bid
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("CreateBid body: %+v\n", req)

	req.OrderID = orderID
	req.ExecutorID = executorID

	// call service and log exact error if any
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		// very important: print full error to stdout so we can see DB error text
		fmt.Printf("ERROR: BidService.Create failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *BidHandler) ListByOrder(c *gin.Context) {
	orderID := c.Param("id")
	list, err := h.svc.ListByOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *BidHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	b, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BidHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BidHandler) Pay(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Pay(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
