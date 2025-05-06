package weakmap

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestMap(t *testing.T) {
	m := Map[int, string]{
		New: func(key int) (*string, error) {
			value := strconv.Itoa(key)
			return &value, nil
		},
	}

	eg := errgroup.Group{}
	for i := range 1024 {
		i := i % 256
		eg.Go(func() error {
			v, err := m.Load(i)
			if err != nil {
				return fmt.Errorf("error loading key %d: %w", i, err)
			}

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

	m.New = func(key int) (*string, error) {
		return nil, errors.New("error")
	}

	if _, err := m.Load(256); err == nil {
		t.Errorf("expected error, got nil")
	}

	m.New = nil
	v, err := m.Load(0)
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}
