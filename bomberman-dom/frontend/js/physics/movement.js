// File for functions that has to do with moving a player
let adjustment = 0;
let width = window.innerWidth;

function debounce(func, wait) {
  let timeout;
  return function (...args) {
    clearTimeout(timeout);
    timeout = setTimeout(() => func.apply(this, args), wait);
  };
}
export function movePlayer(player) {
  console.log("movePlayer function called");
  console.log(width);
  let message;
  const handleResize = function () {
    console.log("resize event");
    console.log(width);
    console.log(window.innerWidth);
    adjustment = window.innerWidth - width;
    let player1 = document.getElementById("player-1");
    console.log(player1.style.left);
    console.log(adjustment);
    adjustment = adjustment / 2;

    width = window.innerWidth;
    // Update players' positions
    let players = document.getElementsByClassName("player");
    for (let i = 0; i < players.length; i++) {
      let player = players[i];
      let playerLeft = parseFloat(player.style.left);
      let newLeft = playerLeft + adjustment;
      player.style.left = `${newLeft}px`;
    }
    console.log(player1.style.left);
  };

  window.addEventListener("resize", debounce(handleResize, 200)); // 200ms delay

  document.addEventListener("keydown", function (event) {
    const chatInput = document.getElementById("chat-input");

    // If chat is focused, don't execute game commands
    if (document.activeElement === chatInput) {
      return;
    }
    // Prevent scrolling for Arrow keys and Space
    if (
      ["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight", "Space"].includes(
        event.code
      )
    ) {
      event.preventDefault();
    }
    switch (event.code) {
      case "ArrowLeft":
      case "KeyA":
        message = {
          command: "move",
          direction: "left",
        };
        break;
      case "ArrowUp":
      case "KeyW":
        message = {
          command: "move",
          direction: "up",
        };
        break;
      case "ArrowRight":
      case "KeyD":
        message = {
          command: "move",
          direction: "right",
        };
        break;
      case "ArrowDown":
      case "KeyS":
        //send info through websocket
        message = {
          command: "move",
          direction: "down",
          // Add other properties to the message as needed
        };
        break;
      case "Space":
        message = {
          command: "place-bomb",
          // Add other properties to the message as needed
        };
        break;
      default:
        return;
    }
    // Check if the WebSocket is still open before sending a message
    if (window.webSocketConnection.readyState === WebSocket.OPEN) {
      console.log("Sending message", message);
      // Send the message through the WebSocket connection
      window.webSocketConnection.send(JSON.stringify(message));
    } else {
      console.log("Can't send message, WebSocket connection is not open");
    }
  });

  function animationLoop() {
    window.webSocketConnection.onmessage = function (event) {
      let message = JSON.parse(event.data);

      if (message.type === "gameGrid") {
        console.log("gameGrid received");
      } else if (message.type === "player") {
        console.log("player received");
        updatePlayerPosition(message.data);
        updateLife(message.data);
      } else if (message.type === "bomb") {
        console.log("bomb received");
        updateBombPlacement(message.data);
      } else if (message.type === "explosion") {
        console.log("explosion received");
        updateExplosion(message.data);
      } else if (message.type === "chat") {
        console.log("chat message received");
        let chatWindow = document.getElementById("chat-messages");
        let cmessage = document.createElement("p");
        cmessage.classList.add("chat-message");
        cmessage.innerHTML = message.data;
        chatWindow.appendChild(cmessage);
      } else if (message.type === "gameOver") {
        console.log("gameOver received");
        console.log(message.data);

        const winnerName = message.data.Winner;
        const modalHTML = `
          <div id="myModal" class="modal">
            <div class="modal-content">
              <p id="winMessage">${winnerName} has won the game!</p>
              <button id="refreshButton">Go Home</button>
            </div>
          </div>
        `;

        // Append the modal HTML to the body of the document
        document.body.innerHTML += modalHTML;

        const modal = document.getElementById("myModal");
        const refreshButton = document.getElementById("refreshButton");

        modal.style.display = "block";

        refreshButton.addEventListener("click", function () {
          location.reload(); // Refresh the page when the button is clicked
        });
      }
    };

    // Call requestAnimationFrame to keep the animation loop running
    requestAnimationFrame(animationLoop);
  }

  animationLoop();
}

function updatePlayerPosition(player) {
  console.log("player: ", player);

  let playerxd = document.getElementById(`player-${player.PlayerID}`);
  let playerGridSquare = document.getElementById(
    `cell-${player.GridPosition[0]}-${player.GridPosition[1]}`
  );
  playerGridSquare.classList.remove(
    "speedy",
    "bombAmountIncrease",
    "bombRangeIncrease"
  );
  let gridPost = playerGridSquare.getBoundingClientRect();
  playerxd.style.left = `${gridPost.left}px`;
  playerxd.style.top = `${gridPost.top}px`;
  // Determine the player class based on the direction
  let direction = "down"; // Default player direction for moving downwards
  switch (player.Direction) {
    case "up":
      direction = "up";
      break;
    case "down":
      direction = "down";
      break;
    case "left":
      direction = "left";
      break;
    case "right":
      direction = "right";
      break;
    default:
      direction = "down"; // Default to "down" if the direction is unknown
      break;
  }

  //update the players background image based on the direction, player1 is down, playerLeft1 is left, playerRight1 is right, playerUpward1 is up (for player 1)
  let playerImg = document.getElementById(`player-${player.PlayerID}`);
  playerImg.style.backgroundImage = `url("/static/images/player${player.PlayerID}${direction}.png")`;

  //update the moving player when pressing the key with the gif image
  let playerGif = document.getElementById(`player-${player.PlayerID}`);
  playerGif.style.backgroundImage = `url("/static/images/player${player.PlayerID}${direction}.gif")`;

  // Remove all previous direction classes and add the current direction class
  playerxd.classList.remove(
    "playerup",
    "playerdown",
    "playerleft",
    "playerright"
  );
  playerxd.classList.add(`player-${direction}`);
}

