package dateutil

import (
	"testing"
	"time"
)

func TestNewKSTIgnorePreviousTimeZone(t *testing.T) {
	cases := []struct {
		in   time.Time
		want string
	}{
		{
			in:   time.Date(2019, 02, 23, -10, 0, 0, 0, time.UTC),
			want: "2019.02.22(금)",
		},
		{
			in:   time.Date(2019, 02, 23, -9, 0, 0, 0, time.UTC),
			want: "2019.02.23(토)",
		},
		{
			in:   time.Date(2019, 02, 23, 0, 0, 0, 0, time.UTC),
			want: "2019.02.23(토)",
		},
		{
			in:   time.Date(2019, 02, 23, 24-9, 0, 0, 0, time.UTC),
			want: "2019.02.24(일)",
		},
		{
			in:   NewKSTDate(2019, 02, 23, 0, 0, 0, 0).Time,
			want: "2019.02.23(토)",
		},
		{
			in:   NewKSTDate(2019, 02, 23, 16, 0, 0, 0).Time,
			want: "2019.02.23(토)",
		},
		{
			in:   NewKSTDate(2019, 02, 23, 0, 0, 0, 0).Time,
			want: "2019.02.23(토)",
		},
		{
			// utc로 이 시간이기에 KST 에서는 9시간 더한 시간이 되기에
			in:   time.Date(2019, 02, 23, 23, 59, 59, 0, time.UTC),
			want: "2019.02.24(일)",
		},
		{
			in: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2019-02-24T00:00:00+09:00")
				return t
			}(),
			want: "2019.02.24(일)",
		},
	}

	for i, c := range cases {
		in, want := c.in, c.want
		got := NewKST(in)

		if want != got.String() {
			t.Errorf("[%d] in: %v, want: %v, got: %v\n", i, in, want, got)
		}
	}
}

func TestCalcThisMonday(t *testing.T) {
	cases := []struct {
		in   KST
		want string
	}{
		{
			in:   NewKSTDate(2019, 02, 18, -10, 0, 0, 0),
			want: "2019.02.11(월)",
		},
		{
			in:   NewKSTDate(2019, 02, 18, -1, 0, 0, 0),
			want: "2019.02.11(월)",
		},
		{
			in:   NewKSTDate(2019, 02, 18, 0, 0, 0, 0),
			want: "2019.02.18(월)",
		},
		{
			in:   NewKSTDate(2019, 02, 24, 0, 0, 0, 0),
			want: "2019.02.18(월)",
		},
		{
			in:   NewKSTDate(2019, 02, 24, 23, 0, 0, 0),
			want: "2019.02.18(월)",
		},
		{
			in:   NewKSTDate(2019, 02, 25, 0, 0, 0, 0),
			want: "2019.02.25(월)",
		},
		{
			in:   NewKSTDate(2019, 02, 25, 8, 0, 0, 0),
			want: "2019.02.25(월)",
		},
		{
			in:   NewKSTDate(2019, 02, 25, 9, 0, 0, 0),
			want: "2019.02.25(월)",
		},
	}

	for i, c := range cases {
		in, want := c.in, c.want
		got := in.CalcThisMonday()
		if want != got.String() {
			t.Errorf("[%d] in: %v, want: %v, got: %v\n", i, in, want, got)
		}
	}
}

func TestKSTString(t *testing.T) {
	cases := []struct {
		in   KST
		want string
	}{
		{
			in:   NewKSTDate(2019, 02, 11, -10, 0, 0, 0),
			want: "2019.02.10(일)",
		},
		{
			in:   NewKSTDate(2019, 02, 11, -1, 0, 0, 0),
			want: "2019.02.10(일)",
		},
		{
			in:   NewKSTDate(2019, 02, 11, 0, 0, 0, 0),
			want: "2019.02.11(월)",
		},
		{
			in:   NewKSTDate(2019, 02, 11, 23, 59, 59, 0),
			want: "2019.02.11(월)",
		},
	}

	for i, c := range cases {
		in, want := c.in, c.want
		got := in.String()
		if got != want {
			t.Errorf("[%d] in: %v, want: %v, got: %v\n", i, in, want, got)
		}
	}
}
