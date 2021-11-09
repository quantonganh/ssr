package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quantonganh/ssr"
	"github.com/quantonganh/ssr/mocks"
)

func TestScanHandler(t *testing.T) {
	scanID := uuid.New()
	finding := ssr.Finding{
		Type:     "sast",
		RuleID:   "G402",
		Location: ssr.Location{
			Path:      "scan.go",
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
	scan := &ssr.Scan{
		ID:           scanID,
		Status:       ssr.Success,
		RepositoryID: 1,
		Findings:     ssr.Findings{finding},
	}

	scanService := new(mocks.ScanService)
	scanService.On("CreateScan", scan).Return(scan, nil)
	scanService.On("GetScan", scanID).Return(scan, nil)
	scanService.On("UpdateScan", scanID, scan.Status, ssr.Findings{finding}).Return(scan, nil)
	scanService.On("DeleteScan", scanID).Return(nil)
	scanService.On("ListScans").Return([]*ssr.Scan{scan}, nil)

	t.Run("create scan", func(t *testing.T) {
		testCreateScanHandler(t, scan, scanService)
	})

	t.Run("get scan", func(t *testing.T) {
		testGetScanHandler(t, scanID, scanService)
	})

	t.Run("update scan", func(t *testing.T) {
		testUpdateScanHandler(t, scanID, scan, scanService)
	})

	t.Run("delete scan", func(t *testing.T) {
		testDeleteScanHandler(t, scanID, scanService)
	})

	t.Run("list scans", func(t *testing.T) {
		testListScansHandler(t, scanService)
	})
}

func testCreateScanHandler(t *testing.T, scan *ssr.Scan, scanService ssr.ScanService) {
	body, err := json.Marshal(scan)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/scans/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	s := NewServer(nil, scanService)
	s.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func testGetScanHandler(t *testing.T, scanID uuid.UUID, scanService ssr.ScanService) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/scans/%s", scanID), nil)
	rr := httptest.NewRecorder()
	s := NewServer(nil, scanService)
	s.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp ssr.Scan
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, uint64(1), resp.RepositoryID)
	assert.Equal(t, ssr.Success, resp.Status)
}

func testUpdateScanHandler(t *testing.T, scanID uuid.UUID, scan *ssr.Scan, scanService ssr.ScanService) {
	body, err := json.Marshal(scan)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/scans/%s", scanID), bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	s := NewServer(nil, scanService)
	s.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func testDeleteScanHandler(t *testing.T, scanID uuid.UUID, scanService ssr.ScanService) {
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/scans/%s", scanID), nil)
	rr := httptest.NewRecorder()
	s := NewServer(nil, scanService)
	s.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func testListScansHandler(t *testing.T, scanService ssr.ScanService) {
	req := httptest.NewRequest(http.MethodGet, "/scans", nil)
	rr := httptest.NewRecorder()
	s := NewServer(nil, scanService)
	s.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, http.StatusOK, rr.Code)
	var scans []*ssr.Scan
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&scans))
	assert.Equal(t, uint64(1), scans[0].RepositoryID)
	assert.Equal(t, ssr.Success, scans[0].Status)
}
