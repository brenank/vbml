package vbml

type lineState struct {
	line   string
	length int
}

func getLinesFromWords(width int, words []string) []string {
	lines := []lineState{{line: "", length: 0}}

	for index, current := range words {
		lastIndex := len(lines) - 1
		lineLength := lines[lastIndex].length
		emptyLine := lines[lastIndex].line == ""

		if current == "\n" && index > 0 && words[index-1] == "\n" {
			lines = append(lines, lineState{line: "", length: 0})
			continue
		}

		if current == "\n" {
			lines[lastIndex] = lineState{line: lines[lastIndex].line, length: width}
			continue
		}

		if isCharacterCodeToken(current) {
			if 1+lineLength > width {
				if current == "{0}" {
					continue
				}
				lines = append(lines, lineState{line: current, length: 1})
				continue
			}

			lines[lastIndex] = lineState{
				line:   lines[lastIndex].line + current,
				length: lineLength + 1,
			}
			continue
		}

		currentLength := stringLength(current)
		if width >= currentLength+lineLength && !emptyLine {
			lines[lastIndex] = lineState{
				line:   lines[lastIndex].line + current,
				length: lineLength + currentLength,
			}
			continue
		}

		lines = append(lines, lineState{line: current, length: currentLength})
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
