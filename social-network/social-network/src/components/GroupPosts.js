import React from "react";
import { Link } from "react-router-dom";

import "../styles/Groups.css";

const fetchGroupPosts = async (groupId) => {
  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  };

  const response = await fetch(
    `http://localhost:6969/api/serve-group-posts?id=${groupId}`,
    requestOptions
  );
  const data = await response.json();
  if (response.status === 200) {
    // console.log("group posts fetched");
    // console.log(data);
    return data;
  } else {
    return response.status;
  }
};

const GroupPosts = () => {
  const [notAuthorizedMessage, setNotAuthorizedMessage] = React.useState("");
  const url = window.location.href;
  const pattern = /groups\/(\d+)/;
  const match = url.match(pattern);
  // console.log(match);
  const groupId = match[1];
  const [groupPosts, setGroupPosts] = React.useState([]);

  React.useEffect(() => {
    const getGroupPosts = async () => {
      const groupPostsFromServer = await fetchGroupPosts(groupId);
      // console.log("group posts from server", groupPostsFromServer);
      if (groupPostsFromServer == 403) {
        // console.log("not authorized");
        setNotAuthorizedMessage(
          "To view group posts, please request group access."
        );
      } else {
        setGroupPosts(groupPostsFromServer);
      }
    };
    getGroupPosts();
    // console.log("group posts");
  }, [groupId]);


  const handleSubmit = async (e) => {
    //check that post does not consist of only spaces
    const postInput = document.getElementById("post-textarea");
    if (postInput.value.trim().length === 0) {
      alert("Post must contain at least one non-space character.");
      return;
    }

    e.preventDefault();

    const post = postInput.value;
    if (post.length < 10 || post.length > 500) {
      alert("Post must be 10-500 characters long.");
      return;
    }

    const picture = document.getElementById("picture");

    const formData = new FormData();
    formData.append("group_id", groupId);
    formData.append("content", post);

    if (picture.files[0]) {
      formData.append("picture", picture.files[0]);
    }

    const requestOptions = {
      method: "POST",
      body: formData,
      credentials: "include",
    };
    const response = await fetch(
      `http://localhost:6969/api/group-posting`,
      requestOptions
    );

    if (response.status === 200) {
      console.log("post created should load new posts");
      postInput.value = "";
      picture.value = "";
      const updatedPosts = await fetchGroupPosts(groupId);
      setGroupPosts(updatedPosts);
    } else {
      alert("Error posting to group.");
    }
    
  };

  //gigachad way of just deducting 3 hours from the time and displaying the shit out of it
  const formatTimestamp = (dateTime) => {
    // Parse the date-time string
    const date = new Date(dateTime);

    // Subtract 3 hours in milliseconds
    date.setTime(date.getTime() - 3 * 60 * 60 * 1000);

    // Format the adjusted time
    const options = {
      weekday: "long",
      day: "numeric",
      month: "long",
      hour: "numeric",
      minute: "numeric",
    };
    return date.toLocaleString(undefined, options);
  };

  return (
    <div>
      {notAuthorizedMessage ? (
        <h3>{notAuthorizedMessage}</h3>
      ) : (
        <div className="group-post-input">
          <h1>Create a Post</h1>
          {/* Rest of the form */}
          <div className="group-post-container">
            <textarea
              className="post-textarea"
              placeholder="What's on your mind?"
              id="post-textarea"
              required
              maxLength="500"
              minLength="10"
              title="Post should be 10-500 characters."
            />
            <input type="file" id="picture" accept="image/*" />
            <button
              type="submit"
              className="group-button-post"
              onClick={handleSubmit}
            >
              Post
            </button>
          </div>
          <h1>Group posts:</h1>
          {groupPosts ? (
            groupPosts.map((groupPost) => {
              const postImageSrc = groupPost.Image
                ? `data:image/jpeg;base64,${groupPost.Image}`
                : null;

              // Return the component from the map function
              return (
                <div className="group-post" key={groupPost.Id}>
                  <h4>{groupPost.FullName}</h4>
                  <div className="post-content">{groupPost.Post}</div>
                  {postImageSrc && (
                    <img src={postImageSrc} alt="Post" className="post-img" />
                  )}
                  <h5>{formatTimestamp(groupPost.CreatedAt)}</h5>
                  <div className="group-post-comment-section">
                    <i className="fa-solid fa-comments fa-lg"></i>
                    <Link to={`/group/${groupId}/group-post/${groupPost.Id}`}>
                      Open Comments
                    </Link>
                  </div>
                </div>
              );
            })
          ) : (
            <h3>No posts yet.</h3>
          )}
        </div>
      )}
    </div>
  );
};

export default GroupPosts;
