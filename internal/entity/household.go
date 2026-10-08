package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Household struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OwnerName string         `gorm:"column:owner_name;type:varchar(255);not null" json:"owner_name" binding:"required"`
	Address   string         `gorm:"column:address;type:text;not null" json:"address" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Household) TableName() string { return "households" }

type CreateHouseholdInput struct {
	OwnerName string `json:"owner_name" binding:"required"`
	Address   string `json:"address" binding:"required"`
}
