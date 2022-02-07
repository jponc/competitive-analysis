package crawler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
	"github.com/jponc/competitive-analysis/api/eventschema"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/sns"
	"github.com/jponc/competitive-analysis/pkg/webscraper"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	webscraperClient *webscraper.Client
	repository       *dbrepository.Repository
	snsClient        *sns.Client
}

func NewService(webscraperClient *webscraper.Client, repository *dbrepository.Repository, snsClient *sns.Client) *Service {
	s := &Service{
		webscraperClient: webscraperClient,
		repository:       repository,
		snsClient:        snsClient,
	}

	return s
}

func (s *Service) WebScraperParseQueryJobURL(ctx context.Context, snsEvent events.SNSEvent) {
	if s.repository == nil {
		log.Fatalf("repository not defined")
	}

	if s.webscraperClient == nil {
		log.Fatalf("webscraperClient not defined")
	}

	if s.snsClient == nil {
		log.Fatalf("snsClient not defined")
	}

	if err := s.repository.Connect(); err != nil {
		log.Fatalf("can't connect to DB")
	}

	// Unmarshal query item
	snsMsg := snsEvent.Records[0].SNS.Message

	var msg eventschema.ParseQueryJobURLMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	url := msg.URL

	log.Infof("Processing URL: %s", url)

	queryJobID, err := uuid.FromString(msg.QueryJobID)
	if err != nil {
		log.Fatalf("unable to convert query job id string to UUID: %v", err)
	}

	queryItems, err := s.repository.GetQueryItemsFromUrl(ctx, queryJobID, url)
	if err != nil {
		log.Fatalf("unable to get query item id's: %v", err)
	}

	var queryItemIDs []uuid.UUID
	for _, queryItem := range *queryItems {
		queryItemIDs = append(queryItemIDs, queryItem.ID)
	}

	// Run scraping
	res, err := s.webscraperClient.Scrape(ctx, url)

	if err == nil {
		// Create links
		for _, queryItemID := range queryItemIDs {
			for _, link := range res.Links {
				err = s.repository.CreateQueryLink(ctx, queryItemID, link.Text, link.LinkURL)
				if err != nil {
					log.Infof("unable to create link: %v", err)
				}
			}
		}

		// Store Body
		err = s.repository.SetQueryItemsProcessedWithBodyAndTitle(ctx, queryJobID, queryItemIDs, res.Body, res.Title)
	} else {
		// don't panic if there's a URL that can't be processed , just continue
		log.Errorf("unable to request cleaned HTML with URL (%s) from webscraper: %v", url, err)
		err = s.repository.SetQueryItemsErrorProcessing(ctx, queryJobID, url)
	}

	if err != nil {
		log.Fatalf("failed to process query items for query job (%s) with url (%s): %v", queryJobID.String(), url, err)
	}

	// Publish DoneProcessingQueryJobURL message
	doneMsg := eventschema.DoneProcessingQueryJobURLMessage{
		QueryJobID: queryJobID.String(),
		URL:        url,
	}

	err = s.snsClient.Publish(ctx, eventschema.DoneProcessingQueryJobURL, doneMsg)
	if err != nil {
		log.Fatalf("failed to publish SNS: %v", err)
	}

	// Send request to textrazor to extract the content
	log.Infof("Done processing (%s), url: (%s)", queryJobID.String(), url)

	if err := s.repository.Close(); err != nil {
		log.Fatalf("can't close DB connection")
	}
}

func (s *Service) CheckCompletedQueryJobs(ctx context.Context, snsEvent events.SNSEvent) {
	if s.repository == nil {
		log.Fatalf("repository not defined")
	}

	if err := s.repository.Connect(); err != nil {
		log.Fatalf("can't connect to DB")
	}

	// Unmarshal msg
	snsMsg := snsEvent.Records[0].SNS.Message

	var msg eventschema.DoneProcessingQueryJobURLMessage
	err := json.Unmarshal([]byte(snsMsg), &msg)
	if err != nil {
		log.Fatalf("unable to unarmarshal message: %v", err)
	}

	queryJobID, err := uuid.FromString(msg.QueryJobID)
	if err != nil {
		log.Fatalf("unable to convert query job id string to UUID: %v", err)
	}

	// Get unprocessed query items count
	unprocessedCount, err := s.repository.GetUnprocessedQueryItemsCount(ctx, queryJobID)
	if err != nil {
		log.Fatalf("unable to get unprocessed query items count for query job (%s): %v", queryJobID.String(), err)
	}

	if unprocessedCount > 0 {
		log.Infof("%s query job still has %d remaining unprocessed query items", queryJobID.String(), unprocessedCount)
		return
	}

	// mark as complete if there are 0 unprocessed query items
	err = s.repository.MarkQueryJobAsComplete(ctx, queryJobID)
	if err != nil {
		log.Fatalf("%s query job cannot be marked as complete: %v", queryJobID.String(), err)
	}

	log.Infof("Marked query job %s as complete", queryJobID.String())

	if err := s.repository.Close(); err != nil {
		log.Fatalf("can't close DB connection")
	}
}
