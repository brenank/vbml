package vbml

func ParseComponent(defaultHeight, defaultWidth int, props map[string]any, component Component) (Board, error) {
	if component.RawCharacters != nil {
		return component.RawCharacters, nil
	}

	width, height := componentDimensions(defaultHeight, defaultWidth, component.Style)
	if component.RandomColors != nil {
		return randomColors(height, width, component.RandomColors.Colors), nil
	}

	template := component.Template
	justify := JustifyLeft
	align := AlignTop
	if component.Style != nil {
		if component.Style.Justify != "" {
			justify = component.Style.Justify
		}
		if component.Style.Align != "" {
			align = component.Style.Align
		}
	}

	text := emojisToCharacterCodes(template)
	text = parseProps(props, text)
	text = SanitizeSpecialCharacters(text)
	words := splitWords(width, text)
	lines := getLinesFromWords(width, words)

	codes := make(Board, 0, len(lines))
	for _, line := range lines {
		row, err := convertCharactersToCharacterCodes(line)
		if err != nil {
			return nil, err
		}
		codes = append(codes, row)
	}

	codes = verticalAlign(height, align, codes)
	codes = horizontalAlign(width, justify, codes)
	return renderComponent(createEmptyBoard(height, width), codes), nil
}

func ParseAbsoluteComponent(defaultHeight, defaultWidth int, props map[string]any, component Component) (AbsoluteComponent, error) {
	characters, err := ParseComponent(defaultHeight, defaultWidth, props, component)
	if err != nil {
		return AbsoluteComponent{}, err
	}

	absoluteComponent := AbsoluteComponent{
		Characters: characters,
		X:          0,
		Y:          0,
	}
	if component.Style != nil && component.Style.AbsolutePosition != nil {
		absoluteComponent.X = component.Style.AbsolutePosition.X
		absoluteComponent.Y = component.Style.AbsolutePosition.Y
	}
	return absoluteComponent, nil
}

func componentDimensions(defaultHeight, defaultWidth int, style *ComponentStyle) (width, height int) {
	width = defaultWidth
	height = defaultHeight
	if style == nil {
		return width, height
	}
	if style.Width != 0 {
		width = style.Width
	}
	if style.Height != 0 {
		height = style.Height
	}
	return width, height
}
