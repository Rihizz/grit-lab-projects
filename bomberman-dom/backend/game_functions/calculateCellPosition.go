package game_functions

// CellPosition is a struct that contains the x and y coordinates of a cell.
type CellPosition struct {
	X         int
	Y         int
	BlockType int
}

// Calculate and store the boundaries based on the grid
func CalculateBoundaries(grid [19][19]int) Boundaries {
	boundaries := Boundaries{
		MinX: 50,
		MaxX: (len(grid[0]) - 1) * 50,
		MinY: 50,
		MaxY: (len(grid) - 1) * 50,
	}
	return boundaries
}

func CalculateCellPosition(xPixels, yPixels int, grid [19][19]int) CellPosition {
	cellPosition := CellPosition{
		X: xPixels / 50,
		Y: yPixels / 50,
	}
	cellPosition.BlockType = grid[cellPosition.Y][cellPosition.X]
	return cellPosition
}
