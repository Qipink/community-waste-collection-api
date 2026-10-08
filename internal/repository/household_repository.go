package repository

import (
	"context"
	"errors"

	"community-waste-collection-api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HouseholdRepository interface {
	Create(ctx context.Context, h *entity.Household) error
	FindAll(ctx context.Context, limit, offset int) ([]entity.Household, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Household, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type householdRepository struct {
	db *gorm.DB
}

func NewHouseholdRepository(db *gorm.DB) HouseholdRepository {
	return &householdRepository{db: db}
}

func (r *householdRepository) Create(ctx context.Context, h *entity.Household) error {
	if err := r.db.WithContext(ctx).Create(h).Error; err != nil {
		return entity.ErrDatabase
	}
	return nil
}

func (r *householdRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.Household, int64, error) {
	var items []entity.Household
	var total int64
	if err := r.db.WithContext(ctx).Model(&entity.Household{}).Count(&total).Error; err != nil {
		return nil, 0, entity.ErrDatabase
	}
	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, entity.ErrDatabase
	}
	return items, total, nil
}

func (r *householdRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Household, error) {
	var h entity.Household
	if err := r.db.WithContext(ctx).First(&h, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrHouseholdNotFound
		}
		return nil, entity.ErrDatabase
	}
	return &h, nil
}

func (r *householdRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&entity.Household{}, "id = ?", id)
	if res.Error != nil {
		return entity.ErrDatabase
	}
	if res.RowsAffected == 0 {
		return entity.ErrHouseholdNotFound
	}
	return nil
}
