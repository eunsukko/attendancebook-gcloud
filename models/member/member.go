package member

import (
	"encoding/json"
)

type ID int

type Member struct {
	ID   ID                `json:"id,omitempty"`
	Info map[string]string `json:"info"`
}

func (m Member) String() string {
	data, _ := json.Marshal(m)
	return string(data)
}
