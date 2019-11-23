package main

import (
	"fmt"

	"github.com/eunsukko/attendancebook-gcloud/models/member"
)

var idCounter = map[string][]member.ID{}

// to figure out who have duplicated name in members.json
func main() {

	accessor := member.GetLoadedInMemoryAccessor()

	for _, id := range accessor.GetExistIDs() {
		curMember, _ := accessor.Get(id)
		name := curMember.Info["name"]
		if idCounter[name] != nil {
			idCounter[name] = append(idCounter[name], id)
		} else {
			idCounter[name] = []member.ID{id}
		}
	}

	//
	for k, v := range idCounter {
		if 1 < len(v) {
			fmt.Printf("%s\n", k)
			for _, id := range v {
				curMember, _ := accessor.Get(id)
				fmt.Printf("%v\n", curMember)
			}
			fmt.Printf("\n")
		}
	}
}
