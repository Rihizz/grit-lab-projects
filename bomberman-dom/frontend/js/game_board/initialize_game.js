// file to handle the sequence of events that occur when the game is initialized
//
// 1. the user clicks the start button
// 2. the start button is removed from the DOM
// 3. a websocket connection is established
// 4. a websocket message is sent to the server
// 5. the server responds with the game board
// 6. the game board is rendered to the DOM
// 7. when multiple users are connected, the game board is updated when the server sends a message
// 8. timer starts
// 9. when the timer reaches 0, the game starts
// 10. the game ends when there is only one player left
// 11. the game board is cleared from the DOM
// 12. the game over screen is rendered to the DOM
import { buildBaseGrid } from "./build_base.js";
import { movePlayer } from "../physics/movement.js";

let timer = 20; // set to 20 for production
let timer2 = 10; // set to 10 for production

let playersConnected = 0;
let timerInterval;
let players = [];
let starting = false;

function updateLobbyDisplay() {
  let lobby = document.getElementById("lobby");
  if (!lobby) {
    // Create lobby if it doesn't exist
    lobby = document.createElement("div");
    lobby.id = "lobby";
    document.body.appendChild(lobby);
  }

  lobby.innerHTML = makeLobbyHTML();
}

function startCountdown() {
  clearInterval(timerInterval); // Clear existing interval if any
  timerInterval = setInterval(() => {
    timer--;

    if (timer <= 0) {
      clearInterval(timerInterval);
      // Start the final countdown
      starting = true;
      updateLobbyDisplay();
      finalCountdown();
    }

    updateLobbyDisplay();
  }, 1000);
}

function finalCountdown() {
  clearInterval(timerInterval); // Clear existing interval if any
  timerInterval = setInterval(() => {
    timer2--;

    if (timer2 <= 0) {
      clearInterval(timerInterval);
      // Transition to start the game
      console.log("start game called");
      startGame();
    }

    updateLobbyDisplay();
  }, 1000);
}

function stopAndResetTimer() {
  clearInterval(timerInterval); // Stops any active timer

  if (starting) {
    timer2 = 10; // reset the second timer if the final countdown was in progress
    starting = false;
  }
  timer = 20; // reset the first timer

  updateLobbyDisplay();
}

async function initializeGame() {
  createChatWindow();
  // remove the start button from the DOM
  let startButton = document.getElementById("start-button");
  startButton.remove();
  let username = document.getElementById("player-name").value;
  let welcomeMessage = document.getElementById("welcome-message");
  welcomeMessage.remove();

  let message = {
    command: "player",
    name: username,
  };

  // establish a websocket connection
  window.webSocketConnection = new WebSocket("ws://localhost:80/ws");
  // send a message to the server
  window.webSocketConnection.onopen = function (event) {
    console.log("WebSocket is open now.");
    window.webSocketConnection.send(JSON.stringify(message));
  };
  window.webSocketConnection.onmessage = async function (event) {
    let receivedMessage = JSON.parse(event.data);
    if (receivedMessage.type === "start") {
      let gameBoard = receivedMessage.data;
      let num = receivedMessage.playerlist.length;
      updateHUD(receivedMessage.playerlist);
      await buildBaseGrid(gameBoard, num);
      setTimeout(function () {
        requestAnimationFrame(movePlayer);
      }, 1000);
    } else if (receivedMessage.type === "player-connected") {
      console.log("player connected");
      playersConnected = receivedMessage.numberOfConns;
      players = receivedMessage.playerlist;
      updateName(players);
      console.log(receivedMessage.playerlist);
      console.log(players);
      updateLobbyDisplay();
      if (playersConnected >= 2) {
        startCountdown();
      }
      if (playersConnected === 4) {
        starting = true;
        updateLobbyDisplay();
        finalCountdown();
      }
    } else if (receivedMessage.type === "player-disconnected") {
      console.log("player disconnected");
      playersConnected = receivedMessage.numberOfConns;
      players = receivedMessage.playerlist;
      updateLobbyDisplay();
      if (playersConnected < 2) {
        stopAndResetTimer();
      }
    } else if (receivedMessage.type === "chat") {
      console.log("chat message received");
      let chatWindow = document.getElementById("chat-messages");
      let message = document.createElement("p");
      message.classList.add("chat-message");
      message.innerHTML = receivedMessage.data;
      chatWindow.appendChild(message);
    } else if (receivedMessage.type === "abandon") {
      console.log("abandon message received");
      window.location.href =
        window.location.origin +
        window.location.pathname +
        "?message=gameAbandoned";
    }
  };
}

