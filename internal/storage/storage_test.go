package storage

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	if id1 < 1000000000000000 {
		t.Logf("Generated ID does not have enough length")
	}
	id2 := GenerateID()
	if id2 < 10e15 {
		t.Logf("Generated ID does not have enough length")
	}
	if id1 == id2 {
		t.Logf("Generated IDs must not be equal")
	}
}
