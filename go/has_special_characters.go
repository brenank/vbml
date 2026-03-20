package vbml

func HasSpecialCharacters(text string) bool {
	if text == "" {
		return false
	}

	for _, character := range text {
		if _, ok := supportedCharacters[string(character)]; !ok {
			return true
		}
	}

	return false
}
