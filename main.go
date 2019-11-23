package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/actions/attendancebookhandlers"
	"github.com/eunsukko/attendancebook-gcloud/actions/memberhandlers"
	"github.com/eunsukko/attendancebook-gcloud/actions/teamhandlers"
	"github.com/gorilla/mux"
)

// 서버에서 특정 요청을 처리중에 shutdown이 되면 어떻게되는 거지??
// 그렇게 되면 데이터가 오염 될 수 있을 듯 한데... (덜 수정된 데이터가 파일로 저장될테니??)
// 그런데 사실... 나중에 디비를 사용하면 문제없을 듯 하고
// team, arrangement 같은 데이터는 일단은 추가를 파일로 하므로 문제 없을 듯 하다
//
// 나중에 데이터 추가를 프로그램에서 하게 될 경우, 이 부분을 좀더 신경써야 할 듯
// (현재 진행중인 요청이 다 끝나기를 기다리고 그 결과를 저장하기로)
//
// 으아.. 진짜로 돌아가는 도중에 꺼져도 복구 가능하게 짜는 것... 그게 실력일 것 같다 (그래서 db를 열심히 쓰는건가..ㅋㅋㅋㅋㅋ)
func processForIntentedShutdown(gracefulStop <-chan os.Signal) {
	sig := <-gracefulStop
	fmt.Printf("caught sig: %+v", sig)
	fmt.Println("Wait for 2 second to finish processing")
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

func init() {
	fmt.Printf("ready to intented shutdown\n")
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go processForIntentedShutdown(gracefulStop)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/members", memberhandlers.GetMembers).Methods("GET")
	r.HandleFunc("/member/{id:[0-9]+}", memberhandlers.GetMember).Methods("GET")

	r.HandleFunc("/teams", teamhandlers.GetTeams).Methods("GET")
	r.HandleFunc("/team/{id:[0-9]+}", teamhandlers.GetTeam).Methods("GET")

	r.HandleFunc("/query/attendancebook", attendancebookhandlers.QueryID)

	// s := r.PathPrefix("/attendancebook").Subrouter()
	// s.HandleFunc("/{id:[0-9]+}", attendancebookhandlers.GetAttendancebook).Methods("GET", "PUT")

	r.HandleFunc("/attendancebook/{id:[0-9]+}", attendancebookhandlers.GetAttendancebook).Methods("GET", "PUT")

	r.HandleFunc("/attendancebooks", attendancebookhandlers.GetAttendancebooks).Methods("GET")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", r))
	fmt.Printf("server closed\n")
}
