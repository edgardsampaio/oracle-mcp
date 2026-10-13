package oracle

import (
	"math"
	"testing"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

func TestNormalizeTypes(t *testing.T) {
	num, err := go_ora.NewNumberFromString("12345678901234567890.5")
	if err != nil {
		t.Fatal(err)
	}
	when := time.Date(2026, 10, 9, 15, 4, 5, 0, time.FixedZone("BRT", -3*3600))
	cases := []struct {
		in   any
		want any
	}{
		{nil, nil},
		{int64(7), int64(7)},
		{float64(42), int64(42)},
		{math.Pi, "3.141592653589793"},
		{when, "2026-10-09T18:04:05Z"},
		{[]byte("texto"), "texto"},
		{go_ora.Clob{String: "clob", Valid: true}, "clob"},
		{go_ora.Clob{}, nil},
		{num, "12345678901234567890.5"},
	}
	for _, tc := range cases {
		got, err := Normalize(tc.in)
		if err != nil {
			t.Fatalf("%#v: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("%#v = %#v, want %#v", tc.in, got, tc.want)
		}
	}
}

func TestCollectLimits(t *testing.T) {
	src := &memRows{
		cols: []string{"id", "nome"},
		data: [][]any{{int64(1), "a"}, {int64(2), "b"}, {int64(3), "c"}},
	}
	got, err := Collect(src, 2, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Rows) != 2 || !got.Truncated {
		t.Fatalf("rows=%d truncated=%v", len(got.Rows), got.Truncated)
	}
	src.i = 0
	got, err = Collect(src, 10, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Rows) == 0 || !got.Truncated {
		t.Fatalf("bytes rows=%d truncated=%v", len(got.Rows), got.Truncated)
	}
}

type memRows struct {
	cols []string
	data [][]any
	i    int
}

func (m *memRows) Columns() ([]string, error) { return m.cols, nil }
func (m *memRows) Next() bool {
	if m.i >= len(m.data) {
		return false
	}
	m.i++
	return true
}
func (m *memRows) Scan(dest ...any) error {
	row := m.data[m.i-1]
	for i := range dest {
		*(dest[i].(*any)) = row[i]
	}
	return nil
}
func (m *memRows) Err() error { return nil }
