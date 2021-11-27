package dbrepository

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
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
			VALUES ($1)
			RETURNING id
		`,
		queryJobID, device, searchEngine, num, country, location)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert query job: %v", err)
	}

	return id, nil

}
