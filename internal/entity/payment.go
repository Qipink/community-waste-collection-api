package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
)

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentPending, PaymentPaid, PaymentFailed:
		return true
	}
	return false
}

type Payment struct {
	ID           uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	HouseholdID  uuid.UUID     `gorm:"column:household_id;type:uuid;not null;index" json:"household_id"`
	WasteID      uuid.UUID     `gorm:"column:waste_id;type:uuid;not null;index" json:"waste_id"`
	Amount       float64       `gorm:"column:amount;type:numeric(12,2);not null" json:"amount"`
	PaymentDate  *time.Time    `gorm:"column:payment_date" json:"payment_date"`
	Status       PaymentStatus `gorm:"column:status;type:varchar(20);not null;default:pending" json:"status"`
	ProofFileURL *string       `gorm:"column:proof_file_url;type:text" json:"proof_file_url"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

func (Payment) TableName() string { return "payments" }

type CreatePaymentInput struct {
	HouseholdID string  `json:"household_id" binding:"required"`
	WasteID     string  `json:"waste_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}
