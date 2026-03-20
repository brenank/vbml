package vbml

func horizontalAlign(width int, justify Justify, codes Board) Board {
	switch justify {
	case JustifyLeft:
		rows := make(Board, len(codes))
		for index, row := range codes {
			rows[index] = removeLeadingBlank(row)
		}
		return rows
	case JustifyRight:
		rows := make(Board, len(codes))
		for index, row := range codes {
			stripped := reverseInts(removeLeadingBlank(reverseInts(row)))
			result := make([]int, width)
			for fillIndex := range result {
				result[fillIndex] = int(CharacterCodeBlank)
			}
			for valueIndex, value := range reverseInts(stripped) {
				if valueIndex >= width {
					break
				}
				result[width-1-valueIndex] = value
			}
			rows[index] = result
		}
		return rows
	case JustifyJustified:
		rows := make(Board, len(codes))
		longestRow := 0
		strippedRows := make(Board, len(codes))
		for index, row := range codes {
			stripped := removeLeadingBlank(row)
			strippedRows[index] = stripped
			if len(stripped) > longestRow {
				longestRow = len(stripped)
			}
		}
		longestRow--
		paddingRight := (width - longestRow) / 2
		paddingLeft := width - (longestRow + (paddingRight + 1))
		padding := paddingRight
		if paddingLeft < padding {
			padding = paddingLeft
		}
		for index, row := range strippedRows {
			result := make([]int, width)
			for column := range result {
				sourceIndex := column - padding
				if sourceIndex >= 0 && sourceIndex < len(row) {
					result[column] = row[sourceIndex]
				}
			}
			rows[index] = result
		}
		return rows
	default:
		rows := make(Board, len(codes))
		for index, row := range codes {
			stripped := reverseInts(removeLeadingBlank(reverseInts(row)))
			paddingLeft := (width - len(stripped)) / 2
			result := make([]int, width)
			for column := range result {
				sourceIndex := column - paddingLeft
				if sourceIndex >= 0 && sourceIndex < len(stripped) {
					result[column] = stripped[sourceIndex]
				}
			}
			rows[index] = result
		}
		return rows
	}
}

func removeLeadingBlank(row []int) []int {
	result := make([]int, 0, len(row))
	skipping := true
	for _, code := range row {
		if code == int(CharacterCodeBlank) && skipping {
			continue
		}
		skipping = false
		result = append(result, code)
	}
	return result
}

func reverseInts(values []int) []int {
	reversed := make([]int, len(values))
	for index := range values {
		reversed[len(values)-1-index] = values[index]
	}
	return reversed
}
