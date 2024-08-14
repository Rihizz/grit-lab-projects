// Description: This file contains the functions for the movement of the pacMan and the ghosts

export { movePacMan, ghostMove };

import { wallArray } from "./maze.js";
//window size
let windowWidth = window.innerWidth;
let windowHeight = window.innerHeight;

// for the pacMan
let audioplayink = false;
let timerInterval;
let xd = 0;
let pause = false;
let pacManLeft = 400;
let pacManTop = 690;
let movePacManLeft = 0;
let movePacManTop = 0;
let running = false;
let gameOver = false;

function movePacMan() {
  displayStartMenu();
  // control the pacMan with the arrow keys
  document.addEventListener("keydown", (event) => {
    let counter = 0;
    if (counter === 0) {
      movePacManLeft = 0;
      movePacManTop = 0;
      counter++;
    }
    if (pause) {
      if (event.code === "Escape") {
        if (!gameOver) {
          pause = false;
          running = true;
          hidePauseMenu();
          ghostMove(ghost);
          xd = 0;
          return;
        }
      } else {
        window.location.reload();
      }
      return;
    }
    running = true;
    if (gameOver === false) {
      switch (event.code) {
        case "ArrowLeft":
        case "KeyA": // left arrow
          movePacManLeft = -4;
          pacMan.style.rotate = "0deg";
          pacMan.style.transform = "scaleX(-1)";
          running = true;
          break;
        case "ArrowUp":
        case "KeyW": // up arrow
          movePacManTop = -4;
          pacMan.style.rotate = "270deg";
          pacMan.style.transform = "scaleY(-1)";
          running = true;
          break;
        case "ArrowRight":
        case "KeyD": // right arrow
          movePacManLeft = 4;
          pacMan.style.transform = "scaleX(1)";
          pacMan.style.rotate = "0deg";
          running = true;
          break;
        case "ArrowDown":
        case "KeyS": // down arrow
          movePacManTop = 4;
          pacMan.style.rotate = "90deg";
          pacMan.style.transform = "scaleY(1)";
          running = true;
          break;
        case "Escape": // escape key
          if (running) {
            displayPauseMenu();
            pause = true;
            running = false;
            break;
          }
        default:
          running = true;
          break;
      }
    }
  });

  function animatePacMan() {
    const jailArray = getJailArray();
    if (xd === 0 && running) {
      const timerDisplay = document.querySelector("#timer");
      timerInterval = setInterval(() => {
        timer++;
        const minutes = Math.floor(timer / 60);
        const seconds = timer % 60;
        timerDisplay.innerHTML = `${minutes < 10 ? "0" + minutes : minutes}:${
          seconds < 10 ? "0" + seconds : seconds
        }`;
      }, 1000);
      xd++;
    }
    if (!running) {
      clearInterval(timerInterval);
    }
    let lastMove = "";
    if (running) {
      pacManLeft += movePacManLeft;
      pacManTop += movePacManTop;
      if (checkBounds(pacManLeft, pacManTop)) {
        let newCoords = checkBounds(pacManLeft, pacManTop);
        pacManLeft = newCoords[0];
        pacManTop = newCoords[1];
      }
      if (checkGhostCollision()) {
        livesGone();
        gameOver = true;
        var audio = new Audio("sounds/xdxdxdxdxd.mp3");
        audio.play();
        let ghostLeft = parseInt(ghost.style.left);
        let ghostTop = parseInt(ghost.style.top);
        running = false;
        const xplode = document.createElement("div");
        document.body.appendChild(xplode);
        xplode.id = "xplos";
        xplode.style.left = ghostLeft - 175 + "px";
        xplode.style.top = ghostTop - 75 + "px";
        xplode.style.width = "200px";
        xplode.style.height = "200px";
        xplode.style.position = "absolute";
        xplode.style.zIndex = "1";
        xplode.style.boxSizing = "border-box";
        xplode.style.transform = "rotate(0deg)";
        xplode.style.transition = "transform 0.1s";

        const xplodeImg = document.createElement("img");
        xplodeImg.src = "pictures/xddd.gif";
        xplodeImg.style.objectFit = "cover";
        xplodeImg.style.borderRadius = "50%";
        xplodeImg.style.width = "400px";
        xplodeImg.style.height = "200px";
        xplode.appendChild(xplodeImg);
        background.style.backgroundImage = "url('pictures/bolis.jpeg')";

        setTimeout(() => {
          displayGameOverMenu();
          // reset the game
          pacManLeft = 450;
          pacManTop = 600;
          movePacManLeft = 0;
          movePacManTop = 0;
          lastMove = "";
          running = false;
          //location.reload();
        }, 1000);
      }
      if (newCheckWall(pacManLeft, pacManTop, wallArray)) {
        let newCoords = newCheckWall(pacManLeft, pacManTop, wallArray);
        pacManLeft = newCoords[0];
        pacManTop = newCoords[1];
      }
      if (newCheckWall(pacManLeft, pacManTop, jailArray)) {
        let newCoords = newCheckWall(pacManLeft, pacManTop, jailArray);
        pacManLeft = newCoords[0];
        pacManTop = newCoords[1];
      }
      if (checkTeleport1(pacManLeft, pacManTop)) {
        let newCoords = [775, 575];
        pacManLeft = newCoords[0];
        pacManTop = newCoords[1];
      }
      if (checkTeleport2(pacManLeft, pacManTop)) {
        let newCoords = [35, 575];
        pacManLeft = newCoords[0];
        pacManTop = newCoords[1];
      }
      if (checkXd(pacManLeft, pacManTop)) {
        var audio69 = new Audio("sounds/xd.mp3");
        if (!audioplayink) {
          audio69.play();
          background.style.backgroundImage = "url('pictures/boom.png')";
          background.style.backgroundSize = "cover";
          background.style.backgroundRepeat = "no-repeat";
          audioplayink = true;
        }
      }
      var audio4 = new Audio("sounds/winer.mp3");
      if (snusCollision()) {
        audio4.play();
        background.style.backgroundImage = "url('pictures/åland.jpeg')";
        setTimeout(() => {
          youWin();
          // reset the game
          pacManLeft = 450;
          pacManTop = 600;
          movePacManLeft = 0;
          movePacManTop = 0;
          lastMove = "";
          snusCount = 0;
          running = false;
        }, 250);
      }

      pacMan.style.left = pacManLeft + "px";
      pacMan.style.top = pacManTop + "px";
      if (movePacManLeft !== 0) {
        lastMove = movePacManLeft > 0 ? "right" : "left";
      } else if (movePacManTop !== 0) {
        lastMove = movePacManTop > 0 ? "down" : "up";
      }
    }
    requestAnimationFrame(animatePacMan);
  }
  animatePacMan();
}

