package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jponc/competitive-analysis/internal/api"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/postgres"
	"github.com/jponc/competitive-analysis/pkg/sns"
	"github.com/jponc/competitive-analysis/pkg/zenserp"

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

	snsClient, err := sns.NewClient(config.AWSRegion, config.SNSPrefix)
	if err != nil {
		log.Fatalf("cannot initialise sns client %v", err)
	}

	httpClient := &http.Client{
		Timeout: time.Duration(1 * time.Minute),
	}

	zenserpClient, err := zenserp.NewClient(config.ZenserpApiKey, httpClient, config.ZenserpBatchWebhookURL)
	if err != nil {
		log.Fatalf("cannot initialise zenserp client %v", err)
	}

	service := api.NewService(dbRepository, snsClient, zenserpClient)
	lambda.Start(service.ZenserpBatchWebhook)
}
