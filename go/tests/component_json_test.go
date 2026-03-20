package tests

import (
	"encoding/json"
	"reflect"
	"testing"

	vbml "github.com/brenank/vbml/go"
)

func TestComponentJSONWithStringTemplate(t *testing.T) {
	input := []byte(`{"template":"Hello World!"}`)

	var component vbml.Component
	if err := json.Unmarshal(input, &component); err != nil {
		t.Fatalf("unmarshal component: %v", err)
	}

	if component.Template != "Hello World!" {
		t.Fatalf("unexpected template: %#v", component.Template)
	}
	if len(component.TemplateParts) != 0 {
		t.Fatalf("expected no template parts, got %#v", component.TemplateParts)
	}

	output, err := json.Marshal(component)
	if err != nil {
		t.Fatalf("marshal component: %v", err)
	}

	assertJSONEqual(t, input, output)
}

func TestComponentJSONWithTemplateParts(t *testing.T) {
	input := []byte(`{"template":[{"template":"rainfall amounts "},{"template":"{{amount}}","wrap":"never"},{"template":" possible."}]}`)

	var component vbml.Component
	if err := json.Unmarshal(input, &component); err != nil {
		t.Fatalf("unmarshal component: %v", err)
	}

	expectedParts := []vbml.TemplatePart{
		{Template: "rainfall amounts ", Wrap: vbml.TemplateWrapNormal},
		{Template: "{{amount}}", Wrap: vbml.TemplateWrapNever},
		{Template: " possible.", Wrap: vbml.TemplateWrapNormal},
	}
	if !reflect.DeepEqual(component.TemplateParts, expectedParts) {
		t.Fatalf("unexpected template parts:\nexpected: %#v\ngot: %#v", expectedParts, component.TemplateParts)
	}
	if component.Template != "" {
		t.Fatalf("expected empty string template, got %#v", component.Template)
	}

	output, err := json.Marshal(component)
	if err != nil {
		t.Fatalf("marshal component: %v", err)
	}

	assertJSONEqual(t, input, output)
}

func TestParseComponentRejectsMutuallyExclusiveTemplateSources(t *testing.T) {
	_, err := vbml.ParseComponent(1, 5, nil, vbml.Component{
		Template: "A",
		TemplateParts: []vbml.TemplatePart{
			{Template: "B"},
		},
	})
	if err == nil {
		t.Fatal("expected mutually exclusive template sources to fail")
	}
	if err.Error() != "component template and templateParts are mutually exclusive" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertJSONEqual(t *testing.T, expected, actual []byte) {
	t.Helper()

	var expectedValue any
	if err := json.Unmarshal(expected, &expectedValue); err != nil {
		t.Fatalf("decode expected json: %v", err)
	}

	var actualValue any
	if err := json.Unmarshal(actual, &actualValue); err != nil {
		t.Fatalf("decode actual json: %v", err)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Fatalf("unexpected json:\nexpected: %#v\ngot: %#v", expectedValue, actualValue)
	}
}
