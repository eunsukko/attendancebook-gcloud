package attendancebook

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/util/dateutil"
	"github.com/eunsukko/attendancebook-gcloud/util/logutil"

	"github.com/eunsukko/attendancebook-gcloud/models/eventlist"
	"github.com/eunsukko/attendancebook-gcloud/models/member"
	"github.com/eunsukko/attendancebook-gcloud/models/team"
	"github.com/eunsukko/attendancebook-gcloud/util/periodutil"
)

type ID int

// Attendancebook gathers info related to attendancebook at a specific period
type Attendancebook struct {
	ID          ID                  `json:"id,omitempty"`
	Team        team.Team           `json:"team"`
	Eventlist   eventlist.Eventlist `json:"eventlist"`
	Period      periodutil.Period   `json:"period"` // Period.To 를 기준으로 accessor 에 저장됨
	Members     []member.Member     `json:"members"`
	Info        map[string]string   `json:"info"` // for adding something (ex. comment... etc)
	Attendances []Attendance        `json:"attendances"`
}

func (ab Attendancebook) String() string {
	jsonData, _ := json.Marshal(ab)
	return fmt.Sprintf("%s", string(jsonData))
}

func NewAttendancebook(t time.Time, teamID team.ID) Attendancebook {

	logutil.Info.Printf("New attendancebook\n")
	logutil.Trace.Printf("t: %v, teamID: %v\n", t, teamID)

	ab := Attendancebook{}
	ab.Team = func(teamID team.ID) team.Team {
		accessor := team.GetLoadedInMemoryAccessor()

		// 지금은 항상 올바른 teamID가 들어온다는 가정이 있음 (클라이언트에서 올바른 teamID만 보이게 할 것이기에)
		curTeam, err := accessor.Get(teamID)

		if err != nil {
			return team.Team{}
		}
		return curTeam
	}(teamID)
	ab.Eventlist = func(t time.Time) eventlist.Eventlist {
		provider := eventlist.GetLoadedInMemoryMapedProvider()
		return provider.GetAt(t)
	}(t)
	// t 가 해당하는 week [해당 월요일, 다음 월요일)
	ab.Period = func(t time.Time) periodutil.Period {
		kst := dateutil.NewKST(t)
		thisMonday := kst.CalcThisMonday()
		return periodutil.Period{
			From: thisMonday.Time,
			To:   thisMonday.Add(7 * 24 * time.Hour),
		}
	}(t)
	ab.Members = func(team team.Team) []member.Member {
		accessor := member.GetLoadedInMemoryAccessor()
		members := []member.Member{}
		for i, id := range team.MemberIDs {
			p := team.MemberValidPeriods[i]
			if p.Include(t) {
				curMember, _ := accessor.Get(id)
				members = append(members, curMember)
			}
		}
		return members
	}(ab.Team)
	ab.Info = make(map[string]string)
	ab.Attendances = func(members []member.Member, eventlist eventlist.Eventlist) []Attendance {
		attendances := make([]Attendance, len(members))
		for i, curMember := range members {
			attendances[i] = Attendance{
				curMember.ID,
				make([]bool, len(eventlist)),
			}
		}
		return attendances
	}(ab.Members, ab.Eventlist)

	return ab
}