function displayStartMenu() {
  let benis = true;
  document.body.innerHTML += `
    <div id="start-menu">
      <h1>Pacman xd</h1>
      <p>Press any key to start</p>
      <p>Movement arrow keys or WASD</p>
      <p>Pause your movement with another key to lure the ghost</p>
      <p>Pause the game with escape</p>
      <p>Collect 10 pucks of Odens to win!</p>
      <p>Don't let the customs man catch you</p>
    </div>
  `;
  document.addEventListener("keydown", function (event) {
    if (event.code && benis) {
      hideStartMenu();
      benis = false;
    }
  });
}

function displayPauseMenu() {
  if (!gameOver) {
    document.body.innerHTML += `
      <div id="pause-menu">
        <p>Game Paused</p>
        <button id="resume-button">Resume</button>
        <button id="restart-button">Restart</button>
      </div>
    `;

    document
      .getElementById("resume-button")
      .addEventListener("click", function () {
        hidePauseMenu();
        ghostMove(ghost);
        pause = false;
        running = true;
        xd = 0;
      });

    document
      .getElementById("restart-button")
      .addEventListener("click", function () {
        window.location.reload();
      });
  } else {
    window.location.reload();
  }
}

function hideStartMenu() {
  ghostMove(ghost);
  document.getElementById("start-menu").remove();
}

function hidePauseMenu() {
  document.getElementById("pause-menu").remove();
}

