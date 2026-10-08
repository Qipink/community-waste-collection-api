package entity

import (
	"time"

	"github.com/google/uuid"
)

type WasteType string
type PickupStatus string

const (
	WasteOrganic    WasteType = "organic"
	WastePlastic    WasteType = "plastic"
	WastePaper      WasteType = "paper"
	WasteElectronic WasteType = "electronic"
)

const (
	PickupPending   PickupStatus = "pending"
	PickupScheduled PickupStatus = "scheduled"
	PickupCompleted PickupStatus = "completed"
	PickupCanceled  PickupStatus = "canceled"
)

func (t WasteType) IsValid() bool {
	switch t {
	case WasteOrganic, WastePlastic, WastePaper, WasteElectronic:
		return true
	}
	return false
}

func (s PickupStatus) IsValid() bool {
	switch s {
	case PickupPending, PickupScheduled, PickupCompleted, PickupCanceled:
		return true
	}
	return false
}

type WastePickup struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	HouseholdID uuid.UUID   `gorm:"column:household_id;type:uuid;not null;index" json:"household_id"`
	Type        WasteType    `gorm:"column:type;type:varchar(20);not null" json:"type"`
	Status      PickupStatus `gorm:"column:status;type:varchar(20);not null;default:pending" json:"status"`
	PickupDate  *time.Time   `gorm:"column:pickup_date" json:"pickup_date"`
	SafetyCheck *bool        `gorm:"column:safety_check" json:"safety_check"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (WastePickup) TableName() string { return "waste_pickups" }

type CreatePickupInput struct {
	HouseholdID string    `json:"household_id" binding:"required"`
	Type        WasteType `json:"type" binding:"required"`
	SafetyCheck *bool     `json:"safety_check"`
}

type SchedulePickupInput struct {
	PickupDate  string `json:"pickup_date" binding:"required"`
	SafetyCheck *bool  `json:"safety_check"`
}

func PaymentAmountForType(t WasteType) float64 {
	if t == WasteElectronic {
		return 100000
	}
	return 50000
}
