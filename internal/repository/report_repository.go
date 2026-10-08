package repository

import (
	"context"

	"gorm.io/gorm"
)

type WasteSummaryRow struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Total  int64  `json:"total"`
}

type PaymentStatusRow struct {
	Status string  `json:"status"`
	Total  int64   `json:"total"`
	Amount float64 `json:"amount"`
}

type ReportRepository interface {
	WasteSummary(ctx context.Context) ([]WasteSummaryRow, error)
	PaymentSummary(ctx context.Context) ([]PaymentStatusRow, float64, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) WasteSummary(ctx context.Context) ([]WasteSummaryRow, error) {
	var rows []WasteSummaryRow
	if err := r.db.WithContext(ctx).
		Table("waste_pickups").
		Select("type, status, COUNT(*) as total").
		Group("type, status").
		Order("type, status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []WasteSummaryRow{}
	}
	return rows, nil
}

func (r *reportRepository) PaymentSummary(ctx context.Context) ([]PaymentStatusRow, float64, error) {
	var rows []PaymentStatusRow
	if err := r.db.WithContext(ctx).
		Table("payments").
		Select("status, COUNT(*) as total, COALESCE(SUM(amount),0) as amount").
		Group("status").
		Order("status").
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	if rows == nil {
		rows = []PaymentStatusRow{}
	}
	var revenue float64
	if err := r.db.WithContext(ctx).
		Table("payments").
		Select("COALESCE(SUM(amount),0)").
		Where("status = ?", "paid").
		Scan(&revenue).Error; err != nil {
		return nil, 0, err
	}
	return rows, revenue, nil
}
