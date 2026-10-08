package service

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"community-waste-collection-api/internal/entity"
	"community-waste-collection-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---------- fake repos ----------

type fakeHouseholdRepo struct {
	data map[uuid.UUID]*entity.Household
}

func (f *fakeHouseholdRepo) Create(ctx context.Context, h *entity.Household) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	f.data[h.ID] = h
	return nil
}

func (f *fakeHouseholdRepo) FindAll(ctx context.Context, limit, offset int) ([]entity.Household, int64, error) {
	return nil, 0, nil
}

func (f *fakeHouseholdRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Household, error) {
	if h, ok := f.data[id]; ok {
		return h, nil
	}
	return nil, entity.ErrHouseholdNotFound
}

func (f *fakeHouseholdRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

type fakePickupRepo struct {
	data map[uuid.UUID]*entity.WastePickup
	db   *gorm.DB
}

func (f *fakePickupRepo) Create(ctx context.Context, p *entity.WastePickup) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	f.data[p.ID] = p
	return nil
}

func (f *fakePickupRepo) FindAll(ctx context.Context, status string, hid *uuid.UUID, limit, offset int) ([]entity.WastePickup, int64, error) {
	return nil, 0, nil
}

func (f *fakePickupRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.WastePickup, error) {
	if p, ok := f.data[id]; ok {
		return p, nil
	}
	return nil, entity.ErrPickupNotFound
}

func (f *fakePickupRepo) Update(ctx context.Context, p *entity.WastePickup) error {
	f.data[p.ID] = p
	return nil
}

func (f *fakePickupRepo) DB() *gorm.DB { return f.db }

type fakePaymentRepo struct {
	pendingCount int64
	items        map[uuid.UUID]*entity.Payment
}

func (f *fakePaymentRepo) Create(ctx context.Context, p *entity.Payment) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if f.items == nil {
		f.items = map[uuid.UUID]*entity.Payment{}
	}
	f.items[p.ID] = p
	return nil
}

func (f *fakePaymentRepo) FindAll(ctx context.Context, filter repository.PaymentFilter, limit, offset int) ([]entity.Payment, int64, error) {
	return nil, 0, nil
}

func (f *fakePaymentRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	if p, ok := f.items[id]; ok {
		return p, nil
	}
	return nil, entity.ErrPaymentNotFound
}

func (f *fakePaymentRepo) Update(ctx context.Context, p *entity.Payment) error {
	f.items[p.ID] = p
	return nil
}

func (f *fakePaymentRepo) CountPendingByHousehold(ctx context.Context, hid uuid.UUID) (int64, error) {
	return f.pendingCount, nil
}

func (f *fakePaymentRepo) DB() *gorm.DB { return nil }

// ---------- tests ----------

func TestPaymentAmountForType(t *testing.T) {
	if got := entity.PaymentAmountForType(entity.WasteOrganic); got != 50000 {
		t.Fatalf("organic harus 50000, dapat %v", got)
	}
	if got := entity.PaymentAmountForType(entity.WasteElectronic); got != 100000 {
		t.Fatalf("electronic harus 100000, dapat %v", got)
	}
}

func TestPickupCreate_BlockedWhenPendingPayment(t *testing.T) {
	hid := uuid.New()
	svc := NewPickupService(
		&fakePickupRepo{data: map[uuid.UUID]*entity.WastePickup{}},
		&fakePaymentRepo{pendingCount: 1},
		&fakeHouseholdRepo{data: map[uuid.UUID]*entity.Household{hid: {ID: hid}}},
	)
	_, err := svc.Create(context.Background(), entity.CreatePickupInput{
		HouseholdID: hid.String(),
		Type:        entity.WasteOrganic,
	})
	if err != entity.ErrHasPendingPayment {
		t.Fatalf("aturan 1: harus ErrHasPendingPayment, dapat %v", err)
	}
}

func TestPickupSchedule_OnlyFromPending(t *testing.T) {
	pid := uuid.New()
	hid := uuid.New()
	repo := &fakePickupRepo{data: map[uuid.UUID]*entity.WastePickup{
		pid: {ID: pid, HouseholdID: hid, Type: entity.WasteOrganic, Status: entity.PickupScheduled},
	}}
	svc := NewPickupService(repo, &fakePaymentRepo{}, &fakeHouseholdRepo{data: map[uuid.UUID]*entity.Household{}})
	_, err := svc.Schedule(context.Background(), pid.String(), entity.SchedulePickupInput{PickupDate: "2026-10-10T10:00:00Z"})
	if err != entity.ErrPickupNotPending {
		t.Fatalf("aturan 2: harus ErrPickupNotPending, dapat %v", err)
	}
}

func TestPickupSchedule_ElectronicRequiresSafetyCheck(t *testing.T) {
	pid := uuid.New()
	hid := uuid.New()
	repo := &fakePickupRepo{data: map[uuid.UUID]*entity.WastePickup{
		pid: {ID: pid, HouseholdID: hid, Type: entity.WasteElectronic, Status: entity.PickupPending},
	}}
	svc := NewPickupService(repo, &fakePaymentRepo{}, &fakeHouseholdRepo{})
	_, err := svc.Schedule(context.Background(), pid.String(), entity.SchedulePickupInput{PickupDate: "2026-10-10T10:00:00Z"})
	if err != entity.ErrSafetyCheckRequired {
		t.Fatalf("aturan 3: harus ErrSafetyCheckRequired, dapat %v", err)
	}
}

