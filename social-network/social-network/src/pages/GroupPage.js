import React, { useState, useEffect } from "react";
import GroupPosts from "../components/GroupPosts";
import "../styles/Groups.css";
import EventContainer from "../components/EventContainer";
import { useAuth } from "../AuthContext"; // import useAuth from AuthContext
import GroupChatModal from "../components/GroupChat";
import { Link } from "react-router-dom";

const fetchGroupData = async (groupNumber) => {
  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  };

  const response = await fetch(
    `http://localhost:6969/api/serve-group-data?id=${groupNumber}`,
    requestOptions
  );

  const data = await response.json();
  if (response.status === 200) {
    // console.log("group data fetched");
    // console.log(data);
    return data;
  } else {
    return response.status;
  }
};
const fetchEventData = async (groupNumber) => {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ groupId: groupNumber }),
  };

  const response = await fetch(
    `http://localhost:6969/api/serve-events`,
    requestOptions
  );

  const data = await response.json();
  if (response.status === 200) {
    // console.log("events data fetched");
    // console.log(data);
    return data.events;
  } else {
    return response.status;
  }
};

const joinGroup = async (groupNumber) => {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      group_id: groupNumber,
    }),
    credentials: "include",
  };
  const response = await fetch(
    "http://localhost:6969/api/join-group",
    requestOptions
  );

  if (response.status === 200) {
    console.log("group joined");
    return "group joined"; // Return success message
  } else {
    console.log("error joining group");
    alert("Error joining group.");
  }
};

