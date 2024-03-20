package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {

	t.Helper()

	if actual != expected {
		t.Errorf("expected %v got %v", expected, actual)
	}
}

func StringContains(t *testing.T, actual, expected string) {

	t.Helper()

	if !strings.Contains(actual, expected) {
		t.Errorf("Got %q, expected to contain : %q", actual, expected)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
