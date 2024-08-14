import React from "react";
import { useNavigate } from "react-router-dom";
import "../styles/CommentPage.css";
import ErrorPage from "./ErrorPage";
import { Link } from "react-router-dom";

const fetchSinglePost = async (postId) => {
  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    timeout: 5000, // Timeout value in milliseconds (adjust as needed)
  };
  try {
    const response = await fetch(
      `http://localhost:6969/api/posts?id=${postId}`,
      requestOptions
    );

    if (response.status !== 200) {
      return response.status;
    }

    const data = await response.json();
    if (response.status === 200) {
      // console.log("posts fetched");
      // console.log(data.posts);
      return data.posts[0];
    } else {
      console.log(response.status);
      return response.status;
    }
  } catch (error) {
    console.error("Error fetching posts:", error);
    return 500;
  }
};

const fetchComments = async (postId) => {
  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    timeout: 5000, // Timeout value in milliseconds (adjust as needed)
  };

  try {
    const response = await fetch(
      `http://localhost:6969/api/serve-comments?id=${postId}`,
      requestOptions
    );

    if (response.status !== 200) {
      return response.status;
    }

    const data = await response.json();
    console.log("comments fetched");
    return data.comments;
  } catch (error) {
    console.error("Error fetching comments:", error);
    return 500;
  }
};

const SinglePostView = () => {
  const navigate = useNavigate();
  const url = window.location.href;
  const pattern = "/(?:profile/[^/]+/)?post/([^/]+)";
  const match = url.match(pattern);
  const postId = match[1];
  const [post, setPost] = React.useState([]);
  const [comments, setComments] = React.useState([]);
  const [error, setError] = React.useState("");
  // console.log("postid", postId);

  const redirectToHome = () => {
    if (url.includes("/profile")) {
      navigate(`/profile/${post.user_id}`);
    } else {
      navigate("/");
    }
  };

  React.useEffect(() => {
    console.log("useEffect triggered");
    const getPost = async () => {
      const postFromServer = await fetchSinglePost(postId);
      console.log("postFromServer", postFromServer);
      switch (postFromServer) {
        case 400:
          setError("Bad request.");
          break;
        case 401:
          setError("Unauthorized.");
          break;
        case 403:
          setError("Forbidden.");
          break;
        case 404:
          setError("Post not found.");
          break;
        case 500:
          setError("Server error.");
          break;
        default:
          setPost(postFromServer);
      }
    };
    getPost();
  }, [postId]);

  React.useEffect(() => {
    const getComments = async () => {
      const commentsFromServer = await fetchComments(postId);
      switch (commentsFromServer) {
        case 400:
          setError("Bad request.");
          break;
        case 401:
          setError("Unauthorized.");
          break;
        case 403:
          setError("Forbidden.");
          break;
        case 404:
          setError("Post not found.");
          break;
        case 500:
          setError("Server error.");
          break;
        default:
          setComments(commentsFromServer);
      }
    };
    getComments();
  }, [postId]);

  const SubmitComment = async (e) => {
    e.preventDefault();
    // console.log("comment submitted");

    // send post to database
    const commentInput = document.getElementById("comment");
    const fileInput = document.getElementById("comment-file");

    if (!postId) {
      alert("Post ID is missing");
      return;
    }
    //check that content does not consist of only spaces
    if (commentInput.value.trim().length === 0) {
      alert("Comment must contain at least one non-space character.");
      return;
    }

    const formData = new FormData();
    formData.append("post_id", postId);
    formData.append("content", commentInput.value);

    if (fileInput && fileInput.files[0]) {
      formData.append("image", fileInput.files[0]);
    }

    const requestOptions = {
      method: "POST",
      body: formData,
      credentials: "include",
    };

    const response = await fetch(
      "http://localhost:6969/api/commenting",
      requestOptions
    );

    if (response.status === 200) {
      // clear form
      document.getElementById("comment").value = "";
      document.getElementById("comment-file").value = null;
      // console.log(response);
      // fetch posts again
      let updatedComments = await fetchComments(postId);
      setComments(updatedComments);
    } else {
      alert("Error posting.");
    }
  };

  // if (!post) {
  //   return <ErrorPage errorType="500" />;
  // }

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

  if (error) {
    return <ErrorPage errorType={error} />;
  }

  return (
    <div className="comment-container">
      <button className="back-button" onClick={redirectToHome}>
        X
      </button>
      <div className="singlepost">
        <div className="og-author">
          <Link to={`/profile/${post.user_id}`}>{post.full_name}</Link>
        </div>
        <div className="og-content">{post.content}</div>
        {post.picture ? (
          <div className="og-image">
            <img
              src={`data:image/jpeg;base64,${post.picture}`}
              className="pic"
            ></img>
          </div>
        ) : null}
        <div className="og-timecreated">{formatTimestamp(post.date)}</div>
      </div>

      <input type="hidden" id="post_id" value={post.id} />

      <textarea
        className="comment-box"
        type="text"
        rows="5"
        placeholder="Comment here..."
        id="comment"
      />
      <input type="file" id="comment-file" accept="image/*" />
      <button className="submit-comment" onClick={SubmitComment}>
        Comment
      </button>

      {comments ? (
        <div>
          <div className="comment-section">Comment Section</div>
          <div className="commentbox-container">
            {comments.map((comment) => (
              <div className="yourcomment" key={comment.id}>
                <div className="commentator">
                  <Link to={`/profile/${comment.user_id}`}>
                    {comment.full_name}
                  </Link>
                </div>

                <div className="new-comment">{comment.content}</div>
                {comment.image && (
                  <div className="comment-image">
                    <img
                      src={`data:image/jpeg;base64,${comment.image}`}
                      className="pic"
                    ></img>
                  </div>
                )}
                <div className="comment-time">
                  {formatTimestamp(comment.created_at)}
                </div>
              </div>
            ))}
          </div>
        </div>
      ) : (
        <div>no comments</div>
      )}
    </div>
  );
};

export default SinglePostView;
