package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/competitive-analysis/internal/api"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/postgres"

	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("cannot initialise config %v", err)
	}

	pgClient, err := postgres.NewClient(config.RDSConnectionURL)
	if err != nil {
		log.Fatalf("cannot initialise pg client: %v", err)
	}

	dbRepository, err := dbrepository.NewRepository(pgClient)
	if err != nil {
		log.Fatalf("cannot initialise repository: %v", err)
	}

	service := api.NewService(dbRepository, nil, nil)
	lambda.Start(service.GetQueryJobs)
}
