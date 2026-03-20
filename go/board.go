package vbml

func createEmptyBoard(rows, columns int) Board {
	board := make(Board, rows)
	for rowIndex := range board {
		board[rowIndex] = make([]int, columns)
	}
	return board
}

func CopyCharacterCodes(characters Board) Board {
	copyBoard := make(Board, len(characters))
	for rowIndex, row := range characters {
		copyRow := make([]int, len(row))
		copy(copyRow, row)
		copyBoard[rowIndex] = copyRow
	}
	return copyBoard
}
