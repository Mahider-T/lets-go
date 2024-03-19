package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, expected, actual T) {

	t.Helper()

	if actual != expected {
		t.Errorf("expected %q got %q", expected, actual)
	}
}
