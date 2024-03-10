package cli

import (
	"testing"
)

func TestConvertMillisecondToConcertDuration(t *testing.T) {
	type test struct {
		ms   int64
		want string
	}
	m := make(map[string]test)
	m["under an hour"] = test{
		ms:   368618,
		want: "6m 8s",
	}
	m["over an hour"] = test{
		ms:   9601071,
		want: "2h 40m",
	}
	for k, v := range m {
		t.Run(k, func(t *testing.T) {
			got := convertMillisecondToConcertDuration(v.ms)
			if v.want != got {
				t.Errorf("got %s want %s", got, v.want)
			}
		})
	}
}
