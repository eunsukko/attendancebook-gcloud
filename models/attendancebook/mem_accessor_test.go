package attendancebook

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/eunsukko/attendancebook-gcloud/models/team"
)

// 동시성도 테스트 하고 싶으면 어떻게 해야 할까용

func TestQueryIDGenerateNewIDAsync(t *testing.T) {
	testCase := 1
	n := 100

	accessor := NewInMemoryAccessor()

	hasDuplicatedID := func(cnt map[ID]int) bool {
		for _, v := range cnt {
			if 1 < v {
				return false
			}
		}
		return true
	}

	for tc := 0; tc < testCase; tc++ {
		var wg = sync.WaitGroup{}
		mu := sync.Mutex{}

		cnt := make(map[ID]int)
		wg.Add(n)
		for i := 1; i <= n; i++ {
			go func(i int) {
				defer wg.Done()
				t := time.Date(0, 0, 0, 0, 0, i, 100, time.UTC)
				teamID := team.ID(i)

				id := accessor.QueryID(t, teamID)

				mu.Lock()
				cnt[id] = cnt[id] + 1
				mu.Unlock()
			}(i)
		}
		wg.Wait()

		if !hasDuplicatedID(cnt) {
			t.Errorf("[%d] gerate duplicated id\n", tc)
		}
	}
}

func TestQueryIDGenerateNewID(t *testing.T) {
	testCase := 1
	n := 100

	accessor := NewInMemoryAccessor()

	hasDuplicatedID := func(cnt map[ID]int) bool {
		for _, v := range cnt {
			if 1 < v {
				return false
			}
		}
		return true
	}

	for tc := 0; tc < testCase; tc++ {
		cnt := make(map[ID]int)
		for i := 1; i <= n; i++ {
			t := time.Date(0, 0, 0, 0, 0, i, 100, time.UTC)
			teamID := team.ID(i)

			id := accessor.QueryID(t, teamID)
			cnt[id] = cnt[id] + 1
		}

		if !hasDuplicatedID(cnt) {
			t.Errorf("[%d] gerate duplicated id\n", tc)
		}
	}
}

func TestQueryIDFindID(t *testing.T) {
	// init for test
	n := 100
	randN := 100

	accessor := NewInMemoryAccessor()
	cases := func() []struct {
		in   bookInfo
		want ID
	} {
		ret := []struct {
			in   bookInfo
			want ID
		}{}

		for i := 1; i <= n; i++ {
			// 지금은 zero value 가 *의 역할을 하기 때문에
			t := time.Date(0, 0, 0, 0, 0, i, 100, time.UTC)
			teamID := team.ID(i)

			id := accessor.QueryID(t, teamID)

			ret = append(ret, struct {
				in   bookInfo
				want ID
			}{bookInfo{t, teamID}, id},
			)
		}

		return ret
	}()

	//
	fmt.Printf("In sequential find")
	for i, c := range cases {
		in, want := c.in, c.want

		got := accessor.QueryID(in.t, in.teamID)

		if want != got {
			t.Errorf("[%d] in: %v, want: %v, got: %v\n", i, in, want, got)
		}
	}

	fmt.Printf("In random find")
	for i := 0; i < randN; i++ {
		c := cases[rand.Intn(len(cases))]

		in, want := c.in, c.want

		got := accessor.QueryID(in.t, in.teamID)

		if want != got {
			t.Errorf("[%d] in: %v, want: %v, got: %v\n", i, in, want, got)
		}
	}
}

func TestQueryIDsFindID(t *testing.T) {
	// init for test
	tN := 10
	teamN := 10

	accessor := NewInMemoryAccessor()
	cases := func() []struct {
		in   bookInfo
		want []ID
	} {
		ret := []struct {
			in   bookInfo
			want []ID
		}{}

		for tV := 1; tV <= tN; tV++ {
			for teamV := 1; teamV <= teamN; teamV++ {
				t := time.Date(0, 0, 0, 0, 0, tV, 100, time.UTC)
				teamID := team.ID(teamV)

				id := accessor.QueryID(t, teamID)

				ret = append(ret, struct {
					in   bookInfo
					want []ID
				}{bookInfo{t, teamID}, []ID{id}},
				)
			}
		}

		return ret
	}()

	for i, c := range cases {
		in, want := c.in, c.want

		got := accessor.QueryIDs(in.t, in.teamID)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("[%d] in: %v, want: %v, got: %v\n", i, in, want, got)
		}
	}
}

type IDSlice []ID

func (ids IDSlice) Len() int {
	return len(ids)
}

func (ids IDSlice) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}

func (ids IDSlice) Less(i, j int) bool {
	return ids[i] < ids[j]
}

func TestQueryIDsFindIDsAtSpecificTime(t *testing.T) {
	// init for test
	tN := 10
	teamN := 10

	accessor := NewInMemoryAccessor()
	cases := func() []struct {
		in   time.Time
		want []ID
	} {
		ret := []struct {
			in   time.Time
			want []ID
		}{}

		for tV := 1; tV <= tN; tV++ {
			t := time.Date(0, 0, 0, 0, 0, tV, 100, time.UTC)
			ids := []ID{}
			for teamV := 1; teamV <= teamN; teamV++ {

				teamID := team.ID(teamV)

				id := accessor.QueryID(t, teamID)

				ids = append(ids, id)
			}
			ret = append(ret, struct {
				in   time.Time
				want []ID
			}{t, ids})
		}

		return ret
	}()

	for i, c := range cases {
		in, want := c.in, c.want

		got := accessor.QueryIDs(in, team.ID(0))

		sort.Sort(IDSlice(want))
		sort.Sort(IDSlice(got))
		if reflect.DeepEqual(want, got) {
			t.Errorf("[%d] in: %v, want: %v, got: %v\n", i, in, want, got)
		}
	}
}
