package vbml

import "strings"

func SanitizeSpecialCharacters(text string) string {
	var builder strings.Builder
	for _, character := range text {
		builder.WriteString(mappingToCharacter(string(character)))
	}
	return builder.String()
}
