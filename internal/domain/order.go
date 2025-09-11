package domain

import "time"

type Order struct {
	ID           string                 `json:"id"`
	OrgID        string                 `json:"org_id,omitempty"`
	ClientUserID string                 `json:"client_user_id,omitempty"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description,omitempty"`
	Category     string                 `json:"category,omitempty"`
	Subcategory  string                 `json:"subcategory,omitempty"`
	Region       string                 `json:"region,omitempty"`
	ModeOnline   bool                   `json:"mode_online"`
	Deadline     *time.Time             `json:"deadline,omitempty"`
	BudgetMin    *int64                 `json:"budget_min,omitempty"`
	BudgetMax    *int64                 `json:"budget_max,omitempty"`
	Currency     string                 `json:"currency,omitempty"`
	Status       string                 `json:"status"`
	Promotion    map[string]interface{} `json:"promotion,omitempty"` // JSONB
	Attachments  map[string]interface{} `json:"attachments,omitempty"`
	ChosenBidID  *string                `json:"chosen_bid_id,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	PublishedAt  *time.Time             `json:"published_at,omitempty"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
