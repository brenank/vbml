package conformance

import (
	"testing"

	vbml "github.com/Vestaboard/vbml/go"
)

type textInput struct {
	Text string `json:"text"`
}

func TestHasSpecialCharacters(t *testing.T) {
	runConformanceSuite[textInput, bool](t, "hasSpecialCharacters", func(input textInput) (bool, error) {
		return vbml.HasSpecialCharacters(input.Text), nil
	})
}
