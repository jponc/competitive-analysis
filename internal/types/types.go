package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type QueryJob struct {
	ID                    uuid.UUID  `db:"id"`
	Keyword               string     `db:"keyword"`
	CompletedAt           *time.Time `db:"completed_at"`
	ZenserpBatchID        *string    `db:"zenserp_batch_id"`
	ZenserpBatchProcessed bool       `db:"zenserp_batch_processed"`
	CreatedAt             time.Time  `db:"created_at"`
}

type QueryLocation struct {
	ID           uuid.UUID `db:"id"`
	QueryJobID   uuid.UUID `db:"query_job_id"`
	Device       string    `db:"device"`
	SearchEngine string    `db:"search_engine"`
	Num          string    `db:"num"`
	Country      string    `db:"country"`
	Location     string    `db:"location"`
	CreatedAt    time.Time `db:"created_at"`
}

type QueryItem struct {
	ID              uuid.UUID  `db:"id"`
	QueryJobID      uuid.UUID  `db:"query_job_id"`
	QueryLocationID uuid.UUID  `db:"query_location_id"`
	Position        string     `db:"position"`
	Title           string     `db:"title"`
	URL             string     `db:"url"`
	Body            string     `db:"body"`
	ProcessedAt     *time.Time `db:"processed_at"`
	CreatedAt       time.Time  `db:"created_at"`
}
