package repository

import (
	"context"
	"errors"

	"community-waste-collection-api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PickupRepository interface {
	Create(ctx context.Context, p *entity.WastePickup) error
	FindAll(ctx context.Context, status string, householdID *uuid.UUID, limit, offset int) ([]entity.WastePickup, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.WastePickup, error)
	Update(ctx context.Context, p *entity.WastePickup) error
	DB() *gorm.DB
}

type pickupRepository struct {
	db *gorm.DB
}

func NewPickupRepository(db *gorm.DB) PickupRepository {
	return &pickupRepository{db: db}
}

func (r *pickupRepository) DB() *gorm.DB { return r.db }

func (r *pickupRepository) Create(ctx context.Context, p *entity.WastePickup) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return entity.ErrDatabase
	}
	return nil
}

func (r *pickupRepository) FindAll(ctx context.Context, status string, householdID *uuid.UUID, limit, offset int) ([]entity.WastePickup, int64, error) {
	var items []entity.WastePickup
	var total int64
	base := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&entity.WastePickup{})
		if status != "" {
			q = q.Where("status = ?", status)
		}
		if householdID != nil {
			q = q.Where("household_id = ?", *householdID)
		}
		return q
	}
	if err := base().Count(&total).Error; err != nil {
		return nil, 0, entity.ErrDatabase
	}
	if err := base().Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, entity.ErrDatabase
	}
	return items, total, nil
}

func (r *pickupRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.WastePickup, error) {
	var p entity.WastePickup
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrPickupNotFound
		}
		return nil, entity.ErrDatabase
	}
	return &p, nil
}

func (r *pickupRepository) Update(ctx context.Context, p *entity.WastePickup) error {
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return entity.ErrDatabase
	}
	return nil
}
