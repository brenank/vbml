package vbml

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

var classicCharacterMap = buildClassicCharacterMap()

const classicUnknownCharacter = -1

func Classic(text string, options ClassicOptions) Board {
	emptyBoard := createEmptyBoard(flagshipBoardHeight, flagshipBoardWidth)
	if text == "" {
		return emptyBoard
	}

	lines := strings.Split(emojisToCharacterCodes(text), "\n")
	chunkedLines := make([][]string, 0, len(lines))
	for _, line := range lines {
		chunkedLines = append(chunkedLines, tokenizeClassicLine(line, options.PreserveDoubleSpaces))
	}

	vestaLines := make([][]int, 0, len(chunkedLines))
	for _, line := range chunkedLines {
		var characters []int
		for _, word := range line {
			characters = append(characters, classicWordToCodes(word, options.PreserveDoubleSpaces)...)
		}
		vestaLines = append(vestaLines, characters)
	}

	wordsLines := make([][][]int, 0, len(vestaLines))
	for _, line := range vestaLines {
		wordsLines = append(wordsLines, splitClassicWords(line, options.PreserveDoubleSpaces))
	}

	contentWidth := flagshipBoardWidth - options.ExtraHPadding
	var wrapping [][][]int
	for _, words := range wordsLines {
		wrapping = append(wrapping, classicMakeLines(words, contentWidth)...)
	}

	formatted := make([][][]int, 0, len(wrapping))
	for _, line := range wrapping {
		var row [][]int
		for index, word := range line {
			row = append(row, word)
			if index != len(line)-1 {
				row = append(row, []int{0})
			}
		}
		formatted = append(formatted, row)
	}

	if len(formatted) == 3 && options.ExtraHPadding == 0 {
		return Classic(text, ClassicOptions{
			ExtraHPadding:        4,
			PreserveDoubleSpaces: options.PreserveDoubleSpaces,
		})
	}

	maxColumns := 0
	for _, line := range formatted {
		width := 0
		for _, word := range line {
			width += len(word)
		}
		if width > maxColumns {
			maxColumns = width
		}
	}

	horizontalPadding := maxInt(int(math.Floor(float64(flagshipBoardWidth-(maxColumns+1))/2)), 0)
	verticalPadding := maxInt(int(math.Floor(float64(flagshipBoardHeight-len(formatted))/2)), 0)

	paddingWords := make([][]int, horizontalPadding)
	for index := range paddingWords {
		paddingWords[index] = []int{0}
	}

	var paddedRows []interface{}
	for index := 0; index < verticalPadding; index++ {
		paddedRows = append(paddedRows, make([]int, flagshipBoardWidth))
	}
	for _, line := range formatted {
		row := make([][]int, 0, len(line)+(horizontalPadding*2))
		row = append(row, paddingWords...)
		row = append(row, line...)
		row = append(row, paddingWords...)
		paddedRows = append(paddedRows, row)
	}
	for index := 0; index < verticalPadding; index++ {
		paddedRows = append(paddedRows, make([]int, flagshipBoardWidth))
	}

	codes := make(Board, 0, flagshipBoardHeight)
	for _, row := range paddedRows {
		codes = append(codes, flattenClassicRow(row)[:minInt(len(flattenClassicRow(row)), flagshipBoardWidth)])
		if len(codes) == flagshipBoardHeight {
			break
		}
	}

	board := createEmptyBoard(flagshipBoardHeight, flagshipBoardWidth)
	for rowIndex := 0; rowIndex < flagshipBoardHeight; rowIndex++ {
		for columnIndex := 0; columnIndex < flagshipBoardWidth; columnIndex++ {
			if rowIndex < len(codes) && columnIndex < len(codes[rowIndex]) {
				if codes[rowIndex][columnIndex] > 0 {
					board[rowIndex][columnIndex] = codes[rowIndex][columnIndex]
				}
			}
		}
	}
	return board
}

func buildClassicCharacterMap() map[string]int {
	characterMap := map[string]int{
		" ":      0,
		"!":      37,
		"@":      38,
		"#":      39,
		"$":      40,
		"(":      41,
		")":      42,
		"-":      44,
		"+":      46,
		"&":      47,
		"=":      48,
		";":      49,
		":":      50,
		"'":      52,
		`"`:      53,
		"‟":      53,
		"\u201c": 53,
		"\u201d": 53,
		"„":      53,
		"¨":      53,
		"\u2019": 52,
		"´":      52,
		"ˋ":      52,
		"ˊ":      52,
		"‚":      52,
		"`":      52,
		"%":      54,
		",":      55,
		".":      56,
		"/":      59,
		"\\":     59,
		"?":      60,
		"°":      62,
		"—":      44,
		"–":      44,
		"¯":      44,
		"~":      44,
		"¸":      55,
		"¦":      50,
		"¿":      60,
		"[":      41,
		"{":      41,
		"]":      42,
		"}":      42,
		"‰":      54,
		"¤":      62,
		"•":      62,
		"·":      62,
		"â":      1,
		"à":      1,
		"å":      1,
		"á":      1,
		"À":      1,
		"Á":      1,
		"Â":      1,
		"Ã":      1,
		"Ä":      1,
		"Å":      1,
		"ã":      1,
		"ç":      3,
		"Ç":      3,
		"¢":      3,
		"Ð":      4,
		"é":      5,
		"ê":      5,
		"ë":      5,
		"è":      5,
		"È":      5,
		"É":      5,
		"Ê":      5,
		"Ë":      5,
		"ƒ":      6,
		"í":      9,
		"ï":      9,
		"î":      9,
		"ì":      9,
		"Ì":      9,
		"Í":      9,
		"Î":      9,
		"Ï":      9,
		"|":      9,
		"£":      12,
		"ñ":      14,
		"Ñ":      14,
		"ó":      15,
		"ô":      15,
		"ö":      15,
		"ò":      15,
		"Ò":      15,
		"Ó":      15,
		"Ô":      15,
		"Õ":      15,
		"Ö":      15,
		"Ø":      15,
		"ð":      15,
		"õ":      15,
		"ø":      15,
		"±":      46,
		"š":      19,
		"Š":      19,
		"§":      19,
		"û":      21,
		"ù":      21,
		"ú":      21,
		"Ù":      21,
		"Ú":      21,
		"Û":      21,
		"ğ":      7,
	}

	for index, letter := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		characterMap[string(letter)] = index + 1
		characterMap[string(unicode.ToLower(letter))] = index + 1
	}
	for index, digit := range "1234567890" {
		characterMap[string(digit)] = index + 27
	}
	for code := 0; code <= 71; code++ {
		characterMap["{"+strconv.Itoa(code)+"}"] = code
	}

	return characterMap
}

