package tests

import (
	"testing"

	vbml "github.com/brenank/vbml/go"
)

func TestCopyCharacterCodes(t *testing.T) {
	original := vbml.Board{{1, 2}}
	copied := vbml.CopyCharacterCodes(original)

	if len(copied) != 1 || len(copied[0]) != 2 {
		t.Fatalf("unexpected copy shape: %#v", copied)
	}
	if copied[0][0] != 1 || copied[0][1] != 2 {
		t.Fatalf("unexpected copy contents: %#v", copied)
	}

	original[0][0] = 3
	if copied[0][0] != 1 {
		t.Fatalf("copy should not share backing data, got %d", copied[0][0])
	}
}
