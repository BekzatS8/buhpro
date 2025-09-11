package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/BekzatS8/buhpro/internal/models"
	"github.com/BekzatS8/buhpro/internal/repository"
)

type BidService struct {
	bidRepo     repository.BidRepo
	paymentRepo repository.PaymentRepo
}

func NewBidService(br repository.BidRepo, pr repository.PaymentRepo) *BidService {
	return &BidService{bidRepo: br, paymentRepo: pr}
}

func (s *BidService) Create(ctx context.Context, b *models.Bid) error {
	// prepare bid
	b.ID = uuid.NewString()
	b.Status = "pending_payment"
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now

	// insert bid
	if err := s.bidRepo.Create(ctx, b); err != nil {
		fmt.Printf("BidService.Create: bidRepo.Create error: %v\n", err)
		return err
	}
	fmt.Printf("BidService.Create: bid inserted OK, id=%s\n", b.ID)

	// prepare payment
	p := &models.Payment{
		ID:          uuid.NewString(),
		UserID:      &b.ExecutorID,
		RelatedType: "bid_fee",
		RelatedID:   &b.ID,
		Provider:    "mock",
		Amount:      500,
		Currency:    "KZT",
		Status:      "initiated",
	}

	fmt.Printf("BidService.Create: about to insert payment: id=%s related_type=%s related_id=%s user=%v\n", p.ID, p.RelatedType, p.RelatedID, p.UserID)

	// insert payment
	if err := s.paymentRepo.Create(ctx, p); err != nil {
		fmt.Printf("BidService.Create: paymentRepo.Create error: %v\n", err)
		// optional: rollback bid (delete) if you want in future
		return err
	}
	fmt.Printf("BidService.Create: payment inserted OK, id=%s\n", p.ID)

	return nil
}
func (s *BidService) Pay(ctx context.Context, bidID string) error {
	now := time.Now()
	if err := s.bidRepo.MarkPaid(ctx, bidID, now); err != nil {
		return err
	}
	return nil
}

// New methods required by handler:

func (s *BidService) ListByOrder(ctx context.Context, orderID string) ([]*models.Bid, error) {
	return s.bidRepo.ListByOrder(ctx, orderID)
}

func (s *BidService) GetByID(ctx context.Context, id string) (*models.Bid, error) {
	return s.bidRepo.GetByID(ctx, id)
}

func (s *BidService) Delete(ctx context.Context, id string) error {
	return s.bidRepo.Delete(ctx, id)
}
