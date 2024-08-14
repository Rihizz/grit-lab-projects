# bomberman

## Objective
In this project, we created a multiplayer version of the classic Bomberman game. The game will allow multiple players to join and battle against each other until only one player remains victorious.

## Instructions
1. Clone the repository
2. Separeate terminal windows for server and client, run ```go run . ```in the backend folder 
3. Open localhost in browser and join the game!

## Map and Players
The game will support 2 to 4 players.
Each player will start with 3 lives, and when they run out of lives, they are eliminated from the game.
The map will be fixed, and all players will see the whole map.
The map will consist of two types of blocks: destroyable blocks (blocks) and indestructible blocks (walls).
Walls will always be placed in the same locations, while the blocks will be randomly generated on the map.
Players will start in the corners of the map to ensure they have space to avoid bomb explosions.


## Power-ups
Destroying a block may release a random power-up on the map. Players can collect these power-ups to gain advantages during the game.

### Power-ups include:
- Bombs with plus: Increase the number of bombs dropped at a time by 1.
- Bombs with arrow ups: Increase the explosion range from the bomb in four directions by 1 block.
- Speedarrows: Increase movement speed.

## Chat
Implement a chat feature that allows players to communicate with each other during the game.

## Game Start and Waiting
When a user opens the game, they will be prompted to enter a nickname for identification purposes.
After selecting a nickname, the user will be taken to a waiting page displaying a player counter, which will increment as more players join.
If there are more than 2 players in the counter and it does not reach 4 players within 20 seconds, a 10-second timer starts to give players time to prepare before the game starts.
If there are 4 players in the counter before 20 seconds, a 10-second timer starts, and the game begins.


## Contributors
Oskar, Santeri, Stefanie, Ville and Wincent July 2023