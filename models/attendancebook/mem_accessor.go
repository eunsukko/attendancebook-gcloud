package attendancebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/util/logutil"

	"github.com/eunsukko/attendancebook-gcloud/models/team"
)

var accessor *InMemoryAccessor

func init() {
	fmt.Println("attendancebook mem_accessor.go's init() called")

	accessor = NewInMemoryAccessor()

	// read from backup file
	fmt.Println("load from backup file")
	fillAccessorFromSavedFile()
}

var savedFileName = "attendancebooks.json"

func newPath(fileName string) string {
	_, curFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("fail getting current file path")
	}
	return path.Join(path.Dir(curFilePath), fileName)
}

func fillAccessorFromSavedFile() {
	savedFilePath := newPath("backup/" + savedFileName)

	jsonData, err := ioutil.ReadFile(savedFilePath)
	if err != nil {
		fmt.Errorf("can't find fild: %v", savedFilePath)
	}

	attendancebooks := []Attendancebook{}
	json.Unmarshal(jsonData, &attendancebooks)

	maxID := ID(3000)
	for _, ab := range attendancebooks {
		fmt.Printf("ab.ID: %v\n", ab.ID)
		newID := accessor.Add(ab)
		maxID = func(a, b ID) ID {
			if a < b {
				return b
			}
			return a
		}(maxID, newID)
	}

	// 좋지 않은 습관..
	// 어떻게 예쁘게 고칠까?
	// 최신 backup 파일에서 가져오는 부분 구현할 때 함께 고치자
	accessor.newID = maxID + 1
}

type bookInfo struct {
	t      time.Time
	teamID team.ID
}

type key string

func (bi bookInfo) GenKey() key {
	return key(fmt.Sprintf("%v:%v", bi.t.Unix(), bi.teamID))
}

// 이런 경우 ...json파일로 dump 할 수 있으려나...
//
type InMemoryAccessor struct {
	newID           ID
	attendancebooks map[ID]Attendancebook

	// ids[i], bookInfos[i]: 동일한 i로 접근, 생성될 때 둘다 append로 맨 뒤에 추가됨
	ids       []ID
	bookInfos []bookInfo // for concurrent query

	toID map[key]ID // map[bookInfo.GenKey()]ID

	sync.RWMutex
}

func GetLoadedAccessor() Accessor {
	return accessor
}

func NewInMemoryAccessor() *InMemoryAccessor {

	accessor := InMemoryAccessor{
		newID:           3000,
		attendancebooks: make(map[ID]Attendancebook),
		ids:             make([]ID, 0, 10),
		bookInfos:       make([]bookInfo, 0, 10),
		toID:            make(map[key]ID),
	}
	return &accessor
}

// id 를 안다는 것은 항상 존재한다는 것 (id 를 알아낼 때 없으면 새 id 가 부여되기에)
func (accessor *InMemoryAccessor) Get(id ID) (Attendancebook, error) {
	accessor.RLock()
	defer accessor.RUnlock()

	return accessor.attendancebooks[id], nil
}

func (accessor *InMemoryAccessor) Put(id ID, ab Attendancebook) error {
	accessor.Lock()
	defer accessor.Unlock()

	accessor.attendancebooks[id] = ab

	return nil
}

// func (accessor *InMemoryAccessor) Post(ab Attendancebook) (Attendancebook.ID, error) {

// }

func (accessor *InMemoryAccessor) GetExistIDs() []ID {
	return accessor.ids
}

func (accessor *InMemoryAccessor) Add(ab Attendancebook) ID {
	accessor.Lock()
	defer accessor.Unlock()

	ab.ID = accessor.newID
	accessor.newID++

	bi := bookInfo{ab.Period.From, ab.Team.ID}
	k := bi.GenKey()

	accessor.toID[k] = ab.ID

	accessor.bookInfos = append(accessor.bookInfos, bi)
	accessor.ids = append(accessor.ids, ab.ID)

	accessor.attendancebooks[ab.ID] = ab

	return ab.ID
}

// queryID 에서 만드는게 좀 마음에 안든다...
// 어떻게 설계하는게 좋은 설계일까??...

// t, teamID를 만족하는 id를 리턴
// 존재하지 않을 경우 빈 출석부를 만들고 해당 id 리턴
func (accessor *InMemoryAccessor) QueryID(t time.Time, teamID team.ID) ID {
	bi := bookInfo{t, teamID}
	k := bi.GenKey()

	accessor.RLock()

	if id, ok := accessor.toID[k]; ok {
		accessor.RUnlock()
		return id
	}
	accessor.RUnlock()

	// bi, k 를 ab에서도 얻어 낼 수 있음이 보장되어야 함
	// 그러려면 일단 가장 먼저 떠오르는 것은 t 가 period의 처음이여야 한다는 것??
	// 생각보다 QueryID를 사용하는데 제한이 많구먼 -> 캐시효과를 볼라니...

	// to reduce time for waiting locking
	ab := NewAttendancebook(t, teamID)

	return accessor.Add(ab)
}

// t, teamID 가 없을경우 (zero value) * 같은 역할
func (accessor *InMemoryAccessor) QueryIDs(t time.Time, teamID team.ID) []ID {
	ids := []ID{}

	accessor.RLock()
	curInfoLen := len(accessor.bookInfos)
	accessor.RUnlock()

	logutil.Trace.Printf("QueryIDs, t:%v, teamID: %v\n", t, teamID)
	for i, info := range accessor.bookInfos[:curInfoLen] {
		logutil.Trace.Printf("bi(%v, %v)\n", info.t.Format(time.RFC3339), info.teamID)
		// teamID가 같은 애들이 더 많기 때문에 먼저 확인함
		if (teamID != team.ID(0)) && teamID != info.teamID {
			continue
		}
		if (t != time.Time{}) && !t.Equal(info.t) {
			continue
		}
		ids = append(ids, accessor.ids[i])
	}

	return ids
}
