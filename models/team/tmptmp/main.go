package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/eunsukko/attendancebook-gcloud/models/member"
)

var lineToken = "==="

type memberInfo struct {
	Name  string   `json:"name,omitempty"`
	Birth string   `json:"birth,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// 원하는 것
// 특정 파일에서 이름들을 계속 읽음
func main() {
	reader := bufio.NewReader(os.Stdin)

	nameToBirths := func() map[string][]string {
		m := map[string][]string{}

		accessor := member.GetLoadedInMemoryAccessor()

		for _, id := range accessor.GetExistIDs() {
			curMember, _ := accessor.Get(id)
			name := curMember.Info["name"]
			birth := curMember.Info["birth"]
			if _, ok := m[name]; ok {
				m[name] = append(m[name], birth)
			} else {
				m[name] = []string{birth}
			}
		}

		return m
	}()

	teamMemberInfos := [][]memberInfo{}

	fmt.Printf("%v\n", nameToBirths)
	fmt.Printf("read start\n")
	for {
		line, _, err := reader.ReadLine()
		if line == nil {
			fmt.Printf("read all \n")
			break
		}
		if err != nil {
			fmt.Errorf("error \n")
		}
		name := string(line)
		if 0 < strings.Count(name, lineToken) {
			teamMemberInfos = append(teamMemberInfos, []memberInfo{})
			continue
		}
		fmt.Printf("%s\n", name)

		mInfos := []memberInfo{}
		for _, birth := range nameToBirths[name] {
			mInfos = append(mInfos, memberInfo{
				Name:  name,
				Birth: birth,
			})
		}

		teamMemberInfos[len(teamMemberInfos)-1] = append(teamMemberInfos[len(teamMemberInfos)-1], mInfos...)
	}
	jsonData, _ := json.Marshal(teamMemberInfos)
	fmt.Printf("%v\n", string(jsonData))
}
