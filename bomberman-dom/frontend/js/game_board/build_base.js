// File for building the base grid for the game board

// Function to build the base grid for the game board
async function buildBaseGrid(gameBoard, playerAmount) {
  // Generate CSS for grid template
  const n = gameBoard.length;
  const gridTemplate = `repeat(${n}, 0fr)`;

  // Create style element
  const style = document.createElement("style");
  style.textContent = `
        #game-board {
            display: grid;
            grid-template-columns: ${gridTemplate};
            grid-template-rows: ${gridTemplate};
            justify-content: center;
        }
        .cell {
            height: 50px;
            width: 50px;
            background-image: url("/static/images/grass.png");
            background-size: cover;
        }
        .edge {
            background-image: url("/static/images/backgroundbrick1.png");
            background-size: cover;
        }
        .steel {
            background-image: url("/static/images/backgroundbrick1.png");
        }
        .speedy {
          background-image: url("/static/images/speedy.png");
        }
        .bombAmountIncrease {
          background-image: url("/static/images/bombAmountIncrease.png");
        }
        .bombRangeIncrease {
          background-image: url("/static/images/powerup.png");
        }
        .safe-zone {
            background-image: url("/static/images/grass.png");
            background-size: cover;
        }
        .brick {
            background-position: center;
            background-size: cover;
            background-image: url("/static/images/brick2.png");
        }
        #player-1 {
            background-image: url("/static/images/player1down.png");
            background-size: cover;
            position: absolute;
            z-index: 1;
            width: 50px;
            height: 50px;
            transition: left 0.3s ease, top 0.3s ease;
        }
        #player-2 {
            background-image: url("/static/images/player2.png");
            background-size: cover;
            position: absolute;
            z-index: 1;
            width: 50px;
            height: 50px;
            transition: left 0.3s ease, top 0.3s ease;
        }
        #player-3 {
            background-image: url("/static/images/player3.png");
            background-size: cover;
            z-index: 1;
            position: absolute;
            width: 50px;
            height: 50px;
            transition: left 0.3s ease, top 0.3s ease;
        }
        #player-4 {
            background-image: url("/static/images/player4.png");
            background-size: cover;
            z-index: 1;
            position: absolute;
            width: 50px;
            height: 50px;
            transition: left 0.3s ease, top 0.3s ease;
        }
    `;

  // Append style element to document head
  document.head.appendChild(style);

  // Generate grid cells based on the gameBoard data from the server
  for (let y = n - 1; y >= 0; y--) {
    for (let x = 0; x < n; x++) {
      let cell = document.createElement("div");

      // Set cell class based on gameBoard value
      switch (gameBoard[y][x]) {
        case 0:
          cell.classList.add("cell");
          break;
        case 1:
          cell.classList.add("cell", "brick");
          break;
        case 2:
          cell.classList.add("cell", "steel");
          break;
        case 3:
          cell.classList.add("cell", "edge");
          break;
        default:
          cell.classList.add("cell");
          break;
      }

      cell.setAttribute("id", `cell-${x}-${y}`);
      document.getElementById("game-board").appendChild(cell);
    }
  }

  const startcell1 = document.getElementById("cell-1-1");
  const startcell2 = document.getElementById("cell-17-17");
  const startcell3 = document.getElementById("cell-1-17");
  const startcell4 = document.getElementById("cell-17-1");

  for (let i = 1; i <= playerAmount; i++) {
    let player = document.createElement("div");
    player.setAttribute("id", `player-${i}`);
    player.classList.add("player");
    // let playerStart = document.getElementsByClassName(`starting-cell-${i}`);
    let playerStart;
    switch (i) {
      case 1:
        playerStart = startcell1;
        break;
      case 2:
        playerStart = startcell2;
        break;
      case 3:
        playerStart = startcell3;
        break;
      case 4:
        playerStart = startcell4;
        break;
    }
    let startRect = playerStart.getBoundingClientRect();
    player.style.top = `${startRect.top}px`;
    player.style.left = `${startRect.left}px`;
    document.getElementById("game-board").appendChild(player);
  }
}

export { buildBaseGrid };
