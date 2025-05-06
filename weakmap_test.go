package weakmap

import (
	"fmt"
	"strconv"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestMap(t *testing.T) {
	m := Map[int, string]{
		New: func(key int) *string {
			value := strconv.Itoa(key)
			return &value
		},
	}

	eg := errgroup.Group{}
	for i := range 1024 {
		eg.Go(func() error {
			v := m.Load(i)
			if v == nil {
				return fmt.Errorf("nil value for key %d", i)
			}

			if expected := strconv.Itoa(i); *v != expected {
				return fmt.Errorf("invalid value for key %d: expected %s, got %s", i, expected, *v)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Errorf("error: %v", err)
	}
}
