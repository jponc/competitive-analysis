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

	if keyword == "" {
		return uuid.Nil, fmt.Errorf("keyword is blank")
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
