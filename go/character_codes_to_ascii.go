package vbml

import "strings"

func CharacterCodesToASCII(characterCodes Board, isWhite bool) string {
	filled := "⬜"
	if isWhite {
		filled = "⬛"
	}

	characterMap := map[int]string{
		0:  "  ",
		1:  "A ",
		2:  "B ",
		3:  "C ",
		4:  "D ",
		5:  "E ",
		6:  "F ",
		7:  "G ",
		8:  "H ",
		9:  "I ",
		10: "J ",
		11: "K ",
		12: "L ",
		13: "M ",
		14: "N ",
		15: "O ",
		16: "P ",
		17: "Q ",
		18: "R ",
		19: "S ",
		20: "T ",
		21: "U ",
		22: "V ",
		23: "W ",
		24: "X ",
		25: "Y ",
		26: "Z ",
		27: "1 ",
		28: "2 ",
		29: "3 ",
		30: "4 ",
		31: "5 ",
		32: "6 ",
		33: "7 ",
		34: "8 ",
		35: "9 ",
		36: "0 ",
		37: "! ",
		38: "@ ",
		39: "# ",
		40: "$ ",
		41: "( ",
		42: ") ",
		43: "  ",
		44: "- ",
		45: "  ",
		46: "+ ",
		47: "& ",
		48: "= ",
		49: "; ",
		50: ": ",
		51: "  ",
		52: "' ",
		53: `" `,
		54: "% ",
		55: ", ",
		56: ". ",
		57: "  ",
		58: "  ",
		59: "/ ",
		60: "? ",
		61: "  ",
		62: "° ",
		63: "🟥",
		64: "🟧",
		65: "🟨",
		66: "🟩",
		67: "🟦",
		68: "🟪",
		69: "⬜",
		70: "⬛",
		71: filled,
	}

	rows := make([]string, 0, len(characterCodes))
	for _, row := range characterCodes {
		var builder strings.Builder
		for _, code := range row {
			builder.WriteString(characterMap[code])
		}
		rows = append(rows, builder.String())
	}

	return strings.Join(rows, "\n\n")
}
