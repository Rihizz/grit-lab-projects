function debounceThrottle(func, delay) {
  let timeout = null;
  let lastRun = 0;
  return function (...args) {
    const context = this;
    const elapsed = Date.now() - lastRun;
    const runCallback = function () {
      lastRun = Date.now();
      func.apply(context, args);
    };
    clearTimeout(timeout);
    if (elapsed >= delay) {
      runCallback();
    } else {
      timeout = setTimeout(runCallback, delay - elapsed);
    }
  };
}
let isTyping = false;
const THROTTLE_INTERVAL = 1000; // milliseconds
let chatOpenWithUser = null;
let OffSet = 10;
let currentUser = null;

const chatSocket = new WebSocket("ws://localhost:80/ws");

document.addEventListener("DOMContentLoaded", function () {
  getUser().then(function (currentUsername) {
    getAndShowUsers().then(function (userElements) {
      attachUserElementClickListeners(userElements, currentUser);
      // const userElements = document.querySelectorAll(".user");
      // const chatSocket = new WebSocket("ws://localhost:80/ws");
      setTimeout(() => {
        updateShit();
      }, 10);

      currentUsername = currentUser;

      // Add event listeners to handle WebSocket events
      chatSocket.addEventListener("open", function (event) {
        console.log("WebSocket connection established");
      });

      chatSocket.addEventListener("message", function (event) {
        console.log("WebSocket message received: " + event.data);
        if (event.data === "USER_JOINED") {
          console.log("USER_JOINED");
          setTimeout(() => {
            updateShit();
          }, 10);
        } else if (event.data === "USER_LEFT") {
          console.log("USER_LEFT");
          updateShit();
        } else if (event.data === "USER_CREATED") {
          console.log("USER_CREATED");
          getAndShowUsers().then(function (newUserElements) {
            userElements = newUserElements; // Update the userElements variable
            attachUserElementClickListeners(userElements, currentUser);
            setTimeout(() => {
              updateShit();
            }, 10);
          });
        } else {
          // Parse the message data
          const data = JSON.parse(event.data);
          console.log(data.command);
          const chatModal = document.getElementById("chat-modal");
          if (data.command === "NEW_MESSAGE") {
            if (data.sender !== currentUsername) {
              console.log("sender is not current user");
              moveUserToTop(data.sender);
            }
            if (
              chatModal !== null &&
              currentUsername === data.receiver &&
              data.sender === chatOpenWithUser
            ) {
              const messagesElement = chatModal.querySelector(".messages");
              console.log(data.Timestamp);
              messagesElement.innerHTML += `<p class="timestamp">${data.Timestamp}<p><strong>${data.sender}</strong>: ${data.text}</p>`;
              messagesElement.scrollTop = messagesElement.scrollHeight;
            } else {
              // If the chat modal is not open, make the senders user element blink
              if (data.sender !== currentUsername) {
                const userElementsArray = Array.from(userElements);
                const senderUserElement = userElementsArray.find(
                  (userElement) => {
                    const username = userElement.querySelector("a").innerText;
                    return username === data.sender;
                  }
                );

                if (senderUserElement) {
                  senderUserElement.classList.add("blink");
                }
              }
              if (data.sender === currentUsername) {
                const messagesElement = chatModal.querySelector(".messages");
                messagesElement.innerHTML += `<p class="timestamp">${data.Timestamp}<p><strong>${data.sender}</strong>: ${data.text}</p>`;
                messagesElement.scrollTop = messagesElement.scrollHeight;
              }
            }
          } else if (data.command === "TYPING") {
            console.log(data.sender, "is typing to", data.receiver);
            console.log(chatModal);
            const typingIndicator = document.getElementById(
              `typing-${data.sender}`
            );
            console.log(typingIndicator);
            typingIndicator.style.display = "flex";
            if (chatModal !== null) {
              console.log("chat modal is open");
              const typingIndicator2 = document.getElementById(
                `typing2-${data.sender}`
              );
              console.log(typingIndicator2);
              typingIndicator2.style.display = "inline";
            }
          } else if (data.command === "STOP_TYPING") {
            console.log(data.sender, "stopped typing to", data.receiver);
            const typingIndicator = document.getElementById(
              `typing-${data.sender}`
            );
            typingIndicator.style.display = "none";
            if (chatModal !== null) {
              const typingIndicator2 = document.getElementById(
                `typing2-${data.sender}`
              );
              typingIndicator2.style.display = "none";
            }
          }
        }
      });

      chatSocket.addEventListener("error", function (event) {
        console.log("WebSocket error occurred: " + event.error);
      });

      chatSocket.addEventListener("close", function (event) {
        console.log("WebSocket connection closed");
      });
    });
  });
});

