package webscraper

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	httpClient *http.Client
}

type Link struct {
	Text    string
	LinkURL string
}

type ScrapeResult struct {
	Title       string
	Description string
	Body        string
	Links       []Link
}

func NewClient(httpClient *http.Client) *Client {
	c := &Client{
		httpClient,
	}

	return c
}

func (c *Client) Scrape(ctx context.Context, link string) (*ScrapeResult, error) {
	res, err := c.httpClient.Get(link)
	if err != nil {
		return nil, fmt.Errorf("failed to get link (%s): %v", link, err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("content type is not text/html: %s", contentType)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Remove
	doc.Find("script").Remove()
	doc.Find("br").Remove()
	doc.Find("img").Remove()
	doc.Find("iframe").Remove()
	doc.Find("style").Remove()

	// Title
	title := doc.Find("title").Text()

	// Links
	links := []Link{}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		linkURL, _ := item.Attr("href")
		text := strings.TrimSpace(item.Text())
		links = append(links, Link{
			Text:    text,
			LinkURL: linkURL,
		})
	})

	// description
	description := ""
	doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("name", "") == "description" {
			description = item.AttrOr("content", "")
		}
	})

	// BODY
	bodyContentsMap := map[string]bool{}
	doc.Find("body *").Each(func(_ int, item *goquery.Selection) {
		b := strings.TrimSpace(item.Text())
		for _, line := range strings.Split(b, "\n") {
			lineCompact := strings.TrimSpace(line)

			if lineCompact != "" {
				bodyContentsMap[lineCompact] = true
			}
		}
	})

	// get body
	bodyContents := []string{}
	for k := range bodyContentsMap {
		bodyContents = append(bodyContents, k)
	}

	scrapeResult := &ScrapeResult{
		Title:       title,
		Description: description,
		Body:        strings.Join(bodyContents, " "),
		Links:       links,
	}

	return scrapeResult, nil
}
