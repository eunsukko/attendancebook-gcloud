package attendancebookhandlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/models/attendancebook"
	"github.com/eunsukko/attendancebook-gcloud/models/team"
	"github.com/eunsukko/attendancebook-gcloud/util/dateutil"
	"github.com/eunsukko/attendancebook-gcloud/util/logutil"
	"github.com/gorilla/mux"
)

var accessor attendancebook.Accessor

func init() {
	accessor = attendancebook.GetLoadedAccessor()
}

// reduce typo
var ContentType = "Content-Type"
var ApplicationJson = "application/json"

type queryJson struct {
	T      time.Time `json:"t"`
	TeamID team.ID   `json:"team_id"`
}

// ab의 id 를 uri를 사용해서 알아내는게 좋은 방법일까??...
// t, teamID가 존재할 것
func QueryID(w http.ResponseWriter, r *http.Request) {
	tmp := map[string]interface{}{}

	jsonData := r.FormValue("json_data")

	logutil.Trace.Printf("json_data: %v\n", jsonData)

	var queryParms queryJson
	json.Unmarshal([]byte(jsonData), &queryParms)

	if (queryParms.T == time.Time{}) || queryParms.TeamID == team.ID(0) {
		// 잘못 데이터가 온 경우
		logutil.Warning.Printf("empty parameters\n")
		logutil.Warning.Printf("queryParms: %v\n", queryParms)
		http.Error(w, "Has Empty Parameters In Query", 400)
		return
	}

	thisMonday := dateutil.NewKST(queryParms.T).CalcThisMonday().Time

	id := accessor.QueryID(thisMonday, queryParms.TeamID)

	w.Header().Set(ContentType, ApplicationJson)
	tmp["id"] = id
	json.NewEncoder(w).Encode(tmp)
}

//
func GetAttendancebook(w http.ResponseWriter, r *http.Request) {
	logutil.Info.Printf("[GetAttendancebook] called\n")
	id, err := getAttendancebookID(r)
	if err != nil {
		logutil.Error.Printf("%v\n", err.Error())
		http.NotFound(w, r)
		return
	}

	logutil.Trace.Printf("[Get] id: %v\n", id)

	ab, _ := accessor.Get(id)

	logutil.Trace.Printf("ab: %v\n", ab)

	json.NewEncoder(w).Encode(ab)
}

// t가 주어져서 해당 t 일때의 모든 ab를 리턴해주기 위함
func GetAttendancebooks(w http.ResponseWriter, r *http.Request) {
	jsonData := r.FormValue("json_data")

	logutil.Trace.Printf("json_data: %v\n", jsonData)

	var queryParms queryJson
	json.Unmarshal([]byte(jsonData), &queryParms)

	if (queryParms.T == time.Time{}) {
		logutil.Warning.Printf("empty parameters\n")

		// return all ab
		abs := []attendancebook.Attendancebook{}
		for _, id := range accessor.GetExistIDs() {
			ab, _ := accessor.Get(id)
			abs = append(abs, ab)
		}

		json.NewEncoder(w).Encode(abs)
		return
	}

	thisMonday := dateutil.NewKST(queryParms.T).CalcThisMonday().Time

	ids := accessor.QueryIDs(thisMonday, team.ID(0))

	abs := make([]attendancebook.Attendancebook, 0, len(ids))

	for _, id := range ids {
		ab, _ := accessor.Get(id)
		abs = append(abs, ab)
	}

	json.NewEncoder(w).Encode(abs)
}

func PutAttendancebook(w http.ResponseWriter, r *http.Request) {
	id, err := getAttendancebookID(r)
	if err != nil {
		logutil.Error.Printf("%v\n", err.Error())
		http.NotFound(w, r)
		return
	}
	jsonAb := r.FormValue("attendancebook")

	logutil.Trace.Printf("jsonAb: %v\n", jsonAb)

	var ab attendancebook.Attendancebook

	json.Unmarshal([]byte(jsonAb), &ab)

	logutil.Trace.Printf("ab: %v\n", ab)

	accessor.Put(id, ab)

	// 다른서버에서는
	// response 를 어떻게 처리하는지 살펴봐야함
	json.NewEncoder(w).Encode("ok")
}

var ErrIDNotFound = errors.New("Wrong integer or missing ID")

// 이제 main.go 에서 id 가 숫자인지 확인하기에 이제는 이 함수가 필요 없음
func getAttendancebookID(r *http.Request) (attendancebook.ID, error) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		return attendancebook.ID(0), ErrIDNotFound
	}

	return attendancebook.ID(id), nil
}
