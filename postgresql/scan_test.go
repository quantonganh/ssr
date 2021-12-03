// +build !integration

package postgresql

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/quantonganh/ssr"
)

const (
	sqlInsertScan = `INSERT INTO "scan" ("id","status","repository_id","findings","queued_at","scanning_at","finished_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`
	sqlSelectScan = `SELECT * FROM "scan" WHERE id = $1 ORDER BY "scan"."id" LIMIT 1`
	sqlUpdateScan = `UPDATE "scan" SET "status"=$1,"findings"=$2 WHERE id = $3`
	sqlDeleteScan = `DELETE FROM "scan" WHERE "scan"."id" = $1`
	sqlListScans = `SELECT * FROM "scan" LIMIT 1`
)

var scanID uuid.UUID

func TestScanService(t *testing.T) {
	t.Run("create scan", testCreateScan)
	t.Run("get scan", func(t *testing.T) {
		testGetScan(t, scanID)
	})
	t.Run("list scans", func(t *testing.T) {
		testListScans(t, scanID)
	})
	t.Run("update scan", func(t *testing.T) {
		testUpdateScan(t, scanID)
	})
	t.Run("delete scan", func(t *testing.T) {
		testDeleteScan(t, scanID)
	})
}

func testCreateScan(t *testing.T) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

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
	now := time.Now()
	scan := &ssr.Scan{
		Status:       ssr.InProgress,
		RepositoryID: 1,
		Findings:     ssr.Findings{finding},
		QueuedAt:     now,
		ScanningAt:   now,
		FinishedAt:   now,
	}
	mock.ExpectExec(sqlInsertScan).
		WithArgs(sqlmock.AnyArg(), scan.Status, scan.RepositoryID, scan.Findings, scan.QueuedAt, scan.ScanningAt, scan.FinishedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	scanService := NewScanService(gormDB)
	scanResult, err := scanService.CreateScan(scan)
	require.NoError(t, err)
	scanID = scanResult.ID
}

func testGetScan(t *testing.T, scanID uuid.UUID) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	require.NoError(t, err)

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
	mock.ExpectQuery(sqlSelectScan).WithArgs(scanID).WillReturnRows(rows)

	scanService := NewScanService(gormDB)
	scanResult, err := scanService.GetScan(scanID)
	require.NoError(t, err)
	assert.Equal(t, ssr.InProgress, scanResult.Status)
	assert.Equal(t, "api.go", scanResult.Findings[0].Location.Path)
	assert.Equal(t, "G402", scanResult.Findings[0].RuleID)
}

func testListScans(t *testing.T, scanID uuid.UUID) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	require.NoError(t, err)

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
	mock.ExpectQuery(regexp.QuoteMeta(sqlListScans)).WillReturnRows(rows)

	scanService := NewScanService(gormDB)
	scans, err := scanService.ListScans(1, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(scans))
	assert.Equal(t, scanID, scans[0].ID)
	assert.Equal(t, ssr.InProgress, scans[0].Status)
}

func testUpdateScan(t *testing.T, scanID uuid.UUID) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

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
	mock.ExpectExec(sqlUpdateScan).WithArgs(scan.Status, ssr.Findings{finding}, scanID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	scanService := NewScanService(gormDB)
	scanResult, err := scanService.UpdateScan(scanID, scan.Status, []ssr.Finding{finding})
	require.NoError(t, err)
	assert.Equal(t, ssr.Success, scanResult.Status)
}

func testDeleteScan(t *testing.T, scanID uuid.UUID) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	mock.ExpectExec(regexp.QuoteMeta(sqlDeleteScan)).WithArgs(scanID).WillReturnResult(sqlmock.NewResult(1, 1))

	scanService := NewScanService(gormDB)
	require.NoError(t, scanService.DeleteScan(scanID))
}

