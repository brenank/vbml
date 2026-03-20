package vbml

func verticalAlign(height int, align Align, codes Board) Board {
	switch align {
	case AlignTop:
		return CopyCharacterCodes(codes)
	case AlignBottom:
		reversedCodes := reverseRows(codes)
		rows := make(Board, height)
		for index := 0; index < height; index++ {
			sourceIndex := height - 1 - index
			if sourceIndex >= 0 && sourceIndex < len(reversedCodes) {
				rows[index] = append([]int(nil), reversedCodes[sourceIndex]...)
				continue
			}
			rows[index] = []int{}
		}
		return rows
	case AlignJustified:
		paddingTop := (height - len(codes) + 1) / 2
		rows := make(Board, height)
		for index := 0; index < height; index++ {
			sourceIndex := index - paddingTop
			if sourceIndex >= 0 && sourceIndex < len(codes) {
				rows[index] = append([]int(nil), codes[sourceIndex]...)
				continue
			}
			rows[index] = []int{}
		}
		return rows
	default:
		paddingTop := (height - len(codes)) / 2
		rows := make(Board, height)
		for index := 0; index < height; index++ {
			sourceIndex := index - paddingTop
			if sourceIndex >= 0 && sourceIndex < len(codes) {
				rows[index] = append([]int(nil), codes[sourceIndex]...)
				continue
			}
			rows[index] = []int{}
		}
		return rows
	}
}

func reverseRows(board Board) Board {
	reversed := make(Board, len(board))
	for index := range board {
		reversed[len(board)-1-index] = append([]int(nil), board[index]...)
	}
	return reversed
}
