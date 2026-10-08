package repository

import (
	"context"
	"errors"
	"time"

	"community-waste-collection-api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentFilter struct {
	Status      string
	HouseholdID *uuid.UUID
	StartDate   *time.Time
	EndDate     *time.Time
}

type PaymentRepository interface {
	Create(ctx context.Context, p *entity.Payment) error
	FindAll(ctx context.Context, f PaymentFilter, limit, offset int) ([]entity.Payment, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	Update(ctx context.Context, p *entity.Payment) error
	CountPendingByHousehold(ctx context.Context, householdID uuid.UUID) (int64, error)
	DB() *gorm.DB
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) DB() *gorm.DB { return r.db }

func (r *paymentRepository) Create(ctx context.Context, p *entity.Payment) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return entity.ErrDatabase
	}
	return nil
}

func applyPaymentFilter(q *gorm.DB, f PaymentFilter) *gorm.DB {
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.HouseholdID != nil {
		q = q.Where("household_id = ?", *f.HouseholdID)
	}
	if f.StartDate != nil {
		q = q.Where("created_at >= ?", *f.StartDate)
	}
	if f.EndDate != nil {
		q = q.Where("created_at <= ?", *f.EndDate)
	}
	return q
}

func (r *paymentRepository) FindAll(ctx context.Context, f PaymentFilter, limit, offset int) ([]entity.Payment, int64, error) {
	var items []entity.Payment
	var total int64
	if err := applyPaymentFilter(r.db.WithContext(ctx).Model(&entity.Payment{}), f).Count(&total).Error; err != nil {
		return nil, 0, entity.ErrDatabase
	}
	if err := applyPaymentFilter(r.db.WithContext(ctx).Model(&entity.Payment{}), f).Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, entity.ErrDatabase
	}
	return items, total, nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	var p entity.Payment
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrPaymentNotFound
		}
		return nil, entity.ErrDatabase
	}
	return &p, nil
}

func (r *paymentRepository) Update(ctx context.Context, p *entity.Payment) error {
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return entity.ErrDatabase
	}
	return nil
}

func (r *paymentRepository) CountPendingByHousehold(ctx context.Context, householdID uuid.UUID) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&entity.Payment{}).
		Where("household_id = ? AND status = ?", householdID, entity.PaymentPending).
		Count(&total).Error; err != nil {
		return 0, entity.ErrDatabase
	}
	return total, nil
}
