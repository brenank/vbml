package vbml

import (
	"errors"
	"strings"
)

const (
	errMutuallyExclusiveTemplateSources = "component template and templateParts are mutually exclusive"
	templateSpaceCell                   = "{0}"
)

type inlineSegment struct {
	Atomic bool
	Cells  []string
}

type templateTokenKind int

const (
	templateTokenWord templateTokenKind = iota
	templateTokenSpace
	templateTokenNewline
)

type templateToken struct {
	kind     templateTokenKind
	segments []inlineSegment
}

type templateLineState struct {
	line   string
	length int
}

func resolveComponentTemplateLines(width int, props map[string]any, component Component) ([]string, error) {
	parts, err := normalizeComponentTemplate(component)
	if err != nil {
		return nil, err
	}

	return resolveTemplateLines(width, props, parts), nil
}

func resolveTemplateLines(width int, props map[string]any, parts []TemplatePart) []string {
	for index, part := range parts {
		parts[index] = preprocessTemplatePart(props, part)
	}

	return wrapTemplateTokens(width, buildTemplateTokens(parts))
}

func normalizeComponentTemplate(component Component) ([]TemplatePart, error) {
	if component.Template != "" && len(component.TemplateParts) > 0 {
		return nil, errors.New(errMutuallyExclusiveTemplateSources)
	}

	if len(component.TemplateParts) > 0 {
		return normalizeTemplateParts(component.TemplateParts), nil
	}

	return []TemplatePart{{
		Template: component.Template,
		Wrap:     TemplateWrapNormal,
	}}, nil
}

func normalizeTemplateParts(parts []TemplatePart) []TemplatePart {
	if len(parts) == 0 {
		return nil
	}

	normalized := make([]TemplatePart, 0, len(parts))
	for _, part := range parts {
		normalized = append(normalized, TemplatePart{
			Template: part.Template,
			Wrap:     normalizeTemplateWrap(part.Wrap),
		})
	}
	return normalized
}

func normalizeTemplateWrap(wrap TemplateWrap) TemplateWrap {
	if wrap == TemplateWrapNever {
		return TemplateWrapNever
	}

	return TemplateWrapNormal
}

func preprocessTemplatePart(props map[string]any, part TemplatePart) TemplatePart {
	return TemplatePart{
		Template: SanitizeSpecialCharacters(
			parseProps(props, emojisToCharacterCodes(part.Template)),
		),
		Wrap: normalizeTemplateWrap(part.Wrap),
	}
}

func tokenizeDisplayCells(template string) []string {
	runes := []rune(template)
	cells := make([]string, 0, len(runes))

	for index := 0; index < len(runes); index++ {
		current := string(runes[index])
		if current == "{" {
			closeIndex := index + 1
			for closeIndex < len(runes) && string(runes[closeIndex]) != "}" {
				closeIndex++
			}
			if closeIndex < len(runes) {
				cells = append(cells, string(runes[index:closeIndex+1]))
				index = closeIndex
				continue
			}
		}

		if current == " " {
			cells = append(cells, templateSpaceCell)
			continue
		}

		cells = append(cells, current)
	}

	return cells
}

func buildTemplateTokens(parts []TemplatePart) []templateToken {
	var tokens []templateToken
	var currentSegments []inlineSegment

	pushWordToken := func() {
		if len(currentSegments) == 0 {
			return
		}
		tokens = append(tokens, createWordToken(currentSegments))
		currentSegments = nil
	}

	pushSegment := func(atomic bool, cells []string) {
		if len(cells) == 0 {
			return
		}
		currentSegments = append(currentSegments, inlineSegment{
			Atomic: atomic,
			Cells:  cloneCells(cells),
		})
	}

	appendNormalCells := func(cells []string) {
		var currentCells []string

		flushCurrentCells := func() {
			pushSegment(false, currentCells)
			currentCells = nil
		}

		for _, cell := range cells {
			switch cell {
			case "\n":
				flushCurrentCells()
				pushWordToken()
				tokens = append(tokens, templateToken{kind: templateTokenNewline})
			case templateSpaceCell:
				flushCurrentCells()
				pushWordToken()
				tokens = append(tokens, templateToken{kind: templateTokenSpace})
			default:
				currentCells = append(currentCells, cell)
			}
		}

		flushCurrentCells()
	}

	appendAtomicCells := func(cells []string) {
		var currentCells []string

		flushCurrentCells := func() {
			pushSegment(true, currentCells)
			currentCells = nil
		}

		for _, cell := range cells {
			if cell == "\n" {
				flushCurrentCells()
				pushWordToken()
				tokens = append(tokens, templateToken{kind: templateTokenNewline})
				continue
			}

			currentCells = append(currentCells, cell)
		}

		flushCurrentCells()
	}

	for _, part := range parts {
		cells := tokenizeDisplayCells(part.Template)
		if normalizeTemplateWrap(part.Wrap) == TemplateWrapNever {
			appendAtomicCells(cells)
			continue
		}
		appendNormalCells(cells)
	}

	pushWordToken()

	return tokens
}

