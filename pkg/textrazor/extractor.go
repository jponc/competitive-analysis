package textrazor

import "strings"

const (
	Entities Extractor = "entities"
	Topics   Extractor = "topics"
)

type Extractor string
type extractorArr []Extractor

func (a extractorArr) ToString() string {
	s := []string{}

	for _, e := range a {
		s = append(s, string(e))
	}

	return strings.Join(s, ",")
}
