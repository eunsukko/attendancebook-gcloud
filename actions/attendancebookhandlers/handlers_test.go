package attendancebookhandlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// 어떻게 공부해야할까?
// https://blog.questionable.services/article/testing-http-handlers-go/
//
// 어떻게 테스트할 수있는 코드를 만드는지.. 그리고 외부에 있는 애를 어떻게 테스트 할 수 있게 하는지
// 꽤 재미있는 예제인듯 + 어떻게 내 것을 검증해가면서 코딩을 할지를..
// 어떤 것을 테스트 해야하는 것인지를 정리해가는 것도 좋을 것 같다
// (용어같은 부분도 내것으로 만들기 위해선)
// https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/mocking
// 잘 짰으면 mocking이 필요한부분이 적다고 얘기하는 듯함

// 테스트에서의 구분: 상태검증(stub + )과 행위검증(mock)
// https://martinfowler.com/articles/mocksArentStubs.html

// 어떤 것을 테스트 할 수 있는지
// 그리고 각각을 테스트할 때 어떤 것들이 (ex, mock(spy, stub))  왜 필요한 건지를 파악하자

// id 관련 테스트
// id 가 올바를때
// id 가 not integer
// id 가 음수
// >> 애러코드로 확인하자

// 특정.... 데이터 셋을 만들어 놓고
// 특정 핸들러가 동작할 때 만들어진 데이터 셋에 따라 올바른 결과를 만드는지 확인

// mux testcode 참조 (https://github.com/gorilla/mux)

// 먼가.. 할 테스트들이 각각의 provider를 만들어야 하는 데... 그것들을 어떻게 쉽게하고 잘 알아보게 하고 확인할지
func TestGetAttendancebook(t *testing.T) {
	tt := []struct {
		routeVariable string
		want          string // Json string
	}{}

	for i, tc := range tt {
		path := fmt.Sprintf("/attendancebook/%s", tc.routeVariable)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		// Nedd to create a router that we can pass the request throught so that the vars will be added to the context
		router := mux.NewRouter()
		router.HandleFunc("/attendancebook/{id:[0-9]+}", GetAttendancebook).Methods("GET")
		router.ServeHTTP(rr, req)

		// json 결과를 테스트 하려면?? 동일한 json인지 테스트 하는 방법이 있으려나?
		// rr.Body
		// if rr.Code == http.StatusOK && !tc.want
		t.Errorf("[%d] want: %q, got: %q\n", i, tc.want, tc.want)
	}

}

// 각 함수마다 무슨 테스트를 해야할 까??
// 일단은 대략은 어떤 식으로 어떤 것을 테스트 해야 할지 살펴봄...

// 일단은 mux 에서 원하지 않는 경로로 오는 경우는 404처리가 되고 있음 + (ex. id 가 숫자가 아닌 경우)
// 그렇다면 핸들러로 호출이 된 상황일 때 내가 해당 함수에게 원하는 책임은 무엇일까?
// 왜 그런 것을 요구할까?

// 이 테스트들에서 어려운 점은... 일단 attendancebook accessor를 만들어 줘야하는 것이다
// 테스트용 데이터를 만들어야 하는데... 그것을 어떻게 쉽게 하는 방법이 있으려나?
// 그러면서도 조금은 더 의미있는 테스트를 할 수 있는 방법이??
// 어떻게 보면 db를 mock 으로 들고 있는 느낌도 있고

// func QueryID(w http.ResponseWriter, r *http.Request)
// 무슨 테스트를 해야할까?

// func GetAttendancebook(w http.ResponseWriter, r *http.Request)
// 무슨 테스트를 해야할까?

// func GetAttendancebooks(w http.ResponseWriter, r *http.Request)
// 무슨 테스트를 해야할까?

// func PutAttendancebook(w http.ResponseWriter, r *http.Request)
// 무슨 테스트를 해야할까?
