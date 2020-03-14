package store

import "testing"

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}
