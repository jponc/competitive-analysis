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
	log "github.com/sirupsen/logrus"
)

type Service struct {
	dbrepository *dbrepository.Repository
	snsClient    *sns.Client
}

func NewService(dbrepository *dbrepository.Repository, snsClient *sns.Client) *Service {
	s := &Service{
		dbrepository: dbrepository,
		snsClient:    snsClient,
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

	id, err := s.dbrepository.CreateQueryJob(ctx, req.Keyword)
	if err != nil {
		log.Errorf("error creating query job: %v", err)
		return lambdaresponses.Respond500()
	}

	msg := eventschema.QueryJobCreatedMessage{
		ID:      id.String(),
		Keyword: req.Keyword,
	}

	err = s.snsClient.Publish(ctx, eventschema.QueryJobCreated, msg)
	if err != nil {
		log.Errorf("failed to publish SNS: %v", err)
		return lambdaresponses.Respond500()
	}

	return lambdaresponses.Respond200(apischema.CreateQueryJobResponse{QueryJobID: id.String()})
}
