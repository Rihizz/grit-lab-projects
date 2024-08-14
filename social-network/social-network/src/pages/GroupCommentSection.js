import React from "react";
import { useNavigate } from "react-router-dom";
import "../styles/GroupCommentView.css";
import ErrorPage from "./ErrorPage";
import { Link } from "react-router-dom";

const fetchSinglePost = async (postId, groupId) => {
  if (!postId || !groupId) {
    return;
  }
  console.log("postId in fetchSinglePost", postId);
  console.log("groupId in fetchSinglePost", groupId);

  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  };

  const response = await fetch(
    `http://localhost:6969/api/serve-group-posts?group-postID=${postId}&id=${groupId}`,
    requestOptions
  );
  if (response.status !== 200) {
    return response.status;
  }
  const data = await response.json();
  if (response.status === 200) {
    // console.log("posts fetched");
    console.log(data);
    return data;
  } else if (response.status === 403) {
    // console.log("error", response.status);
    return response.status;
  }
};

const fetchComments = async (postId, groupId) => {
  if (!postId || !groupId) {
    return;
  }
  console.log("postId in fetchComments", postId);
  console.log("groupId in fetchComments", groupId);
  const requestOptions = {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  };
  // console.log("groupid in fetchComments", groupId);
  const response = await fetch(
    `http://localhost:6969/api/serve-group-comments?id=${postId}&&groupId=${groupId}`,
    requestOptions
  );
  if (response.status !== 200) {
    return response.status;
  }
  const data = await response.json();
  if (response.status === 200) {
    console.log("comments fetched");
    console.log(data.comments);
    return data.comments;
  } else {
    alert("Error fetching comments.");
  }
};

const GroupCommentView = () => {
  const [post, setPost] = React.useState(null);
  const [comments, setComments] = React.useState([]);
  const [error, setError] = React.useState("");
  const [postId, setPostId] = React.useState("");
  const [groupId, setGroupId] = React.useState("");
  const navigate = useNavigate();
  const url = window.location.href;
  const pattern = /(\d+)$/;
  const regex = /group\/(\d+)/;
  const groupIdMatch = url.match(regex);
  const match = url.match(pattern);

  console.log("post", postId);

  const redirectToHome = () => {
    navigate(`/groups/${groupId}`);
  };

  React.useEffect(() => {
    if (match) {
      console.log("match", match);
      setPostId(match[1]);
    } else {
      console.log("no match");
      setError("404");
    }
    if (groupIdMatch) {
      console.log("groupIdMatch", groupIdMatch);
      setGroupId(groupIdMatch[1]);
    } else {
      console.log("no groupIdMatch");
      setError("404");
    }
    const getPost = async () => {
      const postFromServer = await fetchSinglePost(postId, groupId);
      // console.log("postFromServer", postFromServer);
      if (postFromServer === 403) {
        setError("403");
        // console.log("error", error);
        return;
      }
      setPost(postFromServer);
    };
    getPost();
  }, [postId, groupId]);
  console.log("group id", groupId);
  console.log("post id", postId);

  React.useEffect(() => {
    const getComments = async () => {
      const commentsFromServer = await fetchComments(postId, groupId);
      switch (commentsFromServer) {
        case 400:
          setError("400 Bad Request");
          return;
        case 401:
          setError("401 Unauthorized");
          return;
        case 403:
          setError("403 Forbidden");
          return;
        case 404:
          setError("404 Not Found");
          return;
        case 500:
          setError("500 Internal Server Error");
          return;
        default:
          break;
      }
      setComments(commentsFromServer);
    };
    getComments();
  }, [postId, groupId]);

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
    //check that comment does not consist of only spaces
    if (commentInput.value.trim().length === 0) {
      alert("Comment must contain at least one non-space character.");
      return;
    }

    const formData = new FormData();
    formData.append("post_id", postId);
    formData.append("content", commentInput.value);

    if (fileInput && fileInput.files[0]) {
      console.log("got here");
      formData.append("picture", fileInput.files[0]);
    }

    const requestOptions = {
      method: "POST",
      body: formData,
      credentials: "include",
    };

    const response = await fetch(
      "http://localhost:6969/api/group-commenting",
      requestOptions
    );
    if (response.status === 200) {
      // console.log("im in here");
      // clear form
      document.getElementById("comment").value = "";
      // console.log(response);
      // fetch posts again
      let updatedComments = await fetchComments(postId, groupId);
      setComments(updatedComments);
    } else {
      alert("Error posting.");
    }
  };

  if (error !== "") {
    return <ErrorPage errorType={error} />;
  }

  if (!post) {
    console.log("post is null");
    return <div>Loading...</div>;
  }

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

  // console.log("comments", comments);

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

      <input type="hidden" id="post_id" value={post ? post.id : ""} />
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
              <div className="yourcomment" key={comment.comment_id}>
                <div className="commentator">{comment.full_name}</div>
                <div className="new-comment">{comment.content}</div>
                {comment.image ? (
                  <div className="comment-image">
                    <img
                      className="group-comment-image"
                      src={`data:image/jpeg;base64,${comment.image}`}
                      alt="Comment"
                    ></img>
                  </div>
                ) : null}
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

export default GroupCommentView;
