package ssr

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status int64

const (
	Queued Status = iota
	InProgress
	Success
	Failure
)

func (s Status) String() string {
	return [...]string{"Queued", "In Progress", "Success", "Failure"}[s]
}

type Scan struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid"`
	Status     Status   `json:"status"`
	RepositoryID uint64 `json:"repository_id"`
	Findings   Findings `json:"findings"`
	QueuedAt time.Time `json:"queued_at"`
	ScanningAt time.Time `json:"scanning_at"`
	FinishedAt time.Time `json:"finished_at"`
	Repository Repository `gorm:"foreignKey:RepositoryID"`
}

func (Scan) TableName() string {
	return "scan"
}

func (s *Scan) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}

type Findings []Finding

type Finding struct {
	Type string `json:"type"`
	RuleID string `json:"rule_id"`
	Location Location `json:"location"`
	Metadata Metadata `json:"metadata"`
}

type Location struct {
	Path string `json:"path"`
	Positions Positions
}

type Positions struct {
	Begin Begin `json:"begin"`
}

type Begin struct {
	Line int64 `json:"line"`
}

type Metadata struct {
	Description string `json:"description"`
	Severity string `json:"severity"`
}

func (f Findings) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *Findings) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &f)
}

type ScanService interface {
	CreateScan(s *Scan) (*Scan, error)
	GetScan(id uuid.UUID) (*Scan, error)
	ListScans(page, limit int) (scans []*Scan, err error)
	UpdateScan(id uuid.UUID, status Status, findings Findings) (*Scan, error)
	DeleteScan(id uuid.UUID) error
}
