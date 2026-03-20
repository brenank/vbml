package vbml

import rand "math/rand/v2"

func randomColors(rows, columns int, colors []int) Board {
	if rows == 0 {
		rows = flagshipBoardHeight
	}
	if columns == 0 {
		columns = flagshipBoardWidth
	}
	if len(colors) == 0 {
		colors = colorCodes
	}

	board := make(Board, rows)
	for rowIndex := range board {
		row := make([]int, columns)
		for columnIndex := range row {
			row[columnIndex] = colors[rand.IntN(len(colors))]
		}
		board[rowIndex] = row
	}
	return board
}
