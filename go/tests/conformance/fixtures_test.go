package conformance

import (
	"encoding/json"
	"strconv"

	vbml "github.com/brenank/vbml/go"
)

type fixtureAbsolutePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type fixtureComponentStyle struct {
	Justify          string                   `json:"justify"`
	Align            string                   `json:"align"`
	Height           int                      `json:"height"`
	Width            int                      `json:"width"`
	AbsolutePosition *fixtureAbsolutePosition `json:"absolutePosition"`
}

type fixtureCalendarData struct {
	Month           string         `json:"month"`
	Year            string         `json:"year"`
	DefaultDayColor int            `json:"defaultDayColor"`
	Days            map[string]int `json:"days"`
	HideSMTWTFS     bool           `json:"hideSMTWTFS"`
	HideDates       bool           `json:"hideDates"`
	HideMonthYear   bool           `json:"hideMonthYear"`
}

type fixtureRandomColorsData struct {
	Colors []int `json:"colors"`
}

type fixtureBoardStyle struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type fixtureTemplatePart struct {
	Template string `json:"template"`
	Wrap     string `json:"wrap"`
}

type fixtureComponent struct {
	Template      json.RawMessage          `json:"template"`
	RawCharacters vbml.Board               `json:"rawCharacters"`
	Calendar      *fixtureCalendarData     `json:"calendar"`
	RandomColors  *fixtureRandomColorsData `json:"randomColors"`
	Style         *fixtureComponentStyle   `json:"style"`
}

type fixtureInput struct {
	Props      map[string]any     `json:"props"`
	Style      *fixtureBoardStyle `json:"style"`
	Components []fixtureComponent `json:"components"`
}

func (fixture fixtureCalendarData) toRuntime() (vbml.CalendarData, error) {
	month, err := strconv.Atoi(fixture.Month)
	if err != nil {
		return vbml.CalendarData{}, err
	}
	year, err := strconv.Atoi(fixture.Year)
	if err != nil {
		return vbml.CalendarData{}, err
	}

	days := make(map[int]int, len(fixture.Days))
	for key, value := range fixture.Days {
		day, err := strconv.Atoi(key)
		if err != nil {
			return vbml.CalendarData{}, err
		}
		days[day] = value
	}

	return vbml.CalendarData{
		Month:           month,
		Year:            year,
		DefaultDayColor: fixture.DefaultDayColor,
		Days:            days,
		HideSMTWTFS:     fixture.HideSMTWTFS,
		HideDates:       fixture.HideDates,
		HideMonthYear:   fixture.HideMonthYear,
	}, nil
}

func (fixture fixtureComponent) toRuntime() (vbml.Component, error) {
	component := vbml.Component{
		RawCharacters: fixture.RawCharacters,
		RandomColors:  nil,
		Style:         nil,
	}

	template, templateParts, err := fixtureTemplateToRuntime(fixture.Template)
	if err != nil {
		return vbml.Component{}, err
	}
	component.Template = template
	component.TemplateParts = templateParts

	if fixture.RandomColors != nil {
		component.RandomColors = &vbml.RandomColorsData{Colors: fixture.RandomColors.Colors}
	}
	if fixture.Style != nil {
		style := &vbml.ComponentStyle{
			Justify: vbml.Justify(fixture.Style.Justify),
			Align:   vbml.Align(fixture.Style.Align),
			Height:  fixture.Style.Height,
			Width:   fixture.Style.Width,
		}
		if fixture.Style.AbsolutePosition != nil {
			style.AbsolutePosition = &vbml.AbsolutePosition{
				X: fixture.Style.AbsolutePosition.X,
				Y: fixture.Style.AbsolutePosition.Y,
			}
		}
		component.Style = style
	}
	if fixture.Calendar != nil {
		calendar, err := fixture.Calendar.toRuntime()
		if err != nil {
			return vbml.Component{}, err
		}
		component.Calendar = &calendar
	}

	return component, nil
}

func fixtureTemplateToRuntime(raw json.RawMessage) (string, []vbml.TemplatePart, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", nil, nil
	}

	var template string
	if err := json.Unmarshal(raw, &template); err == nil {
		return template, nil, nil
	}

	var fixtureParts []fixtureTemplatePart
	if err := json.Unmarshal(raw, &fixtureParts); err == nil {
		parts := make([]vbml.TemplatePart, 0, len(fixtureParts))
		for _, part := range fixtureParts {
			parts = append(parts, vbml.TemplatePart{
				Template: part.Template,
				Wrap:     vbml.TemplateWrap(part.Wrap),
			})
		}
		return "", parts, nil
	}

	return "", nil, json.Unmarshal(raw, &template)
}

func (fixture fixtureInput) toRuntime() (vbml.Input, error) {
	input := vbml.Input{
		Props: fixture.Props,
	}
	if fixture.Style != nil {
		input.Style = &vbml.BoardStyle{
			Height: fixture.Style.Height,
			Width:  fixture.Style.Width,
		}
	}

	input.Components = make([]vbml.Component, 0, len(fixture.Components))
	for _, componentFixture := range fixture.Components {
		component, err := componentFixture.toRuntime()
		if err != nil {
			return vbml.Input{}, err
		}
		input.Components = append(input.Components, component)
	}

	return input, nil
}
