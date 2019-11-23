package memberhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eunsukko/attendancebook-gcloud/models/member"
	"github.com/eunsukko/attendancebook-gcloud/util/logutil"
	"github.com/gorilla/mux"
)

var accessor member.Accessor

// loat from json file
func init() {
	accessor = member.GetLoadedInMemoryAccessor()
}

// get all members info
func GetMembers(w http.ResponseWriter, r *http.Request) {
	logutil.Info.Printf("[GetMembers] called\n")

	members := make([]member.Member, 0)
	for _, id := range accessor.GetExistIDs() {
		member, err := accessor.Get(id)
		if err != nil {
			fmt.Errorf("%v", err.Error())
		}
		members = append(members, member)
	}
	json.NewEncoder(w).Encode(members)
}

// get the member
func GetMember(w http.ResponseWriter, r *http.Request) {
	logutil.Info.Printf("[GetMember] called\n")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logutil.Error.Printf("%v\n", err.Error())
		http.NotFound(w, r)
		return
	}

	logutil.Trace.Printf("[Get] id: %v\n", id)

	member, _ := accessor.Get(member.ID(id))

	logutil.Trace.Printf("member: %v\n", member)

	json.NewEncoder(w).Encode(member)
}
