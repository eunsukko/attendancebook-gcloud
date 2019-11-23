package arrangement

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/models/team"
	"github.com/eunsukko/attendancebook-gcloud/models/utils/sortutil"
)

type InMemoryMapedProvider struct {
	sync.RWMutex
	snapshots map[int64]Arrangement
}

func (provider *InMemoryMapedProvider) GetAt(t time.Time) Arrangement {
	at := t.Unix()
	// get
	provider.RLock()
	sortedSavedTimes := func() []int64 {
		times := make([]int64, 0, len(provider.snapshots))
		for k := range provider.snapshots {
			times = append(times, k)
		}
		sort.Sort(sortutil.Int64Slice(times))
		return times
	}()
	// 이 경우, provider.snapshots 이 동일키에 대해서 변경되었을때 최신걸로 가져옴
	// (문제가 되지 않고 오히려 좋을듯? 하여서 일단은 그대로 구현)
	// 다른 경우는 해당 키를 가져오기에 (그리고 삭제란게 없기에) 현재의 키들만 가져오면
	// provider.snapshots을 보호할 필요가 없음
	provider.RUnlock()

	//
	arrangement := Arrangement{}
	for i := len(sortedSavedTimes) - 1; 0 <= i; i-- {
		savedTime := sortedSavedTimes[i]
		if savedTime <= at {
			provider.RLock()
			arrangement = provider.snapshots[savedTime]
			provider.RUnlock()
			break
		}
	}

	return arrangement
}

func (provider *InMemoryMapedProvider) SetFrom(from time.Time, arrangement Arrangement) {
	// TODO: sync
	provider.Lock()
	defer provider.Unlock()

	provider.snapshots[from.Unix()] = arrangement
}

var provider Provider

// readFromFile, load to provider
func init() {
	fmt.Println("arrangement's provider init() called")

	provider = NewInMemoryMapedProvider()

	// fill provider
	fillProviderFromSavedFile()
}

var savedFileName = "tArrangements.json"

func newPath(fileName string) string {
	_, curFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("fail getting current file path")
	}
	return path.Join(path.Dir(curFilePath), fileName)
}

func fillProviderFromSavedFile() {
	savedFilePath := newPath(savedFileName)

	fmt.Printf("savedFilePath: %v\n", savedFilePath)

	jsonData, err := ioutil.ReadFile(savedFilePath)
	if err != nil {
		panic("fail reading file")
	}

	tTeamNames := []struct {
		From      time.Time `json:"from"`
		TeamNames []string  `json:"team_names"`
	}{}
	json.Unmarshal(jsonData, &tTeamNames)

	//
	nameToID := map[string]team.ID{}

	teamAccessor := team.GetLoadedInMemoryAccessor()
	for _, id := range teamAccessor.GetExistIDs() {
		team, _ := teamAccessor.Get(id)
		nameToID[team.Name] = id
	}

	for _, tmp := range tTeamNames {
		from, names := tmp.From, tmp.TeamNames

		arrangement := Arrangement{}

		for _, name := range names {
			if id, ok := nameToID[name]; ok {
				arrangement.TeamIDs = append(arrangement.TeamIDs, id)
			} else {
				fmt.Errorf("can not find id using name:[%s]", name)
			}
		}

		provider.SetFrom(from, arrangement)
		fmt.Printf("from: %v, arrangement: %v\n", from, arrangement)
	}
}

func NewInMemoryMapedProvider() *InMemoryMapedProvider {
	return &InMemoryMapedProvider{
		snapshots: map[int64]Arrangement{},
	}
}

func GetLoadedInMemoryMapedProvider() Provider {
	return provider
}
