package crawler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
	"github.com/jponc/competitive-analysis/api/eventschema"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/sns"
	"github.com/jponc/competitive-analysis/pkg/textrazor"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	textRazorClient *textrazor.Client
	repository      *dbrepository.Repository
	snsClient       *sns.Client
}

func NewService(textRazorClient *textrazor.Client, repository *dbrepository.Repository, snsClient *sns.Client) *Service {
	s := &Service{
		textRazorClient: textRazorClient,
		repository:      repository,
		snsClient:       snsClient,
	}

	return s
}

func (s *Service) TextRazorParseQueryJobURL(ctx context.Context, snsEvent events.SNSEvent) {
	if s.repository == nil {
		log.Fatalf("repository not defined")
	}

	if s.textRazorClient == nil {
		log.Fatalf("textRazorClient not defined")
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

	queryJobID, err := uuid.FromString(msg.QueryJobID)
	if err != nil {
		log.Fatalf("unable to convert query job id string to UUID: %v", err)
	}

	res, err := s.textRazorClient.Analyze(ctx, url, []textrazor.Extractor{})
	if err == nil {
		err = s.repository.SetQueryItemsProcessedWithBody(ctx, queryJobID, url, res.CleanedText)
	} else {
		// don't panic if there's a URL that can't be processed by textrazor, just continue
		log.Errorf("unable to request cleaned HTML with URL (%s) from textrazor: %v", url, err)
		err = s.repository.SetQueryItemsErrorProcessing(ctx, queryJobID, url)
	}

	if err != nil {
		log.Fatalf("failed to process  query items for query job (%s) with url (%s)", queryJobID.String(), url)
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
}