func tokenizeClassicLine(line string, preserveDoubleSpaces bool) []string {
	runes := []rune(line)
	var tokens []string

	for index := 0; index < len(runes); {
		switch {
		case isASCIILetter(runes[index]):
			start := index
			for index < len(runes) && isASCIILetter(runes[index]) {
				index++
			}
			tokens = append(tokens, string(runes[start:index]))
		case runes[index] == '{':
			start := index
			index++
			for index < len(runes) && unicode.IsDigit(runes[index]) {
				index++
			}
			if index > start+1 && index < len(runes) && runes[index] == '}' {
				index++
				tokens = append(tokens, string(runes[start:index]))
			} else {
				tokens = append(tokens, string(runes[start]))
				index = start + 1
			}
		case unicode.IsDigit(runes[index]):
			start := index
			for index < len(runes) && unicode.IsDigit(runes[index]) {
				index++
			}
			tokens = append(tokens, string(runes[start:index]))
		case runes[index] == '_':
			index++
		case runes[index] == ' ':
			if preserveDoubleSpaces {
				if index+1 < len(runes) && runes[index+1] == ' ' {
					tokens = append(tokens, "  ")
					index += 2
				} else {
					tokens = append(tokens, " ")
					index++
				}
				continue
			}

			start := index
			for index < len(runes) && unicode.IsSpace(runes[index]) {
				index++
			}
			tokens = append(tokens, string(runes[start:index]))
		default:
			tokens = append(tokens, string(runes[index]))
			index++
		}
	}

	return tokens
}

func classicWordToCodes(word string, preserveDoubleSpaces bool) []int {
	if strings.HasPrefix(word, "{") && strings.HasSuffix(word, "}") {
		if code, ok := classicCharacterMap[word]; ok {
			return []int{code}
		}
		return []int{classicUnknownCharacter}
	}
	if preserveDoubleSpaces && word == "  " {
		return []int{0, 0}
	}

	runes := []rune(word)
	codes := make([]int, 0, len(runes))
	for _, character := range runes {
		if character == 'ä' || character == 'Ä' {
			codes = append(codes, classicCharacterMap["a"], classicCharacterMap["e"])
			continue
		}
		if code, ok := classicCharacterMap[string(character)]; ok {
			codes = append(codes, code)
			continue
		}
		codes = append(codes, classicUnknownCharacter)
	}
	return codes
}

func splitClassicWords(characters []int, preserveDoubleSpaces bool) [][]int {
	var words [][]int
	word := make([]int, 0)
	for index, code := range characters {
		if code == 0 && !preserveDoubleSpaces {
			words = append(words, word)
			word = []int{}
			continue
		}
		if code == 0 && preserveDoubleSpaces {
			nextIsZero := index+1 < len(characters) && characters[index+1] == 0
			prevIsZero := index > 0 && characters[index-1] == 0
			if nextIsZero || prevIsZero {
				word = append(word, code)
			} else {
				words = append(words, word)
				word = []int{}
			}
			continue
		}
		word = append(word, code)
	}
	words = append(words, word)
	return words
}

func classicMakeLines(words [][]int, contentWidth int) [][][]int {
	filtered := make([][]int, 0, len(words))
	for _, word := range words {
		if len(word) > 0 {
			filtered = append(filtered, word)
		}
	}

	var chunked [][]int
	for _, word := range filtered {
		for start := 0; start < maxInt(len(word), 1); start += contentWidth {
			end := start + contentWidth
			if end > len(word) {
				end = len(word)
			}
			chunked = append(chunked, append([]int(nil), word[start:end]...))
		}
	}
	filtered = chunked

	total := 0
	for _, word := range filtered {
		total += len(word)
	}
	if len(filtered) > 0 {
		total += len(filtered) - 1
	}
	if total <= contentWidth {
		return [][][]int{filtered}
	}

	for index := 1; index <= len(filtered); index++ {
		subset := filtered[:index]
		needed := 0
		for _, word := range subset {
			needed += len(word)
		}
		needed += len(subset) - 1
		if needed > contentWidth {
			return append([][][]int{filtered[:index-1]}, classicMakeLines(filtered[index-1:], contentWidth)...)
		}
	}
	return nil
}

func flattenClassicRow(row any) []int {
	switch row := row.(type) {
	case []int:
		return append([]int(nil), row...)
	case [][]int:
		flattened := make([]int, 0)
		for _, word := range row {
			flattened = append(flattened, word...)
		}
		return flattened
	default:
		return nil
	}
}

func isASCIILetter(character rune) bool {
	return ('a' <= character && character <= 'z') || ('A' <= character && character <= 'Z')
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}
