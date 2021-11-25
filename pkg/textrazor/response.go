package textrazor

type EntityArray []Entity
type Entity struct {
	EntityID        string  `json:"entityId"`
	ConfidenceScore float32 `json:"confidenceScore"`
	RelevanceScore  float32 `json:"relevanceScore"`
	MatchedText     string  `json:"matchedText"`
}

type TopicArray []Topic
type Topic struct {
	Label string  `json:"label"`
	Score float32 `json:"score"`
}

type Response struct {
	CleanedText        string      `json:"cleanedText"`
	Language           string      `json:"language"`
	LanguageIsReliable bool        `json:"languageIsReliable"`
	Entities           EntityArray `json:"entities"`
	Topics             TopicArray  `json:"topics"`
}
