import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../AuthContext";
import Notifications from "./Notifications";
import SearchBar from "./SearchBar";

function Navbar() {
  const { loggedIn, logout, ws, userID } = useAuth();
  const navigate = useNavigate();
  const [notifications, setNotifications] = useState([]);

  async function handleLogout() {
    await logout();
    navigate("/");
  }

  // Fetch notifications on component mount
  useEffect(() => {
    async function fetchNotifications() {
      if (!userID) {
        // if userID is not available, don't try to fetch notifications
        return;
      }

      // console.log("Fetching notifications for userid: " + userID);
      const requestOptions = {
        method: "GET",
        headers: { "Content-Type": "application/json" },
        credentials: "include", // include credentials in the request
      };

      const response = await fetch(
        "http://localhost:6969/api/notifications?id=" + userID,
        requestOptions
      );

      if (response.ok) {
        const data = await response.json();
        setNotifications(data);
        // console.log(data);
      } else {
        // handle error
        console.error("Failed to fetch notifications");
      }
    }

    fetchNotifications();
  }, [userID]); // userID as a dependency to useEffect

  useEffect(() => {
    if (!ws) {
      return;
    }

    ws.onmessage = function (event) {
      const newNotification = JSON.parse(event.data);

      if (newNotification.Command === "NOTIFICATION") {
        console.log("Received new notification:", newNotification);
        if (newNotification.UserId === userID) {
          setNotifications((prevNotifications) => {
            console.log("Setting new notifications", {
              prevNotifications,
              newNotification,
            });
            let updatedNotifications = [
              ...(Array.isArray(prevNotifications) ? prevNotifications : []),
              newNotification,
            ];
            console.log("Updated notifications:", updatedNotifications);
            return updatedNotifications;
          });
        }
      }
    };
  }, [ws]);

  return (
    <nav className="navbar">
      <ul className="topnav">
        <li className="home">
          <i className="fa-solid fa-house fa-lg"></i>
          <Link to="/">Home</Link>
        </li>
        <li className="about">
          <i className="fa-solid fa-circle-info fa-lg"></i>
          <Link to="/about">About</Link>
        </li>

        {loggedIn ? (
          <>
            <li className="groups">
              <i className="fa-solid fa-user-group fa-lg"></i>
              <Link to="/groups">Groups</Link>
            </li>
            <li className="profile">
              <i className="fa-solid fa-circle-user fa-lg"></i>
              <Link to="/profile">Profile</Link>
            </li>
          </>
        ) : (
          <li className="login-nav">
            <i className="fa-solid fa-face-grin fa-lg"></i>
            <Link to="/login">Register/Login</Link>
          </li>
        )}
      </ul>

      {loggedIn && (
        <>
          <li className="searchbar">
            <SearchBar />
          </li>
          <ul className="bottomnav">
            {/* <li className="notifications"> */}
            <Notifications notifications={notifications} />
            {/* </li> */}

            <div className="logout-button">
              <i className="fa-solid fa-face-frown-open fa-lg"></i>
              <button className="bals" onClick={handleLogout}>
                Logout
              </button>
            </div>
          </ul>
        </>
      )}
    </nav>
  );
}

export default Navbar;
