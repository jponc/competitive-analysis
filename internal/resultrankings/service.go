package resultrankings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jponc/competitive-analysis/api/eventschema"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/sns"
	"github.com/jponc/competitive-analysis/pkg/zenserp"
)

type Service struct {
	zenserpClient *zenserp.Client
	repository    *dbrepository.Repository
	snsClient     *sns.Client
}

func NewService(zenserpClient *zenserp.Client, repository *dbrepository.Repository, snsClient *sns.Client) *Service {
	s := &Service{
		zenserpClient: zenserpClient,
		repository:    repository,
		snsClient:     snsClient,
	}

	return s
}

func (s *Service) QueryJobZenserp(ctx context.Context, snsEvent events.SNSEvent) {
	if s.zenserpClient == nil {
		log.Fatalf("zenserpClient not defined")
	}

	if s.repository == nil {
		log.Fatalf("repository not defined")
	}

	if s.snsClient == nil {
		log.Fatalf("snsClient not defined")
	}

	if err := s.repository.Connect(); err != nil {
		log.Fatalf("can't connect to DB")
	}

	snsMsg := snsEvent.Records[0].SNS.Message

	var msg eventschema.QueryJobCreatedMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	queryJobID, err := uuid.FromString(msg.ID)
	if err != nil {
		log.Fatalf("unable to convert query job string to UUID: %v", err)
	}

	// Fetch query job
	queryJob, err := s.repository.GetQueryJob(ctx, queryJobID)
	if err != nil {
		log.Fatalf("failed to get query job: %s, %v", queryJobID, err)
	}

	// Fetch query locations
	queryLocations, err := s.repository.GetQueryLocations(ctx, queryJobID)
	if err != nil {
		log.Fatalf("failed to get query locations of query job: %s", queryJobID)
	}

	// Convert query locations to a zenserp jobs
	var zenserpJobs []zenserp.Job
	for _, queryLocation := range *queryLocations {
		zenserpJobs = append(zenserpJobs, zenserp.Job{
			Query:        queryJob.Keyword,
			Num:          queryLocation.Num,
			SearchEngine: queryLocation.SearchEngine,
			Device:       queryLocation.Device,
			Country:      queryLocation.Country,
			Location:     queryLocation.Location,
		})
	}

	// Create zenserp batch
	batchResult, err := s.zenserpClient.Batch(ctx, fmt.Sprintf("%s: %s", queryJob.ID, queryJob.Keyword), zenserpJobs)
	if err != nil {
		log.Fatalf("failed to create zenserp batch: %s", queryJobID)
	}

	// Set batch ID to query job
	err = s.repository.SetZenserpBatchToQueryJob(ctx, queryJobID, batchResult.BatchID)
	if err != nil {
		log.Fatalf("failed to set zenserp batch ID to query job: %v", err)
	}
}

func (s *Service) ZenserpBatchExtractResults(ctx context.Context, snsEvent events.SNSEvent) {
	if s.zenserpClient == nil {
		log.Fatalf("zenserpClient not defined")
	}

	if s.repository == nil {
		log.Fatalf("repository not defined")
	}

	if s.snsClient == nil {
		log.Fatalf("snsClient not defined")
	}

	if err := s.repository.Connect(); err != nil {
		log.Fatalf("can't connect to DB")
	}

	// Unmarshal message
	snsMsg := snsEvent.Records[0].SNS.Message

	var msg eventschema.ZenserpBatchDoneProcessingMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	queryJobID, err := uuid.FromString(msg.QueryJobID)
	if err != nil {
		log.Fatalf("unable to convert query job string to UUID: %v", err)
	}

	// Get QueryLocations so we can pull the ID later based on location
	queryLocations, err := s.repository.GetQueryLocations(ctx, queryJobID)
	if err != nil {
		log.Fatalf("unable to get query locations of %s: %v", queryJobID, err)
	}

	// Get ZenserpBatch
	zenserpBatchID := msg.ZenserpBatchID
	batch, err := s.zenserpClient.GetBatch(ctx, zenserpBatchID)

	urls := map[string]bool{}

	// Create QueryItems based on zenserp batch result and query location
	for _, result := range batch.Results {
		for _, queryLocation := range *queryLocations {
			// Found query location, create query items
			if result.Query.Location == queryLocation.Location {
				for _, resultItem := range result.ResulItems {
					if resultItem.URL == "" {
						// We don't want to process empty URL
						continue
					}

					_, err := s.repository.CreateQueryItem(ctx, queryJobID, queryLocation.ID, resultItem.Position, resultItem.URL, resultItem.Title)
					if err != nil {
						log.Fatalf("unable to create query item for query job (%s) and query location (%s): %v", queryJobID.String(), queryLocation.ID.String(), err)
					}

					// Create a set of urls to be processed later
					urls[resultItem.URL] = true
				}
			}
		}
	}

	// Iterate all urls generated earlier, then publish a message to extract results
	for url := range urls {
		msg := eventschema.ParseQueryJobURLMessage{
			QueryJobID: queryJobID.String(),
			URL:        url,
		}

		err = s.snsClient.Publish(ctx, eventschema.ParseQueryJobURL, msg)
		if err != nil {
			log.Fatalf("failed to publish SNS: %v", err)
		}
	}

	log.Infof("done creating query items: query job id: %s, batchID: %s", queryJobID, zenserpBatchID)
}
