import React from "react";
import "../styles/PostContainer.css";
import PostingForm from "./PostingForm";
import { Link } from "react-router-dom";

// Function to retrieve liked post IDs from localStorage
function getLikedPostsFromStorage() {
  const likedPosts = localStorage.getItem("likedPosts");
  return likedPosts ? JSON.parse(likedPosts) : [];
}

// Function to save liked post IDs to localStorage
function saveLikedPostsToStorage(likedPosts) {
  localStorage.setItem("likedPosts", JSON.stringify(likedPosts));
}

async function fetchPosts() {
  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  };
  try {
  const response = await fetch(
    "http://localhost:6969/api/posts",
    requestOptions
  );
  const data = await response.json();
  if (response.status === 200) {
    // console.log("posts fetched");
    // console.log(data.posts);
    return data.posts;
  } else {
    alert("Error fetching posts.");
  }
  } catch (error) {
    console.error("Error fetching posts:", error);
    return 500;
  }
}

function PostContainer() {
  const [posts, setPosts] = React.useState([]);
  const [likedPosts, setLikedPosts] = React.useState(getLikedPostsFromStorage);

  React.useEffect(() => {
    const getPosts = async () => {
      const posts = await fetchPosts();
      setPosts(posts);
    };
    getPosts();
  }, []);

  React.useEffect(() => {
    // Update localStorage when likedPosts state changes
    saveLikedPostsToStorage(likedPosts);
  }, [likedPosts]);

  // console.log("posts:", posts);

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
    <div className="allposts">
      <div className="post-container">
        <PostingForm fetchPosts={fetchPosts} setPosts={setPosts} />
        {posts.map((post) => {
          const postImageSrc = post.picture
            ? `data:image/jpeg;base64,${post.picture}`
            : null;
          return (
            <div key={post.id} className="post">
              <div className="poster">
                <Link to={`/profile/${post.user_id}`}>{post.full_name}</Link>
              </div>
              <div className="post-content">{post.content}</div>
              {postImageSrc && (
                <img src={postImageSrc} alt="Post" className="post-img" />
              )}

              <div className="post-date">{formatTimestamp(post.date)}</div>

              {/* <i
                onClick={() => {
                  toggleLike(post.id);
                  handleLikeClick(post.id);
                }}
                // className={`fa fa-thumbs-up ${temp ? "liked" : ""}`}
              ></i>
              <div className="likes">
                <span>{post.like_count} </span>
              </div> */}

              <div className="opencomments">
                <i className="fa-solid fa-comments fa-lg"></i>

                <Link to={`/post/${post.id}`}> Open Comments</Link>
              </div>

              {/* <span>{post.likes} </span> */}
            </div>
          );
        })}
      </div>
    </div>
  );
}
export default PostContainer;
