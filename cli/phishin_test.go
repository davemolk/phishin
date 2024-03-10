package cli

import (
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("must supply an argument", func(t *testing.T) {
		i := Run(nil)
		if i != 1 {
			t.Errorf("got %d want 1", i)
		}
	})
	t.Run("must supply PHISHIN_API_KEY as environment var", func(t *testing.T) {
		i := Run([]string{"eras"})
		if i != 1 {
			t.Errorf("got %d want 1", i)
		}
	})
}
