package postgresql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/quantonganh/ssr"
)

const (
	sqlInsertScan = `INSERT INTO scan (id, status, repository_id, findings, queued_at, scanning_at, finished_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
	sqlSelectScan = `SELECT id, status, repository_id, findings, queued_at, scanning_at, finished_at FROM scan WHERE id = $1`
	sqlSelectAllScans = `SELECT id, status, repository_id, findings, queued_at, scanning_at, finished_at FROM scan`
	sqlUpdateScan = `UPDATE scan SET status = $1, findings = $2 WHERE id = $3 RETURNING *`
	sqlDeleteScan = `DELETE FROM scan WHERE id = $1`
)

type scanService struct {
	db *sql.DB
}

func NewScanService(db *sql.DB) ssr.ScanService {
	return &scanService{
		db: db,
	}
}

func (ss *scanService) CreateScan(s *ssr.Scan) (*ssr.Scan, error) {
	var (
		id uuid.UUID
		status ssr.Status
		repoID uint64
		findings ssr.Findings
		queuedAt time.Time
		scanningAt time.Time
		finishedAt time.Time
	)
	err := ss.db.QueryRow(sqlInsertScan, s.ID, s.Status, s.RepositoryID, s.Findings, s.QueuedAt, s.ScanningAt, s.FinishedAt).Scan(&id, &status, &repoID, &findings, &queuedAt, &scanningAt, &finishedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scan")
	}
	return &ssr.Scan{
		ID:           id,
		Status:       status,
		RepositoryID: repoID,
		Findings:     findings,
		QueuedAt:     queuedAt,
		ScanningAt:   scanningAt,
		FinishedAt:   finishedAt,
	}, nil
}

func (ss *scanService) GetScan(id uuid.UUID) (*ssr.Scan, error) {
	var (
		status ssr.Status
		repoID uint64
		findings ssr.Findings
		queuedAt time.Time
		scanningAt time.Time
		finishedAt time.Time
	)
	err := ss.db.QueryRow(sqlSelectScan, id).Scan(&id, &status, &repoID, &findings, &queuedAt, &scanningAt, &finishedAt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select scan: %s", id)
	}

	return &ssr.Scan{
		ID: id,
		Status:     status,
		RepositoryID: repoID,
		Findings:   findings,
		QueuedAt:   queuedAt,
		ScanningAt: scanningAt,
		FinishedAt: finishedAt,
	}, nil
}

func (ss *scanService) ListScans() ([]*ssr.Scan, error) {
	rows, err := ss.db.Query(sqlSelectAllScans)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select all scans")
	}
	defer rows.Close()

	var scans []*ssr.Scan
	for rows.Next() {
		var (
			id uuid.UUID
			status       ssr.Status
			repositoryID string
			findings     ssr.Findings
			queuedAt time.Time
			scanningAt time.Time
			finishedAt time.Time
		)
		err := rows.Scan(&id, &status, &repositoryID, &findings, &queuedAt, &scanningAt, &finishedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan rows")
		}
		scans = append(scans, &ssr.Scan{
			ID: id,
			Status:     status,
			Findings:   findings,
			QueuedAt:   queuedAt,
			ScanningAt: scanningAt,
			FinishedAt: finishedAt,
		})
	}
	return scans, nil
}

func (ss *scanService) UpdateScan(id uuid.UUID, status ssr.Status, findings ssr.Findings) (*ssr.Scan, error) {
	var (
		repoID uint64
		queuedAt time.Time
		scanningAt time.Time
		finishedAt time.Time
	)
	err := ss.db.QueryRow(sqlUpdateScan, status, findings, id).Scan(&id, &status, &repoID, &findings, &queuedAt, &scanningAt, &finishedAt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update scan: %s", id)
	}
	return &ssr.Scan{
		ID:           id,
		Status:       status,
		RepositoryID: repoID,
		Findings:     findings,
		QueuedAt:     queuedAt,
		ScanningAt:   scanningAt,
		FinishedAt:   finishedAt,
	}, nil
}

func (ss *scanService) DeleteScan(id uuid.UUID) error {
	_, err := ss.db.Exec(sqlDeleteScan, id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete scan: %s", id)
	}
	return nil
}
