package vbml

import (
	"encoding/json"
	"errors"
	"fmt"
)

type componentJSON struct {
	Template      any               `json:"template,omitempty"`
	RawCharacters Board             `json:"rawCharacters,omitempty"`
	Calendar      *CalendarData     `json:"calendar,omitempty"`
	RandomColors  *RandomColorsData `json:"randomColors,omitempty"`
	Style         *ComponentStyle   `json:"style,omitempty"`
}

type componentJSONInput struct {
	Template      json.RawMessage   `json:"template"`
	RawCharacters Board             `json:"rawCharacters"`
	Calendar      *CalendarData     `json:"calendar"`
	RandomColors  *RandomColorsData `json:"randomColors"`
	Style         *ComponentStyle   `json:"style"`
}

type templatePartJSON struct {
	Template string       `json:"template"`
	Wrap     TemplateWrap `json:"wrap,omitempty"`
}

func (part TemplatePart) MarshalJSON() ([]byte, error) {
	payload := templatePartJSON{
		Template: part.Template,
	}
	if normalizeTemplateWrap(part.Wrap) == TemplateWrapNever {
		payload.Wrap = TemplateWrapNever
	}

	return json.Marshal(payload)
}

func (part *TemplatePart) UnmarshalJSON(data []byte) error {
	var payload templatePartJSON
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	part.Template = payload.Template
	part.Wrap = normalizeTemplateWrap(payload.Wrap)
	return nil
}

func (component Component) MarshalJSON() ([]byte, error) {
	if component.Template != "" && len(component.TemplateParts) > 0 {
		return nil, errors.New(errMutuallyExclusiveTemplateSources)
	}

	payload := componentJSON{
		RawCharacters: component.RawCharacters,
		Calendar:      component.Calendar,
		RandomColors:  component.RandomColors,
		Style:         component.Style,
	}

	if len(component.TemplateParts) > 0 {
		payload.Template = normalizeTemplateParts(component.TemplateParts)
	} else if component.Template != "" {
		payload.Template = component.Template
	}

	return json.Marshal(payload)
}

func (component *Component) UnmarshalJSON(data []byte) error {
	var payload componentJSONInput
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	*component = Component{
		RawCharacters: payload.RawCharacters,
		Calendar:      payload.Calendar,
		RandomColors:  payload.RandomColors,
		Style:         payload.Style,
	}

	if len(payload.Template) == 0 || string(payload.Template) == "null" {
		return nil
	}

	var template string
	if err := json.Unmarshal(payload.Template, &template); err == nil {
		component.Template = template
		return nil
	}

	var parts []TemplatePart
	if err := json.Unmarshal(payload.Template, &parts); err == nil {
		component.TemplateParts = normalizeTemplateParts(parts)
		return nil
	}

	return fmt.Errorf("invalid template: expected string or array")
}
