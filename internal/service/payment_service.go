package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"community-waste-collection-api/internal/entity"
	"community-waste-collection-api/internal/repository"

	"github.com/google/uuid"
)

type PaymentService interface {
	Create(ctx context.Context, input entity.CreatePaymentInput) (*entity.Payment, error)
	List(ctx context.Context, status string, householdID string, startDate, endDate string, page, limit int) ([]entity.Payment, int64, error)
	Confirm(ctx context.Context, id string, fileName string, file io.Reader) (*entity.Payment, error)
}

type paymentService struct {
	payments   repository.PaymentRepository
	households repository.HouseholdRepository
	pickups    repository.PickupRepository
	uploadDir  string
}

func NewPaymentService(payments repository.PaymentRepository, households repository.HouseholdRepository, pickups repository.PickupRepository, uploadDir string) PaymentService {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	return &paymentService{payments: payments, households: households, pickups: pickups, uploadDir: uploadDir}
}

func (s *paymentService) Create(ctx context.Context, input entity.CreatePaymentInput) (*entity.Payment, error) {
	hid, err := uuid.Parse(input.HouseholdID)
	if err != nil {
		return nil, entity.ErrInvalidUUID
	}
	wid, err := uuid.Parse(input.WasteID)
	if err != nil {
		return nil, entity.ErrInvalidUUID
	}
	if input.Amount <= 0 {
		return nil, entity.ErrValidation
	}
	if _, err := s.households.FindByID(ctx, hid); err != nil {
		return nil, err
	}
	pickup, err := s.pickups.FindByID(ctx, wid)
	if err != nil {
		return nil, err
	}
	if pickup.HouseholdID != hid {
		return nil, entity.ErrValidation
	}
	// Opsi A: nominal manual wajib sama dengan tarif tipe pickup
	// (electronic 100000, lainnya 50000) agar selaras dengan aturan 4.
	if expected := entity.PaymentAmountForType(pickup.Type); input.Amount != expected {
		return nil, fmt.Errorf("%w: tipe %s harus %.0f, bukan %.0f", entity.ErrPaymentAmountMismatch, pickup.Type, expected, input.Amount)
	}
	p := &entity.Payment{
		HouseholdID: hid,
		WasteID:     wid,
		Amount:      input.Amount,
		Status:      entity.PaymentPending,
	}
	if err := s.payments.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func parseDateParam(v string) (*time.Time, error) {
	if v == "" {
		return nil, nil
	}
	for _, f := range []string{"2006-01-02", time.RFC3339, "2006-01-02T15:04:05"} {
		if t, err := time.Parse(f, v); err == nil {
			return &t, nil
		}
	}
	return nil, entity.ErrValidation
}

func (s *paymentService) List(ctx context.Context, status string, householdID string, startDate, endDate string, page, limit int) ([]entity.Payment, int64, error) {
	page, limit = NormalizePagination(page, limit)
	var hid *uuid.UUID
	if householdID != "" {
		parsed, err := uuid.Parse(householdID)
		if err != nil {
			return nil, 0, entity.ErrInvalidUUID
		}
		hid = &parsed
	}
	if status != "" && !entity.PaymentStatus(status).IsValid() {
		return nil, 0, entity.ErrValidation
	}
	start, err := parseDateParam(startDate)
	if err != nil {
		return nil, 0, err
	}
	end, err := parseDateParam(endDate)
	if err != nil {
		return nil, 0, err
	}
	return s.payments.FindAll(ctx, repository.PaymentFilter{
		Status:      status,
		HouseholdID: hid,
		StartDate:   start,
		EndDate:     end,
	}, limit, (page-1)*limit)
}

var allowedExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

var allowedMIME = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func (s *paymentService) Confirm(ctx context.Context, id string, fileName string, file io.Reader) (*entity.Payment, error) {
	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, entity.ErrInvalidUUID
	}
	if strings.TrimSpace(fileName) == "" || file == nil {
		return nil, entity.ErrProofRequired
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	if !allowedExt[ext] {
		return nil, fmt.Errorf("%w: ekstensi %s tidak didukung (jpg/jpeg/png/webp)", entity.ErrValidation, ext)
	}
	head := make([]byte, 512)
	n, _ := file.Read(head)
	if n == 0 {
		return nil, fmt.Errorf("%w: file bukti kosong", entity.ErrValidation)
	}
	head = head[:n]
	if !allowedMIME[http.DetectContentType(head)] {
		return nil, fmt.Errorf("%w: isi file tidak sesuai ekstensi (terdeteksi %s)", entity.ErrValidation, http.DetectContentType(head))
	}
	file = io.MultiReader(bytes.NewReader(head), file)
	p, err := s.payments.FindByID(ctx, pid)
	if err != nil {
		return nil, err
	}
	if p.Status != entity.PaymentPending {
		return nil, entity.ErrPaymentNotPending
	}
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, entity.ErrDatabase
	}
	safe := fmt.Sprintf("%s_%s", uuid.NewString(), filepath.Base(fileName))
	dst := filepath.Join(s.uploadDir, safe)
	out, err := os.Create(dst)
	if err != nil {
		return nil, entity.ErrDatabase
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		return nil, entity.ErrDatabase
	}
	url := "/uploads/" + safe
	now := time.Now()
	p.ProofFileURL = &url
	p.Status = entity.PaymentPaid
	p.PaymentDate = &now
	if err := s.payments.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
