package dbrepositorytest

import (
	"context"
	"os"
	"testing"

	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/postgres"
	"github.com/stretchr/testify/require"
)

type TestRepo struct {
	pgClient     *postgres.Client
	dbRepository *dbrepository.Repository
}

func Init(t *testing.T) *TestRepo {
	pgClient, err := postgres.NewClient(os.Getenv("DB_CONN_URL"))
	require.NoError(t, err)

	dbRepository, err := dbrepository.NewRepository(pgClient)
	require.NoError(t, err)

	return &TestRepo{
		pgClient:     pgClient,
		dbRepository: dbRepository,
	}
}

func (r *TestRepo) GetDBRepository() *dbrepository.Repository {
	return r.dbRepository
}

func (r *TestRepo) CleanDB() {
	ctx := context.Background()

	r.pgClient.Connect()
	r.pgClient.ExecContext(ctx, `DELETE FROM link`)
	r.pgClient.ExecContext(ctx, `DELETE FROM query_item`)
	r.pgClient.ExecContext(ctx, `DELETE FROM query_location`)
	r.pgClient.ExecContext(ctx, `DELETE FROM query_job`)
	r.pgClient.Close()
}