function loadChatHistory1(currentUsername, username, messagesElement) {
  // Send a request to the server to get the chat history for this user
  fetch(`/api/get-messages?receiver=${username}&sender=${currentUsername}`)
    .then((response) => response.json())
    .then((messages) => {
      // Check if there are any messages
      if (messages.length === 0) {
        messagesElement.innerHTML = `<p class="no-messages">No messages yet</p>`;
        return;
      }
      // Add the messages to the chat display
      messages.reverse().forEach((message) => {
        const messageHTML = `
              <div class="message">
                <p class="timestamp">${message.Timestamp}</p>
                <p><strong>${message.Sender}</strong>: ${message.Text}</p>
              </div>
            `;
        messagesElement.insertAdjacentHTML("beforeend", messageHTML);
      });

      // Scroll to the bottom of the chat display
      messagesElement.scrollTop = messagesElement.scrollHeight;
    })
    .catch((error) => {
      console.error(`Error loading chat history: ${error}`);
    });
}

function loadChatMoreHistory1(currentUsername, username, messagesElement) {
  // Send a request to the server to get the chat history for this user
  let scrollHeight = messagesElement.scrollHeight;
  fetch(
    `/api/get-messages?receiver=${username}&sender=${currentUsername}&offset=${OffSet}`
  )
    .then((response) => response.json())
    .then((messages) => {
      // Add the messages to the chat display
      messages.forEach((message) => {
        const messageHTML = `
                <div class="message">
                  <p class="timestamp">${message.Timestamp}</p>
                  <p><strong>${message.Sender}</strong>: ${message.Text}</p>
                </div>
              `;
        messagesElement.insertAdjacentHTML("afterbegin", messageHTML);
      });
      // make sure the scroll position is the same as before
      messagesElement.scrollTop = messagesElement.scrollHeight - scrollHeight;
    })
    .catch((error) => {
      console.error(`Error loading chat history: ${error}`);
    });
  OffSet += 10;
}

async function getUser() {
  try {
    const response = await fetch("/api/current-user", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });
    const data = await response.json();
    if (data != false) {
      JSON.stringify(data);
    }
    currentUser = data;
    return data;
  } catch (error) {
    console.error(`Error getting username: ${error}`);
  }
}

function parseDateSafely(dateString) {
  if (!dateString) return null;
  const sanitizedDateString = dateString.replace(" ", "T");
  const date = new Date(sanitizedDateString);
  return date;
}

async function getAndShowUsers() {
  try {
    const currentUsername = await getUser();
    const response = await fetch(`/api/users`);
    const users = await response.json();

    if (users == null || users.length == 0) {
      return;
    }

    // Fetch the last messages for each user
    const lastMessages = await Promise.all(
      users.map(async (user) => {
        const messages = await fetch(
          `/api/get-messages?receiver=${user.Username}&sender=${currentUsername}`
        ).then((response) => response.json());
        return messages[0];
      })
    );

    // Create an array of users with their last messages
    const usersWithLastMessages = users.map((user, index) => ({
      ...user,
      lastMessage: lastMessages[index],
    }));

    // Sort the users based on the timestamp of the last messages
    usersWithLastMessages.sort((a, b) => {
      const aTimestamp = a.lastMessage
        ? parseDateSafely(a.lastMessage.Timestamp)
        : null;
      const bTimestamp = b.lastMessage
        ? parseDateSafely(b.lastMessage.Timestamp)
        : null;

      if (!aTimestamp && !bTimestamp) return 0;
      if (!aTimestamp) return 1;
      if (!bTimestamp) return -1;
      return bTimestamp - aTimestamp;
    });

    const usersContainer = document.querySelector(".users");
    usersContainer.innerHTML = "";
    usersContainer.innerHTML = generateUsersHTML(usersWithLastMessages);

    // Return the new user elements
    return document.querySelectorAll(".user");
  } catch (error) {
    console.error(error);
  }
}