func TestPaymentConfirm_RequiresFile(t *testing.T) {
	pid := uuid.New()
	hid := uuid.New()
	payRepo := &fakePaymentRepo{items: map[uuid.UUID]*entity.Payment{
		pid: {ID: pid, HouseholdID: hid, Status: entity.PaymentPending},
	}}
	svc := NewPaymentService(payRepo, &fakeHouseholdRepo{}, &fakePickupRepo{}, t.TempDir())
	_, err := svc.Confirm(context.Background(), pid.String(), "", nil)
	if err != entity.ErrProofRequired {
		t.Fatalf("aturan 5: harus ErrProofRequired, dapat %v", err)
	}
}

func TestPaymentConfirm_Success(t *testing.T) {
	pid := uuid.New()
	hid := uuid.New()
	payRepo := &fakePaymentRepo{items: map[uuid.UUID]*entity.Payment{
		pid: {ID: pid, HouseholdID: hid, Status: entity.PaymentPending},
	}}
	svc := NewPaymentService(payRepo, &fakeHouseholdRepo{}, &fakePickupRepo{}, t.TempDir())
	jpegHead := append([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00}, make([]byte, 600)...)
	got, err := svc.Confirm(context.Background(), pid.String(), "bukti.jpg", bytes.NewReader(jpegHead))
	if err != nil {
		t.Fatalf("confirm harus sukses, dapat %v", err)
	}
	if got.Status != entity.PaymentPaid || got.ProofFileURL == nil {
		t.Fatalf("status harus paid + proof terisi")
	}
}

func TestPaymentCreate_AmountMismatchElectronic(t *testing.T) {
	hid := uuid.New()
	wid := uuid.New()
	svc := NewPaymentService(
		&fakePaymentRepo{},
		&fakeHouseholdRepo{data: map[uuid.UUID]*entity.Household{hid: {ID: hid}}},
		&fakePickupRepo{data: map[uuid.UUID]*entity.WastePickup{
			wid: {ID: wid, HouseholdID: hid, Type: entity.WasteElectronic, Status: entity.PickupCompleted},
		}},
		t.TempDir(),
	)
	_, err := svc.Create(context.Background(), entity.CreatePaymentInput{
		HouseholdID: hid.String(),
		WasteID:     wid.String(),
		Amount:      50000,
	})
	if !errors.Is(err, entity.ErrPaymentAmountMismatch) {
		t.Fatalf("electronic + 50000 harus ErrPaymentAmountMismatch, dapat %v", err)
	}
}

func TestPaymentCreate_AmountMatchElectronic(t *testing.T) {
	hid := uuid.New()
	wid := uuid.New()
	svc := NewPaymentService(
		&fakePaymentRepo{},
		&fakeHouseholdRepo{data: map[uuid.UUID]*entity.Household{hid: {ID: hid}}},
		&fakePickupRepo{data: map[uuid.UUID]*entity.WastePickup{
			wid: {ID: wid, HouseholdID: hid, Type: entity.WasteElectronic, Status: entity.PickupCompleted},
		}},
		t.TempDir(),
	)
	got, err := svc.Create(context.Background(), entity.CreatePaymentInput{
		HouseholdID: hid.String(),
		WasteID:     wid.String(),
		Amount:      100000,
	})
	if err != nil {
		t.Fatalf("electronic + 100000 harus sukses, dapat %v", err)
	}
	if got.Amount != 100000 {
		t.Fatalf("amount harus 100000, dapat %v", got.Amount)
	}
}

func TestPaymentConfirm_RejectFakeMIME(t *testing.T) {
	pid := uuid.New()
	hid := uuid.New()
	payRepo := &fakePaymentRepo{items: map[uuid.UUID]*entity.Payment{
		pid: {ID: pid, HouseholdID: hid, Status: entity.PaymentPending},
	}}
	svc := NewPaymentService(payRepo, &fakeHouseholdRepo{}, &fakePickupRepo{}, t.TempDir())
	_, err := svc.Confirm(context.Background(), pid.String(), "bukti.jpg", bytes.NewReader([]byte("teks biasa bukan gambar")))
	if !errors.Is(err, entity.ErrValidation) {
		t.Fatalf("isi palsu harus ErrValidation, dapat %v", err)
	}
}

func TestPaymentConfirm_RejectEmptyFile(t *testing.T) {
	pid := uuid.New()
	hid := uuid.New()
	payRepo := &fakePaymentRepo{items: map[uuid.UUID]*entity.Payment{
		pid: {ID: pid, HouseholdID: hid, Status: entity.PaymentPending},
	}}
	svc := NewPaymentService(payRepo, &fakeHouseholdRepo{}, &fakePickupRepo{}, t.TempDir())
	_, err := svc.Confirm(context.Background(), pid.String(), "bukti.jpg", bytes.NewReader([]byte{}))
	if !errors.Is(err, entity.ErrValidation) {
		t.Fatalf("file kosong harus ErrValidation, dapat %v", err)
	}
}

func TestPaymentConfirm_RejectPDF(t *testing.T) {
	pid := uuid.New()
	hid := uuid.New()
	payRepo := &fakePaymentRepo{items: map[uuid.UUID]*entity.Payment{
		pid: {ID: pid, HouseholdID: hid, Status: entity.PaymentPending},
	}}
	svc := NewPaymentService(payRepo, &fakeHouseholdRepo{}, &fakePickupRepo{}, t.TempDir())
	pdfHead := append([]byte("%PDF-1.4 fake pdf content"), make([]byte, 600)...)
	_, err := svc.Confirm(context.Background(), pid.String(), "bukti.pdf", bytes.NewReader(pdfHead))
	if !errors.Is(err, entity.ErrValidation) {
		t.Fatalf("pdf harus ErrValidation, dapat %v", err)
	}
}
