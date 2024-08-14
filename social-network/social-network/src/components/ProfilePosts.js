import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import "../styles/ProfileCard.css";
import { useAuth } from "../AuthContext";

async function fetchPosts(idProfile) {
  // console.log("fetching posts for user id: " + idProfile);
  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  };
  const response = await fetch(
    "http://localhost:6969/api/pposts?id=" + idProfile,
    requestOptions
  );
  const data = await response.json();
  if (response.status === 200) {
    // console.log("posts fetched");
    return data;
  } else {
    alert("Error fetching posts.");
  }
}

function ProfilePosts({ userId, shouldReloadPosts, setShouldReloadPosts }) {
  const [posts, setPosts] = useState([]);
  // console.log("userId", userId);
  // console.log("shouldReloadPosts", shouldReloadPosts);

  useEffect(() => {
    if (shouldReloadPosts) {
      // console.log("shouldReloadPosts is true");
      const getPosts = async () => {
        const postsFromServer = await fetchPosts(userId);
        setPosts(postsFromServer);
        // console.log("postfromserer", postsFromServer);
        if (postsFromServer === null) {
          setPosts([]);
        }
      };
      getPosts();
      setShouldReloadPosts(false);
    }
  }, [userId, shouldReloadPosts, setShouldReloadPosts]); // Notice the dependency array here

  // console.log("posts:", posts);

  // render something here

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

  // Handle loading state
  if (posts === null) {
    return <div>Loading...</div>;
  }

  // Handle case where user has no posts
  if (posts.length === 0) {
    return <div className="no-post-display">No posts to display.</div>;
  }

  return (
    <div className="allpostXD">
      <div className="post-containerXD">
        {posts.map((post) => {
          const postImageSrc = post.picture
            ? `data:image/jpeg;base64,${post.picture}`
            : null;
          return (
            <div key={post.id} className="ownpost">
              <div className="posterXD">
                <a to={`/profile/${post.user_id}`}>{post.full_name}</a>
              </div>

              {postImageSrc && (
                <img src={postImageSrc} alt="Post" className="post-imgXD" />
              )}
              <div className="post-contentXD">{post.content}</div>
              <div className="post-dateXD">{formatTimestamp(post.date)}</div>
              {/* <div className="likes">
              <span>{post.like_count} </span>
            </div> */}

              <div className="opencomments">
                <i className="fa-solid fa-comments fa-lg"></i>
                <Link to={`/profile/${post.user_id}/post/${post.id}`}>
                  Open Comments
                </Link>
              </div>

              {/* <span>{post.likes} </span> */}
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default ProfilePosts;
