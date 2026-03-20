package vbml

func renderComponent(emptyComponent, codes Board) Board {
	rendered := CopyCharacterCodes(emptyComponent)
	for rowIndex, row := range rendered {
		for columnIndex := range row {
			if rowIndex < len(codes) && columnIndex < len(codes[rowIndex]) {
				rendered[rowIndex][columnIndex] = codes[rowIndex][columnIndex]
			}
		}
	}
	return rendered
}
