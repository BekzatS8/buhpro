package models

import (
	"encoding/json"
	"time"
)

type Bid struct {
	ID               string          `json:"id" db:"id"`
	OrderID          string          `json:"order_id" db:"order_id"`
	ExecutorID       string          `json:"executor_id" db:"executor_id"`
	CoverText        string          `json:"cover_text" db:"cover_text"`
	Price            *int64          `json:"price,omitempty" db:"price"`
	ProposedDeadline *time.Time      `json:"proposed_deadline,omitempty" db:"proposed_deadline"`
	Attachments      json.RawMessage `json:"attachments,omitempty" db:"attachments"` // jsonb (array/object)
	Status           string          `json:"status" db:"status"`
	PaidAt           *time.Time      `json:"paid_at,omitempty" db:"paid_at"`
	VisibleToClient  bool            `json:"visibility_to_client" db:"visibility_to_client"`
	Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}
