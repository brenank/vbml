package vbml

import (
	"fmt"
	"strconv"
	"strings"
)

const variationSelectorSixteen = "\uFE0F"

type characterCodeEntry struct {
	code     CharacterCode
	mappings []string
}

var (
	colorCodes = []int{
		int(CharacterCodeRed),
		int(CharacterCodeOrange),
		int(CharacterCodeYellow),
		int(CharacterCodeGreen),
		int(CharacterCodeBlue),
		int(CharacterCodeViolet),
		int(CharacterCodeWhite),
		int(CharacterCodeBlack),
	}
	characterCodeTable  = buildCharacterCodeTable()
	validCharacterCodes = buildValidCharacterCodes(characterCodeTable)
	supportedCharacters = buildSupportedCharacters(characterCodeTable)
	mappedCharacters    = buildMappedCharacters(characterCodeTable)
)

func buildCharacterCodeTable() []characterCodeEntry {
	return []characterCodeEntry{
		{code: CharacterCodeBlank, mappings: []string{" ", "©", "®", "<", ">", "²", "†", "‡", "ˆ", "Þ", "þ", "µ", "¶", "*", "^", "¬", "«", "»", "›", "³", "¹", "€", "‹", "˜", "÷", "π", "∆", "√", "∫", "∞"}},
		{code: CharacterCodeLetterA, mappings: []string{"A", "a", "â", "à", "å", "á", "À", "Á", "Â", "Ã", "Å", "ã", "ä", "Ä", "∂", "œ", "æ", "Æ"}},
		{code: CharacterCodeLetterB, mappings: []string{"B", "b"}},
		{code: CharacterCodeLetterC, mappings: []string{"C", "c", "ç", "Ç", "¢", "ć", "Ć", "č", "Č"}},
		{code: CharacterCodeLetterD, mappings: []string{"D", "d", "Ð", "ð"}},
		{code: CharacterCodeLetterE, mappings: []string{"E", "e", "é", "ê", "ë", "è", "È", "É", "Ê", "Ë", "€", "£", "∑"}},
		{code: CharacterCodeLetterF, mappings: []string{"F", "f", "ƒ", "ſ"}},
		{code: CharacterCodeLetterG, mappings: []string{"G", "g", "ğ", "Ğ", "ģ", "Ģ", "ġ", "Ġ", "ĝ", "Ĝ"}},
		{code: CharacterCodeLetterH, mappings: []string{"H", "h", "ħ", "Ħ", "ĥ", "Ĥ"}},
		{code: CharacterCodeLetterI, mappings: []string{"I", "i", "í", "ï", "î", "ì", "Ì", "Í", "Î", "Ï", "|", "¡"}},
		{code: CharacterCodeLetterJ, mappings: []string{"J", "j", "ĵ", "Ĵ", "į", "Į"}},
		{code: CharacterCodeLetterK, mappings: []string{"K", "k", "ķ", "Ķ", "ĸ"}},
		{code: CharacterCodeLetterL, mappings: []string{"L", "l", "£", "ł", "Ł", "ļ", "Ļ", "ĺ", "Ĺ", "ľ", "Ľ", "ŀ", "Ŀ"}},
		{code: CharacterCodeLetterM, mappings: []string{"M", "m"}},
		{code: CharacterCodeLetterN, mappings: []string{"N", "n", "ñ", "Ñ", "ń", "Ń", "ň", "Ň", "ņ", "Ņ"}},
		{code: CharacterCodeLetterO, mappings: []string{"O", "o", "ó", "ô", "ò", "Ò", "Ó", "Ô", "Õ", "Ø", "ð", "õ", "ø", "ö", "Ö"}},
		{code: CharacterCodeLetterP, mappings: []string{"P", "p", "Þ", "þ", "¶"}},
		{code: CharacterCodeLetterQ, mappings: []string{"Q", "q"}},
		{code: CharacterCodeLetterR, mappings: []string{"R", "r", "ŕ", "Ŕ", "ř", "Ř", "ŗ", "Ŗ"}},
		{code: CharacterCodeLetterS, mappings: []string{"S", "s", "š", "Š", "§", "ś", "Ś", "ş", "Ş", "ș", "Ș"}},
		{code: CharacterCodeLetterT, mappings: []string{"T", "t", "ť", "Ť", "ţ", "Ţ", "ŧ", "Ŧ"}},
		{code: CharacterCodeLetterU, mappings: []string{"U", "u", "û", "ù", "ú", "Ù", "Ú", "Û", "µ", "ū", "Ū", "ů", "Ů", "ų", "Ų", "Ü"}},
		{code: CharacterCodeLetterV, mappings: []string{"V", "v", "Ʋ", "ʋ"}},
		{code: CharacterCodeLetterW, mappings: []string{"W", "w", "ŵ", "Ŵ", "ẁ", "Ẁ", "ẃ", "Ẃ", "ẅ", "Ẅ"}},
		{code: CharacterCodeLetterX, mappings: []string{"X", "x", "ẍ", "Ẍ"}},
		{code: CharacterCodeLetterY, mappings: []string{"Y", "y", "ý", "ÿ", "Ý", "ŷ", "Ŷ", "ỳ", "Ỳ", "ỹ", "Ỹ", "Ÿ"}},
		{code: CharacterCodeLetterZ, mappings: []string{"Z", "z", "ž", "Ž", "ź", "Ź", "ż", "Ż"}},
		{code: CharacterCodeOne, mappings: []string{"1", "¹"}},
		{code: CharacterCodeTwo, mappings: []string{"2", "²"}},
		{code: CharacterCodeThree, mappings: []string{"3", "³"}},
		{code: CharacterCodeFour, mappings: []string{"4"}},
		{code: CharacterCodeFive, mappings: []string{"5"}},
		{code: CharacterCodeSix, mappings: []string{"6"}},
		{code: CharacterCodeSeven, mappings: []string{"7"}},
		{code: CharacterCodeEight, mappings: []string{"8"}},
		{code: CharacterCodeNine, mappings: []string{"9"}},
		{code: CharacterCodeZero, mappings: []string{"0"}},
		{code: CharacterCodeExclamationMark, mappings: []string{"!", "ǃ"}},
		{code: CharacterCodeAtSign, mappings: []string{"@"}},
		{code: CharacterCodePoundSign, mappings: []string{"#", "№"}},
		{code: CharacterCodeDollarSign, mappings: []string{"$", "¢", "£", "¤", "¥", "₩", "₪", "₫", "€", "₹", "₺", "₽"}},
		{code: CharacterCodeLeftParenthesis, mappings: []string{"(", "[", "{", "⟨", "«"}},
		{code: CharacterCodeRightParenthesis, mappings: []string{")", "]", "}", "⟩", "»"}},
		{code: CharacterCodeHyphen, mappings: []string{"-", "—", "–", "¯", "~", "_"}},
		{code: CharacterCodePlusSign, mappings: []string{"+", "±", "∓", "∔"}},
		{code: CharacterCodeAmpersand, mappings: []string{"&"}},
		{code: CharacterCodeEqualsSign, mappings: []string{"=", "≠", "≈", "≡"}},
		{code: CharacterCodeSemicolon, mappings: []string{";", "；"}},
		{code: CharacterCodeColon, mappings: []string{":", "¦"}},
		{code: CharacterCodeSingleQuote, mappings: []string{"'", "‘", "’", "`", "´", "‚", "‛", "ʹ", "ʻ", "ʽ", "ʾ", "ʿ", "ˈ", "ˊ", "ˋ"}},
		{code: CharacterCodeDoubleQuote, mappings: []string{`"`, "„", "\u201c", "\u201d", "¨", "˝", "ˮ", "˵", "˶", "‟", "\u201f"}},
		{code: CharacterCodePercentSign, mappings: []string{"%", "‰", "‱"}},
		{code: CharacterCodeComma, mappings: []string{",", "¸", "‚", "，", "、", "､"}},
		{code: CharacterCodePeriod, mappings: []string{".", "․", "‥", "…"}},
		{code: CharacterCodeSlash, mappings: []string{"/", "\\", "⁄", "∕", "⧸", "⫻", "⫽", "⧵"}},
		{code: CharacterCodeQuestionMark, mappings: []string{"?", "¿"}},
		{code: CharacterCodeDegreeSign, mappings: []string{"°", "˚", "º", "¤", "•", "·", "∙", "∘", "⚬", "⦿", "⨀", "⨁", "⨂", "❤️", "🧡", "💛", "💚", "💙", "💜", "🖤", "🤍", "🤎", "❤"}},
	}
}

