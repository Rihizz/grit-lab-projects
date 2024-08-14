import React, { useState, useEffect } from "react";
import "../styles/ChatModal.css";
import { useAuth } from "../AuthContext";

async function fetchChatHistory(receiver_id) {
  const requestOption = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include", // send the cookie along with the request
    body: JSON.stringify({ receiver_id }),
  };

  const response = await fetch(
    "http://localhost:6969/api/serve-private-messages",
    requestOption
  );

  const data = await response.json();
  // console.log("data: ", data.messages);
  return data;
}

function ChatModal({ user, onClose }) {
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const { chatWS, nickname, userID } = useAuth();
  const [selectedImage, setSelectedImage] = useState(null);
  // console.log(userID);

  useEffect(() => {
    async function fetchData() {
      const history = await fetchChatHistory(user.id);
      // console.log(history);
      setMessages(history);
    }
    if (chatWS) {
      chatWS.onmessage = (event) => {
        const data = JSON.parse(event.data);
        // console.log("data: ", data);
        if (data.sender_id === user.id || data.receiver_id === user.id) {
          setMessages((prevMessages) => [
            ...prevMessages,
            {
              content: data.message,
              sender_name: data.sender,
              created_at: new Date(),
              image: data.image,
            },
          ]);
        }
      };
    }

    fetchData();
  }, [chatWS, user.id]);

  const modalRef = React.useRef();

  useEffect(() => {
    const handleOutsideClick = (event) => {
      if (!modalRef.current.contains(event.target)) {
        onClose();
      }
    };

    document.addEventListener("mousedown", handleOutsideClick);

    return () => {
      document.removeEventListener("mousedown", handleOutsideClick);
    };
  }, [onClose]);

  const handleInputChange = (event) => {
    setMessage(event.target.value);
  };

  const handleImageChange = (e) => {
    setSelectedImage(e.target.files[0]);
  };

  const handleSendMessage = () => {
    if (!message.trim()) {
      // If the message is empty or contains only spaces
      alert("Message cannot be empty.");
      return;
    }

    // console.log("sender: ", nickname);
    // console.log("receiver: ", user.fullname);
    // Here you would send the message over the WebSocket connection
    if (chatWS.readyState === WebSocket.OPEN) {
      // If an image has been selected
      if (selectedImage) {
        const reader = new FileReader();
        reader.onload = function (event) {
          const imgBase64String = event.target.result;

          const payload = {
            receiver: user.fullname,
            sender: nickname,
            receiver_id: user.id,
            sender_id: userID,
            Command: "NEW_MESSAGE",
            message,
            image: imgBase64String, // Send the image as an ArrayBuffer
          };

          // Send the message over the WebSocket connection
          chatWS.send(JSON.stringify(payload));
        };

        reader.readAsDataURL(selectedImage);
      } else {
        // If no image selected, send message normally
        const payload = {
          receiver: user.fullname,
          sender: nickname,
          receiver_id: user.id,
          sender_id: userID,
          Command: "NEW_MESSAGE",
          message,
        };

        // Send the message over the WebSocket connection
        chatWS.send(JSON.stringify(payload));
      }
    }

    // For now, we just add it to the messages array
    // setMessages((prevMessages) => [
    //   ...prevMessages,
    //   { text: message, from: "me" },
    // ]);
    setMessage("");
    setSelectedImage(null);

    // clear the file input value after posting the message
    const fileInput = document.getElementById("image");
    if (fileInput) {
      fileInput.value = "";
    }
  };

  //send the message when the user presses the enter key
  const handleKeyDown = (event) => {
    if (event.key === "Enter") {
      handleSendMessage();
    }
  };

  //make sure the scroll position is the same as before
  useEffect(() => {
    const chatMessages = document.querySelector(".chat-messages");
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }, [messages]);

  //format the timestamp to be more readable
  const formatTimestamp = (dateTime) => {
    const options = {
      day: "numeric",
      month: "long",
      hour: "numeric",
      minute: "numeric",
    };
    return new Date(dateTime).toLocaleString(undefined, options);
  };

  return (
    <div className="chat-modal" ref={modalRef} id={user.id}>
      <button className="close-x" onClick={onClose}>
        X
      </button>
      <h2>Chat with {user.fullname}</h2>
      <div className="chat-messages">
        {messages.map((message, index) => (
          // console.log(message),
          <div key={index} className={`chat-message ${message.sender_id}`}>
            <div className="chat-username">{message.sender_name}</div>
            <div className="chat-message-text">{message.content}</div>
            <div className="message-image">
              {message.image && (
                <img
                  src={`data:image/jpeg;base64,${message.image}`}
                  alt="alt"
                />
              )}
            </div>
            <div className="chat-timestamp">
              {formatTimestamp(message.created_at)}
            </div>
          </div>
        ))}
      </div>
      <div className="chat-input">
        <input
          className="text-input"
          type="text"
          value={message}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
          placeholder="Type here..."
        />
        <button className="send-button" onClick={handleSendMessage}>
          Send
        </button>

        <div className="chat-image-preview">
          <input
            type="file"
            id="image"
            accept="image/*"
            onChange={handleImageChange}
          />
        </div>
      </div>
    </div>
  );
}

export default ChatModal;
