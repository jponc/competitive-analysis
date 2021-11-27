package zenserp

import (
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	apiKey          string
	baseURL         *url.URL
	httpClient      *http.Client
	batchWebhookURL string
}

// NewClient instantiates a zenserp client
func NewClient(apiKey string, httpClient *http.Client, batchWebhookURL string) (*Client, error) {
	baseURL, err := url.Parse(zenserpBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing zenser Base URL (%w)", err)
	}

	c := &Client{
		apiKey:          apiKey,
		baseURL:         baseURL,
		httpClient:      httpClient,
		batchWebhookURL: batchWebhookURL,
	}

	return c, nil
}
