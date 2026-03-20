package vbml

import "strings"

func splitWords(width int, template string) []string {
	replaced := strings.ReplaceAll(template, " ", "{0}")
	parts := splitCharacterCodeTokens(replaced)

	var expanded []string
	for _, part := range parts {
		if strings.Contains(part, "\n") {
			expanded = append(expanded, splitOnNewlines(part)...)
			continue
		}
		expanded = append(expanded, part)
	}

	var result []string
	for _, part := range expanded {
		switch {
		case part == "":
			continue
		case isCharacterCodeToken(part):
			result = append(result, part)
		case stringLength(part) > width:
			result = append(result, chunkString(part, width)...)
		default:
			result = append(result, part)
		}
	}

	return result
}

func splitCharacterCodeTokens(template string) []string {
	var parts []string
	var current strings.Builder
	runes := []rune(template)

	for index := 0; index < len(runes); index++ {
		if runes[index] != '{' {
			current.WriteRune(runes[index])
			continue
		}

		if current.Len() > 0 {
			parts = append(parts, current.String())
			current.Reset()
		}

		closeIndex := index + 1
		for closeIndex < len(runes) && runes[closeIndex] != '}' {
			closeIndex++
		}
		if closeIndex >= len(runes) {
			current.WriteRune(runes[index])
			continue
		}

		parts = append(parts, string(runes[index:closeIndex+1]))
		index = closeIndex
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func splitOnNewlines(text string) []string {
	var parts []string
	var current strings.Builder

	for _, character := range text {
		if character == '\n' {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			parts = append(parts, "\n")
			continue
		}
		current.WriteRune(character)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func chunkString(word string, width int) []string {
	runes := []rune(word)
	var chunks []string
	for index := 0; index < len(runes); index += width {
		end := index + width
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[index:end]))
	}
	return chunks
}

func isCharacterCodeToken(word string) bool {
	return strings.HasPrefix(word, "{") && strings.HasSuffix(word, "}")
}

func stringLength(text string) int {
	return len([]rune(text))
}
