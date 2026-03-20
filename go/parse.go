package vbml

func Parse(input Input) (Board, error) {
	height, width := inputDimensions(input.Style)
	props := input.Props
	if props == nil {
		props = map[string]any{}
	}

	board := createEmptyBoard(height, width)

	var flowComponents []Board
	var absoluteComponents []AbsoluteComponent
	var calendarComponents []calendarComponent

	for _, component := range input.Components {
		if component.Calendar != nil {
			x := 0
			if component.Style != nil && component.Style.AbsolutePosition != nil {
				x = component.Style.AbsolutePosition.X
			}
			calendarComponents = append(calendarComponents, calendarComponent{
				characters: MakeCalendar(component.Calendar.Month, component.Calendar.Year, *component.Calendar),
				x:          x,
			})
			continue
		}

		if component.Style != nil && component.Style.AbsolutePosition != nil {
			absoluteComponent, err := ParseAbsoluteComponent(height, width, props, component)
			if err != nil {
				return nil, err
			}
			absoluteComponents = append(absoluteComponents, absoluteComponent)
			continue
		}

		parsedComponent, err := ParseComponent(height, width, props, component)
		if err != nil {
			return nil, err
		}
		flowComponents = append(flowComponents, parsedComponent)
	}

	return layoutComponents(board, flowComponents, absoluteComponents, calendarComponents), nil
}

func inputDimensions(style *BoardStyle) (height, width int) {
	height = flagshipBoardHeight
	width = flagshipBoardWidth
	if style == nil {
		return height, width
	}
	if style.Height != 0 {
		height = style.Height
	}
	if style.Width != 0 {
		width = style.Width
	}
	return height, width
}
