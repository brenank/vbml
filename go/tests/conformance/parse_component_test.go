package conformance

import (
	"reflect"
	"testing"

	vbml "github.com/Vestaboard/vbml/go"
)

type parseComponentInput struct {
	Mode      string           `json:"mode"`
	Height    int              `json:"height"`
	Width     int              `json:"width"`
	Props     map[string]any   `json:"props"`
	Component fixtureComponent `json:"component"`
}

type absoluteComponentResult struct {
	Characters vbml.Board `json:"characters"`
	X          int        `json:"x"`
	Y          int        `json:"y"`
}

func TestParseComponent(t *testing.T) {
	for _, testCase := range loadConformanceCases(t, "parseComponent") {
		t.Run(testCase.ID, func(t *testing.T) {
			if testCase.Expected.Skip != "" {
				t.Skip(testCase.Expected.Skip)
			}

			input := decodeRawJSON[parseComponentInput](t, testCase.Input)
			component, err := input.Component.toRuntime()
			if err != nil {
				t.Fatalf("convert fixture component: %v", err)
			}

			if testCase.Expected.Error != "" {
				runParseComponentExpectingError(t, input, component, testCase.Expected.Error)
				return
			}

			if input.Mode == "absolute" {
				result, err := vbml.ParseAbsoluteComponent(input.Height, input.Width, input.Props, component)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				expected := decodeRawJSON[absoluteComponentResult](t, testCase.Expected.Result)
				actual := absoluteComponentResult{
					Characters: result.Characters,
					X:          result.X,
					Y:          result.Y,
				}
				if !reflect.DeepEqual(actual, expected) {
					t.Fatalf("unexpected result:\nexpected: %#v\ngot: %#v", expected, actual)
				}
				return
			}

			result, err := vbml.ParseComponent(input.Height, input.Width, input.Props, component)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expected := decodeRawJSON[vbml.Board](t, testCase.Expected.Result)
			if !reflect.DeepEqual(result, expected) {
				t.Fatalf("unexpected result:\nexpected: %#v\ngot: %#v", expected, result)
			}
		})
	}
}

func runParseComponentExpectingError(
	t *testing.T,
	input parseComponentInput,
	component vbml.Component,
	expectedError string,
) {
	t.Helper()

	if input.Mode == "absolute" {
		_, err := vbml.ParseAbsoluteComponent(input.Height, input.Width, input.Props, component)
		if err == nil {
			t.Fatalf("expected error %q, got nil", expectedError)
		}
		if err.Error() != expectedError {
			t.Fatalf("expected error %q, got %q", expectedError, err.Error())
		}
		return
	}

	_, err := vbml.ParseComponent(input.Height, input.Width, input.Props, component)
	if err == nil {
		t.Fatalf("expected error %q, got nil", expectedError)
	}
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}
