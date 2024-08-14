import React, { useState, useEffect } from "react";
import "../styles/ChatModal.css";
import { useAuth } from "../AuthContext";

async function fetchGroupChatHistory(group) {
  // console.log("goup in fetch: ", group);
  const requestOptions = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ group }),
  };

  const response = fetch(
    "http://localhost:6969/api/serve-group-messages",
    requestOptions
  );

  const data = await (await response).json();
  // console.log("data: ", data.messages);
  return data.messages;
}

function GroupChatModal({ group, onClose }) {
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);
  const { chatWS, nickname, userID } = useAuth();
  const [selectedImage, setSelectedImage] = useState(null);
  // console.log("group: ", group);

  useEffect(() => {
    setMessage("");
    async function fetchData() {
      let groupChatHistory = await fetchGroupChatHistory(group);
      setMessages(groupChatHistory);
      // console.log("messages: ", messages);
    }
    fetchData();
    if (chatWS) {
      chatWS.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.receiver === group && data.command === "GROUP_MESSAGE")
          setMessages((prevMessages) => [
            ...prevMessages,
            {
              content: data.message,
              sender: data.sender,
              image: data.image,
              created_at: new Date(),
            },
          ]);
      };
      // console.log(messages);
    }

    return () => {
      if (chatWS) {
        chatWS.onmessage = null;
      }
    };
  }, [chatWS]);

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
       // Check if the message is empty
      if (!message.trim()) {
        alert("Message cannot be empty.");
        return;
      }
    if (chatWS.readyState === WebSocket.OPEN) {
      // If an image has been selected
      if (selectedImage) {
        // console.log("selectedImage: ");
        const reader = new FileReader();
        reader.onload = function (event) {
          const imgBase64String = event.target.result;

          const payload = {
            receiver: group,
            sender: nickname,
            sender_id: userID,
            Command: "GROUP_MESSAGE",
            message,
            image: imgBase64String, // Send the image as an ArrayBuffer
          };

          // Send the message over the WebSocket connection
          chatWS.send(JSON.stringify(payload));
        };

        reader.readAsDataURL(selectedImage);
      } else {
        // console.log("no image");
        // If no image selected, send message normally
        const payload = {
          receiver: group,
          sender: nickname,
          sender_id: userID,
          Command: "GROUP_MESSAGE",
          message,
        };

        // Send the message over the WebSocket connection
        chatWS.send(JSON.stringify(payload));
      }
    }

    setMessage("");
    setSelectedImage(null);

    // clear the file input value after posting the message
    const fileInput = document.getElementById("image");
    if (fileInput) {
      fileInput.value = "";
    }
  };

  //send the message when the user presses enter
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
    <div className="chat-modal" ref={modalRef}>
      <button className="close-x" onClick={onClose}>
        X
      </button>
      <h2>Group Chat</h2>
      <div className="chat-messages">
        {messages.map((message, index) => (
          <div key={index} className="chat-message">
            <div className="message-sender">{message.sender}</div>
            <div className="message-text">{message.content}</div>
            {message.image && (
              <img src={`data:image/jpeg;base64,${message.image}`} alt="alt" />
            )}
            <div className="message-timestamp">
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
        <button onClick={handleSendMessage}>Send</button>
        <br />
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

export default GroupChatModal;
