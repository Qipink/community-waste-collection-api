package service

import (
	"context"

	"community-waste-collection-api/internal/repository"
)

type ReportService interface {
	WasteSummary(ctx context.Context) ([]repository.WasteSummaryRow, error)
	PaymentSummary(ctx context.Context) ([]repository.PaymentStatusRow, float64, int64, error)
}

type reportService struct {
	reports repository.ReportRepository
}

func NewReportService(r repository.ReportRepository) ReportService {
	return &reportService{reports: r}
}

func (s *reportService) WasteSummary(ctx context.Context) ([]repository.WasteSummaryRow, error) {
	return s.reports.WasteSummary(ctx)
}

func (s *reportService) PaymentSummary(ctx context.Context) ([]repository.PaymentStatusRow, float64, int64, error) {
	rows, revenue, err := s.reports.PaymentSummary(ctx)
	if err != nil {
		return nil, 0, 0, err
	}
	var total int64
	for _, r := range rows {
		total += r.Total
	}
	return rows, revenue, total, nil
}
