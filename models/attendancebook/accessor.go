package attendancebook

import (
	"errors"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/models/team"
)

var ErrAttendancebookNotExist = errors.New("attendance does not exist. check attendancebook.ID")

type Accessor interface {
	Get(id ID) (Attendancebook, error)

	Put(id ID, ab Attendancebook) error
	//
	// Post(ab Attendancebook) (Attendancebook.ID, error)

	GetExistIDs() []ID
	// 일단은 없는 경우 여기서 새로운 attendancebook 객체를 만듬 (id 도 부여)
	QueryID(t time.Time, teamID team.ID) ID // t, teamID 조합으로는 자주 id를 알아내려고 할 것이기에
	QueryIDs(t time.Time, teamID team.ID) []ID
}

// 어떤 것을 질의하게 될까??
// 1. t, teamID 가 확정된 경우 (각 셀에서 출석부 사용)
// 2. t만 알고 있는 경우 (총무단에서 결과 수집)
