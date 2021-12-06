package dbrepository

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jponc/competitive-analysis/internal/types"
	"github.com/jponc/competitive-analysis/pkg/postgres"
)

type Repository struct {
	dbClient *postgres.Client
}

func NewRepository(dbClient *postgres.Client) (*Repository, error) {
	if dbClient == nil {
		return nil, fmt.Errorf("failed to initialise repository: dbClient is nil")
	}

	r := &Repository{
		dbClient,
	}

	return r, nil
}

func (r *Repository) Connect() error {
	return r.dbClient.Connect()
}

func (r *Repository) CreateQueryJob(ctx context.Context, keyword string) (uuid.UUID, error) {
	if r.dbClient == nil {
		return uuid.Nil, fmt.Errorf("dbClient not initialised")
	}

	var id uuid.UUID

	err := r.dbClient.GetContext(
		ctx,
		&id,
		`
			INSERT INTO query_job (keyword)
			VALUES ($1)
			RETURNING id
		`,
		keyword)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert query job: %v", err)
	}

	return id, nil
}

func (r *Repository) GetQueryJob(ctx context.Context, id uuid.UUID) (*types.QueryJob, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	var queryJob types.QueryJob

	err := r.dbClient.GetContext(
		ctx,
		&queryJob,
		`
			SELECT *
			FROM query_job
			WHERE id = $1
		`,
		id)

	if err != nil {
		return nil, fmt.Errorf("failed to get query job: %v", err)
	}

	return &queryJob, nil
}

func (r *Repository) CreateQueryLocation(ctx context.Context, queryJobID, device, searchEngine, num, country, location string) (uuid.UUID, error) {
	if r.dbClient == nil {
		return uuid.Nil, fmt.Errorf("dbClient not initialised")
	}

	var id uuid.UUID

	err := r.dbClient.GetContext(
		ctx,
		&id,
		`
			INSERT INTO query_location (query_job_id, device, search_engine, num, country, location)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`,
		queryJobID, device, searchEngine, num, country, location)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert query location: %v", err)
	}

	return id, nil
}

func (r *Repository) GetQueryLocations(ctx context.Context, queryJobID uuid.UUID) (*[]types.QueryLocation, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	queryLocations := []types.QueryLocation{}

	err := r.dbClient.SelectContext(
		ctx,
		&queryLocations,
		`SELECT * FROM query_location WHERE query_job_id = $1`, queryJobID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get query locations: %w", err)
	}

	return &queryLocations, nil
}

func (r *Repository) SetZenserpBatchToQueryJob(ctx context.Context, queryJobID uuid.UUID, zenserpBatchID string) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.ExecContext(
		ctx,
		`
			UPDATE query_job
			SET zenserp_batch_id = $1
			WHERE id = $2
		`, zenserpBatchID, queryJobID,
	)

	if err != nil {
		return fmt.Errorf("failed to update queryjob with zenserp batch id: %w", err)
	}

	return nil
}

func (r *Repository) GetUnprocessedQueryJobs(ctx context.Context) (*[]types.QueryJob, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	queryJobs := []types.QueryJob{}

	err := r.dbClient.SelectContext(
		ctx,
		&queryJobs,
		`SELECT * FROM query_job WHERE zenserp_batch_processed = false AND zenserp_batch_id IS NOT NULL`,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get query locations: %w", err)
	}

	return &queryJobs, nil
}

func (r *Repository) ProcessQueryJob(ctx context.Context, queryJobID uuid.UUID) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.ExecContext(
		ctx,
		`
			UPDATE query_job
			SET zenserp_batch_processed = true
			WHERE id = $1
		`, queryJobID,
	)

	if err != nil {
		return fmt.Errorf("failed to update queryjob with zenserp batch id: %w", err)
	}

	return nil
}

func (r *Repository) CreateQueryItem(ctx context.Context, queryJobID uuid.UUID, queryLocationID uuid.UUID, position int, url, title string) (uuid.UUID, error) {
	if r.dbClient == nil {
		return uuid.Nil, fmt.Errorf("dbClient not initialised")
	}

	var id uuid.UUID

	err := r.dbClient.GetContext(
		ctx,
		&id,
		`
			INSERT INTO query_item (query_job_id, query_location_id, position, url, title)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, queryJobID, queryLocationID, position, url, title,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create query item: %w", err)
	}

	return id, nil
}

func (r *Repository) GetQueryItem(ctx context.Context, id uuid.UUID) (*types.QueryItem, error) {
	if r.dbClient == nil {
		return nil, fmt.Errorf("dbClient not initialised")
	}

	queryItem := types.QueryItem{}

	err := r.dbClient.GetContext(
		ctx,
		&queryItem,
		`
			SELECT * FROM query_item where id = $1
		`, id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get query item: %w", err)
	}

	return &queryItem, nil
}

func (r *Repository) SetQueryItemsErrorProcessing(ctx context.Context, queryJobID uuid.UUID, url string) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.ExecContext(
		ctx,
		`
			UPDATE query_item
			SET processed_at = now(), error_processing = true
			WHERE query_job_id = $1 and url = $2
		`, queryJobID, url,
	)

	if err != nil {
		return fmt.Errorf("failed to update query item error processing: %w", err)
	}

	return nil
}

func (r *Repository) SetQueryItemsProcessedWithBody(ctx context.Context, queryJobID uuid.UUID, url string, body string) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.ExecContext(
		ctx,
		`
			UPDATE query_item
			SET processed_at = now(), error_processing = false, body = $3
			WHERE query_job_id = $1 and url = $2
		`, queryJobID, url, body,
	)

	if err != nil {
		return fmt.Errorf("failed to update query item error processing: %w", err)
	}

	return nil
}

func (r *Repository) GetUnprocessedQueryItemsCount(ctx context.Context, queryJobID uuid.UUID) (int, error) {
	if r.dbClient == nil {
		return 0, fmt.Errorf("dbClient not initialised")
	}

	var count int

	err := r.dbClient.GetContext(
		ctx,
		&count,
		`
			SELECT COUNT(*)
			FROM query_item
			WHERE query_job_id = $1 AND processed_at IS NULL
		`, queryJobID,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to get count of unprocessed query items (%s): %w", queryJobID, err)
	}

	return count, nil
}

func (r *Repository) MarkQueryJobAsComplete(ctx context.Context, queryJobID uuid.UUID) error {
	if r.dbClient == nil {
		return fmt.Errorf("dbClient not initialised")
	}

	_, err := r.dbClient.ExecContext(
		ctx,
		`
			UPDATE query_job
			SET completed_at = now()
			WHERE id = $1
		`, queryJobID,
	)

	if err != nil {
		return fmt.Errorf("failed to mark query job (%s) as complete: %w", queryJobID.String(), err)
	}

	return nil
}
