package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jponc/competitive-analysis/internal/api"
)

func main() {
	service := api.NewService()
	lambda.Start(service.Healthcheck)
}
