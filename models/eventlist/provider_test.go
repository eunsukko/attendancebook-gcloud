package eventlist

import (
	"reflect"
	"testing"
	"time"
)

func NewDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func Test(t *testing.T) {
	provider := NewInMemoryMapedProvider()

	tes := []struct {
		From      time.Time
		Eventlist Eventlist
	}{
		{
			From:      NewDate(2017, 1, 1),
			Eventlist: Eventlist{"금철", "대예배", "5부예배", "셀모임"},
		},
		{
			From:      NewDate(2018, 12, 24),
			Eventlist: Eventlist{"금철", "대예배", "5부예배", "셀모임", "크리스마스"},
		},
		{
			From:      NewDate(2018, 12, 31),
			Eventlist: Eventlist{"금철", "대예배", "5부예배", "셀모임", "송구영신"},
		},
	}

	for i := range tes {
		provider.SetFrom(tes[i].From, tes[i].Eventlist)
	}

	cases := []struct {
		in   time.Time
		want Eventlist
	}{
		{NewDate(2017, 12, 31), Eventlist{"금철", "대예배", "5부예배", "셀모임"}},
		{NewDate(2018, 12, 23), Eventlist{"금철", "대예배", "5부예배", "셀모임"}},
		//
		{NewDate(2018, 12, 24), Eventlist{"금철", "대예배", "5부예배", "셀모임", "크리스마스"}},
		{NewDate(2018, 12, 25), Eventlist{"금철", "대예배", "5부예배", "셀모임", "크리스마스"}},
		{NewDate(2018, 12, 30), Eventlist{"금철", "대예배", "5부예배", "셀모임", "크리스마스"}},
		//
		{NewDate(2018, 12, 31), Eventlist{"금철", "대예배", "5부예배", "셀모임", "송구영신"}},
		{NewDate(2019, 1, 1), Eventlist{"금철", "대예배", "5부예배", "셀모임", "송구영신"}},
	}

	for _, c := range cases {
		in, want := c.in, c.want
		got := provider.GetAt(in)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("at: %v, got: %v, want: %v\n", in, got, want)
		}
	}
}
