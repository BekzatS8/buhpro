package domain

import "time"

type Payment struct {
	ID                string                 `json:"id"`
	UserID            *string                `json:"user_id,omitempty"`
	OrganizationID    *string                `json:"organization_id,omitempty"`
	RelatedType       string                 `json:"related_type"`
	RelatedID         *string                `json:"related_id,omitempty"`
	Provider          string                 `json:"provider"`
	ProviderPaymentID *string                `json:"provider_payment_id,omitempty"`
	Amount            int64                  `json:"amount"`
	Currency          string                 `json:"currency"`
	Status            string                 `json:"status"`
	Items             map[string]interface{} `json:"items,omitempty"`
	IdempotencyKey    *string                `json:"idempotency_key,omitempty"`
	ExpiresAt         *time.Time             `json:"expires_at,omitempty"`
	WebhookMeta       map[string]interface{} `json:"webhook_meta,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}
