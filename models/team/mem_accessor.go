package team

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/eunsukko/attendancebook-gcloud/models/member"
	"github.com/eunsukko/attendancebook-gcloud/util/periodutil"
)

var accessor *InMemoryAccessor

type InMemoryAccessor struct {
	newID ID
	teams map[ID]Team
	ids   []ID

	mu sync.Mutex
}

func GetLoadedInMemoryAccessor() Accessor {
	return accessor
}

func NewInMemoryAccessor() *InMemoryAccessor {
	return &InMemoryAccessor{
		newID: 8000,
		teams: map[ID]Team{},
		ids:   []ID{},
	}
}

func (accessor InMemoryAccessor) GetExistIDs() []ID {
	return accessor.ids
}

// Add insert new team with new unique ID
func (accessor *InMemoryAccessor) Add(team Team) {
	accessor.mu.Lock()
	defer accessor.mu.Unlock()

	team.ID = accessor.newID
	accessor.ids = append(accessor.ids, accessor.newID)
	accessor.newID++

	accessor.teams[team.ID] = team
}

func (accessor *InMemoryAccessor) Get(id ID) (Team, error) {
	accessor.mu.Lock()
	defer accessor.mu.Unlock()

	return accessor.teams[id], nil
}

func init() {
	fmt.Println("team mem_accessor.go's init() called")
	fmt.Println("laod accessor")

	err := generateTeamsJSONFile()
	if err != nil {
		fmt.Errorf("err: %v", err.Error())
		return
	}

	accessor = NewInMemoryAccessor()

	fillAccessorFromSavedFile()

	// write to file
	// writeToSavedFile()
}

var savedFileName = "teams.json"
var generationFileName = "team_generation.json"

func newPath(fileName string) string {
	_, curFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("fail getting current file path")
	}
	return path.Join(path.Dir(curFilePath), fileName)
}

var errGenerateTeamFail = errors.New("there is ambigous member (ex. member who has same name")

type memberInfo struct {
	Name  string   `json:"name,omitempty"`
	Birth string   `json:"birth,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

type teamGenerationInfo struct {
	Name string `json:"name"`

	MemberInfos       []memberInfo      `json:"member_infos"`
	MemberValidPeriod periodutil.Period `json:"member_valid_period"`
}

func (info teamGenerationInfo) String() string {
	data, _ := json.Marshal(info)
	return string(data)
}

var errDuplicatedName = errors.New("duplicated name")

var errDuplicatedMember = errors.New("duplicated members exist")
var errNotExistMember = errors.New("can not find a member using this information")

func generateTeamsJSONFile() error {
	generationFilePath := newPath(generationFileName)

	jsonData, err := ioutil.ReadFile(generationFilePath)
	if err != nil {
		fmt.Errorf("can find file: %v\n", generationFilePath)
	}

	infos := make([]teamGenerationInfo, 0)
	json.Unmarshal(jsonData, &infos)

	fmt.Printf("infos: %v\n", infos)
	accessor := member.GetLoadedInMemoryAccessor()

	nameToIDs, birthToIDs := func() (map[string][]member.ID, map[string][]member.ID) {
		_nameToIDs, _birthToIDs := map[string][]member.ID{}, map[string][]member.ID{}

		for _, id := range accessor.GetExistIDs() {
			curMember, _ := accessor.Get(id)
			name := curMember.Info["name"]
			birth := curMember.Info["birth"] // 0 인 애들도 존재함
			if _, ok := _nameToIDs[name]; ok {
				_nameToIDs[name] = append(_nameToIDs[name], id)
			} else {
				_nameToIDs[name] = []member.ID{id}
			}

			if _, ok := _birthToIDs[birth]; ok {
				_birthToIDs[birth] = append(_birthToIDs[birth], id)
			} else {
				_birthToIDs[birth] = []member.ID{id}
			}
		}
		return _nameToIDs, _birthToIDs
	}()

	getID := func(mInfo memberInfo) (member.ID, error) {
		name := mInfo.Name
		birth := mInfo.Birth

		ids := []member.ID{}
		if birth == "" {
			ids = nameToIDs[name]
		} else {
			idSet := make(map[member.ID]struct{})
			for _, id := range nameToIDs[name] {
				idSet[id] = struct{}{}
			}

			for _, id := range birthToIDs[birth] {
				if _, ok := idSet[id]; !ok {
					continue
				}
				ids = append(ids, id)
			}
		}

		if len(ids) == 0 {
			return member.ID(-1), errNotExistMember
		} else if 1 < len(ids) {
			return member.ID(-1), errDuplicatedMember
		}
		return ids[0], nil
	}

	// 확실한 memberInfo 를 제공한다고 가정함 (즉 무조건 존재하는 사람 + 동일 이름이 존재할 때 구분할 수 있는 정보 (ex, birth))
	fmt.Printf("generate teams\n")
	//
	wg := sync.WaitGroup{}
	teams := make([]Team, len(infos))
	for i, info := range infos {
		wg.Add(1)
		go func(i int, info teamGenerationInfo) {
			defer wg.Done()

			curTeam := NewTeam()

			curTeam.Name = info.Name
			curTeam.MemberValidPeriods = func() []periodutil.Period {
				periods := make([]periodutil.Period, len(info.MemberInfos))
				for i := range info.MemberInfos {
					periods[i] = info.MemberValidPeriod
				}
				return periods
			}()

			// get ids
			memberIds := make([]member.ID, len(info.MemberInfos))
			for i, mInfo := range info.MemberInfos {
				id, err := getID(mInfo)
				if err != nil {
					fmt.Errorf("%v", err.Error())
					continue
				}

				memberIds[i] = id
				if mInfo.Tags != nil {
					curTeam.Tags[id] = mInfo.Tags
				}
			}
			curTeam.MemberIDs = memberIds

			// 밖에 있는 teams을 동시에 접근하면 무슨 일이 있을까?
			teams[i] = curTeam
		}(i, info)
	}

	wg.Wait()

	fmt.Printf("%v\n", teams)
	savedFilePath := newPath(savedFileName)
	fmt.Printf("write generated teams to savedFilePath(%v)\n", savedFilePath)

	writeToSavedFile(teams)
	return nil
}

func fillAccessorFromSavedFile() {
	savedFilePath := newPath(savedFileName)

	jsonData, err := ioutil.ReadFile(savedFilePath)
	if err != nil {
		fmt.Errorf("can't find file: %v\n", savedFilePath)
	}

	teams := make([]Team, 0)
	json.Unmarshal(jsonData, &teams)

	for _, team := range teams {
		accessor.Add(team)
		fmt.Printf("add team: %v\n", team)
	}
}

func writeToSavedFile(teams []Team) {
	f, err := os.Create(newPath(savedFileName))
	if err != nil {
		panic("fail create teams.json file")
	}
	defer f.Close()

	jsonData, _ := json.Marshal(teams)
	f.Write(jsonData)
}
