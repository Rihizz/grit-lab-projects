// description: This file contains the logic for building the game board and the game elements

export { createPacMan, createGhost, createSnus, wallArray };

// initial possition for the pacMan
let circleLeft = 391;
let circleTop = 690;

// initial for the ghost
let ghostLeft = 391;
let ghostTop = 580;

let wallArray = document.getElementsByClassName("wall");

// create the pacMan element
function createPacMan() {
  const pacMan = document.createElement("div");
  document.body.appendChild(pacMan);
  pacMan.id = "pacMan";
  pacMan.style.left = circleLeft + "px";
  pacMan.style.top = circleTop + "px";
  pacMan.style.width = "35px";
  pacMan.style.height = "35px";
  pacMan.style.borderRadius = "50%";
  pacMan.style.position = "absolute";
  pacMan.style.zIndex = "1";
  pacMan.style.boxSizing = "border-box";
  pacMan.style.transform = "rotate(0deg)";
  pacMan.style.transition = "transform 0.1s";

  const pacManImg = document.createElement("img");
  pacManImg.src = "pictures/autism.png";
  pacManImg.style.width = "35px";
  pacManImg.style.height = "35px";
  pacMan.appendChild(pacManImg);
}

// create the ghost element
function createGhost() {
  const ghost = document.createElement("div");
  document.body.appendChild(ghost);
  ghost.id = "ghost";
  ghost.style.left = ghostLeft + "px";
  ghost.style.top = ghostTop + "px";
  ghost.style.width = "35px";
  ghost.style.height = "35px";
  ghost.style.borderRadius = "50%";
  ghost.style.position = "absolute";
  ghost.style.zIndex = "1";
  ghost.style.boxSizing = "border-box";
  ghost.style.transform = "rotate(0deg)";
  ghost.style.transition = "transform 0.1s";

  const ghostImg = document.createElement("img");
  ghostImg.src = "pictures/borderguard.jpeg";
  ghostImg.style.width = "35px";
  ghostImg.style.height = "35px";
  ghost.appendChild(ghostImg);
}

function createSnus() {
  const snus = document.createElement("div");
  document.body.appendChild(snus);
  snus.id = "snus";
  snus.className = "snus";
  snus.style.left = "320px";
  snus.style.top = "250px";
  snus.style.width = "50px";
  snus.style.height = "50px";
  snus.style.borderRadius = "50%";
  snus.style.position = "absolute";
  snus.style.zIndex = "1";
  snus.style.boxSizing = "border-box";
  snus.style.transform = "rotate(0deg)";
  snus.style.transition = "transform 0.1s";

  const snusImg = document.createElement("img");
  snusImg.src = "pictures/odns.png";
  snusImg.style.width = "50px";
  snusImg.style.height = "45px";
  snus.appendChild(snusImg);
}
