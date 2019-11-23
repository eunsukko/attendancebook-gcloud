package dateutil

import (
	"fmt"
	"time"
)

func init() {
	koreaLoc = time.FixedZone("KST", 9*60*60) // UTC+9
}

// NewDate 는 time.Date 에서 주로 쓰는 year, month, day 만 추려서 time.UTC 시간으로 리턴
func NewDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// NewKSTDate 는 입력받은 시간의 KST를 리턴
func NewKSTDate(year int, month time.Month, day, hour, min, sec, nsec int) KST {
	return KST{time.Date(year, month, day, hour, min, sec, nsec, koreaLoc)}
}

// KST 는 어떤 가정이 있는가?
// 일단 NewKST 를 통해서 만들기 때문에 항상 타임존이 +9임
type KST struct {
	time.Time
}

var koreanDays = [...]string{
	"일", "월", "화", "수", "목", "금", "토",
}
var koreaLoc *time.Location

func (t KST) String() string {
	dayIdx := int(t.Weekday())
	return fmt.Sprintf("%s(%s)", t.Format("2006.01.02"), koreanDays[dayIdx])
}

// CalcThisMonday 는 해당 t 기준으로 지나간 가장 가까운 월요일을 리턴
// time.Truncate 가 Zero time (즉 기준이 utc 기준으로 되어 있어서 utc일때의 시간으로 내가 원하는 월요일을 구하고 다시 kst에 해당 시간으로 리턴을 해줌)
// kst -> utc (여기서 월요일 구하기) -> kst
func (t KST) CalcThisMonday() KST {
	utcTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	tMonday := utcTime.Truncate(7 * 24 * time.Hour)

	return NewKSTDate(tMonday.Year(), tMonday.Month(), tMonday.Day(), tMonday.Hour(), tMonday.Minute(), tMonday.Second(), tMonday.Nanosecond())
}

// NewKST 는 입력받은 t 를 해당 시간의 kst시간의 형태로 만들어서 리턴
// t.UTC() 는 같지만 다른 타임존에 있는 경우를 KST로 통일하기 위함
func NewKST(t time.Time) KST {
	tUTC := t.UTC()
	newT := time.Date(tUTC.Year(), tUTC.Month(), tUTC.Day(), tUTC.Hour(), tUTC.Minute(), tUTC.Second(), tUTC.Nanosecond(), koreaLoc)

	_, offset := newT.Zone()
	newT = newT.Add(time.Duration(offset) * time.Second)

	return KST{newT}
}
