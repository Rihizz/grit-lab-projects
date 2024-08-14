import React, { isValidElement, useEffect, useState } from "react";
import { useAuth } from "../AuthContext";
import "../styles/ProfileCard.css";
import EditProfileModal from "./EditProfileModal";
import ProfilePosts from "./ProfilePosts";
// import Follow from "./Follow";
// import Unfollow from "./Unfollow";

async function Unfollow(userId) {
  // console.log("inside Unfollow.js");
  // console.log("userId", userId);

  const requestOptions = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      id: userId,
    }),
    credentials: "include",
  };

  const response = await fetch(
    "http://localhost:6969/api/unfollow",
    requestOptions
  );

  if (response.ok) {
    // console.log("unfollowed");
    return <></>;
  } else {
    // console.log("unfollow failed");
    return <></>;
  }
}

async function Follow(userId, privacy, userID, ws, nickname) {
  // const { userID, ws, nickname } = useAuth();

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

  return <></>;
}

function ProfileCard(props) {
  const { user, ownProfile, setUser, userId } = props;
  const { userID, ws, nickname } = useAuth();
  const [shouldFollow, setShouldFollow] = useState(false);
  const [isFollowing, setIsFollowing] = useState(false);
  const [shouldReloadPosts, setShouldReloadPosts] = useState(true); // [1
  const [showFollowers, setShowFollowers] = useState(false);
  const [showFollowing, setShowFollowing] = useState(false);
  const [isRequestSent, setRequestSent] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  const [showEditModal, setShowEditModal] = useState(false);
  const [avatarSrc, setAvatarSrc] = useState(
    user.avatar
      ? `data:image/jpeg;base64,${user.avatar}`
      : "path/to/default/avatar.jpg"
  );

  // console.log("user", user);

  // console.log("userid", userId);

  const handleModalClose = () => {
    setShowEditModal(false);
  };

  async function fetchUserData(userId) {
    // console.log("fetching user data for user id: " + userId);
    const requestOption = {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include", // send the cookie along with the request
    };
    const response = await fetch(
      "http://localhost:6969/api/user/" + userId,
      requestOption
    );
    const data = await response.json();
    if (response.status !== 200) {
      throw Error(data.message);
    } else {
      // console.log(data);
      setUser(data.user);
      setAvatarSrc(
        data.user.avatar
          ? `data:image/jpeg;base64,${data.user.avatar}`
          : "path/to/default/avatar.jpg"
      );
    }
  }

  useEffect(() => {
    fetchUserData(userId);
  }, [userId]);

  const handleFollow = async () => {
    // console.log("handleFollow called");
    setShouldFollow(true);
    await Follow(userId, user.privacy, userID, ws, nickname);
    await fetchUserData(userId);
    setShouldReloadPosts(true);
  };

  async function handleUnfollow() {
    // console.log("handleUnfollow called");
    // console.log("userId in handleUnfollow", userId);
    setShouldFollow(false);
    await Unfollow(userId);
    await fetchUserData(userId);
    setShouldReloadPosts(true);
  }

  const handleModalSave = async (updatedData) => {
    const {
      userId,
      email,
      nickname,
      aboutMe,
      avatar,
      newPassword,
      confirmPassword,
      privacy,
    } = updatedData;

    const formData = new FormData();

    formData.append("userId", userId);
    if (email) formData.append("email", email);
    if (nickname) formData.append("nickname", nickname);
    if (aboutMe) formData.append("aboutMe", aboutMe);
    if (newPassword) formData.append("password", newPassword);
    if (confirmPassword) formData.append("confirmPassword", confirmPassword);
    if (privacy) formData.append("privacy", privacy);
    if (avatar) {
      formData.append("avatar", avatar);
    }

    const requestOptions = {
      method: "POST",

      headers: {
        // "Content-Type": "multipart/form-data" should NOT be set manually
        // The browser will automatically set the correct boundary
      },

      body: formData,

      credentials: "include",
    };

    const response = await fetch(
      "http://localhost:6969/api/user/update",
      requestOptions
    );
    if (response.ok) {
      const updatedUser = { ...user };

      if (email) updatedUser.email = email;
      if (nickname) updatedUser.nickname = nickname;
      if (aboutMe) updatedUser.aboutme = aboutMe;
      if (newPassword) updatedUser.newPassword = newPassword;
      if (confirmPassword) updatedUser.confirmPassword = confirmPassword;
      if (privacy) updatedUser.privacy = privacy;

      if (avatar) {
        const reader = new FileReader();
        reader.onloadend = function () {
          // Update the user's avatar data and the avatar source
          updatedUser.avatar = reader.result;
          setAvatarSrc(reader.result);
          setUser(updatedUser);
        };
        reader.readAsDataURL(avatar);
      } else {
        updatedUser.avatar = user.avatar;
        setUser(updatedUser);
      }
      setUser(updatedUser);

      setShowEditModal(false);
    } else {
      //set error message to be passed on to editprofilemodal in case of error
      const responseText = await response.text();
      // console.log("Raw Response Text:", responseText);
      setErrorMessage(responseText);
    }
  };

  //cut out the time from the birthday
  const birthday = user.birthday;
  const birthdayDate = birthday.split("T")[0];
  user.birthday = birthdayDate;

  const handleFollowersClick = () => {
    setShowFollowers(!showFollowers);
  };

  const handleFollowingClick = () => {
    setShowFollowing(!showFollowing);
  };

  return (
    <div className="card">
      <div className="card-body">
        {((!isFollowing && !ownProfile && isRequestSent) || !isRequestSent) && (
          <>
            <h5 className="card-title">
              {user.firstName} {user.lastName}
            </h5>

            <img src={avatarSrc} alt="Avatar" className="card-img" />
          </>
        )}

        {(user.ownProfile || user.isFollowing || user.privacy === "0") && (
          <>
            <p className="card-email">Email: {user.email}</p>

            <p className="card-nickname">Nickname: {user.nickname}</p>

            <p className="card-aboutme">About me: {user.aboutme}</p>

            <p className="card-birthday">Date of birth: {user.birthday}</p>

            <div className="card-followers">
              <span
                className="card-followers-header"
                onClick={handleFollowersClick}
              >
                Followers: {user.followers ? user.followers.length : 0}
              </span>
              {showFollowers && user.followers ? (
                user.followers.map((follower) => (
                  <div className="follower" key={follower.id}>
                    {follower.full_name}
                  </div>
                ))
              ) : (
                <div className="hidden">No followers</div>
              )}
            </div>

            <div className="card-following">
              <span
                className="card-following-header"
                onClick={handleFollowingClick}
              >
                Following: {user.following ? user.following.length : 0}
              </span>
              {showFollowing && user.following ? (
                user.following.map((person) => (
                  <div className="following" key={person.id}>
                    {person.full_name}
                  </div>
                ))
              ) : (
                <div className="hidden">Not following anyone</div>
              )}
            </div>
          </>
        )}

        {ownProfile ? (
          <p className="card-privacy">
            Privacy: {user.privacy == 0 ? "Public" : "Private"}
          </p>
        ) : (
          <></>
        )}
        {ownProfile ? (
          <button
            className="editprofile-button"
            onClick={() => setShowEditModal(true)}
          >
            Edit Profile
          </button>
        ) : (
          <>
            {/* if following make unfollowbutton */}
            <div>
              {user.isFollowing ? (
                <button className="unfollow-button" onClick={handleUnfollow}>
                  Unfollow
                </button>
              ) : (
                <button className="follow-button" onClick={handleFollow}>
                  {user.privacy === "1" ? "Request Follow" : "Follow"}
                </button>
              )}
            </div>
          </>
        )}
      </div>

      <EditProfileModal
        show={showEditModal}
        handleClose={handleModalClose}
        handleSave={handleModalSave}
        userId={userId}
        currentUserData={user}
        errorMessage={errorMessage}
      />

      <ProfilePosts
        userId={userId}
        shouldReloadPosts={shouldReloadPosts}
        setShouldReloadPosts={setShouldReloadPosts}
      />
    </div>
  );
}

export default ProfileCard;
