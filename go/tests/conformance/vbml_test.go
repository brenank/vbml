package conformance

import (
	"testing"

	vbml "github.com/Vestaboard/vbml/go"
)

func TestVBML(t *testing.T) {
	runConformanceSuite[fixtureInput, vbml.Board](t, "vbml", func(input fixtureInput) (vbml.Board, error) {
		runtimeInput, err := input.toRuntime()
		if err != nil {
			return nil, err
		}
		return vbml.Parse(runtimeInput)
	})
}
