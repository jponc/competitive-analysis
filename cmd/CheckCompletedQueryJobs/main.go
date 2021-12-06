package main

import (
	"log"

	"github.com/jponc/competitive-analysis/internal/crawler"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/postgres"

	"github.com/aws/aws-lambda-go/lambda"
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

	service := crawler.NewService(nil, dbRepository, nil)
	lambda.Start(service.CheckCompletedQueryJobs)
}
