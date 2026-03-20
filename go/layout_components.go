package vbml

type calendarComponent struct {
	characters Board
	x          int
}

func layoutComponents(
	board Board,
	components []Board,
	absoluteComponents []AbsoluteComponent,
	calendarComponents []calendarComponent,
) Board {
	boardHeight := len(board)
	boardWidth := 0
	if boardHeight > 0 {
		boardWidth = len(board[0])
	}

	position := struct {
		top    int
		left   int
		height int
	}{}

	for _, component := range components {
		if len(component) == 0 {
			continue
		}

		newLine := position.left+len(component[0]) > boardWidth
		left := position.left
		top := position.top
		if newLine {
			left = 0
			top = position.top + position.height
		}

		for rowIndex, row := range component {
			for columnIndex, bit := range row {
				if rowIndex+top >= boardHeight {
					continue
				}
				if columnIndex+left >= boardWidth {
					continue
				}
				board[rowIndex+top][columnIndex+left] = bit
			}
		}

		position.top = top
		position.left = left + len(component[0])
		position.height = len(component)
	}

	for _, component := range absoluteComponents {
		for rowIndex, row := range component.Characters {
			for columnIndex, bit := range row {
				if component.Y+rowIndex >= boardHeight {
					continue
				}
				if component.X+columnIndex >= boardWidth {
					continue
				}
				board[component.Y+rowIndex][component.X+columnIndex] = bit
			}
		}
	}

	for _, component := range calendarComponents {
		for rowIndex, row := range component.characters {
			for columnIndex, bit := range row {
				if rowIndex >= boardHeight {
					continue
				}
				if component.x+columnIndex >= boardWidth || columnIndex > 12 {
					continue
				}
				board[rowIndex][component.x+columnIndex] = bit
			}
		}
	}

	return board
}
