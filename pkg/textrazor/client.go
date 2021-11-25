package textrazor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const endpoint = "http://api.textrazor.com"

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string, httpClient *http.Client) *Client {
	c := &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
	}

	return c
}

func (c *Client) Analyze(ctx context.Context, siteUrl string, extractors []Extractor) (*Response, error) {
	if c.httpClient == nil {
		return nil, fmt.Errorf("httpClient is not defined")
	}

	if c.apiKey == "" {
		return nil, fmt.Errorf("apiKey is not defined")
	}

	eArr := extractorArr(extractors)

	data := url.Values{}
	data.Set("url", siteUrl)
	data.Set("cleanup.returnCleaned", "true")
	data.Set("extractors", eArr.ToString())

	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("can't initialise http request: %v", err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("x-textrazor-key", c.apiKey)

	res, err := c.httpClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("http request returned non 200: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var jsonResponse struct {
		Response Response `json:"response"`
		Ok       bool     `json:"ok"`
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	if jsonResponse.Ok == false {
		return nil, fmt.Errorf("response is not ok")
	}

	return &jsonResponse.Response, nil
}
