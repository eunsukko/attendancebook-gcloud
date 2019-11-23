package team

import (
	"fmt"

	"github.com/eunsukko/attendancebook-gcloud/models/member"
	"github.com/eunsukko/attendancebook-gcloud/util/dateutil"
)

func Example_Team() {
	team := Team{}
	team.MemberIDs = make([]member.ID, 0)
	team.MemberValidPeriods = make([]Period, 1)
	team.MemberValidPeriods[0] = Period{
		From: dateutil.NewDate(2018, 11, 4),
		To:   dateutil.NewDate(2019, 11, 4),
	}
	fmt.Printf("%v\n", team)

	// Output:
	// .
}
