package scheduler

import (
	"errors"
	"testing"
)

func TestFuncJob(t *testing.T) {
	testingError := errors.New("test error")

	fj := FuncJob(func() error {
		return testingError
	})

	if err := fj.Run(); err != testingError {
		t.Errorf("expected testingError, got: %v", err)
	}

	fj = FuncJob(func() error {
		return nil
	})

	if err := fj.Run(); err != nil {
		t.Errorf("expected nil, got: %v", err)
	}
}
