package conformance

import (
	"testing"

	vbml "github.com/Vestaboard/vbml/go"
)

type characterCodesToASCIIInput struct {
	CharacterCodes vbml.Board `json:"characterCodes"`
	IsWhite        bool       `json:"isWhite"`
}

func TestCharacterCodesToASCII(t *testing.T) {
	runConformanceSuite[characterCodesToASCIIInput, string](t, "characterCodesToAscii", func(input characterCodesToASCIIInput) (string, error) {
		return vbml.CharacterCodesToASCII(input.CharacterCodes, input.IsWhite), nil
	})
}