function makeLobbyHTML() {
  if (!starting) {
    return `
      <div>
          <h1>Waiting for players...</h1>
          <h2>Players connected: ${playersConnected}</h2>
          <h2>Players: ${players}</h2>
          <h2>Game starts in: ${timer} seconds</h2>
      </div>
      `;
  } else {
    return `
        <div>
            <h1>Game starting...</h1>
            <h2>Players connected: ${playersConnected}</h2>
            <h2>Players: ${players}</h2>
            <h2>Game starts in: ${timer2} seconds</h2>
        </div>
        `;
  }
}

function startGame() {
  console.log("inside start game");
  // Hide lobby display
  let lobby = document.getElementById("lobby");
  console.log(lobby);
  if (lobby) {
    lobby.style.display = "none";
  }
  // check that the websocket connection is open
  if (window.webSocketConnection.readyState === WebSocket.OPEN) {
    // send a message to the server
    let message = {
      command: "start",
    };
    window.webSocketConnection.send(JSON.stringify(message));
  }
}

function updateHUD(players) {
  for (let i = 0; i < players.length; i++) {
    let username = players[i];
    let playerNameElement = document.getElementById(`player${i + 1}-name`);
    playerNameElement.innerHTML = username;
  }

  // Hide HUDS for players that are not in the game
  for (let i = players.length; i < 4; i++) {
    let playerElement = document.getElementById(`player${i + 1}`);
    playerElement.style.display = "none";
  }
}

function updateName(players) {
  for (let i = 0; i < players.length; i++) {
    let username = players[i];
    let playerNameElement = document.getElementById(`player${i + 1}-name`);
    playerNameElement.innerHTML = username;
  }
}

function createChatWindow() {
  const chatWindow = document.createElement("div");
  chatWindow.id = "chat-window";

  const chatMessages = document.createElement("div");
  chatMessages.id = "chat-messages";
  chatWindow.appendChild(chatMessages);

  const chatInputContainer = document.createElement("div");
  chatInputContainer.id = "chat-input-container";
  chatWindow.appendChild(chatInputContainer);

  const chatInput = document.createElement("input");
  chatInput.type = "text";
  chatInput.id = "chat-input";
  chatInput.placeholder = "Type a message...";
  chatInputContainer.appendChild(chatInput);

  const sendButton = document.createElement("button");
  sendButton.id = "send-button";
  sendButton.textContent = "Send";
  chatInputContainer.appendChild(sendButton);

  function sendMessage() {
    // If there's a message in the input, display it in the chat messages area
    if (chatInput.value.trim() !== "") {
      // check that the websocket connection is open and send the message
      if (window.webSocketConnection.readyState === WebSocket.OPEN) {
        // send a message to the server
        let message = {
          command: "chat",
          message: chatInput.value,
        };
        window.webSocketConnection.send(JSON.stringify(message));
      }
      // Clear the input after sending
      chatInput.value = "";
      // Scroll to the bottom to see the latest message
      chatMessages.scrollTop = chatMessages.scrollHeight;
    }
  }

  // Send message when the "Send" button is clicked
  sendButton.addEventListener("click", sendMessage);

  // Send message when "Enter" key is pressed in the chat input
  chatInput.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
      event.preventDefault(); // Prevent the default action (new line or form submission)
      sendMessage();
    }
  });

  // Append the chat window to the body (or another parent element if you prefer)
  document.body.appendChild(chatWindow);
}

export { initializeGame };
