package periodutil

import (
	"encoding/json"
	"time"
)

// Period [from, to)
type Period struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func (p Period) String() string {
	data, _ := json.Marshal(p)
	return string(data)
}

func (p Period) Include(t time.Time) bool {
	return t.After(p.From) && t.Before(p.To)
}
