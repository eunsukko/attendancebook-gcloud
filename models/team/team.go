package team

import (
	"encoding/json"

	"github.com/eunsukko/attendancebook-gcloud/models/member"
	"github.com/eunsukko/attendancebook-gcloud/util/periodutil"
)

type ID int

type Team struct {
	ID   ID     `json:"id,omitempty"`
	Name string `json:"name"`

	MemberIDs          []member.ID         `json:"member_ids"`
	MemberValidPeriods []periodutil.Period `json:"member_valid_periods"`

	Tags map[member.ID][]string `json:"tags"`
}

func (team Team) String() string {
	data, _ := json.Marshal(team)
	return string(data)
}

func NewTeam() Team {
	team := Team{}

	team.MemberIDs = make([]member.ID, 0)
	team.MemberValidPeriods = make([]periodutil.Period, 0)

	team.Tags = make(map[member.ID][]string)

	return team
}
