package conformance

import (
	"strconv"

	vbml "github.com/Vestaboard/vbml/go"
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

type fixtureComponent struct {
	Template      string                   `json:"template"`
	RawCharacters vbml.Board               `json:"rawCharacters"`
	Calendar      *fixtureCalendarData     `json:"calendar"`
	RandomColors  *fixtureRandomColorsData `json:"randomColors"`
	Style         *fixtureComponentStyle   `json:"style"`
}

func (fixture fixtureComponent) toRuntime() (vbml.Component, error) {
	component := vbml.Component{
		Template:      fixture.Template,
		RawCharacters: fixture.RawCharacters,
		RandomColors:  nil,
		Style:         nil,
	}

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
		month, err := strconv.Atoi(fixture.Calendar.Month)
		if err != nil {
			return vbml.Component{}, err
		}
		year, err := strconv.Atoi(fixture.Calendar.Year)
		if err != nil {
			return vbml.Component{}, err
		}

		days := make(map[int]int, len(fixture.Calendar.Days))
		for key, value := range fixture.Calendar.Days {
			day, err := strconv.Atoi(key)
			if err != nil {
				return vbml.Component{}, err
			}
			days[day] = value
		}

		component.Calendar = &vbml.CalendarData{
			Month:           month,
			Year:            year,
			DefaultDayColor: fixture.Calendar.DefaultDayColor,
			Days:            days,
			HideSMTWTFS:     fixture.Calendar.HideSMTWTFS,
			HideDates:       fixture.Calendar.HideDates,
			HideMonthYear:   fixture.Calendar.HideMonthYear,
		}
	}

	return component, nil
}
