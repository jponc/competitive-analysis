package resultrankings

import (
	"context"
	"encoding/json"

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

	// res, err := s.zenserpClient.SearchWithCountry(
	//   ctx,
	//   msg.Keyword,
	//   "google.com",
	//   "desktop",
	//   "us",
	//   100,
	// )

	// for _, item := range res.ResulItems {
	//   item.
	// }
}
