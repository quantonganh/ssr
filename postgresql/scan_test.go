// +build !integration

package postgresql

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quantonganh/ssr"
)

func TestScanService(t *testing.T) {
	t.Run("create scan", testCreateScan)
	t.Run("get scan", testGetScan)
	t.Run("list scans", testListScans)
	t.Run("update scan", testUpdateScan)
	t.Run("delete scan", testDeleteScan)
}

func testCreateScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	finding := ssr.Finding{
		Type:     "sast",
		RuleID:   "G402",
		Location: ssr.Location{
			Path:      "",
			Positions: ssr.Positions{
				Begin: ssr.Begin{
					Line: 60,
				},
			},
		},
		Metadata: ssr.Metadata{
			Description: "TLS InsecureSkipVerify set true.",
			Severity:    "HIGH",
		},
	}
	scanID := uuid.New()
	now := time.Now()
	scan := &ssr.Scan{
		ID:           scanID,
		Status:       ssr.InProgress,
		RepositoryID: 1,
		Findings:     ssr.Findings{finding},
		QueuedAt:     now,
		ScanningAt:   now,
		FinishedAt:   now,
	}
	rows := sqlmock.NewRows([]string{"id", "status", "repository_id", "findings", "queued_at", "scanning_at", "finished_at"}).AddRow(scan.ID, scan.Status, scan.RepositoryID, scan.Findings, scan.QueuedAt, scan.ScanningAt, scan.FinishedAt)
	mock.ExpectQuery(regexp.QuoteMeta(sqlInsertScan)).WithArgs(scan.Status, scan.RepositoryID, scan.Findings, scan.QueuedAt, scan.ScanningAt, scan.FinishedAt).WillReturnRows(rows)

	scanService := NewScanService(db)
	scanResult, err := scanService.CreateScan(scan)
	require.NoError(t, err)
	assert.Equal(t, scan.ID, scanResult.ID)
}

func testGetScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	finding := ssr.Finding{
		Type:     "sast",
		RuleID:   "G402",
		Location: ssr.Location{
			Path:      "api.go",
			Positions: ssr.Positions{
				Begin: ssr.Begin{
					Line: 60,
				},
			},
		},
		Metadata: ssr.Metadata{
			Description: "TLS InsecureSkipVerify set true.",
			Severity:    "HIGH",
		},
	}
	scanID := uuid.New()
	now := time.Now()
	scan := &ssr.Scan{
		ID: scanID,
		Status:       ssr.InProgress,
		RepositoryID: 1,
		Findings:     ssr.Findings{finding},
		QueuedAt:     now,
		ScanningAt:   now,
		FinishedAt:   now,
	}
	rows := sqlmock.NewRows([]string{"id", "status", "repository_id", "findings", "queued_at", "scanning_at", "finished_at"}).AddRow(scan.ID, scan.Status, scan.RepositoryID, scan.Findings, scan.QueuedAt, scan.ScanningAt, scan.FinishedAt)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectScan)).WithArgs(scanID).WillReturnRows(rows)

	scanService := NewScanService(db)
	scanResult, err := scanService.GetScan(scanID)
	require.NoError(t, err)
	assert.Equal(t, ssr.InProgress, scanResult.Status)
	assert.Equal(t, "api.go", scanResult.Findings[0].Location.Path)
	assert.Equal(t, "G402", scanResult.Findings[0].RuleID)
}

func testListScans(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	finding := ssr.Finding{
		Type:     "sast",
		RuleID:   "G402",
		Location: ssr.Location{
			Path:      "api.go",
			Positions: ssr.Positions{
				Begin: ssr.Begin{
					Line: 60,
				},
			},
		},
		Metadata: ssr.Metadata{
			Description: "TLS InsecureSkipVerify set true.",
			Severity:    "HIGH",
		},
	}
	scanID := uuid.New()
	now := time.Now()
	scan := &ssr.Scan{
		ID: scanID,
		Status:       ssr.InProgress,
		RepositoryID: 1,
		Findings:     []ssr.Finding{finding},
		QueuedAt:     now,
		ScanningAt:   now,
		FinishedAt:   now,
	}
	rows := sqlmock.NewRows([]string{"id", "status", "repository_id", "findings", "queued_at", "scanning_at", "finished_at"}).AddRow(scan.ID, scan.Status, scan.RepositoryID, scan.Findings, scan.QueuedAt, scan.ScanningAt, scan.FinishedAt)
	queryBuilder := squirrel.Select("id", "status", "repository_id", "findings", "queued_at", "scanning_at", "finished_at").From("scan").PlaceholderFormat(squirrel.Dollar).OrderBy("finished_at DESC, id DESC")
	query, _, err := queryBuilder.ToSql()
	require.NoError(t, err)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	scanService := NewScanService(db)
	scans, _, err := scanService.ListScans(ssr.FetchParam{
		Limit:  1,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, len(scans))
	assert.Equal(t, scanID, scans[0].ID)
	assert.Equal(t, ssr.InProgress, scans[0].Status)
}

func testUpdateScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	finding := ssr.Finding{
		Type:     "sast",
		RuleID:   "G402",
		Location: ssr.Location{
			Path:      "",
			Positions: ssr.Positions{
				Begin: ssr.Begin{
					Line: 60,
				},
			},
		},
		Metadata: ssr.Metadata{
			Description: "TLS InsecureSkipVerify set true.",
			Severity:    "HIGH",
		},
	}
	scanID := uuid.New()
	now := time.Now()
	scan := &ssr.Scan{
		ID:           scanID,
		Status:       ssr.Success,
		RepositoryID: 1,
		Findings:     ssr.Findings{finding},
		QueuedAt:     now,
		ScanningAt:   now,
		FinishedAt:   now,
	}
	rows := sqlmock.NewRows([]string{"id", "status", "repository_id", "findings", "queued_at", "scanning_at", "finished_at"}).AddRow(scan.ID, scan.Status, scan.RepositoryID, scan.Findings, scan.QueuedAt, scan.ScanningAt, scan.FinishedAt)
	mock.ExpectQuery(regexp.QuoteMeta(sqlUpdateScan)).WithArgs(scan.Status, ssr.Findings{finding}, scanID).WillReturnRows(rows)

	scanService := NewScanService(db)
	scanResult, err := scanService.UpdateScan(scanID, scan.Status, []ssr.Finding{finding})
	require.NoError(t, err)
	assert.Equal(t, ssr.Success, scanResult.Status)
}

func testDeleteScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	scanID := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta(sqlDeleteScan)).WithArgs(scanID).WillReturnResult(sqlmock.NewResult(1, 1))

	scanService := NewScanService(db)
	require.NoError(t, scanService.DeleteScan(scanID))
}

