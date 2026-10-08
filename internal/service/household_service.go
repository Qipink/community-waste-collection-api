package service

import (
	"context"
	"strings"

	"community-waste-collection-api/internal/entity"
	"community-waste-collection-api/internal/repository"

	"github.com/google/uuid"
)

type HouseholdService interface {
	CreateHousehold(ctx context.Context, input entity.CreateHouseholdInput) (*entity.Household, error)
	GetAllHouseholds(ctx context.Context, page, limit int) ([]entity.Household, int64, error)
	GetHouseholdByID(ctx context.Context, id string) (*entity.Household, error)
	DeleteHousehold(ctx context.Context, id string) error
}

type householdService struct {
	repo repository.HouseholdRepository
}

func NewHouseholdService(repo repository.HouseholdRepository) HouseholdService {
	return &householdService{repo: repo}
}

func (s *householdService) CreateHousehold(ctx context.Context, input entity.CreateHouseholdInput) (*entity.Household, error) {
	if strings.TrimSpace(input.OwnerName) == "" || strings.TrimSpace(input.Address) == "" {
		return nil, entity.ErrValidation
	}
	h := &entity.Household{
		OwnerName: strings.TrimSpace(input.OwnerName),
		Address:   strings.TrimSpace(input.Address),
	}
	if err := s.repo.Create(ctx, h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *householdService) GetAllHouseholds(ctx context.Context, page, limit int) ([]entity.Household, int64, error) {
	page, limit = NormalizePagination(page, limit)
	return s.repo.FindAll(ctx, limit, (page-1)*limit)
}

func (s *householdService) GetHouseholdByID(ctx context.Context, idStr string) (*entity.Household, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, entity.ErrInvalidUUID
	}
	return s.repo.FindByID(ctx, id)
}

func (s *householdService) DeleteHousehold(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return entity.ErrInvalidUUID
	}
	return s.repo.Delete(ctx, id)
}
