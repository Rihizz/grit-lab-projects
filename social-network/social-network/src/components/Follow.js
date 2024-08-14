import React, { useEffect } from "react";
import { useAuth } from "../AuthContext";

function Follow({ userId, privacy }) {
  const { userID, ws, nickname } = useAuth();

  const handleFollow = async () => {
    // console.log("handleFollow called");
    // console.log("userId", userId);
    // console.log("userID", userID);
    // console.log("privacy", privacy);
    const requestOptions = {
      method: "POST",

      headers: {
        "Content-Type": "application/json",
      },

      body: JSON.stringify({
        followee: userId,
        follower: userID,
        privacy: privacy,
      }),

      credentials: "include",
    };

    const response = await fetch(
      "http://localhost:6969/api/follow",
      requestOptions
    );

    if (response.ok) {
      // console.log("followed");
    } else {
      // console.log("follow failed");
    }

    let notification;

    if (privacy === "0") {
      notification = {
        type: "follow",
        userid: userId,
        sender: userID,
        message: nickname + " followed you",
      };
    } else {
      notification = {
        type: "follow request",
        userid: userId,
        sender: userID,
        message: nickname + " requested to follow you",
      };
    }

    ws.send(
      JSON.stringify({
        command: "NOTIFICATION",
        message: JSON.stringify(notification),
      })
    );
  };

  useEffect(() => {
    handleFollow();
  }, [userId, userID]);

  return <></>;
}

export default Follow;