func wrapTemplateTokens(width int, tokens []templateToken) []string {
	lines := []templateLineState{{
		line:   "",
		length: 0,
	}}
	chunkedTokens := chunkWordTokens(width, tokens)

	for index, token := range chunkedTokens {
		lastIndex := len(lines) - 1
		lineLength := lines[lastIndex].length
		emptyLine := lines[lastIndex].line == ""

		if token.kind == templateTokenNewline && index > 0 && chunkedTokens[index-1].kind == templateTokenNewline {
			lines = append(lines, templateLineState{})
			continue
		}

		if token.kind == templateTokenNewline {
			lines[lastIndex] = templateLineState{
				line:   lines[lastIndex].line,
				length: width,
			}
			continue
		}

		if token.kind == templateTokenSpace {
			if lineLength+1 > width {
				continue
			}

			lines[lastIndex] = templateLineState{
				line:   lines[lastIndex].line + templateSpaceCell,
				length: lineLength + 1,
			}
			continue
		}

		wordWidth := measureWordWidth(token)
		if width >= wordWidth+lineLength && !emptyLine {
			lines[lastIndex] = templateLineState{
				line:   lines[lastIndex].line + wordTokenToString(token),
				length: lineLength + wordWidth,
			}
			continue
		}

		lines = append(lines, templateLineState{
			line:   wordTokenToString(token),
			length: wordWidth,
		})
	}

	startIndex := 0
	if lines[0].line == "" {
		startIndex = 1
	}

	result := make([]string, 0, len(lines)-startIndex)
	for _, line := range lines[startIndex:] {
		result = append(result, line.line)
	}
	return result
}

func chunkWordTokens(width int, tokens []templateToken) []templateToken {
	var chunkedTokens []templateToken

	for _, token := range tokens {
		if token.kind != templateTokenWord || measureWordWidth(token) <= width {
			chunkedTokens = append(chunkedTokens, token)
			continue
		}

		remaining := &token
		for remaining != nil && measureWordWidth(*remaining) > width {
			prefix, remainder := takeWordPrefix(*remaining, width)
			if prefix == nil {
				break
			}

			chunkedTokens = append(chunkedTokens, *prefix)
			remaining = remainder
		}

		if remaining != nil {
			chunkedTokens = append(chunkedTokens, *remaining)
		}
	}

	return chunkedTokens
}

func takeWordPrefix(token templateToken, maxWidth int) (*templateToken, *templateToken) {
	if maxWidth <= 0 {
		remainder := createWordToken(cloneSegments(token.segments))
		return nil, &remainder
	}

	var prefixSegments []inlineSegment
	var remainderSegments []inlineSegment
	consumedWidth := 0

	appendRemainder := func(startIndex int) {
		remainderSegments = append(remainderSegments, cloneSegments(token.segments[startIndex:])...)
	}

	for index, segment := range token.segments {
		segmentWidth := len(segment.Cells)
		if consumedWidth+segmentWidth <= maxWidth {
			prefixSegments = append(prefixSegments, cloneSegment(segment))
			consumedWidth += segmentWidth
			continue
		}

		if segment.Atomic {
			if consumedWidth == 0 {
				prefixSegments = append(prefixSegments, inlineSegment{
					Atomic: true,
					Cells:  cloneCells(segment.Cells[:maxWidth]),
				})

				if segmentWidth > maxWidth {
					remainderSegments = append(remainderSegments, inlineSegment{
						Atomic: true,
						Cells:  cloneCells(segment.Cells[maxWidth:]),
					})
				} else {
					remainderSegments = append(remainderSegments, cloneSegment(segment))
				}

				appendRemainder(index + 1)
			} else {
				remainderSegments = append(remainderSegments, cloneSegment(segment))
				appendRemainder(index + 1)
			}

			return newWordPrefix(prefixSegments), newWordPrefix(remainderSegments)
		}

		remainingWidth := maxWidth - consumedWidth
		if remainingWidth > 0 {
			prefixSegments = append(prefixSegments, inlineSegment{
				Atomic: false,
				Cells:  cloneCells(segment.Cells[:remainingWidth]),
			})
			remainderSegments = append(remainderSegments, inlineSegment{
				Atomic: false,
				Cells:  cloneCells(segment.Cells[remainingWidth:]),
			})
		} else {
			remainderSegments = append(remainderSegments, cloneSegment(segment))
		}

		appendRemainder(index + 1)
		return newWordPrefix(prefixSegments), newWordPrefix(remainderSegments)
	}

	prefix := createWordToken(prefixSegments)
	return &prefix, nil
}

func newWordPrefix(segments []inlineSegment) *templateToken {
	if len(segments) == 0 {
		return nil
	}

	token := createWordToken(segments)
	return &token
}

func createWordToken(segments []inlineSegment) templateToken {
	return templateToken{
		kind:     templateTokenWord,
		segments: cloneSegments(segments),
	}
}

func cloneSegments(segments []inlineSegment) []inlineSegment {
	cloned := make([]inlineSegment, 0, len(segments))
	for _, segment := range segments {
		cloned = append(cloned, cloneSegment(segment))
	}
	return cloned
}

func cloneSegment(segment inlineSegment) inlineSegment {
	return inlineSegment{
		Atomic: segment.Atomic,
		Cells:  cloneCells(segment.Cells),
	}
}

func cloneCells(cells []string) []string {
	return append([]string(nil), cells...)
}

func measureWordWidth(token templateToken) int {
	width := 0
	for _, segment := range token.segments {
		width += len(segment.Cells)
	}
	return width
}

func wordTokenToString(token templateToken) string {
	var builder strings.Builder
	for _, segment := range token.segments {
		for _, cell := range segment.Cells {
			builder.WriteString(cell)
		}
	}
	return builder.String()
}
