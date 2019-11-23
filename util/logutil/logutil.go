package logutil

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger // 기타 모든 로그 (디버깅할 때 정도 보였으면 하는 정보)
	Info    *log.Logger // 중요한 정보 (항상 보였으면 하는 정보, 지금 어떤 부분이 동작하는지 실시간으로 보기위함)
	Warning *log.Logger // 경고성 정보 (돌리기는 할건데.. 예상치못한 결과가 나올 경우)
	Error   *log.Logger // 치명적인 오류 (망했다... 인 경우)
)

func init() {
	file, err := os.OpenFile("errors.txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("can not open errors.txt", err)
	}

	// Trace = log.New(ioutil.Discard,
	Trace = log.New(os.Stdout,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// // out 은 외부에서 log를 사용할 때 어디로 메세지를 배출할지를 한 곳에서 정하기 위함
// // 일단 디버깅일 때와 아닐때 정도로.. 구분하려고 함
// var out io.Writer = ioutil.Discard

// var logMode = "stdout"

// var Logger *log.Logger

// func InitOut() {
// 	switch logMode {
// 	case "file":
// 		// 지금은 파일에는 쓰지 않음
// 		out = os.Stdout
// 	case "stdout":
// 		out = os.Stdout
// 	default:
// 		out = ioutil.Discard
// 	}

// 	Logger = log.New(out, "", log.Lshortfile)
// }

// func init() {
// 	InitOut()
// }
