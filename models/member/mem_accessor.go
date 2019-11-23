package member

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
)

// InMemoryAccessor is a simple in memory database.
type InMemoryAccessor struct {
	newID   ID
	members map[ID]Member
	ids     []ID

	mu sync.Mutex
}

func (accessor InMemoryAccessor) GetExistIDs() []ID {
	return accessor.ids
}

// Add insert new member with new unique ID
func (accessor *InMemoryAccessor) Add(member Member) {
	accessor.mu.Lock()
	defer accessor.mu.Unlock()

	member.ID = accessor.newID
	accessor.ids = append(accessor.ids, accessor.newID)
	accessor.newID++

	accessor.members[member.ID] = member
}

var accessor InMemoryAccessor

func init() {
	fmt.Println("member mem_accessor.go's init() called")
	fmt.Println("laod accessor")

	readFromSavedFile()

	// write to file
	// writeToSavedFile()
}

var savedFileName = "members.json"

func newPath(fileName string) string {
	_, curFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("fail getting current file path")
	}
	return path.Join(path.Dir(curFilePath), fileName)
}

func readFromSavedFile() {
	accessor = InMemoryAccessor{
		newID:   ID(1000),
		members: map[ID]Member{},
		ids:     []ID{},
	}

	savedFilePath := newPath(savedFileName)

	fmt.Printf("savedFilePath: %v\n", savedFilePath)

	jsonData, err := ioutil.ReadFile(savedFilePath)
	if err != nil {
		panic("fail reading file")
	}

	memberList := make([]Member, 0)
	json.Unmarshal(jsonData, &memberList)

	for _, member := range memberList {
		if len(member.Info["name"]) == 0 {
			continue
		}

		//
		if _, ok := member.Info["phoneNumber"]; ok {
			delete(member.Info, "phoneNumber")
		}
		if _, ok := member.Info["phone_number"]; ok {
			delete(member.Info, "phone_number")
		}

		accessor.Add(member)
		// fmt.Printf("put member(%v): %v\n", nb, member)
	}
}

func writeToSavedFile() {
	f, err := os.Create(newPath(savedFileName))
	if err != nil {
		panic("fail create members.json file")
	}
	defer f.Close()

	insertedMemberList := []struct {
		Info map[string]string `json:"info"`
	}{}
	for _, id := range accessor.GetExistIDs() {
		member, _ := accessor.Get(id)
		birth := member.Info["birth"]
		birth = strings.Replace(birth, ".", "", -1)
		member.Info["birth"] = birth
		if len(birth) != 6 {
			fmt.Printf("name: %v, birth: %v\n", member.Info["name"], birth)
		}
		insertedMemberList = append(insertedMemberList, struct {
			Info map[string]string `json:"info"`
		}{member.Info})
	}

	jsonData, _ := json.Marshal(insertedMemberList)
	f.Write(jsonData)
}

// GetLoadedInMemoryAccessor returns the accessor loaded from saved members.json file.
func GetLoadedInMemoryAccessor() Accessor {
	return &accessor
}

// NewInMemoryAccessor returns a new InMemoryAccessor.
func NewInMemoryAccessor() Accessor {
	return &InMemoryAccessor{
		members: map[ID]Member{},
	}
}

// Get returns a member with a given NameBirth.
func (m *InMemoryAccessor) Get(id ID) (Member, error) {
	member, exists := m.members[id]
	if !exists {
		return Member{}, ErrMemberNotExist
	}
	return member, nil
}

// func (m *InMemoryAccessor) Put(id ID, m Member) error {
// }

// func (m *InMemoryAccessor) Post(member Member) (ID, error){
// }

// func (m *InMemoryAccessor) Delete(id ID) error{
// }

// func Example_Member() {
// 	member := Member{}

// 	fmt.Printf("%v\n", member)
// 	fmt.Printf("cnt ids: %d\n", len(accessor.GetExistIDs()))

// 	// Output:
// 	// .
// }
