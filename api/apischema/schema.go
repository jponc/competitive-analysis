package apischema

type HealthcheckResponse struct {
	Status string `json:"status"`
}

type CreateQueryJobRequest struct {
	Keyword string `json:"keyword"`
}

type CreateQueryJobResponse struct {
	QueryJobID string `json:"query_job_id"`
}
