import React, { useState, useContext, createContext, useEffect } from "react";
import ErrorPage from "./pages/ErrorPage";

const AuthContext = createContext();

export function useAuth() {
  return useContext(AuthContext);
}

export function AuthProvider({ children }) {
  const [ws, setWs] = useState(null);
  const [loggedIn, setLoggedIn] = useState(false);
  const [userID, setUserID] = useState(null);
  const [nickname, setNickname] = useState(null);
  const [chatWS, setChatWS] = useState(null);
  const [error, setError] = useState(null);

  // Function to check login status from the API
  async function checkLoginStatus() {
    try {
      const response = await fetch("http://localhost:6969/api/check-login", {
        method: "GET",
        credentials: "include",
      });

      if (response.ok) {
        const data = await response.json();
        setLoggedIn(true);
        setUserID(data.userID);
        //storing user id in session storage
        sessionStorage.setItem("userID", data.userID);
        setNickname(data.nickname);

        // Create a WebSocket connection when the user is logged in
        if (!ws) {
          const newWebSocket = new WebSocket("ws://localhost:6969/ws");
          setWs(newWebSocket);
        }
        if (!chatWS) {
          const newChatWebSocket = new WebSocket("ws://localhost:6969/chat");
          setChatWS(newChatWebSocket);
        }
      } else {
        setLoggedIn(false);
        setUserID(null);
        setNickname(null);
      }
    } catch (error) {
      setLoggedIn(false);
      setUserID(null);
      setNickname(null);
      console.error("Error checking login status:", error);
      setError(500);
    }
  }

  async function logout() {
    try {
      const response = await fetch("http://localhost:6969/api/logout", {
        method: "POST",
        credentials: "include",
      });

      if (response.ok) {
        setLoggedIn(false);
        setUserID(null);
        setNickname(null);

        // Close the WebSocket connection when the user logs out
        if (ws) {
          ws.close();
          setWs(null);
        }
      }
    } catch (error) {
      console.error("Error logging out:", error);
    }
  }

  // Call the checkLoginStatus function when the component mounts
  useEffect(() => {
    checkLoginStatus();
  }, []);

  const value = {
    loggedIn,
    setLoggedIn,
    userID,
    nickname,
    checkLoginStatus,
    logout,
    ws,
    chatWS,
  };

  if (error) {
    console.log("Error in AuthContext:", error);
    return <ErrorPage errorType={error}/>;
  }
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
