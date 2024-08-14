
# Client-side (JavaScript):  

This is where the game is displayed to the users, and where user input is collected and data is sent to backend. The client-side code is responsible for rendering the game map, player characters, and other game elements based on the game state it receives from backend game engine.  

Example of function: 

We receive data that a player placed a bomb in square (5,5) 

We render a bomb on to that square (apply a CSS class?) and start the animation 

Upon receiving data about the bomb timer running out, we play an animation of the explosion and re-render (change CSS classes?) 

NOTE: we will not run timers for bombs in the front-end, only apply animations and re-rendering based on data received from backend 

 

# Server-side (Go):  

The server maintains the game state, processes player inputs, and sends updates to all connected clients. The server is the source of truth for the game state, ensuring all players see the same game world. 

The backend will be responsible for map randomizing and running the game engine. 

 
Example of function: 
We receive data that a player placed a bomb in square (5,5) 
Server assigns a bomb to that square in the state grid, relays this information to clients and starts a 3 second timer 

After the timer has run its course, we check the bombs radius and if there are items/players on the grid within its reach, change the game state accordingly and update clients with new state 

 

# Data transmission (WebSockets):  

To achieve real-time updates, we will have open, two-way communication between the client and the server. WebSockets is a protocol that allows for this. The front-end will have event listeners for user input, that will be sent to the back ends game engine that updates the state and relays the updated state to all players. 

 

# Data storage:  

Game state data will be stored in backend memory and it will consist of the map, player information, bombs and randomized blocks/items. 

When game is over, all data excluding player names will be removed. 

Players will be removed when they disconnect their websockets. 

 
## Game flow

1. Landing page 

Player is given instructions to the game and prompted for a nickname, on submit they will be directed to the lobby 

 

2. Lobby Phase: 

When the lobby has enough players, the server initiates the game start. 

 

3. Game Start: 

The server generates a random map and sends it to all clients. 

The server assigns player positions and sends them to all clients. 

 

4. During the Game: 

Clients send their player inputs (like movements) to the server. 

The server processes these inputs, updates the game state, and sends the updated game state to all clients. 

Clients update their display based on the updated game state. 

 

5. Game End: 

When the game ends, the server sends a game over message to all clients. 

All clients are redirected back to the lobby. 
