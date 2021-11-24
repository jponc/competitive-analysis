package dbrepository

import (
	"fmt"

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
