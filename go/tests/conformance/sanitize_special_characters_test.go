package conformance

import (
	"testing"

	vbml "github.com/brenank/vbml/go"
)

func TestSanitizeSpecialCharacters(t *testing.T) {
	runConformanceSuite[textInput, string](t, "sanitizeSpecialCharacters", func(input textInput) (string, error) {
		return vbml.SanitizeSpecialCharacters(input.Text), nil
	})
}
