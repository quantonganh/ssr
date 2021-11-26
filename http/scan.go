package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/quantonganh/ssr"
)

const (
	defaultLimit = 10
)

func (s *Server) CreateScanHandler(w http.ResponseWriter, r *http.Request) error {
	var scan ssr.Scan
	if err := json.NewDecoder(r.Body).Decode(&scan); err != nil {
		return NewError(err, http.StatusBadRequest, "Bad request: invalid JSON")
	}

	scanResult, err := s.ScanService.CreateScan(&scan)
	if err != nil {
		return err
	}

	response, err := json.Marshal(scanResult)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal scan result")
	}
	_, err = w.Write(response)
	if err != nil {
		return errors.Wrapf(err, "failed to write response body")
	}

	return nil
}

func (s *Server) GetScanHandler(w http.ResponseWriter, r *http.Request) error {
	scanID := mux.Vars(r)["scanID"]
	id, err := uuid.Parse(scanID)
	if err != nil {
		return NewError(err, http.StatusBadRequest, "Bad request: invalid scan ID")
	}

	scan, err := s.ScanService.GetScan(id)
	if err != nil {
		return err
	}

	response, err := json.Marshal(scan)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal scan result")
	}

	_, err = w.Write(response)
	if err != nil {
		return errors.Wrapf(err, "failed to write response body")
	}

	return nil
}

func (s *Server) UpdateScanHandler(w http.ResponseWriter, r *http.Request) error {
	scanID, err := uuid.Parse(mux.Vars(r)["scanID"])
	if err != nil {
		return NewError(err, http.StatusBadRequest, "Bad request: invalid scan ID")
	}

	var scan ssr.Scan
	if err := json.NewDecoder(r.Body).Decode(&scan); err != nil {
		return NewError(err, http.StatusBadRequest, "Bad request: invalid findings")
	}

	scanResult, err := s.ScanService.UpdateScan(scanID, scan.Status, scan.Findings)
	if err != nil {
		return err
	}

	response, err := json.Marshal(scanResult)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal scan result")
	}

	_, err = w.Write(response)
	if err != nil {
		return errors.Wrapf(err, "failed to write response body")
	}

	return nil
}

func (s *Server) DeleteScanHandler(w http.ResponseWriter, r *http.Request) error {
	scanID, err := uuid.Parse(mux.Vars(r)["scanID"])
	if err != nil {
		return NewError(err, http.StatusBadRequest, "Bad request: invalid scan ID")
	}

	if err = s.ScanService.DeleteScan(scanID); err != nil {
		return err
	}

	return nil
}

func (s *Server) ListScansHandler(w http.ResponseWriter, r *http.Request) error {
	limitStr := r.FormValue("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return NewError(err, http.StatusBadRequest, "Bad request: invalid limit parameter")
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
		return err
	}

	response, err := json.Marshal(scans)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal scans result")
	}

	w.Header().Set("X-NextCursor", nextCursor)

	_, err = w.Write(response)
	if err != nil {
		return errors.Wrapf(err, "failed to write response body")
	}

	return nil
}
