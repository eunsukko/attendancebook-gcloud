package teamhandlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/models/arrangement"
	"github.com/eunsukko/attendancebook-gcloud/models/team"
	"github.com/eunsukko/attendancebook-gcloud/util/logutil"
	"github.com/gorilla/mux"
)

var teamAccessor team.Accessor
var arrangementProvider arrangement.Provider

// loat from json file
func init() {
	teamAccessor = team.GetLoadedInMemoryAccessor()
	arrangementProvider = arrangement.GetLoadedInMemoryMapedProvider()
}

// reduce typo
var ContentType = "Content-Type"
var ApplicationJson = "application/json"

func gatherTeams(t time.Time) []team.Team {
	var ids []team.ID
	if (t == time.Time{}) {
		ids = teamAccessor.GetExistIDs()
	} else {
		ids = arrangementProvider.GetAt(t).TeamIDs
	}
	teams := make([]team.Team, 0, len(ids))
	for _, id := range ids {
		curTeam, err := teamAccessor.Get(id)
		if err != nil {
			logutil.Error.Printf("can't get team(id: %v)", id)
		}
		teams = append(teams, curTeam)
	}
	return teams
}

type queryJson struct {
	T time.Time `json:"t"`
}

// 이 경우 mock 을 사용해서 올바르게 arrangement가 사용되었는지만 확인해도 되지 않을까??
func GetTeams(w http.ResponseWriter, r *http.Request) {
	logutil.Info.Printf("[GetTeams] called\n")

	jsonData := r.FormValue("json_data")

	logutil.Trace.Printf("json_data: %v\n", jsonData)

	queryParms := queryJson{}
	json.Unmarshal([]byte(jsonData), &queryParms)

	t := queryParms.T
	logutil.Trace.Printf("t: %v\n", t)

	teams := gatherTeams(t)
	json.NewEncoder(w).Encode(teams)
}

// get the team
func GetTeam(w http.ResponseWriter, r *http.Request) {
	logutil.Info.Printf("[GetTeam] called\n")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logutil.Error.Printf("%v\n", err.Error())
		http.NotFound(w, r)
		return
	}

	logutil.Trace.Printf("[Get] id: %v\n", id)

	curTeam, _ := teamAccessor.Get(team.ID(id))

	logutil.Trace.Printf("team: %v\n", curTeam)

	json.NewEncoder(w).Encode(curTeam)
}
