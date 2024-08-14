package game_functions

import "log"

func CheckBounds(gameBoard [19][19]int, y, x int) bool {
	if x < 0 || x > 18 || y < 0 || y > 18 {
		log.Println("Out of bounds")
		return false
	}
	//log.Println(gameBoard)
	log.Println("Checking bounds for: ", x, y)
	//log.Println(gameBoard[x][y])
	if gameBoard[x][y] == 1 || gameBoard[x][y] == 2 {
		log.Println("Cell is occupied")
		return false
	}
	//check if cell is occupied with 69, 4, 5, 6 or 7 (bomb/players) - if so, return false
	if gameBoard[x][y] == 69 {
		log.Println("Cell is occupied")
		return false
	}

	// check if the cell is occupied with a player
	for _, player := range Players {
		if player.GridPosition[1] == x && player.GridPosition[0] == y {
			log.Println("Cell is occupied")
			return false
		}
	}

	return true
}
