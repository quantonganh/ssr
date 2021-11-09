package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/quantonganh/ssr"
)

const (
	defaultLimit = 10
)

func (s *Server) CreateScanHandler(w http.ResponseWriter, r *http.Request) *appError {
	var scan ssr.Scan
	if err := json.NewDecoder(r.Body).Decode(&scan); err != nil {
		return &appError{
			Error:   err,
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
		}
	}

	scanResult, err := s.ScanService.CreateScan(&scan)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	response, err := json.Marshal(scanResult)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(response)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}

	}

	return nil
}

func (s *Server) GetScanHandler(w http.ResponseWriter, r *http.Request) *appError {
	scanID := mux.Vars(r)["scanID"]
	id, err := uuid.Parse(scanID)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	scan, err := s.ScanService.GetScan(id)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	response, err := json.Marshal(scan)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(response)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}

	}

	return nil
}

func (s *Server) UpdateScanHandler(w http.ResponseWriter, r *http.Request) *appError {
	scanID, err := uuid.Parse(mux.Vars(r)["scanID"])
	if err != nil {
		return &appError{
			Error:   err,
			Message: "Invalid scan ID",
			Code:    http.StatusBadRequest,
		}
	}

	var scan ssr.Scan
	if err := json.NewDecoder(r.Body).Decode(&scan); err != nil {
		return &appError{
			Error:   err,
			Message: "Invalid findings",
			Code:    http.StatusBadRequest,
		}
	}

	scanResult, err := s.ScanService.UpdateScan(scanID, scan.Status, scan.Findings)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusBadRequest,
		}
	}

	response, err := json.Marshal(scanResult)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}

	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(response)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}

	}

	return nil
}

func (s *Server) DeleteScanHandler(w http.ResponseWriter, r *http.Request) *appError {
	scanID, err := uuid.Parse(mux.Vars(r)["scanID"])
	if err != nil {
		return &appError{
			Error:   err,
			Message: "Invalid scan ID",
			Code:    http.StatusBadRequest,
		}
	}

	if err = s.ScanService.DeleteScan(scanID); err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	return nil
}

func (s *Server) ListScansHandler(w http.ResponseWriter, r *http.Request) *appError {
	limitStr := r.FormValue("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return &appError{
			Error:   err,
			Message: "invalid limit parameter",
			Code:    http.StatusInternalServerError,
		}
	}
	if limit == 0 {
		limit = defaultLimit
	}
	cursor := r.FormValue("cursor")
	param := ssr.FetchParam{
		Limit:  uint64(limit),
		Cursor: cursor,
	}

	scans, nextCursor, err := s.ScanService.ListScans(param)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	response, err := json.Marshal(scans)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	w.Header().Set("X-NextCursor", nextCursor)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(response)
	if err != nil {
		return &appError{
			Error:   err,
			Code:    http.StatusInternalServerError,
		}

	}

	return nil
}
