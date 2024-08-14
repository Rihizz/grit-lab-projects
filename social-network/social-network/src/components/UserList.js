import React, { useState, useEffect } from "react";
import { useAuth } from "../AuthContext.js";
import useFetchUserList from "./FetchuserList.js";
import { Link } from "react-router-dom";
import ChatModal from "./ChatModal.js"; // Import the ChatModal component

function UserListItem({ user }) {
  const { chatWS, userID } = useAuth();
  const [showChat, setShowChat] = useState(false);
  const [newMessage, setNewMessage] = useState(false); // State for new message received

  const handleOpenChat = () => {
    setShowChat(true);
    setNewMessage(false); // Reset new message flag when chat is opened
  };

  const handleCloseChat = () => {
    setShowChat(false);
  };

  if (chatWS) {
    chatWS.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.receiver_id === userID && data.sender_id !== userID) {
        setNewMessage(true);
      }
    };
  }

  return (
    <div className="user-list-item">
      <Link to={`/profile/${user.id}`}>{user.fullname}</Link>
      <button
        onClick={handleOpenChat}
        className={`chat-button ${newMessage ? "blink" : ""}`}
      >
        <i className="fa-solid fa-comments"></i>
      </button>
      {showChat && <ChatModal user={user} onClose={handleCloseChat} />}
    </div>
  );
}

function UserList() {
  const { data: userlist, loading, error } = useFetchUserList();

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  // Check if userlist is not null or undefined and if it's empty, display "No contacts yet"
  if (userlist && userlist.length === 0) {
    return <div className="user-list-item">No contacts yet</div>;
  }

  return (
    <div>
      {/* Ensure userlist is not null or undefined before calling map */}
      {userlist &&
        userlist.map((user) => <UserListItem key={user.id} user={user} />)}
    </div>
  );
}

export default UserList;
