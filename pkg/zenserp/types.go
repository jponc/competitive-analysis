package zenserp

type QueryInfo struct {
	Query        string `json:"q"`
	SearchEngine string `json:"search_engine"`
	Device       string `json:"device"`
	URL          string `json:"url"`
}

type ResultItem struct {
	Position    int    `json:"position"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type ResultItemWithoutDescription struct {
	Position int    `json:"position"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

type QueryResult struct {
	Query      QueryInfo    `json:"query"`
	ResulItems []ResultItem `json:"organic"`
}

type QueryResultWithoutDescription struct {
	Query      QueryInfo                      `json:"query"`
	ResulItems []ResultItemWithoutDescription `json:"organic"`
}

type BatchRequest struct {
	WebhookURL string `json:"webhook_url"`
	Name       string `json:"name"`
	Jobs       []Job  `json:"jobs"`
}

type BatchResult struct {
	BatchID string `json:"id"`
}

type Job struct {
	Query        string `json:"q"`
	Num          string `json:"num"`
	SearchEngine string `json:"search_engine"`
	Device       string `json:"device"`
	Country      string `json:"gl"`
	Location     string `json:"location"`
}

type Batch struct {
	ID      string                          `json:"id"`
	Name    string                          `json:"name"`
	State   string                          `json:"state"`
	Results []QueryResultWithoutDescription `json:"jobs"`
}
