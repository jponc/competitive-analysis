//go:build integration
package api_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
	"github.com/jponc/competitive-analysis/api/apischema"
	"github.com/jponc/competitive-analysis/internal/api"
	"github.com/jponc/competitive-analysis/internal/dbrepositorytest"
	"github.com/jponc/competitive-analysis/internal/repository/dbrepository"
	"github.com/stretchr/testify/require"
)

func getResponse(resp events.APIGatewayProxyResponse, _ error) events.APIGatewayProxyResponse {
	return resp
}

type mockSnsClient struct{}

func (m *mockSnsClient) Publish(ctx context.Context, topic string, message interface{}) error {
	return nil
}

func Test_CreateQueryJob(t *testing.T) {
	testRepo := dbrepositorytest.Init(t)
	dbRepository := testRepo.GetDBRepository()

	tests := []struct {
		name                       string
		request                    events.APIGatewayProxyRequest
		expectedResponseStatusCode int
		dbrepository               *dbrepository.Repository
		snsClient                  api.SNSClient
	}{
		{
			name:                       "returns 500 when dbrepository is nil",
			dbrepository:               nil,
			snsClient:                  &mockSnsClient{},
			expectedResponseStatusCode: 500,
		},
		{
			name:                       "returns 500 when snsClient is nil",
			dbrepository:               dbRepository,
			snsClient:                  nil,
			expectedResponseStatusCode: 500,
		},
		{
			name:         "returns 400 when request body is bogus",
			dbrepository: testRepo.GetDBRepository(),
			snsClient:    &mockSnsClient{},
			request: events.APIGatewayProxyRequest{
				Body: `bogus`,
			},
			expectedResponseStatusCode: 400,
		},
		{
			name:         "returns 400 when keyword is empty",
			dbrepository: dbRepository,
			snsClient:    &mockSnsClient{},
			request: events.APIGatewayProxyRequest{
				Body: `{"keyword": ""}`,
			},
			expectedResponseStatusCode: 400,
		},
		{
			name:         "returns 200 and creates a query job",
			dbrepository: dbRepository,
			snsClient:    &mockSnsClient{},
			request: events.APIGatewayProxyRequest{
				Body: `{"keyword": "hello world"}`,
			},
			expectedResponseStatusCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRepo.CleanDB()

			ctx := context.Background()
			service := api.NewService(tt.dbrepository, tt.snsClient, nil)
			resp, _ := service.CreateQueryJob(ctx, tt.request)
			require.Equal(t, tt.expectedResponseStatusCode, resp.StatusCode)

			if resp.StatusCode == 200 {
				dbRepository.Connect()
				responseBody := &apischema.CreateQueryJobResponse{}
				err := json.Unmarshal([]byte(resp.Body), responseBody)
				require.NoError(t, err)

				queryJobID := uuid.FromStringOrNil(responseBody.QueryJobID)
				queryJob, err := dbRepository.GetQueryJob(ctx, queryJobID)
				require.NoError(t, err)
				require.NotNil(t, queryJob)
				dbRepository.Close()
			}
		})
	}

}