func buildValidCharacterCodes(entries []characterCodeEntry) map[int]struct{} {
	valid := make(map[int]struct{}, len(entries))
	for _, entry := range entries {
		valid[int(entry.code)] = struct{}{}
	}
	return valid
}

func buildSupportedCharacters(entries []characterCodeEntry) map[string]struct{} {
	supported := make(map[string]struct{}, len(entries)*2)
	for _, entry := range entries {
		if len(entry.mappings) == 0 {
			continue
		}

		primary := entry.mappings[0]
		supported[strings.ToLower(primary)] = struct{}{}
		supported[strings.ToUpper(primary)] = struct{}{}
	}

	for _, character := range []string{
		"\n",
		"\u201c",
		"\u2018",
		"{",
		"}",
		"⬜",
		"🟥",
		"🟧",
		"🟨",
		"🟩",
		"🟦",
		"🟪",
		"⬛",
		"❤",
	} {
		supported[character] = struct{}{}
	}

	return supported
}

func buildMappedCharacters(entries []characterCodeEntry) map[string]int {
	mapped := make(map[string]int)
	for _, entry := range entries {
		for _, mapping := range entry.mappings {
			mapped[mapping] = int(entry.code)
		}
	}
	return mapped
}

func getCharacterCode(character string) (int, bool) {
	code, ok := mappedCharacters[character]
	return code, ok
}

func validateCharacterCode(code int) (int, error) {
	if _, ok := validCharacterCodes[code]; ok {
		return code, nil
	}
	return 0, fmt.Errorf("Invalid Character Code: %d", code)
}

func convertCharactersToCharacterCodes(characters string) ([]int, error) {
	runes := []rune(characters)
	codes := make([]int, 0, len(runes))

	for index := 0; index < len(runes); index++ {
		current := string(runes[index])
		if current == "{" {
			closeIndex := index + 1
			for closeIndex < len(runes) && string(runes[closeIndex]) != "}" {
				closeIndex++
			}
			if closeIndex >= len(runes) {
				continue
			}

			code, err := strconv.Atoi(string(runes[index+1 : closeIndex]))
			if err != nil {
				continue
			}
			validCode, err := validateCharacterCode(code)
			if err != nil {
				return nil, err
			}
			codes = append(codes, validCode)
			index = closeIndex
			continue
		}

		if current == "}" {
			continue
		}

		if code, ok := getCharacterCode(current); ok {
			codes = append(codes, code)
		}
	}

	return codes, nil
}

func mappingToCharacter(character string) string {
	if character == variationSelectorSixteen {
		return ""
	}

	if mapped, ok := multipleCharacterMappings[character]; ok {
		return mapped
	}

	if _, ok := supportedCharacters[character]; ok {
		return character
	}

	for _, entry := range characterCodeTable {
		for _, mapping := range entry.mappings {
			if mapping == character {
				return strings.ToLower(entry.mappings[0])
			}
		}
	}

	return " "
}
