package postgresql

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/quantonganh/ssr"
)

const (
	sqlInsertScan = `INSERT INTO scan (status, repository_id, findings, queued_at, scanning_at, finished_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	sqlSelectScan = `SELECT id, status, repository_id, findings, queued_at, scanning_at, finished_at FROM scan WHERE id = $1`
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
	err := ss.db.QueryRow(sqlInsertScan, s.Status, s.RepositoryID, s.Findings, s.QueuedAt, s.ScanningAt, s.FinishedAt).Scan(&id, &status, &repoID, &findings, &queuedAt, &scanningAt, &finishedAt)
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

func (ss *scanService) ListScans(param ssr.FetchParam) (scans []*ssr.Scan, nextCursor string, err error) {
	queryBuilder := squirrel.Select("id", "status", "repository_id", "findings", "queued_at", "scanning_at", "finished_at").From("scan").PlaceholderFormat(squirrel.Dollar).OrderBy("finished_at DESC, id DESC")
	if param.Limit > 0 {
		queryBuilder = queryBuilder.Limit(param.Limit)
	}

	if param.Cursor != "" {
		finishedAt, id, decodeErr := decodeCursor(param.Cursor)
		if decodeErr != nil {
			err = errors.New("invalid cursor")
			return
		}

		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{
			"finished_at": finishedAt,
		})
		queryBuilder = queryBuilder.Where(squirrel.Lt{
			"id": id,
		})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return
	}

	rows, err := ss.db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	scans = []*ssr.Scan{}
	var finishedAt time.Time
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
		err = rows.Scan(&id, &status, &repositoryID, &findings, &queuedAt, &scanningAt, &finishedAt)
		if err != nil {
			return
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

	if len(scans) > 0 {
		nextCursor = encodeCursor(finishedAt, scans[len(scans)-1].ID)
	}

	return
}

func decodeCursor(encodedCursor string) (finishedAt time.Time, id uuid.UUID, err error) {
	b, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}

	cursors := strings.Split(string(b), ",")
	if len(cursors) != 2 {
		err = errors.New("invalid cursor")
		return
	}

	finishedAt, err = time.Parse(time.RFC3339Nano, cursors[0])
	if err != nil {
		return
	}
	id, err = uuid.Parse(cursors[1])
	if err != nil {
		return
	}

	return
}

func encodeCursor(finishedAt time.Time, id uuid.UUID) string {
	key := fmt.Sprintf("%s,%s", finishedAt.Format(time.RFC3339Nano), id)
	return base64.StdEncoding.EncodeToString([]byte(key))
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