const GroupPage = () => {
  const { userID, ws, nickname } = useAuth(); // Get the userID
  const [followedUsers, setFollowedUsers] = useState([]);
  const url = window.location.href;
  const pattern = /groups\/(\d+)/;
  const match = url.match(pattern);
  const [selectedUser, setSelectedUser] = useState("");
  const groupNumber = match[1];
  const [groupData, setGroupData] = React.useState([]);
  const [eventsData, setEventsData] = React.useState([]);
  const [groupMembers, setGroupMembers] = React.useState([]);
  const [newEvent, setNewEvent] = React.useState({
    title: "",
    description: "",
    dateTime: "",
  });
  const [errorMessage, setErrorMessage] = React.useState("");

  const [isChatOpen, setIsChatOpen] = useState(false);

  const handleOpenChat = () => {
    setIsChatOpen(true);
  };

  const handleCloseChat = () => {
    setIsChatOpen(false);
  };

  useEffect(() => {
    // Replace the URL with the correct endpoint for fetching followed users
    fetch(`http://localhost:6969/api/followed_users?group=${groupNumber}`, {
      method: "GET",
      credentials: "include",
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error("Failed to fetch followed users");
        }
      })
      .then((data) => setFollowedUsers(data))
      .catch((error) => console.error(error));
  }, [userID, groupNumber]);

  React.useEffect(() => {
    const getGroupData = async () => {
      const groupDataFromServer = await fetchGroupData(groupNumber);
      // console.log("groupDataFromServer", groupDataFromServer);
      setGroupData(groupDataFromServer);
      setGroupMembers(groupDataFromServer.members);
      // console.log("groupDataFromServer.members", groupDataFromServer.members);
    };
    getGroupData();

    const getEventData = async () => {
      const eventsDataFromServer = await fetchEventData(groupNumber);
      if (eventsDataFromServer === 403) {
        setEventsData([]);
      }
      setEventsData(eventsDataFromServer);
    };

    getEventData();
  }, [groupNumber]);

  if (!groupData) {
    return <div>loading...</div>;
  }

  const handleEventChange = (e) => {
    setNewEvent({
      ...newEvent,
      [e.target.name]: e.target.value,
    });
  };

  const handleEventSubmit = async (e) => {
    //check that event title and description do not consist of only spaces
    if (
      newEvent.title.trim().length === 0 ||
      newEvent.description.trim().length === 0
    ) {
      alert("Event title and description cannot be empty.");
      return;
    }
    e.preventDefault();
    const response = await fetch("http://localhost:6969/api/event", {
      method: "POST",
      credentials: "include", // to send the session cookie
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        group_id: groupNumber,
        title: newEvent.title,
        description: newEvent.description,
        date_time: newEvent.dateTime,
      }),
    });

    // If response is OK, re-fetch the events
    if (response.ok) {
      const eventsDataFromServer = await fetchEventData(groupNumber);
      setEventsData(eventsDataFromServer);

      // console.log("Event created successfully");
      // console.log(groupData.members, "group members");

      // Send notification to all group members
      groupData.members.forEach((member) => {
        const notificationData = {
          type: "new event",
          groupid: groupData.id,
          userid: Number(member.id), // assuming member object has 'id' property
          sender: userID,
          message: `${nickname} has created a new event ${newEvent.title} in group ${groupData.name}`,
        };

        ws.send(
          JSON.stringify({
            command: "NOTIFICATION",
            message: JSON.stringify(notificationData),
          })
        );
      });
    } else {
      setErrorMessage("Error creating event");
    }
    setNewEvent({
      title: "",
      description: "",
      dateTime: "",
    });
  };

  const handleGroupJoin = async () => {
    const response = await joinGroup(groupNumber);
    if (response === "group joined") {
      console.log("group joined");
      let notification;
      notification = {
        type: "group request",
        groupid: groupData.id,
        userid: groupData.creator_id,
        sender: userID,
        message: nickname + " requested to join your group " + groupData.name,
      };

      ws.send(
        JSON.stringify({
          command: "NOTIFICATION",
          message: JSON.stringify(notification),
        })
      );
    }
  };

  const groupCreator = groupData.creator_id;

  function handleInviteUser() {
    const inviteData = {
      type: "group invite",
      groupid: groupData.id,
      userid: Number(selectedUser),
      sender: userID,
      message: nickname + " invited you to join " + groupData.name,
    };

    ws.send(
      JSON.stringify({
        command: "NOTIFICATION",
        message: JSON.stringify(inviteData),
      })
    );
    alert("User has been invited to the group.");
  }

  return (
    <div className="group-page">
      <div className="group-page-header-inside">
        <h2>{groupData.name}</h2>
        <p>{groupData.description}</p>
        {groupData.access === false && (
          <>
            <button className="join-group-button" onClick={handleGroupJoin}>
              Request access
            </button>
            <p>
              You are not authorized to view this group. Please request access
              from the group creator.
            </p>
          </>
        )}
        {groupData.access && (
          <>
            <div className="group-chat-modal">
              <button className="group-button" onClick={handleOpenChat}>
                <i className="fa-solid fa-comments"></i> Open Groupchat
              </button>

              {isChatOpen && (
                <GroupChatModal group={groupNumber} onClose={handleCloseChat} />
              )}
            </div>
            {/* <h1>Group Invite</h1> */}
            <div className="group-page-invite">
              <select onChange={(event) => setSelectedUser(event.target.value)}>
                <option value="">Select a user to invite</option>
                {followedUsers && followedUsers.length > 0 ? (
                  followedUsers.map((user) => (
                    <option key={user.id} value={user.id}>
                      {user.fullname}
                    </option>
                  ))
                ) : (
                  <option value="">No users to invite</option>
                )}
              </select>
              <button className="invite-button" onClick={handleInviteUser}>
                Invite to group
              </button>
            </div>
          </>
        )}
      </div>

      {groupData.access && (
        <>
          <div className="group-page-members">
            <h2>Creator of the group:</h2>
            <div>
              <p>
                {groupCreator === userID ? (
                  <Link to={`/profile`}>{nickname}</Link>
                ) : (
                  <Link to={`/profile/${groupCreator}`}>
                    {groupMembers[0].full_name}
                  </Link>
                )}
              </p>
            </div>

            <h2>Members:</h2>
            {groupMembers.map((member) => (
              <div key={member.id}>
                <p>
                  {member.id === userID ? (
                    <Link to={`/profile`}>{member.full_name}</Link>
                  ) : (
                    <Link to={`/profile/${member.id}`}>{member.full_name}</Link>
                  )}
                </p>
              </div>
            ))}
          </div>

          {/* <h1>Group Events</h1> */}
          <div className="group-page-event">
            <form onSubmit={handleEventSubmit}>
              <div className="group-page-event-form">
                <h2>Create event</h2>
                <input
                  className="event-input"
                  type="text"
                  name="title"
                  placeholder="Title"
                  value={newEvent.title}
                  onChange={handleEventChange}
                  required
                  min={1}
                  max={50}
                  title="Event title should be 1-50 characters."
                />
                <input
                  className="event-description"
                  type="text"
                  name="description"
                  placeholder="Description"
                  value={newEvent.description}
                  onChange={handleEventChange}
                  required
                  min={1}
                  max={256}
                  title="Event description should be 1-256 characters."
                />
                <input
                  className="event-date"
                  type="datetime-local"
                  name="dateTime"
                  value={newEvent.dateTime}
                  onChange={handleEventChange}
                  required
                  min={new Date().toISOString().substring(0, 16)}
                />
                <button className="create-event-button" type="submit">
                  Create event
                </button>
              </div>
            </form>
            {errorMessage && <p className="error-message">{errorMessage}</p>}
          </div>

          <div className="group-page-event">
            <h2>Upcoming Events:</h2>
            <EventContainer
              groupId={groupNumber}
              userID={userID}
              eventsData={eventsData}
            />
          </div>

          <div className="group-page-post">
            <h1>Group Posts</h1>
            <div className="post-display">
              <GroupPosts />
            </div>
          </div>
        </>
      )}
    </div>
  );
};

export default GroupPage;
