package eventschema

const (
	QueryJobCreated string = "QueryJobCreated"
)

type QueryJobCreatedMessage struct {
	ID      string `json:"id"`
	Keyword string `json:"keyword"`
}
