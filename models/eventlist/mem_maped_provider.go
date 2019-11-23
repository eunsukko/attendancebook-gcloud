package eventlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/models/utils/sortutil"
)

type InMemoryMapedProvider struct {
	sync.RWMutex
	snapshots map[int64]Eventlist
}

func (provider *InMemoryMapedProvider) GetAt(t time.Time) Eventlist {
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

	fmt.Printf("len sortedSavedTimes: %d\n", len(sortedSavedTimes))
	//
	elist := Eventlist{}
	for i := len(sortedSavedTimes) - 1; 0 <= i; i-- {
		savedTime := sortedSavedTimes[i]
		fmt.Printf("at: %v, savedTime: %v\n", at, savedTime)
		if savedTime <= at {
			provider.RLock()
			elist = provider.snapshots[savedTime]
			provider.RUnlock()
			break
		}
	}

	return elist
}

func (provider *InMemoryMapedProvider) SetFrom(from time.Time, elist Eventlist) {
	// TODO: sync
	provider.Lock()
	defer provider.Unlock()

	provider.snapshots[from.Unix()] = elist
}

var provider Provider

// readFromFile, load to provider
func init() {
	fmt.Println("eventlist's provider init() called")

	provider = NewInMemoryMapedProvider()

	// fill provider
	fillProviderFromSavedFile()
}

var savedFileName = "tEventlists.json"

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

	tEventlists := []struct {
		From      time.Time `json:"from"`
		Eventlist Eventlist `json:"eventlist"`
	}{}
	json.Unmarshal(jsonData, &tEventlists)

	for _, tEventlist := range tEventlists {
		from, elist := tEventlist.From, tEventlist.Eventlist
		provider.SetFrom(from, elist)
		fmt.Printf("from: %v, elist: %v\n", from, elist)
	}
}

func NewInMemoryMapedProvider() Provider {
	return &InMemoryMapedProvider{
		snapshots: map[int64]Eventlist{},
	}
}

func GetLoadedInMemoryMapedProvider() Provider {
	return provider
}