function generateUsersHTML(users) {
  return users
    .map(
      (user) => `
      <div class="user">
      <div class="status-dot" id="${user.Username}"></div>
      <a href="#" aria-label="link to chat with ${user.Username}">${user.Username}</a>
      <div class="typing" id="typing-${user.Username}" style="display: none;"><div class="dot"></div><div class="dot"></div><div class="dot"></div></span>
 
</div>
    </div>
        `
    )
    .join("");
}

function moveUserToTop(username) {
  const userElements = document.querySelectorAll(".user");
  const usersContainer = document.querySelector(".users");

  for (const userElement of userElements) {
    if (userElement.querySelector("a").innerText === username) {
      usersContainer.insertBefore(userElement, usersContainer.firstChild);
      break;
    }
  }
}

function attachUserElementClickListeners(userElements, currentUsername) {
  // Add a click event listener to each user element
  userElements.forEach(function (userElement) {
    userElement.addEventListener("click", function (event) {
      // Prevent the default link click behavior
      event.preventDefault();
      // Get the username from the user element
      const username = userElement.querySelector("a").innerText;
      chatOpenWithUser = username;

      // remove blink class if it exists
      userElement.classList.remove("blink");

      const typing = document.getElementById(`typing-${username}`);
      console.log(typing.style.display);
      let xd = null;
      if (typing.style.display == "flex") {
        console.log("typing");
        xd = "inline";
      } else {
        xd = "none";
      }

      // Create a new chat modal element
      const chatModal = document.createElement("div");
      chatModal.id = "chat-modal";
      chatModal.innerHTML = `
         <div class="modal-content">
           <span class="close">&times;</span>
           <h3>Chatting with <span id="chat-username">${username}</span><span class="typing" id="typing2-${username}" style="display: ${xd};"><div class="dot"></div><div class="dot"></div><div class="dot"></div></span></h3>
           <div class="messages"></div>
           <div class="input-box">
             <input type="text" id="message-input" placeholder="Type your message..." maxlength="100" required>
           </div>
         </div>
       `;

      // Add the chat modal element to the page
      document.body.appendChild(chatModal);

      // Add a click event listener to the close button
      const closeButton = chatModal.querySelector(".close");
      closeButton.addEventListener("click", function () {
        isTyping = false;
        chatModal.remove();
        OffSet = 10;
        chatSocket.send(
          JSON.stringify({
            command: "STOP_TYPING",
            receiver: username,
            sender: currentUsername,
          })
        );
      });

      // Show the chat modal
      chatModal.style.display = "block";

      const messagesElement = chatModal.querySelector(".messages");

      // load chat history with debounce and throttle
      const throttledLoadChatMoreHistory1 = debounceThrottle(() => {
        loadChatMoreHistory1(currentUsername, username, messagesElement);
      }, THROTTLE_INTERVAL);

      // Load the initial chat history
      loadChatHistory1(currentUsername, username, messagesElement);

      // Add a scroll event listener to the chat display
      messagesElement.addEventListener("scroll", () => {
        if (messagesElement.scrollTop === 0) {
          throttledLoadChatMoreHistory1();
        }
      });

      // Add an event listener for when the user presses a key
      const messageInputElement = chatModal.querySelector("#message-input");
      messageInputElement.addEventListener("keyup", function (event) {
        // if the user pressed a key other than the enter key and the input field is not empty
        if (
          event.keyCode !== 13 &&
          messageInputElement.value.trim() !== "" &&
          !isTyping
        ) {
          isTyping = true;
          // Send a typing event to the server
          chatSocket.send(
            JSON.stringify({
              command: "TYPING",
              receiver: username,
              sender: currentUsername,
            })
          );
        } else if (isTyping && messageInputElement.value.trim() === "") {
          isTyping = false;
          // Send a stop typing event to the server
          chatSocket.send(
            JSON.stringify({
              command: "STOP_TYPING",
              receiver: username,
              sender: currentUsername,
            })
          );
        }
        // If the user pressed the enter key
        if (event.keyCode === 13) {
          isTyping = false;
          if (messageInputElement.value.trim() !== "") {
            // Send the message
            chatSocket.send(
              JSON.stringify({
                command: "NEW_MESSAGE",
                text: messageInputElement.value,
                receiver: username,
              })
            );
            chatSocket.send(
              JSON.stringify({
                command: "STOP_TYPING",
                receiver: username,
                sender: currentUsername,
              })
            );
          } else {
            alert("Please enter a message");
          }
          moveUserToTop(username);
          // Clear the input field
          messageInputElement.value = "";
        }
      });
    });
  });
}
