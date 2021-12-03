package postgresql

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/quantonganh/ssr"
)

type scanService struct {
	db *gorm.DB
}

func NewScanService(db *gorm.DB) ssr.ScanService {
	return &scanService{
		db: db,
	}
}

func (ss *scanService) CreateScan(s *ssr.Scan) (*ssr.Scan, error) {
	result := ss.db.Create(&s)
	if err := result.Error; err != nil {
		return nil, errors.Wrap(err, "failed to create scan")
	}
	return s, nil
}

func (ss *scanService) GetScan(id uuid.UUID) (*ssr.Scan, error) {
	var s ssr.Scan
	if err := ss.db.First(&s, "id = ?", id).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to select scan: %s", id)
	}

	return &s, nil
}

func (ss *scanService) ListScans(page, limit int) (scans []*ssr.Scan, err error) {
	if err = ss.db.Scopes(paginate(page, limit)).Find(&scans).Error; err != nil {
		return
	}

	return
}

func paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func (ss *scanService) UpdateScan(id uuid.UUID, status ssr.Status, findings ssr.Findings) (*ssr.Scan, error) {
	var scan ssr.Scan
	if err := ss.db.Model(&scan).Where("id = ?", id).Updates(ssr.Scan{
		Status: status,
		Findings: findings,
	}).Error; err != nil {
		return nil, err
	}
	return &scan, nil
}

func (ss *scanService) DeleteScan(id uuid.UUID) error {
	var scan ssr.Scan
	if err := ss.db.Delete(&scan, id).Error; err != nil {
		return errors.Wrapf(err, "failed to delete scan: %s", id)
	}
	return nil
}
