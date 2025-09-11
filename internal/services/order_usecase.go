package services

import (
	"context"
	"time"

	"github.com/BekzatS8/buhpro/internal/models"
	"github.com/BekzatS8/buhpro/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo   repository.OrderRepo
	paymentRepo repository.PaymentRepo
}

func NewOrderService(or repository.OrderRepo, pr repository.PaymentRepo) *OrderService {
	return &OrderService{orderRepo: or, paymentRepo: pr}
}

func (s *OrderService) Create(ctx context.Context, o *models.Order) error {
	o.ID = uuid.NewString()
	o.Status = "draft"
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now
	return s.orderRepo.Create(ctx, o)
}

func (s *OrderService) GetByID(ctx context.Context, id string) (*models.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) List(ctx context.Context, filters map[string]string, page, perPage int) ([]*models.Order, int, error) {
	return s.orderRepo.List(ctx, filters, page, perPage)
}

func (s *OrderService) Update(ctx context.Context, o *models.Order) error {
	// ensure not published yet
	orig, err := s.orderRepo.GetByID(ctx, o.ID)
	if err != nil {
		return err
	}
	if orig.Status != "draft" && orig.Status != "pending_payment" {
		return ErrOrderImmutable
	}
	return s.orderRepo.Update(ctx, o)
}

func (s *OrderService) Delete(ctx context.Context, id string) error {
	orig, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if orig.Status == "published" || orig.Status == "executor_selected" || orig.Status == "in_progress" {
		return ErrOrderCannotDelete
	}
	return s.orderRepo.Delete(ctx, id)
}

var (
	ErrOrderImmutable    = &ServiceError{"order not editable in current state"}
	ErrOrderCannotDelete = &ServiceError{"order cannot be deleted in current state"}
)

type ServiceError struct{ Msg string }

func (e *ServiceError) Error() string { return e.Msg }

// Publish: create payment record and set order to PENDING_PAYMENT
func (s *OrderService) Publish(ctx context.Context, orderID string, payerID string, amount int64) (*models.Payment, error) {
	// create payment
	p := &models.Payment{
		ID:          uuid.NewString(),
		UserID:      &payerID,
		RelatedType: "order_publish",
		RelatedID:   &orderID,
		Provider:    "mock",
		Amount:      amount,
		Currency:    "KZT",
		Status:      "initiated",
	}
	if err := s.paymentRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	// set order pending payment
	if err := s.orderRepo.SetStatus(ctx, orderID, "pending_payment"); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *OrderService) SelectExecutor(ctx context.Context, orderID, bidID, actorID string) error {
	if err := s.orderRepo.SelectExecutor(ctx, orderID, bidID); err != nil {
		return err
	}
	// audit
	_ = s.orderRepo.AddHistory(ctx, actorID, "select_executor", "order", orderID, map[string]interface{}{"bid_id": bidID})
	return nil
}

func (s *OrderService) Start(ctx context.Context, orderID string) error {
	// minimal checks (in real: validate actor and chosen bid)
	return s.orderRepo.SetStatus(ctx, orderID, "in_progress")
}

func (s *OrderService) Complete(ctx context.Context, orderID, actorID string) error {
	_ = s.orderRepo.SetStatus(ctx, orderID, "client_review")
	_ = s.orderRepo.AddHistory(ctx, actorID, "complete_order", "order", orderID, nil)
	return nil
}

func (s *OrderService) Cancel(ctx context.Context, orderID, actorID string) error {
	_ = s.orderRepo.SetStatus(ctx, orderID, "cancelled")
	_ = s.orderRepo.AddHistory(ctx, actorID, "cancel_order", "order", orderID, nil)
	return nil
}
