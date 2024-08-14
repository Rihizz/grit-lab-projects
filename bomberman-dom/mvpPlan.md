# Minimum Viable Product Plan
Here we want to have a list of questions & answers related to the project planning and framing what the finished product will look like

## Gameplay Mechanics:

What are the movement keys?
> WASD / arrow keys for moving around, space to deploy a bomb. These should be shown in the how to play section in the main screen.

How will player movement be handled? Is it tile-based or continuous?<br>
> Continuous movement on a tile based map<br>
  
How will bomb placement and detonation work?<br>
How will the power-ups be implemented? What is the duration of each power-up?<br>
> Power ups will be found under destroyable blocks and will last untill you perish

How will the game handle a player's death and respawn?<br>
> When a player dies, they are presented with a text "You have died - you will be directed back to the lobby when the game is finished". Then the player sprite will be removed from the map. 

How should the collision detection work between players, bombs, and the environment (blocks and walls)?<br>
> You can not walk through players, bombs or any walls. On contact with a power up player is automatically upgraded with it.

How will the bomb explosion mechanics be implemented? What happens if a bomb's explosion reaches another bomb?<br>
> If a bomb is in the way of an explosion, it will detonate immediatly. An explosion will also be able to destroy a power up.

How can players pick up power-ups? Is it automatic or does the player need to perform an action?<br>
> Walk over & automatic assignment


## Player and Game State Management:

How will the game track the current state (waiting, in-progress, finished)?<br>
How will the game handle player disconnections or reconnections?<br>
How will you represent players in the game state? What information needs to be stored for each player?<br>
How will the game state change over time, especially when power-ups are collected or players are eliminated?<br>
How will the game handle players who want to join a game that's already in progress?<br>

## Map Design:

What size will the game map be, and will it vary between games or remain constant?<br>
> 19x19 in total with the outermost borders being undestructable, leaving the players with a 17x17(289) grid to operate in.<br>


What algorithm will you use for random block placement to ensure a playable map?<br>
> In the grid you have to make sure that you do not fill in the squares next to the spawn area, leaving you with 17x17-8(281) squares to fill in.<br>
>There will also be indestructable blocks in the 17x17 grid, 8x8(64) them to be exact, leaving you with 281-64(217) blocks in total to be filled in with autogeneration.<br>
>We will autogenerate bricks in to the remaining squares one by one with a 50% chance.
>20% of these bricks will be turned in to a power up when destroyed.



How will you ensure that randomly placed blocks don't block players completely?<br>

## Networking:

How will game data be sent between clients and server? Which data needs to be sent and how frequently?<br>
How will you handle player synchronization to ensure fair gameplay?<br>
How will the game handle lag or network issues?<br>

## Chat Functionality:

How will the chat interface look and function?<br>
Will the chat be global or will there be options for private messages between players?<br>

## Game UI:

How will the start page, player count, and game interface look?<br>
How will the game inform players about the countdown before the game starts?<br>
How will you display player lives, current power-ups, and other game information on the screen?<br>
How can players quit the game or join a new one?<br>

## Sound and Graphics:

What kind of sound effects and music will the game have?<br>
How will the game visually represent different events like explosions, power-ups collection, player deaths, etc.?<br>

## Performance:

What strategies will you use to ensure the game runs at 60fps without frame drops?<br>
How will you measure and optimize performance?<br>
How will you ensure that the game runs smoothly even when there's a lot happening (e.g., multiple bombs exploding at once)?<br>
What will your strategy be for stress testing the game to ensure performance standards are met?<br>
