package service

import (
	"context"
	"time"

	"community-waste-collection-api/internal/entity"
	"community-waste-collection-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PickupService interface {
	Create(ctx context.Context, input entity.CreatePickupInput) (*entity.WastePickup, error)
	List(ctx context.Context, status string, householdID string, page, limit int) ([]entity.WastePickup, int64, error)
	Schedule(ctx context.Context, id string, input entity.SchedulePickupInput) (*entity.WastePickup, error)
	Complete(ctx context.Context, id string) (*entity.WastePickup, *entity.Payment, error)
	Cancel(ctx context.Context, id string) (*entity.WastePickup, error)
}

type pickupService struct {
	pickups    repository.PickupRepository
	payments   repository.PaymentRepository
	households repository.HouseholdRepository
}

func NewPickupService(pickups repository.PickupRepository, payments repository.PaymentRepository, households repository.HouseholdRepository) PickupService {
	return &pickupService{pickups: pickups, payments: payments, households: households}
}

func (s *pickupService) Create(ctx context.Context, input entity.CreatePickupInput) (*entity.WastePickup, error) {
	householdID, err := uuid.Parse(input.HouseholdID)
	if err != nil {
		return nil, entity.ErrInvalidUUID
	}
	if !input.Type.IsValid() {
		return nil, entity.ErrValidation
	}
	if _, err := s.households.FindByID(ctx, householdID); err != nil {
		return nil, err
	}
	// Aturan 1: tolak jika ada payment pending.
	pending, err := s.payments.CountPendingByHousehold(ctx, householdID)
	if err != nil {
		return nil, err
	}
	if pending > 0 {
		return nil, entity.ErrHasPendingPayment
	}
	p := &entity.WastePickup{
		HouseholdID: householdID,
		Type:        input.Type,
		Status:      entity.PickupPending,
		SafetyCheck: input.SafetyCheck,
	}
	if err := s.pickups.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *pickupService) List(ctx context.Context, status string, householdID string, page, limit int) ([]entity.WastePickup, int64, error) {
	page, limit = NormalizePagination(page, limit)
	var hid *uuid.UUID
	if householdID != "" {
		parsed, err := uuid.Parse(householdID)
		if err != nil {
			return nil, 0, entity.ErrInvalidUUID
		}
		hid = &parsed
	}
	if status != "" && !entity.PickupStatus(status).IsValid() {
		return nil, 0, entity.ErrValidation
	}
	return s.pickups.FindAll(ctx, status, hid, limit, (page-1)*limit)
}

func parsePickupDate(v string) (time.Time, error) {
	for _, f := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(f, v); err == nil {
			return t, nil
		}
	}
	return time.Time{}, entity.ErrValidation
}

func (s *pickupService) Schedule(ctx context.Context, id string, input entity.SchedulePickupInput) (*entity.WastePickup, error) {
	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, entity.ErrInvalidUUID
	}
	p, err := s.pickups.FindByID(ctx, pid)
	if err != nil {
		return nil, err
	}
	// Aturan 2: hanya dari pending.
	if p.Status != entity.PickupPending {
		return nil, entity.ErrPickupNotPending
	}
	date, err := parsePickupDate(input.PickupDate)
	if err != nil {
		return nil, err
	}
	safety := p.SafetyCheck
	if input.SafetyCheck != nil {
		safety = input.SafetyCheck
	}
	// Aturan 3: electronic wajib safety_check true.
	if p.Type == entity.WasteElectronic {
		if safety == nil || !*safety {
			return nil, entity.ErrSafetyCheckRequired
		}
	}
	p.SafetyCheck = safety
	p.PickupDate = &date
	p.Status = entity.PickupScheduled
	if err := s.pickups.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *pickupService) Complete(ctx context.Context, id string) (*entity.WastePickup, *entity.Payment, error) {
	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, nil, entity.ErrInvalidUUID
	}
	p, err := s.pickups.FindByID(ctx, pid)
	if err != nil {
		return nil, nil, err
	}
	if p.Status == entity.PickupCompleted || p.Status == entity.PickupCanceled {
		return nil, nil, entity.ErrPickupAlreadyFinal
	}
	if p.Status != entity.PickupScheduled {
		return nil, nil, entity.ErrPickupNotScheduled
	}
	p.Status = entity.PickupCompleted

	// Aturan 4: auto-create payment dalam transaksi.
	var payment *entity.Payment
	err = s.pickups.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Save(p).Error; err != nil {
			return entity.ErrDatabase
		}
		payment = &entity.Payment{
			HouseholdID: p.HouseholdID,
			WasteID:     p.ID,
			Amount:      entity.PaymentAmountForType(p.Type),
			Status:      entity.PaymentPending,
		}
		if err := tx.WithContext(ctx).Create(payment).Error; err != nil {
			return entity.ErrDatabase
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return p, payment, nil
}

func (s *pickupService) Cancel(ctx context.Context, id string) (*entity.WastePickup, error) {
	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, entity.ErrInvalidUUID
	}
	p, err := s.pickups.FindByID(ctx, pid)
	if err != nil {
		return nil, err
	}
	if p.Status == entity.PickupCompleted || p.Status == entity.PickupCanceled {
		return nil, entity.ErrPickupAlreadyFinal
	}
	p.Status = entity.PickupCanceled
	if err := s.pickups.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