function displayGameOverMenu() {
  document.body.innerHTML += `
    <div id="game-over-menu">
      <h1>Game Over</h1>
      <p>Press any key to restart</p>
    </div>
  `;
  document.addEventListener("keydown", function (event) {
    if (event.code === "Escape") {
      hideGameOverMenu();
      window.location.reload();
      pause = true;
      running = false;
    } else {
      hideGameOverMenu();
      window.location.reload();
      pause = true;
      running = false;
    }
  });
}

function hideGameOverMenu() {
  document.getElementById("game-over-menu").remove();
}

function youWin() {
  document.body.innerHTML += `
    <div id="you-win-menu">
      <h1>You Win!</h1>
      <p>You brought all the snus back to Åland!</p>
      <p>Press any key to go again. There is always more snus to collect</p>
    </div>
  `;
  document.addEventListener("keydown", function (event) {
    if (event.code === "Escape") {
      hideYouWinMenu();
      window.location.reload();
      pause = true;
      running = false;
    } else {
      hideYouWinMenu();
      window.location.reload();
      pause = true;
      running = false;
    }
  });
}

function hideYouWinMenu() {
  document.getElementById("you-win-menu").remove();
}

function ghostMove(ghost) {
  let ghostLeft = parseInt(ghost.style.left);
  let ghostTop = parseInt(ghost.style.top);
  function animateGhost() {
    if (!pause) {
      if (running) {
        let { ghostLeft: newLeft, ghostTop: newTop } = moveGhostLogic(
          pacManLeft,
          pacManTop,
          ghostLeft,
          ghostTop
        );
        ghostLeft = newLeft;
        ghostTop = newTop;
        ghost.style.left = ghostLeft + "px";
        ghost.style.top = ghostTop + "px";
      }
    }
    requestAnimationFrame(animateGhost);
  }
  animateGhost();
}
let lives = 1;
let snusCount = 0;
let previousIndex = 0;
let timer = 0;
const locations = [
  { x: 320, y: 250 },
  { x: 320, y: 890 },
  { x: 117, y: 250 },
  { x: 117, y: 890 },
  { x: 650, y: 250 },
  { x: 445, y: 250 },
  { x: 650, y: 890 },
  { x: 445, y: 890 },
  { x: 650, y: 570 },
  { x: 117, y: 570 },
  { x: 230, y: 790 },
  { x: 530, y: 790 },
  { x: 530, y: 350 },
  { x: 230, y: 350 },
];

function livesGone() {
  if (!gameOver) {
    const livesDisplay = document.querySelector(".lives");
    lives--;
    livesDisplay.innerHTML = "Lives: " + lives;
  }
}

function snusCollision() {
  const counterDisplay = document.querySelector("#counter");
  const pacman = {
    x: pacManLeft,
    y: pacManTop,
    width: 35,
    height: 35,
  };
  var audio3 = new Audio("sounds/satii.mp3");

  const snus = document.querySelector(".snus"); // get snus element
  const snusRect = snus.getBoundingClientRect();
  if (
    pacman.x < snusRect.right &&
    pacman.x + pacman.width > snusRect.left &&
    pacman.y < snusRect.bottom &&
    pacman.y + pacman.height > snusRect.top
  ) {
    // collision detected!
    snusCount++;
    counterDisplay.innerHTML = "Score: " + snusCount;
    if (snusCount === 10) {
      gameOver = true;
      snus.style.display = "none";
      return true;
    } else {
      let randomIndex = previousIndex;
      while (randomIndex === previousIndex) {
        randomIndex = Math.floor(Math.random() * locations.length);
      }
      audio3.play();
      // move snus to a random location from the list of 10 locations
      previousIndex = randomIndex;
      const location = locations[randomIndex];
      snus.style.left = location.x + "px";
      snus.style.top = location.y + "px";
      return false;
    }
  }
}

