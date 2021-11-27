package main

import (
	"github.com/jponc/competitive-analysis/internal/api"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	service := api.NewService(nil, nil)
	lambda.Start(service.ZenserpBatchWebhook)
}
