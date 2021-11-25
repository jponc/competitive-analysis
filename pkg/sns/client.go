package sns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	awsSns "github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-xray-sdk-go/xray"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	awsSnsClient *awsSns.SNS
	snsPrefix    string
}

// NewClient instantiates a SNS client
func NewClient(awsRegion, snsPrefix string) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create aws session: %v", err)
	}

	awsSnsClient := sns.New(sess)
	xray.AWS(awsSnsClient.Client)

	c := &Client{
		awsSnsClient: awsSnsClient,
		snsPrefix:    snsPrefix,
	}

	return c, nil
}

func (c *Client) Publish(ctx context.Context, topic string, message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal sns message: %v", err)
	}

	input := &awsSns.PublishInput{
		Message:  aws.String(string(msg)),
		TopicArn: aws.String(fmt.Sprintf("%s-%s", c.snsPrefix, topic)),
	}

	result, err := c.awsSnsClient.PublishWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to publish to sns: %v", err)
	}

	log.Infof("successfully processed sns message: (%s) for topic (%s)", result, topic)
	return nil
}