// check that the charachters are within the game board
function checkBounds(placementLeft, placementTop) {
  if (placementLeft < 0) {
    return [0, placementTop];
  } else if (placementLeft + 50 > windowWidth) {
    return [windowWidth - 50, placementTop];
  } else if (placementTop < 0) {
    return [placementLeft, 0];
  } else if (placementTop + 50 > windowHeight) {
    return [placementLeft, windowHeight - 50];
  } else {
    return false;
  }
}

function checkTeleport1(placementLeft, placementTop) {
  const pacman = {
    x: placementLeft,
    y: placementTop,
    width: 35,
    height: 35,
  };
  const teleport1 = document.querySelector(".teleport1"); // get teleport1 element
  const teleport1Rect = teleport1.getBoundingClientRect();
  if (pacman.x < teleport1Rect.right && pacman.x > teleport1Rect.left) {
    return true;
  }
  return false;
}
function checkTeleport2(placementLeft, placementTop) {
  const pacman = {
    x: placementLeft,
    y: placementTop,
    width: 35,
    height: 35,
  };
  const teleport2 = document.querySelector(".teleport2"); // get teleport2 element
  const teleport2Rect = teleport2.getBoundingClientRect();
  if (pacman.x < teleport2Rect.right && pacman.x > teleport2Rect.left) {
    return true;
  }
  return false;
}

function checkXd(placementLeft, placementTop) {
  const pacman = {
    x: placementLeft,
    y: placementTop,
    width: 35,
    height: 35,
  };
  const teleport2 = document.querySelector(".teleport2"); // get teleport2 element
  const teleport2Rect = teleport2.getBoundingClientRect();
  if (pacman.x < teleport2Rect.right + 25 && pacman.x > teleport2Rect.left) {
    return true;
  }
  return false;
}

// checks if the pacMan collides with a wall and returns new coordinates next to the wall if it does
function newCheckWall(placementLeft, placementTop, wallArray) {
  for (const element of wallArray) {
    let wallRect = element.getBoundingClientRect();
    if (
      placementLeft < wallRect.x + wallRect.width &&
      placementLeft + 35 > wallRect.x &&
      placementTop < wallRect.y + wallRect.height &&
      35 + placementTop > wallRect.y
    ) {
      if (placementLeft < wallRect.x) {
        placementLeft = wallRect.x - 35;
      } else if (placementLeft + 35 > wallRect.x + wallRect.width) {
        placementLeft = wallRect.x + wallRect.width;
      } else if (placementTop < wallRect.y) {
        placementTop = wallRect.y - 35;
      } else if (placementTop + 35 > wallRect.y + wallRect.height) {
        placementTop = wallRect.y + wallRect.height;
      }
    }
  }
  return [placementLeft, placementTop];
}

// check if the pacMan has collided with a ghost
function checkGhostCollision() {
  const pacMan = document.getElementById("pacMan");
  const ghost = document.getElementById("ghost");
  let pacManRect = pacMan.getBoundingClientRect();
  let ghostRect = ghost.getBoundingClientRect();
  return !!(
    pacManRect.x < ghostRect.x + ghostRect.width &&
    pacManRect.x + pacManRect.width > ghostRect.x &&
    pacManRect.y < ghostRect.y + ghostRect.height &&
    pacManRect.height + pacManRect.y > ghostRect.y
  );
}

// getting the jail array
function getJailArray() {
  let jailArray = [];
  let jailElements = document.getElementsByClassName("jail");
  for (const element of jailElements) {
    jailArray.push(element);
  }
  let jailDoorElements = document.getElementsByClassName("jailDoor");
  for (const element of jailDoorElements) {
    jailArray.push(element);
  }
  return jailArray;
}

function moveGhostLogic(pacManLeft, pacManTop, ghostLeft, ghostTop) {
  let opposite = pacManTop - ghostTop;
  let adjacent = pacManLeft - ghostLeft;
  let angle = Math.atan(opposite / adjacent);
  if (ghostLeft > pacManLeft) {
    angle = angle + Math.PI;
  }
  let velocity = 2;

  let vx = velocity * Math.cos(angle);
  let vy = velocity * Math.sin(angle);

  ghostLeft = ghostLeft + vx;
  ghostTop = ghostTop + vy;
  return { ghostLeft, ghostTop };
}
