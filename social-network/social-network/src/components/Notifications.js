import React, { useState, useEffect } from "react";

function NotificationItem({ notification, onDelete }) {
  const handleAccept = async () => {
    // console.log("Accept request from", notification.userid);
    // console.log("trying to accept follow request", notification);
    const response = await fetch(
      `http://localhost:6969/api/notif_respond?notif_id=${notification.NotifId}&response=accept&notif_type=${notification.Type}&sender_id=${notification.SenderId}&group_id=${notification.Groupid}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      }
    );

    if (response.ok) {
      // console.log(`accepted`);
      onDelete(notification.NotifId);
    } else {
      console.error(`Error accepting follow request`);
    }
  };

  const handleDelete = async () => {
    const response = await fetch(
      `http://localhost:6969/api/notif_respond?notif_id=${notification.NotifId}&response=delete&notif_type=${notification.Type}&sender_id=${notification.SenderId}&group_id=${notification.Groupid}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      }
    );

    if (response.ok) {
      // console.log(`deleted`);
    } else {
      console.error(`Error deleting notification`);
    }
    onDelete(notification.NotifId);
  };

  const handleDecline = async () => {
    // console.log("Decline request from", notification.userid);
    const response = await fetch(
      `http://localhost:6969/api/notif_respond?notif_id=${notification.NotifId}&response=decline&notif_type=${notification.Type}&sender_id=${notification.SenderId}&group_id=${notification.Groupid}`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      }
    );

    if (response.ok) {
      // console.log(`declined`);
      onDelete(notification.NotifId);
    } else {
      console.error(`Error declining follow request`);
    }
  };

  // console.log("notification", notification);

  if (notification.Type === "follow request") {
    // console.log("follow request");
  }

  return (
    <li>
      {notification.Content}
      {notification.Type === "follow request" ||
      notification.Type === "group request" ||
      notification.Type === "group invite" ? (
        <div className="acceptdecline">
          <button className="icon-button acceptfollower" onClick={handleAccept}>
            <i className="fa-solid fa-check fa-xl"></i>
          </button>
          <button
            className="icon-button declinefollower"
            onClick={handleDecline}
          >
            <i className="fa-solid fa-xmark fa-xl"></i>
          </button>
        </div>
      ) : (
        <button className="icon-button delete-button" onClick={handleDelete}>
          <i className="fa-solid fa-trash fa-xl"></i>
        </button>
      )}
    </li>
  );
}

function Notifications({ notifications }) {
  const [isOpen, setIsOpen] = useState(false);
  const [notifList, setNotifList] = useState(notifications);

  const toggleOpen = () => setIsOpen(!isOpen);

  const handleDeleteNotification = (notifId) => {
    setNotifList((prevNotifList) =>
      prevNotifList.filter((n) => n.NotifId !== notifId)
    );
  };

  useEffect(() => {
    setNotifList(notifications);
  }, [notifications]);

  const modalRef = React.useRef();
  //close modal if clicked on outside of element
  useEffect(() => {
    function handleClickOutside(event) {
      if (modalRef.current && !modalRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [modalRef]);

  return (
    <div className="notificationsXD" ref={modalRef}>
      <div onClick={toggleOpen}>
        Notifications{" "}
        {notifList && notifList.length > 0 && <span>({notifList.length})</span>}
      </div>

      {isOpen && (
        <ul className="notifications-list">
          {notifList && notifList.length > 0 ? (
            notifList.map((notification, i) => (
              <NotificationItem
                key={i}
                notification={notification}
                onDelete={handleDeleteNotification}
              />
            ))
          ) : (
            <li>No notifications yet</li>
          )}
        </ul>
      )}
    </div>
  );
}

export default Notifications;
