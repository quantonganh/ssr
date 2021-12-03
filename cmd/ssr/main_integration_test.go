// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/quantonganh/ssr"
)

var scanID uuid.UUID

func TestScanService(t *testing.T) {
	t.Run("create scan", testCreateScan)

	t.Run("get scan", testGetScan)

	t.Run("update scan", testUpdateScan)

	t.Run("delete scan", testDeleteScan)

	t.Run("list scans", testListScans)
}

func testCreateScan(t *testing.T) {
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
	body, err := json.Marshal(scan)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/scans/1", bytes.NewBuffer(body))
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scanResult ssr.Scan
	err = json.NewDecoder(resp.Body).Decode(&scanResult)
	require.NoError(t, err)
	scanID = scanResult.ID
}

func testGetScan(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/scans/%s", scanID), nil)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scanResult ssr.Scan
	err = json.NewDecoder(resp.Body).Decode(&scanResult)
	require.NoError(t, err)
	assert.Equal(t, ssr.InProgress, scanResult.Status)
}

func testUpdateScan(t *testing.T) {
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
	body, err := json.Marshal(scan)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8080/scans/%s", scanID), bytes.NewBuffer(body))
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scanResult ssr.Scan
	err = json.NewDecoder(resp.Body).Decode(&scanResult)
	assert.Equal(t, ssr.Success, scanResult.Status)
}

func testDeleteScan(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/scans/%s", scanID), nil)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testListScans(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/scans?limit=1&page=1", nil)
	require.NoError(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scans []*ssr.Scan
	err = json.NewDecoder(resp.Body).Decode(&scans)
	require.NoError(t, err)
	assert.Zero(t, len(scans))
}

