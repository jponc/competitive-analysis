package crawler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
	"github.com/jponc/competitive-analysis/api/eventschema"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/jponc/competitive-analysis/pkg/textrazor"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	textRazorClient *textrazor.Client
	repository      *dbrepository.Repository
}

func NewService(textRazorClient *textrazor.Client, repository *dbrepository.Repository) *Service {
	s := &Service{
		textRazorClient: textRazorClient,
		repository:      repository,
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

	// Send request to textrazor to extract the content
	log.Infof("Processing (%s), url: (%s)", queryJobID.String(), url)
}
