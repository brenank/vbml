package conformance

import (
	"testing"

	vbml "github.com/brenank/vbml/go"
)

func TestCalendar(t *testing.T) {
	runConformanceSuite[fixtureCalendarData, vbml.Board](t, "calendar", func(input fixtureCalendarData) (vbml.Board, error) {
		calendar, err := input.toRuntime()
		if err != nil {
			return nil, err
		}
		return vbml.MakeCalendar(calendar.Month, calendar.Year, calendar), nil
	})
}
