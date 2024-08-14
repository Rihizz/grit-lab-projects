package server

import (
	"bomberman-dom/game_functions"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var ConnectedPlayers []string
var gameBoard []byte
var gameGrid [19][19]int
var isGameBoardGenerated bool
var numConnections = 0
var gridMutex = &sync.Mutex{}
var started = false
var gameOver = false

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserConnection struct {
	UserID     int
	Username   string
	Connection *websocket.Conn
}

var Connections = make(map[*websocket.Conn]*UserConnection)
var ConnectionsByName = make(map[string]*websocket.Conn)

type Message struct {
	Command   string `json:"command"`
	Direction string `json:"direction"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	ID        int    `json:"id"`
	Username  string `json:"name"`
	Timestamp string
	Chat      string `json:"message"`
}

type Msg struct {
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
	Playerlist    []string    `json:"playerlist"`
	NumberOfConns int         `json:"numberOfConns"`
}

// function to read the data from the websocket connection
func reader(conn *websocket.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}

		delete(Connections, conn)

		// Decrement the number of connections.
		numConnections = len(Connections)

		// If there are no more connections, reset the game board.
		if numConnections == 0 {
			isGameBoardGenerated = false
			started = false
			gameOver = false
		}
	}()

	// Set up a close handler for the WebSocket connection
	conn.SetCloseHandler(func(code int, text string) error {
		abandon := false
		log.Printf("WebSocket closed with code %d and text: %s", code, text)
		userConn, ok := Connections[conn]
		if ok {
			disconnectedUserID := userConn.UserID
			if disconnectedUserID < numConnections {
				abandon = true
			}
			log.Printf("User with ID %d disconnected", disconnectedUserID)
		}
		delete(Connections, conn)         // Remove the connection from the map.
		numConnections = len(Connections) // Decrement the number of connections.
		if abandon {
			sendAbandonMessage()
		} else {
			sendPlayerDisconnectedMessage()
		}
		return nil
	})
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// Create a map to store the connections by username
		for _, conn := range Connections {
			ConnectionsByName[conn.Username] = conn.Connection
		}

		log.Println("message received: ")
		log.Println(string(p))
		if messageType == 1 {
			// get the value of the message
			var msg Message
			err := json.Unmarshal(p, &msg)
			if err != nil {
				log.Println(err)
			}
			log.Println("message: ", msg)
			log.Println("message command: ", msg.Command)
			log.Println("message text: ", msg.Direction)
			log.Println("connection length: ", len(Connections))
			if msg.Command == "player" {
				log.Println("HAHAXD")
				log.Println("message: ", msg)
				id := Connections[conn].UserID
				assignName(&game_functions.Players, msg.Username, id)
				ConnectedPlayers = append(ConnectedPlayers, msg.Username)
				// assign a name to the userconnection
				Connections[conn].Username = msg.Username
				sendPlayerConnectedMessage()
			}
			if msg.Command == "start" && !started {
				started = true
				plist := getConnectedPlayerNames()
				log.Println("start command received")
				// send the game board to all the players
				for _, conn := range Connections {
					msg := Msg{
						Type:       "start",
						Data:       gameGrid,
						Playerlist: plist,
					}
					err := conn.Connection.WriteJSON(msg)
					if err != nil {
						log.Println(err)
					}
				}
			}
			if msg.Command == "chat" {
				log.Println("chat command received")
				chatMessage := Connections[conn].Username + ": " + msg.Chat
				for _, conn := range Connections {
					msg := Msg{
						Type: "chat",
						Data: chatMessage,
					}
					err := conn.Connection.WriteJSON(msg)
					if err != nil {
						log.Println(err)
					}
				}
			}

			palyer := Connections[conn].UserID
			if game_functions.Players[palyer].Lives != 0 {
				if msg.Command == "move" {
					log.Println("move command received")
					gameGrid = MovePlayer(gameGrid, Connections[conn].UserID, msg.Direction)
				} else if msg.Command == "place-bomb" {
					log.Println("place-bomb command received")
					PlaceBomb(&gameGrid, Connections[conn].UserID)
				}
			}
		}
	}
}

func getConnectedPlayerNames() []string {
	var names []string
	for _, userConn := range Connections {
		names = append(names, userConn.Username)
	}
	return names
}

// function to send player connected message to all connections
func sendPlayerConnectedMessage() {
	log.Println("sending player connected message")
	player := game_functions.Players[len(game_functions.Players)]
	playerList := getConnectedPlayerNames()
	for _, conn := range Connections {
		msg := Msg{
			Type:          "player-connected",
			Data:          player,
			Playerlist:    playerList,
			NumberOfConns: numConnections,
		}
		err := conn.Connection.WriteJSON(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

// function to send player disconnected message to all connections
func sendPlayerDisconnectedMessage() {
	log.Println("sending player disconnected message")
	playerList := getConnectedPlayerNames()
	for _, conn := range Connections {
		msg := Msg{
			Type:          "player-disconnected",
			Playerlist:    playerList,
			NumberOfConns: numConnections,
		}
		err := conn.Connection.WriteJSON(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func sendAbandonMessage() {
	log.Println("sending abandon message")
	for _, conn := range Connections {
		msg := Msg{
			Type: "abandon",
		}
		err := conn.Connection.WriteJSON(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func assignName(players *map[int]game_functions.Player, name string, ID int) {
	player, exists := (*players)[ID]
	if !exists {
		// Handle error or create a new player
		log.Println("player does not exist")
		return
	}
	player.Name = name
	(*players)[ID] = player
	log.Println("player: ", player)
}

// function to set up the websocket endpoint
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	if numConnections >= 4 || started {
		log.Println("Connection attempt rejected. Maximum number of connections reached. or game has already started")
		w.WriteHeader(http.StatusServiceUnavailable) // Or another appropriate status code
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	userconn := &UserConnection{
		UserID:     numConnections + 1,
		Username:   "player " + strconv.Itoa(numConnections+1),
		Connection: ws,
	}
	// Add the connection to the list of active connections.
	Connections[ws] = userconn
	numConnections = len(Connections) // Increment the number of connections.
	log.Printf("User %s with ID %d successfully connected", userconn.Username, userconn.UserID)
	log.Println("connections: ", Connections)

	//Populate Player struct and initialize Speed, Bombs, Direction, Lives and BombRange.
	//Depending on the ID we also set the coordinates of the player.
	//ID 1 = bottom left corner (1,1)
	//ID 2 = bottom right corner (17,1)
	//ID 3 = top left corner (1,17)
	//ID 4 = top right corner (17,17)
	playerID := userconn.UserID
	player := game_functions.Player{
		PlayerID:  userconn.UserID,
		Name:      userconn.Username,
		Speed:     1,
		Lives:     3,
		Bombs:     1,
		BombRange: 1,
		Direction: "down",
		Left:      0,
		Top:       0,
	}
	switch playerID {
	case 1:
		player.GridPosition = [2]int{1, 1}
		player.PixelPosition = [2]int{51, 51}

	case 2:
		player.GridPosition = [2]int{17, 17}
		player.PixelPosition = [2]int{51, 819}
	case 3:
		player.GridPosition = [2]int{1, 17}
		player.PixelPosition = [2]int{819, 51}
	case 4:
		player.GridPosition = [2]int{17, 1}
		player.PixelPosition = [2]int{819, 819}
	}
	//add player to map of players
	game_functions.Players[userconn.UserID] = player
	// log.Println("Players: ", game_functions.Players)

	// generate the game board
	if !isGameBoardGenerated {
		gameGrid, gameBoard, err = game_functions.GenerateGameBoard()
		if err != nil {
			log.Println(err)
		} else {
			isGameBoardGenerated = true
		}
	}

	go reader(ws)
}

// function to set up the routes for the websocket server.
func SetupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

// A function that changes the player's direction & position based on data received through the WS
// NOTE TO VILLE: make helper function for adding powerups to player
func MovePlayer(gameGrid [19][19]int, playerID int, direction string) [19][19]int {
	//Get the player from the map
	player := game_functions.Players[playerID]

	speed := player.Speed
	lastMove := player.LastMove
	//call isAllowedToMove to check if the player is allowed to move
	if !IsAllowedToMove(speed, lastMove) {
		fmt.Println("Player is not allowed to move yet")
		return gameGrid
	}

	//Change the player's direction
	player.Direction = direction
	//Change the player's position
	log.Println("Player position before: ", player.GridPosition)
	log.Println("Player pixel position before: ", player.PixelPosition)
	switch direction {
	case "up":
		if player.GridPosition[1] == 17 || !game_functions.CheckBounds(gameGrid, player.GridPosition[0], player.GridPosition[1]+1) {
			player.Direction = "none"
			break
		}
		//update the starting position in gamegrid to 0 (remove player id from starting point)
		//we also check if a bomb has been planted, in this case the cell value stays the same to block player from walking over it
		if gameGrid[player.GridPosition[1]][player.GridPosition[0]] != 69 {
			gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 0
		}

		player.GridPosition[1]++

		//if value of gameGrid at player.GridPosition is 8, 9 or 10, update player's powerups
		if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 8 {
			if player.Speed < 3 {
				player.Speed++
			}
		} else if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 9 {
			if player.Bombs < 3 {
				player.Bombs++
			}
		} else if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 10 {
			if player.BombRange < 2 {
				player.BombRange++
			}
		}
		//update the value of gameGrid at player.GridPosition to 0
		gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 0

		//update LastMove
		player.LastMove = time.Now()
	case "down":
		if player.GridPosition[1] == 1 || !game_functions.CheckBounds(gameGrid, player.GridPosition[0], player.GridPosition[1]-1) {
			player.Direction = "none"
			break
		}
		//update the starting position in gamegrid to 0 (remove player id from starting point)
		//we also check if a bomb has been planted, in this case the cell value stays the same to block player from walking over it
		if gameGrid[player.GridPosition[1]][player.GridPosition[0]] != 69 {
			gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 0
		}

		player.GridPosition[1]--

		//if value of gameGrid at player.GridPosition is 8, 9 or 10, update player's powerups
		if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 8 {
			if player.Speed < 3 {
				player.Speed++
			}
		} else if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 9 {
			if player.Bombs < 3 {
				player.Bombs++
			}
		} else if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 10 {
			if player.BombRange < 2 {
				player.BombRange++
			}
		}
		//update the value of gameGrid at player.GridPosition to 0
		gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 0

		//update LastMove
		player.LastMove = time.Now()
	case "left":
		if player.GridPosition[0] == 1 || !game_functions.CheckBounds(gameGrid, player.GridPosition[0]-1, player.GridPosition[1]) {
			player.Direction = "none"
			break
		}
		//update the starting position in gamegrid to 0 (remove player id from starting point)
		//we also check if a bomb has been planted, in this case the cell value stays the same to block player from walking over it
		if gameGrid[player.GridPosition[1]][player.GridPosition[0]] != 69 {
			gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 0
		}

		player.GridPosition[0]--

		//if value of gameGrid at player.GridPosition is 8, 9 or 10, update player's powerups
		if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 8 {
			if player.Speed < 3 {
				player.Speed++
			}
		} else if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 9 {
			if player.Bombs < 3 {
				player.Bombs++
			}
		} else if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 10 {
			if player.BombRange < 2 {
				player.BombRange++
			}
		}
		//update the value of gameGrid at player.GridPosition to 0
		gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 0

		//update LastMove
		player.LastMove = time.Now()
	case "right":
		if player.GridPosition[0] == 17 || !game_functions.CheckBounds(gameGrid, player.GridPosition[0]+1, player.GridPosition[1]) {
			player.Direction = "none"
			break
		}
		//update the starting position in gamegrid to 0 (remove player id from starting point)
		//we also check if a bomb has been planted, in this case the cell value stays the same to block player from walking over it
		if gameGrid[player.GridPosition[1]][player.GridPosition[0]] != 69 {
			gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 0
		}

		player.GridPosition[0]++

		//if value of gameGrid at player.GridPosition is 8, 9 or 10, update player's powerups
		if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 8 {
			if player.Speed < 3 {
				player.Speed++
			}
		} else if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 9 {
			if player.Bombs < 3 {
				player.Bombs++
			}
		} else if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 10 {
			if player.BombRange < 2 {
				player.BombRange++
			}
		}
		//update the value of gameGrid at player.GridPosition to 0
		gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 0

		//update LastMove
		player.LastMove = time.Now()
	}

	log.Println("Player position after: ", player.GridPosition)
	log.Println("Player pixel position after: ", player.PixelPosition)
	//print player powerups
	log.Println("Player speed: ", player.Speed)
	log.Println("Player bombs: ", player.Bombs)
	log.Println("Player bomb range: ", player.BombRange)
	//Update the player in the map
	game_functions.Players[playerID] = player

	playerMsg := Msg{
		Type: "player",
		Data: player,
	}

	//send the updated player to all of the clients
	for _, conn := range Connections {
		err := conn.Connection.WriteJSON(playerMsg)
		if err != nil {
			log.Println(err)
		}
	}
	return gameGrid
}

// A function that receives a players speed and when it was last moved and returns a bool
// if speed is 1, the player can move every 0.5 seconds
// if speed is 2, the player can move every 0.3 seconds
// if speed is 3, the player can move every 0.2 seconds
func IsAllowedToMove(speed int, lastMove time.Time) bool {
	//get the current time
	currentTime := time.Now()
	//get the difference between the current time and the last move
	timeDifference := currentTime.Sub(lastMove)
	fmt.Println("Current time: ", currentTime)
	fmt.Println("Last move: ", lastMove)
	fmt.Println("Time difference: ", timeDifference)
	//if the time difference is less than the allowed time, return false
	switch speed {
	case 1:
		if timeDifference.Seconds() < 0.5 {
			return false
		}
	case 2:
		if timeDifference.Seconds() < 0.35 {
			return false
		}
	case 3:
		if timeDifference.Seconds() < 0.2 {
			return false
		}
	}
	//	fmt.Println("Player is allowed to move")
	return true
}

// A function that places a bomb on the game board
func PlaceBomb(gameGrid *[19][19]int, playerID int) {
	gridMutex.Lock()
	log.Println("Game grid: ", gameGrid)
	// Get the player from the map
	player := game_functions.Players[playerID]
	// CHANGE THIS TO THE CONCURRENT BOMB LIMIT!!!!!!!
	if player.Bombs == 0 {
		log.Println("Bomb limit reached")
		defer gridMutex.Unlock()
		return
	}
	// Check if the player is on a bomb
	if gameGrid[player.GridPosition[1]][player.GridPosition[0]] == 69 {
		log.Println("Player is on a bomb")
		defer gridMutex.Unlock()
		return
	}
	// Check the player's bomb range
	bombRange := player.BombRange
	log.Println("Bomb range: ", bombRange)
	player.Bombs--

	// Place a bomb on the game board
	gameGrid[player.GridPosition[1]][player.GridPosition[0]] = 69
	log.Println("Bomb placed at: ", player.GridPosition)
	// Update the player in the map
	game_functions.Players[playerID] = player

	log.Println("Game grid after bomb placement: ", gameGrid)

	// Prepare bomb message
	bombMsg := Msg{
		Type: "bomb",
		Data: struct {
			GridPosition []int
		}{GridPosition: []int{player.GridPosition[0], player.GridPosition[1]}},
	}

	// Send the bomb message to all of the clients
	for _, conn := range Connections {
		err := conn.Connection.WriteJSON(bombMsg)
		if err != nil {
			log.Println(err)
		}
	}

	// Start a timer for the bomb explosion
	go func() {
		timer := time.NewTimer(time.Duration(2) * time.Second)
		<-timer.C
		HandleExplosion(gameGrid, player.GridPosition[0], player.GridPosition[1], bombRange, playerID)
	}()

	defer gridMutex.Unlock()

	player.Direction = "none"

	playerMsg := Msg{
		Type: "player",
		Data: player,
	}

	// Send the updated player to all of the clients
	for _, conn := range Connections {
		err := conn.Connection.WriteJSON(playerMsg)
		if err != nil {
			log.Println(err)
		}
	}

	gameGridMsg := Msg{
		Type: "gameGrid",
		Data: gameGrid,
	}

	// Send the updated game board to all of the clients
	for _, conn := range Connections {
		err := conn.Connection.WriteJSON(gameGridMsg)
		if err != nil {
			log.Println(err)
		}
	}
}

func HandleExplosion(gameGrid *[19][19]int, x int, y int, bombRange int, playerID int) {
	gridMutex.Lock()
	log.Println("---------------------BOOOOOOOOOOOOOOOOOOOOM---------------------")
	printBoard(gameGrid)
	log.Println("Bomb exploded at: ", x, y)
	// Clean the bomb from gameGrid
	gameGrid[y][x] = 0

	player := game_functions.Players[playerID]
	player.Bombs++
	game_functions.Players[playerID] = player

	var explosionCells [][]int
	var affectedCells [][]int
	var affectedPlayers []int
	rand.Seed(time.Now().UnixNano())

	// explosionCells = append(explosionCells, []int{x, y})

	handleExplosionLogic(x, y, gameGrid, &affectedCells, &affectedPlayers, &explosionCells)

	// Explosion to the right
	for i := 1; i <= bombRange; i++ {
		if x+i < 19 {
			if gameGrid[y][x+i] == 2 {
				break
			}
			handleExplosionLogic(x+i, y, gameGrid, &affectedCells, &affectedPlayers, &explosionCells)
		}
	}

	// Explosion to the left
	for i := 1; i <= bombRange; i++ {
		if x-i >= 0 {
			if gameGrid[y][x-i] == 2 {
				break
			}
			handleExplosionLogic(x-i, y, gameGrid, &affectedCells, &affectedPlayers, &explosionCells)
		}
	}

	// Explosion upwards
	for i := 1; i <= bombRange; i++ {
		if y+i < 19 {
			if gameGrid[y+i][x] == 2 {
				break
			}
			handleExplosionLogic(x, y+i, gameGrid, &affectedCells, &affectedPlayers, &explosionCells)
		}
	}

	// Explosion downwards
	for i := 1; i <= bombRange; i++ {
		if y-i >= 0 {
			if gameGrid[y-i][x] == 2 {
				break
			}
			handleExplosionLogic(x, y-i, gameGrid, &affectedCells, &affectedPlayers, &explosionCells)
		}
	}

	gridMutex.Unlock()

	// Log the state of game grid after explosion
	printBoard(gameGrid)
	log.Println("Affected cells: ", affectedCells)
	log.Println("Affected players: ", affectedPlayers)

	explosionMsg := Msg{
		Type: "explosion",
		Data: struct {
			GameGrid        *[19][19]int
			AffectedCells   [][]int
			AffectedPlayers []int
			ExplosionCells  [][]int
		}{
			GameGrid:        gameGrid,
			AffectedCells:   affectedCells,
			AffectedPlayers: affectedPlayers,
			ExplosionCells:  explosionCells,
		},
	}

	//send the updated game board to all of the clients
	for _, conn := range Connections {
		err := conn.Connection.WriteJSON(explosionMsg)
		if err != nil {
			log.Println(err)
		}
	}

	if len(affectedPlayers) > 0 {
		for _, playerID := range affectedPlayers {
			affectedPlayer, exists := game_functions.Players[playerID]
			if exists {
				playerMsg := Msg{
					Type: "player",
					Data: affectedPlayer,
				}

				log.Println(playerMsg)

				//send the updated players to all of the clients
				for _, conn := range Connections {
					err := conn.Connection.WriteJSON(playerMsg)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}

	// Check if the game is over
	if gameOver {
		sendGameOverMessage()
	}
}

// handleExplosionLogic is a helper function to manage the explosion logic at a given cell.
func handleExplosionLogic(x int, y int, gameGrid *[19][19]int, affectedCells *[][]int, affectedPlayers *[]int, explosionCells *[][]int) {
	if gameGrid[y][x] == 1 {
		// 20% chance of a power up, power ups distributed evenly
		if rand.Float32() < 0.4 {
			gameGrid[y][x] = 8 + rand.Intn(3)
			log.Println("Power up placed at: ", x, y)
		} else {
			gameGrid[y][x] = 0
			log.Println("Block destroyed at: ", x, y)
		}
		*affectedCells = append(*affectedCells, []int{x, y})
		*explosionCells = append(*explosionCells, []int{x, y})
	} else if gameGrid[y][x] == 0 {
		*explosionCells = append(*explosionCells, []int{x, y})
	}
	for _, player := range game_functions.Players {
		if player.GridPosition[0] == x && player.GridPosition[1] == y {
			player.Lives--
			if player.Lives == 0 {
				// Player is dead
				player.BombRange = 1
				player.GridPosition[0] = 0
				player.GridPosition[1] = 0
				game_functions.Players[player.PlayerID] = player
				log.Println("Player ", player.PlayerID, " is dead")
				*affectedPlayers = append(*affectedPlayers, player.PlayerID)

				// check if the game is over
				alivePlayers := 0
				for _, player := range game_functions.Players {
					if player.Lives > 0 {
						alivePlayers++
					}
				}
				if alivePlayers == 1 {
					gameOver = true

					// Game is over
					log.Println("Game is over")
					// sendGameOverMessage()
				}

			} else {
				// Player is alive
				player.BombRange = 1
				game_functions.Players[player.PlayerID] = player
				log.Println("Player ", player.PlayerID, " is alive")
				*affectedPlayers = append(*affectedPlayers, player.PlayerID)
			}
		}
	}
}

func sendGameOverMessage() {
	log.Println("in sendGameOverMessage")
	// range over all of the players and find the winner
	winner := ""

	for _, player := range game_functions.Players {
		if player.Lives > 0 {
			winner = player.Name
		}
	}

	gameOverMsg := Msg{
		Type: "gameOver",
		Data: struct {
			Winner string
		}{
			Winner: winner,
		},
	}
	log.Println("Sending game over message in 1 seconds")
	time.Sleep(1 * time.Second)

	//send the gameOverMsg to all of the clients
	for _, conn := range Connections {
		err := conn.Connection.WriteJSON(gameOverMsg)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Game over message sent")
}

func printBoard(board *[19][19]int) {
	// Start printing from the last row to the first
	for i := len(board) - 1; i >= 0; i-- {
		row := board[i]
		for _, cell := range row {
			// Each cell is printed with a fixed width of 3 for better alignment
			fmt.Printf("%3d", cell)
		}
		fmt.Println() // Move to the next line after printing each row
	}
}
