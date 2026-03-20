package conformance

import (
	"testing"

	vbml "github.com/Vestaboard/vbml/go"
)

type classicInput struct {
	Text    string              `json:"text"`
	Options vbml.ClassicOptions `json:"options"`
}

func TestClassic(t *testing.T) {
	runConformanceSuite[classicInput, vbml.Board](t, "classic", func(input classicInput) (vbml.Board, error) {
		return vbml.Classic(input.Text, input.Options), nil
	})
}
