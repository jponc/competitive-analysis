package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jponc/competitive-analysis/internal/crawler"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/postgres"
	"github.com/jponc/competitive-analysis/pkg/textrazor"

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

	httpClient := &http.Client{
		Timeout: time.Duration(1 * time.Minute),
	}

	textrazorClient := textrazor.NewClient(config.TextRazorAPIKey, httpClient)

	service := crawler.NewService(textrazorClient, dbRepository)
	lambda.Start(service.TextRazorParseQueryJobURL)
}
