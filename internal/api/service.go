package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/competitive-analysis/api/apischema"
	"github.com/jponc/competitive-analysis/api/eventschema"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/lambdaresponses"
	"github.com/jponc/competitive-analysis/pkg/sns"
	"github.com/jponc/competitive-analysis/pkg/zenserp"
	log "github.com/sirupsen/logrus"
)

type QueryConfig struct {
	Country      string
	Locations    []string
	Num          string
	Device       string
	SearchEngine string
}

var queryConfigDefaults = QueryConfig{
	Country: "US",
	Locations: []string{
		"Mather,California,United States",
		"Melstone,Montana,United States",
		"Austin County,Texas,United States",
		"Denton,North Carolina,United States",
		"Kingfield,Maine,United States",
	},
	Num:          "100",
	Device:       "desktop",
	SearchEngine: "google.com",
}

type Service struct {
	dbrepository  *dbrepository.Repository
	snsClient     *sns.Client
	zenserpClient *zenserp.Client
}

func NewService(dbrepository *dbrepository.Repository, snsClient *sns.Client, zenserpClient *zenserp.Client) *Service {
	s := &Service{
		dbrepository:  dbrepository,
		snsClient:     snsClient,
		zenserpClient: zenserpClient,
	}

	return s
}

func (s *Service) Healthcheck(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lambdaresponses.Respond200(apischema.HealthcheckResponse{Status: "OK"})
}

func (s *Service) CreateQueryJob(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.dbrepository == nil {
		log.Errorf("dbrepository not defined")
		return lambdaresponses.Respond500()
	}

	if s.snsClient == nil {
		log.Errorf("snsClient not defined")
		return lambdaresponses.Respond500()
	}

	req := &apischema.CreateQueryJobRequest{}

	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil || req.Keyword == "" {
		log.Errorf("failed to Unmarshal or error keyword")
		return lambdaresponses.Respond400(fmt.Errorf("bad request"))
	}

	err = s.dbrepository.Connect()
	if err != nil {
		log.Errorf("error connecting to repository db: %v", err)
		return lambdaresponses.Respond500()
	}

	// Create QueryJob
	queryJobID, err := s.dbrepository.CreateQueryJob(ctx, req.Keyword)
	if err != nil {
		log.Errorf("error creating query job: %v", err)
		return lambdaresponses.Respond500()
	}

	// Create QueryLocations
	for _, location := range queryConfigDefaults.Locations {
		_, err := s.dbrepository.CreateQueryLocation(
			ctx,
			queryJobID.String(),
			queryConfigDefaults.Device,
			queryConfigDefaults.SearchEngine,
			queryConfigDefaults.Num,
			queryConfigDefaults.Country,
			location,
		)

		if err != nil {
			log.Errorf("error creating query location: %v", err)
			return lambdaresponses.Respond500()
		}

		log.Infof("keyword: %s, location: %s", queryJobID, location)
	}

	// Publish QueryJobCreated
	msg := eventschema.QueryJobCreatedMessage{
		ID: queryJobID.String(),
	}

	err = s.snsClient.Publish(ctx, eventschema.QueryJobCreated, msg)
	if err != nil {
		log.Errorf("failed to publish SNS: %v", err)
		return lambdaresponses.Respond500()
	}

	return lambdaresponses.Respond200(apischema.CreateQueryJobResponse{QueryJobID: queryJobID.String()})
}

func (s *Service) ZenserpBatchWebhook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if s.dbrepository == nil {
		log.Errorf("dbrepository not defined")
		return lambdaresponses.Respond500()
	}

	if s.snsClient == nil {
		log.Errorf("snsClient not defined")
		return lambdaresponses.Respond500()
	}

	if s.zenserpClient == nil {
		log.Errorf("zenserpClient not defined")
		return lambdaresponses.Respond500()
	}

	err := s.dbrepository.Connect()
	if err != nil {
		log.Errorf("error connecting to repository db: %v", err)
		return lambdaresponses.Respond500()
	}

	unprocessedQueryJobs, err := s.dbrepository.GetUnprocessedQueryJobs(ctx)
	if err != nil {
		log.Fatalf("failed to query unprocessed query jobs: %v", err)
	}

	for _, queryJob := range *unprocessedQueryJobs {
		batch, err := s.zenserpClient.GetBatch(ctx, *queryJob.ZenserpBatchID)
		if err != nil {
			log.Fatal("failed to get batch %s: %v", queryJob.ZenserpBatchID, err)
		}

		// if zenserp batch state is notified meaning it's done, we send an SNS message to process this batch
		// and mark is processed from the database.
		if batch.State == "notified" {
			// Publish QueryJobCreated
			msg := eventschema.ZenserpBatchDoneProcessingMessage{
				QueryJobID:     queryJob.ID.String(),
				ZenserpBatchID: *queryJob.ZenserpBatchID,
			}

			err = s.snsClient.Publish(ctx, eventschema.ZenserpBatchDoneProcessing, msg)
			if err != nil {
				log.Fatalf("failed to publish SNS: %v", err)
			}

			err = s.dbrepository.ProcessQueryJob(ctx, queryJob.ID)
			if err != nil {
				log.Fatalf("failed to mark query job as processed: %v", err)
			}
		}
	}

	return lambdaresponses.Respond200(apischema.HealthcheckResponse{Status: "OK"})
}
