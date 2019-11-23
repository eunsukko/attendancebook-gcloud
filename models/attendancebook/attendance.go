package attendancebook

import (
	"encoding/json"
	"fmt"

	"github.com/eunsukko/attendancebook-gcloud/models/member"
)

// Attendance represents a member's attendance result
type Attendance struct {
	MemberID member.ID `json:"member_id"`
	Attended []bool    `json:"attended"` // len(Attended) == len(eventlist)
}

func (attendance Attendance) String() string {
	jsonData, _ := json.Marshal(attendance)
	return fmt.Sprintf("%s", string(jsonData))
}
