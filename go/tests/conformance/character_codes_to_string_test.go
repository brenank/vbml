package conformance

import (
	"testing"

	vbml "github.com/brenank/vbml/go"
)

type characterCodesToStringInput struct {
	Characters vbml.Board                         `json:"characters"`
	Options    vbml.CharacterCodesToStringOptions `json:"options"`
}

func TestCharacterCodesToString(t *testing.T) {
	runConformanceSuite[characterCodesToStringInput, string](t, "characterCodesToString", func(input characterCodesToStringInput) (string, error) {
		return vbml.CharacterCodesToString(input.Characters, input.Options), nil
	})
}
