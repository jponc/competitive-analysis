package eventschema

const (
	QueryJobCreated            string = "QueryJobCreated"
	ParseQueryJobURL           string = "ParseQueryJobURL"
	ZenserpBatchDoneProcessing string = "ZenserpBatchDoneProcessing"
)

type QueryJobCreatedMessage struct {
	ID      string `json:"id"`
	Keyword string `json:"keyword"`
}

type ZenserpBatchDoneProcessingMessage struct {
	QueryJobID     string `json:"query_job_id"`
	ZenserpBatchID string `json:"zenserp_batch_id"`
}

type ParseQueryJobURLMessage struct {
	QueryJobID string `json:"query_job_id"`
	URL        string `json:"url"`
}
