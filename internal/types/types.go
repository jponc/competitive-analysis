package types

import (
	"time"

	"github.com/gofrs/uuid"
)

type QueryJob struct {
	ID                    uuid.UUID  `db:"id" json:"id"`
	Keyword               string     `db:"keyword" json:"keyword"`
	CompletedAt           *time.Time `db:"completed_at" json:"completed_at"`
	ZenserpBatchID        *string    `db:"zenserp_batch_id" json:"zenserp_batch_id"`
	ZenserpBatchProcessed bool       `db:"zenserp_batch_processed" json:"zenserp_batch_processed"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
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
	Body            *string    `db:"body"`
	ProcessedAt     *time.Time `db:"processed_at"`
	CreatedAt       time.Time  `db:"created_at"`
	ErrorProcessing bool       `db:"error_processing"`
}

type QueryJobPositionHit struct {
	AvgPosition       float32 `db:"avg_position" json:"avg_position"`
	URL               string  `db:"url" json:"url"`
	LocationHitsCount int     `db:"location_hits_count" json:"location_hits_count"`
}

type UrlInfo struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Body  string `json:"body"`
	Links []Link `json:"links"`
}

type Link struct {
	Text string `db:"text" json:"text"`
	URL  string `db:"url" json:"url"`
}
