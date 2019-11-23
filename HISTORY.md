# 왜 왜 왜??!
이 프로젝트가 어떤 것을 해결하려고 하는지
그리고 현재 상태 (혹은 버전)은 어떤 방식으로 그 문제들을 해결하고 있는지를 정리하기 위한 HISTORY.md

# 20190228(목)
서버에서 handler관련 명세는 해당 테스트로 작성해보는 것은 어떠할까?? (어느 request 에 대한 response 처리를)
(https://blog.questionable.services/article/testing-http-handlers-go/)
로그에 level을 두는 방법? (운영하고 모니터링하면서.. 어디에서 문제가 발생하는 건지.. 조금 빨리 파악하려면)

# 20190227(수)
서버란 것은 어떤 역할을 해주는 걸까?? 가지고 있는 데이터를 특정한 형태로 제공해 준다라는 관점으로 보면... 어느정도를 제공하는게...

에공... uri로 특정정보 (시간 or id)를 어떻게 잘 전달해줄지??
http://jsoneditoronline.org/doc/index.html#query_parameters

일단은 encoded_url 을 통해서 json 데이터를 보냄 (특정 시간 t, teamID)  
- http.Request.FormValue()에서 url 에서 가져오는 경우 알아서 decoding을 해주는듯함

에공... 짜다가보니.. 로그를 잘 남길 방법이 필요할 듯하다
(최소한 디버깅할 때는 보이도록)

http://changelog.ca/log/2015/03/09/golang

지금은 log를 사용하는데...
문제가.. 로그의 수준이 없음.. (즉 어떤 정보는 항상 보이고 싶고 어떤 애들은 특정 때에만, 예를들어 디버깅, 보이고 싶은데)


돌아가는데 문제가 없으려면 어느정도로는 만들어야 할까??
어느정도는 확인 할 수 있어야 할까? 적어도 네트워크 등에 문제가 없을 경우 잘 동작하려면

일단은 서버가 노출하는 api에 대해서 어떤 요청이 있을 것이며 어떤 책임을 져야하는지 내가 확실하게 이해하고 구현해야 할 듯
 > 조금은 더럽더라도.. 일단은 돌아가게 만들어보고 다른 것들을 살핀 후에 더 잘 돌도록 만들까??



# 20190226(화)
attendancebook 관련 개발  
mux 이용한 처리 (uri + 파라미터 사용)
> https://gowebexamples.com/routes-using-gorilla-mux/
클라이언트와 서버에서 주고 받을때...
> https://stackoverflow.com/questions/52789217/flutter-parsing-json-with-datetime-from-golang-rfc3339-formatexception-invalid (시간관련)
간략한 서버 테스트 (json데이터를 보내서 제대로 ab를 생성하고 업데이트하는지)


# 20190224(일)
혼자서 개발하더라도 미래의 나를 위해서 commit message를 잘 써야할 듯하다  
(지금은 먼가.. 메세지에 통일성도 없고.. 어떤걸 남겨야 좋은지도 모르겠다. 밑에 링크를 참고하자)  
https://item4.github.io/2016-11-01/How-to-Write-a-Git-Commit-Message/
https://github.com/spring-projects/spring-framework/blob/30bce7/CONTRIBUTING.md#format-commit-messages


# 20190223(토)
dateutil 작성

사용되는 시간이 주로 kst 시간일 것이디  
그런데 서버시간등은 같은시간이여도 local의 위치가 다를 수도 있을 듯 하다  
따라서 이 프로젝트에서 사용되는 시간이 KST기준으로 돌아가도록 통일하기 위해서 dateutil을 만듬



# 20190221(목)
드디어.. attendancebook 객체를 다뤄보려고 한다

# 20190220(수)

## json 결과를 이쁘게보려면
https://jsonformatter.curiousconcept.com면

# 20190218(월)

# issue 3, https://github.com/eunsukko/attendancebook/issues/3
인터페이스에 대한 고민

## key를 time으로 가지는 구현체에 대한 고민
https://github.com/m3db/build-tools/blob/master/linters/badtime/README.md
https://golang.org/pkg/time/


# 20190215(금)
이 서버에서는 어떤 데이터들을 잘 모아놓고 잘 찾을 수 있게 하는 작업이 주 작업인듯하다.
내부에서는 각 정보를 어떻게 구분하고 (ex. id) 해당 패키지 외부에서는 어떻게 특정 애를 쉽게 추출할지 (id를 몰라도 원하는 id 를 찾아낼 수 있게?)

## 특정 시간에 대한 정보를 얻어내자..!!
오늘은 eventlist, arrangement 를 다뤄보려고 한다.
둘은 공통점이 특정 시간을 포함하는 정보를 얻어내야한다는 것.

- 이름에 대한 고민 (특정 시간의 정보를 나타내는)
- 시간 부분에 대한 고민 (해당 주를 포함하는 모듈을 설계해야할 것 같다 + 테스트)


## 패키지 이름의 중요성
https://blog.golang.org/package-names

## 스스로에게 문제를 주기?? (이슈제공)
무엇이 문제인지 다시 한 번 생각해보고 내가 전에 구현한 설계를 다시 고민해보는?


# 20190214(목)
일단 서버를 중단하기 전에 처리하기 위한 기반 코드를 만들어 본다
(간단하게 종료하려는 signal을 interrupt해서리 그때 .. 필요한 작업(ex. 들고 있는 데이터를 파일에 쓰기)을 적용한다

## team 객체에 대한 고민
어떤 정보들이 있어야 할까?
어떻게 추가하고 관리되어야 할까?
외부에서는 team을 어떻게 접근해야할까?

고민을 했는데 member를 id로 접근하는게 나을듯하다
추가 제거 하는 과정을 위해서 name,birth 로 접근하려고 했는데
그냥 name,birth로 id를 찾기 쉽게 해주는 객체가 있으면 될듯
(그 겍체를 통해서 쉽게 팀을 생성할 수 있도록)

예상되는 시나리오는 크게 3가지이다
1. 팀을 생성할 때: 새로운 해가 될 때 누군가가 각 팀을 구성할 듯 (이때 어느팀에 누가 속하는지 이름들을 알고 있을테니) + 동일한 포함기간으로 구성 // 일단은 이름이 중복되는 인원은 생일로 직접 구분해서 제외하기로 
2. 특정 인원 추가: id, 해당인원이 포함되는 기간등을 가지고 직접 추가해준다 (이정도 인원은 직접 타이핑해서 추가해도 오버헤드가 크지 않을듯)
3. 특정 인원 수정 (제외 or 팀에서 태그 정보 수정) 

### teams.json 을 생성하는 방법



# 20190213(수)

gracefully-shutdown 이라는 개념을 (go 1.8 추가) 정리해보자
이를 통해서 (혹은 interrupt 를 잡는 방식으로) 종료할 때 데이터를 다시 파일에 쓰는 방법을 적용할 수있을듯하다
그러면 일단.. db를 쓰지 않아도 정상 종료시... 현재 존재하는 데이터를 파일에 일단 쓰고 (백업할 방식도 고려는 해보자) 킬 때 다시 복원할 수 있을 듯 하다
https://stackoverflow.com/questions/43631854/gracefully-shutdown-gorilla-server

+ 하다보니 go context 를 살펴보게 되있다.. ㄷㄷ
go context.go 코드 
https://jaehue.github.io/post/how-to-use-golang-context/
https://blog.golang.org/context
https://blog.golang.org/pipelines
https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39


## 개발을 어떻게 진행할지<b>
혼자 하기에... 혼자여도 어떻게 하면 좀 더 체계적으로 할지 고민이다. 
이슈기반으로 개발하기?? (pull request를 이슈기반으로 + commit 에 이슈 번호 넣기)
https://www.popit.kr/%EC%BD%94%EB%93%9C-%EB%A6%AC%EB%B7%B0-%EC%9D%B4%EC%95%BC%EA%B8%B0-1/

> 좋은 개발자로 성장하기 위해서는 여러가지 방법이 있겠지만 필자는 지속적인 글쓰기가 가장 좋은 방법이라고 계속해서 주장해 왔습니다. 어떤 개발자가 좋은 개발자냐 라는 것이 다양한 의견이 있겠지만 좋은 개발자라면 자신의 기술을 체계적으로 설명하고, 정리할 수 있어야 하며, 대중(다른 개발자들)에게 어느 정도 인지도도 있어야 한다고 생각합니다.  이것을 할 수 있는 가장 좋은 방법이 꾸준한 글쓰기 입니다.

## 로그의 필요성?? 
logus??
아직은 어떤 애들이 좋은건지 판단할 실력이 없다.
그러니.. 일단 써본다?? 
https://jaehue.github.io/post/go-my-way-2-database-and-logging/

# 20190212(화)
내가 아는 문제상황을 막 적어본다

문제상황 + 그에 관한 사실(혹은 요구사항)
현재 각 셀에서 각각의 출석부를 가지고 있고 (물리적인 출석부) 이를 총무단에서 관리하고 있음
내가 보기에 물리적인 출석부의 문제를 정리하면
    1. 매주 총무단과 각 셀이 주고 받아야 하는 것 (물리적인 주고받음 + 기다림)
    2. 매달 새로운 출석부를 만들어내야함
    3. 정보(사람, 셀 등)가 추가, 수정되었을 때 출석부에 반영

일단은 이세종목사님께 물어본바로는 각각의 출석데이터가 남을 필요는 없기에 (셀에 몇명이 온건지가 현재는 중요)
'앱을 통한 출석부 + 결과를 모아주는 역할' 을 현재의 목표로 잡음
(그리고 서버를 운영하는 사람이 쉽고 실수하지 않게 정보들을 추가, 수정 할 수 있는 설계)