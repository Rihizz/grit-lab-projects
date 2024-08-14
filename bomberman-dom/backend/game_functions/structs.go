package game_functions

import "time"

// Boundaries of the playable area (17x17 tiles)
// max x = 819px, max y = 819px
// min x = 51px, min y = 51px
type Boundaries struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
}

// Player struct contains the player's name, position, and direction, powerups, etc.
type Player struct {
	Name          string
	PlayerID      int
	Lives         int
	Speed         int
	Direction     string
	Bombs         int
	BombRange     int
	GridPosition  [2]int
	PixelPosition [2]int
	LastMove      time.Time
	Left          int
	Top           int
	Conections    int
}

// Players map contains all players in the game
var Players = make(map[int]Player)