function updateBombPlacement(bomb) {
  console.log("bomb: ", bomb);
  let bombPosition = document.getElementById(
    `cell-${bomb.GridPosition[0]}-${bomb.GridPosition[1]}`
  );

  bombPosition.classList.add("bomb");

  // Create an img element for the bomb image
  let bombImg = document.createElement("img");
  bombImg.src = "/static/images/bomb1.png";
  bombImg.classList.add("bomb");

  bombPosition.appendChild(bombImg);

  let isExploding = false;

  // Timer to switch between bomb images
  setInterval(() => {
    isExploding = !isExploding;
    bombImg.src = isExploding
      ? "/static/images/bomb3.png"
      : "/static/images/bomb1.png";
  }, 450);

  // Automatically remove the bomb after a predetermined time
  setTimeout(() => {
    bombPosition.removeChild(bombImg);
    bombPosition.classList.remove("bomb");
  }, bomb.ExplosionTime || 2000); // Assuming the explosion time is 2 seconds (2000ms). Adjust this value as needed.
}

function updateExplosion(explosion) {
  console.log("explosion: ", explosion);
  console.log("affected cells: ", explosion.AffectedCells);

  // Handles the explosion logic for a single cell, including board modifications
  function explodeCell(cell) {
    let explodedCell = document.getElementById(`cell-${cell[0]}-${cell[1]}`);

    if (explodedCell.classList.contains("brick")) {
      explodedCell.classList.remove("brick");
    }

    // Depending on the value in the GameGrid we can determine what type of powerup to spawn
    switch (explosion.GameGrid[cell[1]][cell[0]]) {
      case 8:
        explodedCell.classList.add("speedy");
        break;
      case 9:
        explodedCell.classList.add("bombAmountIncrease");
        break;
      case 10:
        explodedCell.classList.add("bombRangeIncrease");
        break;
      default:
        break;
    }
  }

  // Handles only the explosion animation for a cell
  function animateExplosion(cell, explodedCell) {
    // Create an img element for the explosion image
    const explosionImg = document.createElement("img");
    explosionImg.src = "/static/images/explosion1.png";
    explosionImg.classList.add("explosion-image");
    explodedCell.appendChild(explosionImg);

    let isExploding = false;
    const explosionImages = [
      "/static/images/explosion1.png",
      "/static/images/explosion2.png",
      "/static/images/explosion3.png",
    ];
    let index = 0;

    const switchImage = () => {
      isExploding = !isExploding;
      explosionImg.src = isExploding
        ? explosionImages[index]
        : "/static/images/explosion1.png";
      index++;

      if (index < explosionImages.length) {
        setTimeout(switchImage, 200);
      } else {
        setTimeout(() => {
          explodedCell.removeChild(explosionImg);
          explodedCell.classList.remove("explosion");
          explodedCell.classList.add("cell");
        }, 500);
      }
    };

    setTimeout(switchImage, 0); // Start the explosion animation immediately

    // Apply the explosion animation
    explodedCell.classList.add("explosion");
    setTimeout(() => {
      explodedCell.classList.remove("explosion");
      explodedCell.classList.add("cell");
    }, 500);
  }

  if (explosion.AffectedCells && explosion.AffectedCells.length > 0) {
    explosion.AffectedCells.forEach(explodeCell);
  }

  if (explosion.ExplosionCells && explosion.ExplosionCells.length > 0) {
    explosion.ExplosionCells.forEach((cell) => {
      const targetCell = document.getElementById(`cell-${cell[0]}-${cell[1]}`);
      animateExplosion(cell, targetCell);
    });
  }
}

// function to remove a life from the player if they get hit by an explosion
function updateLife(player) {
  // set the player lives to the number of lives they have

  // Find the player's life element
  console.log("playerID: ", player.PlayerID);
  const playerLives = document.getElementById(`player${player.PlayerID}-lives`);
  // get the amount of children in the player's life element
  const playerLivesCount = playerLives.childElementCount;

  if (playerLivesCount > player.Lives) {
    // Remove the last child of the player's life element
    const lastHeart = playerLives.querySelector(".heart:last-child");

    if (lastHeart) {
      playerLives.removeChild(lastHeart);
    }
    if (player.Lives <= 0) {
      console.log("player died");
      let playerEl = document.querySelector("#player-" + player.PlayerID);
      console.log("playerEl: ", playerEl);
      // remove player from board
      playerEl.remove();
    }
  }
}
