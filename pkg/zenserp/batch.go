package zenserp

import (
	"context"
	"fmt"
)

func (c *Client) Batch(ctx context.Context, name string, jobs []Job) (*BatchResult, error) {
	batchRequest := &BatchRequest{
		WebhookURL: c.batchWebhookURL,
		Name:       name,
		Jobs:       jobs,
	}

	var batchResult BatchResult

	err := c.postJSON(ctx, batchPath, batchRequest, &batchResult)
	if err != nil {
		return nil, fmt.Errorf("failed to post batch: %w", err)
	}

	return &batchResult, nil
}

func (c *Client) GetBatch(ctx context.Context, batchID string) (*Batch, error) {

	var b Batch

	endpoint := fmt.Sprintf(getBatchPath, batchID)
	err := c.getJSON(ctx, endpoint, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch: %w", err)
	}

	return &b, nil
}
