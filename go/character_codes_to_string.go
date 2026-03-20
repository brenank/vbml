package vbml

import "strings"

var characterCodesToStringMap = map[int]string{
	0:   " ",
	1:   "A",
	2:   "B",
	3:   "C",
	4:   "D",
	5:   "E",
	6:   "F",
	7:   "G",
	8:   "H",
	9:   "I",
	10:  "J",
	11:  "K",
	12:  "L",
	13:  "M",
	14:  "N",
	15:  "O",
	16:  "P",
	17:  "Q",
	18:  "R",
	19:  "S",
	20:  "T",
	21:  "U",
	22:  "V",
	23:  "W",
	24:  "X",
	25:  "Y",
	26:  "Z",
	27:  "1",
	28:  "2",
	29:  "3",
	30:  "4",
	31:  "5",
	32:  "6",
	33:  "7",
	34:  "8",
	35:  "9",
	36:  "0",
	37:  "!",
	38:  "@",
	39:  "#",
	40:  "$",
	41:  "(",
	42:  ")",
	43:  " ",
	44:  "-",
	45:  "",
	46:  "+",
	47:  "&",
	48:  "=",
	49:  ";",
	50:  ":",
	51:  "",
	52:  "'",
	53:  `"`,
	54:  "%",
	55:  ",",
	56:  ".",
	57:  "",
	58:  "",
	59:  "/",
	60:  "?",
	61:  "",
	62:  "°",
	63:  "",
	64:  "",
	65:  "",
	66:  "",
	67:  "",
	68:  "",
	69:  "",
	70:  "",
	71:  " ",
	100: "\n",
}

func CharacterCodesToString(characters Board, options CharacterCodesToStringOptions) string {
	merged := make([]int, 0)
	for index, row := range characters {
		if index == 0 {
			merged = append(merged, row...)
			continue
		}

		separator := 0
		if options.AllowLineBreaks {
			previous := characters[index-1]
			prefix := countEmptyCharactersBeforeFirstWord(previous)
			postfix := countEmptyCharactersBeforeFirstWord(reverseInts(previous))
			firstWord := countFirstWordLength(row)
			if prefix+postfix > firstWord {
				separator = 100
			}
		}

		merged = append(merged, separator)
		merged = append(merged, row...)
	}

	var builder strings.Builder
	for _, code := range merged {
		builder.WriteString(characterCodesToStringMap[code])
	}

	text := strings.TrimSpace(builder.String())
	text = collapseSpaces(text)
	text = strings.ReplaceAll(text, " \n", "\n")
	return text
}

func countEmptyCharactersBeforeFirstWord(row []int) int {
	count := 0
	counting := true
	for _, code := range row {
		if !isBreakableCharacter(code) || !counting {
			counting = false
			continue
		}
		count++
	}
	return count
}

func countFirstWordLength(row []int) int {
	count := 0
	counting := true
	started := false
	for _, code := range row {
		if !counting {
			break
		}
		isCharacter := !isBreakableCharacter(code)
		if isCharacter {
			count++
			started = true
			continue
		}
		if started {
			counting = false
		}
	}
	return count
}

func isBreakableCharacter(code int) bool {
	mapped := characterCodesToStringMap[code]
	return mapped == "" || mapped == " "
}

func collapseSpaces(text string) string {
	parts := strings.Split(text, " ")
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return strings.Join(filtered, " ")
}
