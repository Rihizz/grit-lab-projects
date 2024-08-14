## BACKEND

The game board is a grid of 869x869 pixels. If we can just use the pixel grid instead of making our own 190x190 grid we can save time.
We can keep track of players, blocks, powerups and anything else by keeping track of their cordinates.

At the initialization of the game board we can calculate the position of weadther or not the cell we are moving to is a brick.
this is calculated by the size of the cells(px) and position relative to origo.